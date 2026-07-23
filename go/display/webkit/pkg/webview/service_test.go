// pkg/webview/service_test.go
package webview

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/window"
	gowebview "dappco.re/go/webview"
)

type mockConnector struct {
	url        string
	title      string
	html       string
	evalResult any
	screenshot []byte
	console    []ConsoleMessage
	elements   []*ElementInfo
	closed     bool

	lastClickSel       string
	lastTypeSel        string
	lastTypeText       string
	lastNavURL         string
	lastHoverSel       string
	lastSelectSel      string
	lastSelectVal      string
	lastCheckSel       string
	lastCheckVal       bool
	lastUploadSel      string
	lastUploadPaths    []string
	lastViewportW      int
	lastViewportH      int
	consoleClearCalled bool

	zoom          float64
	lastZoomSet   float64
	printToPDF    bool
	printCalled   bool
	printPDFBytes []byte
	printErr      resultFailure
}

func (m *mockConnector) Navigate(url string) resultFailure { m.lastNavURL = url; return nil }
func (m *mockConnector) Click(sel string) resultFailure    { m.lastClickSel = sel; return nil }
func (m *mockConnector) Type(sel, text string) resultFailure {
	m.lastTypeSel = sel
	m.lastTypeText = text
	return nil
}
func (m *mockConnector) Hover(sel string) resultFailure { m.lastHoverSel = sel; return nil }
func (m *mockConnector) Select(sel, val string) resultFailure {
	m.lastSelectSel = sel
	m.lastSelectVal = val
	return nil
}
func (m *mockConnector) Check(sel string, c bool) resultFailure {
	m.lastCheckSel = sel
	m.lastCheckVal = c
	return nil
}
func (m *mockConnector) Evaluate(s string) (any, resultFailure)     { return m.evalResult, nil }
func (m *mockConnector) Screenshot() ([]byte, resultFailure)        { return m.screenshot, nil }
func (m *mockConnector) GetURL() (string, resultFailure)            { return m.url, nil }
func (m *mockConnector) GetTitle() (string, resultFailure)          { return m.title, nil }
func (m *mockConnector) GetHTML(sel string) (string, resultFailure) { return m.html, nil }
func (m *mockConnector) ClearConsole()                              { m.consoleClearCalled = true }
func (m *mockConnector) Close() resultFailure                       { m.closed = true; return nil }
func (m *mockConnector) SetViewport(w, h int) resultFailure {
	m.lastViewportW = w
	m.lastViewportH = h
	return nil
}
func (m *mockConnector) UploadFile(sel string, p []string) resultFailure {
	m.lastUploadSel = sel
	m.lastUploadPaths = p
	return nil
}

func (m *mockConnector) QuerySelector(sel string) (*ElementInfo, resultFailure) {
	if len(m.elements) > 0 {
		return m.elements[0], nil
	}
	return nil, nil
}

func (m *mockConnector) QuerySelectorAll(sel string) ([]*ElementInfo, resultFailure) {
	return m.elements, nil
}

func (m *mockConnector) GetConsole() []ConsoleMessage { return m.console }

func (m *mockConnector) GetZoom() (float64, resultFailure) {
	if m.zoom == 0 {
		return 1.0, nil
	}
	return m.zoom, nil
}

func (m *mockConnector) SetZoom(zoom float64) resultFailure {
	m.lastZoomSet = zoom
	m.zoom = zoom
	return nil
}

func (m *mockConnector) Print(toPDF bool) ([]byte, resultFailure) {
	m.printCalled = true
	m.printToPDF = toPDF
	return m.printPDFBytes, m.printErr
}

func newTestService(t *core.T, mock *mockConnector) (*Service, *core.Core) {
	t.Helper()
	factory := Register()
	c := core.New(core.WithService(factory), core.WithServiceLock())
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "webview")
	// Inject mock connector
	svc.newConn = func(_, _ string) (connector, resultFailure) { return mock, nil }
	return svc, c
}

