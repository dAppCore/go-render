---
title: CTML
description: Specification for .ctml, the HTML-flavoured markup that parses 1:1 onto go-html's node-tree builder API -- HTML render, terminal render, entitlements, and i18n all work unchanged on a parsed tree.
---

# CTML

`.ctml` is a markup surface for `go-html` node trees. A `.ctml` file parses into exactly the tree the Go builder API (`El`, `Text`, `If`, `Unless`, `Switch`, `Each`, `EachSeq`, `Entitled`, `Layout`, `Responsive`) would produce for the same page written by hand -- there is no separate interpreter, no second rendering path, and no markup-only behaviour. `HTML render` (`Render`), `terminal render` (`RenderTerm`), entitlement gating, and i18n all operate on the parsed tree exactly as they do on a hand-built one, because it is the same tree.

This document is the grammar. A reader who has never seen the parser source should be able to reimplement it from this page alone.

## 1. Design principles

1. **Closed vocabulary, no expression language.** `go-html`'s own architecture doc is explicit that the HTML layer is "not a general-purpose template engine -- no arbitrary expressions, no template inheritance, no macros" (RFC-CORE-006 S:S2). CTML inherits that constraint rather than relaxing it: every dynamic construct in this document resolves to a named lookup (a `Context.Data` key, an item field, a bindings entry), never a parsed expression, operator, or function call.
2. **The parser only ever calls exported constructors.** `go-html`'s terminal renderer (`term.go`) walks a node tree with a type switch; any concrete type it does not recognise falls through to a code path that re-wraps and re-visits the same node -- an infinite loop, not a graceful no-op (traced in `term.go`'s `resolve`/`inline` default cases). Go additionally treats an unexported interface method (`termExpandable.termNodes()`) as identified by *package*, not by spelling, so a type declared outside the `html` package can never satisfy it. Both facts together mean CTML must never introduce a custom `Node` implementation -- every parsed element becomes a real `El`/`Text`/`Raw`/`If`/`Unless`/`Switch`/`Each`/`EachSeq`/`Entitled`/`Layout`/`Responsive` value, so it is automatically safe under both renderers with no CTML-specific cases in either.
3. **Attributes are static.** `elNode.attrs` is a plain `map[string]string` fixed at construction time -- the Go API itself has no mechanism to re-evaluate an attribute value per render. CTML does not invent one: attribute values in a `.ctml` file are literal strings, full stop. (Dynamic *content* is still possible -- see S:S7-8 -- because `Text`, `If`, `Unless`, `Switch` and `Each` each defer part of their work to render time or to a fresh per-item closure; attributes simply are not one of those seams.)
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

- **A run of `CharData`** (trimmed; dropped entirely if it trims to empty -- this is how indentation whitespace between sibling elements disappears) becomes a `Text(trimmed)` node, in document order, *except* inside `<raw>` (S:S6.2) where all content is preserved verbatim and unescaped, and inside `<each>` bodies where an entire run consisting of exactly one `{{path}}` token is a binding reference (S:S8.3), not a literal key.
- **A child element** becomes its own node, in document order, interleaved with any `Text` nodes from adjacent `CharData` runs.

This is why `<p>Hello <strong>world</strong>!</p>` needs no special mixed-content handling: it is `El("p", Text("Hello"), El("strong", Text("world")), Text("!"))` by the same two rules applied three times.

## 3. Element to node mapping

| Element | Maps to | Notes |
|---|---|---|
| any tag not listed below | `El(tag, children...)` | every XML attribute becomes `Attr(node, key, value)`; void tags (`br`, `hr`, `img`, `input`, ...) must self-close |
| `<raw>...</raw>` | `Raw(content)` | content is the concatenation of every `CharData` token inside, unescaped and whitespace-preserved; element children are not permitted (trusted-content escape hatch, matching the Go doc comment on `Raw`) |
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

