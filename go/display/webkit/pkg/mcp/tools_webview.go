// pkg/mcp/tools_webview.go
package mcp

import (
	"context"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"math"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/webview"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- webview_eval ---

type WebviewEvalInput struct {
	Window string `json:"window"`
	Script string `json:"script"`
}

type WebviewEvalOutput struct {
	Result any    `json:"result"`
	Window string `json:"window"`
}

func (s *Subsystem) webviewEval(_ context.Context, _ *mcp.CallToolRequest, input WebviewEvalInput) (*mcp.CallToolResult, WebviewEvalOutput, resultFailure) {
	r := s.core.Action("webview.evaluate").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskEvaluate{Window: input.Window, Script: input.Script}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewEvalOutput{}, e
		}
		return nil, WebviewEvalOutput{}, nil
	}
	return nil, WebviewEvalOutput{Result: r.Value, Window: input.Window}, nil
}

// --- webview_click ---

type WebviewClickInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewClickOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewClick(_ context.Context, _ *mcp.CallToolRequest, input WebviewClickInput) (*mcp.CallToolResult, WebviewClickOutput, resultFailure) {
	r := s.core.Action("webview.click").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskClick{Window: input.Window, Selector: input.Selector}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewClickOutput{}, e
		}
		return nil, WebviewClickOutput{}, nil
	}
	return nil, WebviewClickOutput{Success: true}, nil
}

// --- webview_type ---

type WebviewTypeInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Text     string `json:"text"`
}

type WebviewTypeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewType(_ context.Context, _ *mcp.CallToolRequest, input WebviewTypeInput) (*mcp.CallToolResult, WebviewTypeOutput, resultFailure) {
	r := s.core.Action("webview.type").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskType{Window: input.Window, Selector: input.Selector, Text: input.Text}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewTypeOutput{}, e
		}
		return nil, WebviewTypeOutput{}, nil
	}
	return nil, WebviewTypeOutput{Success: true}, nil
}

// --- webview_navigate ---

type WebviewNavigateInput struct {
	Window string `json:"window"`
	URL    string `json:"url"`
}

type WebviewNavigateOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNavigate(_ context.Context, _ *mcp.CallToolRequest, input WebviewNavigateInput) (*mcp.CallToolResult, WebviewNavigateOutput, resultFailure) {
	r := s.core.Action("webview.navigate").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskNavigate{Window: input.Window, URL: input.URL}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewNavigateOutput{}, e
		}
		return nil, WebviewNavigateOutput{}, nil
	}
	return nil, WebviewNavigateOutput{Success: true}, nil
}

// --- webview_screenshot ---

type WebviewScreenshotInput struct {
	Window string `json:"window"`
}

type WebviewScreenshotOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewScreenshot(_ context.Context, _ *mcp.CallToolRequest, input WebviewScreenshotInput) (*mcp.CallToolResult, WebviewScreenshotOutput, resultFailure) {
	r := s.core.Action("webview.screenshot").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskScreenshot{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScreenshotOutput{}, e
		}
		return nil, WebviewScreenshotOutput{}, nil
	}
	sr, ok := r.Value.(webview.ScreenshotResult)
	if !ok {
		return nil, WebviewScreenshotOutput{}, core.E("mcp.webviewScreenshot", "unexpected result type", nil)
	}
	return nil, WebviewScreenshotOutput{Base64: sr.Base64, MimeType: sr.MimeType}, nil
}

// --- webview_scroll ---

