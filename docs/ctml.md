---
title: CTML
description: Specification for .ctml, the HTML-flavoured markup that parses 1:1 onto go-html's node-tree builder API -- HTML render, terminal render, entitlements, and i18n all work unchanged on a parsed tree.
---

# CTML

`.ctml` is a markup surface for `go-html` node trees. A `.ctml` file parses into exactly the tree the Go builder API (`El`, `Text`, `If`, `Unless`, `Switch`, `Each`, `EachSeq`, `Entitled`, `Layout`, `Responsive`) would produce for the same page written by hand -- there is no separate interpreter, no second rendering path, and no markup-only behaviour. `HTML render` (`Render`), `terminal render` (`RenderTerm`), entitlement gating, and i18n all operate on the parsed tree exactly as they do on a hand-built one, because it is the same tree.

This document is the grammar. A reader who has never seen the parser source should be able to reimplement it from this page alone.

## 1. Design principles

1. **Closed vocabulary, no expression language.** `go-html`'s own architecture doc is explicit that the HTML layer is "not a general-purpose template engine -- no arbitrary expressions, no template inheritance, no macros" (RFC-CORE-006 S:S2). CTML inherits that constraint rather than relaxing it: every dynamic construct in this document resolves to a named lookup (a `Context.Data` key, an item field, a bindings entry), never a parsed expression, operator, or function call.
2. **The parser only ever calls exported constructors.** `go-html`'s terminal renderer (`term.go`) walks a node tree with a type switch; any concrete type it does not recognise falls through to a code path that re-wraps and re-visits the same node -- an infinite loop, not a graceful no-op (traced in `term.go`'s `resolve`/`inline` default cases). Go additionally treats an unexported interface method (`termExpandable.termNodes()`) as identified by *package*, not by spelling, so a type declared outside the `html` package can never satisfy it. Both facts together mean CTML must never introduce a custom `Node` implementation -- every parsed element becomes a real `El`/`Text`/`Raw`/`Verbatim`/`If`/`Unless`/`Switch`/`Each`/`EachSeq`/`Entitled`/`Layout`/`Responsive` value, so it is automatically safe under both renderers with no CTML-specific cases in either. (`Verbatim` (S:S6.5) is a real `html` node like `Raw`, added to package `html` and handled as a first-class case in the terminal type switch -- not a CTML-side `Node` type.)
3. **Attributes are static per constructed node.** `elNode.attrs` is a plain `map[string]string` fixed at construction time -- nothing in the Go API re-evaluates an attribute value per `Render(ctx)` call, and CTML does not invent that. It does allow the same closed-vocabulary `{{path}}` lookup its text runs use (S:S8.3) *inside* an attribute value: the token resolves when the node is **constructed** -- once at document scope, once per row when `<each>` builds that row's node -- yielding one fixed literal string per `elNode`. An `<each>` row can therefore carry a row-scoped `class` or `id` (S:S5) with no attribute ever re-evaluated at render time: the interpolation rides the same fresh-per-item closure `Each` already runs, exactly as a text-run bind does. This stays a named lookup, never an expression (S:S1). (Dynamic *content* is likewise possible -- see S:S7-8 -- because `Text`, `If`, `Unless`, `Switch` and `Each` each defer part of their work to render time or to a fresh per-item closure.)
4. **Text content doubles as its own i18n key and fallback.** `Text(key, args...)` looks the key up through the active `Translator`; every translator seen in this codebase (`go-i18n`, the WASM stub, the test `mapTranslator`) echoes the key verbatim when no catalogue entry exists. CTML exploits that convention: an element's text content *is* the i18n key, so a page reads as plain copy and is translation-ready with zero extra markup.
5. **Data-driven lists bind at parse time, not render time.** `Each(items []T, fn func(T) Node)` fixes `items` when the node is constructed; nothing in the exported API lets a list's *membership* vary per `Render(ctx)` call (only per-item rendering re-runs, because `fn` is invoked fresh each render). CTML is honest about this rather than papering over it: `<each>` sources its rows from a `Bindings` value supplied to `Parse`, not from `Context.Data`. If the underlying data changes, re-parse (or re-bind) -- exactly the constraint a hand-written `Each(repos, fn)` call already has.
6. **Universal item type.** Per-item data for `<each>` is `map[string]any` -- the same "bag of named values" shape `dappco.re/go`'s `Options`/`Option` already establishes as the universal input type elsewhere in CoreGO (AX-6). CTML does not invent a second row type.

## 2. Document grammar

A `.ctml` file is well-formed XML: exactly one root element, matched tags, quoted attribute values, self-closed void elements (`<br/>`, `<hr/>`, `<input/>`). The parser is `encoding/xml`'s `Decoder` (not a hand tokeniser) -- it gives honest `line:column` positions for free via `Decoder.InputPos()`, so every parse error in S:S9 is byte-accurate. Comments (`<!-- ... -->`) are permitted and ignored. There is no DOCTYPE, no processing instruction, no namespace support.

```
document   := element
element    := STag content ETag | EmptyElemTag
content    := (element | CharData)*
```

Two content rules apply uniformly, with no special-casing per tag:

- **A run of `CharData`** that is *entirely* whitespace is dropped -- this is how the indentation between pretty-printed sibling elements disappears. A run that has any non-whitespace content becomes a `Text(text)` node, in document order, where `text` is the run with its leading/trailing edge trimmed *only when that edge contains a newline* -- source-formatting indentation is stripped, but a plain inline space at the edge (as in `"Hello "` immediately before a following element) is significant text-flow spacing and survives untouched. The exceptions: inside `<raw>` (S:S6.2), all content is preserved verbatim and unescaped with no whitespace handling at all; and any `{{path}}` token within a run is a binding reference (S:S8.3), splitting the run into interleaved bind nodes and literal-text `Text` nodes -- so `○ {{tab.label}}` becomes `Text("○ ")` then the bind, and a run that is exactly one token becomes a lone bind. The edge rule is applied to the whole run before this split, not per segment, and empty text segments are dropped. Each bind resolves against the enclosing `<each>` row inside an each body, or against `Bindings.Values` at document scope outside one.
- **A child element** becomes its own node, in document order, interleaved with any `Text`/bind nodes from adjacent `CharData` runs.

This is why `<p>Hello <strong>world</strong>!</p>` needs no special mixed-content handling: it is `El("p", Text("Hello "), El("strong", Text("world")), Text("!"))` by the same two rules applied three times -- note the space kept on `"Hello "`.

## 3. Element to node mapping

| Element | Maps to | Notes |
|---|---|---|
| any tag not listed below | `El(tag, children...)` | every XML attribute becomes `Attr(node, key, value)`; void tags (`br`, `hr`, `img`, `input`, ...) must self-close |
| `<raw>...</raw>` | `Raw(content)` | content is the concatenation of every `CharData` token inside, unescaped and whitespace-preserved; element children are not permitted (trusted-content escape hatch, matching the Go doc comment on `Raw`) |
| `<verbatim value="key"/>` | `Verbatim(Bindings.Values[key])` | self-closing; content is the pre-styled terminal string at `Values[key]` (must be a present `string`, else a parse error -- S:S6.5); terminal render passes it through byte-for-byte, HTML render escapes it |
| `<if cond="key">child</if>` | `If(f, child)` | `f` tests `Context.Data[key]` for truthiness (S:S7.1); exactly one logical child (S:S2 content rules; multiple element children are wrapped, S:S6.1) |
| `<unless cond="key">child</unless>` | `Unless(f, child)` | mirrors `<if>` with the condition inverted at the `Unless` call, not by negating `f` |
| `<switch on="key"><case value="a">...</case>...</switch>` | `Switch(sel, cases)` | `sel` returns the string at `Context.Data[key]` (S:S7.2); each `<case>` becomes one `cases[value]` entry |
| `<case value="...">child</case>` | one `Switch` case | only valid as a direct child of `<switch>` |
| `<entitled feature="name">child</entitled>` | `Entitled(name, child)` | two-argument legacy form -- renders through `Context.Entitlements`, deny-by-default, matching the HTML renderer exactly |
| `<each items="name" as="row">body</each>` | `Each[map[string]any](rows, fn)` | `rows` comes from `Bindings.Sequences[name]` supplied to `Parse`, not from `Context`; see S:S8 |
| `<layout variant="HLCRF">slots</layout>` | `NewLayout(variant).H(...).L(...).C(...).R(...).F(...)` | slot children below; unknown variant letters are ignored, matching `NewLayout`'s own permissive behaviour |
| `<h>`, `<l>`, `<c>`, `<r>`, `<f>` | one `Layout` slot's children | valid only as direct children of `<layout>`; letters match `Layout.H/L/C/R/F` exactly (the same HLCRF vocabulary used throughout this codebase's `slotRegistry`, `blockID`, and variant strings) |
| `<responsive>variants</responsive>` | `NewResponsive().Add(name, layout, media)` | wraps one or more `<variant>` |
| `<variant name="desktop" media="...">layout</variant>` | one `Responsive.Add` entry | valid only as a direct child of `<responsive>`; must contain exactly one `<layout>` |

Sixteen reserved tag names in total: `if`, `unless`, `switch`, `case`, `entitled`, `each`, `raw`, `verbatim`, `layout`, `h`, `l`, `c`, `r`, `f`, `responsive`, `variant`. None collides with a real HTML element name, so every other tag -- `div`, `p`, `h1`, `table`, `button`, an author's own custom element name -- passes straight through to `El` unmodified. This is deliberate: CTML does not maintain an allow-list of "known HTML tags"; anything not reserved is structural, exactly like calling `El` in Go.

## 4. Why wrapper elements, not directive attributes

