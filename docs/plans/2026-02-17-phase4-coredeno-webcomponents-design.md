# Phase 4: CoreDeno + Web Components

**Date:** 2026-02-17
**Status:** Approved (open sections pending dAppServer code review)
**Heritage:** dAppServer prototype (20 repos), Chandler/Dreaming in Code

## Vision

A universal application framework where `.core/view.yml` defines what an app IS.
Run `core` in any directory — it discovers the manifest, verifies its signature,
and boots the application. Like `docker-compose.yml` but for applications.

Philosophical lineage: Mitch Kapor's Chandler (universal configurable app),
rebuilt with Web Components, Deno sandboxing, WASM rendering, and LEM ethics.

## Architecture

```
┌─────────────────────────────────────────────┐
│              WebView2 (Browser)             │
│  ┌───────────┐  ┌──────────┐  ┌──────────┐ │
│  │  Angular   │  │ Web Comp │  │ go-html  │ │
│  │  (shell)   │  │ (modules)│  │  WASM    │ │
│  └─────┬─────┘  └────┬─────┘  └────┬─────┘ │
│        └──────┬───────┘             │       │
│               │ fetch/WS            │       │
└───────────────┼─────────────────────┼───────┘
                │                     │
┌───────────────┼─────────────────────┼───────┐
│         CoreDeno (Deno sidecar)     │       │
│  ┌────────────┴──────────┐    ┌─────┴─────┐ │
│  │  Module Loader        │    │ ITW3→WC   │ │
│  │  + Permission Gates   │    │ Codegen   │ │
│  │  + Dev Server (HMR)   │    │           │ │
│  └────────────┬──────────┘    └───────────┘ │
│               │ gRPC / Unix socket          │
└───────────────┼─────────────────────────────┘
                │
┌───────────────┼─────────────────────────────┐
│         Go Backend (CoreGO)                 │
│  ┌────────┐ ┌┴───────┐ ┌─────────────────┐ │
│  │ Module │ │ gRPC   │ │ MCPBridge       │ │
│  │Registry│ │ Server │ │ (WebView tools) │ │
│  └────────┘ └────────┘ └─────────────────┘ │
└─────────────────────────────────────────────┘
```

Three processes:
- **WebView2**: Angular shell (gradual migration) + Web Components + go-html WASM
- **CoreDeno**: Deno sidecar — module sandbox, I/O fortress, TypeScript toolchain
- **CoreGO**: Framework backbone — lifecycle, services, I/O (core/pkg/io), gRPC server

## Responsibility Split

| Layer | Role |
|-------|------|
| **CoreGO** | Framework (lifecycle, services, I/O via core/pkg/io, module registry, gRPC server) |
| **go-html** | Web Component factory (layout → Shadow DOM, manifest → custom element, WASM client-side registration) |
| **CoreDeno** | Sandbox + toolchain (Deno permissions, TypeScript compilation, dev server, asset serving) |
| **MCPBridge** | Retained for direct WebView tools (window control, display, clipboard, dialogs) |

## CoreDeno Sidecar

### Lifecycle
Go spawns Deno as a managed child process at app startup. Auto-restart on crash.
SIGTERM on app shutdown.

### Communication
- **Channel**: Unix domain socket at `$XDG_RUNTIME_DIR/core/deno.sock`
- **Protocol**: gRPC (proto definitions in `pkg/coredeno/proto/`)
- **Direction**: Bidirectional
  - Deno → Go: I/O requests (file, network, process) gated by permissions
  - Go → Deno: Module lifecycle events, HLCRF re-render triggers

### Deno's Three Roles

**1. Module loader + sandbox**: Reads ITW3 manifests, loads modules with
per-module `--allow-*` permission flags. Modules run in Deno's isolate.

**2. I/O fortress gateway**: All file/network/process I/O routed through Deno's
permission gates before reaching Go via gRPC. A module requesting access outside
its declared paths is denied before Go ever sees the request.

**3. Build/dev toolchain**: TypeScript compilation, module resolution, dev server
with HMR. Replaces Node/npm entirely. In production, pre-compiled bundles
embedded in binary.

### Permission Model
Each module declares required permissions in its manifest:

```yaml
permissions:
  read: ["/data/mining/"]
  write: ["/data/mining/config.json"]
  net: ["pool.lthn.io:3333"]
  run: ["xmrig"]
```

CoreDeno enforces these at the gRPC boundary.

## The .core/ Convention

### Auto-Discovery
Run `core` in any directory. If `.core/view.yml` exists, CoreGO reads it,
validates the ed25519 signature, and boots the application context.

### view.yml Format (successor to .itw3.json)

```yaml
code: photo-browser
name: Photo Browser
version: 0.1.0
sign: <ed25519 signature>

layout: HLCRF
slots:
  H: nav-breadcrumb
  L: folder-tree
  C: photo-grid
  R: metadata-panel
  F: status-bar

permissions:
  read: ["./photos/"]
  net: []
  run: []

modules:
  - core/media
  - core/fs
```