type WebviewScrollInput struct {
	Window string `json:"window"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type WebviewScrollOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewScroll(_ context.Context, _ *mcp.CallToolRequest, input WebviewScrollInput) (*mcp.CallToolResult, WebviewScrollOutput, resultFailure) {
	r := s.core.Action("webview.scroll").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskScroll{Window: input.Window, X: input.X, Y: input.Y}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScrollOutput{}, e
		}
		return nil, WebviewScrollOutput{}, nil
	}
	return nil, WebviewScrollOutput{Success: true}, nil
}

// --- webview_hover ---

type WebviewHoverInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewHoverOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewHover(_ context.Context, _ *mcp.CallToolRequest, input WebviewHoverInput) (*mcp.CallToolResult, WebviewHoverOutput, resultFailure) {
	r := s.core.Action("webview.hover").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskHover{Window: input.Window, Selector: input.Selector}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewHoverOutput{}, e
		}
		return nil, WebviewHoverOutput{}, nil
	}
	return nil, WebviewHoverOutput{Success: true}, nil
}

// --- webview_select ---

type WebviewSelectInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type WebviewSelectOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewSelect(_ context.Context, _ *mcp.CallToolRequest, input WebviewSelectInput) (*mcp.CallToolResult, WebviewSelectOutput, resultFailure) {
	r := s.core.Action("webview.select").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskSelect{Window: input.Window, Selector: input.Selector, Value: input.Value}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewSelectOutput{}, e
		}
		return nil, WebviewSelectOutput{}, nil
	}
	return nil, WebviewSelectOutput{Success: true}, nil
}

// --- webview_check ---

type WebviewCheckInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Checked  bool   `json:"checked"`
}

type WebviewCheckOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewCheck(_ context.Context, _ *mcp.CallToolRequest, input WebviewCheckInput) (*mcp.CallToolResult, WebviewCheckOutput, resultFailure) {
	r := s.core.Action("webview.check").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskCheck{Window: input.Window, Selector: input.Selector, Checked: input.Checked}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewCheckOutput{}, e
		}
		return nil, WebviewCheckOutput{}, nil
	}
	return nil, WebviewCheckOutput{Success: true}, nil
}

// --- webview_upload ---

type WebviewUploadInput struct {
	Window   string   `json:"window"`
	Selector string   `json:"selector"`
	Paths    []string `json:"paths"`
}

type WebviewUploadOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewUpload(_ context.Context, _ *mcp.CallToolRequest, input WebviewUploadInput) (*mcp.CallToolResult, WebviewUploadOutput, resultFailure) {
	r := s.core.Action("webview.upload_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskUploadFile{Window: input.Window, Selector: input.Selector, Paths: input.Paths}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewUploadOutput{}, e
		}
		return nil, WebviewUploadOutput{}, nil
	}
	return nil, WebviewUploadOutput{Success: true}, nil
}

// --- webview_viewport ---

type WebviewViewportInput struct {
	Window string `json:"window"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type WebviewViewportOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewViewport(_ context.Context, _ *mcp.CallToolRequest, input WebviewViewportInput) (*mcp.CallToolResult, WebviewViewportOutput, resultFailure) {
	r := s.core.Action("webview.set_viewport").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskSetViewport{Window: input.Window, Width: input.Width, Height: input.Height}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewViewportOutput{}, e
		}
		return nil, WebviewViewportOutput{}, nil
	}
	return nil, WebviewViewportOutput{Success: true}, nil
}

// --- webview_console ---