Two conventions are common in markup template languages: a dedicated wrapper element (`<if cond="x">...</if>`) or a directive attribute on an ordinary element (`<div ctml-if="x">...</div>`, in the style of Vue's `v-if`). CTML uses wrapper elements because that is what the Go API itself already does -- `ifNode{cond, node}`, `unlessNode{cond, node}`, `entitledNode{feature, node}`, `switchNode{selector, cases}` are all *wrapper types* around a plain node, not a modification of `elNode`. A wrapper element is a direct, no-desugaring transcription of that shape; a directive attribute would require the parser to synthesise a wrapper the Go API doesn't otherwise construct, and to reserve an attribute name on every element rather than a tag name on a closed set of elements.

## 5. Attribute conventions

Every attribute on a non-reserved element maps to `Attr(node, key, value)` verbatim -- `class`, `href`, `id`, `value`, `max`, `placeholder`, `aria-*`, `data-*`, all pass through unchanged, exactly as the terminal renderer already interprets them (`class` for theme tokens, `href` for OSC 8 hyperlinks, `value`/`max` for `<progress>`, and so on -- see `docs/architecture.md` and `term.go`).

**An attribute value may interpolate `{{path}}` tokens.** The value is otherwise a literal string, but the same `{{path}}` binding syntax a text run uses (S:S8.3) is recognised inside it and resolved against the active scope -- the enclosing `<each>` row first, then `Bindings.Values` (S:S8.5). `class="{{row.state}}"`, `class="nav-item {{row.state}}"`, and `id="row-{{row.id}}"` all resolve per row; literal text between tokens survives, a miss renders empty, and a `{{` that opens no valid path token stays literal (S:S8.4). Resolution is at **construction** time, not render time, so each constructed node still carries one fixed attribute string (S:S1.3) -- this is a named lookup, not an expression language, and it works in *any* attribute (there is no `class`/`id` allow-list, exactly as S:S3 keeps no known-HTML-tag allow-list). The resolved value is stringified for display with the same `stringOf` rule an `astBind` uses; it is never routed through the translator, because an attribute is a `class`/`id`/`href`, not an i18n key. This is the seam that lets a data-driven `<each>` carry selection styling -- an active row's `class` -- and a distinct per-row `id`, so each row records its own terminal box (S:S14) rather than colliding on one static id. The binding author owns `id` uniqueness across rows, exactly as the author of a hand-built tree does.

`id` is the one attribute CTML gives an additional meaning to, and only additively: it still passes through to `Attr` as a normal HTML attribute (so it appears in HTML output as `id="..."` as any reader would expect), and it is also the key the terminal box map (S:S12) records a block's rendered bounds under, when the element is block-level. Setting `id` is optional; nothing breaks if it is absent, and the mapping table above is unaffected by it.

Reserved attributes (`cond`, `on`, `value`, `feature`, `items`, `as`, `variant`, `name`, `media`) are consumed by their owning element and are not passed through to `Attr` -- `<if cond="isAdmin">` never emits a literal `cond="isAdmin"` HTML attribute, because `<if>` itself never becomes an `El`.

## 6. Text, i18n, and the `args` attribute

### 6.1 Bare text is the key

Per S:S2's content rules, an element's text content becomes `Text(trimmed)`. There is no separate "translation key" vs "display text" distinction in the markup -- the source string *is* the key, and (per the translator fallback convention, S:S1.4) it is also what renders when no catalogue entry exists yet. A page is legible and correct before a single string is added to a catalogue.

### 6.2 `<raw>` for trusted, unescaped, non-translated content

`<raw>$ core go qa</raw>` maps to `Raw("$ core go qa")` -- unescaped, exempt from i18n, exactly matching the Go escape hatch's own doc comment ("trusted content only"). Whitespace inside `<raw>` is preserved verbatim, which is why `<pre>` bodies should use it.

### 6.3 Multiple children under a single-child wrapper

`If`, `Unless`, `Entitled`, and each `Switch` case all wrap exactly one `Node`. When a `.ctml` author gives `<if>`/`<unless>`/`<entitled>`/`<case>` more than one logical child, the parser wraps them with the same technique available to any Go caller:

```go
func fragment(nodes []html.Node) html.Node {
	if len(nodes) == 1 {
		return nodes[0]
	}
	return html.Each(nodes, func(n html.Node) html.Node { return n })
}
```

`Each` with an identity function renders each item's own `Render` output back to back with no wrapping markup and no extra node type -- it is a plain concatenation, and (because it *is* the real, exported `Each`) it is automatically safe under `RenderTerm` too (`eachNode[T]` already implements the terminal walker's expansion interface for any `T`, including `T = html.Node`). A single child is returned unwrapped, so the common case never pays for this at all.

### 6.4 The `args` attribute

`Text(key, args...)` accepts variadic arguments that flow straight to `Translator.T(key, args...)` -- typically ICU-style positional substitution (`"{0} of {1}"`) or pluralisation counts. Any element whose text content is a *single* run (no element children) may carry an `args="..."` attribute: a comma-separated list of tokens, each either a literal string or a `{{path}}` binding reference (S:S8.3, valid only inside an `<each>` body). `args` on an element with mixed/element content is a parse error -- there is no well-defined single message to attach the arguments to.

```xml
<p args="{{row.count}}">queue.remaining</p>
```

### 6.5 `<verbatim>` for pre-styled terminal content

`<raw>` cannot carry pre-styled ANSI: escape/control bytes are invalid in XML markup, and the terminal renderer routes `Raw` through `StripTags` and whitespace normalisation (which would mangle escape sequences). `<verbatim>` is the escape hatch for content that is *already terminal-ready* -- most importantly Glamour-rendered markdown a caller wants to drop inside composed chrome.

Because the bytes cannot live in markup, the content is supplied through a binding, not as element children:

```xml
<verbatim value="rendered"/>
```

```go
ctml.Parse(src, ctml.Bindings{Values: map[string]any{
	"rendered": glamourOutput, // pre-styled ANSI string
}})
```

`value="key"` names a `Bindings.Values` entry that must be present and a `string`; an absent key or a non-string value is a position-accurate parse error (S:S9). The element is empty (self-closing, or open/close with no content) -- any child element or non-whitespace text is a parse error. It maps to the exported `Verbatim` node in package `html`:

- **Terminal render**: the content is emitted exactly as-is -- no `StripTags`, no whitespace normalisation, no width wrapping. The caller owns fitting the bytes to the target width.
- **HTML render**: the content is HTML-escaped as ordinary text -- a safe default, since raw ANSI/control bytes are meaningless (and unescaped markup would be unsafe) in an HTML sink. Use `<raw>` for trusted HTML.

`Verbatim` is a real `html` node, so the terminal renderer handles it as a first-class case in its type switch, never the default -- the same discipline S:S1.2 requires of every node type (a node the walker does not recognise re-resolves to itself and spins).

## 7. Conditions and selectors: `Context.Data` by key

`If`/`Unless` need a `func(*Context) bool`; `Switch` needs a `func(*Context) string`. Both are genuine per-render closures in the Go API, and both receive only `*Context` -- so both are the one place CTML dynamism can honestly reach *render-time* state, by naming a `Context.Data` key.

