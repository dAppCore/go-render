// pkg/events/service.go
package events

import (
	"context"
	"sort"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
)

// Options holds configuration for the events service (currently empty).
type Options struct{}

// Service bridges Wails custom events into Core IPC.
// Emit/On/Off/OnMultiple/Reset are available as Tasks; QueryListeners reads state.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform

	mu        sync.Mutex
	listeners map[string][]func() // IPC-registered cancels per event name
	counts    map[string]int      // listener counts per event name
}

// OnStartup registers query and action handlers.
func (s *Service) OnStartup(_ context.Context) core.Result {
	s.ensureState()
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("events.emit", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskEmit)
		if err := validateEventName("events.emit", t.Name); err != nil {
			return core.Result{Value: err, OK: false}
		}
		if err := s.requirePlatform("events.emit"); err != nil {
			return core.Result{Value: err, OK: false}
		}
		cancelled := s.platform.Emit(t.Name, t.Data)
		return core.Result{Value: cancelled, OK: true}
	})
	s.Core().Action("events.on", func(ctx context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskOn)
		if err := validateEventName("events.on", t.Name); err != nil {
			return core.Result{Value: err, OK: false}
		}
		if err := s.requirePlatform("events.on"); err != nil {
			return core.Result{Value: err, OK: false}
		}
		cancel := s.platform.On(t.Name, func(event *CustomEvent) {
			if event == nil {
				return
			}
			coreutil.DispatchAction(s.Core(), "events.on", ActionEventFired{Event: *event})
		})
		s.mu.Lock()
		if cancel != nil {
			s.listeners[t.Name] = append(s.listeners[t.Name], cancel)
		}
		s.counts[t.Name]++
		s.mu.Unlock()
		return core.Result{OK: true}
	})
	s.Core().Action("events.off", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskOff)
		if err := validateEventName("events.off", t.Name); err != nil {
			return core.Result{Value: err, OK: false}
		}
		if err := s.requirePlatform("events.off"); err != nil {
			return core.Result{Value: err, OK: false}
		}
		s.platform.Off(t.Name)
		s.mu.Lock()
		for _, cancel := range s.listeners[t.Name] {
			cancel()
		}
		delete(s.listeners, t.Name)
		delete(s.counts, t.Name)
		s.mu.Unlock()
		return core.Result{OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) requirePlatform(method string) resultFailure {
	if s == nil || s.platform == nil {
		return core.E(method, "event platform unavailable", nil)
	}
	return nil
}

func validateEventName(method, name string) resultFailure {
	if core.Trim(name) == "" {
		return core.E(method, "event name must not be empty", nil)
	}
	return nil
}

func (s *Service) ensureState() {
	s.mu.Lock()
	if s.listeners == nil {
		s.listeners = make(map[string][]func())
	}
	if s.counts == nil {
		s.counts = make(map[string]int)
	}
	s.mu.Unlock()
}

// OnShutdown cancels all IPC-registered platform listeners.
func (s *Service) OnShutdown(_ context.Context) core.Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cancels := range s.listeners {
		for _, cancel := range cancels {
			if cancel != nil {
				cancel()
			}
		}
	}
	s.listeners = make(map[string][]func())
	s.counts = make(map[string]int)
	return core.Result{OK: true}
}

// HandleIPCEvents satisfies the core.Service interface (no-op for now).
func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryListeners:
		return core.Result{Value: s.listenerSnapshot(), OK: true}
	default:
		return core.Result{}
	}
}

// listenerSnapshot returns a sorted slice of ListenerInfo for all known event names.
//
//	snapshot := s.listenerSnapshot()
//	for _, info := range snapshot { log(info.EventName, info.Count) }
func (s *Service) listenerSnapshot() []ListenerInfo {
	s.mu.Lock()
	snapshot := make([]ListenerInfo, 0, len(s.counts))
	for name, count := range s.counts {
		snapshot = append(snapshot, ListenerInfo{EventName: name, Count: count})
	}
	s.mu.Unlock()

	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].EventName < snapshot[j].EventName
	})

	return snapshot
}
