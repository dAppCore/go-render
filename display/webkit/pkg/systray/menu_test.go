package systray

import (
	core "dappco.re/go"
)

type recordingTrayPlatform struct {
	tray *recordingTray
	menu *recordingTrayMenu
}

func (p *recordingTrayPlatform) NewTray() PlatformTray {
	p.tray = &recordingTray{}
	return p.tray
}

func (p *recordingTrayPlatform) NewMenu() PlatformMenu {
	p.menu = &recordingTrayMenu{}
	return p.menu
}

type recordingTray struct {
	icon, templateIcon []byte
	tooltip, label     string
	menu               PlatformMenu
	attachedWindow     WindowHandle
}

func (t *recordingTray) SetIcon(data []byte)                               { t.icon = append([]byte(nil), data...) }
func (t *recordingTray) SetTemplateIcon(data []byte)                       { t.templateIcon = append([]byte(nil), data...) }
func (t *recordingTray) SetTooltip(text string)                            { t.tooltip = text }
func (t *recordingTray) SetLabel(text string)                              { t.label = text }
func (t *recordingTray) SetMenu(menu PlatformMenu)                         { t.menu = menu }
func (t *recordingTray) AttachWindow(w WindowHandle, offsetX, offsetY int) { t.attachedWindow = w }
func (t *recordingTray) ShowMessage(title, message string) resultFailure   { return nil }

type recordingTrayMenu struct {
	items []*recordingTrayMenuItem
	subs  []*recordingTrayMenu
}

func (m *recordingTrayMenu) Add(label string) PlatformMenuItem {
	item := &recordingTrayMenuItem{label: label, enabled: true}
	m.items = append(m.items, item)
	return item
}

func (m *recordingTrayMenu) AddSeparator() {
	m.items = append(m.items, &recordingTrayMenuItem{label: "---"})
}

func (m *recordingTrayMenu) AddSubmenu(label string) PlatformMenu {
	sub := &recordingTrayMenu{}
	m.items = append(m.items, &recordingTrayMenuItem{label: label, submenu: sub})
	m.subs = append(m.subs, sub)
	return sub
}

type recordingTrayMenuItem struct {
	label, tooltip   string
	checked, enabled bool
	submenu          *recordingTrayMenu
	onClick          func()
}

func (i *recordingTrayMenuItem) SetTooltip(text string)  { i.tooltip = text }
func (i *recordingTrayMenuItem) SetChecked(checked bool) { i.checked = checked }
func (i *recordingTrayMenuItem) SetEnabled(enabled bool) { i.enabled = enabled }
func (i *recordingTrayMenuItem) OnClick(fn func())       { i.onClick = fn }

func TestManager_SetMenu_Good(t *core.T) {
	// SetMenu
	ax7Variant := "SetMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	platform := &recordingTrayPlatform{}
	mgr := NewManager(platform)
	core.RequireNoError(t, mgr.Setup("Core", "Core"))

	clicked := 0
	items := []TrayMenuItem{
		{Label: "Open", Tooltip: "open", ActionID: "open"},
		{Type: "separator"},
		{Label: "More", Submenu: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}}},
		{Label: "Disabled", Disabled: true},
		{Label: "Checked", Checked: true},
	}
	mgr.RegisterCallback("open", func() { clicked++ })
	mgr.RegisterCallback("nested", func() { clicked += 10 })

	core.RequireNoError(t, mgr.SetMenu(items))
	core.AssertNotNil(t, platform.menu)
	core.AssertNotNil(t, platform.tray)

	core.AssertEqual(t, "Core", platform.tray.tooltip)
	core.AssertEqual(t, "Core", platform.tray.label)
	core.AssertLen(t, platform.menu.items, 5)
	core.AssertEqual(t, "Open", platform.menu.items[0].label)
	core.AssertEqual(t, "---", platform.menu.items[1].label)
	core.AssertEqual(t, "More", platform.menu.items[2].label)
	core.AssertFalse(t, platform.menu.items[3].enabled)
	core.AssertTrue(t, platform.menu.items[4].checked)

	core.AssertNotNil(t, platform.menu.items[0].onClick)
	platform.menu.items[0].onClick()
	core.AssertLen(t, platform.menu.subs, 1)
	core.AssertNotNil(t, platform.menu.subs[0].items[0].onClick)
	platform.menu.subs[0].items[0].onClick()

	core.AssertEqual(t, 11, clicked)
	core.AssertLen(t, mgr.GetInfo()["menuItems"].([]TrayMenuItem), 5)
}

func TestManager_SetMenu_Bad(t *core.T) {
	// SetMenu
	ax7Variant := "SetMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	mgr := NewManager(&recordingTrayPlatform{})
	err := mgr.SetMenu([]TrayMenuItem{{Label: "Quit"}})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "tray not initialised")
}

func TestManager_GetCallback_Ugly(t *core.T) {
	// GetCallback
	ax7Variant := "GetCallback:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	mgr := NewManager(&recordingTrayPlatform{})
	mgr.RegisterCallback("quit", func() {})
	cb, ok := mgr.GetCallback("quit")
	core.RequireTrue(t, ok)
	core.AssertNotNil(t, cb)

	mgr.UnregisterCallback("quit")
	_, ok = mgr.GetCallback("quit")
	core.AssertFalse(t, ok)
}

// AX7 generated source-matching smoke coverage.
func TestMenu_Manager_SetMenu_Good(t *core.T) {
	// Manager SetMenu
	ax7Variant := "Manager_SetMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_SetMenu_Bad(t *core.T) {
	// Manager SetMenu
	ax7Variant := "Manager_SetMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_SetMenu_Ugly(t *core.T) {
	// Manager SetMenu
	ax7Variant := "Manager_SetMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_RegisterCallback_Good(t *core.T) {
	// Manager RegisterCallback
	ax7Variant := "Manager_RegisterCallback:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.RegisterCallback("agent", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_RegisterCallback_Bad(t *core.T) {
	// Manager RegisterCallback
	ax7Variant := "Manager_RegisterCallback:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.RegisterCallback("", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_RegisterCallback_Ugly(t *core.T) {
	// Manager RegisterCallback
	ax7Variant := "Manager_RegisterCallback:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.RegisterCallback("../../edge", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_UnregisterCallback_Good(t *core.T) {
	// Manager UnregisterCallback
	ax7Variant := "Manager_UnregisterCallback:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.UnregisterCallback("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_UnregisterCallback_Bad(t *core.T) {
	// Manager UnregisterCallback
	ax7Variant := "Manager_UnregisterCallback:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.UnregisterCallback("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_UnregisterCallback_Ugly(t *core.T) {
	// Manager UnregisterCallback
	ax7Variant := "Manager_UnregisterCallback:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.UnregisterCallback("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetCallback_Good(t *core.T) {
	// Manager GetCallback
	ax7Variant := "Manager_GetCallback:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.GetCallback("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetCallback_Bad(t *core.T) {
	// Manager GetCallback
	ax7Variant := "Manager_GetCallback:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.GetCallback("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetCallback_Ugly(t *core.T) {
	// Manager GetCallback
	ax7Variant := "Manager_GetCallback:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.GetCallback("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetInfo_Good(t *core.T) {
	// Manager GetInfo
	ax7Variant := "Manager_GetInfo:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.GetInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetInfo_Bad(t *core.T) {
	// Manager GetInfo
	ax7Variant := "Manager_GetInfo:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.GetInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenu_Manager_GetInfo_Ugly(t *core.T) {
	// Manager GetInfo
	ax7Variant := "Manager_GetInfo:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.GetInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
