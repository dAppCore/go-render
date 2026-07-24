//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package framework

import (
	nethttp "net/http"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/render/engine/ts"
)

type loaderContract struct{}

func (loaderContract) Load(core.Context, string) (*ts.Context, error) {
	return &ts.Context{}, nil
}

type rendererContract struct{}

func (rendererContract) Render(core.Context, *nethttp.Request) (*Response, error) {
	return &Response{Status: nethttp.StatusNoContent}, nil
}

func (rendererContract) Close() error {
	return nil
}

func TestFramework_Contracts(t *testing.T) {
	var loader Loader = loaderContract{}
	var renderer Renderer = rendererContract{}
	context, err := loader.Load(core.Background(), "server.mjs")
	core.AssertNoError(t, err)
	core.AssertNotNil(t, context)
	response, err := renderer.Render(core.Background(), &nethttp.Request{})
	core.AssertNoError(t, err)
	core.AssertEqual(t, nethttp.StatusNoContent, response.Status)
	core.AssertNoError(t, renderer.Close())
}

func TestFramework_Response(t *testing.T) {
	response := Response{
		Status: nethttp.StatusAccepted,
		Header: nethttp.Header{"X-Framework": {"Angular"}},
		Body:   []byte("<main>accepted</main>"),
	}
	core.AssertEqual(t, nethttp.StatusAccepted, response.Status)
	core.AssertEqual(t, "Angular", response.Header.Get("X-Framework"))
	core.AssertEqual(t, []byte("<main>accepted</main>"), response.Body)
}
