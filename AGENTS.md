# go-html Agent Guide

This repository is the semantic HTML renderer for the Dappcore Go stack. It
builds composable `Node` trees, HLCRF layouts, responsive variants, grammar
imprints, Web Component code generation, a small WASM bridge, and a Gin-backed
API provider.

## Structure

The root package `dappco.re/go/html` contains the runtime renderer. `node.go`
defines the renderable primitives, `layout.go` owns HLCRF slot composition,
`responsive.go` wraps layout variants, `pipeline.go` connects rendering to the
i18n reversal grammar imprint pipeline, and `context.go` carries locale,
metadata, and translator state. `shadow.go` generates static Web Component
classes from node trees.

The `codegen/` package is the build-time Web Component generator used by
`cmd/codegen/`. The command reads slot maps as JSON and writes JavaScript or
TypeScript definitions. `cmd/wasm/` contains the browser-facing WASM entrypoint
and compatibility layout renderer. `pkg/api/` exposes the render and grammar
checks through the Core provider shape without importing the full provider
runtime.

## Local Rules

Do not edit `third_party/`, `.git/`, `.codex/`, or `BRIEF.md`. The v0.9.0
compliance audit is the source of truth for migration work:

```sh
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```

Tests and examples are file-aware. A public symbol in `foo.go` needs its
`TestFoo_<Symbol>_{Good,Bad,Ugly}` triplet in `foo_test.go` and its
`Example<Symbol>` usage in `foo_example_test.go`. Do not create AX7, versioned,
or monolithic test files.

Direct imports of banned stdlib packages such as `fmt`, `errors`, `strings`,
`path`, `os`, `log`, `encoding/json`, and `bytes` are not accepted in any Go
file. Use `dappco.re/go` wrappers directly, or keep WASM-sized local helpers
small and explicit when the runtime file intentionally avoids the Core import.

## Verification

Before handing work back, run the full repository gate:

```sh
GOWORK=off go mod tidy
GOWORK=off go vet ./...
GOWORK=off go test -count=1 ./...
gofmt -l .
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```

The audit must print `verdict: COMPLIANT` with every counter at zero. Passing
unit tests alone is not a completed compliance pass.
