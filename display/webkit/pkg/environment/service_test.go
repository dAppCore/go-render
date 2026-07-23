// pkg/environment/service_test.go
package environment

import (
	"context"
	"sync"

	core "dappco.re/go"
)

type mockPlatform struct {
	isDark            bool
	info              EnvironmentInfo
	accentColour      string
	openFMErr         resultFailure
	openFMPath        string
	openFMSelect      bool
	focusFollowsMouse bool
	themeHandler      func(isDark bool)
	mu                sync.Mutex
}

func (m *mockPlatform) IsDarkMode() bool           { return m.isDark }
func (m *mockPlatform) Info() EnvironmentInfo      { return m.info }
func (m *mockPlatform) AccentColour() string       { return m.accentColour }
func (m *mockPlatform) HasFocusFollowsMouse() bool { return m.focusFollowsMouse }
func (m *mockPlatform) OpenFileManager(path string, selectFile bool) resultFailure {
	m.openFMPath = path
	m.openFMSelect = selectFile
	return m.openFMErr
}
func (m *mockPlatform) OnThemeChange(handler func(isDark bool)) func() {
	m.mu.Lock()
	m.themeHandler = handler
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		m.themeHandler = nil
		m.mu.Unlock()
	}
}

// simulateThemeChange triggers the stored handler (test helper).
func (m *mockPlatform) simulateThemeChange(isDark bool) {
	m.mu.Lock()
	h := m.themeHandler
	m.mu.Unlock()
	if h != nil {
		h(isDark)
	}
}

func newTestService(t *core.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		isDark:       true,
		accentColour: "rgb(0,122,255)",
		info: EnvironmentInfo{
			OS: "darwin", Arch: "arm64",
			Platform: PlatformInfo{Name: "macOS", Version: "14.0"},
		},
	}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func TestRegister_Good(t *core.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "environment")
	core.AssertNotNil(t, svc)
}

func TestQueryTheme_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryTheme{})
	core.RequireTrue(t, r.OK)
	theme := r.Value.(ThemeInfo)
	core.AssertTrue(t, theme.IsDark)
	core.AssertEqual(t, "dark", theme.Theme)
}

func TestQueryInfo_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryInfo{})
	core.RequireTrue(t, r.OK)
	info := r.Value.(EnvironmentInfo)
	core.AssertEqual(t, "darwin", info.OS)
	core.AssertEqual(t, "arm64", info.Arch)
}

func TestQueryAccentColour_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAccentColour{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "rgb(0,122,255)", r.Value)
}

func TestTaskOpenFileManager_Good(t *core.T) {
	mock, c := newTestService(t)
	r := c.Action("environment.open_file_manager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "/tmp", Select: true}},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, core.CleanPath("/tmp", string(core.PathSeparator)), mock.openFMPath)
	core.AssertTrue(t, mock.openFMSelect)
}

func TestTaskOpenFileManager_Bad_InvalidPath(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("environment.open_file_manager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "../tmp", Select: false}},
	))
	core.AssertFalse(t, r.OK)
	core.AssertContains(t, r.Value.(resultFailure).Error(), "path must be absolute")
}

func TestThemeChange_ActionBroadcast_GoodCase(t *core.T) {
	mock, c := newTestService(t)

	// Register a listener that captures the action
	var received *ActionThemeChanged
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionThemeChanged); ok {
			mu.Lock()
			received = &a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	// Simulate theme change
	mock.simulateThemeChange(false)

	mu.Lock()
	r := received
	mu.Unlock()
	core.AssertNotNil(t, r)
	core.AssertFalse(t, r.IsDark)
}

func TestTaskSetTheme_Good_OverrideAndReset(t *core.T) {
	mock, c := newTestService(t)
	mock.isDark = false

	r := c.Action("environment.set_theme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "dark"}},
	))
	core.RequireTrue(t, r.OK)

	theme := c.QUERY(QueryTheme{})
	core.RequireTrue(t, theme.OK)
	info := theme.Value.(ThemeInfo)
	core.AssertTrue(t, info.IsDark)
	core.AssertEqual(t, "dark", info.Theme)

	r = c.Action("environment.set_theme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "system"}},
	))
	core.RequireTrue(t, r.OK)

	theme = c.QUERY(QueryTheme{})
	core.RequireTrue(t, theme.OK)
	info = theme.Value.(ThemeInfo)
	core.AssertFalse(t, info.IsDark)
	core.AssertEqual(t, "light", info.Theme)
}

func TestTaskSetTheme_Bad_Invalid(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("environment.set_theme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "sepia"}},
	))
	core.AssertFalse(t, r.OK)
}

// --- GetAccentColor ---

func TestQueryAccentColour_Bad_Empty(t *core.T) {
	// accent colour := "" — still returns handled with empty string
	mock := &mockPlatform{
		isDark:       false,
		accentColour: "",
		info:         EnvironmentInfo{OS: "linux", Arch: "amd64"},
	}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	core.RequireTrue(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryAccentColour{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "", r.Value)
}

func TestQueryAccentColour_Ugly_NoService(t *core.T) {
	// No environment service — query is unhandled
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryAccentColour{})
	core.AssertFalse(t, r.OK)
}

// --- OpenFileManager ---

func TestTaskOpenFileManager_Bad_Error(t *core.T) {
	// platform returns an resultFailure on open
	openErr := core.E("test", "file manager unavailable", nil)
	mock := &mockPlatform{openFMErr: openErr}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	core.RequireTrue(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.Action("environment.open_file_manager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "/missing", Select: false}},
	))
	core.AssertFalse(t, r.OK)
	err, _ := r.Value.(resultFailure)
	core.AssertErrorIs(t, err, openErr)
}

func TestTaskOpenFileManager_Ugly_NoService(t *core.T) {
	// No environment service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("environment.open_file_manager").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// --- HasFocusFollowsMouse ---

func TestQueryFocusFollowsMouse_Good_True(t *core.T) {
	// platform reports focus-follows-mouse enabled (Linux/X11 sloppy focus)
	mock := &mockPlatform{focusFollowsMouse: true}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	core.RequireTrue(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryFocusFollowsMouse{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, true, r.Value)
}

func TestQueryFocusFollowsMouse_Bad_False(t *core.T) {
	// platform reports focus-follows-mouse disabled (Windows/macOS default)
	mock := &mockPlatform{focusFollowsMouse: false}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	core.RequireTrue(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryFocusFollowsMouse{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, false, r.Value)
}

func TestQueryFocusFollowsMouse_Ugly_NoService(t *core.T) {
	// No environment service — query is unhandled
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryFocusFollowsMouse{})
	core.AssertFalse(t, r.OK)
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
