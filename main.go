package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
	"github.com/anotherhadi/usbguard-tui/internal/ui"
)

func main() {
	if err := guard.Check(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
