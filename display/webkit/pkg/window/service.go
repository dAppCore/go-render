package window

import (
	"context"
	"sync"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
	"dappco.re/go/render/display/webkit/pkg/screen"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	// specs stores registered Window descriptors keyed by name.
	// taskSetVisibility consults this map on first show to mount the
	// platform window lazily via taskOpenWindow, rather than requiring
	// every window to be explicitly opened before its visibility can
	// be set. Populated via the window.register action.
	specs   map[string]registeredSpec
	specsMu sync.RWMutex
	// pendingEvals routes JS-side replies (delivered via Wails'
	// Events bus on "lthn:eval-reply") back to the taskEvalJS caller
	// blocked on a channel. Keyed by reqId.
	pendingEvals   map[string]chan EvalJSResult
	pendingEvalsMu sync.Mutex
	// nextEvalSeq is the monotonic counter behind generated reqIds.
	// Combined with a per-service prefix so two Service instances in
	// the same process can't collide.
	nextEvalSeq uint64
	evalPrefix  string
}

// registeredSpec pairs a Window descriptor with its WindowKind so the
// service can route first-show through the right platform path (webview
// mount vs systray show).
type registeredSpec struct {
	Window *Window
	Kind   WindowKind
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	// Query config — display registers its handler before us (registration order guarantee).
	// If display is not registered, OK=false and we skip config.
	r := s.Core().QUERY(QueryConfig{})
	if r.OK {
		if windowConfig, ok := r.Value.(map[string]any); ok {
			s.applyConfig(windowConfig)
		}
	}

	s.Core().RegisterQuery(s.handleQuery)
	s.registerTaskActions()
	// Bind the JS-side eval reply listener if the platform supports
	// it. WailsPlatform wires app.Event.On(EvalJSEventName) →
	// CompleteEval; MockPlatform is a no-op and tests call
	// CompleteEval directly.
	if binder, ok := s.platform.(EvalReplyBinder); ok && binder != nil {
		binder.BindEvalReply(func(reqID string, result any, errStr string) {
			s.CompleteEval(reqID, result, errStr)
		})
	}
	return core.Result{OK: true}
}

// EvalReplyBinder is the optional extension a Platform implements
// to route Wails custom-event replies back to the Service's eval
// machinery. The Service calls BindEvalReply once at OnStartup;
// the platform invokes the callback every time a "lthn:eval-reply"
// event arrives on the Wails Events bus.
//
// MockPlatform deliberately does NOT implement this so tests can
// drive Service.CompleteEval directly without an event-bus
// dependency.
type EvalReplyBinder interface {
	BindEvalReply(cb func(reqID string, result any, errStr string))
}

// CustomEventBinder is the optional extension a Platform implements
// to expose a generic Wails custom-event subscription to consumers
// outside pkg/window. Consumers (e.g. pkg/bridge for its console +
// error capture) call Service.SubscribeEvent(name, cb); the service
// delegates to the platform if supported. Each (name, cb) registers
// a fresh listener — no de-dup, no unsubscribe (consumers bind once
// at their own OnStartup and live for the process lifetime).
//
// MockPlatform deliberately does NOT implement this so tests can
// drive the consumer callback directly without an event-bus
// dependency. SubscribeEvent on an unsupported platform is a no-op
// + returns false so consumers can fall back to alternative paths
// if they care.
type CustomEventBinder interface {
	BindCustomEvent(name string, cb func(data any))
}

// SubscribeEvent registers a callback for a Wails custom event with
// the given name. Returns true when the platform supports event
// binding (production WailsPlatform); false on platforms that don't
// (test MockPlatform). Idempotency is the platform's concern — most
// implementations stack listeners, so repeated SubscribeEvent calls
// fan-out per emit.
//
// Usage example (bridge.OnStartup):
//
//	if !windowSvc.SubscribeEvent("lthn:console", s.handleConsoleEvent) {
//	    // platform doesn't support custom events — degrade gracefully
//	}
func (s *Service) SubscribeEvent(name string, cb func(data any)) bool {
	if s == nil || cb == nil || name == "" {
		return false
	}
	binder, ok := s.platform.(CustomEventBinder)
	if !ok || binder == nil {
		return false
	}
	binder.BindCustomEvent(name, cb)
	return true
}

