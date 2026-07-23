// pkg/systray/mock_platform.go
package systray

// MockPlatform is an exported mock for cross-package integration tests.
// Use: platform := systray.NewMockPlatform()
type MockPlatform struct{}

// NewMockPlatform creates a tray platform mock.
// Use: platform := systray.NewMockPlatform()
func NewMockPlatform() *MockPlatform { return &MockPlatform{} }

// NewTray creates a mock tray handle for tests.
// Use: tray := platform.NewTray()
func (m *MockPlatform) NewTray() PlatformTray { return &exportedMockTray{} }

// NewMenu creates a mock tray menu for tests.
// Use: menu := platform.NewMenu()
func (m *MockPlatform) NewMenu() PlatformMenu { return &exportedMockMenu{} }

type exportedMockTray struct {
	icon, templateIcon []byte
	tooltip, label     string
}

func (t *exportedMockTray) SetIcon(data []byte)                               { t.icon = data }
func (t *exportedMockTray) SetTemplateIcon(data []byte)                       { t.templateIcon = data }
func (t *exportedMockTray) SetTooltip(text string)                            { t.tooltip = text }
func (t *exportedMockTray) SetLabel(text string)                              { t.label = text }
func (t *exportedMockTray) SetMenu(menu PlatformMenu)                         {}
func (t *exportedMockTray) AttachWindow(w WindowHandle, offsetX, offsetY int) {}
func (t *exportedMockTray) ShowMessage(title, message string) resultFailure   { return nil }

type exportedMockMenu struct {
	items []exportedMockMenuItem
	subs  []*exportedMockMenu
}

func (m *exportedMockMenu) Add(label string) PlatformMenuItem {
	mi := &exportedMockMenuItem{label: label}
	m.items = append(m.items, *mi)
	return mi
}
func (m *exportedMockMenu) AddSeparator() {}
func (m *exportedMockMenu) AddSubmenu(label string) PlatformMenu {
	m.items = append(m.items, exportedMockMenuItem{label: label})
	sub := &exportedMockMenu{}
	m.subs = append(m.subs, sub)
	return sub
}

type exportedMockMenuItem struct {
	label, tooltip   string
	checked, enabled bool
	onClick          func()
}

func (mi *exportedMockMenuItem) SetTooltip(tip string)   { mi.tooltip = tip }
func (mi *exportedMockMenuItem) SetChecked(checked bool) { mi.checked = checked }
func (mi *exportedMockMenuItem) SetEnabled(enabled bool) { mi.enabled = enabled }
func (mi *exportedMockMenuItem) OnClick(fn func())       { mi.onClick = fn }
