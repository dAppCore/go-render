// pkg/mcp/tools_keybinding.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/keybinding"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- keybinding_add ---

type KeybindingAddInput struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description,omitempty"`
}
type KeybindingAddOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) keybindingAdd(_ context.Context, _ *mcp.CallToolRequest, input KeybindingAddInput) (*mcp.CallToolResult, KeybindingAddOutput, resultFailure) {
	r := s.core.Action("keybinding.add").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: keybinding.TaskAdd{Accelerator: input.Accelerator, Description: input.Description}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, KeybindingAddOutput{}, e
		}
		return nil, KeybindingAddOutput{}, nil
	}
	return nil, KeybindingAddOutput{Success: true}, nil
}

// --- keybinding_remove ---

type KeybindingRemoveInput struct {
	Accelerator string `json:"accelerator"`
}
type KeybindingRemoveOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) keybindingRemove(_ context.Context, _ *mcp.CallToolRequest, input KeybindingRemoveInput) (*mcp.CallToolResult, KeybindingRemoveOutput, resultFailure) {
	r := s.core.Action("keybinding.remove").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: keybinding.TaskRemove{Accelerator: input.Accelerator}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, KeybindingRemoveOutput{}, e
		}
		return nil, KeybindingRemoveOutput{}, nil
	}
	return nil, KeybindingRemoveOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerKeybindingTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "keybinding_add", Description: "Register a keyboard shortcut"}, s.keybindingAdd)
	addTool(s, server, &mcp.Tool{Name: "keybinding_remove", Description: "Unregister a keyboard shortcut"}, s.keybindingRemove)
}
