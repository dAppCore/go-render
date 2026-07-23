// pkg/window/options.go
package window

// WindowOption is the compatibility layer for option-chain callers.
// Prefer a Window literal with Manager.Create.
type WindowOption func(*Window) resultFailure

// Deprecated: use Manager.Create(&Window{Name: "settings", URL: "/", Width: 800, Height: 600}).
func ApplyOptions(options ...WindowOption) (*Window, resultFailure) {
	windowSpec := &Window{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(windowSpec); err != nil {
			return nil, err
		}
	}
	return windowSpec, nil
}

// Compatibility helpers for callers still using option chains.
// Use: window.WithName("main")
func WithName(name string) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.Name = name; return nil }
}

// WithTitle sets the window title.
// Use: window.WithTitle("Core Editor")
func WithTitle(title string) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.Title = title; return nil }
}

// WithURL sets the initial window URL.
// Use: window.WithURL("/editor")
func WithURL(url string) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.URL = url; return nil }
}

// WithHTML sets the initial HTML content.
// Use: window.WithHTML("<main>Ready</main>")
func WithHTML(html string) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.HTML = html; return nil }
}

// WithJS sets the initial preload JavaScript.
// Use: window.WithJS("window.__CORE_READY__ = true")
func WithJS(js string) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.JS = js; return nil }
}

// WithSize sets the initial window size.
// Use: window.WithSize(1280, 800)
func WithSize(width, height int) WindowOption {
	return func(windowSpec *Window) resultFailure {
		windowSpec.Width = width
		windowSpec.Height = height
		return nil
	}
}

// WithPosition sets the initial window position.
// Use: window.WithPosition(160, 120)
func WithPosition(x, y int) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.X = x; windowSpec.Y = y; return nil }
}

// WithMinSize sets the minimum window size.
// Use: window.WithMinSize(640, 480)
func WithMinSize(width, height int) WindowOption {
	return func(windowSpec *Window) resultFailure {
		windowSpec.MinWidth = width
		windowSpec.MinHeight = height
		return nil
	}
}

// WithMaxSize sets the maximum window size.
// Use: window.WithMaxSize(1920, 1080)
func WithMaxSize(width, height int) WindowOption {
	return func(windowSpec *Window) resultFailure {
		windowSpec.MaxWidth = width
		windowSpec.MaxHeight = height
		return nil
	}
}

// WithFrameless toggles the native window frame.
// Use: window.WithFrameless(true)
func WithFrameless(frameless bool) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.Frameless = frameless; return nil }
}

// WithHidden starts the window hidden.
// Use: window.WithHidden(true)
func WithHidden(hidden bool) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.Hidden = hidden; return nil }
}

// WithAlwaysOnTop keeps the window above other windows.
// Use: window.WithAlwaysOnTop(true)
func WithAlwaysOnTop(alwaysOnTop bool) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.AlwaysOnTop = alwaysOnTop; return nil }
}

// WithBackgroundColour sets the window background colour with alpha.
// Use: window.WithBackgroundColour(0, 0, 0, 0)
func WithBackgroundColour(r, g, b, a uint8) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.BackgroundColour = [4]uint8{r, g, b, a}; return nil }
}

// WithFileDrop enables drag-and-drop file handling.
// Use: window.WithFileDrop(true)
func WithFileDrop(enabled bool) WindowOption {
	return func(windowSpec *Window) resultFailure { windowSpec.EnableFileDrop = enabled; return nil }
}
