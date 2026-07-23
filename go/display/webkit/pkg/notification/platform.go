// pkg/notification/platform.go
package notification

// Platform abstracts the native notification backend.
type Platform interface {
	Send(options NotificationOptions) resultFailure
	RequestPermission() (bool, resultFailure)
	CheckPermission() (bool, resultFailure)
	RevokePermission() resultFailure
}

// ClearPlatform is an optional extension for backends that can dismiss
// notifications after they have been shown. An empty id means "clear all".
type ClearPlatform interface {
	Clear(id string) resultFailure
}

// NotificationSeverity indicates the severity for dialog fallback.
type NotificationSeverity int

const (
	SeverityInfo NotificationSeverity = iota
	SeverityWarning
	SeverityError
)

// NotificationAction is a button that can be attached to a notification.
// id := "reply"; action := NotificationAction{ID: id, Title: "Reply", Destructive: false}
type NotificationAction struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Destructive bool   `json:"destructive,omitempty"`
}

// NotificationCategory groups actions under a named category/channel.
// category := NotificationCategory{ID: "message", Actions: []NotificationAction{{ID: "reply", Title: "Reply"}}}
type NotificationCategory struct {
	ID      string               `json:"id"`
	Actions []NotificationAction `json:"actions,omitempty"`
}

// NotificationOptions contains options for sending a notification.
type NotificationOptions struct {
	ID         string               `json:"id,omitempty"`
	Title      string               `json:"title"`
	Message    string               `json:"message"`
	Subtitle   string               `json:"subtitle,omitempty"`
	Severity   NotificationSeverity `json:"severity,omitempty"`
	CategoryID string               `json:"categoryId,omitempty"`
	Actions    []NotificationAction `json:"actions,omitempty"`
}

// PermissionStatus indicates whether notifications are authorised.
type PermissionStatus struct {
	Granted bool `json:"granted"`
}
