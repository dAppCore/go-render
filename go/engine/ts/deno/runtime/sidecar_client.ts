// GoService gRPC client — Deno calls Go for sandboxed I/O, KV, process exec, and core:// resolution.
//
// Usage:
//   const client = createGoServiceClient("unix:///run/core/deno.sock");
//   const { data } = await client.readFile("/workspace/data.json", "my-module");
//   const text = new TextDecoder().decode(data);
//   const { value, found } = await client.storeGet("app", "theme", "my-module");
//   await client.storeSet("app", "theme", "dark", "my-module");
//   const { output, exit_code } = await client.exec("git", ["log", "--oneline"], ".", [], "my-module");
//   const { result_json } = await client.resolveScheme("core://settings/theme", "my-module");
//   client.close();
//
// Transport: Unix domain socket (primary), TCP loopback fallback.
// Proto: proto/core_sidecar.proto (package core.sidecar — canonical contract).

import { grpc, protoLoader } from "./grpc.ts";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const PROTO_PATH = join(__dirname, "..", "proto", "core_sidecar.proto");

let packageDef: protoLoader.PackageDefinition | null = null;

// getProto loads the core_sidecar package definition once and caches it.
function getProto(): GoServiceConstructor {
  if (!packageDef) {
    packageDef = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: Number,
      enums: String,
      defaults: true,
      oneofs: true,
    });
  }
  // deno-lint-ignore no-explicit-any
  const pkg = grpc.loadPackageDefinition(packageDef) as any;
  return pkg.core.sidecar.GoService as GoServiceConstructor;
}

// deno-lint-ignore no-explicit-any
type GoServiceConstructor = new (address: string, credentials: any) => any;

export interface FileEntry {
  name: string;
  is_dir: boolean;
}

export interface GoServiceClient {
  // Exposed for test doubles — callers should use the typed helpers below.
  // deno-lint-ignore no-explicit-any
  raw: any;

  // ReadFile — read file contents through go-io permission gates.
  //   const { data } = await client.readFile("/workspace/data.json", "module-a");
  //   const text = new TextDecoder().decode(data);
  readFile(path: string, moduleCode: string): Promise<{ data: Uint8Array; error: string }>;

  // WriteFile — write file contents through go-io permission gates.
  //   await client.writeFile("/workspace/out.txt", new TextEncoder().encode("hello"), 0o644, "module-a");
  writeFile(
    path: string,
    data: Uint8Array,
    mode: number,
    moduleCode: string,
  ): Promise<{ ok: boolean; error: string }>;

  // ListFiles — list directory entries through go-io permission gates.
  //   const { entries } = await client.listFiles("/workspace", "module-a");
  listFiles(
    path: string,
    moduleCode: string,
  ): Promise<{ entries: FileEntry[]; error: string }>;

  // StoreGet — retrieve a KV value from go-store.
  //   const { value, found } = await client.storeGet("app", "theme", "module-a");
  storeGet(
    group: string,
    key: string,
    moduleCode?: string,
  ): Promise<{ value: string; found: boolean; error: string }>;

  // StoreSet — persist a KV value to go-store.
  //   await client.storeSet("app", "theme", "dark", "module-a");
  storeSet(
    group: string,
    key: string,
    value: string,
    moduleCode?: string,
  ): Promise<{ ok: boolean; error: string }>;

  // Exec — run a process through go-process permission gates.
  //   const { output, exit_code } = await client.exec("git", ["log"], ".", [], "module-a");
  //   const text = new TextDecoder().decode(output);
  exec(
    cmd: string,
    args: string[],
    dir: string,
    env: string[],
    moduleCode: string,
  ): Promise<{
    output: Uint8Array;
    exit_code: number;
    error: string;
  }>;

  // ResolveScheme — resolve a core:// URI via Go route handlers.
  //   const { result_json } = await client.resolveScheme("core://settings", "module-a");
  resolveScheme(
    uri: string,
    moduleCode?: string,
  ): Promise<{ result_json: string; error: string }>;

  // close — release the underlying gRPC channel.
  close(): void;
}

function promisify<T>(
  // deno-lint-ignore no-explicit-any
  rawClient: any,
  method: string,
  request: unknown,
): Promise<T> {
  return new Promise<T>((resolve, reject) => {
    rawClient[method](
      request,
      (err: Error | null, response: T) => {
        if (err) reject(err);
        else resolve(response);
      },
    );
  });
}

// createGoServiceClient — connect to the Go gRPC server.
//
//   // Unix domain socket (preferred — matches Go server default):
//   const client = createGoServiceClient("unix:///run/core/deno.sock");
//
//   // TCP loopback fallback:
//   const client = createGoServiceClient("localhost:50051");
export function createGoServiceClient(address: string): GoServiceClient {
  const GoService = getProto();
  const raw = new GoService(address, grpc.credentials.createInsecure());

  return {
    raw,

    readFile(path: string, moduleCode: string) {
      return promisify(raw, "ReadFile", { path, module_code: moduleCode });
    },

    writeFile(path: string, data: Uint8Array, mode: number, moduleCode: string) {
      return promisify(raw, "WriteFile", {
        path,
        data,
        mode,
        module_code: moduleCode,
      });
    },

    listFiles(path: string, moduleCode: string) {
      return promisify(raw, "ListFiles", { path, module_code: moduleCode });
    },

    storeGet(group: string, key: string, moduleCode = "") {
      return promisify(raw, "StoreGet", {
        group,
        key,
        module_code: moduleCode,
      });
    },

    storeSet(group: string, key: string, value: string, moduleCode = "") {
      return promisify(raw, "StoreSet", {
        group,
        key,
        value,
        module_code: moduleCode,
      });
    },

    exec(cmd: string, args: string[], dir: string, env: string[], moduleCode: string) {
      return promisify(raw, "Exec", { cmd, args, dir, env, module_code: moduleCode });
    },

    resolveScheme(uri: string, moduleCode = "") {
      return promisify(raw, "ResolveScheme", { uri, module_code: moduleCode });
    },

    close() {
      raw.close();
    },
  };
}
