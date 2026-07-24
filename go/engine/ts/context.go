//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	core "dappco.re/go"
	corets "dappco.re/go/ts"
)

const contextConnectTimeout = 10 * core.Second

// Context is a resident CoreTS module whose exported functions can be invoked
// repeatedly without re-importing the module.
type Context struct {
	callMu    core.Mutex
	stateMu   core.RWMutex
	service   *corets.Service
	conn      core.Conn
	reader    *core.BufReader
	identity  string
	workspace string
	nextID    uint64
	closed    bool
	closeErr  error
}

type contextRequest struct {
	ID     uint64 `json:"id"`
	Export string `json:"export"`
	Args   []any  `json:"args"`
}

type contextResponse struct {
	ID    uint64          `json:"id"`
	OK    bool            `json:"ok"`
	Value core.RawMessage `json:"value"`
	Error string          `json:"error"`
}

type contextReady struct {
	Ready bool   `json:"ready"`
	Error string `json:"error"`
}

// Invoke calls an exported function in the resident module. Arguments are
// serialised to JSON and result, when non-nil, receives the decoded result.
func (c *Context) Invoke(ctx core.Context, export string, result any, args ...any) error {
	if c == nil {
		return core.E("ts.Context.Invoke", "context is nil", nil)
	}
	if ctx == nil {
		return core.E("ts.Context.Invoke", "call context is nil", nil)
	}
	if err := ctx.Err(); err != nil {
		return core.E("ts.Context.Invoke", "call context is done", err)
	}
	export = core.Trim(export)
	if export == "" {
		return core.E("ts.Context.Invoke", "export is required", nil)
	}

	c.callMu.Lock()
	defer c.callMu.Unlock()

	c.stateMu.Lock()
	if c.closed || c.conn == nil || c.reader == nil {
		c.stateMu.Unlock()
		return core.E("ts.Context.Invoke", "context is closed", nil)
	}
	connection := c.conn
	reader := c.reader
	c.nextID++
	requestID := c.nextID
	c.stateMu.Unlock()

	requestResult := core.JSONMarshal(contextRequest{
		ID:     requestID,
		Export: export,
		Args:   contextArguments(args),
	})
	if !requestResult.OK {
		return core.E("ts.Context.Invoke", "serialise call", resultError(requestResult))
	}
	writeResult := core.WriteString(connection, core.AsString(requestResult.Value.([]byte))+"\n")
	if !writeResult.OK {
		return core.E("ts.Context.Invoke", "write call", resultError(writeResult))
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return core.E("ts.Context.Invoke", "read call result", err)
	}
	if err := ctx.Err(); err != nil {
		return core.E("ts.Context.Invoke", "call context is done", err)
	}

	var response contextResponse
	responseResult := core.JSONUnmarshalString(line, &response)
	if !responseResult.OK {
		return core.E("ts.Context.Invoke", "decode call result", resultError(responseResult))
	}
	if response.ID != requestID {
		return core.E("ts.Context.Invoke", "call result identity does not match", nil)
	}
	if !response.OK {
		message := core.Trim(response.Error)
		if message == "" {
			message = "export invocation failed"
		}
		return core.E("ts.Context.Invoke", message, nil)
	}
	if result == nil {
		return nil
	}
	decodeResult := core.JSONUnmarshal(response.Value, result)
	if !decodeResult.OK {
		return core.E("ts.Context.Invoke", "decode exported result", resultError(decodeResult))
	}
	return nil
}

// Close unloads the resident module and releases its call connection and
// isolated workspace. It is safe to call Close more than once.
func (c *Context) Close() error {
	if c == nil {
		return core.E("ts.Context.Close", "context is nil", nil)
	}
	return c.close()
}

