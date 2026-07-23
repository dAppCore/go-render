// pkg/browser/wails.go
package browser

import "github.com/wailsapp/wails/v3/pkg/application"

// WailsPlatform implements Platform via Wails v3's app.Browser manager,
// which delegates to github.com/pkg/browser under the hood.
//
//	wp := browser.NewWailsPlatform(app)
//	core.WithService(browser.Register(wp))
type WailsPlatform struct {
	app *application.App
}

// NewWailsPlatform constructs the adapter. nil app makes Send a no-op so
// tests can wire the consumer side without spinning up Wails.
func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) OpenURL(url string) resultFailure {
	if wp == nil || wp.app == nil {
		return nil
	}
	if err := wp.app.Browser.OpenURL(url); err != nil {
		return err
	}
	return nil
}

func (wp *WailsPlatform) OpenFile(path string) resultFailure {
	if wp == nil || wp.app == nil {
		return nil
	}
	if err := wp.app.Browser.OpenFile(path); err != nil {
		return err
	}
	return nil
}
