package ui

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/foam/commands"
	"github.com/anotherhadi/foam/help"
	"github.com/anotherhadi/foam/notify"
	"github.com/anotherhadi/foam/overlay"
	"github.com/anotherhadi/foam/palette"

	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

func (m *Model) global() []commands.Command {
	return []commands.Command{
		{ID: "app.quit", Title: "Quit", Run: func() tea.Cmd { return tea.Quit }},
		{ID: "app.help", Title: "Show key bindings", Help: "help", Run: func() tea.Cmd {
			return overlay.Open(help.New(m.reg.Groups(m.path())))
		}},
		{ID: "app.palette", Title: "Command palette", Run: func() tea.Cmd {
			return overlay.Open(palette.New(m.reg.Active(m.path()), "app.palette"))
		}},
		{ID: "app.refresh", Title: "Refresh the device list", Help: "refresh", Run: func() tea.Cmd {
			return tea.Batch(fetchDevices, fetchDaemonStatus, fetchDefaultPolicy)
		}},
	}
}

func (m *Model) deviceCommands() []commands.Command {
	selected := func() bool { _, ok := m.devices.Selected(); return ok }
	shown := func() bool { return len(m.devices.Shown()) > 0 }
	with := func(fn func(int, bool) error, perm bool, status guard.Status) func() tea.Cmd {
		return func() tea.Cmd {
			dev, ok := m.devices.Selected()
			if !ok {
				return nil
			}
			return m.applyAction(dev, fn, perm, status)
		}
	}
	return append(m.devices.Commands(),
		commands.Command{ID: "device.open", Title: "Open the action menu", Help: "select", When: selected, Run: func() tea.Cmd {
			dev, ok := m.devices.Selected()
			if !ok {
				return nil
			}
			return m.openActions(dev)
		}},
		commands.Command{ID: "device.allow", Title: "Allow the device", Help: "allow", When: selected, Run: with(guard.AllowDevice, false, guard.Allowed)},
		commands.Command{ID: "device.allow-perm", Title: "Allow the device permanently", Help: "allow (perm)", When: selected, Run: with(guard.AllowDevice, true, guard.Allowed)},
		commands.Command{ID: "device.block", Title: "Block the device", Help: "block", When: selected, Run: with(guard.BlockDevice, false, guard.Blocked)},
		commands.Command{ID: "device.block-perm", Title: "Block the device permanently", Help: "block (perm)", When: selected, Run: with(guard.BlockDevice, true, guard.Blocked)},
		commands.Command{ID: "device.reject", Title: "Reject the device", Help: "reject", When: selected, Run: with(guard.RejectDevice, false, guard.Rejected)},
		commands.Command{ID: "device.reject-perm", Title: "Reject the device permanently", Help: "reject (perm)", When: selected, Run: with(guard.RejectDevice, true, guard.Rejected)},
		commands.Command{ID: "device.allow-all", Title: "Allow every device shown", When: shown, Run: func() tea.Cmd {
			return doBulkAction(deviceIDs(m.devices.Shown()), guard.AllowDevice, false)
		}},
		commands.Command{ID: "device.allow-all-perm", Title: "Allow every device shown, permanently", When: shown, Run: m.allowAllPerm},
		commands.Command{ID: "device.copy-rules", Title: "Copy the pending rules, else the current state", Help: "copy rules", Run: m.copyRules},
	)
}

func (m *Model) allowAllPerm() tea.Cmd {
	devs := m.devices.Shown()
	if m.rulesManaged {
		for _, d := range devs {
			m.queueRule(ruleKey(d), guard.NixOSRule(d, guard.Allowed))
		}
		return nil
	}
	return doBulkAction(deviceIDs(devs), guard.AllowDevice, true)
}

func (m *Model) copyRules() tea.Cmd {
	pending := m.pendingRules
	if len(pending) == 0 {
		pending = mergeCurrentStateRules(nil, m.devices.Shown())
	}
	if len(pending) == 0 {
		return notify.Warning("No rules to copy")
	}
	rules := make([]string, len(pending))
	for i, r := range pending {
		rules[i] = r.rule
	}
	text := formatRulesForClipboard(rules, m.rulesManaged)
	m.pendingRules = nil
	return tea.Batch(tea.SetClipboard(text), notify.Success(plural(len(rules), "rule", "rules")+" copied to clipboard"))
}

func mergeCurrentStateRules(pending []pendingRule, devices []guard.Device) []pendingRule {
	for _, d := range devices {
		key := ruleKey(d)
		rule := guard.NixOSRule(d, d.Status)
		if i := slices.IndexFunc(pending, func(r pendingRule) bool { return r.key == key }); i >= 0 {
			pending[i].rule = rule
		} else {
			pending = append(pending, pendingRule{key: key, rule: rule})
		}
	}
	return pending
}

func formatRulesForClipboard(rules []string, nixos bool) string {
	if !nixos {
		return strings.Join(rules, "\n")
	}
	var b strings.Builder
	b.WriteString("# Add to your NixOS configuration:\n")
	b.WriteString("services.usbguard.rules = lib.mkAfter ''\n")
	for _, r := range rules {
		b.WriteString("  " + r + "\n")
	}
	b.WriteString("'';")
	return b.String()
}

// path is the chain of scopes: the global one, then the devices.
func (m *Model) path() []commands.Scope {
	return []commands.Scope{
		{Name: "usbguard-tui", Commands: m.global()},
		{Name: "Devices", Commands: m.deviceCommands()},
	}
}

// allScopes lists every scope, for checking the keys at startup.
func (m *Model) allScopes() []commands.Scope {
	return []commands.Scope{
		{Name: "usbguard-tui", Commands: m.global()},
		{Name: "Devices", Commands: m.deviceCommands()},
		help.New(nil).Scope(),
		palette.New(nil, "").Scope(),
		newActions(guard.Device{}, false, "", nil).Scope(),
	}
}