func (c *Context) close() error {
	c.stateMu.Lock()
	if c.closed {
		err := c.closeErr
		c.stateMu.Unlock()
		return err
	}
	c.closed = true
	connection := c.conn
	service := c.service
	identity := c.identity
	workspace := c.workspace
	c.conn = nil
	c.reader = nil
	c.service = nil
	c.workspace = ""
	c.stateMu.Unlock()

	var closeErr error
	if connection != nil {
		if err := connection.Close(); err != nil {
			closeErr = core.E("ts.Context.Close", "close call connection", err)
		}
	}
	if service != nil && identity != "" {
		response, err := service.UnloadModule(identity)
		if err != nil {
			closeErr = core.ErrorJoin(
				closeErr,
				core.E("ts.Context.Close", "unload resident module", err),
			)
		} else if response == nil || !response.Ok {
			closeErr = core.ErrorJoin(
				closeErr,
				core.E("ts.Context.Close", "resident module did not unload", nil),
			)
		}
	}
	if workspace != "" {
		removeResult := core.RemoveAll(workspace)
		if !removeResult.OK {
			closeErr = core.ErrorJoin(
				closeErr,
				core.E("ts.Context.Close", "remove resident workspace", resultError(removeResult)),
			)
		}
	}

	c.stateMu.Lock()
	c.closeErr = closeErr
	c.stateMu.Unlock()
	return closeErr
}

