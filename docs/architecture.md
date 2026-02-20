# Architecture

`go-html` is an HLCRF DOM compositor with grammar pipeline integration. It provides a pure-Go, type-safe HTML rendering library designed for server-side generation with an optional lightweight WASM client module.

Module path: `forge.lthn.ai/core/go-html`

## Node Interface

All renderable units implement a single interface:

```go
type Node interface {
    Render(ctx *Context) string
}
```

Every node type is a private struct with a public constructor. The API surface is intentionally small: nine public constructors plus `Attr()` and `Render()` helpers.

| Constructor | Description |
|-------------|-------------|
| `El(tag, ...Node)` | HTML element with children |
| `Attr(Node, key, value)` | Set attribute on an El node; chainable |
| `Text(key, ...any)` | Translated, HTML-escaped text via go-i18n |
| `Raw(content)` | Unescaped trusted content |
| `If(cond, Node)` | Conditional render |
| `Unless(cond, Node)` | Inverse conditional render |
| `Each[T](items, fn)` | Type-safe iteration with generics |
| `Switch(selector, cases)` | Runtime dispatch to named cases |
| `Entitled(feature, Node)` | Entitlement-gated render; deny-by-default |

### Safety guarantees

- `Text` nodes are always HTML-escaped. XSS via user-supplied strings fed through `Text()` is not possible.
- `Raw` is an explicit escape hatch for trusted content only. Its name signals intent.
- `Entitled` returns an empty string when no entitlement function is set on the context. Access is denied by default, not granted.
- `El` attributes are sorted alphabetically before output, producing deterministic HTML regardless of insertion order.
- Void elements (`br`, `img`, `input`, etc.) never emit a closing tag.

## HLCRF Layout

The `Layout` type is a compositor for five named slots: **H**eader, **L**eft, **C**ontent, **R**ight, **F**ooter. Each slot maps to a specific semantic HTML element and ARIA role:

| Slot | Element | ARIA role |
|------|---------|-----------|
| H | `<header>` | `banner` |
| L | `<aside>` | `complementary` |
| C | `<main>` | `main` |
| R | `<aside>` | `complementary` |
| F | `<footer>` | `contentinfo` |

A layout variant string selects which slots are rendered and in which order:

```go
NewLayout("HLCRF")   // all five slots
NewLayout("HCF")     // header, content, footer — no sidebars
NewLayout("C")       // content only
```

Each rendered slot receives a deterministic `data-block` attribute encoding its position in the tree. The root layout produces IDs in the form `{slot}-0` (e.g., `H-0`, `C-0`). Nested layouts extend the parent's block ID as a path prefix: a `Layout` placed inside the `L` slot of a root layout will produce inner slot IDs like `L-0-H-0`, `L-0-C-0`.

This path scheme is computed without `fmt.Sprintf` — using simple string concatenation — to keep `fmt` out of the WASM import graph.

### Nested layouts

`Layout` implements `Node`, so it can be placed inside any slot of another layout. At render time, nested layouts are cloned and their `path` field is set to the parent's block ID. This clone-on-render approach avoids shared mutation and is safe for concurrent use.

```go
inner := NewLayout("HCF").H(Raw("nav")).C(Raw("body")).F(Raw("links"))
outer := NewLayout("HLCRF").H(Raw("top")).L(inner).C(Raw("main")).F(Raw("foot"))
```

### Fluent builder

All slot methods return the `*Layout` for chaining. Multiple nodes may be appended to the same slot across multiple calls:

```go
NewLayout("HCF").
    H(El("h1", Text("Title"))).
    C(El("p", Text("Content")), Raw("<hr>")).
    F(El("small", Text("Copyright")))
```

## Responsive Compositor

`Responsive` wraps multiple named `Layout` variants for breakpoint-aware rendering. Each variant renders inside a `<div data-variant="name">` container, giving CSS media queries or JavaScript a stable hook for show/hide logic.

```go
NewResponsive().
    Variant("desktop", NewLayout("HLCRF")...).
    Variant("tablet", NewLayout("HCF")...).
    Variant("mobile", NewLayout("C")...)
```

`Responsive` itself implements `Node` and may be passed to `Imprint()` for cross-variant semantic analysis.

Note: `Responsive.Variant()` accepts only `*Layout`, not arbitrary `Node` values. Arbitrary subtrees must be wrapped in a layout first.

## Rendering Context

`Context` carries per-request state through the entire node tree:

```go
type Context struct {
    Identity     string
    Locale       string
    Entitlements func(feature string) bool
    Data         map[string]any
    service      *i18n.Service  // private; set via NewContextWithService()
}
```

The `service` field is intentionally unexported. Custom i18n adapter injection requires `NewContextWithService(svc)`. This prevents callers from setting it inconsistently after construction.

When `ctx.service` is nil, `Text` nodes fall back to the global `i18n.T()` default service.

## Grammar Pipeline

The grammar pipeline is a server-side-only feature. It is guarded with `//go:build !js` and absent from all WASM builds.

### StripTags

