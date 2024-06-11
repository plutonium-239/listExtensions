package main

import (
	"fmt"
	listx "x-MOD2-DO/tui/listextensions"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	str string
}

func (it item) String() string {
	return it.str
}

type model struct {
	ready    bool
	quitting bool
	list     listx.BasicList
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if !m.ready {
		return "..."
	}
	if m.quitting {
		return "Life is a grave, and I dig it"
	}
	return lipgloss.JoinVertical(lipgloss.Center, m.viewport.View(), m.FooterView())
}

func (m *model) FooterView() string {
	return fmt.Sprintf("Focused:%d, vp.YPos: %d, vp.YOffset: %d", m.list.Focused, m.viewport.YPosition, m.viewport.YOffset)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// l.Items = make([]fmt.Stringer, 0)
			m.quitting = true
			return m, tea.Quit
		default:
			var cmd2 tea.Cmd
			t, cmd := m.list.Update(msg)
			m.list = t.(listx.BasicList)
			m.viewport.SetContent(m.list.View())
			m.viewport, cmd2 = m.viewport.Update(msg)
			return m, tea.Batch(cmd, cmd2)
		}
	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(m.FooterView())
		m.ready = true
		m.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
		m.viewport.SetContent(m.list.View())
		return m, nil
	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
}

func main() {
	list := listx.NewBasicList()
	items := []fmt.Stringer{
		item{"Hello"},
		item{"This is a basic list that is not paginated"},
		item{"like the default list in bubbletea"},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"This is because it is very hard to get the default list working with"},
		item{"a viewport for scrolling behaviour."},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"I think scrolling is a more natural way to navigate than pagination."},
		item{"But that's just my preference"},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"Have fun!"},
		item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."}, item{"."},
		item{"♥💖"},
	}
	list.SetItems(items)
	m := model{list: list}
	err := tea.NewProgram(m).Start()
	if err != nil {
		panic(fmt.Sprintf("error in making progam: %v", err))
	}

}