func (s *Service) applyConfig(configData map[string]any) {
	if width, ok := configData["default_width"]; ok {
		if width, ok := width.(int); ok {
			s.manager.SetDefaultWidth(width)
		}
	}
	if height, ok := configData["default_height"]; ok {
		if height, ok := height.(int); ok {
			s.manager.SetDefaultHeight(height)
		}
	}
	if stateFile, ok := configData["state_file"]; ok {
		if stateFile, ok := stateFile.(string); ok {
			s.manager.State().SetPath(stateFile)
		}
	}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q := q.(type) {
	case QueryWindowList:
		return core.Result{Value: s.queryWindowList(), OK: true}
	case QueryWindowByName:
		return core.Result{Value: s.queryWindowByName(q.Name), OK: true}
	case QueryLayoutList:
		return core.Result{Value: s.manager.Layout().ListLayouts(), OK: true}
	case QueryLayoutGet:
		l, ok := s.manager.Layout().GetLayout(q.Name)
		if !ok {
			return core.Result{Value: (*Layout)(nil), OK: true}
		}
		return core.Result{Value: &l, OK: true}
	case QueryWindowZoom:
		return s.queryWindowZoom(q.Name)
	case QueryWindowBounds:
		return s.queryWindowBounds(q.Name)
	default:
		return core.Result{}
	}
}

func (s *Service) queryWindowList() []WindowInfo {
	names := s.manager.List()
	result := make([]WindowInfo, 0, len(names))
	for _, name := range names {
		if pw, ok := s.manager.Get(name); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			result = append(result, WindowInfo{
				Name: name, Title: pw.Title(), X: x, Y: y, Width: w, Height: h,
				Opacity:     pw.GetOpacity(),
				Maximized:   pw.IsMaximised(),
				Focused:     pw.IsFocused(),
				Visible:     pw.IsVisible(),
				Minimised:   pw.IsMinimised(),
				Fullscreen:  pw.IsFullscreen(),
				AlwaysOnTop: pw.IsAlwaysOnTop(),
			})
		}
	}
	return result
}

func (s *Service) queryWindowByName(name string) *WindowInfo {
	pw, ok := s.manager.Get(name)
	if !ok {
		return nil
	}
	x, y := pw.Position()
	w, h := pw.Size()
	return &WindowInfo{
		Name: name, Title: pw.Title(), X: x, Y: y, Width: w, Height: h,
		Opacity:     pw.GetOpacity(),
		Maximized:   pw.IsMaximised(),
		Focused:     pw.IsFocused(),
		Visible:     pw.IsVisible(),
		Minimised:   pw.IsMinimised(),
		Fullscreen:  pw.IsFullscreen(),
		AlwaysOnTop: pw.IsAlwaysOnTop(),
	}
}

// --- Action Registration ---

