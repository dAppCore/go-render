# Phase 4: CoreDeno + Web Components

**Date:** 2026-02-17
**Status:** Approved
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

### Tier 1: Extract (Core Architecture)

| dAppServer repo | What to extract | Phase 4 target |
|---|---|---|
| `server` | Port 36911 bridge, ZeroMQ IPC (pub/sub + req/rep + push/pull), air-gapped PGP auth, object store, 13 test files with RPC procedures | CoreDeno sidecar, I/O fortress, auth |
| `dappui` | Angular→WC migration, REST+WS+Wails triple, terminal (xterm.js) | Web Component framework, MCPBridge |
| `mod-auth` | PGP zero-knowledge auth (sign→encrypt→verify→JWT), QuasiSalt, roles | Signed manifest verification, identity |
| `mod-io-process` | Process registry, 3-layer I/O streaming (process→ZeroMQ→WS→browser) | `core/pkg/process`, event bus |
| `app-marketplace` | Git-as-database registry, category-as-directory, install pipeline | Module registry, `.core/view.yml` loader |

### Tier 2: Port (Useful Patterns)

| dAppServer repo | What to port | Phase 4 target |
|---|---|---|
| `auth-server` | Keycloak + native PGP fallback | External auth option |
| `mod-docker` | Docker socket client, container CRUD (8 ops) | `core/pkg/process` |
| `app-mining` | CLI Bridge (camelCase→kebab-case), Process-as-Service, API proxy | Generic CLI wrapper |
| `app-directory-browser` | Split-pane layout, lazy tree, filesystem CRUD RPCs | `<core-file-tree>` WC |
| `wails-build-action` | Auto-stack detection, cross-platform signing, Deno CI | Build tooling |

### Tier 3: Reference

| dAppServer repo | Value |
|---|---|
| `depends` | Bitcoin Core hermetic build; libmultiprocess + Cap'n Proto validates process-separation |
| `app-utils-cyberchef` | Purest manifest-only pattern ("manifest IS the app") |
| `devops` | Cross-compilation matrix (9 triples), ancestor of ADR-001 |
| `pwa-native-action` | PWA→Wails native shell proof, ancestor of `core-gui` |
| `docker-images` | C++ cross-compile layers |

### Tier 4: Skip

| dAppServer repo | Reason |
|---|---|
| `server-sdk-python` | Auto-generated, Go replaces |
| `server-sdk-typescript-angular` | Auto-generated, superseded |
| `openvpn` | Unmodified upstream fork |
| `ansible-server-base` | Standard Ansible hardening |
| `.github` | Org profile only |

## Polyfills

dAppServer polyfilled nothing at the browser level. The prototype ran inside
Electron/WebView2 (Chromium), which already supports all required APIs natively:
Custom Elements v1, Shadow DOM v1, ES Modules, `fetch`, `WebSocket`,
`customElements.define()`, `structuredClone()`.

**Decision**: No polyfills needed. WebView2 is Chromium-based. The minimum
Chromium version Wails v3 targets already supports all Web Component APIs.

## Object Store

dAppServer used a file-based JSON key-value store at `data/objects/{group}/{object}.json`.
Six operations discovered from test files:

| Operation | dAppServer endpoint | Phase 4 equivalent |
|-----------|--------------------|--------------------|
| Get | `GET /config/object/{group}/{object}` | `store.get(group, key)` |
| Set | `POST /config/object/{group}/{object}` | `store.set(group, key, value)` |
| Clear | `DELETE /config/object/{group}/{object}` | `store.delete(group, key)` |
| Count | `GET /config/object/{group}/count` | `store.count(group)` |
| Remove group | `DELETE /config/object/{group}` | `store.deleteGroup(group)` |
| Render template | `POST /config/render` | `store.render(template, vars)` |

Used for: installed app registry (`conf/installed-apps.json`), menu state
(`conf/menu.json`), per-module config, user preferences.

**Decision**: Go-managed storage via gRPC. CoreGO owns persistence through
`core/pkg/io`. Modules request storage through the I/O fortress — never
touching the filesystem directly. SQLite backend (already a dependency in
the blockchain layer). IndexedDB reserved for client-side cache only.

## Templated Config Generators

dAppServer's `config.render` endpoint accepted a template string + variable map
and returned the rendered output. Used to generate configs for managed processes
(e.g., xmrig config.json from user-selected pool/wallet parameters).

The pattern in Phase 4:
1. Module declares config templates in `.core/view.yml` under a `config:` key
2. User preferences stored in the object store
3. CoreGO renders templates at process-start time via Go `text/template`
4. Rendered configs written to sandboxed paths the module has `write` permission for

Example from mining module (camelCase→kebab-case CLI arg transformation):
```yaml
config:
  xmrig:
    template: conf/xmrig/config.json.tmpl
    vars:
      pool: "{{ .user.pool }}"
      wallet: "{{ .user.wallet }}"
```

**Decision**: Go `text/template` in CoreGO. Templates live in the module's
`.core/` directory. Variables come from the object store. No Deno involvement —
config rendering is a Go-side I/O fortress operation.

## Git-Based Plugin Marketplace

### dAppServer Pattern (Extracted from `app-marketplace`)

A Git repository serves as the package registry. No server infrastructure needed.

