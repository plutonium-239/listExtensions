package scrollinglist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down key.Binding
	Up   key.Binding

	GotoBottom key.Binding
	GotoTop    key.Binding

	Quit key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
	ShowHideHelp  key.Binding
	ToggleFooter  key.Binding
}

// FullHelp implements help.KeyMap.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Down,
			k.Up,
			k.GotoBottom,
			k.GotoTop,
		},
		{
			k.ToggleFooter,
			k.ShowHideHelp,
			k.Quit,
			k.CloseFullHelp,
		},
	}
}

// ShortHelp implements help.KeyMap.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Down,
		k.Up,
		k.Quit,
		k.ShowFullHelp,
	}
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
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
		ToggleFooter: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle footer"),
		),
		ShowHideHelp: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "show/hide help"),
		),
	}
}
