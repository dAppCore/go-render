// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/events"
)

// EmitEvent fires a Wails custom event onto the WebView's event bus
// via the events.emit action. Wraps the boilerplate consumers were
// hand-rolling for every event broadcast (~10x per typical app's
// sysevents.go + tray click router + keybinding handler etc).
//
//	webkit.EmitEvent(c, "myapp:window:focus", windowName)
//	webkit.EmitEvent(c, "myapp:theme", "dark")
//	webkit.EmitEvent(c, "myapp:notification:response", payload)
//
// Returns the underlying events.emit Result so callers that need the
// success/failure signal can act on it.
func EmitEvent(c *core.Core, name string, data any) core.Result {
	if c == nil {
		return core.Ok(nil)
	}
	return c.Action("events.emit").Run(core.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskEmit{Name: name, Data: data}},
	))
}
