# TODO

## High Priority — WASM Binary Size Fix

~~Current: 6.04 MB raw / 1.58 MB gzip.~~ Fixed: 2.90 MB raw / 830 KB gzip. Target: < 1 MB gzip (Makefile gate: WASM_GZ_LIMIT 1048576).

Root cause: `registerComponents()` pulled in `encoding/json` (~200KB gz), `text/template` (~125KB gz), and `fmt` (~50KB gz). Plus `pipeline.go` linked `go-i18n/reversal` (~250KB gz). These were heavyweight imports for code that doesn't need to run client-side.

### Step 1: Remove `registerComponents()` from WASM

- [x] **Move `cmd/wasm/register.go` out of WASM** — Added `//go:build !js` build tag. The `registerComponents()` JS bridge in `main.go` removed. This removes `encoding/json` and `text/template` from the binary.
- [x] **Move codegen to build-time CLI** — Created `cmd/codegen/main.go` that reads slot config from stdin (JSON) and writes generated JS to stdout. Usage: `echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js`. Consumers pre-generate during build.
- [x] **Update `cmd/wasm/main.go`** — Removed `registerComponents` from the `gohtml` JS object. Only exposes `renderToString`.

### Step 2: Remove Pipeline from WASM

- [x] **Guard `pipeline.go` with build tag** — Added `//go:build !js` to `pipeline.go`. The `Imprint()` and `CompareVariants()` functions use `go-i18n/reversal` which is heavyweight. Server-side analysis only.
- [x] **Update `cmd/wasm/main.go`** — No references to `pipeline.go` functions. `renderToString` never used them.

### Step 3: Minimise `fmt` Usage

- [x] **Replace `fmt.Sprintf` in WASM-linked code** — Replaced `fmt.Sprintf` in `layout.go` `blockID()` with string concatenation. `fmt` eliminated from the WASM import graph.

### Step 4: Verify Size

- [x] **Build and measure** — 2,900,777 bytes raw, 830,314 bytes gzip (842,146 via `make wasm`). Well under 1 MB limit.
- [x] **Document the server/client split** — Updated CLAUDE.md with new architecture: WASM = `renderToString()` only, codegen = build-time CLI.

### Step 5: Tests

- [x] **WASM build gate test** — `TestWASMBinarySize` in `cmd/wasm/size_test.go`: builds WASM, gzips, asserts < 1MB gzip and < 3MB raw. Result: 2.90MB raw, 842KB gzip. `//go:build !js` guarded.
- [x] **Codegen CLI test** — `cmd/codegen/main_test.go`: pipe JSON stdin, verify JS output matches `GenerateBundle()`
- [x] **renderToString still works** — Existing WASM tests for `renderToString` pass (build-tag guarded)
- [x] **Existing tests still pass** — `go test ./...` (non-WASM) all 70+ tests pass, pipeline/codegen tests unaffected

## Medium Priority

- [ ] **TypeScript type definitions** — Add `.d.ts` generation alongside `GenerateBundle()` for Web Component consumers.
- [ ] **Accessibility helpers** — Layout has semantic HTML + ARIA roles but no `aria-label` builder, alt text helpers, or focus management nodes.
- [ ] **Layout variant validation** — `NewLayout("XYZ")` silently produces empty output. Add warning or error for invalid slot characters.

## Low Priority

- [ ] **Browser polyfill documentation** — Document closed Shadow DOM support matrix.
- [ ] **CSS scoping helper** — Optional utility for responsive variant CSS targeting.
