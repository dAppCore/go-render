// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"testing"

	html "dappco.re/go/html/engine/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParse_Pipe_BuiltinNamesGood covers six of the seven builtin pipe names
// (docs/ctml.md S:S8.7) -- ago needs an argument to mean anything and gets
// its own test below. Each renders identically to the hand-built tree a
// caller would get by calling html.FormatValue directly: ctml's pipe wiring
// is what is under test here, not go-i18n's own formatting correctness
// (that is covered directly, with pinned expected strings, in
// dappco.re/go/html's own formatter_test.go).
func TestParse_Pipe_BuiltinNamesGood(t *testing.T) {
	tests := []struct {
		name string
		pipe string
		val  any
	}{
		{"number", "number", 1234567},
		{"decimal", "decimal", 1234.5},
		{"percent", "percent", 0.855},
		{"ordinal", "ordinal", 1},
		{"size", "size", 1536000},
		{"bytes", "bytes", 1536000},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := `<p>{{ n | ` + tc.pipe + ` }}</p>`
			bnd := []Bindings{{Values: map[string]any{"n": tc.val}}}
			want := html.El("p", html.Text(html.FormatValue(tc.pipe, tc.val)))
			assertSameRender(t, src, bnd, want, html.NewContext())
		})
	}
}

func TestParse_Pipe_AgoWithArgGood(t *testing.T) {
	src := `<p>{{ n | ago:minutes }}</p>`
	bnd := []Bindings{{Values: map[string]any{"n": 5}}}
	want := html.El("p", html.Text(html.FormatValue("ago", 5, "minutes")))
	assertSameRender(t, src, bnd, want, html.NewContext())
}

func TestParse_Pipe_WhitespaceVariantsGood(t *testing.T) {
	// {{ x | ago : minutes }} == {{x|ago:minutes}} -- whitespace around | and
	// : is tolerated (docs/ctml.md S:S8.7).
	bnd := Bindings{Values: map[string]any{"n": 5}}

	tight, err := Parse([]byte(`<p>{{n|ago:minutes}}</p>`), bnd)
	require.NoError(t, err)
	spaced, err := Parse([]byte(`<p>{{ n | ago : minutes }}</p>`), bnd)
	require.NoError(t, err)
	mixed, err := Parse([]byte(`<p>{{n | ago:minutes }}</p>`), bnd)
	require.NoError(t, err)

	ctx := html.NewContext()
	want := html.Render(tight, ctx)
	assert.Equal(t, want, html.Render(spaced, ctx), "fully-spaced form renders identically")
	assert.Equal(t, want, html.Render(mixed, ctx), "partially-spaced form renders identically")
}

func TestParse_Pipe_Bad(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantMsg string
	}{
		{"bad: unknown pipe name in a text bind", `<p>{{ n | bogus }}</p>`, "unknown pipe"},
		{"bad: unknown pipe name names the pipe", `<p>{{ n | bogus }}</p>`, "bogus"},
		{"bad: malformed pipe identifier", `<p>{{ n | 1bad }}</p>`, "malformed pipe"},
		{"bad: empty pipe name", `<p>{{ n | }}</p>`, "malformed pipe"},
		{"bad: unknown pipe in a verbatim bind", `<verbatim value="{{ n | bogus }}"/>`, "unknown pipe"},
		{"bad: pipe in an attribute bind", `<div class="{{ n | number }}">x</div>`, "attribute bind must not use a pipe"},
		{"bad: unknown pipe in an attribute bind is still the attribute-bind error", `<div class="{{ n | bogus }}">x</div>`, "attribute bind must not use a pipe"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.src), Bindings{Values: map[string]any{"n": 5}})
			require.Error(t, err)
			var pe *ParseError
			require.ErrorAs(t, err, &pe)
			assert.Contains(t, pe.Msg, tc.wantMsg)
			assert.Greater(t, pe.Line, 0, "line is populated")
			assert.Greater(t, pe.Col, 0, "column is populated")
		})
	}
}

func TestParse_Pipe_ErrorNamesTheLineUgly(t *testing.T) {
	// Three lines: the offending pipe is on line 2.
	src := "<div>\n<p>{{ n | bogus }}</p>\n</div>"
	_, err := Parse([]byte(src), Bindings{Values: map[string]any{"n": 5}})
	require.Error(t, err)
	var pe *ParseError
	require.ErrorAs(t, err, &pe)
	assert.Contains(t, pe.Msg, "bogus")
	assert.Equal(t, 2, pe.Line, "error reported on the line the offending bind is on")
}

