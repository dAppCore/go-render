package deno

import (
	"bufio" // AX-6-exception: streaming stdout framing for the long-lived Deno sidecar.
	"context"
	"io"      // AX-6-exception: stdin/stdout pipe interfaces from exec.Cmd.
	"sync"    // AX-6-exception: manager protects live process and pending RPC state across goroutines.
	"syscall" // AX-6-exception: graceful sidecar shutdown sends SIGTERM.

	core "dappco.re/go"
)

type Options struct {
	Binary string
	Args   []string
	Dir    string
	Env    []string
	Core   *core.Core
}

type Status struct {
	Running   bool   `json:"running"`
	Connected bool   `json:"connected"`
	PID       int    `json:"pid,omitempty"`
	Binary    string `json:"binary,omitempty"`
}

type EvalResult struct {
	Value any `json:"value,omitempty"`
}

type Event struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}

type Manager struct {
	options  Options
	mu       sync.Mutex
	cmd      *core.Cmd
	stdin    io.WriteCloser
	ctx      context.Context
	cancel   context.CancelFunc
	readDone chan struct{}
	pending  map[string]chan rpcMessage
	events   []func(Event)
}

func New(options Options) *Manager {
	if core.Trim(options.Binary) == "" {
		options.Binary = "deno"
	}
	if len(options.Args) == 0 {
		options.Args = []string{"eval", denoBridgeProgram}
	}
	return &Manager{
		options: options,
		pending: make(map[string]chan rpcMessage),
	}
}

func (m *Manager) Start(ctx context.Context) (Status, resultFailure) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd != nil && m.cmd.Process != nil {
		return m.statusLocked(), nil
	}

	sidecarCtx, cancel := context.WithCancel(ctx)
	cmd := commandContext(sidecarCtx, m.options.Binary, m.options.Args...)
	cmd.Dir = m.options.Dir
	cmd.Env = append(core.Environ(), m.options.Env...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return Status{}, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return Status{}, err
	}
	if err := cmd.Start(); err != nil {
		cancel()
		return Status{}, err
	}

	m.cmd = cmd
	m.stdin = stdin
	m.ctx = sidecarCtx
	m.cancel = cancel
	m.readDone = make(chan struct{})
	go m.readLoop(stdout, m.readDone)
	return m.statusLocked(), nil
}

func (m *Manager) Stop(context.Context) (Status, resultFailure) {
	m.mu.Lock()
	if m.cmd == nil || m.cmd.Process == nil {
		m.mu.Unlock()
		return Status{}, nil
	}
	cmd := m.cmd
	done := m.readDone
	cancel := m.cancel
	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil && !isProcessDone(err) {
		status := m.statusLocked()
		m.mu.Unlock()
		return status, err
	}
	if cancel != nil {
		cancel()
	}
	m.mu.Unlock()

	if waitErr := cmd.Wait(); waitErr != nil && m.options.Core != nil {
		if warn := m.options.Core.LogWarn(waitErr, "deno.stop", "sidecar wait returned after termination"); !warn.OK {
			err = waitErr
		}
	}
	if done != nil {
		<-done
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd != cmd {
		return m.statusLocked(), err
	}
	m.cmd = nil
	m.stdin = nil
	m.ctx = nil
	m.cancel = nil
	m.readDone = nil
	for id, ch := range m.pending {
		close(ch)
		delete(m.pending, id)
	}
	return Status{}, nil
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.statusLocked()
}

func (m *Manager) statusLocked() Status {
	status := Status{
		Binary: m.options.Binary,
	}
	if m.cmd != nil && m.cmd.Process != nil {
		status.Running = true
		status.Connected = m.stdin != nil
		status.PID = m.cmd.Process.Pid
	}
	return status
}

func (m *Manager) OnEvent(handler func(Event)) {
	if handler == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, handler)
}

func (m *Manager) Eval(ctx context.Context, code string) (EvalResult, resultFailure) {
	response, err := m.request(ctx, rpcMessage{Type: "eval", Code: code})
	if err != nil {
		return EvalResult{}, err
	}
	return EvalResult{Value: response.Result}, nil
}

func (m *Manager) Emit(name string, data any) resultFailure {
	if core.Trim(name) == "" {
		return core.NewError("event name is required")
	}
	return m.send(rpcMessage{Type: "event", Name: name, Data: data})
}

func (m *Manager) request(ctx context.Context, message rpcMessage) (rpcMessage, resultFailure) {
	m.mu.Lock()
	if m.stdin == nil {
		m.mu.Unlock()
		return rpcMessage{}, core.NewError("deno sidecar is not running")
	}
	message.ID = core.Concat("deno-", core.ID())
	responseCh := make(chan rpcMessage, 1)
	m.pending[message.ID] = responseCh
	payload, err := marshalRPCMessage(message)
	if err != nil {
		delete(m.pending, message.ID)
		m.mu.Unlock()
		return rpcMessage{}, err
	}
	_, err = m.stdin.Write(append(payload, '\n'))
	m.mu.Unlock()
	if err != nil {
		return rpcMessage{}, err
	}
	select {
	case <-ctx.Done():
		m.mu.Lock()
		delete(m.pending, message.ID)
		m.mu.Unlock()
		return rpcMessage{}, ctx.Err()
	case response, ok := <-responseCh:
		if !ok {
			return rpcMessage{}, core.NewError("deno sidecar disconnected")
		}
		if !response.OK {
			return rpcMessage{}, core.NewError(core.Trim(response.Error))
		}
		return response, nil
	}
}

