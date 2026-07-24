//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	"testing"

	core "dappco.re/go"
)

func TestContext_Invoke_Good(t *testing.T) {
	context := fixtureContext(t)
	defer func() { core.AssertNoError(t, context.Close()) }()

	var output int
	err := context.Invoke(core.Background(), "increment", &output, 2)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 2, output)
}

func TestContext_Invoke_Bad(t *testing.T) {
	context := fixtureContext(t)
	defer func() { core.AssertNoError(t, context.Close()) }()

	var output int
	err := context.Invoke(core.Background(), "missing", &output)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "export")
}

func TestContext_Invoke_Ugly(t *testing.T) {
	context := fixtureContext(t)
	defer func() { core.AssertNoError(t, context.Close()) }()

	var first int
	var second int
	core.AssertNoError(t, context.Invoke(core.Background(), "increment", &first, 3))
	core.AssertNoError(t, context.Invoke(core.Background(), "increment", &second, 4))
	core.AssertEqual(t, 3, first)
	core.AssertEqual(t, 7, second)
}

func TestContext_Invoke_WebRequest(t *testing.T) {
	context := fixtureContext(t)
	defer func() { core.AssertNoError(t, context.Close()) }()

	var response WebResponse
	err := context.Invoke(core.Background(), "inspectRequest", &response, WebRequest{
		URL:    "https://example.test/account?tab=profile",
		Method: "POST",
		Headers: [][2]string{
			{"Content-Type", "text/plain"},
			{"X-Trace", "one"},
			{"X-Trace", "two"},
		},
		Body: []byte("hello"),
	})
	core.AssertNoError(t, err)
	core.AssertEqual(t, 202, response.Status)
	core.AssertEqual(t, "Accepted", response.StatusText)
	core.AssertEqual(t, [][2]string{{"content-type", "text/html"}, {"set-cookie", "one=1"}, {"set-cookie", "two=2"}}, response.Headers)
	core.AssertEqual(t, []byte("<main>POST /account hello</main>"), response.Body)
}

func TestContext_Close_Good(t *testing.T) {
	context := fixtureContext(t)
	core.AssertNoError(t, context.Close())

	var output int
	err := context.Invoke(core.Background(), "increment", &output, 1)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "closed")
}

func TestContext_Close_Bad(t *testing.T) {
	var context *Context
	err := context.Close()
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "context is nil")
}

func TestContext_Close_Ugly(t *testing.T) {
	context := fixtureContext(t)
	core.AssertNoError(t, context.Close())
	core.AssertNoError(t, context.Close())
}

func TestContext_contextWrapper_WebValues(t *testing.T) {
	wrapper, err := contextWrapper(
		"data:application/javascript,export%20function%20reqHandler(){}",
		"127.0.0.1",
		8080,
	)
	core.AssertNoError(t, err)
	core.AssertContains(t, wrapper, "new Request")
	core.AssertContains(t, wrapper, "value instanceof Response")
	core.AssertContains(t, wrapper, "getSetCookie")
}

func fixtureContext(t *testing.T) *Context {
	t.Helper()

	client, server := core.NetPipe()
	go serveFixtureContext(server)

	return &Context{
		conn:   client,
		reader: core.NewBufReader(client),
	}
}

func serveFixtureContext(connection core.Conn) {
	defer connection.Close()

	reader := core.NewBufReader(connection)
	count := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var request struct {
			ID     uint64            `json:"id"`
			Export string            `json:"export"`
			Args   []core.RawMessage `json:"args"`
		}
		if result := core.JSONUnmarshalString(line, &request); !result.OK {
			return
		}

		response := map[string]any{
			"id": request.ID,
			"ok": true,
		}
		switch request.Export {
		case "increment":
			if len(request.Args) > 0 {
				var by int
				if result := core.JSONUnmarshal(request.Args[0], &by); !result.OK {
					return
				}
				count += by
			}
			response["value"] = count
		case "inspectRequest":
			var input struct {
				Type   string `json:"__go_render_type"`
				URL    string `json:"url"`
				Method string `json:"method"`
				Body   []byte `json:"body"`
			}
			if len(request.Args) == 0 {
				return
			}
			if result := core.JSONUnmarshal(request.Args[0], &input); !result.OK {
				return
			}
			if input.Type != "Request" || input.URL != "https://example.test/account?tab=profile" {
				response["ok"] = false
				response["error"] = "Web Request was not encoded"
				break
			}
			response["value"] = map[string]any{
				"status":     202,
				"statusText": "Accepted",
				"headers": [][2]string{
					{"content-type", "text/html"},
					{"set-cookie", "one=1"},
					{"set-cookie", "two=2"},
				},
				"body": []byte("<main>" + input.Method + " /account " + core.AsString(input.Body) + "</main>"),
			}
		default:
			response["ok"] = false
			response["error"] = "export is not callable"
		}

		result := core.JSONMarshal(response)
		if !result.OK {
			return
		}
		if write := core.WriteString(connection, core.AsString(result.Value.([]byte))+"\n"); !write.OK {
			return
		}
	}
}
