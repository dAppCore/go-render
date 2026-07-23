package window

import "time"

type WindowInfo struct {
	Name        string  `json:"name"`
	Title       string  `json:"title"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Opacity     float64 `json:"opacity"`
	Maximized   bool    `json:"maximized"`
	Focused     bool    `json:"focused"`
	Visible     bool    `json:"visible"`
	Minimised   bool    `json:"minimised"`
	Fullscreen  bool    `json:"fullscreen"`
	AlwaysOnTop bool    `json:"always_on_top"`
}

type QueryWindowList struct{}

type QueryWindowByName struct{ Name string }

type QueryConfig struct{}

type TaskOpenWindow struct {
	Window  *Window
	Options []WindowOption
}

type TaskCloseWindow struct{ Name string }

type TaskSetPosition struct {
	Name string
	X, Y int
}

type TaskSetSize struct {
	Name          string
	Width, Height int
}

type TaskMaximise struct{ Name string }

type TaskMinimise struct{ Name string }

type TaskFocus struct{ Name string }

type TaskRestore struct{ Name string }

type TaskSetTitle struct {
	Name  string
	Title string
}

type TaskSetAlwaysOnTop struct {
	Name        string
	AlwaysOnTop bool
}

type TaskSetOpacity struct {
	Name    string
	Opacity float64
}

type TaskSetBackgroundColour struct {
	Name  string
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type TaskSetVisibility struct {
	Name    string
	Visible bool
}

// WindowKind classifies how a window is realised by the platform.
// Registered windows declare their kind up-front so taskSetVisibility
// knows whether to lazily mount a WebView (KindWebview) or skip to the
// systray path (KindTray) on first show.
type WindowKind int

const (
	KindWebview WindowKind = iota // standard HTML window — Wails WebviewWindow
	KindTray                      // systray icon — no WebView lifecycle
)

// TaskRegisterWindow stores a Window descriptor in the service registry
// without opening it. taskSetVisibility consults the registry on first
// show, creating the platform window lazily via the existing taskOpenWindow
// path. Issued via the action bus as `window.register`.
//
// Validation enforces the discriminator: KindWebview requires non-empty
// URL, KindTray requires empty URL. Duplicate names rejected.
type TaskRegisterWindow struct {
	Window *Window
	Kind   WindowKind
}

// CloseBehavior governs what happens when the OS / user requests a
// window close (window-control button, Cmd+W, Alt+F4). One of:
//
//	CloseBehaviorDestroy — let the close proceed (default).
//	CloseBehaviorHide    — intercept, hide the window, cancel the
//	                       close. Window stays registered so it can
//	                       be shown again later. Tray-rooted apps use
//	                       this so the popover survives clicking 'x'.
//	CloseBehaviorQuit    — intercept and call app.Quit() instead.
type CloseBehavior string

const (
	CloseBehaviorDestroy CloseBehavior = "destroy"
	CloseBehaviorHide    CloseBehavior = "hide"
	CloseBehaviorQuit    CloseBehavior = "quit"
)

// TaskSetCloseBehavior installs a behaviour-bearing hook on a
// previously-registered window's close event. Declarative — the
// behaviour intent is the payload; the platform decides whether to
// wire RegisterHook, observers, or platform-native cancel surfaces.
//
//	c.Action("window.set_close_behavior").Run(ctx, NewOptions(
//	    Option{Key: "task", Value: TaskSetCloseBehavior{
//	        Name: "tray", Behavior: CloseBehaviorHide,
//	    }},
//	))
type TaskSetCloseBehavior struct {
	Name     string
	Behavior CloseBehavior
}

type TaskFullscreen struct {
	Name       string
	Fullscreen bool
}

type QueryLayoutList struct{}

type QueryLayoutGet struct{ Name string }

type TaskSaveLayout struct{ Name string }

type TaskRestoreLayout struct{ Name string }

type TaskDeleteLayout struct{ Name string }

type TaskTileWindows struct {
	Mode    string   // "left-right", "grid", "left-half", "right-half", etc.
	Windows []string // window names; empty = all
}

type TaskStackWindows struct {
	Windows []string // window names; empty = all
	OffsetX int
	OffsetY int
}

type TaskSnapWindow struct {
	Name     string // window name
	Position string // "left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right", "center"
}

type TaskApplyWorkflow struct {
	Workflow string
	Windows  []string // window names; empty = all
}

type TaskSaveConfig struct{ Config map[string]any }

type ActionWindowOpened struct{ Name string }
type ActionWindowClosed struct{ Name string }

type ActionWindowMoved struct {
	Name string
	X, Y int
}

type ActionWindowResized struct {
	Name          string
	Width, Height int
}

type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }

// Window state transitions — fired when the OS or user changes a
// window's visibility / minimisation / maximisation / fullscreen
// state. Consumers can react to these (e.g. pause rendering on
// hide, reload on show, save layout on state change) without
// polling.
type ActionWindowHidden struct{ Name string }
type ActionWindowShown struct{ Name string }
type ActionWindowMinimised struct{ Name string }
type ActionWindowUnminimised struct{ Name string }
type ActionWindowMaximised struct{ Name string }
type ActionWindowUnmaximised struct{ Name string }
type ActionWindowFullscreened struct{ Name string }
type ActionWindowUnfullscreened struct{ Name string }

// ActionWindowRuntimeReady fires after the WebView JS runtime has
// loaded and the Wails bridge is ready to accept IPC calls. Consumers
// that need to push initial state into a window should wait for this
// (rather than depending on window.create which fires before the
// frontend has mounted).
type ActionWindowRuntimeReady struct{ Name string }

// DropTarget describes the HTML element that received an OS-level
// file drop. Populated from Wails's DropTargetDetails when the drop
// landed on an element carrying the data-file-drop-target attribute;
// nil when the drop landed on a non-target region (e.g. window
// chrome).
type DropTarget struct {
	ID         string            `json:"id,omitempty"`
	X          int               `json:"x"`
	Y          int               `json:"y"`
	ClassList  []string          `json:"classList,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type ActionFilesDropped struct {
	Name  string   `json:"name"` // window name
	Paths []string `json:"paths"`
	// Target is the element that received the drop (nil when the drop
	// landed outside any data-file-drop-target region).
	Target *DropTarget `json:"target,omitempty"`
	// TargetID mirrors Target.ID for back-compat. Empty when Target
	// is nil. Retained so existing consumers keep working without a
	// nil check; new consumers should read Target for full context.
	TargetID string `json:"targetId,omitempty"`
}

