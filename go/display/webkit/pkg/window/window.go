// pkg/window/window.go
package window

import (
	"sync"

	core "dappco.re/go"
)

// Window is CoreGUI's own window descriptor — NOT a Wails type alias.
type Window struct {
	Name                       string
	Title                      string
	URL                        string
	HTML                       string
	JS                         string
	Width, Height              int
	X, Y                       int
	MinWidth, MinHeight        int
	MaxWidth, MaxHeight        int
	Frameless                  bool
	Hidden                     bool
	AlwaysOnTop                bool
	BackgroundColour           [4]uint8
	DisableResize              bool
	EnableFileDrop             bool
	HideOnEscape               bool
	HideOnFocusLost            bool
	DefaultContextMenuDisabled bool
	// HideOnClose: the OS close button hides the window instead of
	// destroying it. The window stays registered + can be re-shown
	// via set_visibility. Tray-rooted apps + steady-state windows
	// (chat, settings) set this so the user can dismiss without
	// paying the cold-start cost on re-open. Applied automatically
	// post-create by gui.Service when the Window is supplied via
	// GuiConfig.Windows.
	HideOnClose bool
	// ContentProtection blocks the OS screen-capture API from
	// recording this window. Wallets / private chat / key reveal
	// surfaces should set true. Applied automatically post-create
	// by gui.Service when the Window is supplied via GuiConfig.Windows.
	// macOS + Windows 10+; no-op on Linux.
	ContentProtection bool
	// ShowDockIcon: when set, gui.OpenWindow fires dock.show_icon
	// before the show sequence so the window comes with a Dock /
	// taskbar presence + Cmd+Tab eligibility. Use for primary "shell"
	// windows in tray-anchored apps that want Dock visibility once
	// the user opens the main UI. Auxiliary windows (settings, about)
	// can stay icon-less to keep the menubar as the canonical entry.
	// No-op on Linux + iOS; macOS handles via NSApp activation
	// policy, Windows via taskbar visibility.
	ShowDockIcon bool
	Mac          MacWindow
	Linux        LinuxWindow
	Windows      WindowsWindow
}

// MacWindow holds macOS-specific window options. Zero values mean
// platform default.
type MacWindow struct {
	WindowLevel             MacWindowLevel
	CollectionBehavior      MacCollectionBehavior
	InvisibleTitleBarHeight int
	DisableBackForwardNav   bool
}

// LinuxWindow holds Linux-specific window options.
type LinuxWindow struct {
	Icon []byte
}

// WindowsWindow holds Windows-specific window options.
type WindowsWindow struct {
	HiddenOnTaskbar bool
}

// MacWindowLevel mirrors application.MacWindowLevel (string-typed
// alias of NSWindow.Level names). Empty string is platform default.
type MacWindowLevel string

const (
	MacWindowLevelDefault     MacWindowLevel = ""
	MacWindowLevelNormal      MacWindowLevel = "normal"
	MacWindowLevelFloating    MacWindowLevel = "floating"
	MacWindowLevelTornOffMenu MacWindowLevel = "tornOffMenu"
	MacWindowLevelModalPanel  MacWindowLevel = "modalPanel"
	MacWindowLevelMainMenu    MacWindowLevel = "mainMenu"
	MacWindowLevelStatus      MacWindowLevel = "status"
	MacWindowLevelPopUpMenu   MacWindowLevel = "popUpMenu"
)

// MacCollectionBehavior is a bitfield mirroring NSWindow.collectionBehavior.
type MacCollectionBehavior uint64

const (
	MacCollectionBehaviorDefault             MacCollectionBehavior = 0
	MacCollectionBehaviorCanJoinAllSpaces    MacCollectionBehavior = 1 << 0
	MacCollectionBehaviorFullScreenAuxiliary MacCollectionBehavior = 1 << 8
	MacCollectionBehaviorIgnoresCycle        MacCollectionBehavior = 1 << 6
)

