---
title: Architecture
description: Internals of the go-html HLCRF DOM compositor, covering the node interface, layout system, responsive wrapper, grammar pipeline, WASM module, and codegen CLI.
---

# Architecture

`go-html` is structured around a single interface, a layout compositor, and a server-side analysis pipeline. Everything renders to `string` -- there is no virtual DOM, no diffing, and no retained state between renders.

## Node Interface

Every renderable unit implements one method:

```go
type Node interface {
    Render(ctx *Context) string
}
```

All concrete node types are unexported structs with exported constructor functions. The public API surface consists of nine node constructors, four accessibility helpers, plus the `Attr()` and `Render()` helpers:

| Constructor | Behaviour |
|-------------|-----------|
| `El(tag, ...Node)` | HTML element with children. Void elements (`br`, `img`, `input`, etc.) never emit a closing tag. |
| `Attr(Node, key, value)` | Sets an attribute on an `El` node. Traverses through `If`, `Unless`, `Entitled`, `Each`, `EachSeq`, `Switch`, `Layout`, and `Responsive` wrappers. Returns the node for chaining. |
| `AriaLabel(Node, label)` | Convenience helper that sets `aria-label` on an element node. |
| `AltText(Node, text)` | Convenience helper that sets `alt` on an element node. |
| `TabIndex(Node, index)` | Convenience helper that sets `tabindex` on an element node. |
| `AutoFocus(Node)` | Convenience helper that sets `autofocus` on an element node. |
| `Role(Node, role)` | Convenience helper that sets `role` on an element node. |
| `Text(key, ...any)` | Translated text via the active context translator. Server builds fall back to global `go-i18n`; JS builds fall back to the key. Output is always HTML-escaped. |
| `Raw(content)` | Unescaped trusted content. Explicit escape hatch. |
| `If(cond, Node)` | Renders the child only when the condition function returns true. |
| `Unless(cond, Node)` | Renders the child only when the condition function returns false. |
| `Each[T](items, fn)` | Iterates a slice and renders each item via a mapping function. Generic over `T`. |
| `EachSeq[T](items, fn)` | Same as `Each` but accepts an `iter.Seq[T]` instead of a slice. |
| `Switch(selector, cases)` | Renders one of several named cases based on a runtime selector function. Returns empty string when no case matches. |
| `Entitled(feature, Node)` | Renders the child only when the context's entitlement function grants the named feature. Deny-by-default: returns empty string when no entitlement function is set. |

### Safety Guarantees

- **XSS prevention**: `Text()` nodes always HTML-escape their output via `html.EscapeString()`. User-supplied strings passed through `Text()` cannot inject HTML.
- **Attribute escaping**: Attribute values are escaped with `html.EscapeString()`, handling `&`, `<`, `>`, `"`, and `'`.
- **Deterministic output**: Attribute keys on `El` nodes are sorted alphabetically before rendering, producing identical output regardless of insertion order.
- **Void elements**: A lookup table of 13 void elements (`area`, `base`, `br`, `col`, `embed`, `hr`, `img`, `input`, `link`, `meta`, `source`, `track`, `wbr`) ensures these never emit a closing tag.
- **Deny-by-default entitlements**: `Entitled` returns an empty string when the context is nil, when no entitlement function is set, or when the function returns false. Content is absent from the DOM, not merely hidden.

## Rendering Context

The `Context` struct carries per-request state through the node tree during rendering:

```go
type Context struct {
    Identity     string                     // e.g. user ID or session identifier
    Locale       string                     // BCP 47 locale string
    Entitlements func(feature string) bool  // feature gate callback
    Data         map[string]any             // arbitrary per-request data
    Metadata     map[string]any             // alias of Data for alternate naming
    service      Translator                 // unexported; set via constructor
}
```

Two constructors are provided:

- `NewContext()` creates a context with sensible defaults and an empty `Data` map.
- `NewContextWithService(svc)` creates a context backed by any translator implementing `T(key, ...any) string` such as `*i18n.Service`.

`Data` and `Metadata` point at the same backing map when the context is created through `NewContext()`. Use whichever name is clearer in the calling code. `SetLocale()` and `SetService()` keep the active translator in sync when either value changes.

The `service` field is intentionally unexported. When nil, server builds fall back to Core's i18n translator while JS builds render the key unchanged. This prevents callers from setting the service inconsistently after construction while keeping the WASM import graph lean.

## HLCRF Layout

The `Layout` type is a compositor for five named slots:

| Slot Letter | Semantic Element | ARIA Role | Accessor |
|-------------|-----------------|-----------|----------|
| H | `<header>` | `banner` | `layout.H(...)` |
| L | `<nav>` | `navigation` | `layout.L(...)` |
| C | `<main>` | `main` | `layout.C(...)` |
| R | `<aside>` | `complementary` | `layout.R(...)` |
| F | `<footer>` | `contentinfo` | `layout.F(...)` |

