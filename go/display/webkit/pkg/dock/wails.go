// pkg/dock/wails.go
package dock

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	wailsdock "github.com/wailsapp/wails/v3/pkg/services/dock"
)

// WailsPlatform implements Platform via Wails v3's dock service.
//
// Wails alpha.83 covers visibility (HideAppIcon/ShowAppIcon) + badge
// (SetBadge/RemoveBadge/GetBadge) on macOS / Windows. Linux's wails
// dock impl is a no-op stub. Wails does NOT expose:
//
//   - SetProgressBar — neither macOS NSDockTile progress nor Windows
//     ITaskbarList3 progress is wrapped. We accept the call so callers
//     don't crash and silently no-op until upstream lands it.
//
//   - Bounce / StopBounce — macOS NSApp.requestUserAttention(:) and the
//     Windows FlashWindowEx flags are not exposed. Same no-op shape;
//     Bounce returns request ID 0 so StopBounce(0) is a valid no-op.
//
//   - IsVisible — Wails has no getter; we track our own state across
//     ShowIcon/HideIcon and seed visible=true (the default state when
//     the app launches in a dock-visible config).
//
//     wp := dock.NewWailsPlatform(app) // also calls app.RegisterService
//     core.WithService(dock.Register(wp))
type WailsPlatform struct {
	app     *application.App
	service *wailsdock.DockService

	mu      sync.RWMutex
	visible bool
}

// NewWailsPlatform creates the Wails DockService and registers it with
// the App so its Startup hook (NSDockTile setup on macOS, taskbar COM
// init on Windows) runs. nil app makes everything a no-op.
func NewWailsPlatform(app *application.App) *WailsPlatform {
	if app == nil {
		return &WailsPlatform{visible: true}
	}
	svc := wailsdock.New()
	app.RegisterService(application.NewService(svc))
	return &WailsPlatform{app: app, service: svc, visible: true}
}

func (wp *WailsPlatform) ShowIcon() resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	wp.service.ShowAppIcon()
	wp.mu.Lock()
	wp.visible = true
	wp.mu.Unlock()
	return nil
}

func (wp *WailsPlatform) HideIcon() resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	wp.service.HideAppIcon()
	wp.mu.Lock()
	wp.visible = false
	wp.mu.Unlock()
	return nil
}

func (wp *WailsPlatform) SetBadge(label string) resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	if err := wp.service.SetBadge(label); err != nil {
		return err
	}
	return nil
}

func (wp *WailsPlatform) RemoveBadge() resultFailure {
	if wp == nil || wp.service == nil {
		return nil
	}
	if err := wp.service.RemoveBadge(); err != nil {
		return err
	}
	return nil
}

// IsVisible reports our locally-tracked visibility state. Wails has no
// getter; we toggle on each ShowIcon/HideIcon call and seed visible=true.
func (wp *WailsPlatform) IsVisible() bool {
	if wp == nil {
		return false
	}
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.visible
}

// SetProgressBar — Wails alpha.83 does not expose dock/taskbar progress;
// no-op. progress is documented as [0.0, 1.0] with -1 to hide; we accept
// the call shape so future upstream support drops in without a Platform
// contract change.
func (wp *WailsPlatform) SetProgressBar(_ float64) resultFailure {
	return nil
}

// Bounce — Wails alpha.83 does not expose macOS requestUserAttention
// or Windows FlashWindowEx; no-op. Returns request ID 0 so a paired
// StopBounce(0) is also a valid no-op.
func (wp *WailsPlatform) Bounce(_ BounceType) (int, resultFailure) {
	return 0, nil
}

// StopBounce — paired no-op for Bounce.
func (wp *WailsPlatform) StopBounce(_ int) resultFailure {
	return nil
}
