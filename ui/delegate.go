package ui

import (
	"fmt"
	"io"
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"

	"github.com/anotherhadi/ilovetui/style"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

func deviceZoneID(dev guard.Device) string {
	return "device-" + strconv.Itoa(dev.ID)
}

var ruleStateIndicators = map[guard.RuleState]string{
	guard.RulePermanent: "● perm",
	guard.RuleTemporary: "○ tmp",
	guard.RuleDefault:   "· default",
}

type deviceDelegate struct{}

func (d deviceDelegate) Height() int                             { return 2 }
func (d deviceDelegate) Spacing() int                            { return 0 }
func (d deviceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d deviceDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	dev, ok := item.(guard.Device)
	if !ok {
		return
	}

	selected := index == m.Index()

	colorMap := statusColors
	if selected {
		colorMap = statusColorsSelected
	}
	clr, ok := colorMap[dev.Status]
	if !ok {
		clr = style.S.Muted
	}

	var nameStyle, descStyle lipgloss.Style
	if selected {
		nameStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.S.Primary).
			Foreground(clr).
			Bold(true).
			PaddingLeft(1)
		descStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.S.Primary).
			Foreground(style.S.Muted).
			PaddingLeft(1)
	} else {
		nameStyle = lipgloss.NewStyle().Foreground(clr).PaddingLeft(2)
		descStyle = lipgloss.NewStyle().Foreground(style.S.Muted).PaddingLeft(2)
	}

	name := dev.Name
	if icon := deviceIcon(dev.Name); icon != "" {
		name = icon + " " + name
	}

	permIndicator := ruleStateIndicators[dev.RuleState]
	width := m.Width()
	line1 := clampToWidth(nameStyle.Render(name), width)
	line2 := clampToWidth(descStyle.Render(fmt.Sprintf("id:%-3d  %s  %s  %s", dev.ID, dev.VidPid, string(dev.Status), permIndicator)), width)
	fmt.Fprint(w, zone.Mark(deviceZoneID(dev), line1+"\n"+line2))
}

type actionItem struct {
	label     string
	fn        func(int, bool) error
	permanent bool
	status    guard.Status
	nixos     bool
}

func (a actionItem) Title() string       { return a.label }
func (a actionItem) Description() string { return "" }
func (a actionItem) FilterValue() string { return a.label }

type actionDelegate struct{}

func (d actionDelegate) Height() int                             { return 1 }
func (d actionDelegate) Spacing() int                            { return 0 }
func (d actionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func actionZoneID(index int) string {
	return "action-" + strconv.Itoa(index)
}

func (d actionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	a, ok := item.(actionItem)
	if !ok {
		return
	}
	var content string
	if index == m.Index() {
		clr, ok := statusColorsSelected[a.status]
		if !ok {
			clr = style.S.Primary
		}
		content = fmt.Sprintf("  %s", lipgloss.NewStyle().Bold(true).Foreground(clr).Render("> "+a.label))
	} else {
		content = fmt.Sprintf("    %s", a.label)
	}
	fmt.Fprint(w, zone.Mark(actionZoneID(index), content))
}
