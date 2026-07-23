// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	"runtime"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/lifecycle"
	"dappco.re/go/render/display/webkit/pkg/systray"
)

// TrayItem is an alias for the systray menu item shape. Re-exported so
// consumers can build TrayConfig without importing the systray package
// directly.
type TrayItem = systray.TrayMenuItem

// TrayConfig declares the system tray surface — icon, tooltip, label,
// menu items, and an optional popover window attachment. Applied by
// webkit.Service.OnStartup after the systray sub-service has started.
//
// Click routing remains caller-owned: register an action handler with
// core.RegisterAction that switches on
// systray.ActionTrayMenuItemClicked.ActionID. Items declare their
// ActionID so the dispatch is decoupled from the menu shape.
//
// Platform branching: when IconTemplate is true, the icon is set via
// the systray.set_template_icon action on darwin (auto-inverted by the
// OS to match menu-bar light/dark), otherwise via systray.set_icon
// uniformly. Tooltip + Label are always set when non-empty.
type TrayConfig struct {
	// Icon is the PNG bytes for the tray glyph. Empty leaves the
	// systray sub-service's bootstrap default in place.
	Icon []byte
	// IconTemplate marks the icon as a macOS template image. The
	// OS auto-inverts it for light/dark menu bars. No-op on other
	// platforms.
	IconTemplate bool
	// Tooltip is the hover text. On macOS this renders as the
	// menu-bar label TEXT next to the icon (there is no separate
	// tooltip surface), so empty here on darwin keeps the tray
	// icon-only. Non-darwin platforms show it as a real tooltip.
	Tooltip string
	// Label is the menu-bar label text. Empty clears display/webkit's
	// systray bootstrap default of "Core".
	Label string
	// Menu is the declarative tray menu. Items support Label, Type
	// ("normal" / "separator" / "checkbox" / "radio"), Checked,
	// Disabled, Tooltip, Submenu, and ActionID for click routing.
	Menu []TrayItem
	// PopoverWindow is the name of a registered window to attach as
	// the tray's popover (clicking the tray icon shows the named
	// window anchored under it). Empty disables the popover
	// attachment; the systray icon stays click-to-show-menu only.
	PopoverWindow string
	// PopoverOffsetY shifts the attached popover N pixels down from
	// the tray icon's anchor point. Useful when the chrome / arrow
	// region needs visual breathing room. Zero = no offset.
	PopoverOffsetY int
	// Routes is the declarative click-routing table. Each entry maps
	// a menu item's ActionID to a window to open, an event to emit,
	// or a quit dispatch. Items whose ActionID is not in the table
	// fall through to any caller-registered ActionTrayMenuItemClicked
	// handler — declarative + bespoke routing can coexist on the
	// same menu (e.g. for prefix-based plugin clicks).
	Routes []TrayRoute
}

// applyTrayConfig fires the systray.set_* + attach_window actions for
// the declared config. Called by webkit.Service.start() once the systray
// sub-service is registered + started.
func applyTrayConfig(c *core.Core, cfg *TrayConfig) {
	if c == nil || cfg == nil {
		return
	}
	ctx := core.Background()
	if cfg.Icon != nil {
		if cfg.IconTemplate && runtime.GOOS == "darwin" {
			c.Action("systray.set_template_icon").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: systray.TaskSetTrayTemplateIcon{Data: cfg.Icon}},
			))
		} else {
			c.Action("systray.set_icon").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: systray.TaskSetTrayIcon{Data: cfg.Icon}},
			))
		}
	}
	c.Action("systray.set_tooltip").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayTooltip{Tooltip: cfg.Tooltip}},
	))
	c.Action("systray.set_label").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayLabel{Label: cfg.Label}},
	))
	if len(cfg.Menu) > 0 {
		c.Action("systray.set_menu").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: systray.TaskSetTrayMenu{Items: cfg.Menu}},
		))
	}
	if cfg.PopoverWindow != "" {
		c.Action("systray.attach_window").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: systray.TaskAttachWindow{Name: cfg.PopoverWindow, OffsetY: cfg.PopoverOffsetY}},
		))
	}
}

// TrayRoute declares the default behaviour for a single tray menu
// item's ActionID. When applyTrayRoutes is installed by webkit.Service,
// any ActionTrayMenuItemClicked carrying a matching ActionID fires
// (in this order):
//
//  1. webkit.OpenWindow(c, OpenWindow) if OpenWindow is non-empty
//  2. webkit.EmitEvent(c, EmitEvent, OpenWindow) if EmitEvent is
//     non-empty (the opened-window name rides as event data so
//     "lthn:tray:open" listeners can switch on it)
//  3. lifecycle.quit dispatch if Quit is true
//
// Items whose ActionID isn't in Routes fall through unhandled — the
// consumer can register an additional RegisterAction handler for
// bespoke routing (e.g. plugin-prefixed ActionIDs whose handler runs
// a registry lookup before opening an ad-hoc window).
//
//	webkit.TrayConfig{
//	    Routes: []webkit.TrayRoute{
//	        {ActionID: "open_app",  OpenWindow: "app",  EmitEvent: "lthn:tray:open"},
//	        {ActionID: "open_chat", OpenWindow: "chat", EmitEvent: "lthn:tray:open"},
//	        {ActionID: "quit",      Quit: true},
//	    },
//	}
type TrayRoute struct {
	ActionID   string
	OpenWindow string
	EmitEvent  string
	Quit       bool
}

// applyTrayRoutes registers a RegisterAction handler that switches on
// ActionTrayMenuItemClicked.ActionID and dispatches via the route
// table. Called by webkit.Service.start() after applyTrayConfig. A nil /
// empty routes table is a no-op; consumers with no declarative routes
// keep the legacy contract (register their own handler externally).
func applyTrayRoutes(c *core.Core, routes []TrayRoute) {
	if c == nil || len(routes) == 0 {
		return
	}
	table := make(map[string]TrayRoute, len(routes))
	for _, r := range routes {
		if r.ActionID == "" {
			continue
		}
		table[r.ActionID] = r
	}
	if len(table) == 0 {
		return
	}
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		click, ok := msg.(systray.ActionTrayMenuItemClicked)
		if !ok {
			return core.Result{OK: true}
		}
		r, ok := table[click.ActionID]
		if !ok {
			return core.Result{OK: true}
		}
		if r.OpenWindow != "" {
			OpenWindow(c, r.OpenWindow)
		}
		if r.EmitEvent != "" {
			EmitEvent(c, r.EmitEvent, r.OpenWindow)
		}
		if r.Quit {
			c.Action("lifecycle.quit").Run(core.Background(), core.NewOptions(
				core.Option{Key: "task", Value: lifecycle.TaskQuit{}},
			))
		}
		return core.Result{OK: true}
	})
}
