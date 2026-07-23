//go:build !js

// SPDX-License-Identifier: EUPL-1.2

// Service registration for the html package — exposes a default
// rendering Context as a Core service so consumers can wire HTML
// rendering through the same plumbing as every other core service.
//
//	c, _ := core.New(
//	    core.WithService(html.NewService(html.Options{
//	        Locale: "en",
//	    })),
//	)
//	svc := core.MustServiceFor[*html.Service](c, "html")
//	out := html.Render(myNode, svc.Context())
//
// Build-tagged !js so the WASM build (per RFC §7's 3.5 MB raw / 1 MB
// gzip budget) doesn't pull in dappco.re/go. The shared rendering
// surface (context.go, node.go, render.go, etc.) stays core-free; this
// file is the server-side service-registration layer only.

package html

import (
	core "dappco.re/go"
)

// Options configures the html service. Empty values fall back to
// package defaults (locale "en", no translator wired).
type Options struct {
	// Locale is the BCP 47 language tag for the default rendering
	// context. Empty → "en".
	Locale string
	// Translator is an optional translation provider. nil → text passes
	// through unchanged via Context.T fallback.
	Translator Translator
}

// Service is the registerable handle for the html package — embeds
// *core.ServiceRuntime[Options] for typed options access and holds a
// pre-built default Context ready for direct rendering.
//
// Usage example: `svc := core.MustServiceFor[*html.Service](c, "html"); out := html.Render(node, svc.Context())`
type Service struct {
	*core.ServiceRuntime[Options]
	ctx *Context
}

// NewService returns a factory that constructs a default rendering
// Context with the given options and produces a *Service ready for
// c.Service() registration. Use through core.WithService.
//
//	core.WithService(html.NewService(html.Options{Locale: "en"}))
func NewService(opts Options) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		locale := opts.Locale
		if locale == "" {
			locale = "en"
		}
		var ctx *Context
		if opts.Translator != nil {
			ctx = NewContextWithService(opts.Translator, locale)
		} else {
			ctx = NewContext(locale)
		}
		return core.Ok(&Service{
			ServiceRuntime: core.NewServiceRuntime(c, opts),
			ctx:            ctx,
		})
	}
}

// Register wires the html service into the Core with default Options —
// the imperative-style alternative to NewService for consumers that
// don't use the WithService factory pattern.
//
//	c := core.New()
//	if r := html.Register(c); !r.OK { return r }
func Register(c *core.Core) core.Result {
	return NewService(Options{})(c)
}

// Context returns the service's pre-built default rendering Context.
// Callers that need per-request state should construct their own
// Context via NewContext / NewContextWithService instead of mutating
// the shared one.
//
//	ctx := svc.Context()
//	out := html.Render(node, ctx)
func (s *Service) Context() *Context {
	if s == nil {
		return nil
	}
	return s.ctx
}
