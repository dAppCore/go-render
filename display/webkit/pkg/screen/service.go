// pkg/screen/service.go
package screen

import (
	"context"

	core "dappco.re/go"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the screen service to a Core instance.
// core.WithService(screen.Register(wailsScreen))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	if s == nil || s.platform == nil {
		switch q.(type) {
		case QueryAll:
			return core.Result{Value: []Screen{}, OK: true}
		case QueryWorkAreas:
			return core.Result{Value: []Rect{}, OK: true}
		default:
			return core.Result{Value: nil, OK: true}
		}
	}
	switch q := q.(type) {
	case QueryAll:
		return core.Result{Value: s.platform.GetAll(), OK: true}
	case QueryPrimary:
		return core.Result{Value: s.platform.GetPrimary(), OK: true}
	case QueryByID:
		return core.Result{Value: s.queryByID(q.ID), OK: true}
	case QueryAtPoint:
		return core.Result{Value: s.queryAtPoint(q.X, q.Y), OK: true}
	case QueryWorkAreas:
		return core.Result{Value: s.queryWorkAreas(), OK: true}
	case QueryCurrent:
		return core.Result{Value: s.platform.GetCurrent(), OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) queryByID(id string) *Screen {
	if s == nil || s.platform == nil {
		return nil
	}
	screens := s.platform.GetAll()
	for i := range screens {
		if screens[i].ID == id {
			return &screens[i]
		}
	}
	return nil
}

func (s *Service) queryAtPoint(x, y int) *Screen {
	if s == nil || s.platform == nil {
		return nil
	}
	screens := s.platform.GetAll()
	for i := range screens {
		b := screens[i].Bounds
		if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
			return &screens[i]
		}
	}
	return nil
}

func (s *Service) queryWorkAreas() []Rect {
	if s == nil || s.platform == nil {
		return nil
	}
	screens := s.platform.GetAll()
	areas := make([]Rect, len(screens))
	for i, scr := range screens {
		areas[i] = scr.WorkArea
	}
	return areas
}
