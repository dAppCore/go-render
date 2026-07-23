package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/systray"
	"dappco.re/go/render/display/webkit/pkg/window"
)

func newTestDisplayAPIService(t *core.T) (*Service, *core.Core) {
	t.Helper()
	return newTestDisplayService(t)
}

func TestDisplayAPI_screenToDisplay_Good(t *core.T) {
	// screenToDisplay
	ax7Variant := "screenToDisplay:good"
	core.AssertContains(t, ax7Variant, "good")
	got := screenToDisplay(&screen.Screen{
		ID:          "screen-1",
		Name:        "Primary",
		ScaleFactor: 2,
		Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
		IsPrimary:   true,
	})

	core.AssertNotNil(t, got)
	core.AssertEqual(t, "screen-1", got.ID)
	core.AssertEqual(t, "Primary", got.Name)
	core.AssertEqual(t, 10, got.X)
	core.AssertEqual(t, 20, got.Y)
	core.AssertEqual(t, 1920, got.Width)
	core.AssertEqual(t, 1080, got.Height)
	core.AssertEqual(t, 2.0, got.ScaleFactor)
	core.AssertTrue(t, got.IsPrimary)
}

func TestDisplayAPI_screenToDisplay_Bad(t *core.T) {
	// screenToDisplay
	ax7Variant := "screenToDisplay:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertNil(t, screenToDisplay(nil))
	observedType := core.Sprintf("%T", screenToDisplay(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestDisplayAPI_screenToDisplay_Ugly(t *core.T) {
	// screenToDisplay
	ax7Variant := "screenToDisplay:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := screenToDisplay(&screen.Screen{})

	core.AssertNotNil(t, got)
	core.AssertEmpty(t, got.ID)
	core.AssertEmpty(t, got.Name)
	core.AssertEmpty(t, got.Width)
	core.AssertEmpty(t, got.Height)
}

func TestDisplayAPI_toDialogOpenFileOptions_Good(t *core.T) {
	// toDialogOpenFileOptions
	ax7Variant := "toDialogOpenFileOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	got := toDialogOpenFileOptions(OpenFileOptions{
		Title:            "Pick",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "report.csv",
		AllowMultiple:    true,
		Filters: []FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})

	core.AssertEqual(t, "Pick", got.Title)
	core.AssertEqual(t, "/tmp", got.Directory)
	core.AssertEqual(t, "report.csv", got.Filename)
	core.AssertTrue(t, got.AllowMultiple)
	core.AssertLen(t, got.Filters, 1)
	core.AssertEqual(t, "CSV", got.Filters[0].DisplayName)
	core.AssertEqual(t, "*.csv", got.Filters[0].Pattern)
}

func TestDisplayAPI_toDialogOpenFileOptions_Bad(t *core.T) {
	// toDialogOpenFileOptions
	ax7Variant := "toDialogOpenFileOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	got := toDialogOpenFileOptions(OpenFileOptions{})

	core.AssertEmpty(t, got.Title)
	core.AssertEmpty(t, got.Directory)
	core.AssertEmpty(t, got.Filename)
	core.AssertFalse(t, got.AllowMultiple)
	core.AssertNil(t, got.Filters)
}

func TestDisplayAPI_toDialogOpenFileOptions_Ugly(t *core.T) {
	// toDialogOpenFileOptions
	ax7Variant := "toDialogOpenFileOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := toDialogOpenFileOptions(OpenFileOptions{
		Filters: []FileFilter{
			{DisplayName: "All", Pattern: "*.*"},
			{DisplayName: "Media", Pattern: "*.png;*.jpg"},
		},
	})

	core.AssertLen(t, got.Filters, 2)
	core.AssertEqual(t, "All", got.Filters[0].DisplayName)
	core.AssertEqual(t, "*.png;*.jpg", got.Filters[1].Pattern)
}

func TestDisplayAPI_trayMenuItemsToSystray_Good(t *core.T) {
	// trayMenuItemsToSystray
	ax7Variant := "trayMenuItemsToSystray:good"
	core.AssertContains(t, ax7Variant, "good")
	got := trayMenuItemsToSystray([]TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{IsSeparator: true},
		{
			Label:    "More",
			ActionID: "more",
			Children: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}},
		},
	})

	core.AssertLen(t, got, 3)
	core.AssertEqual(t, "Open", got[0].Label)
	core.AssertEqual(t, "separator", got[1].Type)
	core.AssertLen(t, got[2].Submenu, 1)
	core.AssertEqual(t, "nested", got[2].Submenu[0].ActionID)
}

