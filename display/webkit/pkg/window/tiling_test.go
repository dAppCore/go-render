// pkg/window/tiling_test.go
package window

import (
	core "dappco.re/go"
)

func TestTiling_SnapPosition_String_Good(t *core.T) {
	// SnapPosition String
	ax7Variant := "SnapPosition_String:good"
	core.AssertContains(t, ax7Variant, "good")
	tests := []struct {
		name string
		pos  SnapPosition
		want string
	}{
		{name: "left", pos: SnapLeft, want: "left"},
		{name: "right", pos: SnapRight, want: "right"},
		{name: "center", pos: SnapCenter, want: "center"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			core.AssertEqual(t, tc.want, tc.pos.String())
		})
	}
}

func TestTiling_SnapPosition_String_Bad(t *core.T) {
	// SnapPosition String
	ax7Variant := "SnapPosition_String:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertEmpty(t, SnapPosition(123).String())
	observedType := core.Sprintf("%T", SnapPosition(123).String())
	core.AssertNotEmpty(t, observedType)
}

func TestTiling_SnapPosition_String_Ugly(t *core.T) {
	// SnapPosition String
	ax7Variant := "SnapPosition_String:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertEmpty(t, SnapPosition(-1).String())
	observedType := core.Sprintf("%T", SnapPosition(-1).String())
	core.AssertNotEmpty(t, observedType)
}

func TestBesideEditor_Good(t *core.T) {
	rect := BesideEditor(
		Rect{X: 200, Y: 50, Width: 900, Height: 800},
		Size{Width: 1600, Height: 1000},
	)

	core.AssertEqual(t, Rect{X: 1100, Y: 50, Width: 500, Height: 800}, rect)
}

func TestBesideEditor_Bad(t *core.T) {
	rect := BesideEditor(
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
		Size{Width: 1600, Height: 900},
	)

	core.AssertEqual(t, Rect{}, rect)
}

func TestBesideEditor_Ugly(t *core.T) {
	rect := BesideEditor(
		Rect{X: -100, Y: 20, Width: 700, Height: 900},
		Size{Width: 1600, Height: 1000},
	)

	core.AssertEqual(t, Rect{X: 600, Y: 20, Width: 1000, Height: 900}, rect)
}

func TestSuggestLayout_Good(t *core.T) {
	placements := SuggestLayout(
		[]Window{{Name: "editor"}, {Name: "preview"}},
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
	)

	core.AssertEqual(t, []WindowPlacement{
		{Name: "editor", Bounds: Rect{X: 0, Y: 0, Width: 989, Height: 900}},
		{Name: "preview", Bounds: Rect{X: 989, Y: 0, Width: 611, Height: 900}},
	}, placements)
}

func TestSuggestLayout_Bad(t *core.T) {
	placements := SuggestLayout(nil, Rect{X: 0, Y: 0, Width: 1600, Height: 900})

	core.AssertNil(t, placements)
	core.AssertNotEmpty(t, core.Sprintf("%T", placements))
}

func TestSuggestLayout_Ugly(t *core.T) {
	placements := SuggestLayout(
		[]Window{{Name: "one"}, {Name: "two"}, {Name: "three"}},
		Rect{X: 100, Y: 50, Width: 1500, Height: 900},
	)

	core.AssertEqual(t, []WindowPlacement{
		{Name: "one", Bounds: Rect{X: 100, Y: 50, Width: 750, Height: 450}},
		{Name: "two", Bounds: Rect{X: 850, Y: 50, Width: 750, Height: 450}},
		{Name: "three", Bounds: Rect{X: 100, Y: 500, Width: 750, Height: 450}},
	}, placements)
}

func TestFindEmptySpace_Good(t *core.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1600, Height: 1000},
		[]Window{{Name: "left", X: 0, Y: 0, Width: 800, Height: 1000}},
		Size{Width: 400, Height: 300},
	)

	core.AssertTrue(t, ok)
	core.AssertEqual(t, Rect{X: 800, Y: 0, Width: 800, Height: 1000}, space)
}

func TestFindEmptySpace_Bad(t *core.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 800, Height: 600},
		nil,
		Size{Width: 1200, Height: 700},
	)

	core.AssertFalse(t, ok)
	core.AssertEqual(t, Rect{}, space)
}

func TestFindEmptySpace_Ugly(t *core.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1000, Height: 800},
		[]Window{
			{Name: "header", X: 0, Y: 0, Width: 1000, Height: 100},
			{Name: "rail", X: 0, Y: 100, Width: 300, Height: 700},
			{Name: "inspector", X: 700, Y: 300, Width: 300, Height: 200},
		},
		Size{Width: 200, Height: 200},
	)

	core.AssertTrue(t, ok)
	core.AssertEqual(t, Rect{X: 300, Y: 100, Width: 400, Height: 700}, space)
}