### 7.1 `cond` truthiness

`cond="key"` compiles to a closure that looks up `ctx.Data[key]` and applies a fixed, minimal coercion -- no operators, no comparisons:

| Value at the key | Truthy when |
|---|---|
| missing / `nil` | never (false) |
| `bool` | its own value |
| `string` | non-empty |
| any numeric type | non-zero |
| slice or map | non-empty (`len > 0`) |
| anything else | present (true) |

### 7.2 `switch on`

`on="key"` compiles to a closure returning `ctx.Data[key]` coerced to `string` (used as-is if already a string; any other type is treated as unmatched -- i.e. the same as an absent key -- rather than guessing a stringification). A `<case value="">` catches the unmatched/absent case explicitly, matching `Switch`'s own "no case matches -> empty string" behaviour when no such case is present.

## 8. `<each>`: data-driven lists

### 8.1 Why the item source is a `Bindings` value, not `Context.Data`

S:S1.5 already states the constraint; this section is the mechanism. `Parse` takes an optional `Bindings`:

```go
type Bindings struct {
	Sequences map[string][]map[string]any
	Values    map[string]any
}
```

`Sequences` supplies the row lists for `<each>` (this section); `Values` supplies the document-wide scalars a lone `{{path}}` outside any `<each>` resolves against (S:S8.5). Both are supplied at parse time and both are optional -- an absent name in either map is data absence, not a parse error.

`<each items="repos" as="row">` looks `"repos"` up in `Bindings.Sequences` **at parse time** and calls the real, exported generic constructor:

```go
html.Each(bindings.Sequences["repos"], func(row map[string]any) html.Node {
	// body, built once per Parse call, closes over `row`
})
```

Because `fn` is a closure stored on the returned `*eachNode[map[string]any]` and only *invoked* inside `Render`, per-item content (translations, entitlement checks inside the body) still re-evaluates correctly on every render -- only the row *count and membership* is fixed at parse time. An `items` name absent from `Bindings` renders as an empty list (zero rows), not an error -- a `.ctml` file is valid to parse standalone (for the round-trip/structure tests in S:S1) even before its data is wired up.

### 8.2 `as` scoping

`as="row"` names the current item for the body subtree. Nested `<each>` elements each introduce their own name; references are always fully qualified (`row.field`, `order.field`), so nesting is unambiguous without a shadowing rule to specify or misremember.

### 8.3 `{{path}}` binding references

A `{{path}}` token is recognised in three places: as one or more tokens within a text run (S:S2, interleaved with literal text), as one comma-separated token inside an `args=` attribute (S:S6.4), or interpolated within any other attribute value (S:S5, interleaved with literal text the same way a text run is). `path` is an identifier followed by zero or more dotted field steps (`greeting`, `row.name`, `user.address.city`), each step indexing one level into a `map[string]any` (or a nested map found there). There is no other operator -- no filters, no arithmetic, no comparisons.

Resolution is by scope, innermost first: if `path`'s root identifier names an enclosing `<each as="...">`, the token resolves against that row (the nearest such each wins, so nested loops stay unambiguous); otherwise it resolves against `Bindings.Values` at document scope (S:S8.5). Row scope always wins over `Values` -- a `Values` key never shadows a row name.

- As whole text-run content: `<li>{{row.name}}</li>` compiles to `Text(stringOf(resolve("row.name")))` -- resolved fresh on every invocation of the body closure, HTML-escaped (or term-styled) exactly like any other `Text` node, via the real constructor.
- As an `args` token: the resolved value (kept as its native Go type -- `string`, `int`, `float64`, `bool`) is passed straight through as one positional argument, so a translator can format it (pluralise, localise a number) rather than CTML stringifying it up front.

Every miss renders as the empty string -- an absent field, an absent `Values` key, or a nil map alike: a `{{path}}` that resolves to nothing is data absence, not a document defect, mirroring `Options.Get`'s own miss-is-empty convention elsewhere in CoreGO. There is therefore **no unbound-reference parse error**: because `Values` is a document-wide scope, every syntactically valid `{{path}}` has a resolution target, and supplying the data stays a parse-time-optional concern (S:S8.1) -- exactly as an absent `items` sequence renders as an empty list rather than failing to parse.

### 8.4 What `{{path}}` deliberately does not do

Mixed static-and-dynamic text in one run *is* supported: `Cost: {{row.amount}}` splits into `Text("Cost: ")` and the bind (S:S2), and `○ {{tab.label}}` into `Text("○ ")` and the bind. What a `{{path}}` still does not do is compute: it is a named lookup and nothing else -- no filters, no arithmetic, no comparisons, no function calls. For a *translated* message, still prefer `args` with a catalogue entry (`<span args="{{row.amount}}">cost.line</span>`, catalogue `"Cost: {0}"`) over interpolating mid-sentence, so word order stays translator-controlled rather than baked into the markup.

There is also no escape syntax for a literal `{{`. Because the vocabulary is closed, a *well-formed* `{{path}}` is always a lookup -- there is no way to write one that renders as literal braces. A `{{` that does not open a valid path token (`{{oops!}}`, `{{not a path}}`) is not a token at all and stays literal text; content that genuinely needs literal `{{ident}}` braces belongs in `<raw>`.

### 8.5 Scalar values: `{{path}}` outside an `<each>`

`Bindings.Values` is a flat `map[string]any` of document-wide scalars, the companion to `Sequences`. A lone dynamic value -- a title, a count, a username -- rides `Values` directly rather than having to be wrapped in a one-row `<each>` just to be referenced:

```xml
<h1>{{page.title}}</h1>
<p args="{{unread}}">inbox.count</p>
```

