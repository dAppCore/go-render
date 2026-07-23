// pkg/mcp/tools_lifecycle.go
package mcp

import (
	"context"

	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
	"dappco.re/go/render/display/webkit/pkg/lifecycle"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- app_quit ---

type AppQuitInput struct{}
type AppQuitOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) appQuit(_ context.Context, _ *mcp.CallToolRequest, _ AppQuitInput) (*mcp.CallToolResult, AppQuitOutput, resultFailure) {
	// Broadcast the will-terminate action which triggers application shutdown
	coreutil.DispatchAction(s.core, "mcp.appQuit", lifecycle.ActionWillTerminate{})
	return nil, AppQuitOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerLifecycleTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "app_quit", Description: "Quit the application"}, s.appQuit)
}