// registerTaskActions registers all window task handlers as named Core actions.
func (s *Service) registerTaskActions() {
	c := s.Core()
	c.Action("window.open", func(_ context.Context, opts core.Options) core.Result {
		t := taskOpenWindowFromOptions(opts)
		return s.taskOpenWindow(t)
	})
	c.Action("gui.window.create", func(_ context.Context, opts core.Options) core.Result {
		t := taskOpenWindowFromOptions(opts)
		return s.taskOpenWindow(t)
	})
	c.Action("gui.window.open", func(_ context.Context, opts core.Options) core.Result {
		t := taskOpenWindowFromOptions(opts)
		return s.taskOpenWindow(t)
	})
	c.Action("window.close", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskCloseWindow]("window.close", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskCloseWindow(t.Name))
	})
	c.Action("window.register", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskRegisterWindow]("window.register", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskRegisterWindow(t.Window, t.Kind))
	})
	c.Action("window.set_position", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetPosition]("window.set_position", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetPosition(t.Name, t.X, t.Y))
	})
	c.Action("window.set_size", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetSize]("window.set_size", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetSize(t.Name, t.Width, t.Height))
	})
	c.Action("window.maximise", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskMaximise]("window.maximise", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskMaximise(t.Name))
	})
	c.Action("window.minimise", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskMinimise]("window.minimise", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskMinimise(t.Name))
	})
	c.Action("window.focus", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskFocus]("window.focus", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskFocus(t.Name))
	})
	c.Action("window.restore", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskRestore]("window.restore", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskRestore(t.Name))
	})
	c.Action("window.set_title", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetTitle]("window.set_title", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetTitle(t.Name, t.Title))
	})
	c.Action("window.set_always_on_top", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetAlwaysOnTop]("window.set_always_on_top", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetAlwaysOnTop(t.Name, t.AlwaysOnTop))
	})
	c.Action("window.set_opacity", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetOpacity]("window.set_opacity", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetOpacity(t.Name, t.Opacity))
	})
	c.Action("window.set_background_colour", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetBackgroundColour]("window.set_background_colour", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetBackgroundColour(t.Name, t.Red, t.Green, t.Blue, t.Alpha))
	})
	c.Action("window.set_visibility", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetVisibility]("window.set_visibility", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetVisibility(t.Name, t.Visible))
	})
	c.Action("window.set_close_behavior", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetCloseBehavior]("window.set_close_behavior", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetCloseBehavior(t.Name, t.Behavior))
	})
	c.Action("window.fullscreen", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskFullscreen]("window.fullscreen", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskFullscreen(t.Name, t.Fullscreen))
	})
	c.Action("window.save_layout", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSaveLayout]("window.save_layout", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSaveLayout(t.Name))
	})
	c.Action("window.restore_layout", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskRestoreLayout]("window.restore_layout", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskRestoreLayout(t.Name))
	})
	c.Action("window.delete_layout", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskDeleteLayout]("window.delete_layout", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		s.manager.Layout().DeleteLayout(t.Name)
		return core.Result{OK: true}
	})
	c.Action("window.tile_windows", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskTileWindows]("window.tile_windows", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskTileWindows(t.Mode, t.Windows))
	})
	c.Action("window.stack_windows", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskStackWindows]("window.stack_windows", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskStackWindows(t.Windows, t.OffsetX, t.OffsetY))
	})
	c.Action("window.snap_window", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSnapWindow]("window.snap_window", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSnapWindow(t.Name, t.Position))
	})
	c.Action("window.apply_workflow", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskApplyWorkflow]("window.apply_workflow", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskApplyWorkflow(t.Workflow, t.Windows))
	})
	c.Action("window.layout_beside_editor", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskLayoutBesideEditor]("window.layout_beside_editor", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		result, err := s.taskLayoutBesideEditor(t)
		return core.Result{}.New(result, err)
	})
	c.Action("window.layout_suggest", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskLayoutSuggest]("window.layout_suggest", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: s.taskLayoutSuggest(t), OK: true}
	})
	c.Action("window.find_space", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskScreenFindSpace]("window.find_space", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: s.taskScreenFindSpace(t), OK: true}
	})
	c.Action("window.arrange_pair", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskWindowArrangePair]("window.arrange_pair", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		result, err := s.taskWindowArrangePair(t)
		return core.Result{}.New(result, err)
	})
	c.Action("window.set_zoom", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetZoom]("window.set_zoom", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetZoom(t.Name, t.Magnification))
	})
	c.Action("window.zoom_in", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskZoomIn]("window.zoom_in", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskZoomIn(t.Name))
	})
	c.Action("window.zoom_out", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskZoomOut]("window.zoom_out", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskZoomOut(t.Name))
	})
	c.Action("window.zoom_reset", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskZoomReset]("window.zoom_reset", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskZoomReset(t.Name))
	})
	c.Action("window.set_url", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetURL]("window.set_url", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetURL(t.Name, t.URL))
	})
	c.Action("window.set_html", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetHTML]("window.set_html", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetHTML(t.Name, t.HTML))
	})
	c.Action("window.exec_js", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskExecJS]("window.exec_js", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskExecJS(t.Name, t.JS))
	})
	c.Action("window.eval_js", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskEvalJS]("window.eval_js", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		res, evalErr := s.taskEvalJS(t.Name, t.JS, t.Timeout)
		if evalErr != nil {
			return core.Result{Value: evalErr, OK: false}
		}
		return core.Result{Value: res, OK: true}
	})
	c.Action("window.toggle_fullscreen", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskToggleFullscreen]("window.toggle_fullscreen", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskToggleFullscreen(t.Name))
	})
	c.Action("window.toggle_maximise", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskToggleMaximise]("window.toggle_maximise", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskToggleMaximise(t.Name))
	})
	c.Action("window.set_bounds", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetBounds]("window.set_bounds", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetBounds(t.Name, t.X, t.Y, t.Width, t.Height))
	})
	c.Action("window.set_content_protection", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskSetContentProtection]("window.set_content_protection", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskSetContentProtection(t.Name, t.Protection))
	})
	c.Action("window.flash", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskFlash]("window.flash", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskFlash(t.Name, t.Enabled))
	})
	c.Action("window.print", func(_ context.Context, opts core.Options) core.Result {
		t, err := taskFromOptions[TaskPrint]("window.print", opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.taskPrint(t.Name))
	})
}

func taskFromOptions[T any](action string, opts core.Options) (T, resultFailure) {
	var zero T
	task := opts.Get("task")
	if !task.OK {
		return zero, core.E(action, "missing task payload", nil)
	}
	switch value := task.Value.(type) {
	case T:
		return value, nil
	case map[string]any:
		var decoded T
		if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
			return decoded, nil
		}
	}
	return zero, core.E(action, "invalid task payload", nil)
}

func taskOpenWindowFromOptions(opts core.Options) TaskOpenWindow {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskOpenWindow:
			return value
		case map[string]any:
			var decoded TaskOpenWindow
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}

	var decoded TaskOpenWindow
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskOpenWindow{}
}

func optsToMap(opts core.Options) map[string]any {
	items := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		items[item.Key] = item.Value
	}
	return items
}

