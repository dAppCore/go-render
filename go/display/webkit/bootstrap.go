// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/browser"
	"dappco.re/go/render/display/webkit/pkg/chat"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/container"
	"dappco.re/go/render/display/webkit/pkg/contextmenu"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/display"
	"dappco.re/go/render/display/webkit/pkg/dock"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/events"
	"dappco.re/go/render/display/webkit/pkg/keybinding"
	"dappco.re/go/render/display/webkit/pkg/lifecycle"
	"dappco.re/go/render/display/webkit/pkg/menu"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"dappco.re/go/render/display/webkit/pkg/p2p"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/systray"
	"dappco.re/go/render/display/webkit/pkg/webview"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Bootstrap returns the full set of [core.CoreOption] a desktop consumer
// needs to wire the GUI service stack in one append, instead of hand-wiring
// every sub-service.
//
// What's included:
//   - "gui"          — the top-level gui shell service
//   - "window"       — window manager + layout/state persistence (NewLayoutManager
//     uses DIR_CONFIG / Core/layouts.json by default)
//   - "display"      — screens, work areas, dialogs, clipboard, system tray
//   - "webview"      — JS evaluation, console capture, DOM queries (CDP-backed)
//   - "menu"         — native application menu
//   - "systray"      — system tray icon + menu
//   - "browser"      — open external URLs/files in the OS default app
//     (forge/mantis/docs links; native target=_blank replacement)
//   - "notification" — native OS notifications (macOS Notification Center,
//     Windows Toast, Linux D-Bus). macOS requires bundle+sign to fire.
//   - "lifecycle"    — app lifecycle events (started, will-terminate,
//     did-become-active, did-resign-active, opened-with-file). Subscribers
//     receive the corresponding c.Action dispatches.
//   - "dialog"       — native file open/save/directory pickers + info/
//     warning/error/question message dialogs. File pickers block until
//     the user resolves the dialog; message dialogs adapt Wails's async
//     button-callback model to a synchronous "which label was clicked"
//     return.
//   - "contextmenu"  — register named native context menus. Frontend
//     opts an element in via the CSS custom property
//     --custom-contextmenu: <name>; (with optional
//     --custom-contextmenu-data for per-element payload). Item clicks
//     dispatch ActionItemClicked on the consumer's core.
//   - "keybinding"   — register global accelerators at runtime
//     (Cmd+S, Ctrl+P, F1, etc.). Shares the same Wails key-binding map
//     as application.Options.KeyBindings — runtime registration wins
//     over boot-time when accelerators collide.
//   - "dock"         — macOS dock + Windows taskbar icon visibility +
//     badge label. Progress bar / bounce are accepted but no-op until
//     Wails exposes them upstream.
//   - "environment"  — OS / arch / debug / platform info, dark-mode
//     query + ThemeChanged subscription, accent colour, OpenFileManager,
//     focus-follows-mouse (Linux).
//   - "screen"       — multi-monitor info: GetAll / GetPrimary /
//     GetCurrent (containing-window fallback to primary).
//   - "clipboard"    — system clipboard text read/write. Wails alpha.83
//     does not expose an image clipboard API; image ops fall through
//     to the gui Service's "platform unavailable" branch.
//   - "events"       — Wails custom event bus: Emit / On / Off /
//     OnMultiple / Reset. CustomEvent translates between gui and Wails
//     types so consumers never see Wails directly.
//   - "chat"         — LLM chat surface. Defaults: APIURL
//     "http://localhost:8090", store at $DIR_HOME/.core/gui/chat.db.
//     Auto-builds its own gui_mcp subsystem as ToolExecutor if no
//     consumer-supplied one. Override defaults via chat.Register(opts...).
//   - "container"    — TIM (Trusted In-Memory) container manager.
//     Wired with empty TIMOptions; container Image/Name configured at
//     runtime via the core config service. Operations on an unconfigured
//     manager surface clear "no image" errors.
//   - "p2p"          — peer-to-peer routing service with TCP driver.
//     Wired with empty Options (no listen address, no peers); add
//     bootstrap peers at runtime via the p2p config service. Subscribe/
//     Publish are no-ops until ListenAddr is set.
//
// The wails [*application.App] is the only boundary the consumer touches —
// after that, everything runs through the canonical Core IPC pattern
// (c.Action/c.QUERY) so consumers don't need direct wails imports.
//
//	app := application.New(opts) // consumer creates the app
//	coreOpts := []core.CoreOption{ /* your services */ }
//	coreOpts = append(coreOpts, webkit.Bootstrap(app)...)
//	c, _ := core.New(coreOpts...)
//
// Packages NOT wired into Bootstrap (deliberately):
//   - mcp — already wired by the IDE consumer as "gui_mcp"; that
//     subsystem is what surfaces every other service's MCP tools and
//     must be owned by the consumer (it attaches to the IDE's MCP
//     bridge service, not to the gui shell).
//   - marketplace — pure library; the marketplace_list / fetch /
//     verify / install MCP tools delegate to display.marketplace.*
//     actions registered by the display service above.
//   - deno — sidecar Manager (deno.New(Options) *Manager), not a Core
//     service. Used by feature code that wants to spawn a deno
//     subprocess (typed extensions, plugin runtime). Wire per-use, not
//     globally.
//   - preload — JS-injection assets + TrustedOriginPolicy helpers used
//     by the window service. Not a registerable service.
func Bootstrap(app *application.App) []core.CoreOption {
	return BootstrapWithConfig(app, BootstrapConfig{})
}

