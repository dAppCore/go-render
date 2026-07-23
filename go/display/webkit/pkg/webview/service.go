// pkg/webview/service.go
package webview

import (
	"context"
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
	"dappco.re/go/render/display/webkit/pkg/window"
	gowebview "dappco.re/go/webview"
)

// connector abstracts go-webview for testing. The real implementation wraps
// *gowebview.Webview, converting go-webview types to our own types at the boundary.
type connector interface {
	Navigate(url string) resultFailure
	Click(selector string) resultFailure
	Type(selector, text string) resultFailure
	Hover(selector string) resultFailure
	Select(selector, value string) resultFailure
	Check(selector string, checked bool) resultFailure
	Evaluate(script string) (any, resultFailure)
	Screenshot() ([]byte, resultFailure)
	GetURL() (string, resultFailure)
	GetTitle() (string, resultFailure)
	GetHTML(selector string) (string, resultFailure)
	QuerySelector(selector string) (*ElementInfo, resultFailure)
	QuerySelectorAll(selector string) ([]*ElementInfo, resultFailure)
	GetConsole() []ConsoleMessage
	ClearConsole()
	SetViewport(width, height int) resultFailure
	UploadFile(selector string, paths []string) resultFailure
	GetZoom() (float64, resultFailure)
	SetZoom(zoom float64) resultFailure
	Print(toPDF bool) ([]byte, resultFailure)
	Close() resultFailure
}

type Options struct {
	DebugURL     string        // Chrome debug endpoint (default: "http://localhost:9222")
	Timeout      time.Duration // Operation timeout (default: 30s)
	ConsoleLimit int           // Max console messages per window (default: 1000)
}

type Service struct {
	*core.ServiceRuntime[Options]
	options      Options
	connections  map[string]connector
	mu           sync.RWMutex
	diagMu       sync.RWMutex
	exceptions   map[string][]ExceptionInfo
	newConn      func(debugURL, windowName string) (connector, resultFailure) // injectable for tests
	watcherSetup func(conn connector, windowName string)                      // called after connection creation
}

// Register binds the webview service to a Core instance.
// core.WithService(webview.Register())
// core.WithService(webview.Register(func(o *Options) { o.DebugURL = "http://localhost:9223" }))
func Register(optionFns ...func(*Options)) func(*core.Core) core.Result {
	o := Options{
		DebugURL:     "http://localhost:9222",
		Timeout:      30 * time.Second,
		ConsoleLimit: 1000,
	}
	for _, fn := range optionFns {
		fn(&o)
	}
	return func(c *core.Core) core.Result {
		svc := &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, o),
			options:        o,
			connections:    make(map[string]connector),
			exceptions:     make(map[string][]ExceptionInfo),
			newConn:        defaultNewConn(o),
		}
		svc.watcherSetup = svc.defaultWatcherSetup
		return core.Result{Value: svc, OK: true}
	}
}

// defaultNewConn creates real go-webview connections.
func defaultNewConn(options Options) func(string, string) (connector, resultFailure) {
	return func(debugURL, windowName string) (connector, resultFailure) {
		windowName = core.Trim(windowName)
		if windowName == "" {
			return nil, core.E("webview.connect", "window name is required", nil)
		}
		// Enumerate targets, match by exact title/URL to avoid attaching to the wrong page.
		targets, err := gowebview.ListTargets(debugURL)
		if err != nil {
			return nil, err
		}
		if exactWindowTargetWSURL(targets, windowName) == "" {
			return nil, core.E("webview.connect", "no page target matched window name", nil)
		}
		res := gowebview.New(
			gowebview.WithDebugURL(debugURL),
			gowebview.WithTimeout(options.Timeout),
			gowebview.WithConsoleLimit(options.ConsoleLimit),
		)
		if !res.OK {
			return nil, res.Err()
		}
		return &realConnector{wv: res.Value.(*gowebview.Webview), debugURL: debugURL}, nil
	}
}

func exactWindowTargetMatch(candidate, windowName string) bool {
	return core.Trim(candidate) == windowName
}

