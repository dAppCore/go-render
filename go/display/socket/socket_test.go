//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package socket

import (
	nethttp "net/http"
	"testing"

	core "dappco.re/go"
	httpdisplay "dappco.re/go/html/display/http"
)

type socketRenderer struct {
	output []byte
	entry  string
}

func (r *socketRenderer) Render(_ core.Context, entry string, _ any) ([]byte, error) {
	r.entry = entry
	return r.output, nil
}

func TestSocket_Serve_Good(t *testing.T) {
	listenResult := core.NetListen("tcp", "127.0.0.1:0")
	core.RequireTrue(t, listenResult.OK, listenResult.Error())
	listener := listenResult.Value.(core.Listener)
	renderer := &socketRenderer{output: []byte("<main>socket</main>")}
	done := make(chan error, 1)
	go func() {
		done <- Serve(renderer, listener, httpdisplay.WithEntry("server.ts"))
	}()

	client := &nethttp.Client{Timeout: 2 * core.Second}
	response, err := client.Get("http://" + listener.Addr().String() + "/")
	core.RequireNoError(t, err)
	readResult := core.ReadAll(response.Body)
	core.RequireTrue(t, readResult.OK, readResult.Error())
	core.AssertEqual(t, "<main>socket</main>", readResult.Value.(string))
	core.AssertEqual(t, "server.ts", renderer.entry)

	core.AssertNoError(t, listener.Close())
	core.AssertNoError(t, <-done)
}

func TestSocket_Serve_Bad(t *testing.T) {
	listenResult := core.NetListen("tcp", "127.0.0.1:0")
	core.RequireTrue(t, listenResult.OK, listenResult.Error())
	listener := listenResult.Value.(core.Listener)
	defer listener.Close()

	err := Serve(nil, listener)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "engine is nil")
	core.AssertContains(t, err.Error(), "socket.Serve")
}

func TestSocket_Serve_Ugly(t *testing.T) {
	renderer := &socketRenderer{output: []byte("unused")}
	err := Serve(renderer, nil, httpdisplay.WithEntry("server.ts"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "listener is nil")
	core.AssertEqual(t, "", renderer.entry)
	core.AssertEqual(t, "unused", core.AsString(renderer.output))
}
