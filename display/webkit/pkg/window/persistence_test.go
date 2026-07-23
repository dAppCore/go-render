// pkg/window/persistence_test.go
package window

import (
	core "dappco.re/go"
	"time"
)

// --- StateManager Persistence Tests ---

func TestStateManager_SetAndGet_GoodCase(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	state := WindowState{
		X: 150, Y: 250, Width: 1024, Height: 768,
		Maximized: true, Screen: "primary", URL: "/app",
	}
	sm.SetState("editor", state)

	got, ok := sm.GetState("editor")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 150, got.X)
	core.AssertEqual(t, 250, got.Y)
	core.AssertEqual(t, 1024, got.Width)
	core.AssertEqual(t, 768, got.Height)
	core.AssertTrue(t, got.Maximized)
	core.AssertEqual(t, "primary", got.Screen)
	core.AssertEqual(t, "/app", got.URL)
	core.AssertNotEmpty(t, got.UpdatedAt, "UpdatedAt should be set by SetState")
}

func TestStateManager_UpdatePosition_Good(t *core.T) {
	// UpdatePosition
	ax7Variant := "UpdatePosition:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{X: 0, Y: 0, Width: 800, Height: 600})

	sm.UpdatePosition("win", 300, 400)

	got, ok := sm.GetState("win")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 300, got.X)
	core.AssertEqual(t, 400, got.Y)
	// Width/Height should remain unchanged
	core.AssertEqual(t, 800, got.Width)
	core.AssertEqual(t, 600, got.Height)
}

func TestStateManager_UpdateSize_Good(t *core.T) {
	// UpdateSize
	ax7Variant := "UpdateSize:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{X: 100, Y: 200, Width: 800, Height: 600})

	sm.UpdateSize("win", 1920, 1080)

	got, ok := sm.GetState("win")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 1920, got.Width)
	core.AssertEqual(t, 1080, got.Height)
	// Position should remain unchanged
	core.AssertEqual(t, 100, got.X)
	core.AssertEqual(t, 200, got.Y)
}

func TestStateManager_UpdateMaximized_Good(t *core.T) {
	// UpdateMaximized
	ax7Variant := "UpdateMaximized:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{Width: 800, Height: 600, Maximized: false})

	sm.UpdateMaximized("win", true)

	got, ok := sm.GetState("win")
	core.RequireTrue(t, ok)
	core.AssertTrue(t, got.Maximized)

	sm.UpdateMaximized("win", false)

	got, ok = sm.GetState("win")
	core.RequireTrue(t, ok)
	core.AssertFalse(t, got.Maximized)
}

func TestStateManager_CaptureState_Good(t *core.T) {
	// CaptureState
	ax7Variant := "CaptureState:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	pw := &mockWindow{
		name: "captured", x: 75, y: 125,
		width: 1440, height: 900, maximised: true,
	}

	sm.CaptureState(pw)

	got, ok := sm.GetState("captured")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 75, got.X)
	core.AssertEqual(t, 125, got.Y)
	core.AssertEqual(t, 1440, got.Width)
	core.AssertEqual(t, 900, got.Height)
	core.AssertTrue(t, got.Maximized)
	core.AssertNotEmpty(t, got.UpdatedAt)
}

func TestStateManager_ApplyState_Good(t *core.T) {
	// ApplyState
	ax7Variant := "ApplyState:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("target", WindowState{X: 55, Y: 65, Width: 700, Height: 500})

	w := &Window{Name: "target", Width: 1280, Height: 800, X: 0, Y: 0}
	sm.ApplyState(w)

	core.AssertEqual(t, 55, w.X)
	core.AssertEqual(t, 65, w.Y)
	core.AssertEqual(t, 700, w.Width)
	core.AssertEqual(t, 500, w.Height)
}

func TestStateManager_ApplyState_Good_NoState(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())

	w := &Window{Name: "untouched", Width: 1280, Height: 800, X: 10, Y: 20}
	sm.ApplyState(w)

	// Window should remain unchanged when no state is saved
	core.AssertEqual(t, 10, w.X)
	core.AssertEqual(t, 20, w.Y)
	core.AssertEqual(t, 1280, w.Width)
	core.AssertEqual(t, 800, w.Height)
}

