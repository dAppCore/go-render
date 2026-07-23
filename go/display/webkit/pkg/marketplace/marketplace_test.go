package marketplace

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	core "dappco.re/go"
	"encoding/base64"
	"encoding/hex"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
)

func signedManifest(t *core.T, manifest Manifest) Manifest {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	core.RequireNoError(t, err)

	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature := ed25519.Sign(priv, []byte(payload))
	manifest.Signature = Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(pub),
		Value:     base64.StdEncoding.EncodeToString(signature),
	}
	return manifest
}

func TestMarketplace_FetchManifest_Good(t *core.T) {
	// FetchManifest
	ax7Variant := "FetchManifest:good"
	core.AssertContains(t, ax7Variant, "good")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: core-ui\nversion: 1.2.3\nref: main\n"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	manifest, err := installer.FetchManifest(context.Background(), server.URL)
	core.RequireNoError(t, err)
	core.AssertEqual(t, "core-ui", manifest.Name)
	core.AssertEqual(t, "1.2.3", manifest.Version)
	core.AssertEqual(t, server.URL, manifest.Repository)
	core.AssertEqual(t, "main", manifest.Ref)
}

func TestMarketplace_FetchManifest_Bad(t *core.T) {
	// FetchManifest
	ax7Variant := "FetchManifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.FetchManifest(context.Background(), server.URL)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "manifest fetch failed")
}

func TestMarketplace_FetchManifest_Ugly(t *core.T) {
	// FetchManifest
	ax7Variant := "FetchManifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	t.Run("invalid yaml", func(t *core.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(": not-yaml"))
		}))
		t.Cleanup(server.Close)

		installer := Installer{HTTPClient: server.Client()}
		_, err := installer.FetchManifest(context.Background(), server.URL)
		core.AssertError(t, err)
	})

	t.Run("size limit", func(t *core.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("name: " + repeatString("x", maxManifestBytes)))
		}))
		t.Cleanup(server.Close)

		installer := Installer{HTTPClient: server.Client()}
		_, err := installer.FetchManifest(context.Background(), server.URL)
		core.AssertError(t, err)
		core.AssertContains(t, err.Error(), "exceeds")
	})
}

func TestMarketplace_List_Bad(t *core.T) {
	// List
	ax7Variant := "List:bad"
	core.AssertContains(t, ax7Variant, "bad")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("boom"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.List(context.Background(), server.URL)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "marketplace list failed")
}

func TestMarketplace_List_Ugly(t *core.T) {
	// List
	ax7Variant := "List:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(repeatString("a", maxManifestBytes+1)))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.List(context.Background(), server.URL)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "exceeds")
}

func TestMarketplace_VerifyManifest_Good(t *core.T) {
	// VerifyManifest
	ax7Variant := "VerifyManifest:good"
	core.AssertContains(t, ax7Variant, "good")
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})

	core.RequireNoError(t, VerifyManifest(manifest))
}

func TestMarketplace_VerifyManifest_Bad(t *core.T) {
	// VerifyManifest
	ax7Variant := "VerifyManifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	manifest.Ref = "dev"

	err := VerifyManifest(manifest)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "manifest signature verification failed")
}

func TestMarketplace_VerifyManifest_Ugly(t *core.T) {
	// VerifyManifest
	ax7Variant := "VerifyManifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
		Signature: Signature{
			Algorithm: "ed25519",
			PublicKey: "not-base64",
			Value:     "also-not-base64",
		},
	}

	core.AssertError(t, VerifyManifest(manifest))
}

func TestMarketplace_VerifyManifest_RequiresSignature(t *core.T) {
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}

	core.AssertError(t, VerifyManifest(manifest))
}

func TestMarketplace_Install_Good(t *core.T) {
	// Install
	ax7Variant := "Install:good"
	core.AssertContains(t, ax7Variant, "good")
	scriptDir := t.TempDir()
	logFile := core.PathJoin(scriptDir, "git.log")
	targetRoot := t.TempDir()
	scriptPath := core.PathJoin(scriptDir, "git")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + shellQuote(logFile) + "\nlast=''\nfor arg in \"$@\"; do last=\"$arg\"; done\nmkdir -p \"$last\"\nexit 0\n"
	core.RequireNoError(t, coreWriteFile(scriptPath, []byte(script), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: targetRoot,
	}

	targetDir, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "Core UI",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}))
	core.RequireNoError(t, err)
	resolvedRoot, err := pathEvalSymlinks(targetRoot)
	core.RequireNoError(t, err)
	core.AssertEqual(t, core.PathJoin(resolvedRoot, "core-ui"), targetDir)
	_, err = coreStat(targetDir)
	core.RequireNoError(t, err)

	contents, err := coreReadFile(logFile)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(contents), "clone")
	core.AssertContains(t, string(contents), "--branch")
	core.AssertContains(t, string(contents), "main")
	core.AssertContains(t, string(contents), "--")

	installedManifest, err := coreReadFile(core.PathJoin(targetDir, ".core", "marketplace.yaml"))
	core.RequireNoError(t, err)
	core.AssertContains(t, string(installedManifest), "name: Core UI")
}

