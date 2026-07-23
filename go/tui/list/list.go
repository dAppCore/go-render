// Package list re-exports charmbracelet/bubbles/list through go-html's tui
// seam. A consumer swaps the import path (bubbles/list → html/tui/list) and
// keeps every list.Model / list.New / list.Item reference unchanged.
package list

import "github.com/charmbracelet/bubbles/list"

type (
	Model           = list.Model
	Item            = list.Item
	ItemDelegate    = list.ItemDelegate
	DefaultDelegate = list.DefaultDelegate
	DefaultItem     = list.DefaultItem
	FilterState     = list.FilterState
	FilterFunc      = list.FilterFunc
	Rank            = list.Rank
	Styles          = list.Styles
)

var (
	New                = list.New
	NewDefaultDelegate = list.NewDefaultDelegate
	DefaultFilter      = list.DefaultFilter
	DefaultStyles      = list.DefaultStyles
)

// FilterState identities.
const (
	Unfiltered    = list.Unfiltered
	Filtering     = list.Filtering
	FilterApplied = list.FilterApplied
)