### Variant String

The variant string passed to `NewLayout()` determines which slots render and in which order:

```go
NewLayout("HLCRF")  // all five slots
NewLayout("HCF")    // header, content, footer (no sidebars)
NewLayout("C")      // content only
NewLayout("LC")     // left sidebar and content
```

Slot letters not present in the variant string are ignored, even if nodes have been appended to those slots. Unrecognised characters (lowercase, digits, special characters) are silently skipped -- no error is returned.

### Deterministic Block IDs

Each rendered slot receives a `data-block` attribute encoding its position in the layout tree. At the root level, IDs use the slot letter itself:

```html
<header role="banner" data-block="H">...</header>
<main role="main" data-block="C">...</main>
<footer role="contentinfo" data-block="F">...</footer>
```

Block IDs are constructed by simple string concatenation (no `fmt.Sprintf`) to keep the `fmt` package out of the WASM import graph.

### Nested Layouts

`Layout` implements `Node`, so a layout can be placed inside any slot of another layout. At render time, nested layouts retain the parent's block ID as a prefix. This produces hierarchical paths:

```go
inner := html.NewLayout("HCF").
    H(html.Raw("nav")).
    C(html.Raw("body")).
    F(html.Raw("links"))

outer := html.NewLayout("HLCRF").
    H(html.Raw("top")).
    L(inner).              // inner layout nested in the Left slot
    C(html.Raw("main")).
    F(html.Raw("foot"))
```

The inner layout's slots render with prefixed block IDs: `L.0`, `L.0.1`, `L.0.2`. At 10 levels of nesting, the deepest block ID becomes `C.0.0.0.0.0.0.0.0.0` (tested in `edge_test.go`).

The clone-on-render approach means the original layout is never mutated. This is safe for concurrent use.

### Fluent Builder

All slot methods return `*Layout` for chaining. Multiple nodes can be appended to the same slot across multiple calls:

```go
html.NewLayout("HCF").
    H(html.El("h1", html.Text("page.title"))).
    C(html.El("p", html.Text("intro"))).
    C(html.El("p", html.Text("body"))).       // appends to the same C slot
    F(html.El("small", html.Text("footer")))
```

### Block ID Parsing

`ParseBlockID()` in `path.go` extracts the slot letter sequence from a `data-block` attribute value:

```go
ParseBlockID("L.0.C.0")       // returns ['L', 'C']
ParseBlockID("L-0-C-0")       // legacy hyphenated form, also returns ['L', 'C']
ParseBlockID("C.0.C.0.C.0")    // returns ['C', 'C', 'C']
ParseBlockID("H")              // returns ['H']
ParseBlockID("")               // returns nil
```

This enables server-side or client-side code to locate a specific block in the rendered tree by its structural path.

## Responsive Compositor

`Responsive` wraps multiple named `Layout` variants for breakpoint-aware rendering:

```go
html.NewResponsive().
    Variant("desktop", html.NewLayout("HLCRF").
        H(html.Raw("header")).L(html.Raw("nav")).C(html.Raw("main")).
        R(html.Raw("aside")).F(html.Raw("footer"))).
    Variant("tablet", html.NewLayout("HCF").
        H(html.Raw("header")).C(html.Raw("main")).F(html.Raw("footer"))).
    Variant("mobile", html.NewLayout("C").
        C(html.Raw("main")))
```

Each variant renders inside a `<div data-variant="name">` container. Variants render in insertion order. When supplied, `Responsive.Add(name, layout, media)` also emits `data-media="..."` on the wrapper so downstream CSS can reflect the breakpoint hint. CSS media queries or JavaScript can target these containers for show/hide logic.

`VariantSelector(name)` returns a CSS attribute selector for a specific responsive variant, making stylesheet targeting less error-prone than hand-writing the attribute selector repeatedly.

`Responsive` implements `Node`, so it can be passed to `Render()` or `Imprint()`. The `Variant()` method accepts `*Layout` specifically, not arbitrary `Node` values.

Each variant maintains independent block ID namespaces -- nesting a layout inside a responsive variant does not conflict with the same layout structure in another variant.

## Grammar Pipeline (Server-Side Only)

The grammar pipeline is excluded from WASM builds via `//go:build !js` on `pipeline.go`. It bridges the rendering layer to the semantic analysis layer.

### StripTags

```go
func StripTags(html string) string
```

Converts rendered HTML to plain text. Tag boundaries are collapsed into single spaces; the result is trimmed. The implementation is a single-pass rune scanner with no regular expressions and no allocations beyond the output `strings.Builder`. It does not handle `<script>` or `<style>` content because `go-html` never generates those elements.

### Imprint