func (s *Service) primaryScreenArea() (int, int, int, int) {
	const fallbackX = 0
	const fallbackY = 0
	const fallbackWidth = 1920
	const fallbackHeight = 1080

	r := s.Core().QUERY(screen.QueryPrimary{})
	if !r.OK {
		return fallbackX, fallbackY, fallbackWidth, fallbackHeight
	}

	primary, ok := r.Value.(*screen.Screen)
	if !ok || primary == nil {
		return fallbackX, fallbackY, fallbackWidth, fallbackHeight
	}

	x := primary.WorkArea.X
	y := primary.WorkArea.Y
	width := primary.WorkArea.Width
	height := primary.WorkArea.Height
	if width <= 0 || height <= 0 {
		x = primary.Bounds.X
		y = primary.Bounds.Y
		width = primary.Bounds.Width
		height = primary.Bounds.Height
	}
	if width <= 0 || height <= 0 {
		return fallbackX, fallbackY, fallbackWidth, fallbackHeight
	}

	return x, y, width, height
}

func (s *Service) taskOpenWindow(t TaskOpenWindow) core.Result {
	spec, err := s.buildWindowSpec(t)
	if err != nil {
		return core.Result{Value: err, OK: false}
	}
	if err := s.prepareWindowSpec(spec); err != nil {
		return core.Result{Value: err, OK: false}
	}

	pw, err := s.manager.Create(spec)
	if err != nil {
		return core.Result{Value: err, OK: false}
	}
	x, y := pw.Position()
	w, h := pw.Size()
	info := WindowInfo{
		Name: pw.Name(), Title: pw.Title(), X: x, Y: y, Width: w, Height: h,
		Opacity:     pw.GetOpacity(),
		Maximized:   pw.IsMaximised(),
		Focused:     pw.IsFocused(),
		Visible:     pw.IsVisible(),
		Minimised:   pw.IsMinimised(),
		Fullscreen:  pw.IsFullscreen(),
		AlwaysOnTop: pw.IsAlwaysOnTop(),
	}

	// Attach platform event listeners that convert to IPC actions
	s.trackWindow(pw)

	// Broadcast to all listeners
	coreutil.DispatchAction(s.Core(), "window.create", ActionWindowOpened{Name: pw.Name()})
	return core.Result{Value: info, OK: true}
}

// trackWindow attaches platform event listeners that emit IPC actions.
func (s *Service) trackWindow(pw PlatformWindow) {
	pw.OnWindowEvent(func(e WindowEvent) {
		switch e.Type {
		case "focus":
			coreutil.DispatchAction(s.Core(), "window.focus", ActionWindowFocused{Name: e.Name})
		case "blur":
			coreutil.DispatchAction(s.Core(), "window.blur", ActionWindowBlurred{Name: e.Name})
		case "move":
			if data := e.Data; data != nil {
				x, _ := data["x"].(int)
				y, _ := data["y"].(int)
				coreutil.DispatchAction(s.Core(), "window.move", ActionWindowMoved{Name: e.Name, X: x, Y: y})
			}
			// Auto-persist OS-driven moves — without this, dragging a
			// window only saves when the window is explicitly closed
			// (which never happens for HideOnClose tray-rooted apps).
			// CaptureState writes via a 500ms debounced timer so rapid
			// drags coalesce into one save.
			if pw, ok := s.manager.Get(e.Name); ok {
				s.manager.State().CaptureState(pw)
			}
		case "resize":
			if data := e.Data; data != nil {
				w, _ := data["w"].(int)
				h, _ := data["h"].(int)
				coreutil.DispatchAction(s.Core(), "window.resize", ActionWindowResized{Name: e.Name, Width: w, Height: h})
			}
			// Auto-persist OS-driven resize — same rationale as move.
			if pw, ok := s.manager.Get(e.Name); ok {
				s.manager.State().CaptureState(pw)
			}
		case "close":
			coreutil.DispatchAction(s.Core(), "window.closeEvent", ActionWindowClosed{Name: e.Name})
		case "hide":
			coreutil.DispatchAction(s.Core(), "window.hide", ActionWindowHidden{Name: e.Name})
		case "show":
			coreutil.DispatchAction(s.Core(), "window.show", ActionWindowShown{Name: e.Name})
		case "minimise":
			coreutil.DispatchAction(s.Core(), "window.minimise", ActionWindowMinimised{Name: e.Name})
		case "unminimise":
			coreutil.DispatchAction(s.Core(), "window.unminimise", ActionWindowUnminimised{Name: e.Name})
		case "maximise":
			coreutil.DispatchAction(s.Core(), "window.maximise", ActionWindowMaximised{Name: e.Name})
		case "unmaximise":
			coreutil.DispatchAction(s.Core(), "window.unmaximise", ActionWindowUnmaximised{Name: e.Name})
		case "fullscreen":
			coreutil.DispatchAction(s.Core(), "window.fullscreen", ActionWindowFullscreened{Name: e.Name})
		case "unfullscreen":
			coreutil.DispatchAction(s.Core(), "window.unfullscreen", ActionWindowUnfullscreened{Name: e.Name})
		case "ready":
			coreutil.DispatchAction(s.Core(), "window.ready", ActionWindowRuntimeReady{Name: e.Name})
		}
	})
	pw.OnFileDrop(func(paths []string, target *DropTarget) {
		event := ActionFilesDropped{
			Name:   pw.Name(),
			Paths:  paths,
			Target: target,
		}
		if target != nil {
			event.TargetID = target.ID
		}
		coreutil.DispatchAction(s.Core(), "window.fileDrop", event)
	})
}

