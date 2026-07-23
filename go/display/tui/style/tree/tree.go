// Package tree is go-html's static styled tree renderer, over
// lipgloss/v2/tree — build a Tree of nested Root/Child nodes, style the
// enumerator/indenter/item per node, then paint it once with String(). It is
// a STATIC renderer: unlike the interactive tui/table and tui/list (bubbles
// Models driven by an Update/View loop and a selection cursor), a Tree has no
// loop and no cursor — it renders and is done. Swap the import path
// (lipgloss/v2/tree → html/tui/style/tree) and keep every tree.New /
// tree.Root / tree.Tree reference unchanged. Child, Enumerator, Indenter,
// *Style, *StyleFunc, Width, Offset, Hide and the rest are Tree methods, not
// package functions, so they need no re-export here — they come along for
// free since Tree is a genuine alias.
package tree

import "charm.land/lipgloss/v2/tree"

type (
	Node         = tree.Node
	Leaf         = tree.Leaf
	Tree         = tree.Tree
	Children     = tree.Children
	NodeChildren = tree.NodeChildren
	Filter       = tree.Filter
	Enumerator   = tree.Enumerator
	Indenter     = tree.Indenter
	StyleFunc    = tree.StyleFunc
	Style        = tree.Style
)

var (
	New               = tree.New
	Root              = tree.Root
	NewLeaf           = tree.NewLeaf
	NewStringData     = tree.NewStringData
	NewFilter         = tree.NewFilter
	DefaultEnumerator = tree.DefaultEnumerator
	RoundedEnumerator = tree.RoundedEnumerator
	DefaultIndenter   = tree.DefaultIndenter
)
