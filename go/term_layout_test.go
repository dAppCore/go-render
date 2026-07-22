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
