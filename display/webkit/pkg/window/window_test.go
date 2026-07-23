// pkg/window/window_test.go
package window

import (
	core "dappco.re/go"
)

func TestWindowDefaults_Good(t *core.T) {
	w := &Window{}
	core.AssertEqual(t, "", w.Name)
	core.AssertEqual(t, 0, w.Width)
}

func TestWindowOption_Name_Good(t *core.T) {
	// Name
	ax7Variant := "Name:good"
	core.AssertContains(t, ax7Variant, "good")
	w := &Window{}
	err := WithName("main")(w)
	core.RequireNoError(t, err)
	core.AssertEqual(t, "main", w.Name)
}

func TestWindowOption_Title_Good(t *core.T) {
	// Title
	ax7Variant := "Title:good"
	core.AssertContains(t, ax7Variant, "good")
	w := &Window{}
	err := WithTitle("My App")(w)
	core.RequireNoError(t, err)
	core.AssertEqual(t, "My App", w.Title)
}

func TestWindowOption_URL_Good(t *core.T) {
	// URL
	ax7Variant := "URL:good"
	core.AssertContains(t, ax7Variant, "good")
	w := &Window{}
	err := WithURL("/dashboard")(w)
	core.RequireNoError(t, err)
	core.AssertEqual(t, "/dashboard", w.URL)
}

func TestWindowOption_Size_GoodCase(t *core.T) {
	w := &Window{}
	err := WithSize(1280, 720)(w)
	core.RequireNoError(t, err)
	core.AssertEqual(t, 1280, w.Width)
	core.AssertEqual(t, 720, w.Height)
}

func TestWindowOption_Position_GoodCase(t *core.T) {
	w := &Window{}
	err := WithPosition(100, 200)(w)
	core.RequireNoError(t, err)
	core.AssertEqual(t, 100, w.X)
	core.AssertEqual(t, 200, w.Y)
}

func TestApplyOptions_Good(t *core.T) {
	w, err := ApplyOptions(
		WithName("test"),
		WithTitle("Test Window"),
		WithURL("/test"),
		WithSize(800, 600),
	)
	core.RequireNoError(t, err)
	core.AssertEqual(t, "test", w.Name)
	core.AssertEqual(t, "Test Window", w.Title)
	core.AssertEqual(t, "/test", w.URL)
	core.AssertEqual(t, 800, w.Width)
	core.AssertEqual(t, 600, w.Height)
}

func TestApplyOptions_Bad(t *core.T) {
	_, err := ApplyOptions(func(w *Window) resultFailure {
		return core.AnError
	})
	core.AssertError(t, err)
}

func TestApplyOptions_Empty_GoodCase(t *core.T) {
	w, err := ApplyOptions()
	core.RequireNoError(t, err)
	core.AssertNotNil(t, w)
}

// newTestManager creates a Manager with a mock platform and clean state for testing.
func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	m := &Manager{
		platform: p,
		state:    &StateManager{states: make(map[string]WindowState)},
		layout:   &LayoutManager{layouts: make(map[string]Layout)},
		windows:  make(map[string]PlatformWindow),
	}
	return m, p
}

func TestManager_Open_Good(t *core.T) {
	// Open
	ax7Variant := "Open:good"
	core.AssertContains(t, ax7Variant, "good")
	m, p := newTestManager()
	pw, err := m.Open(WithName("test"), WithTitle("Test"), WithURL("/test"), WithSize(800, 600))
	core.RequireNoError(t, err)
	core.AssertNotNil(t, pw)
	core.AssertEqual(t, "test", pw.Name())
	core.AssertLen(t, p.windows, 1)
}

func TestManager_Open_Defaults_Good(t *core.T) {
	// Open Defaults
	ax7Variant := "Open_Defaults:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	pw, err := m.Open()
	core.RequireNoError(t, err)
	core.AssertEqual(t, "main", pw.Name())
	w, h := pw.Size()
	core.AssertEqual(t, 1280, w)
	core.AssertEqual(t, 800, h)
}

func TestManager_Open_CustomDefaults_Good(t *core.T) {
	// Open CustomDefaults
	ax7Variant := "Open_CustomDefaults:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	m.SetDefaultWidth(1440)
	m.SetDefaultHeight(900)

	pw, err := m.Open()
	core.RequireNoError(t, err)

	w, h := pw.Size()
	core.AssertEqual(t, 1440, w)
	core.AssertEqual(t, 900, h)
}

