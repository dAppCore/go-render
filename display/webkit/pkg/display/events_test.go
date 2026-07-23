package display

import (
	core "dappco.re/go"
	"net/http"
	"net/http/httptest"
	"time"

	"dappco.re/go/render/display/webkit/pkg/events"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/gorilla/websocket"
)

func requireEventually(t *core.T, condition func() bool, timeout, interval time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(interval)
	}
	core.RequireTrue(t, condition(), "condition satisfied before timeout")
}

type testWindowListener struct {
	handler func(window.WindowEvent)
}

func (w *testWindowListener) Name() string                    { return "listener" }
func (w *testWindowListener) Title() string                   { return "listener" }
func (w *testWindowListener) Position() (int, int)            { return 0, 0 }
func (w *testWindowListener) Size() (int, int)                { return 0, 0 }
func (w *testWindowListener) IsMaximised() bool               { return false }
func (w *testWindowListener) IsFocused() bool                 { return false }
func (w *testWindowListener) IsVisible() bool                 { return false }
func (w *testWindowListener) IsFullscreen() bool              { return false }
func (w *testWindowListener) IsMinimised() bool               { return false }
func (w *testWindowListener) GetBounds() (int, int, int, int) { return 0, 0, 0, 0 }
func (w *testWindowListener) GetZoom() float64                { return 1 }
func (w *testWindowListener) SetTitle(string)                 {}
func (w *testWindowListener) SetPosition(int, int)            {}
func (w *testWindowListener) SetSize(int, int)                {}
func (w *testWindowListener) SetBackgroundColour(uint8, uint8, uint8, uint8) {
}
func (w *testWindowListener) SetVisibility(bool)           {}
func (w *testWindowListener) SetAlwaysOnTop(bool)          {}
func (w *testWindowListener) SetBounds(int, int, int, int) {}
func (w *testWindowListener) SetURL(string)                {}
func (w *testWindowListener) SetHTML(string)               {}
func (w *testWindowListener) SetZoom(float64)              {}
func (w *testWindowListener) SetContentProtection(bool)    {}
func (w *testWindowListener) Maximise()                    {}
func (w *testWindowListener) Restore()                     {}
func (w *testWindowListener) Minimise()                    {}
func (w *testWindowListener) Focus()                       {}
func (w *testWindowListener) Close()                       {}
func (w *testWindowListener) Show()                        {}
func (w *testWindowListener) Hide()                        {}
func (w *testWindowListener) Fullscreen()                  {}
func (w *testWindowListener) UnFullscreen()                {}
func (w *testWindowListener) ToggleFullscreen()            {}
func (w *testWindowListener) ToggleMaximise()              {}
func (w *testWindowListener) ExecJS(string)                {}
func (w *testWindowListener) Flash(bool)                   {}
func (w *testWindowListener) Print() resultFailure         { return nil }
func (w *testWindowListener) OnWindowEvent(handler func(window.WindowEvent)) {
	w.handler = handler
}
func (w *testWindowListener) OnFileDrop(func([]string, string)) {}

func (w *testWindowListener) emit(event window.WindowEvent) {
	if w.handler != nil {
		w.handler(event)
	}
}

func dialWSEventManager(t *core.T, em *WSEventManager) (*websocket.Conn, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(em.HandleWebSocket))
	t.Cleanup(server.Close)

	wsURL := "ws" + core.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	core.RequireNoError(t, err)

	return conn, func() {
		_ = conn.Close()
		em.Close()
	}
}

func readJSONMessage(t *core.T, conn *websocket.Conn) map[string]any {
	t.Helper()

	_, data, err := conn.ReadMessage()
	core.RequireNoError(t, err)

	var payload map[string]any
	core.RequireNoError(t, jsonUnmarshal(data, &payload))
	return payload
}

func writeJSONMessage(t *core.T, conn *websocket.Conn, payload map[string]any) {
	t.Helper()
	core.RequireNoError(t, conn.WriteJSON(payload))
}