### Signed Application Loading
The `sign` field contains an ed25519 signature. CoreGO verifies before loading.
Unsigned or tampered manifests are rejected. The I/O fortress operates at the
application boundary — the entire app load chain is authenticated.

## Web Component Lifecycle

1. **Discovery** → `core` reads `.core/view.yml`, verifies signature
2. **Resolve** → CoreGO checks module registry for declared components
3. **Codegen** → go-html generates Web Component class definitions from manifest
4. **Permission binding** → CoreDeno wraps component I/O calls with per-module gates
5. **Composition** → HLCRF layout assembles slots, each a custom element with Shadow DOM
6. **Hot reload** → Dev mode: Deno watches files, WASM re-renders affected slots only

### HLCRF Slot Composition

```
┌──────────────────────────────────┐
│ <nav-breadcrumb>    (H - shadow) │
├────────┬───────────────┬─────────┤
│ <folder│ <photo-grid>  │<metadata│
│ -tree> │   (C-shadow)  │ -panel> │
│(L-shad)│               │(R-shad) │
├────────┴───────────────┴─────────┤
│ <status-bar>        (F - shadow) │
└──────────────────────────────────┘
```

Each slot is a custom element with closed Shadow DOM. Isolation by design —
one module cannot reach into another's shadow tree.

### go-html WASM Integration
- **Server-side (Go)**: go-html reads manifests, generates WC class definitions
- **Client-side (WASM)**: go-html WASM in browser dynamically registers custom
  elements at runtime via `customElements.define()`
- Same code, two targets. Server pre-renders for initial load, client handles
  dynamic re-renders when slots change.

## Angular Migration Path

**Phase 4a** (current): Web Components load inside Angular's `<router-outlet>`.
Angular sees custom elements via `CUSTOM_ELEMENTS_SCHEMA`. No Angular code needed
for new modules.

**Phase 4b**: ApplicationFrame becomes a go-html Web Component (HLCRF outer frame).
Angular router replaced by lightweight hash-based router mapping URLs to
`.core/view.yml` slot configurations.

**Phase 4c**: Angular removed. WebView2 loads:
1. go-html WASM (layout engine + WC factory)
2. Thin router (~50 lines)
3. CoreDeno-served module bundles
4. Web Awesome (design system — already vanilla custom elements)

## dAppServer Heritage

20 repos at `github.com/dAppServer/` — the original client-side server concept
and browser↔Go communications bridge. Extract and port, not copy.

| dAppServer repo | CoreDeno equivalent | Action |
|---|---|---|
| `server` | CoreDeno sidecar architecture | Extract |
| `dappui` | Web Component framework | Extract |
| `auth-server` + `mod-auth` | Signed manifest verification | Extract |
| `mod-docker` | Process management (core/pkg/process) | Port |
| `mod-io-process` | I/O fortress gRPC bridge | Extract |
| `server-sdk-typescript-angular` | Deno TypeScript bindings | Port |
| `app-marketplace` | Module registry/discovery | Extract |
| `app-directory-browser` | File browser Web Component | Port |
| `devops` | CI/deployment patterns | Reference |
| `app-utils-cyberchef` | Utility module example | Reference |
| `app-mining` | Mining module (already in ITW3) | Reference |
| `server-sdk-python` | Not needed (Go replaces) | Skip |

## Open Sections (Pending dAppServer Code Review)

These sections will be refined after studying the dAppServer codebase:

### Polyfills
What browser APIs did dAppServer polyfill? Which are still needed in WebView2?
WebView2 is Chromium-based so most modern APIs are available natively.

### Object Store
dAppServer's data persistence model. How does it map to CoreDeno? Options include
IndexedDB in WebView2, SQLite via Deno, or Go-managed storage via gRPC.

### Templated Config Generators
The pattern for generating configs from templates. Informs `.core/view.yml`
codegen pipeline and how complex configurations are composed from simpler parts.

### Git-Based Plugin Marketplace
Module discovery, versioning, and distribution via Git repositories. To be
designed after dAppServer `app-marketplace` code review.

## Deliverables

| Component | Location | Language |
|---|---|---|
| CoreDeno sidecar manager | `core/pkg/coredeno/` | Go |
| gRPC proto definitions | `core/pkg/coredeno/proto/` | Protobuf |
| gRPC server (Go side) | `core/pkg/coredeno/server.go` | Go |
| Deno client runtime | `core-deno/` (new repo) | TypeScript |
| ITW3 → WC codegen | `go-html/codegen/` | Go |
| .core/view.yml loader | `core/pkg/manifest/` | Go |
| Manifest signing/verify | `core/pkg/manifest/sign.go` | Go |
| WASM WC registration | `go-html/cmd/wasm/` (extend) | Go |

## Not In Scope (Future Phases)

- LEM auto-loading from signed manifests
- Module marketplace server infrastructure
- Offline-first sync / object store
- Full Angular removal (Phase 4c)
