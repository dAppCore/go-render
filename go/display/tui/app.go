// SPDX-Licence-Identifier: EUPL-1.2

// app.go: the tui manager. App runs a .ctml UI as a Bubble Tea program, so a
// consumer supplies only the markup and gets a live, cross-platform terminal
// screen -- the "active thing" over go-render's passive engine/display packages.
// App owns the runtime shell (window size, altscreen + mouse, quit); the app
// supplies the .ctml (and, once the hooks prop lands, interactive panels +
// domain messages). It is Path A of the manager design: a root Bubble Tea
// Model, reusing the widget palette as-is.
//
// Example -- render a .ctml document as a terminal screen in one call:
//
//	node, _ := ctml.Parse([]byte(`<h1>Welcome</h1>`))
//	tui.Run(tui.NewApp(node))
package tui

import (
	html "dappco.re/go/render/engine/html"
)

// App is the root manager Model: each frame it renders a .ctml node tree to the
// terminal at the current width, tracks the window size, and quits on the quit
// keys (ctrl+c / q by default). It satisfies Model, so tui.Run -- or any Bubble
// Tea Program -- drives it.
type App struct {
	node   html.Node
	ctx    *html.Context
	opts   html.TermOptions
	width  int
	height int
	quit   func(KeyPressMsg) bool
}

// AppOption configures an App at construction.
type AppOption func(*App)

// WithTheme sets the terminal theme the .ctml renders through. Defaults to the
// house theme (html.DefaultTermTheme, applied by the renderer when Theme is nil).
func WithTheme(t *html.TermTheme) AppOption { return func(a *App) { a.opts.Theme = t } }

// WithRenderContext sets the render Context (bindings, locale). Defaults to
// html.NewContext(). Named to avoid colliding with WithContext, the re-exported
// Bubble Tea program option.
func WithRenderContext(ctx *html.Context) AppOption { return func(a *App) { a.ctx = ctx } }

// WithQuitKeys replaces the quit predicate. The default quits on ctrl+c or q;
// pass a predicate to change or disable that (e.g. ctrl+c only).
func WithQuitKeys(pred func(KeyPressMsg) bool) AppOption { return func(a *App) { a.quit = pred } }

// NewApp builds a manager around a parsed .ctml node tree (see
// dappco.re/go/render/engine/ctml.Parse). Options tune the theme, render context,
// and quit keys; the defaults render the house theme and quit on ctrl+c / q.
func NewApp(node html.Node, opts ...AppOption) *App {
	a := &App{
		node: node,
		ctx:  html.NewContext(),
		quit: defaultQuit,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

// defaultQuit quits on ctrl+c or a bare q -- the two idioms every terminal user
// reaches for. It reads the key's String form ("ctrl+c", "q") rather than the
// Code+Mod pair so it stays legible.
func defaultQuit(m KeyPressMsg) bool {
	switch m.String() {
	case "ctrl+c", "q":
		return true
	}
	// ctrl+c can also arrive as Code 'c' with the ctrl modifier, depending on the
	// terminal's keyboard mode -- catch that shape too.
	return m.Code == 'c' && m.Mod.Contains(ModCtrl)
}

// Init requests the initial terminal size so the first View renders at the real
// width rather than zero.
func (a *App) Init() Cmd { return RequestWindowSize }

// Update tracks the window size and quits on the quit keys; every other message
// (app-domain messages, later) leaves the model unchanged for now.
func (a *App) Update(msg Msg) (Model, Cmd) {
	switch m := msg.(type) {
	case WindowSizeMsg:
		a.width, a.height = m.Width, m.Height
	case KeyPressMsg:
		if a.quit(m) {
			return a, Quit
		}
	}
	return a, nil
}

// View renders the .ctml at the current width, in the alternate screen with
// mouse cell-motion enabled -- altscreen and mouse are View fields in Bubble Tea
// v2, not program options, so the runtime shell sets them here.
func (a *App) View() View {
	opts := a.opts
	opts.Width = a.width
	frame := html.RenderTerm(a.node, a.ctx, opts)

	v := NewView(frame)
	v.AltScreen = true
	v.MouseMode = MouseModeCellMotion
	return v
}