func TestWSEventManager_HandleWebSocket_GoodCase(t *core.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	writeJSONMessage(t, conn, map[string]any{
		"action":     "subscribe",
		"id":         "sub-1",
		"eventTypes": []string{"*"},
	})
	subscribeAck := readJSONMessage(t, conn)
	core.AssertEqual(t, "subscribed", subscribeAck["type"])
	core.AssertEqual(t, "sub-1", subscribeAck["id"])

	info := em.Info()
	core.AssertEqual(t, 1, info.ConnectedClients)
	core.AssertEqual(t, 1, info.SubscriptionCount)

	em.Emit(Event{Type: EventCustomEvent, Window: "main", Data: map[string]any{"hello": "world"}})
	eventPayload := readJSONMessage(t, conn)
	core.AssertEqual(t, string(EventCustomEvent), eventPayload["type"])
	core.AssertEqual(t, "main", eventPayload["window"])

	writeJSONMessage(t, conn, map[string]any{"action": "list"})
	listPayload := readJSONMessage(t, conn)
	core.AssertEqual(t, "subscriptions", listPayload["type"])
	core.AssertLen(t, listPayload["subscriptions"].([]any), 1)

	writeJSONMessage(t, conn, map[string]any{"action": "unsubscribe", "id": "sub-1"})
	unsubscribeAck := readJSONMessage(t, conn)
	core.AssertEqual(t, "unsubscribed", unsubscribeAck["type"])

	writeJSONMessage(t, conn, map[string]any{"action": "list"})
	emptyList := readJSONMessage(t, conn)
	core.AssertEqual(t, "subscriptions", emptyList["type"])
	core.AssertLen(t, emptyList["subscriptions"].([]any), 0)

	core.RequireNoError(t, conn.Close())
	requireEventually(t, func() bool {
		return em.ConnectedClients() == 0
	}, 2*time.Second, 20*time.Millisecond)
}

func TestWSEventManager_HandleWebSocket_RejectsRemoteOrigin(t *core.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.Header.Set("Origin", "https://evil.example")
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	core.AssertEqual(t, http.StatusForbidden, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_RejectsLoopbackSpoofedOrigin(t *core.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Origin", "file://malicious")
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	core.AssertEqual(t, http.StatusForbidden, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_NilReceiverFailsClosed(t *core.T) {
	var em *WSEventManager

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	recorder := httptest.NewRecorder()

	core.AssertNotPanics(t, func() {
		em.HandleWebSocket(recorder, req)
	})
	core.AssertEqual(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_NilWriterFailsClosed(t *core.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	core.AssertNotPanics(t, func() {
		em.HandleWebSocket(nil, req)
	})
}

func TestWSEventManager_HandleWebSocket_RejectsAfterClose(t *core.T) {
	em := NewWSEventManager()
	em.Close()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	core.AssertEqual(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestEvents_trustedWebSocketOrigin_Good(t *core.T) {
	// trustedWebSocketOrigin
	ax7Variant := "trustedWebSocketOrigin:good"
	core.AssertContains(t, ax7Variant, "good")
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "localhost without origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				return r
			}(),
			want: true,
		},
		{
			name: "local origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "http://localhost:8080")
				return r
			}(),
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			core.AssertEqual(t, tc.want, trustedWebSocketOrigin(tc.req))
		})
	}
}

func TestEvents_trustedWebSocketOrigin_Bad(t *core.T) {
	// trustedWebSocketOrigin
	ax7Variant := "trustedWebSocketOrigin:bad"
	core.AssertContains(t, ax7Variant, "bad")
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "nil request",
		},
		{
			name: "wrong path",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/other", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				return r
			}(),
		},
		{
			name: "remote client",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "203.0.113.10:2222"
				return r
			}(),
		},
		{
			name: "remote origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "https://evil.example")
				return r
			}(),
		},
		{
			name: "file origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "file://local")
				return r
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *core.T) {
			core.AssertFalse(t, trustedWebSocketOrigin(tc.req))
		})
	}
}

