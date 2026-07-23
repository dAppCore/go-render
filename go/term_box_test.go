//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTermBoxes_MatchesRenderTerm(t *testing.T) {
	page := termTestPage()
	ctx := termTestPageContext()
	opts := TermOptions{Width: 120}

	plain := page.RenderTerm(ctx, opts)
	boxed, boxes := page.RenderTermBoxes(ctx, opts)

	assert.Equal(t, plain, boxed, "recording boxes must not perturb the rendered string")
	assert.NotEmpty(t, boxes)
}

func TestLayout_RenderTermBoxes_SideBySide(t *testing.T) {
	page := termTestPage()
	_, boxes := page.RenderTermBoxes(termTestPageContext(), TermOptions{Width: 120})

	for _, id := range []string{"H", "L", "C", "R", "F"} {
		require.Contains(t, boxes, id, "slot %q recorded", id)
	}

	h, l, c, r, f := boxes["H"], boxes["L"], boxes["C"], boxes["R"], boxes["F"]

	assert.Equal(t, 0, h.Row, "H starts at the top")
	assert.Equal(t, 0, h.Col)
	assert.Equal(t, 120, h.Width)

	assert.Equal(t, l.Row, c.Row, "L and C share a row at 120 columns")
	assert.Equal(t, c.Row, r.Row, "C and R share a row at 120 columns")
	assert.Equal(t, l.Row, h.Row+h.Height, "the middle band starts right after H")

	assert.Equal(t, 0, l.Col, "L starts at the left edge")
	assert.Less(t, l.Col+l.Width, c.Col, "C starts to the right of L, with a gap")
	assert.Less(t, c.Col+c.Width, r.Col, "R starts to the right of C, with a gap")

	assert.Equal(t, f.Row, l.Row+l.Height, "F starts right after the middle band")
	assert.Equal(t, 0, f.Col)

	for id, box := range boxes {
		assert.Greater(t, box.Width, 0, "%s has a positive width", id)
		assert.Greater(t, box.Height, 0, "%s has a positive height", id)
		assert.Same(t, page, box.Node, "%s carries the Layout it belongs to", id)
	}
}

func TestLayout_RenderTermBoxes_NarrowStacks(t *testing.T) {
	page := termTestPage()
	_, boxes := page.RenderTermBoxes(termTestPageContext(), TermOptions{Width: 60})

	l, c, r := boxes["L"], boxes["C"], boxes["R"]
	require.NotZero(t, l.Height)
	require.NotZero(t, c.Height)
	require.NotZero(t, r.Height)

	assert.Equal(t, 0, l.Col, "stacked columns all start at the left edge")
	assert.Equal(t, 0, c.Col)
	assert.Equal(t, 0, r.Col)
	assert.Equal(t, 60, l.Width)
	assert.Equal(t, 60, c.Width)
	assert.Equal(t, 60, r.Width)

	assert.Less(t, l.Row, c.Row, "L stacks above C below the 80-column threshold")
	assert.Less(t, c.Row, r.Row, "C stacks above R below the 80-column threshold")
	assert.Equal(t, c.Row, l.Row+l.Height, "C starts exactly where L ends")
	assert.Equal(t, r.Row, c.Row+c.Height, "R starts exactly where C ends")

	for _, box := range boxes {
		assert.Same(t, page, box.Node)
	}
}

func TestLayout_RenderTermBoxes_Nested(t *testing.T) {
	inner := NewLayout("HCF").
		H(Text("in.h")).
		C(El("p", Text("in.c"))).
		F(Text("in.f"))
	outer := NewLayout("HLCF").
		H(Text("out.h")).
		L(El("ul", El("li", Text("nav")))).
		C(El("h2", Text("out.title")), inner).
		F(Text("out.f"))
	ctx := termTestContext(map[string]string{
		"out.h": "Outer header", "nav": "Menu", "out.title": "Outer content",
		"in.h": "Inner header", "in.c": "Inner body", "in.f": "Inner footer",
		"out.f": "Outer footer",
	})

	_, boxes := outer.RenderTermBoxes(ctx, TermOptions{Width: 92})

	for _, id := range []string{"H", "L", "C", "F"} {
		require.Contains(t, boxes, id, "outer slot %q recorded with clean keys", id)
	}
	for _, id := range []string{"L1.H", "L1.C", "L1.F"} {
		require.Contains(t, boxes, id, "nested layout slot %q recorded with a disambiguating prefix", id)
	}
	assert.NotContains(t, boxes, "L1.L", "the nested layout has no L slot to record")

	outerC := boxes["C"]
	innerH := boxes["L1.H"]
	innerF := boxes["L1.F"]

	assert.GreaterOrEqual(t, innerH.Row, outerC.Row, "the nested frame renders inside the outer C slot's row range")
	assert.Less(t, innerF.Row, outerC.Row+outerC.Height, "the nested frame's footer still falls inside the outer C slot")
	assert.Greater(t, innerH.Col, 0, "the nested frame is indented into the outer content column, not at column zero")
	assert.Same(t, outer, outerC.Node)
	assert.Same(t, inner, innerH.Node)
}

