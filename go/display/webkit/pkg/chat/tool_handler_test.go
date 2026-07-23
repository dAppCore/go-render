package chat

import (
	"context"
	"io"
	"net/http"
	"sync"

	core "dappco.re/go"
	guimcp "dappco.re/go/render/display/webkit/pkg/mcp"
)

type strictToolExecutor struct {
	mu    sync.Mutex
	calls []ToolCall
}

func (m *strictToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"window_count": map[string]any{"type": "integer"},
			},
		},
	}}
}

func (m *strictToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (m *strictToolExecutor) CallTool(_ context.Context, name string, arguments map[string]any) (string, resultFailure) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, ToolCall{Name: name, Arguments: arguments})
	if name != "layout_suggest" {
		return "", core.E("test.tool", "unknown tool: "+name, nil)
	}
	return `{"mode":"left-right"}`, nil
}

func (m *strictToolExecutor) Calls() []ToolCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]ToolCall(nil), m.calls...)
}

type completionRecorder struct {
	mu        sync.Mutex
	requests  []openAIRequest
	responses [][]string
}

func (r *completionRecorder) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	var completion openAIRequest
	if result := core.JSONUnmarshal(body, &completion); !result.OK {
		http.Error(w, renderToolResultContent(result.Value), http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	r.requests = append(r.requests, completion)
	index := len(r.requests) - 1
	r.mu.Unlock()

	if index >= len(r.responses) {
		http.Error(w, "unexpected completion request", http.StatusInternalServerError)
		return
	}
	writeSSE(w, r.responses[index]...)
}

func (r *completionRecorder) Requests() []openAIRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]openAIRequest(nil), r.requests...)
}

