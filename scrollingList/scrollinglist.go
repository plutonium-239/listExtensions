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
	"github.com/charmbracelet/x/ansi"
)

type ScrollingList struct {
	originalItems []fmt.Stringer
	lengths       []int
	focusedID     int

	// len(preRendered) NOT EQUAL TO len(lengths) = len(originalItems)
	// It is supposed to serve as a "flattened" version of the texts in original items
	preRendered               []string
	itemIDs                   []int
	focused                   int
	firstVisible, lastVisible int

	// Styles can be customized
	UnfocusedStyle   lipgloss.Style
	FocusedStyle     lipgloss.Style
	FocusedLineStyle lipgloss.Style
	FooterStyle      lipgloss.Style
	TitleStyle       lipgloss.Style

	// Alignment of text within the list
	ListAlignment lipgloss.Position
	// Alignment of the list and footer w.r.t container
	GlobalAlignment lipgloss.Position

	// Footer can be toggled and customized, can be multiline
	ShowFooter   bool
	CustomFooter func(*ScrollingList) string
	// Title
	Title string
	// Title can be toggled
	ShowTitle bool

	Help help.Model
	// Help can be toggled
	ShowHelp bool

	// The number of lines to keep between focused item and screen edges
	NumLinesFromBorder int

	// Size of the whole list (scroll height is handled automatically based on footer/help)
	Width, Height int
	scrollHeight  int

	initialized      bool
	initializedItems bool
	// KeyMap
	KeyMap KeyMap

	status string
}

func NewScrollingList() ScrollingList {
	InitStyles()
	return ScrollingList{
		KeyMap:             DefaultKeyMap(),
		UnfocusedStyle:     DefaultUnfocusedStyle,
		FocusedStyle:       DefaultFocusedStyle,
		FocusedLineStyle:   DefaultFocusedLineStyle,
		FooterStyle:        DefaultFooterStyle,
		TitleStyle:         DefaultTitleStyle,
		ListAlignment:      lipgloss.Center,
		GlobalAlignment:    lipgloss.Center,
		NumLinesFromBorder: 5,
		ShowFooter:         true,
		ShowTitle:          true,
		ShowHelp:           true,
		Help:               help.New(),
		status:             "uninitialized",
	}
}

// BubbleTea Init method, does nothing
func (sl ScrollingList) Init() tea.Cmd {
	return nil
}

func (sl *ScrollingList) place(view string, truncate bool) string {
	if truncate && len(strings.Split(view, "\n")) == 1 {
		// TODO: TRUNCATE FOR MULTILINE
		view = ansi.Truncate(view, sl.Width, "…")
	}
	return lipgloss.PlaceHorizontal(sl.Width, sl.GlobalAlignment, view)
}

// BubbleTea View method, returns the actual final rendering result
func (sl ScrollingList) View() string {

	if !sl.initialized {
		return "..."
	}
	views := make([]string, 0)

	if sl.ShowTitle {
		views = append(views, sl.TitleView())
	}

	views = append(views, lipgloss.JoinVertical(sl.ListAlignment, sl.VisibleLines()...))
	if sl.ShowFooter {
		// return lipgloss.PlaceHorizontal(sl.Width, sl.GlobalAlignment, lipgloss.JoinVertical(sl.GlobalAlignment, views...))
		views = append(views, sl.FooterView())
	}
	if sl.ShowHelp {
		views = append(views, sl.place(sl.HelpView(), false))
	}
	if !sl.ShowFooter && !sl.ShowHelp {
		views[0] = sl.place(views[0], false)
	}
	return lipgloss.JoinVertical(sl.GlobalAlignment, views...)
}

// returns the visible lines, with partial first/last elements s.t. output fits on screen
func (sl *ScrollingList) VisibleLines() []string {
	if !sl.initializedItems {
		return []string{}
	}
	potentialRendered := sl.preRendered[sl.firstVisible : sl.lastVisible+1]

	// potentialRendered[sl.Focused] = sl.FocusedStyle.Render(potentialRendered[sl.Focused])
	// potentialLength := 0
	styledRendered := make([]string, len(potentialRendered))
	for i := sl.firstVisible; i <= sl.lastVisible; i++ {
		// potentialLength += sl.lengths[i]
		styledRendered[i-sl.firstVisible] = sl.styleSingle(i)
	}
	// result := lipgloss.JoinVertical(lipgloss.Left, styledRendered...)
	result := styledRendered

	return result
}

func (sl *ScrollingList) styleSingle(index int) string {
	if index == sl.focused {
		return sl.FocusedLineStyle.Render(sl.preRendered[index])
	}
	if sl.itemIDs[index] == sl.focusedID {
		return sl.FocusedStyle.Render(sl.preRendered[index])
	}
	return sl.UnfocusedStyle.Render(sl.preRendered[index])
}