func TestLayout_RenderTermBoxes_FitSlots(t *testing.T) {
	// Friction 1 residual: a content-packed strip rides L/C/R slots. FitSlots
	// sizes each slot to its own content and packs them edge-to-edge, so the
	// recorded boxes tile the row exactly rather than sitting in the fixed
	// 24/28 budgets with C filling a huge middle.
	page := NewLayout("LCR").L(Text("brand")).C(Text("mid")).R(Text("tail"))
	ctx := termTestContext(map[string]string{"brand": "Brand", "mid": "Mid", "tail": "Tail"})

	_, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 100, FitSlots: true})
	l, c, r := boxes["L"], boxes["C"], boxes["R"]
	require.NotZero(t, l.Width)
	require.NotZero(t, c.Width)
	require.NotZero(t, r.Width)

	assert.Equal(t, 0, l.Col, "L opens at the left edge")
	assert.Equal(t, l.Col+l.Width, c.Col, "C abuts L exactly -- fit mode drops the inter-slot gutter")
	assert.Equal(t, c.Col+c.Width, r.Col, "R abuts C exactly, so the three boxes tile the row")
	assert.Equal(t, l.Row, c.Row, "the strip is a single row")
	assert.Equal(t, c.Row, r.Row)

	assert.Less(t, l.Width, termSidebarWidth, "L is content-sized, narrower than the fixed L budget")
	assert.Less(t, r.Width, termAsideWidth, "R is content-sized, narrower than the fixed R budget")
	assert.Less(t, r.Col+r.Width, 100, "the packed strip does not fill the frame")

	for _, box := range boxes {
		assert.Same(t, page, box.Node)
	}

	// Default (no FitSlots) keeps the fixed budgets and the inter-slot gutter,
	// unchanged by the new option.
	_, fixed := page.RenderTermBoxes(ctx, TermOptions{Width: 100})
	assert.Equal(t, termSidebarWidth, fixed["L"].Width, "default L keeps the fixed budget")
	assert.Equal(t, termAsideWidth, fixed["R"].Width, "default R keeps the fixed budget")
	assert.Less(t, fixed["L"].Col+fixed["L"].Width, fixed["C"].Col, "default keeps a gutter between L and C")
}

// glyphColumn finds the first line containing label and returns the display
// column it starts at, after ANSI stripping -- used to check a slot's rendered
// glyphs fall inside its recorded box.
func glyphColumn(lines []string, label string) (int, bool) {
	for _, ln := range lines {
		s := termStripANSI(ln)
		if i := strings.Index(s, label); i >= 0 {
			return lipgloss.Width(s[:i]), true
		}
	}
	return 0, false
}