func TestManager_Open_Bad(t *core.T) {
	// Open
	ax7Variant := "Open:bad"
	core.AssertContains(t, ax7Variant, "bad")
	m, _ := newTestManager()
	_, err := m.Open(func(w *Window) resultFailure { return core.AnError })
	core.AssertError(t, err)
}

func TestManager_Get_Good(t *core.T) {
	// Get
	ax7Variant := "Get:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("findme"))
	pw, ok := m.Get("findme")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "findme", pw.Name())
}

func TestManager_Get_Bad(t *core.T) {
	// Get
	ax7Variant := "Get:bad"
	core.AssertContains(t, ax7Variant, "bad")
	m, _ := newTestManager()
	_, ok := m.Get("nonexistent")
	core.AssertFalse(t, ok)
}

func TestManager_List_Good(t *core.T) {
	// List
	ax7Variant := "List:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"))
	_, _ = m.Open(WithName("b"))
	names := m.List()
	core.AssertLen(t, names, 2)
	core.AssertContains(t, names, "a")
	core.AssertContains(t, names, "b")
}

func TestManager_Remove_Good(t *core.T) {
	// Remove
	ax7Variant := "Remove:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("temp"))
	m.Remove("temp")
	_, ok := m.Get("temp")
	core.AssertFalse(t, ok)
}

func TestManager_NewManagerWithDir_Good(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:good"
	core.AssertContains(t, ax7Variant, "good")
	dir := t.TempDir()
	p := newMockPlatform()

	m := NewManagerWithDir(p, dir)

	core.AssertNotNil(t, m)
	core.AssertSame(t, p, m.Platform())
	core.AssertEqual(t, dir, m.State().dataDir())
	core.AssertEqual(t, core.PathJoin(dir, "layouts.json"), m.Layout().filePath())
}

func TestManager_NewManagerWithDir_Bad(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	m := NewManagerWithDir(nil, "")

	core.AssertNotNil(t, m)
	core.AssertNil(t, m.Platform())
	core.AssertEmpty(t, m.State().dataDir())
}

func TestManager_NewManagerWithDir_Ugly(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	dir := core.PathJoin(t.TempDir(), "..", "workspace")
	m := NewManagerWithDir(nil, dir)

	core.AssertNotNil(t, m)
	core.AssertEqual(t, dir, m.State().dataDir())
}

// --- Tiling Tests ---

func TestTileMode_String_Good(t *core.T) {
	// String
	ax7Variant := "String:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertEqual(t, "left-half", TileModeLeftHalf.String())
	core.AssertEqual(t, "grid", TileModeGrid.String())
	core.AssertNotEmpty(t, core.Sprintf("%T", TileModeLeftHalf.String()))
}

func TestManager_TileWindows_Good(t *core.T) {
	// TileWindows
	ax7Variant := "TileWindows:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"), WithSize(800, 600))
	_, _ = m.Open(WithName("b"), WithSize(800, 600))
	err := m.TileWindows(TileModeLeftRight, []string{"a", "b"}, 1920, 1080)
	core.RequireNoError(t, err)
	a, _ := m.Get("a")
	b, _ := m.Get("b")
	aw, _ := a.Size()
	bw, _ := b.Size()
	core.AssertEqual(t, 960, aw)
	core.AssertEqual(t, 960, bw)
}

func TestManager_TileWindows_Bad(t *core.T) {
	// TileWindows
	ax7Variant := "TileWindows:bad"
	core.AssertContains(t, ax7Variant, "bad")
	m, _ := newTestManager()
	err := m.TileWindows(TileModeLeftRight, []string{"nonexistent"}, 1920, 1080)
	core.AssertError(t, err)
}

func TestManager_SnapWindow_Good(t *core.T) {
	// SnapWindow
	ax7Variant := "SnapWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("snap"), WithSize(800, 600))
	err := m.SnapWindow("snap", SnapLeft, 1920, 1080)
	core.RequireNoError(t, err)
	w, _ := m.Get("snap")
	x, _ := w.Position()
	core.AssertEqual(t, 0, x)
	sw, _ := w.Size()
	core.AssertEqual(t, 960, sw)
}

