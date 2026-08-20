package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui/style"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

var statusColors = map[guard.Status]color.Color{
	guard.Allowed:  style.S.Success,
	guard.Blocked:  style.S.Error,
	guard.Rejected: style.S.Warning,
}

var statusColorsSelected = map[guard.Status]color.Color{
	guard.Allowed:  style.S.Success,
	guard.Blocked:  style.S.Error,
	guard.Rejected: style.S.Warning,
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(style.S.Primary).
			PaddingLeft(1)

	daemonActiveStyle = lipgloss.NewStyle().Foreground(style.S.Success)
	daemonOtherStyle  = lipgloss.NewStyle().Foreground(style.S.Muted)

	mutedStyle = lipgloss.NewStyle().Foreground(style.S.Muted)

	infoLabelStyle = lipgloss.NewStyle().
			Foreground(style.S.Muted).
			PaddingLeft(1).
			Width(16)

	warnStyle = lipgloss.NewStyle().Foreground(style.S.Warning)
)
