//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func termTestPage() *Layout {
	return NewLayout("HLCRF").
		H(Text("head")).
		L(El("ul", El("li", Text("navA")), El("li", Text("navB")))).
		C(El("h2", Text("title")), El("p", Text("body"))).
		R(El("p", Text("aside"))).
		F(Text("foot"))
}

func termTestPageContext() *Context {
	return termTestContext(map[string]string{
		"head": "Header band", "navA": "Overview", "navB": "Settings",
		"title": "Content title", "body": "Content body copy.",
		"aside": "Side detail", "foot": "Footer status",
	})
}

func TestTermLayout_RenderTerm(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	tests := []struct {
		name     string
		layout   *Layout
		width    int
		contains []string
		absent   []string
	}{
		{
			name:   "good: full frame renders every slot wide",
			layout: termTestPage(),
			width:  120,
			contains: []string{
				"Header band", "Overview", "Settings",
				"Content title", "Content body copy.", "Side detail", "Footer status",
				"╭", "╰", "─",
			},
		},
		{
			name:     "good: HCF variant skips sides",
			layout:   NewLayout("HCF").H(Text("head")).C(El("p", Text("body"))).F(Text("foot")),
			width:    100,
			contains: []string{"Header band", "Content body copy.", "Footer status"},
			absent:   []string{"╭"},
		},
		{
			name:     "good: empty slots render no bands",
			layout:   NewLayout("HLCRF").C(El("p", Text("body"))),
			width:    100,
			contains: []string{"Content body copy."},
			absent:   []string{"╭"},
		},
		{
			name:     "bad: unknown variant characters are ignored",
			layout:   NewLayout("HXC").C(El("p", Text("body"))),
			width:    100,
			contains: []string{"Content body copy."},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.layout.RenderTerm(termTestPageContext(), TermOptions{Width: tc.width})
			for _, want := range tc.contains {
				assert.Contains(t, out, want)
			}
			for _, unwanted := range tc.absent {
				assert.NotContains(t, out, unwanted)
			}
		})
	}
}

func TestTermLayout_RenderTerm_SideBySide(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	out := termTestPage().RenderTerm(termTestPageContext(), TermOptions{Width: 120})
	lines := strings.Split(out, "\n")

	titleLine := ""
	for _, line := range lines {
		if strings.Contains(line, "Content title") {
			titleLine = line
			break
		}
	}
	require.NotEmpty(t, titleLine, "content title line present")
	assert.Contains(t, titleLine, "╭", "sidebar top border shares the content title row at 120 columns")
	assert.Equal(t, 2, strings.Count(titleLine, "╭"), "left and right boxes open on the same row")
}

func TestTermLayout_RenderTerm_NarrowStacks(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	out := termTestPage().RenderTerm(termTestPageContext(), TermOptions{Width: 60})
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Overview") {
			assert.NotContains(t, line, "Content", "below 80 columns the frame stacks vertically")
		}
	}
	assert.Contains(t, out, "Content body copy.")
	assert.Contains(t, out, "Side detail")
}

func TestTermLayout_RenderTerm_FitSlots(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	page := NewLayout("LCR").L(Text("brand")).C(Text("mid")).R(Text("tail"))
	ctx := termTestContext(map[string]string{"brand": "Brand", "mid": "Mid", "tail": "Tail"})

	fixed := page.RenderTerm(ctx, TermOptions{Width: 100})
	fit := page.RenderTerm(ctx, TermOptions{Width: 100, FitSlots: true})

	// Both carry all the content; fit packs it into far fewer columns than the
	// fixed 24 + gutter + C + gutter + 28 frame.
	for _, want := range []string{"Brand", "Mid", "Tail"} {
		assert.Contains(t, fit, want)
	}
	assert.Less(t, lipgloss.Width(fit), lipgloss.Width(fixed), "fit output is content-width, not frame-width")

	// FitSlots is opt-in: the default render is exactly what it was before the
	// option existed.
	assert.Equal(t, fixed, page.RenderTerm(ctx, TermOptions{Width: 100}), "default render is unchanged")
}

