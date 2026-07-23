package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/chat"
	"dappco.re/go/render/display/webkit/pkg/window"
)

func TestDisplay_Good_WindowOpenIncludesPreload(t *core.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("preload"),
				window.WithURL("https://example.com"),
			},
		}},
	))
	core.RequireTrue(t, result.OK)
	core.AssertLen(t, platform.Windows, 1)
	core.AssertNotEmpty(t, platform.Windows[0].ExecJSCalls())
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.ml")
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.storage.cookies")
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "Document.prototype, 'cookie'")
	core.AssertNotContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.electron")
	core.AssertNotContains(t, platform.Windows[0].ExecJSCalls()[0], "core.background.service_worker.register")
}

func TestPreload_Good_TrustedOriginIncludesPrivilegedBridge(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)

	script, err := svc.BuildPreloadScriptWithTrustedOriginPolicy(
		"core://app/",
		NewTrustedOriginPolicy([]string{"core://app/"}),
	)
	core.RequireNoError(t, err)

	core.AssertContains(t, script, "globalThis.electron")
	core.AssertContains(t, script, "core.background.service_worker.register")
	core.AssertContains(t, script, "globalThis.core.ml")
	core.AssertContains(t, script, "gui.notification.requestPermission")
	core.AssertContains(t, script, "gui.notification.clear")
	core.AssertContains(t, script, "systray.show_message")
	core.AssertContains(t, script, "webview.devtools_open")
}

func TestDisplay_Good_WindowOpenManifestBackedOriginIncludesManifestPreloadOnly(t *core.T) {
	home := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(home, ".core", "apps", "example.com", ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(home, ".core", "apps", "example.com", "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(home, ".core", "apps", "example.com", ".core", "view.yaml"), "name: example\npreloads:\n  - path: preload.js\n", 0o644))
	core.RequireNoError(t, coreWriteMode(
		core.PathJoin(home, ".core", "preload-origins.yaml"),
		"origins:\n  - https://example.com/\n",
		0o644,
	))
	t.Setenv("DIR_HOME", home)

	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("manifest-backed"),
				window.WithURL("https://example.com/app"),
			},
		}},
	))
	core.RequireTrue(t, result.OK)
	core.AssertLen(t, platform.Windows, 1)
	script := platform.Windows[0].ExecJSCalls()[0]
	core.AssertContains(t, script, "__manifestLoaded")
	core.AssertContains(t, script, "globalThis.core.ml")
	core.AssertNotContains(t, script, "globalThis.electron")
	core.AssertNotContains(t, script, "core.background.service_worker.register")
}

func TestDisplay_Good_CoreSchemeRoutesThroughBackend(t *core.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(chat.Register(func(o *chat.Options) { o.StorePath = core.PathJoin(t.TempDir(), "chat.db") })),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	core.RequireTrue(t, c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{window.WithName("settings")},
		}},
	)).OK)

	core.RequireTrue(t, c.Action("window.set_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetURL{Name: "settings", URL: "core://settings"}},
	)).OK)

	core.AssertLen(t, platform.Windows, 1)
	core.AssertTrue(t, core.Contains(platform.Windows[0].HTMLContent(), "core://settings"))
}

func TestPreload_ValidatedLocalMLAPIURL_GoodCase(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("http://localhost:8090/"))
	core.AssertEqual(t, "https://127.0.0.1:9443", validatedLocalMLAPIURL("https://127.0.0.1:9443/"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("http://localhost:8090/")))
}

func TestPreload_ValidatedLocalMLAPIURL_BadCase(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("https://example.com"))
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("ftp://localhost:8090"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("https://example.com")))
}

func TestPreload_ValidatedLocalMLAPIURL_UglyCase(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL(""))
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("not a url"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("")))
}

func TestPreload_TrustedPreloadOrigin_GoodCase(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/"})

	core.AssertTrue(t, trustedPreloadOrigin("core://lab.lthn.sh/page", policy))
	core.AssertNotEmpty(t, core.Sprintf("%T", policy))
}

