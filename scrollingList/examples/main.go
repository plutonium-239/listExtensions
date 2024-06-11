package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
		item{"\"Top 10 staTues tHat CriEd bloOd\""},
		item{""},
		item{"I guess that some of us"},
		item{"Are just born with tragedy"},
		item{"In our blood"},
		item{"It's just basic chemistry"},
		item{""},
		item{"OR maybe all these open wounds"},
		item{"Is how the light gets into you"},
		item{"Cos I'm starting to realise"},
		item{""},
		item{"No one's gonna come and"},
		item{"Rescue me"},
		item{"I'm drowning in my sleep"},
		item{"The scars have got too deep"},
		item{"And no amount of love could"},
		item{"Set you free"},
		item{"The fight's inside"},
		item{"I'll take myself to hell and back"},
		item{"Tonight we go to war"},
		item{""},
		item{"Someone put a gun"},
		item{"Straight to our hearts"},
		item{"And paint the walls with our love"},
		item{"There's a world of hurt in us"},
		item{""},
		item{"And maybe once we spill our guts"},
		item{"We can stitch ourselves back up"},
		item{"Yeah I'm starting to realise"},
		item{""},
		item{"No one's gonna come and"},
		item{"Rescue me"},
		item{"I'm drowning in my sleep"},
		item{"The scars have got too deep"},
		item{"And no amount of love could"},
		item{"Set you free"},
		item{"The fight's inside"},
		item{"I'll take myself to hell and back"},
		item{"Tonight we go to war"},
		item{""},
		item{"Cos I know"},
		item{"That you're low"},
		item{"But once you hit the bottom"},
		item{"At least there's nowhere to go but up"},
		item{"Yeah I know"},
		item{"That you're low"},
		item{"But once you hit the bottom"},
		item{"AT least there's nowhere to go but up"},
		item{""},
		item{"Yeah"},
		item{"The hardest thing you'll ever know"},
		item{"Is there's no love like your own"},
		item{"No, there's no love like your own"},
		item{""},
		item{"Cos no one's gonna come and"},
		item{"Rescue me"},
		item{"I'm drowning in my sleep"},
		item{"The scars have got too deep"},
		item{"And no amount of love could"},
		item{"Set you free"},
		item{"The fight's inside"},
		item{"I'll take myself to hell and back"},
		item{""},
		item{"No one's gonna come and"},
		item{"Rescue me"},
		item{"I'm drowning in my sleep"},
		item{"The scars have got too deep"},
		item{"And no amount of love could"},
		item{"Set you free"},
		item{"The fight's inside"},
		item{"I'll take myself to hell and back"},
		item{"Tonight we go to war"},
	}
	list.SetItems(items)
	err := tea.NewProgram(
		list,
		tea.WithAltScreen(),
	).Start()
	if err != nil {
		panic(fmt.Sprintf("error in making progam: %v", err))
	}

}
