// pkg/mcp/tools_dock.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/dock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- dock_show ---

type DockShowInput struct{}
type DockShowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockShow(_ context.Context, _ *mcp.CallToolRequest, _ DockShowInput) (*mcp.CallToolResult, DockShowOutput, resultFailure) {
	r := s.core.Action("dock.show_icon").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockShowOutput{}, e
		}
		return nil, DockShowOutput{}, nil
	}
	return nil, DockShowOutput{Success: true}, nil
}

// --- dock_hide ---

type DockHideInput struct{}
type DockHideOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockHide(_ context.Context, _ *mcp.CallToolRequest, _ DockHideInput) (*mcp.CallToolResult, DockHideOutput, resultFailure) {
	r := s.core.Action("dock.hide_icon").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockHideOutput{}, e
		}
		return nil, DockHideOutput{}, nil
	}
	return nil, DockHideOutput{Success: true}, nil
}

// --- dock_badge ---

type DockBadgeInput struct {
	Label string `json:"label"`
}
type DockBadgeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockBadge(_ context.Context, _ *mcp.CallToolRequest, input DockBadgeInput) (*mcp.CallToolResult, DockBadgeOutput, resultFailure) {
	r := s.core.Action("dock.set_badge").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dock.TaskSetBadge{Label: input.Label}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockBadgeOutput{}, e
		}
		return nil, DockBadgeOutput{}, nil
	}
	return nil, DockBadgeOutput{Success: true}, nil
}

// --- dock_remove_badge ---

type DockRemoveBadgeInput struct{}

type DockRemoveBadgeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockRemoveBadge(_ context.Context, _ *mcp.CallToolRequest, _ DockRemoveBadgeInput) (*mcp.CallToolResult, DockRemoveBadgeOutput, resultFailure) {
	r := s.core.Action("dock.remove_badge").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockRemoveBadgeOutput{}, e
		}
		return nil, DockRemoveBadgeOutput{}, nil
	}
	return nil, DockRemoveBadgeOutput{Success: true}, nil
}

// --- dock_info ---

type DockInfoInput struct{}

type DockInfoOutput struct {
	Visible bool `json:"visible"`
}

func (s *Subsystem) dockInfo(_ context.Context, _ *mcp.CallToolRequest, _ DockInfoInput) (*mcp.CallToolResult, DockInfoOutput, resultFailure) {
	r := s.core.QUERY(dock.QueryVisible{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockInfoOutput{}, e
		}
		return nil, DockInfoOutput{}, nil
	}
	visible, ok := r.Value.(bool)
	if !ok {
		return nil, DockInfoOutput{}, core.E("mcp.dockInfo", "unexpected result type", nil)
	}
	return nil, DockInfoOutput{Visible: visible}, nil
}

// --- dock_set_progress_bar ---

type DockSetProgressBarInput struct {
	Progress float64 `json:"progress"`
}

type DockSetProgressBarOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockSetProgressBar(_ context.Context, _ *mcp.CallToolRequest, input DockSetProgressBarInput) (*mcp.CallToolResult, DockSetProgressBarOutput, resultFailure) {
	r := s.core.Action("dock.set_progress_bar").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dock.TaskSetProgressBar{Progress: input.Progress}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockSetProgressBarOutput{}, e
		}
		return nil, DockSetProgressBarOutput{}, nil
	}
	return nil, DockSetProgressBarOutput{Success: true}, nil
}

// --- dock_bounce ---

type DockBounceInput struct {
	BounceType dock.BounceType `json:"bounceType"`
}

type DockBounceOutput struct {
	RequestID int `json:"requestId"`
}

func (s *Subsystem) dockBounce(_ context.Context, _ *mcp.CallToolRequest, input DockBounceInput) (*mcp.CallToolResult, DockBounceOutput, resultFailure) {
	r := s.core.Action("dock.bounce").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dock.TaskBounce{BounceType: input.BounceType}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockBounceOutput{}, e
		}
		return nil, DockBounceOutput{}, nil
	}
	requestID, ok := r.Value.(int)
	if !ok {
		return nil, DockBounceOutput{}, core.E("mcp.dockBounce", "unexpected result type", nil)
	}
	return nil, DockBounceOutput{RequestID: requestID}, nil
}

// --- dock_stop_bounce ---

type DockStopBounceInput struct {
	RequestID int `json:"requestId"`
}

type DockStopBounceOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockStopBounce(_ context.Context, _ *mcp.CallToolRequest, input DockStopBounceInput) (*mcp.CallToolResult, DockStopBounceOutput, resultFailure) {
	r := s.core.Action("dock.stop_bounce").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dock.TaskStopBounce{RequestID: input.RequestID}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockStopBounceOutput{}, e
		}
		return nil, DockStopBounceOutput{}, nil
	}
	return nil, DockStopBounceOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerDockTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "dock_show", Description: "Show the dock/taskbar icon"}, s.dockShow)
	addTool(s, server, &mcp.Tool{Name: "dock_hide", Description: "Hide the dock/taskbar icon"}, s.dockHide)
	addTool(s, server, &mcp.Tool{Name: "dock_badge", Description: "Set the dock/taskbar badge label"}, s.dockBadge)
	addTool(s, server, &mcp.Tool{Name: "dock_remove_badge", Description: "Remove the dock/taskbar badge label"}, s.dockRemoveBadge)
	addTool(s, server, &mcp.Tool{Name: "dock_info", Description: "Get the current dock/taskbar visibility"}, s.dockInfo)
	addTool(s, server, &mcp.Tool{Name: "dock_set_progress_bar", Description: "Set the dock/taskbar progress indicator"}, s.dockSetProgressBar)
	addTool(s, server, &mcp.Tool{Name: "dock_bounce", Description: "Request dock/taskbar attention"}, s.dockBounce)
	addTool(s, server, &mcp.Tool{Name: "dock_stop_bounce", Description: "Cancel a dock/taskbar attention request"}, s.dockStopBounce)
}
