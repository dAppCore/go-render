// Package spinner re-exports charmbracelet/bubbles/spinner through go-html's
// tui seam — an animated activity indicator. Swap the import path
// (bubbles/spinner → html/tui/spinner) and keep every spinner.Model /
// spinner.New / spinner.MiniDot reference unchanged. Tick is a Model method
// (m.Tick()), not a package function, so it needs no re-export here — it
// comes along for free since Model is a genuine alias.
package spinner

import "charm.land/bubbles/v2/spinner"

type (
	Model   = spinner.Model
	Spinner = spinner.Spinner
	TickMsg = spinner.TickMsg
)

var (
	New         = spinner.New
	MiniDot     = spinner.MiniDot
	WithSpinner = spinner.WithSpinner
	Ellipsis    = spinner.Ellipsis
)
