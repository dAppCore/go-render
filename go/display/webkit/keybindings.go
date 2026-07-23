// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/events"
	"dappco.re/go/render/display/webkit/pkg/keybinding"
)

// Keybinding declares one global accelerator + its emit-on-trigger
// event. webkit.Service registers the accelerator via keybinding.add and
// wires a single shared click router that emits EventName when the
// accelerator fires.
//
// Accelerator format mirrors the keybinding package — modifier+key
// segments joined with "+" (e.g. "Cmd+J", "Ctrl+Shift+M", "Escape").
// Cross-platform consumers register both "Cmd+X" and "Ctrl+X" for the
// same EventName.
//
//	cfg.Keybindings = []webkit.Keybinding{
//	    {Accelerator: "Cmd+J", Description: "Toggle tray popover", EventName: "lthn:key:popover"},
//	    {Accelerator: "Ctrl+J", Description: "Toggle tray popover", EventName: "lthn:key:popover"},
//	}
type Keybinding struct {
	Accelerator string
	Description string
	EventName   string
}

// applyKeybindings registers every declared accelerator + installs the
// shared trigger router. Called by webkit.Service.start() after the
// keybinding sub-service has started.
func applyKeybindings(c *core.Core, bindings []Keybinding) {
	if c == nil || len(bindings) == 0 {
		return
	}
	ctx := core.Background()
	// Index for the trigger router: accelerator → event name. Built
	// once at registration; the router does an O(1) lookup per event.
	eventByAccel := make(map[string]string, len(bindings))
	for _, b := range bindings {
		if b.Accelerator == "" {
			continue
		}
		eventByAccel[b.Accelerator] = b.EventName
		c.Action("keybinding.add").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: keybinding.TaskAdd{
				Accelerator: b.Accelerator,
				Description: b.Description,
			}},
		))
	}

	c.RegisterAction(func(c *core.Core, msg core.Message) core.Result {
		triggered, ok := msg.(keybinding.ActionTriggered)
		if !ok {
			return core.Ok(nil)
		}
		event := eventByAccel[triggered.Accelerator]
		if event == "" {
			return core.Ok(nil)
		}
		return c.Action("events.emit").Run(core.Background(), core.NewOptions(
			core.Option{Key: "task", Value: events.TaskEmit{Name: event, Data: ""}},
		))
	})
}