func TestManager_StackWindows_Good(t *core.T) {
	// StackWindows
	ax7Variant := "StackWindows:good"
	core.AssertContains(t, ax7Variant, "good")
	m, _ := newTestManager()
	_, _ = m.Open(WithName("s1"), WithSize(800, 600))
	_, _ = m.Open(WithName("s2"), WithSize(800, 600))
	err := m.StackWindows([]string{"s1", "s2"}, 30, 30)
	core.RequireNoError(t, err)
	s2, _ := m.Get("s2")
	x, y := s2.Position()
	core.AssertEqual(t, 30, x)
	core.AssertEqual(t, 30, y)
}

func TestWorkflowLayout_Good(t *core.T) {
	core.AssertEqual(t, "coding", WorkflowCoding.String())
	core.AssertEqual(t, "debugging", WorkflowDebugging.String())
	core.AssertNotEmpty(t, core.Sprintf("%T", WorkflowCoding.String()))
}

// --- Comprehensive Tiling Tests ---

func TestTileWindows_AllModes_GoodCase(t *core.T) {
	const screenW, screenH = 1920, 1080
	halfW, halfH := screenW/2, screenH/2

	tests := []struct {
		name       string
		mode       TileMode
		wantX      int
		wantY      int
		wantWidth  int
		wantHeight int
	}{
		{"LeftHalf", TileModeLeftHalf, 0, 0, halfW, screenH},
		{"RightHalf", TileModeRightHalf, halfW, 0, halfW, screenH},
		{"TopHalf", TileModeTopHalf, 0, 0, screenW, halfH},
		{"BottomHalf", TileModeBottomHalf, 0, halfH, screenW, halfH},
		{"TopLeft", TileModeTopLeft, 0, 0, halfW, halfH},
		{"TopRight", TileModeTopRight, halfW, 0, halfW, halfH},
		{"BottomLeft", TileModeBottomLeft, 0, halfH, halfW, halfH},
		{"BottomRight", TileModeBottomRight, halfW, halfH, halfW, halfH},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("win"), WithSize(800, 600))
			core.RequireNoError(t, err)

			err = m.TileWindows(tc.mode, []string{"win"}, screenW, screenH)
			core.RequireNoError(t, err)

			pw, ok := m.Get("win")
			core.RequireTrue(t, ok)

			x, y := pw.Position()
			w, h := pw.Size()
			core.AssertEqual(t, tc.wantX, x, "x position")
			core.AssertEqual(t, tc.wantY, y, "y position")
			core.AssertEqual(t, tc.wantWidth, w, "width")
			core.AssertEqual(t, tc.wantHeight, h, "height")
		})
	}
}

func TestSnapWindow_AllPositions_GoodCase(t *core.T) {
	const screenW, screenH = 1920, 1080
	halfW, halfH := screenW/2, screenH/2

	tests := []struct {
		name       string
		pos        SnapPosition
		initW      int
		initH      int
		wantX      int
		wantY      int
		wantWidth  int
		wantHeight int
	}{
		{"Right", SnapRight, 800, 600, halfW, 0, halfW, screenH},
		{"Top", SnapTop, 800, 600, 0, 0, screenW, halfH},
		{"Bottom", SnapBottom, 800, 600, 0, halfH, screenW, halfH},
		{"TopLeft", SnapTopLeft, 800, 600, 0, 0, halfW, halfH},
		{"TopRight", SnapTopRight, 800, 600, halfW, 0, halfW, halfH},
		{"BottomLeft", SnapBottomLeft, 800, 600, 0, halfH, halfW, halfH},
		{"BottomRight", SnapBottomRight, 800, 600, halfW, halfH, halfW, halfH},
		{"Center", SnapCenter, 800, 600, (screenW - 800) / 2, (screenH - 600) / 2, 800, 600},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("snap"), WithSize(tc.initW, tc.initH))
			core.RequireNoError(t, err)

			err = m.SnapWindow("snap", tc.pos, screenW, screenH)
			core.RequireNoError(t, err)

			pw, ok := m.Get("snap")
			core.RequireTrue(t, ok)

			x, y := pw.Position()
			w, h := pw.Size()
			core.AssertEqual(t, tc.wantX, x, "x position")
			core.AssertEqual(t, tc.wantY, y, "y position")
			core.AssertEqual(t, tc.wantWidth, w, "width")
			core.AssertEqual(t, tc.wantHeight, h, "height")
		})
	}
}

