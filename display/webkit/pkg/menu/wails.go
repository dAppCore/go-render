// pkg/menu/wails.go
package menu

import "github.com/wailsapp/wails/v3/pkg/application"

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) NewMenu() PlatformMenu {
	return &wailsMenu{menu: application.NewMenu()}
}

func (wp *WailsPlatform) SetApplicationMenu(menu PlatformMenu) {
	if wp == nil || wp.app == nil {
		return
	}
	if wm, ok := menu.(*wailsMenu); ok && wm != nil && wm.menu != nil {
		wp.app.Menu.SetApplicationMenu(wm.menu)
	}
}

type wailsMenu struct {
	menu *application.Menu
}

func (wm *wailsMenu) Add(label string) PlatformMenuItem {
	return &wailsMenuItem{item: wm.menu.Add(label)}
}

func (wm *wailsMenu) AddSeparator() {
	wm.menu.AddSeparator()
}

func (wm *wailsMenu) AddSubmenu(label string) PlatformMenu {
	sub := wm.menu.AddSubmenu(label)
	return &wailsMenu{menu: sub}
}

func (wm *wailsMenu) AddRole(role MenuRole) {
	switch role {
	case RoleAppMenu:
		wm.menu.AddRole(application.AppMenu)
	case RoleFileMenu:
		wm.menu.AddRole(application.FileMenu)
	case RoleEditMenu:
		wm.menu.AddRole(application.EditMenu)
	case RoleViewMenu:
		wm.menu.AddRole(application.ViewMenu)
	case RoleWindowMenu:
		wm.menu.AddRole(application.WindowMenu)
	case RoleHelpMenu:
		wm.menu.AddRole(application.HelpMenu)
	}
}

type wailsMenuItem struct {
	item *application.MenuItem
}

func (mi *wailsMenuItem) SetAccelerator(accel string) PlatformMenuItem {
	mi.item.SetAccelerator(accel)
	return mi
}

func (mi *wailsMenuItem) SetTooltip(text string) PlatformMenuItem {
	mi.item.SetTooltip(text)
	return mi
}

func (mi *wailsMenuItem) SetChecked(checked bool) PlatformMenuItem {
	mi.item.SetChecked(checked)
	return mi
}

func (mi *wailsMenuItem) SetEnabled(enabled bool) PlatformMenuItem {
	mi.item.SetEnabled(enabled)
	return mi
}

func (mi *wailsMenuItem) OnClick(fn func()) PlatformMenuItem {
	mi.item.OnClick(func(ctx *application.Context) { fn() })
	return mi
}
