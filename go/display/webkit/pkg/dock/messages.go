// pkg/dock/messages.go
package dock

// --- Queries (read-only) ---

// QueryVisible returns whether the dock icon is visible. Result: bool
type QueryVisible struct{}

// --- Tasks (side-effects) ---

// TaskShowIcon shows the dock/taskbar icon. Result: nil
type TaskShowIcon struct{}

// TaskHideIcon hides the dock/taskbar icon. Result: nil
type TaskHideIcon struct{}

// TaskSetBadge sets the dock/taskbar badge label.
// Empty string "" shows the default system badge indicator.
// Numeric "3", "99" shows unread count. Text "New", "Paused" shows brief status.
// Result: nil
type TaskSetBadge struct{ Label string }

// TaskRemoveBadge removes the dock/taskbar badge. Result: nil
type TaskRemoveBadge struct{}

// TaskSetProgressBar updates the progress indicator on the dock/taskbar icon.
// Progress is clamped to [0.0, 1.0]. Pass -1.0 to hide the indicator.
// c.PERFORM(dock.TaskSetProgressBar{Progress: 0.75})  // 75% complete
// c.PERFORM(dock.TaskSetProgressBar{Progress: -1.0})  // hide indicator
// Result: nil
type TaskSetProgressBar struct{ Progress float64 }

// TaskBounce requests user attention by animating the dock icon.
// Result: int (requestID for use with TaskStopBounce)
// c.PERFORM(dock.TaskBounce{BounceType: dock.BounceInformational})
type TaskBounce struct{ BounceType BounceType }

// TaskStopBounce cancels a pending attention request.
// c.PERFORM(dock.TaskStopBounce{RequestID: id})
// Result: nil
type TaskStopBounce struct{ RequestID int }

// --- Actions (broadcasts) ---

// ActionVisibilityChanged is broadcast after a successful TaskShowIcon or TaskHideIcon.
type ActionVisibilityChanged struct{ Visible bool }

// ActionProgressChanged is broadcast after a successful TaskSetProgressBar.
type ActionProgressChanged struct{ Progress float64 }

// ActionBounceStarted is broadcast after a successful TaskBounce.
type ActionBounceStarted struct {
	RequestID  int        `json:"requestId"`
	BounceType BounceType `json:"bounceType"`
}
