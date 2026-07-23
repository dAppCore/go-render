package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/screen"
)

func newScreenToolsTestSubsystem(t *core.T, query func(core.Query) core.Result) *Subsystem {
	t.Helper()
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if query != nil {
			return query(q)
		}
		return core.Result{}
	})
	return New(c)
}

func TestToolsScreen_screenGet_Good(t *core.T) {
	// screenGet
	ax7Variant := "screenGet:good"
	core.AssertContains(t, ax7Variant, "good")
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{
				Value: &screen.Screen{
					ID:   "primary",
					Name: "Main",
					Bounds: screen.Rect{
						X: 0, Y: 0, Width: 1920, Height: 1080,
					},
					WorkArea: screen.Rect{
						X: 0, Y: 0, Width: 1920, Height: 1040,
					},
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	_, out, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "primary"})
	core.RequireNoError(t, err)
	core.AssertNotNil(t, out.Screen)
	core.AssertEqual(t, "primary", out.Screen.ID)
	core.AssertEqual(t, "Main", out.Screen.Name)
}

func TestToolsScreen_screenGet_Bad(t *core.T) {
	// screenGet
	ax7Variant := "screenGet:bad"
	core.AssertContains(t, ax7Variant, "bad")
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(screen.QueryByID); ok {
			return core.Result{OK: false, Value: "screen backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "missing"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "screen query failed")
}

func TestToolsScreen_screenGet_Ugly(t *core.T) {
	// screenGet
	ax7Variant := "screenGet:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(screen.QueryByID); ok {
			return core.Result{OK: true, Value: core.NewError("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "broken"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "unexpected result type")
}
