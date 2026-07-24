//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package angular_test

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
	httpdisplay "dappco.re/go/render/display/http"
	"dappco.re/go/render/display/http/framework/angular"
	tsengine "dappco.re/go/render/engine/ts"
	corets "dappco.re/go/ts"
)

func ExampleNew() {
	engine, cleanup := exampleEngine()
	defer cleanup()

	serverBundle := inlineModule(`
		let requests = 0;
		export async function reqHandler(request) {
			const path = new URL(request.url).pathname;
			if (path === "/unhandled") return null;
			requests += 1;
			return new Response(
				"<main>" + requests + " " + request.method + " " + path + "</main>",
				{
					status: 201,
					headers: {
						"content-type": "text/html; charset=utf-8",
						"x-request-count": String(requests),
					},
				},
			);
		}
	`)
	app, err := angular.New(core.Background(), engine, serverBundle)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			panic(err)
		}
	}()

	handler := httpdisplay.Handler(nil, "", httpdisplay.WithFramework(app))
	for _, target := range []string{
		"https://example.test/first",
		"https://example.test/second",
		"https://example.test/unhandled",
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(nethttp.MethodGet, target, nil)
		handler.ServeHTTP(recorder, request)
		core.Println(
			recorder.Code,
			recorder.Header().Get("X-Request-Count"),
			core.Trim(recorder.Body.String()),
		)
	}

	// Output:
	// 201 1 <main>1 GET /first</main>
	// 201 2 <main>2 GET /second</main>
	// 404  404 page not found
}

func TestAngularExampleSidecar(t *testing.T) {
	if core.Env("DENO_SOCKET") == "" {
		return
	}
	if err := runExampleSidecar(); err != nil {
		t.Fatal(err)
	}
}

type fixtureRPCRequest struct {
	ID          uint64 `json:"id"`
	Method      string `json:"method"`
	Code        string `json:"code"`
	EntryPoint  string `json:"entry_point"`
	Permissions struct {
		Net []string `json:"net"`
	} `json:"permissions"`
}

type fixtureCallRequest struct {
	ID     uint64            `json:"id"`
	Export string            `json:"export"`
	Args   []core.RawMessage `json:"args"`
}

func exampleEngine() (*tsengine.Engine, func()) {
	tempResult := core.MkdirTemp("", "go-render-angular-example-")
	if !tempResult.OK {
		panic(resultError(tempResult))
	}
	tempDir := tempResult.Value.(string)

	executableResult := core.Executable()
	if !executableResult.OK {
		panic(resultError(executableResult))
	}
	engine, err := tsengine.New(corets.Options{
		DenoPath:       executableResult.Value.(string),
		SocketPath:     core.PathJoin(tempDir, "core.sock"),
		DenoSocketPath: core.PathJoin(tempDir, "deno.sock"),
		StoreDBPath:    ":memory:",
		SidecarArgs:    []string{"-test.run=^TestAngularExampleSidecar$"},
	})
	if err != nil {
		if removeResult := core.RemoveAll(tempDir); !removeResult.OK {
			panic(resultError(removeResult))
		}
		panic(err)
	}
	return engine, func() {
		if err := engine.Close(); err != nil {
			panic(err)
		}
		if removeResult := core.RemoveAll(tempDir); !removeResult.OK {
			panic(resultError(removeResult))
		}
	}
}

func runExampleSidecar() error {
	socketPath := core.Env("DENO_SOCKET")
	core.Remove(socketPath)
	listenResult := core.NetListen("unix", socketPath)
	if !listenResult.OK {
		return core.E("angular.example.sidecar", "listen on CoreTS socket", resultError(listenResult))
	}
	listener := listenResult.Value.(core.Listener)
	defer listener.Close()

	connection, err := listener.Accept()
	if err != nil {
		return core.E("angular.example.sidecar", "accept CoreTS connection", err)
	}
	defer connection.Close()
	return serveExampleRPC(connection)
}

func serveExampleRPC(connection core.Conn) error {
	reader := core.NewBufReader(connection)
	bridges := make(map[string]core.Conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil
		}
		var request fixtureRPCRequest
		if result := core.JSONUnmarshalString(line, &request); !result.OK {
			return core.E("angular.example.sidecar", "decode CoreTS request", resultError(result))
		}

		result := map[string]any{"ok": true}
		switch request.Method {
		case "Ping":
		case "LoadModule":
			bridge, loadErr := loadExampleModule(request)
			if loadErr != nil {
				result["ok"] = false
				result["error"] = loadErr.Error()
			} else {
				bridges[request.Code] = bridge
			}
		case "UnloadModule":
			if bridge := bridges[request.Code]; bridge != nil {
				if closeErr := bridge.Close(); closeErr != nil {
					result["ok"] = false
				}
				delete(bridges, request.Code)
			}
		case "ModuleStatus":
			result["code"] = request.Code
			if bridges[request.Code] == nil {
				result["status"] = "STOPPED"
			} else {
				result["status"] = "RUNNING"
			}
		default:
			result["ok"] = false
			result["error"] = "unknown CoreTS method"
		}
		if err := writeRPCResult(connection, request.ID, result); err != nil {
			return err
		}
	}
}

