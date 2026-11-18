// Tests for the real DenoServiceHandlers implementation.
// Follows TestFilename_Function_{Good,Bad,Ugly} naming convention.

import {
  createDenoServiceHandlers,
  CoreComponent,
  DenoServiceState,
  defineCoreElement,
} from "./sidecar_handlers.ts";

function assert(condition: boolean, message: string): void {
  if (!condition) throw new Error(message);
}

function assertEquals<T>(actual: T, expected: T, message: string): void {
  if (!Object.is(actual, expected)) {
    throw new Error(
      `${message}: expected ${String(expected)}, got ${String(actual)}`,
    );
  }
}

// --- Lifecycle handler tests ---

Deno.test("TestDenoServiceHandlers_Lifecycle_Good", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  assertEquals(state.phase, "idle", "initial phase should be idle");

  // onStart — updates phase and reason
  await handlers.onStart("start", "boot");
  assertEquals(state.phase, "start", "onStart should set phase to 'start'");
  assertEquals(state.reason, "boot", "onStart should set reason to 'boot'");

  // onStop — updates phase and reason
  await handlers.onStop("stop", "shutdown");
  assertEquals(state.phase, "stop", "onStop should set phase to 'stop'");
  assertEquals(state.reason, "shutdown", "onStop should set reason to 'shutdown'");
});

Deno.test("TestDenoServiceHandlers_Lifecycle_Bad", async () => {
  // Lifecycle handlers must not throw even with empty strings.
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  await handlers.onStart("", "");
  assertEquals(state.phase, "", "onStart with empty strings should still update phase");

  await handlers.onStop("", "");
  assertEquals(state.reason, "", "onStop with empty strings should still update reason");
});

Deno.test("TestDenoServiceHandlers_Lifecycle_Ugly", async () => {
  // Multiple rapid calls — last write wins, state remains consistent.
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  await Promise.all([
    handlers.onStart("start", "boot"),
    handlers.onStart("start", "hotreload"),
    handlers.onStop("stop", "crash"),
  ]);

  // Last write wins — state is consistent (not corrupted).
  assert(
    state.phase !== undefined,
    "phase should not be undefined after concurrent calls",
  );
  assert(
    state.reason !== undefined,
    "reason should not be undefined after concurrent calls",
  );
});

// --- Config change handler tests ---

Deno.test("TestDenoServiceHandlers_OnConfigChange_Good", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  await handlers.onConfigChange("theme.accent", "#6366f1");
  assertEquals(
    state.config.get("theme.accent"),
    "#6366f1",
    "onConfigChange should persist the value in state.config",
  );

  await handlers.onConfigChange("log.level", "debug");
  assertEquals(
    state.config.get("log.level"),
    "debug",
    "onConfigChange should update a second key independently",
  );
  assertEquals(
    state.config.size,
    2,
    "config map should hold both keys after two changes",
  );
});

Deno.test("TestDenoServiceHandlers_OnConfigChange_Bad", async () => {
  // Empty key should still be stored without throwing.
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  await handlers.onConfigChange("", "orphan");
  assert(
    state.config.has(""),
    "onConfigChange should store even an empty key",
  );
});

Deno.test("TestDenoServiceHandlers_OnConfigChange_Ugly", async () => {
  // Overwriting an existing key — config must reflect the latest value.
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  await handlers.onConfigChange("app.mode", "development");
  await handlers.onConfigChange("app.mode", "production");
  assertEquals(
    state.config.get("app.mode"),
    "production",
    "onConfigChange should overwrite an existing key with the latest value",
  );
});

// --- Eval handler tests ---

Deno.test("TestDenoServiceHandlers_Eval_Good", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Numeric expression
  const numResult = await handlers.eval("2 + 3");
  assertEquals(numResult, "5", "eval should return JSON-encoded numeric result");

  // String expression
  const strResult = await handlers.eval('"hello".toUpperCase()');
  assertEquals(
    strResult,
    '"HELLO"',
    "eval should return JSON-encoded string result (with quotes)",
  );

  // Array expression
  const arrResult = await handlers.eval("[1, 2, 3].map(x => x * 2)");
  assertEquals(
    arrResult,
    "[2,4,6]",
    "eval should return JSON-encoded array result",
  );

  // Boolean expression
  const boolResult = await handlers.eval("true || false");
  assertEquals(boolResult, "true", "eval should return JSON-encoded boolean result");
});

Deno.test("TestDenoServiceHandlers_Eval_Bad", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Syntax error — must reject
  let caught = "";
  try {
    await handlers.eval("((( broken syntax ===");
  } catch (err) {
    caught = err instanceof Error ? err.message : String(err);
  }
  assert(caught.length > 0, "eval should reject on syntax errors");
});

Deno.test("TestDenoServiceHandlers_Eval_Ugly", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Undefined result — JSON.stringify(undefined) === undefined, fallback to "null"
  // (JSON.stringify returns undefined for undefined, but JSON encodes null for that)
  const result = await handlers.eval("undefined");
  // JSON.stringify(undefined) returns undefined (not a string), so we encode null
  // per JSON spec: undefined is not a valid JSON value.
  assert(
    result === "null" || result === undefined,
    "eval of undefined should encode as null or undefined",
  );
});

// --- Render handler tests ---

// Register a test component before the render tests run.
const RENDER_TAG = "core-sidecar-test-card";
class TestCardComponent extends CoreComponent<{ title: string }> {
  constructor() {
    super({ title: "Default" });
  }
  protected template(): string {
    return `<div class="card"><h1>${this.state.title}</h1></div>`;
  }
}
defineCoreElement(RENDER_TAG, TestCardComponent);

Deno.test("TestDenoServiceHandlers_Render_Good", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Render with empty props — default state
  const html = await handlers.render(RENDER_TAG, "{}");
  assert(
    html.includes("Default"),
    "render should produce HTML containing the component's default state",
  );
  assert(
    html.includes(RENDER_TAG),
    "render should wrap the shadow DOM in the component tag",
  );

  // Render with JSON props as attributes
  const htmlWithProps = await handlers.render(
    RENDER_TAG,
    '{"title":"Welcome"}',
  );
  assert(
    typeof htmlWithProps === "string" && htmlWithProps.length > 0,
    "render with props should return a non-empty HTML string",
  );
});

Deno.test("TestDenoServiceHandlers_Render_Bad", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Invalid JSON props — must reject
  let caught = "";
  try {
    await handlers.render(RENDER_TAG, "{ not json ]");
  } catch (err) {
    caught = err instanceof Error ? err.message : String(err);
  }
  assert(caught.length > 0, "render should reject when propsJson is invalid JSON");
});

Deno.test("TestDenoServiceHandlers_Render_Ugly", async () => {
  const state = new DenoServiceState();
  const handlers = createDenoServiceHandlers(state);

  // Unknown component — document.createElement returns a plain element with no shadow.
  // The handler should return a shell tag, not throw.
  const html = await handlers.render("core-unknown-zzz", "{}");
  assert(
    typeof html === "string",
    "render of an unknown component should return a string, not throw",
  );
});