Fifteen reserved tag names in total: `if`, `unless`, `switch`, `case`, `entitled`, `each`, `raw`, `layout`, `h`, `l`, `c`, `r`, `f`, `responsive`, `variant`. None collides with a real HTML element name, so every other tag -- `div`, `p`, `h1`, `table`, `button`, an author's own custom element name -- passes straight through to `El` unmodified. This is deliberate: CTML does not maintain an allow-list of "known HTML tags"; anything not reserved is structural, exactly like calling `El` in Go.

## 4. Why wrapper elements, not directive attributes

Two conventions are common in markup template languages: a dedicated wrapper element (`<if cond="x">...</if>`) or a directive attribute on an ordinary element (`<div ctml-if="x">...</div>`, in the style of Vue's `v-if`). CTML uses wrapper elements because that is what the Go API itself already does -- `ifNode{cond, node}`, `unlessNode{cond, node}`, `entitledNode{feature, node}`, `switchNode{selector, cases}` are all *wrapper types* around a plain node, not a modification of `elNode`. A wrapper element is a direct, no-desugaring transcription of that shape; a directive attribute would require the parser to synthesise a wrapper the Go API doesn't otherwise construct, and to reserve an attribute name on every element rather than a tag name on a closed set of elements.

## 5. Attribute conventions

Every attribute on a non-reserved element maps to `Attr(node, key, value)` verbatim -- `class`, `href`, `id`, `value`, `max`, `placeholder`, `aria-*`, `data-*`, all pass through unchanged, exactly as the terminal renderer already interprets them (`class` for theme tokens, `href` for OSC 8 hyperlinks, `value`/`max` for `<progress>`, and so on -- see `docs/architecture.md` and `term.go`).

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
}
```

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

A `{{path}}` token is recognised in exactly two places: as the entire trimmed content of a text run inside an `<each as="name">` body, or as one comma-separated token inside an `args=` attribute (S:S6.4), also inside such a body. `path` is `name` followed by one or more dotted field steps (`row.name`, `row.address.city`), each step indexing one level into the current item's `map[string]any` (or a nested map found there). There is no other operator -- no filters, no arithmetic, no comparisons.

- As whole text-run content: `<li>{{row.name}}</li>` compiles the body closure to call `Text(stringOf(row["name"]))` -- resolved fresh on every invocation of `fn`, HTML-escaped (or term-styled) exactly like any other `Text` node, via the real constructor.
- As an `args` token: the resolved value (kept as its native Go type -- `string`, `int`, `float64`, `bool`) is passed straight through as one positional argument, so a translator can format it (pluralise, localise a number) rather than CTML stringifying it up front.

A `{{path}}` token whose root name has no enclosing `<each as="...">`, or whose remaining path segments are not found in the item map at render time, is a **parse-time error** for the former (S:S9 -- the binding is structurally unreachable, caught while walking the body) and renders as an empty string for the latter (a missing field is data absence, not a document defect, mirroring `Options.Get`'s own miss-is-empty convention elsewhere in CoreGO).

### 8.4 What `{{path}}` deliberately does not do

Mixed static-and-dynamic text in one run (`Cost: {{row.amount}}`) is not supported as a single token; wrap the dynamic part in its own element (`Cost: <b>{{row.amount}}</b>`) or use `args` with a catalogue message (`<span args="{{row.amount}}">cost.line</span>`, catalogue entry `"Cost: {0}"`). This keeps word order translator-controlled rather than baked into the markup, and keeps the token grammar to "the whole run or nothing" -- no partial-run scanning to specify.

## 9. Errors

`Parse` returns `(html.Node, error)`. Every error is a `*ctml.ParseError`:

```go
type ParseError struct {
	Line, Col int
	Msg       string
	Cause     error // wrapped via core.E; nil for structural (non-XML) errors
}
```

`Line`/`Col` come from `xml.Decoder.InputPos()` at the point the offending token was read (malformed XML) or the point the offending element/attribute was being interpreted (a reserved element used wrongly, an unbound `{{path}}`, `args` on mixed content, and so on) -- both classes of error carry real, non-zero positions, not just a wrapped stdlib message.

## 10. Package API

```go
package ctml // dappco.re/go/html/ctml

type Bindings struct {
	Sequences map[string][]map[string]any
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

<!-- box-map-section -->

<!-- corecommand-section -->