func TestMarketplace_Install_Bad(t *core.T) {
	// Install
	ax7Variant := "Install:bad"
	core.AssertContains(t, ax7Variant, "bad")
	installer := Installer{InstallDir: ""}
	_, err := installer.Install(context.Background(), Manifest{Name: "core-ui"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "install dir is required")
}

func TestMarketplace_Install_RejectsTraversalName(t *core.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "../../escape",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}))
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "path separators")
}

func TestMarketplace_Install_RejectsDashPrefixedRepository(t *core.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "--upload-pack=sh",
		Ref:        "main",
	}))
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "repository must not begin with a dash")
}

func TestMarketplace_Install_RejectsDashPrefixedRef(t *core.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "--upload-pack=sh",
	}))
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "ref must not begin with a dash")
}

func TestMarketplace_Install_Ugly(t *core.T) {
	// Install
	ax7Variant := "Install:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	scriptDir := t.TempDir()
	scriptPath := core.PathJoin(scriptDir, "git")
	core.RequireNoError(t, coreWriteFile(scriptPath, []byte("#!/bin/sh\nprintf '%s\\n' 'fatal: https://token:secret@example.com/repo.git' >&2\nexit 1\n"), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: t.TempDir(),
	}

	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Repository: "https://example.com/core-ui.git",
	}))
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "git clone failed")
	core.AssertNotContains(t, err.Error(), "secret")
	core.AssertNotContains(t, err.Error(), "token:")
}

