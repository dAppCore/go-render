# Findings

## Code Quality

- **53 tests, 100% pass** — excellent coverage ratios across all packages
- **Zero TODOs/FIXMEs** in codebase — clean
- **`go vet` clean** — no static analysis warnings
- **Safe-by-default design** — XSS prevention verified in render_test.go, HTML escaping on all Text nodes, void elements self-close, entitlements deny-by-default

## Architecture Strengths

- Clean minimal API: 9 public constructors + Node interface
- Type-safe generics: `Each[T]` for iteration
- Deterministic output: sorted attributes, reproducible block IDs
- Fluent builder pattern: `NewLayout("HLCRF").H(node).C(node).F(node)`
- Pipeline bridges rendering to privacy layer (GrammarImprint via go-i18n reversal)

## Known Issues

1. **WASM size blocker** — 6.0 MB raw / 1.58 MB gzip. Root cause: stdlib imports (json, encoding, text/template) bloat the WASM binary. Makefile rejects at 1 MB gzip threshold.
2. **No WASM main_test.go** — cmd/wasm/ has register_test.go but no integration test for the JS exports.
3. **Layout accepts invalid variants silently** — `NewLayout("XYZ")` renders nothing, no error returned.
4. **Context.service is private** — Must use `NewContextWithService()`. Limits custom i18n adapter injection.
5. **Responsive only accepts *Layout** — Cannot nest arbitrary nodes in variants, must wrap in Layout first.

## Coverage Gaps

| File | Lines | Tests |
|------|-------|-------|
| node.go | 254 | 206 lines of tests (81%) |
| layout.go | 119 | 116 lines (97%) |
| pipeline.go | 83 | 128 lines (154%) |
| responsive.go | 39 | 89 lines (228%) |
| codegen.go | 90 | 54 lines (60%) |
| cmd/wasm/main.go | 78 | **0 lines (0%)** |
