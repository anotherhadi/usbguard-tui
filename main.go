package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
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

	p := tea.NewProgram(newApp())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
