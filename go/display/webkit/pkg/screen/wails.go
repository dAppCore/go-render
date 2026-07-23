// pkg/screen/wails.go
package screen

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's ScreenManager.
//
// Wails surfaces:
//   - app.Screen.GetAll()      → []*application.Screen
//   - app.Screen.GetPrimary()  → *application.Screen
//   - app.Screen.GetByID(id)   → *application.Screen
//   - app.Screen.GetByIndex(i) → *application.Screen
//
// Wails has no first-class "current" screen concept (no focused-screen
// tracker). Our gui Platform docs allow falling back to primary when
// current is unset; the closest faithful read would be "the screen
// containing the current window" — see GetCurrent below for the
// best-effort implementation.
//
//	wp := screen.NewWailsPlatform(app)
//	core.WithService(screen.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) GetAll() []Screen {
	if wp == nil || wp.app == nil || wp.app.Screen == nil {
		return nil
	}
	wScreens := wp.app.Screen.GetAll()
	out := make([]Screen, 0, len(wScreens))
	for _, ws := range wScreens {
		if ws == nil {
			continue
		}
		out = append(out, fromWailsScreen(ws))
	}
	return out
}

func (wp *WailsPlatform) GetPrimary() *Screen {
	if wp == nil || wp.app == nil || wp.app.Screen == nil {
		return nil
	}
	ws := wp.app.Screen.GetPrimary()
	if ws == nil {
		return nil
	}
	s := fromWailsScreen(ws)
	return &s
}

// GetCurrent — best-effort: Wails has no focused-screen API. We use
// the Window manager's Current() window's GetScreen() if available,
// otherwise fall back to GetPrimary (gui Platform docs sanction this
// fallback explicitly).
func (wp *WailsPlatform) GetCurrent() *Screen {
	if wp == nil || wp.app == nil {
		return nil
	}
	if wp.app.Window != nil {
		if w := wp.app.Window.Current(); w != nil {
			if ws, err := w.GetScreen(); err == nil && ws != nil {
				s := fromWailsScreen(ws)
				return &s
			}
		}
	}
	return wp.GetPrimary()
}

// fromWailsScreen maps Wails's *application.Screen onto our gui Screen
// shape. Wails uses float32 for ScaleFactor / Rotation; we widen to
// float64 to match the gui Platform contract (consistent with the
// JSON output the consumer round-trips through the bridge).
func fromWailsScreen(ws *application.Screen) Screen {
	return Screen{
		ID:               ws.ID,
		Name:             ws.Name,
		ScaleFactor:      float64(ws.ScaleFactor),
		Size:             Size{Width: ws.Size.Width, Height: ws.Size.Height},
		Bounds:           Rect{X: ws.Bounds.X, Y: ws.Bounds.Y, Width: ws.Bounds.Width, Height: ws.Bounds.Height},
		PhysicalBounds:   Rect{X: ws.PhysicalBounds.X, Y: ws.PhysicalBounds.Y, Width: ws.PhysicalBounds.Width, Height: ws.PhysicalBounds.Height},
		WorkArea:         Rect{X: ws.WorkArea.X, Y: ws.WorkArea.Y, Width: ws.WorkArea.Width, Height: ws.WorkArea.Height},
		PhysicalWorkArea: Rect{X: ws.PhysicalWorkArea.X, Y: ws.PhysicalWorkArea.Y, Width: ws.PhysicalWorkArea.Width, Height: ws.PhysicalWorkArea.Height},
		IsPrimary:        ws.IsPrimary,
		Rotation:         float64(ws.Rotation),
	}
}
