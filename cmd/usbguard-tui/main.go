// Command usbguard-tui is a terminal UI for managing USB devices via usbguard.
package main

import (
	"errors"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	fconfig "github.com/anotherhadi/foam/config"
	"github.com/anotherhadi/foam/theme"

	"github.com/anotherhadi/usbguard-tui/internal/config"
	"github.com/anotherhadi/usbguard-tui/internal/guard"
	"github.com/anotherhadi/usbguard-tui/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-version" || os.Args[1] == "-v") {
		fmt.Println("usbguard-tui", version)
		return
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "usbguard-tui:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Args[1:])
	if errors.Is(err, fconfig.ErrExit) {
		return nil
	}
	if err != nil {
		return err
	}
	theme.Set(cfg.Theme)

	if err := guard.Check(); err != nil {
		return err
	}
	if _, err := guard.ListDevices(); err != nil {
		return err
	}

	m, err := ui.New(cfg)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m).Run()
	return err
}
