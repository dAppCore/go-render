// pkg/mcp/tools_events.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/events"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- event_emit ---

type EventEmitInput struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}
type EventEmitOutput struct {
	Cancelled bool `json:"cancelled"`
}

func (s *Subsystem) eventEmit(_ context.Context, _ *mcp.CallToolRequest, input EventEmitInput) (*mcp.CallToolResult, EventEmitOutput, resultFailure) {
	r := s.core.Action("events.emit").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskEmit{Name: input.Name, Data: input.Data}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventEmitOutput{}, e
		}
		return nil, EventEmitOutput{}, nil
	}
	cancelled, ok := r.Value.(bool)
	if !ok {
		return nil, EventEmitOutput{}, core.E("mcp.eventEmit", "unexpected result type", nil)
	}
	return nil, EventEmitOutput{Cancelled: cancelled}, nil
}

// --- event_on ---

type EventOnInput struct {
	Name string `json:"name"`
}
type EventOnOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) eventOn(_ context.Context, _ *mcp.CallToolRequest, input EventOnInput) (*mcp.CallToolResult, EventOnOutput, resultFailure) {
	r := s.core.Action("events.on").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskOn{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventOnOutput{}, e
		}
		return nil, EventOnOutput{}, nil
	}
	return nil, EventOnOutput{Success: true}, nil
}

// --- event_subscribe ---

type EventSubscribeInput struct {
	Name string `json:"name"`
}
type EventSubscribeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) eventSubscribe(ctx context.Context, req *mcp.CallToolRequest, input EventSubscribeInput) (*mcp.CallToolResult, EventSubscribeOutput, resultFailure) {
	result, output, err := s.eventOn(ctx, req, EventOnInput{Name: input.Name})
	if err != nil {
		return nil, EventSubscribeOutput{}, err
	}
	if result != nil {
		return result, EventSubscribeOutput{}, nil
	}
	return nil, EventSubscribeOutput{Success: output.Success}, nil
}

// --- event_off ---

type EventOffInput struct {
	Name string `json:"name"`
}
type EventOffOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) eventOff(_ context.Context, _ *mcp.CallToolRequest, input EventOffInput) (*mcp.CallToolResult, EventOffOutput, resultFailure) {
	r := s.core.Action("events.off").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskOff{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventOffOutput{}, e
		}
		return nil, EventOffOutput{}, nil
	}
	return nil, EventOffOutput{Success: true}, nil
}

// --- event_unsubscribe ---

type EventUnsubscribeInput struct {
	Name string `json:"name"`
}
type EventUnsubscribeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) eventUnsubscribe(ctx context.Context, req *mcp.CallToolRequest, input EventUnsubscribeInput) (*mcp.CallToolResult, EventUnsubscribeOutput, resultFailure) {
	result, output, err := s.eventOff(ctx, req, EventOffInput{Name: input.Name})
	if err != nil {
		return nil, EventUnsubscribeOutput{}, err
	}
	if result != nil {
		return result, EventUnsubscribeOutput{}, nil
	}
	return nil, EventUnsubscribeOutput{Success: output.Success}, nil
}

// --- event_list ---

type EventListInput struct{}
type EventListOutput struct {
	Listeners []events.ListenerInfo `json:"listeners"`
}

func (s *Subsystem) eventList(_ context.Context, _ *mcp.CallToolRequest, _ EventListInput) (*mcp.CallToolResult, EventListOutput, resultFailure) {
	r := s.core.QUERY(events.QueryListeners{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventListOutput{}, e
		}
		return nil, EventListOutput{}, nil
	}
	listenerInfos, ok := r.Value.([]events.ListenerInfo)
	if !ok {
		return nil, EventListOutput{}, core.E("mcp.eventList", "unexpected result type", nil)
	}
	return nil, EventListOutput{Listeners: listenerInfos}, nil
}

// --- event_info ---

type EventInfoInput struct{}
type EventInfoOutput struct {
	Info events.ServerInfo `json:"info"`
}

func (s *Subsystem) eventInfo(_ context.Context, _ *mcp.CallToolRequest, _ EventInfoInput) (*mcp.CallToolResult, EventInfoOutput, resultFailure) {
	r := s.core.QUERY(events.QueryServerInfo{})
	if !r.OK {
		return nil, EventInfoOutput{}, nil
	}
	info, ok := r.Value.(events.ServerInfo)
	if !ok {
		return nil, EventInfoOutput{}, core.E("mcp.eventInfo", "unexpected result type", nil)
	}
	return nil, EventInfoOutput{Info: info}, nil
}

// --- Registration ---

func (s *Subsystem) registerEventsTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "event_emit", Description: "Fire a named custom event with optional data"}, s.eventEmit)
	addTool(s, server, &mcp.Tool{Name: "event_on", Description: "Register a listener for a named custom event"}, s.eventOn)
	addTool(s, server, &mcp.Tool{Name: "event_subscribe", Description: "Register a listener for a named custom event"}, s.eventSubscribe)
	addTool(s, server, &mcp.Tool{Name: "event_off", Description: "Remove all listeners for a named custom event"}, s.eventOff)
	addTool(s, server, &mcp.Tool{Name: "event_unsubscribe", Description: "Remove all listeners for a named custom event"}, s.eventUnsubscribe)
	addTool(s, server, &mcp.Tool{Name: "event_list", Description: "Query all registered event listeners"}, s.eventList)
	addTool(s, server, &mcp.Tool{Name: "event_info", Description: "Get WebSocket event server information"}, s.eventInfo)
}
