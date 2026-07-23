// SPDX-Licence-Identifier: EUPL-1.2

package main

import (
	"net/http"

	core "dappco.re/go"
	webkit "dappco.re/go/render/display/webkit"
	"dappco.re/go/render/display/webkit/pkg/deno"
	"dappco.re/go/render/display/webkit/pkg/window"
)

func main() {
	if result := core.Stat("ui/dist/gui-ui/browser/index.html"); !result.OK {
		panic(core.E("cmd/webkit", "Angular production build is required: ui/dist/gui-ui/browser/index.html", result.Err()))
	}

	var manager *deno.Manager
	c := core.New(core.WithName("gui", func(c *core.Core) core.Result {
		manager = deno.New(deno.Options{Core: c})
		return webkit.NewService(webkit.GuiConfig{
			Mode: webkit.ModeSingleWindow,
			Name: "go-render-webkit",
			Assets: webkit.AssetOptions{
				Handler: http.FileServer(http.Dir("ui/dist/gui-ui/browser")),
			},
			Bindings: []webkit.Binding{webkit.Bind(manager)},
		})(c)
	}))

	service, ok := core.ServiceFor[*webkit.Service](c, "gui")
	if !ok || service == nil {
		panic(core.E("cmd/webkit", "gui service is unavailable", nil))
	}
	if manager == nil {
		panic(core.E("cmd/webkit", "deno manager is unavailable", nil))
	}

	ctx := core.Background()
	if result := service.OnStartup(ctx); !result.OK {
		panic(core.E("cmd/webkit", "start gui service", result.Err()))
	}
	app := service.App()
	if app == nil {
		panic(core.E("cmd/webkit", "Wails application is unavailable", nil))
	}
	if !webkit.OpenAdhocWindow(c, &window.Window{
		Name:   "main",
		Title:  "go-render-webkit",
		URL:    "/",
		Width:  1280,
		Height: 800,
	}) {
		panic(core.E("cmd/webkit", "open main window", nil))
	}

	if _, err := manager.Start(ctx); err != nil {
		panic(core.E("cmd/webkit", "start deno manager", err))
	}
	runResult := service.Run()
	_, stopErr := manager.Stop(ctx)
	if !runResult.OK {
		panic(core.E("cmd/webkit", "run gui service", runResult.Err()))
	}
	if stopErr != nil {
		panic(core.E("cmd/webkit", "stop deno manager", stopErr))
	}
}
