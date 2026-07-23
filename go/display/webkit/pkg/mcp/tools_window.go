// pkg/mcp/tools_window.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- window_list ---

type WindowListInput struct{}
type WindowListOutput struct {
	Windows []window.WindowInfo `json:"windows"`
}

func (s *Subsystem) windowList(_ context.Context, _ *mcp.CallToolRequest, _ WindowListInput) (*mcp.CallToolResult, WindowListOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowListOutput{}, e
		}
		return nil, WindowListOutput{}, nil
	}
	windows, ok := r.Value.([]window.WindowInfo)
	if !ok {
		return nil, WindowListOutput{}, core.E("mcp.windowList", "unexpected result type", nil)
	}
	return nil, WindowListOutput{Windows: windows}, nil
}

// --- window_get ---

type WindowGetInput struct {
	Name string `json:"name"`
}
type WindowGetOutput struct {
	Window *window.WindowInfo `json:"window"`
}

func (s *Subsystem) windowGet(_ context.Context, _ *mcp.CallToolRequest, input WindowGetInput) (*mcp.CallToolResult, WindowGetOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowGetOutput{}, e
		}
		return nil, WindowGetOutput{}, nil
	}
	info, ok := r.Value.(*window.WindowInfo)
	if !ok {
		return nil, WindowGetOutput{}, core.E("mcp.windowGet", "unexpected result type", nil)
	}
	return nil, WindowGetOutput{Window: info}, nil
}

// --- window_focused ---

type WindowFocusedInput struct{}
type WindowFocusedOutput struct {
	Window string `json:"window"`
}

func (s *Subsystem) windowFocused(_ context.Context, _ *mcp.CallToolRequest, _ WindowFocusedInput) (*mcp.CallToolResult, WindowFocusedOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowFocusedOutput{}, e
		}
		return nil, WindowFocusedOutput{}, nil
	}
	windows, ok := r.Value.([]window.WindowInfo)
	if !ok {
		return nil, WindowFocusedOutput{}, core.E("mcp.windowFocused", "unexpected result type", nil)
	}
	for _, w := range windows {
		if w.Focused {
			return nil, WindowFocusedOutput{Window: w.Name}, nil
		}
	}
	return nil, WindowFocusedOutput{}, nil
}

// --- window_create ---

