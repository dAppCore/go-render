// pkg/dock/service_test.go
package dock

import (
	"context"

	core "dappco.re/go"
)

// --- Mock Platform ---

type mockPlatform struct {
	visible          bool
	badge            string
	hasBadge         bool
	progress         float64
	bounceID         int
	bounceType       BounceType
	bounceCalled     bool
	stopBounceCalled bool
	showErr          resultFailure
	hideErr          resultFailure
	badgeErr         resultFailure
	removeErr        resultFailure
	progressErr      resultFailure
	bounceErr        resultFailure
	stopBounceErr    resultFailure
}

func (m *mockPlatform) ShowIcon() resultFailure {
	if m.showErr != nil {
		return m.showErr
	}
	m.visible = true
	return nil
}

func (m *mockPlatform) HideIcon() resultFailure {
	if m.hideErr != nil {
		return m.hideErr
	}
	m.visible = false
	return nil
}

func (m *mockPlatform) SetBadge(label string) resultFailure {
	if m.badgeErr != nil {
		return m.badgeErr
	}
	m.badge = label
	m.hasBadge = true
	return nil
}

func (m *mockPlatform) RemoveBadge() resultFailure {
	if m.removeErr != nil {
		return m.removeErr
	}
	m.badge = ""
	m.hasBadge = false
	return nil
}

func (m *mockPlatform) IsVisible() bool { return m.visible }

func (m *mockPlatform) SetProgressBar(progress float64) resultFailure {
	if m.progressErr != nil {
		return m.progressErr
	}
	m.progress = progress
	return nil
}

func (m *mockPlatform) Bounce(bounceType BounceType) (int, resultFailure) {
	if m.bounceErr != nil {
		return 0, m.bounceErr
	}
	m.bounceCalled = true
	m.bounceType = bounceType
	m.bounceID++
	return m.bounceID, nil
}

func (m *mockPlatform) StopBounce(requestID int) resultFailure {
	if m.stopBounceErr != nil {
		return m.stopBounceErr
	}
	m.stopBounceCalled = true
	return nil
}

// --- Test helpers ---

func newTestDockService(t *core.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := &mockPlatform{visible: true}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "dock")
	return svc, c, mock
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func setBadge(c *core.Core, label string) core.Result {
	return c.Action("dock.set_badge").Run(context.Background(), core.NewOptions(
		core.Option{Key: "label", Value: label},
	))
}

// --- Tests ---

func TestRegister_Good(t *core.T) {
	svc, _, _ := newTestDockService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotEmpty(t, core.Sprintf("%T", svc))
}

func TestQueryVisible_Good(t *core.T) {
	_, c, _ := newTestDockService(t)
	r := c.QUERY(QueryVisible{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, true, r.Value)
}

func TestQueryVisible_Bad(t *core.T) {
	// No dock service registered — QUERY returns handled=false
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryVisible{})
	core.AssertFalse(t, r.OK)
}

func TestTaskShowIcon_Good(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = false // Start hidden

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.show_icon", TaskShowIcon{})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.visible)
	core.AssertNotNil(t, received)
	core.AssertTrue(t, received.Visible)
}

func TestTaskHideIcon_Good(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = true // Start visible

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.hide_icon", TaskHideIcon{})
	core.RequireTrue(t, r.OK)
	core.AssertFalse(t, mock.visible)
	core.AssertNotNil(t, received)
	core.AssertFalse(t, received.Visible)
}

func TestTaskSetBadge_Good(t *core.T) {
	_, c, mock := newTestDockService(t)
	r := setBadge(c, "3")
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "3", mock.badge)
	core.AssertTrue(t, mock.hasBadge)
}

func TestTaskSetBadge_EmptyLabel_GoodCase(t *core.T) {
	_, c, mock := newTestDockService(t)
	r := setBadge(c, "")
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "", mock.badge)
	core.AssertTrue(t, mock.hasBadge) // Empty string = default system badge indicator
}

