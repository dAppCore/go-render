// Package help re-exports charmbracelet/bubbles/help through go-html's tui
// seam — a key-binding help line/box. Swap the import path (bubbles/help →
// html/tui/help) and keep every help.Model / help.New reference unchanged.
package help

import "github.com/charmbracelet/bubbles/help"

type (
	Model  = help.Model
	KeyMap = help.KeyMap
)

var New = help.New
