# CLAUDE.md

Agent instructions for `go-html`. Module path: `forge.lthn.ai/core/go-html`

## Commands

```bash
go test ./...                                                      # Run all tests
go test -run TestName ./...                                        # Single test
go test -short ./...                                               # Skip slow WASM build test
go test -bench . ./...                                             # Benchmarks
go vet ./...                                                       # Static analysis
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o gohtml.wasm ./cmd/wasm/  # WASM build
echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ # Codegen CLI
```

## Architecture

See `docs/architecture.md` for full detail. Summary:

- **Node interface**: `Render(ctx *Context) string` — El, Text, Raw, If, Unless, Each[T], Switch, Entitled
- **HLCRF Layout**: Header/Left/Content/Right/Footer compositor with ARIA roles and deterministic `data-block` IDs
- **Responsive**: Multi-variant breakpoint wrapper (`data-variant` attributes)
- **Pipeline**: Render → StripTags → go-i18n/reversal Tokenise → GrammarImprint (server-side only)
- **Codegen**: Web Component classes with closed Shadow DOM, generated at build time by `cmd/codegen/`
- **WASM**: `cmd/wasm/` exports `renderToString()` only — 2.90 MB raw / 842 KB gzip

## Server/Client Split

Files guarded with `//go:build !js` are excluded from WASM:

- `pipeline.go` — Imprint/CompareVariants use `go-i18n/reversal` (server-side only)
- `cmd/wasm/register.go` — encoding/json + codegen (replaced by `cmd/codegen/` CLI)

Never import `encoding/json`, `text/template`, or `fmt` in WASM-linked code. Use string concatenation instead of `fmt.Sprintf` in `layout.go` and any other file without a `!js` guard.

## Key Files

| File | Purpose |
|------|---------|
| `node.go` | All node types (El, Text, Raw, If, Unless, Each, Switch, Entitled) |
| `layout.go` | HLCRF compositor |
| `pipeline.go` | StripTags, Imprint, CompareVariants (!js only) |
| `responsive.go` | Multi-variant breakpoint wrapper |
| `context.go` | Rendering context (Identity, Locale, Entitlements, i18n Service) |
| `codegen/codegen.go` | Web Component class generation |
| `cmd/wasm/main.go` | WASM entry point (renderToString only) |
| `cmd/codegen/main.go` | Build-time CLI for WC bundle generation |
| `cmd/wasm/size_test.go` | WASM binary size gate (< 1 MB gzip, < 3 MB raw) |

## Dependencies

- `forge.lthn.ai/core/go-i18n` (replace directive → `../go-i18n`)
- `go-i18n` and `go-inference` must be present alongside this repo for builds

## Coding Standards

- UK English (colour, organisation, centre)
- All types annotated
- Tests use `testify` assert/require
- Licence: EUPL-1.2 — add `// SPDX-Licence-Identifier: EUPL-1.2` to new files
- Safe-by-default: HTML escaping on Text nodes, void element handling, entitlement deny-by-default
- Deterministic output: sorted attributes, reproducible paths
- Commits: conventional commits + `Co-Authored-By: Virgil <virgil@lethean.io>`

## Test Conventions

No specific suffix pattern. Use table-driven subtests with `t.Run()`. Integration tests that use `Text` nodes must call `i18n.SetDefault(svc)` before rendering.
