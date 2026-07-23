package window

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/screen"
)

type mockScreenPlatform struct {
	screens []screen.Screen
}

func (m *mockScreenPlatform) GetAll() []screen.Screen { return m.screens }

func (m *mockScreenPlatform) GetPrimary() *screen.Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}

func (m *mockScreenPlatform) GetCurrent() *screen.Screen {
	return m.GetPrimary()
}

func newTestWindowServiceWithScreen(t *core.T, screens []screen.Screen) (*Service, *core.Core) {
	t.Helper()

	c := core.New(
		core.WithService(screen.Register(&mockScreenPlatform{screens: screens})),
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func TestTaskTileWindows_Good_UsesPrimaryScreenSize(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		},
	})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}}).OK)
	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}}).OK)

	r := taskRun(c, "window.tile_windows", TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	core.RequireTrue(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "left"})
	core.RequireTrue(t, r2.OK)
	left := r2.Value.(*WindowInfo)
	core.AssertEqual(t, 0, left.X)
	core.AssertEqual(t, 1000, left.Width)
	core.AssertEqual(t, 1000, left.Height)

	r3 := c.QUERY(QueryWindowByName{Name: "right"})
	core.RequireTrue(t, r3.OK)
	right := r3.Value.(*WindowInfo)
	core.AssertEqual(t, 1000, right.X)
	core.AssertEqual(t, 1000, right.Width)
	core.AssertEqual(t, 1000, right.Height)
}

func TestTaskSnapWindow_Good_UsesPrimaryScreenSize(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		},
	})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("snap"), WithSize(400, 300)}}).OK)

	r := taskRun(c, "window.snap_window", TaskSnapWindow{Name: "snap", Position: "left"})
	core.RequireTrue(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "snap"})
	core.RequireTrue(t, r2.OK)
	info := r2.Value.(*WindowInfo)
	core.AssertEqual(t, 0, info.X)
	core.AssertEqual(t, 0, info.Y)
	core.AssertEqual(t, 1000, info.Width)
	core.AssertEqual(t, 1000, info.Height)
}

func TestTaskTileWindows_Good_UsesPrimaryWorkAreaOrigin(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 100, Y: 50, Width: 2000, Height: 1000},
		},
	})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}}).OK)
	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}}).OK)

	r := taskRun(c, "window.tile_windows", TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	core.RequireTrue(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "left"})
	core.RequireTrue(t, r2.OK)
	left := r2.Value.(*WindowInfo)
	core.AssertEqual(t, 100, left.X)
	core.AssertEqual(t, 50, left.Y)
	core.AssertEqual(t, 1000, left.Width)
	core.AssertEqual(t, 1000, left.Height)

	r3 := c.QUERY(QueryWindowByName{Name: "right"})
	core.RequireTrue(t, r3.OK)
	right := r3.Value.(*WindowInfo)
	core.AssertEqual(t, 1100, right.X)
	core.AssertEqual(t, 50, right.Y)
	core.AssertEqual(t, 1000, right.Width)
	core.AssertEqual(t, 1000, right.Height)
}
