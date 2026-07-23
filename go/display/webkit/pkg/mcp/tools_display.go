// pkg/mcp/tools_display.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- scheme_resolve ---

type SchemeResolveInput struct {
	URL string `json:"url"`
}

type SchemeResolveOutput struct {
	URL         string `json:"url"`
	Route       string `json:"route"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

func (s *Subsystem) schemeResolve(_ context.Context, _ *mcp.CallToolRequest, input SchemeResolveInput) (*mcp.CallToolResult, SchemeResolveOutput, resultFailure) {
	result := s.core.Action("display.resolve_scheme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, SchemeResolveOutput{}, err
		}
		return nil, SchemeResolveOutput{}, core.E("mcp.schemeResolve", "display.resolveScheme failed", nil)
	}

	payload, ok := result.Value.(map[string]any)
	if !ok {
		return nil, SchemeResolveOutput{}, core.E("mcp.schemeResolve", "unexpected result type", nil)
	}

	output := SchemeResolveOutput{
		URL:         stringValue(payload, "url"),
		Route:       stringValue(payload, "route"),
		ContentType: stringValue(payload, "content_type"),
		Body:        stringValue(payload, "body"),
	}
	if output.URL == "" {
		output.URL = input.URL
	}
	return nil, output, nil
}

func (s *Subsystem) registerDisplayTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "scheme_resolve",
		Description: `Resolve a core:// route or page URL through the display service. Example: {"url":"core://store?q=theme"}`,
	}, s.schemeResolve)
}
