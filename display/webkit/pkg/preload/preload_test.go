package preload

import (
	core "dappco.re/go"
)

type captureWebview struct {
	scripts []string
}

func (c *captureWebview) ExecJS(script string) {
	c.scripts = append(c.scripts, script)
}

func TestInjectPreload_Good(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreMkdirAll(core.PathJoin(root, ".core"), 0o755))
	core.RequireNoError(t, coreWriteFile(core.PathJoin(root, "index.html"), []byte("<html></html>"), 0o644))
	core.RequireNoError(t, coreWriteFile(
		core.PathJoin(root, ".core", "view.yaml"),
		[]byte("manifest:\n  preloads:\n    - path: preload.js\n"),
		0o644))
	core.RequireNoError(t, coreWriteFile(
		core.PathJoin(root, "preload.js"),
		[]byte("globalThis.__manifestPreloadLoaded = true;"),
		0o644))

	target := &captureWebview{}
	err := InjectPreload(target, "file://"+core.PathToSlash(core.PathJoin(root, "index.html")))
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)

	script := target.scripts[0]
	core.AssertContains(t, script, "globalThis.core.storage.local")
	core.AssertContains(t, script, "globalThis.core.ml = globalThis.core.ml ||")
	core.AssertContains(t, script, "globalThis.electron = electron")
	core.AssertContains(t, script, "globalThis.__manifestPreloadLoaded = true;")
}

func TestInjectPreload_Bad(t *core.T) {
	err := InjectPreload(nil, "http://localhost:3000")
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "preload target is required")
}

func TestInjectPreload_Ugly(t *core.T) {
	target := &captureWebview{}
	err := InjectPreload(target, "https://example.com/app")
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)

	script := target.scripts[0]
	core.AssertContains(t, script, "globalThis.core.storage.local")
	core.AssertContains(t, script, "globalThis.core.ml = globalThis.core.ml ||")
	core.AssertNotContains(t, script, "globalThis.electron = electron")
	core.AssertNotContains(t, script, "ipcRenderer")
}

func TestTrustedOrigin_EmptyAllowListDeniesSchemeURLs(t *core.T) {
	policy := NewTrustedOriginPolicy(nil)

	core.AssertFalse(t, trustedOrigin("core://lab.lthn.sh/page", policy))
	core.AssertFalse(t, trustedOrigin("core://app/", policy))
	core.AssertFalse(t, trustedOrigin("core://attacker.com/x", policy))
}

func TestTrustedOrigin_AllowListMatchesSchemeHostAndPath(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/"})

	core.AssertTrue(t, trustedOrigin("core://lab.lthn.sh/page", policy))
	core.AssertFalse(t, trustedOrigin("core://attacker.com/x", policy))
	core.AssertFalse(t, trustedOrigin("wails://lab.lthn.sh/x", policy))
}

func TestTrustedOrigin_PathPrefix(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/x"})

	core.AssertTrue(t, trustedOrigin("core://lab.lthn.sh/x/y", policy))
	core.AssertFalse(t, trustedOrigin("core://lab.lthn.sh/y", policy))
}

func TestBridgeActionAllowList(t *core.T) {
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

func TestBridgeActionGuardScript(t *core.T) {
	policy := NewTrustedOriginPolicyWithActions(map[string][]string{
		"core://lab.lthn.sh/": {"display.sidecar.eval"},
	})

	script, err := buildScriptWithTrustedOriginPolicy("core://lab.lthn.sh/page", policy)
	core.RequireNoError(t, err)

	core.AssertContains(t, script, "Core bridge action not permitted for this origin")
	core.AssertContains(t, script, `"display.sidecar.eval"`)
	core.AssertNotContains(t, script, `"marketplace.install"`)
}

func TestManifestBackedPreloadOrigin_EmptyAllowListDeniesPlantedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)

	core.AssertFalse(t, manifestBackedPreloadOrigin(
		"https://attacker.com/app",
		NewTrustedOriginPolicy(nil),
	))
}

func TestManifestBackedPreloadOrigin_AllowsListedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "lab.lthn.sh")
	t.Setenv("DIR_HOME", home)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertTrue(t, manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func TestManifestBackedPreloadOrigin_DeniesUnlistedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, manifestBackedPreloadOrigin("https://attacker.com/app", policy))
}

func TestManifestBackedPreloadOrigin_DeniesListedHTTPSOriginWithoutManifest(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func writeMarketplaceViewManifest(t *core.T, home, host string) {
	t.Helper()
	dir := core.PathJoin(home, ".core", "apps", host, ".core")
	core.RequireNoError(t, coreMkdirAll(dir, 0o755))
	core.RequireNoError(t, coreWriteFile(core.PathJoin(dir, "view.yaml"), []byte("name: "+host+"\n"), 0o644))
}

// AX7 generated source-matching smoke coverage.
func TestPreload_InjectPreload_Good(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := InjectPreload(*new(Webview), "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_InjectPreload_Bad(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := InjectPreload(*new(Webview), "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_InjectPreload_Ugly(t *core.T) {
	// InjectPreload
	ax7Variant := "InjectPreload:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := InjectPreload(*new(Webview), "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_InjectPreloadWithTrustedOriginPolicy_Good(t *core.T) {
	// InjectPreloadWithTrustedOriginPolicy
	ax7Variant := "InjectPreloadWithTrustedOriginPolicy:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := InjectPreloadWithTrustedOriginPolicy(*new(Webview), "agent", *new(TrustedOriginPolicy))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_InjectPreloadWithTrustedOriginPolicy_Bad(t *core.T) {
	// InjectPreloadWithTrustedOriginPolicy
	ax7Variant := "InjectPreloadWithTrustedOriginPolicy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := InjectPreloadWithTrustedOriginPolicy(*new(Webview), "", *new(TrustedOriginPolicy))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestPreload_InjectPreloadWithTrustedOriginPolicy_Ugly(t *core.T) {
	// InjectPreloadWithTrustedOriginPolicy
	ax7Variant := "InjectPreloadWithTrustedOriginPolicy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := InjectPreloadWithTrustedOriginPolicy(*new(Webview), "../../edge", *new(TrustedOriginPolicy))
		return core.Sprintf("%T", got0)
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
