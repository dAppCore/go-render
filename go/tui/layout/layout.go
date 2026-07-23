// SPDX-Licence-Identifier: EUPL-1.2

// Package layout re-exports charmbracelet/ultraviolet/layout through
// go-html's tui seam — a Cassowary constraint solver that partitions a
// rectangular area into non-overlapping sub-rectangles. Swap the import path
// (ultraviolet/layout → html/tui/layout) and keep every layout.Vertical /
// layout.Horizontal / layout.Split reference unchanged.
//
// A Layout takes a Direction (Vertical or Horizontal flow) and an ordered
// list of Constraints — Len (fixed cells), Min/Max (bounded), Percent/Ratio
// (proportional to the full area), or Fill (greedy, shares leftover space by
// weight) — and Split partitions a Rect into a Splitted slice of sub-Rects,
// one per constraint, in the same order. Padding insets the area before
// solving; Flex controls how any leftover space is distributed once every
// constraint is satisfied (FlexStart, the zero value, is the default unless
// WithFlex overrides it).
//
// This is the one primitive nothing else in go-html's tui/ stack has —
// lipgloss and bubbletea style and animate cells, neither partitions space by
// constraint — which makes this package the seed of the emerging
// render/"active" layer (go-render), not another parity wrap.
//
// WithDirection, WithPadding, WithFlex, WithSpacing, WithConstraints,
// SplitWithSpacers, Split and String are Layout/constraint methods, not
// package functions, so they need no re-export here — they come along for
// free since Layout and the constraint kinds are genuine aliases. Splitted's
// Assign method rides free the same way, as do Rect's methods (Dx, Dy, …) —
// Rect is a genuine alias through uv.Rectangle down to image.Rectangle.
package layout

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/layout"
)

// Layout is the direction + constraint set a Split call resolves. Constraint
// is the sealed interface every constraint kind (Len, Min, Max, Percent,
// Ratio, Fill) implements. Direction picks Vertical or Horizontal flow. Flex
// controls how leftover space is distributed once constraints are satisfied.
// Rect is the rectangle type Split consumes and produces. Splitted is the
// []Rect result of a Split call, in constraint order. Padding is the inset a
// Layout applies to its area before solving.
type (
	Layout     = layout.Layout
	Constraint = layout.Constraint
	Direction  = layout.Direction
	Flex       = layout.Flex
	Rect       = uv.Rectangle
	Splitted   = layout.Splitted
	Padding    = layout.Padding
)

// Constraint kinds — each a defined type over int (Ratio is a Num/Den
// struct) implementing Constraint by conversion: Len(3), Min(1), Max(9),
// Percent(50), Fill(1), Ratio{Num: 1, Den: 2}. Priority order when the
// solver cannot satisfy every constraint (highest first): Min, Max, Len,
// Percent, Ratio, Fill.
type (
	Min     = layout.Min
	Max     = layout.Max
	Len     = layout.Len
	Percent = layout.Percent
	Ratio   = layout.Ratio
	Fill    = layout.Fill
)

// Direction identities, selecting how a Layout's segments flow.
const (
	DirectionVertical   = layout.DirectionVertical
	DirectionHorizontal = layout.DirectionHorizontal
)

// Flex identities, controlling how leftover space is distributed once every
// constraint is satisfied. FlexStart is the zero value and so the default
// for New/Vertical/Horizontal unless WithFlex overrides it.
const (
	FlexStart        = layout.FlexStart
	FlexLegacy       = layout.FlexLegacy
	FlexEnd          = layout.FlexEnd
	FlexCenter       = layout.FlexCenter
	FlexSpaceBetween = layout.FlexSpaceBetween
	FlexSpaceEvenly  = layout.FlexSpaceEvenly
	FlexSpaceAround  = layout.FlexSpaceAround
)

// New builds a Layout, Vertical/Horizontal are its directional shorthand, and
// Pad builds a Padding from CSS-style shorthand argument counts (0/1/2/4
// sides).
var (
	New        = layout.New
	Vertical   = layout.Vertical
	Horizontal = layout.Horizontal
	Pad        = layout.Pad
)
