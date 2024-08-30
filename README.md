# listExtensions: bubbletea list bubbles extensions

[![Go Reference](https://pkg.go.dev/badge/github.com/plutonium-239/listExtensions.svg)](https://pkg.go.dev/github.com/plutonium-239/listExtensions)

Provides 2 bubbles:
1. `basiclist`: Basic list with no pagination, meant for use with viewport (see examples)
2. `scrollinglist`: Scrollable list (no pagination), (see examples)

Demo (`scrollinglist`):
[![asciicast](https://asciinema.org/a/663534.svg)](https://asciinema.org/a/663534)

Usage:
```bash
go get github.com/plutonium-239/listExtensions
```

```go
import (
	listx "github.com/plutonium-239/listExtensions"
)

newlist := listx.NewScrollingList()
```

Please refer to the short [example script](/scrollingList/examples/main.go) on more details.

> [!NOTE]
> The following is an AI generated summary of the customizability options generated from my code comments


# Scrolling List Component

ScrollingList is a bubbletea component that displays a list of items, with a focused item and scrolling capabilities.
The component is highly customizable, and allows for different styles, alignments, and toggling of footer and title.

## Initialization

To initialize a ScrollingList, simply call NewScrollingList().

## Customization

The ScrollingList struct has many customization options, some of which are:
- UnfocusedStyle: The style for unfocused items
- FocusedStyle: The style for the focused item
- FocusedLineStyle: The style for the line of the focused item
- FooterStyle: The style for the footer
- TitleStyle: The style for the title
- ListAlignment: The alignment of the list
- GlobalAlignment: The alignment of the list and footer w.r.t container
- MainVerticalAlignment: The vertical alignment of main text when not enough lines to fill screen
- ShowFooter: Whether to show the footer
- CustomFooter: Custom text for the footer
- ShowTitle: Whether to show the title
- ShowHelp: Whether to show the help text
- NumLinesFromBorder: The number of lines to keep between focused item and screen edges
- Width: The width of the list
- Height: The height of the list

## Usage

The ScrollingList component can be used like any other bubbletea component. It has Init(), View(), Update(), and SetSize() methods.
It supports keybindings for scrolling up, down, page up, page down, goto top, goto bottom, and toggling the footer, title, and help text.

The component can be used to display a list of items, with a focused item and scrolling capabilities. It can be customized to fit any style or design.
