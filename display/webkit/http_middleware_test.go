// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWailsHTTPMiddleware_Good covers the three routing arms:
// /wails/custom.js → empty 200, /wails/* → next, anything else → engine.
func TestWailsHTTPMiddleware_Good(t *testing.T) {
	engine := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("engine"))
	})
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("next"))
	})

	mw := WailsHTTPMiddleware(engine)(next)

	cases := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
		wantCT     string
	}{
		{"custom_js", "/wails/custom.js", http.StatusOK, "/* no user overrides */\n", "application/javascript"},
		{"wails_runtime", "/wails/runtime.js", http.StatusAccepted, "next", ""},
		{"engine_route", "/api/v1/foo", http.StatusTeapot, "engine", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()
			mw.ServeHTTP(w, r)
			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if got := w.Body.String(); got != tc.wantBody {
				t.Fatalf("body = %q, want %q", got, tc.wantBody)
			}
			if tc.wantCT != "" && w.Header().Get("Content-Type") != tc.wantCT {
				t.Fatalf("content-type = %q, want %q", w.Header().Get("Content-Type"), tc.wantCT)
			}
		})
	}
}

// TestWailsHTTPMiddleware_Bad covers the nil engine input — a
// well-known caller mistake. Panics from inside HandlerFunc are the
// stdlib idiom; the test asserts the panic so a future refactor can't
// silently swallow misuse.
func TestWailsHTTPMiddleware_Bad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on nil engine, got none")
		}
	}()
	mw := WailsHTTPMiddleware(nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	mw.ServeHTTP(httptest.NewRecorder(), r)
}
