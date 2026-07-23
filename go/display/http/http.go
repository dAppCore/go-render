//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package httpdisplay

import (
	nethttp "net/http"

	core "dappco.re/go"
	tsengine "dappco.re/go/render/engine/ts"
)

// Option configures Handler.
type Option func(*handlerOptions)

type handlerOptions struct {
	entry string
}

// WithEntry selects the TypeScript or JavaScript SSR entry module.
func WithEntry(entry string) Option {
	return func(options *handlerOptions) {
		options.entry = entry
	}
}

// Handler returns an HTTP handler that serves files found under assetsDir and
// server-renders every other request through engine.
//
// Static files retain their standard content types and caching semantics.
// Rendered responses receive text/html; charset=utf-8.
func Handler(engine tsengine.Renderer, assetsDir string, opts ...Option) nethttp.Handler {
	options := handlerOptions{}
	for _, option := range opts {
		if option != nil {
			option(&options)
		}
	}

	var assets nethttp.Handler
	if assetsDir != "" {
		assets = nethttp.FileServer(nethttp.Dir(assetsDir))
	}

	return nethttp.HandlerFunc(func(writer nethttp.ResponseWriter, request *nethttp.Request) {
		if assets != nil && isAssetRequest(assetsDir, request) {
			assets.ServeHTTP(writer, request)
			return
		}
		if engine == nil {
			nethttp.Error(writer, "render engine unavailable", nethttp.StatusInternalServerError)
			return
		}
		if core.Trim(options.entry) == "" {
			nethttp.Error(writer, "render entry unavailable", nethttp.StatusInternalServerError)
			return
		}

		input := map[string]any{
			"method":  request.Method,
			"url":     request.URL.String(),
			"path":    request.URL.Path,
			"query":   request.URL.Query(),
			"headers": request.Header.Clone(),
			"host":    request.Host,
		}
		output, err := engine.Render(request.Context(), options.entry, input)
		if err != nil {
			nethttp.Error(writer, "render failed", nethttp.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.WriteHeader(nethttp.StatusOK)
		_, _ = writer.Write(output)
	})
}

func isAssetRequest(assetsDir string, request *nethttp.Request) bool {
	if request.Method != nethttp.MethodGet && request.Method != nethttp.MethodHead {
		return false
	}

	file, err := nethttp.Dir(assetsDir).Open(request.URL.Path)
	if err != nil {
		return false
	}
	defer file.Close()

	info, err := file.Stat()
	return err == nil && !info.IsDir()
}