func TestParse_Pipe_UnpipedNearMissStillLenientGood(t *testing.T) {
	// A "|"-free near-miss keeps S:S8.4's existing leniency -- pipes only
	// tighten the rule for tokens that actually contain a "|".
	got, err := Parse([]byte(`<p>See {{oops!}} for details</p>`))
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	assert.Contains(t, out, "oops!")
}

func TestParse_Verbatim_PipeGood(t *testing.T) {
	src := `<verbatim value="{{ n | number }}"/>`
	bnd := []Bindings{{Values: map[string]any{"n": 1234567}}}
	want := html.Verbatim(html.FormatValue("number", 1234567))
	assertSameRender(t, src, bnd, want, html.NewContext())
}

func TestParse_Verbatim_PipeRowScopedGood(t *testing.T) {
	src := `<each items="turns" as="turn"><verbatim value="{{turn.n|number}}"/></each>`
	bnd := []Bindings{{Sequences: map[string][]map[string]any{
		"turns": {{"n": 1000}, {"n": 2000000}},
	}}}
	want := html.Each([]map[string]any{{"n": 1000}, {"n": 2000000}}, func(row map[string]any) html.Node {
		return html.Verbatim(html.FormatValue("number", row["n"]))
	})
	got := assertSameRender(t, src, bnd, want, html.NewContext())

	out := html.Render(got, html.NewContext())
	assert.Contains(t, out, html.FormatValue("number", 1000))
	assert.Contains(t, out, html.FormatValue("number", 2000000))
}

func TestParse_Pipe_MissRendersEmptyGood(t *testing.T) {
	// A missing bind stays data absence (S:S8.3) even with a pipe attached --
	// it must not render a formatted zero.
	src := `<p>{{ missing | number }}</p>`
	got, err := Parse([]byte(src), Bindings{Values: map[string]any{"present": 1}})
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	assert.Equal(t, "<p></p>", out, "a missing piped bind renders empty, not a formatted zero")
}

func TestParse_Verbatim_PipeMissRendersEmptyGood(t *testing.T) {
	src := `<verbatim value="{{ missing | number }}"/>`
	got, err := Parse([]byte(src), Bindings{Values: map[string]any{"present": 1}})
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	assert.Equal(t, "", out, "a missing piped verbatim bind renders an empty Verbatim")
}

// TestParse_Pipe_Width_FitSlotsMeasuresPostFormatStringGood is the width
// requirement: a formatted pipe changes the rendered string's width, so
// FitSlots (docs/ctml.md S:S15.1) must measure the POST-format string, not
// the raw bound value. The theme's fixed slot chrome is identical for both
// renders, so it cancels out of the width delta -- only the string-length
// difference the pipe introduces should remain.
func TestParse_Pipe_Width_FitSlotsMeasuresPostFormatStringGood(t *testing.T) {
	raw := `<layout variant="C"><c>{{ n }}</c></layout>`
	piped := `<layout variant="C"><c>{{ n | number }}</c></layout>`
	bnd := Bindings{Values: map[string]any{"n": 1234567}}

	rawLayout, err := ParseLayout([]byte(raw), bnd)
	require.NoError(t, err)
	pipedLayout, err := ParseLayout([]byte(piped), bnd)
	require.NoError(t, err)

	ctx := html.NewContext()
	_, rawBoxes := rawLayout.RenderTermBoxes(ctx, html.TermOptions{Width: 100, FitSlots: true})
	_, pipedBoxes := pipedLayout.RenderTermBoxes(ctx, html.TermOptions{Width: 100, FitSlots: true})

	rawWidth, pipedWidth := rawBoxes["C"].Width, pipedBoxes["C"].Width
	require.NotZero(t, rawWidth, "the raw slot recorded a box")
	require.NotZero(t, pipedWidth, "the piped slot recorded a box")

	gotDelta := pipedWidth - rawWidth
	wantDelta := len(html.FormatValue("number", 1234567)) - len("1234567")
	require.Greater(t, wantDelta, 0, "sanity: grouping commas actually widen the string")
	assert.Equal(t, wantDelta, gotDelta,
		"FitSlots measured the piped bind's POST-format string width, not the raw value's -- "+
			"the C slot's box width must widen by exactly the formatted string's extra length")
}
