# TODO

## High Priority — WASM Binary Size Fix

Current: 6.04 MB raw / 1.58 MB gzip. Target: < 1 MB gzip (Makefile gate: WASM_GZ_LIMIT 1048576).

Root cause: `registerComponents()` pulls in `encoding/json` (~200KB gz), `text/template` (~125KB gz), and `fmt` (~50KB gz). Plus `pipeline.go` links `go-i18n/reversal` (~250KB gz). These are heavyweight imports for code that doesn't need to run client-side.

### Step 1: Remove `registerComponents()` from WASM

- [ ] **Move `cmd/wasm/register.go` out of WASM** — Add `//go:build !js` build tag OR delete it from `cmd/wasm/`. The `registerComponents()` JS bridge in `main.go` must also be removed. This removes `encoding/json` and `text/template` from the binary.
- [ ] **Move codegen to build-time CLI** — Create `cmd/codegen/main.go` that reads slot config from stdin (JSON) and writes generated JS to stdout. This replaces the WASM-based registration. Usage: `echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js`. Consumers pre-generate during build.
- [ ] **Update `cmd/wasm/main.go`** — Remove `registerComponents` from the `gohtml` JS object. Only expose `renderToString`. Remove the `encoding/json` and `codegen` imports.

### Step 2: Remove Pipeline from WASM

- [ ] **Guard `pipeline.go` with build tag** — Add `//go:build !js` to `pipeline.go`. The `Imprint()` and `CompareVariants()` functions use `go-i18n/reversal` which is heavyweight. These are server-side analysis functions, not needed in browser rendering.
- [ ] **Update `cmd/wasm/main.go`** — Ensure no references to `pipeline.go` functions. Currently `renderToString` doesn't use them, so this should be clean.

### Step 3: Minimise `fmt` Usage

- [ ] **Replace `fmt.Errorf` in WASM-linked code** — In any source files compiled into WASM (node.go, layout.go, responsive.go, context.go, render.go), replace `fmt.Errorf("...: %w", err)` with `errors.New("...")` or manual string concatenation where wrapping isn't needed. Goal: eliminate `fmt` from the WASM import graph entirely if possible.

### Step 4: Verify Size

- [ ] **Build and measure** — Run `GOOS=js GOARCH=wasm go build -o gohtml.wasm ./cmd/wasm/` then `gzip -9 -c gohtml.wasm | wc -c`. Must be < 1,048,576 bytes. Update Makefile if the gate passes.
- [ ] **Document the server/client split** — Update CLAUDE.md with the new architecture: WASM = `renderToString()` only, codegen = build-time CLI.

### Step 5: Tests

- [ ] **WASM build gate test** — `TestWASMBinarySize` in `cmd/wasm/main_test.go`: build WASM, gzip, assert < 1MB
- [ ] **Codegen CLI test** — `cmd/codegen/main_test.go`: pipe JSON stdin → verify JS output matches `GenerateBundle()`
- [ ] **renderToString still works** — Existing WASM tests for `renderToString` pass (may need JS runtime like `wasmedge` or build-tag guarded)
- [ ] **Existing tests still pass** — `go test ./...` (non-WASM) still passes, pipeline/codegen tests unaffected

## Medium Priority

- [ ] **TypeScript type definitions** — Add `.d.ts` generation alongside `GenerateBundle()` for Web Component consumers.
- [ ] **Accessibility helpers** — Layout has semantic HTML + ARIA roles but no `aria-label` builder, alt text helpers, or focus management nodes.
- [ ] **Layout variant validation** — `NewLayout("XYZ")` silently produces empty output. Add warning or error for invalid slot characters.

## Low Priority

- [ ] **Browser polyfill documentation** — Document closed Shadow DOM support matrix.
- [ ] **CSS scoping helper** — Optional utility for responsive variant CSS targeting.
