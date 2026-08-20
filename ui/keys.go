package ui

import "charm.land/bubbles/v2/key"

type listKeyMap struct {
	Open    key.Binding
	Filter  key.Binding
	Refresh key.Binding
	Quit    key.Binding
	Help    key.Binding

	Allow      key.Binding
	AllowPerm  key.Binding
	Block      key.Binding
	BlockPerm  key.Binding
	Reject     key.Binding
	RejectPerm key.Binding

	AllowAll     key.Binding
	AllowAllPerm key.Binding
	CopyRules    key.Binding

	Up        key.Binding
	Down      key.Binding
	GoToStart key.Binding
	GoToEnd   key.Binding
	PrevPage  key.Binding
	NextPage  key.Binding
}

func (k listKeyMap) globalBindings() []key.Binding {
	return []key.Binding{
		k.Open, k.Filter, k.Refresh, k.Quit,
		k.Allow, k.AllowPerm, k.Block, k.BlockPerm, k.Reject, k.RejectPerm,
		k.AllowAll, k.AllowAllPerm, k.CopyRules,
		k.Up, k.Down, k.GoToStart, k.GoToEnd, k.PrevPage, k.NextPage,
	}
}

func (k listKeyMap) shortHelpBindings() []key.Binding {
	return []key.Binding{
		k.Open, k.Allow, k.Block, k.Reject, k.Filter, k.Quit,
	}
}

var listKeys = listKeyMap{
	Open:    key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter/tab", "select action")),
	Filter:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Quit:    key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc", "quit")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),

	Allow:        key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "allow")),
	AllowPerm:    key.NewBinding(key.WithKeys("A"), key.WithHelp("A", "allow (perm)")),
	Block:        key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "block")),
	BlockPerm:    key.NewBinding(key.WithKeys("B"), key.WithHelp("B", "block (perm)")),
	Reject:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "reject")),
	RejectPerm:   key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "reject (perm)")),
	AllowAll:     key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "allow all")),
	AllowAllPerm: key.NewBinding(key.WithKeys("W"), key.WithHelp("W", "allow all (perm)")),
	CopyRules:    key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy rules")),

	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	GoToStart: key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g/home", "go to start")),
	GoToEnd:   key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G/end", "go to end")),
	PrevPage:  key.NewBinding(key.WithKeys("left", "pgup", "h"), key.WithHelp("←/pgup/h", "prev page")),
	NextPage:  key.NewBinding(key.WithKeys("right", "pgdown", "l"), key.WithHelp("→/pgdn/l", "next page")),
}

var cancelKey = key.NewBinding(key.WithKeys("esc", "q", "ctrl+c"), key.WithHelp("esc/q", "cancel"))