func (s *Service) taskCloseWindow(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskClose", "window not found: "+name, nil)
	}
	// Persist state BEFORE closing (spec requirement)
	s.manager.State().CaptureState(pw)
	pw.Close()
	s.manager.Remove(name)
	coreutil.DispatchAction(s.Core(), "window.taskClose", ActionWindowClosed{Name: name})
	return nil
}

func (s *Service) taskSetPosition(name string, x, y int) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetPosition", "window not found: "+name, nil)
	}
	pw.SetPosition(x, y)
	s.manager.State().UpdatePosition(name, x, y)
	return nil
}

func (s *Service) taskSetSize(name string, width, height int) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetSize", "window not found: "+name, nil)
	}
	pw.SetSize(width, height)
	s.manager.State().UpdateSize(name, width, height)
	return nil
}

func (s *Service) taskMaximise(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskMaximise", "window not found: "+name, nil)
	}
	pw.Maximise()
	s.manager.State().UpdateMaximized(name, true)
	return nil
}

func (s *Service) taskMinimise(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskMinimise", "window not found: "+name, nil)
	}
	pw.Minimise()
	return nil
}

func (s *Service) taskFocus(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskFocus", "window not found: "+name, nil)
	}
	pw.Focus()
	return nil
}

func (s *Service) taskRestore(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskRestore", "window not found: "+name, nil)
	}
	pw.Restore()
	s.manager.State().UpdateMaximized(name, false)
	return nil
}

func (s *Service) taskSetTitle(name, title string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetTitle", "window not found: "+name, nil)
	}
	pw.SetTitle(title)
	return nil
}

func (s *Service) taskSetAlwaysOnTop(name string, alwaysOnTop bool) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetAlwaysOnTop", "window not found: "+name, nil)
	}
	pw.SetAlwaysOnTop(alwaysOnTop)
	return nil
}

func (s *Service) taskSetOpacity(name string, opacity float64) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetOpacity", "window not found: "+name, nil)
	}
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}
	pw.SetOpacity(opacity)
	return nil
}

func (s *Service) taskSetBackgroundColour(name string, red, green, blue, alpha uint8) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetBackgroundColour", "window not found: "+name, nil)
	}
	pw.SetBackgroundColour(red, green, blue, alpha)
	return nil
}

func (s *Service) taskSetVisibility(name string, visible bool) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		// First-show path: the window isn't mounted yet. Hiding a
		// never-shown window is a no-op success; showing requires a
		// registered spec to know what to mount.
		if !visible {
			return nil
		}
		s.specsMu.RLock()
		spec, hasSpec := s.specs[name]
		s.specsMu.RUnlock()
		if !hasSpec {
			return core.E("window.taskSetVisibility", "unregistered window: "+name, nil)
		}
		switch spec.Kind {
		case KindTray:
			// Tray windows are realised by pkg/systray, not by this
			// service. Visibility on a tray-classified name is accepted
			// without WebView mount — the systray service owns the
			// actual icon lifecycle.
			return nil
		case KindWebview:
			// Reuse taskOpenWindow for the full create + register flow
			// (buildWindowSpec → prepareWindowSpec → manager.Create).
			openResult := s.taskOpenWindow(TaskOpenWindow{Window: spec.Window})
			if !openResult.OK {
				if e, isErr := openResult.Value.(error); isErr {
					return e
				}
				return core.E("window.taskSetVisibility", "failed to open registered window: "+name, nil)
			}
			refreshed, hit := s.manager.Get(name)
			if !hit {
				return core.E("window.taskSetVisibility", "window mount completed but manager has no entry for "+name, nil)
			}
			pw = refreshed
		default:
			return core.E("window.taskSetVisibility", "unknown WindowKind for "+name, nil)
		}
	}
	pw.SetVisibility(visible)
	return nil
}

