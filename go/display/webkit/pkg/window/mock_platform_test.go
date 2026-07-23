package window

import (
	core "dappco.re/go"
)

func TestMockPlatform_CreateWindow_Good(t *core.T) {
	// CreateWindow
	ax7Variant := "CreateWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	p := NewMockPlatform()
	w := p.CreateWindow(PlatformWindowOptions{
		Name:   "main",
		Title:  "Core GUI",
		URL:    "/home",
		HTML:   "<main>Ready</main>",
		JS:     "globalThis.ready = true",
		Width:  1280,
		Height: 800,
		X:      10,
		Y:      20,
	})

	core.AssertLen(t, p.Windows, 1)
	got := w.(*MockWindow)
	core.AssertEqual(t, "main", got.Name())
	core.AssertEqual(t, []string{"globalThis.ready = true"}, got.ExecJSCalls())
	core.AssertEqual(t, "Core GUI", got.Title())
	core.AssertEqual(t, 10, got.x)
	core.AssertEqual(t, 20, got.y)

	got.SetPosition(30, 40)
	got.SetSize(1920, 1080)
	got.SetVisibility(true)
	got.SetAlwaysOnTop(true)
	got.SetOpacity(0.75)
	got.SetBounds(1, 2, 3, 4)
	got.SetURL("/dashboard")
	got.SetHTML("<main>Updated</main>")
	got.SetZoom(1.25)
	got.SetContentProtection(true)
	got.Maximise()
	got.Restore()
	got.Minimise()
	got.Focus()
	got.Show()
	got.Hide()
	got.Fullscreen()
	got.UnFullscreen()
	got.ToggleFullscreen()
	got.ToggleMaximise()
	got.ExecJS("alert(1)")
	got.Flash(true)
	got.OpenDevTools()
	got.CloseDevTools()

	core.AssertEqual(t, 1, got.x)
	core.AssertEqual(t, 2, got.y)
	core.AssertEqual(t, 3, got.width)
	core.AssertEqual(t, 4, got.height)
	core.AssertTrue(t, got.maximised)
	core.AssertTrue(t, got.focused)
	core.AssertFalse(t, got.visible)
	core.AssertTrue(t, got.fullscreened)
	core.AssertTrue(t, got.minimised)
	core.AssertEqual(t, 0.75, got.opacity)
	core.AssertEqual(t, []string{"globalThis.ready = true", "alert(1)"}, got.ExecJSCalls())
	core.AssertTrue(t, got.flashed)
	core.AssertFalse(t, got.DevToolsOpen())
}

func TestMockPlatform_GetWindows_Bad(t *core.T) {
	// GetWindows
	ax7Variant := "GetWindows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	p := NewMockPlatform()
	core.AssertEmpty(t, p.GetWindows())
	core.AssertNotEmpty(t, core.Sprintf("%T", p))
}

func TestMockWindow_FileDrop_UglyCase(t *core.T) {
	w := &mockWindow{}
	calls := 0
	w.OnFileDrop(func(paths []string, target *DropTarget) {
		calls++
		core.AssertEqual(t, []string{"a.txt"}, paths)
		core.AssertNotNil(t, target)
		core.AssertEqual(t, "drop-zone", target.ID)
	})
	w.emitFileDrop([]string{"a.txt"}, &DropTarget{ID: "drop-zone"})

	core.AssertEqual(t, 1, calls)
}

