package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/anotherhadi/ilovetui/app"
	"github.com/anotherhadi/ilovetui/minsize"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/notification"
	"github.com/anotherhadi/usbguard-tui/ui"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	minWidth  = 40
	minHeight = 15
)

type appModel struct {
	core          ui.Model
	modals        modal.Model
	notif         notification.Model
	minsize       minsize.Model
	width, height int
}

func newApp() appModel {
	zone.NewGlobal()
	return appModel{
		core:    ui.New(),
		modals:  modal.New(),
		notif:   notification.New(),
		minsize: minsize.New(minWidth, minHeight),
	}
}

func (a appModel) Init() tea.Cmd {
	return tea.Batch(a.core.Init(), a.modals.Init(), a.notif.Init())
}

func (a appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(app.QuitMsg); ok {
		return a, tea.Quit
	}
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		a.width, a.height = sizeMsg.Width, sizeMsg.Height
	}

	a.core.ModalOpen = a.modals.Open()
	var coreCmd tea.Cmd
	a.core, coreCmd = a.core.Update(msg)

	var modalCmd, notifCmd tea.Cmd
	a.modals, modalCmd = a.modals.Update(msg)
	a.notif, notifCmd = a.notif.Update(msg)

	return a, tea.Batch(coreCmd, modalCmd, notifCmd)
}

func (a appModel) View() tea.View {
	if !a.minsize.Fits(a.width, a.height) {
		return tea.View{
			Content:     a.minsize.View(a.width, a.height),
			AltScreen:   true,
			WindowTitle: "USBGuard TUI",
		}
	}

	bg := lipgloss.Place(a.width, a.height, lipgloss.Left, lipgloss.Top, a.core.View())
	bg = a.modals.Render(bg)
	bg = a.notif.Render(bg)
	bg = zone.Scan(bg)
	return tea.View{
		Content:     bg,
		AltScreen:   true,
		WindowTitle: "USBGuard TUI",
		MouseMode:   tea.MouseModeCellMotion,
	}
}
