// pkg/mcp/tools_screen.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- screen_list ---

type ScreenListInput struct{}
type ScreenListOutput struct {
	Screens []screen.Screen `json:"screens"`
}

func (s *Subsystem) screenList(_ context.Context, _ *mcp.CallToolRequest, _ ScreenListInput) (*mcp.CallToolResult, ScreenListOutput, resultFailure) {
	r := s.core.QUERY(screen.QueryAll{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenListOutput{}, e
		}
		return nil, ScreenListOutput{}, core.E("mcp.screenList", "screen query failed", nil)
	}
	screens, ok := r.Value.([]screen.Screen)
	if !ok {
		return nil, ScreenListOutput{}, core.E("mcp.screenList", "unexpected result type", nil)
	}
	return nil, ScreenListOutput{Screens: screens}, nil
}

// --- screen_get ---

type ScreenGetInput struct {
	ID string `json:"id"`
}
type ScreenGetOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenGet(_ context.Context, _ *mcp.CallToolRequest, input ScreenGetInput) (*mcp.CallToolResult, ScreenGetOutput, resultFailure) {
	r := s.core.QUERY(screen.QueryByID{ID: input.ID})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenGetOutput{}, e
		}
		return nil, ScreenGetOutput{}, core.E("mcp.screenGet", "screen query failed", nil)
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, ScreenGetOutput{}, core.E("mcp.screenGet", "unexpected result type", nil)
	}
	return nil, ScreenGetOutput{Screen: scr}, nil
}

// --- screen_primary ---

type ScreenPrimaryInput struct{}
type ScreenPrimaryOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenPrimary(_ context.Context, _ *mcp.CallToolRequest, _ ScreenPrimaryInput) (*mcp.CallToolResult, ScreenPrimaryOutput, resultFailure) {
	r := s.core.QUERY(screen.QueryPrimary{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenPrimaryOutput{}, e
		}
		return nil, ScreenPrimaryOutput{}, nil
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, ScreenPrimaryOutput{}, core.E("mcp.screenPrimary", "unexpected result type", nil)
	}
	return nil, ScreenPrimaryOutput{Screen: scr}, nil
}

// --- screen_at_point ---

type ScreenAtPointInput struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type ScreenAtPointOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenAtPoint(_ context.Context, _ *mcp.CallToolRequest, input ScreenAtPointInput) (*mcp.CallToolResult, ScreenAtPointOutput, resultFailure) {
	r := s.core.QUERY(screen.QueryAtPoint{X: input.X, Y: input.Y})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenAtPointOutput{}, e
		}
		return nil, ScreenAtPointOutput{}, nil
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, ScreenAtPointOutput{}, core.E("mcp.screenAtPoint", "unexpected result type", nil)
	}
	return nil, ScreenAtPointOutput{Screen: scr}, nil
}

// --- screen_work_areas ---

type ScreenWorkAreasInput struct{}
type ScreenWorkAreasOutput struct {
	WorkAreas []screen.Rect `json:"workAreas"`
}

func (s *Subsystem) screenWorkAreas(_ context.Context, _ *mcp.CallToolRequest, _ ScreenWorkAreasInput) (*mcp.CallToolResult, ScreenWorkAreasOutput, resultFailure) {
	r := s.core.QUERY(screen.QueryWorkAreas{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenWorkAreasOutput{}, e
		}
		return nil, ScreenWorkAreasOutput{}, nil
	}
	areas, ok := r.Value.([]screen.Rect)
	if !ok {
		return nil, ScreenWorkAreasOutput{}, core.E("mcp.screenWorkAreas", "unexpected result type", nil)
	}
	return nil, ScreenWorkAreasOutput{WorkAreas: areas}, nil
}

// --- screen_work_area ---

type ScreenWorkAreaInput struct {
	ID string `json:"id,omitempty"`
}
type ScreenWorkAreaOutput struct {
	WorkArea screen.Rect `json:"workArea"`
}

func (s *Subsystem) screenWorkArea(_ context.Context, _ *mcp.CallToolRequest, input ScreenWorkAreaInput) (*mcp.CallToolResult, ScreenWorkAreaOutput, resultFailure) {
	var query core.Query = screen.QueryPrimary{}
	if input.ID != "" {
		query = screen.QueryByID{ID: input.ID}
	}

	r := s.core.QUERY(query)
	if !r.OK {
		return nil, ScreenWorkAreaOutput{}, nil
	}
	scr, _ := r.Value.(*screen.Screen)
	if scr == nil {
		return nil, ScreenWorkAreaOutput{}, nil
	}

	workArea := scr.WorkArea
	if workArea.IsEmpty() {
		workArea = scr.Bounds
	}
	return nil, ScreenWorkAreaOutput{WorkArea: workArea}, nil
}

// --- screen_for_window ---

type ScreenForWindowInput struct {
	Name string `json:"name"`
}
type ScreenForWindowOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenForWindow(_ context.Context, _ *mcp.CallToolRequest, input ScreenForWindowInput) (*mcp.CallToolResult, ScreenForWindowOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if !r.OK {
		return nil, ScreenForWindowOutput{}, nil
	}
	info, _ := r.Value.(*window.WindowInfo)
	if info == nil {
		return nil, ScreenForWindowOutput{}, nil
	}
	centerX := info.X + info.Width/2
	centerY := info.Y + info.Height/2
	r2 := s.core.QUERY(screen.QueryAtPoint{X: centerX, Y: centerY})
	if !r2.OK {
		return nil, ScreenForWindowOutput{}, nil
	}
	scr, _ := r2.Value.(*screen.Screen)
	return nil, ScreenForWindowOutput{Screen: scr}, nil
}

// --- Registration ---

func (s *Subsystem) registerScreenTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "screen_list", Description: "List all connected displays/screens"}, s.screenList)
	addTool(s, server, &mcp.Tool{Name: "screen_get", Description: "Get information about a specific screen"}, s.screenGet)
	addTool(s, server, &mcp.Tool{Name: "screen_primary", Description: "Get the primary screen"}, s.screenPrimary)
	addTool(s, server, &mcp.Tool{Name: "screen_at_point", Description: "Get the screen at a specific point"}, s.screenAtPoint)
	addTool(s, server, &mcp.Tool{Name: "screen_work_area", Description: "Get the usable work area for a screen"}, s.screenWorkArea)
	addTool(s, server, &mcp.Tool{Name: "screen_work_areas", Description: "Get work areas for all screens"}, s.screenWorkAreas)
	addTool(s, server, &mcp.Tool{Name: "screen_for_window", Description: "Get the screen containing a window"}, s.screenForWindow)
}
