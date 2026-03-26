---
title: Development Guide
description: How to build, test, and contribute to go-html, including WASM builds, benchmarks, coding standards, and test patterns.
---

# Development Guide

## Prerequisites

- **Go 1.26** or later. The module uses Go 1.26 features (e.g. `range` over integers, `iter.Seq`).
- **go-i18n** cloned alongside this repository at `../go-i18n` relative to the repo root. The `go.mod` `replace` directive points there.
- **go-inference** also resolved via `replace` directive at `../go-inference`. It is an indirect dependency pulled in by `go-i18n`.
- **Go workspace** (`go.work`): this module is part of a shared workspace. Run `go work sync` after cloning.

No additional tools are required for server-side development. WASM builds require the standard Go cross-compilation support (`GOOS=js GOARCH=wasm`), included in all official Go distributions.

## Directory Layout

```
go-html/
  node.go              Node interface and all node types
  layout.go            HLCRF compositor
  pipeline.go          StripTags, Imprint, CompareVariants (!js only)
  responsive.go        Multi-variant breakpoint wrapper
  context.go           Rendering context
  render.go            Render() convenience function
  path.go              ParseBlockID() for data-block path decoding
  codegen/
    codegen.go         Web Component JS generation (server-side)
    codegen_test.go    Tests for codegen
    bench_test.go      Codegen benchmarks
  cmd/
    codegen/
      main.go          Build-time CLI (stdin JSON, stdout JS)
      main_test.go     CLI integration tests
    wasm/
      main.go          WASM entry point (js+wasm build only)
      register.go      buildComponentJS helper (!js only)
      register_test.go Tests for register helper
      size_test.go     WASM binary size gate test (!js only)
  dist/                WASM build output (gitignored)
  docs/                This documentation
    plans/             Phase design documents (historical)
  Makefile             WASM build with size checking
  .core/build.yaml     Build system configuration
```

## Running Tests

```bash
# All tests
go test ./...

# Single test by name
go test -run TestElNode_Render .

# Skip the slow WASM build test
go test -short ./...

# Verbose output
go test -v ./...

# Tests for a specific package
go test ./codegen/
go test ./cmd/codegen/
go test ./cmd/wasm/
```

The WASM size gate test (`TestWASMBinarySize_WithinBudget`) builds the WASM binary as a subprocess. It is slow and is skipped under `-short`. It is also guarded with `//go:build !js` so it cannot run within the WASM environment itself.

### Test Dependencies

Tests use the `testify` library (`assert` and `require` packages). Integration tests and benchmarks that exercise `Text` nodes must initialise the `go-i18n` default service before rendering:

```go
svc, _ := i18n.New()
i18n.SetDefault(svc)
```

The `bench_test.go` file does this in an `init()` function. Individual integration tests do so explicitly.

## Benchmarks

```bash
# All benchmarks
go test -bench . ./...

# Specific benchmark
go test -bench BenchmarkRender_FullPage .

# With memory allocation statistics
go test -bench . -benchmem ./...

# Extended benchmark duration
go test -bench . -benchtime=5s ./...
```

Available benchmark groups:

| Group | Variants |
|-------|----------|
| `BenchmarkRender_*` | Depth 1, 3, 5, 7 element trees; full page with layout |
| `BenchmarkLayout_*` | Content-only, HCF, HLCRF, nested, 50-child slot |
| `BenchmarkEach_*` | 10, 100, 1000 items |
| `BenchmarkResponsive_*` | Three-variant compositor |
| `BenchmarkStripTags_*` | Short and long HTML inputs |
| `BenchmarkImprint_*` | Small and large page trees |
| `BenchmarkCompareVariants_*` | Two and three variant comparison |
| `BenchmarkGenerateClass` | Single Web Component class generation |
| `BenchmarkGenerateBundle_*` | Small (2-slot) and full (5-slot) bundles |
| `BenchmarkTagToClassName` | Kebab-to-PascalCase conversion |
| `BenchmarkGenerateRegistration` | `customElements.define()` call generation |

## WASM Build

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o gohtml.wasm ./cmd/wasm/
```

Strip flags (`-s -w`) are required. Without them, the binary is approximately 50% larger.

The Makefile `wasm` target performs the build and checks the output size:

```bash
make wasm
```

The Makefile enforces a 1 MB gzip transfer limit and a 3 MB raw size limit. Current measured output: approximately 2.90 MB raw, 842 KB gzip.

To verify the gzip size manually:

```bash
gzip -c -9 gohtml.wasm | wc -c
```

## Codegen CLI

The codegen CLI reads a JSON slot map from stdin and writes a Web Component JS bundle to stdout:

```bash
echo '{"H":"site-header","C":"app-content","F":"site-footer"}' \
    | go run ./cmd/codegen/ \
    > components.js