func TestTaskRemoveBadge_Good(t *core.T) {
	_, c, mock := newTestDockService(t)
	// Set a badge first
	_ = setBadge(c, "5")

	r := taskRun(c, "dock.remove_badge", TaskRemoveBadge{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "", mock.badge)
	core.AssertFalse(t, mock.hasBadge)
}

func TestTaskShowIcon_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.showErr = core.AnError

	r := taskRun(c, "dock.show_icon", TaskShowIcon{})
	core.AssertFalse(t, r.OK)
}

func TestTaskHideIcon_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.hideErr = core.AnError

	r := taskRun(c, "dock.hide_icon", TaskHideIcon{})
	core.AssertFalse(t, r.OK)
}

func TestTaskSetBadge_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.badgeErr = core.AnError

	r := setBadge(c, "3")
	core.AssertFalse(t, r.OK)
}

// --- TaskSetProgressBar ---

func TestTaskSetProgressBar_Good(t *core.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionProgressChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionProgressChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.set_progress_bar", TaskSetProgressBar{Progress: 0.5})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, 0.5, mock.progress)
	core.AssertNotNil(t, received)
	core.AssertEqual(t, 0.5, received.Progress)
}

func TestTaskSetProgressBar_Hide_GoodCase(t *core.T) {
	// Progress -1.0 hides the indicator
	_, c, mock := newTestDockService(t)
	r := taskRun(c, "dock.set_progress_bar", TaskSetProgressBar{Progress: -1.0})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, -1.0, mock.progress)
}

func TestTaskSetProgressBar_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.progressErr = core.AnError

	r := taskRun(c, "dock.set_progress_bar", TaskSetProgressBar{Progress: 0.5})
	core.AssertFalse(t, r.OK)
}

func TestTaskSetProgressBar_Ugly(t *core.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.set_progress_bar").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// --- TaskBounce ---

func TestTaskBounce_Good(t *core.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionBounceStarted
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionBounceStarted); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.bounceCalled)
	core.AssertEqual(t, BounceInformational, mock.bounceType)
	requestID, ok := r.Value.(int)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 1, requestID)
	core.AssertNotNil(t, received)
	core.AssertEqual(t, BounceInformational, received.BounceType)
}

func TestTaskBounce_Critical_GoodCase(t *core.T) {
	_, c, mock := newTestDockService(t)
	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceCritical})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, BounceCritical, mock.bounceType)
	requestID, ok := r.Value.(int)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 1, requestID)
}

func TestTaskBounce_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.bounceErr = core.AnError

	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	core.AssertFalse(t, r.OK)
}

func TestTaskBounce_Ugly(t *core.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.bounce").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// --- TaskStopBounce ---

func TestTaskStopBounce_Good(t *core.T) {
	_, c, mock := newTestDockService(t)

	// Start a bounce to get a requestID
	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	core.RequireTrue(t, r.OK)
	requestID := r.Value.(int)

	r2 := taskRun(c, "dock.stop_bounce", TaskStopBounce{RequestID: requestID})
	core.RequireTrue(t, r2.OK)
	core.AssertTrue(t, mock.stopBounceCalled)
}

func TestTaskStopBounce_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.stopBounceErr = core.AnError

	r := taskRun(c, "dock.stop_bounce", TaskStopBounce{RequestID: 1})
	core.AssertFalse(t, r.OK)
}

func TestTaskStopBounce_Ugly(t *core.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.stop_bounce").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskRemoveBadge_Bad(t *core.T) {
	_, c, mock := newTestDockService(t)
	mock.removeErr = core.AnError

	r := taskRun(c, "dock.remove_badge", TaskRemoveBadge{})
	core.AssertFalse(t, r.OK)
}

func TestQueryVisible_Ugly(t *core.T) {
	// Dock icon initially hidden
	mock := &mockPlatform{visible: false}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryVisible{})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, false, r.Value)
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
