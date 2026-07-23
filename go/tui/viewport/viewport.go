// Package viewport re-exports charmbracelet/bubbles/viewport through go-html's
// tui seam — a scrollable content pane. Swap the import path
// (bubbles/viewport → html/tui/viewport) and keep every viewport.Model /
// viewport.New reference unchanged.
package viewport

import "github.com/charmbracelet/bubbles/viewport"

type Model = viewport.Model

var New = viewport.New
