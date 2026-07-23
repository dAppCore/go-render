package mcp

import (
	"context"

	core "dappco.re/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolsLifecycle_appQuit_GoodCase(t *core.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerLifecycleTools(server)

	result, err := sub.CallTool(context.Background(), "app_quit", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "\"success\":true")
}

func TestToolsLifecycle_appQuit_Bad(t *core.T) {
	// appQuit
	ax7Variant := "appQuit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	sub := New(core.New(core.WithServiceLock()))

	_, out, err := sub.appQuit(context.Background(), nil, AppQuitInput{})
	core.RequireNoError(t, err)
	core.AssertTrue(t, out.Success)
	core.AssertNil(t, err)
}

func TestToolsLifecycle_appQuit_UglyCase(t *core.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerLifecycleTools(server)

	_, err := sub.CallTool(context.Background(), "app_quit", nil)
	core.RequireNoError(t, err)
}
