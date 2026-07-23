// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"strings"
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/html/engine/html"
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
			// A run with multiple {{path}} tokens splits into interleaved
			// bind and literal-text nodes (S:S2, S:S8.4): the "/" between the
			// two tokens becomes its own Text node, each token its own bind.
			name: "good: multiple bind tokens in one run interpolate with literal text between",
			src:  `<each items="repos" as="row"><li>{{row.name}}/{{row.status}}</li></each>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"repos": {{"name": "go-html", "status": "green"}},
			}}},
			want: html.Each([]map[string]any{{"name": "go-html", "status": "green"}}, func(row map[string]any) html.Node {
				return html.El("li", html.Text(row["name"].(string)), html.Text("/"), html.Text(row["status"].(string)))
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
			// Nested <each> scopes resolve independently: {{g.name}} binds to
			// the outer row and {{r.n}} to the inner, each its own text run.
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

func TestParse_Values_Good(t *testing.T) {
	tests := []struct {
		name string
		src  string
		bnd  []Bindings
		want html.Node
	}{
		{
			// A lone {{path}} outside any <each> resolves against
			// Bindings.Values -- no one-row <each> needed to carry it.
			name: "good: whole-run value renders",
			src:  `<p>{{greeting}}</p>`,
			bnd:  []Bindings{{Values: map[string]any{"greeting": "Hello"}}},
			want: html.El("p", html.Text("Hello")),
		},
		{
			// Dotted paths index one level into a nested Values map, exactly
			// like a row field path (lookupPath semantics).
			name: "good: nested value path resolves through a nested map",
			src:  `<p>{{user.name}}</p>`,
			bnd:  []Bindings{{Values: map[string]any{"user": map[string]any{"name": "Ada"}}}},
			want: html.El("p", html.Text("Ada")),
		},
		{
			// The {{path}} grammar is uniform across text runs and args
			// tokens, so an args token outside an <each> resolves Values too.
			name: "good: args token outside each resolves against Values",
			src:  `<p args="{{count}}">queue.remaining</p>`,
			bnd:  []Bindings{{Values: map[string]any{"count": 3}}},
			want: html.El("p", html.Text("queue.remaining", 3)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertSameRender(t, tc.src, tc.bnd, tc.want, html.NewContext())
		})
	}
}

func TestParse_Values_Bad_MissingKeyRendersEmpty(t *testing.T) {
	// A Values miss is data absence, not a document defect: it parses and
	// renders as empty text, matching the row field-miss behaviour (S:S8.3).
	got, err := Parse([]byte(`<p>{{absent}}</p>`), Bindings{Values: map[string]any{"present": "x"}})
	require.NoError(t, err, "a Values miss parses -- absence is not a parse error")
	assert.Equal(t, "<p></p>", html.Render(got, html.NewContext()), "a missing Values key renders as empty text")
}

func TestParse_MixedInterpolation(t *testing.T) {
	tests := []struct {
		name string
		src  string
		bnd  []Bindings
		want html.Node
	}{
		{
			// The tab-strip friction: a marker glued to a bind in one run.
			// The run splits into Text("○ ") + bind, rather than staying a
			// single literal key. The bind resolves against Values here.
			name: "good: text before a bind splits (tab-strip shape)",
			src:  `<span>○ {{tab.label}}</span>`,
			bnd:  []Bindings{{Values: map[string]any{"tab": map[string]any{"label": "Editor"}}}},
			want: html.El("span", html.Text("○ "), html.Text("Editor")),
		},
		{
			// Text before, between and after several tokens all survive as
			// their own Text nodes; every token is resolved against the row.
			name: "good: text around multiple bind tokens all survives",
			src:  `<each items="rows" as="row"><li>id-{{row.id}}: {{row.name}}!</li></each>`,
			bnd: []Bindings{{Sequences: map[string][]map[string]any{
				"rows": {{"id": "7", "name": "go-html"}},
			}}},
			want: html.Each([]map[string]any{{"id": "7", "name": "go-html"}}, func(row map[string]any) html.Node {
				return html.El("li",
					html.Text("id-"), html.Text(row["id"].(string)),
					html.Text(": "), html.Text(row["name"].(string)), html.Text("!"))
			}),
		},
		{
			// A "{{" that opens no valid path token is literal text: no
			// escape is invented (S:S8.4), the closed vocabulary means a
			// well-formed {{path}} is always a lookup and nothing else is.
			name: "bad: an invalid {{ token }} stays literal, no escape invented",
			src:  `<p>a {{not a path}} b</p>`,
			want: html.El("p", html.Text("a {{not a path}} b")),
		},
		{
			// A valid bind and an invalid brace-run coexist in one run: the
			// valid token interpolates, the invalid one stays literal text.
			name: "ugly: a valid bind and an invalid brace-run coexist in one run",
			src:  `<p>{{greeting}} {{nope!}}</p>`,
			bnd:  []Bindings{{Values: map[string]any{"greeting": "Hi"}}},
			want: html.El("p", html.Text("Hi"), html.Text(" {{nope!}}")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertSameRender(t, tc.src, tc.bnd, tc.want, html.NewContext())
		})
	}
}

func TestParse_Each_RowScopedAttributes(t *testing.T) {
	// Friction 2+8: a picker-style <each> where the active row carries a
	// different class AND every row records its own terminal box. class/id
	// values interpolate {{path}} per row (S:S5): a pure-bind class="{{row.cls}}"
	// and a mixed id="row-{{row.id}}" both resolve against the row, so the
	// downstream no longer needs a three-sequence rowsBefore/rowsActive/rowsAfter
	// split, and each row's distinct id lets its rendered rectangle record in
	// the box map (S:S14) instead of three rows colliding on one static id.
	src := `<layout variant="C"><c><each items="rows" as="row">` +
		`<div id="row-{{row.id}}" class="{{row.cls}}">{{row.label}}</div>` +
		`</each></c></layout>`
	bnd := Bindings{Sequences: map[string][]map[string]any{
		"rows": {
			{"id": "0", "cls": "nav-item", "label": "Alpha"},
			{"id": "1", "cls": "nav-item active", "label": "Beta"},
			{"id": "2", "cls": "nav-item", "label": "Gamma"},
		},
	}}
	layout, err := ParseLayout([]byte(src), bnd)
	require.NoError(t, err)

	// HTML render: the active row alone carries the active class, and every
	// row's id interpolated its own row value.
	out := html.Render(layout, html.NewContext())
	assert.Equal(t, 1, strings.Count(out, `class="nav-item active"`), "only the active row is styled active")
	assert.Equal(t, 2, strings.Count(out, `class="nav-item"`), "the two inactive rows keep the base class")
	for _, id := range []string{`id="row-0"`, `id="row-1"`, `id="row-2"`} {
		assert.Contains(t, out, id, "each row interpolated its own id")
	}

	// Terminal box map: each row records its own box under its distinct id,
	// stacked top to bottom -- no collision on a shared static id.
	_, boxes := layout.RenderTermBoxes(html.NewContext(), html.TermOptions{Width: 40})
	require.Contains(t, boxes, "row-0")
	require.Contains(t, boxes, "row-1")
	require.Contains(t, boxes, "row-2")
	assert.Less(t, boxes["row-0"].Row, boxes["row-1"].Row, "rows stack top to bottom in the box map")
	assert.Less(t, boxes["row-1"].Row, boxes["row-2"].Row)
	for _, id := range []string{"row-0", "row-1", "row-2"} {
		assert.Greater(t, boxes[id].Height, 0, "%s recorded a positive height", id)
	}
}

func TestParse_AttrInterpolation_StaticUnchanged(t *testing.T) {
	// A static attribute (no {{path}}) keeps the literal fast path: parsing it
	// is byte-identical to the hand-built Attr chain, so the interpolation seam
	// adds nothing to the common case (S:S5).
	src := `<a href="/docs" class="link" id="docs-link">docs.label</a>`
	want := html.Attr(html.Attr(html.Attr(html.El("a", html.Text("docs.label")), "href", "/docs"), "class", "link"), "id", "docs-link")
	assertSameRender(t, src, nil, want, html.NewContext())
}

func TestParse_AttrInterpolation_Values(t *testing.T) {
	// Outside an <each> an interpolated attribute resolves against Values, the
	// same document scope a lone {{path}} text token uses (S:S8.5); a missing
	// key renders empty, and a "{{" opening no valid token stays literal.
	tests := []struct {
		name string
		src  string
		bnd  []Bindings
		want html.Node
	}{
		{
			name: "good: whole-value bind resolves against Values",
			src:  `<a href="{{link}}">go</a>`,
			bnd:  []Bindings{{Values: map[string]any{"link": "/repos"}}},
			want: html.Attr(html.El("a", html.Text("go")), "href", "/repos"),
		},
		{
			name: "good: mixed literal-and-bind attribute value",
			src:  `<div class="card {{tone}}">x</div>`,
			bnd:  []Bindings{{Values: map[string]any{"tone": "ok"}}},
			want: html.Attr(html.El("div", html.Text("x")), "class", "card ok"),
		},
		{
			name: "bad: a missing Values key renders the bind as empty",
			src:  `<div class="card {{absent}}">x</div>`,
			bnd:  []Bindings{{Values: map[string]any{"present": "x"}}},
			want: html.Attr(html.El("div", html.Text("x")), "class", "card "),
		},
		{
			name: "ugly: a near-miss brace token stays literal in the attribute",
			src:  `<div data-tpl="{{not a path}}">x</div>`,
			want: html.Attr(html.El("div", html.Text("x")), "data-tpl", "{{not a path}}"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertSameRender(t, tc.src, tc.bnd, tc.want, html.NewContext())
		})
	}
}

func TestParse_Values_Ugly_RowWinsOverValuesInsideEach(t *testing.T) {
	// Inside an <each> body a row-prefixed {{path}} resolves to the row even
	// when Values also carries that root name -- Values does not shadow rows.
	src := `<each items="rows" as="row"><li>{{row.name}}</li></each>`
	bnd := Bindings{
		Sequences: map[string][]map[string]any{"rows": {{"name": "from-row"}}},
		Values:    map[string]any{"row": map[string]any{"name": "from-values"}},
	}
	got, err := Parse([]byte(src), bnd)
	require.NoError(t, err)
	out := html.Render(got, html.NewContext())
	assert.Contains(t, out, "from-row")
	assert.NotContains(t, out, "from-values")
}

func TestParse_Verbatim_Good(t *testing.T) {
	// <verbatim value="key"/> wires its content from Bindings.Values[key] at
	// parse time; the parsed tree renders identically to a hand-built
	// html.Verbatim under both renderers (HTML escapes, term passes through).
	ansi := "\x1b[1mpre-styled\x1b[0m <not-a-tag>"
	src := `<div><verbatim value="banner"/></div>`
	bnd := []Bindings{{Values: map[string]any{"banner": ansi}}}
	want := html.El("div", html.Verbatim(ansi))
	got := assertSameRender(t, src, bnd, want, html.NewContext())
	// And the term bytes reach output untouched via the ctml path.
	assert.Contains(t, html.RenderTerm(got, html.NewContext()), ansi)
}

func TestParse_Verbatim_Bad(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		bnd     Bindings
		wantMsg string
	}{
		{"bad: missing value attribute", `<verbatim/>`, Bindings{Values: map[string]any{"k": "x"}}, "requires a value attribute"},
		{"bad: key absent from Values", `<verbatim value="missing"/>`, Bindings{Values: map[string]any{"k": "x"}}, "no such key in Bindings.Values"},
		{"bad: value is not a string", `<verbatim value="n"/>`, Bindings{Values: map[string]any{"n": 42}}, "is not a string"},
		{"bad: child text content is rejected", `<verbatim value="k">oops</verbatim>`, Bindings{Values: map[string]any{"k": "x"}}, "cannot contain text content"},
		{"bad: child element is rejected", `<verbatim value="k"><b/></verbatim>`, Bindings{Values: map[string]any{"k": "x"}}, "cannot contain child elements"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.src), tc.bnd)
			require.Error(t, err)
			var pe *ParseError
			require.ErrorAs(t, err, &pe)
			assert.Contains(t, pe.Msg, tc.wantMsg)
			assert.Greater(t, pe.Line, 0, "line is populated")
			assert.Greater(t, pe.Col, 0, "column is populated")
		})
	}
}

func TestParse_Verbatim_RowScoped_Good(t *testing.T) {
	// Round 4: inside an <each>, <verbatim value="{{turn.body}}"/> binds per row at
	// materialise time -- mirroring how id="row-{{row.id}}" resolves per row (S:S5)
	// -- so a chat transcript, each turn's pre-styled ANSI as a row, is expressible.
	// Before this the parse-time closure returned one fixed value for every row.
	turn0 := "\x1b[1mAda\x1b[0m: hello"
	turn1 := "\x1b[1mGrace\x1b[0m: hi there"
	src := `<div><each items="turns" as="turn"><verbatim value="{{turn.body}}"/></each></div>`
	bnd := []Bindings{{Sequences: map[string][]map[string]any{
		"turns": {{"body": turn0}, {"body": turn1}},
	}}}
	want := html.El("div", html.Each(
		[]map[string]any{{"body": turn0}, {"body": turn1}},
		func(row map[string]any) html.Node { return html.Verbatim(row["body"].(string)) },
	))
	got := assertSameRender(t, src, bnd, want, html.NewContext())

	// Each row reached output with its OWN pre-styled bytes, byte-for-byte.
	out := html.RenderTerm(got, html.NewContext())
	assert.Contains(t, out, turn0)
	assert.Contains(t, out, turn1)
}

func TestParse_Verbatim_BracedBinding(t *testing.T) {
	// A whole {{path}} value defers to materialise time (S:S6.5): it resolves
	// against the enclosing <each> row or Values, and -- unlike a plain key, which
	// errors at parse (TestParse_Verbatim_Bad) -- a miss renders an empty Verbatim,
	// the same miss-is-empty rule every other {{path}} follows (S:S8.3), because the
	// target may only exist at bind time.
	ansi := "\x1b[1mbanner\x1b[0m"
	tests := []struct {
		name string
		src  string
		bnd  []Bindings
		want html.Node
	}{
		{
			name: "good: document-scope braced value resolves against Values",
			src:  `<verbatim value="{{banner}}"/>`,
			bnd:  []Bindings{{Values: map[string]any{"banner": ansi}}},
			want: html.Verbatim(ansi),
		},
		{
			name: "good: a missing braced key renders empty, not a parse error",
			src:  `<verbatim value="{{absent}}"/>`,
			bnd:  []Bindings{{Values: map[string]any{"present": ansi}}},
			want: html.Verbatim(""),
		},
		{
			name: "good: a missing row field renders an empty verbatim",
			src:  `<each items="turns" as="turn"><verbatim value="{{turn.body}}"/></each>`,
			bnd:  []Bindings{{Sequences: map[string][]map[string]any{"turns": {{"other": "x"}}}}},
			want: html.Each([]map[string]any{{"other": "x"}}, func(_ map[string]any) html.Node { return html.Verbatim("") }),
		},
		{
			name: "ugly: a nested-map value renders empty (stringOf has no representation)",
			src:  `<verbatim value="{{obj}}"/>`,
			bnd:  []Bindings{{Values: map[string]any{"obj": map[string]any{"k": "v"}}}},
			want: html.Verbatim(""),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertSameRender(t, tc.src, tc.bnd, tc.want, html.NewContext())
		})
	}
}