func TestLayout_RenderTermBoxes_FitSlots_ThemeChrome(t *testing.T) {
	// Friction (round 3): FitSlots measures each L/R slot's chrome from the active
	// theme's slot style instead of assuming the default rounded, padded slot's
	// +4. So a borderless or space-glyph theme records boxes that tile the
	// *visible* glyphs exactly, rather than drifting a column or two wide of them
	// as the old fixed constant did.
	page := NewLayout("LCR").L(Text("brand")).C(Text("mid")).R(Text("tail"))
	ctx := termTestContext(map[string]string{"brand": "Brand", "mid": "Mid", "tail": "Tail"})

	borderless := DefaultTermTheme()
	borderless.Sidebar = lipgloss.NewStyle()
	borderless.Aside = lipgloss.NewStyle()

	spaceGlyph := lipgloss.NewStyle().Padding(0, 1).
		Border(lipgloss.Border{Left: " ", Right: " "}, false, true, false, true)
	spaceTheme := DefaultTermTheme()
	spaceTheme.Sidebar = spaceGlyph
	spaceTheme.Aside = spaceGlyph

	tests := []struct {
		name     string
		theme    *TermTheme
		wantRows int
	}{
		{
			name:     "good: default bordered theme records exact boxes (regression pin)",
			theme:    DefaultTermTheme(),
			wantRows: 3,
		},
		{
			name:     "good: borderless theme boxes match the visible glyphs",
			theme:    borderless,
			wantRows: 1,
		},
		{
			name:     "good: space-glyph border theme keeps the downstream idiom exact",
			theme:    spaceTheme,
			wantRows: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 100, FitSlots: true, Theme: tc.theme})
			l, c, r := boxes["L"], boxes["C"], boxes["R"]

			lines := strings.Split(out, "\n")
			assert.Len(t, lines, tc.wantRows, "the strip renders on the expected number of rows")

			// One packed strip: every slot box shares the strip's top row.
			assert.Equal(t, l.Row, c.Row, "L and C share the strip's top row")
			assert.Equal(t, c.Row, r.Row, "C and R share the strip's top row")

			// The boxes tile edge-to-edge with no inter-slot gutter.
			assert.Equal(t, l.Col+l.Width, c.Col, "C abuts L exactly")
			assert.Equal(t, c.Col+c.Width, r.Col, "R abuts C exactly")

			// No drift: the recorded extent equals the rendered strip's own visible
			// width -- the whole point of measuring chrome from the active theme.
			visible := 0
			for _, ln := range lines {
				if w := lipgloss.Width(termStripANSI(ln)); w > visible {
					visible = w
				}
			}
			assert.Equal(t, visible, r.Col+r.Width, "the recorded boxes span exactly the visible strip -- no chrome drift")

			// Each slot's own glyphs fall inside its recorded box.
			for _, probe := range []struct {
				box   Box
				label string
			}{{l, "Brand"}, {c, "Mid"}, {r, "Tail"}} {
				col, ok := glyphColumn(lines, probe.label)
				require.True(t, ok, "%s is rendered", probe.label)
				assert.GreaterOrEqual(t, col, probe.box.Col, "%s starts at or after its slot's left edge", probe.label)
				assert.Less(t, col, probe.box.Col+probe.box.Width, "%s ends before its slot's right edge", probe.label)
			}
		})
	}

	// The borderless theme is content-sized to zero chrome, so its slots and its
	// whole strip are strictly narrower than the default bordered theme's -- proof
	// the overhead is measured, not the fixed +4.
	_, bordered := page.RenderTermBoxes(ctx, TermOptions{Width: 100, FitSlots: true})
	_, bare := page.RenderTermBoxes(ctx, TermOptions{Width: 100, FitSlots: true, Theme: borderless})
	assert.Less(t, bare["L"].Width, bordered["L"].Width, "a borderless L is narrower than a bordered L")
	assert.Less(t, bare["R"].Col+bare["R"].Width, bordered["R"].Col+bordered["R"].Width, "the borderless strip is narrower overall")
}

func TestLayout_RenderTermBoxes_SlotWidthsForPresizing(t *testing.T) {
	// Round 4, friction 1: a downstream pre-sizing content reads each slot's outer
	// width from the box map it already holds for mouse resolution (S:S14), rather
	// than hardcoding the fixed 24/28 budgets. The recorded widths are the render-
	// time source of truth: outer width from the box, inner width by subtracting
	// the slot's themed chrome (S:S15.5). No accessor for the budgets is exported.
	page := NewLayout("LCR").
		L(El("p", Text("l"))).C(El("p", Text("c"))).R(El("p", Text("r")))
	ctx := termTestContext(map[string]string{"l": "L", "c": "C", "r": "R"})
	theme := DefaultTermTheme()

	_, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 100, Theme: theme})
	require.Contains(t, boxes, "L")
	require.Contains(t, boxes, "C")
	require.Contains(t, boxes, "R")

	// The fixed side budgets are readable straight off the boxes -- no magic constant.
	assert.Equal(t, termSidebarWidth, boxes["L"].Width, "L outer width is its fixed budget, read from the box")
	assert.Equal(t, termAsideWidth, boxes["R"].Width, "R outer width is its fixed budget, read from the box")

	// The three slot boxes plus the two single-space gutters span the frame, so the
	// box widths reconstruct the layout arithmetic a caller would otherwise redo.
	assert.Equal(t, 100, boxes["L"].Width+1+boxes["C"].Width+1+boxes["R"].Width, "slot boxes plus gutters span the frame width")

	// The pre-sizing recipe: inner width = the box's outer width minus the slot's
	// themed chrome (S:S15.5) -- what a host fits pre-styled content to.
	assert.Equal(t, termSidebarWidth-4, boxes["L"].Width-termChrome(theme.Sidebar), "L inner width from box width minus themed chrome")
	assert.Equal(t, termAsideWidth-4, boxes["R"].Width-termChrome(theme.Aside), "R inner width from box width minus themed chrome")
}