// taskRegisterWindow stores a Window descriptor in the service registry.
// Called by consumers at boot to declare windows that may be opened
// lazily by taskSetVisibility. Validates the WindowKind discriminator:
// KindWebview requires non-empty URL; KindTray requires empty URL.
// Duplicate names rejected — consumers re-register at every boot, so
// any duplicate indicates a coding error.
func (s *Service) taskRegisterWindow(w *Window, kind WindowKind) resultFailure {
	if w == nil {
		return core.E("window.taskRegisterWindow", "nil Window", nil)
	}
	if w.Name == "" {
		return core.E("window.taskRegisterWindow", "Window.Name is required", nil)
	}
	if kind == KindWebview && w.URL == "" {
		return core.E("window.taskRegisterWindow", "KindWebview requires non-empty URL", nil)
	}
	if kind == KindTray && w.URL != "" {
		return core.E("window.taskRegisterWindow", "KindTray must have empty URL", nil)
	}
	s.specsMu.Lock()
	defer s.specsMu.Unlock()
	if _, exists := s.specs[w.Name]; exists {
		return core.E("window.taskRegisterWindow", "window already registered: "+w.Name, nil)
	}
	s.specs[w.Name] = registeredSpec{Window: w, Kind: kind}
	return nil
}

// taskSetCloseBehavior installs the requested CloseBehavior on a
// tracked window. Looks up by name, delegates to the platform.
func (s *Service) taskSetCloseBehavior(name string, behavior CloseBehavior) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetCloseBehavior", "window not found: "+name, nil)
	}
	pw.SetCloseBehavior(behavior)
	return nil
}

func (s *Service) taskFullscreen(name string, fullscreen bool) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskFullscreen", "window not found: "+name, nil)
	}
	if fullscreen {
		pw.Fullscreen()
	} else {
		pw.UnFullscreen()
	}
	return nil
}

func (s *Service) taskSaveLayout(name string) resultFailure {
	windows := s.queryWindowList()
	states := make(map[string]WindowState, len(windows))
	for _, w := range windows {
		states[w.Name] = WindowState{
			X: w.X, Y: w.Y, Width: w.Width, Height: w.Height,
			Maximized: w.Maximized,
		}
	}
	return s.manager.Layout().SaveLayout(name, states)
}

func (s *Service) taskRestoreLayout(name string) resultFailure {
	layout, ok := s.manager.Layout().GetLayout(name)
	if !ok {
		return core.E("window.taskRestoreLayout", "layout not found: "+name, nil)
	}
	for winName, state := range layout.Windows {
		pw, found := s.manager.Get(winName)
		if !found {
			continue
		}
		pw.SetPosition(state.X, state.Y)
		pw.SetSize(state.Width, state.Height)
		if state.Maximized {
			pw.Maximise()
		} else {
			pw.Restore()
		}
		s.manager.State().CaptureState(pw)
	}
	return nil
}

var tileModeMap = map[string]TileMode{
	"left-half": TileModeLeftHalf, "right-half": TileModeRightHalf,
	"top-half": TileModeTopHalf, "bottom-half": TileModeBottomHalf,
	"top-left": TileModeTopLeft, "top-right": TileModeTopRight,
	"bottom-left": TileModeBottomLeft, "bottom-right": TileModeBottomRight,
	"left-right": TileModeLeftRight, "grid": TileModeGrid,
}

func (s *Service) taskTileWindows(mode string, names []string) resultFailure {
	tm, ok := tileModeMap[mode]
	if !ok {
		return core.E("window.taskTileWindows", "unknown tile mode: "+mode, nil)
	}
	if len(names) == 0 {
		names = s.manager.List()
	}
	originX, originY, screenWidth, screenHeight := s.primaryScreenArea()
	return s.manager.TileWindows(tm, names, screenWidth, screenHeight, originX, originY)
}

func (s *Service) taskStackWindows(names []string, offsetX, offsetY int) resultFailure {
	if len(names) == 0 {
		names = s.manager.List()
	}
	originX, originY, _, _ := s.primaryScreenArea()
	return s.manager.StackWindows(names, offsetX, offsetY, originX, originY)
}

var snapPosMap = map[string]SnapPosition{
	"left": SnapLeft, "right": SnapRight,
	"top": SnapTop, "bottom": SnapBottom,
	"top-left": SnapTopLeft, "top-right": SnapTopRight,
	"bottom-left": SnapBottomLeft, "bottom-right": SnapBottomRight,
	"center": SnapCenter, "centre": SnapCenter,
}

func (s *Service) taskSnapWindow(name, position string) resultFailure {
	pos, ok := snapPosMap[position]
	if !ok {
		return core.E("window.taskSnapWindow", "unknown snap position: "+position, nil)
	}
	originX, originY, screenWidth, screenHeight := s.primaryScreenArea()
	return s.manager.SnapWindow(name, pos, screenWidth, screenHeight, originX, originY)
}

var workflowLayoutMap = map[string]WorkflowLayout{
	"coding":       WorkflowCoding,
	"debugging":    WorkflowDebugging,
	"presenting":   WorkflowPresenting,
	"side-by-side": WorkflowSideBySide,
}

