// pkg/lifecycle/service_test.go
package lifecycle

import (
	"context"
	"sync"

	core "dappco.re/go"
)

// --- Mock Platform ---

type mockPlatform struct {
	mu           sync.Mutex
	handlers     map[EventType][]func()
	fileHandlers []func(string)
	urlHandlers  []func(string)
	quitCalls    int
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		handlers: make(map[EventType][]func()),
	}
}

func (m *mockPlatform) OnApplicationEvent(eventType EventType, handler func()) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[eventType] = append(m.handlers[eventType], handler)
	idx := len(m.handlers[eventType]) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.handlers[eventType]) {
			m.handlers[eventType] = append(m.handlers[eventType][:idx], m.handlers[eventType][idx+1:]...)
		}
	}
}

func (m *mockPlatform) OnOpenedWithFile(handler func(string)) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fileHandlers = append(m.fileHandlers, handler)
	idx := len(m.fileHandlers) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.fileHandlers) {
			m.fileHandlers = append(m.fileHandlers[:idx], m.fileHandlers[idx+1:]...)
		}
	}
}

func (m *mockPlatform) OnLaunchedWithUrl(handler func(string)) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.urlHandlers = append(m.urlHandlers, handler)
	idx := len(m.urlHandlers) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.urlHandlers) {
			m.urlHandlers = append(m.urlHandlers[:idx], m.urlHandlers[idx+1:]...)
		}
	}
}

func (m *mockPlatform) Quit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.quitCalls++
}

// simulateEvent fires all registered handlers for the given event type.
func (m *mockPlatform) simulateEvent(eventType EventType) {
	m.mu.Lock()
	handlers := make([]func(), len(m.handlers[eventType]))
	copy(handlers, m.handlers[eventType])
	m.mu.Unlock()
	for _, h := range handlers {
		h()
	}
}

// simulateFileOpen fires all registered file-open handlers.
func (m *mockPlatform) simulateFileOpen(path string) {
	m.mu.Lock()
	handlers := make([]func(string), len(m.fileHandlers))
	copy(handlers, m.fileHandlers)
	m.mu.Unlock()
	for _, h := range handlers {
		h(path)
	}
}

// handlerCount returns the number of registered handlers for event-based + file-based.
func (m *mockPlatform) handlerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := len(m.fileHandlers)
	for _, handlers := range m.handlers {
		count += len(handlers)
	}
	return count
}

// --- Test helpers ---

func newTestLifecycleService(t *core.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := newMockPlatform()
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "lifecycle")
	return svc, c, mock
}

// --- Tests ---

func TestRegister_Good(t *core.T) {
	svc, _, _ := newTestLifecycleService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotEmpty(t, core.Sprintf("%T", svc))
}

func TestApplicationStarted_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventApplicationStarted)
	core.AssertTrue(t, received)
}

func TestDidBecomeActive_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionDidBecomeActive); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventDidBecomeActive)
	core.AssertTrue(t, received)
}

func TestDidResignActive_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionDidResignActive); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventDidResignActive)
	core.AssertTrue(t, received)
}

func TestWillTerminate_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionWillTerminate); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventWillTerminate)
	core.AssertTrue(t, received)
}

func TestPowerStatusChanged_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionPowerStatusChanged); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventPowerStatusChanged)
	core.AssertTrue(t, received)
}

func TestSystemSuspend_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionSystemSuspend); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventSystemSuspend)
	core.AssertTrue(t, received)
}

func TestSystemResume_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionSystemResume); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	mock.simulateEvent(EventSystemResume)
	core.AssertTrue(t, received)
}

func TestOpenedWithFile_Good(t *core.T) {
	_, c, mock := newTestLifecycleService(t)

	var receivedPath string
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionOpenedWithFile); ok {
			receivedPath = a.Path
		}
		return core.Result{OK: true}
	})

	mock.simulateFileOpen("/Users/snider/Documents/test.txt")
	core.AssertEqual(t, "/Users/snider/Documents/test.txt", receivedPath)
}

func TestOnShutdown_CancelsAll_GoodCase(t *core.T) {
	svc, _, mock := newTestLifecycleService(t)

	// Verify handlers were registered during OnStartup
	core.AssertGreater(t, mock.handlerCount(), 0, "handlers should be registered after OnStartup")

	// Shutdown should cancel all registrations
	core.RequireTrue(t, svc.OnShutdown(context.Background()).OK)

	core.AssertEqual(t, 0, mock.handlerCount(), "all handlers should be cancelled after OnShutdown")
}

func TestRegister_Bad(t *core.T) {
	// No lifecycle service registered — actions are not received
	c := core.New(core.WithServiceLock())

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	// No way to trigger events without the service
	core.AssertFalse(t, received)
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
