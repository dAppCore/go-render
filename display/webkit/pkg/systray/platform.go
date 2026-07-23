// pkg/systray/platform.go
package systray

// Platform abstracts the system tray backend.
// Use: var p systray.Platform
type Platform interface {
	NewTray() PlatformTray
	NewMenu() PlatformMenu // Menu factory for building tray menus
}

// PlatformTray is a live tray handle from the backend.
// Use: var tray systray.PlatformTray
type PlatformTray interface {
	SetIcon(data []byte)
	SetTemplateIcon(data []byte)
	SetTooltip(text string)
	SetLabel(text string)
	SetMenu(menu PlatformMenu)
	AttachWindow(w WindowHandle, offsetX, offsetY int)
	ShowMessage(title, message string) resultFailure
}

// PlatformMenu is a tray menu built by the backend.
// Use: var menu systray.PlatformMenu
type PlatformMenu interface {
	Add(label string) PlatformMenuItem
	AddSeparator()
	AddSubmenu(label string) PlatformMenu
}

// PlatformMenuItem is a single item in a tray menu.
// Use: var item systray.PlatformMenuItem
type PlatformMenuItem interface {
	SetTooltip(text string)
	SetChecked(checked bool)
	SetEnabled(enabled bool)
	OnClick(fn func())
}

// WindowHandle is a minimal cross-package window reference.
// Defined locally to avoid circular imports (display imports systray).
// Concrete panel operations are invoked dynamically because Wails windows use
// fluent Show/Hide methods while the internal window abstraction uses void
// methods.
// Use: var w systray.WindowHandle
type WindowHandle interface {
	Name() string
}
