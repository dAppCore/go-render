// Tests for the GoService client (Deno → Go direction).
// Follows TestFilename_Function_{Good,Bad,Ugly} naming convention.

import { createGoServiceClient } from "./sidecar_client.ts";

function assert(condition: boolean, message: string): void {
  if (!condition) throw new Error(message);
}

function assertEquals<T>(actual: T, expected: T, message: string): void {
  if (!Object.is(actual, expected)) {
    throw new Error(`${message}: expected ${String(expected)}, got ${String(actual)}`);
  }
}

// Shared RPC stub type used by every test.
type StubRawClient = Record<
  string,
  (request: unknown, callback: (error: Error | null, response: unknown) => void) => void
> & { close(): void };

Deno.test("TestSidecarClient_createGoServiceClient_Good", async () => {
  const client = createGoServiceClient("unix:///tmp/core-test.sock");
  const calls: Array<{ method: string; request: unknown }> = [];
  let closed = false;
  const raw = client.raw as StubRawClient;

  const fileBytes = new TextEncoder().encode("file-data");
  const outputBytes = new TextEncoder().encode("ok\n");

  // Stub all methods
  raw.ReadFile = (request, callback) => {
    calls.push({ method: "ReadFile", request });
    callback(null, { data: fileBytes, error: "" });
  };
  raw.WriteFile = (request, callback) => {
    calls.push({ method: "WriteFile", request });
    callback(null, { ok: true, error: "" });
  };
  raw.ListFiles = (request, callback) => {
    calls.push({ method: "ListFiles", request });
    callback(null, { entries: [{ name: "a.ts", is_dir: false }], error: "" });
  };
  raw.StoreGet = (request, callback) => {
    calls.push({ method: "StoreGet", request });
    callback(null, { value: "dark", found: true, error: "" });
  };
  raw.StoreSet = (request, callback) => {
    calls.push({ method: "StoreSet", request });
    callback(null, { ok: true, error: "" });
  };
  raw.Exec = (request, callback) => {
    calls.push({ method: "Exec", request });
    callback(null, { output: outputBytes, exit_code: 0, error: "" });
  };
  raw.ResolveScheme = (request, callback) => {
    calls.push({ method: "ResolveScheme", request });
    callback(null, { result_json: '{"route":"settings"}', error: "" });
  };
  raw.close = () => { closed = true; };

  // readFile — response uses bytes data field
  const readResult = await client.readFile("/workspace/data.json", "module-a");
  assertEquals(
    new TextDecoder().decode(readResult.data),
    "file-data",
    "readFile() data should decode to the expected string",
  );
  assertEquals(readResult.error, "", "readFile() error should be empty on success");

  // writeFile — sends bytes data + mode
  const writeResult = await client.writeFile(
    "/workspace/out.txt",
    new TextEncoder().encode("hello"),
    0o644,
    "module-a",
  );
  assertEquals(writeResult.ok, true, "writeFile() should proxy the WriteFile RPC");
  assertEquals(writeResult.error, "", "writeFile() error should be empty on success");
  // Verify the wire request included mode
  const writeReq = calls.find((c) => c.method === "WriteFile")?.request as Record<string, unknown>;
  assertEquals(writeReq?.mode as number, 0o644, "writeFile() should pass mode to the RPC");

  // listFiles — FileEntry has no size field
  const list = await client.listFiles("/workspace", "module-a");
  assertEquals(
    list.entries.length,
    1,
    "listFiles() should return the FileEntry array",
  );
  assertEquals(list.entries[0]!.name, "a.ts", "listFiles() entry name");
  assertEquals(list.entries[0]!.is_dir, false, "listFiles() entry is_dir");

  // storeGet — default module_code
  const storeGetResult = await client.storeGet("app", "theme");
  assertEquals(storeGetResult.value, "dark", "storeGet() value");
  assertEquals(storeGetResult.found, true, "storeGet() found");
  assertEquals(storeGetResult.error, "", "storeGet() error should be empty on success");
  assertEquals(
    (calls.find((c) => c.method === "StoreGet")?.request as Record<string, unknown>)?.module_code,
    "",
    "storeGet() should default module_code to empty string",
  );

  // storeSet
  const storeSetResult = await client.storeSet("app", "theme", "dark", "module-b");
  assertEquals(storeSetResult.ok, true, "storeSet() should proxy the StoreSet RPC");
  assertEquals(storeSetResult.error, "", "storeSet() error should be empty on success");

  // exec — canonical fields: cmd, args, dir, env; response: output bytes + exit_code + error
  const execResult = await client.exec("git", ["log", "--oneline"], ".", [], "module-a");
  assertEquals(execResult.exit_code, 0, "exec() exit_code");
  assertEquals(
    new TextDecoder().decode(execResult.output),
    "ok\n",
    "exec() output should decode to the expected string",
  );
  assertEquals(execResult.error, "", "exec() error should be empty on success");
  // Verify the wire request used canonical field names
  const execReq = calls.find((c) => c.method === "Exec")?.request as Record<string, unknown>;
  assert("cmd" in execReq!, "exec() request should use 'cmd' field (not 'command')");
  assertEquals(execReq?.dir as string, ".", "exec() should pass dir to the RPC");

  // resolveScheme — response uses result_json (not payload/content_type)
  const schemeResult = await client.resolveScheme("core://settings");
  assertEquals(
    schemeResult.result_json,
    '{"route":"settings"}',
    "resolveScheme() should return result_json",
  );
  assertEquals(schemeResult.error, "", "resolveScheme() error should be empty on success");

  // close
  client.close();
  assert(closed, "close() should delegate to the raw client");

  // Verify dispatch order
  const methodOrder = calls.map((c) => c.method);
  assertEquals(
    JSON.stringify(methodOrder),
    JSON.stringify([
      "ReadFile",
      "WriteFile",
      "ListFiles",
      "StoreGet",
      "StoreSet",
      "Exec",
      "ResolveScheme",
    ]),
    "createGoServiceClient() should map each helper to the expected RPC method",
  );
});

Deno.test("TestSidecarClient_createGoServiceClient_Bad", async () => {
  const client = createGoServiceClient("unix:///tmp/core-test.sock");
  const raw = client.raw as StubRawClient;

  raw.StoreGet = (_request, callback) => {
    callback(new Error("connection refused"), undefined);
  };

  let message = "";
  try {
    await client.storeGet("app", "theme");
  } catch (err) {
    message = err instanceof Error ? err.message : String(err);
  }

  assertEquals(
    message,
    "connection refused",
    "storeGet() should propagate RPC errors as rejected promises",
  );
});

Deno.test("TestSidecarClient_createGoServiceClient_Ugly", () => {
  // createGoServiceClient() must always expose the raw client even before any calls.
  const client = createGoServiceClient("unix:///tmp/core-test.sock");
  client.close();
  assert(client.raw !== undefined, "createGoServiceClient() should always expose the raw client");
});
