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

type updatePlatform interface {
	Update(options NotificationOptions) resultFailure
}

type categoryPlatform interface {
	RegisterCategory(category NotificationCategory) resultFailure
}

type responsePlatform interface {
	OnResponse(callback func(notificationID, actionID, userText string))
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
	ID               string               `json:"id"`
	Actions          []NotificationAction `json:"actions,omitempty"`
	HasReplyField    bool                 `json:"hasReplyField,omitempty"`
	ReplyPlaceholder string               `json:"replyPlaceholder,omitempty"`
	ReplyButtonTitle string               `json:"replyButtonTitle,omitempty"`
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
	Data       map[string]any       `json:"data,omitempty"`
	// Silent suppresses the delivery sound. SoundName selects a
	// platform sound when Silent is false.
	Silent    bool   `json:"silent,omitempty"`
	SoundName string `json:"soundName,omitempty"`
	// AttachmentPaths contains absolute media paths. Wails uses the
	// first image on Windows/Linux and supports multiple attachments
	// on macOS.
	AttachmentPaths []string `json:"attachmentPaths,omitempty"`
	ThreadID        string   `json:"threadId,omitempty"`
	// InterruptionLevel is passive, active, timeSensitive, or critical.
	InterruptionLevel string `json:"interruptionLevel,omitempty"`
	// ScheduleDelaySeconds and ScheduleAtUnix are mutually exclusive.
	// Windows/Linux schedules are in-process; macOS schedules natively.
	ScheduleDelaySeconds int   `json:"scheduleDelaySeconds,omitempty"`
	ScheduleAtUnix       int64 `json:"scheduleAtUnix,omitempty"`
	// Update replaces an in-flight notification with the same ID where
	// the platform supports replacement.
	Update bool `json:"update,omitempty"`
}

// PermissionStatus indicates whether notifications are authorised.
type PermissionStatus struct {
	Granted bool `json:"granted"`
}
