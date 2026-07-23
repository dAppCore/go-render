// pkg/systray/mock_test.go
package systray

type mockPlatform struct {
	trays []*mockTray
	menus []*mockTrayMenu
}

func newMockPlatform() *mockPlatform { return &mockPlatform{} }

func (p *mockPlatform) NewTray() PlatformTray {
	t := &mockTray{}
	p.trays = append(p.trays, t)
	return t
}

func (p *mockPlatform) NewMenu() PlatformMenu {
	m := &mockTrayMenu{}
	p.menus = append(p.menus, m)
	return m
}

type mockTrayMenu struct {
	items []string
	subs  []*mockTrayMenu
}

func (m *mockTrayMenu) Add(label string) PlatformMenuItem {
	m.items = append(m.items, label)
	return &mockTrayMenuItem{}
}
func (m *mockTrayMenu) AddSeparator() { m.items = append(m.items, "---") }
func (m *mockTrayMenu) AddSubmenu(label string) PlatformMenu {
	m.items = append(m.items, label)
	sub := &mockTrayMenu{}
	m.subs = append(m.subs, sub)
	return sub
}

type mockTrayMenuItem struct{}

func (mi *mockTrayMenuItem) SetTooltip(text string)  {}
func (mi *mockTrayMenuItem) SetChecked(checked bool) {}
func (mi *mockTrayMenuItem) SetEnabled(enabled bool) {}
func (mi *mockTrayMenuItem) OnClick(fn func())       {}

type mockTray struct {
	icon, templateIcon []byte
	tooltip, label     string
	menu               PlatformMenu
	attachedWindow     WindowHandle
	lastMessageTitle   string
	lastMessageBody    string
}

func (t *mockTray) SetIcon(data []byte)                               { t.icon = data }
func (t *mockTray) SetTemplateIcon(data []byte)                       { t.templateIcon = data }
func (t *mockTray) SetTooltip(text string)                            { t.tooltip = text }
func (t *mockTray) SetLabel(text string)                              { t.label = text }
func (t *mockTray) SetMenu(menu PlatformMenu)                         { t.menu = menu }
func (t *mockTray) AttachWindow(w WindowHandle, offsetX, offsetY int) { t.attachedWindow = w }
func (t *mockTray) ShowMessage(title, message string) resultFailure {
	t.lastMessageTitle = title
	t.lastMessageBody = message
	return nil
}
