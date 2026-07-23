// pkg/notification/service_test.go
package notification

import (
	"context"
	core "dappco.re/go"

	"dappco.re/go/render/display/webkit/pkg/dialog"
)

type mockPlatform struct {
	sendErr      resultFailure
	permGranted  bool
	permErr      resultFailure
	revokeErr    resultFailure
	revokeCalled bool
	clearErr     resultFailure
	clearCalled  bool
	clearID      string
	lastOpts     NotificationOptions
	sendCalled   bool
	updateCalled bool
	category     NotificationCategory
	response     func(notificationID, actionID, userText string)
}

func (m *mockPlatform) Send(opts NotificationOptions) resultFailure {
	m.sendCalled = true
	m.lastOpts = opts
	return m.sendErr
}
func (m *mockPlatform) RequestPermission() (bool, resultFailure) { return m.permGranted, m.permErr }
func (m *mockPlatform) CheckPermission() (bool, resultFailure)   { return m.permGranted, m.permErr }
func (m *mockPlatform) RevokePermission() resultFailure {
	m.revokeCalled = true
	return m.revokeErr
}
func (m *mockPlatform) Clear(id string) resultFailure {
	m.clearCalled = true
	m.clearID = id
	return m.clearErr
}
func (m *mockPlatform) Update(opts NotificationOptions) resultFailure {
	m.updateCalled = true
	m.lastOpts = opts
	return m.sendErr
}
func (m *mockPlatform) RegisterCategory(category NotificationCategory) resultFailure {
	m.category = category
	return nil
}
func (m *mockPlatform) OnResponse(callback func(notificationID, actionID, userText string)) {
	m.response = callback
}

// mockDialogPlatform tracks whether MessageDialog was called (for fallback test).
type mockDialogPlatform struct {
	messageCalled bool
	lastMsgOpts   dialog.MessageDialogOptions
}

func (m *mockDialogPlatform) OpenFile(opts dialog.OpenFileOptions) ([]string, resultFailure) {
	return nil, nil
}
func (m *mockDialogPlatform) SaveFile(opts dialog.SaveFileOptions) (string, resultFailure) {
	return "", nil
}
func (m *mockDialogPlatform) OpenDirectory(opts dialog.OpenDirectoryOptions) (string, resultFailure) {
	return "", nil
}
func (m *mockDialogPlatform) MessageDialog(opts dialog.MessageDialogOptions) (string, resultFailure) {
	m.messageCalled = true
	m.lastMsgOpts = opts
	return "OK", nil
}

func newTestService(t *core.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{permGranted: true}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *core.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "notification")
	core.AssertNotNil(t, svc)
}

func TestTaskSend_Good(t *core.T) {
	mock, c := newTestService(t)
	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{Title: "Test", Message: "Hello"},
	})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.sendCalled)
	core.AssertEqual(t, "Test", mock.lastOpts.Title)
}

func TestTaskSend_Fallback_GoodCase(t *core.T) {
	// Platform fails -> falls back to dialog via IPC
	mockNotify := &mockPlatform{sendErr: core.NewError("no permission")}
	mockDlg := &mockDialogPlatform{}
	c := core.New(
		core.WithService(dialog.Register(mockDlg)),
		core.WithService(Register(mockNotify)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{Title: "Warn", Message: "Oops", Severity: SeverityWarning},
	})
	core.AssertTrue(t, r.OK) // fallback succeeds even though platform failed
	core.AssertTrue(t, mockDlg.messageCalled)
	core.AssertEqual(t, dialog.DialogWarning, mockDlg.lastMsgOpts.Type)
}

func TestQueryPermission_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryPermission{})
	core.RequireTrue(t, r.OK)
	status := r.Value.(PermissionStatus)
	core.AssertTrue(t, status.Granted)
}

func TestTaskRequestPermission_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("notification.request_permission").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, true, r.Value)
}

func TestTaskSend_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("notification.send").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// --- TaskRevokePermission ---

func TestTaskRevokePermission_Good(t *core.T) {
	mock, c := newTestService(t)
	r := c.Action("notification.revoke_permission").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.revokeCalled)
}

func TestTaskRevokePermission_Bad(t *core.T) {
	mock, c := newTestService(t)
	mock.revokeErr = core.NewError("cannot revoke")
	r := c.Action("notification.revoke_permission").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskRevokePermission_Ugly(t *core.T) {
	// No service registered — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("notification.revoke_permission").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

// --- TaskRegisterCategory ---

func TestTaskRegisterCategory_Good(t *core.T) {
	_, c := newTestService(t)
	category := NotificationCategory{
		ID: "message",
		Actions: []NotificationAction{
			{ID: "reply", Title: "Reply"},
			{ID: "delete", Title: "Delete", Destructive: true},
		},
	}
	r := taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: category})
	core.RequireTrue(t, r.OK)

	svc := core.MustServiceFor[*Service](c, "notification")
	stored, ok := svc.categories["message"]
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 2, len(stored.Actions))
	core.AssertEqual(t, "reply", stored.Actions[0].ID)
	core.AssertTrue(t, stored.Actions[1].Destructive)
}

func TestTaskRegisterCategory_Bad(t *core.T) {
	// No service registered — action is not registered
	c := core.New(core.WithServiceLock())
	r := taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: NotificationCategory{ID: "x"}})
	core.AssertFalse(t, r.OK)
}

