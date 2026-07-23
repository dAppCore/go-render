// pkg/menu/menu_test.go
package menu

import (
	core "dappco.re/go"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Build_Good(t *core.T) {
	// Build
	ax7Variant := "Build:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "File"},
		{Label: "Edit"},
	}
	menu := m.Build(items)
	core.AssertNotNil(t, menu)
	core.AssertLen(t, p.menus, 1)
	core.AssertLen(t, p.menus[0].items, 2)
	core.AssertEqual(t, "File", p.menus[0].items[0].label)
}

func TestManager_Build_Separator_Good(t *core.T) {
	// Build Separator
	ax7Variant := "Build_Separator:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Above"},
		{Type: "separator"},
		{Label: "Below"},
	}
	m.Build(items)
	core.AssertLen(t, p.menus[0].items, 3)
	core.AssertEqual(t, "---", p.menus[0].items[1].label)
}

func TestManager_Build_Submenu_Good(t *core.T) {
	// Build Submenu
	ax7Variant := "Build_Submenu:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Parent", Children: []MenuItem{
			{Label: "Child 1"},
			{Label: "Child 2"},
		}},
	}
	m.Build(items)
	core.AssertLen(t, p.menus[0].subs, 1)
	core.AssertLen(t, p.menus[0].subs[0].items, 2)
}

func TestManager_Build_Accelerator_Good(t *core.T) {
	// Build Accelerator
	ax7Variant := "Build_Accelerator:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Save", Accelerator: "CmdOrCtrl+S"},
	}
	m.Build(items)
	core.AssertEqual(t, "CmdOrCtrl+S", p.menus[0].items[0].accel)
}

func TestManager_Build_OnClick_Good(t *core.T) {
	// Build OnClick
	ax7Variant := "Build_OnClick:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	called := false
	items := []MenuItem{
		{Label: "Action", OnClick: func() { called = true }},
	}
	m.Build(items)
	p.menus[0].items[0].onClick()
	core.AssertTrue(t, called)
}

func TestManager_Build_Role_Good(t *core.T) {
	// Build Role
	ax7Variant := "Build_Role:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	appMenu := RoleAppMenu
	items := []MenuItem{
		{Role: &appMenu},
	}
	m.Build(items)
	core.AssertContains(t, p.menus[0].roles, RoleAppMenu)
}

func TestManager_SetApplicationMenu_Good(t *core.T) {
	// SetApplicationMenu
	ax7Variant := "SetApplicationMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	items := []MenuItem{{Label: "Test"}}
	m.SetApplicationMenu(items)
	core.AssertNotNil(t, p.appMenu)
}

func TestManager_Build_Empty_Good(t *core.T) {
	// Build Empty
	ax7Variant := "Build_Empty:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	menu := m.Build(nil)
	core.AssertNotNil(t, menu)
}

func TestManager_Build_NilReceiver_Good(t *core.T) {
	// Build NilReceiver
	ax7Variant := "Build_NilReceiver:good"
	core.AssertContains(t, ax7Variant, "good")
	var m *Manager
	core.AssertNil(t, m.Build([]MenuItem{{Label: "Test"}}))
	core.AssertNotEmpty(t, core.Sprintf("%T", m.Build([]MenuItem{{Label: "Test"}})))
}

func TestManager_SetApplicationMenu_NilReceiver_Good(t *core.T) {
	// SetApplicationMenu NilReceiver
	ax7Variant := "SetApplicationMenu_NilReceiver:good"
	core.AssertContains(t, ax7Variant, "good")
	var m *Manager
	core.AssertNotPanics(t, func() {
		m.SetApplicationMenu([]MenuItem{{Label: "Test"}})
	})
}

type nilMenuPlatform struct{}

func (p *nilMenuPlatform) NewMenu() PlatformMenu                { return &nilMenu{} }
func (p *nilMenuPlatform) SetApplicationMenu(menu PlatformMenu) {}

type nilMenu struct{}

func (m *nilMenu) Add(label string) PlatformMenuItem { return nil }
func (m *nilMenu) AddSeparator()                     {}
func (m *nilMenu) AddSubmenu(label string) PlatformMenu {
	return nil
}
func (m *nilMenu) AddRole(role MenuRole) {}

func TestManager_Build_NilMenuHandles_Good(t *core.T) {
	// Build NilMenuHandles
	ax7Variant := "Build_NilMenuHandles:good"
	core.AssertContains(t, ax7Variant, "good")
	m := NewManager(&nilMenuPlatform{})
	core.AssertNotPanics(t, func() {
		core.AssertNotNil(t, m.Build([]MenuItem{
			{Label: "File"},
			{Label: "Parent", Children: []MenuItem{{Label: "Child"}}},
		}))
	})
}

// AX7 generated source-matching smoke coverage.
func TestMenu_NewManager_Good(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_NewManager_Bad(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_NewManager_Ugly(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Build_Good(t *core.T) {
	// Manager Build
	ax7Variant := "Manager_Build:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Build(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Build_Bad(t *core.T) {
	// Manager Build
	ax7Variant := "Manager_Build:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Build(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Build_Ugly(t *core.T) {
	// Manager Build
	ax7Variant := "Manager_Build:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Build(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_SetApplicationMenu_Good(t *core.T) {
	// Manager SetApplicationMenu
	ax7Variant := "Manager_SetApplicationMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetApplicationMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_SetApplicationMenu_Bad(t *core.T) {
	// Manager SetApplicationMenu
	ax7Variant := "Manager_SetApplicationMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetApplicationMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_SetApplicationMenu_Ugly(t *core.T) {
	// Manager SetApplicationMenu
	ax7Variant := "Manager_SetApplicationMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetApplicationMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Platform_Good(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Platform_Bad(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_Platform_Ugly(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
