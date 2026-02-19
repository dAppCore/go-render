# CLAUDE.md

## Project

`go-html` is an HLCRF DOM compositor with grammar pipeline. Module path: `forge.lthn.ai/core/go-html`

## Commands

```bash
go test ./...                    # Run all tests
go test -run TestName ./...      # Single test
go test -bench . ./...           # Benchmarks
go vet ./...                     # Static analysis
GOOS=js GOARCH=wasm go build -o gohtml.wasm ./cmd/wasm/  # WASM build
```

## Architecture

- **Node interface**: `Render(ctx *Context) string` — El, Text, Raw, If, Unless, Each[T], Switch, Entitled
- **HLCRF Layout**: Header/Left/Content/Right/Footer compositor with ARIA roles and deterministic `data-block` IDs
- **Responsive**: Multi-variant breakpoint wrapper (`data-variant` attributes)
- **Pipeline**: Render → StripTags → go-i18n/reversal Tokenise → GrammarImprint
- **Codegen**: Web Component classes with closed Shadow DOM
- **WASM**: `cmd/wasm/` exports `renderToString()` and `registerComponents()` to JS

## Dependencies

- `forge.lthn.ai/core/go-i18n` (replace directive → `../go-i18n`)
- go-i18n must be present alongside this repo for builds

## Coding Standards

- UK English (colour, organisation, centre)
- All types annotated
- Tests use `testify` assert/require
- Licence: EUPL-1.2
- Safe-by-default: HTML escaping on Text nodes, void element handling, entitlement deny-by-default
- Deterministic output: sorted attributes, reproducible paths

## Test Conventions

No specific suffix pattern — use table-driven subtests with `t.Run()`.

## Key Files

| File | Purpose |
|------|---------|
| `node.go` | All node types (El, Text, Raw, If, Unless, Each, Switch, Entitled) |
| `layout.go` | HLCRF compositor |
| `pipeline.go` | StripTags, Imprint, CompareVariants |
| `responsive.go` | Multi-variant breakpoint wrapper |
| `context.go` | Rendering context (Identity, Locale, Entitlements, i18n Service) |
| `codegen/codegen.go` | Web Component class generation |
| `cmd/wasm/main.go` | WASM entry point |