// ToPlatformOptions converts a Window to PlatformWindowOptions for the backend.
func (w *Window) ToPlatformOptions() PlatformWindowOptions {
	return PlatformWindowOptions{
		Name: w.Name, Title: w.Title, URL: w.URL, HTML: w.HTML, JS: w.JS,
		Width: w.Width, Height: w.Height, X: w.X, Y: w.Y,
		MinWidth: w.MinWidth, MinHeight: w.MinHeight,
		MaxWidth: w.MaxWidth, MaxHeight: w.MaxHeight,
		Frameless: w.Frameless, Hidden: w.Hidden,
		AlwaysOnTop: w.AlwaysOnTop, BackgroundColour: w.BackgroundColour,
		DisableResize: w.DisableResize, EnableFileDrop: w.EnableFileDrop,
		HideOnEscape:               w.HideOnEscape,
		HideOnFocusLost:            w.HideOnFocusLost,
		DefaultContextMenuDisabled: w.DefaultContextMenuDisabled,
		Mac:                        w.Mac,
		Linux:                      w.Linux,
		Windows:                    w.Windows,
	}
}

// Manager manages window lifecycle through a Platform backend.
type Manager struct {
	platform      Platform
	state         *StateManager
	layout        *LayoutManager
	windows       map[string]PlatformWindow
	defaultWidth  int
	defaultHeight int
	mu            sync.RWMutex
}

// NewManager creates a window Manager with the given platform backend.
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManager(),
		layout:   NewLayoutManager(),
		windows:  make(map[string]PlatformWindow),
	}
}

// NewManagerWithDir creates a window Manager with a custom config directory for state/layout persistence.
// Useful for testing or when the default config directory is not appropriate.
func NewManagerWithDir(platform Platform, configDir string) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManagerWithDir(configDir),
		layout:   NewLayoutManagerWithDir(configDir),
		windows:  make(map[string]PlatformWindow),
	}
}

func (m *Manager) SetDefaultWidth(width int) {
	if width > 0 {
		m.defaultWidth = width
	}
}

func (m *Manager) SetDefaultHeight(height int) {
	if height > 0 {
		m.defaultHeight = height
	}
}

// Open creates a window from compatibility options.
// Use: manager.Open(window.WithName("main"), window.WithURL("/"), window.WithSize(1280, 800))
func (m *Manager) Open(options ...WindowOption) (PlatformWindow, resultFailure) {
	windowSpec, err := ApplyOptions(options...)
	if err != nil {
		return nil, core.E("window.Manager.Open", "failed to apply options", err)
	}
	return m.Create(windowSpec)
}

// Create creates a window from a Window descriptor.
func (m *Manager) Create(w *Window) (PlatformWindow, resultFailure) {
	if w.Name == "" {
		w.Name = "main"
	}
	if w.Title == "" {
		w.Title = "Core"
	}
	if w.Width == 0 {
		if m.defaultWidth > 0 {
			w.Width = m.defaultWidth
		} else {
			w.Width = 1280
		}
	}
	if w.Height == 0 {
		if m.defaultHeight > 0 {
			w.Height = m.defaultHeight
		} else {
			w.Height = 800
		}
	}
	if w.URL == "" {
		w.URL = "/"
	}

	// Apply saved state if available
	m.state.ApplyState(w)

	pw := m.platform.CreateWindow(w.ToPlatformOptions())

	m.mu.Lock()
	m.windows[w.Name] = pw
	m.mu.Unlock()

	return pw, nil
}

// Get returns a tracked window by name.
func (m *Manager) Get(name string) (PlatformWindow, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pw, ok := m.windows[name]
	return pw, ok
}

// List returns all tracked window names.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.windows))
	for name := range m.windows {
		names = append(names, name)
	}
	return names
}

// Remove stops tracking a window by name.
func (m *Manager) Remove(name string) {
	m.mu.Lock()
	delete(m.windows, name)
	m.mu.Unlock()
}

// Platform returns the underlying platform for direct access.
func (m *Manager) Platform() Platform {
	return m.platform
}

// State returns the state manager for window persistence.
func (m *Manager) State() *StateManager {
	return m.state
}

// Layout returns the layout manager.
func (m *Manager) Layout() *LayoutManager {
	return m.layout
}
