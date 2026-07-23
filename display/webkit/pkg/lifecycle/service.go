package lifecycle

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	cancels  []func()
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	eventActions := map[EventType]func(){
		EventApplicationStarted: func() { coreutil.DispatchAction(s.Core(), "lifecycle.applicationStarted", ActionApplicationStarted{}) },
		EventWillTerminate:      func() { coreutil.DispatchAction(s.Core(), "lifecycle.willTerminate", ActionWillTerminate{}) },
		EventDidBecomeActive:    func() { coreutil.DispatchAction(s.Core(), "lifecycle.didBecomeActive", ActionDidBecomeActive{}) },
		EventDidResignActive:    func() { coreutil.DispatchAction(s.Core(), "lifecycle.didResignActive", ActionDidResignActive{}) },
		EventPowerStatusChanged: func() { coreutil.DispatchAction(s.Core(), "lifecycle.powerStatusChanged", ActionPowerStatusChanged{}) },
		EventSystemSuspend:      func() { coreutil.DispatchAction(s.Core(), "lifecycle.systemSuspend", ActionSystemSuspend{}) },
		EventSystemResume:       func() { coreutil.DispatchAction(s.Core(), "lifecycle.systemResume", ActionSystemResume{}) },
	}

	for eventType, handler := range eventActions {
		cancel := s.platform.OnApplicationEvent(eventType, handler)
		s.cancels = append(s.cancels, cancel)
	}

	cancel := s.platform.OnOpenedWithFile(func(path string) {
		coreutil.DispatchAction(s.Core(), "lifecycle.openedWithFile", ActionOpenedWithFile{Path: path})
	})
	s.cancels = append(s.cancels, cancel)

	urlCancel := s.platform.OnLaunchedWithUrl(func(url string) {
		coreutil.DispatchAction(s.Core(), "lifecycle.launchedWithUrl", ActionLaunchedWithUrl{URL: url})
	})
	s.cancels = append(s.cancels, urlCancel)

	s.Core().Action("lifecycle.quit", func(_ context.Context, _ core.Options) core.Result {
		if s.platform == nil {
			return core.Result{Value: core.E("lifecycle.quit", "platform unavailable", nil), OK: false}
		}
		s.platform.Quit()
		return core.Result{OK: true}
	})

	return core.Result{OK: true}
}

func (s *Service) OnShutdown(_ context.Context) core.Result {
	for _, cancel := range s.cancels {
		cancel()
	}
	s.cancels = nil
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}
