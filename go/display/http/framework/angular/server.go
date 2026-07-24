//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

// Package angular adapts Angular's non-Node RequestHandlerFunction SSR
// contract to the Go HTTP display.
package angular

import (
	nethttp "net/http"

	core "dappco.re/go"
	"dappco.re/go/render/display/http/framework"
	"dappco.re/go/render/engine/ts"
)

const requestHandlerExport = "reqHandler"

type module interface {
	Invoke(core.Context, string, any, ...any) error
	Close() error
}

type renderer struct {
	mu       core.Mutex
	module   module
	closed   bool
	closeErr error
}

type requestBody struct {
	reader core.Reader
}

// New loads serverBundle once into a resident CoreTS context and returns an
// Angular RequestHandlerFunction renderer.
func New(ctx core.Context, loader framework.Loader, serverBundle string) (framework.Renderer, error) {
	if ctx == nil {
		return nil, core.E("angular.New", "context is nil", nil)
	}
	if loader == nil {
		return nil, core.E("angular.New", "CoreTS loader is nil", nil)
	}
	serverBundle = core.Trim(serverBundle)
	if serverBundle == "" {
		return nil, core.E("angular.New", "server bundle is required", nil)
	}

	context, err := loader.Load(ctx, serverBundle)
	if err != nil {
		return nil, core.E("angular.New", "load Angular server bundle", err)
	}
	if context == nil {
		return nil, core.E("angular.New", "load Angular server bundle returned no context", nil)
	}
	return &renderer{module: context}, nil
}

func (r *renderer) Render(ctx core.Context, request *nethttp.Request) (*framework.Response, error) {
	if r == nil {
		return nil, core.E("angular.renderer.Render", "renderer is nil", nil)
	}
	if ctx == nil {
		return nil, core.E("angular.renderer.Render", "context is nil", nil)
	}
	if request == nil {
		return nil, core.E("angular.renderer.Render", "request is nil", nil)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed || r.module == nil {
		return nil, core.E("angular.renderer.Render", "renderer is closed", nil)
	}

	webRequest, err := webRequestFromHTTP(request)
	if err != nil {
		return nil, err
	}
	var webResponse *ts.WebResponse
	if err := r.module.Invoke(ctx, requestHandlerExport, &webResponse, webRequest); err != nil {
		return nil, core.E("angular.renderer.Render", "invoke Angular request handler", err)
	}
	if webResponse == nil {
		return nil, nil
	}

	header := make(nethttp.Header)
	for _, value := range webResponse.Headers {
		header.Add(value[0], value[1])
	}
	return &framework.Response{
		Status: webResponse.Status,
		Header: header,
		Body:   append([]byte(nil), webResponse.Body...),
	}, nil
}

func (r *renderer) Close() error {
	if r == nil {
		return core.E("angular.renderer.Close", "renderer is nil", nil)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return r.closeErr
	}
	r.closed = true
	if r.module == nil {
		r.closeErr = core.E("angular.renderer.Close", "resident module is nil", nil)
		return r.closeErr
	}
	if err := r.module.Close(); err != nil {
		r.closeErr = core.E("angular.renderer.Close", "close resident Angular module", err)
	}
	r.module = nil
	return r.closeErr
}

func webRequestFromHTTP(request *nethttp.Request) (ts.WebRequest, error) {
	if request.URL == nil {
		return ts.WebRequest{}, core.E("angular.webRequestFromHTTP", "request URL is nil", nil)
	}

	requestURL := request.URL.String()
	if !request.URL.IsAbs() {
		host := core.Trim(request.Host)
		if host == "" {
			return ts.WebRequest{}, core.E("angular.webRequestFromHTTP", "request host is required", nil)
		}
		scheme := "http"
		if request.TLS != nil {
			scheme = "https"
		}
		requestURL = core.Concat(scheme, "://", host, request.URL.RequestURI())
	}

	headers := make([][2]string, 0, len(request.Header)+1)
	hasHost := false
	for name, values := range request.Header {
		if core.EqualFold(name, "Host") {
			hasHost = true
		}
		for _, value := range values {
			headers = append(headers, [2]string{name, value})
		}
	}
	if !hasHost && request.Host != "" {
		headers = append(headers, [2]string{"Host", request.Host})
	}

	body, err := readRequestBody(request)
	if err != nil {
		return ts.WebRequest{}, err
	}
	return ts.WebRequest{
		URL:     requestURL,
		Method:  request.Method,
		Headers: headers,
		Body:    body,
	}, nil
}

func readRequestBody(request *nethttp.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	readResult := core.ReadAll(request.Body)
	if !readResult.OK {
		return nil, core.E("angular.readRequestBody", "read request body", resultError(readResult))
	}
	content, ok := readResult.Value.(string)
	if !ok {
		return nil, core.E("angular.readRequestBody", "request body is not text-compatible bytes", nil)
	}
	body := append([]byte(nil), core.AsBytes(content)...)
	request.Body = &requestBody{reader: core.NewBufferReader(body)}
	return body, nil
}

func (b *requestBody) Read(buffer []byte) (int, error) {
	return b.reader.Read(buffer)
}

func (b *requestBody) Close() error {
	return nil
}

func resultError(result core.Result) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.E("angular", result.Error(), nil)
}
