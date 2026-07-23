# html
**Import:** `dappco.re/go/core/html`
**Files:** 13

## Types

### `Context`
`type Context struct`

Context carries rendering state through the node tree.
Usage example: ctx := NewContext()

Fields:
- `Identity string`
- `Locale string`
- `Entitlements func(feature string) bool`
- `Data map[string]any`
- Unexported fields are present.

Methods:
None.

### `Layout`
`type Layout struct`

Layout is an HLCRF compositor. Arranges nodes into semantic HTML regions
with deterministic path-based IDs.
Usage example: page := NewLayout("HCF").H(Text("title")).C(Text("body"))

Fields:
- No exported fields.
- Unexported fields are present.

Methods:
- `func (l *Layout) C(nodes ...Node) *Layout`
  C appends nodes to the Content (main) slot.
  Usage example: NewLayout("C").C(Text("body"))
- `func (l *Layout) F(nodes ...Node) *Layout`
  F appends nodes to the Footer slot.
  Usage example: NewLayout("CF").F(Text("footer"))
- `func (l *Layout) H(nodes ...Node) *Layout`
  H appends nodes to the Header slot.
  Usage example: NewLayout("HCF").H(Text("title"))
- `func (l *Layout) L(nodes ...Node) *Layout`
  L appends nodes to the Left aside slot.
  Usage example: NewLayout("LC").L(Text("nav"))
- `func (l *Layout) R(nodes ...Node) *Layout`
  R appends nodes to the Right aside slot.
  Usage example: NewLayout("CR").R(Text("ads"))
- `func (l *Layout) Render(ctx *Context) string`
  Render produces the semantic HTML for this layout.
  Usage example: html := NewLayout("C").C(Text("body")).Render(NewContext())
  Only slots present in the variant string are rendered.

### `Node`
`type Node interface`

Node is anything renderable.
Usage example: var n Node = El("div", Text("welcome"))

Members:
- `Render(ctx *Context) string`

Methods:
None.

### `Responsive`
`type Responsive struct`

Responsive wraps multiple Layout variants for breakpoint-aware rendering.
Usage example: r := NewResponsive().Variant("mobile", NewLayout("C"))
Each variant is rendered inside a container with data-variant for CSS targeting.

Fields:
- No exported fields.
- Unexported fields are present.

Methods:
- `func (r *Responsive) Render(ctx *Context) string`
  Render produces HTML with each variant in a data-variant container.
  Usage example: html := NewResponsive().Variant("mobile", NewLayout("C")).Render(NewContext())
- `func (r *Responsive) Variant(name string, layout *Layout) *Responsive`
  Variant adds a named layout variant (e.g., "desktop", "tablet", "mobile").
  Usage example: NewResponsive().Variant("desktop", NewLayout("HLCRF"))
  Variants render in insertion order.

### `Translator`
`type Translator interface`

Translator provides Text() lookups for a rendering context.
Usage example: ctx := NewContextWithService(myTranslator)

The default server build uses go-i18n. Alternate builds, including WASM,
can provide any implementation with the same T() method.

Members:
- `T(key string, args ...any) string`

Methods:
None.

## Functions

### `Attr`
`func Attr(n Node, key, value string) Node`

Attr sets an attribute on an El node. Returns the node for chaining.
Usage example: Attr(El("a", Text("docs")), "href", "/docs")
It recursively traverses through wrappers like If, Unless, and Entitled.

### `CompareVariants`
`func CompareVariants(r *Responsive, ctx *Context) map[string]float64`

CompareVariants runs the imprint pipeline on each responsive variant independently
and returns pairwise similarity scores. Key format: "name1:name2".
Usage example: scores := CompareVariants(NewResponsive(), NewContext())

### `Each`
`func Each[T any](items []T, fn func(T) Node) Node`

Each iterates items and renders each via fn.
Usage example: Each([]string{"a", "b"}, func(v string) Node { return Text(v) })

### `EachSeq`
`func EachSeq[T any](items iter.Seq[T], fn func(T) Node) Node`

EachSeq iterates an iter.Seq and renders each via fn.
Usage example: EachSeq(slices.Values([]string{"a", "b"}), func(v string) Node { return Text(v) })

### `El`
`func El(tag string, children ...Node) Node`

El creates an HTML element node with children.
Usage example: El("section", Text("welcome"))

### `Entitled`
`func Entitled(feature string, node Node) Node`

Entitled renders child only when entitlement is granted. Absent, not hidden.
Usage example: Entitled("beta", Text("preview"))
If no entitlement function is set on the context, access is denied by default.

### `If`
`func If(cond func(*Context) bool, node Node) Node`

If renders child only when condition is true.
Usage example: If(func(ctx *Context) bool { return ctx.Identity != "" }, Text("hi"))

### `Imprint`
`func Imprint(node Node, ctx *Context) reversal.GrammarImprint`

Imprint renders a node tree to HTML, strips tags, tokenises the text,
and returns a GrammarImprint — the full render-reverse pipeline.
Usage example: imp := Imprint(Text("welcome"), NewContext())

### `NewContext`
`func NewContext() *Context`

NewContext creates a new rendering context with sensible defaults.
Usage example: html := Render(Text("welcome"), NewContext())

### `NewContextWithService`
`func NewContextWithService(svc Translator) *Context`

NewContextWithService creates a rendering context backed by a specific translator.
Usage example: ctx := NewContextWithService(myTranslator)

### `NewLayout`
`func NewLayout(variant string) *Layout`

NewLayout creates a new Layout with the given variant string.
Usage example: page := NewLayout("HLCRF")
The variant determines which slots are rendered (e.g., "HLCRF", "HCF", "C").

### `NewResponsive`
`func NewResponsive() *Responsive`

NewResponsive creates a new multi-variant responsive compositor.
Usage example: r := NewResponsive()

### `ParseBlockID`
`func ParseBlockID(id string) []byte`

ParseBlockID extracts the slot sequence from a data-block ID.
Usage example: slots := ParseBlockID("L-0-C-0")
"L-0-C-0" → ['L', 'C']

### `Raw`
`func Raw(content string) Node`

Raw creates a node that renders without escaping (escape hatch for trusted content).
Usage example: Raw("<strong>trusted</strong>")

### `Render`
`func Render(node Node, ctx *Context) string`

Render is a convenience function that renders a node tree to HTML.
Usage example: html := Render(El("main", Text("welcome")), NewContext())

### `StripTags`
`func StripTags(html string) string`

StripTags removes HTML tags from rendered output, returning plain text.
Usage example: text := StripTags("<main>Hello <strong>world</strong></main>")
Tag boundaries are collapsed into single spaces; result is trimmed.
Does not handle script/style element content (go-html does not generate these).

### `Switch`
`func Switch(selector func(*Context) string, cases map[string]Node) Node`

Switch renders based on runtime selector value.
Usage example: Switch(func(ctx *Context) string { return ctx.Locale }, map[string]Node{"en": Text("hello")})

### `Text`
`func Text(key string, args ...any) Node`

Text creates a node that renders through the go-i18n grammar pipeline.
Usage example: Text("welcome", "Ada")
Output is HTML-escaped by default. Safe-by-default path.

### `Unless`
`func Unless(cond func(*Context) bool, node Node) Node`

Unless renders child only when condition is false.
Usage example: Unless(func(ctx *Context) bool { return ctx.Identity == "" }, Text("welcome"))