func (sl *ScrollingList) TitleView() string {
	return sl.TitleStyle.Render(sl.place(sl.Title, true))
}

// returns the rendered footer
func (sl *ScrollingList) FooterView() string {
	// ShowFooter handling is done in View
	var text string
	if sl.CustomFooter != nil {
		text = sl.CustomFooter(sl)
	} else {
		text = fmt.Sprintf("Focused:%d (ID=%d), First:%d, Last:%d | Status: %s | main(w,h): (%d,%d) scrollHeight: %d",
			sl.focused, sl.focusedID, sl.firstVisible, sl.lastVisible, sl.status, sl.Width, sl.Height, sl.scrollHeight)
	}
	return sl.FooterStyle.Render(sl.place(text, true))
}

// returns the rendered help
func (sl *ScrollingList) HelpView() string {
	return sl.Help.View(sl.KeyMap)
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
		case key.Matches(msg, sl.KeyMap.PageDown):
			sl.PageDown()
			return sl, nil
		case key.Matches(msg, sl.KeyMap.PageUp):
			sl.PageUp()
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
		case key.Matches(msg, sl.KeyMap.ToggleTitle):
			sl.ShowTitle = !sl.ShowTitle
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		case key.Matches(msg, sl.KeyMap.ShowHideHelp):
			sl.ShowHelp = !sl.ShowHelp
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		case key.Matches(msg, sl.KeyMap.ShowFullHelp, sl.KeyMap.CloseFullHelp):
			sl.Help.ShowAll = !sl.Help.ShowAll
			sl.SetSize(sl.Width, sl.Height)
			return sl, nil
		}
	case tea.WindowSizeMsg:
		sl.SetSize(msg.Width, msg.Height)
		// TODO: find a way to persist through screen resize, this starts by deciding whether to keep focused closer to top or bottom (or middle?)
		// sl.lastVisible, _ = PrefixSumBreak(sl.lengths[:len(sl.lengths)-1], sl.scrollHeight)
		sl.initOrReinit()
		return sl, nil
	}
	return sl, nil
}

// Replace all items
func (sl *ScrollingList) SetItems(items []fmt.Stringer) {
	sl.originalItems = items
	sl.preRendered = make([]string, 0)
	sl.itemIDs = make([]int, 0)
	sl.lengths = make([]int, len(items))
	for i := range items {
		fulltextSplit := strings.Split(sl.originalItems[i].String(), "\n")
		for _, line := range fulltextSplit {
			sl.preRendered = append(sl.preRendered, line)
			sl.itemIDs = append(sl.itemIDs, i)
		}
		// sl.preRendered =
		sl.lengths[i] = len(fulltextSplit)
	}
	sl.initializedItems = true
	if sl.initialized {
		// sl.lastVisible, _ = PrefixSumBreak(sl.lengths[:len(sl.lengths)-1], sl.scrollHeight)
		// sl.lastVisible = sl.scrollHeight - 1
		sl.initOrReinit()
	}
}

func (sl *ScrollingList) initOrReinit() {
	if !sl.initializedItems {
		return
	}
	sl.setLast(sl.firstVisible + (sl.scrollHeight - 1))
	if diff := sl.lastVisible - sl.firstVisible - sl.scrollHeight + 1; diff < 0 {
		// i.e. sl.setlast chose len(lines) as minimum
		sl.setFirst(sl.firstVisible + diff)
	}
	sl.ensureFocusedInBounds()
}

func (sl *ScrollingList) ensureFocusedInBounds() {
	sl.focused = max(sl.firstVisible, sl.focused)
	sl.focused = min(sl.lastVisible, sl.focused)
}

func (sl *ScrollingList) setFocusedID() {
	if !sl.initializedItems {
		return
	}
	sl.focusedID = sl.itemIDs[sl.focused]
}

// Replace one item at index int with given item
func (sl *ScrollingList) SetItemAt(item fmt.Stringer, index int) {
	if !sl.initializedItems {
		return
	}
	sl.status = fmt.Sprintf("ID%d updated", index)
	oldLen := sl.lengths[index]
	sl.originalItems[index] = item
	sl.SetItems(sl.originalItems)
	diff := sl.lengths[index] - oldLen

	// TODO: Implement a == case properly (so that it stays same distance away from itemID start)
	if index <= sl.focusedID {
		sl.focused += diff
	} else {
		// Is there enough (i.e. -diff)no. of lines after lastvisible to fill up screen?
		if sl.lastVisible >= len(sl.itemIDs)-1-(-diff) {
			// Ans: No, need to set lastvisible to very last line
			diffFromEnd := sl.lastVisible - sl.focused
			sl.lastVisible = len(sl.itemIDs) - 1
			sl.focused = sl.lastVisible - diffFromEnd
		}
	}

	sl.ensureFocusedInBounds()
}

