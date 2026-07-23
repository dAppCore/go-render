// pkg/notification/wails.go
package notification

import (
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
	wOpts := wailsnotif.NotificationOptions{
		ID:         coalesceID(opts.ID),
		Title:      opts.Title,
		Subtitle:   opts.Subtitle,
		Body:       opts.Message,
		CategoryID: opts.CategoryID,
	}
	if len(opts.Actions) > 0 {
		if err := wp.service.SendNotificationWithActions(wOpts); err != nil {
			return err
		}
		return nil
	}
	if err := wp.service.SendNotification(wOpts); err != nil {
		return err
	}
	return nil
}

func (wp *WailsPlatform) RequestPermission() (bool, resultFailure) {
	if wp == nil || wp.service == nil {
		return false, nil
	}
	granted, err := wp.service.RequestNotificationAuthorization()
	if err != nil {
		return false, err
	}
	return granted, nil
}

func (wp *WailsPlatform) CheckPermission() (bool, resultFailure) {
	if wp == nil || wp.service == nil {
		return false, nil
	}
	granted, err := wp.service.CheckNotificationAuthorization()
	if err != nil {
		return false, err
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
			return err
		}
		return nil
	}
	if err := wp.service.RemoveDeliveredNotification(id); err != nil {
		return err
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