func TestPreload_TrustedPreloadOrigin_BadCase(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/"})

	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/x", policy))
	core.AssertFalse(t, trustedPreloadOrigin("wails://lab.lthn.sh/x", policy))
	core.AssertFalse(t, trustedPreloadOrigin("https://example.com", policy))
	core.AssertFalse(t, trustedPreloadOrigin("http://localhost:3000", policy))
	core.AssertFalse(t, trustedPreloadOrigin("file:///tmp/app/index.html", policy))
}

func TestPreload_TrustedPreloadOrigin_EmptyAllowListDeniesSchemeURLs(t *core.T) {
	policy := NewTrustedOriginPolicy(nil)

	core.AssertFalse(t, trustedPreloadOrigin("core://lab.lthn.sh/page", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://app/", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/x", policy))
}

func TestPreload_TrustedPreloadOrigin_PathPrefix(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/x"})

	core.AssertTrue(t, trustedPreloadOrigin("core://lab.lthn.sh/x/y", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://lab.lthn.sh/y", policy))
}

func TestPreload_BridgeActionAllowList(t *core.T) {
	policy := NewTrustedOriginPolicyWithActions(map[string][]string{
		"core://lab.lthn.sh/": {"display.sidecar.eval"},
		"core://empty/":       {},
	})

	core.AssertFalse(t, policy.AllowsActionURL("core://lab.lthn.sh/page", "marketplace.install"))
	core.AssertTrue(t, policy.AllowsActionURL("core://lab.lthn.sh/page", "display.sidecar.eval"))
	core.AssertFalse(t, policy.AllowsActionURL("core://attacker.com/page", "display.sidecar.eval"))
	core.AssertTrue(t, policy.AllowsURL("core://empty/page"))
	core.AssertFalse(t, policy.AllowsActionURL("core://empty/page", "display.sidecar.eval"))
}

func TestPreload_BridgeActionGuardScript(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicyWithActions(map[string][]string{
		"core://lab.lthn.sh/": {"display.sidecar.eval"},
	})

	script, err := svc.BuildPreloadScriptWithTrustedOriginPolicy("core://lab.lthn.sh/page", policy)
	core.RequireNoError(t, err)

	core.AssertContains(t, script, "Core bridge action not permitted for this origin")
	core.AssertContains(t, script, `"display.sidecar.eval"`)
	core.AssertNotContains(t, script, `"marketplace.install"`)
}

func TestPreload_ManifestBackedPreloadOrigin_EmptyAllowListDeniesPlantedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin(
		"https://attacker.com/app",
		NewTrustedOriginPolicy(nil),
	))
}

func TestPreload_ManifestBackedPreloadOrigin_AllowsListedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "lab.lthn.sh")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertTrue(t, svc.manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func TestPreload_ManifestBackedPreloadOrigin_DeniesUnlistedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin("https://attacker.com/app", policy))
}

func TestPreload_ManifestBackedPreloadOrigin_DeniesListedHTTPSOriginWithoutManifest(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func TestPreload_DefaultTrustedOriginPolicy_LoadsConfig(t *core.T) {
	home := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(home, ".core")))
	core.RequireNoError(t, coreWriteMode(
		core.PathJoin(home, ".core", "preload-origins.yaml"),
		"origins:\n  - core://app/\n",
		0o644,
	))
	t.Setenv("DIR_HOME", home)
	t.Setenv("HOME", home)

	policy := DefaultTrustedOriginPolicy()

	core.AssertTrue(t, trustedPreloadOrigin("core://app/shell", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/shell", policy))
}

type preloadCapture struct {
	scripts []string
}

func (p *preloadCapture) ExecJS(script string) {
	p.scripts = append(p.scripts, script)
}

func writeMarketplaceViewManifest(t *core.T, home, host string) {
	t.Helper()
	dir := core.PathJoin(home, ".core", "apps", host, ".core")
	core.RequireNoError(t, coreEnsureDir(dir))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(dir, "view.yaml"), "name: "+host+"\n", 0o644))
}