type WebviewConsoleInput struct {
	Window string `json:"window"`
	Level  string `json:"level,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type WebviewConsoleOutput struct {
	Messages []webview.ConsoleMessage `json:"messages"`
}

func (s *Subsystem) webviewConsole(_ context.Context, _ *mcp.CallToolRequest, input WebviewConsoleInput) (*mcp.CallToolResult, WebviewConsoleOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryConsole{Window: input.Window, Level: input.Level, Limit: input.Limit})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewConsoleOutput{}, e
		}
		return nil, WebviewConsoleOutput{}, nil
	}
	msgs, ok := r.Value.([]webview.ConsoleMessage)
	if !ok {
		return nil, WebviewConsoleOutput{}, core.E("mcp.webviewConsole", "unexpected result type", nil)
	}
	return nil, WebviewConsoleOutput{Messages: msgs}, nil
}

// --- webview_console_clear ---

type WebviewConsoleClearInput struct {
	Window string `json:"window"`
}

type WebviewConsoleClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewConsoleClear(_ context.Context, _ *mcp.CallToolRequest, input WebviewConsoleClearInput) (*mcp.CallToolResult, WebviewConsoleClearOutput, resultFailure) {
	r := s.core.Action("webview.clear_console").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskClearConsole{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewConsoleClearOutput{}, e
		}
		return nil, WebviewConsoleClearOutput{}, nil
	}
	return nil, WebviewConsoleClearOutput{Success: true}, nil
}

// --- webview_query ---

type WebviewQueryInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewQueryOutput struct {
	Element *webview.ElementInfo `json:"element"`
}

func (s *Subsystem) webviewQuery(_ context.Context, _ *mcp.CallToolRequest, input WebviewQueryInput) (*mcp.CallToolResult, WebviewQueryOutput, resultFailure) {
	r := s.core.QUERY(webview.QuerySelector{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewQueryOutput{}, e
		}
		return nil, WebviewQueryOutput{}, nil
	}
	el, ok := r.Value.(*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryOutput{}, core.E("mcp.webviewQuery", "unexpected result type", nil)
	}
	return nil, WebviewQueryOutput{Element: el}, nil
}

// --- webview_query_all ---

type WebviewQueryAllInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewQueryAllOutput struct {
	Elements []*webview.ElementInfo `json:"elements"`
}

func (s *Subsystem) webviewQueryAll(_ context.Context, _ *mcp.CallToolRequest, input WebviewQueryAllInput) (*mcp.CallToolResult, WebviewQueryAllOutput, resultFailure) {
	r := s.core.QUERY(webview.QuerySelectorAll{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewQueryAllOutput{}, e
		}
		return nil, WebviewQueryAllOutput{}, nil
	}
	els, ok := r.Value.([]*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryAllOutput{}, core.E("mcp.webviewQueryAll", "unexpected result type", nil)
	}
	return nil, WebviewQueryAllOutput{Elements: els}, nil
}

// --- webview_dom_tree ---

type WebviewDOMTreeInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector,omitempty"`
}

type WebviewDOMTreeOutput struct {
	HTML string `json:"html"`
}

func (s *Subsystem) webviewDOMTree(_ context.Context, _ *mcp.CallToolRequest, input WebviewDOMTreeInput) (*mcp.CallToolResult, WebviewDOMTreeOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryDOMTree{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewDOMTreeOutput{}, e
		}
		return nil, WebviewDOMTreeOutput{}, nil
	}
	html, ok := r.Value.(string)
	if !ok {
		return nil, WebviewDOMTreeOutput{}, core.E("mcp.webviewDOMTree", "unexpected result type", nil)
	}
	return nil, WebviewDOMTreeOutput{HTML: html}, nil
}

// --- webview_url ---

type WebviewURLInput struct {
	Window string `json:"window"`
}

type WebviewURLOutput struct {
	URL string `json:"url"`
}

func (s *Subsystem) webviewURL(_ context.Context, _ *mcp.CallToolRequest, input WebviewURLInput) (*mcp.CallToolResult, WebviewURLOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryURL{Window: input.Window})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewURLOutput{}, e
		}
		return nil, WebviewURLOutput{}, nil
	}
	url, ok := r.Value.(string)
	if !ok {
		return nil, WebviewURLOutput{}, core.E("mcp.webviewURL", "unexpected result type", nil)
	}
	return nil, WebviewURLOutput{URL: url}, nil
}

// --- webview_title ---

type WebviewTitleInput struct {
	Window string `json:"window"`
}

type WebviewTitleOutput struct {
	Title string `json:"title"`
}

func (s *Subsystem) webviewTitle(_ context.Context, _ *mcp.CallToolRequest, input WebviewTitleInput) (*mcp.CallToolResult, WebviewTitleOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryTitle{Window: input.Window})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewTitleOutput{}, e
		}
		return nil, WebviewTitleOutput{}, nil
	}
	title, ok := r.Value.(string)
	if !ok {
		return nil, WebviewTitleOutput{}, core.E("mcp.webviewTitle", "unexpected result type", nil)
	}
	return nil, WebviewTitleOutput{Title: title}, nil
}

