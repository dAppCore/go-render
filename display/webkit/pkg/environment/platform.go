// pkg/environment/platform.go
package environment

// Platform abstracts environment and theme backend queries.
type Platform interface {
	IsDarkMode() bool
	Info() EnvironmentInfo
	AccentColour() string
	OpenFileManager(path string, selectFile bool) resultFailure
	HasFocusFollowsMouse() bool
	OnThemeChange(handler func(isDark bool)) func() // returns cancel func
}

// EnvironmentInfo contains system environment details.
type EnvironmentInfo struct {
	OS       string       `json:"os,omitempty"`
	Arch     string       `json:"arch"`
	Debug    bool         `json:"debug"`
	Platform PlatformInfo `json:"platform"`
}

// PlatformInfo contains platform-specific details.
type PlatformInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ThemeInfo contains the current theme state.
type ThemeInfo struct {
	IsDark bool   `json:"isDark"`
	Theme  string `json:"theme"` // "dark" or "light"
}
