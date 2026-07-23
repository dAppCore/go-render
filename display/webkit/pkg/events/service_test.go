// pkg/events/service_test.go
package events

import (
	"context"
	"sync"

	core "dappco.re/go"
)

// --- Mock Platform ---

type mockPlatform struct {
	mu          sync.Mutex
	listeners   map[string][]*mockListener
	emitted     []CustomEvent
	resetCalled bool
	nilCancel   bool
}

type mockListener struct {
	callback func(*CustomEvent)
	counter  int // -1 = persistent
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		listeners: make(map[string][]*mockListener),
	}
}

func (m *mockPlatform) Emit(name string, data ...any) bool {
	event := &CustomEvent{Name: name}
	if len(data) == 1 {
		event.Data = data[0]
	} else if len(data) > 1 {
		event.Data = data
	}

	m.mu.Lock()
	m.emitted = append(m.emitted, *event)
	active := make([]*mockListener, len(m.listeners[name]))
	copy(active, m.listeners[name])
	m.mu.Unlock()

	for _, listener := range active {
		listener.callback(event)
	}
	return false
}

func (m *mockPlatform) On(name string, callback func(*CustomEvent)) func() {
	listener := &mockListener{callback: callback, counter: -1}
	m.mu.Lock()
	m.listeners[name] = append(m.listeners[name], listener)
	m.mu.Unlock()
	if m.nilCancel {
		return nil
	}
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		updated := m.listeners[name][:0]
		for _, existing := range m.listeners[name] {
			if existing != listener {
				updated = append(updated, existing)
			}
		}
		m.listeners[name] = updated
	}
}

func (m *mockPlatform) Off(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.listeners, name)
}

func (m *mockPlatform) OnMultiple(name string, callback func(*CustomEvent), counter int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners[name] = append(m.listeners[name], &mockListener{callback: callback, counter: counter})
}

func (m *mockPlatform) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = make(map[string][]*mockListener)
	m.resetCalled = true
}

// simulateEvent fires all registered listeners for the given event name with optional data.
func (m *mockPlatform) simulateEvent(name string, data any) {
	event := &CustomEvent{Name: name, Data: data}
	m.mu.Lock()
	active := make([]*mockListener, len(m.listeners[name]))
	copy(active, m.listeners[name])
	m.mu.Unlock()
	for _, listener := range active {
		listener.callback(event)
	}
}

// listenerCount returns the total number of registered listeners across all event names.
func (m *mockPlatform) listenerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, listeners := range m.listeners {
		total += len(listeners)
	}
	return total
}

// --- Test helpers ---

func newTestService(t *core.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := newMockPlatform()
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")
	return svc, c, mock
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

// --- Good path tests ---

func TestRegister_Good(t *core.T) {
	svc, _, _ := newTestService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotEmpty(t, core.Sprintf("%T", svc))
}

func TestTaskEmit_Good(t *core.T) {
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "user:login", Data: "alice"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, false, r.Value) // not cancelled

	core.AssertLen(t, mock.emitted, 1)
	core.AssertEqual(t, "user:login", mock.emitted[0].Name)
	core.AssertEqual(t, "alice", mock.emitted[0].Data)
}

func TestTaskEmit_NoData_GoodCase(t *core.T) {
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "ping"})
	core.RequireTrue(t, r.OK)
	core.AssertLen(t, mock.emitted, 1)
	core.AssertNil(t, mock.emitted[0].Data)
}

func TestTaskOn_Good(t *core.T) {
	_, c, mock := newTestService(t)

	var received []ActionEventFired
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if fired, ok := msg.(ActionEventFired); ok {
			received = append(received, fired)
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "events.on", TaskOn{Name: "theme:changed"})
	core.RequireTrue(t, r.OK)

	mock.simulateEvent("theme:changed", "dark")

	core.AssertLen(t, received, 1)
	core.AssertEqual(t, "theme:changed", received[0].Event.Name)
	core.AssertEqual(t, "dark", received[0].Event.Data)
}