func TestStateManager_ListStates_Good(t *core.T) {
	// ListStates
	ax7Variant := "ListStates:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{Width: 100})
	sm.SetState("beta", WindowState{Width: 200})
	sm.SetState("gamma", WindowState{Width: 300})

	names := sm.ListStates()
	core.AssertLen(t, names, 3)
	core.AssertContains(t, names, "alpha")
	core.AssertContains(t, names, "beta")
	core.AssertContains(t, names, "gamma")
}

func TestStateManager_Clear_Good(t *core.T) {
	// Clear
	ax7Variant := "Clear:good"
	core.AssertContains(t, ax7Variant, "good")
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("a", WindowState{Width: 100})
	sm.SetState("b", WindowState{Width: 200})
	sm.SetState("c", WindowState{Width: 300})

	sm.Clear()

	names := sm.ListStates()
	core.AssertEmpty(t, names)

	_, ok := sm.GetState("a")
	core.AssertFalse(t, ok)
}

func TestStateManager_Persistence_GoodCase(t *core.T) {
	dir := t.TempDir()

	// First manager: write state and force sync to disk
	sm1 := NewStateManagerWithDir(dir)
	sm1.SetState("persist-win", WindowState{
		X: 42, Y: 84, Width: 500, Height: 300,
		Maximized: true, Screen: "secondary", URL: "/settings",
	})
	sm1.ForceSync()

	// Second manager: load from the same directory
	sm2 := NewStateManagerWithDir(dir)

	got, ok := sm2.GetState("persist-win")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 42, got.X)
	core.AssertEqual(t, 84, got.Y)
	core.AssertEqual(t, 500, got.Width)
	core.AssertEqual(t, 300, got.Height)
	core.AssertTrue(t, got.Maximized)
	core.AssertEqual(t, "secondary", got.Screen)
	core.AssertEqual(t, "/settings", got.URL)
	core.AssertNotEmpty(t, got.UpdatedAt)
}

func TestStateManager_SetPath_Good(t *core.T) {
	// SetPath
	ax7Variant := "SetPath:good"
	core.AssertContains(t, ax7Variant, "good")
	dir := t.TempDir()
	path := core.JoinPath(dir, "custom", "window-state.json")

	sm := NewStateManagerWithDir(dir)
	sm.SetPath(path)
	sm.SetState("custom", WindowState{Width: 640, Height: 480})
	sm.ForceSync()

	content, err := coreReadFile(path)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), "custom")
}

// --- LayoutManager Persistence Tests ---

func TestLayoutManager_SaveAndGet_GoodCase(t *core.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	windows := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 540},
		"browser":  {X: 960, Y: 540, Width: 960, Height: 540},
	}

	err := lm.SaveLayout("coding", windows)
	core.RequireNoError(t, err)

	layout, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "coding", layout.Name)
	core.AssertLen(t, layout.Windows, 3)
	core.AssertEqual(t, 960, layout.Windows["editor"].Width)
	core.AssertEqual(t, 1080, layout.Windows["editor"].Height)
	core.AssertEqual(t, 960, layout.Windows["terminal"].X)
	core.AssertNotEmpty(t, layout.CreatedAt)
	core.AssertNotEmpty(t, layout.UpdatedAt)
	core.AssertEqual(t, layout.CreatedAt, layout.UpdatedAt, "CreatedAt and UpdatedAt should match on first save")
}

func TestLayoutManager_SaveLayout_EmptyName_Bad(t *core.T) {
	// SaveLayout EmptyName
	ax7Variant := "SaveLayout_EmptyName:bad"
	core.AssertContains(t, ax7Variant, "bad")
	lm := NewLayoutManagerWithDir(t.TempDir())
	err := lm.SaveLayout("", map[string]WindowState{
		"win": {Width: 800},
	})
	core.AssertError(t, err)
}

