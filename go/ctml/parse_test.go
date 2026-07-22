// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"strings"
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertSameRender parses src and checks it renders byte-identically to
// want -- built via the Go API -- under both Render and RenderTerm, with
// the same Context. Both sides go through the same process-wide lipgloss
// colour profile, so this holds regardless of what that profile resolves
// to; only relative equality is asserted.
func assertSameRender(t *testing.T, src string, bindings []Bindings, want html.Node, ctx *html.Context) html.Node {
	t.Helper()
	got, err := Parse([]byte(src), bindings...)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, html.Render(want, ctx), html.Render(got, ctx), "HTML render")
	assert.Equal(t, html.RenderTerm(want, ctx), html.RenderTerm(got, ctx), "term render")
	return got
}

func TestParse_Good(t *testing.T) {
	tests := []struct {
		name string
		src  string
		bnd  []Bindings
		want html.Node
		ctx  *html.Context
	}{
		{
			name: "good: bare text is the i18n key",
			src:  `<p>page.title</p>`,
			want: html.El("p", html.Text("page.title")),
			ctx:  html.NewContext(),
		},
		{
			name: "good: attributes pass through to Attr, sorted deterministically like El",
			src:  `<a href="/docs" class="link" id="docs-link">docs.label</a>`,
			want: html.Attr(html.Attr(html.Attr(html.El("a", html.Text("docs.label")), "href", "/docs"), "class", "link"), "id", "docs-link"),
			ctx:  html.NewContext(),
		},
		{
			name: "good: void element self-closes",
			src:  `<div><br/><hr/></div>`,
			want: html.El("div", html.El("br"), html.El("hr")),
			ctx:  html.NewContext(),
		},
		{
			name: "good: mixed text and element children interleave in order, inline spacing preserved",
			src:  `<p>Hello <strong>world</strong>!</p>`,
			want: html.El("p", html.Text("Hello "), html.El("strong", html.Text("world")), html.Text("!")),
			ctx:  html.NewContext(),
		},
		{
			name: "good: raw is unescaped and whitespace-preserved",
			src:  "<pre><raw>$ core go qa\ncore build</raw></pre>",
			want: html.El("pre", html.Raw("$ core go qa\ncore build")),
			ctx:  html.NewContext(),
		},
		{
			name: "good: if renders its child when the data key is truthy",
			src:  `<if cond="admin"><p>admin.panel</p></if>`,
			want: html.If(func(ctx *html.Context) bool { return true }, html.El("p", html.Text("admin.panel"))),
			ctx:  &html.Context{Data: map[string]any{"admin": true}},
		},
		{
			name: "good: if omits its child when the data key is falsy",
			src:  `<if cond="admin"><p>admin.panel</p></if>`,
			want: html.If(func(ctx *html.Context) bool { return false }, html.El("p", html.Text("admin.panel"))),
			ctx:  html.NewContext(),
		},
		{
			name: "good: unless inverts the same truthiness rule",
			src:  `<unless cond="guest"><p>member.panel</p></unless>`,
			want: html.Unless(func(ctx *html.Context) bool { return false }, html.El("p", html.Text("member.panel"))),
			ctx:  html.NewContext(),
		},
		{
			name: "good: if wraps multiple children in a fragment",
			src:  `<if cond="on"><p>a</p><p>b</p></if>`,
			want: html.If(func(*html.Context) bool { return true }, html.Each([]html.Node{html.El("p", html.Text("a")), html.El("p", html.Text("b"))}, func(n html.Node) html.Node { return n })),
			ctx:  &html.Context{Data: map[string]any{"on": true}},
		},
		{
			name: "good: switch selects the matching case",
			src:  `<switch on="locale"><case value="en"><p>hi</p></case><case value="fr"><p>bonjour</p></case></switch>`,
			want: html.Switch(func(*html.Context) string { return "fr" }, map[string]html.Node{"en": html.El("p", html.Text("hi")), "fr": html.El("p", html.Text("bonjour"))}),
			ctx:  &html.Context{Data: map[string]any{"locale": "fr"}},
		},
		{
			name: "good: entitled renders when granted",
			src:  `<entitled feature="ops"><p>demo.ops</p></entitled>`,
			want: html.Entitled("ops", html.El("p", html.Text("demo.ops"))),
			ctx:  &html.Context{Entitlements: func(f string) bool { return f == "ops" }},
		},
		{
			name: "good: entitled is absent when denied (deny-by-default)",
			src:  `<entitled feature="ops"><p>demo.ops</p></entitled>`,
			want: html.Entitled("ops", html.El("p", html.Text("demo.ops"))),
			ctx:  html.NewContext(),
		},
		{
			name: "good: each expands bound rows with a whole-run bind and args",
			src:  `<ul><each items="repos" as="row"><li args="{{row.status}}">repo.name</li></each></ul>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"repos": {{"status": "green"}, {"status": "amber"}},
			}}},
			want: html.El("ul", html.Each([]map[string]any{{"status": "green"}, {"status": "amber"}}, func(row map[string]any) html.Node {
				return html.El("li", html.Text("repo.name", row["status"]))
			})),
			ctx: html.NewContext(),
		},
		{
			// Two valid {{}} tokens glued in one run are not "the entire
			// trimmed content is exactly one token" (S:S8.3), so the whole
			// run stays literal text rather than partially resolving.
			name: "good: two bind tokens glued in one run stay literal, not partially resolved",
			src:  `<each items="repos" as="row"><li>{{row.name}}/{{row.status}}</li></each>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"repos": {{"name": "go-html", "status": "green"}},
			}}},
			want: html.Each([]map[string]any{{"name": "go-html", "status": "green"}}, func(row map[string]any) html.Node {
				return html.El("li", html.Text("{{row.name}}/{{row.status}}"))
			}),
			ctx: html.NewContext(),
		},
		{
			name: "good: each item field as whole text-run content",
			src:  `<ul><each items="repos" as="row"><li>{{row.name}}</li></each></ul>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"repos": {{"name": "go-html"}, {"name": "go-io"}},
			}}},
			want: html.El("ul", html.Each([]map[string]any{{"name": "go-html"}, {"name": "go-io"}}, func(row map[string]any) html.Node {
				return html.El("li", html.Text(row["name"].(string)))
			})),
			ctx: html.NewContext(),
		},
		{
			name: "good: each with an unbound items name renders as an empty list",
			src:  `<ul><each items="missing" as="row"><li>{{row.name}}</li></each></ul>`,
			want: html.El("ul", html.Each([]map[string]any{}, func(row map[string]any) html.Node { return html.El("li") })),
			ctx:  html.NewContext(),
		},
		{
			// {{g.name}}/{{r.n}} glued in one run is deliberately NOT a
			// binding (S:S8.4 -- a run is the whole token or nothing), so
			// each field gets its own element to stay a whole-run bind.
			name: "good: nested each resolves both outer and inner fields",
			src:  `<each items="groups" as="g"><div>{{g.name}}<each items="repos" as="r"><span>{{r.n}}</span></each></div></each>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"groups": {{"name": "core"}},
				"repos":  {{"n": "go-html"}},
			}}},
			want: html.Each([]map[string]any{{"name": "core"}}, func(g map[string]any) html.Node {
				return html.El("div", html.Text(g["name"].(string)), html.Each([]map[string]any{{"n": "go-html"}}, func(r map[string]any) html.Node {
					return html.El("span", html.Text(r["n"].(string)))
				}))
			}),
			ctx: html.NewContext(),
		},
		{
			name: "good: nested field path resolves through a nested map",
			src:  `<each items="rows" as="row"><p>{{row.addr.city}}</p></each>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"rows": {{"addr": map[string]any{"city": "London"}}},
			}}},
			want: html.Each([]map[string]any{{"addr": map[string]any{"city": "London"}}}, func(row map[string]any) html.Node {
				return html.El("p", html.Text("London"))
			}),
			ctx: html.NewContext(),
		},
		{
			name: "good: responsive picks the matching variant",
			src:  `<responsive><variant name="desktop"><layout variant="C"><c><p>wide</p></c></layout></variant><variant name="mobile"><layout variant="C"><c><p>narrow</p></c></layout></variant></responsive>`,
			want: html.NewResponsive().
				Variant("desktop", html.NewLayout("C").C(html.El("p", html.Text("wide")))).
				Variant("mobile", html.NewLayout("C").C(html.El("p", html.Text("narrow")))),
			ctx: html.NewContext(),
		},
		{
			name: "good: responsive variant carries an optional media hint",
			src:  `<responsive><variant name="desktop" media="(min-width: 1024px)"><layout variant="C"><c><p>wide</p></c></layout></variant></responsive>`,
			want: html.NewResponsive().Add("desktop", html.NewLayout("C").C(html.El("p", html.Text("wide"))), "(min-width: 1024px)"),
			ctx:  html.NewContext(),
		},
		{
			name: "ugly: nil-shaped empty document content still renders calmly",
			src:  `<div></div>`,
			want: html.El("div"),
			ctx:  html.NewContext(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertSameRender(t, tc.src, tc.bnd, tc.want, tc.ctx)
		})
	}
}

func TestParseLayout_Good(t *testing.T) {
	src := `<layout variant="HLCRF">
		<h><span>demo.brand</span></h>
		<l><ul><li>demo.menu.a</li></ul></l>
		<c><h1>demo.title</h1></c>
		<r><p>demo.side</p></r>
		<f><span>demo.footer</span></f>
	</layout>`

	want := html.NewLayout("HLCRF").
		H(html.El("span", html.Text("demo.brand"))).
		L(html.El("ul", html.El("li", html.Text("demo.menu.a")))).
		C(html.El("h1", html.Text("demo.title"))).
		R(html.El("p", html.Text("demo.side"))).
		F(html.El("span", html.Text("demo.footer")))

	layout, err := ParseLayout([]byte(src))
	require.NoError(t, err)
	require.NotNil(t, layout)

	ctx := html.NewContext()
	assert.Equal(t, want.Render(ctx), layout.Render(ctx))
	assert.Equal(t, want.RenderTerm(ctx), layout.RenderTerm(ctx))
	assert.Equal(t, want.RenderTerm(ctx, html.TermOptions{Width: 60}), layout.RenderTerm(ctx, html.TermOptions{Width: 60}),
		"narrow width stacks the same way on both trees")
}

func TestParseLayout_Bad(t *testing.T) {
	_, err := ParseLayout([]byte(`<div><p>not a layout</p></div>`))
	require.Error(t, err)
	var pe *ParseError
	require.ErrorAs(t, err, &pe)
	assert.Contains(t, pe.Msg, "must be <layout>")
}

func TestParse_Bad(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantMsg string
	}{
		{"bad: malformed xml reports a position", `<div><p>oops</div>`, "malformed XML"},
		{"bad: unclosed document reports unexpected end", `<div><p>oops`, "unexpected end of document"},
		{"bad: if without cond", `<if><p>x</p></if>`, "requires a cond attribute"},
		{"bad: unless without cond", `<unless><p>x</p></unless>`, "requires a cond attribute"},
		{"bad: switch without on", `<switch><case value="a"><p>x</p></case></switch>`, "requires an on attribute"},
		{"bad: case without value", `<switch on="k"><case><p>x</p></case></switch>`, "requires a value attribute"},
		{"bad: entitled without feature", `<entitled><p>x</p></entitled>`, "requires a feature attribute"},
		{"bad: each without items", `<each as="row"><p>x</p></each>`, "requires an items attribute"},
		{"bad: each without as", `<each items="repos"><p>x</p></each>`, "requires an as attribute"},
		{"bad: layout without variant", `<layout><c><p>x</p></c></layout>`, "requires a variant attribute"},
		{"bad: case outside switch", `<case value="a"><p>x</p></case>`, "only valid as a direct child of <switch>"},
		{"bad: slot outside layout", `<c><p>x</p></c>`, "only valid as a direct child of <layout>"},
		{"bad: variant outside responsive", `<variant name="a"><layout variant="C"><c><p>x</p></c></layout></variant>`, "only valid as a direct child of <responsive>"},
		{"bad: switch rejects a non-case child", `<switch on="k"><p>x</p></switch>`, "only accepts <case> children"},
		{"bad: layout rejects an unknown slot letter", `<layout variant="C"><x><p>x</p></x></layout>`, "only accepts <h> <l> <c> <r> <f> children"},
		{"bad: responsive rejects a non-variant child", `<responsive><p>x</p></responsive>`, "only accepts <variant> children"},
		{"bad: raw rejects element children", `<raw><p>x</p></raw>`, "cannot contain child elements"},
		{"bad: duplicate case value", `<switch on="k"><case value="a"><p>1</p></case><case value="a"><p>2</p></case></switch>`, "duplicate"},
		{"bad: variant requires exactly one layout child", `<responsive><variant name="a"><p>x</p></variant></responsive>`, "requires a <layout> child"},
		{"bad: variant rejects two layout children", `<responsive><variant name="a"><layout variant="C"><c><p>x</p></c></layout><layout variant="C"><c><p>y</p></c></layout></variant></responsive>`, "requires exactly one <layout> child"},
		{"bad: args on element content is rejected", `<p args="x"><b>hi</b></p>`, "not an element child"},
		{"bad: args on multi-run content is rejected", `<p args="x">Hello <b>world</b></p>`, "requires exactly one text child"},
		{"bad: unbound each reference outside any each", `<p>{{row.name}}</p>`, "unbound reference"},
		{"bad: unbound each reference after its scope closes", `<div><each items="a" as="row"><p>{{row.n}}</p></each><p>{{row.n}}</p></div>`, "unbound reference"},
		{"bad: unbound reference in args", `<p args="{{row.n}}">k</p>`, "unbound reference"},
		{"bad: trailing content after the root element", `<div><p>x</p></div><div>y</div>`, "unexpected content after root"},
		{"bad: text before the root element", `stray<div><p>x</p></div>`, "unexpected text before root"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.src))
			require.Error(t, err)
			var pe *ParseError
			require.ErrorAs(t, err, &pe)
			assert.Contains(t, pe.Msg, tc.wantMsg)
			assert.Greater(t, pe.Line, 0, "line is populated")
			assert.Greater(t, pe.Col, 0, "column is populated")
		})
	}
}

func TestParse_Bad_ErrorPositions(t *testing.T) {
	// Three lines: the offending <if> with no cond starts on line 2.
	src := "<div>\n<if><p>x</p></if>\n</div>"
	_, err := Parse([]byte(src))
	require.Error(t, err)
	var pe *ParseError
	require.ErrorAs(t, err, &pe)
	assert.Equal(t, 2, pe.Line, "error reported on the line the offending element starts")
}

func TestParse_Bad_LenientOnNearMissBindingTokens(t *testing.T) {
	// A near-miss like {{oops!}} is not a valid identifier path, so it is
	// left as literal text rather than rejected -- S:S8.4.
	got, err := Parse([]byte(`<p>See {{oops!}} for details</p>`))
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	assert.Contains(t, out, "oops!")
}

func TestParse_Ugly_EmptyDocument(t *testing.T) {
	_, err := Parse(nil)
	require.Error(t, err)
	var pe *ParseError
	require.ErrorAs(t, err, &pe)
	assert.Contains(t, pe.Msg, "unexpected end of document")
}

func TestParse_Ugly_UnwrapsToCoreErr(t *testing.T) {
	_, err := Parse([]byte(`<if><p>x</p></if>`))
	require.Error(t, err)
	var coreErr *core.Err
	assert.ErrorAs(t, err, &coreErr, "ParseError.Unwrap chains to a core.Err, matching the repo's error idiom")
}

func TestParseError_Error_Good(t *core.T) {
	err := &ParseError{Line: 5, Col: 12, Msg: "oh no"}
	core.AssertEqual(t, "ctml:5:12: oh no", err.Error())
}

func TestParseError_Error_Bad(t *core.T) {
	var err *ParseError
	core.AssertEqual(t, "", err.Error())
}

func TestParseError_Error_Ugly(t *core.T) {
	err := &ParseError{}
	core.AssertEqual(t, "ctml:0:0: ", err.Error())
}

func TestParseError_Unwrap_Good(t *core.T) {
	cause := core.NewError("boom")
	err := &ParseError{Cause: cause}
	core.AssertEqual(t, cause, err.Unwrap())
}

func TestParseError_Unwrap_Bad(t *core.T) {
	var err *ParseError
	core.AssertEqual(t, nil, err.Unwrap())
}

func TestParseError_Unwrap_Ugly(t *core.T) {
	err := &ParseError{}
	core.AssertEqual(t, nil, err.Unwrap())
}

func TestParse_Each_RowsCloseOverFreshPerItem(t *testing.T) {
	// Guards against the classic closure-capture bug: each rendered <li>
	// must see its own row, not the last row in the slice.
	src := `<ul><each items="repos" as="row"><li>{{row.name}}</li></each></ul>`
	bnd := Bindings{Sequences: map[string][]map[string]any{
		"repos": {{"name": "alpha"}, {"name": "beta"}, {"name": "gamma"}},
	}}
	got, err := Parse([]byte(src), bnd)
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	require.True(t, strings.Contains(out, "alpha") && strings.Contains(out, "beta") && strings.Contains(out, "gamma"))
	assert.Equal(t, 1, strings.Count(out, "alpha"))
	assert.Equal(t, 1, strings.Count(out, "beta"))
	assert.Equal(t, 1, strings.Count(out, "gamma"))
}
