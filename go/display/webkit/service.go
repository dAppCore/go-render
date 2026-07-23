// SPDX-License-Identifier: EUPL-1.2

// Package webkit is the canonical entry-point for the go-render display/webkit module.
// One service per repo — `c.Service("gui")` returns this.
//
// The Service owns the Wails application lifecycle for consumers that
// want a desktop runtime. Consumers register sibling domain services
// on Core first, then register webkit.Service last with options. The
// Service's OnStartup constructs the wails App from those options and
// registers every webkit sub-service (window, systray, lifecycle,
// menu, dialog, contextmenu, keybinding, dock, environment, screen,
// clipboard, events, browser, notification, webview, display, chat,
// container, p2p). Run() blocks on the wails event loop until the user
// or OS quits.
//
//	cfg := webkit.GuiConfig{
//	    Name:        "lthn",
//	    Description: "Lethean Desktop",
//	    Icon:        appIconBytes,
//	    Assets:      webkit.AssetOptions{Handler: ginEngine, Middleware: ginMW},
//	    Mac:         webkit.MacOptions{ActivationPolicy: webkit.ActivationPolicyAccessory},
//	    Bindings:    []any{runnerSvc, serverSvc, /* …domain services… */},
//	}
//	c, _ := core.New( /* domain services */, core.WithName("gui", webkit.NewService(cfg)))
//	if r := c.Service("gui").Value.(*webkit.Service).Run(); !r.OK { … }

package webkit

