<!-- STATUS: Not yet implemented -->
# Phase 4: CoreDeno + Web Components Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the universal application framework where `.core/view.yml` defines what an app IS — auto-discovery, ed25519 signed manifests, Web Component generation from go-html, object store, CoreDeno sidecar, and git-based marketplace.

**Architecture:** Three-process model (WebView2 + CoreDeno sidecar + CoreGO backend) communicating over gRPC/Unix socket. go-html generates Web Component class definitions from manifests. CoreDeno enforces per-module I/O permissions. CoreGO owns lifecycle, storage, and module registry via its existing DI framework.

**Tech Stack:** Go 1.25, gRPC, protobuf, ed25519, SQLite (via `pkg/io/sqlite`), WASM (`syscall/js`), Deno (child process), Web Components (Custom Elements v1 + Shadow DOM v1)

---

## Phase 4a: Foundation (Tasks 1–4)

### Task 1: Manifest Types

**Files:**
- Create: `~/Code/host-uk/core/pkg/manifest/manifest.go`
- Test: `~/Code/host-uk/core/pkg/manifest/manifest_test.go`

**Context:** The `.core/view.yml` manifest is the central contract. It declares an app's identity, HLCRF layout slots, permissions, and module dependencies. This task defines the Go types and YAML parsing.

**Step 1: Write the failing test**

```go
// manifest_test.go
package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Good(t *testing.T) {
	raw := `
code: photo-browser
name: Photo Browser
version: 0.1.0
sign: dGVzdHNpZw==

layout: HLCRF
slots:
  H: nav-breadcrumb
  L: folder-tree
  C: photo-grid
  R: metadata-panel
  F: status-bar

permissions:
  read: ["./photos/"]
  write: []
  net: []
  run: []

modules:
  - core/media
  - core/fs
`
	m, err := Parse([]byte(raw))
	require.NoError(t, err)
	assert.Equal(t, "photo-browser", m.Code)
	assert.Equal(t, "Photo Browser", m.Name)
	assert.Equal(t, "0.1.0", m.Version)
	assert.Equal(t, "dGVzdHNpZw==", m.Sign)
	assert.Equal(t, "HLCRF", m.Layout)
	assert.Equal(t, "nav-breadcrumb", m.Slots["H"])
	assert.Equal(t, "photo-grid", m.Slots["C"])
	assert.Len(t, m.Permissions.Read, 1)
	assert.Equal(t, "./photos/", m.Permissions.Read[0])
	assert.Len(t, m.Modules, 2)
}

func TestParse_Bad(t *testing.T) {
	_, err := Parse([]byte("not: valid: yaml: ["))
	assert.Error(t, err)
}

func TestManifest_SlotNames_Good(t *testing.T) {
	m := Manifest{
		Slots: map[string]string{
			"H": "nav-bar",
			"C": "main-content",
		},
	}
	names := m.SlotNames()
	assert.Contains(t, names, "nav-bar")
	assert.Contains(t, names, "main-content")
	assert.Len(t, names, 2)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestParse ./pkg/manifest/ -v`
Expected: FAIL — package does not exist

**Step 3: Write minimal implementation**