// --- webview_list ---

type WebviewListInput struct{}

type WebviewListOutput struct {
	Windows []window.WindowInfo `json:"windows"`
}

func (s *Subsystem) webviewList(_ context.Context, _ *mcp.CallToolRequest, _ WebviewListInput) (*mcp.CallToolResult, WebviewListOutput, resultFailure) {
	r := s.core.QUERY(window.QueryWindowList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewListOutput{}, e
		}
		return nil, WebviewListOutput{}, nil
	}
	windows, ok := r.Value.([]window.WindowInfo)
	if !ok {
		return nil, WebviewListOutput{}, core.E("mcp.webviewList", "unexpected result type", nil)
	}
	return nil, WebviewListOutput{Windows: windows}, nil
}

// --- webview_errors ---

type WebviewErrorsInput struct {
	Window string `json:"window"`
	Limit  int    `json:"limit,omitempty"`
}

type WebviewErrorsOutput struct {
	Errors []webview.ExceptionInfo `json:"errors,omitempty"`
}

func (s *Subsystem) webviewErrors(_ context.Context, _ *mcp.CallToolRequest, input WebviewErrorsInput) (*mcp.CallToolResult, WebviewErrorsOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryExceptions{Window: input.Window, Limit: input.Limit})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewErrorsOutput{}, e
		}
		return nil, WebviewErrorsOutput{}, nil
	}
	errors, ok := r.Value.([]webview.ExceptionInfo)
	if !ok {
		return nil, WebviewErrorsOutput{}, core.E("mcp.webviewErrors", "unexpected result type", nil)
	}
	return nil, WebviewErrorsOutput{Errors: errors}, nil
}

// --- webview_clear_console ---

type WebviewClearConsoleInput struct {
	Window string `json:"window"`
}

type WebviewClearConsoleOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewClearConsole(_ context.Context, _ *mcp.CallToolRequest, input WebviewClearConsoleInput) (*mcp.CallToolResult, WebviewClearConsoleOutput, resultFailure) {
	_, out, err := s.webviewConsoleClear(context.Background(), nil, WebviewConsoleClearInput{Window: input.Window})
	return nil, WebviewClearConsoleOutput{Success: out.Success}, err
}

// --- webview_element_info ---

type WebviewElementInfoInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewElementInfoOutput struct {
	Element *webview.ElementInfo `json:"element"`
}

func (s *Subsystem) webviewElementInfo(_ context.Context, _ *mcp.CallToolRequest, input WebviewElementInfoInput) (*mcp.CallToolResult, WebviewElementInfoOutput, resultFailure) {
	_, out, err := s.webviewQuery(context.Background(), nil, WebviewQueryInput{Window: input.Window, Selector: input.Selector})
	return nil, WebviewElementInfoOutput{Element: out.Element}, err
}

// --- webview_highlight ---

type WebviewHighlightInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Colour   string `json:"colour,omitempty"`
}

type WebviewHighlightOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewHighlight(_ context.Context, _ *mcp.CallToolRequest, input WebviewHighlightInput) (*mcp.CallToolResult, WebviewHighlightOutput, resultFailure) {
	result, err := s.evaluateWebview(input.Window, webview.HighlightScript(input.Selector, input.Colour))
	if err != nil {
		return nil, WebviewHighlightOutput{}, err
	}
	success, _ := result.(bool)
	return nil, WebviewHighlightOutput{Success: success}, nil
}

// --- webview_computed_style ---

type WebviewComputedStyleInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewComputedStyleOutput struct {
	Styles map[string]any `json:"styles"`
}