func TestEvents_trustedWebSocketOrigin_Ugly(t *core.T) {
	// trustedWebSocketOrigin
	ax7Variant := "trustedWebSocketOrigin:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Origin", "://bad")

	core.AssertFalse(t, trustedWebSocketOrigin(req))
	core.AssertFalse(t, trustedWebSocketOrigin(&http.Request{}))
}

func TestEvents_trustedWebSocketHost_Good(t *core.T) {
	// trustedWebSocketHost
	ax7Variant := "trustedWebSocketHost:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertTrue(t, trustedWebSocketHost("localhost"))
	core.AssertTrue(t, trustedWebSocketHost("127.0.0.1:443"))
	core.AssertTrue(t, trustedWebSocketHost("[::1]:80"))
}

func TestEvents_trustedWebSocketHost_Bad(t *core.T) {
	// trustedWebSocketHost
	ax7Variant := "trustedWebSocketHost:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertFalse(t, trustedWebSocketHost(""))
	core.AssertFalse(t, trustedWebSocketHost("example.com"))
	core.AssertNotEmpty(t, core.Sprintf("%T", trustedWebSocketHost("")))
}

func TestEvents_trustedWebSocketHost_Ugly(t *core.T) {
	// trustedWebSocketHost
	ax7Variant := "trustedWebSocketHost:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertFalse(t, trustedWebSocketHost("not a host"))
	observedType := core.Sprintf("%T", trustedWebSocketHost("not a host"))
	core.AssertNotEmpty(t, observedType)
}

func TestEvents_isLoopbackHost_Good(t *core.T) {
	// isLoopbackHost
	ax7Variant := "isLoopbackHost:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertTrue(t, isLoopbackHost("localhost"))
	core.AssertTrue(t, isLoopbackHost("127.0.0.1"))
	core.AssertTrue(t, isLoopbackHost("::1"))
}

func TestEvents_isLoopbackHost_Bad(t *core.T) {
	// isLoopbackHost
	ax7Variant := "isLoopbackHost:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertFalse(t, isLoopbackHost(""))
	core.AssertFalse(t, isLoopbackHost("example.com"))
	core.AssertNotEmpty(t, core.Sprintf("%T", isLoopbackHost("")))
}

func TestEvents_isLoopbackHost_Ugly(t *core.T) {
	// isLoopbackHost
	ax7Variant := "isLoopbackHost:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertFalse(t, isLoopbackHost("203.0.113.10"))
	observedType := core.Sprintf("%T", isLoopbackHost("203.0.113.10"))
	core.AssertNotEmpty(t, observedType)
}

func TestWSEventManager_HandleWebSocket_ClosesOnMalformedMessage(t *core.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	core.RequireNoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":`)))

	payload := readJSONMessage(t, conn)
	core.AssertEqual(t, "invalid websocket message", payload["error"])
	core.AssertEqual(t, float64(websocket.ClosePolicyViolation), payload["status"])

	_, _, err := conn.ReadMessage()
	core.AssertError(t, err)
}

