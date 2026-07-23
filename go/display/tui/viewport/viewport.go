// Package viewport re-exports charmbracelet/bubbles/viewport through go-html's
// tui seam — a scrollable content pane. Swap the import path
// (bubbles/viewport → html/tui/viewport) and keep every viewport.Model /
// viewport.New reference unchanged.
package viewport

import "charm.land/bubbles/v2/viewport"

type (
	Model  = viewport.Model
	Option = viewport.Option
)

var (
	New        = viewport.New
	WithWidth  = viewport.WithWidth
	WithHeight = viewport.WithHeight
)