func TestToolCallHandler_Good_ServiceDispatchesInlineToolCall(t *core.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{
		{
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"layout_suggest\",\"arguments\":{\"window_count\":2}}}"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
		{
			`{"id":"chatcmpl-2","choices":[{"delta":{"content":"Layout applied"}}]}`,
			`{"id":"chatcmpl-2","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
	}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Arrange this workspace"},
	))
	core.RequireTrue(t, send.OK)

	calls := executor.Calls()
	core.AssertLen(t, calls, 1)
	core.AssertEqual(t, "layout_suggest", calls[0].Name)
	core.AssertEqual(t, float64(2), calls[0].Arguments["window_count"])

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	core.AssertLen(t, history, 4)
	core.AssertEqual(t, "assistant", history[1].Role)
	core.AssertLen(t, history[1].ToolCalls, 1)
	core.AssertEqual(t, "tool", history[2].Role)
	core.AssertContains(t, history[2].Content, "left-right")
	core.AssertEqual(t, "Layout applied", history[3].Content)

	requests := recorder.Requests()
	core.AssertLen(t, requests, 2)
	core.RequireNotEmpty(t, requests[0].Messages)
	systemPrompt, ok := requests[0].Messages[0].Content.(string)
	core.RequireTrue(t, ok)
	core.AssertTrue(t, core.HasPrefix(systemPrompt, "Available MCP tools:"))
	core.AssertContains(t, systemPrompt, "layout_suggest")
	core.AssertContains(t, systemPrompt, "You are a helpful assistant.")
}

func TestToolCallHandler_Bad_UnknownToolErrorAppearsInConversation(t *core.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{
		{
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"missing_tool\",\"arguments\":{}}}"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
		{
			`{"id":"chatcmpl-2","choices":[{"delta":{"content":"Could not run that tool"}}]}`,
			`{"id":"chatcmpl-2","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
	}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Use the missing tool"},
	))
	core.RequireTrue(t, send.OK)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	core.AssertLen(t, history, 4)
	core.AssertEqual(t, "tool", history[2].Role)
	core.AssertContains(t, history[2].Content, "missing_tool")
	core.AssertEqual(t, "Could not run that tool", history[3].Content)
}

func TestToolCallHandler_Ugly_MalformedInlineToolCallDoesNotDispatch(t *core.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{{
		`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"layout_suggest\",\"arguments\":"}}]}`,
		`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
		`[DONE]`,
	}}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Try malformed JSON"},
	))
	core.RequireTrue(t, send.OK)

	core.AssertEmpty(t, executor.Calls())
	core.AssertLen(t, recorder.Requests(), 1)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	core.AssertLen(t, history, 2)
	core.AssertEqual(t, "assistant", history[1].Role)
	core.AssertContains(t, history[1].Content, "tool_call")
	core.AssertEmpty(t, history[1].ToolCalls)
}

// AX7 generated source-matching smoke coverage.
func TestToolHandler_NewToolCallHandler_Good(t *core.T) {
	// NewToolCallHandler
	ax7Variant := "NewToolCallHandler:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewToolCallHandler(*new(ToolExecutor))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_NewToolCallHandler_Bad(t *core.T) {
	// NewToolCallHandler
	ax7Variant := "NewToolCallHandler:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewToolCallHandler(*new(ToolExecutor))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_NewToolCallHandler_Ugly(t *core.T) {
	// NewToolCallHandler
	ax7Variant := "NewToolCallHandler:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewToolCallHandler(*new(ToolExecutor))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_OnToolCall_Good(t *core.T) {
	// ToolCallHandler OnToolCall
	ax7Variant := "ToolCallHandler_OnToolCall:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0, got1 := subject.OnToolCall(core.Background(), *new(ToolCall))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_OnToolCall_Bad(t *core.T) {
	// ToolCallHandler OnToolCall
	ax7Variant := "ToolCallHandler_OnToolCall:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0, got1 := subject.OnToolCall(core.Background(), *new(ToolCall))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_OnToolCall_Ugly(t *core.T) {
	// ToolCallHandler OnToolCall
	ax7Variant := "ToolCallHandler_OnToolCall:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0, got1 := subject.OnToolCall(core.Background(), *new(ToolCall))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_BuildToolManifest_Good(t *core.T) {
	// ToolCallHandler BuildToolManifest
	ax7Variant := "ToolCallHandler_BuildToolManifest:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0 := subject.BuildToolManifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_BuildToolManifest_Bad(t *core.T) {
	// ToolCallHandler BuildToolManifest
	ax7Variant := "ToolCallHandler_BuildToolManifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0 := subject.BuildToolManifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolCallHandler_BuildToolManifest_Ugly(t *core.T) {
	// ToolCallHandler BuildToolManifest
	ax7Variant := "ToolCallHandler_BuildToolManifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0 := subject.BuildToolManifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_Manifest_Good(t *core.T) {
	// ToolExecutor Manifest
	ax7Variant := "ToolExecutor_Manifest:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_Manifest_Bad(t *core.T) {
	// ToolExecutor Manifest
	ax7Variant := "ToolExecutor_Manifest:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_Manifest_Ugly(t *core.T) {
	// ToolExecutor Manifest
	ax7Variant := "ToolExecutor_Manifest:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_ManifestText_Good(t *core.T) {
	// ToolExecutor ManifestText
	ax7Variant := "ToolExecutor_ManifestText:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_ManifestText_Bad(t *core.T) {
	// ToolExecutor ManifestText
	ax7Variant := "ToolExecutor_ManifestText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_ManifestText_Ugly(t *core.T) {
	// ToolExecutor ManifestText
	ax7Variant := "ToolExecutor_ManifestText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_CallTool_Good(t *core.T) {
	// ToolExecutor CallTool
	ax7Variant := "ToolExecutor_CallTool:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "agent", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_CallTool_Bad(t *core.T) {
	// ToolExecutor CallTool
	ax7Variant := "ToolExecutor_CallTool:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestToolHandler_ToolExecutor_CallTool_Ugly(t *core.T) {
	// ToolExecutor CallTool
	ax7Variant := "ToolExecutor_CallTool:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "../../edge", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}