func (e *Engine) loadContext(ctx core.Context, entry string) (_ *Context, returnErr error) {
	if e == nil {
		return nil, core.E("ts.Engine.Load", "engine is nil", nil)
	}
	if ctx == nil {
		return nil, core.E("ts.Engine.Load", "context is nil", nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, core.E("ts.Engine.Load", "load context is done", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed || e.service == nil {
		return nil, core.E("ts.Engine.Load", "engine is closed", nil)
	}

	entryPoint, entryReads, err := e.entryPoint("ts.Engine.Load", entry)
	if err != nil {
		return nil, err
	}
	randomResult := core.RandomString(16)
	if !randomResult.OK {
		return nil, core.E("ts.Engine.Load", "create context identity", resultError(randomResult))
	}
	identity := contextModulePrefix + randomResult.Value.(string)
	workspaceResult := core.MkdirTemp(e.tempDir, identity+"-")
	if !workspaceResult.OK {
		return nil, core.E("ts.Engine.Load", "create resident workspace", resultError(workspaceResult))
	}
	workspace := workspaceResult.Value.(string)

	listenResult := core.NetListen("tcp", "127.0.0.1:0")
	if !listenResult.OK {
		removeResult := core.RemoveAll(workspace)
		if !removeResult.OK {
			return nil, core.E("ts.Engine.Load", "clean up after call listener failure", resultError(removeResult))
		}
		return nil, core.E("ts.Engine.Load", "listen for resident calls", resultError(listenResult))
	}
	listener := listenResult.Value.(core.Listener)
	moduleLoaded := false
	keepWorkspace := false
	defer func() {
		if closeErr := listener.Close(); closeErr != nil && returnErr == nil {
			returnErr = core.E("ts.Engine.Load", "close resident call listener", closeErr)
		}
		if keepWorkspace {
			return
		}
		if moduleLoaded {
			if response, unloadErr := e.service.UnloadModule(identity); unloadErr != nil {
				core.Warn("TypeScript resident module cleanup failed", "module", identity, "err", unloadErr)
			} else if response == nil || !response.Ok {
				core.Warn("TypeScript resident module did not unload during cleanup", "module", identity)
			}
		}
		if removeResult := core.RemoveAll(workspace); !removeResult.OK {
			core.Warn("TypeScript resident workspace cleanup failed", "path", workspace, "err", resultError(removeResult))
		}
	}()

	tcpAddress, ok := listener.Addr().(*core.TCPAddr)
	if !ok {
		return nil, core.E("ts.Engine.Load", "resident call listener is not TCP", nil)
	}
	wrapper, err := contextWrapper(entryPoint, tcpAddress.IP.String(), tcpAddress.Port)
	if err != nil {
		return nil, err
	}
	wrapperPath := core.PathJoin(workspace, "context.js")
	writeResult := core.WriteFile(wrapperPath, core.AsBytes(wrapper), 0o600)
	if !writeResult.OK {
		return nil, core.E("ts.Engine.Load", "write resident wrapper", resultError(writeResult))
	}

	address := listener.Addr().String()
	loadResponse, err := e.service.LoadModule(identity, fileURL(wrapperPath), corets.ModulePermissions{
		Read: append([]string{workspace}, entryReads...),
		Net:  []string{address},
	})
	if err != nil {
		return nil, core.E("ts.Engine.Load", "load resident module", err)
	}
	if loadResponse == nil || !loadResponse.Ok {
		message := "TypeScript resident module was rejected"
		if loadResponse != nil && loadResponse.Error != "" {
			message = loadResponse.Error
		}
		return nil, core.E("ts.Engine.Load", message, nil)
	}
	moduleLoaded = true

	deadline := core.Now().Add(contextConnectTimeout)
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		deadline = contextDeadline
	}
	if tcpListener, ok := listener.(*core.TCPListener); ok {
		if err := tcpListener.SetDeadline(deadline); err != nil {
			return nil, core.E("ts.Engine.Load", "set resident connection deadline", err)
		}
	}
	connection, err := listener.Accept()
	if err != nil {
		return nil, core.E("ts.Engine.Load", "accept resident call connection", err)
	}
	reader := core.NewBufReader(connection)
	line, err := reader.ReadString('\n')
	if err != nil {
		if closeErr := connection.Close(); closeErr != nil {
			core.Warn("TypeScript resident connection cleanup failed", "err", closeErr)
		}
		return nil, core.E("ts.Engine.Load", "read resident ready frame", err)
	}
	var ready contextReady
	readyResult := core.JSONUnmarshalString(line, &ready)
	if !readyResult.OK {
		if closeErr := connection.Close(); closeErr != nil {
			core.Warn("TypeScript resident connection cleanup failed", "err", closeErr)
		}
		return nil, core.E("ts.Engine.Load", "decode resident ready frame", resultError(readyResult))
	}
	if !ready.Ready {
		if closeErr := connection.Close(); closeErr != nil {
			core.Warn("TypeScript resident connection cleanup failed", "err", closeErr)
		}
		message := core.Trim(ready.Error)
		if message == "" {
			message = "resident module did not become ready"
		}
		return nil, core.E("ts.Engine.Load", message, nil)
	}

	resident := &Context{
		service:   e.service,
		conn:      connection,
		reader:    reader,
		identity:  identity,
		workspace: workspace,
	}
	if e.contexts == nil {
		e.contexts = make(map[*Context]struct{})
	}
	e.contexts[resident] = struct{}{}
	keepWorkspace = true
	return resident, nil
}

func contextWrapper(entryPoint, hostname string, port int) (string, error) {
	entryJSON, err := marshalJSON("ts.Engine.Load", "serialise module entry", entryPoint)
	if err != nil {
		return "", err
	}
	hostnameJSON, err := marshalJSON("ts.Engine.Load", "serialise call hostname", hostname)
	if err != nil {
		return "", err
	}

	return core.Concat(
		"// Generated by dappco.re/go/render/engine/ts.\n",
		"let residentModule;\n",
		"let residentConnection;\n",
		"const residentEncoder = new TextEncoder();\n",
		"const residentDecoder = new TextDecoder();\n",
		"export async function init() {\n",
		"  residentModule = await import(", entryJSON, ");\n",
		"  residentConnection = await Deno.connect({ transport: \"tcp\", hostname: ", hostnameJSON, ", port: ", core.Itoa(port), " });\n",
		"  await writeResidentFrame({ ready: true });\n",
		"  void serveResidentFrames().catch(() => {});\n",
		"}\n",
		"async function serveResidentFrames() {\n",
		"  let pending = \"\";\n",
		"  const buffer = new Uint8Array(64 * 1024);\n",
		"  for (;;) {\n",
		"    const count = await residentConnection.read(buffer);\n",
		"    if (count === null) return;\n",
		"    pending += residentDecoder.decode(buffer.subarray(0, count), { stream: true });\n",
		"    for (;;) {\n",
		"      const newline = pending.indexOf(\"\\n\");\n",
		"      if (newline < 0) break;\n",
		"      const line = pending.slice(0, newline);\n",
		"      pending = pending.slice(newline + 1);\n",
		"      if (line.length === 0) continue;\n",
		"      await dispatchResidentCall(line);\n",
		"    }\n",
		"  }\n",
		"}\n",
		"async function dispatchResidentCall(line) {\n",
		"  let request;\n",
		"  try {\n",
		"    request = JSON.parse(line);\n",
		"    const callable = residentModule[request.export];\n",
		"    if (typeof callable !== \"function\") throw new Error(\"export is not callable: \" + request.export);\n",
		"    const args = (request.args ?? []).map(decodeResidentArgument);\n",
		"    const value = await callable(...args);\n",
		"    await writeResidentFrame({ id: request.id, ok: true, value: await encodeResidentValue(value) });\n",
		"  } catch (error) {\n",
		"    await writeResidentFrame({\n",
		"      id: request?.id ?? 0,\n",
		"      ok: false,\n",
		"      error: error instanceof Error ? error.message : String(error),\n",
		"    });\n",
		"  }\n",
		"}\n",
		"function decodeResidentArgument(value) {\n",
		"  if (!value || value.__go_render_type !== \"Request\") return value;\n",
		"  const method = String(value.method || \"GET\").toUpperCase();\n",
		"  const init = { method, headers: new Headers(value.headers ?? []) };\n",
		"  if (value.body && method !== \"GET\" && method !== \"HEAD\") init.body = residentBase64ToBytes(value.body);\n",
		"  return new Request(value.url, init);\n",
		"}\n",
		"async function encodeResidentValue(value) {\n",
		"  if (!(value instanceof Response)) return value === undefined ? null : value;\n",
		"  const headers = [];\n",
		"  const hasSetCookie = typeof value.headers.getSetCookie === \"function\";\n",
		"  for (const header of value.headers.entries()) {\n",
		"    if (header[0].toLowerCase() !== \"set-cookie\" || !hasSetCookie) headers.push(header);\n",
		"  }\n",
		"  if (hasSetCookie) {\n",
		"    for (const cookie of value.headers.getSetCookie()) headers.push([\"set-cookie\", cookie]);\n",
		"  }\n",
		"  const body = new Uint8Array(await value.arrayBuffer());\n",
		"  return { status: value.status, statusText: value.statusText, headers, body: residentBytesToBase64(body) };\n",
		"}\n",
		"function residentBase64ToBytes(value) {\n",
		"  const binary = atob(value);\n",
		"  const bytes = new Uint8Array(binary.length);\n",
		"  for (let index = 0; index < binary.length; index += 1) bytes[index] = binary.charCodeAt(index);\n",
		"  return bytes;\n",
		"}\n",
		"function residentBytesToBase64(bytes) {\n",
		"  const chunks = [];\n",
		"  const chunkSize = 0x6000;\n",
		"  for (let offset = 0; offset < bytes.length; offset += chunkSize) {\n",
		"    let binary = \"\";\n",
		"    for (const byte of bytes.subarray(offset, offset + chunkSize)) binary += String.fromCharCode(byte);\n",
		"    chunks.push(btoa(binary));\n",
		"  }\n",
		"  return chunks.join(\"\");\n",
		"}\n",
		"async function writeResidentFrame(value) {\n",
		"  const data = residentEncoder.encode(JSON.stringify(value) + \"\\n\");\n",
		"  let offset = 0;\n",
		"  while (offset < data.length) offset += await residentConnection.write(data.subarray(offset));\n",
		"}\n",
	), nil
}

func contextArguments(args []any) []any {
	encoded := make([]any, len(args))
	for index, argument := range args {
		switch request := argument.(type) {
		case WebRequest:
			encoded[index] = webRequestValue(request)
		case *WebRequest:
			if request == nil {
				encoded[index] = nil
			} else {
				encoded[index] = webRequestValue(*request)
			}
		default:
			encoded[index] = argument
		}
	}
	return encoded
}
