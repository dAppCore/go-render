// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	"runtime"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/menu"
)

// MenuItem aliases the menu sub-package's MenuItem so consumers can
// declare application menus without importing the menu package
// directly.
type MenuItem = menu.MenuItem

// MenuRole aliases menu.MenuRole. Use the Role* constants below in
// MenuItem.Role to get standard platform-managed items (Quit / Cut /
// Window list / etc.) without listing each by hand.
type MenuRole = menu.MenuRole

// Role constants for application menu items. Mirrors menu.Role*.
// Wrapped in pointer literals via &webkit.RoleAppMenu when building a
// MenuItem.
var (
	RoleAppMenu    = menu.RoleAppMenu
	RoleFileMenu   = menu.RoleFileMenu
	RoleEditMenu   = menu.RoleEditMenu
	RoleViewMenu   = menu.RoleViewMenu
	RoleWindowMenu = menu.RoleWindowMenu
	RoleHelpMenu   = menu.RoleHelpMenu
)

// applyAppMenu fires menu.set_app_menu when AppMenu is non-empty.
// Auto-gated to darwin: the macOS menu bar IS the application menu
// even for accessory apps; other platforms either have no global menu
// (Linux varies) or expose the menu per-window (Windows) — both cases
// surface via the per-window menu APIs, not the global app menu.
//
// Called by webkit.Service.start() after the menu sub-service starts.
func applyAppMenu(c *core.Core, items []MenuItem) {
	if c == nil || len(items) == 0 {
		return
	}
	if runtime.GOOS != "darwin" {
		return
	}
	c.Action("menu.set_app_menu").Run(core.Background(), core.NewOptions(
		core.Option{Key: "task", Value: menu.TaskSetAppMenu{Items: items}},
	))
}
