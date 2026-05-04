package ui

import "charm.land/bubbles/v2/key"

type listKeyMap struct {
	Open    key.Binding
	Filter  key.Binding
	Refresh key.Binding
	Quit    key.Binding
	Help    key.Binding
	// shown only in full help
	Allow      key.Binding
	AllowPerm  key.Binding
	Block      key.Binding
	BlockPerm  key.Binding
	Reject     key.Binding
	RejectPerm key.Binding
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Filter, k.Refresh, k.Quit, k.Help}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Filter, k.Refresh, k.Quit},
		{k.Allow, k.AllowPerm, k.Block, k.BlockPerm, k.Reject, k.RejectPerm},
	}
}

var listKeys = listKeyMap{
	Open:       key.NewBinding(key.WithKeys("enter", "tab", "l"), key.WithHelp("enter/l", "select action")),
	Filter:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Refresh:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Quit:       key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc", "quit")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
	Allow:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "allow")),
	AllowPerm:  key.NewBinding(key.WithKeys("A"), key.WithHelp("A", "allow (perm)")),
	Block:      key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "block")),
	BlockPerm:  key.NewBinding(key.WithKeys("B"), key.WithHelp("B", "block (perm)")),
	Reject:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "reject")),
	RejectPerm: key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "reject (perm)")),
}

var cancelKey = key.NewBinding(key.WithKeys("esc", "q", "ctrl+c"), key.WithHelp("esc/q", "cancel"))