```go
func Imprint(node Node, ctx *Context) reversal.GrammarImprint
```

Runs the full render-to-analysis pipeline:

1. Renders the node tree to HTML via `node.Render(ctx)`.
2. Strips HTML tags via `StripTags()` to extract plain text.
3. Tokenises the text via `go-i18n/reversal.NewTokeniser().Tokenise()`.
4. Wraps tokens in a `reversal.GrammarImprint` for structural analysis.

The resulting `GrammarImprint` exposes `TokenCount`, `UniqueVerbs`, and a `Similar()` method for pairwise semantic similarity scoring.

A nil context is handled gracefully: `Imprint` creates a default context internally.

### CompareVariants

```go
func CompareVariants(r *Responsive, ctx *Context) map[string]float64
```

Runs `Imprint` independently on each named layout variant in a `Responsive` and returns pairwise similarity scores. Keys are formatted as `"name1:name2"`.

This enables detection of semantically divergent responsive variants -- for example, a mobile layout that strips critical information present in the desktop variant. Same-content variants with different layout structures (e.g. `HLCRF` vs `HCF`) score above 0.8 similarity.

A single-variant `Responsive` produces an empty score map (no pairs to compare).

## WASM Module

The WASM entry point at `cmd/wasm/main.go` is compiled with `GOOS=js GOARCH=wasm` and exposes a single JavaScript function:

```js
gohtml.renderToString(variant, locale, slots)
```

**Parameters:**

- `variant` (string): HLCRF variant string, e.g. `"HCF"`.
- `locale` (string): BCP 47 locale string for i18n, e.g. `"en-GB"`.
- `slots` (object): Optional keys `H`, `L`, `C`, `R`, `F` containing HTML strings.

Slot content is injected via `Raw()`. The caller is responsible for sanitisation -- the WASM module is a rendering engine for trusted content produced server-side or by the application's own templates.

### Size Budget

The WASM binary has a size gate enforced by `cmd/wasm/size_test.go`:

| Metric | Limit | Current |
|--------|-------|---------|
| Raw binary | 3.5 MB | ~2.90 MB |
| Gzip compressed | 1 MB | ~842 KB |

The test builds the WASM binary as a subprocess and is skipped under `go test -short`. The Makefile `wasm` target performs the same build with size checking.

### Server/Client Split

The binary split is enforced by Go build tags:

| File | Build Tag | Reason for WASM Exclusion |
|------|-----------|--------------------------|
| `pipeline.go` | `!js` | Imports `go-i18n/reversal` |
| `cmd/wasm/register.go` | `!js` | Imports `encoding/json` and `text/template` |

The WASM binary includes only: node types, layout, responsive, context, render, path, and `go-i18n` core translation. No codegen, no pipeline, no JSON, no templates, no `fmt`.

## Codegen CLI

`cmd/codegen/main.go` generates Web Component JavaScript bundles from HLCRF slot assignments at build time:

```bash
echo '{"H":"nav-bar","C":"main-content","F":"page-footer"}' | go run ./cmd/codegen/ > components.js
```

The `codegen` package (`codegen/codegen.go`) generates ES2022 class definitions with closed Shadow DOM. For each custom element tag, it produces:

1. A class extending `HTMLElement` with a private `#shadow` field.
2. `constructor()` attaching a closed shadow root (`mode: "closed"`).
3. `connectedCallback()` dispatching a `wc-ready` custom event with the tag name and slot.
4. `render(html)` method that sets shadow content from a `<template>` clone.
5. A `customElements.define()` registration call.

Tag names must contain a hyphen (Web Components specification requirement). `TagToClassName()` converts kebab-case to PascalCase: `nav-bar` becomes `NavBar`, `my-super-widget` becomes `MySuperWidget`.

`GenerateBundle()` deduplicates tags -- if the same tag is assigned to multiple slots, only one class definition is emitted.

The codegen CLI uses `encoding/json` and `text/template`, which are excluded from the WASM build. Consumers generate the JS bundle at build time and serve it as a static asset.

## Data Flow Summary

```
                         Server-Side
                    +-------------------+
                    |                   |
  Node tree ------->  Render(ctx)      |-----> HTML string
                    |                   |
                    |  StripTags()      |-----> plain text
                    |                   |
                    |  Imprint()        |-----> GrammarImprint
                    |                   |         .TokenCount
                    |  CompareVariants()|         .UniqueVerbs
                    |                   |         .Similar()
                    +-------------------+

                         WASM Client
                    +-------------------+
                    |                   |
  JS call --------->  renderToString() |-----> HTML string
  (variant, locale, |                   |
   slots object)    +-------------------+

                         Build Time
                    +-------------------+
                    |                   |
  JSON slot map --->  cmd/codegen/     |-----> Web Component JS
  (stdin)           |                   |       (stdout)
                    +-------------------+
```