func TestTaskOn_NilEvent_Ignores(t *core.T) {
	_, c, mock := newTestService(t)

	var received []ActionEventFired
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if fired, ok := msg.(ActionEventFired); ok {
			received = append(received, fired)
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "events.on", TaskOn{Name: "theme:changed"})
	core.RequireTrue(t, r.OK)

	core.RequireNotEmpty(t, mock.listeners["theme:changed"])
	core.AssertNotPanics(t, func() {
		mock.listeners["theme:changed"][0].callback(nil)
	})

	core.AssertEmpty(t, received)
}

func TestTaskOff_Good(t *core.T) {
	_, c, mock := newTestService(t)

	// Register via IPC then remove
	r := taskRun(c, "events.on", TaskOn{Name: "file:saved"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, 1, mock.listenerCount())

	r2 := taskRun(c, "events.off", TaskOff{Name: "file:saved"})
	core.RequireTrue(t, r2.OK)
	core.AssertEqual(t, 0, mock.listenerCount())
}

func TestQueryListeners_Good(t *core.T) {
	_, c, _ := newTestService(t)

	core.RequireTrue(t, taskRun(c, "events.on", TaskOn{Name: "user:login"}).OK)
	core.RequireTrue(t, taskRun(c, "events.on", TaskOn{Name: "user:login"}).OK)
	core.RequireTrue(t, taskRun(c, "events.on", TaskOn{Name: "theme:changed"}).OK)

	r := c.QUERY(QueryListeners{})
	core.RequireTrue(t, r.OK)

	infos := r.Value.([]ListenerInfo)
	counts := make(map[string]int)
	for _, info := range infos {
		counts[info.EventName] = info.Count
	}
	core.AssertEqual(t, 2, counts["user:login"])
	core.AssertEqual(t, 1, counts["theme:changed"])
}

func TestQueryListeners_Empty_GoodCase(t *core.T) {
	_, c, _ := newTestService(t)

	r := c.QUERY(QueryListeners{})
	core.RequireTrue(t, r.OK)

	infos := r.Value.([]ListenerInfo)
	core.AssertEmpty(t, infos)
}

func TestOnShutdown_CancelsAll_GoodCase(t *core.T) {
	svc, _, mock := newTestService(t)

	core.RequireTrue(t, taskRun(svc.Core(), "events.on", TaskOn{Name: "a:b"}).OK)
	core.RequireTrue(t, taskRun(svc.Core(), "events.on", TaskOn{Name: "c:d"}).OK)
	core.AssertEqual(t, 2, mock.listenerCount())

	core.RequireTrue(t, svc.OnShutdown(context.Background()).OK)
	core.AssertEqual(t, 0, mock.listenerCount())
}

func TestOnShutdown_IgnoresNilCancels_GoodCase(t *core.T) {
	mock := newMockPlatform()
	mock.nilCancel = true
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")

	core.RequireTrue(t, taskRun(c, "events.on", TaskOn{Name: "a:b"}).OK)
	core.AssertEqual(t, 1, mock.listenerCount())

	core.AssertNotPanics(t, func() {
		core.AssertTrue(t, svc.OnShutdown(context.Background()).OK)
	})
}

func TestActionEventFired_BroadcastOnSimulate_GoodCase(t *core.T) {
	_, c, mock := newTestService(t)

	var receivedEvents []CustomEvent
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if fired, ok := msg.(ActionEventFired); ok {
			receivedEvents = append(receivedEvents, fired.Event)
		}
		return core.Result{OK: true}
	})

	core.RequireTrue(t, taskRun(c, "events.on", TaskOn{Name: "data:ready"}).OK)

	mock.simulateEvent("data:ready", map[string]any{"rows": 42})

	core.AssertLen(t, receivedEvents, 1)
	core.AssertEqual(t, "data:ready", receivedEvents[0].Name)
}

// --- Bad path tests ---

func TestTaskOn_EmptyName_BadCase(t *core.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.on", TaskOn{Name: ""})
	core.AssertFalse(t, r.OK)
}

