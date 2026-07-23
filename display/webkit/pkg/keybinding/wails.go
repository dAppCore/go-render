// pkg/keybinding/wails.go
package keybinding

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's KeyBindingManager.
//
// Wails owns one global key-binding map keyed by accelerator string;
// when a window receives a key event matching a binding, the registered
// callback fires with that window as its argument. Our gui Platform
// passes a callback that takes no window — we wrap it to ignore the
// argument so the consumer doesn't need to know about Wails Windows.
//
// Boot-time bindings configured via application.Options.KeyBindings
// (e.g. ⌘R reload in core-ide's wails_boot.go) live in the same map,
// so this adapter can replace them at runtime — last writer wins.
//
//	wp := keybinding.NewWailsPlatform(app)
//	core.WithService(keybinding.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

// Add registers a global accelerator. Wails accelerator syntax matches
// our doc'd format ("Cmd+S" / "Ctrl+S" / "Shift+F1" / "Cmd+Alt+I" etc.).
func (wp *WailsPlatform) Add(accelerator string, handler func()) resultFailure {
	if wp == nil || wp.app == nil || handler == nil {
		return nil
	}
	wp.app.KeyBinding.Add(accelerator, func(_ application.Window) {
		handler()
	})
	return nil
}

// Remove unregisters by accelerator. No-op if absent.
func (wp *WailsPlatform) Remove(accelerator string) resultFailure {
	if wp == nil || wp.app == nil {
		return nil
	}
	wp.app.KeyBinding.Remove(accelerator)
	return nil
}

// Process triggers the registered handler programmatically. Wails
// requires a Window argument so its callbacks can know the source —
// our wrapped callback ignores it; we still pass the currently focused
// window when possible so any future Wails-aware handlers (e.g. boot-
// time bindings sharing the same map) get the expected context.
func (wp *WailsPlatform) Process(accelerator string) bool {
	if wp == nil || wp.app == nil {
		return false
	}
	var w application.Window
	if wp.app.Window != nil {
		w = wp.app.Window.Current()
	}
	return wp.app.KeyBinding.Process(accelerator, w)
}

// GetAll returns just the accelerator strings — the gui Platform
// contract is intentionally narrower than Wails's *KeyBinding (which
// also exposes the callback) because the callback is opaque to the
// consumer / MCP surface.
func (wp *WailsPlatform) GetAll() []string {
	if wp == nil || wp.app == nil {
		return nil
	}
	all := wp.app.KeyBinding.GetAll()
	out := make([]string, 0, len(all))
	for _, kb := range all {
		if kb == nil {
			continue
		}
		out = append(out, kb.Accelerator)
	}
	return out
}
