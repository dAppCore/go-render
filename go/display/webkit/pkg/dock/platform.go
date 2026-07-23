// pkg/dock/platform.go
package dock

// BounceType controls how the dock icon attracts attention.
// bounce := dock.BounceInformational — single bounce
// bounce := dock.BounceCritical       — continuous bounce until focused
type BounceType int

const (
	// BounceInformational performs a single bounce to indicate a background event.
	BounceInformational BounceType = iota
	// BounceCritical bounces continuously until the application becomes active.
	BounceCritical
)

// Platform abstracts the dock/taskbar backend (Wails v3).
// macOS: dock icon show/hide, badge, progress bar, bounce.
// Windows: taskbar badge + progress bar (show/hide and bounce not supported).
// Linux: not supported — adapter returns nil for all operations.
type Platform interface {
	ShowIcon() resultFailure
	HideIcon() resultFailure
	SetBadge(label string) resultFailure
	RemoveBadge() resultFailure
	IsVisible() bool
	// SetProgressBar sets a progress indicator on the dock/taskbar icon.
	// progress is clamped to [0.0, 1.0]. Pass -1.0 to hide the indicator.
	SetProgressBar(progress float64) resultFailure
	// Bounce requests user attention by animating the dock icon.
	// Returns a request ID that can be passed to StopBounce.
	Bounce(bounceType BounceType) (int, resultFailure)
	// StopBounce cancels a pending attention request by its ID.
	StopBounce(requestID int) resultFailure
}
