//go:build compliance

package chat

import core "dappco.re/go"

func ExampleNewStreamRenderer() {
	core.Println("NewStreamRenderer")
	// Output:
	// NewStreamRenderer
}

func ExampleStreamRenderer_Render() {
	core.Println("StreamRenderer_Render")
	// Output:
	// StreamRenderer_Render
}

func ExampleStreamRenderer_ToolCalls() {
	core.Println("StreamRenderer_ToolCalls")
	// Output:
	// StreamRenderer_ToolCalls
}

func ExampleStreamRenderer_Thinking() {
	core.Println("StreamRenderer_Thinking")
	// Output:
	// StreamRenderer_Thinking
}

func ExampleStreamRenderer_Message() {
	core.Println("StreamRenderer_Message")
	// Output:
	// StreamRenderer_Message
}
