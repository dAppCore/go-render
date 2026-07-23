// pkg/window/service_config_test.go
package window

import (
	core "dappco.re/go"
)

func newConfigTestWindowService(t *core.T) (*Service, *Manager) {
	t.Helper()

	mgr := &Manager{
		platform: newMockPlatform(),
		state:    NewStateManagerWithDir(t.TempDir()),
		layout:   NewLayoutManagerWithDir(t.TempDir()),
		windows:  make(map[string]PlatformWindow),
	}

	return &Service{manager: mgr}, mgr
}

func TestServiceConfig_applyConfig_Good(t *core.T) {
	// applyConfig
	ax7Variant := "applyConfig:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, mgr := newConfigTestWindowService(t)

	stateFile := core.PathJoin(t.TempDir(), "window-state.json")
	svc.applyConfig(map[string]any{
		"default_width":  1440,
		"default_height": 900,
		"state_file":     stateFile,
	})

	core.AssertEqual(t, stateFile, mgr.State().filePath())

	pw, err := mgr.Open(WithName("main"))
	core.RequireNoError(t, err)

	width, height := pw.Size()
	core.AssertEqual(t, 1440, width)
	core.AssertEqual(t, 900, height)

	mgr.State().SetState("main", WindowState{Width: width, Height: height})
	mgr.State().ForceSync()

	content, err := coreReadFile(stateFile)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), `"main"`)
}

func TestServiceConfig_applyConfig_Bad(t *core.T) {
	// applyConfig
	ax7Variant := "applyConfig:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, mgr := newConfigTestWindowService(t)
	mgr.SetDefaultWidth(1111)
	mgr.SetDefaultHeight(2222)
	initialPath := mgr.State().filePath()

	svc.applyConfig(map[string]any{
		"default_width":  "wide",
		"default_height": true,
		"state_file":     123,
	})

	core.AssertEqual(t, 1111, mgr.defaultWidth)
	core.AssertEqual(t, 2222, mgr.defaultHeight)
	core.AssertEqual(t, initialPath, mgr.State().filePath())
}

func TestServiceConfig_applyConfig_Ugly(t *core.T) {
	// applyConfig
	ax7Variant := "applyConfig:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, mgr := newConfigTestWindowService(t)
	initialPath := mgr.State().filePath()

	core.AssertNotPanics(t, func() {
		svc.applyConfig(nil)
	})

	core.AssertEqual(t, initialPath, mgr.State().filePath())
}
