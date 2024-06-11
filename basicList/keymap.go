package basiclist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down key.Binding
	Up   key.Binding

	GotoBottom key.Binding
	GotoTop    key.Binding

	Quit key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move 1 line down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move 1 line up"),
		),

		GotoBottom: key.NewBinding(
			key.WithKeys("ctrl+down"),
			key.WithHelp("C^↓", "goto bottom"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("ctrl+up"),
			key.WithHelp("C^↑", "goto top"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/C^c", "quit"),
		),

		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close full help"),
		),
	}
}
