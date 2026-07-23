// pkg/mcp/tools_environment.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- theme_get ---

type ThemeGetInput struct{}
type ThemeGetOutput struct {
	Theme environment.ThemeInfo `json:"theme"`
}

func (s *Subsystem) themeGet(_ context.Context, _ *mcp.CallToolRequest, _ ThemeGetInput) (*mcp.CallToolResult, ThemeGetOutput, resultFailure) {
	result := s.core.QUERY(environment.QueryTheme{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ThemeGetOutput{}, err
		}
		return nil, ThemeGetOutput{}, core.E("mcp.themeGet", "theme query failed", nil)
	}
	theme, ok := result.Value.(environment.ThemeInfo)
	if !ok {
		return nil, ThemeGetOutput{}, core.E("mcp.themeGet", "unexpected result type", nil)
	}
	return nil, ThemeGetOutput{Theme: theme}, nil
}

// --- theme_system ---

type ThemeSystemInput struct{}
type ThemeSystemOutput struct {
	Info environment.EnvironmentInfo `json:"info"`
}

func (s *Subsystem) themeSystem(_ context.Context, _ *mcp.CallToolRequest, _ ThemeSystemInput) (*mcp.CallToolResult, ThemeSystemOutput, resultFailure) {
	result := s.core.QUERY(environment.QueryInfo{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ThemeSystemOutput{}, err
		}
		return nil, ThemeSystemOutput{}, nil
	}
	info, ok := result.Value.(environment.EnvironmentInfo)
	if !ok {
		return nil, ThemeSystemOutput{}, core.E("mcp.themeSystem", "unexpected result type", nil)
	}
	return nil, ThemeSystemOutput{Info: info}, nil
}

// --- theme_set ---

type ThemeSetInput struct {
	Theme string `json:"theme"`
}

type ThemeSetOutput struct {
	Theme environment.ThemeInfo `json:"theme"`
}

func (s *Subsystem) themeSet(_ context.Context, _ *mcp.CallToolRequest, input ThemeSetInput) (*mcp.CallToolResult, ThemeSetOutput, resultFailure) {
	result := s.core.Action("environment.set_theme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: environment.TaskSetTheme{Theme: input.Theme}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ThemeSetOutput{}, err
		}
		return nil, ThemeSetOutput{}, nil
	}
	theme, ok := result.Value.(environment.ThemeInfo)
	if !ok {
		return nil, ThemeSetOutput{}, core.E("mcp.themeSet", "unexpected result type", nil)
	}
	return nil, ThemeSetOutput{Theme: theme}, nil
}

// --- Registration ---

func (s *Subsystem) registerEnvironmentTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "theme_get", Description: "Get the current application theme"}, s.themeGet)
	addTool(s, server, &mcp.Tool{
		Name:        "theme_set",
		Description: `Set the application theme to dark, light, or system. Example: {"theme":"dark"}`,
	}, s.themeSet)
	addTool(s, server, &mcp.Tool{Name: "theme_system", Description: "Get system environment and theme information"}, s.themeSystem)
}
