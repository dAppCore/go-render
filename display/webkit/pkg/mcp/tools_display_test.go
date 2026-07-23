package mcp

import (
	"context"

	core "dappco.re/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newDisplayToolTestSubsystem(t *core.T, handler func(core.Options) core.Result) *Subsystem {
	t.Helper()
	c := core.New(core.WithServiceLock())
	if handler != nil {
		c.Action("display.resolve_scheme", func(_ context.Context, opts core.Options) core.Result {
			return handler(opts)
		})
	}
	return New(c)
}

func TestToolsDisplay_schemeResolve_GoodCase(t *core.T) {
	sub := newDisplayToolTestSubsystem(t, func(opts core.Options) core.Result {
		return core.Result{
			Value: map[string]any{
				"url":          opts.String("url"),
				"route":        "store",
				"content_type": "text/html",
				"body":         "<html>core://store</html>",
			},
			OK: true,
		}
	})
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerDisplayTools(server)

	result, err := sub.CallTool(context.Background(), "scheme_resolve", map[string]any{"url": "core://store?q=alpha"})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "core://store?q=alpha")
	core.AssertContains(t, result, "\"route\":\"store\"")
	core.AssertContains(t, result, "\"content_type\":\"text/html\"")
}

func TestToolsDisplay_schemeResolve_Bad(t *core.T) {
	// schemeResolve
	ax7Variant := "schemeResolve:bad"
	core.AssertContains(t, ax7Variant, "bad")
	sub := newDisplayToolTestSubsystem(t, func(core.Options) core.Result {
		return core.Result{Value: core.NewError("display offline"), OK: false}
	})

	_, _, err := sub.schemeResolve(context.Background(), nil, SchemeResolveInput{URL: "core://store"})
	core.AssertError(t, err)
	core.AssertEqual(t, "display offline", err.Error())
}

func TestToolsDisplay_schemeResolve_Ugly(t *core.T) {
	// schemeResolve
	ax7Variant := "schemeResolve:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	sub := newDisplayToolTestSubsystem(t, func(core.Options) core.Result {
		return core.Result{Value: map[string]any{
			"route":        "store",
			"content_type": "text/html",
			"body":         "<html>fallback</html>",
		}, OK: true}
	})

	_, out, err := sub.schemeResolve(context.Background(), nil, SchemeResolveInput{URL: "core://store?q=beta"})
	core.RequireNoError(t, err)
	core.AssertEqual(t, "core://store?q=beta", out.URL)
	core.AssertEqual(t, "store", out.Route)
}
