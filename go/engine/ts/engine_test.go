//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	"testing"

	core "dappco.re/go"
	corets "dappco.re/go/ts"
)

func TestEngine_New_Good(t *testing.T) {
	engine, err := New(smokeOptions(t))
	core.RequireNoError(t, err)
	core.AssertNotNil(t, engine)
	core.AssertNotNil(t, engine.service)
	core.AssertTrue(t, engine.service.Sidecar().IsRunning())
	core.AssertNoError(t, engine.Close())
}

func TestEngine_New_Bad(t *testing.T) {
	engine, err := New(corets.Options{})
	core.AssertNil(t, engine)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "SidecarArgs")
}

func TestEngine_New_Ugly(t *testing.T) {
	engine, err := New(corets.Options{
		DenoPath:    core.PathJoin(t.TempDir(), "missing-deno"),
		SidecarArgs: []string{"run", "runtime.ts"},
	})
	core.AssertNil(t, engine)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "start CoreTS service")
}

func TestEngine_Render_Good(t *testing.T) {
	engine := newSmokeEngine(t)
	defer func() { core.AssertNoError(t, engine.Close()) }()

	entry := inlineModule(`export function render(input: {name: string}) {
		return "<main>Hello, " + input.name + "</main>";
	}`)
	output, err := engine.Render(core.Background(), entry, map[string]any{"name": "Deno"})
	core.AssertNoError(t, err)
	core.AssertEqual(t, "<main>Hello, Deno</main>", core.AsString(output))
}

func TestEngine_Render_Bad(t *testing.T) {
	engine := &Engine{closed: true}
	output, err := engine.Render(core.Background(), "entry.ts", nil)
	core.AssertNil(t, output)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "closed")
}

func TestEngine_Render_Ugly(t *testing.T) {
	engine := newSmokeEngine(t)
	defer func() { core.AssertNoError(t, engine.Close()) }()

	entry := inlineModule(`export default function() {
		return new Uint8Array([0, 1, 2, 255]);
	}`)
	output, err := engine.Render(core.Background(), entry, nil)
	core.AssertNoError(t, err)
	core.AssertEqual(t, []byte{0, 1, 2, 255}, output)
}

func TestEngine_Close_Good(t *testing.T) {
	engine := newSmokeEngine(t)
	sidecar := engine.service.Sidecar()
	core.AssertTrue(t, sidecar.IsRunning())
	core.AssertNoError(t, engine.Close())
	core.AssertFalse(t, sidecar.IsRunning())
}

func TestEngine_Close_Bad(t *testing.T) {
	var engine *Engine
	err := engine.Close()
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "engine is nil")
	core.AssertContains(t, err.Error(), "Engine.Close")
}

func TestEngine_Close_Ugly(t *testing.T) {
	engine := newSmokeEngine(t)
	core.AssertNoError(t, engine.Close())
	core.AssertNoError(t, engine.Close())
	core.AssertTrue(t, engine.closed)
}

func newSmokeEngine(t *testing.T) *Engine {
	t.Helper()

	engine, err := New(smokeOptions(t))
	core.RequireNoError(t, err)
	return engine
}

func smokeOptions(t *testing.T) corets.Options {
	t.Helper()

	denoPath := findDeno(t)
	tempDir := t.TempDir()
	runtimePath := core.PathJoin(tempDir, "runtime.ts")
	writeResult := core.WriteFile(runtimePath, core.AsBytes(smokeRuntime), 0o600)
	core.RequireTrue(t, writeResult.OK, writeResult.Error())

	return corets.Options{
		DenoPath:       denoPath,
		SocketPath:     core.PathJoin(tempDir, "core.sock"),
		DenoSocketPath: core.PathJoin(tempDir, "deno.sock"),
		StoreDBPath:    ":memory:",
		SidecarArgs: []string{
			"run",
			"-A",
			"--unstable-worker-options",
			runtimePath,
		},
	}
}

func findDeno(t *testing.T) string {
	t.Helper()

	directories := core.Split(core.Env("PATH"), string(core.PathListSeparator))
	if home := core.Env("DIR_HOME"); home != "" {
		directories = append(directories, core.PathJoin(home, ".deno", "bin"))
	}
	for _, directory := range directories {
		for _, name := range []string{"deno", "deno.exe"} {
			candidate := core.PathJoin(directory, name)
			if core.Stat(candidate).OK {
				return candidate
			}
		}
	}
	t.Skip("deno not installed")
	return ""
}

func inlineModule(source string) string {
	return "data:application/typescript;base64," + core.Base64Encode(core.AsBytes(source))
}

const smokeRuntime = `
const socketPath = Deno.env.get("DENO_SOCKET");
if (!socketPath) throw new Error("DENO_SOCKET is required");
try { await Deno.remove(socketPath); } catch (error) {
	if (!(error instanceof Deno.errors.NotFound)) throw error;
}

const listener = Deno.listen({ transport: "unix", path: socketPath });
const encoder = new TextEncoder();
const decoder = new TextDecoder();

for await (const connection of listener) {
	void serve(connection);
}

async function serve(connection: Deno.Conn) {
	let pending = "";
	const buffer = new Uint8Array(64 * 1024);
	try {
		for (;;) {
			const count = await connection.read(buffer);
			if (count === null) return;
			pending += decoder.decode(buffer.subarray(0, count), { stream: true });
			for (;;) {
				const newline = pending.indexOf("\n");
				if (newline < 0) break;
				const line = pending.slice(0, newline);
				pending = pending.slice(newline + 1);
				if (line.length === 0) continue;
				const request = JSON.parse(line);
				const result = await dispatch(request);
				await connection.write(encoder.encode(JSON.stringify({
					jsonrpc: "2.0",
					id: request.id,
					result,
				}) + "\n"));
			}
		}
	} finally {
		connection.close();
	}
}

async function dispatch(request: Record<string, unknown>) {
	switch (request.method) {
		case "Ping":
			return { ok: true };
		case "LoadModule":
			return await loadModule(
				String(request.code ?? ""),
				String(request.entry_point ?? ""),
				(request.permissions ?? {}) as Record<string, string[]>,
			);
		case "UnloadModule":
			return { ok: true };
		case "ModuleStatus":
			return { code: String(request.code ?? ""), status: "RUNNING" };
		default:
			throw new Error("unknown method: " + request.method);
	}
}

async function loadModule(
	code: string,
	entryPoint: string,
	permissions: Record<string, string[]>,
) {
	const workerSource = ` + "`" + `
self.onmessage = async (event) => {
	try {
		const module = await import(event.data);
		if (typeof module.init === "function") await module.init({});
		self.postMessage({ ok: true });
	} catch (error) {
		self.postMessage({
			ok: false,
			error: error instanceof Error ? error.message : String(error),
		});
	}
};
` + "`" + `;
	const workerURL = "data:application/javascript;base64," + btoa(workerSource);
	const worker = new Worker(workerURL, {
		type: "module",
		name: code,
		deno: {
			permissions: {
				read: permissions.read ?? [],
				write: permissions.write ?? [],
				net: permissions.net ?? [],
				run: permissions.run ?? [],
				env: false,
				sys: false,
				ffi: false,
			},
		},
	});
	try {
		return await new Promise((resolve) => {
			worker.onmessage = (event) => resolve(event.data);
			worker.onerror = (event) => resolve({ ok: false, error: event.message });
			worker.postMessage(entryPoint);
		});
	} finally {
		worker.terminate();
	}
}
`