func (m *Manager) readLoop(stdout io.Reader, done chan struct{}) {
	defer close(done)
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var message rpcMessage
		if result := core.JSONUnmarshal(scanner.Bytes(), &message); !result.OK {
			continue
		}
		m.handleMessage(message)
	}
}

func (m *Manager) handleMessage(message rpcMessage) {
	switch message.Type {
	case "result":
		m.mu.Lock()
		ch := m.pending[message.ID]
		delete(m.pending, message.ID)
		m.mu.Unlock()
		if ch != nil {
			ch <- message
			close(ch)
		}
	case "event":
		m.mu.Lock()
		handlers := append([]func(Event){}, m.events...)
		m.mu.Unlock()
		for _, handler := range handlers {
			handler(Event{Name: message.Name, Data: message.Data})
		}
	case "action":
		m.handleAction(message)
	}
}

func (m *Manager) handleAction(message rpcMessage) {
	response := rpcMessage{Type: "result", ID: message.ID}
	if m.options.Core == nil {
		response.Error = "core is unavailable"
		if err := m.send(response); err != nil {
			return
		}
		return
	}
	opts := core.NewOptions()
	for key, value := range message.Options {
		opts.Set(key, value)
	}
	m.mu.Lock()
	ctx := m.ctx
	m.mu.Unlock()
	if ctx == nil {
		ctx = context.Background()
	}
	result := m.options.Core.Action(message.Name).Run(ctx, opts)
	response.OK = result.OK
	if result.OK {
		response.Result = result.Value
	} else if err, ok := result.Value.(error); ok {
		response.Error = err.Error()
	} else {
		response.Error = core.Sprint(result.Value)
	}
	if err := m.send(response); err != nil && m.options.Core != nil {
		if warn := m.options.Core.LogWarn(err, "deno.handleAction", "failed to send action response"); !warn.OK {
			return
		}
	}
}

func (m *Manager) send(message rpcMessage) resultFailure {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stdin == nil {
		return core.NewError("deno sidecar is not running")
	}
	payload, err := marshalRPCMessage(message)
	if err != nil {
		return err
	}
	_, err = m.stdin.Write(append(payload, '\n'))
	return err
}

func marshalRPCMessage(message rpcMessage) ([]byte, resultFailure) {
	result := core.JSONMarshal(message)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, err
		}
		return nil, core.NewError("failed to marshal deno sidecar message")
	}
	payload, ok := result.Value.([]byte)
	if !ok {
		return nil, core.NewError("failed to marshal deno sidecar message")
	}
	return payload, nil
}

type rpcMessage struct {
	Type    string         `json:"type"`
	ID      string         `json:"id,omitempty"`
	Name    string         `json:"name,omitempty"`
	Code    string         `json:"code,omitempty"`
	Data    any            `json:"data,omitempty"`
	Options map[string]any `json:"options,omitempty"`
	OK      bool           `json:"ok,omitempty"`
	Result  any            `json:"result,omitempty"`
	Error   string         `json:"error,omitempty"`
}

const denoBridgeProgram = `const encoder = new TextEncoder();
const decoder = new TextDecoder();
globalThis.core = {
  emit(name, data) {
    return send({ type: "event", name, data });
  },
  action(name, options = {}) {
    return request({ type: "action", name, options });
  },
};
const pending = new Map();
async function send(message) {
  await Deno.stdout.write(encoder.encode(JSON.stringify(message) + "\n"));
}
function request(message) {
  const id = crypto.randomUUID();
  return new Promise((resolve, reject) => {
    pending.set(id, { resolve, reject });
    send({ ...message, id }).catch((error) => {
      pending.delete(id);
      reject(error);
    });
  });
}
async function handle(message) {
  if (message.type === "eval") {
    try {
      const value = await (0, eval)(message.code);
      await send({ type: "result", id: message.id, ok: true, result: value });
    } catch (error) {
      await send({ type: "result", id: message.id, ok: false, error: String(error?.stack ?? error) });
    }
    return;
  }
  if (message.type === "event") {
    globalThis.dispatchEvent(new CustomEvent(message.name || "core.event", { detail: message.data ?? null }));
    await send({ type: "result", id: message.id, ok: true, result: null });
    return;
  }
  if (message.type === "result") {
    const pendingRequest = pending.get(message.id);
    if (!pendingRequest) return;
    pending.delete(message.id);
    if (message.ok) {
      pendingRequest.resolve(message.result);
    } else {
      pendingRequest.reject(new Error(message.error || "deno sidecar request failed"));
    }
  }
}
let buffer = "";
while (true) {
  const chunk = new Uint8Array(4096);
  const read = await Deno.stdin.read(chunk);
  if (read === null) break;
  buffer += decoder.decode(chunk.subarray(0, read));
  let newline = buffer.indexOf("\n");
  while (newline >= 0) {
    const line = buffer.slice(0, newline).trim();
    buffer = buffer.slice(newline + 1);
    if (line) {
      try {
        await handle(JSON.parse(line));
      } catch (error) {
        console.error("Failed to parse message:", line, error);
      }
    }
    newline = buffer.indexOf("\n");
  }
}`
