// actions is an overlay offering the allow/block/reject actions for one
// device, replacing the old action modal.
package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/foam/commands"
	"github.com/anotherhadi/foam/draw"
	"github.com/anotherhadi/foam/list"
	"github.com/anotherhadi/foam/overlay"
	"github.com/anotherhadi/foam/theme"

	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type actionItem struct {
	label     string
	fn        func(int, bool) error
	permanent bool
	status    guard.Status
	nixos     bool // queue a NixOS rule instead of applying it
}

// actionItems lists the device actions. With rulesManaged, the permanent
// variants queue a NixOS rule instead of writing a read-only rules file.
func actionItems(rulesManaged bool) []actionItem {
	if rulesManaged {
		return []actionItem{
			{"allow", guard.AllowDevice, false, guard.Allowed, false},
			{"allow (perm)", nil, true, guard.Allowed, true},
			{"block", guard.BlockDevice, false, guard.Blocked, false},
			{"block (perm)", nil, true, guard.Blocked, true},
			{"reject", guard.RejectDevice, false, guard.Rejected, false},
			{"reject (perm)", nil, true, guard.Rejected, true},
		}
	}
	return []actionItem{
		{"allow", guard.AllowDevice, false, guard.Allowed, false},
		{"allow (permanent)", guard.AllowDevice, true, guard.Allowed, false},
		{"block", guard.BlockDevice, false, guard.Blocked, false},
		{"block (permanent)", guard.BlockDevice, true, guard.Blocked, false},
		{"reject", guard.RejectDevice, false, guard.Rejected, false},
		{"reject (permanent)", guard.RejectDevice, true, guard.Rejected, false},
	}
}

type actions struct {
	dev          guard.Device
	rulesManaged bool
	hint         string
	items        []actionItem
	list         list.Model
	pick         func(actionItem) tea.Cmd
}

// newActions offers the actions for dev; pick gets the chosen one, and its
// command runs once the overlay is closed. hint, when set, is shown below
// the list (the NixOS perm-rules notice).
func newActions(dev guard.Device, rulesManaged bool, hint string, pick func(actionItem) tea.Cmd) *actions {
	items := actionItems(rulesManaged)
	m := &actions{dev: dev, rulesManaged: rulesManaged, hint: hint, items: items, pick: pick}
	m.list.SetLen(len(items))
	return m
}

func (m *actions) choose() tea.Cmd {
	if len(m.items) == 0 || m.pick == nil {
		return nil
	}
	return m.pick(m.items[m.list.Selected()])
}

func (m *actions) Scope() commands.Scope {
	return commands.Scope{Name: "Actions", Commands: []commands.Command{
		{ID: "actions.close", Help: "close", Run: func() tea.Cmd { return overlay.Close }},
		{ID: "actions.pick", Help: "select", Run: m.choose},
		{ID: "actions.up", Run: func() tea.Cmd { m.list.Move(-1); return nil }},
		{ID: "actions.down", Run: func() tea.Cmd { m.list.Move(1); return nil }},
	}}
}

func (m *actions) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		mm := msg.Mouse()
		mm.Y-- // the border comes before the first row
		if i := m.list.Mouse(tea.MouseClickMsg(mm)); i >= 0 {
			return m.choose()
		}
	case tea.MouseWheelMsg:
		mm := msg.Mouse()
		mm.Y--
		m.list.Mouse(tea.MouseWheelMsg(mm))
	}
	return nil
}

func (m *actions) View(maxW, maxH int) string {
	w := min(32, maxW)
	inner := max(w-2, 0)
	m.list.SetHeight(max(min(len(m.items), maxH-2), 0))

	var lines []string
	start, end := m.list.Visible()
	for i := start; i < end; i++ {
		it := m.items[i]
		row := " " + statusColor(it.status).Render(it.label)
		if i == m.list.Selected() {
			row = draw.Selected(row, inner, true)
		} else {
			row = draw.Pad(row, inner)
		}
		lines = append(lines, row)
	}
	if m.hint != "" {
		lines = append(lines, draw.Rule(inner), theme.Subtle().Render(" "+m.hint))
	}
	return draw.Frame(strings.Join(lines, "\n"), w, len(lines)+2, draw.FrameOpts{Title: m.dev.Name, Active: true})
}
