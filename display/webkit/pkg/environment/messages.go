// pkg/environment/messages.go
package environment

// QueryTheme returns the current theme. Result: ThemeInfo
type QueryTheme struct{}

// QueryInfo returns environment information. Result: EnvironmentInfo
type QueryInfo struct{}

// QueryAccentColour returns the system accent colour. Result: string
type QueryAccentColour struct{}

// TaskOpenFileManager opens the system file manager. Result: error only
type TaskOpenFileManager struct {
	Path   string `json:"path,omitempty"`
	Select bool   `json:"select"`
}

// TaskSetTheme overrides the application theme.
// Theme may be "dark", "light", or "system" to follow the platform again.
type TaskSetTheme struct {
	Theme string `json:"theme"`
}

// QueryFocusFollowsMouse returns whether the platform uses focus-follows-mouse. Result: bool
type QueryFocusFollowsMouse struct{}

// ActionThemeChanged is broadcast when the system theme changes.
type ActionThemeChanged struct {
	IsDark bool `json:"isDark"`
}
