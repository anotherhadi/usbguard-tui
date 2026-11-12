// Package config loads usbguard-tui's configuration: foam's shared defaults,
// then usbguard-tui's own, then the user's files, the environment and flags.
package config

import (
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/anotherhadi/foam"
	"github.com/anotherhadi/foam/commands"
	fconfig "github.com/anotherhadi/foam/config"
	"github.com/anotherhadi/foam/theme"
)

//go:embed default.yaml
var Defaults []byte

type Config struct {
	Theme   theme.Theme       `koanf:",squash"`
	Keys    commands.Bindings `koanf:"keys"`
	Refresh time.Duration     `koanf:"refresh"`
}

const usage = `usbguard-tui: a terminal UI for usbguard

Usage: usbguard-tui [options]`

// Load returns fconfig.ErrExit after --help or --print-config.
func Load(args []string) (*Config, error) {
	cfg, _, err := fconfig.Load[Config](fconfig.Options{
		Name:  "usbguard-tui",
		Usage: usage,
		Base:  [][]byte{foam.Defaults, Defaults},
		Files: []string{
			filepath.Join(fconfig.Dir("foam"), "config.yaml"),
			filepath.Join(fconfig.Dir("usbguard-tui"), "config.yaml"),
		},
	}, args)
	if err != nil {
		return nil, err
	}
	if err := checkKeys(cfg.Keys); err != nil {
		return nil, fmt.Errorf("invalid config:\n%w", err)
	}
	return cfg, nil
}

// checkKeys rejects bindings for unknown commands, otherwise a typo in keys
// would be silently ignored.
func checkKeys(b commands.Bindings) error {
	base, err := fconfig.LoadBytes[Config](foam.Defaults, Defaults)
	if err != nil {
		return err
	}
	known, err := base.Keys.Flatten()
	if err != nil {
		return err
	}
	ids, err := b.Flatten()
	if err != nil {
		return err
	}
	var errs []error
	for id := range ids {
		if _, ok := known[id]; !ok {
			errs = append(errs, fmt.Errorf("keys.%s: no such command", id))
		}
	}
	slices.SortFunc(errs, func(a, b error) int { return strings.Compare(a.Error(), b.Error()) })
	return errors.Join(errs...)
}