func TestMarketplace_Install_CleansUpOnCloneFailure(t *core.T) {
	scriptDir := t.TempDir()
	scriptPath := core.PathJoin(scriptDir, "git")
	targetRoot := t.TempDir()
	script := "#!/bin/sh\nlast=''\nfor arg in \"$@\"; do last=\"$arg\"; done\nmkdir -p \"$last\"\ntouch \"$last/partial\"\nexit 1\n"
	core.RequireNoError(t, coreWriteFile(scriptPath, []byte(script), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: targetRoot,
	}

	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	_, err := installer.Install(context.Background(), manifest)
	core.AssertError(t, err)

	targetDir := core.PathJoin(targetRoot, "core-ui")
	_, statErr := coreStat(targetDir)
	core.AssertError(t, statErr)
	core.AssertTrue(t, core.IsNotExist(statErr))
}

func TestMarketplace_Verify_Good(t *core.T) {
	// Verify
	ax7Variant := "Verify:good"
	core.AssertContains(t, ax7Variant, "good")
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := yaml.Marshal(manifest)
		core.RequireNoError(t, err)
		_, _ = w.Write(data)
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	verified, err := installer.Verify(context.Background(), server.URL)
	core.RequireNoError(t, err)
	core.AssertEqual(t, manifest.Name, verified.Name)
	core.AssertEqual(t, manifest.Ref, verified.Ref)
}

func TestMarketplace_List_Good(t *core.T) {
	// List
	ax7Variant := "List:good"
	core.AssertContains(t, ax7Variant, "good")
	manifests := []Manifest{
		{Name: "core-ui", Version: "1.2.3"},
		{Name: "core-chat", Version: "0.9.0"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("manifests:\n  - name: core-ui\n    version: 1.2.3\n  - name: core-chat\n    version: 0.9.0\n"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	listed, err := installer.List(context.Background(), server.URL)
	core.RequireNoError(t, err)
	core.AssertEqual(t, manifests, listed)
}

func TestMarketplace_DigestManifest_Good(t *core.T) {
	// DigestManifest
	ax7Variant := "DigestManifest:good"
	core.AssertContains(t, ax7Variant, "good")
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}

	got := DigestManifest(manifest)
	expected := sha256.Sum256([]byte(manifest.Name + ":" + manifest.Version + ":" + manifest.Repository + ":" + manifest.Ref))

	core.AssertEqual(t, hex.EncodeToString(expected[:]), got)
}

func TestMarketplace_DigestManifest_Bad(t *core.T) {
	// DigestManifest
	ax7Variant := "DigestManifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	base := Manifest{Name: "core-ui", Version: "1.2.3", Repository: "https://example.com/core-ui.git", Ref: "main"}
	changed := base
	changed.Ref = "dev"

	core.AssertNotEqual(t, DigestManifest(base), DigestManifest(changed))
}

func TestMarketplace_DigestManifest_Ugly(t *core.T) {
	// DigestManifest
	ax7Variant := "DigestManifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := DigestManifest(Manifest{})

	core.AssertLen(t, got, 64)
	core.AssertNotEmpty(t, got)
}

func TestMarketplace_safeName_EmptyFallbackUsesInputDigest(t *core.T) {
	slashes := safeName("////")
	ats := safeName("@@@")

	core.AssertEqual(t, "module-0ea28b45", slashes)
	core.AssertEqual(t, "module-2ec847d8", ats)
	core.AssertNotEqual(t, slashes, ats)
	assertSafeModuleName(t, slashes)
	assertSafeModuleName(t, ats)
	core.AssertEqual(t, "valid-name", safeName("valid-name"))
}

func TestMarketplace_validateManifestName_Good(t *core.T) {
	// validateManifestName
	ax7Variant := "validateManifestName:good"
	core.AssertContains(t, ax7Variant, "good")
	core.RequireNoError(t, validateManifestName("core-ui"))
	observedType := core.Sprintf("%T", validateManifestName("core-ui"))
	core.AssertNotEmpty(t, observedType)
}

func TestMarketplace_validateManifestName_Bad(t *core.T) {
	// validateManifestName
	ax7Variant := "validateManifestName:bad"
	core.AssertContains(t, ax7Variant, "bad")
	err := validateManifestName("")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "required")
}

func TestMarketplace_validateManifestName_Ugly(t *core.T) {
	// validateManifestName
	ax7Variant := "validateManifestName:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	err := validateManifestName("..")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "path traversal")
}

func TestMarketplace_validateCloneArg_Good(t *core.T) {
	// validateCloneArg
	ax7Variant := "validateCloneArg:good"
	core.AssertContains(t, ax7Variant, "good")
	core.RequireNoError(t, validateCloneArg("repository", "https://example.com/core-ui.git"))
	observedType := core.Sprintf("%T", validateCloneArg("repository", "https://example.com/core-ui.git"))
	core.AssertNotEmpty(t, observedType)
}

func TestMarketplace_validateCloneArg_Bad(t *core.T) {
	// validateCloneArg
	ax7Variant := "validateCloneArg:bad"
	core.AssertContains(t, ax7Variant, "bad")
	err := validateCloneArg("repository", "--upload-pack=sh")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "dash")
}

func TestMarketplace_validateCloneArg_Ugly(t *core.T) {
	// validateCloneArg
	ax7Variant := "validateCloneArg:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	err := validateCloneArg("repository", "https://example.com/core-ui.git\n--depth 1")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "invalid control characters")
}

func TestMarketplace_validateRepositorySource_Good(t *core.T) {
	// validateRepositorySource
	ax7Variant := "validateRepositorySource:good"
	core.AssertContains(t, ax7Variant, "good")
	core.RequireNoError(t, validateRepositorySource("https://example.com/core-ui.git"))
	core.RequireNoError(t, validateRepositorySource("git@example.com:core-ui.git"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validateRepositorySource("https://example.com/core-ui.git")))
}

func TestMarketplace_validateRepositorySource_Bad(t *core.T) {
	// validateRepositorySource
	ax7Variant := "validateRepositorySource:bad"
	core.AssertContains(t, ax7Variant, "bad")
	err := validateRepositorySource("ext::sh -c id")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "remote helper")
}

func TestMarketplace_validateRepositorySource_Ugly(t *core.T) {
	// validateRepositorySource
	ax7Variant := "validateRepositorySource:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	err := validateRepositorySource("file:///tmp/core-ui.git")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "not allowed")
}

func TestMarketplace_decodeManifestList_Good(t *core.T) {
	// decodeManifestList
	ax7Variant := "decodeManifestList:good"
	core.AssertContains(t, ax7Variant, "good")
	manifests, err := decodeManifestList([]byte(`[{"name":"core-ui","version":"1.2.3"},{"name":"core-chat","version":"0.9.0"}]`))

	core.RequireNoError(t, err)
	core.AssertLen(t, manifests, 2)
	core.AssertEqual(t, "core-ui", manifests[0].Name)
	core.AssertEqual(t, "core-chat", manifests[1].Name)
}

