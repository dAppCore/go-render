package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/menu"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MenuGetInput struct{}

type MenuOutput struct {
	Items []map[string]any `json:"items"`
}

type MenuSetInput struct {
	Items []map[string]any `json:"items"`
}

func (s *Subsystem) menuGet(_ context.Context, _ *mcp.CallToolRequest, _ MenuGetInput) (*mcp.CallToolResult, MenuOutput, resultFailure) {
	items, err := s.queryMenuItems()
	if err != nil {
		return nil, MenuOutput{}, err
	}
	return nil, MenuOutput{Items: items}, nil
}

func (s *Subsystem) menuSet(_ context.Context, _ *mcp.CallToolRequest, input MenuSetInput) (*mcp.CallToolResult, MenuOutput, resultFailure) {
	items, err := decodeMenuItems(input.Items)
	if err != nil {
		return nil, MenuOutput{}, err
	}
	r := s.core.Action("menu.set_app_menu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: menu.TaskSetAppMenu{Items: items}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, MenuOutput{}, e
		}
		return nil, MenuOutput{}, core.E("mcp.menuSet", "menu.setAppMenu failed", nil)
	}
	snapshot, err := s.queryMenuItems()
	if err != nil {
		return nil, MenuOutput{}, err
	}
	return nil, MenuOutput{Items: snapshot}, nil
}

func (s *Subsystem) queryMenuItems() ([]map[string]any, resultFailure) {
	r := s.core.QUERY(menu.QueryGetAppMenu{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, e
		}
		return nil, core.E("mcp.menuGet", "menu query failed", nil)
	}
	items, ok := r.Value.([]menu.MenuItem)
	if !ok {
		return nil, core.E("mcp.menuGet", "unexpected result type", nil)
	}
	return encodeMenuItems(items), nil
}

func encodeMenuItems(items []menu.MenuItem) []map[string]any {
	if len(items) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		spec := map[string]any{}
		if item.Label != "" {
			spec["label"] = item.Label
		}
		if item.Accelerator != "" {
			spec["accelerator"] = item.Accelerator
		}
		if item.Type != "" {
			spec["type"] = item.Type
		}
		if item.Checked {
			spec["checked"] = item.Checked
		}
		if item.Disabled {
			spec["disabled"] = item.Disabled
		}
		if item.Tooltip != "" {
			spec["tooltip"] = item.Tooltip
		}
		if children := encodeMenuItems(item.Children); len(children) > 0 {
			rawChildren := make([]any, len(children))
			for i, child := range children {
				rawChildren[i] = child
			}
			spec["children"] = rawChildren
		}
		if item.Role != nil {
			spec["role"] = encodeMenuRole(*item.Role)
		}
		out = append(out, spec)
	}
	return out
}

func decodeMenuItems(items []map[string]any) ([]menu.MenuItem, resultFailure) {
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]menu.MenuItem, 0, len(items))
	for _, item := range items {
		roleName, _ := item["role"].(string)
		role, err := decodeMenuRole(roleName)
		if err != nil {
			return nil, err
		}
		children, err := decodeMenuChildren(item["children"])
		if err != nil {
			return nil, err
		}
		out = append(out, menu.MenuItem{
			Label:       stringValue(item, "label"),
			Accelerator: stringValue(item, "accelerator"),
			Type:        stringValue(item, "type"),
			Checked:     boolValue(item, "checked"),
			Disabled:    boolValue(item, "disabled"),
			Tooltip:     stringValue(item, "tooltip"),
			Children:    children,
			Role:        role,
		})
	}
	return out, nil
}

func decodeMenuChildren(value any) ([]menu.MenuItem, resultFailure) {
	switch children := value.(type) {
	case nil:
		return nil, nil
	case []any:
		items := make([]map[string]any, 0, len(children))
		for _, child := range children {
			childMap, ok := child.(map[string]any)
			if !ok {
				return nil, core.E("mcp.decodeMenuChildren", "child menu item must be an object", nil)
			}
			items = append(items, childMap)
		}
		return decodeMenuItems(items)
	case []map[string]any:
		return decodeMenuItems(children)
	default:
		return nil, core.E("mcp.decodeMenuChildren", "children must be an array", nil)
	}
}

func stringValue(item map[string]any, key string) string {
	value, _ := item[key].(string)
	return value
}

func boolValue(item map[string]any, key string) bool {
	value, _ := item[key].(bool)
	return value
}

func encodeMenuRole(role menu.MenuRole) string {
	switch role {
	case menu.RoleAppMenu:
		return "app"
	case menu.RoleFileMenu:
		return "file"
	case menu.RoleEditMenu:
		return "edit"
	case menu.RoleViewMenu:
		return "view"
	case menu.RoleWindowMenu:
		return "window"
	case menu.RoleHelpMenu:
		return "help"
	default:
		return ""
	}
}

func decodeMenuRole(role string) (*menu.MenuRole, resultFailure) {
	switch core.Trim(core.Lower(role)) {
	case "":
		return nil, nil
	case "app":
		value := menu.RoleAppMenu
		return &value, nil
	case "file":
		value := menu.RoleFileMenu
		return &value, nil
	case "edit":
		value := menu.RoleEditMenu
		return &value, nil
	case "view":
		value := menu.RoleViewMenu
		return &value, nil
	case "window":
		value := menu.RoleWindowMenu
		return &value, nil
	case "help":
		value := menu.RoleHelpMenu
		return &value, nil
	default:
		return nil, core.E("mcp.decodeMenuRole", "unknown menu role: "+role, nil)
	}
}

func (s *Subsystem) registerMenuTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "menu_get", Description: "Get the current application menu structure"}, s.menuGet)
	addTool(s, server, &mcp.Tool{Name: "menu_set", Description: "Set the application menu structure"}, s.menuSet)
}