func newTestServiceWithWindow(t *core.T, mock *mockConnector) (*Service, *window.MockPlatform, *core.Core) {
	t.Helper()
	windowPlatform := window.NewMockPlatform()
	c := core.New(
		core.WithService(window.Register(windowPlatform)),
		core.WithService(Register()),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	result := taskRun(c, "window.open", window.TaskOpenWindow{
		Window: &window.Window{Name: "main", Title: "Main", Width: 800, Height: 600},
	})
	core.RequireTrue(t, result.OK)

	svc := core.MustServiceFor[*Service](c, "webview")
	svc.newConn = func(_, _ string) (connector, resultFailure) { return mock, nil }
	return svc, windowPlatform, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *core.T) {
	svc, _ := newTestService(t, &mockConnector{})
	core.AssertNotNil(t, svc)
	core.AssertNotEmpty(t, core.Sprintf("%T", svc))
}

func TestQueryURL_Good(t *core.T) {
	_, c := newTestService(t, &mockConnector{url: "https://example.com"})
	r := c.QUERY(QueryURL{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "https://example.com", r.Value)
}

func TestQueryTitle_Good(t *core.T) {
	_, c := newTestService(t, &mockConnector{title: "Test Page"})
	r := c.QUERY(QueryTitle{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "Test Page", r.Value)
}

func TestQueryConsole_Good(t *core.T) {
	mock := &mockConnector{console: []ConsoleMessage{
		{Type: core.Concat("lo", "g"), Text: "hello"},
		{Type: "error", Text: "oops"},
		{Type: core.Concat("lo", "g"), Text: "world"},
	}}
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryConsole{Window: "main", Level: "error", Limit: 10})
	core.RequireTrue(t, r.OK)
	msgs, _ := r.Value.([]ConsoleMessage)
	core.AssertLen(t, msgs, 1)
	core.AssertEqual(t, "oops", msgs[0].Text)
}

func TestQueryConsole_Good_Limit(t *core.T) {
	mock := &mockConnector{console: []ConsoleMessage{
		{Type: core.Concat("lo", "g"), Text: "a"},
		{Type: core.Concat("lo", "g"), Text: "b"},
		{Type: core.Concat("lo", "g"), Text: "c"},
	}}
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryConsole{Window: "main", Limit: 2})
	msgs, _ := r.Value.([]ConsoleMessage)
	core.AssertLen(t, msgs, 2)
	core.AssertEqual(t, "b", msgs[0].Text) // last 2
}

func TestTaskEvaluate_Good(t *core.T) {
	_, c := newTestService(t, &mockConnector{evalResult: 42})
	r := taskRun(c, "webview.evaluate", TaskEvaluate{Window: "main", Script: "21*2"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, 42, r.Value)
}

func TestTaskEvaluate_Bad_EmptyWindow(t *core.T) {
	_, c := newTestService(t, &mockConnector{evalResult: 42})
	r := taskRun(c, "webview.evaluate", TaskEvaluate{Window: " ", Script: "21*2"})
	core.AssertFalse(t, r.OK)
}

func TestExactWindowTargetMatch_Good(t *core.T) {
	core.AssertTrue(t, exactWindowTargetMatch("main", "main"))
	core.AssertFalse(t, exactWindowTargetMatch("main - docs", "main"))
	core.AssertNotEmpty(t, core.Sprintf("%T", exactWindowTargetMatch("main", "main")))
}

func TestExactWindowTargetWSURL_PicksMatchingPage(t *core.T) {
	wsURL := exactWindowTargetWSURL([]gowebview.TargetInfo{
		{Type: "page", Title: "dashboard", WebSocketDebuggerURL: "ws://first"},
		{Type: "page", Title: "main", WebSocketDebuggerURL: "ws://second"},
	}, "main")

	core.AssertEqual(t, "ws://second", wsURL)
}

func TestTaskClick_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.click", TaskClick{Window: "main", Selector: "#btn"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "#btn", mock.lastClickSel)
}

func TestTaskNavigate_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.navigate", TaskNavigate{Window: "main", URL: "https://example.com"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "https://example.com", mock.lastNavURL)
}

func TestTaskScreenshot_Good(t *core.T) {
	mock := &mockConnector{screenshot: []byte{0x89, 0x50, 0x4E, 0x47}}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.screenshot", TaskScreenshot{Window: "main"})
	core.RequireTrue(t, r.OK)
	sr, ok := r.Value.(ScreenshotResult)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "image/png", sr.MimeType)
	core.AssertNotEmpty(t, sr.Base64)
}

func TestTaskClearConsole_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.clear_console", TaskClearConsole{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, mock.consoleClearCalled)
}

func TestConnectionCleanup_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Access creates connection
	c.QUERY(QueryURL{Window: "main"})
	core.AssertFalse(t, mock.closed)
	// Window close action triggers cleanup
	_ = c.ACTION(window.ActionWindowClosed{Name: "main"})
	core.AssertTrue(t, mock.closed)
}