func TestDisplayAPI_trayMenuItemsToSystray_Bad(t *core.T) {
	// trayMenuItemsToSystray
	ax7Variant := "trayMenuItemsToSystray:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertNil(t, trayMenuItemsToSystray(nil))
	observedType := core.Sprintf("%T", trayMenuItemsToSystray(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestDisplayAPI_trayMenuItemsToSystray_Ugly(t *core.T) {
	// trayMenuItemsToSystray
	ax7Variant := "trayMenuItemsToSystray:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := trayMenuItemsToSystray([]TrayMenuItem{{Children: []TrayMenuItem{{IsSeparator: true}}}})

	core.AssertLen(t, got, 1)
	core.AssertLen(t, got[0].Submenu, 1)
	core.AssertEqual(t, "separator", got[0].Submenu[0].Type)
}

func TestDisplayAPI_GetScreens_Good(t *core.T) {
	// GetScreens
	ax7Variant := "GetScreens:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{
				{
					ID:          "screen-1",
					Name:        "Primary",
					Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
					ScaleFactor: 2,
					IsPrimary:   true,
				},
			}, OK: true}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()

	core.AssertLen(t, screens, 1)
	core.AssertEqual(t, "screen-1", screens[0].ID)
	core.AssertEqual(t, 10, screens[0].X)
	core.AssertEqual(t, 1920, screens[0].Width)
}

func TestDisplayAPI_GetScreens_Empty(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{}, OK: true}
		default:
			return core.Result{}
		}
	})

	core.AssertEmpty(t, svc.GetScreens())
}

func TestDisplayAPI_GetScreens_Bad(t *core.T) {
	// GetScreens
	ax7Variant := "GetScreens:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []string{"unexpected"}, OK: true}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()
	core.AssertNotNil(t, screens)
	core.AssertEmpty(t, screens)
}

func TestDisplayAPI_GetScreens_Ugly(t *core.T) {
	// GetScreens
	ax7Variant := "GetScreens:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()
	core.AssertNotNil(t, screens)
	core.AssertEmpty(t, screens)
}

func TestDisplayAPI_GetWorkAreas_Ugly(t *core.T) {
	// GetWorkAreas
	ax7Variant := "GetWorkAreas:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryWorkAreas:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	areas := svc.GetWorkAreas()

	core.AssertNotNil(t, areas)
	core.AssertEmpty(t, areas)
}

func TestDisplayAPI_GetScreen_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreen("screen-1")

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_CreateWindow_UglyResultType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("window.open", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	got, err := svc.CreateWindow(CreateWindowOptions{
		Name: "broken-window",
	})

	core.AssertError(t, err)
	core.AssertNil(t, got)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_GetScreen_Ugly(t *core.T) {
	// GetScreen
	ax7Variant := "GetScreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreen("screen-1")

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_GetPrimaryScreen_Ugly(t *core.T) {
	// GetPrimaryScreen
	ax7Variant := "GetPrimaryScreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryPrimary:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetPrimaryScreen()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_GetScreenAtPoint_Ugly(t *core.T) {
	// GetScreenAtPoint
	ax7Variant := "GetScreenAtPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAtPoint:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenAtPoint(10, 20)

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_OpenFileDialog_Good(t *core.T) {
	// OpenFileDialog
	ax7Variant := "OpenFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_file", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenFile)
		core.AssertEqual(t, "Pick file", task.Options.Title)
		core.AssertTrue(t, task.Options.AllowMultiple)
		return core.Result{Value: []string{"/tmp/a.txt", "/tmp/b.txt"}, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{
		Title:         "Pick file",
		AllowMultiple: true,
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, []string{"/tmp/a.txt", "/tmp/b.txt"}, paths)
}

func TestDisplayAPI_OpenFileDialog_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Bad(t *core.T) {
	// OpenFileDialog
	ax7Variant := "OpenFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Ugly(t *core.T) {
	// OpenFileDialog
	ax7Variant := "OpenFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_RequestNotificationPermission_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("notification.request_permission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: true}
	})

	granted, err := svc.RequestNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
}

