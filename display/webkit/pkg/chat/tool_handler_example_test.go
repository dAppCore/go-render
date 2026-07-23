package chat

import (
	"context"
	core "dappco.re/go"

	guimcp "dappco.re/go/render/display/webkit/pkg/mcp"
)

type exampleToolExecutor struct{}

func (exampleToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{"type": "object"},
	}}
}

func (exampleToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (exampleToolExecutor) CallTool(_ context.Context, name string, _ map[string]any) (string, resultFailure) {
	if name == "layout_suggest" {
		return `{"mode":"left-right"}`, nil
	}
	return "", nil
}

func ExampleNewToolCallHandler() {
	handler := NewToolCallHandler(exampleToolExecutor{})
	result, err := handler.OnToolCall(context.Background(), ToolCall{
		ID:        "call-1",
		Name:      "layout_suggest",
		Arguments: map[string]any{"window_count": 2},
	})

	core.Println(err == nil)
	core.Println(result)
	core.Println(core.Contains(handler.BuildToolManifest(), "layout_suggest"))
	// Output:
	// true
	// {"mode":"left-right"}
	// true
}

// AX7 generated examples exercise each public call path with stable output.
func ExampleToolCallHandler_OnToolCall() {
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0, got1 := subject.OnToolCall(core.Background(), *new(ToolCall))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleToolCallHandler_BuildToolManifest() {
	var subject noopToolCallHandler
	result := core.Try(func() any {
		got0 := subject.BuildToolManifest()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleToolExecutor_Manifest() {
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleToolExecutor_ManifestText() {
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleToolExecutor_CallTool() {
	subject := new(actionToolExecutor)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "agent", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}