```go
ctml.Parse(src, ctml.Bindings{Values: map[string]any{
	"page":   map[string]any{"title": "Dashboard"},
	"unread": 4,
}})
```

A whole-run `{{path}}` (or an `args` token) whose root identifier matches no enclosing `<each>` resolves against `Values` with the same `lookupPath` semantics a row uses: `{{unread}}` is `Values["unread"]`, `{{page.title}}` indexes one level into the nested map, and a native Go value passed as an `args` token keeps its type for the translator. Inside an `<each>` body a row-prefixed path still resolves to the row (S:S8.3) -- `Values` does not shadow rows -- so the two scopes never collide by accident. Like `Sequences`, `Values` may be omitted when parsing a document standalone; every unmatched key simply renders empty.

## 9. Errors

`Parse` returns `(html.Node, error)`. Every error is a `*ctml.ParseError`:

```go
type ParseError struct {
	Line, Col int
	Msg       string
	Cause     error // wrapped via core.E; nil for structural (non-XML) errors
}
```

`Line`/`Col` come from `xml.Decoder.InputPos()` at the point the offending token was read (malformed XML) or the point the offending element/attribute was being interpreted (a reserved element used wrongly, `args` on mixed content, a `<verbatim>` whose `value` key is absent or non-string, and so on) -- both classes of error carry real, non-zero positions, not just a wrapped stdlib message.

## 10. Package API

```go
package ctml // dappco.re/go/html/ctml

type Bindings struct {
	Sequences map[string][]map[string]any
	Values    map[string]any
}

type ParseError struct {
	Line, Col int
	Msg       string
	Cause     error
}
func (e *ParseError) Error() string
func (e *ParseError) Unwrap() error

// Parse parses src into the node tree the Go builder API would produce for
// the same page. bindings is optional; an absent <each> source renders as
// an empty list rather than failing to parse.
func Parse(src []byte, bindings ...Bindings) (html.Node, error)

// ParseLayout is Parse specialised to documents whose root is <layout>,
// returning the concrete *html.Layout so callers can chain further
// builder calls (Responsive.Add, additional slot appends) or call
// Layout-specific methods (RenderTerm) directly.
func ParseLayout(src []byte, bindings ...Bindings) (*html.Layout, error)
```

## 11. Worked example: a static page

```xml
<layout variant="HCF">
  <h><span class="brand">demo.brand</span></h>
  <c>
    <h1>demo.title</h1>
    <p>demo.intro</p>
  </c>
  <f><span>demo.footer</span></f>
</layout>
```

Equivalent to:

```go
html.NewLayout("HCF").
	H(html.Attr(html.El("span", html.Text("demo.brand")), "class", "brand")).
	C(
		html.El("h1", html.Text("demo.title")),
		html.El("p", html.Text("demo.intro")),
	).
	F(html.El("span", html.Text("demo.footer")))
```

`Render(tree, ctx)` and `layout.RenderTerm(ctx, opts)` both work unchanged -- `ParseLayout` returns the concrete `*html.Layout`.

## 12. Worked example: a data-driven list

```xml
<div>
  <h2>demo.h.repos</h2>
  <ul>
    <each items="repos" as="row">
      <li args="{{row.name}},{{row.status}}">repo.row</li>
    </each>
  </ul>
</div>
```

Parsed with:

```go
bindings := ctml.Bindings{Sequences: map[string][]map[string]any{
	"repos": {
		{"name": "go-html", "status": "green"},
		{"name": "go-io", "status": "green"},
	},
}}
tree, err := ctml.Parse(src, bindings)
```

Equivalent to:

```go
html.El("div",
	html.El("h2", html.Text("demo.h.repos")),
	html.El("ul",
		html.Each([]map[string]any{
			{"name": "go-html", "status": "green"},
			{"name": "go-io", "status": "green"},
		}, func(row map[string]any) html.Node {
			return html.El("li", html.Text("repo.row", row["name"], row["status"]))
		}),
	),
)
```

## 13. Worked example: an HLCRF layout

```xml
<layout variant="HLCRF">
  <h><span>demo.brand</span></h>
  <l>
    <ul>
      <li>demo.menu.a</li>
      <li>demo.menu.b</li>
    </ul>
  </l>
  <c>
    <h1>demo.title</h1>
    <entitled feature="ops">
      <p class="ok">demo.ops</p>
    </entitled>
  </c>
  <r><p>demo.side.h</p></r>
  <f><span>demo.footer</span></f>
</layout>
```

This is structurally identical to `cmd/termdemo`'s hand-built page (`go run ./cmd/termdemo/ -w 110`) minus its table/progress-bar detail -- `NewLayout("HLCRF").H(...).L(...).C(...).R(...).F(...)`, rendering through both `Render` and `RenderTerm` unchanged, with `<entitled feature="ops">` gating exactly like `Entitled("ops", ...)` does today.

## 14. Box map and mouse resolution

`RenderTermBoxes(n, ctx, opts...)` (and `Layout.RenderTermBoxes`, `Responsive.RenderTermBoxes`) render exactly like `RenderTerm`, plus a `BoxMap` recording the terminal rectangle of every identified block:

```go
type Box struct {
	Row, Col      int
	Width, Height int
	Node          Node
}
type BoxMap map[string]Box
```

