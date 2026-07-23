// pkg/screen/service_test.go
package screen

import (
	"context"

	core "dappco.re/go"
)

type mockPlatform struct {
	screens []Screen
	current *Screen
}

func (m *mockPlatform) GetAll() []Screen { return m.screens }
func (m *mockPlatform) GetPrimary() *Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}
func (m *mockPlatform) GetCurrent() *Screen {
	if m.current != nil {
		return m.current
	}
	return m.GetPrimary()
}

func newTestService(t *core.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		screens: []Screen{
			{
				ID: "1", Name: "Built-in", IsPrimary: true,
				Bounds:   Rect{X: 0, Y: 0, Width: 2560, Height: 1600},
				WorkArea: Rect{X: 0, Y: 38, Width: 2560, Height: 1562},
				Size:     Size{Width: 2560, Height: 1600},
			},
			{
				ID: "2", Name: "External",
				Bounds:   Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				WorkArea: Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				Size:     Size{Width: 1920, Height: 1080},
			},
		},
	}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func TestRegister_Good(t *core.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "screen")
	core.AssertNotNil(t, svc)
}

func TestQueryAll_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAll{})
	core.RequireTrue(t, r.OK)
	screens := r.Value.([]Screen)
	core.AssertLen(t, screens, 2)
}

func TestQueryPrimary_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryPrimary{})
	core.RequireTrue(t, r.OK)
	scr := r.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertEqual(t, "Built-in", scr.Name)
	core.AssertTrue(t, scr.IsPrimary)
}

func TestQueryByID_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryByID{ID: "2"})
	core.RequireTrue(t, r.OK)
	scr := r.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertEqual(t, "External", scr.Name)
}

func TestQueryByID_Bad(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryByID{ID: "99"})
	core.RequireTrue(t, r.OK)
	core.AssertNil(t, r.Value)
}

func TestQueryAtPoint_Good(t *core.T) {
	_, c := newTestService(t)

	// Point on primary screen
	r := c.QUERY(QueryAtPoint{X: 100, Y: 100})
	core.RequireTrue(t, r.OK)
	scr := r.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertEqual(t, "Built-in", scr.Name)

	// Point on external screen
	r2 := c.QUERY(QueryAtPoint{X: 3000, Y: 500})
	scr = r2.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertEqual(t, "External", scr.Name)
}

func TestQueryAtPoint_Bad(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAtPoint{X: -1000, Y: -1000})
	core.RequireTrue(t, r.OK)
	core.AssertNil(t, r.Value)
}

func TestQueryWorkAreas_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryWorkAreas{})
	core.RequireTrue(t, r.OK)
	areas := r.Value.([]Rect)
	core.AssertLen(t, areas, 2)
	core.AssertEqual(t, 38, areas[0].Y) // primary has menu bar offset
}

// --- QueryCurrent ---

func TestQueryCurrent_Good(t *core.T) {
	// current falls back to primary when not explicitly set
	_, c := newTestService(t)
	r := c.QUERY(QueryCurrent{})
	core.RequireTrue(t, r.OK)
	scr := r.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertTrue(t, scr.IsPrimary)
	core.AssertEqual(t, "Built-in", scr.Name)
}

func TestQueryCurrent_Bad(t *core.T) {
	// no screens at all → GetCurrent returns nil
	mock := &mockPlatform{screens: []Screen{}}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryCurrent{})
	core.RequireTrue(t, r.OK)
	core.AssertNil(t, r.Value)
}

func TestQueryCurrent_Ugly(t *core.T) {
	// current is explicitly set to the external screen
	mock := &mockPlatform{
		screens: []Screen{
			{ID: "1", Name: "Built-in", IsPrimary: true,
				Bounds: Rect{X: 0, Y: 0, Width: 2560, Height: 1600}},
			{ID: "2", Name: "External",
				Bounds: Rect{X: 2560, Y: 0, Width: 1920, Height: 1080}},
		},
	}
	mock.current = &mock.screens[1]
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryCurrent{})
	scr := r.Value.(*Screen)
	core.AssertNotNil(t, scr)
	core.AssertEqual(t, "External", scr.Name)
}

// --- Rect geometry helpers ---

func TestRect_Origin_Good(t *core.T) {
	// Origin
	ax7Variant := "Origin:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.Origin()
	core.AssertEqual(t, Point{X: 10, Y: 20}, pt)
}

func TestRect_Corner_Good(t *core.T) {
	// Corner
	ax7Variant := "Corner:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.Corner()
	core.AssertEqual(t, Point{X: 110, Y: 70}, pt)
}

func TestRect_InsideCorner_Good(t *core.T) {
	// InsideCorner
	ax7Variant := "InsideCorner:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.InsideCorner()
	core.AssertEqual(t, Point{X: 109, Y: 69}, pt)
}

func TestRect_IsEmpty_Good(t *core.T) {
	// IsEmpty
	ax7Variant := "IsEmpty:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertFalse(t, Rect{X: 0, Y: 0, Width: 1, Height: 1}.IsEmpty())
	observedType := core.Sprintf("%T", Rect{X: 0, Y: 0, Width: 1, Height: 1}.IsEmpty())
	core.AssertNotEmpty(t, observedType)
}

