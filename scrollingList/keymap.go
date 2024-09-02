package scrollinglist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down key.Binding
	Up   key.Binding

	PageDown key.Binding
	PageUp   key.Binding

	GotoBottom key.Binding
	GotoTop    key.Binding

	Quit       key.Binding
	Edit       key.Binding
	CancelEdit key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
	ShowHideHelp  key.Binding
	ToggleFooter  key.Binding
	ToggleTitle   key.Binding
}

// FullHelp implements help.KeyMap.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Down,
			k.Up,
			k.PageDown,
			k.PageUp,
			k.GotoBottom,
			k.GotoTop,
		},
		{
			k.Edit,
			k.CancelEdit,
			k.ToggleTitle,
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
		// doesnt work without specifying keys
		key.NewBinding(key.WithKeys("up/down"), key.WithHelp("↑/↓", "up/down")),
		k.Edit,
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

		PageDown: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "page up"),
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
			key.WithKeys("q", "Q", "ctrl+c", "ctrl+C"),
			key.WithHelp("q/C^c", "quit"),
		),
		Edit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↩", "edit/save"),
		),
		CancelEdit: key.NewBinding(
			key.WithKeys("escape", "esc"),
			key.WithHelp("esc", "cancel edit"),
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
		ToggleTitle: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T(shift+t)", "toggle title"),
		),
		ShowHideHelp: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "show/hide help"),
		),
	}
}
