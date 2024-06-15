// Package scrollinglist provides a scrollalbe list bubbletea component
package scrollinglist

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
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
	items   []fmt.Stringer
	focused int

	lengths     []int
	preRendered []string

	// Styles can be customized
	UnfocusedStyle lipgloss.Style
	FocusedStyle   lipgloss.Style

	// Alignment of text within the list
	ListAlignment lipgloss.Position
	// Alignment of the list and footer w.r.t container
	GlobalAlignment lipgloss.Position

	// Footer can be toggled and customized, can be multiline
	ShowFooter bool
	CustomFooter func() string

	help help.Model
	// Help can be toggled
	ShowHelp bool

	firstVisible, lastVisible int
	firstPartial, lastPartial int

	// The number of lines to keep between focused item and screen edges
	NumLinesFromBorder int

	// Size of the whole list (scroll height is handled automatically based on footer/help)
	Width, Height int
	scrollHeight int

	initialized bool
	// KeyMap
	KeyMap KeyMap

	status string
}

func NewScrollingList() ScrollingList {
	return ScrollingList{
		KeyMap:             DefaultKeyMap(),
		UnfocusedStyle:     DefaultUnfocusedStyle,
		FocusedStyle:       DefaultFocusedStyle,
		ListAlignment:      lipgloss.Center,
		GlobalAlignment:    lipgloss.Center,
		NumLinesFromBorder: 5,
		ShowFooter:       true,
		ShowHelp: true,
		help: help.New(),
	}
}

// BubbleTea Init method, does nothing
func (sl ScrollingList) Init() tea.Cmd {
	return nil
}

func (sl *ScrollingList) place(view string) string {
	return lipgloss.PlaceHorizontal(sl.Width, sl.GlobalAlignment, view)
}

// BubbleTea View method, returns the actual final rendering result
func (sl ScrollingList) View() string {

	if !sl.initialized {
		return "..."
	}
	views := make([]string, 0)

	views = append(views, lipgloss.JoinVertical(sl.ListAlignment, sl.VisibleLines()...))
	if sl.ShowFooter {
		// return lipgloss.PlaceHorizontal(sl.Width, sl.GlobalAlignment, lipgloss.JoinVertical(sl.GlobalAlignment, views...))
		views = append(views, sl.place(sl.FooterView()))
	}
	if sl.ShowHelp {
		views = append(views, sl.place(sl.HelpView()))
	}
	if !sl.ShowFooter && !sl.ShowHelp {
		views[0] = sl.place(views[0])
	}
	return lipgloss.JoinVertical(sl.GlobalAlignment, views...)
}

// returns the visible lines, with partial first/last elements s.t. output fits on screen
func (sl *ScrollingList) VisibleLines() []string {
	if len(sl.items) == 0 {
		return []string{}
	}
	potentialRendered := sl.preRendered[sl.firstVisible+1 : sl.lastVisible]

	// potentialRendered[sl.Focused] = sl.FocusedStyle.Render(potentialRendered[sl.Focused])
	potentialLength := 0
	styledRendered := make([]string, len(potentialRendered))
	for i := sl.firstVisible + 1; i < sl.lastVisible; i++ {
		potentialLength += sl.lengths[i]
		styledRendered[i-sl.firstVisible-1] = sl.styleSingle(i)
	}
	// result := lipgloss.JoinVertical(lipgloss.Left, styledRendered...)
	result := styledRendered

	// reqd_lines := sl.Height - potentialLength
	// if reqd_lines > 0 {
	if sl.firstPartial >= 0 {
		s_split := strings.Split(sl.styleSingle(sl.firstVisible), "\n")
		if sl.firstPartial == 0 {
			sl.firstPartial = len(s_split)
		}
		s := strings.Join(s_split[len(s_split)-sl.firstPartial:], "\n")
		// result = lipgloss.JoinVertical(lipgloss.Left, s, result)
		result = slices.Insert(result, 0, s)
	}
	if sl.lastPartial >= 0 {
		s_split := strings.Split(sl.styleSingle(sl.lastVisible), "\n")
		if sl.lastPartial == 0 {
			sl.lastPartial = len(s_split)
		}
		s := strings.Join(s_split[:sl.lastPartial], "\n")
		// result = lipgloss.JoinVertical(lipgloss.Left, result, s)
		result = append(result, s)
	}
	// }
	// potentialRendered = append(potentialRendered, )

	return result
}

func (sl *ScrollingList) styleSingle(index int) string {
	if index == sl.focused {
		return sl.FocusedStyle.Render(sl.preRendered[index])
	}
	return sl.UnfocusedStyle.Render(sl.preRendered[index])
}

// returns the rendered footer
func (sl *ScrollingList) FooterView() string {
	// ShowFooter handling is done in View
	if sl.CustomFooter != nil {
		return sl.CustomFooter()
	}
	return fmt.Sprintf("Focused:%d, FirstVisible:%d, LastVisible:%d | Status: %s | main(w,h): (%d,%d) scrollHeight: %d", 
	sl.focused, sl.firstVisible, sl.lastVisible, sl.status, sl.Width, sl.Height, sl.scrollHeight)
}

// returns the rendered help
func (sl *ScrollingList) HelpView() string {
	return sl.help.View(sl.KeyMap)
}

