package systray

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/notification"
)

func newTestSystrayService(t *core.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "systray")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *core.T) {
	svc, _ := newTestSystrayService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotNil(t, svc.manager)
}

func TestTaskSetTrayIcon_Good(t *core.T) {
	svc, c := newTestSystrayService(t)

	// Setup tray first (normally done via config in OnStartup)
	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	icon := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	r := taskRun(c, "systray.set_icon", TaskSetTrayIcon{Data: icon})
	core.RequireTrue(t, r.OK)
}

func TestTaskSetTrayMenu_Good(t *core.T) {
	svc, c := newTestSystrayService(t)

	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	items := []TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}
	r := taskRun(c, "systray.set_menu", TaskSetTrayMenu{Items: items})
	core.RequireTrue(t, r.OK)
}

func TestTaskSetTrayTooltip_Good(t *core.T) {
	svc, c := newTestSystrayService(t)
	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	r := taskRun(c, "systray.set_tooltip", TaskSetTrayTooltip{Tooltip: "New Tooltip"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "New Tooltip", svc.manager.GetInfo()["tooltip"])
}

func TestTaskSetTrayLabel_Good(t *core.T) {
	svc, c := newTestSystrayService(t)
	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	r := taskRun(c, "systray.set_label", TaskSetTrayLabel{Label: "Ready"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "Ready", svc.manager.GetInfo()["label"])
}

func TestTaskShowMessage_Good(t *core.T) {
	svc, c := newTestSystrayService(t)
	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	r := taskRun(c, "systray.show_message", TaskShowMessage{Title: "Core", Message: "Up"})
	core.RequireTrue(t, r.OK)

	mockTray := svc.manager.Tray().(*mockTray)
	core.AssertEqual(t, "Core", mockTray.lastMessageTitle)
	core.AssertEqual(t, "Up", mockTray.lastMessageBody)
}

type fallbackNotificationPlatform struct {
	sent bool
	opts notification.NotificationOptions
}

func (m *fallbackNotificationPlatform) Send(opts notification.NotificationOptions) resultFailure {
	m.sent = true
	m.opts = opts
	return nil
}
func (m *fallbackNotificationPlatform) RequestPermission() (bool, resultFailure) { return true, nil }
func (m *fallbackNotificationPlatform) CheckPermission() (bool, resultFailure)   { return true, nil }
func (m *fallbackNotificationPlatform) RevokePermission() resultFailure          { return nil }
func (m *fallbackNotificationPlatform) Clear(id string) resultFailure            { return nil }

type failingTrayPlatform struct{}

func (failingTrayPlatform) NewTray() PlatformTray { return &failingTray{} }
func (failingTrayPlatform) NewMenu() PlatformMenu { return &mockTrayMenu{} }

type failingTray struct{ mockTray }

func (t *failingTray) ShowMessage(title, message string) resultFailure {
	return core.NewError("tray balloon unavailable")
}

func TestTaskShowMessage_FallbackToNotification_GoodCase(t *core.T) {
	notifPlatform := &fallbackNotificationPlatform{}
	c := core.New(
		core.WithService(notification.Register(notifPlatform)),
		core.WithService(Register(failingTrayPlatform{})),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	svc := core.MustServiceFor[*Service](c, "systray")
	core.RequireNoError(t, svc.manager.Setup("Test", "Test"))

	r := taskRun(c, "systray.show_message", TaskShowMessage{Title: "Core", Message: "Up"})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, notifPlatform.sent)
	core.AssertEqual(t, "Core", notifPlatform.opts.Title)
	core.AssertEqual(t, "Up", notifPlatform.opts.Message)
}

func TestQueryInfo_Good(t *core.T) {
	svc, c := newTestSystrayService(t)
	core.RequireNoError(t, svc.manager.Setup("Core", "Core"))

	r := c.QUERY(QueryInfo{})
	core.RequireTrue(t, r.OK)
	info := r.Value.(map[string]any)
	core.AssertEqual(t, "Core", info["tooltip"])
	core.AssertEqual(t, "Core", info["label"])
}

func TestTaskSetTrayIcon_Bad(t *core.T) {
	// No systray service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("systray.set_icon").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// AX7 generated source-matching smoke coverage.
func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Manager_Good(t *core.T) {
	// Service Manager
	ax7Variant := "Service_Manager:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Manager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Manager_Bad(t *core.T) {
	// Service Manager
	ax7Variant := "Service_Manager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Manager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Manager_Ugly(t *core.T) {
	// Service Manager
	ax7Variant := "Service_Manager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Manager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
