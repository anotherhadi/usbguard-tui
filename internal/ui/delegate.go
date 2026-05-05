package ui

import (
	"fmt"
	"image/color"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

// deviceDelegate renders device list items with status colors.
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

	var clr color.Color
	if selected {
		var ok bool
		clr, ok = statusColorsSelected[dev.Status]
		if !ok {
			clr = colorMuted
		}
	} else {
		var ok bool
		clr, ok = statusColors[dev.Status]
		if !ok {
			clr = colorMuted
		}
	}

	var nameStyle, descStyle lipgloss.Style
	if selected {
		nameStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(colorAccent).
			Foreground(clr).
			Bold(true).
			PaddingLeft(1)
		descStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(colorAccent).
			Foreground(colorMuted).
			PaddingLeft(1)
	} else {
		nameStyle = lipgloss.NewStyle().Foreground(clr).PaddingLeft(2)
		descStyle = lipgloss.NewStyle().Foreground(colorMuted).PaddingLeft(2)
	}

	permIndicator := "○ tmp"
	if dev.Permanent {
		permIndicator = "● perm"
	}
	fmt.Fprintf(w, "%s\n%s",
		nameStyle.Render(dev.Name),
		descStyle.Render(fmt.Sprintf("id:%-3d  %s  %s  %s", dev.ID, dev.VidPid, string(dev.Status), permIndicator)),
	)
}

// actionItem represents a device policy action in the select popup.
type actionItem struct {
	label     string
	fn        func(int, bool) error
	permanent bool
	status    guard.Status
}

func (a actionItem) Title() string       { return a.label }
func (a actionItem) Description() string { return "" }
func (a actionItem) FilterValue() string { return a.label }

// actionDelegate renders single-line action items.
type actionDelegate struct{}

func (d actionDelegate) Height() int                             { return 1 }
func (d actionDelegate) Spacing() int                            { return 0 }
func (d actionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d actionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	a, ok := item.(actionItem)
	if !ok {
		return
	}
	if index == m.Index() {
		clr, ok := statusColorsSelected[a.status]
		if !ok {
			clr = colorAccent
		}
		fmt.Fprintf(w, "  %s", lipgloss.NewStyle().Bold(true).Foreground(clr).Render("> "+a.label))
	} else {
		fmt.Fprintf(w, "    %s", a.label)
	}
}
