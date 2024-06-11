// Package scrollinglist provides a scrollalbe list bubbletea component
package scrollinglist

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var DefaultUnfocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.HiddenBorder()).
	BorderTop(false).BorderBottom(false).BorderRight(false).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#C2B8C2"))

var DefaultFocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("197")).
	BorderTop(false).BorderBottom(false).BorderRight(false).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#EE6FF8"))

type ScrollingList struct {
	Items   []fmt.Stringer
	Focused int

	lengths     []int
	preRendered []string

	UnfocusedStyle  lipgloss.Style
	FocusedStyle    lipgloss.Style
	ListAlignment   lipgloss.Position
	GlobalAlignment lipgloss.Position

	FirstVisible       int
	FirstPartial       int
	LastVisible        int
	LastPartial        int
	NumLinesFromBorder int

	Height int
	Width  int

	initialized bool
	keymap      KeyMap

	status string
}

func NewScrollingList() ScrollingList {
	return ScrollingList{
		keymap:             DefaultKeyMap(),
		UnfocusedStyle:     DefaultUnfocusedStyle,
		FocusedStyle:       DefaultFocusedStyle,
		ListAlignment:      lipgloss.Center,
		GlobalAlignment:    lipgloss.Center,
		NumLinesFromBorder: 5,
	}
}

func (sl ScrollingList) Init() tea.Cmd {
	return nil
}

func (sl ScrollingList) View() string {
	if !sl.initialized {
		return "..."
	}
	return lipgloss.JoinVertical(sl.GlobalAlignment,
		lipgloss.JoinVertical(sl.ListAlignment, sl.VisibleLines()...),
		lipgloss.PlaceHorizontal(sl.Width, sl.GlobalAlignment, sl.FooterView()),
	)
}

func (sl *ScrollingList) VisibleLines() []string {
	if len(sl.Items) == 0 {
		return []string{}
	}
	potentialRendered := sl.preRendered[sl.FirstVisible+1 : sl.LastVisible]

	// potentialRendered[sl.Focused] = sl.FocusedStyle.Render(potentialRendered[sl.Focused])
	potentialLength := 0
	styledRendered := make([]string, len(potentialRendered))
	for i := sl.FirstVisible + 1; i < sl.LastVisible; i++ {
		potentialLength += sl.lengths[i]
		styledRendered[i-sl.FirstVisible-1] = sl.styleSingle(i)
	}
	// result := lipgloss.JoinVertical(lipgloss.Left, styledRendered...)
	result := styledRendered

	// reqd_lines := sl.Height - potentialLength
	// if reqd_lines > 0 {
	if sl.FirstPartial >= 0 {
		s_split := strings.Split(sl.styleSingle(sl.FirstVisible), "\n")
		if sl.FirstPartial == 0 {
			sl.FirstPartial = len(s_split)
		}
		s := strings.Join(s_split[len(s_split)-sl.FirstPartial:], "\n")
		// result = lipgloss.JoinVertical(lipgloss.Left, s, result)
		result = slices.Insert(result, 0, s)
	}
	if sl.LastPartial >= 0 {
		s_split := strings.Split(sl.styleSingle(sl.LastVisible), "\n")
		if sl.LastPartial == 0 {
			sl.LastPartial = len(s_split)
		}
		s := strings.Join(s_split[:sl.LastPartial], "\n")
		// result = lipgloss.JoinVertical(lipgloss.Left, result, s)
		result = append(result, s)
	}
	// }
	// potentialRendered = append(potentialRendered, )

	return result
}

func (sl *ScrollingList) styleSingle(index int) string {
	if index == sl.Focused {
		return sl.FocusedStyle.Render(sl.preRendered[index])
	}
	return sl.UnfocusedStyle.Render(sl.preRendered[index])
}

func (sl *ScrollingList) FooterView() string {
	return fmt.Sprintf("Focused:%d, FirstVisible:%d, LastVisible:%d | Status: %s", sl.Focused, sl.FirstVisible, sl.LastVisible, sl.status)
}