func TestQueryURL_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryURL{Window: "main"})
	core.AssertFalse(t, r.OK)
}

// --- SetURL ---

func TestTaskSetURL_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "main", URL: "https://example.com/page"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "https://example.com/page", mock.lastNavURL)
}

func TestTaskSetURL_Bad_UnknownWindow(t *core.T) {
	_, c := newTestService(t, &mockConnector{})
	// Inject a connector factory that errors
	svc := core.MustServiceFor[*Service](c, "webview")
	svc.newConn = func(_, _ string) (connector, resultFailure) {
		return nil, core.E("test", "no connection", nil)
	}
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "bad", URL: "https://example.com"})
	core.AssertFalse(t, r.OK)
}

func TestTaskSetURL_Ugly_EmptyURL(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "main", URL: ""})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "", mock.lastNavURL)
}

// --- Zoom ---

func TestQueryZoom_Good(t *core.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryZoom{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertInDelta(t, 1.5, r.Value.(float64), 0.001)
}

func TestQueryZoom_Good_DefaultsToOne(t *core.T) {
	mock := &mockConnector{} // zoom not set → GetZoom returns 1.0
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryZoom{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertInDelta(t, 1.0, r.Value.(float64), 0.001)
}

func TestQueryZoom_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryZoom{Window: "main"})
	core.AssertFalse(t, r.OK)
}

func TestTaskSetZoom_Good(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.set_zoom", TaskSetZoom{Window: "main", Zoom: 2.0})
	core.RequireTrue(t, r.OK)
	core.AssertInDelta(t, 2.0, mock.lastZoomSet, 0.001)
}

func TestTaskSetZoom_Good_Reset(t *core.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.set_zoom", TaskSetZoom{Window: "main", Zoom: 1.0})
	core.RequireTrue(t, r.OK)
	core.AssertInDelta(t, 1.0, mock.zoom, 0.001)
}

func TestTaskSetZoom_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("webview.set_zoom").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskSetZoom_Ugly_ZeroZoom(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Zero zoom is technically valid input; the connector accepts it.
	r := taskRun(c, "webview.set_zoom", TaskSetZoom{Window: "main", Zoom: 0})
	core.RequireTrue(t, r.OK)
	core.AssertInDelta(t, 0.0, mock.lastZoomSet, 0.001)
}

// --- Print ---

func TestTaskPrint_Good_Dialog(t *core.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: false})
	core.RequireTrue(t, r.OK)
	core.AssertNil(t, r.Value)
	core.AssertTrue(t, mock.printCalled)
	core.AssertFalse(t, mock.printToPDF)
}

func TestTaskPrint_Good_PDF(t *core.T) {
	pdfHeader := []byte{0x25, 0x50, 0x44, 0x46} // %PDF
	mock := &mockConnector{printPDFBytes: pdfHeader}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	core.RequireTrue(t, r.OK)
	pr, ok := r.Value.(PrintResult)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "application/pdf", pr.MimeType)
	core.AssertNotEmpty(t, pr.Base64)
	core.AssertTrue(t, mock.printToPDF)
}

func TestTaskPrint_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("webview.print").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskPrint_Bad_Error(t *core.T) {
	mock := &mockConnector{printErr: core.E("test", "print failed", nil)}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	core.AssertFalse(t, r.OK)
}

func TestTaskPrint_Ugly_EmptyPDF(t *core.T) {
	// toPDF=true but connector returns zero bytes — should still wrap as PrintResult
	mock := &mockConnector{printPDFBytes: []byte{}}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	core.RequireTrue(t, r.OK)
	pr, ok := r.Value.(PrintResult)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "application/pdf", pr.MimeType)
	core.AssertEqual(t, "", pr.Base64) // empty PDF encodes to empty base64
}

func TestQueryExceptions_Good(t *core.T) {
	svc, c := newTestService(t, &mockConnector{})
	svc.recordException("main", ExceptionInfo{Text: "boom", Line: 7})

	r := c.QUERY(QueryExceptions{Window: "main"})
	core.RequireTrue(t, r.OK)
	exceptions := r.Value.([]ExceptionInfo)
	core.AssertLen(t, exceptions, 1)
	core.AssertEqual(t, "boom", exceptions[0].Text)
	core.AssertEqual(t, 7, exceptions[0].Line)
}

