// pkg/contextmenu/platform.go
package contextmenu

// Platform abstracts the context menu backend (Wails v3).
// The Add callback must broadcast ActionItemClicked via s.Core().ACTION()
// when a menu item is clicked — the adapter translates MenuItemDef.ActionID
// to a callback that does this.
type Platform interface {
	// Add registers a context menu by name.
	// The onItemClick callback is called with (menuName, actionID, data)
	// when any item in the menu is clicked. The adapter creates per-item
	// OnClick handlers that call this with the appropriate ActionID.
	Add(name string, menu ContextMenuDef, onItemClick func(menuName, actionID, data string)) resultFailure

	// Remove unregisters a context menu by name.
	Remove(name string) resultFailure

	// Get returns a context menu definition by name, or false if not found.
	Get(name string) (*ContextMenuDef, bool)

	// GetAll returns all registered context menu definitions.
	GetAll() map[string]ContextMenuDef
}

// ContextMenuDef describes a context menu and its items.
type ContextMenuDef struct {
	Name  string        `json:"name"`
	Items []MenuItemDef `json:"items"`
}

// MenuItemDef describes a single item in a context menu.
// Items may be nested (submenu children via Items field).
type MenuItemDef struct {
	Label       string        `json:"label"`
	Type        string        `json:"type,omitempty"` // "" (normal), "separator", "checkbox", "radio", "submenu"
	Accelerator string        `json:"accelerator,omitempty"`
	Enabled     *bool         `json:"enabled,omitempty"` // nil = true (default)
	Checked     bool          `json:"checked,omitempty"`
	ActionID    string        `json:"actionId,omitempty"` // identifies which item was clicked
	Items       []MenuItemDef `json:"items,omitempty"`    // submenu children (recursive)
}