// --- Zoom ---

type QueryWindowZoom struct{ Name string }

type TaskSetZoom struct {
	Name          string
	Magnification float64
}

type TaskZoomIn struct{ Name string }

type TaskZoomOut struct{ Name string }

type TaskZoomReset struct{ Name string }

// --- Content ---

type TaskSetURL struct {
	Name string
	URL  string
}

type TaskSetHTML struct {
	Name string
	HTML string
}

type TaskExecJS struct {
	Name string
	JS   string
}

// TaskEvalJS evaluates JS in the named window and waits for the
// result via Wails' public Events bus. The handler wraps the JS
// body in an IIFE that imports @wailsio/runtime and emits
// "lthn:eval-reply" with {reqId, result/error}; the service's
// listener completes a pending channel keyed by reqId and the
// caller gets the value back synchronously.
//
// Replaces the older HTTP-fetchback path that timed out silently
// against the bridge HTTP server because the WebView's cross-
// origin POSTs failed the bridge's DNS-rebind defence. Uses
// Wails' built-in Events bus so no cgo + no CORS + no script-
// message-handler vendoring.
//
// Timeout defaults to 5 seconds when zero. Result and Err are
// mutually exclusive in the response — non-empty Err carries the
// JS-side exception string or the timeout sentinel.
type TaskEvalJS struct {
	Name    string
	JS      string
	Timeout time.Duration
}

// EvalJSResult is the return wrapper for TaskEvalJS. The action
// returns this struct in core.Result.Value on success and on
// JS-side errors (the OK flag distinguishes); only platform-
// level failures (window not found, listener not registered)
// land as core.Fail.
type EvalJSResult struct {
	ReqID  string `json:"reqId"`
	Result any    `json:"result,omitempty"`
	Err    string `json:"err,omitempty"`
}

// --- State toggles ---

type TaskToggleFullscreen struct{ Name string }

type TaskToggleMaximise struct{ Name string }

// --- Bounds ---

type QueryWindowBounds struct{ Name string }

type WindowBounds struct {
	X, Y, Width, Height int
}

type TaskSetBounds struct {
	Name                string
	X, Y, Width, Height int
}

// --- Content protection ---

type TaskSetContentProtection struct {
	Name       string
	Protection bool
}

// --- Flash ---

type TaskFlash struct {
	Name    string
	Enabled bool
}

// --- Print ---

type TaskPrint struct{ Name string }

// --- Smart layout ---

type TaskLayoutBesideEditor struct {
	Name   string  `json:"name"`
	Editor string  `json:"editor,omitempty"`
	Side   string  `json:"side,omitempty"`
	Ratio  float64 `json:"ratio,omitempty"`
}

type LayoutBesideEditorResult struct {
	Editor       string       `json:"editor"`
	EditorBounds WindowBounds `json:"editor_bounds"`
	WindowBounds WindowBounds `json:"window_bounds"`
	Side         string       `json:"side"`
	ScreenID     string       `json:"screen_id,omitempty"`
}

type TaskLayoutSuggest struct {
	ScreenID    string `json:"screen_id,omitempty"`
	WindowCount int    `json:"window_count,omitempty"`
}

type LayoutSuggestion struct {
	Mode     string `json:"mode"`
	Reason   string `json:"reason"`
	ScreenID string `json:"screen_id,omitempty"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type TaskScreenFindSpace struct {
	ScreenID string `json:"screen_id,omitempty"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Padding  int    `json:"padding,omitempty"`
}

type ScreenSpace struct {
	ScreenID string `json:"screen_id,omitempty"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type TaskWindowArrangePair struct {
	Primary   string  `json:"primary"`
	Secondary string  `json:"secondary"`
	ScreenID  string  `json:"screen_id,omitempty"`
	Ratio     float64 `json:"ratio,omitempty"`
}

type PairArrangement struct {
	Primary     WindowBounds `json:"primary"`
	Secondary   WindowBounds `json:"secondary"`
	Orientation string       `json:"orientation"`
	ScreenID    string       `json:"screen_id,omitempty"`
}
