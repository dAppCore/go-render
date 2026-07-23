// pkg/environment/wails.go
package environment

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// WailsPlatform implements Platform via Wails v3's EnvironmentManager
// + the ThemeChanged application event.
//
// Wails surfaces:
//
//   - app.Env.Info()              → OS / Arch / Debug / OSInfo / PlatformInfo
//
//   - app.Env.IsDarkMode()
//
//   - app.Env.GetAccentColor()    → "rgb(r,g,b)" string
//
//   - app.Env.OpenFileManager(path, selectFile)
//
//   - app.Env.HasFocusFollowsMouse()  (Linux only — Wails returns false elsewhere)
//
//   - app.Event.OnApplicationEvent(events.Common.ThemeChanged, …)
//
//     wp := environment.NewWailsPlatform(app)
//     core.WithService(environment.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) IsDarkMode() bool {
	if wp == nil || wp.app == nil || wp.app.Env == nil {
		return false
	}
	return wp.app.Env.IsDarkMode()
}

// Info maps Wails's EnvironmentInfo (which carries pointer-typed OSInfo
// + map-typed PlatformInfo) onto our flatter cross-platform shape.
// PlatformInfo is reduced to {Name, Version} — the gui Platform contract
// is intentionally narrower so the consumer can serialise it as JSON
// without leaking platform-specific keys.
func (wp *WailsPlatform) Info() EnvironmentInfo {
	if wp == nil || wp.app == nil || wp.app.Env == nil {
		return EnvironmentInfo{}
	}
	wInfo := wp.app.Env.Info()
	out := EnvironmentInfo{
		OS:    wInfo.OS,
		Arch:  wInfo.Arch,
		Debug: wInfo.Debug,
	}
	if wInfo.OSInfo != nil {
		out.Platform = PlatformInfo{
			Name:    wInfo.OSInfo.Name,
			Version: wInfo.OSInfo.Version,
		}
	}
	return out
}

func (wp *WailsPlatform) AccentColour() string {
	if wp == nil || wp.app == nil || wp.app.Env == nil {
		return ""
	}
	return wp.app.Env.GetAccentColor()
}

func (wp *WailsPlatform) OpenFileManager(path string, selectFile bool) resultFailure {
	if wp == nil || wp.app == nil || wp.app.Env == nil {
		return nil
	}
	if err := wp.app.Env.OpenFileManager(path, selectFile); err != nil {
		return err
	}
	return nil
}

func (wp *WailsPlatform) HasFocusFollowsMouse() bool {
	if wp == nil || wp.app == nil || wp.app.Env == nil {
		return false
	}
	return wp.app.Env.HasFocusFollowsMouse()
}

// OnThemeChange registers a handler for the cross-platform ThemeChanged
// event. The event Context exposes IsDarkMode(); we forward that value
// to the handler. Returns the cancel func from Wails OnApplicationEvent.
func (wp *WailsPlatform) OnThemeChange(handler func(isDark bool)) func() {
	if wp == nil || wp.app == nil || wp.app.Event == nil || handler == nil {
		return func() {}
	}
	return wp.app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(evt *application.ApplicationEvent) {
		if evt == nil {
			return
		}
		ctx := evt.Context()
		if ctx == nil {
			handler(wp.IsDarkMode())
			return
		}
		handler(ctx.IsDarkMode())
	})
}