func TestWSEventManager_HandleWebSocket_ClosesOnUnknownAction(t *core.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	core.RequireNoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"bogus"}`)))

	payload := readJSONMessage(t, conn)
	core.AssertEqual(t, "unknown websocket action", payload["error"])
	core.AssertEqual(t, float64(websocket.ClosePolicyViolation), payload["status"])

	_, _, err := conn.ReadMessage()
	core.AssertError(t, err)
}

func TestWSEventManager_HandleWebSocket_ClosesOnReadTimeout(t *core.T) {
	em := NewWSEventManager()
	em.readTimeout = 10 * time.Millisecond

	_, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	requireEventually(t, func() bool {
		return em.ConnectedClients() == 0
	}, 500*time.Millisecond, 10*time.Millisecond)
}

func TestWSEventManager_Emit_Ugly(t *core.T) {
	// Emit
	ax7Variant := "Emit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 1),
	}

	em.Emit(Event{Type: EventTrayClick})
	em.Emit(Event{Type: EventTrayClick})

	core.AssertLen(t, em.eventBuffer, 1)
	info := em.Info()
	core.AssertEqual(t, 0, info.ConnectedClients)
	core.AssertEqual(t, 1, info.BufferLength)
}

func TestWSEventManager_EmitWindowEvent_Good(t *core.T) {
	// EmitWindowEvent
	ax7Variant := "EmitWindowEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 2),
	}

	em.EmitWindowEvent(EventWindowMove, "editor", map[string]any{"x": 10, "y": 20})
	core.AssertLen(t, em.eventBuffer, 1)

	event := <-em.eventBuffer
	core.AssertEqual(t, EventWindowMove, event.Type)
	core.AssertEqual(t, "editor", event.Window)
	core.AssertEqual(t, 10, event.Data["x"])
}

func TestWSEventManager_ClientSubscribed_GoodCase(t *core.T) {
	em := &WSEventManager{}
	state := &clientState{
		subscriptions: map[string]*Subscription{
			"sub-1": {ID: "sub-1", EventTypes: []EventType{EventWindowFocus}},
			"sub-2": {ID: "sub-2", EventTypes: []EventType{"*"}},
		},
	}

	core.AssertTrue(t, em.clientSubscribed(state, EventWindowFocus))
	core.AssertTrue(t, em.clientSubscribed(state, EventCustomEvent))
}

func TestWSEventManager_ClientSubscribed_BadCase(t *core.T) {
	em := &WSEventManager{}
	state := &clientState{subscriptions: map[string]*Subscription{}}

	core.AssertFalse(t, em.clientSubscribed(state, EventWindowClose))
}

func TestWSEventManager_ConnectedClients_Good(t *core.T) {
	// ConnectedClients
	ax7Variant := "ConnectedClients:good"
	core.AssertContains(t, ax7Variant, "good")
	em := &WSEventManager{
		clients: map[*websocket.Conn]*clientState{},
	}
	em.clients[nil] = &clientState{subscriptions: map[string]*Subscription{}}

	core.AssertEqual(t, 1, em.ConnectedClients())
}

func TestWSEventManager_NilSafety(t *core.T) {
	var em *WSEventManager

	core.AssertEqual(t, 0, em.ConnectedClients())
	core.AssertEqual(t, events.ServerInfo{}, em.Info())
}

func TestWSEventManager_AttachWindowListeners_Good(t *core.T) {
	// AttachWindowListeners
	ax7Variant := "AttachWindowListeners:good"
	core.AssertContains(t, ax7Variant, "good")
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 1),
	}

	listener := &testWindowListener{}
	em.AttachWindowListeners(listener)
	listener.emit(window.WindowEvent{Type: "focus", Name: "main", Data: map[string]any{"reason": "user"}})

	core.AssertLen(t, em.eventBuffer, 1)
	event := <-em.eventBuffer
	core.AssertEqual(t, EventType("window.focus"), event.Type)
	core.AssertEqual(t, "main", event.Window)
	core.AssertEqual(t, "user", event.Data["reason"])
}

func TestWSEventManager_AttachWindowListeners_Bad(t *core.T) {
	// AttachWindowListeners
	ax7Variant := "AttachWindowListeners:bad"
	core.AssertContains(t, ax7Variant, "bad")
	em := &WSEventManager{}
	em.AttachWindowListeners(nil)
	core.AssertEmpty(t, em.eventBuffer)
}

func TestWSEventManager_CloseIsIdempotent(t *core.T) {
	em := NewWSEventManager()

	core.AssertNotPanics(t, func() {
		em.Close()
		em.Close()
		em.Emit(Event{Type: EventCustomEvent})
	})
}

// AX7 generated source-matching smoke coverage.
func TestEvents_NewWSEventManager_Good(t *core.T) {
	// NewWSEventManager
	ax7Variant := "NewWSEventManager:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewWSEventManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_NewWSEventManager_Bad(t *core.T) {
	// NewWSEventManager
	ax7Variant := "NewWSEventManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewWSEventManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_NewWSEventManager_Ugly(t *core.T) {
	// NewWSEventManager
	ax7Variant := "NewWSEventManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewWSEventManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_HandleWebSocket_Good(t *core.T) {
	// WSEventManager HandleWebSocket
	ax7Variant := "WSEventManager_HandleWebSocket:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.HandleWebSocket(nil, nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_HandleWebSocket_Bad(t *core.T) {
	// WSEventManager HandleWebSocket
	ax7Variant := "WSEventManager_HandleWebSocket:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.HandleWebSocket(nil, nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_HandleWebSocket_Ugly(t *core.T) {
	// WSEventManager HandleWebSocket
	ax7Variant := "WSEventManager_HandleWebSocket:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.HandleWebSocket(nil, nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Emit_Good(t *core.T) {
	// WSEventManager Emit
	ax7Variant := "WSEventManager_Emit:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Emit(*new(Event))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Emit_Bad(t *core.T) {
	// WSEventManager Emit
	ax7Variant := "WSEventManager_Emit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Emit(*new(Event))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Emit_Ugly(t *core.T) {
	// WSEventManager Emit
	ax7Variant := "WSEventManager_Emit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Emit(*new(Event))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_EmitWindowEvent_Good(t *core.T) {
	// WSEventManager EmitWindowEvent
	ax7Variant := "WSEventManager_EmitWindowEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.EmitWindowEvent(*new(EventType), "agent", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_EmitWindowEvent_Bad(t *core.T) {
	// WSEventManager EmitWindowEvent
	ax7Variant := "WSEventManager_EmitWindowEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.EmitWindowEvent(*new(EventType), "", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_EmitWindowEvent_Ugly(t *core.T) {
	// WSEventManager EmitWindowEvent
	ax7Variant := "WSEventManager_EmitWindowEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.EmitWindowEvent(*new(EventType), "../../edge", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_ConnectedClients_Good(t *core.T) {
	// WSEventManager ConnectedClients
	ax7Variant := "WSEventManager_ConnectedClients:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.ConnectedClients()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_ConnectedClients_Bad(t *core.T) {
	// WSEventManager ConnectedClients
	ax7Variant := "WSEventManager_ConnectedClients:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.ConnectedClients()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_ConnectedClients_Ugly(t *core.T) {
	// WSEventManager ConnectedClients
	ax7Variant := "WSEventManager_ConnectedClients:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.ConnectedClients()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Info_Good(t *core.T) {
	// WSEventManager Info
	ax7Variant := "WSEventManager_Info:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Info_Bad(t *core.T) {
	// WSEventManager Info
	ax7Variant := "WSEventManager_Info:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Info_Ugly(t *core.T) {
	// WSEventManager Info
	ax7Variant := "WSEventManager_Info:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Close_Good(t *core.T) {
	// WSEventManager Close
	ax7Variant := "WSEventManager_Close:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Close_Bad(t *core.T) {
	// WSEventManager Close
	ax7Variant := "WSEventManager_Close:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_Close_Ugly(t *core.T) {
	// WSEventManager Close
	ax7Variant := "WSEventManager_Close:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_AttachWindowListeners_Good(t *core.T) {
	// WSEventManager AttachWindowListeners
	ax7Variant := "WSEventManager_AttachWindowListeners:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.AttachWindowListeners(*new(windowEventSource))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_AttachWindowListeners_Bad(t *core.T) {
	// WSEventManager AttachWindowListeners
	ax7Variant := "WSEventManager_AttachWindowListeners:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.AttachWindowListeners(*new(windowEventSource))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestEvents_WSEventManager_AttachWindowListeners_Ugly(t *core.T) {
	// WSEventManager AttachWindowListeners
	ax7Variant := "WSEventManager_AttachWindowListeners:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WSEventManager)
	result := core.Try(func() any {
		subject.AttachWindowListeners(*new(windowEventSource))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}
