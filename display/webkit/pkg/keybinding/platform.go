// pkg/keybinding/platform.go
package keybinding

// Platform abstracts the keyboard shortcut backend (Wails v3).
type Platform interface {
	// Add registers a global keyboard shortcut with the given accelerator string.
	// The handler is called when the shortcut is triggered.
	// Accelerator syntax is platform-aware: "Cmd+S" (macOS), "Ctrl+S" (Windows/Linux).
	// Special keys: F1-F12, Escape, Enter, Space, Tab, Backspace, Delete, arrow keys.
	Add(accelerator string, handler func()) resultFailure

	// Remove unregisters a previously registered keyboard shortcut.
	Remove(accelerator string) resultFailure

	// Process triggers the registered handler for the given accelerator programmatically.
	// Returns true if a handler was found and invoked, false if not registered.
	//
	//	handled := platform.Process("Ctrl+S")
	Process(accelerator string) bool

	// GetAll returns all currently registered accelerator strings.
	// Used for adapter-level reconciliation only — not read by QueryList.
	GetAll() []string
}
