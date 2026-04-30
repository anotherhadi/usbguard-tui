package ui

import (
	"github.com/anotherhadi/usbguard-tui/internal/guard"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorAllowed          = lipgloss.Color("28")
	colorAllowedSelected  = lipgloss.Color("42")
	colorBlocked          = lipgloss.Color("124")
	colorBlockedSelected  = lipgloss.Color("196")
	colorRejected         = lipgloss.Color("130")
	colorRejectedSelected = lipgloss.Color("214")
	colorMuted            = lipgloss.Color("240")
	colorAccent           = lipgloss.Color("99")
)

var statusColors = map[guard.Status]lipgloss.Color{
	guard.Allowed:  colorAllowed,
	guard.Blocked:  colorBlocked,
	guard.Rejected: colorRejected,
}

var statusColorsSelected = map[guard.Status]lipgloss.Color{
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
