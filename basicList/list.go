package basiclist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This package defines a basic list bubble that is not paginated and simply extends to the full width
// This is to use in conjunction with a viewport, for an easily scrollable list

var DefaultUnfocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.HiddenBorder()).
	BorderTop(false).BorderBottom(false).BorderRight(false).
	Foreground(lipgloss.Color("#C2B8C2"))

var DefaultFocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("197")).
	BorderTop(false).BorderBottom(false).BorderRight(false).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#EE6FF8"))

type BasicList struct {
	Items   []fmt.Stringer
	Focused int

	UnfocusedStyle lipgloss.Style
	FocusedStyle   lipgloss.Style

	keymap KeyMap
}

func NewBasicList() BasicList {
	return BasicList{
		UnfocusedStyle: DefaultUnfocusedStyle,
		FocusedStyle:   DefaultFocusedStyle,
	}
}

func (l BasicList) Init() tea.Cmd {
	return nil
}

func (l BasicList) View() string {
	itemStrings := make([]string, 0)
	for i, item := range l.Items {
		s := item.String()
		if i == l.Focused {
			s = l.FocusedStyle.Render(s)
		} else {
			s = l.UnfocusedStyle.Render(s)
		}
		itemStrings = append(itemStrings, s)
	}

	// return lipgloss.JoinVertical(lipgloss.Left, itemStrings...)
	return strings.Join(itemStrings, "\n")
}

func (l BasicList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keymap.Down):
			if l.Focused < len(l.Items)-1 {
				l.Focused++
			}
			return l, nil
		case key.Matches(msg, l.keymap.Up):
			if l.Focused > 0 {
				l.Focused--
			}
			return l, nil
		case key.Matches(msg, l.keymap.GotoBottom):
			l.Focused = len(l.Items) - 1
			return l, nil
		case key.Matches(msg, l.keymap.GotoTop):
			l.Focused = 0
			return l, nil
		case key.Matches(msg, l.keymap.Quit):
			return l, tea.Quit
		}
	}
	return l, nil
}

func (l *BasicList) SetItems(items []fmt.Stringer) {
	l.Items = items
}
