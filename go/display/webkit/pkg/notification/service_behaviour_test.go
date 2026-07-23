// pkg/notification/service_behaviour_test.go
package notification

import (
	"context"

	core "dappco.re/go"
)

// notification.send accepts direct field options (no "task" wrapper),
// decoding them through decodeOptions / notificationOptionsFrom.
//
//	c.Action("notification.send").Run(ctx, core.NewOptions(
//	    core.Option{Key: "title", Value: "Direct"}))
func TestServiceBehaviour_SendDirectOptions_Good(t *core.T) {
	mock, c := newTestService(t)

	r := c.Action("notification.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "title", Value: "Direct"},
		core.Option{Key: "message", Value: "No task wrapper"},
	))
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.sendCalled)
	core.AssertEqual(t, "Direct", mock.lastOpts.Title)
	// An ID is synthesised when none was supplied.
	core.AssertNotEmpty(t, mock.lastOpts.ID)
}

// A registered category's actions are applied to a notification that names the
// category but supplies no actions of its own.
func TestServiceBehaviour_CategoryActions_Good(t *core.T) {
	mock, c := newTestService(t)

	reg := taskRun(c, "notification.registerCategory", TaskRegisterCategory{
		Category: NotificationCategory{
			ID:      "message",
			Actions: []NotificationAction{{ID: "reply", Title: "Reply"}},
		},
	})
	core.RequireTrue(t, reg.OK)

	send := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{Title: "Msg", Message: "Hi", CategoryID: "message"},
	})
	core.RequireTrue(t, send.OK)
	core.AssertLen(t, mock.lastOpts.Actions, 1)
	core.AssertEqual(t, "reply", mock.lastOpts.Actions[0].ID)
}

// notification.clear removes the active notification through the platform's
// ClearPlatform interface.
func TestServiceBehaviour_Clear_Good(t *core.T) {
	mock, c := newTestService(t)

	send := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n-1", Title: "T", Message: "M"},
	})
	core.RequireTrue(t, send.OK)

	clearResult := c.Action("notification.clear").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: "n-1"},
	))
	core.RequireTrue(t, clearResult.OK)
	core.AssertTrue(t, mock.clearCalled)
}

// decodeOptions returns the zero value (no error) when there are no options to
// decode.
func TestServiceBehaviour_decodeOptions_Empty(t *core.T) {
	got, err := decodeOptions[NotificationOptions](core.NewOptions())
	core.AssertNil(t, err)
	core.AssertEqual(t, "", got.Title)
}

// notificationOptionsFrom resolves both the TaskSend and bare
// NotificationOptions task shapes.
func TestServiceBehaviour_notificationOptionsFrom_Good(t *core.T) {
	fromTask, err := notificationOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskSend{Options: NotificationOptions{Title: "A"}}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, "A", fromTask.Title)

	fromOpts, err := notificationOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: NotificationOptions{Title: "B"}},
	))
	core.AssertNil(t, err)
	core.AssertEqual(t, "B", fromOpts.Title)
}

func TestServiceBehaviour_toWailsNotificationOptions_Good(t *core.T) {
	got := toWailsNotificationOptions(NotificationOptions{
		ID:                   "n-7",
		Title:                "Ready",
		Message:              "Done",
		Silent:               true,
		SoundName:            "Ping",
		AttachmentPaths:      []string{"/tmp/result.png"},
		ThreadID:             "jobs",
		InterruptionLevel:    "active",
		ScheduleDelaySeconds: 15,
	})

	core.AssertEqual(t, "n-7", got.ID)
	core.RequireTrue(t, got.Sound != nil)
	core.AssertTrue(t, got.Sound.Silent)
	core.AssertEqual(t, "Ping", got.Sound.Name)
	core.AssertLen(t, got.Attachments, 1)
	core.AssertEqual(t, "/tmp/result.png", got.Attachments[0].Path)
	core.RequireTrue(t, got.Schedule != nil)
	core.AssertEqual(t, 15, got.Schedule.DelaySeconds)
}
