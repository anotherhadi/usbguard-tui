package ui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type state int

const (
	stateList state = iota
	statePopup
)

type (
	tickMsg         time.Time
	devicesMsg      []guard.Device
	daemonStatusMsg string
	actionMsg       struct{ err error }
)

type Model struct {
	state        state
	list         list.Model
	actionList   list.Model
	help         help.Model
	daemonStatus string
	width        int
	height       int
	notice       string
	selectedDev  *guard.Device
	rulesManaged bool
}

func New() Model {
	l := list.New(nil, deviceDelegate{}, 0, 0)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)
	l.SetShowTitle(false)
	l.DisableQuitKeybindings()
	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up"))
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down"))

	l.Styles = list.DefaultStyles(true)
	filterStyles := textinput.DefaultStyles(true)
	filterStyles.Focused.Prompt = filterStyles.Focused.Prompt.Foreground(colorAccent)
	filterStyles.Blurred.Prompt = filterStyles.Blurred.Prompt.Foreground(colorAccent)
	l.Styles.Filter = filterStyles

	h := help.New()
	h.Styles = help.DefaultStyles(true)

	rulesManaged := guard.IsRulesManaged()
	notice := ""
	if rulesManaged {
		notice = "Rules managed by NixOS config: permanent actions not available."
		listKeys.AllowPerm.SetEnabled(false)
		listKeys.BlockPerm.SetEnabled(false)
		listKeys.RejectPerm.SetEnabled(false)
	}

	return Model{
		state:        stateList,
		list:         l,
		actionList:   makeActionList(rulesManaged),
		help:         h,
		rulesManaged: rulesManaged,
		notice:       notice,
	}
}

