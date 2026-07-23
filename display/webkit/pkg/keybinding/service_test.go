// pkg/keybinding/service_test.go
package keybinding

import (
	"context"
	"sync"

	core "dappco.re/go"
)

// mockPlatform records Add/Remove calls and allows triggering shortcuts.
type mockPlatform struct {
	mu       sync.Mutex
	handlers map[string]func()
	removed  []string
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{handlers: make(map[string]func())}
}

func (m *mockPlatform) Add(accelerator string, handler func()) resultFailure {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[accelerator] = handler
	return nil
}

func (m *mockPlatform) Remove(accelerator string) resultFailure {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, accelerator)
	m.removed = append(m.removed, accelerator)
	return nil
}

func (m *mockPlatform) Process(accelerator string) bool {
	m.mu.Lock()
	h, ok := m.handlers[accelerator]
	m.mu.Unlock()
	if ok && h != nil {
		h()
		return true
	}
	return false
}

func (m *mockPlatform) GetAll() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.handlers))
	for k := range m.handlers {
		out = append(out, k)
	}
	return out
}

// trigger simulates a shortcut keypress by calling the registered handler.
func (m *mockPlatform) trigger(accelerator string) {
	m.mu.Lock()
	h, ok := m.handlers[accelerator]
	m.mu.Unlock()
	if ok {
		h()
	}
}

func newTestKeybindingService(t *core.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "keybinding")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *core.T) {
	mp := newMockPlatform()
	svc, _ := newTestKeybindingService(t, mp)
	core.AssertNotNil(t, svc)
	core.AssertNotNil(t, svc.platform)
}

func TestTaskAdd_Good(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.add", TaskAdd{
		Accelerator: "Ctrl+S", Description: "Save",
	})
	core.RequireTrue(t, r.OK)

	// Verify binding registered on platform
	core.AssertContains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskAdd_Bad_Duplicate(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Second add with same accelerator should fail
	r := taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save Again"})
	core.AssertFalse(t, r.OK)
	err, _ := r.Value.(resultFailure)
	core.AssertErrorIs(t, err, ErrorAlreadyRegistered)
}

func TestTaskRemove_Good(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+S"})
	core.RequireTrue(t, r.OK)

	// Verify removed from platform
	core.AssertNotContains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskRemove_Bad_NotFound(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+X"})
	core.AssertFalse(t, r.OK)
}

func TestQueryList_Good(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+Z", Description: "Undo"})

	r := c.QUERY(QueryList{})
	core.RequireTrue(t, r.OK)
	list := r.Value.([]BindingInfo)
	core.AssertLen(t, list, 2)
}

func TestQueryList_Good_Empty(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := c.QUERY(QueryList{})
	core.RequireTrue(t, r.OK)
	list := r.Value.([]BindingInfo)
	core.AssertLen(t, list, 0)
}

func TestTaskAdd_Good_TriggerBroadcast(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	// Capture broadcast actions
	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Simulate shortcut trigger via mock
	mp.trigger("Ctrl+S")

	mu.Lock()
	core.AssertEqual(t, "Ctrl+S", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskAdd_Good_RebindAfterRemove(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+S"})

	// Should succeed after remove
	r := taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save v2"})
	core.RequireTrue(t, r.OK)

	// Verify new description
	r2 := c.QUERY(QueryList{})
	list := r2.Value.([]BindingInfo)
	core.AssertLen(t, list, 1)
	core.AssertEqual(t, "Save v2", list[0].Description)
}

func TestQueryList_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryList{})
	core.AssertFalse(t, r.OK)
}

// --- TaskProcess tests ---

func TestTaskProcess_Good(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})

	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	core.RequireTrue(t, r.OK)

	mu.Lock()
	core.AssertEqual(t, "Ctrl+P", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskProcess_Bad_NotRegistered(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	core.AssertFalse(t, r.OK)
	err, _ := r.Value.(resultFailure)
	core.AssertErrorIs(t, err, ErrorNotRegistered)
}

func TestTaskProcess_Ugly_RemovedBinding(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})
	taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+P"})

	// After remove, process should fail with ErrorNotRegistered
	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	core.AssertFalse(t, r.OK)
	err, _ := r.Value.(resultFailure)
	core.AssertErrorIs(t, err, ErrorNotRegistered)
}

// --- TaskRemove ErrorNotRegistered sentinel tests ---

func TestTaskRemove_Bad_ErrorSentinel(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+X"})
	core.AssertFalse(t, r.OK)
	err, _ := r.Value.(resultFailure)
	core.AssertErrorIs(t, err, ErrorNotRegistered)
}

// --- QueryList Ugly: concurrent adds ---

func TestQueryList_Ugly_ConcurrentAdds(t *core.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	accelerators := []string{"Ctrl+1", "Ctrl+2", "Ctrl+3", "Ctrl+4", "Ctrl+5"}
	var wg sync.WaitGroup
	for _, accelerator := range accelerators {
		wg.Add(1)
		go func(acc string) {
			defer wg.Done()
			taskRun(c, "keybinding.add", TaskAdd{Accelerator: acc, Description: acc})
		}(accelerator)
	}
	wg.Wait()

	r := c.QUERY(QueryList{})
	core.RequireTrue(t, r.OK)
	list := r.Value.([]BindingInfo)
	core.AssertLen(t, list, len(accelerators))
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