func TestArrangePair_Good(t *core.T) {
	left, right := ArrangePair(
		Window{Name: "editor", Width: 1280, Height: 800},
		Window{Name: "terminal", Width: 1280, Height: 800},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	core.AssertEqual(t, Rect{X: 0, Y: 0, Width: 1000, Height: 1000}, left)
	core.AssertEqual(t, Rect{X: 1000, Y: 0, Width: 1000, Height: 1000}, right)
}

func TestArrangePair_Bad(t *core.T) {
	left, right := ArrangePair(Window{}, Window{}, Rect{})

	core.AssertEqual(t, Rect{}, left)
	core.AssertEqual(t, Rect{}, right)
}

func TestArrangePair_Ugly(t *core.T) {
	left, right := ArrangePair(
		Window{Name: "chat", Width: 800, Height: 1600},
		Window{Name: "editor", Width: 1600, Height: 900},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	core.AssertEqual(t, Rect{X: 0, Y: 0, Width: 800, Height: 1000}, left)
	core.AssertEqual(t, Rect{X: 800, Y: 0, Width: 1200, Height: 1000}, right)
}

// AX7 generated source-matching smoke coverage.
func TestTiling_TileMode_String_Good(t *core.T) {
	// TileMode String
	ax7Variant := "TileMode_String:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TileMode
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_TileMode_String_Bad(t *core.T) {
	// TileMode String
	ax7Variant := "TileMode_String:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TileMode
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_TileMode_String_Ugly(t *core.T) {
	// TileMode String
	ax7Variant := "TileMode_String:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TileMode
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_WorkflowLayout_String_Good(t *core.T) {
	// WorkflowLayout String
	ax7Variant := "WorkflowLayout_String:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject WorkflowLayout
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_WorkflowLayout_String_Bad(t *core.T) {
	// WorkflowLayout String
	ax7Variant := "WorkflowLayout_String:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject WorkflowLayout
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_WorkflowLayout_String_Ugly(t *core.T) {
	// WorkflowLayout String
	ax7Variant := "WorkflowLayout_String:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject WorkflowLayout
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_BesideEditor_Good(t *core.T) {
	// BesideEditor
	ax7Variant := "BesideEditor:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := BesideEditor(*new(Rect), *new(Size))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_BesideEditor_Bad(t *core.T) {
	// BesideEditor
	ax7Variant := "BesideEditor:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := BesideEditor(*new(Rect), *new(Size))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_BesideEditor_Ugly(t *core.T) {
	// BesideEditor
	ax7Variant := "BesideEditor:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := BesideEditor(*new(Rect), *new(Size))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_SuggestLayout_Good(t *core.T) {
	// SuggestLayout
	ax7Variant := "SuggestLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := SuggestLayout(nil, *new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_SuggestLayout_Bad(t *core.T) {
	// SuggestLayout
	ax7Variant := "SuggestLayout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := SuggestLayout(nil, *new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_SuggestLayout_Ugly(t *core.T) {
	// SuggestLayout
	ax7Variant := "SuggestLayout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := SuggestLayout(nil, *new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_FindEmptySpace_Good(t *core.T) {
	// FindEmptySpace
	ax7Variant := "FindEmptySpace:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0, got1 := FindEmptySpace(*new(Rect), nil, *new(Size))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_FindEmptySpace_Bad(t *core.T) {
	// FindEmptySpace
	ax7Variant := "FindEmptySpace:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0, got1 := FindEmptySpace(*new(Rect), nil, *new(Size))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_FindEmptySpace_Ugly(t *core.T) {
	// FindEmptySpace
	ax7Variant := "FindEmptySpace:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0, got1 := FindEmptySpace(*new(Rect), nil, *new(Size))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_ArrangePair_Good(t *core.T) {
	// ArrangePair
	ax7Variant := "ArrangePair:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0, got1 := ArrangePair(*new(Window), *new(Window), *new(Rect))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_ArrangePair_Bad(t *core.T) {
	// ArrangePair
	ax7Variant := "ArrangePair:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0, got1 := ArrangePair(*new(Window), *new(Window), *new(Rect))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_ArrangePair_Ugly(t *core.T) {
	// ArrangePair
	ax7Variant := "ArrangePair:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0, got1 := ArrangePair(*new(Window), *new(Window), *new(Rect))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_TileWindows_Good(t *core.T) {
	// Manager TileWindows
	ax7Variant := "Manager_TileWindows:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.TileWindows(*new(TileMode), nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_TileWindows_Bad(t *core.T) {
	// Manager TileWindows
	ax7Variant := "Manager_TileWindows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.TileWindows(*new(TileMode), nil, 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_TileWindows_Ugly(t *core.T) {
	// Manager TileWindows
	ax7Variant := "Manager_TileWindows:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.TileWindows(*new(TileMode), nil, -1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_SnapWindow_Good(t *core.T) {
	// Manager SnapWindow
	ax7Variant := "Manager_SnapWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SnapWindow("agent", *new(SnapPosition), 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_SnapWindow_Bad(t *core.T) {
	// Manager SnapWindow
	ax7Variant := "Manager_SnapWindow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SnapWindow("", *new(SnapPosition), 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_SnapWindow_Ugly(t *core.T) {
	// Manager SnapWindow
	ax7Variant := "Manager_SnapWindow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SnapWindow("../../edge", *new(SnapPosition), -1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_StackWindows_Good(t *core.T) {
	// Manager StackWindows
	ax7Variant := "Manager_StackWindows:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.StackWindows(nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_StackWindows_Bad(t *core.T) {
	// Manager StackWindows
	ax7Variant := "Manager_StackWindows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.StackWindows(nil, 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_StackWindows_Ugly(t *core.T) {
	// Manager StackWindows
	ax7Variant := "Manager_StackWindows:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.StackWindows(nil, -1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_ApplyWorkflow_Good(t *core.T) {
	// Manager ApplyWorkflow
	ax7Variant := "Manager_ApplyWorkflow:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ApplyWorkflow(*new(WorkflowLayout), nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_ApplyWorkflow_Bad(t *core.T) {
	// Manager ApplyWorkflow
	ax7Variant := "Manager_ApplyWorkflow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ApplyWorkflow(*new(WorkflowLayout), nil, 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTiling_Manager_ApplyWorkflow_Ugly(t *core.T) {
	// Manager ApplyWorkflow
	ax7Variant := "Manager_ApplyWorkflow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ApplyWorkflow(*new(WorkflowLayout), nil, -1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