func exactWindowTargetWSURL(targets []gowebview.TargetInfo, windowName string) string {
	for _, t := range targets {
		if t.Type != "page" || t.WebSocketDebuggerURL == "" {
			continue
		}
		if exactWindowTargetMatch(t.Title, windowName) || exactWindowTargetMatch(t.URL, windowName) {
			return t.WebSocketDebuggerURL
		}
	}
	return ""
}

// defaultWatcherSetup wires up console/exception watchers on real connectors.
// It broadcasts ActionConsoleMessage and ActionException via the Core IPC bus.
func (s *Service) defaultWatcherSetup(conn connector, windowName string) {
	rc, ok := conn.(*realConnector)
	if !ok {
		return // test mocks don't need watchers
	}

	cw := gowebview.NewConsoleWatcher(rc.wv)
	cw.AddHandler(func(msg gowebview.ConsoleMessage) {
		coreutil.DispatchAction(s.Core(), "webview.console", ActionConsoleMessage{
			Window: windowName,
			Message: ConsoleMessage{
				Type:      msg.Type,
				Text:      msg.Text,
				Timestamp: msg.Timestamp,
				URL:       msg.URL,
				Line:      msg.Line,
				Column:    msg.Column,
			},
		})
	})

	ew := gowebview.NewExceptionWatcher(rc.wv)
	ew.AddHandler(func(exc gowebview.ExceptionInfo) {
		info := ExceptionInfo{
			Text:       exc.Text,
			URL:        exc.URL,
			Line:       exc.LineNumber,
			Column:     exc.ColumnNumber,
			StackTrace: exc.StackTrace,
			Timestamp:  exc.Timestamp,
		}
		s.recordException(windowName, info)
		coreutil.DispatchAction(s.Core(), "webview.exception", ActionException{
			Window:    windowName,
			Exception: info,
		})
	})
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.registerTaskActions()
	return core.Result{OK: true}
}

// OnShutdown closes all CDP connections.
func (s *Service) OnShutdown(_ context.Context) core.Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, conn := range s.connections {
		conn.Close()
		delete(s.connections, name)
	}
	return core.Result{OK: true}
}

// HandleIPCEvents listens for window close events to clean up connections.
func (s *Service) HandleIPCEvents(_ *core.Core, msg core.Message) core.Result {
	switch m := msg.(type) {
	case window.ActionWindowClosed:
		s.mu.Lock()
		if conn, ok := s.connections[m.Name]; ok {
			conn.Close()
			delete(s.connections, m.Name)
		}
		s.mu.Unlock()
		s.diagMu.Lock()
		delete(s.exceptions, m.Name)
		s.diagMu.Unlock()
	}
	return core.Result{OK: true}
}

