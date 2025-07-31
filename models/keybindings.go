package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// Shared keymap and help instance
var (
	CommonKeys = NewKeyMap()
	CommonHelp = help.New()
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Search key.Binding
	Quit   key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Search, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Search, k.Quit},
	}
}

func NewKeyMap() KeyMap {
	return KeyMap{
		Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "up")),
		Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "down")),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}
