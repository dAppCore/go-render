# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Agent instructions for `go-html`. Module path: `forge.lthn.ai/core/go-html`

## Commands

```bash
go test ./...                                                      # Run all tests
go test -run TestName ./...                                        # Single test
go test -short ./...                                               # Skip slow WASM build test
go test -bench . ./...                                             # Benchmarks
go test -bench . -benchmem ./...                                   # Benchmarks with alloc stats
go vet ./...                                                       # Static analysis
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o gohtml.wasm ./cmd/wasm/  # WASM build
make wasm                                                          # WASM build with size gate
echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ # Codegen CLI
```

## Architecture

See `docs/architecture.md` for full detail. Summary:

- **Node interface**: `Render(ctx *Context) string` — El, Text, Raw, If, Unless, Each[T], EachSeq[T], Switch, Entitled
- **HLCRF Layout**: Header/Left/Content/Right/Footer compositor with ARIA roles and deterministic `data-block` IDs. Variant string (e.g. "HCF", "HLCRF", "C") controls which slots render. Layouts nest via clone-on-render (thread-safe).
- **Responsive**: Multi-variant breakpoint wrapper (`data-variant` attributes), renders all variants in insertion order
- **Pipeline**: Render → StripTags → go-i18n/reversal Tokenise → GrammarImprint (server-side only)
- **Codegen**: Web Component classes with closed Shadow DOM, generated at build time by `cmd/codegen/`
- **WASM**: `cmd/wasm/` exports `renderToString()` only — size gate: < 3.5 MB raw, < 1 MB gzip

## Server/Client Split

Files guarded with `//go:build !js` are excluded from WASM:

- `pipeline.go` — Imprint/CompareVariants use `go-i18n/reversal` (server-side only)
- `cmd/wasm/register.go` — encoding/json + codegen (replaced by `cmd/codegen/` CLI)

**Critical WASM constraint**: Never import `encoding/json`, `text/template`, or `fmt` in WASM-linked code (files without a `!js` build tag). Use string concatenation instead of `fmt.Sprintf` in `layout.go`, `node.go`, `responsive.go`, `render.go`, `path.go`, and `context.go`. The `fmt` package alone adds ~500 KB to the WASM binary.

## Dependencies

- `forge.lthn.ai/core/go-i18n` (replace directive → `../go-i18n`)
- `forge.lthn.ai/core/go-inference` (indirect, via go-i18n)
- Both `go-i18n` and `go-inference` must be cloned alongside this repo for builds
- Go 1.26+ required (uses `range` over integers, `iter.Seq`, `maps.Keys`, `slices.Collect`)

## Coding Standards

- UK English (colour, organisation, centre, behaviour, licence, serialise)
- All types annotated; use `any` not `interface{}`
- Tests use `testify` assert/require
- Licence: EUPL-1.2 — add `// SPDX-Licence-Identifier: EUPL-1.2` to new files
- Safe-by-default: HTML escaping via `html.EscapeString()` on Text nodes and attribute values, void element handling, entitlement deny-by-default
- Deterministic output: sorted attributes on El nodes, reproducible block ID paths
- Errors: use `log.E("scope", "message", err)` from `go-log`, never `fmt.Errorf`
- File I/O: use `coreio.Local` from `go-io`, never `os.ReadFile`/`os.WriteFile`
- Commits: conventional commits + `Co-Authored-By: Virgil <virgil@lethean.io>`

## Test Conventions

Use table-driven subtests with `t.Run()`. Integration tests that use `Text` nodes must initialise i18n before rendering:

```go
svc, _ := i18n.New()
i18n.SetDefault(svc)
```