`StripTags(html string) string` converts rendered HTML to plain text. Tag boundaries are collapsed to single spaces; the result is trimmed. The implementation is a single-pass rune scanner: no regex, no allocations beyond the output builder. It does not attempt to elide `<script>` or `<style>` content because `go-html` never generates those elements.

### Imprint

`Imprint(node Node, ctx *Context) reversal.GrammarImprint` runs the full render-to-analysis pipeline:

1. Call `node.Render(ctx)` to produce HTML.
2. Pass HTML through `StripTags` to extract plain text.
3. Pass plain text through `go-i18n/reversal.Tokeniser` to produce a token sequence.
4. Wrap tokens in a `reversal.GrammarImprint` for structural analysis.

The resulting `GrammarImprint` exposes `TokenCount`, `UniqueVerbs`, and a `Similar()` method for pairwise semantic similarity scoring. This bridges the rendering layer to the privacy and analytics layers of the Lethean stack.

### CompareVariants

`CompareVariants(r *Responsive, ctx *Context) map[string]float64` runs `Imprint` on each named layout variant in a `Responsive` and returns pairwise similarity scores. Keys are `"name1:name2"`. This enables detection of semantically divergent responsive variants — for example, a mobile layout that strips critical information that appears in the desktop variant.

## Server/Client Split

The binary split is enforced by Go build tags.

| File | Build tag | Reason for exclusion from WASM |
|------|-----------|-------------------------------|
| `pipeline.go` | `//go:build !js` | Imports `go-i18n/reversal` (~250 KB gzip) |
| `cmd/wasm/register.go` | `//go:build !js` | Imports `encoding/json` (~200 KB gzip) and `text/template` (~125 KB gzip) |

The WASM binary includes only: node types, layout, responsive, context, render, path, and go-i18n core (translation). No codegen, no pipeline, no JSON, no templates, no `fmt`.

## WASM Module

The WASM entry point is `cmd/wasm/main.go`, compiled with `GOOS=js GOARCH=wasm`.

It exposes a single JavaScript function on `window.gohtml`:

```js
gohtml.renderToString(variant, locale, slots)
```

- `variant`: HLCRF variant string, e.g. `"HCF"`.
- `locale`: BCP 47 locale string for i18n, e.g. `"en-GB"`.
- `slots`: object with optional keys `H`, `L`, `C`, `R`, `F` containing HTML strings.

Slot content is injected via `Raw()`. The caller is responsible for sanitisation. This is intentional: the WASM module is a rendering engine for trusted content produced server-side or by the application's own templates.

### Size gate

`cmd/wasm/size_test.go` contains `TestWASMBinarySize_Good`, a build-gated test that:

1. Builds the WASM binary with `-ldflags=-s -w`.
2. Gzip-compresses the output at best compression.
3. Asserts the compressed size is below 1,048,576 bytes (1 MB).
4. Asserts the raw size is below 3,145,728 bytes (3 MB).

This test is skipped under `go test -short`. It is guarded with `//go:build !js` so it does not run within the WASM environment itself. Current measured size: 2.90 MB raw, 842 KB gzip.

## Codegen CLI

`cmd/codegen/main.go` is a build-time tool for generating Web Component JavaScript bundles from HLCRF slot assignments. It reads a JSON slot map from stdin and writes the generated JS to stdout.

```bash
echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
```

The `codegen` package generates ES2022 class definitions with closed Shadow DOM. The generated pattern per component:

- A class extending `HTMLElement` with a private `#shadow` field.
- `constructor()` attaches a closed shadow root (`mode: "closed"`).
- `connectedCallback()` dispatches a `wc-ready` custom event with the tag name and slot.
- `render(html)` sets shadow content from a `<template>` clone.
- `customElements.define()` registration.

Closed Shadow DOM provides style isolation. Content is set via the DOM API, never via `innerHTML` directly on the element.

Tag names must contain a hyphen (Web Components specification requirement). `TagToClassName()` converts kebab-case tags to PascalCase class names: `nav-bar` becomes `NavBar`.

The codegen CLI uses `encoding/json` and `text/template`, which are excluded from the WASM build. Consumers generate the JS bundle at build time, not at runtime.

## Block ID Path Scheme

`path.go` exports `ParseBlockID(id string) []byte`, which extracts the slot letter sequence from a `data-block` attribute value.

Format: slots are separated by `-0-`. The sequence `L-0-C-0` decodes to `['L', 'C']`, meaning the content slot of a layout nested inside the left slot.

This scheme is deterministic and human-readable. It enables server-side or client-side code to locate a specific block in the rendered tree by path.

## Dependency Graph

```
go-html
├── forge.lthn.ai/core/go-i18n          (direct, all builds)
│   └── forge.lthn.ai/core/go-inference (indirect)
├── forge.lthn.ai/core/go-i18n/reversal (server builds only, !js)
└── github.com/stretchr/testify         (test only)
```

Both `go-i18n` and `go-html` are developed in parallel. The `go.mod` uses a `replace` directive pointing to `../go-i18n`. Both repositories must be present on the local filesystem for builds and tests.
