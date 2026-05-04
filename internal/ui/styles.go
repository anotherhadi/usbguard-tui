package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

var (
	colorAllowed          color.Color = lipgloss.Color("28")
	colorAllowedSelected  color.Color = lipgloss.Color("42")
	colorBlocked          color.Color = lipgloss.Color("124")
	colorBlockedSelected  color.Color = lipgloss.Color("196")
	colorRejected         color.Color = lipgloss.Color("130")
	colorRejectedSelected color.Color = lipgloss.Color("214")
	colorMuted            color.Color = lipgloss.Color("240")
	colorAccent           color.Color = lipgloss.Color("99")
)

var statusColors = map[guard.Status]color.Color{
	guard.Allowed:  colorAllowed,
	guard.Blocked:  colorBlocked,
	guard.Rejected: colorRejected,
}

var statusColorsSelected = map[guard.Status]color.Color{
	guard.Allowed:  colorAllowedSelected,
	guard.Blocked:  colorBlockedSelected,
	guard.Rejected: colorRejectedSelected,
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			PaddingLeft(1)

	daemonActiveStyle = lipgloss.NewStyle().Foreground(colorAllowedSelected)
	daemonOtherStyle  = lipgloss.NewStyle().Foreground(colorMuted)

	mutedStyle = lipgloss.NewStyle().Foreground(colorMuted)

	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 3)

	popupTitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)

	keyHintStyle = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(colorRejected)
	errStyle     = lipgloss.NewStyle().Foreground(colorBlocked).Bold(true)
)
