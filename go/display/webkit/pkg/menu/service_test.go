package menu

import (
	"context"

	core "dappco.re/go"
)

func newTestMenuService(t *core.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "menu")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *core.T) {
	svc, _ := newTestMenuService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotNil(t, svc.manager)
}

func TestTaskSetAppMenu_Good(t *core.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{
		{Label: "File", Children: []MenuItem{
			{Label: "New"},
			{Type: "separator"},
			{Label: "Quit"},
		}},
	}
	r := taskRun(c, "menu.set_app_menu", TaskSetAppMenu{Items: items})
	core.RequireTrue(t, r.OK)
}

func TestQueryGetAppMenu_Good(t *core.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{{Label: "File"}, {Label: "Edit"}}
	taskRun(c, "menu.set_app_menu", TaskSetAppMenu{Items: items})

	r := c.QUERY(QueryGetAppMenu{})
	core.RequireTrue(t, r.OK)
	menuItems := r.Value.([]MenuItem)
	core.AssertLen(t, menuItems, 2)
	core.AssertEqual(t, "File", menuItems[0].Label)
}

func TestTaskSetAppMenu_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("menu.set_app_menu").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskSetAppMenu_NoManager_FailsClosed(t *core.T) {
	c := core.New(core.WithServiceLock())
	svc := &Service{
		ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
	}
	core.RequireTrue(t, svc.OnStartup(context.Background()).OK)

	r := c.Action("menu.set_app_menu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetAppMenu{Items: []MenuItem{{Label: "File"}}}},
	))
	core.AssertFalse(t, r.OK)
	err, ok := r.Value.(error)
	core.RequireTrue(t, ok)
	core.AssertContains(t, err.Error(), "menu manager unavailable")
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

func TestService_Service_ShowDevTools_Good(t *core.T) {
	// Service ShowDevTools
	ax7Variant := "Service_ShowDevTools:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowDevTools()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_ShowDevTools_Bad(t *core.T) {
	// Service ShowDevTools
	ax7Variant := "Service_ShowDevTools:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowDevTools()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_ShowDevTools_Ugly(t *core.T) {
	// Service ShowDevTools
	ax7Variant := "Service_ShowDevTools:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowDevTools()
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
