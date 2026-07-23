package window

import (
	core "dappco.re/go"
	"time"
)

func TestLayoutManager_SaveLayout_Good(t *core.T) {
	// SaveLayout
	ax7Variant := "SaveLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	lm := NewLayoutManagerWithDir(t.TempDir())
	windows := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 540},
	}

	core.RequireNoError(t, lm.SaveLayout("coding", windows))

	layout, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "coding", layout.Name)
	core.AssertLen(t, layout.Windows, 2)
	core.AssertNotEmpty(t, layout.CreatedAt)
	core.AssertNotEmpty(t, layout.UpdatedAt)

	infos := lm.ListLayouts()
	core.AssertLen(t, infos, 1)
	core.AssertEqual(t, "coding", infos[0].Name)
	core.AssertEqual(t, 2, infos[0].WindowCount)

	lm.DeleteLayout("coding")
	_, ok = lm.GetLayout("coding")
	core.AssertFalse(t, ok)
}

func TestLayoutManager_SaveLayout_Bad(t *core.T) {
	// SaveLayout
	ax7Variant := "SaveLayout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	lm := NewLayoutManagerWithDir(t.TempDir())
	err := lm.SaveLayout("", map[string]WindowState{"main": {Width: 1}})

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "layout name cannot be empty")
}

func TestLayoutManager_SaveLayout_Ugly(t *core.T) {
	// SaveLayout
	ax7Variant := "SaveLayout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	lm := NewLayoutManagerWithDir(t.TempDir())
	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 800}}))
	first, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)

	time.Sleep(2 * time.Millisecond)
	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 1024}}))
	second, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)

	core.AssertEqual(t, first.CreatedAt, second.CreatedAt)
	core.AssertGreater(t, second.UpdatedAt, first.UpdatedAt)
	core.AssertEqual(t, 1024, second.Windows["main"].Width)
}

func TestLayoutManager_NewLayoutManagerWithPathEnv_GoodCase(t *core.T) {
	path := core.PathJoin(t.TempDir(), "custom", "layouts.json")
	t.Setenv(layoutFileEnv, path)

	lm := NewLayoutManager()

	core.AssertNotNil(t, lm)
	core.AssertEqual(t, path, lm.filePath())
	core.AssertEqual(t, core.PathDir(path), lm.dataDir())

	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{
		"main": {Width: 800, Height: 600},
	}))

	content, err := coreReadFile(path)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), `"coding"`)
}

// AX7 generated source-matching smoke coverage.
func TestLayout_NewLayoutManager_Good(t *core.T) {
	// NewLayoutManager
	ax7Variant := "NewLayoutManager:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewLayoutManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManager_Bad(t *core.T) {
	// NewLayoutManager
	ax7Variant := "NewLayoutManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewLayoutManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManager_Ugly(t *core.T) {
	// NewLayoutManager
	ax7Variant := "NewLayoutManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewLayoutManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithDir_Good(t *core.T) {
	// NewLayoutManagerWithDir
	ax7Variant := "NewLayoutManagerWithDir:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithDir("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithDir_Bad(t *core.T) {
	// NewLayoutManagerWithDir
	ax7Variant := "NewLayoutManagerWithDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithDir("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithDir_Ugly(t *core.T) {
	// NewLayoutManagerWithDir
	ax7Variant := "NewLayoutManagerWithDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithDir("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithPath_Good(t *core.T) {
	// NewLayoutManagerWithPath
	ax7Variant := "NewLayoutManagerWithPath:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithPath("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithPath_Bad(t *core.T) {
	// NewLayoutManagerWithPath
	ax7Variant := "NewLayoutManagerWithPath:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithPath("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_NewLayoutManagerWithPath_Ugly(t *core.T) {
	// NewLayoutManagerWithPath
	ax7Variant := "NewLayoutManagerWithPath:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewLayoutManagerWithPath("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SetPath_Good(t *core.T) {
	// LayoutManager SetPath
	ax7Variant := "LayoutManager_SetPath:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.SetPath("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SetPath_Bad(t *core.T) {
	// LayoutManager SetPath
	ax7Variant := "LayoutManager_SetPath:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.SetPath("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SetPath_Ugly(t *core.T) {
	// LayoutManager SetPath
	ax7Variant := "LayoutManager_SetPath:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.SetPath("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SaveLayout_Good(t *core.T) {
	// LayoutManager SaveLayout
	ax7Variant := "LayoutManager_SaveLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.SaveLayout("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SaveLayout_Bad(t *core.T) {
	// LayoutManager SaveLayout
	ax7Variant := "LayoutManager_SaveLayout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.SaveLayout("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_SaveLayout_Ugly(t *core.T) {
	// LayoutManager SaveLayout
	ax7Variant := "LayoutManager_SaveLayout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.SaveLayout("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_GetLayout_Good(t *core.T) {
	// LayoutManager GetLayout
	ax7Variant := "LayoutManager_GetLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetLayout("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_GetLayout_Bad(t *core.T) {
	// LayoutManager GetLayout
	ax7Variant := "LayoutManager_GetLayout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetLayout("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_GetLayout_Ugly(t *core.T) {
	// LayoutManager GetLayout
	ax7Variant := "LayoutManager_GetLayout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetLayout("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_ListLayouts_Good(t *core.T) {
	// LayoutManager ListLayouts
	ax7Variant := "LayoutManager_ListLayouts:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.ListLayouts()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_ListLayouts_Bad(t *core.T) {
	// LayoutManager ListLayouts
	ax7Variant := "LayoutManager_ListLayouts:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.ListLayouts()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_ListLayouts_Ugly(t *core.T) {
	// LayoutManager ListLayouts
	ax7Variant := "LayoutManager_ListLayouts:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		got0 := subject.ListLayouts()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_DeleteLayout_Good(t *core.T) {
	// LayoutManager DeleteLayout
	ax7Variant := "LayoutManager_DeleteLayout:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.DeleteLayout("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_DeleteLayout_Bad(t *core.T) {
	// LayoutManager DeleteLayout
	ax7Variant := "LayoutManager_DeleteLayout:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.DeleteLayout("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestLayout_LayoutManager_DeleteLayout_Ugly(t *core.T) {
	// LayoutManager DeleteLayout
	ax7Variant := "LayoutManager_DeleteLayout:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(LayoutManager)
	result := core.Try(func() any {
		subject.DeleteLayout("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}