import (
	"context"
	"net/http"

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

// GuiConfig configures the gui service. All fields are optional except
// when a consumer needs the corresponding wails surface — leaving e.g.
// Mac zeroed yields the wails defaults.
type GuiConfig struct {
	// Mode is the meta-decision driving platform-policy defaults
	// (Mac.ActivationPolicy, terminate-after-last-window behaviour,
	// Windows.DisableQuitOnLastWindowClosed). One of ModeTray /
	// ModeSingleWindow / ModeMultiWindow. ModeDefault (zero value)
	// leaves all fields at their package zeros. Explicit fields below
	// always win — Mode only fills gaps.
	Mode Mode

	// Name is the application name shown in the default about box and
	// platform UI. Default: "core-gui".
	Name string

	// Description is the application description shown by the OS.
	Description string

	// Icon is the default application icon (PNG bytes). Empty leaves
	// the wails default in place.
	Icon []byte

	// Assets binds the WebView's HTTP handler + middleware chain.
	// The handler is the consumer's HTTP server (Gin, std mux, etc.);
	// the middleware lets the consumer carve out `/wails/*` and other
	// platform-internal routes before the user handler sees them.
	Assets AssetOptions

	// Mac carries macOS-specific application options (activation
	// policy, terminate-after-last-window behaviour).
	Mac MacOptions

	// Windows carries Windows-specific application options
	// (disable-quit-on-last-window, WebView2 feature flags).
	Windows WindowsOptions

	// SingleInstance enables OS-level single-instance enforcement.
	// Nil disables single-instance — every launch creates a new
	// process.
	SingleInstance *SingleInstanceOptions

	// Bindings are domain services that the consumer wants exposed as
	// Wails IPC bindings. Each entry is constructed via webkit.Bind[T];
	// bindings land at the service's package path in the frontend.
	Bindings []Binding

	// WindowRegistry is the declarative window list. Each entry is
	// registered + pre-created (hidden) during OnStartup, with
	// HideOnClose / ContentProtection auto-applied post-create. The
	// caller can open any registered window later via the
	// `window.open` / `window.set_visibility` actions. Empty slice =
	// no pre-create; consumers can still issue ad-hoc TaskOpenWindow.
	//
	// Named WindowRegistry (not Windows) to avoid clashing with the
	// Windows field above which carries the Microsoft Windows OS
	// options.
	WindowRegistry []*window.Window

	// WindowStatePath is the on-disk path the window service uses to
	// persist per-window position/size/maximised state. Empty leaves
	// the window service defaults in place (DIR_CONFIG/Core/
	// window_state.json). Apps that own their own conf directory
	// (lthn → ~/Lethean/conf/window_state.json) override here.
	WindowStatePath string

	// WindowLayoutPath is the on-disk path for the named-layout
	// store. Empty = window service default. Layouts (multi-window
	// snapshots) persist here.
	WindowLayoutPath string

	// Tray declares the system tray surface (icon, tooltip, label,
	// menu items + popover window attachment). Nil = no tray
	// configured; the systray sub-service is still registered but
	// the consumer drives its config imperatively. Click routing
	// stays in the consumer: when a tray menu item is clicked, an
	// ActionTrayMenuItemClicked message lands on the action bus
	// carrying the item's ActionID — register a handler with
	// core.RegisterAction to dispatch.
	Tray *TrayConfig

	// Keybindings declares global hotkeys. webkit.Service registers each
	// accelerator + installs a shared trigger router that emits the
	// configured EventName when the accelerator fires.
	Keybindings []Keybinding

	// ContextMenus declares named right-click surfaces. webkit.Service
	// registers each menu + installs a relay that emits the
	// configured EventTemplate on item click.
	ContextMenus []ContextMenu

	// AppMenu declares the macOS application menu bar. Auto-gated to
	// darwin — other platforms either lack a global app menu (Linux
	// varies) or surface menus per-window (Windows). Role-based items
	// (e.g. &webkit.RoleAppMenu) get platform-managed entries (Quit / Cut
	// / Window list) for free.
	AppMenu []MenuItem

	// ShouldQuit returns false to veto an OS/user quit request. Nil
	// means "always allow quit" (wails default).
	ShouldQuit func() bool

	// OnShutdown fires after the user accepts quit, before the wails
	// event loop tears down. Drain in-flight work here.
	OnShutdown func()

	// PostShutdown fires after the wails event loop has fully stopped.
	// Last chance to close anything that held a ref into the loop.
	PostShutdown func()

	// OnPanic captures uncaught panics from Go-side service methods
	// (binding adapters etc.). Renamed from wails's PanicHandler so
	// every Options callback uses the On* prefix uniformly.
	OnPanic func(PanicDetails)

	// Chat / Container / P2P — sub-service config overrides forwarded
	// to chat.Register / container.NewService / p2p.NewService. Blank
	// fields keep package defaults. Mirrors BootstrapConfig.
	Chat      ChatConfig
	Container ContainerConfig
	P2P       P2PConfig
}

// AssetOptions binds the WebView's HTTP handler + middleware chain.
// The handler serves the consumer's frontend (embedded vite dist,
// reverse-proxied dev server, etc.); the middleware lets the consumer
// carve out platform-internal routes (e.g. /wails/*) before the user
// handler sees them.
type AssetOptions struct {
	Handler        http.Handler
	Middleware     MiddlewareFunc
	DisableLogging bool
}

// MiddlewareFunc wraps a downstream HTTP handler. Mirrors the stdlib
// idiom — implementations call next.ServeHTTP for requests they don't
// intercept.
type MiddlewareFunc func(next http.Handler) http.Handler

// MacOptions carries macOS-specific application behaviour.
type MacOptions struct {
	// ApplicationShouldTerminateAfterLastWindowClosed: when false, the
	// process survives the last window closing. Required for menu-bar
	// / accessory apps where the tray IS the lifetime anchor.
	ApplicationShouldTerminateAfterLastWindowClosed bool
	// ActivationPolicy controls the macOS app activation kind
	// (Regular, Accessory, Prohibited). See the ActivationPolicy
	// constants.
	ActivationPolicy ActivationPolicy
}

// WindowsOptions carries Windows-specific application behaviour.
type WindowsOptions struct {
	// DisableQuitOnLastWindowClosed: when true, closing the last
	// window keeps the process alive. Required for tray-anchored
	// lifetime, mirroring Mac.ApplicationShouldTerminateAfterLastWindowClosed.
	DisableQuitOnLastWindowClosed bool
	// EnabledFeatures opts into named WebView2 features
	// (e.g. "msWebView2EnableDraggableRegions" for --wails-draggable
	// CSS support).
	EnabledFeatures []string
	// DisabledFeatures explicitly disables named WebView2 feature flags.
	DisabledFeatures []string
	// AdditionalBrowserArgs passes process-wide switches to WebView2.
	// Entries include their "--" prefix.
	AdditionalBrowserArgs []string
	// UseVisualHosting selects WebView2's DirectComposition visual
	// hosting. It avoids the monitor-scale re-marshalling path that can
	// stall remote-desktop sessions and is configured before app.Run.
	UseVisualHosting bool
	// WndClass overrides the native Windows window class.
	WndClass string
	// WebviewUserDataPath overrides WebView2's user-data directory.
	WebviewUserDataPath string
	// WebviewBrowserPath selects a fixed WebView2 runtime directory.
	WebviewBrowserPath string
}

// ActivationPolicy is the macOS activation kind. Mirrors the wails
// enum but stays string-typed for read-as-docs call sites.
type ActivationPolicy int

const (
	// ActivationPolicyRegular: standard app with Dock icon + Cmd-Tab
	// presence. wails default.
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory: menu-bar / tray-only, no Dock icon,
	// no Cmd-Tab entry.
	ActivationPolicyAccessory
	// ActivationPolicyProhibited: hidden from Dock + Cmd-Tab + Mission
	// Control. Background-only.
	ActivationPolicyProhibited
)

// SingleInstanceOptions enables OS-level single-instance enforcement.
// A second launch hands off URL/file/args to the first instance via
// OnSecondInstanceLaunch and then exits.
type SingleInstanceOptions struct {
	// UniqueID is the OS-level instance key. macOS uses the Bundle
	// Identifier; Windows uses a mutex name. Pick a stable value per
	// app/build channel so OS-level flock behaviour is correct.
	UniqueID string
	// EncryptionKey enables AES-256-GCM on the inter-instance channel.
	// Without it second-instance payloads MUST be treated as
	// untrusted. Per-install random key recommended.
	EncryptionKey [32]byte
	// AdditionalData is forwarded to OnSecondInstanceLaunch alongside
	// the launch args. Useful for app/version stamping so the receiver
	// can audit which build sent the payload.
	AdditionalData map[string]string
	// OnSecondInstanceLaunch fires in the FIRST instance when a second
	// launch attempt occurs. The data carries the second instance's
	// args/workingDir/additionalData; the first instance typically
	// re-broadcasts onto its own bus + brings a window to the front.
	OnSecondInstanceLaunch func(SecondInstanceData)
}

// SecondInstanceData carries the second-launch context delivered to
// OnSecondInstanceLaunch on the first running instance.
type SecondInstanceData struct {
	Args           []string
	WorkingDir     string
	AdditionalData map[string]string
}

// PanicDetails describes an uncaught panic from a Go-side binding or
// lifecycle hook. Mirrors wails's shape minus the time field (the
// consumer's logger stamps its own time).
type PanicDetails struct {
	Error          error
	StackTrace     string
	FullStackTrace string
}

// Service is the registerable handle for the gui repo. It owns the
// wails *application.App, the registered sub-service set (window /
// systray / menu / ...), and the Run loop.
//
// Usage example: `svc := core.MustServiceFor[*webkit.Service](c, "gui")`
type Service struct {
	*core.ServiceRuntime[GuiConfig]
	registrations core.Once
	app           *application.App
}

// NewService returns the factory for c.Service() registration.
//
// Usage example: `c, _ := core.New(core.WithName("gui", webkit.NewService(webkit.GuiConfig{})))`
func NewService(config GuiConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Ok(&Service{
			ServiceRuntime: core.NewServiceRuntime(c, config),
		})
	}
}

