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