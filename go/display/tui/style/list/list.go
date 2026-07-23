// Package list is go-html's static styled list renderer, over
// lipgloss/v2/list — a (possibly nested) List of items with an Enumerator
// prefix (Bullet, Arabic, Roman, Alphabet, Dash, Asterisk, …), built up and
// then painted once with String(). It is the STATIC renderer: distinct from
// the interactive tui/list (a bubbles Model with a filterable, key-driven
// selection cursor wired through an Update/View loop) — this List has no
// loop and no cursor. A List item that is itself a *List renders as a nested
// sub-list — there is no separate sub-list constructor, nest one with
// Item/Items. Swap the import path (lipgloss/v2/list → html/tui/style/list)
// and keep every list.New / list.List / list.Bullet reference unchanged.
// Item, Items, Enumerator, Indenter, *Style, *StyleFunc, Offset, Hide and the
// rest are List methods, not package functions, so they need no re-export
// here — they come along for free since List is a genuine alias.
package list

import "charm.land/lipgloss/v2/list"

type (
	List       = list.List
	Items      = list.Items
	StyleFunc  = list.StyleFunc
	Enumerator = list.Enumerator
	Indenter   = list.Indenter
)

var (
	New      = list.New
	Alphabet = list.Alphabet
	Arabic   = list.Arabic
	Roman    = list.Roman
	Bullet   = list.Bullet
	Asterisk = list.Asterisk
	Dash     = list.Dash
)