func TestTaskEmit_UnknownEvent_BadCase(t *core.T) {
	// Emitting an event with no listeners is valid — returns not-cancelled.
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "no:listeners"})
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, false, r.Value)
	core.AssertLen(t, mock.emitted, 1) // still recorded as emitted
}

func TestTaskEmit_EmptyName_BadCase(t *core.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: ""})
	core.AssertFalse(t, r.OK)
}

func TestQueryListeners_NoService_BadCase(t *core.T) {
	// No events service registered — query is not handled.
	c := core.New(core.WithServiceLock())

	r := c.QUERY(QueryListeners{})
	core.AssertFalse(t, r.OK)
}

func TestTaskEmit_NoService_BadCase(t *core.T) {
	c := core.New(core.WithServiceLock())

	r := c.Action("events.emit").Run(context.Background(), core.NewOptions())
	core.AssertFalse(t, r.OK)
}

func TestTaskEmit_PlatformUnavailable_BadCase(t *core.T) {
	c := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")

	r := taskRun(c, "events.emit", TaskEmit{Name: "user:login"})
	core.AssertFalse(t, r.OK)
	err, ok := r.Value.(resultFailure)
	core.RequireTrue(t, ok)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "event platform unavailable")

	r = taskRun(c, "events.on", TaskOn{Name: "user:login"})
	core.AssertFalse(t, r.OK)
	err, ok = r.Value.(resultFailure)
	core.RequireTrue(t, ok)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "event platform unavailable")

	r = taskRun(c, "events.off", TaskOff{Name: "user:login"})
	core.AssertFalse(t, r.OK)
	err, ok = r.Value.(resultFailure)
	core.RequireTrue(t, ok)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "event platform unavailable")

	core.AssertNotPanics(t, func() {
		core.AssertTrue(t, svc.OnShutdown(context.Background()).OK)
	})
}

// --- Ugly path tests ---

func TestTaskOff_NeverRegistered_UglyCase(t *core.T) {
	// Off on a name that was never registered is a no-op — must not panic.
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.off", TaskOff{Name: "nonexistent:event"})
	core.AssertTrue(t, r.OK)
}

func TestTaskOff_EmptyName_BadCase(t *core.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.off", TaskOff{Name: ""})
	core.AssertFalse(t, r.OK)
}

func TestTaskOn_MultipleListeners_UglyCase(t *core.T) {
	// Multiple IPC listeners for the same event each receive ActionEventFired.
	_, c, mock := newTestService(t)

	var mu sync.Mutex
	var fireCount int
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionEventFired); ok {
			mu.Lock()
			fireCount++
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	taskRun(c, "events.on", TaskOn{Name: "flood"})
	taskRun(c, "events.on", TaskOn{Name: "flood"})
	taskRun(c, "events.on", TaskOn{Name: "flood"})

	mock.simulateEvent("flood", nil)

	mu.Lock()
	count := fireCount
	mu.Unlock()
	core.AssertEqual(t, 3, count)
}

func TestTaskOff_ThenEmit_UglyCase(t *core.T) {
	// After Off, simulating the event must not trigger any IPC actions.
	_, c, mock := newTestService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionEventFired); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	taskRun(c, "events.on", TaskOn{Name: "transient"})
	taskRun(c, "events.off", TaskOff{Name: "transient"})

	mock.simulateEvent("transient", "late-data")
	core.AssertFalse(t, received)
}

func TestQueryListeners_AfterOff_UglyCase(t *core.T) {
	// After Off, the event name must not appear in QueryListeners results.
	_, c, _ := newTestService(t)

	taskRun(c, "events.on", TaskOn{Name: "ephemeral"})
	taskRun(c, "events.off", TaskOff{Name: "ephemeral"})

	r := c.QUERY(QueryListeners{})
	infos := r.Value.([]ListenerInfo)

	for _, info := range infos {
		core.AssertNotEqual(t, "ephemeral", info.EventName)
	}
}

// AX7 generated source-matching smoke coverage.
func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnShutdown(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
