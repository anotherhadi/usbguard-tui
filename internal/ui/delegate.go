package ui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

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
		clr = ilovetui.S.Muted
	}

	var nameStyle, descStyle lipgloss.Style
	if selected {
		nameStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ilovetui.S.Primary).
			Foreground(clr).
			Bold(true).
			PaddingLeft(1)
		descStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ilovetui.S.Primary).
			Foreground(ilovetui.S.Muted).
			PaddingLeft(1)
	} else {
		nameStyle = lipgloss.NewStyle().Foreground(clr).PaddingLeft(2)
		descStyle = lipgloss.NewStyle().Foreground(ilovetui.S.Muted).PaddingLeft(2)
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

func (d actionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	a, ok := item.(actionItem)
	if !ok {
		return
	}
	if index == m.Index() {
		clr, ok := statusColorsSelected[a.status]
		if !ok {
			clr = ilovetui.S.Primary
		}
		fmt.Fprintf(w, "  %s", lipgloss.NewStyle().Bold(true).Foreground(clr).Render("> "+a.label))
	} else {
		fmt.Fprintf(w, "    %s", a.label)
	}
}
