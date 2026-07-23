// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	"net/http"
	"strings"
)

// WailsHTTPMiddleware returns a MiddlewareFunc that splits HTTP traffic
// between the host app's engine (Gin / chi / stdlib mux / anything else)
// and the Wails runtime. Requests under /wails/* go to the runtime's
// next handler; everything else is delegated to engine.
//
// One carve-out: /wails/custom.js. Wails fetches this URL at boot to
// allow user-supplied JS overrides; apps that ship none would otherwise
// see a 404 spammed into the console on every reload. The middleware
// intercepts that exact path and returns an empty 200 with the right
// content-type, so the runtime continues with no overrides applied.
//
// Wire it via GuiConfig.Bindings + AssetOptions:
//
//	webkit.GuiConfig{
//	    Assets: webkit.AssetOptions{
//	        Handler:    engine,                              // host HTTP handler
//	        Middleware: webkit.WailsHTTPMiddleware(engine),     // wails carve-out
//	    },
//	}
//
// Replaces the per-app `ginMiddleware` helper every Wails+HTTP consumer
// would otherwise hand-roll (see wails-v3 examples/gin-routing).
func WailsHTTPMiddleware(engine http.Handler) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/wails/custom.js" {
				w.Header().Set("Content-Type", "application/javascript")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("/* no user overrides */\n"))
				return
			}
			if strings.HasPrefix(r.URL.Path, "/wails") {
				next.ServeHTTP(w, r)
				return
			}
			engine.ServeHTTP(w, r)
		})
	}
}
