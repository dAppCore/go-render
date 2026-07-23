// Package spinner re-exports charmbracelet/bubbles/spinner through go-html's
// tui seam — an animated activity indicator. Swap the import path
// (bubbles/spinner → html/tui/spinner) and keep every spinner.Model /
// spinner.New / spinner.MiniDot reference unchanged.
package spinner

import "github.com/charmbracelet/bubbles/spinner"

type (
	Model   = spinner.Model
	Spinner = spinner.Spinner
	TickMsg = spinner.TickMsg
)

var (
	New     = spinner.New
	Tick    = spinner.Tick
	MiniDot = spinner.MiniDot
)
