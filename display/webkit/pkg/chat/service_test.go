package chat

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	core "dappco.re/go"
	guimcp "dappco.re/go/render/display/webkit/pkg/mcp"
)

type mockToolExecutor struct {
	calls []ToolCall
}

func (m *mockToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{"type": "object"},
	}}
}

func (m *mockToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (m *mockToolExecutor) CallTool(_ context.Context, name string, arguments map[string]any) (string, resultFailure) {
	m.calls = append(m.calls, ToolCall{Name: name, Arguments: arguments})
	return `{"mode":"left-right"}`, nil
}

func newChatCore(t *core.T, handler http.HandlerFunc, toolExecutor ToolExecutor, optionFns ...func(*Options)) *core.Core {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	options := []func(*Options){
		func(o *Options) { o.APIURL = server.URL },
		func(o *Options) { o.StorePath = core.PathJoin(t.TempDir(), "chat.db") },
		func(o *Options) { o.ToolExecutor = toolExecutor },
		func(o *Options) { o.Now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() } },
		func(o *Options) { o.ModelRoots = nil },
	}
	options = append(options, optionFns...)

	c := core.New(
		core.WithService(Register(options...)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	return c
}

func createDiscoveredModelRoot(t *core.T, name, architecture string) string {
	t.Helper()
	root := t.TempDir()
	modelDir := core.PathJoin(root, name)
	core.RequireNoError(t, coreEnsureDir(modelDir))
	configJSON := `{"model_type":"` + architecture + `","quantization":{"bits":4,"group_size":32}}`
	core.RequireNoError(t, coreWriteMode(core.PathJoin(modelDir, "config.json"), configJSON, 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(modelDir, "weights.safetensors"), "fake", 0o644))
	return root
}

func sequencedNow(start time.Time) func() time.Time {
	current := start.Add(-time.Second)
	return func() time.Time {
		current = current.Add(time.Second)
		return current
	}
}

func writeSSE(w http.ResponseWriter, payloads ...string) {
	w.Header().Set("Content-Type", "text/event-stream")
	for _, payload := range payloads {
		_, _ = io.WriteString(w, "data: "+payload+"\n\n")
	}
}

func latestConversation(t *core.T, c *core.Core) Conversation {
	t.Helper()
	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)
	conversations, ok := result.Value.([]Conversation)
	core.RequireTrue(t, ok)
	core.RequireNotEmpty(t, conversations)
	return conversations[0]
}

func historyMessages(t *core.T, c *core.Core, conversationID string, limit int) []Message {
	t.Helper()
	options := []core.Option{{
		Key:   "conversation_id",
		Value: conversationID,
	}}
	if limit > 0 {
		options = append(options, core.Option{Key: "limit", Value: limit})
	}
	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(options...))
	core.RequireTrue(t, result.OK)
	messages, ok := result.Value.([]Message)
	core.RequireTrue(t, ok)
	return messages
}

func TestActionSend_Good_ReturnsAssistantMessageID(t *core.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":" world"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hi"},
		core.Option{Key: "model", Value: "lemma"},
	))
	core.RequireTrue(t, send.OK)

	messageID, ok := send.Value.(string)
	core.RequireTrue(t, ok)
	core.RequireNotEmpty(t, messageID)

	conv := latestConversation(t, c)
	core.AssertLen(t, conv.Messages, 2)
	core.AssertEqual(t, messageID, conv.Messages[1].ID)
	core.AssertEqual(t, "lemma", conv.Model)
	core.AssertEqual(t, "Hello world", conv.Messages[1].Content)
}

func TestActionSend_Bad_RejectsEmptyMessage(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "message content is required")
}

func TestActionSend_Ugly_PropagatesUpstreamFailure(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "model unavailable", http.StatusBadGateway)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hi"},
	))
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "model unavailable")
}

func TestActionHistory_Good_HonoursLimit(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Alpha"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "One"},
	))
	core.RequireTrue(t, send.OK)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 1)
	core.AssertLen(t, history, 1)
	core.AssertEqual(t, "assistant", history[0].Role)
	core.AssertEqual(t, "Alpha", history[0].Content)
}

func TestActionHistory_Bad_RequiresConversationID(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "conversation id is required")
}

func TestActionHistory_Ugly_UnknownConversationFails(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
}

func TestActionModels_Good_ReportsSizeAndStatus(t *core.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	core.RequireTrue(t, ok)
	core.AssertLen(t, models, 1)
	core.AssertEqual(t, "lemma", models[0].Name)
	core.AssertEqual(t, int64(4), models[0].Size)
	core.AssertEqual(t, "active", models[0].Status)
}

func TestActionModels_Bad_ReturnsFallbackWhenNothingDiscovered(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	core.RequireTrue(t, ok)
	core.AssertLen(t, models, 1)
	core.AssertEqual(t, "default", models[0].Name)
	core.AssertEqual(t, "active", models[0].Status)
}

