// SPDX-Licence-Identifier: EUPL-1.2

package layout_test

import (
	"image"
	"testing"

	"dappco.re/go/render/display/tui/layout"
)

// TestNew builds a Layout for each Direction and checks the fields New sets
// directly — the Direction it was given and the Constraints slice length —
// with no solver involved.
func TestNew(t *testing.T) {
	h := layout.New(layout.DirectionHorizontal, layout.Len(1), layout.Fill(1))
	if h.Direction != layout.DirectionHorizontal {
		t.Fatalf("New(DirectionHorizontal, ...).Direction = %v, want DirectionHorizontal", h.Direction)
	}
	if len(h.Constraints) != 2 {
		t.Fatalf("len(New(...).Constraints) = %d, want 2", len(h.Constraints))
	}

	v := layout.New(layout.DirectionVertical, layout.Len(1))
	if v.Direction != layout.DirectionVertical {
		t.Fatalf("New(DirectionVertical, ...).Direction = %v, want DirectionVertical", v.Direction)
	}
	if len(v.Constraints) != 1 {
		t.Fatalf("len(New(...).Constraints) = %d, want 1", len(v.Constraints))
	}
}

// TestVertical proves the shorthand delegates to New(DirectionVertical, …)
// exactly — same Direction, same Constraints count.
func TestVertical(t *testing.T) {
	got := layout.Vertical(layout.Len(2), layout.Fill(1))
	want := layout.New(layout.DirectionVertical, layout.Len(2), layout.Fill(1))
	if got.Direction != want.Direction {
		t.Fatalf("Vertical(...).Direction = %v, want %v (New(DirectionVertical, ...))", got.Direction, want.Direction)
	}
	if len(got.Constraints) != len(want.Constraints) {
		t.Fatalf("len(Vertical(...).Constraints) = %d, want %d", len(got.Constraints), len(want.Constraints))
	}
}

// TestHorizontal proves the shorthand delegates to New(DirectionHorizontal,
// …) exactly — same Direction, same Constraints count.
func TestHorizontal(t *testing.T) {
	got := layout.Horizontal(layout.Len(2), layout.Fill(1))
	want := layout.New(layout.DirectionHorizontal, layout.Len(2), layout.Fill(1))
	if got.Direction != want.Direction {
		t.Fatalf("Horizontal(...).Direction = %v, want %v (New(DirectionHorizontal, ...))", got.Direction, want.Direction)
	}
	if len(got.Constraints) != len(want.Constraints) {
		t.Fatalf("len(Horizontal(...).Constraints) = %d, want %d", len(got.Constraints), len(want.Constraints))
	}
}

// TestDirection proves DirectionVertical and DirectionHorizontal steer which
// axis Split partitions: the same two constraints (Len(3), Fill(1)) stack
// top-to-bottom under Vertical and left-to-right under Horizontal.
func TestDirection(t *testing.T) {
	area := image.Rect(0, 0, 10, 10)
	cs := []layout.Constraint{layout.Len(3), layout.Fill(1)}

	v := layout.New(layout.DirectionVertical, cs...).Split(area)
	wantV := image.Rect(0, 0, 10, 3)
	if v[0] != wantV {
		t.Fatalf("DirectionVertical segment 0 = %v, want %v (stacks top-to-bottom)", v[0], wantV)
	}

	h := layout.New(layout.DirectionHorizontal, cs...).Split(area)
	wantH := image.Rect(0, 0, 3, 10)
	if h[0] != wantH {
		t.Fatalf("DirectionHorizontal segment 0 = %v, want %v (stacks left-to-right)", h[0], wantH)
	}
}

// TestLayout builds a Layout as a bare struct literal — proving the alias is
// a genuine usable struct, not only reachable through New/Vertical/
// Horizontal — and Splits it.
func TestLayout(t *testing.T) {
	l := layout.Layout{
		Direction:   layout.DirectionHorizontal,
		Constraints: []layout.Constraint{layout.Percent(50), layout.Percent(50)},
	}

	got := l.Split(image.Rect(0, 0, 10, 10))
	want := layout.Splitted{image.Rect(0, 0, 5, 10), image.Rect(5, 0, 10, 10)}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Split() = %v, want %v", got, want)
	}
}

// TestLen fixes a segment to exactly the given number of cells: Len(3)
// against a Fill(1) sibling over a 10-tall rect leaves the Len segment at
// exactly 3 and the Fill segment absorbing the remaining 7.
func TestLen(t *testing.T) {
	got := layout.Vertical(layout.Len(3), layout.Fill(1)).Split(image.Rect(0, 0, 10, 10))

	wantTop := image.Rect(0, 0, 10, 3)
	wantBottom := image.Rect(0, 3, 10, 10)
	if len(got) != 2 || got[0] != wantTop || got[1] != wantBottom {
		t.Fatalf("Split() = %v, want [%v %v]", got, wantTop, wantBottom)
	}
}

