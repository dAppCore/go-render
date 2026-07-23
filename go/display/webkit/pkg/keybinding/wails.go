// pkg/keybinding/wails.go
package keybinding

import (
	"sync"

	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's GlobalShortcutManager.
//
// Global shortcuts are registered with the operating system and fire even
// when no Wails window is focused. The local callback map preserves the
// adapter's Process method for deterministic tests and command dispatch.
//
//	wp := keybinding.NewWailsPlatform(app)
//	core.WithService(keybinding.Register(wp))
type WailsPlatform struct {
	app      *application.App
	mu       sync.RWMutex
	handlers map[string]func()
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app, handlers: make(map[string]func())}
}

// Add registers a system-wide accelerator. Wails accelerator syntax matches
// our doc'd format ("Cmd+S" / "Ctrl+S" / "Shift+F1" / "Cmd+Alt+I" etc.).
func (wp *WailsPlatform) Add(accelerator string, handler func()) resultFailure {
	if wp == nil || wp.app == nil || handler == nil {
		return nil
	}
	if wp.app.GlobalShortcut == nil {
		return core.E("keybinding.WailsPlatform.Add", "Wails global shortcut manager is unavailable", nil)
	}
	if err := wp.app.GlobalShortcut.Register(accelerator, handler); err != nil {
		return core.E("keybinding.WailsPlatform.Add", "failed to register global shortcut", err)
	}
	wp.mu.Lock()
	wp.handlers[accelerator] = handler
	wp.mu.Unlock()
	return nil
}

// Remove unregisters by accelerator. No-op if absent.
func (wp *WailsPlatform) Remove(accelerator string) resultFailure {
	if wp == nil || wp.app == nil {
		return nil
	}
	if wp.app.GlobalShortcut == nil {
		return core.E("keybinding.WailsPlatform.Remove", "Wails global shortcut manager is unavailable", nil)
	}
	if err := wp.app.GlobalShortcut.Unregister(accelerator); err != nil {
		return core.E("keybinding.WailsPlatform.Remove", "failed to unregister global shortcut", err)
	}
	wp.mu.Lock()
	delete(wp.handlers, accelerator)
	wp.mu.Unlock()
	return nil
}

// Process triggers the registered handler programmatically.
func (wp *WailsPlatform) Process(accelerator string) bool {
	if wp == nil {
		return false
	}
	wp.mu.RLock()
	handler, ok := wp.handlers[accelerator]
	wp.mu.RUnlock()
	if !ok || handler == nil {
		return false
	}
	handler()
	return true
}

// GetAll returns Wails's canonical, sorted accelerator strings.
func (wp *WailsPlatform) GetAll() []string {
	if wp == nil || wp.app == nil || wp.app.GlobalShortcut == nil {
		return nil
	}
	return wp.app.GlobalShortcut.GetAll()
}
