// deviceList is the device list: fuzzy-filtered, one row per device.
package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/foam/commands"
	"github.com/anotherhadi/foam/draw"
	"github.com/anotherhadi/foam/fuzzy"
	"github.com/anotherhadi/foam/input"
	"github.com/anotherhadi/foam/list"
	"github.com/anotherhadi/foam/theme"

	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type deviceList struct {
	all       []guard.Device
	shown     []int // indexes into all
	list      list.Model
	filter    input.Model
	filtering bool
}

func newDeviceList() *deviceList {
	m := &deviceList{filter: input.New("/ ")}
	m.filter.Placeholder = "filter by name, id or status"
	return m
}

func (m *deviceList) SetDevices(devs []guard.Device) {
	m.all = devs
	m.refilter()
}

func (m *deviceList) refilter() {
	sel, ok := m.Selected()
	texts := make([]string, len(m.all))
	for i, d := range m.all {
		texts[i] = d.Name + " " + d.VidPid + " " + string(d.Status)
	}
	m.shown = fuzzy.Filter(m.filter.Value(), texts)
	m.list.SetLen(len(m.shown))
	if ok {
		for i, k := range m.shown {
			if m.all[k].ID == sel.ID {
				m.list.Select(i)
				return
			}
		}
	}
	m.list.Select(m.list.Selected())
}

func (m *deviceList) Selected() (guard.Device, bool) {
	if len(m.shown) == 0 {
		return guard.Device{}, false
	}
	return m.all[m.shown[m.list.Selected()]], true
}

// Shown ignores the filter. The slice belongs to the list: do not modify it.
func (m *deviceList) Shown() []guard.Device {
	out := make([]guard.Device, len(m.shown))
	for i, k := range m.shown {
		out[i] = m.all[k]
	}
	return out
}

func (m *deviceList) Filter() string { return m.filter.Value() }

func (m *deviceList) Capturing() bool { return m.filtering }

func (m *deviceList) Input(msg tea.KeyPressMsg) (bool, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", "ctrl+c":
		m.filtering = false
		return true, nil
	}
	if !m.filter.Update(msg) {
		return false, nil
	}
	m.refilter()
	return true, nil
}

func (m *deviceList) Paste(s string) tea.Cmd {
	if !m.filtering {
		return nil
	}
	m.filter.Update(tea.PasteMsg{Content: s})
	m.refilter()
	return nil
}

func (m *deviceList) Commands() []commands.Command {
	return append(m.list.Commands(),
		commands.Command{ID: "device.filter", Help: "filter", Run: func() tea.Cmd { m.filtering = true; return nil }},
		commands.Command{ID: "device.clear-filter", Title: "Clear the filter", When: func() bool { return m.filter.Value() != "" },
			Run: func() tea.Cmd { m.filter.Reset(); m.refilter(); return nil }},
	)
}

// headerHeight is the number of lines the filter box takes above the rows.
func (m *deviceList) headerHeight() int {
	if m.filtering || m.filter.Value() != "" {
		return filterBoxHeight
	}
	return 0
}

// Mouse: a click on the already-selected row opens it only if focused.
func (m *deviceList) Mouse(ev tea.MouseMsg, focused bool) (dev guard.Device, open bool) {
	mm := ev.Mouse()
	if _, click := ev.(tea.MouseClickMsg); click && mm.Y >= 0 && mm.Y < m.headerHeight() {
		m.filtering = true
		return guard.Device{}, false
	}
	mm.Y -= m.headerHeight()
	switch ev.(type) {
	case tea.MouseClickMsg:
		if i, activate := m.list.Reselect(tea.MouseClickMsg(mm), focused); i >= 0 {
			return m.all[m.shown[i]], activate
		}
	case tea.MouseWheelMsg:
		m.list.Mouse(tea.MouseWheelMsg(mm))
	}
	return guard.Device{}, false
}

// filterBoxHeight is the filter box's border, its input row, and its border.
const filterBoxHeight = 3

func (m *deviceList) View(w, h int, focused bool) string {
	var lines []string
	if m.filtering || m.filter.Value() != "" {
		f := m.filter
		if !m.filtering {
			f.Update(tea.KeyPressMsg{Code: tea.KeyHome})
		}
		box := draw.Frame(f.View(max(w-2, 0), m.filtering), w, filterBoxHeight,
			draw.FrameOpts{Title: "Filter", Active: focused && m.filtering})
		lines = append(lines, strings.Split(box, "\n")...)
	}
	m.list.SetHeight(max(h-len(lines), 0))
	start, end := m.list.Visible()
	for i := start; i < end; i++ {
		d := m.all[m.shown[i]]
		row := deviceRow(d, w)
		if i == m.list.Selected() {
			row = draw.Selected(row, w, focused)
		}
		lines = append(lines, row)
	}
	switch {
	case len(m.all) == 0:
		lines = append(lines, theme.Subtle().Render(" No USB device"))
	case len(m.shown) == 0:
		lines = append(lines, theme.Subtle().Render(" No match"))
	}
	return draw.Fit(strings.Join(lines, "\n"), w, h)
}

func deviceRow(d guard.Device, w int) string {
	clr := statusColor(d.Status)
	name := deviceIcon(d.Name) + d.Name
	right := fmt.Sprintf("id:%-3d  %-9s  %-6s  %s", d.ID, d.VidPid, string(d.Status), ruleStateIndicators[d.RuleState])
	rightW := ansi.StringWidth(right)
	nameW := max(w-rightW-2, 1)
	pad := max(w-nameW-rightW-1, 1)
	line := " " + clr.Render(draw.Pad(ansi.Truncate(name, nameW, "…"), nameW)) + strings.Repeat(" ", pad) + theme.Muted().Render(right)
	return ansi.Truncate(line, w, "")
}
