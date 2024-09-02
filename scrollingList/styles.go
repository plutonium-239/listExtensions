package scrollinglist

import "github.com/charmbracelet/lipgloss"

var (
	DefaultUnfocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder(), false, false, false, true).
				PaddingLeft(1).
				Foreground(lipgloss.AdaptiveColor{Light: "#3d473d", Dark: "#C2B8C2"})

	DefaultFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder(), false, false, false, true).
				// BorderForeground(lipgloss.Color("#ff0077")).
				PaddingLeft(1).
				Foreground(lipgloss.Color("#EE6FF8"))

	DefaultFocusedLineStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#ff0077")).
				PaddingLeft(1).
				Foreground(lipgloss.Color("#EE6FF8"))

	RightArrowBorder = lipgloss.ThickBorder()

	DefaultFooterStyle = lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Light: "#e2e2e2", Dark: "#1D1D1D"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#984b9e", Dark: "#f7a3fb"})

	DefaultTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#e2e2e2", Dark: "#1D1D1D"}).
				Background(lipgloss.Color("#ff0077"))

	DefaultEditingStyle = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
				Background(lipgloss.AdaptiveColor{Light: "#e2e2e2", Dark: "#1D1D1D"})
	// Border(lipgloss.DoubleBorder(), false, true).
	// BorderForeground(lipgloss.AdaptiveColor{Light: "#e2e2e2", Dark: "#1D1D1D"}).
	// BorderBackground(lipgloss.Color("#ff0077"))
)

func InitStyles() {
	RightArrowBorder.Left = ">"
	DefaultFocusedLineStyle = DefaultFocusedLineStyle.Border(RightArrowBorder, false, false, false, true)
	// DefaultFocusedStyle.Border()
}
