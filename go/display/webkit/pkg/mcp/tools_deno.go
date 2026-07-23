package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/deno"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DenoStatusInput struct{}

type DenoStatusOutput struct {
	Status deno.Status `json:"status"`
}

func (s *Subsystem) denoStatus(_ context.Context, _ *mcp.CallToolRequest, _ DenoStatusInput) (*mcp.CallToolResult, DenoStatusOutput, resultFailure) {
	result := s.core.Action("core.deno.sidecar.status").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, DenoStatusOutput{}, err
		}
		return nil, DenoStatusOutput{}, core.E("mcp.denoStatus", "core.deno.sidecar.status failed", nil)
	}
	status, ok := result.Value.(deno.Status)
	if !ok {
		return nil, DenoStatusOutput{}, core.E("mcp.denoStatus", "unexpected result type", nil)
	}
	return nil, DenoStatusOutput{Status: status}, nil
}

type DenoStartInput struct{}
type DenoStartOutput struct {
	Status deno.Status `json:"status"`
}

func (s *Subsystem) denoStart(_ context.Context, _ *mcp.CallToolRequest, _ DenoStartInput) (*mcp.CallToolResult, DenoStartOutput, resultFailure) {
	result := s.core.Action("core.deno.sidecar.start").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, DenoStartOutput{}, err
		}
		return nil, DenoStartOutput{}, core.E("mcp.denoStart", "core.deno.sidecar.start failed", nil)
	}
	status, ok := result.Value.(deno.Status)
	if !ok {
		return nil, DenoStartOutput{}, core.E("mcp.denoStart", "unexpected result type", nil)
	}
	return nil, DenoStartOutput{Status: status}, nil
}

type DenoStopInput struct{}
type DenoStopOutput struct {
	Status deno.Status `json:"status"`
}

func (s *Subsystem) denoStop(_ context.Context, _ *mcp.CallToolRequest, _ DenoStopInput) (*mcp.CallToolResult, DenoStopOutput, resultFailure) {
	result := s.core.Action("core.deno.sidecar.stop").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, DenoStopOutput{}, err
		}
		return nil, DenoStopOutput{}, core.E("mcp.denoStop", "core.deno.sidecar.stop failed", nil)
	}
	status, ok := result.Value.(deno.Status)
	if !ok {
		return nil, DenoStopOutput{}, core.E("mcp.denoStop", "unexpected result type", nil)
	}
	return nil, DenoStopOutput{Status: status}, nil
}

type DenoEvalInput struct {
	Code string `json:"code"`
}

type DenoEvalOutput struct {
	Result deno.EvalResult `json:"result"`
}

func (s *Subsystem) denoEval(_ context.Context, _ *mcp.CallToolRequest, input DenoEvalInput) (*mcp.CallToolResult, DenoEvalOutput, resultFailure) {
	result := s.core.Action("core.deno.sidecar.eval").Run(context.Background(), core.NewOptions(
		core.Option{Key: "code", Value: input.Code},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, DenoEvalOutput{}, err
		}
		return nil, DenoEvalOutput{}, core.E("mcp.denoEval", "core.deno.sidecar.eval failed", nil)
	}
	value, ok := result.Value.(deno.EvalResult)
	if !ok {
		return nil, DenoEvalOutput{}, core.E("mcp.denoEval", "unexpected result type", nil)
	}
	return nil, DenoEvalOutput{Result: value}, nil
}

func (s *Subsystem) registerDenoTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "deno_status", Description: "Inspect the CoreDeno sidecar process and IPC connection state"}, s.denoStatus)
	addTool(s, server, &mcp.Tool{Name: "deno_start", Description: "Start the CoreDeno sidecar process"}, s.denoStart)
	addTool(s, server, &mcp.Tool{Name: "deno_stop", Description: "Stop the CoreDeno sidecar process"}, s.denoStop)
	addTool(s, server, &mcp.Tool{Name: "deno_eval", Description: `Evaluate JavaScript inside the CoreDeno sidecar. Example: {"code":"await core.action('display.models.state')"} `}, s.denoEval)
}