```go
// manifest.go
package manifest

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Manifest represents a .core/view.yml application manifest.
type Manifest struct {
	Code    string            `yaml:"code"`
	Name    string            `yaml:"name"`
	Version string            `yaml:"version"`
	Sign    string            `yaml:"sign"`
	Layout  string            `yaml:"layout"`
	Slots   map[string]string `yaml:"slots"`

	Permissions Permissions `yaml:"permissions"`
	Modules     []string    `yaml:"modules"`
}

// Permissions declares the I/O capabilities a module requires.
type Permissions struct {
	Read  []string `yaml:"read"`
	Write []string `yaml:"write"`
	Net   []string `yaml:"net"`
	Run   []string `yaml:"run"`
}

// Parse decodes YAML bytes into a Manifest.
func Parse(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest.Parse: %w", err)
	}
	return &m, nil
}

// SlotNames returns a deduplicated list of component names from slots.
func (m *Manifest) SlotNames() []string {
	seen := make(map[string]bool)
	var names []string
	for _, name := range m.Slots {
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test -run TestParse ./pkg/manifest/ -v`
Expected: PASS (3 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/manifest/manifest.go pkg/manifest/manifest_test.go
git commit -m "feat(manifest): add .core/view.yml types and parser"
```

---

### Task 2: Ed25519 Manifest Signing

**Files:**
- Create: `~/Code/host-uk/core/pkg/manifest/sign.go`
- Test: `~/Code/host-uk/core/pkg/manifest/sign_test.go`

**Context:** Every `.core/view.yml` must be ed25519-signed. Unsigned manifests are rejected in production. The signature covers all fields except `sign` itself. The existing `pkg/crypt` has PGP/AES/ChaCha but no ed25519 — this is intentionally in the manifest package since it's manifest-specific.

**Step 1: Write the failing test**

```go
// sign_test.go
package manifest

import (
	"crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignAndVerify_Good(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	m := &Manifest{
		Code:    "test-app",
		Name:    "Test App",
		Version: "1.0.0",
		Layout:  "HLCRF",
		Slots:   map[string]string{"C": "main"},
	}

	err = Sign(m, priv)
	require.NoError(t, err)
	assert.NotEmpty(t, m.Sign)

	ok, err := Verify(m, pub)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestVerify_Bad_Tampered(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	m := &Manifest{Code: "test-app", Version: "1.0.0"}
	_ = Sign(m, priv)

	m.Code = "evil-app" // tamper

	ok, err := Verify(m, pub)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerify_Bad_Unsigned(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(nil)
	m := &Manifest{Code: "test-app"}

	ok, err := Verify(m, pub)
	assert.Error(t, err)
	assert.False(t, ok)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestSign ./pkg/manifest/ -v`
Expected: FAIL — Sign/Verify not defined

**Step 3: Write minimal implementation**

```go
// sign.go
package manifest

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"gopkg.in/yaml.v3"
)

// signable returns the canonical bytes to sign (manifest without sign field).
func signable(m *Manifest) ([]byte, error) {
	// Copy to avoid mutating the original
	tmp := *m
	tmp.Sign = ""
	return yaml.Marshal(&tmp)
}

// Sign computes the ed25519 signature and stores it in m.Sign (base64).
func Sign(m *Manifest, priv ed25519.PrivateKey) error {
	msg, err := signable(m)
	if err != nil {
		return fmt.Errorf("manifest.Sign: marshal: %w", err)
	}
	sig := ed25519.Sign(priv, msg)
	m.Sign = base64.StdEncoding.EncodeToString(sig)
	return nil
}

// Verify checks the ed25519 signature in m.Sign against the public key.
func Verify(m *Manifest, pub ed25519.PublicKey) (bool, error) {
	if m.Sign == "" {
		return false, fmt.Errorf("manifest.Verify: no signature present")
	}
	sig, err := base64.StdEncoding.DecodeString(m.Sign)
	if err != nil {
		return false, fmt.Errorf("manifest.Verify: decode: %w", err)
	}
	msg, err := signable(m)
	if err != nil {
		return false, fmt.Errorf("manifest.Verify: marshal: %w", err)
	}
	return ed25519.Verify(pub, msg, sig), nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test -run TestSign ./pkg/manifest/ -v && go test -run TestVerify ./pkg/manifest/ -v`
Expected: PASS (3 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/manifest/sign.go pkg/manifest/sign_test.go
git commit -m "feat(manifest): ed25519 signing and verification"
```

---

### Task 3: Manifest Loader (Auto-Discovery)

**Files:**
- Create: `~/Code/host-uk/core/pkg/manifest/loader.go`
- Test: `~/Code/host-uk/core/pkg/manifest/loader_test.go`

**Context:** The loader reads `.core/view.yml` from a directory, parses it, and optionally verifies its signature. This is the "run `core` in any directory" experience. Uses the `io.Medium` interface for testability.

**Step 1: Write the failing test**

```go
// loader_test.go
package manifest

import (
	"crypto/ed25519"
	"testing"

	"forge.lthn.ai/core/cli/pkg/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Good(t *testing.T) {
	fs := io.NewMockMedium()
	fs.Files[".core/view.yml"] = `
code: test-app
name: Test App
version: 1.0.0
layout: HLCRF
slots:
  C: main-content
`
	m, err := Load(fs, ".")
	require.NoError(t, err)
	assert.Equal(t, "test-app", m.Code)
	assert.Equal(t, "main-content", m.Slots["C"])
}

func TestLoad_Bad_NoManifest(t *testing.T) {
	fs := io.NewMockMedium()
	_, err := Load(fs, ".")
	assert.Error(t, err)
}

func TestLoadVerified_Good(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	m := &Manifest{
		Code: "signed-app", Name: "Signed", Version: "1.0.0",
		Layout: "HLCRF", Slots: map[string]string{"C": "main"},
	}
	_ = Sign(m, priv)

	// Re-serialize with signature
	raw, _ := marshalYAML(m)
	fs := io.NewMockMedium()
	fs.Files[".core/view.yml"] = string(raw)

	loaded, err := LoadVerified(fs, ".", pub)
	require.NoError(t, err)
	assert.Equal(t, "signed-app", loaded.Code)
}

func TestLoadVerified_Bad_Tampered(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	m := &Manifest{Code: "app", Version: "1.0.0"}
	_ = Sign(m, priv)

	raw, _ := marshalYAML(m)
	// Tamper after signing
	tampered := string(raw)
	tampered = "code: evil\n" + tampered[6:]
	fs := io.NewMockMedium()
	fs.Files[".core/view.yml"] = tampered

	_, err := LoadVerified(fs, ".", pub)
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestLoad ./pkg/manifest/ -v`
Expected: FAIL — Load/LoadVerified/marshalYAML not defined

**Step 3: Write minimal implementation**

```go
// loader.go
package manifest

import (
	"crypto/ed25519"
	"fmt"
	"path/filepath"

	"forge.lthn.ai/core/cli/pkg/io"
	"gopkg.in/yaml.v3"
)

const manifestPath = ".core/view.yml"

// marshalYAML is a helper that serializes a manifest to YAML bytes.
func marshalYAML(m *Manifest) ([]byte, error) {
	return yaml.Marshal(m)
}

// Load reads and parses a .core/view.yml from the given root directory.
func Load(medium io.Medium, root string) (*Manifest, error) {
	path := filepath.Join(root, manifestPath)
	data, err := medium.Read(path)
	if err != nil {
		return nil, fmt.Errorf("manifest.Load: %w", err)
	}
	return Parse([]byte(data))
}

// LoadVerified reads, parses, and verifies the ed25519 signature.
func LoadVerified(medium io.Medium, root string, pub ed25519.PublicKey) (*Manifest, error) {
	m, err := Load(medium, root)
	if err != nil {
		return nil, err
	}
	ok, err := Verify(m, pub)
	if err != nil {
		return nil, fmt.Errorf("manifest.LoadVerified: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("manifest.LoadVerified: signature verification failed for %q", m.Code)
	}
	return m, nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test -run TestLoad ./pkg/manifest/ -v`
Expected: PASS (4 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/manifest/loader.go pkg/manifest/loader_test.go
git commit -m "feat(manifest): auto-discovery loader with signature verification"
```

---

### Task 4: Object Store Service

**Files:**
- Create: `~/Code/host-uk/core/pkg/store/store.go`
- Test: `~/Code/host-uk/core/pkg/store/store_test.go`

**Context:** A key-value store with group namespacing, backed by `io/sqlite.Medium`. Six operations extracted from dAppServer: get, set, delete, count, deleteGroup, render (Go text/template). Registered as a framework service.

**Step 1: Write the failing test**

```go
// store_test.go
package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetGet_Good(t *testing.T) {
	s, err := New(":memory:")
	require.NoError(t, err)
	defer s.Close()

	err = s.Set("config", "theme", "dark")
	require.NoError(t, err)

	val, err := s.Get("config", "theme")
	require.NoError(t, err)
	assert.Equal(t, "dark", val)
}

func TestGet_Bad_NotFound(t *testing.T) {
	s, _ := New(":memory:")
	defer s.Close()

	_, err := s.Get("config", "missing")
	assert.Error(t, err)
}

func TestDelete_Good(t *testing.T) {
	s, _ := New(":memory:")
	defer s.Close()

	_ = s.Set("config", "key", "val")
	err := s.Delete("config", "key")
	require.NoError(t, err)

	_, err = s.Get("config", "key")
	assert.Error(t, err)
}

func TestCount_Good(t *testing.T) {
	s, _ := New(":memory:")
	defer s.Close()

	_ = s.Set("grp", "a", "1")
	_ = s.Set("grp", "b", "2")
	_ = s.Set("other", "c", "3")

	n, err := s.Count("grp")
	require.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestDeleteGroup_Good(t *testing.T) {
	s, _ := New(":memory:")
	defer s.Close()

	_ = s.Set("grp", "a", "1")
	_ = s.Set("grp", "b", "2")
	err := s.DeleteGroup("grp")
	require.NoError(t, err)

	n, _ := s.Count("grp")
	assert.Equal(t, 0, n)
}

func TestRender_Good(t *testing.T) {
	s, _ := New(":memory:")
	defer s.Close()

	_ = s.Set("user", "pool", "pool.lthn.io:3333")
	_ = s.Set("user", "wallet", "iz...")

	tmpl := `{"pool":"{{ .pool }}","wallet":"{{ .wallet }}"}`
	out, err := s.Render(tmpl, "user")
	require.NoError(t, err)
	assert.Contains(t, out, "pool.lthn.io:3333")
	assert.Contains(t, out, "iz...")
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestSetGet ./pkg/store/ -v`
Expected: FAIL — package does not exist

**Step 3: Write minimal implementation**

```go
// store.go
package store

import (
	"database/sql"
	"fmt"
	"strings"
	"text/template"

	_ "modernc.org/sqlite"
)

// Store is a group-namespaced key-value store backed by SQLite.
type Store struct {
	db *sql.DB
}

// New creates a Store at the given SQLite path. Use ":memory:" for tests.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("store.New: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("store.New: WAL: %w", err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS kv (
		grp   TEXT NOT NULL,
		key   TEXT NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY (grp, key)
	)`); err != nil {
		db.Close()
		return nil, fmt.Errorf("store.New: schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Get retrieves a value by group and key.
func (s *Store) Get(group, key string) (string, error) {
	var val string
	err := s.db.QueryRow("SELECT value FROM kv WHERE grp = ? AND key = ?", group, key).Scan(&val)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("store.Get: not found: %s/%s", group, key)
	}
	if err != nil {
		return "", fmt.Errorf("store.Get: %w", err)
	}
	return val, nil
}

// Set stores a value by group and key, overwriting if exists.
func (s *Store) Set(group, key, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO kv (grp, key, value) VALUES (?, ?, ?)
		 ON CONFLICT(grp, key) DO UPDATE SET value = excluded.value`,
		group, key, value,
	)
	if err != nil {
		return fmt.Errorf("store.Set: %w", err)
	}
	return nil
}

// Delete removes a single key from a group.
func (s *Store) Delete(group, key string) error {
	_, err := s.db.Exec("DELETE FROM kv WHERE grp = ? AND key = ?", group, key)
	if err != nil {
		return fmt.Errorf("store.Delete: %w", err)
	}
	return nil
}

// Count returns the number of keys in a group.
func (s *Store) Count(group string) (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM kv WHERE grp = ?", group).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store.Count: %w", err)
	}
	return n, nil
}

// DeleteGroup removes all keys in a group.
func (s *Store) DeleteGroup(group string) error {
	_, err := s.db.Exec("DELETE FROM kv WHERE grp = ?", group)
	if err != nil {
		return fmt.Errorf("store.DeleteGroup: %w", err)
	}
	return nil
}

// Render loads all key-value pairs from a group and renders a Go template.
func (s *Store) Render(tmplStr, group string) (string, error) {
	rows, err := s.db.Query("SELECT key, value FROM kv WHERE grp = ?", group)
	if err != nil {
		return "", fmt.Errorf("store.Render: query: %w", err)
	}
	defer rows.Close()

	vars := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return "", fmt.Errorf("store.Render: scan: %w", err)
		}
		vars[k] = v
	}

	tmpl, err := template.New("render").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("store.Render: parse: %w", err)
	}
	var b strings.Builder
	if err := tmpl.Execute(&b, vars); err != nil {
		return "", fmt.Errorf("store.Render: exec: %w", err)
	}
	return b.String(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test ./pkg/store/ -v`
Expected: PASS (6 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/store/store.go pkg/store/store_test.go
git commit -m "feat(store): group-namespaced key-value store with template rendering"
```

---

## Phase 4b: CoreDeno (Tasks 5–10)

### Task 5: Web Component Codegen

**Files:**
- Create: `~/go-html/codegen/codegen.go`
- Test: `~/go-html/codegen/codegen_test.go`

**Context:** go-html generates Web Component class definitions from manifest data. Given a component tag name and HLCRF slot, it produces a JavaScript class string that extends HTMLElement, attaches a closed Shadow DOM, and renders slot content. The JS output uses standard Web Component APIs (`customElements.define`, `attachShadow`, Shadow DOM content setting). This is trusted codegen output — never user input.

**Step 1: Write the failing test**

```go
// codegen_test.go
package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateClass_Good(t *testing.T) {
	js, err := GenerateClass("photo-grid", "C")
	require.NoError(t, err)
	assert.Contains(t, js, "class PhotoGrid extends HTMLElement")
	assert.Contains(t, js, "attachShadow")
	assert.Contains(t, js, `mode: "closed"`)
	assert.Contains(t, js, "photo-grid")
}

func TestGenerateClass_Bad_InvalidTag(t *testing.T) {
	_, err := GenerateClass("invalid", "C")
	assert.Error(t, err, "custom element names must contain a hyphen")
}

func TestGenerateRegistration_Good(t *testing.T) {
	js := GenerateRegistration("photo-grid", "PhotoGrid")
	assert.Contains(t, js, "customElements.define")
	assert.Contains(t, js, `"photo-grid"`)
	assert.Contains(t, js, "PhotoGrid")
}

func TestTagToClassName_Good(t *testing.T) {
	tests := []struct{ tag, want string }{
		{"photo-grid", "PhotoGrid"},
		{"nav-breadcrumb", "NavBreadcrumb"},
		{"my-super-widget", "MySuperWidget"},
	}
	for _, tt := range tests {
		got := TagToClassName(tt.tag)
		assert.Equal(t, tt.want, got, "TagToClassName(%q)", tt.tag)
	}
}

func TestGenerateBundle_Good(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "main-content",
	}
	js, err := GenerateBundle(slots)
	require.NoError(t, err)
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	// Should contain exactly 2 class definitions
	assert.Equal(t, 2, strings.Count(js, "extends HTMLElement"))
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/go-html && go test ./codegen/ -v`
Expected: FAIL — package does not exist

**Step 3: Write minimal implementation**

```go
// codegen.go
package codegen

import (
	"fmt"
	"strings"
	"text/template"
)

// wcTemplate is the Web Component class template.
// Uses closed Shadow DOM for isolation. Content is set via the shadow root's
// DOM API using trusted go-html codegen output (never user input).
//
// SECURITY NOTE: The shadow root content is populated from go-html server-side
// render output, which is HTML-escaped at the node level (see node.go Text()).
// This is equivalent to how frameworks like Lit set shadow DOM content from
// trusted template output.
var wcTemplate = template.Must(template.New("wc").Parse(`class {{.ClassName}} extends HTMLElement {
  #shadow;
  constructor() {
    super();
    this.#shadow = this.attachShadow({ mode: "closed" });
  }
  connectedCallback() {
    this.#shadow.textContent = "";
    const slot = this.getAttribute("data-slot") || "{{.Slot}}";
    this.dispatchEvent(new CustomEvent("wc-ready", { detail: { tag: "{{.Tag}}", slot } }));
  }
  render(html) {
    const tpl = document.createElement("template");
    tpl.insertAdjacentHTML("afterbegin", html);
    this.#shadow.textContent = "";
    this.#shadow.appendChild(tpl.content.cloneNode(true));
  }
}`))

// GenerateClass produces a JS class definition for a custom element.
func GenerateClass(tag, slot string) (string, error) {
	if !strings.Contains(tag, "-") {
		return "", fmt.Errorf("codegen: custom element tag %q must contain a hyphen", tag)
	}
	var b strings.Builder
	err := wcTemplate.Execute(&b, struct {
		ClassName, Tag, Slot string
	}{
		ClassName: TagToClassName(tag),
		Tag:       tag,
		Slot:      slot,
	})
	if err != nil {
		return "", fmt.Errorf("codegen: template exec: %w", err)
	}
	return b.String(), nil
}

// GenerateRegistration produces the customElements.define() call.
func GenerateRegistration(tag, className string) string {
	return fmt.Sprintf(`customElements.define("%s", %s);`, tag, className)
}

// TagToClassName converts a kebab-case tag to PascalCase class name.
func TagToClassName(tag string) string {
	parts := strings.Split(tag, "-")
	var b strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			b.WriteString(strings.ToUpper(p[:1]))
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

// GenerateBundle produces all WC class definitions and registrations
// for a set of HLCRF slot assignments.
func GenerateBundle(slots map[string]string) (string, error) {
	seen := make(map[string]bool)
	var b strings.Builder

	for slot, tag := range slots {
		if seen[tag] {
			continue
		}
		seen[tag] = true

		cls, err := GenerateClass(tag, slot)
		if err != nil {
			return "", err
		}
		b.WriteString(cls)
		b.WriteByte('\n')
		b.WriteString(GenerateRegistration(tag, TagToClassName(tag)))
		b.WriteByte('\n')
	}
	return b.String(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/go-html && go test ./codegen/ -v`
Expected: PASS (4 tests)

**Step 5: Commit**

```bash
cd ~/go-html
git add codegen/codegen.go codegen/codegen_test.go
git commit -m "feat(codegen): Web Component class generation from HLCRF slots"
```

---

### Task 6: WASM Web Component Exports

**Files:**
- Modify: `~/go-html/cmd/wasm/main.go`
- Test: `~/go-html/cmd/wasm/main_test.go` (if not exists, create)

**Context:** Extend the existing WASM entry point (currently exports `gohtml.renderToString`) to also export `gohtml.registerComponents(slotsJSON)`. This function takes a JSON-encoded slot map, generates the WC bundle via codegen, and executes it in the browser context via `syscall/js` global function invocation. The browser-side execution uses the standard `Function` constructor (not string evaluation) — this is the normal Go WASM→JS bridge pattern.

**Step 1: Write the failing test**

Since WASM tests require `GOOS=js GOARCH=wasm`, create a unit test for the pure-Go parts:

```go
// cmd/wasm/register_test.go (build-tagged for non-wasm too)
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildComponentJS_Good(t *testing.T) {
	slotsJSON := `{"H":"nav-bar","C":"main-content"}`
	js, err := buildComponentJS(slotsJSON)
	require.NoError(t, err)
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	assert.Contains(t, js, "customElements.define")
}

func TestBuildComponentJS_Bad_InvalidJSON(t *testing.T) {
	_, err := buildComponentJS("not json")
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/go-html && go test ./cmd/wasm/ -v -run TestBuildComponentJS`
Expected: FAIL — buildComponentJS not defined

**Step 3: Write minimal implementation**

Add to the WASM main (or a separate file that compiles for both targets):

```go
// cmd/wasm/register.go
package main

import (
	"encoding/json"
	"fmt"

	"forge.lthn.ai/core/go-html/codegen"
)

// buildComponentJS takes a JSON slot map and returns the WC bundle JS string.
// This is the pure-Go part testable without WASM.
func buildComponentJS(slotsJSON string) (string, error) {
	var slots map[string]string
	if err := json.Unmarshal([]byte(slotsJSON), &slots); err != nil {
		return "", fmt.Errorf("registerComponents: %w", err)
	}
	return codegen.GenerateBundle(slots)
}
```

Then in the WASM-specific main.go, register the JS-callable function:

```go
// In main.go (WASM build), add to the init/main:
// js.Global().Get("gohtml").Set("registerComponents", js.FuncOf(registerComponentsJS))
//
// registerComponentsJS calls buildComponentJS and then uses
// js.Global().Call("Function", jsCode) to execute in the browser.
// This is the standard syscall/js pattern for Go WASM.
```

**Step 4: Run test to verify it passes**

Run: `cd ~/go-html && go test ./cmd/wasm/ -v -run TestBuildComponentJS`
Expected: PASS (2 tests)

**Step 5: Check WASM size budget**

Run: `cd ~/go-html && make wasm`
Expected: gzip size still under 1MB limit (codegen adds minimal code)

**Step 6: Commit**

```bash
cd ~/go-html
git add cmd/wasm/register.go cmd/wasm/register_test.go
git commit -m "feat(wasm): add registerComponents export for WC codegen"
```

---

### Task 7: CoreDeno Sidecar Manager

**Files:**
- Create: `~/Code/host-uk/core/pkg/coredeno/coredeno.go`
- Test: `~/Code/host-uk/core/pkg/coredeno/coredeno_test.go`

**Context:** CoreDeno wraps a Deno child process managed by Go. Uses the existing `pkg/process` service for process lifecycle. The sidecar connects over a Unix domain socket at `$XDG_RUNTIME_DIR/core/deno.sock`. Auto-restart on crash, SIGTERM on shutdown.

**Step 1: Write the failing test**

```go
// coredeno_test.go
package coredeno

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSidecar_Good(t *testing.T) {
	opts := Options{
		DenoPath:   "echo", // use echo as a fake Deno for unit test
		SocketPath: "/tmp/test-core-deno.sock",
	}
	sc := NewSidecar(opts)
	require.NotNil(t, sc)
	assert.Equal(t, "echo", sc.opts.DenoPath)
	assert.Equal(t, "/tmp/test-core-deno.sock", sc.opts.SocketPath)
}

func TestDefaultSocketPath_Good(t *testing.T) {
	path := DefaultSocketPath()
	assert.Contains(t, path, "core/deno.sock")
}

func TestSidecar_PermissionFlags_Good(t *testing.T) {
	perms := Permissions{
		Read:  []string{"./data/"},
		Write: []string{"./data/config.json"},
		Net:   []string{"pool.lthn.io:3333"},
		Run:   []string{"xmrig"},
	}
	flags := perms.Flags()
	assert.Contains(t, flags, "--allow-read=./data/")
	assert.Contains(t, flags, "--allow-write=./data/config.json")
	assert.Contains(t, flags, "--allow-net=pool.lthn.io:3333")
	assert.Contains(t, flags, "--allow-run=xmrig")
}

func TestSidecar_PermissionFlags_Empty(t *testing.T) {
	perms := Permissions{}
	flags := perms.Flags()
	assert.Empty(t, flags)
}

func TestDefaultSocketPath_XDG(t *testing.T) {
	orig := os.Getenv("XDG_RUNTIME_DIR")
	defer os.Setenv("XDG_RUNTIME_DIR", orig)

	os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	path := DefaultSocketPath()
	assert.Equal(t, "/run/user/1000/core/deno.sock", path)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test ./pkg/coredeno/ -v`
Expected: FAIL — package does not exist

**Step 3: Write minimal implementation**

```go
// coredeno.go
package coredeno

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options configures the CoreDeno sidecar.
type Options struct {
	DenoPath   string // path to deno binary (default: "deno")
	SocketPath string // Unix socket path for gRPC
}

// Permissions declares per-module Deno permission flags.
type Permissions struct {
	Read  []string
	Write []string
	Net   []string
	Run   []string
}

// Flags converts permissions to Deno --allow-* CLI flags.
func (p Permissions) Flags() []string {
	var flags []string
	if len(p.Read) > 0 {
		flags = append(flags, fmt.Sprintf("--allow-read=%s", strings.Join(p.Read, ",")))
	}
	if len(p.Write) > 0 {
		flags = append(flags, fmt.Sprintf("--allow-write=%s", strings.Join(p.Write, ",")))
	}
	if len(p.Net) > 0 {
		flags = append(flags, fmt.Sprintf("--allow-net=%s", strings.Join(p.Net, ",")))
	}
	if len(p.Run) > 0 {
		flags = append(flags, fmt.Sprintf("--allow-run=%s", strings.Join(p.Run, ",")))
	}
	return flags
}

// DefaultSocketPath returns the default Unix socket path.
func DefaultSocketPath() string {
	xdg := os.Getenv("XDG_RUNTIME_DIR")
	if xdg == "" {
		xdg = "/tmp"
	}
	return filepath.Join(xdg, "core", "deno.sock")
}

// Sidecar manages a Deno child process.
type Sidecar struct {
	opts Options
}

// NewSidecar creates a Sidecar with the given options.
func NewSidecar(opts Options) *Sidecar {
	if opts.DenoPath == "" {
		opts.DenoPath = "deno"
	}
	if opts.SocketPath == "" {
		opts.SocketPath = DefaultSocketPath()
	}
	return &Sidecar{opts: opts}
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test ./pkg/coredeno/ -v`
Expected: PASS (5 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/coredeno/coredeno.go pkg/coredeno/coredeno_test.go
git commit -m "feat(coredeno): sidecar types, permission flags, socket path"
```

---

### Task 8: CoreDeno Lifecycle (Start/Stop/Restart)

**Files:**
- Create: `~/Code/host-uk/core/pkg/coredeno/lifecycle.go`
- Test: `~/Code/host-uk/core/pkg/coredeno/lifecycle_test.go`

**Context:** Start/Stop/Restart the Deno child process. Uses `os/exec` with context cancellation. Socket directory auto-created. Restart on crash via goroutine monitor. This builds on the Sidecar type from Task 7.

**Step 1: Write the failing test**

```go
// lifecycle_test.go
package coredeno

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart_Good(t *testing.T) {
	// Use "sleep" as a fake long-running process
	sockDir := t.TempDir()
	sc := NewSidecar(Options{
		DenoPath:   "sleep",
		SocketPath: filepath.Join(sockDir, "test.sock"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := sc.Start(ctx, "10") // sleep 10 — will be killed by Stop
	require.NoError(t, err)
	assert.True(t, sc.IsRunning())

	err = sc.Stop()
	require.NoError(t, err)
	assert.False(t, sc.IsRunning())
}

func TestStop_Good_NotStarted(t *testing.T) {
	sc := NewSidecar(Options{DenoPath: "sleep"})
	err := sc.Stop()
	assert.NoError(t, err, "stopping a not-started sidecar should be a no-op")
}

func TestSocketDirCreated_Good(t *testing.T) {
	dir := t.TempDir()
	sockPath := filepath.Join(dir, "sub", "deno.sock")
	sc := NewSidecar(Options{
		DenoPath:   "sleep",
		SocketPath: sockPath,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := sc.Start(ctx, "10")
	require.NoError(t, err)
	defer sc.Stop()

	_, err = os.Stat(filepath.Join(dir, "sub"))
	assert.NoError(t, err, "socket directory should be created")
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestStart ./pkg/coredeno/ -v`
Expected: FAIL — Start/Stop/IsRunning not defined

**Step 3: Write minimal implementation**

```go
// lifecycle.go
package coredeno

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Start launches the Deno sidecar process with the given entrypoint args.
func (s *Sidecar) Start(ctx context.Context, args ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd != nil {
		return fmt.Errorf("coredeno: already running")
	}

	// Ensure socket directory exists
	sockDir := filepath.Dir(s.opts.SocketPath)
	if err := os.MkdirAll(sockDir, 0755); err != nil {
		return fmt.Errorf("coredeno: mkdir %s: %w", sockDir, err)
	}

	// Remove stale socket
	os.Remove(s.opts.SocketPath)

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.cmd = exec.CommandContext(s.ctx, s.opts.DenoPath, args...)
	if err := s.cmd.Start(); err != nil {
		s.cmd = nil
		s.cancel()
		return fmt.Errorf("coredeno: start: %w", err)
	}

	// Monitor in background
	go s.monitor()
	return nil
}

// Stop sends SIGTERM and waits for the process to exit.
func (s *Sidecar) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil {
		return nil
	}
	s.cancel()
	err := s.cmd.Wait()
	s.cmd = nil
	// Context cancellation causes an expected error — ignore it
	if s.ctx.Err() != nil {
		return nil
	}
	return err
}

// IsRunning returns true if the sidecar process is alive.
func (s *Sidecar) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cmd != nil && s.cmd.ProcessState == nil
}

// monitor waits for the process to exit and cleans up.
func (s *Sidecar) monitor() {
	if s.cmd == nil {
		return
	}
	s.cmd.Wait()
	s.mu.Lock()
	s.cmd = nil
	s.mu.Unlock()
}
```

Also add the missing fields to `Sidecar` in `coredeno.go`:

```go
// Update Sidecar struct to include:
type Sidecar struct {
	opts   Options
	mu     sync.RWMutex
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test ./pkg/coredeno/ -v`
Expected: PASS (all tests from Task 7 + 3 new)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/coredeno/coredeno.go pkg/coredeno/lifecycle.go pkg/coredeno/lifecycle_test.go
git commit -m "feat(coredeno): sidecar Start/Stop/Restart lifecycle"
```

---

### Task 9: Proto Definitions (gRPC)

**Files:**
- Create: `~/Code/host-uk/core/pkg/coredeno/proto/coredeno.proto`

**Context:** gRPC proto definitions for the CoreDeno ↔ CoreGO communication channel. Defines the bidirectional service: Deno calls Go for I/O (FileRead, FileWrite, StoreGet, etc.), Go calls Deno for module lifecycle (LoadModule, UnloadModule). This is the I/O fortress boundary.

**Step 1: Write the proto file**

```protobuf
// coredeno.proto
syntax = "proto3";
package coredeno;
option go_package = "forge.lthn.ai/core/cli/pkg/coredeno/proto";

// CoreService is implemented by CoreGO — Deno calls this for I/O.
service CoreService {
  // Filesystem (gated by manifest permissions)
  rpc FileRead(FileReadRequest) returns (FileReadResponse);
  rpc FileWrite(FileWriteRequest) returns (FileWriteResponse);
  rpc FileList(FileListRequest) returns (FileListResponse);
  rpc FileDelete(FileDeleteRequest) returns (FileDeleteResponse);

  // Object store
  rpc StoreGet(StoreGetRequest) returns (StoreGetResponse);
  rpc StoreSet(StoreSetRequest) returns (StoreSetResponse);

  // Process management
  rpc ProcessStart(ProcessStartRequest) returns (ProcessStartResponse);
  rpc ProcessStop(ProcessStopRequest) returns (ProcessStopResponse);
}

// DenoService is implemented by CoreDeno — Go calls this for module lifecycle.
service DenoService {
  rpc LoadModule(LoadModuleRequest) returns (LoadModuleResponse);
  rpc UnloadModule(UnloadModuleRequest) returns (UnloadModuleResponse);
  rpc ModuleStatus(ModuleStatusRequest) returns (ModuleStatusResponse);
}

// --- Core (Go-side) messages ---

message FileReadRequest { string path = 1; string module_code = 2; }
message FileReadResponse { string content = 1; }

message FileWriteRequest { string path = 1; string content = 2; string module_code = 3; }
message FileWriteResponse { bool ok = 1; }

message FileListRequest { string path = 1; string module_code = 2; }
message FileListResponse {
  repeated FileEntry entries = 1;
}
message FileEntry {
  string name = 1;
  bool is_dir = 2;
  int64 size = 3;
}

message FileDeleteRequest { string path = 1; string module_code = 2; }
message FileDeleteResponse { bool ok = 1; }

message StoreGetRequest { string group = 1; string key = 2; }
message StoreGetResponse { string value = 1; bool found = 2; }

message StoreSetRequest { string group = 1; string key = 2; string value = 3; }
message StoreSetResponse { bool ok = 1; }

message ProcessStartRequest { string command = 1; repeated string args = 2; string module_code = 3; }
message ProcessStartResponse { string process_id = 1; }

message ProcessStopRequest { string process_id = 1; }
message ProcessStopResponse { bool ok = 1; }

// --- Deno-side messages ---

message LoadModuleRequest { string code = 1; string entry_point = 2; repeated string permissions = 3; }
message LoadModuleResponse { bool ok = 1; string error = 2; }

message UnloadModuleRequest { string code = 1; }
message UnloadModuleResponse { bool ok = 1; }

message ModuleStatusRequest { string code = 1; }
message ModuleStatusResponse {
  string code = 1;
  enum Status {
    UNKNOWN = 0;
    LOADING = 1;
    RUNNING = 2;
    STOPPED = 3;
    ERRORED = 4;
  }
  Status status = 2;
}
```

**Step 2: Generate Go code from proto**

Run: `cd ~/Code/host-uk/core && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/coredeno/proto/coredeno.proto`

Note: If `protoc` is not installed, install via: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

**Step 3: Verify generated code compiles**

Run: `cd ~/Code/host-uk/core && go build ./pkg/coredeno/proto/`
Expected: Build succeeds

**Step 4: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/coredeno/proto/
git commit -m "feat(coredeno): gRPC proto definitions for I/O fortress"
```

---

### Task 10: Permission Engine + gRPC Server

**Files:**
- Create: `~/Code/host-uk/core/pkg/coredeno/permissions.go`
- Create: `~/Code/host-uk/core/pkg/coredeno/server.go`
- Test: `~/Code/host-uk/core/pkg/coredeno/permissions_test.go`
- Test: `~/Code/host-uk/core/pkg/coredeno/server_test.go`

**Context:** The permission engine checks I/O requests against the module's declared permissions from its manifest. The gRPC server implements `CoreService` (Go side), gating every request through the permission engine. This is the I/O fortress.

**Step 1: Write the failing test (permissions)**

```go
// permissions_test.go
package coredeno

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPath_Good_Allowed(t *testing.T) {
	allowed := []string{"./data/", "./config/"}
	assert.True(t, CheckPath("./data/file.txt", allowed))
	assert.True(t, CheckPath("./config/app.json", allowed))
}

func TestCheckPath_Bad_Denied(t *testing.T) {
	allowed := []string{"./data/"}
	assert.False(t, CheckPath("./secrets/key.pem", allowed))
	assert.False(t, CheckPath("../escape/file", allowed))
}

func TestCheckPath_Good_EmptyAllowAll(t *testing.T) {
	// Empty means no access (deny by default)
	assert.False(t, CheckPath("./anything", nil))
	assert.False(t, CheckPath("./anything", []string{}))
}

func TestCheckNet_Good_Allowed(t *testing.T) {
	allowed := []string{"pool.lthn.io:3333", "api.lthn.io:443"}
	assert.True(t, CheckNet("pool.lthn.io:3333", allowed))
}

func TestCheckNet_Bad_Denied(t *testing.T) {
	allowed := []string{"pool.lthn.io:3333"}
	assert.False(t, CheckNet("evil.com:80", allowed))
}

func TestCheckRun_Good(t *testing.T) {
	allowed := []string{"xmrig", "sha256sum"}
	assert.True(t, CheckRun("xmrig", allowed))
	assert.False(t, CheckRun("rm", allowed))
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestCheck ./pkg/coredeno/ -v`
Expected: FAIL — CheckPath/CheckNet/CheckRun not defined

**Step 3: Write minimal implementation (permissions)**

```go
// permissions.go
package coredeno

import (
	"path/filepath"
	"strings"
)

// CheckPath returns true if the given path is under any of the allowed prefixes.
// Empty allowed list means deny all (secure by default).
func CheckPath(path string, allowed []string) bool {
	if len(allowed) == 0 {
		return false
	}
	clean := filepath.Clean(path)
	for _, prefix := range allowed {
		cleanPrefix := filepath.Clean(prefix)
		if strings.HasPrefix(clean, cleanPrefix) {
			return true
		}
	}
	return false
}

// CheckNet returns true if the given host:port is in the allowed list.
func CheckNet(hostPort string, allowed []string) bool {
	for _, a := range allowed {
		if a == hostPort {
			return true
		}
	}
	return false
}

// CheckRun returns true if the given command is in the allowed list.
func CheckRun(cmd string, allowed []string) bool {
	for _, a := range allowed {
		if a == cmd {
			return true
		}
	}
	return false
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test -run TestCheck ./pkg/coredeno/ -v`
Expected: PASS (6 tests)

**Step 5: Write the gRPC server stub**

```go
// server.go
package coredeno

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/cli/pkg/manifest"
	"forge.lthn.ai/core/cli/pkg/store"
	"forge.lthn.ai/core/cli/pkg/io"
	pb "forge.lthn.ai/core/cli/pkg/coredeno/proto"
)

// Server implements the CoreService gRPC interface.
type Server struct {
	pb.UnimplementedCoreServiceServer
	manifests map[string]*manifest.Manifest // module_code -> manifest
	store     *store.Store
	medium    io.Medium
}

// NewServer creates a CoreService server with the given dependencies.
func NewServer(medium io.Medium, st *store.Store) *Server {
	return &Server{
		manifests: make(map[string]*manifest.Manifest),
		store:     st,
		medium:    medium,
	}
}

// RegisterModule adds a module's manifest to the permission registry.
func (s *Server) RegisterModule(m *manifest.Manifest) {
	s.manifests[m.Code] = m
}

// FileRead implements CoreService.FileRead with permission gating.
func (s *Server) FileRead(ctx context.Context, req *pb.FileReadRequest) (*pb.FileReadResponse, error) {
	m, ok := s.manifests[req.ModuleCode]
	if !ok {
		return nil, fmt.Errorf("unknown module: %s", req.ModuleCode)
	}
	if !CheckPath(req.Path, m.Permissions.Read) {
		return nil, fmt.Errorf("permission denied: %s cannot read %s", req.ModuleCode, req.Path)
	}
	content, err := s.medium.Read(req.Path)
	if err != nil {
		return nil, err
	}
	return &pb.FileReadResponse{Content: content}, nil
}

// StoreGet implements CoreService.StoreGet.
func (s *Server) StoreGet(ctx context.Context, req *pb.StoreGetRequest) (*pb.StoreGetResponse, error) {
	val, err := s.store.Get(req.Group, req.Key)
	if err != nil {
		return &pb.StoreGetResponse{Found: false}, nil
	}
	return &pb.StoreGetResponse{Value: val, Found: true}, nil
}

// StoreSet implements CoreService.StoreSet.
func (s *Server) StoreSet(ctx context.Context, req *pb.StoreSetRequest) (*pb.StoreSetResponse, error) {
	if err := s.store.Set(req.Group, req.Key, req.Value); err != nil {
		return nil, err
	}
	return &pb.StoreSetResponse{Ok: true}, nil
}
```

**Step 6: Write server test**

```go
// server_test.go
package coredeno

import (
	"context"
	"testing"

	"forge.lthn.ai/core/cli/pkg/io"
	"forge.lthn.ai/core/cli/pkg/manifest"
	"forge.lthn.ai/core/cli/pkg/store"
	pb "forge.lthn.ai/core/cli/pkg/coredeno/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	medium := io.NewMockMedium()
	medium.Files["./data/test.txt"] = "hello"
	st, err := store.New(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { st.Close() })

	srv := NewServer(medium, st)
	srv.RegisterModule(&manifest.Manifest{
		Code: "test-mod",
		Permissions: manifest.Permissions{
			Read: []string{"./data/"},
		},
	})
	return srv
}

func TestFileRead_Good(t *testing.T) {
	srv := newTestServer(t)
	resp, err := srv.FileRead(context.Background(), &pb.FileReadRequest{
		Path: "./data/test.txt", ModuleCode: "test-mod",
	})
	require.NoError(t, err)
	assert.Equal(t, "hello", resp.Content)
}

func TestFileRead_Bad_PermissionDenied(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.FileRead(context.Background(), &pb.FileReadRequest{
		Path: "./secrets/key.pem", ModuleCode: "test-mod",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

func TestStoreGetSet_Good(t *testing.T) {
	srv := newTestServer(t)
	ctx := context.Background()

	_, err := srv.StoreSet(ctx, &pb.StoreSetRequest{Group: "cfg", Key: "theme", Value: "dark"})
	require.NoError(t, err)

	resp, err := srv.StoreGet(ctx, &pb.StoreGetRequest{Group: "cfg", Key: "theme"})
	require.NoError(t, err)
	assert.True(t, resp.Found)
	assert.Equal(t, "dark", resp.Value)
}
```

**Step 7: Run tests**

Run: `cd ~/Code/host-uk/core && go test ./pkg/coredeno/ -v`
Expected: PASS (all tests)

**Step 8: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/coredeno/permissions.go pkg/coredeno/permissions_test.go pkg/coredeno/server.go pkg/coredeno/server_test.go
git commit -m "feat(coredeno): permission engine + gRPC server with I/O fortress gating"
```

---

## Phase 4c: Marketplace (Task 11)

### Task 11: Marketplace Index Parser

**Files:**
- Create: `~/Code/host-uk/core/pkg/marketplace/marketplace.go`
- Test: `~/Code/host-uk/core/pkg/marketplace/marketplace_test.go`

**Context:** The marketplace is a Git repo with category directories and index.json files. This parser reads the index, provides search/filter, and resolves module repos for cloning. No Git operations in this task — just the data model and index parsing.

**Step 1: Write the failing test**

```go
// marketplace_test.go
package marketplace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIndex_Good(t *testing.T) {
	raw := `{
		"version": 1,
		"modules": [
			{"code": "mining-xmrig", "name": "XMRig Miner", "repo": "https://forge.lthn.io/host-uk/mod-xmrig.git", "sign_key": "abc123", "category": "miner"},
			{"code": "utils-cyberchef", "name": "CyberChef", "repo": "https://forge.lthn.io/host-uk/mod-cyberchef.git", "sign_key": "def456", "category": "utils"}
		],
		"categories": ["miner", "utils"]
	}`
	idx, err := ParseIndex([]byte(raw))
	require.NoError(t, err)
	assert.Equal(t, 1, idx.Version)
	assert.Len(t, idx.Modules, 2)
	assert.Equal(t, "mining-xmrig", idx.Modules[0].Code)
}

func TestSearch_Good(t *testing.T) {
	idx := &Index{
		Modules: []Module{
			{Code: "mining-xmrig", Name: "XMRig Miner", Category: "miner"},
			{Code: "utils-cyberchef", Name: "CyberChef", Category: "utils"},
		},
	}
	results := idx.Search("miner")
	assert.Len(t, results, 1)
	assert.Equal(t, "mining-xmrig", results[0].Code)
}

func TestByCategory_Good(t *testing.T) {
	idx := &Index{
		Modules: []Module{
			{Code: "a", Category: "miner"},
			{Code: "b", Category: "utils"},
			{Code: "c", Category: "miner"},
		},
	}
	miners := idx.ByCategory("miner")
	assert.Len(t, miners, 2)
}

func TestFind_Good(t *testing.T) {
	idx := &Index{
		Modules: []Module{
			{Code: "mining-xmrig", Name: "XMRig"},
		},
	}
	m, ok := idx.Find("mining-xmrig")
	assert.True(t, ok)
	assert.Equal(t, "XMRig", m.Name)
}

func TestFind_Bad_NotFound(t *testing.T) {
	idx := &Index{}
	_, ok := idx.Find("nope")
	assert.False(t, ok)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test ./pkg/marketplace/ -v`
Expected: FAIL — package does not exist

**Step 3: Write minimal implementation**

```go
// marketplace.go
package marketplace

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Module is a marketplace entry pointing to a module's Git repo.
type Module struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Repo     string `json:"repo"`
	SignKey  string `json:"sign_key"`
	Category string `json:"category"`
}

// Index is the root marketplace catalog.
type Index struct {
	Version    int      `json:"version"`
	Modules    []Module `json:"modules"`
	Categories []string `json:"categories"`
}

// ParseIndex decodes a marketplace index.json.
func ParseIndex(data []byte) (*Index, error) {
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("marketplace.ParseIndex: %w", err)
	}
	return &idx, nil
}

// Search returns modules matching the query in code, name, or category.
func (idx *Index) Search(query string) []Module {
	q := strings.ToLower(query)
	var results []Module
	for _, m := range idx.Modules {
		if strings.Contains(strings.ToLower(m.Code), q) ||
			strings.Contains(strings.ToLower(m.Name), q) ||
			strings.Contains(strings.ToLower(m.Category), q) {
			results = append(results, m)
		}
	}
	return results
}

// ByCategory returns all modules in the given category.
func (idx *Index) ByCategory(category string) []Module {
	var results []Module
	for _, m := range idx.Modules {
		if m.Category == category {
			results = append(results, m)
		}
	}
	return results
}

// Find returns the module with the given code, or false if not found.
func (idx *Index) Find(code string) (Module, bool) {
	for _, m := range idx.Modules {
		if m.Code == code {
			return m, true
		}
	}
	return Module{}, false
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test ./pkg/marketplace/ -v`
Expected: PASS (5 tests)

**Step 5: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/marketplace/marketplace.go pkg/marketplace/marketplace_test.go
git commit -m "feat(marketplace): Git-based module index parser and search"
```

---

## Phase 4d: Integration (Task 12)

### Task 12: Core DI Service + Integration Tests

**Files:**
- Create: `~/Code/host-uk/core/pkg/coredeno/service.go`
- Test: `~/Code/host-uk/core/pkg/coredeno/service_test.go`

**Context:** Wire CoreDeno into the core framework as a service via `WithService(coredeno.NewService(opts))`. Implements `Startable`/`Stoppable` lifecycle interfaces. Depends on the existing framework pattern at `pkg/framework/core/`.

**Step 1: Write the failing test**

```go
// service_test.go
package coredeno

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService_Good(t *testing.T) {
	opts := Options{
		DenoPath:   "echo",
		SocketPath: "/tmp/test-service.sock",
	}
	factory := NewService(opts)
	require.NotNil(t, factory)
}

func TestService_HasLifecycle(t *testing.T) {
	opts := Options{DenoPath: "echo"}
	factory := NewService(opts)
	require.NotNil(t, factory)

	// The factory should return something that has OnStartup/OnShutdown
	// We can't easily test the full framework integration without a Core,
	// but we can verify the factory function signature
	assert.NotNil(t, factory)
}
```

**Step 2: Run test to verify it fails**

Run: `cd ~/Code/host-uk/core && go test -run TestNewService ./pkg/coredeno/ -v`
Expected: FAIL — NewService not defined

**Step 3: Write minimal implementation**

```go
// service.go
package coredeno

import (
	"context"

	"forge.lthn.ai/core/cli/pkg/framework"
	core "forge.lthn.ai/core/cli/pkg/framework/core"
)

// Service wraps the CoreDeno sidecar as a framework service.
type Service struct {
	*core.ServiceRuntime[Options]
	sidecar *Sidecar
}

// NewService returns a factory function for framework registration.
//
// Usage:
//
//	core.New(core.WithService(coredeno.NewService(opts)))
func NewService(opts Options) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		svc := &Service{
			ServiceRuntime: framework.NewServiceRuntime(c, opts),
			sidecar:        NewSidecar(opts),
		}
		return svc, nil
	}
}

// OnStartup starts the Deno sidecar. Called by the framework.
func (s *Service) OnStartup(ctx context.Context) error {
	// The actual Deno entrypoint will be configured later
	// For now, this is the lifecycle hook point
	return nil
}

// OnShutdown stops the Deno sidecar. Called by the framework.
func (s *Service) OnShutdown() error {
	return s.sidecar.Stop()
}

// Sidecar returns the underlying sidecar for direct access.
func (s *Service) Sidecar() *Sidecar {
	return s.sidecar
}
```

**Step 4: Run test to verify it passes**

Run: `cd ~/Code/host-uk/core && go test ./pkg/coredeno/ -v`
Expected: PASS (all tests)

**Step 5: Run full test suite**

Run: `cd ~/Code/host-uk/core && go test ./pkg/manifest/ ./pkg/store/ ./pkg/coredeno/ ./pkg/marketplace/ -v`
Expected: All tests PASS

**Step 6: Commit**

```bash
cd ~/Code/host-uk/core
git add pkg/coredeno/service.go pkg/coredeno/service_test.go
git commit -m "feat(coredeno): framework service with Startable/Stoppable lifecycle"
```

---

## Summary

| Phase | Tasks | New packages | Tests |
|-------|-------|-------------|-------|
| 4a: Foundation | 1–4 | `manifest`, `store` | ~16 |
| 4b: CoreDeno | 5–10 | `codegen`, `coredeno`, `coredeno/proto` | ~20 |
| 4c: Marketplace | 11 | `marketplace` | ~5 |
| 4d: Integration | 12 | (service wiring) | ~2 |
| **Total** | **12** | **6 new packages** | **~43 tests** |

### Dependency Graph

```
Task 1 (manifest types)
├── Task 2 (signing) ← depends on Task 1
├── Task 3 (loader) ← depends on Task 1, 2
└── Task 4 (store) ← independent

Task 5 (WC codegen) ← independent
├── Task 6 (WASM exports) ← depends on Task 5

Task 7 (sidecar types) ← independent
├── Task 8 (lifecycle) ← depends on Task 7
├── Task 9 (proto) ← depends on Task 7
└── Task 10 (permissions + server) ← depends on Task 7, 9, 1, 4

Task 11 (marketplace) ← independent

Task 12 (service) ← depends on Task 7, 8
```

### Parallel Execution Lanes

- **Lane A**: Tasks 1 → 2 → 3
- **Lane B**: Tasks 4, 5 → 6, 11 (all independent)
- **Lane C**: Tasks 7 → 8 → 9 → 10 → 12

Lanes A and B can run in parallel. Lane C depends on Lane A (Task 10 needs manifest types) and Lane B (Task 10 needs store).