// Delete one item at index int
func (sl *ScrollingList) DeleteItemAt(index int) {
	if !sl.initializedItems {
		return
	}
	sl.status = fmt.Sprintf("ID%d deleted", index)
	oldLen := sl.lengths[index]
	sl.originalItems = slices.Concat(sl.originalItems[:index], sl.originalItems[index+1:])
	sl.SetItems(sl.originalItems)

	// Go to last line of item with ID = index - 1
	if index == sl.focusedID {
		for i := sl.focused; i >= sl.firstVisible; i-- {
			if sl.itemIDs[i] == sl.focusedID-1 {
				sl.focused = i
				sl.focusedID -= 1
				break
			}
		}
	} else if index < sl.focusedID {
		sl.focused -= oldLen
		sl.ensureFocusedInBounds()
	}
}

// Move 1 item up
func (sl *ScrollingList) Prev() {
	if sl.focused <= 0 {
		return
	}
	sl.focused--
	sl.setFocusedID()
	sl.status = "Going Up"
	if sl.focused < sl.firstVisible+sl.NumLinesFromBorder && sl.firstVisible > 0 {
		sl.firstVisible--
		sl.lastVisible--
	}
}

// Move 1 item down
func (sl *ScrollingList) Next() {
	if sl.focused >= len(sl.itemIDs)-1 {
		return
	}
	sl.focused++
	sl.setFocusedID()
	sl.status = "Going Down"
	if sl.focused > sl.lastVisible-sl.NumLinesFromBorder && sl.lastVisible < len(sl.itemIDs)-1 {
		sl.lastVisible++
		sl.firstVisible++
	}
}

// Move 1 page up
func (sl *ScrollingList) PageUp() {
	h := sl.scrollHeight - 1
	sl.setFirst(sl.firstVisible - h)
	sl.focused -= h
	sl.focused = max(sl.focused, sl.firstVisible)
	sl.setLast(sl.firstVisible + h)
	sl.setFocusedID()
	sl.status = "Page Up"
}

// Move 1 page down
func (sl *ScrollingList) PageDown() {
	h := sl.scrollHeight - 1
	sl.setLast(sl.lastVisible + h)
	sl.focused += h
	sl.focused = min(sl.focused, sl.lastVisible)
	sl.setFirst(sl.lastVisible - h)
	sl.setFocusedID()
	sl.status = "Page Down"
}

// Go to first item
func (sl *ScrollingList) GotoTop() {
	sl.firstVisible = 0
	sl.status = "At Top"
	sl.setLast(sl.firstVisible + (sl.scrollHeight - 1))
	sl.focused = sl.firstVisible
	sl.setFocusedID()
}

// Go to last item
func (sl *ScrollingList) GotoBottom() {
	sl.lastVisible = len(sl.itemIDs) - 1
	sl.status = "At Bottom"
	sl.setFirst(sl.lastVisible - (sl.scrollHeight - 1))
	sl.focused = sl.lastVisible
	sl.setFocusedID()
}

// Set the width and height for list and set initialized
func (sl *ScrollingList) SetSize(width, height int) {
	sl.initialized = true
	sl.Height = height
	sl.scrollHeight = height
	if sl.ShowTitle {
		sl.scrollHeight -= lipgloss.Height(sl.TitleView())
	}
	if sl.ShowFooter {
		sl.scrollHeight -= lipgloss.Height(sl.FooterView())
	}
	if sl.ShowHelp {
		sl.scrollHeight -= lipgloss.Height(sl.HelpView())
	}
	sl.Width = width
	sl.Help.Width = width
	if sl.lastVisible < len(sl.itemIDs)-1-sl.NumLinesFromBorder {
		sl.setLast(sl.firstVisible + (sl.scrollHeight - 1))
	} else {
		sl.setFirst(sl.lastVisible - (sl.scrollHeight - 1))
	}
}

// Returns the currently focused line
func (sl *ScrollingList) GetFocusedLine() int {
	return sl.focused
}

// Returns the currently focused item id
func (sl *ScrollingList) GetFocused() int {
	return sl.focusedID
}

// Returns the currently visible first and last indices
func (sl *ScrollingList) GetCurrentVisibleRange() (int, int) {
	return sl.firstVisible, sl.lastVisible
}

// Returns the current status
func (sl *ScrollingList) GetStatus() string {
	return sl.status
}

// Convenience funcitions
func (sl *ScrollingList) setFirst(value int) {
	sl.firstVisible = max(value, 0)
}
func (sl *ScrollingList) setLast(value int) {
	sl.lastVisible = min(value, len(sl.itemIDs)-1)
}