func TestTermLayout_RegionInnerContentWidth(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	// S:S15.5: each region renders content into (region width - its horizontal
	// chrome). A full-width band (H/C/F) reserves its (0,1) gutter -> width-2; a
	// bordered L/R box reserves the default rounded border + (0,1) padding ->
	// width-4. Pin the boundary with an unbreakable token: exactly the inner width
	// fits one line, one column wider wraps to a second.
	countTokenLines := func(out, token string) int {
		n := 0
		for _, ln := range strings.Split(out, "\n") {
			if strings.Contains(ln, token) {
				n++
			}
		}
		return n
	}

	tests := []struct {
		name  string
		outer int
		inner int
		build func() *Layout
	}{
		{
			name:  "good: H band inner content width is band width minus 2",
			outer: 40, inner: 38,
			build: func() *Layout { return NewLayout("H").H(El("p", Text("x"))) },
		},
		{
			name:  "good: C content inner width is band width minus 2",
			outer: 40, inner: 38,
			build: func() *Layout { return NewLayout("C").C(El("p", Text("x"))) },
		},
		{
			name:  "good: F band inner content width is band width minus 2",
			outer: 40, inner: 38,
			build: func() *Layout { return NewLayout("F").F(El("p", Text("x"))) },
		},
		{
			name:  "good: L boxed slot inner width is the fixed sidebar budget minus 4",
			outer: 100, inner: termSidebarWidth - 4,
			build: func() *Layout { return NewLayout("LC").L(El("p", Text("x"))).C(Text("c")) },
		},
		{
			name:  "good: R boxed slot inner width is the fixed aside budget minus 4",
			outer: 100, inner: termAsideWidth - 4,
			build: func() *Layout { return NewLayout("CR").C(Text("c")).R(El("p", Text("x"))) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctxFit := termTestContext(map[string]string{"x": strings.Repeat("x", tc.inner), "c": "c"})
			ctxOver := termTestContext(map[string]string{"x": strings.Repeat("x", tc.inner+1), "c": "c"})

			fit := termStripANSI(tc.build().RenderTerm(ctxFit, TermOptions{Width: tc.outer}))
			over := termStripANSI(tc.build().RenderTerm(ctxOver, TermOptions{Width: tc.outer}))

			assert.Equal(t, 1, countTokenLines(fit, "x"), "a token exactly the inner width fits one line")
			assert.Equal(t, 2, countTokenLines(over, "x"), "one column wider wraps to a second line")
		})
	}
}

func TestTermLayout_FitToInner_ByteExactVerbatim(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	// Friction-4 re-probe, pinned. On this base (measured chrome from round 3 +
	// the S:S15.5 inner-width contract) the original ANSI corruption is closed: a
	// pre-styled Verbatim line the host fitted to a slot's INNER width (outer -
	// chrome) passes through byte-exact on one row for C, L and R alike. One
	// column wider trips the slot's Width() word-wrap onto a second row, splitting
	// the ANSI -- so the corruption the reporter saw on v0.13.0 was a host fitting
	// to the slot's nominal width, not a renderer defect. Fit to inner, bytes live.
	theme := DefaultTermTheme()
	ansiOf := func(n int) string { return "\x1b[1m" + strings.Repeat("A", n) + "\x1b[0m" }
	rows := func(out string) int {
		n := 0
		for _, ln := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			if strings.Contains(termStripANSI(ln), "A") {
				n++
			}
		}
		return n
	}

	tests := []struct {
		name  string
		outer int
		inner int
		build func(body Node) *Layout
	}{
		{
			name:  "good: C content passes byte-exact at its inner width",
			outer: 40, inner: 40 - termChrome(theme.Content),
			build: func(body Node) *Layout { return NewLayout("C").C(body) },
		},
		{
			name:  "good: L boxed slot passes byte-exact at its inner width",
			outer: 100, inner: termSidebarWidth - termChrome(theme.Sidebar),
			build: func(body Node) *Layout { return NewLayout("LC").L(body).C(Text("c")) },
		},
		{
			name:  "good: R boxed slot passes byte-exact at its inner width",
			outer: 100, inner: termAsideWidth - termChrome(theme.Aside),
			build: func(body Node) *Layout { return NewLayout("CR").C(Text("c")).R(body) },
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := termTestContext(map[string]string{"c": "c"})
			fit := tc.build(Verbatim(ansiOf(tc.inner))).RenderTerm(ctx, TermOptions{Width: tc.outer})
			over := tc.build(Verbatim(ansiOf(tc.inner + 1))).RenderTerm(ctx, TermOptions{Width: tc.outer})

			assert.Equal(t, 1, rows(fit), "a line fitted to the inner width rides one row")
			assert.Contains(t, fit, ansiOf(tc.inner), "the pre-styled ANSI survives byte-exact at the inner width")
			assert.Equal(t, 2, rows(over), "one column over the inner width word-wraps to a second row")
		})
	}
}

func TestTermChrome(t *testing.T) {
	// termChrome measures the horizontal columns a slot style spends on border and
	// padding -- what FitSlots adds to a slot's content width so its recorded box
	// spans the rendered glyphs. The default rounded, (0,1)-padded slot is 4; a
	// stripped theme is less, and the measurement is what keeps its boxes exact.
	tests := []struct {
		name  string
		style lipgloss.Style
		want  int
	}{
		{
			name:  "good: rounded border with (0,1) padding is 4 (the default slot)",
			style: lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()),
			want:  4,
		},
		{
			name:  "good: borderless and unpadded is 0",
			style: lipgloss.NewStyle(),
			want:  0,
		},
		{
			name:  "good: space-glyph left/right border with padding still counts the border",
			style: lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.Border{Left: " ", Right: " "}, false, true, false, true),
			want:  4,
		},
		{
			name:  "bad: padding only, no border",
			style: lipgloss.NewStyle().Padding(0, 2),
			want:  4,
		},
		{
			name:  "bad: full border only, no padding",
			style: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
			want:  2,
		},
		{
			name:  "ugly: a hidden border still occupies its columns",
			style: lipgloss.NewStyle().Border(lipgloss.HiddenBorder()),
			want:  2,
		},
		{
			name:  "ugly: a single left border is one column",
			style: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true),
			want:  1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, termChrome(tc.style))
		})
	}
}