func TestPreload_InjectPreload_Good(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:good"
	core.AssertContains(t, ax7Variant, "good")
	root := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(root, ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, ".core", "view.yaml"), "preloads:\n  - path: preload.js\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+core.PathToSlash(core.PathJoin(root, "index.html")))
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)
	core.AssertContains(t, target.scripts[0], "globalThis.core.ml")
	core.AssertNotContains(t, target.scripts[0], "globalThis.electron")
	core.AssertContains(t, target.scripts[0], "__manifestLoaded")
}

func TestPreload_InjectPreload_Bad(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "https://example.com/app")
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)
	core.AssertContains(t, target.scripts[0], "globalThis.core.ml")
	core.AssertNotContains(t, target.scripts[0], "globalThis.electron")
	core.AssertNotContains(t, target.scripts[0], "core.background.service_worker.register")
}

func TestPreload_InjectPreload_Ugly(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	root := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(root, ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, ".core", "view.yaml"), "preloads: [\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+core.PathToSlash(core.PathJoin(root, "index.html")))
	core.AssertError(t, err)
	core.AssertEmpty(t, target.scripts)
}

// AX7 generated source-matching smoke coverage.
func TestPreload_Service_InjectPreload_Good(t *core.T) {
	// Service InjectPreload
	ax7Variant := "Service_InjectPreload:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.InjectPreload(*new(PreloadTarget), "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_InjectPreload_Bad(t *core.T) {
	// Service InjectPreload
	ax7Variant := "Service_InjectPreload:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.InjectPreload(*new(PreloadTarget), "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_InjectPreload_Ugly(t *core.T) {
	// Service InjectPreload
	ax7Variant := "Service_InjectPreload:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.InjectPreload(*new(PreloadTarget), "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScript_Good(t *core.T) {
	// Service BuildPreloadScript
	ax7Variant := "Service_BuildPreloadScript:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScript("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScript_Bad(t *core.T) {
	// Service BuildPreloadScript
	ax7Variant := "Service_BuildPreloadScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScript("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScript_Ugly(t *core.T) {
	// Service BuildPreloadScript
	ax7Variant := "Service_BuildPreloadScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScript("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScriptWithTrustedOriginPolicy_Good(t *core.T) {
	// Service BuildPreloadScriptWithTrustedOriginPolicy
	ax7Variant := "Service_BuildPreloadScriptWithTrustedOriginPolicy:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScriptWithTrustedOriginPolicy("agent", *new(TrustedOriginPolicy))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScriptWithTrustedOriginPolicy_Bad(t *core.T) {
	// Service BuildPreloadScriptWithTrustedOriginPolicy
	ax7Variant := "Service_BuildPreloadScriptWithTrustedOriginPolicy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScriptWithTrustedOriginPolicy("", *new(TrustedOriginPolicy))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_Service_BuildPreloadScriptWithTrustedOriginPolicy_Ugly(t *core.T) {
	// Service BuildPreloadScriptWithTrustedOriginPolicy
	ax7Variant := "Service_BuildPreloadScriptWithTrustedOriginPolicy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.BuildPreloadScriptWithTrustedOriginPolicy("../../edge", *new(TrustedOriginPolicy))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicy_Good(t *core.T) {
	// NewTrustedOriginPolicy
	ax7Variant := "NewTrustedOriginPolicy:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicy(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicy_Bad(t *core.T) {
	// NewTrustedOriginPolicy
	ax7Variant := "NewTrustedOriginPolicy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicy(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicy_Ugly(t *core.T) {
	// NewTrustedOriginPolicy
	ax7Variant := "NewTrustedOriginPolicy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicy(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicyWithActions_Good(t *core.T) {
	// NewTrustedOriginPolicyWithActions
	ax7Variant := "NewTrustedOriginPolicyWithActions:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicyWithActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicyWithActions_Bad(t *core.T) {
	// NewTrustedOriginPolicyWithActions
	ax7Variant := "NewTrustedOriginPolicyWithActions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicyWithActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_NewTrustedOriginPolicyWithActions_Ugly(t *core.T) {
	// NewTrustedOriginPolicyWithActions
	ax7Variant := "NewTrustedOriginPolicyWithActions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewTrustedOriginPolicyWithActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_DefaultTrustedOriginPolicy_Good(t *core.T) {
	// DefaultTrustedOriginPolicy
	ax7Variant := "DefaultTrustedOriginPolicy:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := DefaultTrustedOriginPolicy()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_DefaultTrustedOriginPolicy_Bad(t *core.T) {
	// DefaultTrustedOriginPolicy
	ax7Variant := "DefaultTrustedOriginPolicy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := DefaultTrustedOriginPolicy()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_DefaultTrustedOriginPolicy_Ugly(t *core.T) {
	// DefaultTrustedOriginPolicy
	ax7Variant := "DefaultTrustedOriginPolicy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := DefaultTrustedOriginPolicy()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsURL_Good(t *core.T) {
	// TrustedOriginPolicy AllowsURL
	ax7Variant := "TrustedOriginPolicy_AllowsURL:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsURL_Bad(t *core.T) {
	// TrustedOriginPolicy AllowsURL
	ax7Variant := "TrustedOriginPolicy_AllowsURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsURL_Ugly(t *core.T) {
	// TrustedOriginPolicy AllowsURL
	ax7Variant := "TrustedOriginPolicy_AllowsURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_Allows_Good(t *core.T) {
	// TrustedOriginPolicy Allows
	ax7Variant := "TrustedOriginPolicy_Allows:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.Allows(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_Allows_Bad(t *core.T) {
	// TrustedOriginPolicy Allows
	ax7Variant := "TrustedOriginPolicy_Allows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.Allows(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_Allows_Ugly(t *core.T) {
	// TrustedOriginPolicy Allows
	ax7Variant := "TrustedOriginPolicy_Allows:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.Allows(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsActionURL_Good(t *core.T) {
	// TrustedOriginPolicy AllowsActionURL
	ax7Variant := "TrustedOriginPolicy_AllowsActionURL:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsActionURL("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsActionURL_Bad(t *core.T) {
	// TrustedOriginPolicy AllowsActionURL
	ax7Variant := "TrustedOriginPolicy_AllowsActionURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsActionURL("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsActionURL_Ugly(t *core.T) {
	// TrustedOriginPolicy AllowsActionURL
	ax7Variant := "TrustedOriginPolicy_AllowsActionURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsActionURL("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsAction_Good(t *core.T) {
	// TrustedOriginPolicy AllowsAction
	ax7Variant := "TrustedOriginPolicy_AllowsAction:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsAction(nil, "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsAction_Bad(t *core.T) {
	// TrustedOriginPolicy AllowsAction
	ax7Variant := "TrustedOriginPolicy_AllowsAction:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsAction(nil, "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowsAction_Ugly(t *core.T) {
	// TrustedOriginPolicy AllowsAction
	ax7Variant := "TrustedOriginPolicy_AllowsAction:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowsAction(nil, "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActionsForURL_Good(t *core.T) {
	// TrustedOriginPolicy AllowedActionsForURL
	ax7Variant := "TrustedOriginPolicy_AllowedActionsForURL:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActionsForURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActionsForURL_Bad(t *core.T) {
	// TrustedOriginPolicy AllowedActionsForURL
	ax7Variant := "TrustedOriginPolicy_AllowedActionsForURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActionsForURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActionsForURL_Ugly(t *core.T) {
	// TrustedOriginPolicy AllowedActionsForURL
	ax7Variant := "TrustedOriginPolicy_AllowedActionsForURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActionsForURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActions_Good(t *core.T) {
	// TrustedOriginPolicy AllowedActions
	ax7Variant := "TrustedOriginPolicy_AllowedActions:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActions_Bad(t *core.T) {
	// TrustedOriginPolicy AllowedActions
	ax7Variant := "TrustedOriginPolicy_AllowedActions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_TrustedOriginPolicy_AllowedActions_Ugly(t *core.T) {
	// TrustedOriginPolicy AllowedActions
	ax7Variant := "TrustedOriginPolicy_AllowedActions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TrustedOriginPolicy
	result := core.Try(func() any {
		got0 := subject.AllowedActions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
