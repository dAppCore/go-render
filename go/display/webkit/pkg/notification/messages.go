package notification

// QueryPermission returns current notification permission status. Result: PermissionStatus
type QueryPermission struct{}

// TaskSend sends a native notification, falling back to dialog on failure.
type TaskSend struct{ Options NotificationOptions }

// TaskRequestPermission requests notification permission from the OS. Result: bool (granted)
type TaskRequestPermission struct{}

// TaskRevokePermission revokes previously granted notification permission. Result: nil
type TaskRevokePermission struct{}

// TaskRegisterCategory registers a notification category with its actions.
// c.PERFORM(notification.TaskRegisterCategory{Category: notification.NotificationCategory{ID: "message", Actions: actions}})
type TaskRegisterCategory struct{ Category NotificationCategory }

// TaskClear dismisses a notification by id or all notifications when ID is empty.
type TaskClear struct {
	ID string `json:"id,omitempty"`
}

// ActionNotificationClicked is broadcast when the user clicks a notification body.
type ActionNotificationClicked struct{ ID string }

// ActionNotificationActionTriggered is broadcast when the user activates a notification action button.
//
//	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
//	  if a, ok := msg.(notification.ActionNotificationActionTriggered); ok { ... }
//	  return nil
//	})
type ActionNotificationActionTriggered struct {
	NotificationID string `json:"notificationId"`
	ActionID       string `json:"actionId"`
	UserText       string `json:"userText,omitempty"`
}

// ActionNotificationDismissed is broadcast when the user dismisses a notification.
type ActionNotificationDismissed struct{ ID string }