**Identification reuses two existing addressing schemes rather than inventing a third.** A Layout slot (H/L/C/R/F) is keyed by the same `blockID` string the HTML renderer already writes to its `data-block` attribute -- `"H"`, `"C"`, or `"C.1"` for a second top-level Content slot. An arbitrary element is keyed by its own `id` attribute, when set -- `id` already passes through to `Attr` as an ordinary HTML attribute (S:S5), so this is additive: setting it changes nothing about how an element renders, it just makes that element's rendered rectangle appear in the box map too. Because an `id` value may itself interpolate a `{{path}}` (S:S5), the rows of an `<each>` can each carry a distinct `id` -- `id="row-{{row.id}}"` -- so every row records its own box instead of colliding on one shared static id. The binding author owns id uniqueness across rows, exactly as a hand-built tree's author does; two blocks recorded under the same id resolve last-writer-wins, like any other map key. (A row's box records when the row element is block-level and reaches the block walker directly -- an `<each>` of `<div>`/`<section>` rows in a slot or at document top level; rows nested inside a `<ul>`/`<table>` render through the list/table path, which does not record per-item boxes.)

**Nested layouts get a disambiguating prefix.** The HTML renderer's `data-block` scheme threads a `Layout.path` string through clone-on-render specifically to keep nested slot IDs unique; the terminal renderer has no equivalent mechanism (it never clones a `Layout` with a path). Reusing bare `"H"`/`"C"`/etc for a nested layout's own slots would collide with the outer layout's. Rather than re-deriving the HTML side's path machinery, the box map disambiguates locally: the outermost layout in a render keeps the clean, `data-block`-matching keys; any layout nested inside a slot gets an `L<n>.` prefix (`"L1.H"`, `"L1.C"`), where `n` is a simple per-render counter. A click resolver only ever needs the prefix to tell two blocks apart, not to reconstruct a path.

**Box positions are computed alongside the existing width/column arithmetic, not by re-parsing the rendered string.** `term_layout.go`'s frame assembly already computes each band's line count and each column's width (`sidebarWidth`, `contentWidth`, `asideWidth`) to lay the page out; box recording reads those same values as they're computed. Recording is opt-in and zero-cost on the existing paths: the renderer carries an optional recorder field, `nil` on `RenderTerm`/`Layout.RenderTerm`/`Responsive.RenderTerm`, and every new code path is guarded on it -- confirmed by running the full existing terminal-render test suite unchanged after this addition (byte-identical output).

**Element-id boxes inside a themed slot (Header/Sidebar/Aside/Footer/Card) approximate their origin from that slot's own outer box.** A slot's border and padding shift its content in by a cell or two; the box map does not walk lipgloss's own layout maths to correct for that, so an id'd element a few rows into a bordered sidebar may be off by a cell versus its exact rendered position. The slot's own box is always exact. This is a deliberate scope cut, not an oversight -- pixel-perfect accounting would have to track every themed style's border/padding independently and re-derive it whenever `term_theme.go` changes.

### Mouse resolution

`go/teabox` resolves a coordinate against a `BoxMap` with no dependency on `bubbletea` -- `go.mod` does not carry it, and the resolver only needs two integers, not the `tea.MouseMsg` type itself:

```go
func Resolve(boxes html.BoxMap, x, y int) (Hit, bool)
func ResolveNode(boxes html.BoxMap, x, y int) (html.Node, bool)
```

A caller wires a real `tea.MouseMsg` at the one line that needs the type, in their own code: `teabox.Resolve(boxes, int(msg.X), int(msg.Y))`. When boxes overlap -- a nested layout's slot always renders inside its enclosing slot's rectangle -- the smallest-area match wins, so a click on a card inside a content slot resolves to the card, not the whole page; an exact-area tie resolves to the lexicographically smaller block ID, so the result is deterministic rather than dependent on Go's random map iteration order.

## 15. Terminal layout geometry

The terminal renderer (`term.go`, `term_layout.go`) composes an HLCRF `Layout` into a styled ANSI frame. A handful of its geometry decisions are load-bearing for a downstream TUI and are specified here so a caller can rely on them.

### 15.1 Fixed and content-sized slot widths

By default the wide (>= 80 column) middle band gives L a fixed 24-column budget and R a fixed 28, with C filling the remainder and a single-space gutter either side of C; below 80 columns the three stack vertically at full width. This is the right default for a page -- a sidebar wants a stable width -- but it cannot express a *content-packed strip*: a brand plus a few short cells that should ride L/C/R as one tight row rather than being flung to the far edges of a wide frame.

`TermOptions.FitSlots` is the opt-in for that strip. When set, each present L/C/R slot is measured to its own rendered content width (the widest line once padding and styling are discounted) and the slots are packed **edge-to-edge, left to right, with no inter-slot gutter**, on **one row whatever the terminal width** (the narrow-width stacking is bypassed -- a strip is meant to stay a strip). Slot chrome overhead is **measured from the active theme's slot style**, not assumed: an L/R box adds its `Sidebar`/`Aside` style's border and padding columns, and the C content adds its own structural `(0,1)` gutter (S:S15.2, which is not themed). For the default theme that measures to 4 columns for a bordered L/R box (rounded border 2 + `(0,1)` padding 2) and 2 for C. Measuring rather than assuming a fixed 4 is what lets a **stripped (borderless) or space-glyph slot theme record boxes that tile the visible glyphs exactly** -- a borderless theme whose chrome is 0, or a space-glyph left/right border with `(0,1)` padding whose chrome is still 4, both land their boxes on the rendered strip rather than a column or two wide of it. The recorded slot boxes (S:S14) tile the row at these true content-sized origins and widths -- `C.Col == L.Col + L.Width`, `R.Col == C.Col + C.Width` -- so mouse resolution stays exact whatever border and padding the active theme's slots carry.

`FitSlots` is `false` by default and changes nothing about any existing render. It is a terminal-render option, so it rides `TermOptions` (as `Width` and `Theme` do); it is not stored on the shared `Layout`, which is also the WASM-linked HTML compositor. The caller owns keeping content narrow enough for the target width -- fit slots size to their content and can, with wide content, exceed the frame -- the same ownership boundary as `id` uniqueness (S:S5).

### 15.2 The content gutter and band alignment

`renderTermContent` renders the C slot inside `(0,1)` padding -- one column of gutter to the left and right. This is deliberate alignment, not an artefact. The Header and Footer band styles (`term_theme.go`) already carry the same `(0,1)` padding, so H and F text sit one column in from the frame edge; the C gutter puts C content on that same column, so a page's header, body, and footer text line up vertically down the left margin. Removing the C padding would pull C content alone to column 0 -- one column left of every band, a worse misalignment, not a fix. The gutter has been part of the renderer since its first commit, alongside the band padding it matches; it is also the `+2` chrome overhead a content-sized C slot carries (S:S15.1).

### 15.3 Slot-junction trimming

Every slot and band renders its content through the same block assembler (`blocks`), which trims trailing blank lines before returning. A block-level node -- a `<p>`, an `<h2>`, a list -- emits a trailing blank line as its own paragraph spacing, but at a slot boundary that trailing gap is trimmed, so stacked bands (H over the middle band over F) and side-by-side slots abut cleanly instead of carrying a double blank line at the junction. This is layout-native geometry, not a defect: the frame owns the spacing *between* regions, the content inside a region does not reach across the boundary. To reintroduce vertical breathing room at a junction, give the relevant band style vertical padding in the theme -- e.g. `theme.Footer = theme.Footer.Padding(1, 1)` -- which the renderer honours as part of that band's own box (and its recorded height), rather than relying on an in-content trailing blank line that the junction trim will always remove.

### 15.4 Whitespace and the leading gutter

A leading space-only gutter on a **block-open line** -- indenting a paragraph's first line with blank columns alone -- is not expressible in terminal output, by design at two layers. At parse, a `CharData` run that is entirely whitespace is dropped (S:S2) so source indentation between elements disappears; a standalone gutter run before a child element goes with it. At terminal render, the block assembler flushes each inline paragraph through a single `strings.TrimSpace` over the *whole flattened run* -- the inline flush in `blocks`, and the `<p>`/`<address>` case in `blockEl`. `TrimSpace` bites only the run's two outer edges, so it removes the **block-open line's leading gutter** (and the final line's trailing gutter), which is what collapses a leading gutter that *did* survive parse by being glued to its text (`  text`). This is HTML-faithful: a block strips its own opening whitespace. A whitespace entity does not slip past either layer on that opening edge: `&#32;`, `&#160;` (NBSP) and `&#8199;` (figure space) all decode to whitespace *before* either sees them, and every Unicode space character is `unicode.IsSpace`, so `strings.TrimSpace` removes it the same as a plain space; `<raw>  text</raw>` preserves the spaces through parse, but the paragraph flush still trims them at the block-open edge.

