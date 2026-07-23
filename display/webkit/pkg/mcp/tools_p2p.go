package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/p2p"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type P2PPublishInput struct {
	Topic    string         `json:"topic"`
	Route    string         `json:"route,omitempty"`
	SenderID string         `json:"sender_id,omitempty"`
	Payload  map[string]any `json:"payload,omitempty"`
}

type P2PPublishOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) p2pPublish(_ context.Context, _ *mcp.CallToolRequest, input P2PPublishInput) (*mcp.CallToolResult, P2PPublishOutput, resultFailure) {
	result := s.core.Action("p2p.publish").Run(context.Background(), core.NewOptions(
		core.Option{Key: "topic", Value: input.Topic},
		core.Option{Key: "route", Value: input.Route},
		core.Option{Key: "sender_id", Value: input.SenderID},
		core.Option{Key: "payload", Value: input.Payload},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, P2PPublishOutput{}, err
		}
		return nil, P2PPublishOutput{}, core.E("mcp.p2pPublish", "p2p.publish failed", nil)
	}
	return nil, P2PPublishOutput{Success: true}, nil
}

type P2PStateInput struct{}

type P2PStateOutput struct {
	State p2p.State `json:"state"`
}

func (s *Subsystem) p2pState(_ context.Context, _ *mcp.CallToolRequest, _ P2PStateInput) (*mcp.CallToolResult, P2PStateOutput, resultFailure) {
	result := s.core.Action("p2p.state").Run(context.Background(), core.Options{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, P2PStateOutput{}, err
		}
		return nil, P2PStateOutput{}, core.E("mcp.p2pState", "p2p.state failed", nil)
	}
	state, ok := result.Value.(p2p.State)
	if !ok {
		return nil, P2PStateOutput{}, core.E("mcp.p2pState", "unexpected result type", nil)
	}
	return nil, P2PStateOutput{State: state}, nil
}

func (s *Subsystem) registerP2PTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "p2p_publish",
		Description: `Publish a P2P envelope over the configured transport. Example: {"topic":"display","route":"chat.sync","payload":{"message":"hello"}}`,
	}, s.p2pPublish)
	addTool(s, server, &mcp.Tool{
		Name:        "p2p_state",
		Description: "Inspect the configured P2P node state, listen address, and observed peers",
	}, s.p2pState)
}