// getConn returns the connector for a window, creating it if needed.
func (s *Service) getConn(windowName string) (connector, resultFailure) {
	windowName = core.Trim(windowName)
	if windowName == "" {
		return nil, core.E("webview.getConn", "window name is required", nil)
	}
	s.mu.RLock()
	if conn, ok := s.connections[windowName]; ok {
		s.mu.RUnlock()
		return conn, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	// Double-check after acquiring write lock
	if conn, ok := s.connections[windowName]; ok {
		return conn, nil
	}
	conn, err := s.newConn(s.options.DebugURL, windowName)
	if err != nil {
		return nil, err
	}
	s.connections[windowName] = conn
	if s.watcherSetup != nil {
		s.watcherSetup(conn, windowName)
	}
	return conn, nil
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q := q.(type) {
	case QueryURL:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		url, err := conn.GetURL()
		return core.Result{}.New(url, err)
	case QueryTitle:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		title, err := conn.GetTitle()
		return core.Result{}.New(title, err)
	case QueryConsole:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		msgs := conn.GetConsole()
		// Filter by level if specified
		if q.Level != "" {
			var filtered []ConsoleMessage
			for _, m := range msgs {
				if m.Type == q.Level {
					filtered = append(filtered, m)
				}
			}
			msgs = filtered
		}
		// Apply limit
		if q.Limit > 0 && len(msgs) > q.Limit {
			msgs = msgs[len(msgs)-q.Limit:]
		}
		return core.Result{Value: msgs, OK: true}
	case QuerySelector:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		el, err := conn.QuerySelector(q.Selector)
		return core.Result{}.New(el, err)
	case QuerySelectorAll:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		els, err := conn.QuerySelectorAll(q.Selector)
		return core.Result{}.New(els, err)
	case QueryDOMTree:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		selector := q.Selector
		if selector == "" {
			selector = "html"
		}
		html, err := conn.GetHTML(selector)
		return core.Result{}.New(html, err)
	case QueryExceptions:
		return core.Result{Value: s.exceptionLog(q.Window, q.Limit), OK: true}
	case QueryZoom:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		zoom, err := conn.GetZoom()
		return core.Result{}.New(zoom, err)
	default:
		return core.Result{}
	}
}

// registerTaskActions registers all webview task handlers as named Core actions.
func (s *Service) registerTaskActions() {
	c := s.Core()
	register := func(names []string, handler func(context.Context, core.Options) core.Result) {
		for _, name := range names {
			c.Action(name, handler)
		}
	}
	register([]string{"webview.evaluate", "gui.webview.eval"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskEvaluate)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		result, err := conn.Evaluate(t.Script)
		return core.Result{}.New(result, err)
	})
	register([]string{"webview.click", "gui.webview.click"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClick)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Click(t.Selector))
	})
	register([]string{"webview.type", "gui.webview.type"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskType)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Type(t.Selector, t.Text))
	})
	register([]string{"webview.navigate", "gui.webview.navigate"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskNavigate)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Navigate(t.URL))
	})
	register([]string{"webview.screenshot", "gui.webview.screenshot"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskScreenshot)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		png, err := conn.Screenshot()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: ScreenshotResult{
			Base64:   base64.StdEncoding.EncodeToString(png),
			MimeType: "image/png",
		}, OK: true}
	})
	register([]string{"webview.scroll", "gui.webview.scroll"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskScroll)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		_, err = conn.Evaluate("window.scrollTo(" + strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y) + ")")
		return core.Result{Value: nil, OK: true}.New(err)
	})
	register([]string{"webview.hover", "gui.webview.hover"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskHover)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Hover(t.Selector))
	})
	register([]string{"webview.select", "gui.webview.select"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSelect)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Select(t.Selector, t.Value))
	})
	register([]string{"webview.check", "gui.webview.check"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskCheck)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Check(t.Selector, t.Checked))
	})
	register([]string{"webview.upload_file", "gui.webview.uploadFile"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskUploadFile)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.UploadFile(t.Selector, t.Paths))
	})
	register([]string{"webview.set_viewport", "gui.webview.setViewport"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetViewport)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.SetViewport(t.Width, t.Height))
	})
	register([]string{"webview.clear_console", "gui.webview.clearConsole"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClearConsole)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conn.ClearConsole()
		return core.Result{OK: true}
	})
	register([]string{"webview.setURL", "gui.webview.setURL"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetURL)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Navigate(t.URL))
	})
	register([]string{"webview.set_zoom", "gui.webview.setZoom"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetZoom)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.SetZoom(t.Zoom))
	})
	register([]string{"webview.print", "gui.webview.print"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskPrint)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		pdfBytes, err := conn.Print(t.ToPDF)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if !t.ToPDF {
			return core.Result{OK: true}
		}
		return core.Result{Value: PrintResult{
			Base64:   base64.StdEncoding.EncodeToString(pdfBytes),
			MimeType: "application/pdf",
		}, OK: true}
	})
	register([]string{"webview.devtools_open", "gui.webview.devtoolsOpen"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskDevToolsOpen)
		return core.Result{Value: nil, OK: true}.New(s.devToolsOpen(t.Window))
	})
	register([]string{"webview.devtools_close", "gui.webview.devtoolsClose"}, func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskDevToolsClose)
		return core.Result{Value: nil, OK: true}.New(s.devToolsClose(t.Window))
	})
}

func (s *Service) recordException(windowName string, info ExceptionInfo) {
	s.diagMu.Lock()
	defer s.diagMu.Unlock()

	log := append(s.exceptions[windowName], info)
	if limit := s.options.ConsoleLimit; limit > 0 && len(log) > limit {
		log = log[len(log)-limit:]
	}
	s.exceptions[windowName] = log
}