// BubbleTea Update method, handles key presses and window resize
func (sl ScrollingList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sl.KeyMap.Quit):
			return sl, tea.Quit
		case key.Matches(msg, sl.KeyMap.Down):
			sl.Next()
			return sl, nil
		case key.Matches(msg, sl.KeyMap.Up):
			sl.Prev()
			return sl, nil
		case key.Matches(msg, sl.KeyMap.GotoBottom):
			sl.GotoBottom()
			return sl, nil
		case key.Matches(msg, sl.KeyMap.GotoTop):
			sl.GotoTop()
			return sl, nil
		case key.Matches(msg, sl.KeyMap.ToggleFooter):
			sl.ShowFooter = !sl.ShowFooter
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		case key.Matches(msg, sl.KeyMap.ShowHideHelp):
			sl.ShowHelp = !sl.ShowHelp
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		case key.Matches(msg, sl.KeyMap.ShowFullHelp, sl.KeyMap.CloseFullHelp):
			sl.help.ShowAll = !sl.help.ShowAll
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		}
	case tea.WindowSizeMsg:
		sl.SetSize(msg.Width, msg.Height)
		// TODO: find a way to persist through screen resize, this starts by deciding whether to keep focused closer to top or bottom (or middle?)
		sl.lastVisible, sl.lastPartial = PrefixSumBreak(sl.lengths[:len(sl.lengths)-1], sl.scrollHeight)
		return sl, nil
	}
	return sl, nil
}

// Replace all items
func (sl *ScrollingList) SetItems(items []fmt.Stringer) {
	sl.items = items
	sl.preRendered = make([]string, len(items))
	sl.lengths = make([]int, len(items))
	for i := range items {
		sl.recalculateForIndex(i)
	}
	sl.firstPartial = sl.lengths[0]
	if sl.initialized {
		sl.lastVisible, sl.lastPartial = PrefixSumBreak(sl.lengths[:len(sl.lengths)-1], sl.scrollHeight)
	}
}

// Replace one item at index int with given item
func (sl *ScrollingList) SetItemAt(item fmt.Stringer, index int) {
	sl.items[index] = item
	sl.recalculateForIndex(index)
}

// update internal preRendered and lengths for the item at index
func (sl *ScrollingList) recalculateForIndex(index int) {
	sl.preRendered[index] = sl.items[index].String()
	sl.lengths[index] = len(strings.Split(sl.preRendered[index], "\n"))
}

// Move 1 item up
func (sl *ScrollingList) Prev() {
	if sl.focused <= 0 {
		return
	}
	sl.focused--
	sl.status = "Going Up"
	if sl.focused < sl.firstVisible+sl.NumLinesFromBorder && sl.firstVisible > 0 {
		sl.firstVisible--
		var update int
		update, sl.lastPartial = SuffixSumBreak(sl.lengths[sl.firstVisible:sl.lastVisible], sl.lengths[sl.firstVisible])
		sl.lastVisible -= update
	}
}

// Move 1 item down
func (sl *ScrollingList) Next() {
	if sl.focused >= len(sl.items)-1 {
		return
	}
	sl.focused++
	sl.status = "Going Down"
	if sl.focused > sl.lastVisible-sl.NumLinesFromBorder && sl.lastVisible < len(sl.items)-1 {
		sl.lastVisible++
		var update int
		update, sl.firstPartial = PrefixSumBreak(sl.lengths[sl.firstVisible:sl.lastVisible], sl.lengths[sl.lastVisible])
		sl.firstVisible += 1 + update
	}
}

// Go to last item
func (sl *ScrollingList) GotoBottom() {
	sl.lastVisible = len(sl.items) - 1
	sl.lastPartial = 0
	sl.status = "At Bottom"
	var update int
	update, sl.firstPartial = SuffixSumBreak(sl.lengths, sl.Height)
	sl.firstVisible = len(sl.lengths) - update
	sl.focused = sl.lastVisible - sl.NumLinesFromBorder
}

// Go to first item
func (sl *ScrollingList) GotoTop() {
	sl.firstVisible = 0
	sl.firstPartial = 0
	sl.status = "At Top"
	sl.lastVisible, sl.lastPartial = PrefixSumBreak(sl.lengths[:len(sl.lengths)-1], sl.scrollHeight)
	sl.focused = sl.firstVisible + sl.NumLinesFromBorder
}

// Set the width and height for list and set initialized
func (sl *ScrollingList) SetSize(width, height int) {
	sl.initialized = true
	sl.Height = height
	sl.scrollHeight = height
	if sl.ShowFooter {
		sl.scrollHeight -= lipgloss.Height(sl.FooterView())
	}
	if sl.ShowHelp {
		sl.scrollHeight -= lipgloss.Height(sl.HelpView())
	}
	sl.Width = width
	sl.help.Width = width
	if sl.lastVisible < len(sl.lengths) - 1 - sl.NumLinesFromBorder {
		var update int
		update, sl.lastPartial = PrefixSumBreak(sl.lengths[sl.firstVisible:], sl.scrollHeight)
		sl.lastVisible = sl.firstVisible + update
	} else {
		var update int
		update, sl.firstPartial = SuffixSumBreak(sl.lengths[:sl.lastVisible], sl.scrollHeight)
		sl.firstVisible = sl.lastVisible - update
	}
}

// Returns the currently focused index
func (sl *ScrollingList) GetFocused() int {
	return sl.focused
}

// Returns the currently visible first and last indices
func (sl *ScrollingList) GetCurrentVisibleRange() (int, int) {
	return sl.firstVisible, sl.lastVisible
}

// TODO: Add title similar to default list
// TODO: help