// BootstrapConfig lets the consumer override the chat / container / p2p
// service options without rewriting the rest of Bootstrap. Pass blank
// values to keep package defaults; only set the fields you actually
// want to override (the IDE's settings panel writes here).
type BootstrapConfig struct {
	Chat      ChatConfig
	Container ContainerConfig
	P2P       P2PConfig
}

// ChatConfig — overrides for the gui chat service.
//
// APIURL blank → http://localhost:8090. StorePath blank → package default
// $DIR_HOME/.core/gui/chat.db. ToolExecutor stays auto-built (gui_mcp
// subsystem) when nil; that's what the IDE wants 99% of the time.
type ChatConfig struct {
	APIURL    string
	StorePath string
}

// ContainerConfig — overrides for the gui container service (TIM).
//
// All-blank means the manager starts in "no container configured" mode;
// /dev/tim surfaces this as empty state. The Exec field is wired by
// the container service itself when nil — it falls back to the core
// process service.
type ContainerConfig struct {
	Image   string
	Name    string
	DataDir string
	Command []string
}

// P2PConfig — overrides for the gui p2p service (TCP driver).
//
// All-blank means no listener bound, no peers — Subscribe/Publish are
// no-ops and /dev/p2p surfaces "no listener" empty state. NodeID is
// auto-assigned by the TCP driver when blank.
type P2PConfig struct {
	ListenAddr string
	PeerAddrs  []string
	NodeID     string
}

// BootstrapWithConfig is Bootstrap with consumer-supplied overrides for
// chat / container / p2p. Use this when the consumer (e.g. core/ide)
// wants to drive those services from its own config layer rather than
// take package defaults.
//
//	cfg := webkit.BootstrapConfig{
//	    Chat:      webkit.ChatConfig{APIURL: "http://localhost:11434"},
//	    Container: webkit.ContainerConfig{Image: "alpine", Name: "tim-1"},
//	    P2P:       webkit.P2PConfig{ListenAddr: "127.0.0.1:9100"},
//	}
//	coreOpts = append(coreOpts, webkit.BootstrapWithConfig(app, cfg)...)
func BootstrapWithConfig(app *application.App, cfg BootstrapConfig) []core.CoreOption {
	if app == nil {
		return nil
	}
	return []core.CoreOption{
		core.WithName("gui", NewService(GuiConfig{})),
		core.WithService(window.Register(window.NewWailsPlatform(app))),
		core.WithService(display.Register(app)),
		core.WithService(webview.Register()),
		core.WithService(menu.Register(menu.NewWailsPlatform(app))),
		core.WithService(systray.Register(systray.NewWailsPlatform(app))),
		core.WithService(browser.Register(browser.NewWailsPlatform(app))),
		core.WithService(notification.Register(notification.NewWailsPlatform(app))),
		core.WithService(lifecycle.Register(lifecycle.NewWailsPlatform(app))),
		core.WithService(dialog.Register(dialog.NewWailsPlatform(app))),
		core.WithService(contextmenu.Register(contextmenu.NewWailsPlatform(app))),
		core.WithService(keybinding.Register(keybinding.NewWailsPlatform(app))),
		core.WithService(dock.Register(dock.NewWailsPlatform(app))),
		core.WithService(environment.Register(environment.NewWailsPlatform(app))),
		core.WithService(screen.Register(screen.NewWailsPlatform(app))),
		core.WithService(clipboard.Register(clipboard.NewWailsPlatform(app))),
		core.WithService(events.Register(events.NewWailsPlatform(app))),
		core.WithService(chat.Register(applyChatConfig(cfg.Chat)...)),
		core.WithService(func(c *core.Core) core.Result {
			return core.Result{Value: container.NewService(c, container.TIMOptions{
				Image:   cfg.Container.Image,
				Name:    cfg.Container.Name,
				DataDir: cfg.Container.DataDir,
				Command: cfg.Container.Command,
			}), OK: true}
		}),
		core.WithService(func(c *core.Core) core.Result {
			return core.Result{Value: p2p.NewService(c, p2p.Options{
				ListenAddr: cfg.P2P.ListenAddr,
				PeerAddrs:  cfg.P2P.PeerAddrs,
				NodeID:     cfg.P2P.NodeID,
			}), OK: true}
		}),
	}
}

// applyChatConfig translates the consumer's ChatConfig into the chat
// package's optional-fn slice. Only emits an opt fn for non-blank
// fields so package defaults stay in effect for anything the consumer
// didn't override.
func applyChatConfig(cfg ChatConfig) []func(*chat.Options) {
	var fns []func(*chat.Options)
	if cfg.APIURL != "" {
		apiURL := cfg.APIURL
		fns = append(fns, func(o *chat.Options) { o.APIURL = apiURL })
	}
	if cfg.StorePath != "" {
		storePath := cfg.StorePath
		fns = append(fns, func(o *chat.Options) { o.StorePath = storePath })
	}
	return fns
}
