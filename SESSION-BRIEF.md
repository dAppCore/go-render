# Session Brief: core/go-html

**Repo**: `forge.lthn.ai/core/go-html` (clone at `/tmp/core-go-html`)
**Module**: `forge.lthn.ai/core/go-html`
**Status**: Current tests pass; WASM build is within budget and codegen emits JS plus TypeScript defs
**Wiki**: https://forge.lthn.ai/core/go-html/wiki (6 pages)

## What This Is

HLCRF DOM compositor with grammar pipeline. Renders semantic HTML from composable node trees with:
- **Node interface**: El, Text, Raw, If, Unless, Each[T], Switch, Entitled
- **HLCRF Layout**: Header/Left/Content/Right/Footer with ARIA roles
- **Responsive**: Multi-variant breakpoint rendering
- **Pipeline**: Render → strip tags → tokenise via go-i18n/reversal → GrammarImprint
- **WASM target**: `cmd/wasm/` exposes `renderToString()` to JS
- **Codegen**: Web Component classes with closed Shadow DOM plus `.d.ts` generation

## Current State

| Area | Status |
|------|--------|
| Core (node, layout, responsive, pipeline) | SOLID — all tested, clean API |
| Tests | Passing |
| go vet | Clean |
| TODOs/FIXMEs | None |
| WASM build | PASS — within the 1 MB gzip gate |
| Codegen | Working — generates WC classes and `.d.ts` definitions |

## Dependencies

- `forge.lthn.ai/core/go-i18n` (replace directive → `../go-i18n`)
- `github.com/stretchr/testify` v1.11.1
- `golang.org/x/text` v0.33.0

## Priority Work

No active blockers recorded here. See `docs/history.md` for the remaining design choices and deferred ideas that were captured during earlier implementation phases.

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