func TestTermLayout_ContentSlotChrome(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	// Round 4: the C slot's alignment gutter is the themeable Content style
	// (S:S15.2), symmetric with the Sidebar/Aside/Header/Footer band styles. The
	// default keeps the (0,1) gutter that aligns C with the H/F bands; a theme may
	// zero or widen it, and renderTermContent picks the themed style up so the
	// content lands at exactly that gutter column.
	tests := []struct {
		name       string
		content    lipgloss.Style
		wantColumn int
	}{
		{"good: default (0,1) gutter aligns C content one column in", DefaultTermTheme().Content, 1},
		{"good: a zero Content style renders C content at column 0", lipgloss.NewStyle(), 0},
		{"good: a wider Content padding pushes C content further in", lipgloss.NewStyle().Padding(0, 3), 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			theme := DefaultTermTheme()
			theme.Content = tc.content
			page := NewLayout("C").C(El("p", Text("mark")))
			ctx := termTestContext(map[string]string{"mark": "MARK"})
			out := page.RenderTerm(ctx, TermOptions{Width: 40, Theme: theme})
			col, ok := glyphColumn(strings.Split(out, "\n"), "MARK")
			require.True(t, ok, "content is rendered")
			assert.Equal(t, tc.wantColumn, col, "C content sits at the themed gutter column")
		})
	}
}

func TestTermLayout_ContentSlotChrome_ByteExactVerbatim(t *testing.T) {
	// The downstream's ask: a zero-chrome Content slot passes a pre-styled ANSI
	// panel body byte-exact at the slot's full width. With the default (0,1)
	// gutter the same full-width line lands two columns over budget and the slot's
	// Width() word-wraps it, splitting the ANSI -- which is exactly why the lever
	// is needed for a host that pre-fits to the slot's nominal width (S:S15.5). No
	// asciiProfile here: the verbatim ANSI must reach output to be asserted intact.
	ansi := "\x1b[1m" + strings.Repeat("A", 40) + "\x1b[0m"
	rows := func(out string) int {
		n := 0
		for _, ln := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			if strings.Contains(termStripANSI(ln), "A") {
				n++
			}
		}
		return n
	}
	page := NewLayout("C").C(Verbatim(ansi))

	zero := DefaultTermTheme()
	zero.Content = lipgloss.NewStyle()
	fit := page.RenderTerm(NewContext(), TermOptions{Width: 40, Theme: zero})
	assert.Equal(t, 1, rows(fit), "zero-chrome C passes a full-width pre-styled line on one row")
	assert.Contains(t, fit, ansi, "the pre-styled ANSI survives byte-exact")

	wrapped := page.RenderTerm(NewContext(), TermOptions{Width: 40})
	assert.Equal(t, 2, rows(wrapped), "the default (0,1) gutter word-wraps a full-width line")
}

func TestTermTheme_Content_DefaultByteIdentity(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	// Regression pin (round 4): making the C gutter a theme field must not change
	// the shipped theme's output. The default Content is exactly (0,1), so its
	// measured chrome stays the documented 2 (S:S15.5) and a full HLCRF frame under
	// DefaultTermTheme is byte-identical to one whose Content is an explicit (0,1).
	assert.Equal(t, 2, termChrome(DefaultTermTheme().Content), "default C chrome is the documented 2")

	page := termTestPage()
	ctx := termTestPageContext()
	explicit := DefaultTermTheme()
	explicit.Content = lipgloss.NewStyle().Padding(0, 1)

	assert.Equal(t,
		page.RenderTerm(ctx, TermOptions{Width: 100}),
		page.RenderTerm(ctx, TermOptions{Width: 100, Theme: explicit}),
		"default-theme frame is byte-identical to an explicit (0,1) Content",
	)
}

func TestTermLayout_Responsive_RenderTerm(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	wide := NewLayout("C").C(El("p", Text("wide")))
	narrow := NewLayout("C").C(El("p", Text("narrow")))
	resp := NewResponsive().
		Variant("desktop", wide).
		Variant("mobile", narrow)
	ctx := termTestContext(map[string]string{"wide": "wide copy", "narrow": "narrow copy"})

	tests := []struct {
		name   string
		width  int
		want   string
		unwant string
	}{
		{name: "good: desktop at 120", width: 120, want: "wide copy", unwant: "narrow copy"},
		{name: "good: mobile below 80", width: 60, want: "narrow copy", unwant: "wide copy"},
		{name: "good: tablet falls back to desktop", width: 100, want: "wide copy", unwant: "narrow copy"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := resp.RenderTerm(ctx, TermOptions{Width: tc.width})
			assert.Contains(t, out, tc.want)
			assert.NotContains(t, out, tc.unwant)
		})
	}

	t.Run("ugly: empty responsive renders empty", func(t *testing.T) {
		assert.Equal(t, "", NewResponsive().RenderTerm(ctx))
	})
}
