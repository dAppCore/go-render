package contextmenu

import core "dappco.re/go"

var ErrorMenuNotFound = core.E("contextmenu", "menu not found", nil)

// QueryGet returns a named context menu definition. Result: *ContextMenuDef (nil if not found)
type QueryGet struct {
	Name string `json:"name"`
}

// QueryList returns all registered context menus. Result: map[string]ContextMenuDef
type QueryList struct{}

// QueryGetAll returns all registered context menus. Equivalent to QueryList.
// Result: map[string]ContextMenuDef
type QueryGetAll struct{}

// TaskAdd registers a named context menu. Replaces if already exists.
type TaskAdd struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskRemove unregisters a context menu by name. Error: ErrorMenuNotFound if missing.
type TaskRemove struct {
	Name string `json:"name"`
}

// TaskUpdate replaces an existing context menu's definition. Error: ErrorMenuNotFound if missing.
type TaskUpdate struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskDestroy removes a context menu and releases all associated resources.
// Error: ErrorMenuNotFound if missing.
type TaskDestroy struct {
	Name string `json:"name"`
}

// ActionItemClicked is broadcast when a context menu item is clicked.
type ActionItemClicked struct {
	MenuName string `json:"menuName"`
	ActionID string `json:"actionId"`
	Data     string `json:"data,omitempty"`
}