**Interior lines are different -- the one asymmetry to know.** Because the flush trims only the run's outer edge, a line *after* an interior `<br>` is a continuation line whose **leading whitespace survives**: `<p>head<br>  tail</p>` keeps the two-space gutter on `  tail`, and `<p>  first<br>  second</p>` collapses the first line's gutter but keeps the second's. "Nothing survives" holds for the block-open line only; a post-`<br>` continuation gutter rides through untouched. So the blessed **marker-glyph** idiom -- a non-space leading character the flush keeps, `○ ` or `· ` before the text, exactly the tab-strip shape S:S8.4 already uses -- is needed **only at block-open rows** (a paragraph's first line, a band-opening row); a continuation line after `<br>` already carries its plain-space gutter and needs no marker, which is why a caller can mark just the opening rows and pass exact bytes elsewhere. The other route to a visible gutter, at any row, is **structural padding**: the C content gutter (S:S15.2), or a themed slot's own padding (`Sidebar`, `Card`). Reach for a marker glyph on the opening row, a plain gutter (or marker) on continuation rows, or a padded region -- never a run of spaces expecting the block-open edge to keep it.

### 15.5 Band and slot inner content width

A downstream doing row-budget maths -- how many columns a region gives before text wraps, and so how many rows a block of copy takes before a frame has to split -- needs the inner content width each region renders into, not just its outer width. Every region reserves its own chrome from the width it is handed, so a band of outer width `W` does **not** render content at `W`:

| Region | Inner content width | Where the chrome comes from |
|---|---|---|
| H, F bands | `W - 2` | the band's `(0,1)` horizontal padding gutter (S:S15.2) -- one column each side |
| C content | `W - 2` | its structural `(0,1)` gutter (S:S15.2), the same alignment gutter, not themed |
| L, R boxed slots | `W - 4` | the default rounded border (2) + `(0,1)` padding (2) |

The general rule is one line: **inner content width = region width − the region's horizontal chrome (border + padding)** -- the same `termChrome` FitSlots measures for L/R (S:S15.1). For the default theme that is `2` for the full-width bands (H/C/F -- padding only, no left/right border) and `4` for the bordered L/R boxes. Two consequences worth stating outright:

- The C gutter is **structural and fixed at 2** whatever the `TermTheme` is, because C's `(0,1)` alignment gutter is not a theme field (S:S15.2). The `- 4` for L/R is **theme-dependent**: a theme whose `Sidebar`/`Aside` style drops the border, or widens the padding, changes their inner width by exactly that much -- and FitSlots already tiles their boxes to that measured chrome (S:S15.1). The H/F bands render at a fixed `W - 2` matching the default band gutter; a theme that re-pads a band away from `(0,1)` owns keeping its content-width expectation in step, the same caller-owns-content boundary as FitSlots.
- The `W` a slot is handed is itself the layout's arithmetic, not the terminal width: at `>= 80` columns L is a fixed `24` and R a fixed `28` (so inner `20` and `24`), C fills the remainder; below 80 the slots stack at full width; under FitSlots each slot is content-sized (S:S15.1). Compose the two -- outer width from the layout, minus the chrome above -- to get the wrap width for a region.

