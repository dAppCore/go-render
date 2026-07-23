// pkg/mcp/tools_clipboard.go
package mcp

import (
	"context"
	"encoding/base64"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- clipboard_read ---

type ClipboardReadInput struct{}
type ClipboardReadOutput struct {
	Content string `json:"content"`
}

func (s *Subsystem) clipboardRead(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadInput) (*mcp.CallToolResult, ClipboardReadOutput, resultFailure) {
	r := s.core.QUERY(clipboard.QueryText{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardReadOutput{}, e
		}
		return nil, ClipboardReadOutput{}, core.E("mcp.clipboardRead", "clipboard query failed", nil)
	}
	content, ok := r.Value.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardReadOutput{}, core.E("mcp.clipboardRead", "unexpected result type", nil)
	}
	return nil, ClipboardReadOutput{Content: content.Text}, nil
}

// --- clipboard_write ---

type ClipboardWriteInput struct {
	Text string `json:"text"`
}
type ClipboardWriteOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardWrite(_ context.Context, _ *mcp.CallToolRequest, input ClipboardWriteInput) (*mcp.CallToolResult, ClipboardWriteOutput, resultFailure) {
	r := s.core.Action("clipboard.set_text").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: clipboard.TaskSetText{Text: input.Text}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardWriteOutput{}, e
		}
		return nil, ClipboardWriteOutput{}, nil
	}
	return nil, ClipboardWriteOutput{Success: true}, nil
}

// --- clipboard_has ---

type ClipboardHasInput struct{}
type ClipboardHasOutput struct {
	HasContent bool `json:"hasContent"`
}

func (s *Subsystem) clipboardHas(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardHasInput) (*mcp.CallToolResult, ClipboardHasOutput, resultFailure) {
	r := s.core.QUERY(clipboard.QueryText{})
	if !r.OK {
		return nil, ClipboardHasOutput{}, nil
	}
	content, ok := r.Value.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardHasOutput{}, core.E("mcp.clipboardHas", "unexpected result type", nil)
	}
	return nil, ClipboardHasOutput{HasContent: content.HasContent}, nil
}

// --- clipboard_clear ---

type ClipboardClearInput struct{}
type ClipboardClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardClear(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardClearInput) (*mcp.CallToolResult, ClipboardClearOutput, resultFailure) {
	r := s.core.Action("clipboard.clear").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardClearOutput{}, e
		}
		return nil, ClipboardClearOutput{}, nil
	}
	return nil, ClipboardClearOutput{Success: true}, nil
}

// --- clipboard_read_image ---

type ClipboardReadImageInput struct{}
type ClipboardReadImageOutput struct {
	Base64 string `json:"base64"`
}

func (s *Subsystem) clipboardReadImage(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadImageInput) (*mcp.CallToolResult, ClipboardReadImageOutput, resultFailure) {
	r := s.core.QUERY(clipboard.QueryImage{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardReadImageOutput{}, e
		}
		return nil, ClipboardReadImageOutput{}, nil
	}
	content, ok := r.Value.(clipboard.ImageContent)
	if !ok {
		return nil, ClipboardReadImageOutput{}, core.E("mcp.clipboardReadImage", "unexpected result type", nil)
	}
	if !content.HasImage {
		return nil, ClipboardReadImageOutput{}, nil
	}
	return nil, ClipboardReadImageOutput{Base64: base64.StdEncoding.EncodeToString(content.Data)}, nil
}

// --- clipboard_write_image ---

type ClipboardWriteImageInput struct {
	Base64 string `json:"base64"`
}
type ClipboardWriteImageOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardWriteImage(_ context.Context, _ *mcp.CallToolRequest, input ClipboardWriteImageInput) (*mcp.CallToolResult, ClipboardWriteImageOutput, resultFailure) {
	maxEncodedLen := ((clipboard.MaxImageBytes + 2) / 3) * 4
	if len(input.Base64) == 0 || len(input.Base64) > maxEncodedLen {
		return nil, ClipboardWriteImageOutput{}, core.E("mcp.clipboardWriteImage", "clipboard image exceeds maximum size", nil)
	}
	data, err := base64.StdEncoding.DecodeString(input.Base64)
	if err != nil {
		return nil, ClipboardWriteImageOutput{}, core.E("mcp.clipboardWriteImage", "invalid base64 image data", err)
	}
	r := s.core.Action("clipboard.set_image").Run(context.Background(), core.NewOptions(
		core.Option{Key: "data", Value: data},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardWriteImageOutput{}, e
		}
		return nil, ClipboardWriteImageOutput{}, nil
	}
	return nil, ClipboardWriteImageOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "clipboard_read", Description: "Read the current clipboard content"}, s.clipboardRead)
	addTool(s, server, &mcp.Tool{Name: "clipboard_write", Description: "Write text to the clipboard"}, s.clipboardWrite)
	addTool(s, server, &mcp.Tool{Name: "clipboard_has", Description: "Check if the clipboard has content"}, s.clipboardHas)
	addTool(s, server, &mcp.Tool{
		Name:        "clipboard_read_image",
		Description: `Read image data from the clipboard as base64 PNG bytes. Example: {}`,
	}, s.clipboardReadImage)
	addTool(s, server, &mcp.Tool{
		Name:        "clipboard_write_image",
		Description: `Write base64 image data to the clipboard. Example: {"base64":"iVBORw0KGgoAAA..."}`,
	}, s.clipboardWriteImage)
	addTool(s, server, &mcp.Tool{Name: "clipboard_clear", Description: "Clear the clipboard"}, s.clipboardClear)
}
