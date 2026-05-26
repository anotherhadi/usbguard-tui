package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

var statusColors = map[guard.Status]color.Color{
	guard.Allowed:  ilovetui.S.Success,
	guard.Blocked:  ilovetui.S.Error,
	guard.Rejected: ilovetui.S.Warning,
}

var statusColorsSelected = map[guard.Status]color.Color{
	guard.Allowed:  ilovetui.S.Success,
	guard.Blocked:  ilovetui.S.Error,
	guard.Rejected: ilovetui.S.Warning,
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ilovetui.S.Primary).
			PaddingLeft(1)

	daemonActiveStyle = lipgloss.NewStyle().Foreground(ilovetui.S.Success)
	daemonOtherStyle  = lipgloss.NewStyle().Foreground(ilovetui.S.Muted)

	mutedStyle = lipgloss.NewStyle().Foreground(ilovetui.S.Muted)

	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ilovetui.S.Primary).
			Padding(1, 3)

	popupTitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)

	warnStyle = lipgloss.NewStyle().Foreground(ilovetui.S.Warning)
)
