// Package textinput re-exports charmbracelet/bubbles/textinput through
// go-html's tui seam — a single-line text field. Swap the import path
// (bubbles/textinput → html/tui/textinput) and keep every textinput.Model /
// textinput.New reference unchanged.
package textinput

import "charm.land/bubbles/v2/textinput"

type Model = textinput.Model

var (
	New   = textinput.New
	Blink = textinput.Blink
)