type WindowCreateInput struct {
	Name   string `json:"name"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
type WindowCreateOutput struct {
	Window window.WindowInfo `json:"window"`
}

func (s *Subsystem) windowCreate(_ context.Context, _ *mcp.CallToolRequest, input WindowCreateInput) (*mcp.CallToolResult, WindowCreateOutput, resultFailure) {
	r := s.core.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   input.Name,
				Title:  input.Title,
				URL:    input.URL,
				Width:  input.Width,
				Height: input.Height,
				X:      input.X,
				Y:      input.Y,
			},
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowCreateOutput{}, e
		}
		return nil, WindowCreateOutput{}, core.E("mcp.windowCreate", "window.open failed", nil)
	}
	info, ok := r.Value.(window.WindowInfo)
	if !ok {
		return nil, WindowCreateOutput{}, core.E("mcp.windowCreate", "unexpected result type", nil)
	}
	return nil, WindowCreateOutput{Window: info}, nil
}

// --- window_close ---

type WindowCloseInput struct {
	Name string `json:"name"`
}
type WindowCloseOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowClose(_ context.Context, _ *mcp.CallToolRequest, input WindowCloseInput) (*mcp.CallToolResult, WindowCloseOutput, resultFailure) {
	r := s.core.Action("window.close").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskCloseWindow{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowCloseOutput{}, e
		}
		return nil, WindowCloseOutput{}, nil
	}
	return nil, WindowCloseOutput{Success: true}, nil
}

// --- window_position ---

type WindowPositionInput struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}
type WindowPositionOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowPosition(_ context.Context, _ *mcp.CallToolRequest, input WindowPositionInput) (*mcp.CallToolResult, WindowPositionOutput, resultFailure) {
	r := s.core.Action("window.set_position").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetPosition{Name: input.Name, X: input.X, Y: input.Y}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowPositionOutput{}, e
		}
		return nil, WindowPositionOutput{}, nil
	}
	return nil, WindowPositionOutput{Success: true}, nil
}

// --- window_size ---

type WindowSizeInput struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type WindowSizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowSize(_ context.Context, _ *mcp.CallToolRequest, input WindowSizeInput) (*mcp.CallToolResult, WindowSizeOutput, resultFailure) {
	r := s.core.Action("window.set_size").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetSize{Name: input.Name, Width: input.Width, Height: input.Height}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowSizeOutput{}, e
		}
		return nil, WindowSizeOutput{}, nil
	}
	return nil, WindowSizeOutput{Success: true}, nil
}

// --- window_bounds ---

type WindowBoundsInput struct {
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type WindowBoundsOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowBounds(_ context.Context, _ *mcp.CallToolRequest, input WindowBoundsInput) (*mcp.CallToolResult, WindowBoundsOutput, resultFailure) {
	r := s.core.Action("window.set_bounds").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetBounds{
			Name: input.Name, X: input.X, Y: input.Y, Width: input.Width, Height: input.Height,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowBoundsOutput{}, e
		}
		return nil, WindowBoundsOutput{}, nil
	}
	return nil, WindowBoundsOutput{Success: true}, nil
}

// --- window_maximize ---

type WindowMaximizeInput struct {
	Name string `json:"name"`
}
type WindowMaximizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowMaximize(_ context.Context, _ *mcp.CallToolRequest, input WindowMaximizeInput) (*mcp.CallToolResult, WindowMaximizeOutput, resultFailure) {
	r := s.core.Action("window.maximise").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskMaximise{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowMaximizeOutput{}, e
		}
		return nil, WindowMaximizeOutput{}, nil
	}
	return nil, WindowMaximizeOutput{Success: true}, nil
}

// --- window_minimize ---

type WindowMinimizeInput struct {
	Name string `json:"name"`
}
type WindowMinimizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowMinimize(_ context.Context, _ *mcp.CallToolRequest, input WindowMinimizeInput) (*mcp.CallToolResult, WindowMinimizeOutput, resultFailure) {
	r := s.core.Action("window.minimise").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskMinimise{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowMinimizeOutput{}, e
		}
		return nil, WindowMinimizeOutput{}, nil
	}
	return nil, WindowMinimizeOutput{Success: true}, nil
}

// --- window_restore ---

type WindowRestoreInput struct {
	Name string `json:"name"`
}
type WindowRestoreOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowRestore(_ context.Context, _ *mcp.CallToolRequest, input WindowRestoreInput) (*mcp.CallToolResult, WindowRestoreOutput, resultFailure) {
	r := s.core.Action("window.restore").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRestore{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowRestoreOutput{}, e
		}
		return nil, WindowRestoreOutput{}, nil
	}
	return nil, WindowRestoreOutput{Success: true}, nil
}

// --- window_focus ---

type WindowFocusInput struct {
	Name string `json:"name"`
}
type WindowFocusOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowFocus(_ context.Context, _ *mcp.CallToolRequest, input WindowFocusInput) (*mcp.CallToolResult, WindowFocusOutput, resultFailure) {
	r := s.core.Action("window.focus").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFocus{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowFocusOutput{}, e
		}
		return nil, WindowFocusOutput{}, nil
	}
	return nil, WindowFocusOutput{Success: true}, nil
}

// --- focus_set ---

type FocusSetInput struct {
	Name string `json:"name"`
}
type FocusSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) focusSet(ctx context.Context, req *mcp.CallToolRequest, input FocusSetInput) (*mcp.CallToolResult, FocusSetOutput, resultFailure) {
	result, output, err := s.windowFocus(ctx, req, WindowFocusInput{Name: input.Name})
	if err != nil {
		return nil, FocusSetOutput{}, err
	}
	if result != nil {
		return result, FocusSetOutput{}, nil
	}
	return nil, FocusSetOutput{Success: output.Success}, nil
}

// --- window_title ---

type WindowTitleInput struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
type WindowTitleOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowTitle(_ context.Context, _ *mcp.CallToolRequest, input WindowTitleInput) (*mcp.CallToolResult, WindowTitleOutput, resultFailure) {
	r := s.core.Action("window.set_title").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetTitle{Name: input.Name, Title: input.Title}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowTitleOutput{}, e
		}
		return nil, WindowTitleOutput{}, nil
	}
	return nil, WindowTitleOutput{Success: true}, nil
}

// --- window_title_set ---

type WindowTitleSetInput struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
type WindowTitleSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowTitleSet(ctx context.Context, req *mcp.CallToolRequest, input WindowTitleSetInput) (*mcp.CallToolResult, WindowTitleSetOutput, resultFailure) {
	result, output, err := s.windowTitle(ctx, req, WindowTitleInput{Name: input.Name, Title: input.Title})
	if err != nil {
		return nil, WindowTitleSetOutput{}, err
	}
	if result != nil {
		return result, WindowTitleSetOutput{}, nil
	}
	return nil, WindowTitleSetOutput{Success: output.Success}, nil
}

// --- window_title_get ---

type WindowTitleGetInput struct {
	Name string `json:"name"`
}
type WindowTitleGetOutput struct {
	Title string `json:"title"`
}

func (s *Subsystem) windowTitleGet(_ context.Context, _ *mcp.CallToolRequest, input WindowTitleGetInput) (*mcp.CallToolResult, WindowTitleGetOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if !r.OK {
		return nil, WindowTitleGetOutput{}, nil
	}
	info, _ := r.Value.(*window.WindowInfo)
	if info == nil {
		return nil, WindowTitleGetOutput{}, nil
	}
	return nil, WindowTitleGetOutput{Title: info.Title}, nil
}

// --- window_visibility ---

type WindowVisibilityInput struct {
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
}
type WindowVisibilityOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowVisibility(_ context.Context, _ *mcp.CallToolRequest, input WindowVisibilityInput) (*mcp.CallToolResult, WindowVisibilityOutput, resultFailure) {
	r := s.core.Action("window.set_visibility").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetVisibility{Name: input.Name, Visible: input.Visible}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowVisibilityOutput{}, e
		}
		return nil, WindowVisibilityOutput{}, nil
	}
	return nil, WindowVisibilityOutput{Success: true}, nil
}

// --- window_always_on_top ---

type WindowAlwaysOnTopInput struct {
	Name        string `json:"name"`
	AlwaysOnTop bool   `json:"alwaysOnTop"`
}
type WindowAlwaysOnTopOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowAlwaysOnTop(_ context.Context, _ *mcp.CallToolRequest, input WindowAlwaysOnTopInput) (*mcp.CallToolResult, WindowAlwaysOnTopOutput, resultFailure) {
	r := s.core.Action("window.set_always_on_top").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetAlwaysOnTop{Name: input.Name, AlwaysOnTop: input.AlwaysOnTop}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowAlwaysOnTopOutput{}, e
		}
		return nil, WindowAlwaysOnTopOutput{}, nil
	}
	return nil, WindowAlwaysOnTopOutput{Success: true}, nil
}

// --- window_opacity ---

type WindowOpacityInput struct {
	Name    string  `json:"name"`
	Opacity float64 `json:"opacity"`
}
type WindowOpacityOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowOpacity(_ context.Context, _ *mcp.CallToolRequest, input WindowOpacityInput) (*mcp.CallToolResult, WindowOpacityOutput, resultFailure) {
	r := s.core.Action("window.set_opacity").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetOpacity{Name: input.Name, Opacity: input.Opacity}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowOpacityOutput{}, e
		}
		return nil, WindowOpacityOutput{}, nil
	}
	return nil, WindowOpacityOutput{Success: true}, nil
}

// --- window_background_colour ---

type WindowBackgroundColourInput struct {
	Name  string `json:"name"`
	Red   uint8  `json:"red"`
	Green uint8  `json:"green"`
	Blue  uint8  `json:"blue"`
	Alpha uint8  `json:"alpha"`
}
type WindowBackgroundColourOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowBackgroundColour(_ context.Context, _ *mcp.CallToolRequest, input WindowBackgroundColourInput) (*mcp.CallToolResult, WindowBackgroundColourOutput, resultFailure) {
	r := s.core.Action("window.set_background_colour").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetBackgroundColour{
			Name: input.Name, Red: input.Red, Green: input.Green, Blue: input.Blue, Alpha: input.Alpha,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowBackgroundColourOutput{}, e
		}
		return nil, WindowBackgroundColourOutput{}, nil
	}
	return nil, WindowBackgroundColourOutput{Success: true}, nil
}

// --- window_fullscreen ---

type WindowFullscreenInput struct {
	Name       string `json:"name"`
	Fullscreen bool   `json:"fullscreen"`
}
type WindowFullscreenOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowFullscreen(_ context.Context, _ *mcp.CallToolRequest, input WindowFullscreenInput) (*mcp.CallToolResult, WindowFullscreenOutput, resultFailure) {
	r := s.core.Action("window.fullscreen").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFullscreen{Name: input.Name, Fullscreen: input.Fullscreen}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowFullscreenOutput{}, e
		}
		return nil, WindowFullscreenOutput{}, nil
	}
	return nil, WindowFullscreenOutput{Success: true}, nil
}

// --- window_zoom_set ---

type WindowZoomSetInput struct {
	Name          string  `json:"name"`
	Magnification float64 `json:"magnification"`
}
type WindowZoomSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowZoomSet(_ context.Context, _ *mcp.CallToolRequest, input WindowZoomSetInput) (*mcp.CallToolResult, WindowZoomSetOutput, resultFailure) {
	r := s.core.Action("window.set_zoom").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetZoom{Name: input.Name, Magnification: input.Magnification}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowZoomSetOutput{}, e
		}
		return nil, WindowZoomSetOutput{}, nil
	}
	return nil, WindowZoomSetOutput{Success: true}, nil
}

// --- window_zoom_in ---

type WindowZoomInInput struct {
	Name string `json:"name"`
}
type WindowZoomInOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowZoomIn(_ context.Context, _ *mcp.CallToolRequest, input WindowZoomInInput) (*mcp.CallToolResult, WindowZoomInOutput, resultFailure) {
	r := s.core.Action("window.zoom_in").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskZoomIn{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowZoomInOutput{}, e
		}
		return nil, WindowZoomInOutput{}, nil
	}
	return nil, WindowZoomInOutput{Success: true}, nil
}

// --- window_zoom_out ---

type WindowZoomOutInput struct {
	Name string `json:"name"`
}
type WindowZoomOutOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowZoomOut(_ context.Context, _ *mcp.CallToolRequest, input WindowZoomOutInput) (*mcp.CallToolResult, WindowZoomOutOutput, resultFailure) {
	r := s.core.Action("window.zoom_out").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskZoomOut{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowZoomOutOutput{}, e
		}
		return nil, WindowZoomOutOutput{}, nil
	}
	return nil, WindowZoomOutOutput{Success: true}, nil
}

// --- window_zoom_reset ---

type WindowZoomResetInput struct {
	Name string `json:"name"`
}
type WindowZoomResetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowZoomReset(_ context.Context, _ *mcp.CallToolRequest, input WindowZoomResetInput) (*mcp.CallToolResult, WindowZoomResetOutput, resultFailure) {
	r := s.core.Action("window.zoom_reset").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskZoomReset{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowZoomResetOutput{}, e
		}
		return nil, WindowZoomResetOutput{}, nil
	}
	return nil, WindowZoomResetOutput{Success: true}, nil
}

// --- window_url_set ---

type WindowURLSetInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type WindowURLSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowURLSet(_ context.Context, _ *mcp.CallToolRequest, input WindowURLSetInput) (*mcp.CallToolResult, WindowURLSetOutput, resultFailure) {
	r := s.core.Action("window.set_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetURL{Name: input.Name, URL: input.URL}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowURLSetOutput{}, e
		}
		return nil, WindowURLSetOutput{}, nil
	}
	return nil, WindowURLSetOutput{Success: true}, nil
}

// --- window_html_set ---

type WindowHTMLSetInput struct {
	Name string `json:"name"`
	HTML string `json:"html"`
}
type WindowHTMLSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowHTMLSet(_ context.Context, _ *mcp.CallToolRequest, input WindowHTMLSetInput) (*mcp.CallToolResult, WindowHTMLSetOutput, resultFailure) {
	r := s.core.Action("window.set_html").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetHTML{Name: input.Name, HTML: input.HTML}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowHTMLSetOutput{}, e
		}
		return nil, WindowHTMLSetOutput{}, nil
	}
	return nil, WindowHTMLSetOutput{Success: true}, nil
}

// --- window_exec_js ---

type WindowExecJSInput struct {
	Name string `json:"name"`
	JS   string `json:"js"`
}
type WindowExecJSOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowExecJS(_ context.Context, _ *mcp.CallToolRequest, input WindowExecJSInput) (*mcp.CallToolResult, WindowExecJSOutput, resultFailure) {
	r := s.core.Action("window.exec_js").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskExecJS{Name: input.Name, JS: input.JS}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowExecJSOutput{}, e
		}
		return nil, WindowExecJSOutput{}, nil
	}
	return nil, WindowExecJSOutput{Success: true}, nil
}

// --- window_toggle_fullscreen ---

type WindowToggleFullscreenInput struct {
	Name string `json:"name"`
}
type WindowToggleFullscreenOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowToggleFullscreen(_ context.Context, _ *mcp.CallToolRequest, input WindowToggleFullscreenInput) (*mcp.CallToolResult, WindowToggleFullscreenOutput, resultFailure) {
	r := s.core.Action("window.toggle_fullscreen").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskToggleFullscreen{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowToggleFullscreenOutput{}, e
		}
		return nil, WindowToggleFullscreenOutput{}, nil
	}
	return nil, WindowToggleFullscreenOutput{Success: true}, nil
}

// --- window_toggle_maximise ---

type WindowToggleMaximiseInput struct {
	Name string `json:"name"`
}
type WindowToggleMaximiseOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowToggleMaximise(_ context.Context, _ *mcp.CallToolRequest, input WindowToggleMaximiseInput) (*mcp.CallToolResult, WindowToggleMaximiseOutput, resultFailure) {
	r := s.core.Action("window.toggle_maximise").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskToggleMaximise{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowToggleMaximiseOutput{}, e
		}
		return nil, WindowToggleMaximiseOutput{}, nil
	}
	return nil, WindowToggleMaximiseOutput{Success: true}, nil
}

// --- window_set_content_protection ---

type WindowSetContentProtectionInput struct {
	Name       string `json:"name"`
	Protection bool   `json:"protection"`
}
type WindowSetContentProtectionOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowSetContentProtection(_ context.Context, _ *mcp.CallToolRequest, input WindowSetContentProtectionInput) (*mcp.CallToolResult, WindowSetContentProtectionOutput, resultFailure) {
	r := s.core.Action("window.set_content_protection").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetContentProtection{Name: input.Name, Protection: input.Protection}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowSetContentProtectionOutput{}, e
		}
		return nil, WindowSetContentProtectionOutput{}, nil
	}
	return nil, WindowSetContentProtectionOutput{Success: true}, nil
}

// --- window_flash ---

type WindowFlashInput struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
type WindowFlashOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowFlash(_ context.Context, _ *mcp.CallToolRequest, input WindowFlashInput) (*mcp.CallToolResult, WindowFlashOutput, resultFailure) {
	r := s.core.Action("window.flash").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFlash{Name: input.Name, Enabled: input.Enabled}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowFlashOutput{}, e
		}
		return nil, WindowFlashOutput{}, nil
	}
	return nil, WindowFlashOutput{Success: true}, nil
}

// --- window_print ---

type WindowPrintInput struct {
	Name string `json:"name"`
}
type WindowPrintOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowPrint(_ context.Context, _ *mcp.CallToolRequest, input WindowPrintInput) (*mcp.CallToolResult, WindowPrintOutput, resultFailure) {
	r := s.core.Action("window.print").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskPrint{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowPrintOutput{}, e
		}
		return nil, WindowPrintOutput{}, nil
	}
	return nil, WindowPrintOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerWindowTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "window_list", Description: "List all application windows"}, s.windowList)
	addTool(s, server, &mcp.Tool{Name: "window_get", Description: "Get information about a specific window"}, s.windowGet)
	addTool(s, server, &mcp.Tool{Name: "window_focused", Description: "Get the currently focused window"}, s.windowFocused)
	addTool(s, server, &mcp.Tool{
		Name:        "window_create",
		Description: `Create a new application window. Example: {"name":"preview","title":"Preview","url":"https://example.com","x":100,"y":100,"width":1200,"height":800}`,
	}, s.windowCreate)
	addTool(s, server, &mcp.Tool{Name: "window_close", Description: "Close an application window"}, s.windowClose)
	addTool(s, server, &mcp.Tool{Name: "window_position", Description: "Set the position of a window"}, s.windowPosition)
	addTool(s, server, &mcp.Tool{Name: "window_size", Description: "Set the size of a window"}, s.windowSize)
	addTool(s, server, &mcp.Tool{Name: "window_bounds", Description: "Set both position and size of a window"}, s.windowBounds)
	addTool(s, server, &mcp.Tool{Name: "window_maximize", Description: "Maximise a window"}, s.windowMaximize)
	addTool(s, server, &mcp.Tool{Name: "window_minimize", Description: "Minimise a window"}, s.windowMinimize)
	addTool(s, server, &mcp.Tool{Name: "window_restore", Description: "Restore a maximised or minimised window"}, s.windowRestore)
	addTool(s, server, &mcp.Tool{Name: "window_focus", Description: "Bring a window to the front"}, s.windowFocus)
	addTool(s, server, &mcp.Tool{Name: "focus_set", Description: "Set focus to a specific window"}, s.focusSet)
	addTool(s, server, &mcp.Tool{Name: "window_title", Description: "Set the title of a window"}, s.windowTitle)
	addTool(s, server, &mcp.Tool{
		Name:        "window_title_set",
		Description: `Set the title of a window. Example: {"name":"main","title":"Core GUI"}`,
	}, s.windowTitleSet)
	addTool(s, server, &mcp.Tool{Name: "window_title_get", Description: "Get the title of a window"}, s.windowTitleGet)
	addTool(s, server, &mcp.Tool{
		Name:        "window_visibility",
		Description: `Show or hide a window. Example: {"name":"main","visible":false}`,
	}, s.windowVisibility)
	addTool(s, server, &mcp.Tool{Name: "window_always_on_top", Description: "Pin a window above others"}, s.windowAlwaysOnTop)
	addTool(s, server, &mcp.Tool{
		Name:        "window_opacity",
		Description: `Set a window's opacity. Example: {"name":"main","opacity":0.85}`,
	}, s.windowOpacity)
	addTool(s, server, &mcp.Tool{Name: "window_background_colour", Description: "Set a window background colour"}, s.windowBackgroundColour)
	addTool(s, server, &mcp.Tool{Name: "window_fullscreen", Description: "Set a window to fullscreen mode"}, s.windowFullscreen)
	addTool(s, server, &mcp.Tool{Name: "window_zoom_set", Description: "Set a window zoom level"}, s.windowZoomSet)
	addTool(s, server, &mcp.Tool{Name: "window_zoom_in", Description: "Zoom a window in"}, s.windowZoomIn)
	addTool(s, server, &mcp.Tool{Name: "window_zoom_out", Description: "Zoom a window out"}, s.windowZoomOut)
	addTool(s, server, &mcp.Tool{Name: "window_zoom_reset", Description: "Reset a window zoom level"}, s.windowZoomReset)
	addTool(s, server, &mcp.Tool{Name: "window_url_set", Description: "Set a window URL"}, s.windowURLSet)
	addTool(s, server, &mcp.Tool{Name: "window_html_set", Description: "Set a window HTML document"}, s.windowHTMLSet)
	addTool(s, server, &mcp.Tool{Name: "window_exec_js", Description: "Execute JavaScript in a window"}, s.windowExecJS)
	addTool(s, server, &mcp.Tool{Name: "window_toggle_fullscreen", Description: "Toggle fullscreen mode for a window"}, s.windowToggleFullscreen)
	addTool(s, server, &mcp.Tool{Name: "window_toggle_maximise", Description: "Toggle maximised state for a window"}, s.windowToggleMaximise)
	addTool(s, server, &mcp.Tool{Name: "window_set_content_protection", Description: "Enable or disable content protection for a window"}, s.windowSetContentProtection)
	addTool(s, server, &mcp.Tool{Name: "window_flash", Description: "Flash a window to request attention"}, s.windowFlash)
	addTool(s, server, &mcp.Tool{Name: "window_print", Description: "Open the native print dialog for a window"}, s.windowPrint)
}