// Register builds the gui service with default GuiConfig — imperative
// alternative to NewService.
//
// Usage example: `r := webkit.Register(c)`
func Register(c *core.Core) core.Result {
	return NewService(GuiConfig{})(c)
}

// App returns the constructed wails *application.App for callers that
// need direct access (legacy bridges, native integrations not yet
// surfaced via display/webkit sub-services). Returns nil before OnStartup
// runs.
//
// Usage example: `app := svc.App()`
func (s *Service) App() *application.App {
	if s == nil {
		return nil
	}
	return s.app
}

// OnStartup implements core.Startable. Builds the wails application
// from GuiConfig, registers every display/webkit sub-service on the Core,
// and starts each sub-service's OnStartup. Idempotent via core.Once.
//
// OnStartup does NOT block on the event loop — Run() is the blocking
// call. This split lets consumers register additional post-app actions
// (tray menu config, key bindings, window pre-create) between
// OnStartup and Run.
//
// Usage example: `r := svc.OnStartup(ctx)`
func (s *Service) OnStartup(ctx context.Context) core.Result {
	if s == nil {
		return core.Ok(nil)
	}
	var startErr core.Result
	s.registrations.Do(func() {
		startErr = s.start(ctx)
	})
	return startErr
}

// Run blocks on the wails event loop until the user or OS quits. Must
// be called after OnStartup. Returns the wails app.Run() error wrapped
// as core.Result.
//
// Usage example: `r := svc.Run()`
func (s *Service) Run() core.Result {
	if s == nil || s.app == nil {
		return core.Fail(core.E("gui.Service.Run", "OnStartup has not been called", nil))
	}
	if err := s.app.Run(); err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}

