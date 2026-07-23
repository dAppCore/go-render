// pkg/systray/tray.go
package systray

import (
	_ "embed"
	"reflect"
	"sync"

	core "dappco.re/go"
)

//go:embed assets/apptray.png
var defaultIcon []byte

// Manager manages the system tray lifecycle.
// Use: manager := systray.NewManager(platform)
type Manager struct {
	platform        Platform
	tray            PlatformTray
	panelWindow     WindowHandle
	callbacks       map[string]func()
	tooltip         string
	label           string
	hasIcon         bool
	hasTemplateIcon bool
	menuItems       []TrayMenuItem
	mu              sync.RWMutex
}

// NewManager creates a systray Manager.
// systray.NewManager(systray.NewWailsPlatform(app)).Setup("Core", "Core")
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform:  platform,
		callbacks: make(map[string]func()),
	}
}

// Setup creates the system tray with default icon and tooltip.
// systray.NewManager(systray.NewWailsPlatform(app)).Setup("Core", "Core")
func (m *Manager) Setup(tooltip, label string) resultFailure {
	m.tray = m.platform.NewTray()
	if m.tray == nil {
		return core.E("systray.Setup", "platform returned nil tray", nil)
	}
	m.tray.SetTemplateIcon(defaultIcon)
	m.tray.SetTooltip(tooltip)
	m.tray.SetLabel(label)
	m.tooltip = tooltip
	m.label = label
	m.hasTemplateIcon = true
	return nil
}

// SetIcon sets the tray icon.
// Use: _ = manager.SetIcon(iconBytes)
func (m *Manager) SetIcon(data []byte) resultFailure {
	if m.tray == nil {
		return core.E("systray.SetIcon", "tray not initialised", nil)
	}
	m.tray.SetIcon(data)
	m.hasIcon = len(data) > 0
	return nil
}

// SetTemplateIcon sets the template icon (macOS).
// Use: _ = manager.SetTemplateIcon(iconBytes)
func (m *Manager) SetTemplateIcon(data []byte) resultFailure {
	if m.tray == nil {
		return core.E("systray.SetTemplateIcon", "tray not initialised", nil)
	}
	m.tray.SetTemplateIcon(data)
	m.hasTemplateIcon = len(data) > 0
	return nil
}

// SetTooltip sets the tray tooltip.
// Use: _ = manager.SetTooltip("Core is ready")
func (m *Manager) SetTooltip(text string) resultFailure {
	if m.tray == nil {
		return core.E("systray.SetTooltip", "tray not initialised", nil)
	}
	m.tray.SetTooltip(text)
	m.tooltip = text
	return nil
}

// SetLabel sets the tray label.
// Use: _ = manager.SetLabel("Core")
func (m *Manager) SetLabel(text string) resultFailure {
	if m.tray == nil {
		return core.E("systray.SetLabel", "tray not initialised", nil)
	}
	m.tray.SetLabel(text)
	m.label = text
	return nil
}

// AttachWindow attaches a panel window to the tray, with an optional
// pixel offset relative to the tray icon (positive Y values move the
// popover down; offset semantics are platform-dependent — Wails treats
// the value as a single offset along the platform's natural axis).
// Use: _ = manager.AttachWindow(windowHandle, 0, 5)
func (m *Manager) AttachWindow(w WindowHandle, offsetX, offsetY int) resultFailure {
	if m.tray == nil {
		return core.E("systray.AttachWindow", "tray not initialised", nil)
	}
	m.mu.Lock()
	m.panelWindow = w
	m.mu.Unlock()
	m.tray.AttachWindow(w, offsetX, offsetY)
	return nil
}

// ShowMessage displays a tray message if the backend supports it.
func (m *Manager) ShowMessage(title, message string) resultFailure {
	if m.tray == nil {
		return core.E("systray.ShowMessage", "tray not initialised", nil)
	}
	return m.tray.ShowMessage(title, message)
}

// ShowPanel reveals the attached tray panel window.
func (m *Manager) ShowPanel() resultFailure {
	m.mu.RLock()
	panel := m.panelWindow
	m.mu.RUnlock()
	if panel == nil {
		return core.E("systray.ShowPanel", "panel window not attached", nil)
	}
	return invokePanelMethod(panel, "Show")
}

// HidePanel hides the attached tray panel window.
func (m *Manager) HidePanel() resultFailure {
	m.mu.RLock()
	panel := m.panelWindow
	m.mu.RUnlock()
	if panel == nil {
		return core.E("systray.HidePanel", "panel window not attached", nil)
	}
	return invokePanelMethod(panel, "Hide")
}

// Tray returns the underlying platform tray for direct access.
// Use: tray := manager.Tray()
func (m *Manager) Tray() PlatformTray {
	return m.tray
}

// IsActive returns whether a tray has been created.
// Use: active := manager.IsActive()
func (m *Manager) IsActive() bool {
	return m.tray != nil
}

func invokePanelMethod(panel WindowHandle, method string) resultFailure {
	value := reflect.ValueOf(panel)
	if !value.IsValid() {
		return core.E("systray.invokePanelMethod", "panel window is invalid", nil)
	}

	target := value.MethodByName(method)
	if !target.IsValid() {
		return core.E("systray.invokePanelMethod", "panel window does not support "+method, nil)
	}

	target.Call(nil)
	return nil
}