func (s *Service) exceptionLog(windowName string, limit int) []ExceptionInfo {
	s.diagMu.RLock()
	defer s.diagMu.RUnlock()

	log := append([]ExceptionInfo(nil), s.exceptions[windowName]...)
	if limit > 0 && len(log) > limit {
		log = log[len(log)-limit:]
	}
	return log
}

func (s *Service) devToolsOpen(windowName string) resultFailure {
	return s.withWindowHandle(windowName, func(handle any) resultFailure {
		if opener, ok := handle.(interface{ OpenDevTools() }); ok {
			opener.OpenDevTools()
			return nil
		}
		return core.E("webview.devToolsOpen", "window does not support developer tools", nil)
	})
}

func (s *Service) devToolsClose(windowName string) resultFailure {
	return s.withWindowHandle(windowName, func(handle any) resultFailure {
		if closer, ok := handle.(interface{ CloseDevTools() }); ok {
			closer.CloseDevTools()
			return nil
		}
		return nil
	})
}

func (s *Service) withWindowHandle(windowName string, fn func(handle any) resultFailure) resultFailure {
	windowService, ok := core.ServiceFor[*window.Service](s.Core(), "window")
	if !ok {
		return core.E("webview.withWindowHandle", "window service unavailable", nil)
	}
	handle, ok := windowService.Manager().Get(windowName)
	if !ok {
		return core.E("webview.withWindowHandle", "window not found: "+windowName, nil)
	}
	return fn(handle)
}

// realConnector wraps *gowebview.Webview, converting types at the boundary.
// debugURL is retained so that PDF printing can issue a Page.printToPDF CDP call
// via a fresh CDPClient, since go-webview v0.1.7 does not expose a PrintToPDF helper.
type realConnector struct {
	wv       *gowebview.Webview
	debugURL string // Chrome debug HTTP endpoint (e.g., http://localhost:9222) for direct CDP calls
}

func (r *realConnector) Navigate(url string) resultFailure           { return r.wv.Navigate(url) }
func (r *realConnector) Click(sel string) resultFailure              { return r.wv.Click(sel) }
func (r *realConnector) Type(sel, text string) resultFailure         { return r.wv.Type(sel, text) }
func (r *realConnector) Evaluate(script string) (any, resultFailure) { return r.wv.Evaluate(script) }
func (r *realConnector) Screenshot() ([]byte, resultFailure)         { return r.wv.Screenshot() }
func (r *realConnector) GetURL() (string, resultFailure)             { return r.wv.GetURL() }
func (r *realConnector) GetTitle() (string, resultFailure)           { return r.wv.GetTitle() }
func (r *realConnector) GetHTML(sel string) (string, resultFailure)  { return r.wv.GetHTML(sel) }
func (r *realConnector) ClearConsole()                               { r.wv.ClearConsole() }
func (r *realConnector) Close() resultFailure                        { return r.wv.Close() }
func (r *realConnector) SetViewport(w, h int) resultFailure          { return r.wv.SetViewport(w, h) }
func (r *realConnector) UploadFile(sel string, p []string) resultFailure {
	return r.wv.UploadFile(sel, p)
}

// GetZoom returns the current CSS zoom level as a float64.
// zoom, _ := conn.GetZoom()  // 1.0 = 100%, 1.5 = 150%
func (r *realConnector) GetZoom() (float64, resultFailure) {
	raw, err := r.wv.Evaluate("parseFloat(document.documentElement.style.zoom) || 1.0")
	if err != nil {
		return 0, core.E("realConnector.GetZoom", "failed to get zoom", err)
	}
	switch v := raw.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 1.0, nil
	}
}

// SetZoom sets the CSS zoom level on the document root element.
// conn.SetZoom(1.5)  // 150%
// conn.SetZoom(1.0)  // reset to normal
func (r *realConnector) SetZoom(zoom float64) resultFailure {
	script := "document.documentElement.style.zoom = '" + strconv.FormatFloat(zoom, 'g', -1, 64) + "'; undefined"
	_, err := r.wv.Evaluate(script)
	if err != nil {
		return core.E("realConnector.SetZoom", "failed to set zoom", err)
	}
	return nil
}

