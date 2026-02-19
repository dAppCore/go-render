# TODO

## High Priority

- [ ] **Fix WASM binary size** — 1.58 MB gzip exceeds 1 MB limit. Move `buildComponentJS()` and JSON parsing to server-side. WASM should only expose `Render()`. Consider: pre-parsed slots, inline template execution, or TinyGo.
- [ ] **Add WASM integration tests** — No `cmd/wasm/main_test.go`. Need JS↔Go round-trip verification.

## Medium Priority

- [ ] **Performance benchmarks** — Add `BenchmarkRender`, `BenchmarkImprint`, `BenchmarkCompareVariants`, `BenchmarkLayout` with varying tree depths.
- [ ] **TypeScript type definitions** — Add `.d.ts` generation alongside `GenerateBundle()` for Web Component consumers.
- [ ] **Accessibility helpers** — Layout has semantic HTML + ARIA roles but no `aria-label` builder, alt text helpers, or focus management nodes.
- [ ] **Layout variant validation** — `NewLayout("XYZ")` silently produces empty output. Add warning or error for invalid slot characters.

## Low Priority

- [ ] **Unicode/RTL edge cases** — Test emoji, RTL text, zero-width characters in Text nodes.
- [ ] **Deep nesting stress test** — Verify performance with deeply nested Layouts and large `Each[T]` lists.
- [ ] **Browser polyfill documentation** — Document closed Shadow DOM support matrix.
- [ ] **CSS scoping helper** — Optional utility for responsive variant CSS targeting.
