// Real DenoServiceHandlers — lifecycle, config, component SSR, expression eval.
//
// Usage (wire into main.ts or any entrypoint):
//   import { createDenoServiceHandlers, type DenoServiceState } from "./sidecar_handlers.ts";
//   const state = new DenoServiceState();
//   const handlers = createDenoServiceHandlers(state);
//   const grpcServer = await startDenoServiceServer(address, handlers);
//
// Render uses CoreTS's existing component system (defineCoreElement / CoreComponent).
// Components must be registered via defineCoreElement() before Render is called;
// registration is the app layer's responsibility (same as in a browser context).
//
// Eval uses the Function constructor — synchronous TS/JS expression evaluation
// inside the Deno runtime. The result is JSON-encoded and returned in result_json.
//
// Lifecycle handlers update the shared DenoServiceState. OnConfigChange applies
// key/value pairs to the mutable config map on the state object so that any
// code holding a reference to the state will see live updates.

import type { DenoServiceHandlers } from "./sidecar_server.ts";
import { defineCoreElement, CoreComponent } from "../src/components.ts";

export { defineCoreElement, CoreComponent };

// DenoServiceState — shared runtime + config state for the sidecar.
//
//   const state = new DenoServiceState();
//   state.config.get("theme.accent")  // live view after OnConfigChange
//   state.phase                        // "start" | "stop" | "idle"
export class DenoServiceState {
  // phase — last lifecycle phase received from Go ("idle" until OnStart fires).
  phase: string = "idle";
  // reason — last lifecycle reason received from Go.
  reason: string = "";
  // config — live key/value config map. OnConfigChange updates it in place.
  readonly config: Map<string, string> = new Map();
}

// renderComponentToHTML — SSR a registered CoreComponent to an HTML string.
//
//   defineCoreElement("my-card", MyCard);
//   const html = renderComponentToHTML("my-card", { title: "Hello" });
//   // => '<my-card ...>...</my-card>' wrapping the shadow DOM HTML
function renderComponentToHTML(component: string, props: Record<string, unknown>): string {
  // createElement goes through the fallback document path installed by components.ts
  // when running outside a browser, which invokes the registered ctor.
  const el = document.createElement(component) as HTMLElement & {
    shadowRoot: { innerHTML: string } | null;
    setState?: (patch: Record<string, unknown>) => void;
    connectedCallback?: () => void;
  };

  // Inject props as attributes and via setState if the component supports it.
  for (const [key, value] of Object.entries(props)) {
    const strValue = typeof value === "string" ? value : JSON.stringify(value);
    el.setAttribute(key, strValue);
  }

  // Trigger connectedCallback so the component renders (the fallback document's
  // appendChild does this via queueMicrotask, but we call it directly for sync SSR).
  if (typeof el.connectedCallback === "function") {
    el.connectedCallback();
  }

  const shadow = el.shadowRoot;
  if (!shadow) {
    // Component rendered without a shadow root — return an empty shell.
    return `<${component}></${component}>`;
  }

  return `<${component}>${shadow.innerHTML}</${component}>`;
}

// evalExpression — evaluate a JS/TS expression and return the JSON-encoded result.
//
//   const json = await evalExpression("1 + 1")          // => "2"
//   const json = await evalExpression("Math.PI.toFixed(4)") // => '"3.1416"'
//
// Uses the Function constructor (indirect eval) so the expression runs in the
// global scope, not in the caller's module scope — same security boundary as
// a script tag. Async expressions are awaited.
async function evalExpression(expression: string): Promise<string> {
  // deno-lint-ignore no-eval
  const fn = new Function(`return (async () => { return (${expression}); })()`);
  const result = await (fn() as Promise<unknown>);
  return JSON.stringify(result);
}

// createDenoServiceHandlers — build DenoServiceHandlers backed by real implementations.
//
//   const state = new DenoServiceState();
//   const handlers = createDenoServiceHandlers(state);
//   await startDenoServiceServer("unix:///run/core/deno.grpc.sock", handlers);
export function createDenoServiceHandlers(
  state: DenoServiceState,
): DenoServiceHandlers {
  return {
    // onStart — update runtime phase and log. Go calls this when the application boots.
    //   state.phase === "start" after this fires
    onStart(phase: string, reason: string): void {
      state.phase = phase;
      state.reason = reason;
      console.error(`CoreDeno/gRPC: onStart phase=${phase} reason=${reason}`);
    },

    // onStop — update runtime phase and log. Go calls this on graceful shutdown.
    //   state.phase === "stop" after this fires
    onStop(phase: string, reason: string): void {
      state.phase = phase;
      state.reason = reason;
      console.error(`CoreDeno/gRPC: onStop phase=${phase} reason=${reason}`);
    },

    // onConfigChange — apply a single key/value pair to the live config map.
    //   state.config.get("theme.accent") reflects the new value immediately
    onConfigChange(key: string, value: string): void {
      state.config.set(key, value);
      console.error(`CoreDeno/gRPC: onConfigChange ${key}=${value}`);
    },

    // render — server-side render a registered CoreComponent to HTML.
    //
    //   defineCoreElement("core-hero", HeroComponent);  // app layer registers
    //   const html = await handlers.render("core-hero", '{"title":"Welcome"}');
    //   // => "<core-hero><h1>Welcome</h1></core-hero>"
    //
    // Returns an error string (not a rejection) for unknown components, consistent
    // with the RenderResponse proto contract.
    render(component: string, propsJson: string): Promise<string> {
      let props: Record<string, unknown> = {};
      if (propsJson && propsJson !== "{}") {
        try {
          props = JSON.parse(propsJson) as Record<string, unknown>;
        } catch {
          return Promise.reject(
            new Error(`render: invalid props JSON for component "${component}"`),
          );
        }
      }
      return Promise.resolve(renderComponentToHTML(component, props));
    },

    // eval — evaluate a TS/JS expression and return the JSON-encoded result.
    //
    //   const json = await handlers.eval("2 + 3")      // => "5"
    //   const json = await handlers.eval('"hello".toUpperCase()') // => '"HELLO"'
    eval(expression: string): Promise<string> {
      return evalExpression(expression);
    },
  };
}