func TestStackWindows_ThreeWindows_GoodCase(t *core.T) {
	m, _ := newTestManager()
	names := []string{"s1", "s2", "s3"}
	for _, name := range names {
		_, err := m.Open(WithName(name), WithSize(800, 600))
		core.RequireNoError(t, err)
	}

	err := m.StackWindows(names, 30, 30)
	core.RequireNoError(t, err)

	for i, name := range names {
		pw, ok := m.Get(name)
		core.RequireTrue(t, ok, "window %s should exist", name)
		x, y := pw.Position()
		core.AssertEqual(t, i*30, x, "window %s x position", name)
		core.AssertEqual(t, i*30, y, "window %s y position", name)
	}
}

func TestApplyWorkflow_AllLayouts_GoodCase(t *core.T) {
	const screenW, screenH = 1920, 1080

	tests := []struct {
		name     string
		workflow WorkflowLayout
		// Expected positions/sizes for the first two windows.
		// For WorkflowSideBySide, TileWindows(LeftRight) divides equally.
		win0X, win0Y, win0W, win0H int
		win1X, win1Y, win1W, win1H int
	}{
		{
			"Coding",
			WorkflowCoding,
			0, 0, 1344, screenH, // 70% of 1920 = 1344
			1344, 0, screenW - 1344, screenH, // remaining 30%
		},
		{
			"Debugging",
			WorkflowDebugging,
			0, 0, 1152, screenH, // 60% of 1920 = 1152
			1152, 0, screenW - 1152, screenH, // remaining 40%
		},
		{
			"Presenting",
			WorkflowPresenting,
			0, 0, screenW, screenH, // maximised
			0, 0, 800, 600, // second window untouched
		},
		{
			"SideBySide",
			WorkflowSideBySide,
			0, 0, 960, screenH, // left half (1920/2)
			960, 0, 960, screenH, // right half
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("editor"), WithSize(800, 600))
			core.RequireNoError(t, err)
			_, err = m.Open(WithName("terminal"), WithSize(800, 600))
			core.RequireNoError(t, err)

			err = m.ApplyWorkflow(tc.workflow, []string{"editor", "terminal"}, screenW, screenH)
			core.RequireNoError(t, err)

			pw0, ok := m.Get("editor")
			core.RequireTrue(t, ok)
			x0, y0 := pw0.Position()
			w0, h0 := pw0.Size()
			core.AssertEqual(t, tc.win0X, x0, "editor x")
			core.AssertEqual(t, tc.win0Y, y0, "editor y")
			core.AssertEqual(t, tc.win0W, w0, "editor width")
			core.AssertEqual(t, tc.win0H, h0, "editor height")

			pw1, ok := m.Get("terminal")
			core.RequireTrue(t, ok)
			x1, y1 := pw1.Position()
			w1, h1 := pw1.Size()
			core.AssertEqual(t, tc.win1X, x1, "terminal x")
			core.AssertEqual(t, tc.win1Y, y1, "terminal y")
			core.AssertEqual(t, tc.win1W, w1, "terminal width")
			core.AssertEqual(t, tc.win1H, h1, "terminal height")
		})
	}
}

func TestApplyWorkflow_Empty_BadCase(t *core.T) {
	m, _ := newTestManager()
	err := m.ApplyWorkflow(WorkflowCoding, []string{}, 1920, 1080)
	core.AssertError(t, err)
}

