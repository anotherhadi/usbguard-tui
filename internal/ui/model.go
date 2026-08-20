package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/helpbar"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/notification"
	"github.com/anotherhadi/ilovetui/style"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type (
	tickMsg          time.Time
	devicesMsg       []guard.Device
	daemonStatusMsg  string
	defaultPolicyMsg guard.Status
	actionMsg        struct{ err error }
	nixRuleMsg       struct{ rule string }
)

type deviceSummary struct {
	total, allowed, blocked, rejected int
}

type Model struct {
	list          list.Model
	help          helpbar.Model
	modals        modal.Model
	notif         notification.Model
	daemonStatus  string
	defaultPolicy guard.Status
	deviceCounts  deviceSummary
	width         int
	height        int
	rulesManaged  bool
	rulesWritable bool
	pendingRules  []string
	fatalShown    bool
}

func (m Model) PendingRules() []string { return m.pendingRules }

func New() Model {
	l := bubbles.NewList(nil, 0, 0)
	l.SetDelegate(deviceDelegate{})
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.DisableQuitKeybindings()
	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up"))
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down"))

	h := helpbar.New(
		helpbar.WithToggle(listKeys.Help),
		helpbar.WithGlobal(listKeys.globalBindings()...),
	)

	rulesManaged := guard.IsRulesManaged()
	rulesWritable, _ := guard.RulesWritable()

	return Model{
		list:          l,
		help:          h,
		modals:        modal.New(),
		notif:         notification.New(),
		rulesManaged:  rulesManaged,
		rulesWritable: rulesWritable,
	}
}

func makeActionList(rulesManaged bool) list.Model {
	var items []list.Item
	if rulesManaged {
		items = []list.Item{
			actionItem{"allow", guard.AllowDevice, false, guard.Allowed, false},
			actionItem{"allow (perm)", nil, true, guard.Allowed, true},
			actionItem{"block", guard.BlockDevice, false, guard.Blocked, false},
			actionItem{"block (perm)", nil, true, guard.Blocked, true},
			actionItem{"reject", guard.RejectDevice, false, guard.Rejected, false},
			actionItem{"reject (perm)", nil, true, guard.Rejected, true},
		}
	} else {
		items = []list.Item{
			actionItem{"allow", guard.AllowDevice, false, guard.Allowed, false},
			actionItem{"allow (permanent)", guard.AllowDevice, true, guard.Allowed, false},
			actionItem{"block", guard.BlockDevice, false, guard.Blocked, false},
			actionItem{"block (permanent)", guard.BlockDevice, true, guard.Blocked, false},
			actionItem{"reject", guard.RejectDevice, false, guard.Rejected, false},
			actionItem{"reject (permanent)", guard.RejectDevice, true, guard.Rejected, false},
		}
	}
	l := bubbles.NewList(items, 24, len(items))
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()
	l.SetFilteringEnabled(false)

	l.SetDelegate(actionDelegate{})
	return l
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy, tickCmd(), m.modals.Init(), m.notif.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	next, cmd := m.updateMain(msg)
	nm := next.(Model)

	var modalCmd, notifCmd tea.Cmd
	nm.modals, modalCmd = nm.modals.Update(msg)
	nm.notif, notifCmd = nm.notif.Update(msg)

	return nm, tea.Batch(cmd, modalCmd, notifCmd)
}

