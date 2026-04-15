# Project History

## Phase 1: Core Node Types (initial scaffolding)

Commits: `d7bb0b2` through `c724094`

The module was scaffolded with the Go module path `forge.lthn.ai/core/go-html`. The foundational work established:

- `d7bb0b2` — Module scaffold, `Node` interface with `Render(ctx *Context) string`.
- `3e76e72` — `Text` node wired to `go-i18n` grammar pipeline with HTML escaping.
- `c724094` — Conditional nodes (`If`, `Unless`), entitlement gating (`Entitled`, deny-by-default), runtime dispatch (`Switch`), and type-safe iteration (`Each[T]`).

The `Raw` escape hatch was present from the first commit. The decision to make `Text` always escape and `Raw` never escape was made at this stage and has not changed.

## Phase 2: HLCRF Layout and Pipeline

Commits: `946ea8d` through `ef77793`

- `946ea8d` — `Layout` type with HLCRF slot registry. Semantic HTML elements (`<header>`, `<main>`, `<aside>`, `<footer>`) and ARIA roles assigned per slot.
- `d75988a` — Nested layout path chains. Block IDs computed as `{slot}-0` at root, extended with `{parent}-{slot}-0` for nested layouts.
- `40da0d8` — Deterministic attribute sorting and thread-safe nested layout cloning (clone-on-render pattern).
- `f49ddbf` — `Attr()` helper for setting element attributes with chaining.
- `e041f76` — `Responsive` multi-variant compositor with `data-variant` containers.
- `8ac5123` — `StripTags` single-pass rune scanner for HTML-to-text stripping.
- `76cef5a` — `Imprint()` full render-reverse-imprint pipeline using `go-i18n/reversal`.
- `ef77793` — `CompareVariants()` pairwise semantic similarity scoring across responsive variants.

## Phase 3: WASM Entry Point

Commits: `456adce` through `9bc1fa7`

- `456adce` — Makefile with `wasm` target. Size gate: `WASM_GZ_LIMIT = 1048576` (1 MB). Initial measurement revealed the binary was already too large at this stage.
- `5acf63c` — WASM entry point `cmd/wasm/main.go` with `renderToString` exported to `window.gohtml`.
- `2fab89e` — Integration tests refactored to use `Imprint` pipeline.
- `e34c5c9` — WASM browser test harness added.
- `18d2933` — WASM binary size reporting improvements.
- `9bc1fa7` — Variant name escaping in `Responsive`, single-pass `StripTags` optimisation, WASM security contract documented in source.

## Phase 4: Codegen and Web Components

Commits: `937c08d` through `ab7ab92`

- `37b50ab`, `496513e` — Phase 4 design documents and implementation plan.
- `937c08d` — `codegen` package with `GenerateClass`, `GenerateBundle`, `TagToClassName`. Web Component classes with closed Shadow DOM.
- `dcd55a4` — `registerComponents` export added to `cmd/wasm/main.go`, bridging JSON slot config to WC bundle JS. This was the source of the subsequent binary size problem.
- `ab7ab92` — Transitive `replace` directive added for `go-inference`.

## WASM Binary Size Reduction

Commits: `6abda8b`, `4c65737`, `aae5d21`

The initial WASM binary measured 6.04 MB raw / 1.58 MB gzip — 58% over the 1 MB gzip limit set in the Makefile. The root causes were three heavyweight stdlib imports pulled in by `registerComponents()` in the WASM binary:

| Import | Approx. gzip contribution |
|--------|--------------------------|
| `encoding/json` | ~200 KB |
| `text/template` | ~125 KB |
| `fmt` (via `layout.go`) | ~50 KB |
| `go-i18n/reversal` (via `pipeline.go`) | ~250 KB |

**Total bloat**: ~625 KB gzip over the core rendering requirement.

The fix was applied in three distinct steps:

### Step 1: Remove registerComponents from WASM (`4c65737`)

