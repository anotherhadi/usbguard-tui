package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
	"github.com/anotherhadi/usbguard-tui/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-version" || os.Args[1] == "-v") {
		fmt.Println("usbguard-tui", version)
		return
	}

	if err := guard.Check(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if fm, ok := m.(ui.Model); ok {
		if rules := fm.PendingRules(); len(rules) > 0 {
			fmt.Println("# Add to your NixOS configuration:")
			fmt.Println("services.usbguard.rules = lib.mkAfter ''")
			for _, rule := range rules {
				fmt.Println(" ", rule)
			}
			fmt.Println("'';")
		}
	}
}
