// pkg/display/interfaces.go
package display

import (
	"net/url"

	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// RouteSchemeHandler dispatches a parsed route URL through Core.
//
//	result := handler.Handle(parsedURL)
type RouteSchemeHandler interface {
	Handle(url *url.URL) core.Result
}

// SchemeHandlerProvider exposes the active route scheme handler.
//
//	handler := svc.SchemeHandler()
type SchemeHandlerProvider interface {
	SchemeHandler() RouteSchemeHandler
}

// App abstracts the Wails application for the orchestrator.
// After Spec D cleanup, only Quit() and Logger() remain —
// all other Wails Manager APIs are accessed via IPC.
type App interface {
	Logger() Logger
	Quit()
}

// Logger wraps Wails logging.
type Logger interface {
	Info(message string, args ...any)
}

// wailsApp wraps *application.App for the App interface.
type wailsApp struct {
	app *application.App
}

func newWailsApp(app *application.App) App {
	return &wailsApp{app: app}
}

func (w *wailsApp) Logger() Logger { return w.app.Logger }
func (w *wailsApp) Quit()          { w.app.Quit() }
