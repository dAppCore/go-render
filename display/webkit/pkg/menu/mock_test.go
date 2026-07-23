// pkg/menu/mock_test.go
package menu

type mockPlatform struct {
	menus   []*mockMenu
	appMenu PlatformMenu
}

func newMockPlatform() *mockPlatform { return &mockPlatform{} }

func (p *mockPlatform) NewMenu() PlatformMenu {
	m := &mockMenu{}
	p.menus = append(p.menus, m)
	return m
}

func (p *mockPlatform) SetApplicationMenu(menu PlatformMenu) { p.appMenu = menu }

type mockMenu struct {
	items []*mockMenuItem
	subs  []*mockMenu
	roles []MenuRole
}

func (m *mockMenu) Add(label string) PlatformMenuItem {
	mi := &mockMenuItem{label: label}
	m.items = append(m.items, mi)
	return mi
}

func (m *mockMenu) AddSeparator() {
	m.items = append(m.items, &mockMenuItem{label: "---"})
}

func (m *mockMenu) AddSubmenu(label string) PlatformMenu {
	sub := &mockMenu{}
	m.subs = append(m.subs, sub)
	m.items = append(m.items, &mockMenuItem{label: label, isSubmenu: true})
	return sub
}

func (m *mockMenu) AddRole(role MenuRole) { m.roles = append(m.roles, role) }

type mockMenuItem struct {
	label, accel, tooltip string
	checked, enabled      bool
	isSubmenu             bool
	onClick               func()
}

func (mi *mockMenuItem) SetAccelerator(accel string) PlatformMenuItem {
	mi.accel = accel
	return mi
}
func (mi *mockMenuItem) SetTooltip(text string) PlatformMenuItem {
	mi.tooltip = text
	return mi
}
func (mi *mockMenuItem) SetChecked(checked bool) PlatformMenuItem {
	mi.checked = checked
	return mi
}
func (mi *mockMenuItem) SetEnabled(enabled bool) PlatformMenuItem {
	mi.enabled = enabled
	return mi
}
func (mi *mockMenuItem) OnClick(fn func()) PlatformMenuItem {
	mi.onClick = fn
	return mi
}