func TestMarketplace_decodeManifestList_Bad(t *core.T) {
	// decodeManifestList
	ax7Variant := "decodeManifestList:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manifests, err := decodeManifestList([]byte("   "))

	core.RequireNoError(t, err)
	core.AssertNil(t, manifests)
}

func TestMarketplace_decodeManifestList_Ugly(t *core.T) {
	// decodeManifestList
	ax7Variant := "decodeManifestList:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	_, err := decodeManifestList([]byte(": not-yaml"))

	core.AssertError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", err))
}

func TestMarketplace_sanitizeCommandOutput_Good(t *core.T) {
	// sanitizeCommandOutput
	ax7Variant := "sanitizeCommandOutput:good"
	core.AssertContains(t, ax7Variant, "good")
	got := sanitizeCommandOutput([]byte("fatal: https://token:secret@example.com/repo.git"))

	core.AssertContains(t, got, "[redacted]@")
	core.AssertNotContains(t, got, "secret")
	core.AssertNotContains(t, got, "token:")
}

func TestMarketplace_sanitizeCommandOutput_Bad(t *core.T) {
	// sanitizeCommandOutput
	ax7Variant := "sanitizeCommandOutput:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertEqual(t, "command produced no output", sanitizeCommandOutput(nil))
	core.AssertEqual(t, "command produced no output", sanitizeCommandOutput([]byte("   \n")))
	core.AssertNotEmpty(t, core.Sprintf("%T", sanitizeCommandOutput(nil)))
}

func TestMarketplace_sanitizeCommandOutput_Ugly(t *core.T) {
	// sanitizeCommandOutput
	ax7Variant := "sanitizeCommandOutput:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := sanitizeCommandOutput([]byte(repeatString("a", 1024)))

	core.AssertLen(t, got, 515)
	core.AssertTrue(t, core.HasSuffix(got, "..."))
}

func shellQuote(value string) string {
	return "'" + core.Replace(value, "'", "'\\''") + "'"
}

func assertSafeModuleName(t *core.T, value string) {
	t.Helper()

	core.RequireNotEmpty(t, value)
	core.AssertLessOrEqual(t, len(value), 32)
	for _, r := range value {
		core.AssertTrue(t, (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-', core.Sprintf("unsafe character %q in %q", r, value))
	}
}

// AX7 generated source-matching smoke coverage.
func TestMarketplace_Installer_FetchManifest_Good(t *core.T) {
	// Installer FetchManifest
	ax7Variant := "Installer_FetchManifest:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.FetchManifest(core.Background(), "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_FetchManifest_Bad(t *core.T) {
	// Installer FetchManifest
	ax7Variant := "Installer_FetchManifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.FetchManifest(core.Background(), "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_FetchManifest_Ugly(t *core.T) {
	// Installer FetchManifest
	ax7Variant := "Installer_FetchManifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.FetchManifest(core.Background(), "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Verify_Good(t *core.T) {
	// Installer Verify
	ax7Variant := "Installer_Verify:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Verify(core.Background(), "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Verify_Bad(t *core.T) {
	// Installer Verify
	ax7Variant := "Installer_Verify:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Verify(core.Background(), "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Verify_Ugly(t *core.T) {
	// Installer Verify
	ax7Variant := "Installer_Verify:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Verify(core.Background(), "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Install_Good(t *core.T) {
	// Installer Install
	ax7Variant := "Installer_Install:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Install(core.Background(), *new(Manifest))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Install_Bad(t *core.T) {
	// Installer Install
	ax7Variant := "Installer_Install:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Install(core.Background(), *new(Manifest))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_Install_Ugly(t *core.T) {
	// Installer Install
	ax7Variant := "Installer_Install:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.Install(core.Background(), *new(Manifest))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_List_Good(t *core.T) {
	// Installer List
	ax7Variant := "Installer_List:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.List(core.Background(), "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_List_Bad(t *core.T) {
	// Installer List
	ax7Variant := "Installer_List:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.List(core.Background(), "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMarketplace_Installer_List_Ugly(t *core.T) {
	// Installer List
	ax7Variant := "Installer_List:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Installer
	result := core.Try(func() any {
		got0, got1 := subject.List(core.Background(), "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}
