// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	"context"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/systray"
)

// TestApplyTrayRoutes_Good registers a route table and fires three
// ActionTrayMenuItemClicked messages (matching Quit + matching
// EmitEvent + non-matching). Asserts dispatch reaches the expected
// targets via stub action handlers. The OpenWindow path is not
// exercised here — it depends on a registered webkit.Service, which the
// existing webkit.OpenWindow tests cover end-to-end.
func TestApplyTrayRoutes_Good(t *testing.T) {
	c := core.New()

	var quitCount, eventCount int
	c.Action("lifecycle.quit", func(_ context.Context, _ core.Options) core.Result {
		quitCount++
		return core.Result{OK: true}
	})
	c.Action("events.emit", func(_ context.Context, _ core.Options) core.Result {
		eventCount++
		return core.Result{OK: true}
	})

	applyTrayRoutes(c, []TrayRoute{
		{ActionID: "emit_only", EmitEvent: "lthn:tray:open"},
		{ActionID: "quit", Quit: true},
		{ActionID: "", OpenWindow: "ignored"}, // empty ActionID — must be filtered
	})

	c.ACTION(systray.ActionTrayMenuItemClicked{ActionID: "emit_only"})
	c.ACTION(systray.ActionTrayMenuItemClicked{ActionID: "unknown"})
	c.ACTION(systray.ActionTrayMenuItemClicked{ActionID: "quit"})

	if eventCount != 1 {
		t.Errorf("emit_only route fired events.emit %d times, want 1", eventCount)
	}
	if quitCount != 1 {
		t.Errorf("quit route fired lifecycle.quit %d times, want 1", quitCount)
	}
}

// TestApplyTrayRoutes_Bad covers no-op paths: nil core, empty / nil
// routes table, and a route with empty ActionID. None should panic or
// install an active handler.
func TestApplyTrayRoutes_Bad(t *testing.T) {
	applyTrayRoutes(nil, []TrayRoute{{ActionID: "x"}}) // nil core, no panic
	c := core.New()
	applyTrayRoutes(c, nil)                                     // nil slice
	applyTrayRoutes(c, []TrayRoute{})                           // empty slice
	applyTrayRoutes(c, []TrayRoute{{ActionID: "", Quit: true}}) // only empty-ID entries
}
