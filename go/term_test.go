//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mapTranslator resolves Text keys from a fixed map, mirroring how tests drive
// the translator seam without a locale catalogue.
type mapTranslator map[string]string

func (m mapTranslator) T(key string, _ ...any) string {
	if value, ok := m[key]; ok {
		return value
	}
	return key
}

func termTestContext(entries map[string]string) *Context {
	return NewContextWithService(mapTranslator(entries), "en-GB")
}

func termTestEntitledContext() *Context {
	ctx := termTestContext(map[string]string{"secret": "ops panel"})
	ctx.Entitlements = func(feature string) bool { return feature == "ops" }
	return ctx
}

// asciiProfile forces deterministic unstyled output for structural assertions;
// the returned restore function reinstates the detected profile.
func asciiProfile() func() {
	previous := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.Ascii)
	return func() { lipgloss.SetColorProfile(previous) }
}

func TestTerm_RenderTerm(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	tests := []struct {
		name     string
		node     Node
		ctx      *Context
		opts     TermOptions
		contains []string
		absent   []string
	}{
		{
			name:     "good: heading renders with rule",
			node:     El("h1", Text("page.title")),
			ctx:      termTestContext(map[string]string{"page.title": "Dashboard"}),
			contains: []string{"Dashboard", "─"},
		},
		{
			name:     "good: paragraph wraps to width",
			node:     El("p", Text("page.body")),
			ctx:      termTestContext(map[string]string{"page.body": strings.Repeat("alpha beta ", 12)}),
			opts:     TermOptions{Width: 24},
			contains: []string{"alpha beta"},
		},
		{
			name: "good: unordered list bullets items",
			node: El("ul",
				El("li", Text("one")),
				El("li", Text("two")),
			),
			ctx:      termTestContext(map[string]string{"one": "First item", "two": "Second item"}),
			contains: []string{"• First item", "• Second item"},
		},
		{
			name: "good: ordered list numbers items",
			node: El("ol",
				El("li", Text("one")),
				El("li", Text("two")),
			),
			ctx:      termTestContext(map[string]string{"one": "First", "two": "Second"}),
			contains: []string{"1. First", "2. Second"},
		},
		{
			name: "good: each expands rows",
			node: El("ul", Each([]string{"lint", "build"}, func(item string) Node {
				return El("li", Text(item))
			})),
			ctx:      termTestContext(map[string]string{"lint": "lint clean", "build": "build green"}),
			contains: []string{"• lint clean", "• build green"},
		},
		{
			name:     "good: pre keeps raw shape",
			node:     El("pre", Raw("core go qa\ncore build")),
			ctx:      termTestContext(nil),
			contains: []string{"core go qa", "core build"},
		},
		{
			name:     "good: progress renders bar and percent",
			node:     Attr(Attr(El("progress"), "value", "3"), "max", "4"),
			ctx:      termTestContext(nil),
			contains: []string{"█", "░", "75%"},
		},
		{
			name:     "good: blockquote carries the bar",
			node:     El("blockquote", El("p", Text("quote"))),
			ctx:      termTestContext(map[string]string{"quote": "ship the smallest artifact"}),
			contains: []string{"│ ship the smallest artifact"},
		},
		{
			name:     "good: entitled renders when granted",
			node:     Entitled("ops", El("p", Text("secret"))),
			ctx:      termTestEntitledContext(),
			contains: []string{"ops panel"},
		},
		{
			name:   "bad: entitled absent without grant",
			node:   Entitled("ops", El("p", Text("secret"))),
			ctx:    termTestContext(map[string]string{"secret": "ops panel"}),
			absent: []string{"ops panel"},
		},
		{
			name:   "bad: if false renders nothing",
			node:   If(func(*Context) bool { return false }, El("p", Text("hidden"))),
			ctx:    termTestContext(map[string]string{"hidden": "should not appear"}),
			absent: []string{"should not appear"},
		},
		{
			name:     "good: switch picks the matching case",
			node:     Switch(func(*Context) string { return "b" }, map[string]Node{"a": El("p", Text("left")), "b": El("p", Text("right"))}),
			ctx:      termTestContext(map[string]string{"left": "case a", "right": "case b"}),
			contains: []string{"case b"},
			absent:   []string{"case a"},
		},
		{
			name:     "good: raw content is tag-stripped",
			node:     El("p", Raw("plain <strong>bold</strong> done")),
			ctx:      termTestContext(nil),
			contains: []string{"plain bold done"},
			absent:   []string{"<strong>"},
		},
		{
			name:     "ugly: nil children and unknown tags stay calm",
			node:     El("widget", nil, El("span"), Text("page.body")),
			ctx:      termTestContext(map[string]string{"page.body": "still here"}),
			contains: []string{"still here"},
		},
		{
			name:     "ugly: zero width clamps instead of panicking",
			node:     El("p", Text("page.body")),
			ctx:      termTestContext(map[string]string{"page.body": "narrow"}),
			opts:     TermOptions{Width: 1},
			contains: []string{"narrow"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := RenderTerm(tc.node, tc.ctx, tc.opts)
			for _, want := range tc.contains {
				assert.Contains(t, out, want)
			}
			for _, unwanted := range tc.absent {
				assert.NotContains(t, out, unwanted)
			}
		})
	}
}

