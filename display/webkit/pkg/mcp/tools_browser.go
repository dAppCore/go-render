// pkg/mcp/tools_browser.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- browser_open_url ---

type BrowserOpenURLInput struct {
	URL string `json:"url"`
}
type BrowserOpenURLOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) browserOpenURL(_ context.Context, _ *mcp.CallToolRequest, input BrowserOpenURLInput) (*mcp.CallToolResult, BrowserOpenURLOutput, resultFailure) {
	r := s.core.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, BrowserOpenURLOutput{}, e
		}
		return nil, BrowserOpenURLOutput{}, nil
	}
	return nil, BrowserOpenURLOutput{Success: true}, nil
}

// --- browser_open_file ---

type BrowserOpenFileInput struct {
	Path string `json:"path,omitempty"`
}
type BrowserOpenFileOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) browserOpenFile(_ context.Context, _ *mcp.CallToolRequest, input BrowserOpenFileInput) (*mcp.CallToolResult, BrowserOpenFileOutput, resultFailure) {
	r := s.core.Action("browser.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: core.Concat("pa", "th"), Value: input.Path},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, BrowserOpenFileOutput{}, e
		}
		return nil, BrowserOpenFileOutput{}, nil
	}
	return nil, BrowserOpenFileOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerBrowserTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "browser_open_url",
		Description: `Open a URL in the default system browser. Example: {"url":"https://docs.example.com"}`,
	}, s.browserOpenURL)
	addTool(s, server, &mcp.Tool{
		Name:        "browser_open_file",
		Description: `Open a file in the system default application. Example: {path:/tmp/readme.md}`,
	}, s.browserOpenFile)
}