// Print triggers window.print() or exports to PDF via Page.printToPDF.
// When toPDF is false the browser print dialog is opened (via window.print()) and nil bytes are returned.
// When toPDF is true a fresh CDPClient is opened against the stored WebSocket URL to issue
// Page.printToPDF, which returns raw PDF bytes.
func (r *realConnector) Print(toPDF bool) ([]byte, resultFailure) {
	if !toPDF {
		_, err := r.wv.Evaluate("window.print(); undefined")
		if err != nil {
			return nil, core.E("realConnector.Print", "failed to open print dialog", err)
		}
		return nil, nil
	}

	if r.debugURL == "" {
		return nil, core.E("realConnector.Print", "no debug URL stored; cannot issue Page.printToPDF", nil)
	}

	// Open a dedicated CDPClient for the single Page.printToPDF call.
	// NewCDPClient connects to the first page target at the debug endpoint.
	clientRes := gowebview.NewCDPClient(r.debugURL)
	if !clientRes.OK {
		return nil, core.E("realConnector.Print", "failed to connect for PDF export", clientRes.Err())
	}
	client := clientRes.Value.(*gowebview.CDPClient)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := client.Call(ctx, "Page.printToPDF", map[string]any{
		"printBackground": true,
	})
	if err != nil {
		return nil, core.E("realConnector.Print", "Page.printToPDF failed", err)
	}

	dataStr, ok := result["data"].(string)
	if !ok {
		return nil, core.E("realConnector.Print", "Page.printToPDF returned no data", nil)
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil, core.E("realConnector.Print", "failed to decode PDF data", err)
	}

	return pdfBytes, nil
}

func (r *realConnector) Hover(sel string) resultFailure {
	return gowebview.NewActionSequence().Add(&gowebview.HoverAction{Selector: sel}).Execute(context.Background(), r.wv)
}

func (r *realConnector) Select(sel, val string) resultFailure {
	return gowebview.NewActionSequence().Add(&gowebview.SelectAction{Selector: sel, Value: val}).Execute(context.Background(), r.wv)
}

func (r *realConnector) Check(sel string, checked bool) resultFailure {
	return gowebview.NewActionSequence().Add(&gowebview.CheckAction{Selector: sel, Checked: checked}).Execute(context.Background(), r.wv)
}

func (r *realConnector) QuerySelector(sel string) (*ElementInfo, resultFailure) {
	el, err := r.wv.QuerySelector(sel)
	if err != nil {
		return nil, err
	}
	return convertElementInfo(el), nil
}

func (r *realConnector) QuerySelectorAll(sel string) ([]*ElementInfo, resultFailure) {
	els, err := r.wv.QuerySelectorAll(sel)
	if err != nil {
		return nil, err
	}
	result := make([]*ElementInfo, len(els))
	for i, el := range els {
		result[i] = convertElementInfo(el)
	}
	return result, nil
}

func (r *realConnector) GetConsole() []ConsoleMessage {
	raw := r.wv.GetConsole()
	msgs := make([]ConsoleMessage, len(raw))
	for i, m := range raw {
		msgs[i] = ConsoleMessage{
			Type: m.Type, Text: m.Text, Timestamp: m.Timestamp,
			URL: m.URL, Line: m.Line, Column: m.Column,
		}
	}
	return msgs
}

func convertElementInfo(el *gowebview.ElementInfo) *ElementInfo {
	if el == nil {
		return nil
	}
	info := &ElementInfo{
		TagName:    el.TagName,
		Attributes: el.Attributes,
		InnerText:  el.InnerText,
		InnerHTML:  el.InnerHTML,
	}
	if el.BoundingBox != nil {
		info.BoundingBox = &BoundingBox{
			X: el.BoundingBox.X, Y: el.BoundingBox.Y,
			Width: el.BoundingBox.Width, Height: el.BoundingBox.Height,
		}
	}
	return info
}
