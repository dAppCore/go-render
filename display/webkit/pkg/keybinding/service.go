// pkg/keybinding/service.go
package keybinding

import (
	"context"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform           Platform
	mu                 sync.RWMutex
	registeredBindings map[string]BindingInfo
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("keybinding.add", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskAdd)
		return core.Result{Value: nil, OK: true}.New(s.taskAdd(t))
	})
	s.Core().Action("keybinding.remove", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskRemove)
		return core.Result{Value: nil, OK: true}.New(s.taskRemove(t))
	})
	s.Core().Action("keybinding.process", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskProcess)
		return core.Result{Value: nil, OK: true}.New(s.taskProcess(t))
	})
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

// --- Query Handlers ---

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryList:
		return core.Result{Value: s.queryList(), OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) queryList() []BindingInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]BindingInfo, 0, len(s.registeredBindings))
	for _, info := range s.registeredBindings {
		result = append(result, info)
	}
	return result
}

func (s *Service) taskAdd(t TaskAdd) resultFailure {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.registeredBindings[t.Accelerator]; exists {
		return ErrorAlreadyRegistered
	}

	// Register on platform with a callback that broadcasts ActionTriggered
	err := s.platform.Add(t.Accelerator, func() {
		coreutil.DispatchAction(s.Core(), "keybinding.taskAdd", ActionTriggered{Accelerator: t.Accelerator})
	})
	if err != nil {
		return core.E("keybinding.taskAdd", "platform add failed", err)
	}

	s.registeredBindings[t.Accelerator] = BindingInfo{
		Accelerator: t.Accelerator,
		Description: t.Description,
	}
	return nil
}

func (s *Service) taskRemove(t TaskRemove) resultFailure {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.registeredBindings[t.Accelerator]; !exists {
		return core.E("keybinding.taskRemove", "not registered: "+t.Accelerator, ErrorNotRegistered)
	}

	err := s.platform.Remove(t.Accelerator)
	if err != nil {
		return core.E("keybinding.taskRemove", "platform remove failed", err)
	}

	delete(s.registeredBindings, t.Accelerator)
	return nil
}

// taskProcess triggers the registered handler for the given accelerator programmatically.
// Broadcasts ActionTriggered if handled; returns ErrorNotRegistered if the accelerator is unknown.
//
//	c.Action("keybinding.process").Run(ctx, core.NewOptions(core.Option{Key:"task", Value:keybinding.TaskProcess{Accelerator:"Ctrl+S"}}))
func (s *Service) taskProcess(t TaskProcess) resultFailure {
	s.mu.RLock()
	_, exists := s.registeredBindings[t.Accelerator]
	s.mu.RUnlock()
	if !exists {
		return core.E("keybinding.taskProcess", "not registered: "+t.Accelerator, ErrorNotRegistered)
	}

	handled := s.platform.Process(t.Accelerator)
	if !handled {
		return core.E("keybinding.taskProcess", "platform did not handle: "+t.Accelerator, nil)
	}

	return nil
}
