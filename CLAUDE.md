# CLAUDE.md

Agent instructions for **go-render** (formerly go-html / go-help): a terminal-UI
**display server**. Library module `dappco.re/go/html` lives at `go/`; the binaries live
at `cli/` (module `dappco.re/go/html/cli`, a sibling — Library-as-package vs Instantiated
code). The module path stays `dappco.re/go/html` until the `go/render` graduation.

## Layout (a display server)

```
go/                     library (dappco.re/go/html)
├── engine/             produce renderable content
│   ├── html/           node model · HLCRF layout · terminal render (RenderTerm / RenderTermBoxes) · WASM · pipeline · grammar
│   ├── ctml/           the .ctml parser: ctml.Parse(src) -> html.Node
│   ├── teabox/         (x,y) -> node hit-test (resolve against a BoxMap)
│   └── codegen/        Web Component bundle generation from slot maps
├── display/            paint to a surface
│   ├── tui/            charm-free tui layer: widget re-exports + the tui manager (tui.Run / tui.App)
│   └── ctmltest/       the _test.ctml harness (in-process render + Snapshot / Image backends)
└── pkg/api/            standalone HTTP handlers/provider
cli/                    instantiated binaries (dappco.re/go/html/cli): codegen · termdemo · wasm
# host/{win,pi,riscv} — later
```

**Architecture:** the `Node` interface (`Render(ctx *Context) string` — El, Text, Raw,
If, Each[T], Switch, Entitled). **HLCRF** = the Header/Left/Content/Right/Footer
compositor (variant strings like "HCF"/"HLCRF"/"C", deterministic `data-block` IDs).
`.ctml` is the declarative markup; `engine/ctml.Parse` yields a node tree the terminal
renderer or the WASM/web path paints.

## Build / test

```bash
cd go && GOWORK=off go test ./...      # library — deps resolve from go.mod TAGS
cd go && GOWORK=off go vet ./...
go build ./cli/...                     # binaries — the go.work (use ./go ./cli) resolves the local library
GOOS=js GOARCH=wasm go build ./cli/wasm/   # WASM build (via the workspace)
```

**go.work holds the repo's own modules only** (`./go ./cli`); every other dependency
resolves from a released tag. **No external submodules in go.work** (banned) and **no
`replace` in go.mod**.

## The tui manager (display/tui)

`tui.Run(tui.NewApp(node))` turns a parsed `.ctml` into a live, cross-platform terminal
screen — the "active thing" over the passive engines. `tui.App` is a root Bubble Tea
Model (render the `.ctml` via `engine/html.RenderTerm`, track window size, set
altscreen + mouse on the v2 `View`, quit on ctrl+c/q). The `tui/*` subpackages re-export
charmbracelet (`charm.land/*`) so consumers import `dappco.re/go/html/display/tui/*` and
never touch charmbracelet directly — harnesses included.

## WASM constraint (engine/html)

Files WITHOUT a `//go:build !js` tag are WASM-linked — **never** import `encoding/json`,
`text/template`, or `fmt` in them (string concatenation instead; `fmt` alone adds ~500 KB).
The terminal renderer (`engine/html/term*.go`, lipgloss) and `pipeline.go` (go-i18n
reversal) are `!js`-tagged, server-side only. WASM size gate: < 3.5 MB raw, < 1 MB gzip.

## Coding standards

- UK English (colour, organisation, behaviour, licence, serialise).
- Licence EUPL-1.2 — `// SPDX-Licence-Identifier: EUPL-1.2` on new files.
- Errors via `core.E` / `log.E`, never `fmt.Errorf`. File I/O via `coreio.Local` (go-io),
  never `os.ReadFile` / `os.WriteFile`.
- All types annotated; `any` not `interface{}`. Deterministic output (sorted attributes,
  reproducible block IDs). Safe-by-default (HTML escaping, deny-by-default entitlement).
- One test per symbol; match neighbouring files' test style (plain stdlib `testing` in
  `display/tui`, testify where existing files use it). Ship an `Example` + `// Output:`
  with a public feature. Integration tests using `Text` nodes init i18n first
  (`svc,_ := i18n.New(); i18n.SetDefault(svc)`).
- Commits: conventional + `Co-Authored-By: Virgil <virgil@lethean.io>`. Go 1.26+.
```

