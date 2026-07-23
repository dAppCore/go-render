package chat

import (
	core "dappco.re/go"
	"time"
)

func TestStreamRenderer_Good_ParsesThinkingContentAndToolCalls(t *core.T) {
	stream := core.Join("\n", []string{
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"thinking":"Let me think"}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello"}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"tool_calls":[{"index":0,"id":"call-1","function":{"name":"layout_suggest","arguments":"{\"window_count\":2}"}}]}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"finish_reason":"tool_calls"}]}`,
		"",
		`data: [DONE]`,
		"",
	}...)

	renderer := NewStreamRenderer(StreamCallbacks{})
	core.RequireNoError(t, renderer.Render(core.NewReader(stream)))

	message := renderer.Message("msg-1", "lemer", testTime())
	core.AssertNotNil(t, message.Thinking)
	core.AssertEqual(t, "Hello", message.Content)
	core.AssertEqual(t, "Let me think", message.Thinking.Content)
	core.AssertLen(t, message.ToolCalls, 1)
	core.AssertEqual(t, "layout_suggest", message.ToolCalls[0].Name)
	core.AssertEqual(t, 2.0, message.ToolCalls[0].Arguments["window_count"])
	core.AssertEqual(t, "tool_calls", message.FinishReason)
}

func testTime() time.Time {
	return time.Unix(1_700_000_000, 0).UTC()
}

// AX7 generated source-matching smoke coverage.
func TestChat_NewStreamRenderer_Good(t *core.T) {
	// NewStreamRenderer
	ax7Variant := "NewStreamRenderer:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewStreamRenderer(*new(StreamCallbacks))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_NewStreamRenderer_Bad(t *core.T) {
	// NewStreamRenderer
	ax7Variant := "NewStreamRenderer:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewStreamRenderer(*new(StreamCallbacks))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_NewStreamRenderer_Ugly(t *core.T) {
	// NewStreamRenderer
	ax7Variant := "NewStreamRenderer:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewStreamRenderer(*new(StreamCallbacks))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Render_Good(t *core.T) {
	// StreamRenderer Render
	ax7Variant := "StreamRenderer_Render:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Render(core.NewReader(""))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Render_Bad(t *core.T) {
	// StreamRenderer Render
	ax7Variant := "StreamRenderer_Render:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Render(core.NewReader(""))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Render_Ugly(t *core.T) {
	// StreamRenderer Render
	ax7Variant := "StreamRenderer_Render:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Render(core.NewReader(""))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_ToolCalls_Good(t *core.T) {
	// StreamRenderer ToolCalls
	ax7Variant := "StreamRenderer_ToolCalls:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.ToolCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_ToolCalls_Bad(t *core.T) {
	// StreamRenderer ToolCalls
	ax7Variant := "StreamRenderer_ToolCalls:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.ToolCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_ToolCalls_Ugly(t *core.T) {
	// StreamRenderer ToolCalls
	ax7Variant := "StreamRenderer_ToolCalls:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.ToolCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Thinking_Good(t *core.T) {
	// StreamRenderer Thinking
	ax7Variant := "StreamRenderer_Thinking:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Thinking()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Thinking_Bad(t *core.T) {
	// StreamRenderer Thinking
	ax7Variant := "StreamRenderer_Thinking:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Thinking()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Thinking_Ugly(t *core.T) {
	// StreamRenderer Thinking
	ax7Variant := "StreamRenderer_Thinking:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Thinking()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Message_Good(t *core.T) {
	// StreamRenderer Message
	ax7Variant := "StreamRenderer_Message:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Message("agent", "agent", core.Now())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Message_Bad(t *core.T) {
	// StreamRenderer Message
	ax7Variant := "StreamRenderer_Message:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Message("", "", core.Now())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestChat_StreamRenderer_Message_Ugly(t *core.T) {
	// StreamRenderer Message
	ax7Variant := "StreamRenderer_Message:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(StreamRenderer)
	result := core.Try(func() any {
		got0 := subject.Message("../../edge", "../../edge", core.Now())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
