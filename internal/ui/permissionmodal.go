package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/style"
)

type permissionModal struct{}

func newPermissionModal() tea.Model { return permissionModal{} }

func permissionModalStyle() modal.Styles {
	s := modal.DefaultStyles()
	s.Title = s.Title.Foreground(style.S.Error)
	s.Border = s.Border.BorderForeground(style.S.Error)
	return s
}

func (permissionModal) Init() tea.Cmd { return nil }

func (m permissionModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg, tea.MouseClickMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (permissionModal) View() tea.View {
	textStyle := lipgloss.NewStyle().Foreground(style.S.Text)
	hintStyle := lipgloss.NewStyle().Foreground(style.S.Muted)

	parts := []string{
		textStyle.Render("Can't reach the usbguard daemon"),
		textStyle.Render("(permission denied)."),
		textStyle.Render("Join the usbguard group, or use sudo."),
		hintStyle.Render("press any key to quit"),
	}
	return tea.NewView(strings.Join(parts, "\n"))
}
