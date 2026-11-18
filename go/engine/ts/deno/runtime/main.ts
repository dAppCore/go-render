// CoreDeno Runtime Entry Point
// Connects to CoreGO via gRPC over Unix socket.
// Implements DenoService for module lifecycle management.

// Must be first import — patches http2 before @grpc/grpc-js loads.
import "./polyfill.ts";

import type { CoreClient } from "./client.ts";
import { CoreDevServer } from "./dev.ts";
import { CoreSidecar } from "./sidecar.ts";
import { type DenoServer, startDenoServer } from "./server.ts";
import { ModuleRegistry } from "./modules.ts";
import {
  loadSharedLocale,
  resolvePreferredLocale,
  setLocale,
  setLocaleBridge,
} from "../src/i18n.ts";
import { createGoServiceClient, type GoServiceClient } from "./sidecar_client.ts";
import { type DenoServiceServer, startDenoServiceServer } from "./sidecar_server.ts";
import { createDenoServiceHandlers, DenoServiceState } from "./sidecar_handlers.ts";

if (import.meta.main) {
  await bootstrapCoreDenoRuntime();
}

async function bootstrapCoreDenoRuntime(): Promise<void> {
  const coreSocket = Deno.env.get("CORE_SOCKET");
  if (!coreSocket) {
    console.error("FATAL: CORE_SOCKET environment variable not set");
    Deno.exit(1);
  }

  const denoSocket = Deno.env.get("DENO_SOCKET");
  if (!denoSocket) {
    console.error("FATAL: DENO_SOCKET environment variable not set");
    Deno.exit(1);
  }

  console.error(`CoreDeno: CORE_SOCKET=${coreSocket}`);
  console.error(`CoreDeno: DENO_SOCKET=${denoSocket}`);

  // 1. Create module registry
  const registry = new ModuleRegistry();

  // Optional dev server: watch source trees and trigger HMR reload hooks.
  let devServer: CoreDevServer | null = null;
  const devRoot = Deno.env.get("CORE_DEV_ROOT");
  if (devRoot) {
    try {
      devServer = new CoreDevServer({
        root: devRoot,
        hmrPath: Deno.env.get("CORE_HMR_PATH") ?? "/_core/hmr",
        onReload: () => {
          void registry.reloadAll().catch((err) => {
            console.error(`CoreDeno: HMR reload failed: ${err}`);
          });
        },
      });
      await devServer.start();
      console.error(`CoreDeno: dev server watching ${devRoot}`);
    } catch (err) {
      console.error(`CoreDeno: dev server unavailable: ${err}`);
    }
  }

  function stopDevServer(): void {
    try {
      devServer?.stop();
    } catch {
      // Best-effort cleanup during bootstrap/shutdown failures.
    }
  }

  let denoServer: DenoServer | null = null;
  const sidecar = new CoreSidecar({
    socketPath: coreSocket,
    logger: console,
  });

  // Optional gRPC sidecar bridge (CORE_GRPC_SOCKET + DENO_GRPC_SOCKET).
  // When both env vars are set the GoService client and DenoService gRPC server
  // start alongside the existing JSON-RPC path. If either var is absent the
  // bridge is silently skipped — existing deployments are unaffected.
  const coreGrpcSocket = Deno.env.get("CORE_GRPC_SOCKET");
  const denoGrpcSocket = Deno.env.get("DENO_GRPC_SOCKET");
  let grpcGoClient: GoServiceClient | null = null;
  let grpcDenoServer: DenoServiceServer | null = null;
  const grpcState = new DenoServiceState();

  function shutdownRuntime(): void {
    stopDevServer();
    sidecar.shutdown();
    try {
      denoServer?.close();
    } catch {
      // Best-effort cleanup during bootstrap/shutdown failures.
    }
    denoServer = null;
    try {
      grpcDenoServer?.close();
    } catch {
      // Best-effort cleanup during bootstrap/shutdown failures.
    }
    grpcDenoServer = null;
    try {
      grpcGoClient?.close();
    } catch {
      // Best-effort cleanup during bootstrap/shutdown failures.
    }
    grpcGoClient = null;
  }

  // 2. Start DenoService server (Go calls us here via JSON-RPC over Unix socket)
  try {
    denoServer = await startDenoServer(denoSocket, registry);
    console.error("CoreDeno: DenoService server started");
  } catch (err) {
    console.error(`FATAL: failed to start DenoService server: ${err}`);
    stopDevServer();
    Deno.exit(1);
  }

  // 2b. Start gRPC sidecar bridge (additive — only when both socket vars are set).
  if (coreGrpcSocket && denoGrpcSocket) {
    try {
      grpcGoClient = createGoServiceClient(coreGrpcSocket);
      console.error(`CoreDeno: GoService gRPC client → ${coreGrpcSocket}`);
    } catch (err) {
      console.error(`CoreDeno: GoService gRPC client unavailable: ${err}`);
    }
    try {
      const handlers = createDenoServiceHandlers(grpcState);
      grpcDenoServer = await startDenoServiceServer(denoGrpcSocket, handlers);
      console.error(`CoreDeno: DenoService gRPC server ← ${denoGrpcSocket}`);
    } catch (err) {
      console.error(`CoreDeno: DenoService gRPC server unavailable: ${err}`);
    }
  }

  let coreClient: CoreClient;
  try {
    // 3. Connect to CoreService (we call Go here) with retry
    coreClient = await sidecar.connect();
    console.error("CoreDeno: CoreService client connected");
    setLocaleBridge({
      localeGet(locale: string) {
        return coreClient.localeGet(locale);
      },
    });

    const preferredLocale = resolvePreferredLocale({
      CORE_LOCALE: Deno.env.get("CORE_LOCALE") ?? undefined,
      LANG: Deno.env.get("LANG") ?? undefined,
    });
    setLocale(preferredLocale);
    try {
      await loadSharedLocale(preferredLocale, { bridge: coreClient });
    } catch (err) {
      console.error(`CoreDeno: locale preload unavailable: ${err}`);
    }

    try {
      await sidecar.healthCheck();
    } catch (err) {
      console.error(`FATAL: failed CoreService health check: ${err}`);
      shutdownRuntime();
      Deno.exit(1);
    }
  } catch (err) {
    console.error(`FATAL: ${err}`);
    stopDevServer();
    denoServer?.close();
    Deno.exit(1);
  }

  // 4. Inject CoreClient into registry for I/O bridge
  registry.setCoreClient(coreClient);

  // 5. Signal readiness
  console.error("CoreDeno: ready");

  // 6. Keep alive until SIGTERM
  const shutdownController = new AbortController();
  Deno.addSignalListener("SIGTERM", () => {
    console.error("CoreDeno: shutting down");
    shutdownController.abort();
  });

  try {
    await new Promise((_resolve, reject) => {
      shutdownController.signal.addEventListener(
        "abort",
        () => reject(new Error("shutdown")),
      );
    });
  } catch {
    // Clean shutdown
    shutdownRuntime();
  }
}
