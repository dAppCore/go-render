[![Go Reference](https://pkg.go.dev/badge/dappco.re/go/core/html.svg)](https://pkg.go.dev/dappco.re/go/core/html)
[![License: EUPL-1.2](https://img.shields.io/badge/License-EUPL--1.2-blue.svg)](LICENSE.md)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](go.mod)

# go-html

HLCRF DOM compositor with grammar pipeline integration for server-side HTML generation and optional WASM client rendering. Provides a type-safe node tree (El, Text, Raw, If, Each, Switch, Entitled, AriaLabel, AltText), a five-slot Header/Left/Content/Right/Footer layout compositor with deterministic `data-block` path IDs and ARIA roles, a responsive multi-variant wrapper, a server-side grammar pipeline (StripTags, GrammarImprint via go-i18n reversal, CompareVariants), a build-time Web Component codegen CLI, and a WASM module (2.90 MB raw, 842 KB gzip) exposing `renderToString()`.

**Module**: `dappco.re/go/core/html`
**Licence**: EUPL-1.2
**Language**: Go 1.26

## Quick Start

```go
import html "dappco.re/go/core/html"

page := html.NewLayout("HCF").
    H(html.El("nav", html.Text("i18n.label.navigation"))).
    C(html.El("main",
        html.El("h1", html.Text("i18n.label.welcome")),
        html.Each(items, func(item Item) html.Node {
            return html.El("li", html.Text(item.Name))
        }),
    )).
    F(html.El("footer", html.Text("i18n.label.copyright")))

rendered := page.Render(html.NewContext("en-GB"))
```

## Documentation

- [Architecture](docs/architecture.md) — node interface, HLCRF layout, responsive compositor, grammar pipeline, WASM module, codegen CLI
- [Development Guide](docs/development.md) — building, testing, WASM build, server/client split rules
- [Project History](docs/history.md) — completed phases and known limitations

## Build & Test

```bash
go test ./...
go test -bench . ./...
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o gohtml.wasm ./cmd/wasm/
go build ./...
```

## Licence

European Union Public Licence 1.2 — see [LICENCE](LICENCE) for details.
