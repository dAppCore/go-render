package display

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/marketplace"
	"gopkg.in/yaml.v3"
)

func TestMarketplace_marketplaceRegistryURL_Good(t *core.T) {
	// marketplaceRegistryURL
	ax7Variant := "marketplaceRegistryURL:good"
	core.AssertContains(t, ax7Variant, "good")
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "")

	opts := core.NewOptions(
		core.Option{Key: "url", Value: "  https://override.example/registry  "},
	)

	core.AssertEqual(t, "https://override.example/registry", marketplaceRegistryURL(opts))
}

func TestMarketplace_marketplaceRegistryURL_Bad(t *core.T) {
	// marketplaceRegistryURL
	ax7Variant := "marketplaceRegistryURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "")

	core.AssertEmpty(t, marketplaceRegistryURL(core.NewOptions()))
	core.AssertNotEmpty(t, core.Sprintf("%T", marketplaceRegistryURL(core.NewOptions())))
}

func TestMarketplace_marketplaceRegistryURL_Ugly(t *core.T) {
	// marketplaceRegistryURL
	ax7Variant := "marketplaceRegistryURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "  https://env.example/registry  ")

	core.AssertEqual(t, "https://env.example/registry", marketplaceRegistryURL(core.NewOptions()))
	core.AssertNotEmpty(t, core.Sprintf("%T", marketplaceRegistryURL(core.NewOptions())))
}

func TestMarketplace_marketplaceInstallRoot_Good(t *core.T) {
	// marketplaceInstallRoot
	ax7Variant := "marketplaceInstallRoot:good"
	core.AssertContains(t, ax7Variant, "good")
	root := marketplaceInstallRoot("  /tmp/custom/apps  ")

	core.AssertEqual(t, "/tmp/custom/apps", root)
	core.AssertNotEmpty(t, core.Sprintf("%T", root))
}

func TestMarketplace_marketplaceInstallRoot_Bad(t *core.T) {
	// marketplaceInstallRoot
	ax7Variant := "marketplaceInstallRoot:bad"
	core.AssertContains(t, ax7Variant, "bad")
	t.Setenv("DIR_HOME", "")

	root := marketplaceInstallRoot("")
	core.AssertNotContains(t, root, core.TempDir())
	core.RequireTrue(t, core.HasSuffix(root, core.PathJoin("core", "apps")) || core.HasSuffix(root, core.PathJoin(".core", "apps")))
}

func TestMarketplace_marketplaceInstallRoot_Ugly(t *core.T) {
	// marketplaceInstallRoot
	ax7Variant := "marketplaceInstallRoot:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	t.Setenv("DIR_HOME", "  /Users/tester  ")

	core.RequireTrue(t, core.HasSuffix(marketplaceInstallRoot(""), core.PathJoin(".core", "apps")))
	core.AssertNotEmpty(t, core.Sprintf("%T", core.HasSuffix(marketplaceInstallRoot(""), core.PathJoin(".core", "apps"))))
}

func TestMarketplace_registerMarketplaceActions_GoodCase(t *core.T) {
	_, c := newTestDisplayService(t)

	registry := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("manifests:\n  - name: core-ui\n    version: 1.2.3\n"))
	}))
	t.Cleanup(registry.Close)
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", registry.URL)

	listResult := c.Action("display.marketplace.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, listResult.OK)

	payload, ok := listResult.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, registry.URL, payload["registry_url"])
	manifests, ok := payload["manifests"].([]marketplace.Manifest)
	core.RequireTrue(t, ok)
	core.AssertLen(t, manifests, 1)
	core.AssertEqual(t, "core-ui", manifests[0].Name)

	manifest := signedMarketplaceManifest(t, marketplace.Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	manifestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := yaml.Marshal(manifest)
		core.RequireNoError(t, err)
		_, _ = w.Write(data)
	}))
	t.Cleanup(manifestServer.Close)

	fetchResult := c.Action("display.marketplace.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
	))
	core.RequireTrue(t, fetchResult.OK)
	fetched, ok := fetchResult.Value.(marketplace.Manifest)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, manifest.Name, fetched.Name)
	core.AssertEqual(t, manifest.Ref, fetched.Ref)

	verifyResult := c.Action("display.marketplace.verify").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
	))
	core.RequireTrue(t, verifyResult.OK)
	verified, ok := verifyResult.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, marketplace.DigestManifest(manifest), verified["digest"])

	installDir := t.TempDir()
	var gitArgs []string
	previousGitRunner := marketplaceGitRunner
	marketplaceGitRunner = func(_ context.Context, _ string, args ...string) ([]byte, error) {
		gitArgs = append([]string(nil), args...)
		core.RequireNotEmpty(t, args)
		core.RequireNoError(t, coreMkdirAll(args[len(args)-1], 0o755))
		return nil, nil
	}
	t.Cleanup(func() { marketplaceGitRunner = previousGitRunner })

	installResult := c.Action("display.marketplace.install").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
		core.Option{Key: "install_dir", Value: installDir},
		core.Option{Key: "git_binary", Value: "git"},
	))
	core.RequireTrue(t, installResult.OK)
	installed, ok := installResult.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, installDir, installed["install_dir"])
	resolvedInstallDir, err := pathEvalSymlinks(installDir)
	core.RequireNoError(t, err)
	core.AssertEqual(t, core.PathJoin(resolvedInstallDir, "core-ui"), installed["target_dir"])

	core.AssertContains(t, gitArgs, "clone")
	core.AssertContains(t, gitArgs, "--branch")
	core.AssertContains(t, gitArgs, "--")
}

func TestMarketplace_registerMarketplaceActions_BadCase(t *core.T) {
	_, c := newTestDisplayService(t)

	result := c.Action("display.marketplace.fetch").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "manifest url is required")
}

func TestMarketplace_registerMarketplaceActions_UglyCase(t *core.T) {
	_, c := newTestDisplayService(t)

	result := c.Action("display.marketplace.install").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "manifest url is required")
}

func signedMarketplaceManifest(t *core.T, manifest marketplace.Manifest) marketplace.Manifest {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	core.RequireNoError(t, err)

	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature := ed25519.Sign(priv, []byte(payload))
	manifest.Signature = marketplace.Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(pub),
		Value:     base64.StdEncoding.EncodeToString(signature),
	}
	return manifest
}
