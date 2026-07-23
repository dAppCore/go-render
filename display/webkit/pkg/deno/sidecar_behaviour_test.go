package deno

import (
	"bufio"
	"bytes"
	"io"

	core "dappco.re/go"
)

// drainStdin attaches an in-memory pipe as the manager's stdin and returns a
// channel yielding each newline-framed RPC message the manager writes.
//
//	m, lines := managerWithStdin(t)
//	m.Emit("ready", nil)
//	core.AssertContains(t, <-lines, `"type":"event"`)
func managerWithStdin(t *core.T, opts Options) (*Manager, <-chan string) {
	t.Helper()
	m := New(opts)
	reader, writer := io.Pipe()
	m.stdin = writer

	lines := make(chan string, 8)
	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()
	return m, lines
}

// isProcessDone recognises the family of "already terminated" errors so Stop
// treats a dead sidecar as a clean shutdown rather than a failure.
//
//	isProcessDone(core.NewError("os: process already finished")) // true
func TestSidecarBehaviour_isProcessDone_Good(t *core.T) {
	core.AssertTrue(t, isProcessDone(core.NewError("os: process already finished")))
	core.AssertTrue(t, isProcessDone(core.NewError("waitid: no child processes")))
	core.AssertTrue(t, isProcessDone(core.NewError("PROCESS ALREADY FINISHED")))
}

// isProcessDone reports false for unrelated errors so genuine failures surface.
func TestSidecarBehaviour_isProcessDone_Bad(t *core.T) {
	core.AssertFalse(t, isProcessDone(core.NewError("permission denied")))
}

// isProcessDone treats a nil error as not-done.
func TestSidecarBehaviour_isProcessDone_Ugly(t *core.T) {
	core.AssertFalse(t, isProcessDone(nil))
}

// marshalRPCMessage serialises a message to JSON bytes.
//
//	payload, _ := marshalRPCMessage(rpcMessage{Type: "eval", Code: "1+1"})
func TestSidecarBehaviour_marshalRPCMessage_Good(t *core.T) {
	payload, err := marshalRPCMessage(rpcMessage{Type: "eval", ID: "deno-1", Code: "1+1"})
	core.AssertNil(t, err)
	core.AssertContains(t, string(payload), `"type":"eval"`)
	core.AssertContains(t, string(payload), `"code":"1+1"`)
}

// marshalRPCMessage round-trips back into an equivalent message.
func TestSidecarBehaviour_marshalRPCMessage_Bad(t *core.T) {
	payload, err := marshalRPCMessage(rpcMessage{Type: "event", Name: "ready"})
	core.AssertNil(t, err)
	var back rpcMessage
	core.AssertTrue(t, core.JSONUnmarshal(payload, &back).OK)
	core.AssertEqual(t, "event", back.Type)
	core.AssertEqual(t, "ready", back.Name)
}

// marshalRPCMessage omits empty optional fields (omitempty tags hold).
func TestSidecarBehaviour_marshalRPCMessage_Ugly(t *core.T) {
	payload, err := marshalRPCMessage(rpcMessage{Type: "result"})
	core.AssertNil(t, err)
	core.AssertNotContains(t, string(payload), `"id"`)
	core.AssertNotContains(t, string(payload), `"name"`)
}

// handleMessage of type "result" delivers to the matching pending channel.
//
//	m.pending["deno-7"] = ch
//	m.handleMessage(rpcMessage{Type: "result", ID: "deno-7", OK: true})
func TestSidecarBehaviour_handleMessage_Good(t *core.T) {
	m := New(Options{})
	ch := make(chan rpcMessage, 1)
	m.pending["deno-7"] = ch

	m.handleMessage(rpcMessage{Type: "result", ID: "deno-7", OK: true, Result: 42})

	got := <-ch
	core.AssertTrue(t, got.OK)
	core.AssertEqual(t, 42, got.Result)
	// The pending entry is consumed.
	_, still := m.pending["deno-7"]
	core.AssertFalse(t, still)
}

