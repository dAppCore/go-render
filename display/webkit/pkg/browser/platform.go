// pkg/browser/platform.go
package browser

// Platform abstracts the system browser/file-opener backend.
type Platform interface {
	// OpenURL opens the given URL in the default system browser.
	OpenURL(url string) resultFailure

	// OpenFile opens the given file path with the system default application.
	OpenFile(path string) resultFailure
}