func loadExampleModule(request fixtureRPCRequest) (core.Conn, error) {
	if len(request.Permissions.Net) != 1 {
		return nil, core.E("angular.example.sidecar", "resident module needs one call address", nil)
	}
	if err := validateExampleWrapper(request.EntryPoint); err != nil {
		return nil, err
	}
	dialResult := core.NetDial("tcp", request.Permissions.Net[0])
	if !dialResult.OK {
		return nil, core.E("angular.example.sidecar", "connect resident bridge", resultError(dialResult))
	}
	connection := dialResult.Value.(core.Conn)
	if result := core.WriteString(connection, "{\"ready\":true}\n"); !result.OK {
		connection.Close()
		return nil, core.E("angular.example.sidecar", "write resident ready frame", resultError(result))
	}
	go serveExampleCalls(connection)
	return connection, nil
}

func validateExampleWrapper(entryPoint string) error {
	parseResult := core.URLParse(entryPoint)
	if !parseResult.OK {
		return core.E("angular.example.sidecar", "parse resident wrapper URL", resultError(parseResult))
	}
	wrapperURL := parseResult.Value.(*core.URL)
	if !core.HasSuffix(wrapperURL.Path, "context.js") {
		return core.E("angular.example.sidecar", "resident wrapper is not JavaScript", nil)
	}
	readResult := core.ReadFile(wrapperURL.Path)
	if !readResult.OK {
		return core.E("angular.example.sidecar", "read resident wrapper", resultError(readResult))
	}
	wrapper := core.AsString(readResult.Value.([]byte))
	if !core.Contains(wrapper, "residentModule = await import(") ||
		!core.Contains(wrapper, "new Request") ||
		!core.Contains(wrapper, "value instanceof Response") {
		return core.E("angular.example.sidecar", "resident wrapper lacks Web request bridge", nil)
	}
	return nil
}

func serveExampleCalls(connection core.Conn) {
	reader := core.NewBufReader(connection)
	requests := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		var request fixtureCallRequest
		if result := core.JSONUnmarshalString(line, &request); !result.OK {
			return
		}
		response := map[string]any{
			"id": request.ID,
			"ok": true,
		}
		if request.Export != "reqHandler" || len(request.Args) != 1 {
			response["ok"] = false
			response["error"] = "reqHandler export is required"
		} else {
			response["value"] = exampleAngularResponse(request.Args[0], &requests)
		}
		result := core.JSONMarshal(response)
		if !result.OK {
			return
		}
		if writeResult := core.WriteString(connection, core.AsString(result.Value.([]byte))+"\n"); !writeResult.OK {
			return
		}
	}
}

func exampleAngularResponse(raw core.RawMessage, requests *int) any {
	var request struct {
		Type   string `json:"__go_render_type"`
		URL    string `json:"url"`
		Method string `json:"method"`
	}
	if result := core.JSONUnmarshal(raw, &request); !result.OK || request.Type != "Request" {
		return nil
	}
	parseResult := core.URLParse(request.URL)
	if !parseResult.OK {
		return nil
	}
	path := parseResult.Value.(*core.URL).Path
	if path == "/unhandled" {
		return nil
	}
	*requests++
	count := core.Itoa(*requests)
	return map[string]any{
		"status":     nethttp.StatusCreated,
		"statusText": "Created",
		"headers": [][2]string{
			{"content-type", "text/html; charset=utf-8"},
			{"x-request-count", count},
		},
		"body": []byte("<main>" + count + " " + request.Method + " " + path + "</main>"),
	}
}

func writeRPCResult(connection core.Conn, id uint64, result map[string]any) error {
	responseResult := core.JSONMarshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	})
	if !responseResult.OK {
		return core.E("angular.example.sidecar", "encode CoreTS response", resultError(responseResult))
	}
	writeResult := core.WriteString(connection, core.AsString(responseResult.Value.([]byte))+"\n")
	if !writeResult.OK {
		return core.E("angular.example.sidecar", "write CoreTS response", resultError(writeResult))
	}
	return nil
}

func inlineModule(source string) string {
	return "data:application/javascript;base64," + core.Base64Encode(core.AsBytes(source))
}

func resultError(result core.Result) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.E("angular.example", result.Error(), nil)
}