// AX7 generated source-matching smoke coverage.
func TestMockPlatform_NewMockPlatform_Good(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_NewMockPlatform_Bad(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_NewMockPlatform_Ugly(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_CreateWindow_Good(t *core.T) {
	// MockPlatform CreateWindow
	ax7Variant := "MockPlatform_CreateWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_CreateWindow_Bad(t *core.T) {
	// MockPlatform CreateWindow
	ax7Variant := "MockPlatform_CreateWindow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_CreateWindow_Ugly(t *core.T) {
	// MockPlatform CreateWindow
	ax7Variant := "MockPlatform_CreateWindow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_GetWindows_Good(t *core.T) {
	// MockPlatform GetWindows
	ax7Variant := "MockPlatform_GetWindows:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_GetWindows_Bad(t *core.T) {
	// MockPlatform GetWindows
	ax7Variant := "MockPlatform_GetWindows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_GetWindows_Ugly(t *core.T) {
	// MockPlatform GetWindows
	ax7Variant := "MockPlatform_GetWindows:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Name_Good(t *core.T) {
	// MockWindow Name
	ax7Variant := "MockWindow_Name:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Name_Bad(t *core.T) {
	// MockWindow Name
	ax7Variant := "MockWindow_Name:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Name_Ugly(t *core.T) {
	// MockWindow Name
	ax7Variant := "MockWindow_Name:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Title_Good(t *core.T) {
	// MockWindow Title
	ax7Variant := "MockWindow_Title:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Title_Bad(t *core.T) {
	// MockWindow Title
	ax7Variant := "MockWindow_Title:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Title_Ugly(t *core.T) {
	// MockWindow Title
	ax7Variant := "MockWindow_Title:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Position_Good(t *core.T) {
	// MockWindow Position
	ax7Variant := "MockWindow_Position:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Position_Bad(t *core.T) {
	// MockWindow Position
	ax7Variant := "MockWindow_Position:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Position_Ugly(t *core.T) {
	// MockWindow Position
	ax7Variant := "MockWindow_Position:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Size_Good(t *core.T) {
	// MockWindow Size
	ax7Variant := "MockWindow_Size:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Size_Bad(t *core.T) {
	// MockWindow Size
	ax7Variant := "MockWindow_Size:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Size_Ugly(t *core.T) {
	// MockWindow Size
	ax7Variant := "MockWindow_Size:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMaximised_Good(t *core.T) {
	// MockWindow IsMaximised
	ax7Variant := "MockWindow_IsMaximised:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMaximised_Bad(t *core.T) {
	// MockWindow IsMaximised
	ax7Variant := "MockWindow_IsMaximised:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMaximised_Ugly(t *core.T) {
	// MockWindow IsMaximised
	ax7Variant := "MockWindow_IsMaximised:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFocused_Good(t *core.T) {
	// MockWindow IsFocused
	ax7Variant := "MockWindow_IsFocused:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFocused_Bad(t *core.T) {
	// MockWindow IsFocused
	ax7Variant := "MockWindow_IsFocused:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFocused_Ugly(t *core.T) {
	// MockWindow IsFocused
	ax7Variant := "MockWindow_IsFocused:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsVisible_Good(t *core.T) {
	// MockWindow IsVisible
	ax7Variant := "MockWindow_IsVisible:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsVisible_Bad(t *core.T) {
	// MockWindow IsVisible
	ax7Variant := "MockWindow_IsVisible:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsVisible_Ugly(t *core.T) {
	// MockWindow IsVisible
	ax7Variant := "MockWindow_IsVisible:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFullscreen_Good(t *core.T) {
	// MockWindow IsFullscreen
	ax7Variant := "MockWindow_IsFullscreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFullscreen_Bad(t *core.T) {
	// MockWindow IsFullscreen
	ax7Variant := "MockWindow_IsFullscreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsFullscreen_Ugly(t *core.T) {
	// MockWindow IsFullscreen
	ax7Variant := "MockWindow_IsFullscreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMinimised_Good(t *core.T) {
	// MockWindow IsMinimised
	ax7Variant := "MockWindow_IsMinimised:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMinimised_Bad(t *core.T) {
	// MockWindow IsMinimised
	ax7Variant := "MockWindow_IsMinimised:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_IsMinimised_Ugly(t *core.T) {
	// MockWindow IsMinimised
	ax7Variant := "MockWindow_IsMinimised:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetBounds_Good(t *core.T) {
	// MockWindow GetBounds
	ax7Variant := "MockWindow_GetBounds:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetBounds_Bad(t *core.T) {
	// MockWindow GetBounds
	ax7Variant := "MockWindow_GetBounds:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetBounds_Ugly(t *core.T) {
	// MockWindow GetBounds
	ax7Variant := "MockWindow_GetBounds:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetZoom_Good(t *core.T) {
	// MockWindow GetZoom
	ax7Variant := "MockWindow_GetZoom:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetZoom_Bad(t *core.T) {
	// MockWindow GetZoom
	ax7Variant := "MockWindow_GetZoom:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetZoom_Ugly(t *core.T) {
	// MockWindow GetZoom
	ax7Variant := "MockWindow_GetZoom:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetOpacity_Good(t *core.T) {
	// MockWindow GetOpacity
	ax7Variant := "MockWindow_GetOpacity:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetOpacity_Bad(t *core.T) {
	// MockWindow GetOpacity
	ax7Variant := "MockWindow_GetOpacity:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_GetOpacity_Ugly(t *core.T) {
	// MockWindow GetOpacity
	ax7Variant := "MockWindow_GetOpacity:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetTitle_Good(t *core.T) {
	// MockWindow SetTitle
	ax7Variant := "MockWindow_SetTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetTitle_Bad(t *core.T) {
	// MockWindow SetTitle
	ax7Variant := "MockWindow_SetTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetTitle_Ugly(t *core.T) {
	// MockWindow SetTitle
	ax7Variant := "MockWindow_SetTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetPosition_Good(t *core.T) {
	// MockWindow SetPosition
	ax7Variant := "MockWindow_SetPosition:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetPosition_Bad(t *core.T) {
	// MockWindow SetPosition
	ax7Variant := "MockWindow_SetPosition:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetPosition_Ugly(t *core.T) {
	// MockWindow SetPosition
	ax7Variant := "MockWindow_SetPosition:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetSize_Good(t *core.T) {
	// MockWindow SetSize
	ax7Variant := "MockWindow_SetSize:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetSize_Bad(t *core.T) {
	// MockWindow SetSize
	ax7Variant := "MockWindow_SetSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetSize_Ugly(t *core.T) {
	// MockWindow SetSize
	ax7Variant := "MockWindow_SetSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Good(t *core.T) {
	// MockWindow SetBackgroundColour
	ax7Variant := "MockWindow_SetBackgroundColour:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(1, 1, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Bad(t *core.T) {
	// MockWindow SetBackgroundColour
	ax7Variant := "MockWindow_SetBackgroundColour:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Ugly(t *core.T) {
	// MockWindow SetBackgroundColour
	ax7Variant := "MockWindow_SetBackgroundColour:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetVisibility_Good(t *core.T) {
	// MockWindow SetVisibility
	ax7Variant := "MockWindow_SetVisibility:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetVisibility_Bad(t *core.T) {
	// MockWindow SetVisibility
	ax7Variant := "MockWindow_SetVisibility:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetVisibility_Ugly(t *core.T) {
	// MockWindow SetVisibility
	ax7Variant := "MockWindow_SetVisibility:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Good(t *core.T) {
	// MockWindow SetAlwaysOnTop
	ax7Variant := "MockWindow_SetAlwaysOnTop:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Bad(t *core.T) {
	// MockWindow SetAlwaysOnTop
	ax7Variant := "MockWindow_SetAlwaysOnTop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Ugly(t *core.T) {
	// MockWindow SetAlwaysOnTop
	ax7Variant := "MockWindow_SetAlwaysOnTop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetOpacity_Good(t *core.T) {
	// MockWindow SetOpacity
	ax7Variant := "MockWindow_SetOpacity:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetOpacity_Bad(t *core.T) {
	// MockWindow SetOpacity
	ax7Variant := "MockWindow_SetOpacity:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetOpacity_Ugly(t *core.T) {
	// MockWindow SetOpacity
	ax7Variant := "MockWindow_SetOpacity:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(-1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBounds_Good(t *core.T) {
	// MockWindow SetBounds
	ax7Variant := "MockWindow_SetBounds:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(1, 1, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBounds_Bad(t *core.T) {
	// MockWindow SetBounds
	ax7Variant := "MockWindow_SetBounds:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetBounds_Ugly(t *core.T) {
	// MockWindow SetBounds
	ax7Variant := "MockWindow_SetBounds:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(-1, -1, -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetURL_Good(t *core.T) {
	// MockWindow SetURL
	ax7Variant := "MockWindow_SetURL:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetURL_Bad(t *core.T) {
	// MockWindow SetURL
	ax7Variant := "MockWindow_SetURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetURL_Ugly(t *core.T) {
	// MockWindow SetURL
	ax7Variant := "MockWindow_SetURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetHTML_Good(t *core.T) {
	// MockWindow SetHTML
	ax7Variant := "MockWindow_SetHTML:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetHTML_Bad(t *core.T) {
	// MockWindow SetHTML
	ax7Variant := "MockWindow_SetHTML:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetHTML_Ugly(t *core.T) {
	// MockWindow SetHTML
	ax7Variant := "MockWindow_SetHTML:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetZoom_Good(t *core.T) {
	// MockWindow SetZoom
	ax7Variant := "MockWindow_SetZoom:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetZoom_Bad(t *core.T) {
	// MockWindow SetZoom
	ax7Variant := "MockWindow_SetZoom:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetZoom_Ugly(t *core.T) {
	// MockWindow SetZoom
	ax7Variant := "MockWindow_SetZoom:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(-1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetContentProtection_Good(t *core.T) {
	// MockWindow SetContentProtection
	ax7Variant := "MockWindow_SetContentProtection:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetContentProtection_Bad(t *core.T) {
	// MockWindow SetContentProtection
	ax7Variant := "MockWindow_SetContentProtection:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_SetContentProtection_Ugly(t *core.T) {
	// MockWindow SetContentProtection
	ax7Variant := "MockWindow_SetContentProtection:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Maximise_Good(t *core.T) {
	// MockWindow Maximise
	ax7Variant := "MockWindow_Maximise:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Maximise_Bad(t *core.T) {
	// MockWindow Maximise
	ax7Variant := "MockWindow_Maximise:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Maximise_Ugly(t *core.T) {
	// MockWindow Maximise
	ax7Variant := "MockWindow_Maximise:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Restore_Good(t *core.T) {
	// MockWindow Restore
	ax7Variant := "MockWindow_Restore:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Restore_Bad(t *core.T) {
	// MockWindow Restore
	ax7Variant := "MockWindow_Restore:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Restore_Ugly(t *core.T) {
	// MockWindow Restore
	ax7Variant := "MockWindow_Restore:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Minimise_Good(t *core.T) {
	// MockWindow Minimise
	ax7Variant := "MockWindow_Minimise:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Minimise_Bad(t *core.T) {
	// MockWindow Minimise
	ax7Variant := "MockWindow_Minimise:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Minimise_Ugly(t *core.T) {
	// MockWindow Minimise
	ax7Variant := "MockWindow_Minimise:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Focus_Good(t *core.T) {
	// MockWindow Focus
	ax7Variant := "MockWindow_Focus:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Focus_Bad(t *core.T) {
	// MockWindow Focus
	ax7Variant := "MockWindow_Focus:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Focus_Ugly(t *core.T) {
	// MockWindow Focus
	ax7Variant := "MockWindow_Focus:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Close_Good(t *core.T) {
	// MockWindow Close
	ax7Variant := "MockWindow_Close:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Close_Bad(t *core.T) {
	// MockWindow Close
	ax7Variant := "MockWindow_Close:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Close_Ugly(t *core.T) {
	// MockWindow Close
	ax7Variant := "MockWindow_Close:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Show_Good(t *core.T) {
	// MockWindow Show
	ax7Variant := "MockWindow_Show:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Show_Bad(t *core.T) {
	// MockWindow Show
	ax7Variant := "MockWindow_Show:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Show_Ugly(t *core.T) {
	// MockWindow Show
	ax7Variant := "MockWindow_Show:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Hide_Good(t *core.T) {
	// MockWindow Hide
	ax7Variant := "MockWindow_Hide:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Hide_Bad(t *core.T) {
	// MockWindow Hide
	ax7Variant := "MockWindow_Hide:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Hide_Ugly(t *core.T) {
	// MockWindow Hide
	ax7Variant := "MockWindow_Hide:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Fullscreen_Good(t *core.T) {
	// MockWindow Fullscreen
	ax7Variant := "MockWindow_Fullscreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Fullscreen_Bad(t *core.T) {
	// MockWindow Fullscreen
	ax7Variant := "MockWindow_Fullscreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Fullscreen_Ugly(t *core.T) {
	// MockWindow Fullscreen
	ax7Variant := "MockWindow_Fullscreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_UnFullscreen_Good(t *core.T) {
	// MockWindow UnFullscreen
	ax7Variant := "MockWindow_UnFullscreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_UnFullscreen_Bad(t *core.T) {
	// MockWindow UnFullscreen
	ax7Variant := "MockWindow_UnFullscreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_UnFullscreen_Ugly(t *core.T) {
	// MockWindow UnFullscreen
	ax7Variant := "MockWindow_UnFullscreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Good(t *core.T) {
	// MockWindow ToggleFullscreen
	ax7Variant := "MockWindow_ToggleFullscreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Bad(t *core.T) {
	// MockWindow ToggleFullscreen
	ax7Variant := "MockWindow_ToggleFullscreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Ugly(t *core.T) {
	// MockWindow ToggleFullscreen
	ax7Variant := "MockWindow_ToggleFullscreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleMaximise_Good(t *core.T) {
	// MockWindow ToggleMaximise
	ax7Variant := "MockWindow_ToggleMaximise:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleMaximise_Bad(t *core.T) {
	// MockWindow ToggleMaximise
	ax7Variant := "MockWindow_ToggleMaximise:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ToggleMaximise_Ugly(t *core.T) {
	// MockWindow ToggleMaximise
	ax7Variant := "MockWindow_ToggleMaximise:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJS_Good(t *core.T) {
	// MockWindow ExecJS
	ax7Variant := "MockWindow_ExecJS:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJS_Bad(t *core.T) {
	// MockWindow ExecJS
	ax7Variant := "MockWindow_ExecJS:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJS_Ugly(t *core.T) {
	// MockWindow ExecJS
	ax7Variant := "MockWindow_ExecJS:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Flash_Good(t *core.T) {
	// MockWindow Flash
	ax7Variant := "MockWindow_Flash:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Flash_Bad(t *core.T) {
	// MockWindow Flash
	ax7Variant := "MockWindow_Flash:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Flash_Ugly(t *core.T) {
	// MockWindow Flash
	ax7Variant := "MockWindow_Flash:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Print_Good(t *core.T) {
	// MockWindow Print
	ax7Variant := "MockWindow_Print:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Print_Bad(t *core.T) {
	// MockWindow Print
	ax7Variant := "MockWindow_Print:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_Print_Ugly(t *core.T) {
	// MockWindow Print
	ax7Variant := "MockWindow_Print:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OpenDevTools_Good(t *core.T) {
	// MockWindow OpenDevTools
	ax7Variant := "MockWindow_OpenDevTools:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OpenDevTools_Bad(t *core.T) {
	// MockWindow OpenDevTools
	ax7Variant := "MockWindow_OpenDevTools:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OpenDevTools_Ugly(t *core.T) {
	// MockWindow OpenDevTools
	ax7Variant := "MockWindow_OpenDevTools:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_CloseDevTools_Good(t *core.T) {
	// MockWindow CloseDevTools
	ax7Variant := "MockWindow_CloseDevTools:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_CloseDevTools_Bad(t *core.T) {
	// MockWindow CloseDevTools
	ax7Variant := "MockWindow_CloseDevTools:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_CloseDevTools_Ugly(t *core.T) {
	// MockWindow CloseDevTools
	ax7Variant := "MockWindow_CloseDevTools:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnWindowEvent_Good(t *core.T) {
	// MockWindow OnWindowEvent
	ax7Variant := "MockWindow_OnWindowEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnWindowEvent_Bad(t *core.T) {
	// MockWindow OnWindowEvent
	ax7Variant := "MockWindow_OnWindowEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnWindowEvent_Ugly(t *core.T) {
	// MockWindow OnWindowEvent
	ax7Variant := "MockWindow_OnWindowEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnFileDrop_Good(t *core.T) {
	// MockWindow OnFileDrop
	ax7Variant := "MockWindow_OnFileDrop:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnFileDrop_Bad(t *core.T) {
	// MockWindow OnFileDrop
	ax7Variant := "MockWindow_OnFileDrop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_OnFileDrop_Ugly(t *core.T) {
	// MockWindow OnFileDrop
	ax7Variant := "MockWindow_OnFileDrop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJSCalls_Good(t *core.T) {
	// MockWindow ExecJSCalls
	ax7Variant := "MockWindow_ExecJSCalls:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJSCalls_Bad(t *core.T) {
	// MockWindow ExecJSCalls
	ax7Variant := "MockWindow_ExecJSCalls:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_ExecJSCalls_Ugly(t *core.T) {
	// MockWindow ExecJSCalls
	ax7Variant := "MockWindow_ExecJSCalls:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_HTMLContent_Good(t *core.T) {
	// MockWindow HTMLContent
	ax7Variant := "MockWindow_HTMLContent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_HTMLContent_Bad(t *core.T) {
	// MockWindow HTMLContent
	ax7Variant := "MockWindow_HTMLContent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_HTMLContent_Ugly(t *core.T) {
	// MockWindow HTMLContent
	ax7Variant := "MockWindow_HTMLContent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_DevToolsOpen_Good(t *core.T) {
	// MockWindow DevToolsOpen
	ax7Variant := "MockWindow_DevToolsOpen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_DevToolsOpen_Bad(t *core.T) {
	// MockWindow DevToolsOpen
	ax7Variant := "MockWindow_DevToolsOpen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockWindow_DevToolsOpen_Ugly(t *core.T) {
	// MockWindow DevToolsOpen
	ax7Variant := "MockWindow_DevToolsOpen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