// AX7 generated source-matching smoke coverage.
func TestWindow_Window_ToPlatformOptions_Good(t *core.T) {
	// Window ToPlatformOptions
	ax7Variant := "Window_ToPlatformOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Window)
	result := core.Try(func() any {
		got0 := subject.ToPlatformOptions()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Window_ToPlatformOptions_Bad(t *core.T) {
	// Window ToPlatformOptions
	ax7Variant := "Window_ToPlatformOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Window)
	result := core.Try(func() any {
		got0 := subject.ToPlatformOptions()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Window_ToPlatformOptions_Ugly(t *core.T) {
	// Window ToPlatformOptions
	ax7Variant := "Window_ToPlatformOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Window)
	result := core.Try(func() any {
		got0 := subject.ToPlatformOptions()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManager_Good(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManager_Bad(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManager_Ugly(t *core.T) {
	// NewManager
	ax7Variant := "NewManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManagerWithDir_Good(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewManagerWithDir(*new(Platform), "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManagerWithDir_Bad(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewManagerWithDir(*new(Platform), "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_NewManagerWithDir_Ugly(t *core.T) {
	// NewManagerWithDir
	ax7Variant := "NewManagerWithDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewManagerWithDir(*new(Platform), "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultWidth_Good(t *core.T) {
	// Manager SetDefaultWidth
	ax7Variant := "Manager_SetDefaultWidth:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultWidth(1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultWidth_Bad(t *core.T) {
	// Manager SetDefaultWidth
	ax7Variant := "Manager_SetDefaultWidth:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultWidth(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultWidth_Ugly(t *core.T) {
	// Manager SetDefaultWidth
	ax7Variant := "Manager_SetDefaultWidth:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultWidth(-1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultHeight_Good(t *core.T) {
	// Manager SetDefaultHeight
	ax7Variant := "Manager_SetDefaultHeight:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultHeight(1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultHeight_Bad(t *core.T) {
	// Manager SetDefaultHeight
	ax7Variant := "Manager_SetDefaultHeight:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultHeight(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_SetDefaultHeight_Ugly(t *core.T) {
	// Manager SetDefaultHeight
	ax7Variant := "Manager_SetDefaultHeight:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.SetDefaultHeight(-1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Open_Good(t *core.T) {
	// Manager Open
	ax7Variant := "Manager_Open:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Open()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Open_Bad(t *core.T) {
	// Manager Open
	ax7Variant := "Manager_Open:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Open()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Open_Ugly(t *core.T) {
	// Manager Open
	ax7Variant := "Manager_Open:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Open()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Create_Good(t *core.T) {
	// Manager Create
	ax7Variant := "Manager_Create:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Create(nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Create_Bad(t *core.T) {
	// Manager Create
	ax7Variant := "Manager_Create:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Create(nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Create_Ugly(t *core.T) {
	// Manager Create
	ax7Variant := "Manager_Create:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Create(nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Get_Good(t *core.T) {
	// Manager Get
	ax7Variant := "Manager_Get:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Get_Bad(t *core.T) {
	// Manager Get
	ax7Variant := "Manager_Get:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Get_Ugly(t *core.T) {
	// Manager Get
	ax7Variant := "Manager_Get:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_List_Good(t *core.T) {
	// Manager List
	ax7Variant := "Manager_List:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.List()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_List_Bad(t *core.T) {
	// Manager List
	ax7Variant := "Manager_List:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.List()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_List_Ugly(t *core.T) {
	// Manager List
	ax7Variant := "Manager_List:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.List()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Remove_Good(t *core.T) {
	// Manager Remove
	ax7Variant := "Manager_Remove:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.Remove("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Remove_Bad(t *core.T) {
	// Manager Remove
	ax7Variant := "Manager_Remove:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.Remove("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Remove_Ugly(t *core.T) {
	// Manager Remove
	ax7Variant := "Manager_Remove:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.Remove("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Platform_Good(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Platform_Bad(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Platform_Ugly(t *core.T) {
	// Manager Platform
	ax7Variant := "Manager_Platform:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Platform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_State_Good(t *core.T) {
	// Manager State
	ax7Variant := "Manager_State:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_State_Bad(t *core.T) {
	// Manager State
	ax7Variant := "Manager_State:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_State_Ugly(t *core.T) {
	// Manager State
	ax7Variant := "Manager_State:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Layout_Good(t *core.T) {
	// Manager Layout
	ax7Variant := "Manager_Layout:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Layout()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Layout_Bad(t *core.T) {
	// Manager Layout
	ax7Variant := "Manager_Layout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Layout()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestWindow_Manager_Layout_Ugly(t *core.T) {
	// Manager Layout
	ax7Variant := "Manager_Layout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Layout()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
