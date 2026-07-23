package chat

import (
	"bufio"
	"io"
	"slices"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/textutil"
)

type StreamCallbacks struct {
	OnStart          func(streamID string)
	OnToken          func(content string)
	OnThinkingStart  func(state ThinkingState)
	OnThinkingAppend func(content string)
	OnThinkingEnd    func(state ThinkingState)
	OnToolCall       func(call ToolCall)
	OnFinish         func(reason string)
}

type StreamRenderer struct {
	callbacks      StreamCallbacks
	now            func() time.Time
	content        streamBuilder
	thinking       streamBuilder
	thinkingState  ThinkingState
	toolCalls      map[int]*streamToolCall
	parsedToolArgs map[int]map[string]any
	toolOrder      []int
	streamID       string
	started        bool
	finishReason   string
}

type streamToolCall struct {
	ID        string
	Name      string
	Arguments streamBuilder
}

type streamBuilder struct {
	value string
}

func (b *streamBuilder) writeString(value string) {
	b.value += value
}

func (b *streamBuilder) string() string {
	return b.value
}

func (b *streamBuilder) len() int {
	return len(b.value)
}

type streamChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			Reasoning string `json:"reasoning"`
			Thinking  string `json:"thinking"`
			Thought   string `json:"thought"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func NewStreamRenderer(callbacks StreamCallbacks) *StreamRenderer {
	return &StreamRenderer{
		callbacks:      callbacks,
		now:            time.Now,
		toolCalls:      make(map[int]*streamToolCall),
		parsedToolArgs: make(map[int]map[string]any),
	}
}

func (r *StreamRenderer) Render(reader io.Reader) resultFailure {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 4096), 1024*1024)

	var dataLines []string
	flush := func() error {
		if len(dataLines) == 0 {
			return nil
		}
		payload := core.Join("\n", dataLines...)
		dataLines = nil
		return r.handleData(payload)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if core.Trim(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		if core.HasPrefix(line, "data:") {
			dataLines = append(dataLines, core.Trim(core.TrimPrefix(line, "data:")))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}

	r.completeThinking()
	for _, call := range r.ToolCalls() {
		if r.callbacks.OnToolCall != nil {
			r.callbacks.OnToolCall(call)
		}
	}
	if r.callbacks.OnFinish != nil {
		r.callbacks.OnFinish(r.finishReason)
	}
	return nil
}

func (r *StreamRenderer) handleData(payload string) resultFailure {
	if payload == "" {
		return nil
	}
	if payload == "[DONE]" {
		return nil
	}

	var chunk streamChunk
	result := core.JSONUnmarshalString(payload, &chunk)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("chat.StreamRenderer.handleData", "failed to decode stream chunk", nil)
	}
	if !r.started {
		r.started = true
		r.streamID = chunk.ID
		if r.callbacks.OnStart != nil {
			r.callbacks.OnStart(chunk.ID)
		}
	}

	for _, choice := range chunk.Choices {
		delta := choice.Delta
		thought := textutil.FirstNonEmpty(delta.Thinking, delta.Reasoning, delta.Thought)
		if thought != "" {
			r.appendThinking(thought)
		}
		if delta.Content != "" {
			r.completeThinking()
			r.content.writeString(delta.Content)
			if r.callbacks.OnToken != nil {
				r.callbacks.OnToken(delta.Content)
			}
		}
		for _, call := range delta.ToolCalls {
			r.appendToolCall(call.Index, call.ID, call.Function.Name, call.Function.Arguments)
		}
		if choice.FinishReason != "" {
			r.finishReason = choice.FinishReason
		}
	}
	return nil
}

func (r *StreamRenderer) appendThinking(content string) {
	if content == "" {
		return
	}
	if !r.thinkingState.Active {
		r.thinkingState = ThinkingState{
			Active:    true,
			StartedAt: r.now(),
		}
		if r.callbacks.OnThinkingStart != nil {
			r.callbacks.OnThinkingStart(r.thinkingState)
		}
	}
	r.thinking.writeString(content)
	if r.callbacks.OnThinkingAppend != nil {
		r.callbacks.OnThinkingAppend(content)
	}
}

func (r *StreamRenderer) completeThinking() {
	if !r.thinkingState.Active {
		return
	}
	r.thinkingState.Active = false
	r.thinkingState.Content = r.thinking.string()
	r.thinkingState.EndedAt = r.now()
	r.thinkingState.DurationMS = r.thinkingState.EndedAt.Sub(r.thinkingState.StartedAt).Milliseconds()
	if r.callbacks.OnThinkingEnd != nil {
		r.callbacks.OnThinkingEnd(r.thinkingState)
	}
}

func (r *StreamRenderer) appendToolCall(index int, id, name, arguments string) {
	call := r.toolCalls[index]
	if call == nil {
		call = &streamToolCall{}
		r.toolCalls[index] = call
		r.toolOrder = append(r.toolOrder, index)
	}
	if id != "" {
		call.ID = id
	}
	if name != "" {
		call.Name = name
	}
	if arguments != "" {
		call.Arguments.writeString(arguments)
		delete(r.parsedToolArgs, index)
	}
}

func (r *StreamRenderer) ToolCalls() []ToolCall {
	order := slices.Clone(r.toolOrder)
	slices.Sort(order)
	result := make([]ToolCall, 0, len(order))
	for _, index := range order {
		call := r.toolCalls[index]
		if call == nil {
			continue
		}
		arguments, ok := r.parsedToolArgs[index]
		if !ok {
			arguments = map[string]any{}
			raw := core.Trim(call.Arguments.string())
			if raw != "" {
				if decode := core.JSONUnmarshalString(raw, &arguments); !decode.OK {
					arguments = map[string]any{"raw": raw}
				}
			}
			r.parsedToolArgs[index] = cloneArguments(arguments)
		}
		result = append(result, ToolCall{
			ID:        call.ID,
			Name:      call.Name,
			Arguments: cloneArguments(arguments),
		})
	}
	return result
}

func (r *StreamRenderer) Thinking() *ThinkingState {
	if r.thinking.len() == 0 && !r.thinkingState.Active {
		return nil
	}
	state := r.thinkingState
	state.Content = r.thinking.string()
	return &state
}

func (r *StreamRenderer) Message(messageID, model string, createdAt time.Time) ChatMessage {
	return ChatMessage{
		ID:           messageID,
		Role:         "assistant",
		Content:      r.content.string(),
		CreatedAt:    createdAt,
		Model:        model,
		Thinking:     r.Thinking(),
		ToolCalls:    r.ToolCalls(),
		FinishReason: r.finishReason,
	}
}