func (s *Subsystem) webviewComputedStyle(_ context.Context, _ *mcp.CallToolRequest, input WebviewComputedStyleInput) (*mcp.CallToolResult, WebviewComputedStyleOutput, resultFailure) {
	result, err := s.evaluateWebview(input.Window, webview.ComputedStyleScript(input.Selector))
	if err != nil {
		return nil, WebviewComputedStyleOutput{}, err
	}
	styles, err := decodeJSONLike[map[string]any](result)
	if err != nil {
		return nil, WebviewComputedStyleOutput{}, err
	}
	return nil, WebviewComputedStyleOutput{Styles: styles}, nil
}

// --- webview_source ---

type WebviewSourceInput struct {
	Window string `json:"window"`
}

type WebviewSourceOutput struct {
	HTML string `json:"html"`
}

func (s *Subsystem) webviewSource(_ context.Context, _ *mcp.CallToolRequest, input WebviewSourceInput) (*mcp.CallToolResult, WebviewSourceOutput, resultFailure) {
	r := s.core.QUERY(webview.QueryDOMTree{Window: input.Window})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewSourceOutput{}, e
		}
		return nil, WebviewSourceOutput{}, nil
	}
	html, ok := r.Value.(string)
	if !ok {
		return nil, WebviewSourceOutput{}, core.E("mcp.webviewSource", "unexpected result type", nil)
	}
	return nil, WebviewSourceOutput{HTML: html}, nil
}

// --- webview_screenshot_element ---

type WebviewScreenshotElementInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewScreenshotElementOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewScreenshotElement(_ context.Context, _ *mcp.CallToolRequest, input WebviewScreenshotElementInput) (*mcp.CallToolResult, WebviewScreenshotElementOutput, resultFailure) {
	r := s.core.QUERY(webview.QuerySelector{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScreenshotElementOutput{}, e
		}
		return nil, WebviewScreenshotElementOutput{}, nil
	}
	element, ok := r.Value.(*webview.ElementInfo)
	if !ok {
		return nil, WebviewScreenshotElementOutput{}, core.E("mcp.webviewScreenshotElement", "unexpected result type", nil)
	}
	if element == nil || element.BoundingBox == nil {
		return nil, WebviewScreenshotElementOutput{}, core.E("mcp.webviewScreenshotElement", "element not found or has no bounding box", nil)
	}

	r = s.core.Action("webview.screenshot").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskScreenshot{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScreenshotElementOutput{}, e
		}
		return nil, WebviewScreenshotElementOutput{}, nil
	}
	screenshotResult, ok := r.Value.(webview.ScreenshotResult)
	if !ok {
		return nil, WebviewScreenshotElementOutput{}, core.E("mcp.webviewScreenshotElement", "unexpected screenshot result type", nil)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(screenshotResult.Base64)
	if err != nil {
		return nil, WebviewScreenshotElementOutput{}, err
	}
	cropped, err := cropPNGToBoundingBox(imageBytes, element.BoundingBox)
	if err != nil {
		return nil, WebviewScreenshotElementOutput{}, err
	}
	return nil, WebviewScreenshotElementOutput{
		Base64:   base64.StdEncoding.EncodeToString(cropped),
		MimeType: "image/png",
	}, nil
}

// --- webview_pdf ---

type WebviewPDFInput struct {
	Window string `json:"window"`
}

type WebviewPDFOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewPDF(_ context.Context, _ *mcp.CallToolRequest, input WebviewPDFInput) (*mcp.CallToolResult, WebviewPDFOutput, resultFailure) {
	r := s.core.Action("webview.print").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskPrint{Window: input.Window, ToPDF: true}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewPDFOutput{}, e
		}
		return nil, WebviewPDFOutput{}, nil
	}
	result, ok := r.Value.(webview.PrintResult)
	if !ok {
		return nil, WebviewPDFOutput{}, core.E("mcp.webviewPDF", "unexpected result type", nil)
	}
	return nil, WebviewPDFOutput{Base64: result.Base64, MimeType: result.MimeType}, nil
}

// --- webview_print ---

type WebviewPrintInput struct {
	Window string `json:"window"`
}

type WebviewPrintOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewPrint(_ context.Context, _ *mcp.CallToolRequest, input WebviewPrintInput) (*mcp.CallToolResult, WebviewPrintOutput, resultFailure) {
	r := s.core.Action("webview.print").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskPrint{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewPrintOutput{}, e
		}
		return nil, WebviewPrintOutput{}, nil
	}
	return nil, WebviewPrintOutput{Success: true}, nil
}

// --- webview_network ---

type WebviewNetworkInput struct {
	Window string `json:"window"`
	Limit  int    `json:"limit,omitempty"`
}

type WebviewNetworkOutput struct {
	Requests []map[string]any `json:"requests"`
}

func (s *Subsystem) webviewNetwork(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkInput) (*mcp.CallToolResult, WebviewNetworkOutput, resultFailure) {
	result, err := s.evaluateWebview(input.Window, webview.NetworkLogScript(input.Limit))
	if err != nil {
		return nil, WebviewNetworkOutput{}, err
	}
	requests, err := decodeJSONLike[[]map[string]any](result)
	if err != nil {
		return nil, WebviewNetworkOutput{}, err
	}
	return nil, WebviewNetworkOutput{Requests: requests}, nil
}

// --- webview_network_clear ---

type WebviewNetworkClearInput struct {
	Window string `json:"window"`
}

type WebviewNetworkClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNetworkClear(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkClearInput) (*mcp.CallToolResult, WebviewNetworkClearOutput, resultFailure) {
	_, err := s.evaluateWebview(input.Window, webview.NetworkClearScript())
	if err != nil {
		return nil, WebviewNetworkClearOutput{}, err
	}
	return nil, WebviewNetworkClearOutput{Success: true}, nil
}

// --- webview_network_inject ---

type WebviewNetworkInjectInput struct {
	Window string `json:"window"`
}

type WebviewNetworkInjectOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNetworkInject(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkInjectInput) (*mcp.CallToolResult, WebviewNetworkInjectOutput, resultFailure) {
	_, err := s.evaluateWebview(input.Window, webview.NetworkInitScript())
	if err != nil {
		return nil, WebviewNetworkInjectOutput{}, err
	}
	return nil, WebviewNetworkInjectOutput{Success: true}, nil
}

// --- webview_performance ---

type WebviewPerformanceInput struct {
	Window string `json:"window"`
}

type WebviewPerformanceOutput struct {
	Metrics map[string]any `json:"metrics"`
}

func (s *Subsystem) webviewPerformance(_ context.Context, _ *mcp.CallToolRequest, input WebviewPerformanceInput) (*mcp.CallToolResult, WebviewPerformanceOutput, resultFailure) {
	result, err := s.evaluateWebview(input.Window, webview.PerformanceScript())
	if err != nil {
		return nil, WebviewPerformanceOutput{}, err
	}
	metrics, err := decodeJSONLike[map[string]any](result)
	if err != nil {
		return nil, WebviewPerformanceOutput{}, err
	}
	return nil, WebviewPerformanceOutput{Metrics: metrics}, nil
}

// --- webview_resources ---

type WebviewResourcesInput struct {
	Window string `json:"window"`
}

type WebviewResourcesOutput struct {
	Resources []map[string]any `json:"resources"`
}

func (s *Subsystem) webviewResources(_ context.Context, _ *mcp.CallToolRequest, input WebviewResourcesInput) (*mcp.CallToolResult, WebviewResourcesOutput, resultFailure) {
	result, err := s.evaluateWebview(input.Window, webview.ResourcesScript())
	if err != nil {
		return nil, WebviewResourcesOutput{}, err
	}
	resources, err := decodeJSONLike[[]map[string]any](result)
	if err != nil {
		return nil, WebviewResourcesOutput{}, err
	}
	return nil, WebviewResourcesOutput{Resources: resources}, nil
}

// --- webview_devtools_open ---

type WebviewDevToolsOpenInput struct {
	Window string `json:"window"`
}

type WebviewDevToolsOpenOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewDevToolsOpen(_ context.Context, _ *mcp.CallToolRequest, input WebviewDevToolsOpenInput) (*mcp.CallToolResult, WebviewDevToolsOpenOutput, resultFailure) {
	r := s.core.Action("webview.devtools_open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskDevToolsOpen{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewDevToolsOpenOutput{}, e
		}
		return nil, WebviewDevToolsOpenOutput{}, nil
	}
	return nil, WebviewDevToolsOpenOutput{Success: true}, nil
}

// --- webview_devtools_close ---

type WebviewDevToolsCloseInput struct {
	Window string `json:"window"`
}

type WebviewDevToolsCloseOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewDevToolsClose(_ context.Context, _ *mcp.CallToolRequest, input WebviewDevToolsCloseInput) (*mcp.CallToolResult, WebviewDevToolsCloseOutput, resultFailure) {
	r := s.core.Action("webview.devtools_close").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskDevToolsClose{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewDevToolsCloseOutput{}, e
		}
		return nil, WebviewDevToolsCloseOutput{}, nil
	}
	return nil, WebviewDevToolsCloseOutput{Success: true}, nil
}

func (s *Subsystem) evaluateWebview(windowName, script string) (any, resultFailure) {
	r := s.core.Action("webview.evaluate").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskEvaluate{Window: windowName, Script: script}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, e
		}
		return nil, core.E("mcp.evaluateWebview", "webview evaluation failed", nil)
	}
	return r.Value, nil
}

func decodeJSONLike[T any](value any) (T, resultFailure) {
	var out T
	result := core.JSONUnmarshalString(core.JSONMarshalString(value), &out)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return out, err
		}
		return out, core.E("mcp.decodeJSONLike", "failed to decode result", nil)
	}
	return out, nil
}

func cropPNGToBoundingBox(pngData []byte, bbox *webview.BoundingBox) ([]byte, resultFailure) {
	img, err := png.Decode(core.NewBuffer(pngData))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rect := image.Rect(
		maxInt(bounds.Min.X, int(math.Floor(bbox.X))),
		maxInt(bounds.Min.Y, int(math.Floor(bbox.Y))),
		minInt(bounds.Max.X, int(math.Ceil(bbox.X+bbox.Width))),
		minInt(bounds.Max.Y, int(math.Ceil(bbox.Y+bbox.Height))),
	)
	if rect.Empty() {
		return nil, core.E("mcp.cropPNGToBoundingBox", "element bounding box is empty", nil)
	}

	var cropped image.Image
	if subImager, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}); ok {
		cropped = subImager.SubImage(rect)
	} else {
		target := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
		draw.Draw(target, target.Bounds(), img, rect.Min, draw.Src)
		cropped = target
	}

	out := core.NewBuffer()
	if err := png.Encode(out, cropped); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- Registration ---

