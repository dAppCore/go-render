package window

type mockPlatform struct {
	windows []*mockWindow
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{}
}

func (m *mockPlatform) CreateWindow(options PlatformWindowOptions) PlatformWindow {
	w := &mockWindow{
		name: options.Name, title: options.Title, url: options.URL, html: options.HTML,
		width: options.Width, height: options.Height,
		x: options.X, y: options.Y,
		opacity: 1.0,
	}
	if options.JS != "" {
		w.execJSCalls = append(w.execJSCalls, options.JS)
	}
	m.windows = append(m.windows, w)
	return w
}

func (m *mockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.windows))
	for i, w := range m.windows {
		out[i] = w
	}
	return out
}

type mockWindow struct {
	name, title, url, html string
	width, height, x, y    int
	maximised, focused     bool
	visible, alwaysOnTop   bool
	backgroundColour       [4]uint8
	closed                 bool
	minimised              bool
	fullscreened           bool
	zoom                   float64
	opacity                float64
	contentProtection      bool
	flashed                bool
	devToolsOpen           bool
	execJSCalls            []string
	eventHandlers          []func(WindowEvent)
	fileDropHandlers       []func(paths []string, target *DropTarget)
	closeBehavior          CloseBehavior
}

func (w *mockWindow) Name() string                    { return w.name }
func (w *mockWindow) Title() string                   { return w.title }
func (w *mockWindow) Position() (int, int)            { return w.x, w.y }
func (w *mockWindow) Size() (int, int)                { return w.width, w.height }
func (w *mockWindow) IsMaximised() bool               { return w.maximised }
func (w *mockWindow) IsFocused() bool                 { return w.focused }
func (w *mockWindow) IsVisible() bool                 { return w.visible }
func (w *mockWindow) IsFullscreen() bool              { return w.fullscreened }
func (w *mockWindow) IsMinimised() bool               { return w.minimised }
func (w *mockWindow) IsAlwaysOnTop() bool             { return w.alwaysOnTop }
func (w *mockWindow) GetBounds() (int, int, int, int) { return w.x, w.y, w.width, w.height }
func (w *mockWindow) GetZoom() float64 {
	if w.zoom == 0 {
		return 1.0
	}
	return w.zoom
}
func (w *mockWindow) GetOpacity() float64                  { return w.opacity }
func (w *mockWindow) SetTitle(title string)                { w.title = title }
func (w *mockWindow) SetPosition(x, y int)                 { w.x = x; w.y = y }
func (w *mockWindow) SetSize(width, height int)            { w.width = width; w.height = height }
func (w *mockWindow) SetBackgroundColour(r, g, b, a uint8) { w.backgroundColour = [4]uint8{r, g, b, a} }
func (w *mockWindow) SetVisibility(visible bool)           { w.visible = visible }
func (w *mockWindow) SetAlwaysOnTop(alwaysOnTop bool)      { w.alwaysOnTop = alwaysOnTop }
func (w *mockWindow) SetOpacity(opacity float64)           { w.opacity = opacity }
func (w *mockWindow) SetBounds(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}
func (w *mockWindow) SetURL(url string)                    { w.url = url }
func (w *mockWindow) SetHTML(html string)                  { w.html = html }
func (w *mockWindow) SetZoom(magnification float64)        { w.zoom = magnification }
func (w *mockWindow) SetContentProtection(protection bool) { w.contentProtection = protection }
func (w *mockWindow) Maximise()                            { w.maximised = true }
func (w *mockWindow) Restore()                             { w.maximised = false }
func (w *mockWindow) Minimise()                            { w.minimised = true }
func (w *mockWindow) Focus()                               { w.focused = true }
func (w *mockWindow) Close()                               { w.closed = true }
func (w *mockWindow) Show()                                { w.visible = true }
func (w *mockWindow) Hide()                                { w.visible = false }
func (w *mockWindow) Fullscreen()                          { w.fullscreened = true }
func (w *mockWindow) UnFullscreen()                        { w.fullscreened = false }
func (w *mockWindow) ToggleFullscreen()                    { w.fullscreened = !w.fullscreened }
func (w *mockWindow) ToggleMaximise()                      { w.maximised = !w.maximised }
func (w *mockWindow) ExecJS(js string)                     { w.execJSCalls = append(w.execJSCalls, js) }
func (w *mockWindow) Flash(enabled bool)                   { w.flashed = enabled }
func (w *mockWindow) Print() resultFailure                 { return nil }
func (w *mockWindow) OpenDevTools()                        { w.devToolsOpen = true }
func (w *mockWindow) CloseDevTools()                       { w.devToolsOpen = false }
func (w *mockWindow) OnWindowEvent(handler func(WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}
func (w *mockWindow) OnFileDrop(handler func(paths []string, target *DropTarget)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}

func (w *mockWindow) SetCloseBehavior(behavior CloseBehavior) {
	w.closeBehavior = behavior
}

// emit fires a test event to all registered handlers.
func (w *mockWindow) emit(e WindowEvent) {
	for _, h := range w.eventHandlers {
		h(e)
	}
}

// emitFileDrop simulates a file drop on the window. Pass nil target
// for legacy zero-context drops, or a DropTarget with the element
// metadata the consumer expects to receive.
func (w *mockWindow) emitFileDrop(paths []string, target *DropTarget) {
	for _, h := range w.fileDropHandlers {
		h(paths, target)
	}
}

// recordingBinder is a test Platform that composes mockPlatform's
// window-creation surface AND implements CustomEventBinder by
// recording every (name, cb) pair the Service registers. Tests
// inspect bindings to prove the wiring without depending on
// app.Event.On (Wails-only).
type recordingBinder struct {
	mockPlatform
	bindings []recordedBinding
}

type recordedBinding struct {
	name string
	cb   func(any)
}

// BindCustomEvent satisfies CustomEventBinder.
func (r *recordingBinder) BindCustomEvent(name string, cb func(data any)) {
	r.bindings = append(r.bindings, recordedBinding{name: name, cb: cb})
}
