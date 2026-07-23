// Package table re-exports charmbracelet/bubbles/table through go-html's tui
// seam — a selectable, scrollable data table. Swap the import path
// (bubbles/table → html/tui/table) and keep every table.Model / table.New /
// table.Column / table.Row reference unchanged. Methods such as SelectedRow,
// SetRows, MoveUp/MoveDown, Focus/Blur and GotoTop/GotoBottom are Model
// methods, not package functions, so they need no re-export here — they come
// along for free since Model is a genuine alias.
package table

import "charm.land/bubbles/v2/table"

type (
	Model  = table.Model
	Column = table.Column
	Row    = table.Row
	Styles = table.Styles
	KeyMap = table.KeyMap
	Option = table.Option
)

var (
	New           = table.New
	DefaultStyles = table.DefaultStyles
	DefaultKeyMap = table.DefaultKeyMap
	WithColumns   = table.WithColumns
	WithRows      = table.WithRows
	WithHeight    = table.WithHeight
	WithWidth     = table.WithWidth
	WithFocused   = table.WithFocused
	WithStyles    = table.WithStyles
	WithKeyMap    = table.WithKeyMap
)
