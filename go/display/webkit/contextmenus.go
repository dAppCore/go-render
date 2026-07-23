// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/contextmenu"
	"dappco.re/go/render/display/webkit/pkg/events"
)

// ContextMenu declares one named right-click surface. The Lit
// frontend opts an element in via CSS custom property
// `--custom-contextmenu: <Name>`; per-click context flows via
// `--custom-contextmenu-data: <value>` and surfaces at the click
// handler as ActionItemClicked.Data.
//
// When EventTemplate is non-empty, webkit.Service installs a relay that
// emits the configured event on every item click. The template
// supports {menu} and {action} placeholders — e.g.
// "lthn:context:{menu}:{action}" yields
// "lthn:context:message:copy" for the "copy" action of the "message"
// menu. The {menu} placeholder uses the menu Name with the optional
// MenuPrefixStrip removed (handy when menu names are prefixed
// "lthn-message" but the event slot wants just "message").
//
//	cfg.ContextMenus = []webkit.ContextMenu{
//	    {
//	        Name:             "lthn-message",
//	        EventTemplate:    "lthn:context:{menu}:{action}",
//	        MenuPrefixStrip:  "lthn-",
//	        Items: []webkit.ContextMenuItem{
//	            {Label: "Copy", ActionID: "copy"},
//	            {Type: "separator"},
//	            {Label: "Regenerate", ActionID: "regenerate"},
//	        },
//	    },
//	}
type ContextMenu struct {
	Name            string
	Items           []ContextMenuItem
	EventTemplate   string
	MenuPrefixStrip string
}

// ContextMenuItem aliases the contextmenu MenuItemDef shape so
// consumers don't need to import the contextmenu sub-package.
type ContextMenuItem = contextmenu.MenuItemDef

// applyContextMenus registers each declared menu + installs the shared
// relay that emits the configured event on item click. Called by
// webkit.Service.start() after the contextmenu sub-service has started.
func applyContextMenus(c *core.Core, menus []ContextMenu) {
	if c == nil || len(menus) == 0 {
		return
	}
	ctx := core.Background()

	// Index for the relay: menu name → (template, prefixStrip). Built
	// once at registration; the relay does O(1) lookup per click.
	type relayEntry struct {
		template    string
		prefixStrip string
	}
	relayByMenu := make(map[string]relayEntry, len(menus))
	for _, m := range menus {
		if m.Name == "" {
			continue
		}
		relayByMenu[m.Name] = relayEntry{template: m.EventTemplate, prefixStrip: m.MenuPrefixStrip}
		c.Action("contextmenu.add").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: contextmenu.TaskAdd{
				Name: m.Name,
				Menu: contextmenu.ContextMenuDef{Name: m.Name, Items: m.Items},
			}},
		))
	}

	c.RegisterAction(func(c *core.Core, msg core.Message) core.Result {
		clicked, ok := msg.(contextmenu.ActionItemClicked)
		if !ok {
			return core.Ok(nil)
		}
		entry, ok := relayByMenu[clicked.MenuName]
		if !ok || entry.template == "" {
			return core.Ok(nil)
		}
		menu := clicked.MenuName
		if entry.prefixStrip != "" {
			menu = core.TrimPrefix(menu, entry.prefixStrip)
		}
		event := core.Replace(entry.template, "{menu}", menu)
		event = core.Replace(event, "{action}", clicked.ActionID)
		return c.Action("events.emit").Run(core.Background(), core.NewOptions(
			core.Option{Key: "task", Value: events.TaskEmit{Name: event, Data: clicked.Data}},
		))
	})
}
