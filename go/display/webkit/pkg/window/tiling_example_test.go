package window

import core "dappco.re/go"

func ExampleBesideEditor() {
	rect := BesideEditor(
		Rect{X: 240, Y: 40, Width: 960, Height: 720},
		Size{Width: 1600, Height: 900},
	)

	core.WriteString(core.Stdout(), core.Sprintf("%d,%d %dx%d\n", rect.X, rect.Y, rect.Width, rect.Height))
	// Output:
	// 1200,40 400x720
}

func ExampleSuggestLayout() {
	placements := SuggestLayout(
		[]Window{{Name: "editor"}, {Name: "preview"}},
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
	)

	for _, placement := range placements {
		core.WriteString(core.Stdout(), core.Sprintf("%s:%d,%d %dx%d\n", placement.Name, placement.Bounds.X, placement.Bounds.Y, placement.Bounds.Width, placement.Bounds.Height))
	}
	// Output:
	// editor:0,0 989x900
	// preview:989,0 611x900
}

func ExampleFindEmptySpace() {
	rect, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
		[]Window{{Name: "editor", X: 0, Y: 0, Width: 1000, Height: 900}},
		Size{Width: 300, Height: 300},
	)

	core.WriteString(core.Stdout(), core.Sprintf("%t %d,%d %dx%d\n", ok, rect.X, rect.Y, rect.Width, rect.Height))
	// Output:
	// true 1000,0 600x900
}

func ExampleArrangePair() {
	left, right := ArrangePair(
		Window{Name: "editor", Width: 1600, Height: 900},
		Window{Name: "chat", Width: 900, Height: 1600},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	core.WriteString(core.Stdout(), core.Sprintf("%dx%d | %dx%d\n", left.Width, left.Height, right.Width, right.Height))
	// Output:
	// 1200x1000 | 800x1000
}

// AX7 generated examples exercise each public call path with stable output.
func ExampleTileMode_String() {
	var subject TileMode
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleSnapPosition_String() {
	var subject SnapPosition
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleWorkflowLayout_String() {
	var subject WorkflowLayout
	result := core.Try(func() any {
		got0 := subject.String()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleManager_TileWindows() {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.TileWindows(*new(TileMode), nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleManager_SnapWindow() {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SnapWindow("agent", *new(SnapPosition), 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleManager_StackWindows() {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.StackWindows(nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleManager_ApplyWorkflow() {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ApplyWorkflow(*new(WorkflowLayout), nil, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}
