---
title: go-html
description: HLCRF DOM compositor with grammar pipeline integration for type-safe server-side HTML generation and optional WASM client rendering.
---

# go-html

`go-html` is a pure-Go library for building HTML documents as type-safe node trees and rendering them to string output. It provides a five-slot layout compositor (Header, Left, Content, Right, Footer -- abbreviated HLCRF), a responsive multi-variant wrapper, a server-side grammar analysis pipeline, a Web Component code generator, and an optional WASM module for client-side rendering.

**Module path:** `forge.lthn.ai/core/go-html`
**Go version:** 1.26
**Licence:** EUPL-1.2

## Quick Start

```go
package main

import html "forge.lthn.ai/core/go-html"

func main() {
    page := html.NewLayout("HCF").
        H(html.El("nav", html.Text("nav.label"))).
        C(html.El("article",
            html.El("h1", html.Text("page.title")),
            html.Each(items, func(item Item) html.Node {
                return html.El("li", html.Text(item.Name))
            }),
        )).
        F(html.El("footer", html.Text("footer.copyright")))

    output := page.Render(html.NewContext())
}
```

This builds a Header-Content-Footer layout with semantic HTML elements (`<header>`, `<main>`, `<footer>`), ARIA roles, and deterministic `data-block` path identifiers. Text nodes pass through the `go-i18n` translation layer and are HTML-escaped by default.

## Package Layout

| Path | Purpose |
|------|---------|
| `node.go` | `Node` interface and all node types: `El`, `Text`, `Raw`, `If`, `Unless`, `Each`, `EachSeq`, `Switch`, `Entitled` |
| `layout.go` | HLCRF compositor with semantic HTML elements and ARIA roles |
| `responsive.go` | Multi-variant breakpoint wrapper (`data-variant` containers) |
| `context.go` | Rendering context: identity, locale, entitlements, i18n service |
| `render.go` | `Render()` convenience function |
| `path.go` | `ParseBlockID()` for decoding `data-block` path attributes |
| `pipeline.go` | `StripTags`, `Imprint`, `CompareVariants` (server-side only, `!js` build tag) |
| `codegen/codegen.go` | Web Component class generation (closed Shadow DOM) |
| `cmd/codegen/main.go` | Build-time CLI: JSON slot map on stdin, JS bundle on stdout |
| `cmd/wasm/main.go` | WASM entry point exporting `renderToString()` to JavaScript |

## Key Concepts

**Node tree** -- All renderable units implement `Node`, a single-method interface: `Render(ctx *Context) string`. The library composes nodes into trees using `El()` for elements, `Text()` for translated text, and control-flow constructors (`If`, `Unless`, `Each`, `Switch`, `Entitled`).

**HLCRF Layout** -- A five-slot compositor that maps to semantic HTML: `<header>` (H), `<aside>` (L/R), `<main>` (C), `<footer>` (F). The variant string controls which slots render: `"HLCRF"` for all five, `"HCF"` for three, `"C"` for content only. Layouts nest: placing a `Layout` inside another layout's slot produces hierarchical `data-block` paths like `L-0-C-0`.

**Responsive variants** -- `Responsive` wraps multiple `Layout` instances with named breakpoints (e.g. `"desktop"`, `"mobile"`). Each variant renders inside a `<div data-variant="name">` container for CSS or JavaScript targeting.

**Grammar pipeline** -- Server-side only. `Imprint()` renders a node tree to HTML, strips tags, tokenises the plain text via `go-i18n/reversal`, and returns a `GrammarImprint` for semantic analysis. `CompareVariants()` computes pairwise similarity scores across responsive variants.

**Web Component codegen** -- `cmd/codegen/` generates ES2022 Web Component classes with closed Shadow DOM from a JSON slot-to-tag mapping. This is a build-time tool, not used at runtime.

## Dependencies

```
forge.lthn.ai/core/go-html
  forge.lthn.ai/core/go-i18n          (direct, all builds)
    forge.lthn.ai/core/go-inference    (indirect, via go-i18n)
  forge.lthn.ai/core/go-i18n/reversal (server builds only, !js)
  github.com/stretchr/testify          (test only)
```

Both `go-i18n` and `go-inference` must be present on the local filesystem. The `go.mod` uses `replace` directives pointing to sibling directories (`../go-i18n`, `../go-inference`).

## Further Reading

- [Architecture](architecture.md) -- Node interface, HLCRF layout internals, responsive compositor, grammar pipeline, WASM module, codegen CLI
- [Development](development.md) -- Building, testing, benchmarks, WASM builds, coding standards, contribution guide
