package scrollinglist

import "github.com/charmbracelet/lipgloss"


var DefaultUnfocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.HiddenBorder(), false, false, false, true).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#C2B8C2"))

var DefaultFocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(lipgloss.Color("197")).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#EE6FF8"))

var DefaultFocusedLineStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder(), false, false, false, true).
	BorderForeground(lipgloss.Color("197")).
	PaddingLeft(1).
	Foreground(lipgloss.Color("#EE6FF8"))

var	RightArrowBorder = lipgloss.ThickBorder()

func InitStyles() {
	RightArrowBorder.Left = ">"
	DefaultFocusedLineStyle = DefaultFocusedLineStyle.Border(RightArrowBorder, false, false, false, true)
	// DefaultFocusedStyle.Border()
}