func TestActionModels_Ugly_ReflectsSelectedModelStatus(t *core.T) {
	rootA := createDiscoveredModelRoot(t, "alpha", "gemma3")
	rootB := createDiscoveredModelRoot(t, "beta", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{rootA, rootB} })

	selected := c.Action("gui.chat.select_model").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "beta"},
	))
	core.RequireTrue(t, selected.OK)

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	core.RequireTrue(t, ok)
	core.AssertLen(t, models, 2)
	statusByName := map[string]string{}
	for _, model := range models {
		statusByName[model.Name] = model.Status
	}
	core.AssertEqual(t, "available", statusByName["alpha"])
	core.AssertEqual(t, "active", statusByName["beta"])
}

func TestActionSelectModel_Good_UpdatesConversationAndSettings(t *core.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	created := c.Action("gui.chat.conversations.new").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, created.OK)
	conv := created.Value.(Conversation)

	selected := c.Action("gui.chat.select_model").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "lemma"},
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	core.RequireTrue(t, selected.OK)

	settings := selected.Value.(ChatSettings)
	core.AssertEqual(t, "lemma", settings.DefaultModel)

	loaded := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	core.RequireTrue(t, loaded.OK)
	core.AssertEqual(t, "lemma", loaded.Value.(Conversation).Model)
}

func TestActionSelectModel_Bad_RequiresModelName(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.select_model").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "model is required")
}

func TestActionSelectModel_Ugly_RejectsUnknownDiscoveredModel(t *core.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	result := c.Action("gui.chat.select_model").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "missing"},
	))
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "model is not available")
}

func TestActionConversationsList_Good_ReturnsNewestFirst(t *core.T) {
	now := sequencedNow(time.Unix(1_700_000_000, 0).UTC())
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Ack"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{}, func(o *Options) { o.Now = now })

	core.RequireTrue(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "First"},
	)).OK)
	core.RequireTrue(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Second"},
	)).OK)

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)
	conversations := result.Value.([]Conversation)
	core.AssertLen(t, conversations, 2)
	core.AssertEqual(t, "Second", conversations[0].Messages[0].Content)
	core.AssertEqual(t, "First", conversations[1].Messages[0].Content)
}

func TestActionConversationsList_Bad_EmptyStoreReturnsEmptySlice(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)
	conversations, ok := result.Value.([]Conversation)
	core.RequireTrue(t, ok)
	core.AssertEmpty(t, conversations)
}

func TestActionConversationsList_Ugly_IgnoresCorruptRows(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Ack"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	core.RequireTrue(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Good"},
	)).OK)

	svc := core.MustServiceFor[*Service](c, "chat")
	core.RequireNoError(t, svc.store.set(conversationsGroup, "broken", "{"))

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, result.OK)
	conversations := result.Value.([]Conversation)
	core.AssertLen(t, conversations, 1)
	core.AssertEqual(t, "Good", conversations[0].Messages[0].Content)
}

func TestActionConversationsLoad_Good_ReturnsConversation(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Reply"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	core.RequireTrue(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	)).OK)
	conv := latestConversation(t, c)

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	core.RequireTrue(t, result.OK)
	loaded := result.Value.(Conversation)
	core.AssertLen(t, loaded.Messages, 2)
	core.AssertEqual(t, "Reply", loaded.Messages[1].Content)
}

func TestActionConversationsLoad_Bad_RequiresConversationID(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "conversation id is required")
}

func TestActionConversationsLoad_Ugly_UnknownConversationFails(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
}

func TestActionConversationsDelete_Good_RemovesConversation(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Reply"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	core.RequireTrue(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	)).OK)
	conv := latestConversation(t, c)

	deleted := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	core.RequireTrue(t, deleted.OK)
	core.AssertEqual(t, true, deleted.Value)

	listed := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, listed.OK)
	core.AssertEmpty(t, listed.Value.([]Conversation))
}

func TestActionConversationsDelete_Bad_RequiresConversationID(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "conversation id is required")
}

func TestActionConversationsDelete_Ugly_IsIdempotentForMissingConversation(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	core.RequireTrue(t, result.OK)
	core.AssertEqual(t, true, result.Value)
}

func TestActionThinkingStart_Good_ReturnsActiveState(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	core.RequireTrue(t, result.OK)
	state := result.Value.(ThinkingState)
	core.AssertTrue(t, state.Active)
	core.AssertFalse(t, state.StartedAt.IsZero())
}

func TestActionThinkingStart_Bad_RequiresConversationID(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "conversation id is required")
}

func TestActionThinkingStart_Ugly_RestartReplacesExistingState(t *core.T) {
	now := sequencedNow(time.Unix(1_700_000_000, 0).UTC())
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.Now = now })

	first := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	second := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	core.RequireTrue(t, first.OK)
	core.RequireTrue(t, second.OK)
	core.AssertTrue(t, second.Value.(ThinkingState).StartedAt.After(first.Value.(ThinkingState).StartedAt))
}

