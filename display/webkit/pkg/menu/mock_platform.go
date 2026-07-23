// pkg/menu/mock_platform.go
package menu

// MockPlatform is an exported mock for cross-package integration tests.
type MockPlatform struct{}

func NewMockPlatform() *MockPlatform { return &MockPlatform{} }

func (m *MockPlatform) NewMenu() PlatformMenu                { return &exportedMockPlatformMenu{} }
func (m *MockPlatform) SetApplicationMenu(menu PlatformMenu) {}

type exportedMockPlatformMenu struct{}

func (m *exportedMockPlatformMenu) Add(label string) PlatformMenuItem {
	return &exportedMockPlatformMenuItem{}
}
func (m *exportedMockPlatformMenu) AddSeparator() {}
func (m *exportedMockPlatformMenu) AddSubmenu(label string) PlatformMenu {
	return &exportedMockPlatformMenu{}
}
func (m *exportedMockPlatformMenu) AddRole(role MenuRole) {}

type exportedMockPlatformMenuItem struct{}

func (mi *exportedMockPlatformMenuItem) SetAccelerator(acc string) PlatformMenuItem { return mi }
func (mi *exportedMockPlatformMenuItem) SetTooltip(tip string) PlatformMenuItem     { return mi }
func (mi *exportedMockPlatformMenuItem) SetChecked(checked bool) PlatformMenuItem   { return mi }
func (mi *exportedMockPlatformMenuItem) SetEnabled(enabled bool) PlatformMenuItem   { return mi }
func (mi *exportedMockPlatformMenuItem) OnClick(fn func()) PlatformMenuItem         { return mi }
