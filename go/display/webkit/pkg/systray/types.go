// pkg/systray/types.go
package systray

// TrayMenuItem describes a menu item for dynamic tray menus.
type TrayMenuItem struct {
	Label    string         `json:"label"`
	Type     string         `json:"type"` // "normal", "separator", "checkbox", "radio"
	Checked  bool           `json:"checked,omitempty"`
	Disabled bool           `json:"disabled,omitempty"`
	Tooltip  string         `json:"tooltip,omitempty"`
	Submenu  []TrayMenuItem `json:"submenu,omitempty"`
	ActionID string         `json:"action_id,omitempty"`
}
