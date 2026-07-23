// pkg/clipboard/wails.go
package clipboard

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's ClipboardManager.
//
// Wails alpha.83 only exposes text — no image API. We satisfy the
// base Platform interface (Text / SetText) and skip the optional
// ImagePlatform extension; consumers calling clipboard.set_image /
// clipboard.get_image will see the gui service's "platform unavailable
// for images" path until upstream lands an image API.
//
//	wp := clipboard.NewWailsPlatform(app)
//	core.WithService(clipboard.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

// Text reads the current clipboard text. Wails returns ("", false) when
// the clipboard is empty or holds non-text content; we surface that
// shape directly so the gui Service's empty-clipboard handling still
// fires the "no content" branch.
func (wp *WailsPlatform) Text() (string, bool) {
	if wp == nil || wp.app == nil || wp.app.Clipboard == nil {
		return "", false
	}
	return wp.app.Clipboard.Text()
}

// SetText writes text to the clipboard. Wails's bool return reports
// whether the underlying platform call succeeded — we forward it
// unchanged so callers can branch on success.
func (wp *WailsPlatform) SetText(text string) bool {
	if wp == nil || wp.app == nil || wp.app.Clipboard == nil {
		return false
	}
	return wp.app.Clipboard.SetText(text)
}
