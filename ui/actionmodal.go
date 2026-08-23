package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"

	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/style"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

type actionModal struct {
	list         list.Model
	dev          guard.Device
	rulesManaged bool
}

func newActionModal(dev guard.Device, rulesManaged bool) tea.Model {
	return actionModal{
		list:         makeActionList(rulesManaged),
		dev:          dev,
		rulesManaged: rulesManaged,
	}
}

func actionModalStyle(status guard.Status) modal.Styles {
	s := modal.DefaultStyles()
	if clr, ok := statusColors[status]; ok {
		s.Title = s.Title.Foreground(clr)
	}
	return s
}

func (a actionModal) Init() tea.Cmd { return nil }

func (a actionModal) confirmCmd(it actionItem) tea.Cmd {
	if it.nixos {
		rule := guard.NixOSRule(a.dev, it.status)
		key := a.dev.VidPid
		return tea.Batch(modal.Close(), func() tea.Msg { return nixRuleMsg{key: key, rule: rule} })
	}
	return tea.Batch(modal.Close(), doAction(a.dev.ID, it.fn, it.permanent))
}

func (a actionModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, cancelKey):
			return a, modal.Close()
		case key.Matches(msg, listKeys.Open):
			item := a.list.SelectedItem()
			if item == nil {
				return a, nil
			}
			return a, a.confirmCmd(item.(actionItem))
		}

	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			a.list.CursorUp()
		case tea.MouseWheelDown:
			a.list.CursorDown()
		}
		return a, nil

	case tea.MouseClickMsg:
		if msg.Mouse().Button != tea.MouseLeft {
			return a, nil
		}
		for i, item := range a.list.VisibleItems() {
			if zone.Get(actionZoneID(i)).InBounds(msg) {
				a.list.Select(i)
				return a, a.confirmCmd(item.(actionItem))
			}
		}
		return a, nil
	}

	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)
	return a, cmd
}

func (a actionModal) View() tea.View {
	hintStyle := lipgloss.NewStyle().Foreground(style.S.Muted)
	parts := []string{a.list.View(), ""}
	if a.rulesManaged {
		parts = append(parts, hintStyle.Render(fmt.Sprintf("[NixOS: perm rules queued, press %s to copy]", listKeys.CopyRules.Help().Key)))
	}
	parts = append(parts, hintStyle.Render("↑↓ navigate  enter confirm  esc cancel"))
	return tea.NewView(strings.Join(parts, "\n"))
}
