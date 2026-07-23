// pkg/notification/wails.go
package notification

import (
	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
	wailsnotif "github.com/wailsapp/wails/v3/pkg/services/notifications"
)

// WailsPlatform implements Platform via Wails v3's notifications service.
//
// macOS requires a properly bundled + signed application — Wails refuses
// to send otherwise. Windows uses Toast notifications, Linux uses D-Bus.
//
//	wp := notification.NewWailsPlatform(app) // also calls app.RegisterService
//	core.WithService(notification.Register(wp))
type WailsPlatform struct {
	app     *application.App
	service *wailsnotif.NotificationService
}

// NewWailsPlatform creates the singleton Wails NotificationService and
// registers it with the App so its Startup hook runs (macOS bundle check,
// permission delegate init). nil app makes Send a no-op.
func NewWailsPlatform(app *application.App) *WailsPlatform {
	if app == nil {
		return &WailsPlatform{}
	}
	svc := wailsnotif.New()
	app.RegisterService(application.NewService(svc))
	return &WailsPlatform{app: app, service: svc}
}

func (wp *WailsPlatform) Send(opts NotificationOptions) resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	wOpts := toWailsNotificationOptions(opts)
	if len(opts.Actions) > 0 {
		if err := wp.service.SendNotificationWithActions(wOpts); err != nil {
			return core.E("notification.WailsPlatform.Send", "failed to send notification with actions", err)
		}
		return nil
	}
	if err := wp.service.SendNotification(wOpts); err != nil {
		return core.E("notification.WailsPlatform.Send", "failed to send notification", err)
	}
	return nil
}

func (wp *WailsPlatform) Update(opts NotificationOptions) resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	if err := wp.service.UpdateNotification(toWailsNotificationOptions(opts)); err != nil {
		return core.E("notification.WailsPlatform.Update", "failed to update notification", err)
	}
	return nil
}

func (wp *WailsPlatform) RegisterCategory(category NotificationCategory) resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	actions := make([]wailsnotif.NotificationAction, 0, len(category.Actions))
	for _, action := range category.Actions {
		actions = append(actions, wailsnotif.NotificationAction{
			ID:          action.ID,
			Title:       action.Title,
			Destructive: action.Destructive,
		})
	}
	wailsCategory := wailsnotif.NotificationCategory{
		ID:               category.ID,
		Actions:          actions,
		HasReplyField:    category.HasReplyField,
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}
	if err := wp.service.RegisterNotificationCategory(wailsCategory); err != nil {
		return core.E("notification.WailsPlatform.RegisterCategory", "failed to register notification category", err)
	}
	return nil
}

func (wp *WailsPlatform) OnResponse(callback func(notificationID, actionID, userText string)) {
	if wp == nil || wp.service == nil || callback == nil {
		return
	}
	wp.service.OnNotificationResponse(func(result wailsnotif.NotificationResult) {
		if result.Error != nil {
			return
		}
		response := result.Response
		callback(response.ID, response.ActionIdentifier, response.UserText)
	})
}

func toWailsNotificationOptions(opts NotificationOptions) wailsnotif.NotificationOptions {
	wOpts := wailsnotif.NotificationOptions{
		ID:                coalesceID(opts.ID),
		Title:             opts.Title,
		Subtitle:          opts.Subtitle,
		Body:              opts.Message,
		CategoryID:        opts.CategoryID,
		Data:              opts.Data,
		ThreadID:          opts.ThreadID,
		InterruptionLevel: opts.InterruptionLevel,
	}
	if opts.Silent || opts.SoundName != "" {
		wOpts.Sound = &wailsnotif.NotificationSound{
			Silent: opts.Silent,
			Name:   opts.SoundName,
		}
	}
	for _, path := range opts.AttachmentPaths {
		wOpts.Attachments = append(wOpts.Attachments, wailsnotif.NotificationAttachment{Path: path})
	}
	if opts.ScheduleDelaySeconds != 0 || opts.ScheduleAtUnix != 0 {
		wOpts.Schedule = &wailsnotif.NotificationSchedule{
			DelaySeconds: opts.ScheduleDelaySeconds,
			At:           opts.ScheduleAtUnix,
		}
	}
	return wOpts
}

func (wp *WailsPlatform) RequestPermission() (bool, resultFailure) {
	if wp == nil || wp.service == nil {
		return false, nil
	}
	granted, err := wp.service.RequestNotificationAuthorization()
	if err != nil {
		return false, core.E("notification.WailsPlatform.RequestPermission", "failed to request notification authorisation", err)
	}
	return granted, nil
}

func (wp *WailsPlatform) CheckPermission() (bool, resultFailure) {
	if wp == nil || wp.service == nil {
		return false, nil
	}
	granted, err := wp.service.CheckNotificationAuthorization()
	if err != nil {
		return false, core.E("notification.WailsPlatform.CheckPermission", "failed to check notification authorisation", err)
	}
	return granted, nil
}

// RevokePermission is a no-op — neither macOS nor Windows expose a
// programmatic revoke API; the user manages it via System Settings.
func (wp *WailsPlatform) RevokePermission() resultFailure {
	return nil
}

// Clear implements ClearPlatform. Empty id clears all delivered.
func (wp *WailsPlatform) Clear(id string) resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	if id == "" {
		if err := wp.service.RemoveAllDeliveredNotifications(); err != nil {
			return core.E("notification.WailsPlatform.Clear", "failed to clear delivered notifications", err)
		}
		return nil
	}
	if err := wp.service.RemoveDeliveredNotification(id); err != nil {
		return core.E("notification.WailsPlatform.Clear", "failed to clear delivered notification", err)
	}
	return nil
}

// coalesceID — Wails NotificationOptions.ID is required (validated by the
// service); fall back to a stable default if the consumer didn't supply one.
func coalesceID(id string) string {
	if id != "" {
		return id
	}
	return "core.notification"
}
