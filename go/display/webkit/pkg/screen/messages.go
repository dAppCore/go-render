// pkg/screen/messages.go
package screen

// QueryAll returns all screens. Result: []Screen
type QueryAll struct{}

// QueryPrimary returns the primary screen. Result: *Screen (nil if not found)
type QueryPrimary struct{}

// QueryByID returns a screen by ID. Result: *Screen (nil if not found)
type QueryByID struct{ ID string }

// QueryAtPoint returns the screen containing a point. Result: *Screen (nil if none)
type QueryAtPoint struct{ X, Y int }

// QueryWorkAreas returns work areas for all screens. Result: []Rect
type QueryWorkAreas struct{}

// QueryCurrent returns the most recently active screen. Result: *Screen (nil if no screens registered)
//
//	result, _, _ := c.QUERY(screen.QueryCurrent{})
//	current := result.(*screen.Screen)
type QueryCurrent struct{}

// ActionScreensChanged is broadcast when displays change (future).
type ActionScreensChanged struct{ Screens []Screen }