No helper is exported for the subtraction: a region's chrome is a property of the `TermTheme` the caller already holds, the default contract above (`- 2` for bands, `- 4` for L/R) is stable for the shipped theme, and a single documented rule is a steadier thing to depend on than one more one-line accessor (S:S1 keeps the surface closed).

## 16. CoreCommand-derived default TUI (exploratory)

`CoreCommand` (`dappco.re/go`, the `core` module already in `go.mod`) is a declarative command tree: `Command{Name, Description, Path, Action CommandAction, Managed string, Flags Options, Hidden bool}`, `CommandAction = func(Options) Result`, registered and fetched via `(*Core).Command(path string, command ...Command) Result` (zero variadic args = lookup), listed flat via `(*Core).Commands() []string` in registration order. There is no separate subcommand-tree type: the flat, path-keyed registry *is* the tree (`"deploy/to/homelab"`), and registering a nested path auto-creates placeholder ancestor entries so the tree is always walkable from any leaf.

This splits the brief's ask -- "flags -> input chips, subcommands -> a list" -- into one half that is well-founded on real, exercised CoreCommand surface, and one half that is not. Both are covered below; only the first shipped as code this pass.

### Subcommands -> a list (shipped)

`ctml.SubcommandList(c *core.Core, root string, paths []string) html.Node` generates a `<ul>` of a path's *direct* children (`root=""` for the top level, `root="deploy"` for `deploy/*`), skipping `Hidden` commands, with each item's label being that command's `I18nKey()` -- `"deploy/to/homelab"` with no `Description` becomes the key `"cmd.deploy.to.homelab.description"`, exactly the derivation `Command.I18nKey()` itself already documents. This works because every piece it depends on is real, tested, load-bearing CoreCommand surface (`command_example_test.go` in `dappco.re/go` exercises `I18nKey`, `Core.Command` registration and lookup, and `Managed`/`Hidden` directly) -- the generator adds no new concept, it just walks a path prefix and renders what `I18nKey()` already promises, the same i18n-key-as-text convention S:S6.1 uses everywhere else in this document.

### Flags -> input chips (design question, not code)

`Command.Flags` is typed `Options` -- a slice of `Option{Key string, Value any}`, nothing else. Searching `dappco.re/go` v0.12.0 end to end: no example sets `Command.Flags`, no test reads it, and no internal consumer (`Cli`, `Command.Run`) touches it either -- it is declared, documented as "declared flags" in a doc comment, and otherwise unexercised. The CLI argument parser that *does* exist (`ParseFlag`, `IsFlag` in `utils.go`) works directly on `os.Args`-shaped strings (`"--port=8080"` -> `("port", "8080", true)`) and has no visible connection to `Command.Flags` at all.

A flags-to-chips generator needs, per flag, at minimum: a display label, an input kind (checkbox/number/text/select), and ideally whether it's required and what a select's choices are. `Option{Key, Value}` supplies exactly two of those by inference -- `Key` as a fallback label (no separate label field exists, unlike `Command`'s real `I18nKey()`), and `Value`'s runtime type as a weak signal for input kind (`bool` -> checkbox, numeric -> number, else -> text; a `reflect.TypeOf(opt.Value)` switch, mirroring the coercion style already established in S:S7.1's `cond` truthiness). There is no source for required/optional, help text, or enumerated choices, and `Value` has to double as both "the declared default" and "the currently-set value" since the type carries only one slot for it.

Building a generator against a field with zero real-world usage to learn from would mean guessing at a convention `dappco.re/go` itself hasn't settled on -- and `dappco.re/go` is explicitly the one repo in this stack that does not take unattended changes; extending `Option`/`Command`'s own shape is a decision for that repo's maintainer, not something to freelance from a `go-html`-scoped lane. Three shapes the decision could reasonably land on, for whoever makes that call:

1. **Reuse `Command`'s own convention on `Option`.** Give `Option` (or a new `Flag` type used only in `Command.Flags`) a `Description` field with the same derive-from-key-if-empty behaviour `Command.I18nKey()` already has, so a flag named `"port"` on command path `"serve"` defaults to key `"cmd.serve.flag.port.description"`. Minimal type growth, consistent with the one labelling convention this whole codebase already uses.
2. **A small closed `Kind` enum on `Option`** (`KindText`, `KindNumber`, `KindBool`, `KindSelect`), replacing type-sniffing `Value` with an explicit declaration -- steadier than `reflect.TypeOf` (a default of `0` is ambiguous between "int flag, zero default" and "unset"), at the cost of every caller that already builds `Option{Key, Value}` literals now needing to also set `Kind`.
3. **Leave `Option` alone; add an opt-in side table** (`map[string]FlagMeta` keyed by flag name, supplied to the generator directly, independent of `Command.Flags`) -- zero changes to `dappco.re/go`, but the metadata lives apart from the command registration it describes, which is exactly the kind of drift-prone split CoreGO's SPOR principle warns against elsewhere in this stack.

Whichever shape is chosen, the generator side is small: once a flag reliably carries a label and an input kind, `ctml.SubcommandList`'s pattern (walk the declarative structure, emit `El`/`Text`/`Attr` calls, no new node types) extends directly -- the open question is entirely in `dappco.re/go`'s `Option` shape, not in anything `go-html` owns.
