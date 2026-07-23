package keybinding

import core "dappco.re/go"

var ErrorAlreadyRegistered = core.E("keybinding", "accelerator already registered", nil)
var ErrorNotRegistered = core.E("keybinding", "accelerator not registered", nil)

// BindingInfo describes a registered global key binding.
type BindingInfo struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// QueryList returns all registered key bindings. Result: []BindingInfo
type QueryList struct{}

// TaskAdd registers a global key binding. Error: ErrorAlreadyRegistered if accelerator taken.
type TaskAdd struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// TaskRemove unregisters a global key binding by accelerator. Error: ErrorNotRegistered if not found.
type TaskRemove struct {
	Accelerator string `json:"accelerator"`
}

// TaskProcess triggers a registered key binding programmatically.
// Returns ActionTriggered if the accelerator was handled, ErrorNotRegistered if not found.
//
//	c.PERFORM(keybinding.TaskProcess{Accelerator: "Ctrl+S"})
type TaskProcess struct {
	Accelerator string `json:"accelerator"`
}

// ActionTriggered is broadcast when a registered key binding fires.
type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
