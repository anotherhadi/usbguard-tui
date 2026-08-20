package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/anotherhadi/ilovetui/app"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/notification"
	"github.com/anotherhadi/usbguard-tui/ui"
)

type appModel struct {
	core   ui.Model
	modals modal.Model
	notif  notification.Model
}

func newApp() appModel {
	return appModel{
		core:   ui.New(),
		modals: modal.New(),
		notif:  notification.New(),
	}
}

func (a appModel) Init() tea.Cmd {
	return tea.Batch(a.core.Init(), a.modals.Init(), a.notif.Init())
}

func (a appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(app.QuitMsg); ok {
		return a, tea.Quit
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
	bg := a.core.View()
	bg = a.modals.Render(bg)
	bg = a.notif.Render(bg)
	return tea.View{
		Content:     bg,
		AltScreen:   true,
		WindowTitle: "USBGuard TUI",
		MouseMode:   tea.MouseModeCellMotion,
	}
}