func (s *Service) taskApplyWorkflow(workflow string, names []string) resultFailure {
	layout, ok := workflowLayoutMap[workflow]
	if !ok {
		return core.E("window.taskApplyWorkflow", "unknown workflow layout: "+workflow, nil)
	}
	if len(names) == 0 {
		names = s.manager.List()
	}
	originX, originY, screenWidth, screenHeight := s.primaryScreenArea()
	return s.manager.ApplyWorkflow(layout, names, screenWidth, screenHeight, originX, originY)
}

// --- Zoom ---

func (s *Service) queryWindowZoom(name string) core.Result {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.Result{Value: core.E("window.queryWindowZoom", "window not found: "+name, nil), OK: false}
	}
	return core.Result{Value: pw.GetZoom(), OK: true}
}

func (s *Service) taskSetZoom(name string, magnification float64) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetZoom", "window not found: "+name, nil)
	}
	pw.SetZoom(magnification)
	return nil
}

func (s *Service) taskZoomIn(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskZoomIn", "window not found: "+name, nil)
	}
	current := pw.GetZoom()
	pw.SetZoom(current + 0.1)
	return nil
}

func (s *Service) taskZoomOut(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskZoomOut", "window not found: "+name, nil)
	}
	current := pw.GetZoom()
	next := current - 0.1
	if next < 0.1 {
		next = 0.1
	}
	pw.SetZoom(next)
	return nil
}

func (s *Service) taskZoomReset(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskZoomReset", "window not found: "+name, nil)
	}
	pw.SetZoom(1.0)
	return nil
}

// --- Content ---

func (s *Service) taskSetURL(name, url string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetURL", "window not found: "+name, nil)
	}
	if core.HasPrefix(url, "core://") {
		resolved, ok, err := s.resolveCoreScheme(url)
		if err != nil {
			return err
		}
		if !ok {
			return core.E("window.taskSetURL", "core scheme handler unavailable for "+url, nil)
		}
		pw.SetHTML(resolved.Body)
		preload := s.buildPreload(url)
		if preload != "" {
			pw.ExecJS(preload)
		}
		return nil
	}
	pw.SetURL(url)
	preload := s.buildPreload(url)
	if preload != "" {
		pw.ExecJS(preload)
	}
	return nil
}

func (s *Service) taskSetHTML(name, html string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetHTML", "window not found: "+name, nil)
	}
	pw.SetHTML(html)
	return nil
}

func (s *Service) taskExecJS(name, js string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskExecJS", "window not found: "+name, nil)
	}
	pw.ExecJS(js)
	return nil
}

// EvalJSEventName is the Wails event name the JS-side wrap emits a
// reply on. Exported so the WailsPlatform binder can register the
// app.Event.On listener with the same literal the wrap uses.
const EvalJSEventName = "lthn:eval-reply"

// taskEvalJS fires the JS body in the named window, then waits for
// a "lthn:eval-reply" event to land via CompleteEval. The JS body
// is wrapped in an IIFE that imports @wailsio/runtime and emits
// the reply keyed by reqId; CompleteEval routes that reply back to
// the channel this method blocks on.
//
// Returns EvalJSResult on platform success (whether the JS itself
// succeeded or raised — Err carries the JS-side exception in the
// latter case). Returns an error only on platform-level failure
// (window missing, timeout).
func (s *Service) taskEvalJS(name, js string, timeout time.Duration) (EvalJSResult, error) {
	pw, ok := s.manager.Get(name)
	if !ok {
		return EvalJSResult{}, core.E("window.taskEvalJS", "window not found: "+name, nil)
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	reqID := s.nextEvalID()
	ch := make(chan EvalJSResult, 1)
	s.pendingEvalsMu.Lock()
	if s.pendingEvals == nil {
		s.pendingEvals = make(map[string]chan EvalJSResult)
	}
	s.pendingEvals[reqID] = ch
	s.pendingEvalsMu.Unlock()
	defer func() {
		s.pendingEvalsMu.Lock()
		delete(s.pendingEvals, reqID)
		s.pendingEvalsMu.Unlock()
	}()

	pw.ExecJS(wrapEvalScript(reqID, js))

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		return EvalJSResult{}, core.E("window.taskEvalJS", "eval timeout: "+name, nil)
	}
}

// CompleteEval routes a JS-side reply from the Wails event listener
// back to the taskEvalJS caller. WailsPlatform's app.Event.On
// listener calls this with the {reqId, result, err} payload it
// receives on EvalJSEventName.
//
// Safe to call from any goroutine. Returns true when the reqId
// matched a pending caller; false when the reply arrived after
// timeout (or for an unknown id — the latter only happens if a
// third party emits the event, which would be an injection bug).
func (s *Service) CompleteEval(reqID string, result any, errStr string) bool {
	s.pendingEvalsMu.Lock()
	ch, ok := s.pendingEvals[reqID]
	s.pendingEvalsMu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- EvalJSResult{ReqID: reqID, Result: result, Err: errStr}:
		return true
	default:
		// Buffer full — caller already took a reply; ignore.
		return false
	}
}

