package window

import (
	core "dappco.re/go"

	"dappco.re/go/render/display/webkit/pkg/screen"
)

func TestTaskLayoutBesideEditor_Good(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("cursor"), WithTitle("Cursor"), WithPosition(0, 0), WithSize(1400, 1000),
	}}).OK)
	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("assistant"), WithPosition(10, 10), WithSize(400, 500),
	}}).OK)

	r := taskRun(c, "window.layout_beside_editor", TaskLayoutBesideEditor{Name: "assistant"})
	core.RequireTrue(t, r.OK)
	result := r.Value.(LayoutBesideEditorResult)
	core.AssertEqual(t, "cursor", result.Editor)
	core.AssertEqual(t, "right", result.Side)
	core.AssertEqual(t, 1500, result.WindowBounds.X)
	core.AssertEqual(t, 500, result.WindowBounds.Width)
}

func TestTaskLayoutSuggest_Good(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	r := taskRun(c, "window.layout_suggest", TaskLayoutSuggest{WindowCount: 2})
	core.RequireTrue(t, r.OK)
	suggestion := r.Value.(LayoutSuggestion)
	core.AssertEqual(t, "left-right", suggestion.Mode)
}

func TestTaskScreenFindSpace_Good(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("left"), WithPosition(0, 0), WithSize(1000, 1000),
	}}).OK)

	r := taskRun(c, "window.find_space", TaskScreenFindSpace{Width: 400, Height: 400})
	core.RequireTrue(t, r.OK)
	space := r.Value.(ScreenSpace)
	core.AssertGreaterOrEqual(t, space.X, 1000)
	core.AssertGreaterOrEqual(t, space.Width, 400)
}

func TestTaskWindowArrangePair_Good(t *core.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("primary")}}).OK)
	core.RequireTrue(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("secondary")}}).OK)

	r := taskRun(c, "window.arrange_pair", TaskWindowArrangePair{Primary: "primary", Secondary: "secondary", Ratio: 0.6})
	core.RequireTrue(t, r.OK)
	arrangement := r.Value.(PairArrangement)
	core.AssertEqual(t, "horizontal", arrangement.Orientation)
	core.AssertEqual(t, 1200, arrangement.Primary.Width)
	core.AssertEqual(t, 1200, arrangement.Secondary.X)
}