// OnShutdown implements core.Stoppable. The Service does not actively
// stop the wails event loop — Wails' own quit flow handles that. This
// hook is here for symmetry + future expansion.
//
// Usage example: `r := svc.OnShutdown(ctx)`
func (s *Service) OnShutdown(context.Context) core.Result {
	return core.Ok(nil)
}

// start builds the wails App + registers sub-services. Called once
// inside OnStartup's core.Once.Do.
func (s *Service) start(ctx context.Context) core.Result {
	cfg := s.Options()
	applyModeDefaults(&cfg)
	s.app = application.New(buildWailsOptions(cfg))

	// Register sub-services on the Core, mirroring BootstrapWithConfig
	// but applied post-hoc to the existing Core (since the App was
	// just built, we need to attach factories that take it as a ref).
	registrations := []struct {
		name    string
		factory func(*core.Core) core.Result
	}{
		{"display", display.Register(s.app)},
		{"window", window.Register(window.NewWailsPlatform(s.app))},
		{"webview", webview.Register()},
		{"menu", menu.Register(menu.NewWailsPlatform(s.app))},
		{"systray", systray.Register(systray.NewWailsPlatform(s.app))},
		{"browser", browser.Register(browser.NewWailsPlatform(s.app))},
		{"notification", notification.Register(notification.NewWailsPlatform(s.app))},
		{"lifecycle", lifecycle.Register(lifecycle.NewWailsPlatform(s.app))},
		{"dialog", dialog.Register(dialog.NewWailsPlatform(s.app))},
		{"contextmenu", contextmenu.Register(contextmenu.NewWailsPlatform(s.app))},
		{"keybinding", keybinding.Register(keybinding.NewWailsPlatform(s.app))},
		{"dock", dock.Register(dock.NewWailsPlatform(s.app))},
		{"environment", environment.Register(environment.NewWailsPlatform(s.app))},
		{"screen", screen.Register(screen.NewWailsPlatform(s.app))},
		{"clipboard", clipboard.Register(clipboard.NewWailsPlatform(s.app))},
		{"events", events.Register(events.NewWailsPlatform(s.app))},
		{"chat", chat.Register(applyChatConfig(cfg.Chat)...)},
		{"container", containerFactory(cfg.Container)},
		{"p2p", p2pFactory(cfg.P2P)},
	}
	c := s.Core()
	for _, r := range registrations {
		if rr := registerSubservice(ctx, c, r.name, r.factory); !rr.OK {
			return rr
		}
	}

	// Apply WindowStatePath / WindowLayoutPath if the caller supplied
	// overrides. The window service is now registered + started, so
	// we can point its persistence layer at the consumer's chosen
	// directory before any registry pre-create runs.
	if cfg.WindowStatePath != "" || cfg.WindowLayoutPath != "" {
		if winSvc, ok := core.ServiceFor[*window.Service](c, "window"); ok {
			if cfg.WindowStatePath != "" {
				winSvc.Manager().State().SetPath(cfg.WindowStatePath)
			}
			if cfg.WindowLayoutPath != "" {
				winSvc.Manager().Layout().SetPath(cfg.WindowLayoutPath)
			}
		}
	}

	// Window registry — register each declared window with the window
	// service + pre-create as hidden so the first show is instant.
	// HideOnClose + ContentProtection auto-apply post-create via the
	// existing TaskSetCloseBehavior / TaskSetContentProtection actions.
	for _, w := range cfg.WindowRegistry {
		if w == nil || w.Name == "" {
			continue
		}
		if rr := registerAndPreCreateWindow(c, w); !rr.OK {
			return rr
		}
	}

	// Tray config — apply declared icon / tooltip / label / menu /
	// popover attachment. Click routing can be declared via
	// cfg.Tray.Routes; entries not in the table fall through to the
	// consumer's own ActionTrayMenuItemClicked handler.
	applyTrayConfig(c, cfg.Tray)
	if cfg.Tray != nil {
		applyTrayRoutes(c, cfg.Tray.Routes)
	}

	// Keybindings — register each accelerator + install the shared
	// trigger router. EventName per binding fires on trigger.
	applyKeybindings(c, cfg.Keybindings)

	// Context menus — register each named menu + install the shared
	// relay. EventTemplate per menu fires on item click with the
	// click's context data forwarded.
	applyContextMenus(c, cfg.ContextMenus)

	// App menu — fire menu.set_app_menu on darwin. No-op elsewhere.
	applyAppMenu(c, cfg.AppMenu)

	return core.Ok(nil)
}

