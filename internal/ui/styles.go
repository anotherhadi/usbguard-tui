package ui

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/foam/theme"

	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

func statusColor(s guard.Status) theme.Color {
	switch s {
	case guard.Allowed:
		return theme.Success()
	case guard.Blocked:
		return theme.Danger()
	case guard.Rejected:
		return theme.Warning()
	}
	return theme.Muted()
}

var ruleStateIndicators = map[guard.RuleState]string{
	guard.RulePermanent: "● perm",
	guard.RuleTemporary: "○ tmp",
	guard.RuleDefault:   "· default",
}

func infoLabel(s string) string {
	return lipgloss.NewStyle().Foreground(theme.Muted().Value()).PaddingLeft(1).Width(16).Render(s)
}

func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}
