// pkg/mcp/tools_tray.go
package mcp

import (
	"context"
	"encoding/base64"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/systray"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- tray_set_icon ---

type TraySetIconInput struct {
	Base64 string `json:"base64"`
}
type TraySetIconOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetIcon(_ context.Context, _ *mcp.CallToolRequest, input TraySetIconInput) (*mcp.CallToolResult, TraySetIconOutput, resultFailure) {
	if input.Base64 == "" {
		return nil, TraySetIconOutput{}, core.E("mcp.traySetIcon", "base64 icon data is required", nil)
	}
	data, err := base64.StdEncoding.DecodeString(input.Base64)
	if err != nil {
		return nil, TraySetIconOutput{}, core.E("mcp.traySetIcon", "invalid base64 icon data", err)
	}
	r := s.core.Action("systray.set_icon").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayIcon{Data: data}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TraySetIconOutput{}, e
		}
		return nil, TraySetIconOutput{}, nil
	}
	return nil, TraySetIconOutput{Success: true}, nil
}

// --- tray_set_tooltip ---

type TraySetTooltipInput struct {
	Tooltip string `json:"tooltip"`
}
type TraySetTooltipOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetTooltip(_ context.Context, _ *mcp.CallToolRequest, input TraySetTooltipInput) (*mcp.CallToolResult, TraySetTooltipOutput, resultFailure) {
	r := s.core.Action("systray.set_tooltip").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayTooltip{Tooltip: input.Tooltip}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TraySetTooltipOutput{}, e
		}
		return nil, TraySetTooltipOutput{}, nil
	}
	return nil, TraySetTooltipOutput{Success: true}, nil
}

// --- tray_set_label ---

type TraySetLabelInput struct {
	Label string `json:"label"`
}
type TraySetLabelOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetLabel(_ context.Context, _ *mcp.CallToolRequest, input TraySetLabelInput) (*mcp.CallToolResult, TraySetLabelOutput, resultFailure) {
	r := s.core.Action("systray.set_label").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayLabel{Label: input.Label}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TraySetLabelOutput{}, e
		}
		return nil, TraySetLabelOutput{}, nil
	}
	return nil, TraySetLabelOutput{Success: true}, nil
}

// --- tray_set_menu ---

type TraySetMenuInput struct {
	Items []systray.TrayMenuItem `json:"items"`
}

type TraySetMenuOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetMenu(_ context.Context, _ *mcp.CallToolRequest, input TraySetMenuInput) (*mcp.CallToolResult, TraySetMenuOutput, resultFailure) {
	r := s.core.Action("systray.set_menu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayMenu{Items: input.Items}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TraySetMenuOutput{}, e
		}
		return nil, TraySetMenuOutput{}, nil
	}
	return nil, TraySetMenuOutput{Success: true}, nil
}

// --- tray_show_message ---

type TrayShowMessageInput struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type TrayShowMessageOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) trayShowMessage(_ context.Context, _ *mcp.CallToolRequest, input TrayShowMessageInput) (*mcp.CallToolResult, TrayShowMessageOutput, resultFailure) {
	r := s.core.Action("systray.show_message").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskShowMessage{Title: input.Title, Message: input.Message}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TrayShowMessageOutput{}, e
		}
		return nil, TrayShowMessageOutput{}, nil
	}
	return nil, TrayShowMessageOutput{Success: true}, nil
}

// --- tray_info ---

type TrayInfoInput struct{}
type TrayInfoOutput struct {
	Config map[string]any `json:"config"`
}

func (s *Subsystem) trayInfo(_ context.Context, _ *mcp.CallToolRequest, _ TrayInfoInput) (*mcp.CallToolResult, TrayInfoOutput, resultFailure) {
	r := s.core.QUERY(systray.QueryInfo{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TrayInfoOutput{}, e
		}
		return nil, TrayInfoOutput{}, nil
	}
	config, ok := r.Value.(map[string]any)
	if !ok {
		return nil, TrayInfoOutput{}, core.E("mcp.trayInfo", "unexpected result type", nil)
	}
	return nil, TrayInfoOutput{Config: config}, nil
}

// --- Registration ---

func (s *Subsystem) registerTrayTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "tray_set_icon",
		Description: `Set the system tray icon from base64 PNG data. Example: {"base64":"iVBORw0KGgoAAA..."}`,
	}, s.traySetIcon)
	addTool(s, server, &mcp.Tool{Name: "tray_set_tooltip", Description: "Set the system tray tooltip"}, s.traySetTooltip)
	addTool(s, server, &mcp.Tool{Name: "tray_set_label", Description: "Set the system tray label"}, s.traySetLabel)
	addTool(s, server, &mcp.Tool{
		Name:        "tray_set_menu",
		Description: "Set the system tray menu items",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"items": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "object"},
				},
			},
			"required": []string{"items"},
		},
	}, s.traySetMenu)
	addTool(s, server, &mcp.Tool{
		Name:        "tray_show_message",
		Description: `Show a tray balloon notification. Example: {"title":"Sync complete","message":"Files are up to date"}`,
	}, s.trayShowMessage)
	addTool(s, server, &mcp.Tool{Name: "tray_info", Description: "Get system tray configuration"}, s.trayInfo)
}