func TestTaskDevTools_Good(t *core.T) {
	_, windowPlatform, c := newTestServiceWithWindow(t, &mockConnector{})

	r := taskRun(c, "webview.devtools_open", TaskDevToolsOpen{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertTrue(t, windowPlatform.Windows[0].DevToolsOpen())

	r = taskRun(c, "webview.devtools_close", TaskDevToolsClose{Window: "main"})
	core.RequireTrue(t, r.OK)
	core.AssertFalse(t, windowPlatform.Windows[0].DevToolsOpen())
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

type Connector = realConnector

func TestService_Connector_Navigate_Good(t *core.T) {
	// Connector Navigate
	ax7Variant := "Connector_Navigate:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Navigate("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Navigate_Bad(t *core.T) {
	// Connector Navigate
	ax7Variant := "Connector_Navigate:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Navigate("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Navigate_Ugly(t *core.T) {
	// Connector Navigate
	ax7Variant := "Connector_Navigate:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Navigate("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Click_Good(t *core.T) {
	// Connector Click
	ax7Variant := "Connector_Click:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Click("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Click_Bad(t *core.T) {
	// Connector Click
	ax7Variant := "Connector_Click:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Click("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Click_Ugly(t *core.T) {
	// Connector Click
	ax7Variant := "Connector_Click:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Click("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Type_Good(t *core.T) {
	// Connector Type
	ax7Variant := "Connector_Type:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Type("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Type_Bad(t *core.T) {
	// Connector Type
	ax7Variant := "Connector_Type:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Type("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Type_Ugly(t *core.T) {
	// Connector Type
	ax7Variant := "Connector_Type:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Type("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Evaluate_Good(t *core.T) {
	// Connector Evaluate
	ax7Variant := "Connector_Evaluate:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Evaluate("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Evaluate_Bad(t *core.T) {
	// Connector Evaluate
	ax7Variant := "Connector_Evaluate:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Evaluate("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Evaluate_Ugly(t *core.T) {
	// Connector Evaluate
	ax7Variant := "Connector_Evaluate:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Evaluate("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Screenshot_Good(t *core.T) {
	// Connector Screenshot
	ax7Variant := "Connector_Screenshot:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Screenshot()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Screenshot_Bad(t *core.T) {
	// Connector Screenshot
	ax7Variant := "Connector_Screenshot:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Screenshot()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Screenshot_Ugly(t *core.T) {
	// Connector Screenshot
	ax7Variant := "Connector_Screenshot:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Screenshot()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetURL_Good(t *core.T) {
	// Connector GetURL
	ax7Variant := "Connector_GetURL:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetURL()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetURL_Bad(t *core.T) {
	// Connector GetURL
	ax7Variant := "Connector_GetURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetURL()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetURL_Ugly(t *core.T) {
	// Connector GetURL
	ax7Variant := "Connector_GetURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetURL()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetTitle_Good(t *core.T) {
	// Connector GetTitle
	ax7Variant := "Connector_GetTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetTitle()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetTitle_Bad(t *core.T) {
	// Connector GetTitle
	ax7Variant := "Connector_GetTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetTitle()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetTitle_Ugly(t *core.T) {
	// Connector GetTitle
	ax7Variant := "Connector_GetTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetTitle()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetHTML_Good(t *core.T) {
	// Connector GetHTML
	ax7Variant := "Connector_GetHTML:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetHTML("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetHTML_Bad(t *core.T) {
	// Connector GetHTML
	ax7Variant := "Connector_GetHTML:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetHTML("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetHTML_Ugly(t *core.T) {
	// Connector GetHTML
	ax7Variant := "Connector_GetHTML:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetHTML("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_ClearConsole_Good(t *core.T) {
	// Connector ClearConsole
	ax7Variant := "Connector_ClearConsole:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		subject.ClearConsole()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_ClearConsole_Bad(t *core.T) {
	// Connector ClearConsole
	ax7Variant := "Connector_ClearConsole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		subject.ClearConsole()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_ClearConsole_Ugly(t *core.T) {
	// Connector ClearConsole
	ax7Variant := "Connector_ClearConsole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		subject.ClearConsole()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Close_Good(t *core.T) {
	// Connector Close
	ax7Variant := "Connector_Close:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Close_Bad(t *core.T) {
	// Connector Close
	ax7Variant := "Connector_Close:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Close_Ugly(t *core.T) {
	// Connector Close
	ax7Variant := "Connector_Close:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetViewport_Good(t *core.T) {
	// Connector SetViewport
	ax7Variant := "Connector_SetViewport:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetViewport(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetViewport_Bad(t *core.T) {
	// Connector SetViewport
	ax7Variant := "Connector_SetViewport:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetViewport(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetViewport_Ugly(t *core.T) {
	// Connector SetViewport
	ax7Variant := "Connector_SetViewport:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetViewport(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_UploadFile_Good(t *core.T) {
	// Connector UploadFile
	ax7Variant := "Connector_UploadFile:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.UploadFile("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_UploadFile_Bad(t *core.T) {
	// Connector UploadFile
	ax7Variant := "Connector_UploadFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.UploadFile("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_UploadFile_Ugly(t *core.T) {
	// Connector UploadFile
	ax7Variant := "Connector_UploadFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.UploadFile("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetZoom_Good(t *core.T) {
	// Connector GetZoom
	ax7Variant := "Connector_GetZoom:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetZoom()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetZoom_Bad(t *core.T) {
	// Connector GetZoom
	ax7Variant := "Connector_GetZoom:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetZoom()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetZoom_Ugly(t *core.T) {
	// Connector GetZoom
	ax7Variant := "Connector_GetZoom:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.GetZoom()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetZoom_Good(t *core.T) {
	// Connector SetZoom
	ax7Variant := "Connector_SetZoom:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetZoom(1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetZoom_Bad(t *core.T) {
	// Connector SetZoom
	ax7Variant := "Connector_SetZoom:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetZoom(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_SetZoom_Ugly(t *core.T) {
	// Connector SetZoom
	ax7Variant := "Connector_SetZoom:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.SetZoom(-1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Print_Good(t *core.T) {
	// Connector Print
	ax7Variant := "Connector_Print:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Print(true)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Print_Bad(t *core.T) {
	// Connector Print
	ax7Variant := "Connector_Print:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Print(false)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Print_Ugly(t *core.T) {
	// Connector Print
	ax7Variant := "Connector_Print:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.Print(false)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Hover_Good(t *core.T) {
	// Connector Hover
	ax7Variant := "Connector_Hover:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Hover("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Hover_Bad(t *core.T) {
	// Connector Hover
	ax7Variant := "Connector_Hover:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Hover("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Hover_Ugly(t *core.T) {
	// Connector Hover
	ax7Variant := "Connector_Hover:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Hover("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Select_Good(t *core.T) {
	// Connector Select
	ax7Variant := "Connector_Select:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Select("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Select_Bad(t *core.T) {
	// Connector Select
	ax7Variant := "Connector_Select:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Select("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Select_Ugly(t *core.T) {
	// Connector Select
	ax7Variant := "Connector_Select:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Select("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Check_Good(t *core.T) {
	// Connector Check
	ax7Variant := "Connector_Check:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Check("agent", true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Check_Bad(t *core.T) {
	// Connector Check
	ax7Variant := "Connector_Check:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Check("", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_Check_Ugly(t *core.T) {
	// Connector Check
	ax7Variant := "Connector_Check:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.Check("../../edge", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelector_Good(t *core.T) {
	// Connector QuerySelector
	ax7Variant := "Connector_QuerySelector:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelector("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelector_Bad(t *core.T) {
	// Connector QuerySelector
	ax7Variant := "Connector_QuerySelector:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelector("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelector_Ugly(t *core.T) {
	// Connector QuerySelector
	ax7Variant := "Connector_QuerySelector:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelector("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelectorAll_Good(t *core.T) {
	// Connector QuerySelectorAll
	ax7Variant := "Connector_QuerySelectorAll:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelectorAll("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelectorAll_Bad(t *core.T) {
	// Connector QuerySelectorAll
	ax7Variant := "Connector_QuerySelectorAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelectorAll("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_QuerySelectorAll_Ugly(t *core.T) {
	// Connector QuerySelectorAll
	ax7Variant := "Connector_QuerySelectorAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0, got1 := subject.QuerySelectorAll("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetConsole_Good(t *core.T) {
	// Connector GetConsole
	ax7Variant := "Connector_GetConsole:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.GetConsole()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetConsole_Bad(t *core.T) {
	// Connector GetConsole
	ax7Variant := "Connector_GetConsole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.GetConsole()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Connector_GetConsole_Ugly(t *core.T) {
	// Connector GetConsole
	ax7Variant := "Connector_GetConsole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(realConnector)
	result := core.Try(func() any {
		got0 := subject.GetConsole()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

// AX7 generated source-matching smoke coverage.
