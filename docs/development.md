# Development Guide

## Prerequisites

- Go 1.25 or later (Go workspace required).
- `go-i18n` repository cloned alongside this one: `../go-i18n` relative to the repository root. The `go.mod` `replace` directive points there.
- `go-inference` also resolved via `replace` directive at `../go-inference`. It is an indirect dependency pulled in by `go-i18n`.
- `testify` is the only external test dependency; it is fetched by the Go module system.

No additional tools are required for server-side development. WASM builds require a standard Go installation with `GOOS=js GOARCH=wasm` cross-compilation support, which is included in all official Go distributions.

## Directory Layout

```
go-html/
├── node.go              Node interface and all node types
├── layout.go            HLCRF compositor
├── pipeline.go          StripTags, Imprint, CompareVariants (!js only)
├── responsive.go        Multi-variant breakpoint wrapper
├── context.go           Rendering context
├── render.go            Render() convenience function
├── path.go              ParseBlockID() for data-block path decoding
├── codegen/
│   └── codegen.go       Web Component JS generation (server-side)
├── cmd/
│   ├── codegen/
│   │   └── main.go      Build-time CLI (stdin JSON → stdout JS)
│   └── wasm/
│       ├── main.go      WASM entry point (js+wasm build only)
│       ├── register.go  buildComponentJS helper (!js only)
│       └── size_test.go WASM binary size gate test (!js only)
└── docs/
    └── plans/           Phase design documents (historical)
```

## Running Tests

```bash
# All tests
go test ./...

# Single test by name
go test -run TestWASMBinarySize_Good ./cmd/wasm/

# Skip slow WASM build test
go test -short ./...

# Tests with verbose output
go test -v ./...
```

Tests use `testify` assert and require helpers. Test names follow Go's standard `TestFunctionName` convention. Subtests use `t.Run()` with descriptive names.

The WASM size gate test (`TestWASMBinarySize_Good`) builds the WASM binary as a subprocess and is therefore slow. It is skipped automatically under `-short`. It is also guarded with `//go:build !js` so it cannot run under `GOARCH=wasm`.

## Benchmarks

```bash
# All benchmarks
go test -bench . ./...

# Specific benchmark
go test -bench BenchmarkRender_FullPage ./...

# With memory allocations
go test -bench . -benchmem ./...

# Fixed iteration count
go test -bench . -benchtime=5s ./...
```

Benchmarks are organised by operation:

| Group | Variants |
|-------|---------|
| `BenchmarkRender_*` | Depth 1, 3, 5, 7 trees; full page |
| `BenchmarkLayout_*` | Content-only, HCF, HLCRF, nested, many children |
| `BenchmarkEach_*` | 10, 100, 1000 items |
| `BenchmarkResponsive_*` | Three-variant compositor |
| `BenchmarkStripTags_*` | Short and long HTML inputs |
| `BenchmarkImprint_*` | Small and large page trees |
| `BenchmarkCompareVariants_*` | Two and three variant comparison |

## WASM Build

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o gohtml.wasm ./cmd/wasm/
```

Strip flags (`-s -w`) are required. Without them the binary is approximately 50% larger.

The Makefile target `make wasm` performs the build and measures the gzip size:

```bash
make wasm
```

The Makefile enforces a 1 MB gzip limit (`WASM_GZ_LIMIT = 1048576`). The build fails if this limit is exceeded.

To verify the size manually:

```bash
gzip -c -9 gohtml.wasm | wc -c
```

Current measured output: 2.90 MB raw, 842 KB gzip.

## Codegen CLI

The codegen CLI reads a JSON slot map from stdin and writes a Web Component JS bundle to stdout. It is a build-time tool, not intended for runtime use.

```bash
# Generate components for a two-slot layout
echo '{"H":"site-header","C":"app-content","F":"site-footer"}' \
    | go run ./cmd/codegen/ \
    > components.js
```

The JSON keys are HLCRF slot letters (`H`, `L`, `C`, `R`, `F`). The values are custom element tag names (must contain a hyphen). Duplicate tag values are deduplicated.

To test the CLI:

```bash
go test ./cmd/codegen/
```

## Static Analysis

```bash
go vet ./...
```

The codebase passes `go vet` with no warnings.

## Coding Standards

### Language

UK English throughout: colour, organisation, centre, behaviour, licence (noun), serialise. American spellings are not used.

### Types

All exported and unexported functions carry full parameter and return type annotations. The `any` alias is used in preference to `interface{}`.

### Error handling

Errors are wrapped with context using `fmt.Errorf("pkg.Function: %w", err)`. The codegen package prefixes all errors with `codegen:`.

### HTML safety

- Use `Text()` for any user-supplied or translated content. It escapes HTML.
- Use `Raw()` only for content you control or have sanitised upstream.
- Never construct HTML by string concatenation in application code.

### Determinism

Output must be deterministic. Attributes are sorted before rendering. `map` iteration in `codegen.GenerateBundle()` may produce non-deterministic class order across runs — this is acceptable because Web Component registration order does not affect correctness.

### Build tags

Files excluded from WASM use `//go:build !js` as the first line, before the `package` declaration. Files compiled only under WASM use `//go:build js && wasm`. Do not use the older `// +build` syntax.

### Licence

All files carry the EUPL-1.2 SPDX identifier:

```go
// SPDX-Licence-Identifier: EUPL-1.2
```

### Commit format

Conventional commits with lowercase type and optional scope:

```
feat(codegen): add TypeScript type definition generation
fix(wasm): correct slot injection for empty strings
test: add edge case for Unicode surrogate pairs
docs: update architecture with pipeline diagram
```

Commits include a co-author trailer:

```
Co-Authored-By: Virgil <virgil@lethean.io>
```

## Test Patterns

### Standard unit test

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

### Table-driven subtest

```go
func TestStripTags(t *testing.T) {
    cases := []struct {
        name  string
        input string
        want  string
    }{
        {"empty", "", ""},
        {"plain", "hello", "hello"},
        {"single tag", "<p>hello</p>", "hello"},
        {"nested", "<div><p>a</p><p>b</p></div>", "a b"},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := StripTags(tc.input)
            if got != tc.want {
                t.Errorf("StripTags(%q) = %q, want %q", tc.input, got, tc.want)
            }
        })
    }
}
```

### Integration test with i18n

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

Integration tests that exercise the full pipeline (`Imprint`, `CompareVariants`) must initialise the i18n default service before calling `Text` nodes. The `bench_test.go` `init()` function does this for benchmarks; individual integration tests must do so explicitly.

## Known Limitations

- `NewLayout("XYZ")` silently produces empty output when given unrecognised slot letters. There is no warning or error. Valid slot letters are `H`, `L`, `C`, `R`, `F`.
- `Responsive.Variant()` accepts only `*Layout`, not arbitrary `Node` values. Arbitrary subtrees must be wrapped in a single-slot layout.
- `Context.service` is private. Custom i18n adapter injection requires `NewContextWithService()`. There is no way to set or swap the service after construction.
- `cmd/wasm/main.go` has no integration test for the JS exports. The `size_test.go` file tests binary size only; it does not exercise `renderToString` behaviour.