func TestTerm_DefinitionTermWrapsToWidth(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	// A term far longer than the render width: it must wrap onto further
	// lines like every sibling block, not overflow as one clipped line.
	term := strings.TrimSpace(strings.Repeat("alpha ", 20)) // 20 words, ~119 cols
	const width = 24

	tests := []struct {
		name string
		node Node
	}{
		{
			name: "good: standalone dt wraps",
			node: El("dt", Text("term")),
		},
		{
			name: "good: dt inside dl wraps",
			node: El("dl", El("dt", Text("term")), El("dd", Text("def"))),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := termTestContext(map[string]string{"term": term, "def": "definition"})
			out := termStripANSI(RenderTerm(tc.node, ctx, TermOptions{Width: width}))

			lines := strings.Split(out, "\n")
			assert.Greater(t, len(lines), 1, "a term wider than the render width must wrap to multiple lines")
			for _, line := range lines {
				assert.LessOrEqual(t, lipgloss.Width(line), width, "no rendered line exceeds the render width")
			}
			assert.Equal(t, 20, strings.Count(out, "alpha"), "wrapping loses no content")
		})
	}
}

func TestTerm_Verbatim_PassesAnsiByteExact(t *testing.T) {
	// A Verbatim node passes its (already terminal-ready) bytes through the
	// terminal renderer untouched: no StripTags (the <not-a-tag> survives),
	// no whitespace normalisation (the double space survives), no width
	// wrapping (a width of 4 does not fold the long line). Pre-styled ANSI
	// -- e.g. Glamour-rendered markdown -- reaches the terminal byte-for-byte.
	ansi := "\x1b[1mBOLD\x1b[0m  \x1b[38;5;42mgreen\x1b[0m\n<not-a-tag> kept"
	out := RenderTerm(Verbatim(ansi), NewContext(), TermOptions{Width: 4})
	assert.Equal(t, ansi, out, "verbatim content is emitted byte-for-byte")
}

func TestTerm_Verbatim_InsideComposedChrome(t *testing.T) {
	restore := asciiProfile()
	defer restore()
	// The intended shape: pre-styled ANSI placed inside composed chrome. The
	// verbatim bytes survive intact even though the surrounding block renders
	// through the normal styled path.
	ansi := "\x1b[31mred\x1b[0m"
	out := RenderTerm(El("div", El("h2", Text("head")), Verbatim(ansi)), termTestContext(map[string]string{"head": "Section"}))
	assert.Contains(t, out, ansi, "verbatim bytes survive inside a composed block")
	assert.Contains(t, out, "Section", "sibling chrome still renders")
}

func TestTerm_RenderTerm_Table(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	node := El("table",
		El("thead", El("tr", El("th", Text("repo")), El("th", Text("status")))),
		El("tbody",
			El("tr", El("td", Text("go-html")), El("td", Text("green"))),
			El("tr", El("td", Text("go-io")), El("td", Text("green"))),
		),
	)
	ctx := termTestContext(map[string]string{
		"repo": "Repo", "status": "Status",
		"go-html": "go-html", "go-io": "go-io", "green": "green",
	})

	out := RenderTerm(node, ctx)
	assert.Contains(t, out, "Repo")
	assert.Contains(t, out, "Status")
	assert.Contains(t, out, "go-html")
	assert.Contains(t, out, "go-io")
	assert.Contains(t, out, "─")
	require.Greater(t, strings.Count(out, "\n"), 3, "table renders a bordered multi-line frame")
}

func TestTerm_RenderTerm_Hyperlink(t *testing.T) {
	restore := asciiProfile()
	defer restore()

	linked := Attr(El("a", Text("site")), "href", "https://dappco.re")
	ctx := termTestContext(map[string]string{"site": "dappco.re"})

	out := RenderTerm(linked, ctx)
	assert.Contains(t, out, "\x1b]8;;https://dappco.re\x1b\\", "OSC 8 hyperlink opens")
	assert.Contains(t, out, "dappco.re")

	theme := DefaultTermTheme()
	theme.Hyperlinks = false
	plain := RenderTerm(linked, ctx, TermOptions{Theme: theme})
	assert.NotContains(t, plain, "\x1b]8;;", "hyperlinks disabled by theme")
	assert.Contains(t, plain, "dappco.re")
}

func TestTerm_RenderTerm_ClassTokens(t *testing.T) {
	previous := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	defer lipgloss.SetColorProfile(previous)

	node := El("p", Attr(El("span", Text("state")), "class", "error"))
	ctx := termTestContext(map[string]string{"state": "failed"})

	out := RenderTerm(node, ctx)
	assert.Contains(t, out, "\x1b[", "class token styles emit ANSI under a colour profile")
	assert.Contains(t, termStripANSI(out), "failed")
}

func TestTerm_termStripANSI(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "good: plain passes through", in: "plain", want: "plain"},
		{name: "good: sgr removed", in: "\x1b[1mbold\x1b[0m", want: "bold"},
		{name: "good: osc8 removed", in: "\x1b]8;;https://x\x1b\\text\x1b]8;;\x1b\\", want: "text"},
		{name: "ugly: truncated escape survives", in: "tail\x1b[", want: "tail"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, termStripANSI(tc.in))
		})
	}
}
