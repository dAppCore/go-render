package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ContainerDetectInput struct{}
type ContainerDetectOutput struct {
	Runtime container.ContainerRuntime `json:"runtime"`
}

func (s *Subsystem) containerDetect(_ context.Context, _ *mcp.CallToolRequest, _ ContainerDetectInput) (*mcp.CallToolResult, ContainerDetectOutput, resultFailure) {
	result := s.core.Action("container.runtime.detect").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ContainerDetectOutput{}, err
		}
		return nil, ContainerDetectOutput{}, core.E("mcp.containerDetect", "container.runtime.detect failed", nil)
	}
	runtime, ok := result.Value.(container.ContainerRuntime)
	if !ok {
		return nil, ContainerDetectOutput{}, core.E("mcp.containerDetect", "unexpected result type", nil)
	}
	return nil, ContainerDetectOutput{Runtime: runtime}, nil
}

type TIMStateInput struct{}
type TIMStateOutput struct {
	State container.TIMState `json:"state"`
}

func (s *Subsystem) timStatus(_ context.Context, _ *mcp.CallToolRequest, _ TIMStateInput) (*mcp.CallToolResult, TIMStateOutput, resultFailure) {
	result := s.core.Action("tim.status").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, TIMStateOutput{}, err
		}
		return nil, TIMStateOutput{}, core.E("mcp.timStatus", "tim.status failed", nil)
	}
	state, ok := result.Value.(container.TIMState)
	if !ok {
		return nil, TIMStateOutput{}, core.E("mcp.timStatus", "unexpected result type", nil)
	}
	return nil, TIMStateOutput{State: state}, nil
}

type TIMStartInput struct{}
type TIMStartOutput struct {
	State container.TIMState `json:"state"`
}

func (s *Subsystem) timStart(_ context.Context, _ *mcp.CallToolRequest, _ TIMStartInput) (*mcp.CallToolResult, TIMStartOutput, resultFailure) {
	result := s.core.Action("tim.start").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, TIMStartOutput{}, err
		}
		return nil, TIMStartOutput{}, core.E("mcp.timStart", "tim.start failed", nil)
	}
	state, ok := result.Value.(container.TIMState)
	if !ok {
		return nil, TIMStartOutput{}, core.E("mcp.timStart", "unexpected result type", nil)
	}
	return nil, TIMStartOutput{State: state}, nil
}

type TIMStopInput struct{}
type TIMStopOutput struct {
	State container.TIMState `json:"state"`
}

func (s *Subsystem) timStop(_ context.Context, _ *mcp.CallToolRequest, _ TIMStopInput) (*mcp.CallToolResult, TIMStopOutput, resultFailure) {
	result := s.core.Action("tim.stop").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, TIMStopOutput{}, err
		}
		return nil, TIMStopOutput{}, core.E("mcp.timStop", "tim.stop failed", nil)
	}
	state, ok := result.Value.(container.TIMState)
	if !ok {
		return nil, TIMStopOutput{}, core.E("mcp.timStop", "unexpected result type", nil)
	}
	return nil, TIMStopOutput{State: state}, nil
}

func (s *Subsystem) registerContainerTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "container_detect_runtime", Description: "Detect the preferred isolated workload runtime on this host"}, s.containerDetect)
	addTool(s, server, &mcp.Tool{Name: "tim_status", Description: "Inspect the TIM container state"}, s.timStatus)
	addTool(s, server, &mcp.Tool{Name: "tim_start", Description: "Start the TIM container"}, s.timStart)
	addTool(s, server, &mcp.Tool{Name: "tim_stop", Description: "Stop the TIM container"}, s.timStop)
}
