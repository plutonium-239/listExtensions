package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	listx "github.com/plutonium-239/listExtensions/scrollingList"
)

type item struct {
	str string
}

func (it item) String() string {
	return it.str
}

func main() {
	list := listx.NewScrollingList()
	items := []fmt.Stringer{
		// item{"\"Top 10 staTues tHat CriEd bloOd\""},
		// item{""},
		item{"I guess that some of us\nAre just born with tragedy\nIn our blood\nIt's just basic chemistry"},
		item{""},
		item{"OR maybe all these open wounds\nIs how the light gets into you\nCos I'm starting to realise"},
		item{""},
		item{"No one's gonna come and\nRescue me\nI'm drowning in my sleep\nThe scars have got too deep"},
		item{"And no amount of love could\nSet you free\nThe fight's inside\nI'll take myself to hell and back"},
		item{"Tonight we go to war"},
		item{""},
		item{"Someone put a gun\nStraight to our hearts\nAnd paint the walls with our love\nThere's a world of hurt in us"},
		item{""},
		item{"And maybe once we spill our guts\nWe can stitch ourselves back up\nYeah I'm starting to realise"},
		item{""},
		item{"No one's gonna come and\nRescue me\nI'm drowning in my sleep\nThe scars have got too deep"},
		item{"And no amount of love could\nSet you free\nThe fight's inside\nI'll take myself to hell and back"},
		item{"Tonight we go to war"},
		item{""},
		item{"Cos I know\nThat you're low"},
		item{"But once you hit the bottom\nAt least there's nowhere to go but up"},
		item{"Yeah I know\nThat you're low"},
		item{"But once you hit the bottom\nAT least there's nowhere to go but up"},
		item{""},
		item{"Yeah"},
		item{"The hardest thing you'll ever know"},
		item{"Is there's no love like your own"},
		item{"No, there's no love like your own"},
		item{""},
		item{"Cos no one's gonna come and\nRescue me\nI'm drowning in my sleep\nThe scars have got too deep"},
		item{"And no amount of love could\nSet you free\nThe fight's inside"},
		item{"I'll take myself to hell and back"},
		item{""},
		item{"No one's gonna come and\nRescue me\nI'm drowning in my sleep\nThe scars have got too deep"},
		item{"And no amount of love could\nSet you free\nThe fight's inside\nI'll take myself to hell and back"},
		item{"Tonight we go to war"},
	}
	list.SetItems(items)
	list.ListAlignment = lipgloss.Center
	list.Title = "\"Top 10 staTues tHat CriEd bloOd\"\nBring Me The Horizon"
	m := model{list: list}
	// m.list.CustomFooter = MakeCustomFooter
	// m.list.MainVerticalAlignment = lipgloss.Center
	err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Start()
	if err != nil {
		panic(fmt.Sprintf("error in making progam: %v", err))
	}

}

func MakeCustomFooter(list *listx.ScrollingList) string {
	x := fmt.Sprintf("Focused: %d (ID=%d) | Status: %s", list.GetFocusedLine(), list.GetFocused(), list.GetStatus())
	return x
}

type model struct {
	list  listx.ScrollingList
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return m.list.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "x" {
			m.list.SetItemAt(item{"changed"}, m.list.GetFocused())
			return m, nil
		} else if msg.String() == "delete" {
			m.list.DeleteItemAt(m.list.GetFocused())
			return m, nil
		}
	case tea.WindowSizeMsg:
		t, cmd := m.list.Update(tea.WindowSizeMsg{Width: msg.Width/4, Height: msg.Height})
		m.list = t.(listx.ScrollingList)
		return m, cmd
	}
	t, cmd := m.list.Update(msg)
	m.list = t.(listx.ScrollingList)
	return m, cmd
}