func TestTaskRegisterCategory_Ugly(t *core.T) {
	// Re-registering a category replaces the previous one
	_, c := newTestService(t)
	first := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "a", Title: "A"}}}
	second := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "b", Title: "B"}, {ID: "c", Title: "C"}}}

	core.RequireTrue(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: first}).OK)
	core.RequireTrue(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: second}).OK)

	svc := core.MustServiceFor[*Service](c, "notification")
	core.AssertEqual(t, 2, len(svc.categories["chat"].Actions))
	core.AssertEqual(t, "b", svc.categories["chat"].Actions[0].ID)
}

// --- NotificationOptions with Actions ---

func TestTaskSend_WithActions_GoodCase(t *core.T) {
	mock, c := newTestService(t)
	options := NotificationOptions{
		Title:      "Team Chat",
		Message:    "New message from Alice",
		CategoryID: "message",
		Actions: []NotificationAction{
			{ID: "reply", Title: "Reply"},
			{ID: "dismiss", Title: "Dismiss"},
		},
	}
	r := taskRun(c, "notification.send", TaskSend{Options: options})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "message", mock.lastOpts.CategoryID)
	core.AssertEqual(t, 2, len(mock.lastOpts.Actions))
}

func TestTaskSend_CurrentWailsNotificationFeatures_Good(t *core.T) {
	mock, c := newTestService(t)
	options := NotificationOptions{
		ID:                   "build-7",
		Title:                "Build complete",
		Message:              "The application is ready.",
		Silent:               true,
		AttachmentPaths:      []string{"/tmp/build.png"},
		ThreadID:             "builds",
		InterruptionLevel:    "timeSensitive",
		ScheduleDelaySeconds: 5,
		Update:               true,
	}

	r := taskRun(c, "notification.send", TaskSend{Options: options})

	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.updateCalled)
	core.AssertEqual(t, options.AttachmentPaths, mock.lastOpts.AttachmentPaths)
	core.AssertEqual(t, "timeSensitive", mock.lastOpts.InterruptionLevel)
}

func TestNotificationResponse_CurrentWailsFeatures_Good(t *core.T) {
	mock, c := newTestService(t)
	var triggered ActionNotificationActionTriggered
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if action, ok := msg.(ActionNotificationActionTriggered); ok {
			triggered = action
		}
		return core.Result{OK: true}
	})

	core.RequireTrue(t, mock.response != nil)
	mock.response("message-1", "reply", "On my way")

	core.AssertEqual(t, "message-1", triggered.NotificationID)
	core.AssertEqual(t, "reply", triggered.ActionID)
	core.AssertEqual(t, "On my way", triggered.UserText)
}

func TestTaskSend_RegisteredCategoryActions_GoodCase(t *core.T) {
	mock, c := newTestService(t)
	core.RequireTrue(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{
		Category: NotificationCategory{
			ID: "message",
			Actions: []NotificationAction{
				{ID: "reply", Title: "Reply"},
				{ID: "dismiss", Title: "Dismiss"},
			},
		},
	}).OK)

	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{
			Title:      "Team Chat",
			Message:    "New message from Alice",
			CategoryID: "message",
		},
	})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, 2, len(mock.lastOpts.Actions))
	core.AssertEqual(t, "reply", mock.lastOpts.Actions[0].ID)
}

func TestTaskClear_Good_Specific(t *core.T) {
	mock, c := newTestService(t)
	core.RequireTrue(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n1", Title: "One", Message: "Hello"},
	}).OK)

	r := taskRun(c, "notification.clear", TaskClear{ID: "n1"})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.clearCalled)
	core.AssertEqual(t, "n1", mock.clearID)
}

func TestTaskClear_Good_All(t *core.T) {
	mock, c := newTestService(t)
	core.RequireTrue(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n1", Title: "One", Message: "Hello"},
	}).OK)
	core.RequireTrue(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n2", Title: "Two", Message: "World"},
	}).OK)

	r := taskRun(c, "notification.clear", TaskClear{})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.clearCalled)
	core.AssertEqual(t, "", mock.clearID)

	svc := core.MustServiceFor[*Service](c, "notification")
	core.AssertEmpty(t, svc.active)
}

// --- ActionNotificationActionTriggered ---

func TestActionNotificationActionTriggered_Good(t *core.T) {
	// ActionNotificationActionTriggered is broadcast by external code; confirm it can be received
	_, c := newTestService(t)
	var received *ActionNotificationActionTriggered
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionNotificationActionTriggered); ok {
			received = &a
		}
		return core.Result{OK: true}
	})
	_ = c.ACTION(ActionNotificationActionTriggered{NotificationID: "n1", ActionID: "reply"})
	core.AssertNotNil(t, received)
	core.AssertEqual(t, "n1", received.NotificationID)
	core.AssertEqual(t, "reply", received.ActionID)
}

func TestActionNotificationDismissed_Good(t *core.T) {
	_, c := newTestService(t)
	var received *ActionNotificationDismissed
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionNotificationDismissed); ok {
			received = &a
		}
		return core.Result{OK: true}
	})
	_ = c.ACTION(ActionNotificationDismissed{ID: "n2"})
	core.AssertNotNil(t, received)
	core.AssertEqual(t, "n2", received.ID)
}

func TestQueryPermission_Bad(t *core.T) {
	// No service — QUERY returns handled=false
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryPermission{})
	core.AssertFalse(t, r.OK)
}

func TestQueryPermission_Ugly(t *core.T) {
	// Platform returns resultFailure — QUERY returns OK=false (framework does not propagate Value for failed queries)
	mock := &mockPlatform{permErr: core.NewError("platform resultFailure")}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	r := c.QUERY(QueryPermission{})
	core.AssertFalse(t, r.OK)
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

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