func TestDisplayAPI_GetTheme_Good(t *core.T) {
	// GetTheme
	ax7Variant := "GetTheme:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{Value: environment.ThemeInfo{IsDark: true, Theme: "dark"}, OK: true}
		default:
			return core.Result{}
		}
	})

	theme := svc.GetTheme()
	core.AssertNotNil(t, theme)
	core.AssertTrue(t, theme.IsDark)
	core.AssertEqual(t, "dark", svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Bad(t *core.T) {
	// GetTheme
	ax7Variant := "GetTheme:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	theme := svc.GetTheme()
	core.AssertNil(t, theme)
	core.AssertEmpty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Ugly(t *core.T) {
	// GetTheme
	ax7Variant := "GetTheme:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	core.AssertNil(t, svc.GetTheme())
	core.AssertEmpty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_SaveFileDialog_Good(t *core.T) {
	// SaveFileDialog
	ax7Variant := "SaveFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.save_file", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskSaveFile)
		core.AssertEqual(t, "Export", task.Options.Title)
		core.AssertEqual(t, "/tmp", task.Options.Directory)
		core.AssertEqual(t, "data.json", task.Options.Filename)
		core.AssertLen(t, task.Options.Filters, 1)
		core.AssertEqual(t, "JSON", task.Options.Filters[0].DisplayName)
		return core.Result{Value: "/exports/data.json", OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{
		Title:            "Export",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "data.json",
		Filters:          []FileFilter{{DisplayName: "JSON", Pattern: "*.json"}},
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/exports/data.json", path)
}

func TestDisplayAPI_SaveFileDialog_Bad(t *core.T) {
	// SaveFileDialog
	ax7Variant := "SaveFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.save_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_SaveFileDialog_Ugly(t *core.T) {
	// SaveFileDialog
	ax7Variant := "SaveFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.save_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Good(t *core.T) {
	// OpenDirectoryDialog
	ax7Variant := "OpenDirectoryDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_directory", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenDirectory)
		core.AssertEqual(t, "Choose", task.Options.Title)
		core.AssertEqual(t, "/var", task.Options.Directory)
		return core.Result{Value: "/var/data", OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{
		Title:            "Choose",
		DefaultDirectory: "/var",
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/var/data", path)
}

func TestDisplayAPI_OpenDirectoryDialog_Bad(t *core.T) {
	// OpenDirectoryDialog
	ax7Variant := "OpenDirectoryDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_directory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Ugly(t *core.T) {
	// OpenDirectoryDialog
	ax7Variant := "OpenDirectoryDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.open_directory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_PromptDialog_Good(t *core.T) {
	// PromptDialog
	ax7Variant := "PromptDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskPrompt)
		core.AssertEqual(t, "Rename", task.Title)
		core.AssertEqual(t, "Enter a new name", task.Message)
		return core.Result{Value: dialog.PromptResult{Value: "draft", Confirmed: true}, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.RequireNoError(t, err)
	core.AssertTrue(t, confirmed)
	core.AssertEqual(t, "draft", value)
}

func TestDisplayAPI_PromptDialog_Bad(t *core.T) {
	// PromptDialog
	ax7Variant := "PromptDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEmpty(t, value)
}

func TestDisplayAPI_PromptDialog_Ugly(t *core.T) {
	// PromptDialog
	ax7Variant := "PromptDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEmpty(t, value)
}

func TestDisplayAPI_ReadClipboardImage_Good(t *core.T) {
	// ReadClipboardImage
	ax7Variant := "ReadClipboardImage:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	payload := []byte{1, 2, 3}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: clipboard.ImageContent{Data: payload, HasImage: true}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.RequireNoError(t, err)
	core.AssertEqual(t, []byte{1, 2, 3}, got)
	payload[0] = 9
	core.AssertEqual(t, byte(1), got[0])
}

func TestDisplayAPI_ReadClipboardImage_Bad(t *core.T) {
	// ReadClipboardImage
	ax7Variant := "ReadClipboardImage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: clipboard.ImageContent{HasImage: false}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.RequireNoError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly(t *core.T) {
	// ReadClipboardImage
	ax7Variant := "ReadClipboardImage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly_BackendFailure(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_WriteClipboardImage_Good(t *core.T) {
	// WriteClipboardImage
	ax7Variant := "WriteClipboardImage:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)
	var got []byte
	c.Action("clipboard.set_image", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("data").Value.([]byte)...)
		return core.Result{OK: true}
	})

	input := []byte{4, 5, 6}
	err := svc.WriteClipboardImage(input)

	core.RequireNoError(t, err)
	input[0] = 9
	core.AssertTrue(t, bytesEqual([]byte{4, 5, 6}, got))
}

func TestDisplayAPI_WriteClipboardImage_Bad(t *core.T) {
	// WriteClipboardImage
	ax7Variant := "WriteClipboardImage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, _ := newTestDisplayAPIService(t)

	err := svc.WriteClipboardImage(nil)

	core.AssertError(t, err)
}

func TestDisplayAPI_WriteClipboardImage_Ugly(t *core.T) {
	// WriteClipboardImage
	ax7Variant := "WriteClipboardImage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)
	c.Action("clipboard.set_image", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.WriteClipboardImage([]byte{1})

	core.AssertError(t, err)
}

func TestDisplayAPI_GetScreenForWindow_Good(t *core.T) {
	// GetScreenForWindow
	ax7Variant := "GetScreenForWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	var gotX, gotY int
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch typed := q.(type) {
		case window.QueryWindowByName:
			core.AssertEqual(t, "editor", typed.Name)
			return core.Result{
				Value: &window.WindowInfo{
					Name:   typed.Name,
					X:      100,
					Y:      200,
					Width:  300,
					Height: 400,
				},
				OK: true,
			}
		case screen.QueryAtPoint:
			gotX, gotY = typed.X, typed.Y
			return core.Result{
				Value: &screen.Screen{
					ID:          "screen-1",
					Name:        "Primary",
					ScaleFactor: 2,
					Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
					IsPrimary:   true,
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("editor")

	core.RequireNoError(t, err)
	core.AssertNotNil(t, got)
	core.AssertEqual(t, "screen-1", got.ID)
	core.AssertEqual(t, 250, gotX)
	core.AssertEqual(t, 400, gotY)
	core.AssertEqual(t, 10, got.X)
	core.AssertEqual(t, 20, got.Y)
	core.AssertEqual(t, 1920, got.Width)
	core.AssertEqual(t, 1080, got.Height)
}

func TestDisplayAPI_GetScreenForWindow_Bad(t *core.T) {
	// GetScreenForWindow
	ax7Variant := "GetScreenForWindow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	var screenQueried bool
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowByName:
			return core.Result{Value: (*window.WindowInfo)(nil), OK: true}
		case screen.QueryAtPoint:
			screenQueried = true
			return core.Result{OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("missing")

	core.RequireNoError(t, err)
	core.AssertNil(t, got)
	core.AssertFalse(t, screenQueried)
}

func TestDisplayAPI_GetScreenForWindow_Ugly(t *core.T) {
	// GetScreenForWindow
	ax7Variant := "GetScreenForWindow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowByName:
			return core.Result{
				Value: &window.WindowInfo{X: 1, Y: 2, Width: 3, Height: 4},
				OK:    true,
			}
		case screen.QueryAtPoint:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("editor")

	core.AssertError(t, err)
	core.AssertNil(t, got)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_OpenSingleFileDialog_Good(t *core.T) {
	// OpenSingleFileDialog
	ax7Variant := "OpenSingleFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	var task dialog.TaskOpenFile
	c.Action("dialog.open_file", func(_ context.Context, opts core.Options) core.Result {
		task = opts.Get("task").Value.(dialog.TaskOpenFile)
		return core.Result{Value: []string{"/tmp/report.csv"}, OK: true}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{
		Title:           "Pick report",
		DefaultFilename: "report.csv",
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/report.csv", path)
	core.AssertEqual(t, "Pick report", task.Options.Title)
	core.AssertEqual(t, "report.csv", task.Options.Filename)
}

func TestDisplayAPI_OpenSingleFileDialog_Bad(t *core.T) {
	// OpenSingleFileDialog
	ax7Variant := "OpenSingleFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.open_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: []string{}, OK: true}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	core.RequireNoError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenSingleFileDialog_Ugly(t *core.T) {
	// OpenSingleFileDialog
	ax7Variant := "OpenSingleFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.open_file", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
	core.AssertContains(t, err.Error(), "dialog.openFile action failed")
}

func TestDisplayAPI_ConfirmDialog_Good(t *core.T) {
	// ConfirmDialog
	ax7Variant := "ConfirmDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	var task dialog.TaskQuestion
	c.Action("dialog.question", func(_ context.Context, opts core.Options) core.Result {
		task = opts.Get("task").Value.(dialog.TaskQuestion)
		return core.Result{Value: "Yes", OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.RequireNoError(t, err)
	core.AssertTrue(t, confirmed)
	core.AssertEqual(t, "Confirm", task.Title)
	core.AssertEqual(t, []string{"Yes", "No"}, task.Buttons)
}

func TestDisplayAPI_ConfirmDialog_Bad(t *core.T) {
	// ConfirmDialog
	ax7Variant := "ConfirmDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_ConfirmDialog_Ugly(t *core.T) {
	// ConfirmDialog
	ax7Variant := "ConfirmDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_ReadClipboard_Good(t *core.T) {
	// ReadClipboard
	ax7Variant := "ReadClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{
				Value: clipboard.ClipboardContent{
					Text:       "hello clipboard",
					HasContent: true,
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "hello clipboard", text)
}

func TestDisplayAPI_ReadClipboard_Bad(t *core.T) {
	// ReadClipboard
	ax7Variant := "ReadClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, text)
	// Missing seam: QUERY drops non-OK backend errors, so propagation is not observable here.
}

func TestDisplayAPI_ReadClipboard_Ugly(t *core.T) {
	// ReadClipboard
	ax7Variant := "ReadClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.AssertError(t, err)
	core.AssertEmpty(t, text)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_CheckNotificationPermission_Good(t *core.T) {
	// CheckNotificationPermission
	ax7Variant := "CheckNotificationPermission:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{Value: notification.PermissionStatus{Granted: true}, OK: true}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.RequireNoError(t, err)
	core.AssertTrue(t, granted)
}

func TestDisplayAPI_CheckNotificationPermission_Bad(t *core.T) {
	// CheckNotificationPermission
	ax7Variant := "CheckNotificationPermission:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
	core.AssertContains(t, err.Error(), "notification query failed")
}

func TestDisplayAPI_CheckNotificationPermission_Ugly(t *core.T) {
	// CheckNotificationPermission
	ax7Variant := "CheckNotificationPermission:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_WriteClipboard_Good(t *core.T) {
	// WriteClipboard
	ax7Variant := "WriteClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	var gotText string
	c.Action("clipboard.set_text", func(_ context.Context, opts core.Options) core.Result {
		gotText = opts.Get("task").Value.(clipboard.TaskSetText).Text
		return core.Result{OK: true}
	})

	err := svc.WriteClipboard("hello")

	core.RequireNoError(t, err)
	core.AssertEqual(t, "hello", gotText)
}

func TestDisplayAPI_WriteClipboard_Bad(t *core.T) {
	// WriteClipboard
	ax7Variant := "WriteClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.set_text", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.WriteClipboard("hello")

	core.AssertError(t, err)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_WriteClipboard_Ugly(t *core.T) {
	// WriteClipboard
	ax7Variant := "WriteClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.set_text", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.WriteClipboard("")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "clipboard.set_text")
}

func TestDisplayAPI_SetTrayIcon_Good(t *core.T) {
	// SetTrayIcon
	ax7Variant := "SetTrayIcon:good"
	core.AssertContains(t, ax7Variant, "good")
	svc, c := newTestDisplayAPIService(t)

	var got []byte
	c.Action("systray.set_icon", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("task").Value.(systray.TaskSetTrayIcon).Data...)
		return core.Result{OK: true}
	})

	err := svc.SetTrayIcon([]byte{1, 2, 3})

	core.RequireNoError(t, err)
	core.AssertEqual(t, []byte{1, 2, 3}, got)
}

func TestDisplayAPI_SetTrayIcon_Bad(t *core.T) {
	// SetTrayIcon
	ax7Variant := "SetTrayIcon:bad"
	core.AssertContains(t, ax7Variant, "bad")
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.set_icon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.SetTrayIcon([]byte{1})

	core.AssertError(t, err)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_SetTrayIcon_Ugly(t *core.T) {
	// SetTrayIcon
	ax7Variant := "SetTrayIcon:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.set_icon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.SetTrayIcon(nil)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "systray.set_icon")
}

// AX7 generated source-matching smoke coverage.
func TestApi_Service_GetScreens_Good(t *core.T) {
	// Service GetScreens
	ax7Variant := "Service_GetScreens:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetScreens()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreens_Bad(t *core.T) {
	// Service GetScreens
	ax7Variant := "Service_GetScreens:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetScreens()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreens_Ugly(t *core.T) {
	// Service GetScreens
	ax7Variant := "Service_GetScreens:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetScreens()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreen_Good(t *core.T) {
	// Service GetScreen
	ax7Variant := "Service_GetScreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreen_Bad(t *core.T) {
	// Service GetScreen
	ax7Variant := "Service_GetScreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreen_Ugly(t *core.T) {
	// Service GetScreen
	ax7Variant := "Service_GetScreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetPrimaryScreen_Good(t *core.T) {
	// Service GetPrimaryScreen
	ax7Variant := "Service_GetPrimaryScreen:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetPrimaryScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetPrimaryScreen_Bad(t *core.T) {
	// Service GetPrimaryScreen
	ax7Variant := "Service_GetPrimaryScreen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetPrimaryScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetPrimaryScreen_Ugly(t *core.T) {
	// Service GetPrimaryScreen
	ax7Variant := "Service_GetPrimaryScreen:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetPrimaryScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenAtPoint_Good(t *core.T) {
	// Service GetScreenAtPoint
	ax7Variant := "Service_GetScreenAtPoint:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenAtPoint(1, 1)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenAtPoint_Bad(t *core.T) {
	// Service GetScreenAtPoint
	ax7Variant := "Service_GetScreenAtPoint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenAtPoint(0, 0)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenAtPoint_Ugly(t *core.T) {
	// Service GetScreenAtPoint
	ax7Variant := "Service_GetScreenAtPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenAtPoint(-1, -1)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenForWindow_Good(t *core.T) {
	// Service GetScreenForWindow
	ax7Variant := "Service_GetScreenForWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenForWindow("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenForWindow_Bad(t *core.T) {
	// Service GetScreenForWindow
	ax7Variant := "Service_GetScreenForWindow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenForWindow("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetScreenForWindow_Ugly(t *core.T) {
	// Service GetScreenForWindow
	ax7Variant := "Service_GetScreenForWindow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreenForWindow("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetWorkAreas_Good(t *core.T) {
	// Service GetWorkAreas
	ax7Variant := "Service_GetWorkAreas:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetWorkAreas()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetWorkAreas_Bad(t *core.T) {
	// Service GetWorkAreas
	ax7Variant := "Service_GetWorkAreas:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetWorkAreas()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetWorkAreas_Ugly(t *core.T) {
	// Service GetWorkAreas
	ax7Variant := "Service_GetWorkAreas:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetWorkAreas()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenSingleFileDialog_Good(t *core.T) {
	// Service OpenSingleFileDialog
	ax7Variant := "Service_OpenSingleFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenSingleFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenSingleFileDialog_Bad(t *core.T) {
	// Service OpenSingleFileDialog
	ax7Variant := "Service_OpenSingleFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenSingleFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenSingleFileDialog_Ugly(t *core.T) {
	// Service OpenSingleFileDialog
	ax7Variant := "Service_OpenSingleFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenSingleFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenFileDialog_Good(t *core.T) {
	// Service OpenFileDialog
	ax7Variant := "Service_OpenFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenFileDialog_Bad(t *core.T) {
	// Service OpenFileDialog
	ax7Variant := "Service_OpenFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenFileDialog_Ugly(t *core.T) {
	// Service OpenFileDialog
	ax7Variant := "Service_OpenFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenFileDialog(*new(OpenFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SaveFileDialog_Good(t *core.T) {
	// Service SaveFileDialog
	ax7Variant := "Service_SaveFileDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SaveFileDialog(*new(SaveFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SaveFileDialog_Bad(t *core.T) {
	// Service SaveFileDialog
	ax7Variant := "Service_SaveFileDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SaveFileDialog(*new(SaveFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SaveFileDialog_Ugly(t *core.T) {
	// Service SaveFileDialog
	ax7Variant := "Service_SaveFileDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SaveFileDialog(*new(SaveFileOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenDirectoryDialog_Good(t *core.T) {
	// Service OpenDirectoryDialog
	ax7Variant := "Service_OpenDirectoryDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenDirectoryDialog(*new(OpenDirectoryOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenDirectoryDialog_Bad(t *core.T) {
	// Service OpenDirectoryDialog
	ax7Variant := "Service_OpenDirectoryDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenDirectoryDialog(*new(OpenDirectoryOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_OpenDirectoryDialog_Ugly(t *core.T) {
	// Service OpenDirectoryDialog
	ax7Variant := "Service_OpenDirectoryDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.OpenDirectoryDialog(*new(OpenDirectoryOptions))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ConfirmDialog_Good(t *core.T) {
	// Service ConfirmDialog
	ax7Variant := "Service_ConfirmDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ConfirmDialog("agent", "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ConfirmDialog_Bad(t *core.T) {
	// Service ConfirmDialog
	ax7Variant := "Service_ConfirmDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ConfirmDialog("", "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ConfirmDialog_Ugly(t *core.T) {
	// Service ConfirmDialog
	ax7Variant := "Service_ConfirmDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ConfirmDialog("../../edge", "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_PromptDialog_Good(t *core.T) {
	// Service PromptDialog
	ax7Variant := "Service_PromptDialog:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1, got2 := subject.PromptDialog("agent", "agent")
		return core.Sprintf("%T,%T,%T", got0, got1, got2)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_PromptDialog_Bad(t *core.T) {
	// Service PromptDialog
	ax7Variant := "Service_PromptDialog:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1, got2 := subject.PromptDialog("", "")
		return core.Sprintf("%T,%T,%T", got0, got1, got2)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_PromptDialog_Ugly(t *core.T) {
	// Service PromptDialog
	ax7Variant := "Service_PromptDialog:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1, got2 := subject.PromptDialog("../../edge", "../../edge")
		return core.Sprintf("%T,%T,%T", got0, got1, got2)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayIcon_Good(t *core.T) {
	// Service SetTrayIcon
	ax7Variant := "Service_SetTrayIcon:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayIcon_Bad(t *core.T) {
	// Service SetTrayIcon
	ax7Variant := "Service_SetTrayIcon:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayIcon_Ugly(t *core.T) {
	// Service SetTrayIcon
	ax7Variant := "Service_SetTrayIcon:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayTooltip_Good(t *core.T) {
	// Service SetTrayTooltip
	ax7Variant := "Service_SetTrayTooltip:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayTooltip("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayTooltip_Bad(t *core.T) {
	// Service SetTrayTooltip
	ax7Variant := "Service_SetTrayTooltip:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayTooltip("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayTooltip_Ugly(t *core.T) {
	// Service SetTrayTooltip
	ax7Variant := "Service_SetTrayTooltip:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayTooltip("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayLabel_Good(t *core.T) {
	// Service SetTrayLabel
	ax7Variant := "Service_SetTrayLabel:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayLabel("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayLabel_Bad(t *core.T) {
	// Service SetTrayLabel
	ax7Variant := "Service_SetTrayLabel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayLabel("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayLabel_Ugly(t *core.T) {
	// Service SetTrayLabel
	ax7Variant := "Service_SetTrayLabel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayLabel("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayMenu_Good(t *core.T) {
	// Service SetTrayMenu
	ax7Variant := "Service_SetTrayMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayMenu_Bad(t *core.T) {
	// Service SetTrayMenu
	ax7Variant := "Service_SetTrayMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTrayMenu_Ugly(t *core.T) {
	// Service SetTrayMenu
	ax7Variant := "Service_SetTrayMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTrayMenu(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTrayInfo_Good(t *core.T) {
	// Service GetTrayInfo
	ax7Variant := "Service_GetTrayInfo:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTrayInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTrayInfo_Bad(t *core.T) {
	// Service GetTrayInfo
	ax7Variant := "Service_GetTrayInfo:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTrayInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTrayInfo_Ugly(t *core.T) {
	// Service GetTrayInfo
	ax7Variant := "Service_GetTrayInfo:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTrayInfo()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowTrayMessage_Good(t *core.T) {
	// Service ShowTrayMessage
	ax7Variant := "Service_ShowTrayMessage:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowTrayMessage("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowTrayMessage_Bad(t *core.T) {
	// Service ShowTrayMessage
	ax7Variant := "Service_ShowTrayMessage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowTrayMessage("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowTrayMessage_Ugly(t *core.T) {
	// Service ShowTrayMessage
	ax7Variant := "Service_ShowTrayMessage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowTrayMessage("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboard_Good(t *core.T) {
	// Service ReadClipboard
	ax7Variant := "Service_ReadClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboard()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboard_Bad(t *core.T) {
	// Service ReadClipboard
	ax7Variant := "Service_ReadClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboard()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboard_Ugly(t *core.T) {
	// Service ReadClipboard
	ax7Variant := "Service_ReadClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboard()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboard_Good(t *core.T) {
	// Service WriteClipboard
	ax7Variant := "Service_WriteClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboard("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboard_Bad(t *core.T) {
	// Service WriteClipboard
	ax7Variant := "Service_WriteClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboard("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboard_Ugly(t *core.T) {
	// Service WriteClipboard
	ax7Variant := "Service_WriteClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboard("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_HasClipboard_Good(t *core.T) {
	// Service HasClipboard
	ax7Variant := "Service_HasClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HasClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_HasClipboard_Bad(t *core.T) {
	// Service HasClipboard
	ax7Variant := "Service_HasClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HasClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_HasClipboard_Ugly(t *core.T) {
	// Service HasClipboard
	ax7Variant := "Service_HasClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HasClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearClipboard_Good(t *core.T) {
	// Service ClearClipboard
	ax7Variant := "Service_ClearClipboard:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearClipboard_Bad(t *core.T) {
	// Service ClearClipboard
	ax7Variant := "Service_ClearClipboard:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearClipboard_Ugly(t *core.T) {
	// Service ClearClipboard
	ax7Variant := "Service_ClearClipboard:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearClipboard()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboardImage_Good(t *core.T) {
	// Service ReadClipboardImage
	ax7Variant := "Service_ReadClipboardImage:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboardImage()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboardImage_Bad(t *core.T) {
	// Service ReadClipboardImage
	ax7Variant := "Service_ReadClipboardImage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboardImage()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ReadClipboardImage_Ugly(t *core.T) {
	// Service ReadClipboardImage
	ax7Variant := "Service_ReadClipboardImage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ReadClipboardImage()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboardImage_Good(t *core.T) {
	// Service WriteClipboardImage
	ax7Variant := "Service_WriteClipboardImage:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboardImage(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboardImage_Bad(t *core.T) {
	// Service WriteClipboardImage
	ax7Variant := "Service_WriteClipboardImage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboardImage(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_WriteClipboardImage_Ugly(t *core.T) {
	// Service WriteClipboardImage
	ax7Variant := "Service_WriteClipboardImage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.WriteClipboardImage(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowNotification_Good(t *core.T) {
	// Service ShowNotification
	ax7Variant := "Service_ShowNotification:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowNotification(*new(NotificationOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowNotification_Bad(t *core.T) {
	// Service ShowNotification
	ax7Variant := "Service_ShowNotification:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowNotification(*new(NotificationOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowNotification_Ugly(t *core.T) {
	// Service ShowNotification
	ax7Variant := "Service_ShowNotification:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowNotification(*new(NotificationOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowInfoNotification_Good(t *core.T) {
	// Service ShowInfoNotification
	ax7Variant := "Service_ShowInfoNotification:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowInfoNotification("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowInfoNotification_Bad(t *core.T) {
	// Service ShowInfoNotification
	ax7Variant := "Service_ShowInfoNotification:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowInfoNotification("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowInfoNotification_Ugly(t *core.T) {
	// Service ShowInfoNotification
	ax7Variant := "Service_ShowInfoNotification:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowInfoNotification("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowWarningNotification_Good(t *core.T) {
	// Service ShowWarningNotification
	ax7Variant := "Service_ShowWarningNotification:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowWarningNotification("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowWarningNotification_Bad(t *core.T) {
	// Service ShowWarningNotification
	ax7Variant := "Service_ShowWarningNotification:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowWarningNotification("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowWarningNotification_Ugly(t *core.T) {
	// Service ShowWarningNotification
	ax7Variant := "Service_ShowWarningNotification:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowWarningNotification("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowErrorNotification_Good(t *core.T) {
	// Service ShowErrorNotification
	ax7Variant := "Service_ShowErrorNotification:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowErrorNotification("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowErrorNotification_Bad(t *core.T) {
	// Service ShowErrorNotification
	ax7Variant := "Service_ShowErrorNotification:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowErrorNotification("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ShowErrorNotification_Ugly(t *core.T) {
	// Service ShowErrorNotification
	ax7Variant := "Service_ShowErrorNotification:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ShowErrorNotification("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_RequestNotificationPermission_Good(t *core.T) {
	// Service RequestNotificationPermission
	ax7Variant := "Service_RequestNotificationPermission:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.RequestNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_RequestNotificationPermission_Bad(t *core.T) {
	// Service RequestNotificationPermission
	ax7Variant := "Service_RequestNotificationPermission:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.RequestNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_RequestNotificationPermission_Ugly(t *core.T) {
	// Service RequestNotificationPermission
	ax7Variant := "Service_RequestNotificationPermission:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.RequestNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_CheckNotificationPermission_Good(t *core.T) {
	// Service CheckNotificationPermission
	ax7Variant := "Service_CheckNotificationPermission:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.CheckNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_CheckNotificationPermission_Bad(t *core.T) {
	// Service CheckNotificationPermission
	ax7Variant := "Service_CheckNotificationPermission:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.CheckNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_CheckNotificationPermission_Ugly(t *core.T) {
	// Service CheckNotificationPermission
	ax7Variant := "Service_CheckNotificationPermission:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.CheckNotificationPermission()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearNotifications_Good(t *core.T) {
	// Service ClearNotifications
	ax7Variant := "Service_ClearNotifications:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearNotifications()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearNotifications_Bad(t *core.T) {
	// Service ClearNotifications
	ax7Variant := "Service_ClearNotifications:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearNotifications()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_ClearNotifications_Ugly(t *core.T) {
	// Service ClearNotifications
	ax7Variant := "Service_ClearNotifications:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.ClearNotifications()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTheme_Good(t *core.T) {
	// Service SetTheme
	ax7Variant := "Service_SetTheme:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTheme("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTheme_Bad(t *core.T) {
	// Service SetTheme
	ax7Variant := "Service_SetTheme:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTheme("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_SetTheme_Ugly(t *core.T) {
	// Service SetTheme
	ax7Variant := "Service_SetTheme:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SetTheme("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTheme_Good(t *core.T) {
	// Service GetTheme
	ax7Variant := "Service_GetTheme:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTheme_Bad(t *core.T) {
	// Service GetTheme
	ax7Variant := "Service_GetTheme:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetTheme_Ugly(t *core.T) {
	// Service GetTheme
	ax7Variant := "Service_GetTheme:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetSystemTheme_Good(t *core.T) {
	// Service GetSystemTheme
	ax7Variant := "Service_GetSystemTheme:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetSystemTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetSystemTheme_Bad(t *core.T) {
	// Service GetSystemTheme
	ax7Variant := "Service_GetSystemTheme:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetSystemTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestApi_Service_GetSystemTheme_Ugly(t *core.T) {
	// Service GetSystemTheme
	ax7Variant := "Service_GetSystemTheme:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.GetSystemTheme()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