func makeActionList(rulesManaged bool) list.Model {
	items := []list.Item{
		actionItem{"allow", guard.AllowDevice, false, guard.Allowed},
		actionItem{"block", guard.BlockDevice, false, guard.Blocked},
		actionItem{"reject", guard.RejectDevice, false, guard.Rejected},
	}
	if !rulesManaged {
		items = []list.Item{
			actionItem{"allow", guard.AllowDevice, false, guard.Allowed},
			actionItem{"allow (permanent)", guard.AllowDevice, true, guard.Allowed},
			actionItem{"block", guard.BlockDevice, false, guard.Blocked},
			actionItem{"block (permanent)", guard.BlockDevice, true, guard.Blocked},
			actionItem{"reject", guard.RejectDevice, false, guard.Rejected},
			actionItem{"reject (permanent)", guard.RejectDevice, true, guard.Rejected},
		}
	}
	l := list.New(items, actionDelegate{}, 24, len(items))
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()
	l.SetFilteringEnabled(false)
	return l
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchDevices, fetchDaemonStatus, tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		m.list.SetSize(msg.Width, m.listHeight())
		m.updateActionListSize()
		return m, nil

	case tickMsg:
		return m, tea.Batch(fetchDevices, fetchDaemonStatus, tickCmd())

	case devicesMsg:
		items := make([]list.Item, len(msg))
		for i, d := range msg {
			items[i] = d
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case daemonStatusMsg:
		m.daemonStatus = string(msg)
		return m, nil

	case actionMsg:
		m.state = stateList
		m.selectedDev = nil
		if msg.err != nil {
			switch msg.err {
			case guard.ErrReadOnly:
				m.notice = "Rules file is not writable: permanent changes are not supported."
			case guard.ErrPermission:
				m.notice = "Permission denied. Run with appropriate privileges."
			default:
				m.notice = msg.err.Error()
			}
		} else {
			m.notice = m.defaultNotice()
		}
		return m, fetchDevices

	case tea.KeyPressMsg:
		if m.state == statePopup {
			return m.updatePopup(msg)
		}
		return m.updateList(msg)

	case tea.MouseClickMsg:
		if m.state == statePopup {
			var cmd tea.Cmd
			m.actionList, cmd = m.actionList.Update(msg)
			return m, cmd
		}

	case tea.MouseWheelMsg:
		if m.state == statePopup {
			switch msg.Button {
			case tea.MouseWheelUp:
				m.actionList.CursorUp()
			case tea.MouseWheelDown:
				m.actionList.CursorDown()
			}
		} else {
			switch msg.Button {
			case tea.MouseWheelUp:
				m.list.CursorUp()
			case tea.MouseWheelDown:
				m.list.CursorDown()
			}
		}
		return m, nil
	}

	if m.state == stateList {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if !m.list.SettingFilter() {
		id := m.selectedDevID()
		switch {
		case key.Matches(msg, listKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, listKeys.Refresh):
			m.notice = m.defaultNotice()
			return m, tea.Batch(fetchDevices, fetchDaemonStatus)
		case key.Matches(msg, listKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.list.SetSize(m.width, m.listHeight())
			return m, nil
		case key.Matches(msg, listKeys.Open):
			if item := m.list.SelectedItem(); item != nil {
				d := item.(guard.Device)
				m.selectedDev = &d
				m.updateActionListSize()
				m.actionList.Select(0)
				m.state = statePopup
				return m, nil
			}
		case id >= 0 && key.Matches(msg, listKeys.Allow):
			return m, doAction(id, guard.AllowDevice, false)
		case id >= 0 && !m.rulesManaged && key.Matches(msg, listKeys.AllowPerm):
			return m, doAction(id, guard.AllowDevice, true)
		case id >= 0 && key.Matches(msg, listKeys.Block):
			return m, doAction(id, guard.BlockDevice, false)
		case id >= 0 && !m.rulesManaged && key.Matches(msg, listKeys.BlockPerm):
			return m, doAction(id, guard.BlockDevice, true)
		case id >= 0 && key.Matches(msg, listKeys.Reject):
			return m, doAction(id, guard.RejectDevice, false)
		case id >= 0 && !m.rulesManaged && key.Matches(msg, listKeys.RejectPerm):
			return m, doAction(id, guard.RejectDevice, true)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updatePopup(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, cancelKey):
		m.state = stateList
		m.selectedDev = nil
		return m, nil
	case key.Matches(msg, listKeys.Open):
		if item := m.actionList.SelectedItem(); item != nil {
			a := item.(actionItem)
			return m, doAction(m.selectedDev.ID, a.fn, a.permanent)
		}
	}
	var cmd tea.Cmd
	m.actionList, cmd = m.actionList.Update(msg)
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
	notice := m.renderNotice()
	listView := strings.TrimRight(m.list.View(), "\n")
	helpView := strings.TrimRight(m.help.View(listKeys), "\n")
	bg := strings.Join([]string{header, listView, notice, helpView}, "\n")

	if m.state == statePopup && m.selectedDev != nil {
		return placeOverlay(bg, m.renderActionSelect(), m.width, m.height)
	}
	return bg
}

func (m Model) renderHeader() string {
	title := headerStyle.Render("USBGuard-tui")
	switch m.daemonStatus {
	case "active":
		return title + mutedStyle.Render(" - ") + daemonActiveStyle.Render("active")
	case "":
		return title
	default:
		return title + mutedStyle.Render(" - ") + daemonOtherStyle.Render(m.daemonStatus)
	}
}

func (m Model) renderNotice() string {
	if m.notice == "" {
		return ""
	}
	return warnStyle.Render(m.notice)
}

func (m Model) renderActionSelect() string {
	dev := m.selectedDev
	color := statusColors[dev.Status]
	innerW := m.actionListInnerWidth()

	title := popupTitleStyle.Foreground(color).Width(innerW).Render(dev.Name)
	hint := lipgloss.NewStyle().Foreground(colorMuted).Width(innerW).Render("↑↓ navigate  enter confirm  esc cancel")

	content := strings.Join([]string{title, m.actionList.View(), "", hint}, "\n")
	return popupStyle.Width(innerW).Render(content)
}

func (m Model) popupOuterWidth() int {
	w := m.width - 6
	if w > 60 {
		w = 60
	}
	if w < 32 {
		w = 32
	}
	return w
}

func (m Model) actionListInnerWidth() int {
	return m.popupOuterWidth() - 8 // border(2) + padding_h(6)
}

func (m Model) defaultNotice() string {
	if m.rulesManaged {
		return "Rules managed by NixOS config: permanent actions not available."
	}
	return ""
}

func (m Model) actionItemCount() int {
	if m.rulesManaged {
		return 3
	}
	return 6
}

// updateActionListSize sizes the action list and toggles pagination based on available space.
// When there is enough room for all items: pagination is hidden and height is set exactly,
// avoiding the phantom line that bubbles/list reserves when showPagination=true.
// When space is limited: pagination is shown naturally by bubbles/list.
func (m *Model) updateActionListSize() {
	items := m.actionItemCount()
	innerW := m.actionListInnerWidth()
	// popup overhead: border(2) + padding_v(2) + title(1) + blank(1) + hint(1) = 7
	available := m.height - 7 - 2 // 2 lines margin
	if available >= items {
		m.actionList.SetShowPagination(false)
		m.actionList.SetSize(innerW, items)
	} else {
		m.actionList.SetShowPagination(true)
		h := available
		if h < 2 {
			h = 2
		}
		m.actionList.SetSize(innerW, h)
	}
}

func (m Model) listHeight() int {
	helpH := lipgloss.Height(strings.TrimRight(m.help.View(listKeys), "\n"))
	return m.height - 1 - helpH - 1 // header - help - notice
}

func (m Model) selectedDevID() int {
	if item := m.list.SelectedItem(); item != nil {
		return item.(guard.Device).ID
	}
	return -1
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

func doAction(id int, fn func(int, bool) error, permanent bool) tea.Cmd {
	return func() tea.Msg {
		return actionMsg{err: fn(id, permanent)}
	}
}