func (m Model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		m.list.SetSize(msg.Width, m.listHeight())
		return m, nil

	case tickMsg:
		return m, tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy, tickCmd())

	case devicesMsg:
		items := make([]list.Item, len(msg))
		summary := deviceSummary{total: len(msg)}
		for i, d := range msg {
			items[i] = d
			switch d.Status {
			case guard.Allowed:
				summary.allowed++
			case guard.Blocked:
				summary.blocked++
			case guard.Rejected:
				summary.rejected++
			}
		}
		m.deviceCounts = summary
		cmd := m.list.SetItems(items)
		return m, cmd

	case daemonStatusMsg:
		m.daemonStatus = string(msg)
		return m, nil

	case defaultPolicyMsg:
		m.defaultPolicy = guard.Status(msg)
		return m, nil

	case nixRuleMsg:
		m.pendingRules = append(m.pendingRules, msg.rule)
		return m, nil

	case actionMsg:
		if msg.err != nil {

			if msg.err == guard.ErrPermission {

				if m.fatalShown {
					return m, nil
				}
				m.fatalShown = true
				return m, modal.Show("Permission Error", newPermissionModal(),
					modal.WithModalStyle(permissionModalStyle()))
			}
			return m, errorToast(msg.err)
		}
		return m, fetchDevices

	case tea.KeyPressMsg:
		if m.modals.Open() {
			return m, nil
		}
		return m.updateList(msg)

	case tea.MouseWheelMsg:
		if m.modals.Open() {
			return m, nil
		}
		switch msg.Button {
		case tea.MouseWheelUp:
			m.list.CursorUp()
		case tea.MouseWheelDown:
			m.list.CursorDown()
		}
		return m, nil
	}

	if m.modals.Open() {
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateList(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if !m.list.SettingFilter() {
		dev, hasSelection := m.selectedDevice()
		switch {
		case key.Matches(msg, listKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, listKeys.Refresh):
			return m, tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy)
		case key.Matches(msg, listKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.list.SetSize(m.width, m.listHeight())
			return m, nil
		case key.Matches(msg, listKeys.Open):
			if hasSelection {
				return m, modal.Show(dev.Name, newActionModal(dev, m.rulesManaged),
					modal.WithModalStyle(actionModalStyle(dev.Status)))
			}
		case key.Matches(msg, listKeys.AllowAll):
			return m, doBulkAction(m.visibleDeviceIDs(), guard.AllowDevice, false)
		case key.Matches(msg, listKeys.AllowAllPerm):
			if m.rulesManaged {
				return m, queueNixOSRules(m.visibleDevices(), guard.Allowed)
			}
			return m, doBulkAction(m.visibleDeviceIDs(), guard.AllowDevice, true)
		case key.Matches(msg, listKeys.PrintAll):
			return m, queueCurrentStateRules(m.visibleDevices())
		}
		if hasSelection {
			if cmd := m.deviceActionCmd(msg, dev); cmd != nil {
				return m, cmd
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	return tea.View{
		Content:     m.renderContent(),
		AltScreen:   true,
		WindowTitle: "USBGuard TUI",
		MouseMode:   tea.MouseModeCellMotion,
	}
}

func (m Model) renderContent() string {
	header := m.renderHeader()
	listView := strings.TrimRight(m.list.View(), "\n")
	helpView := strings.TrimRight(m.help.View(), "\n")
	bg := strings.Join([]string{header, listView, helpView}, "\n")

	bg = m.modals.Render(bg)
	bg = m.notif.Render(bg)
	return bg
}

func (m Model) renderHeader() string {
	title := headerStyle.Render("USBGuard-tui")
	lines := []string{
		title,
		m.renderServiceLine(),
		m.renderPolicyLine(),
		m.renderDevicesLine(),
		m.renderRulesLine(),
	}
	if pending := m.renderPendingRulesLine(); pending != "" {
		lines = append(lines, pending)
	}

	for i, l := range lines {
		lines[i] = clampToWidth(l, m.width)
	}
	return strings.Join(lines, "\n")
}

func clampToWidth(s string, width int) string {
	if width <= 0 {
		return s
	}
	return ansi.Truncate(s, width, "…")
}

func (m Model) renderServiceLine() string {
	label := infoLabelStyle.Render("Service")
	switch m.daemonStatus {
	case "active":
		return label + daemonActiveStyle.Render("active")
	case "":
		return label + mutedStyle.Render("checking...")
	default:
		return label + daemonOtherStyle.Render(m.daemonStatus)
	}
}

func (m Model) renderPolicyLine() string {
	label := infoLabelStyle.Render("Default policy")
	if m.defaultPolicy == "" {
		return label + mutedStyle.Render("unknown")
	}
	clr, ok := statusColors[m.defaultPolicy]
	if !ok {
		clr = style.S.Muted
	}
	return label + lipgloss.NewStyle().Foreground(clr).Render(string(m.defaultPolicy))
}

func (m Model) renderDevicesLine() string {
	label := infoLabelStyle.Render("Devices")
	if m.deviceCounts.total == 0 {
		return label + mutedStyle.Render("0")
	}
	parts := []string{mutedStyle.Render(fmt.Sprintf("%d total", m.deviceCounts.total))}
	if m.deviceCounts.allowed > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(statusColors[guard.Allowed]).
			Render(fmt.Sprintf("%d allow", m.deviceCounts.allowed)))
	}
	if m.deviceCounts.blocked > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(statusColors[guard.Blocked]).
			Render(fmt.Sprintf("%d block", m.deviceCounts.blocked)))
	}
	if m.deviceCounts.rejected > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(statusColors[guard.Rejected]).
			Render(fmt.Sprintf("%d reject", m.deviceCounts.rejected)))
	}
	return label + strings.Join(parts, mutedStyle.Render("  ·  "))
}

