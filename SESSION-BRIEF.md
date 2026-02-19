# Session Brief: core/go-html

**Repo**: `forge.lthn.ai/core/go-html` (clone at `/tmp/core-go-html`)
**Module**: `forge.lthn.ai/core/go-html`
**Status**: 19 Go files, 1,591 LOC, 53 tests ALL PASS
**Wiki**: https://forge.lthn.ai/core/go-html/wiki (6 pages)

## What This Is

HLCRF DOM compositor with grammar pipeline. Renders semantic HTML from composable node trees with:
- **Node interface**: El, Text, Raw, If, Unless, Each[T], Switch, Entitled
- **HLCRF Layout**: Header/Left/Content/Right/Footer with ARIA roles
- **Responsive**: Multi-variant breakpoint rendering
- **Pipeline**: Render → strip tags → tokenise via go-i18n/reversal → GrammarImprint
- **WASM target**: `cmd/wasm/` exposes renderToString() and registerComponents() to JS
- **Codegen**: Web Component classes with closed Shadow DOM

## Current State

| Area | Status |
|------|--------|
| Core (node, layout, responsive, pipeline) | SOLID — all tested, clean API |
| Tests | 53/53 pass, excellent coverage ratios |
| go vet | Clean |
| TODOs/FIXMEs | None |
| WASM build | FAILS — 1.58 MB gzip exceeds 1 MB Makefile limit |
| Codegen | Working — generates WC classes |

## Dependencies

- `forge.lthn.ai/core/go-i18n` (replace directive → `../go-i18n`)
- `github.com/stretchr/testify` v1.11.1
- `golang.org/x/text` v0.33.0

## Priority Work

### High (blockers)
1. **Fix WASM size** — Move `buildComponentJS()` / JSON parsing to server-side. WASM should only do `Render()`. Current: 6.0 MB raw / 1.58 MB gzip.
2. **WASM integration tests** — No `cmd/wasm/main_test.go` exists. Can't test JS↔Go round-trip.

### Medium (completeness)
3. **Performance benchmarks** — No `BenchmarkRender()` or `BenchmarkImprint()`. Add them.
4. **TypeScript type definitions** — Codegen only produces JS. Add `.d.ts` generation for WC bundle.
5. **Accessibility helpers** — Layout has semantic HTML + ARIA roles, but no aria-label builder or alt text helpers.
6. **Layout variant validation** — `NewLayout("XYZ")` silently produces empty output. Could warn.

### Low (hardening)
7. **Unicode/RTL edge cases** — Test emoji, RTL text in Text nodes
8. **Deep nesting stress test** — Circular or very deep Layout nesting
9. **Large Each[T]** — Test with thousands of items
10. **Browser polyfill docs** — Closed Shadow DOM support matrix

## File Map

```
/tmp/core-go-html/
├── node.go (254)        + node_test.go (206)
├── layout.go (119)      + layout_test.go (116)
├── pipeline.go (83)     + pipeline_test.go (128)
├── responsive.go (39)   + responsive_test.go (89)
├── context.go (27)
├── render.go (9)        + render_test.go (97)
├── path.go (22)         + path_test.go (86)
├── integration_test.go (52)
├── cmd/wasm/
│   ├── main.go (78)     — WASM entry point
│   ├── register.go (18) + register_test.go (24)
├── codegen/
│   ├── codegen.go (90)  + codegen_test.go (54)
├── go.mod
└── Makefile
```

## Conventions

- UK English (colour, organisation)
- `declare(strict_types=1)` equivalent: all types annotated
- Tests: testify assert/require
- Licence: EUPL-1.2