func (sl ScrollingList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sl.keymap.Quit):
			// l.Items = make([]fmt.Stringer, 0)
			return sl, tea.Quit
		case key.Matches(msg, sl.keymap.Down):
			sl.Next()
			return sl, cmd
		case key.Matches(msg, sl.keymap.Up):
			sl.Prev()
			return sl, cmd
		case key.Matches(msg, sl.keymap.GotoBottom):
			sl.GotoBottom()
			return sl, cmd
		case key.Matches(msg, sl.keymap.GotoTop):
			sl.GotoTop()
			return sl, cmd
			// default:
			// 	return sl, cmd
		}
	case tea.WindowSizeMsg:
		sl.initialized = true
		// sl.Height = 10
		footerHeight := lipgloss.Height(sl.FooterView())
		sl.Height = msg.Height - footerHeight
		sl.Width = msg.Width
		sl.LastVisible, sl.LastPartial = PrefixSumBreak(sl.lengths, sl.Height)
		return sl, nil
	default:
		var cmd tea.Cmd

		return sl, cmd
	}
	return sl, cmd
}

// Replace all items
func (sl *ScrollingList) SetItems(items []fmt.Stringer) {
	sl.Items = items
	sl.preRendered = make([]string, len(items))
	sl.lengths = make([]int, len(items))
	for i, _ := range items {
		sl.UpdateItem(i)
	}
	sl.FirstPartial = sl.lengths[0]
	if sl.initialized {
		sl.LastVisible, sl.LastPartial = PrefixSumBreak(sl.lengths, sl.Height)
	}
}

// Updates/Replaces the item at one index
func (sl *ScrollingList) UpdateItem(index int) {
	sl.preRendered[index] = sl.Items[index].String()
	sl.lengths[index] = len(strings.Split(sl.preRendered[index], "\n"))
}

func (sl *ScrollingList) Prev() {
	if sl.Focused <= 0 {
		return
	}
	sl.Focused--
	sl.status = "Going Up"
	if sl.Focused < sl.FirstVisible+sl.NumLinesFromBorder && sl.FirstVisible > 0 {
		sl.FirstVisible--
		// TODO: Set LastPartial
		// sl.status = fmt.Sprintf("Execd suffixsum with %d:%d, %d", sl.FirstVisible, sl.LastVisible, sl.lengths[sl.FirstVisible])
		var update int
		update, sl.LastPartial = SuffixSumBreak(sl.lengths[sl.FirstVisible:sl.LastVisible], sl.lengths[sl.FirstVisible])
		sl.LastVisible -= update
		// sl.status += fmt.Sprintf(", final LastVisible = %d", sl.LastVisible)
	}
}

func (sl *ScrollingList) Next() {
	if sl.Focused >= len(sl.Items)-1 {
		return
	}
	sl.Focused++
	sl.status = "Going Down"
	if sl.Focused > sl.LastVisible-sl.NumLinesFromBorder && sl.LastVisible < len(sl.Items)-1 {
		sl.LastVisible++
		// sl.status = fmt.Sprintf("Execd prefixsum with %d:%d, %d", sl.FirstVisible, sl.LastVisible, sl.lengths[sl.LastVisible])
		var update int
		update, sl.FirstPartial = PrefixSumBreak(sl.lengths[sl.FirstVisible:sl.LastVisible], sl.lengths[sl.LastVisible])
		sl.FirstVisible += 1 + update
		// sl.status += fmt.Sprintf(", final FirstVisible = %d", sl.FirstVisible)
	}
}

func (sl *ScrollingList) GotoBottom() {
	sl.LastVisible = len(sl.Items) - 1
	sl.LastPartial = 0
	sl.status = "At Bottom"
	var update int
	update, sl.FirstPartial = SuffixSumBreak(sl.lengths, sl.Height)
	sl.FirstVisible = len(sl.lengths) - update
	sl.Focused = sl.LastVisible - sl.NumLinesFromBorder
}

func (sl *ScrollingList) GotoTop() {
	sl.FirstVisible = 0
	sl.FirstPartial = 0
	sl.status = "At Top"
	// var update int
	sl.LastVisible, sl.LastPartial = PrefixSumBreak(sl.lengths, sl.Height)
	sl.Focused = sl.FirstPartial + sl.NumLinesFromBorder
}