func (m Model) renderRulesLine() string {
	label := infoLabelStyle.Render("Rules")
	switch {
	case m.rulesManaged:
		return label + warnStyle.Render("read-only (NixOS managed)")
	case !m.rulesWritable:
		return label + warnStyle.Render("read-only")
	default:
		return label + daemonActiveStyle.Render("writable")
	}
}

func (m Model) renderPendingRulesLine() string {
	count := len(m.pendingRules)
	if count == 0 {
		return ""
	}
	noun := "rule"
	if count > 1 {
		noun = "rules"
	}
	label := infoLabelStyle.Render("Pending rules")
	return label + warnStyle.Render(fmt.Sprintf("%d %s queued (printed on exit)", count, noun))
}

func (m Model) listHeight() int {
	headerH := lipgloss.Height(m.renderHeader())
	helpH := m.help.Height()
	return m.height - headerH - helpH
}

func (m Model) selectedDevice() (guard.Device, bool) {
	if item := m.list.SelectedItem(); item != nil {
		return item.(guard.Device), true
	}
	return guard.Device{}, false
}

func (m Model) visibleDevices() []guard.Device {
	items := m.list.VisibleItems()
	devices := make([]guard.Device, len(items))
	for i, item := range items {
		devices[i] = item.(guard.Device)
	}
	return devices
}

func (m Model) visibleDeviceIDs() []int {
	items := m.list.VisibleItems()
	ids := make([]int, len(items))
	for i, item := range items {
		ids[i] = item.(guard.Device).ID
	}
	return ids
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchDevices() tea.Msg {
	devices, err := guard.ListDevices()
	if err != nil {
		return actionMsg{err: err}
	}
	return devicesMsg(devices)
}

func fetchDaemonStatus() tea.Msg {
	return daemonStatusMsg(guard.DaemonStatus())
}

func fetchDefaultPolicy() tea.Msg {
	return defaultPolicyMsg(guard.DefaultPolicy())
}

func errorToast(err error) tea.Cmd {
	msg := err.Error()
	if err == guard.ErrReadOnly {
		msg = "Rules file is not writable: permanent changes are not supported."
	}
	return notification.Show("Error", msg, notification.Error, notification.WithID("action-error"))
}

func doAction(id int, fn func(int, bool) error, permanent bool) tea.Cmd {
	return func() tea.Msg {
		return actionMsg{err: fn(id, permanent)}
	}
}

func doBulkAction(ids []int, fn func(int, bool) error, permanent bool) tea.Cmd {
	return func() tea.Msg {
		for _, id := range ids {
			if err := fn(id, permanent); err != nil {
				return actionMsg{err: err}
			}
		}
		return actionMsg{}
	}
}

func queueNixOSRules(devices []guard.Device, status guard.Status) tea.Cmd {
	cmds := make([]tea.Cmd, len(devices))
	for i, d := range devices {
		rule := guard.NixOSRule(d, status)
		cmds[i] = func() tea.Msg { return nixRuleMsg{rule: rule} }
	}
	return tea.Batch(cmds...)
}

func queueCurrentStateRules(devices []guard.Device) tea.Cmd {
	cmds := make([]tea.Cmd, len(devices))
	for i, d := range devices {
		rule := guard.NixOSRule(d, d.Status)
		cmds[i] = func() tea.Msg { return nixRuleMsg{rule: rule} }
	}
	return tea.Batch(cmds...)
}

type actionBinding struct {
	binding key.Binding
	fn      func(int, bool) error
	perm    bool
	status  guard.Status
}

var deviceActionBindings = []actionBinding{
	{listKeys.Allow, guard.AllowDevice, false, guard.Allowed},
	{listKeys.AllowPerm, guard.AllowDevice, true, guard.Allowed},
	{listKeys.Block, guard.BlockDevice, false, guard.Blocked},
	{listKeys.BlockPerm, guard.BlockDevice, true, guard.Blocked},
	{listKeys.Reject, guard.RejectDevice, false, guard.Rejected},
	{listKeys.RejectPerm, guard.RejectDevice, true, guard.Rejected},
}

func (m Model) deviceActionCmd(msg tea.KeyPressMsg, dev guard.Device) tea.Cmd {
	for _, b := range deviceActionBindings {
		if !key.Matches(msg, b.binding) {
			continue
		}
		if b.perm && m.rulesManaged {
			return queueNixOSRules([]guard.Device{dev}, b.status)
		}
		return doAction(dev.ID, b.fn, b.perm)
	}
	return nil
}