```
marketplace/                    # Git repo
├── index.json                  # Root: {version, apps[], dirs[]}
├── miner/
│   └── index.json              # Category: {version, apps[], dirs[]}
├── utils/
│   └── index.json              # Category: {version, apps[], dirs[]}
└── ...
```

Each `index.json` entry points to a raw `.itw3.json` URL in the plugin's own repo:
```json
{"code": "utils-cyberchef", "name": "CyberChef", "type": "bin",
 "pkg": "https://raw.githubusercontent.com/dAppServer/app-utils-cyberchef/main/.itw3.json"}
```

Install pipeline: browse index → fetch manifest → download zip from `app.url` →
extract → run hooks (rename, etc.) → register in object store → add menu entry.

### Phase 4 Evolution

Replace GitHub-specific URLs with Forgejo-compatible Git operations:

1. **Registry**: A Git repo (`host-uk/marketplace`) with category directories
   and `index.json` files. Cloned/pulled by CoreGO at startup and periodically.
2. **Manifests**: Each module's `.core/view.yml` is the manifest (replaces `.itw3.json`).
   The marketplace index points to the module's Git repo, not a raw file URL.
3. **Distribution**: Git clone of the module repo (not zip downloads). CoreGO
   clones into a managed modules directory with depth=1.
4. **Verification**: ed25519 signature in `view.yml` verified before loading.
   The marketplace index includes the expected signing public key.
5. **Install hooks**: Declared in `view.yml` under `hooks:`. Executed by CoreGO
   in the I/O fortress (rename, template render, permission grant).
6. **Updates**: `git pull` on the module repo. Signature re-verified after pull.
   If signature fails, rollback to previous commit.
7. **Discovery**: `core marketplace list`, `core marketplace search <query>`,
   `core marketplace install <code>`.

```yaml
# marketplace/index.json
version: 1
modules:
  - code: utils-cyberchef
    name: CyberChef Data Toolkit
    repo: https://forge.lthn.io/host-uk/mod-cyberchef.git
    sign_key: <ed25519 public key>
    category: utils
categories:
  - miner
  - utils
  - network
```

### RPC Surface (from dAppServer test extraction)

| Operation | CLI | RPC |
|-----------|-----|-----|
| Browse | `core marketplace list` | `marketplace.list(category?)` |
| Search | `core marketplace search <q>` | `marketplace.search(query)` |
| Install | `core marketplace install <code>` | `marketplace.install(code)` |
| Remove | `core marketplace remove <code>` | `marketplace.remove(code)` |
| Installed | `core marketplace installed` | `marketplace.installed()` |
| Update | `core marketplace update <code>` | `marketplace.update(code)` |
| Update all | `core marketplace update` | `marketplace.updateAll()` |

## Complete RPC Surface (Archaeological Extraction)

All procedures discovered from dAppServer test files and controllers:

### Auth
- `auth.create(username, password)` — PGP key generation + QuasiSalt hash
- `auth.login(username, encryptedPayload)` — Zero-knowledge PGP verify → JWT
- `auth.delete(username)` — Remove account

### Crypto
- `crypto.pgp.generateKeyPair(name, email, passphrase)` → {publicKey, privateKey}
- `crypto.pgp.encrypt(data, publicKey)` → encryptedData
- `crypto.pgp.decrypt(data, privateKey, passphrase)` → plaintext
- `crypto.pgp.sign(data, privateKey, passphrase)` → signature
- `crypto.pgp.verify(data, signature, publicKey)` → boolean

### Filesystem
- `fs.list(path, detailed?)` → FileEntry[]
- `fs.read(path)` → content
- `fs.write(path, content)` → boolean
- `fs.delete(path)` → boolean
- `fs.rename(from, to)` → boolean
- `fs.mkdir(path)` → boolean
- `fs.isDir(path)` → boolean
- `fs.ensureDir(path)` → boolean

### Process
- `process.run(command, args, options)` → ProcessHandle
- `process.add(request)` → key
- `process.start(key)` → boolean
- `process.stop(key)` → boolean
- `process.kill(key)` → boolean
- `process.list()` → string[]
- `process.get(key)` → ProcessInfo
- `process.stdout.subscribe(key)` → stream
- `process.stdin.write(key, data)` → void

### Object Store
- `store.get(group, key)` → value
- `store.set(group, key, value)` → void
- `store.delete(group, key)` → void
- `store.count(group)` → number
- `store.deleteGroup(group)` → void
- `store.render(template, vars)` → string

### IPC / Event Bus
- `ipc.pub.subscribe(channel)` → stream
- `ipc.pub.publish(channel, message)` → void
- `ipc.req.send(channel, message)` → response
- `ipc.push.send(message)` → void

### Marketplace
- `marketplace.list(category?)` → ModuleEntry[]
- `marketplace.search(query)` → ModuleEntry[]
- `marketplace.install(code)` → boolean
- `marketplace.remove(code)` → boolean
- `marketplace.installed()` → InstalledModule[]
- `marketplace.update(code)` → boolean

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
- Marketplace server infrastructure (Git-based registry is sufficient)
- Offline-first sync (IndexedDB client cache)
- Full Angular removal (Phase 4c)
- GlusterFS distributed storage (dAppServer aspiration, not needed yet)
- Multi-chain support (Phase 4 is Lethean-only)
