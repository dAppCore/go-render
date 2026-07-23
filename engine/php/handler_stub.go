//go:build !frankenphp

package php

import (
	"net/http"

	core "dappco.re/go"
)

// Handler implements http.Handler when the embedded FrankenPHP runtime is not built.
type Handler struct {
	docRoot     string
	laravelRoot string
}

// HandlerConfig configures the FrankenPHP handler.
type HandlerConfig struct {
	NumThreads int
	NumWorkers int
	PHPIni     map[string]string
}

// NewHandler returns a handler placeholder unless built with -tags frankenphp.
func NewHandler(laravelRoot string, cfg HandlerConfig) (*Handler, func(), error) { // Result boundary
	if cfg.NumThreads == 0 {
		cfg.NumThreads = 4
	}
	if cfg.NumWorkers == 0 {
		cfg.NumWorkers = 2
	}

	handler := &Handler{
		docRoot:     core.PathJoin(laravelRoot, "public"),
		laravelRoot: laravelRoot,
	}
	cleanup := func() {
		// No resources are allocated when embedded FrankenPHP is not built.
	}
	return handler, cleanup, core.Errorf("embedded FrankenPHP support is not built; rebuild with -tags frankenphp")
}

// LaravelRoot returns the path to the extracted Laravel application.
func (h *Handler) LaravelRoot() string {
	return h.laravelRoot
}

// DocRoot returns the path to the document root (public/).
func (h *Handler) DocRoot() string {
	return h.docRoot
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "embedded FrankenPHP support is not built", http.StatusNotImplemented)
}