// registerAndPreCreateWindow fires the standard registry + pre-create
// sequence for a single declared window:
//
//   - window.register with KindWebview so taskSetVisibility can lazy-
//     mount the window on first show
//   - window.open as hidden so the platform window exists before the
//     first show click (no cold-start render delay)
//   - window.set_close_behavior to CloseBehaviorHide if HideOnClose set
//   - window.set_content_protection if ContentProtection set
//
// Failures inside one step bubble up but pre-create remains best-effort
// for the steps after registration — if the consumer set HideOnClose
// and the close-behavior action fails, the window still opens; the
// consumer can retry the behaviour action later if needed.
func registerAndPreCreateWindow(c *core.Core, w *window.Window) core.Result {
	ctx := core.Background()
	if r := c.Action("window.register").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRegisterWindow{Window: w, Kind: window.KindWebview}},
	)); !r.OK {
		return r
	}
	// Pre-create hidden — clone the descriptor so Hidden=true doesn't
	// mutate the caller-supplied Window pointer.
	openSpec := *w
	openSpec.Hidden = true
	if r := c.Action("window.open").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{Window: &openSpec}},
	)); !r.OK {
		return r
	}
	if w.HideOnClose {
		c.Action("window.set_close_behavior").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: window.TaskSetCloseBehavior{
				Name:     w.Name,
				Behavior: window.CloseBehaviorHide,
			}},
		))
	}
	if w.ContentProtection {
		c.Action("window.set_content_protection").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: window.TaskSetContentProtection{
				Name:       w.Name,
				Protection: true,
			}},
		))
	}
	return core.Ok(nil)
}

// registerSubservice mirrors lthn/desktop's registerStartedCoreGUIService
// helper — runs the factory, registers the result on the Core, and
// fires OnStartup on the returned service if it implements Startable.
// Idempotent: returns OK if the service is already registered.
func registerSubservice(ctx context.Context, c *core.Core, name string, factory func(*core.Core) core.Result) core.Result {
	if c.Service(name).OK {
		return core.Ok(nil)
	}
	r := factory(c)
	if !r.OK {
		return core.Fail(core.E("gui.Service.start", "register "+name+" failed", r.Value.(error)))
	}
	if r.Value != nil {
		if rr := c.RegisterService(name, r.Value); !rr.OK {
			return rr
		}
	}
	svc := c.Service(name)
	if !svc.OK {
		return core.Ok(nil)
	}
	startable, ok := svc.Value.(core.Startable)
	if !ok {
		return core.Ok(nil)
	}
	if sr := startable.OnStartup(ctx); !sr.OK {
		return core.Fail(core.E("gui.Service.start", "startup "+name+" failed", sr.Value.(error)))
	}
	return core.Ok(nil)
}

// containerFactory returns the registration factory for the container
// sub-service with config-driven options.
func containerFactory(cfg ContainerConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: container.NewService(c, container.TIMOptions{
			Image:   cfg.Image,
			Name:    cfg.Name,
			DataDir: cfg.DataDir,
			Command: cfg.Command,
		}), OK: true}
	}
}

// p2pFactory returns the registration factory for the p2p sub-service
// with config-driven options.
func p2pFactory(cfg P2PConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: p2p.NewService(c, p2p.Options{
			ListenAddr: cfg.ListenAddr,
			PeerAddrs:  cfg.PeerAddrs,
			NodeID:     cfg.NodeID,
		}), OK: true}
	}
}