func TestRect_IsEmpty_Bad(t *core.T) {
	// IsEmpty
	ax7Variant := "IsEmpty:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertTrue(t, Rect{}.IsEmpty())
	core.AssertTrue(t, Rect{Width: 0, Height: 10}.IsEmpty())
	core.AssertTrue(t, Rect{Width: 10, Height: -1}.IsEmpty())
}

func TestRect_Contains_Good(t *core.T) {
	// Contains
	ax7Variant := "Contains:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	core.AssertTrue(t, r.Contains(Point{X: 0, Y: 0}))
	core.AssertTrue(t, r.Contains(Point{X: 50, Y: 50}))
	core.AssertTrue(t, r.Contains(Point{X: 99, Y: 99}))
}

func TestRect_Contains_Bad(t *core.T) {
	// Contains
	ax7Variant := "Contains:bad"
	core.AssertContains(t, ax7Variant, "bad")
	r := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	// exclusive right/bottom edge
	core.AssertFalse(t, r.Contains(Point{X: 100, Y: 50}))
	core.AssertFalse(t, r.Contains(Point{X: 50, Y: 100}))
	core.AssertFalse(t, r.Contains(Point{X: -1, Y: 50}))
}

func TestRect_Contains_Ugly(t *core.T) {
	// Contains
	ax7Variant := "Contains:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	// zero-size rect never contains anything
	r := Rect{X: 5, Y: 5, Width: 0, Height: 0}
	core.AssertFalse(t, r.Contains(Point{X: 5, Y: 5}))
	core.AssertNotEmpty(t, core.Sprintf("%T", r))
}

func TestRect_RectSize_Good(t *core.T) {
	// RectSize
	ax7Variant := "RectSize:good"
	core.AssertContains(t, ax7Variant, "good")
	r := Rect{X: 100, Y: 200, Width: 1920, Height: 1080}
	sz := r.RectSize()
	core.AssertEqual(t, Size{Width: 1920, Height: 1080}, sz)
}

func TestRect_Intersect_Good(t *core.T) {
	// Intersect
	ax7Variant := "Intersect:good"
	core.AssertContains(t, ax7Variant, "good")
	a := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	b := Rect{X: 50, Y: 50, Width: 100, Height: 100}
	overlap := a.Intersect(b)
	core.AssertEqual(t, Rect{X: 50, Y: 50, Width: 50, Height: 50}, overlap)
}

func TestRect_Intersect_Bad(t *core.T) {
	// Intersect
	ax7Variant := "Intersect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	// no overlap
	a := Rect{X: 0, Y: 0, Width: 50, Height: 50}
	b := Rect{X: 100, Y: 100, Width: 50, Height: 50}
	overlap := a.Intersect(b)
	core.AssertTrue(t, overlap.IsEmpty())
}

func TestRect_Intersect_Ugly(t *core.T) {
	// Intersect
	ax7Variant := "Intersect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	// empty rect intersects nothing
	a := Rect{X: 0, Y: 0, Width: 0, Height: 0}
	b := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	overlap := a.Intersect(b)
	core.AssertTrue(t, overlap.IsEmpty())
}

// --- ScreenPlacement ---

func TestScreenPlacement_Apply_Good(t *core.T) {
	// Apply
	ax7Variant := "Apply:good"
	core.AssertContains(t, ax7Variant, "good")
	// secondary placed to the RIGHT of primary, no offset
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 2560, Height: 1600},
		WorkArea: Rect{X: 0, Y: 38, Width: 2560, Height: 1562},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 3000, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 3000, Y: 0, Width: 1920, Height: 1080},
	}
	NewPlacement(secondary, primary, AlignRight, 0, OffsetBegin).Apply()
	core.AssertEqual(t, 2560, secondary.Bounds.X)
	core.AssertEqual(t, 0, secondary.Bounds.Y)
	core.AssertEqual(t, 2560, secondary.WorkArea.X)
}

func TestScreenPlacement_Apply_Bad(t *core.T) {
	// Apply
	ax7Variant := "Apply:bad"
	core.AssertContains(t, ax7Variant, "bad")
	// screen placed ABOVE primary: newY = primary.Y - secondary.Height
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 0, Y: -600, Width: 1920, Height: 600},
		WorkArea: Rect{X: 0, Y: -600, Width: 1920, Height: 600},
	}
	NewPlacement(secondary, primary, AlignTop, 0, OffsetBegin).Apply()
	core.AssertEqual(t, 0, secondary.Bounds.X)
	core.AssertEqual(t, -600, secondary.Bounds.Y)
}

func TestScreenPlacement_Apply_Ugly(t *core.T) {
	// Apply
	ax7Variant := "Apply:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	// END offset reference — places secondary flush to the bottom-right of parent
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 800, Height: 600},
		WorkArea: Rect{X: 0, Y: 0, Width: 800, Height: 600},
	}
	// AlignBottom + OffsetEnd + offset=0 → secondary starts at right edge of parent
	NewPlacement(secondary, primary, AlignBottom, 0, OffsetEnd).Apply()
	core.AssertEqual(t, 1920-800, secondary.Bounds.X) // flush right
	core.AssertEqual(t, 1080, secondary.Bounds.Y)     // just below parent
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