// TestFill proves weighted Fill segments share leftover space proportionally
// to their integer weight: Fill(1) and Fill(2) over a 9-wide rect split 1:2,
// giving exactly 3 and 6 cells with no rounding ambiguity.
func TestFill(t *testing.T) {
	got := layout.Horizontal(layout.Fill(1), layout.Fill(2)).Split(image.Rect(0, 0, 9, 5))

	wantLeft := image.Rect(0, 0, 3, 5)
	wantRight := image.Rect(3, 0, 9, 5)
	if len(got) != 2 || got[0] != wantLeft || got[1] != wantRight {
		t.Fatalf("Split() = %v, want [%v %v]", got, wantLeft, wantRight)
	}
}

// TestPercent sizes a segment as a fraction of the total area: two
// Percent(50) constraints over a 10-wide rect split evenly, 5 and 5.
func TestPercent(t *testing.T) {
	got := layout.Horizontal(layout.Percent(50), layout.Percent(50)).Split(image.Rect(0, 0, 10, 10))

	wantLeft := image.Rect(0, 0, 5, 10)
	wantRight := image.Rect(5, 0, 10, 10)
	if len(got) != 2 || got[0] != wantLeft || got[1] != wantRight {
		t.Fatalf("Split() = %v, want [%v %v]", got, wantLeft, wantRight)
	}
}