```

JSON keys are HLCRF slot letters (`H`, `L`, `C`, `R`, `F`). Values are custom element tag names (must contain a hyphen per the Web Components specification). Duplicate tag values are deduplicated.

To test the CLI:

```bash
go test ./cmd/codegen/
```

## Static Analysis

```bash
go vet ./...
```

The repository also includes a `.golangci.yml` configuration for `golangci-lint`.

## Coding Standards

### Language

UK English throughout: colour, organisation, centre, behaviour, licence (noun), serialise. American spellings are not used.

### Type Annotations

All exported and unexported functions carry full parameter and return type annotations. The `any` alias is used in preference to `interface{}`.

### HTML Safety

- Use `Text()` for any user-supplied or translated content. It escapes HTML automatically.
- Use `Raw()` only for content you control or have sanitised upstream. Its name explicitly signals "no escaping".
- Never construct HTML by string concatenation in application code.

### Error Handling

Errors are wrapped with context using `fmt.Errorf()`. The codegen package prefixes all errors with `codegen:`.

### Determinism

Output must be deterministic. `El` node attributes are sorted alphabetically before rendering. `map` iteration order in `codegen.GenerateBundle()` may vary across runs -- this is acceptable because Web Component registration order does not affect correctness.

### Build Tags

Files excluded from WASM use `//go:build !js` as the first line, before the `package` declaration. Files compiled only under WASM use `//go:build js && wasm`. The older `// +build` syntax is not used.

The `fmt` package must never be imported in files without a `!js` build tag, as it significantly inflates the WASM binary. Use string concatenation instead of `fmt.Sprintf` in layout and node code.

### Licence

All new files should carry the EUPL-1.2 SPDX identifier:

```go
// SPDX-Licence-Identifier: EUPL-1.2
```

### Commit Format

Conventional commits with lowercase type and optional scope:

```
feat(codegen): add TypeScript type definition generation
fix(wasm): correct slot injection for empty strings
test: add edge case for Unicode surrogate pairs
docs: update architecture with pipeline diagram
```

Include a co-author trailer:

```
Co-Authored-By: Virgil <virgil@lethean.io>
```

## Test Patterns

### Standard Unit Test

```go
func TestElNode_Render(t *testing.T) {
    ctx := NewContext()
    node := El("div", Raw("content"))
    got := node.Render(ctx)
    want := "<div>content</div>"
    if got != want {
        t.Errorf("El(\"div\", Raw(\"content\")).Render() = %q, want %q", got, want)
    }
}
```

### Table-Driven Subtest

```go
func TestStripTags_Unicode(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"emoji in tags", "<span>\U0001F680</span>", "\U0001F680"},
        {"RTL in tags", "<div>\u0645\u0631\u062D\u0628\u0627</div>", "\u0645\u0631\u062D\u0628\u0627"},
        {"CJK in tags", "<p>\u4F60\u597D\u4E16\u754C</p>", "\u4F60\u597D\u4E16\u754C"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := StripTags(tt.input)
            if got != tt.want {
                t.Errorf("StripTags(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}
```

### Integration Test with i18n

```go
func TestIntegration_RenderThenReverse(t *testing.T) {
    svc, _ := i18n.New()
    i18n.SetDefault(svc)
    ctx := NewContext()

    page := NewLayout("HCF").
        H(El("h1", Text("Building project"))).
        C(El("p", Text("Files deleted successfully"))).
        F(El("small", Text("Completed")))

    imp := Imprint(page, ctx)

    if imp.UniqueVerbs == 0 {
        t.Error("reversal found no verbs in rendered page")
    }
}
```

### Codegen Tests with Testify

```go
func TestGenerateClass_ValidTag(t *testing.T) {
    js, err := GenerateClass("photo-grid", "C")
    require.NoError(t, err)
    assert.Contains(t, js, "class PhotoGrid extends HTMLElement")
    assert.Contains(t, js, "attachShadow")
    assert.Contains(t, js, `mode: "closed"`)
}
```

## Known Limitations

- `NewLayout("XYZ")` silently produces empty output for unrecognised slot letters. Valid letters are `H`, `L`, `C`, `R`, `F`. There is no error or warning.
- `Responsive.Variant()` accepts only `*Layout`, not arbitrary `Node` values. Arbitrary subtrees must be wrapped in a single-slot layout first.
- `Context.service` is unexported. Custom translation injection requires `NewContextWithService()`. There is no way to swap the translator after construction.
- The WASM module has no integration test for the JavaScript exports. `size_test.go` tests binary size only; it does not exercise `renderToString` behaviour from JavaScript.
- `codegen.GenerateBundle()` iterates a `map`, so the order of class definitions in the output is non-deterministic. This does not affect correctness but may cause cosmetic diffs between runs.
