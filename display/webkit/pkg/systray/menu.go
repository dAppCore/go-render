// pkg/systray/menu.go
package systray

import core "dappco.re/go"

// SetMenu sets a dynamic menu on the tray from TrayMenuItem descriptors.
// Use: _ = m.SetMenu([]TrayMenuItem{{Label: "Quit", ActionID: "quit"}})
func (m *Manager) SetMenu(items []TrayMenuItem) resultFailure {
	if m.tray == nil {
		return core.E("systray.SetMenu", "tray not initialised", nil)
	}
	m.mu.Lock()
	m.menuItems = append([]TrayMenuItem(nil), items...)
	m.mu.Unlock()
	menu := m.platform.NewMenu()
	m.buildMenu(menu, items)
	m.tray.SetMenu(menu)
	return nil
}

// buildMenu recursively builds a PlatformMenu from TrayMenuItem descriptors.
func (m *Manager) buildMenu(menu PlatformMenu, items []TrayMenuItem) {
	for _, item := range items {
		if item.Type == "separator" {
			menu.AddSeparator()
			continue
		}
		if len(item.Submenu) > 0 {
			sub := menu.AddSubmenu(item.Label)
			m.buildMenu(sub, item.Submenu)
			continue
		}
		mi := menu.Add(item.Label)
		if item.Tooltip != "" {
			mi.SetTooltip(item.Tooltip)
		}
		if item.Disabled {
			mi.SetEnabled(false)
		}
		if item.Checked {
			mi.SetChecked(true)
		}
		if item.ActionID != "" {
			actionID := item.ActionID
			mi.OnClick(func() {
				if cb, ok := m.GetCallback(actionID); ok {
					cb()
				}
			})
		}
	}
}

// RegisterCallback registers a callback for a menu action ID.
// Use: m.RegisterCallback("quit", func() { _ = app.Quit() })
func (m *Manager) RegisterCallback(actionID string, callback func()) {
	m.mu.Lock()
	m.callbacks[actionID] = callback
	m.mu.Unlock()
}

// UnregisterCallback removes a callback.
// Use: m.UnregisterCallback("quit")
func (m *Manager) UnregisterCallback(actionID string) {
	m.mu.Lock()
	delete(m.callbacks, actionID)
	m.mu.Unlock()
}

// GetCallback returns the callback for an action ID.
// Use: callback, ok := m.GetCallback("quit")
func (m *Manager) GetCallback(actionID string) (func(), bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cb, ok := m.callbacks[actionID]
	return cb, ok
}

// GetInfo returns tray status information.
// Use: info := m.GetInfo()
func (m *Manager) GetInfo() map[string]any {
	return map[string]any{
		"active":          m.IsActive(),
		"tooltip":         m.tooltip,
		"label":           m.label,
		"hasIcon":         m.hasIcon,
		"hasTemplateIcon": m.hasTemplateIcon,
		"menuItems":       append([]TrayMenuItem(nil), m.menuItems...),
	}
}