`cmd/wasm/register.go` received a `//go:build !js` build tag, completely excluding it from the WASM compilation unit. The `registerComponents` entry on the `gohtml` JS object was removed from `cmd/wasm/main.go`. The codegen function was moved to a standalone build-time CLI at `cmd/codegen/main.go`. This eliminated `encoding/json` and `text/template` from the WASM import graph.

### Step 2: Remove pipeline from WASM

`pipeline.go` received a `//go:build !js` build tag. The `Imprint()` and `CompareVariants()` functions depend on `go-i18n/reversal`, which is a heavyweight analysis library. These functions are server-side analysis tools and have no use in a client-side rendering module. The `renderToString` function in the WASM entry point never called them, so removal was non-breaking.

### Step 3: Eliminate fmt from WASM

`layout.go`'s `blockID()` method had used `fmt.Sprintf` for string construction. Replacing this with direct string concatenation (`l.path + string(slot) + "-0"`) removed `fmt` from the WASM import graph entirely.

**Result**: 2.90 MB raw, 842 KB gzip. 47% reduction in gzip size. Well within the 1 MB limit.

### Size gate test (`aae5d21`)

`cmd/wasm/size_test.go` was added to prevent regression. `TestWASMBinarySize_WithinBudget` builds the WASM binary in a temp directory, gzip-compresses it, and asserts:

- Gzip size < 1,048,576 bytes (1 MB).
- Raw size < 3,145,728 bytes (3 MB).

The test is skipped under `go test -short` and is guarded with `//go:build !js`.

## Test Coverage Milestones

- `7efd2ab` — Benchmarks added across all subsystems. Unicode edge case tests. Stress tests.
- `ab7ab92` — 53 passing tests across the package and sub-packages.
- `aae5d21` — 70+ tests passing (server-side); WASM size gate and codegen CLI tests added.

## Known Limitations (as of current HEAD)

These are not regressions; they are design choices or deferred work recorded for future consideration.

1. **Invalid layout variants are still non-fatal at render time.** `NewLayout("XYZ")` produces empty output, but `VariantError()` exposes the validation result without changing the `NewLayout` API.

2. **No WASM integration test.** `cmd/wasm/size_test.go` tests binary size only. The `renderToString` behaviour is tested by building and running the WASM binary in a browser, not by an automated test. A `syscall/js`-compatible test harness would be needed.

3. **Responsive accepts only Layout.** `Responsive.Variant()` takes `*Layout` rather than `Node`. The rationale is that `CompareVariants` in the pipeline needs access to the slot structure. Accepting `Node` would require a different approach to variant analysis.

4. **Context.service is private.** The i18n service is still unexported, but callers can now swap it explicitly with `Context.SetService()`. This keeps the field encapsulated while allowing controlled mutation.

5. **TypeScript definitions are generated.** `codegen.GenerateTypeScriptDefinitions()` and the `cmd/codegen -types` flag emit `.d.ts` companions for generated Web Components.

6. **CSS scoping helper added.** `VariantSelector(name)` returns a reusable `data-variant` attribute selector for stylesheet targeting. The `Responsive` rendering model remains unchanged.

7. **Browser polyfill matrix not documented.** Closed Shadow DOM is well-supported but older browsers require polyfills. The support matrix is not documented.

## Future Considerations

These items were captured during the WASM size reduction work and expert review sessions. They are not committed work items.

- **TypeScript type definitions** alongside `GenerateBundle()` for typed Web Component consumers.
- **Accessibility helpers** — `aria-label` builder, `alt` text helpers, and focus management helpers (`TabIndex`, `AutoFocus`). The layout has semantic HTML and ARIA roles but no API for fine-grained accessibility attributes beyond `Attr()`.
- **Responsive CSS helpers** — `VariantSelector(name)` makes `data-variant` targeting explicit and reusable in stylesheets.
- **Layout variant validation** — return a warning or sentinel error from `NewLayout` when the variant string contains unrecognised slot characters.
- **Daemon mode for codegen** — watch mode for regenerating the JS bundle when slot config changes, for development workflows.