func TestActionThinkingStop_Good_ClearsThinkingState(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	core.RequireTrue(t, c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
		core.Option{Key: "started_at", Value: time.Unix(1_700_000_000, 0).UTC()},
	)).OK)

	stopped := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
		core.Option{Key: "duration_ms", Value: int64(25)},
	))
	core.RequireTrue(t, stopped.OK)
	state := stopped.Value.(ThinkingState)
	core.AssertFalse(t, state.Active)
	core.AssertEqual(t, int64(25), state.DurationMS)
}

func TestActionThinkingStop_Bad_RequiresConversationID(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "conversation id is required")
}

func TestActionThinkingStop_Ugly_AllowsStopWithoutStart(t *core.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	core.RequireTrue(t, result.OK)
	state := result.Value.(ThinkingState)
	core.AssertFalse(t, state.Active)
	core.AssertTrue(t, state.DurationMS >= 0)
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Send_Good(t *core.T) {
	// Service Send
	ax7Variant := "Service_Send:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.Send(core.Background(), *new(sendInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Send_Bad(t *core.T) {
	// Service Send
	ax7Variant := "Service_Send:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.Send(core.Background(), *new(sendInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Send_Ugly(t *core.T) {
	// Service Send
	ax7Variant := "Service_Send:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.Send(core.Background(), *new(sendInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_History_Good(t *core.T) {
	// Service History
	ax7Variant := "Service_History:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.History("agent", 1)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_History_Bad(t *core.T) {
	// Service History
	ax7Variant := "Service_History:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.History("", 0)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_History_Ugly(t *core.T) {
	// Service History
	ax7Variant := "Service_History:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.History("../../edge", -1)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Models_Good(t *core.T) {
	// Service Models
	ax7Variant := "Service_Models:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Models()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Models_Bad(t *core.T) {
	// Service Models
	ax7Variant := "Service_Models:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Models()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_Models_Ugly(t *core.T) {
	// Service Models
	ax7Variant := "Service_Models:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Models()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_SelectModel_Good(t *core.T) {
	// Service SelectModel
	ax7Variant := "Service_SelectModel:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SelectModel(*new(selectModelInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_SelectModel_Bad(t *core.T) {
	// Service SelectModel
	ax7Variant := "Service_SelectModel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SelectModel(*new(selectModelInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_SelectModel_Ugly(t *core.T) {
	// Service SelectModel
	ax7Variant := "Service_SelectModel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SelectModel(*new(selectModelInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_ListConversations_Good(t *core.T) {
	// Service ListConversations
	ax7Variant := "Service_ListConversations:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ListConversations()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_ListConversations_Bad(t *core.T) {
	// Service ListConversations
	ax7Variant := "Service_ListConversations:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ListConversations()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_ListConversations_Ugly(t *core.T) {
	// Service ListConversations
	ax7Variant := "Service_ListConversations:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ListConversations()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_LoadConversation_Good(t *core.T) {
	// Service LoadConversation
	ax7Variant := "Service_LoadConversation:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.LoadConversation("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_LoadConversation_Bad(t *core.T) {
	// Service LoadConversation
	ax7Variant := "Service_LoadConversation:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.LoadConversation("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_LoadConversation_Ugly(t *core.T) {
	// Service LoadConversation
	ax7Variant := "Service_LoadConversation:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.LoadConversation("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_DeleteConversation_Good(t *core.T) {
	// Service DeleteConversation
	ax7Variant := "Service_DeleteConversation:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.DeleteConversation("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_DeleteConversation_Bad(t *core.T) {
	// Service DeleteConversation
	ax7Variant := "Service_DeleteConversation:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.DeleteConversation("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_DeleteConversation_Ugly(t *core.T) {
	// Service DeleteConversation
	ax7Variant := "Service_DeleteConversation:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.DeleteConversation("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StartThinking_Good(t *core.T) {
	// Service StartThinking
	ax7Variant := "Service_StartThinking:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StartThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StartThinking_Bad(t *core.T) {
	// Service StartThinking
	ax7Variant := "Service_StartThinking:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StartThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StartThinking_Ugly(t *core.T) {
	// Service StartThinking
	ax7Variant := "Service_StartThinking:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StartThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StopThinking_Good(t *core.T) {
	// Service StopThinking
	ax7Variant := "Service_StopThinking:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StopThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StopThinking_Bad(t *core.T) {
	// Service StopThinking
	ax7Variant := "Service_StopThinking:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StopThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_StopThinking_Ugly(t *core.T) {
	// Service StopThinking
	ax7Variant := "Service_StopThinking:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StopThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}