func TestLayout_RenderTermBoxes_FitSlots_OneRowBelowStackThreshold(t *testing.T) {
	// FitSlots is meant for a strip that rides one row whatever the width, so it
	// bypasses the narrow-width stacking the default middle band applies below 80
	// columns.
	page := NewLayout("LCR").L(Text("a")).C(Text("b")).R(Text("c"))
	ctx := termTestContext(map[string]string{"a": "A", "b": "B", "c": "C"})

	_, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 60, FitSlots: true})
	l, c, r := boxes["L"], boxes["C"], boxes["R"]
	assert.Equal(t, l.Row, c.Row, "fit mode stays one row below the stack threshold")
	assert.Equal(t, c.Row, r.Row)
	assert.Equal(t, l.Col+l.Width, c.Col, "boxes still tile edge-to-edge")
	assert.Equal(t, c.Col+c.Width, r.Col)
}

func TestRenderTermBoxes_ElementID(t *testing.T) {
	page := NewLayout("C").C(
		Attr(El("div", El("p", Text("a"))), "id", "card-a"),
		Attr(El("div", El("p", Text("b"))), "id", "card-b"),
	)
	ctx := termTestContext(map[string]string{"a": "Card A", "b": "Card B"})

	out, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 40})
	require.Contains(t, out, "Card A")
	require.Contains(t, out, "Card B")

	require.Contains(t, boxes, "card-a")
	require.Contains(t, boxes, "card-b")
	assert.Less(t, boxes["card-a"].Row, boxes["card-b"].Row, "card-a renders above card-b")
	assert.Greater(t, boxes["card-a"].Height, 0)
	assert.Greater(t, boxes["card-b"].Height, 0)
}

func TestRenderTermBoxes_Ugly(t *testing.T) {
	t.Run("ugly: nil node returns an empty box map, not a panic", func(t *testing.T) {
		out, boxes := RenderTermBoxes(nil, NewContext())
		assert.Equal(t, "", out)
		assert.Empty(t, boxes)
	})

	t.Run("ugly: nil layout returns an empty box map, not a panic", func(t *testing.T) {
		var l *Layout
		out, boxes := l.RenderTermBoxes(NewContext())
		assert.Equal(t, "", out)
		assert.Empty(t, boxes)
	})

	t.Run("ugly: nil responsive returns an empty box map, not a panic", func(t *testing.T) {
		var r *Responsive
		out, boxes := r.RenderTermBoxes(NewContext())
		assert.Equal(t, "", out)
		assert.Empty(t, boxes)
	})

	t.Run("ugly: empty responsive returns an empty box map, not a panic", func(t *testing.T) {
		out, boxes := NewResponsive().RenderTermBoxes(NewContext())
		assert.Equal(t, "", out)
		assert.Empty(t, boxes)
	})

	t.Run("ugly: element with an empty id is not recorded", func(t *testing.T) {
		page := NewLayout("C").C(Attr(El("div", Text("x")), "id", ""))
		ctx := termTestContext(map[string]string{"x": "x"})
		_, boxes := page.RenderTermBoxes(ctx, TermOptions{Width: 40})
		assert.NotContains(t, boxes, "")
	})
}

func TestResponsive_RenderTermBoxes(t *testing.T) {
	wide := NewLayout("C").C(El("p", Text("wide")))
	narrow := NewLayout("C").C(El("p", Text("narrow")))
	resp := NewResponsive().Variant("desktop", wide).Variant("mobile", narrow)
	ctx := termTestContext(map[string]string{"wide": "wide copy", "narrow": "narrow copy"})

	out, boxes := resp.RenderTermBoxes(ctx, TermOptions{Width: 120})
	assert.Contains(t, out, "wide copy")
	require.Contains(t, boxes, "C")
	assert.Same(t, wide, boxes["C"].Node)
}
