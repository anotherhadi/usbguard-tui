// Package ui is usbguard-tui's root model: the device list, the header
// stats, the action overlay and the key/mouse routing.
package ui

import (
	"errors"
	"fmt"
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/foam/commands"
	"github.com/anotherhadi/foam/notify"
	"github.com/anotherhadi/foam/overlay"

	"github.com/anotherhadi/usbguard-tui/internal/config"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type (
	tickMsg          time.Time
	devicesMsg       []guard.Device
	daemonStatusMsg  string
	defaultPolicyMsg guard.Status
	actionMsg        struct{ err error }
	nixRuleMsg       struct{ key, rule string }
)

type pendingRule struct {
	key  string
	rule string
}

type deviceSummary struct {
	total, allowed, blocked, rejected int
}

type Model struct {
	cfg *config.Config
	reg *commands.Registry

	overlays overlay.Stack
	notes    notify.Center

	devices *deviceList

	daemonStatus  string
	defaultPolicy guard.Status
	counts        deviceSummary
	rulesManaged  bool
	rulesWritable bool
	pendingRules  []pendingRule

	w, h       int
	toastRects []toastRect
}

type toastRect struct {
	X, Y, W, H int
	id         int
}

func (r toastRect) contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

func New(cfg *config.Config) (*Model, error) {
	reg, err := commands.New(cfg.Keys)
	if err != nil {
		return nil, fmt.Errorf("invalid config:\n%w", err)
	}
	rulesWritable, _ := guard.RulesWritable()
	m := &Model{
		cfg:           cfg,
		reg:           reg,
		devices:       newDeviceList(),
		rulesManaged:  guard.IsRulesManaged(),
		rulesWritable: rulesWritable,
	}
	m.notes.Max = 50
	if err := reg.Check(m.allScopes()); err != nil {
		return nil, fmt.Errorf("invalid config:\n%w", err)
	}
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy, tick(m.cfg.Refresh))
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, m.update(msg)
}

func (m *Model) update(msg tea.Msg) tea.Cmd {
	if m.overlays.Handle(msg) {
		return nil
	}
	if cmd, ok := m.notes.Update(msg); ok {
		return cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		return m.key(msg)
	case tea.PasteMsg:
		if m.overlays.Top() == nil && m.devices.Capturing() {
			return m.devices.Paste(msg.Content)
		}
	case tea.MouseMsg:
		return m.mouse(msg)
	case tickMsg:
		return tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy, tick(m.cfg.Refresh))
	case devicesMsg:
		m.setDevices(msg)
	case daemonStatusMsg:
		m.daemonStatus = string(msg)
	case defaultPolicyMsg:
		m.defaultPolicy = guard.Status(msg)
	case nixRuleMsg:
		m.queueRule(msg.key, msg.rule)
	case actionMsg:
		if msg.err != nil {
			return actionErr(msg.err)
		}
		return fetchDevices
	}
	return nil
}

func (m *Model) setDevices(devs []guard.Device) {
	var s deviceSummary
	s.total = len(devs)
	for _, d := range devs {
		switch d.Status {
		case guard.Allowed:
			s.allowed++
		case guard.Blocked:
			s.blocked++
		case guard.Rejected:
			s.rejected++
		}
	}
	m.counts = s
	m.devices.SetDevices(devs)
}

func (m *Model) key(msg tea.KeyPressMsg) tea.Cmd {
	if !minSize.Fits(m.w, m.h) {
		if b, ok := m.reg.Match(m.path()[:1], msg); ok && b.ID == "app.quit" {
			return b.Run()
		}
		return nil
	}
	if m.overlays.Top() != nil {
		return m.overlays.Key(m.reg, msg)
	}
	if m.devices.Capturing() {
		if ok, cmd := m.devices.Input(msg); ok {
			return cmd
		}
	}
	if b, ok := m.reg.Match(m.path(), msg); ok {
		return b.Run()
	}
	return nil
}

func (m *Model) mouse(msg tea.MouseMsg) tea.Cmd {
	if !minSize.Fits(m.w, m.h) {
		return nil
	}
	if m.overlays.Top() != nil {
		return m.overlays.Mouse(msg)
	}
	mm := msg.Mouse()
	if _, click := msg.(tea.MouseClickMsg); click {
		for _, t := range m.toastRects {
			if t.contains(mm.X, mm.Y) {
				m.notes.DismissID(t.id)
				return nil
			}
		}
	}
	mm.Y -= m.bodyTop()
	var local tea.MouseMsg
	switch msg.(type) {
	case tea.MouseClickMsg:
		local = tea.MouseClickMsg(mm)
	case tea.MouseWheelMsg:
		local = tea.MouseWheelMsg(mm)
	default:
		return nil
	}
	if dev, open := m.devices.Mouse(local, true); open {
		return m.openActions(dev)
	}
	return nil
}

func (m *Model) keyHint(id string) string {
	if ks := m.reg.KeysOf(commands.Command{ID: id}); len(ks) > 0 {
		return ks[0].Canonical()
	}
	return "the palette"
}

func (m *Model) openActions(dev guard.Device) tea.Cmd {
	hint := ""
	if m.rulesManaged {
		hint = "perm rules are queued · " + m.keyHint("device.copy-rules") + " to copy"
	}
	pick := func(it actionItem) tea.Cmd {
		if it.nixos {
			m.queueRule(ruleKey(dev), guard.NixOSRule(dev, it.status))
			return overlay.Close
		}
		return tea.Sequence(overlay.Close, doAction(dev.ID, it.fn, it.permanent))
	}
	return overlay.Open(newActions(dev, m.rulesManaged, hint, pick))
}

// ruleKey identifies the physical device the pending rule is for: its VID:PID
// alone is shared by every unit of that model, so two different units of the
// same model must not collide into a single queued rule.
func ruleKey(dev guard.Device) string {
	if dev.Hash == "" {
		return dev.VidPid
	}
	return dev.VidPid + "#" + dev.Hash
}

func (m *Model) queueRule(key, rule string) {
	if i := slices.IndexFunc(m.pendingRules, func(r pendingRule) bool { return r.key == key }); i >= 0 {
		m.pendingRules[i].rule = rule
		return
	}
	m.pendingRules = append(m.pendingRules, pendingRule{key: key, rule: rule})
}

func (m *Model) applyAction(dev guard.Device, fn func(int, bool) error, perm bool, status guard.Status) tea.Cmd {
	if perm && m.rulesManaged {
		m.queueRule(ruleKey(dev), guard.NixOSRule(dev, status))
		return nil
	}
	return doAction(dev.ID, fn, perm)
}

func deviceIDs(devs []guard.Device) []int {
	out := make([]int, len(devs))
	for i, d := range devs {
		out[i] = d.ID
	}
	return out
}

func fetchDevices() tea.Msg {
	devices, err := guard.ListDevices()
	if err != nil {
		return actionMsg{err: err}
	}
	return devicesMsg(devices)
}

func fetchDaemonStatus() tea.Msg { return daemonStatusMsg(guard.DaemonStatus()) }

func fetchDefaultPolicy() tea.Msg { return defaultPolicyMsg(guard.DefaultPolicy()) }

func actionErr(err error) tea.Cmd {
	if errors.Is(err, guard.ErrReadOnly) {
		return notify.Error("Rules file is not writable: permanent changes are not supported.")
	}
	return notify.Err(err)
}

func doAction(id int, fn func(int, bool) error, permanent bool) tea.Cmd {
	return func() tea.Msg { return actionMsg{err: fn(id, permanent)} }
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
