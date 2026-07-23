// Package table is go-html's static styled table renderer, over
// lipgloss/v2/table — a fixed-layout Table with borders and a per-cell
// StyleFunc(row, col), built up and then painted once with String() (or
// Render()). It is the STATIC renderer: distinct from the interactive
// tui/table (a bubbles Model with a scrolling, key-driven selection cursor
// wired through an Update/View loop) — this Table has no loop and no cursor,
// it renders and is done. Swap the import path (lipgloss/v2/table →
// html/tui/style/table) and keep every table.New / table.Table /
// table.StyleFunc reference unchanged. Headers, Row, Rows, Border*, Width,
// Height, StyleFunc and the rest are Table methods, not package functions, so
// they need no re-export here — they come along for free since Table is a
// genuine alias.
package table

import "charm.land/lipgloss/v2/table"

type (
	Table      = table.Table
	StyleFunc  = table.StyleFunc
	Data       = table.Data
	StringData = table.StringData
	Filter     = table.Filter
)

// HeaderRow is the row index a StyleFunc receives when styling the header.
const HeaderRow = table.HeaderRow

var (
	New           = table.New
	DefaultStyles = table.DefaultStyles
	NewStringData = table.NewStringData
	NewFilter     = table.NewFilter
	DataToMatrix  = table.DataToMatrix
)