// handleMessage of type "event" fans out to every registered handler.
func TestSidecarBehaviour_handleMessage_Bad(t *core.T) {
	m := New(Options{})
	var seen Event
	m.OnEvent(func(e Event) { seen = e })

	m.handleMessage(rpcMessage{Type: "event", Name: "ready", Data: "now"})

	core.AssertEqual(t, "ready", seen.Name)
	core.AssertEqual(t, "now", seen.Data)
}

// handleMessage ignores unknown types and a result for an unknown ID.
func TestSidecarBehaviour_handleMessage_Ugly(t *core.T) {
	m := New(Options{})
	core.AssertNotPanics(t, func() {
		m.handleMessage(rpcMessage{Type: "unrecognised"})
		m.handleMessage(rpcMessage{Type: "result", ID: "missing", OK: true})
	})
}

// handleAction runs the named Core action and writes an OK result back to stdin.
//
//	c.Action("gui.ping", func(...) core.Result { return core.Ok("pong") })
//	m.handleAction(rpcMessage{Type: "action", ID: "deno-1", Name: "gui.ping"})
func TestSidecarBehaviour_handleAction_Good(t *core.T) {
	c := core.New()
	c.Action("gui.ping", func(_ core.Context, _ core.Options) core.Result {
		return core.Ok("pong")
	})
	m, lines := managerWithStdin(t, Options{Core: c})

	go m.handleAction(rpcMessage{Type: "action", ID: "deno-1", Name: "gui.ping"})

	line := <-lines
	core.AssertContains(t, line, `"type":"result"`)
	core.AssertContains(t, line, `"id":"deno-1"`)
	core.AssertContains(t, line, `"ok":true`)
	core.AssertContains(t, line, "pong")
}

// handleAction reports the action error string when the action fails.
func TestSidecarBehaviour_handleAction_Bad(t *core.T) {
	c := core.New()
	c.Action("gui.boom", func(_ core.Context, _ core.Options) core.Result {
		return core.Fail(core.NewError("kaboom"))
	})
	m, lines := managerWithStdin(t, Options{Core: c})

	go m.handleAction(rpcMessage{Type: "action", ID: "deno-2", Name: "gui.boom"})

	line := <-lines
	core.AssertContains(t, line, `"id":"deno-2"`)
	core.AssertContains(t, line, "kaboom")
	core.AssertNotContains(t, line, `"ok":true`)
}

// handleAction reports "core is unavailable" when no Core is wired.
func TestSidecarBehaviour_handleAction_Ugly(t *core.T) {
	m, lines := managerWithStdin(t, Options{})

	go m.handleAction(rpcMessage{Type: "action", ID: "deno-3", Name: "gui.ping"})

	line := <-lines
	core.AssertContains(t, line, "core is unavailable")
}

// readLoop frames newline-delimited JSON from stdout and dispatches each
// message; an "event" line fans out through the registered handler and a
// malformed line is skipped without halting the loop.
//
//	go m.readLoop(stdout, done)
func TestSidecarBehaviour_readLoop_Good(t *core.T) {
	m := New(Options{})
	events := make(chan Event, 1)
	m.OnEvent(func(e Event) { events <- e })

	stdout := bytes.NewBufferString(
		"not-json\n" +
			`{"type":"event","name":"ready","data":"go"}` + "\n",
	)
	done := make(chan struct{})
	go m.readLoop(stdout, done)

	got := <-events
	core.AssertEqual(t, "ready", got.Name)
	core.AssertEqual(t, "go", got.Data)
	// readLoop closes done when the reader is exhausted.
	<-done
}

// readLoop exits cleanly on an empty stream, closing the done channel.
func TestSidecarBehaviour_readLoop_Bad(t *core.T) {
	m := New(Options{})
	done := make(chan struct{})
	go m.readLoop(bytes.NewBufferString(""), done)
	<-done
}

// send fails fast when the sidecar is not running (no stdin attached).
func TestSidecarBehaviour_send_NotRunning(t *core.T) {
	m := New(Options{})
	err := m.send(rpcMessage{Type: "event", Name: "x"})
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "not running")
}
