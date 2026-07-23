// pkg/systray/wails.go
package systray

import (
	core "dappco.re/go"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform using Wails v3.
// Use: platform := systray.NewWailsPlatform(app)
type WailsPlatform struct {
	app *application.App
}

// NewWailsPlatform creates a Wails-backed tray platform.
// Use: platform := systray.NewWailsPlatform(app)
func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

// NewTray creates a Wails system tray handle.
// Use: tray := platform.NewTray()
func (wp *WailsPlatform) NewTray() PlatformTray {
	// alpha.91 may surface a nil SystemTray on platforms without one;
	// guard before construction. Keep `app:` populated because
	// AttachWindow below needs wt.app.Window.GetByName.
	if wp == nil || wp.app == nil || wp.app.SystemTray == nil {
		return nil
	}
	return &wailsTray{app: wp.app, tray: wp.app.SystemTray.New()}
}

// NewMenu creates a Wails tray menu handle.
// Use: menu := platform.NewMenu()
func (wp *WailsPlatform) NewMenu() PlatformMenu {
	if wp == nil || wp.app == nil {
		return &wailsTrayMenu{menu: application.NewMenu()}
	}
	return &wailsTrayMenu{menu: wp.app.NewMenu()}
}

type wailsTray struct {
	app  *application.App
	tray *application.SystemTray
}

func (wt *wailsTray) SetIcon(data []byte)         { wt.tray.SetIcon(data) }
func (wt *wailsTray) SetTemplateIcon(data []byte) { wt.tray.SetTemplateIcon(data) }
func (wt *wailsTray) SetTooltip(text string)      { wt.tray.SetTooltip(text) }
func (wt *wailsTray) SetLabel(text string)        { wt.tray.SetLabel(text) }

func (wt *wailsTray) SetMenu(menu PlatformMenu) {
	if wm, ok := menu.(*wailsTrayMenu); ok {
		wt.tray.SetMenu(wm.menu)
	}
}

// AttachWindow anchors a previously-created window (by name) to the
// tray. Looks the window up in wails's WindowManager and chains
// WindowOffset(offsetY) when offsetY > 0. offsetX is currently unused —
// wails3 alpha.91 exposes a single offset along the platform's natural
// axis only.
func (wt *wailsTray) AttachWindow(w WindowHandle, offsetX, offsetY int) {
	_ = offsetX
	if wt == nil || wt.tray == nil || wt.app == nil || w == nil {
		return
	}
	window, ok := wt.app.Window.GetByName(w.Name())
	if !ok || window == nil {
		return
	}
	attached := wt.tray.AttachWindow(window)
	if attached != nil && offsetY != 0 {
		attached.WindowOffset(offsetY)
	}
}

func (wt *wailsTray) ShowMessage(title, message string) resultFailure {
	_ = title
	_ = message
	return core.E("systray.wailsTray.ShowMessage", "tray balloon messages are not supported by this backend", nil)
}

// wailsTrayMenu wraps *application.Menu for the PlatformMenu interface.
type wailsTrayMenu struct {
	menu *application.Menu
}

func (m *wailsTrayMenu) Add(label string) PlatformMenuItem {
	return &wailsTrayMenuItem{item: m.menu.Add(label)}
}

func (m *wailsTrayMenu) AddSeparator() {
	m.menu.AddSeparator()
}

func (m *wailsTrayMenu) AddSubmenu(label string) PlatformMenu {
	return &wailsTrayMenu{menu: m.menu.AddSubmenu(label)}
}

// wailsTrayMenuItem wraps *application.MenuItem for the PlatformMenuItem interface.
type wailsTrayMenuItem struct {
	item *application.MenuItem
}

func (mi *wailsTrayMenuItem) SetTooltip(text string)  { mi.item.SetTooltip(text) }
func (mi *wailsTrayMenuItem) SetChecked(checked bool) { mi.item.SetChecked(checked) }
func (mi *wailsTrayMenuItem) SetEnabled(enabled bool) { mi.item.SetEnabled(enabled) }
func (mi *wailsTrayMenuItem) OnClick(fn func()) {
	mi.item.OnClick(func(ctx *application.Context) { fn() })
}
