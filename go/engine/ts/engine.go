//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	core "dappco.re/go"
	corets "dappco.re/go/ts"
)

const (
	renderModulePrefix = "go-render-"
	renderGroupPrefix  = "go-render-ts-"
)

// Renderer is the content-producing contract shared by the server-side
// display packages.
type Renderer interface {
	Render(ctx core.Context, entry string, input any) ([]byte, error)
}

// Engine owns a CoreTS service and its managed Deno sidecar.
type Engine struct {
	mu      core.RWMutex
	app     *core.Core
	service *corets.Service
	tempDir string
	appRoot string
	closed  bool
}

// New starts a CoreTS service and its configured Deno sidecar.
//
// SidecarArgs must launch a CoreTS-compatible Deno runtime. A typical runtime
// uses: run, -A, --unstable-worker-options, and CoreTS's runtime/main.ts.
func New(opts corets.Options) (*Engine, error) {
	if len(opts.SidecarArgs) == 0 {
		return nil, core.E("ts.New", "CoreTS SidecarArgs are required", nil)
	}

	tempResult := core.MkdirTemp("", renderGroupPrefix)
	if !tempResult.OK {
		return nil, core.E("ts.New", "create render workspace", resultError(tempResult))
	}
	tempDir := tempResult.Value.(string)

	if opts.SocketPath == "" {
		opts.SocketPath = core.PathJoin(tempDir, "core.sock")
	}
	if opts.DenoSocketPath == "" {
		opts.DenoSocketPath = core.PathJoin(tempDir, "deno.sock")
	}

	app := core.New(core.WithService(corets.NewServiceFactory(opts)))
	service, ok := core.ServiceFor[*corets.Service](app, "ts")
	if !ok || service == nil {
		removeResult := core.RemoveAll(tempDir)
		if !removeResult.OK {
			return nil, core.E("ts.New", "remove render workspace after CoreTS registration failure", resultError(removeResult))
		}
		return nil, core.E("ts.New", "register CoreTS service", nil)
	}

	startResult := app.ServiceStartup(core.Background(), nil)
	if !startResult.OK {
		shutdownResult := app.ServiceShutdown(core.Background())
		if !shutdownResult.OK {
			return nil, core.E("ts.New", "clean up after CoreTS startup failure", resultError(shutdownResult))
		}
		removeResult := core.RemoveAll(tempDir)
		if !removeResult.OK {
			return nil, core.E("ts.New", "remove render workspace after CoreTS startup failure", resultError(removeResult))
		}
		return nil, core.E("ts.New", "start CoreTS service", resultError(startResult))
	}
	if service.DenoClient() == nil {
		shutdownResult := app.ServiceShutdown(core.Background())
		if !shutdownResult.OK {
			return nil, core.E("ts.New", "stop CoreTS service without a Deno client", resultError(shutdownResult))
		}
		removeResult := core.RemoveAll(tempDir)
		if !removeResult.OK {
			return nil, core.E("ts.New", "remove render workspace without a Deno client", resultError(removeResult))
		}
		return nil, core.E("ts.New", "CoreTS Deno client is not connected", nil)
	}

	return &Engine{
		app:     app,
		service: service,
		tempDir: tempDir,
		appRoot: canonicalPath(opts.AppRoot),
	}, nil
}

