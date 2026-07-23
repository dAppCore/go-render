//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package socket

import (
	"net"
	nethttp "net/http"

	core "dappco.re/go"
	httpdisplay "dappco.re/go/html/display/http"
	tsengine "dappco.re/go/html/engine/ts"
)

// Serve serves the HTTP SSR handler over listener until the listener closes.
// The listener may use TCP or a Unix-domain socket.
func Serve(engine tsengine.Renderer, listener net.Listener, opts ...httpdisplay.Option) error {
	if engine == nil {
		return core.E("socket.Serve", "render engine is nil", nil)
	}
	if listener == nil {
		return core.E("socket.Serve", "listener is nil", nil)
	}

	err := nethttp.Serve(listener, httpdisplay.Handler(engine, "", opts...))
	if err == nil || core.Is(err, nethttp.ErrServerClosed) || core.Is(err, net.ErrClosed) {
		return nil
	}
	return core.E("socket.Serve", "serve render listener", err)
}