// TestRatio sizes a segment as a Num/Den fraction of the total area: four
// Ratio{1, 4} constraints over a 50-wide rect give 13/12/13/12 — the
// documented example from upstream's Ratio doc comment, since 1/4 of 50
// cannot be represented as a whole number of cells and the solver rounds
// each segment to the nearest cell independently.
func TestRatio(t *testing.T) {
	quarter := layout.Ratio{Num: 1, Den: 4}
	got := layout.Horizontal(quarter, quarter, quarter, quarter).Split(image.Rect(0, 0, 50, 10))

	want := layout.Splitted{
		image.Rect(0, 0, 13, 10),
		image.Rect(13, 0, 25, 10),
		image.Rect(25, 0, 38, 10),
		image.Rect(38, 0, 50, 10),
	}
	if len(got) != len(want) {
		t.Fatalf("len(Split()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Split()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

// TestMin proves Min enforces a genuine floor: paired against a Percent(100)
// sibling that alone would claim the whole 50-wide area, Min(20) still
// clamps its own segment to at least 20, forcing Percent(100) down to 30
// under FlexLegacy (the mode Min's own doc example assumes).
func TestMin(t *testing.T) {
	got := layout.Horizontal(layout.Percent(100), layout.Min(20)).
		WithFlex(layout.FlexLegacy).
		Split(image.Rect(0, 0, 50, 10))

	wantLeft := image.Rect(0, 0, 30, 10)
	wantRight := image.Rect(30, 0, 50, 10)
	if len(got) != 2 || got[0] != wantLeft || got[1] != wantRight {
		t.Fatalf("Split() = %v, want [%v %v] (Min(20) floors the right segment)", got, wantLeft, wantRight)
	}
}

// TestMax proves Max enforces a genuine ceiling: paired against a
// Percent(0) sibling, Max(20) still caps its own segment at 20 even though
// nothing else claims the leftover, forcing Percent(0) to absorb the other
// 30 under FlexLegacy (the mode Max's own doc example assumes).
func TestMax(t *testing.T) {
	got := layout.Horizontal(layout.Percent(0), layout.Max(20)).
		WithFlex(layout.FlexLegacy).
		Split(image.Rect(0, 0, 50, 10))

	wantLeft := image.Rect(0, 0, 30, 10)
	wantRight := image.Rect(30, 0, 50, 10)
	if len(got) != 2 || got[0] != wantLeft || got[1] != wantRight {
		t.Fatalf("Split() = %v, want [%v %v] (Max(20) caps the right segment)", got, wantLeft, wantRight)
	}
}

// TestFlex proves WithFlex actually changes how leftover space is
// distributed: two Len(3) segments over a 10-tall rect leave 4 cells
// unclaimed under the default FlexStart (segments stay exactly 3 and 3,
// gap trails after), but FlexLegacy assigns that surplus to the final
// segment (3 and 7, covering the whole area).
func TestFlex(t *testing.T) {
	area := image.Rect(0, 0, 10, 10)
	cs := []layout.Constraint{layout.Len(3), layout.Len(3)}

	start := layout.Vertical(cs...).Split(area)
	wantStart := layout.Splitted{image.Rect(0, 0, 10, 3), image.Rect(0, 3, 10, 6)}
	if len(start) != 2 || start[0] != wantStart[0] || start[1] != wantStart[1] {
		t.Fatalf("default Flex Split() = %v, want %v (FlexStart leaves surplus trailing)", start, wantStart)
	}

	legacy := layout.Vertical(cs...).WithFlex(layout.FlexLegacy).Split(area)
	wantLegacy := layout.Splitted{image.Rect(0, 0, 10, 3), image.Rect(0, 3, 10, 10)}
	if len(legacy) != 2 || legacy[0] != wantLegacy[0] || legacy[1] != wantLegacy[1] {
		t.Fatalf("WithFlex(FlexLegacy) Split() = %v, want %v (surplus assigned to the last segment)", legacy, wantLegacy)
	}
}

// TestPadding insets the solved area from a Padding set directly on the
// Layout struct (rather than through the Pad shorthand): a 1-cell inset on
// every side shrinks a 10x10 rect to the 8x8 area from (1,1) to (9,9).
func TestPadding(t *testing.T) {
	l := layout.Layout{
		Direction:   layout.DirectionVertical,
		Constraints: []layout.Constraint{layout.Fill(1)},
		Padding:     layout.Padding{Top: 1, Right: 1, Bottom: 1, Left: 1},
	}

	got := l.Split(image.Rect(0, 0, 10, 10))
	want := image.Rect(1, 1, 9, 9)
	if len(got) != 1 || got[0] != want {
		t.Fatalf("Split() = %v, want [%v] (Padding insets the area before solving)", got, want)
	}
}

// TestPad covers every shorthand argument count Pad accepts, following the
// CSS-style convention documented on Pad: zero, uniform, vertical/
// horizontal, and top/right/bottom/left.
func TestPad(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want layout.Padding
	}{
		{"zero", nil, layout.Padding{}},
		{"uniform", []int{1}, layout.Padding{Top: 1, Right: 1, Bottom: 1, Left: 1}},
		{"vertical-horizontal", []int{1, 2}, layout.Padding{Top: 1, Right: 2, Bottom: 1, Left: 2}},
		{"top-right-bottom-left", []int{1, 2, 3, 4}, layout.Padding{Top: 1, Right: 2, Bottom: 3, Left: 4}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := layout.Pad(tc.args...); got != tc.want {
				t.Fatalf("Pad(%v) = %+v, want %+v", tc.args, got, tc.want)
			}
		})
	}
}

// TestSplitted_Assign mirrors Assign's own doc example: it stores each
// resulting Rect into the corresponding pointer in order, and silently
// skips nil pointers rather than assigning through them.
func TestSplitted_Assign(t *testing.T) {
	var top, bottom layout.Rect

	layout.Vertical(layout.Len(3), layout.Fill(1)).
		Split(image.Rect(0, 0, 10, 10)).
		Assign(&top, &bottom)

	wantTop := image.Rect(0, 0, 10, 3)
	wantBottom := image.Rect(0, 3, 10, 10)
	if top != wantTop || bottom != wantBottom {
		t.Fatalf("Assign() top=%v bottom=%v, want top=%v bottom=%v", top, bottom, wantTop, wantBottom)
	}

	bottom = image.Rectangle{}
	layout.Vertical(layout.Len(3), layout.Fill(1)).
		Split(image.Rect(0, 0, 10, 10)).
		Assign(nil, &bottom)
	if bottom != wantBottom {
		t.Fatalf("Assign(nil, &bottom) = %v, want %v (nil pointer skipped, second still assigned)", bottom, wantBottom)
	}
}

// TestRect proves Rect is a genuine image.Rectangle alias — an
// image.Rect(...) value assigns to it with no conversion, its Dx/Dy methods
// ride free, and it round-trips through Split unchanged as both the input
// area and (for a single unconstrained-by-siblings Fill) the output.
func TestRect(t *testing.T) {
	var r layout.Rect = image.Rect(1, 2, 11, 9)
	if r.Dx() != 10 || r.Dy() != 7 {
		t.Fatalf("Rect{%v}.Dx()/.Dy() = %d/%d, want 10/7", r, r.Dx(), r.Dy())
	}

	got := layout.Vertical(layout.Fill(1)).Split(r)
	if len(got) != 1 || got[0] != r {
		t.Fatalf("Split(single Fill) = %v, want [%v] unchanged", got, r)
	}
}

// TestConstraint proves the sealed Constraint interface accepts every
// constraint kind at once — Len, Min, Max, Percent, Ratio and Fill in one
// heterogeneous slice — and that Split resolves all six without error,
// returning exactly one Rect per Constraint in order.
func TestConstraint(t *testing.T) {
	cs := []layout.Constraint{
		layout.Len(1),
		layout.Min(1),
		layout.Max(9),
		layout.Percent(1),
		layout.Ratio{Num: 1, Den: 8},
		layout.Fill(1),
	}

	got := layout.Vertical(cs...).Split(image.Rect(0, 0, 10, 60))
	if len(got) != len(cs) {
		t.Fatalf("len(Split()) = %d, want %d (one Rect per Constraint, every kind accepted)", len(got), len(cs))
	}
}