func TestLayoutManager_SaveLayout_Update_Good(t *core.T) {
	// SaveLayout Update
	ax7Variant := "SaveLayout_Update:good"
	core.AssertContains(t, ax7Variant, "good")
	lm := NewLayoutManagerWithDir(t.TempDir())

	// First save
	err := lm.SaveLayout("evolving", map[string]WindowState{
		"win1": {Width: 800, Height: 600},
	})
	core.RequireNoError(t, err)

	first, ok := lm.GetLayout("evolving")
	core.RequireTrue(t, ok)
	originalCreatedAt := first.CreatedAt
	originalUpdatedAt := first.UpdatedAt

	// Small delay to ensure UpdatedAt differs
	time.Sleep(2 * time.Millisecond)

	// Second save with same name but different windows
	err = lm.SaveLayout("evolving", map[string]WindowState{
		"win1": {Width: 1024, Height: 768},
		"win2": {Width: 640, Height: 480},
	})
	core.RequireNoError(t, err)

	updated, ok := lm.GetLayout("evolving")
	core.RequireTrue(t, ok)

	// CreatedAt should be preserved from the original save
	core.AssertEqual(t, originalCreatedAt, updated.CreatedAt, "CreatedAt should be preserved on update")
	// UpdatedAt should be newer
	core.AssertGreaterOrEqual(t, updated.UpdatedAt, originalUpdatedAt, "UpdatedAt should advance on update")
	// Windows should reflect the second save
	core.AssertLen(t, updated.Windows, 2)
	core.AssertEqual(t, 1024, updated.Windows["win1"].Width)
}

func TestLayoutManager_ListLayouts_Good(t *core.T) {
	// ListLayouts
	ax7Variant := "ListLayouts:good"
	core.AssertContains(t, ax7Variant, "good")
	lm := NewLayoutManagerWithDir(t.TempDir())
	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{
		"editor": {Width: 960}, "terminal": {Width: 960},
	}))
	core.RequireNoError(t, lm.SaveLayout("presenting", map[string]WindowState{
		"slides": {Width: 1920},
	}))
	core.RequireNoError(t, lm.SaveLayout("debugging", map[string]WindowState{
		"code": {Width: 640}, "debugger": {Width: 640}, "console": {Width: 640},
	}))

	infos := lm.ListLayouts()
	core.AssertLen(t, infos, 3)

	// Build a lookup map for assertions regardless of order
	byName := make(map[string]LayoutInfo)
	for _, info := range infos {
		byName[info.Name] = info
	}

	core.AssertEqual(t, 2, byName["coding"].WindowCount)
	core.AssertEqual(t, 1, byName["presenting"].WindowCount)
	core.AssertEqual(t, 3, byName["debugging"].WindowCount)
}

func TestLayoutManager_DeleteLayout_Good(t *core.T) {
	// DeleteLayout
	ax7Variant := "DeleteLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	lm := NewLayoutManagerWithDir(t.TempDir())
	core.RequireNoError(t, lm.SaveLayout("temporary", map[string]WindowState{
		"win": {Width: 800},
	}))

	// Verify it exists
	_, ok := lm.GetLayout("temporary")
	core.RequireTrue(t, ok)

	lm.DeleteLayout("temporary")

	// Verify it is gone
	_, ok = lm.GetLayout("temporary")
	core.AssertFalse(t, ok)

	// Verify list is empty
	core.AssertEmpty(t, lm.ListLayouts())
}

func TestLayoutManager_Persistence_GoodCase(t *core.T) {
	dir := t.TempDir()

	// First manager: save layout to disk
	lm1 := NewLayoutManagerWithDir(dir)
	err := lm1.SaveLayout("persisted", map[string]WindowState{
		"main":    {X: 0, Y: 0, Width: 1280, Height: 800},
		"sidebar": {X: 1280, Y: 0, Width: 640, Height: 800},
	})
	core.RequireNoError(t, err)

	// Second manager: load from the same directory
	lm2 := NewLayoutManagerWithDir(dir)

	layout, ok := lm2.GetLayout("persisted")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "persisted", layout.Name)
	core.AssertLen(t, layout.Windows, 2)
	core.AssertEqual(t, 1280, layout.Windows["main"].Width)
	core.AssertEqual(t, 800, layout.Windows["main"].Height)
	core.AssertEqual(t, 640, layout.Windows["sidebar"].Width)
	core.AssertNotEmpty(t, layout.CreatedAt)
	core.AssertNotEmpty(t, layout.UpdatedAt)
}
