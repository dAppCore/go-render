// DenoService gRPC server — Go calls Deno for lifecycle events, SSR, and TS eval.
//
// Usage:
//   const handlers: DenoServiceHandlers = {
//     onStart:         (phase, reason) => { console.error("started", phase, reason); },
//     onStop:          (phase, reason) => { console.error("stopped", phase, reason); },
//     onConfigChange:  (key, value) => { applyConfig(key, value); },
//     render:          async (component, props) => renderToHTML(component, props),
//     eval:            async (expression) => evalExpression(expression),
//   };
//   const server = await startDenoServiceServer("/run/core/deno.grpc.sock", handlers);
//   // ... later ...
//   server.close();
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
// deno-lint-ignore no-explicit-any
function getProto(): any {
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
  return pkg.core.sidecar;
}

export interface DenoServiceHandlers {
  // onStart — Go notifies Deno the application is starting.
  //   handlers.onStart("start", "boot")
  onStart(phase: string, reason: string): void | Promise<void>;

  // onStop — Go notifies Deno the application is stopping.
  //   handlers.onStop("stop", "shutdown")
  onStop(phase: string, reason: string): void | Promise<void>;

  // onConfigChange — Go pushes a single config value change to Deno.
  //   handlers.onConfigChange("theme.accent", "#6366f1")
  onConfigChange(key: string, value: string): void | Promise<void>;

  // render — Go requests server-side rendering of a component.
  //   const html = await handlers.render("dashboard", '{"user":"snider"}')
  render(component: string, propsJson: string): Promise<string>;

  // eval — Go requests evaluation of a TypeScript expression.
  //   const json = await handlers.eval("1 + 1")
  eval(expression: string): Promise<string>;
}

export interface DenoServiceServer {
  // close — drain pending calls and release the listening socket.
  close(): void;
}

// startDenoServiceServer — bind the gRPC server and begin accepting calls from Go.
//
//   // Unix domain socket (matches Go DenoClient default):
//   const srv = await startDenoServiceServer("unix:///run/core/deno.grpc.sock", handlers);
//
//   // TCP loopback fallback:
//   const srv = await startDenoServiceServer("0.0.0.0:50052", handlers);
export function startDenoServiceServer(
  address: string,
  handlers: DenoServiceHandlers,
): Promise<DenoServiceServer> {
  const proto = getProto();
  const server = new grpc.Server();

  server.addService(proto.DenoService.service, {
    OnStart: wrapUnary(async (req: LifecycleEventMsg) => {
      await handlers.onStart(req.phase ?? "", req.reason ?? "");
      return { ok: true };
    }),

    OnStop: wrapUnary(async (req: LifecycleEventMsg) => {
      await handlers.onStop(req.phase ?? "", req.reason ?? "");
      return { ok: true };
    }),

    OnConfigChange: wrapUnary(async (req: ConfigChangeEventMsg) => {
      await handlers.onConfigChange(req.key ?? "", req.value ?? "");
      return { ok: true };
    }),

    Render: wrapUnary(async (req: RenderRequestMsg) => {
      try {
        const html = await handlers.render(req.component ?? "", req.props ?? "");
        return { html, error: "" };
      } catch (err) {
        throw asGrpcError(err);
      }
    }),

    Eval: wrapUnary(async (req: EvalRequestMsg) => {
      try {
        const resultJson = await handlers.eval(req.expression ?? "");
        return { result_json: resultJson, error: "" };
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        return { result_json: "", error: message };
      }
    }),
  });

  return new Promise<DenoServiceServer>((resolve, reject) => {
    server.bindAsync(
      address,
      grpc.ServerCredentials.createInsecure(),
      (err: Error | null, _port: number) => {
        if (err) {
          reject(err);
          return;
        }
        resolve({
          close() {
            server.tryShutdown((_err) => {
              // Best-effort — ignore errors during shutdown.
            });
          },
        });
      },
    );
  });
}

// --- Internal helpers ---

interface LifecycleEventMsg {
  phase?: string;
  reason?: string;
}

interface ConfigChangeEventMsg {
  key?: string;
  value?: string;
}

interface RenderRequestMsg {
  component?: string;
  props?: string;
}

interface EvalRequestMsg {
  expression?: string;
}

// wrapUnary adapts an async handler to the grpc-js ServerUnaryCall callback style.
// deno-lint-ignore no-explicit-any
function wrapUnary<TReq, TRes>(handler: (req: TReq) => Promise<TRes>): any {
  // deno-lint-ignore no-explicit-any
  return (call: any, callback: (err: Error | null, res?: TRes) => void) => {
    handler(call.request as TReq).then(
      (res) => callback(null, res),
      (err) => callback(asGrpcError(err)),
    );
  };
}

function asGrpcError(err: unknown): Error {
  if (err instanceof Error) return err;
  return new Error(String(err));
}