// Render loads a permission-scoped wrapper module in CoreTS, invokes the
// entry's render export with input, and returns its HTML or byte output.
//
// The entry may be a local .ts/.tsx/.js module or a data: module. It must
// export render(input), a default function, or a default object with a render
// method. String, []byte/Uint8Array, ArrayBuffer, Blob, Response, and
// JSON-serialisable results are supported.
func (e *Engine) Render(ctx core.Context, entry string, input any) ([]byte, error) {
	if e == nil {
		return nil, core.E("ts.Engine.Render", "engine is nil", nil)
	}
	if ctx == nil {
		return nil, core.E("ts.Engine.Render", "context is nil", nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, core.E("ts.Engine.Render", "render context is done", err)
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.closed || e.service == nil {
		return nil, core.E("ts.Engine.Render", "engine is closed", nil)
	}

	entryPoint, entryReads, err := e.entryPoint(entry)
	if err != nil {
		return nil, err
	}
	inputJSON, err := marshalJSON("ts.Engine.Render", "serialise render input", input)
	if err != nil {
		return nil, err
	}

	randomResult := core.RandomString(16)
	if !randomResult.OK {
		return nil, core.E("ts.Engine.Render", "create module identity", resultError(randomResult))
	}
	identity := renderModulePrefix + randomResult.Value.(string)
	workspaceResult := core.MkdirTemp(e.tempDir, identity+"-")
	if !workspaceResult.OK {
		return nil, core.E("ts.Engine.Render", "create isolated render workspace", resultError(workspaceResult))
	}
	renderDir := workspaceResult.Value.(string)
	defer func() {
		if removeResult := core.RemoveAll(renderDir); !removeResult.OK {
			core.Warn("TypeScript render workspace cleanup failed", "path", renderDir, "err", resultError(removeResult))
		}
	}()
	wrapperPath := core.PathJoin(renderDir, "render.ts")
	outputPath := core.PathJoin(renderDir, "render.out")

	wrapper, err := renderWrapper(entryPoint, inputJSON, outputPath)
	if err != nil {
		return nil, err
	}
	writeResult := core.WriteFile(wrapperPath, core.AsBytes(wrapper), 0o600)
	if !writeResult.OK {
		return nil, core.E("ts.Engine.Render", "write render wrapper", resultError(writeResult))
	}

	wrapperEntry := fileURL(wrapperPath)
	permissions := corets.ModulePermissions{
		Read:  append([]string{renderDir}, entryReads...),
		Write: []string{renderDir},
	}
	loadResponse, err := e.service.LoadModule(identity, wrapperEntry, permissions)
	if err != nil {
		return nil, core.E("ts.Engine.Render", "load TypeScript render module", err)
	}
	defer func() {
		unloadResponse, unloadErr := e.service.UnloadModule(identity)
		if unloadErr != nil {
			core.Warn("TypeScript render module unload failed", "module", identity, "err", unloadErr)
			return
		}
		if unloadResponse == nil || !unloadResponse.Ok {
			core.Warn("TypeScript render module did not unload", "module", identity)
		}
	}()
	if loadResponse == nil || !loadResponse.Ok {
		message := "TypeScript render module was rejected"
		if loadResponse != nil && loadResponse.Error != "" {
			message = loadResponse.Error
		}
		return nil, core.E("ts.Engine.Render", message, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, core.E("ts.Engine.Render", "render context is done", err)
	}

	readResult := core.ReadFile(outputPath)
	if !readResult.OK {
		return nil, core.E("ts.Engine.Render", "read rendered output", resultError(readResult))
	}
	return readResult.Value.([]byte), nil
}

// Close stops the CoreTS service and sidecar and removes the render workspace.
// It is safe to call Close more than once.
func (e *Engine) Close() error {
	if e == nil {
		return core.E("ts.Engine.Close", "engine is nil", nil)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}
	e.closed = true

	var closeErr error
	if e.app != nil {
		shutdownResult := e.app.ServiceShutdown(core.Background())
		if !shutdownResult.OK {
			closeErr = core.E("ts.Engine.Close", "stop CoreTS service", resultError(shutdownResult))
		}
	}
	if e.tempDir != "" {
		removeResult := core.RemoveAll(e.tempDir)
		if !removeResult.OK && closeErr == nil {
			closeErr = core.E("ts.Engine.Close", "remove render workspace", resultError(removeResult))
		}
	}

	e.app = nil
	e.service = nil
	e.tempDir = ""
	return closeErr
}

func (e *Engine) entryPoint(entry string) (string, []string, error) {
	entry = core.Trim(entry)
	if entry == "" {
		return "", nil, core.E("ts.Engine.Render", "entry is required", nil)
	}
	if core.HasPrefix(entry, "data:") {
		return entry, nil, nil
	}

	path := entry
	if core.HasPrefix(entry, "file:") {
		parseResult := core.URLParse(entry)
		if !parseResult.OK {
			return "", nil, core.E("ts.Engine.Render", "parse file entry URL", resultError(parseResult))
		}
		parsed := parseResult.Value.(*core.URL)
		if parsed.Scheme != "file" || (parsed.Host != "" && parsed.Host != "localhost") {
			return "", nil, core.E("ts.Engine.Render", "entry must be a local file URL", nil)
		}
		path = parsed.Path
	} else if !core.PathIsAbs(entry) {
		parseResult := core.URLParse(entry)
		if parseResult.OK && parseResult.Value.(*core.URL).Scheme != "" {
			return "", nil, core.E("ts.Engine.Render", "remote module entries are not permitted", nil)
		}
	}

	absoluteResult := core.PathAbs(path)
	if !absoluteResult.OK {
		return "", nil, core.E("ts.Engine.Render", "resolve module entry", resultError(absoluteResult))
	}
	absolute := absoluteResult.Value.(string)
	resolvedResult := core.PathEvalSymlinks(absolute)
	if !resolvedResult.OK {
		return "", nil, core.E("ts.Engine.Render", "resolve module entry symlinks", resultError(resolvedResult))
	}
	resolved := resolvedResult.Value.(string)

	statResult := core.Stat(resolved)
	if !statResult.OK {
		return "", nil, core.E("ts.Engine.Render", "stat module entry", resultError(statResult))
	}
	if statResult.Value.(core.FsFileInfo).IsDir() {
		return "", nil, core.E("ts.Engine.Render", "module entry is a directory", nil)
	}

	readRoot := core.PathDir(resolved)
	if e.appRoot != "" {
		if !corets.CheckPath(resolved, []string{e.appRoot}) {
			return "", nil, core.E("ts.Engine.Render", "module entry is outside CoreTS AppRoot", nil)
		}
		readRoot = e.appRoot
	}
	return fileURL(resolved), []string{readRoot}, nil
}

func renderWrapper(entryPoint, inputJSON, outputPath string) (string, error) {
	entryJSON, err := marshalJSON("ts.Engine.Render", "serialise module entry", entryPoint)
	if err != nil {
		return "", err
	}
	outputJSON, err := marshalJSON("ts.Engine.Render", "serialise output path", outputPath)
	if err != nil {
		return "", err
	}

	return core.Concat(
		"// Generated by dappco.re/go/html/engine/ts.\n",
		"export async function init() {\n",
		"  const module = await import(", entryJSON, ");\n",
		"  let render = module.render;\n",
		"  if (typeof render !== \"function\" && typeof module.default === \"function\") render = module.default;\n",
		"  if (typeof render !== \"function\" && module.default && typeof module.default.render === \"function\") render = module.default.render.bind(module.default);\n",
		"  if (typeof render !== \"function\") throw new Error(\"render entry must export render(input) or a default renderer\");\n",
		"  const value = await render(", inputJSON, ");\n",
		"  const encoder = new TextEncoder();\n",
		"  let output;\n",
		"  if (typeof value === \"string\") output = encoder.encode(value);\n",
		"  else if (value instanceof Uint8Array) output = value;\n",
		"  else if (value instanceof ArrayBuffer) output = new Uint8Array(value);\n",
		"  else if (typeof Blob !== \"undefined\" && value instanceof Blob) output = new Uint8Array(await value.arrayBuffer());\n",
		"  else if (typeof Response !== \"undefined\" && value instanceof Response) output = new Uint8Array(await value.arrayBuffer());\n",
		"  else if (value === undefined || value === null) output = new Uint8Array();\n",
		"  else output = encoder.encode(JSON.stringify(value));\n",
		"  await Deno.writeFile(", outputJSON, ", output, { create: true, createNew: true, write: true });\n",
		"}\n",
	), nil
}

func marshalJSON(operation, message string, value any) (string, error) {
	result := core.JSONMarshal(value)
	if !result.OK {
		return "", core.E(operation, message, resultError(result))
	}
	return core.AsString(result.Value.([]byte)), nil
}

func fileURL(path string) string {
	return (&core.URL{Scheme: "file", Path: core.PathToSlash(path)}).String()
}

func canonicalPath(path string) string {
	if path == "" {
		return ""
	}
	absoluteResult := core.PathAbs(path)
	if !absoluteResult.OK {
		return path
	}
	absolute := absoluteResult.Value.(string)
	resolvedResult := core.PathEvalSymlinks(absolute)
	if resolvedResult.OK {
		return resolvedResult.Value.(string)
	}
	return absolute
}

func resultError(result core.Result) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.E("ts", result.Error(), nil)
}
