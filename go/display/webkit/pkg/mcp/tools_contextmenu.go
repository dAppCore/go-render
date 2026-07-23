// pkg/mcp/tools_contextmenu.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/contextmenu"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- contextmenu_add ---

// ContextMenuAddInput uses map[string]any for the menu definition because
// contextmenu.ContextMenuDef contains self-referencing MenuItemDef (Items []MenuItemDef)
// which the MCP SDK schema generator cannot handle (cycle detection panic).
type ContextMenuAddInput struct {
	Name string         `json:"name"`
	Menu map[string]any `json:"menu"`
}
type ContextMenuAddOutput struct {
	Success bool `json:"success"`
}

func jsonBytesFromResult(op, message string, result core.Result) ([]byte, resultFailure) {
	if !result.OK {
		if err, ok := result.Value.(error); ok && err != nil {
			return nil, core.E(op, message, err)
		}
		return nil, core.E(op, message, nil)
	}

	data, ok := result.Value.([]byte)
	if !ok {
		return nil, core.E(op, core.Sprintf("%s: unexpected helper result type %T", message, result.Value), nil)
	}
	return data, nil
}

func resultError(result core.Result) resultFailure {
	if err, ok := result.Value.(error); ok && err != nil {
		return err
	}
	return nil
}

func (s *Subsystem) contextMenuAdd(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuAddInput) (*mcp.CallToolResult, ContextMenuAddOutput, resultFailure) {
	// Convert map[string]any to ContextMenuDef via JSON round-trip
	menuJSON, err := jsonBytesFromResult("mcp.contextMenuAdd", "failed to marshal menu definition", core.JSONMarshal(input.Menu))
	if err != nil {
		return nil, ContextMenuAddOutput{}, err
	}

	var menuDef contextmenu.ContextMenuDef
	unmarshalResult := core.JSONUnmarshal(menuJSON, &menuDef)
	if !unmarshalResult.OK {
		return nil, ContextMenuAddOutput{}, core.E("mcp.contextMenuAdd", "failed to unmarshal menu definition", resultError(unmarshalResult))
	}
	r := s.core.Action("contextmenu.add").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: contextmenu.TaskAdd{Name: input.Name, Menu: menuDef}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuAddOutput{}, e
		}
		return nil, ContextMenuAddOutput{}, nil
	}
	return nil, ContextMenuAddOutput{Success: true}, nil
}

// --- contextmenu_remove ---

type ContextMenuRemoveInput struct {
	Name string `json:"name"`
}
type ContextMenuRemoveOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) contextMenuRemove(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuRemoveInput) (*mcp.CallToolResult, ContextMenuRemoveOutput, resultFailure) {
	r := s.core.Action("contextmenu.remove").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: contextmenu.TaskRemove{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuRemoveOutput{}, e
		}
		return nil, ContextMenuRemoveOutput{}, nil
	}
	return nil, ContextMenuRemoveOutput{Success: true}, nil
}

// --- contextmenu_get ---

type ContextMenuGetInput struct {
	Name string `json:"name"`
}
type ContextMenuGetOutput struct {
	Menu map[string]any `json:"menu"`
}

func (s *Subsystem) contextMenuGet(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuGetInput) (*mcp.CallToolResult, ContextMenuGetOutput, resultFailure) {
	r := s.core.QUERY(contextmenu.QueryGet{Name: input.Name})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuGetOutput{}, e
		}
		return nil, ContextMenuGetOutput{}, nil
	}
	menu, ok := r.Value.(*contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuGetOutput{}, core.E("mcp.contextMenuGet", "unexpected result type", nil)
	}
	if menu == nil {
		return nil, ContextMenuGetOutput{}, nil
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menuJSON, err := jsonBytesFromResult("mcp.contextMenuGet", "failed to marshal context menu", core.JSONMarshal(menu))
	if err != nil {
		return nil, ContextMenuGetOutput{}, err
	}

	var menuMap map[string]any
	unmarshalResult := core.JSONUnmarshal(menuJSON, &menuMap)
	if !unmarshalResult.OK {
		return nil, ContextMenuGetOutput{}, core.E("mcp.contextMenuGet", "failed to unmarshal context menu", resultError(unmarshalResult))
	}
	return nil, ContextMenuGetOutput{Menu: menuMap}, nil
}

// --- contextmenu_list ---

type ContextMenuListInput struct{}
type ContextMenuListOutput struct {
	Menus map[string]any `json:"menus"`
}

func (s *Subsystem) contextMenuList(_ context.Context, _ *mcp.CallToolRequest, _ ContextMenuListInput) (*mcp.CallToolResult, ContextMenuListOutput, resultFailure) {
	r := s.core.QUERY(contextmenu.QueryList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuListOutput{}, e
		}
		return nil, ContextMenuListOutput{}, nil
	}
	menus, ok := r.Value.(map[string]contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuListOutput{}, core.E("mcp.contextMenuList", "unexpected result type", nil)
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menusJSON, err := jsonBytesFromResult("mcp.contextMenuList", "failed to marshal context menus", core.JSONMarshal(menus))
	if err != nil {
		return nil, ContextMenuListOutput{}, err
	}

	var menusMap map[string]any
	unmarshalResult := core.JSONUnmarshal(menusJSON, &menusMap)
	if !unmarshalResult.OK {
		return nil, ContextMenuListOutput{}, core.E("mcp.contextMenuList", "failed to unmarshal context menus", resultError(unmarshalResult))
	}
	return nil, ContextMenuListOutput{Menus: menusMap}, nil
}

// --- Registration ---

func (s *Subsystem) registerContextMenuTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "contextmenu_add", Description: "Register a context menu"}, s.contextMenuAdd)
	addTool(s, server, &mcp.Tool{Name: "contextmenu_remove", Description: "Unregister a context menu"}, s.contextMenuRemove)
	addTool(s, server, &mcp.Tool{Name: "contextmenu_get", Description: "Get a context menu by name"}, s.contextMenuGet)
	addTool(s, server, &mcp.Tool{Name: "contextmenu_list", Description: "List all registered context menus"}, s.contextMenuList)
}
