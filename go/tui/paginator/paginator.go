// Package paginator re-exports charmbracelet/bubbles/paginator through
// go-html's tui seam — a page indicator (dots or arabic) with
// keystroke-driven paging navigation. Swap the import path
// (bubbles/paginator → html/tui/paginator) and keep every paginator.Model /
// paginator.New reference unchanged. Update, View, NextPage, PrevPage,
// OnFirstPage, OnLastPage, GetSliceBounds, SetTotalPages and ItemsOnPage are
// Model methods, not package functions, so they need no re-export here —
// they come along for free since Model is a genuine alias.
package paginator

import "charm.land/bubbles/v2/paginator"

type (
	Model  = paginator.Model
	Type   = paginator.Type
	KeyMap = paginator.KeyMap
	Option = paginator.Option
)

var (
	New            = paginator.New
	DefaultKeyMap  = paginator.DefaultKeyMap
	WithTotalPages = paginator.WithTotalPages
	WithPerPage    = paginator.WithPerPage
)

// Type identities — the pagination rendering styles.
const (
	Arabic = paginator.Arabic
	Dots   = paginator.Dots
)
