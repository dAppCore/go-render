//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

// Package framework defines the HTTP-neutral boundary implemented by
// server-rendering framework adapters.
package framework

import (
	nethttp "net/http"

	core "dappco.re/go"
	"dappco.re/go/render/engine/ts"
)

// Loader creates resident CoreTS contexts for framework server bundles.
type Loader interface {
	Load(core.Context, string) (*ts.Context, error)
}

// Renderer translates an HTTP request through a server-rendering framework.
// A nil response means that the framework did not handle the request.
type Renderer interface {
	Render(core.Context, *nethttp.Request) (*Response, error)
	Close() error
}

// Response is a framework-rendered HTTP response.
type Response struct {
	Status int
	Header nethttp.Header
	Body   []byte
}