func (s *Subsystem) registerWebviewTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "webview_eval", Description: "Execute JavaScript in a webview"}, s.webviewEval)
	addTool(s, server, &mcp.Tool{Name: "webview_list", Description: "List webview windows with geometry"}, s.webviewList)
	addTool(s, server, &mcp.Tool{Name: "webview_click", Description: "Click an element in a webview"}, s.webviewClick)
	addTool(s, server, &mcp.Tool{Name: "webview_type", Description: "Type text into an element in a webview"}, s.webviewType)
	addTool(s, server, &mcp.Tool{Name: "webview_navigate", Description: "Navigate a webview to a URL"}, s.webviewNavigate)
	addTool(s, server, &mcp.Tool{Name: "webview_screenshot", Description: "Capture a webview screenshot as base64 PNG"}, s.webviewScreenshot)
	addTool(s, server, &mcp.Tool{Name: "webview_screenshot_element", Description: "Capture a specific DOM element as base64 PNG"}, s.webviewScreenshotElement)
	addTool(s, server, &mcp.Tool{Name: "webview_scroll", Description: "Scroll a webview to an absolute position"}, s.webviewScroll)
	addTool(s, server, &mcp.Tool{Name: "webview_hover", Description: "Hover over an element in a webview"}, s.webviewHover)
	addTool(s, server, &mcp.Tool{Name: "webview_select", Description: "Select an option in a select element"}, s.webviewSelect)
	addTool(s, server, &mcp.Tool{Name: "webview_check", Description: "Check or uncheck a checkbox"}, s.webviewCheck)
	addTool(s, server, &mcp.Tool{Name: "webview_upload", Description: "Upload files to a file input element"}, s.webviewUpload)
	addTool(s, server, &mcp.Tool{Name: "webview_viewport", Description: "Set the webview viewport dimensions"}, s.webviewViewport)
	addTool(s, server, &mcp.Tool{Name: "webview_console", Description: "Get captured console messages from a webview"}, s.webviewConsole)
	addTool(s, server, &mcp.Tool{Name: "webview_console_clear", Description: "Clear captured console messages"}, s.webviewConsoleClear)
	addTool(s, server, &mcp.Tool{Name: "webview_clear_console", Description: "Clear captured console messages"}, s.webviewClearConsole)
	addTool(s, server, &mcp.Tool{Name: "webview_errors", Description: "Get captured JavaScript errors and exceptions"}, s.webviewErrors)
	addTool(s, server, &mcp.Tool{Name: "webview_query", Description: "Find a single DOM element by CSS selector"}, s.webviewQuery)
	addTool(s, server, &mcp.Tool{Name: "webview_element_info", Description: "Get detailed information about a DOM element"}, s.webviewElementInfo)
	addTool(s, server, &mcp.Tool{Name: "webview_query_all", Description: "Find all DOM elements matching a CSS selector"}, s.webviewQueryAll)
	addTool(s, server, &mcp.Tool{Name: "webview_dom_tree", Description: "Get HTML content of a webview"}, s.webviewDOMTree)
	addTool(s, server, &mcp.Tool{Name: "webview_source", Description: "Get the full HTML source of a webview"}, s.webviewSource)
	addTool(s, server, &mcp.Tool{Name: "webview_highlight", Description: "Highlight a DOM element inside the webview"}, s.webviewHighlight)
	addTool(s, server, &mcp.Tool{Name: "webview_computed_style", Description: "Get computed CSS styles for a DOM element"}, s.webviewComputedStyle)
	addTool(s, server, &mcp.Tool{Name: "webview_url", Description: "Get the current URL of a webview"}, s.webviewURL)
	addTool(s, server, &mcp.Tool{Name: "webview_title", Description: "Get the current page title of a webview"}, s.webviewTitle)
	addTool(s, server, &mcp.Tool{Name: "webview_pdf", Description: "Export the current page as PDF"}, s.webviewPDF)
	addTool(s, server, &mcp.Tool{Name: "webview_print", Description: "Open the native print dialog for the current page"}, s.webviewPrint)
	addTool(s, server, &mcp.Tool{Name: "webview_network", Description: "Get recent network activity for the page"}, s.webviewNetwork)
	addTool(s, server, &mcp.Tool{Name: "webview_network_clear", Description: "Clear the injected webview network log"}, s.webviewNetworkClear)
	addTool(s, server, &mcp.Tool{Name: "webview_network_inject", Description: "Inject fetch and XHR interception for detailed network logging"}, s.webviewNetworkInject)
	addTool(s, server, &mcp.Tool{Name: "webview_performance", Description: "Get page performance metrics"}, s.webviewPerformance)
	addTool(s, server, &mcp.Tool{Name: "webview_resources", Description: "List loaded page resources"}, s.webviewResources)
	addTool(s, server, &mcp.Tool{Name: "webview_devtools_open", Description: "Open native developer tools for the window"}, s.webviewDevToolsOpen)
	addTool(s, server, &mcp.Tool{Name: "webview_devtools_close", Description: "Close native developer tools for the window when supported"}, s.webviewDevToolsClose)
}
