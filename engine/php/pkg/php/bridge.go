package php

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"

	core "dappco.re/go"
)

// BridgeHandler is the interface that the host application implements to
// respond to PHP-initiated requests via the native bridge.
type BridgeHandler interface {
	// HandleBridgeCall processes a named bridge call with JSON args.
	// Returns a JSON-serializable response.
	HandleBridgeCall(method string, args json.RawMessage) (any, error)
}

// Bridge provides a localhost HTTP API that PHP code can call
// to access native desktop capabilities (file dialogs, notifications, etc.).
//
// Livewire renders server-side in PHP, so it can't call Wails bindings
// (window.go.*) directly. Instead, PHP makes HTTP requests to this bridge.
// The bridge port is injected into Laravel's .env as NATIVE_BRIDGE_URL.
type Bridge struct {
	server  *http.Server
	port    int
	handler BridgeHandler
}

// NewBridge creates and starts the bridge on a random available port.
// The handler processes incoming PHP requests via HandleBridgeCall.
func NewBridge(handler BridgeHandler) (*Bridge, error) { // Result boundary
	mux := http.NewServeMux()
	bridge := &Bridge{handler: handler}

	mux.HandleFunc("GET /bridge/health", func(w http.ResponseWriter, r *http.Request) {
		bridgeJSON(w, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /bridge/call", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string          `json:"method"`
			Args   json.RawMessage `json:"args"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r := core.JSONUnmarshal(body, &req); !r.OK {
			http.Error(w, r.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler.HandleBridgeCall(req.Method, req.Args)
		if err != nil {
			bridgeJSON(w, map[string]any{"error": err.Error()})
			return
		}
		bridgeJSON(w, map[string]any{"result": result})
	})

	// Listen on a random available port (localhost only)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, core.E("php.NewBridge", "listen", err)
	}

	bridge.port = listener.Addr().(*net.TCPAddr).Port
	bridge.server = &http.Server{Handler: mux}

	go func() {
		if err := bridge.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			core.Println(core.Sprintf("go-php: bridge error: %v", err))
		}
	}()

	core.Println(core.Sprintf("go-php: bridge listening on http://127.0.0.1:%d", bridge.port))
	return bridge, nil
}

// Port returns the port the bridge is listening on.
func (b *Bridge) Port() int {
	return b.port
}

// URL returns the full base URL of the bridge.
func (b *Bridge) URL() string {
	return core.Sprintf("http://127.0.0.1:%d", b.port)
}

// Shutdown gracefully stops the bridge server.
func (b *Bridge) Shutdown(ctx context.Context) error { // Result boundary
	return b.server.Shutdown(ctx)
}

func bridgeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	r := core.JSONMarshal(v)
	if !r.OK {
		core.Println(core.Sprintf("go-php: bridge encode error: %v", r.Error()))
		return
	}
	if _, err := w.Write(r.Value.([]byte)); err != nil {
		core.Println(core.Sprintf("go-php: bridge write error: %v", err))
	}
}
