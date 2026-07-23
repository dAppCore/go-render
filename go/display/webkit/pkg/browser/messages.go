// pkg/browser/messages.go
package browser

// --- Tasks (all side-effects, no queries or actions) ---

// TaskOpenURL opens a URL in the default system browser. Result: nil
type TaskOpenURL struct {
	URL string `json:"url"`
}

// TaskOpenFile opens a file with the system default application. Result: nil
type TaskOpenFile struct {
	Path string `json:"path,omitempty"`
}