// nextEvalID returns a monotonic per-service reqId. The prefix
// keeps two Service instances in the same process from colliding;
// it's lazily initialised on first use so tests that build a
// Service via NewService don't pay the syscall cost up front.
func (s *Service) nextEvalID() string {
	s.pendingEvalsMu.Lock()
	if s.evalPrefix == "" {
		r := core.RandomString(8)
		if r.OK {
			s.evalPrefix = r.Value.(string)
		} else {
			s.evalPrefix = "ev"
		}
	}
	s.nextEvalSeq++
	seq := s.nextEvalSeq
	prefix := s.evalPrefix
	s.pendingEvalsMu.Unlock()
	return prefix + "-" + core.Sprintf("%d", seq)
}

// wrapEvalScript renders the IIFE template that evaluates body,
// awaits a Promise return, catches exceptions, and emits the result
// on the Wails Events bus keyed by reqId. Body is injected as a
// JSON string literal so it can contain any source verbatim.
//
// The wrap is injected via pw.ExecJS, which runs as a CLASSIC script
// in the WebView main world — it has NO module resolver, so it cannot
// import("@wailsio/runtime") (a bare specifier) at runtime. Instead it
// reuses window.__lthnEmit, the already-resolved emitter the frontend
// bridge shim (frontend/index.html, a type="module" script) publishes
// after its own runtime import settles. One resolution, reused by
// console/error capture AND eval replies. If the shim hasn't published
// the global yet (eval fired before first paint), the reply is dropped
// and the caller times out — a console.error breadcrumb is left in
// DevTools rather than vanishing silently.
func wrapEvalScript(reqID, body string) string {
	const tpl = `(function(){
  var __id=%s;
  var __name=%q;
  var __post=function(payload){
    try {
      if (typeof window.__lthnEmit === "function") { window.__lthnEmit(__name, payload); }
      else { try { console.error("lthn-eval: window.__lthnEmit unavailable — bridge shim not loaded"); } catch(_){} }
    } catch(e){ try { console.error("lthn-eval emit:", e); } catch(_){} }
  };
  try {
    var __r = (0,eval)(%s);
    if (__r && typeof __r.then === "function") {
      __r.then(
        function(v){ __post({ reqId: __id, result: v }); },
        function(e){ __post({ reqId: __id, err: String(e) + (e && e.stack ? "\n"+e.stack : "") }); }
      );
    } else {
      __post({ reqId: __id, result: __r });
    }
  } catch(e) {
    __post({ reqId: __id, err: String(e) + (e && e.stack ? "\n"+e.stack : "") });
  }
})();`
	return core.Sprintf(tpl, jsString(reqID), EvalJSEventName, jsString(body))
}

// jsString produces a JSON-encoded JS string literal — escapes
// quotes, newlines, and unicode so the body can be safely embedded
// in the IIFE template above. Equivalent to JSON.stringify(str)
// on the JS side.
func jsString(s string) string {
	return core.JSONMarshalString(s)
}

// --- State toggles ---

func (s *Service) taskToggleFullscreen(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskToggleFullscreen", "window not found: "+name, nil)
	}
	pw.ToggleFullscreen()
	return nil
}

func (s *Service) taskToggleMaximise(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskToggleMaximise", "window not found: "+name, nil)
	}
	pw.ToggleMaximise()
	return nil
}

// --- Bounds ---

func (s *Service) queryWindowBounds(name string) core.Result {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.Result{Value: core.E("window.queryWindowBounds", "window not found: "+name, nil), OK: false}
	}
	x, y, width, height := pw.GetBounds()
	return core.Result{Value: WindowBounds{X: x, Y: y, Width: width, Height: height}, OK: true}
}

func (s *Service) taskSetBounds(name string, x, y, width, height int) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetBounds", "window not found: "+name, nil)
	}
	pw.SetBounds(x, y, width, height)
	s.manager.State().UpdatePosition(name, x, y)
	s.manager.State().UpdateSize(name, width, height)
	return nil
}

// --- Content protection ---

func (s *Service) taskSetContentProtection(name string, protection bool) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskSetContentProtection", "window not found: "+name, nil)
	}
	pw.SetContentProtection(protection)
	return nil
}

// --- Flash ---

func (s *Service) taskFlash(name string, enabled bool) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskFlash", "window not found: "+name, nil)
	}
	pw.Flash(enabled)
	return nil
}

// --- Print ---

func (s *Service) taskPrint(name string) resultFailure {
	pw, ok := s.manager.Get(name)
	if !ok {
		return core.E("window.taskPrint", "window not found: "+name, nil)
	}
	return pw.Print()
}

// Manager returns the underlying window Manager for direct access.
func (s *Service) Manager() *Manager {
	return s.manager
}
