package window

// MockPlatform is an exported mock for cross-package integration tests.
// For internal tests, use the unexported mockPlatform in mock_test.go.
type MockPlatform struct {
	Windows []*MockWindow
}

func NewMockPlatform() *MockPlatform {
	return &MockPlatform{}
}

func (m *MockPlatform) CreateWindow(options PlatformWindowOptions) PlatformWindow {
	w := &MockWindow{
		name: options.Name, title: options.Title, url: options.URL, html: options.HTML,
		width: options.Width, height: options.Height,
		x: options.X, y: options.Y,
		opacity:     1.0,
		execJSCalls: nil,
	}
	if options.JS != "" {
		w.execJSCalls = append(w.execJSCalls, options.JS)
	}
	m.Windows = append(m.Windows, w)
	return w
}

func (m *MockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.Windows))
	for i, w := range m.Windows {
		out[i] = w
	}
	return out
}

type MockWindow struct {
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

func (w *MockWindow) Name() string                    { return w.name }
func (w *MockWindow) Title() string                   { return w.title }
func (w *MockWindow) Position() (int, int)            { return w.x, w.y }
func (w *MockWindow) Size() (int, int)                { return w.width, w.height }
func (w *MockWindow) IsMaximised() bool               { return w.maximised }
func (w *MockWindow) IsFocused() bool                 { return w.focused }
func (w *MockWindow) IsVisible() bool                 { return w.visible }
func (w *MockWindow) IsFullscreen() bool              { return w.fullscreened }
func (w *MockWindow) IsMinimised() bool               { return w.minimised }
func (w *MockWindow) IsAlwaysOnTop() bool             { return w.alwaysOnTop }
func (w *MockWindow) GetBounds() (int, int, int, int) { return w.x, w.y, w.width, w.height }
func (w *MockWindow) GetZoom() float64 {
	if w.zoom == 0 {
		return 1.0
	}
	return w.zoom
}
func (w *MockWindow) GetOpacity() float64                  { return w.opacity }
func (w *MockWindow) SetTitle(title string)                { w.title = title }
func (w *MockWindow) SetPosition(x, y int)                 { w.x = x; w.y = y }
func (w *MockWindow) SetSize(width, height int)            { w.width = width; w.height = height }
func (w *MockWindow) SetBackgroundColour(r, g, b, a uint8) { w.backgroundColour = [4]uint8{r, g, b, a} }
func (w *MockWindow) SetVisibility(visible bool)           { w.visible = visible }
func (w *MockWindow) SetAlwaysOnTop(alwaysOnTop bool)      { w.alwaysOnTop = alwaysOnTop }
func (w *MockWindow) SetOpacity(opacity float64)           { w.opacity = opacity }
func (w *MockWindow) SetBounds(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}
func (w *MockWindow) SetURL(url string)                    { w.url = url }
func (w *MockWindow) SetHTML(html string)                  { w.html = html }
func (w *MockWindow) SetZoom(magnification float64)        { w.zoom = magnification }
func (w *MockWindow) SetContentProtection(protection bool) { w.contentProtection = protection }
func (w *MockWindow) Maximise()                            { w.maximised = true }
func (w *MockWindow) Restore()                             { w.maximised = false }
func (w *MockWindow) Minimise()                            { w.minimised = true }
func (w *MockWindow) Focus()                               { w.focused = true }
func (w *MockWindow) Close()                               { w.closed = true }
func (w *MockWindow) Show()                                { w.visible = true }
func (w *MockWindow) Hide()                                { w.visible = false }
func (w *MockWindow) Fullscreen()                          { w.fullscreened = true }
func (w *MockWindow) UnFullscreen()                        { w.fullscreened = false }
func (w *MockWindow) ToggleFullscreen()                    { w.fullscreened = !w.fullscreened }
func (w *MockWindow) ToggleMaximise()                      { w.maximised = !w.maximised }
func (w *MockWindow) ExecJS(js string)                     { w.execJSCalls = append(w.execJSCalls, js) }
func (w *MockWindow) Flash(enabled bool)                   { w.flashed = enabled }
func (w *MockWindow) Print() resultFailure                 { return nil }
func (w *MockWindow) OpenDevTools()                        { w.devToolsOpen = true }
func (w *MockWindow) CloseDevTools()                       { w.devToolsOpen = false }
func (w *MockWindow) OnWindowEvent(handler func(WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}
func (w *MockWindow) OnFileDrop(handler func(paths []string, target *DropTarget)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}

// SetCloseBehavior records the requested behaviour. Tests verify
// via CloseBehavior().
func (w *MockWindow) SetCloseBehavior(behavior CloseBehavior) {
	w.closeBehavior = behavior
}

// CloseBehavior returns the most-recently-set close behaviour.
// Empty until SetCloseBehavior is called.
func (w *MockWindow) CloseBehavior() CloseBehavior {
	return w.closeBehavior
}

func (w *MockWindow) ExecJSCalls() []string {
	return append([]string(nil), w.execJSCalls...)
}

func (w *MockWindow) HTMLContent() string {
	return w.html
}

func (w *MockWindow) DevToolsOpen() bool {
	return w.devToolsOpen
}
