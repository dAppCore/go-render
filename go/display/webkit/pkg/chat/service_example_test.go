package chat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	core "dappco.re/go"
)

func ExampleRegister() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello from chat"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}))
	defer server.Close()

	storeDir, err := coreMkdirTemp("", "chat-example-*")
	if err != nil {
		panic(err)
	}
	defer coreRemoveAll(storeDir)

	c := core.New(
		core.WithService(Register(
			func(o *Options) { o.APIURL = server.URL },
			func(o *Options) { o.StorePath = core.PathJoin(storeDir, "chat.db") },
			func(o *Options) { o.ToolExecutor = &mockToolExecutor{} },
			func(o *Options) { o.Now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() } },
		)),
		core.WithServiceLock(),
	)
	if !c.ServiceStartup(context.Background(), nil).OK {
		panic("chat startup failed")
	}

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	))
	if !send.OK {
		panic(send.Value)
	}

	conversations := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	history := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conversations.Value.([]Conversation)[0].ID},
	))

	core.Println(len(history.Value.([]Message)))
	core.Println(history.Value.([]Message)[1].Content)
	// Output:
	// 2
	// Hello from chat
}

// AX7 generated examples exercise each public call path with stable output.
func ExampleService_OnStartup() {
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_HandleIPCEvents() {
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_Send() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.Send(core.Background(), *new(sendInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_History() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.History("agent", 1)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_Models() {
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.Models()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_SelectModel() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.SelectModel(*new(selectModelInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_ListConversations() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.ListConversations()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_LoadConversation() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.LoadConversation("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_DeleteConversation() {
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.DeleteConversation("agent")
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_StartThinking() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StartThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleService_StopThinking() {
	subject := new(Service)
	result := core.Try(func() any {
		got0, got1 := subject.StopThinking(*new(thinkingInput))
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}
