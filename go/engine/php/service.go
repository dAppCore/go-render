// SPDX-License-Identifier: EUPL-1.2

// Service registration for the php repo. Single repo-level service
// that exposes the PHP project introspection surface (Laravel /
// FrankenPHP / package-manager detection) as IPC actions so the
// CorePHP toolchain can be supervised by go-process and queried over
// IPC.
//
// Bridge HTTP server, FrankenPHP serve, AddPHPCommands command-tree
// helpers stay on direct package use — they own their own lifecycles
// (Bridge.NewBridge starts a listener; AddPHPCommands mutates a
// *core.Core's command tree) and the Service is the canonical
// registration entry-point alongside them.
//
//	c, _ := core.New(
//	    core.WithName("php", php.NewService(php.PhpConfig{})),
//	)
//	r := c.Action("php.detect.laravel").Run(ctx, core.NewOptions(
//	    core.Option{Key: "dir", Value: "/srv/app"},
//	))

package php

import (
	"context"

	core "dappco.re/go"
)

// PhpConfig configures the php service. Empty config gives the default
// detection surface (per-call dir argument; no global root needed).
//
// Usage example: `cfg := php.PhpConfig{}`
type PhpConfig struct{}

// Service is the registerable handle for the php repo — embeds
// *core.ServiceRuntime[PhpConfig] for typed options access. The
// detection surface is exposed as IPC actions; FrankenPHP serve +
// Bridge HTTP listener stay on direct package use because they
// manage their own listeners.
//
// Usage example: `svc := core.MustServiceFor[*php.CoreService](c, "php"); _ = svc`
type CoreService struct {
	*core.ServiceRuntime[PhpConfig]
	registrations core.Once
}

// NewService returns a factory that produces a *Service ready for
// c.Service() registration.
//
// Usage example: `c, _ := core.New(core.WithName("php", php.NewService(php.PhpConfig{})))`
func NewService(config PhpConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Ok(&CoreService{
			ServiceRuntime: core.NewServiceRuntime(c, config),
		})
	}
}

// Register builds the php service with default PhpConfig and returns
// the service Result directly — the imperative-style alternative to
// NewService for consumers wiring services without WithName options.
//
// Usage example: `r := php.Register(c); svc := r.Value.(*php.CoreService)`
func Register(c *core.Core) core.Result {
	return NewService(PhpConfig{})(c)
}

// OnStartup registers the php action handlers on the attached Core.
// Implements core.Startable. Idempotent via core.Once.
//
// Usage example: `r := svc.OnStartup(ctx)`
func (s *CoreService) OnStartup(context.Context) core.Result {
	if s == nil {
		return core.Ok(nil)
	}
	s.registrations.Do(func() {
		c := s.Core()
		if c == nil {
			return
		}
		c.Action("php.detect.laravel", s.handleDetectLaravel)
		c.Action("php.detect.frankenphp", s.handleDetectFrankenPHP)
		c.Action("php.detect.services", s.handleDetectServices)
		c.Action("php.detect.package_manager", s.handleDetectPackageManager)
		c.Action("php.laravel.app_name", s.handleLaravelAppName)
		c.Action("php.laravel.app_url", s.handleLaravelAppURL)
	})
	return core.Ok(nil)
}

// OnShutdown is a no-op — detection helpers are stateless and the
// Bridge / FrankenPHP serve listeners (when wired separately) own
// their own shutdown via direct method use. Implements core.Stoppable.
//
// Usage example: `r := svc.OnShutdown(ctx)`
func (s *CoreService) OnShutdown(context.Context) core.Result {
	return core.Ok(nil)
}

// handleDetectLaravel — `php.detect.laravel` action handler. Reads
// opts.dir and returns bool in r.Value indicating whether the
// directory contains a Laravel project.
//
// Usage example: `r := c.Action("php.detect.laravel").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleDetectLaravel(_ core.Context, opts core.Options) core.Result {
	return core.Ok(IsLaravelProject(opts.String("dir")))
}

// handleDetectFrankenPHP — `php.detect.frankenphp` action handler.
// Reads opts.dir and returns bool in r.Value indicating whether the
// directory is configured for FrankenPHP.
//
// Usage example: `r := c.Action("php.detect.frankenphp").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleDetectFrankenPHP(_ core.Context, opts core.Options) core.Result {
	return core.Ok(IsFrankenPHPProject(opts.String("dir")))
}

// handleDetectServices — `php.detect.services` action handler. Reads
// opts.dir and returns []DetectedService in r.Value covering the
// detected runtime services (Vite, Horizon, Reverb, Redis, etc.).
//
// Usage example: `r := c.Action("php.detect.services").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleDetectServices(_ core.Context, opts core.Options) core.Result {
	return core.Ok(DetectServices(opts.String("dir")))
}

// handleDetectPackageManager — `php.detect.package_manager` action
// handler. Reads opts.dir and returns the JS package-manager name
// (npm/pnpm/yarn/bun) in r.Value, or empty string when absent.
//
// Usage example: `r := c.Action("php.detect.package_manager").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleDetectPackageManager(_ core.Context, opts core.Options) core.Result {
	return core.Ok(DetectPackageManager(opts.String("dir")))
}

// handleLaravelAppName — `php.laravel.app_name` action handler.
// Reads opts.dir and returns the Laravel APP_NAME from .env in
// r.Value.
//
// Usage example: `r := c.Action("php.laravel.app_name").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleLaravelAppName(_ core.Context, opts core.Options) core.Result {
	return core.Ok(GetLaravelAppName(opts.String("dir")))
}

// handleLaravelAppURL — `php.laravel.app_url` action handler. Reads
// opts.dir and returns the Laravel APP_URL from .env in r.Value.
//
// Usage example: `r := c.Action("php.laravel.app_url").Run(ctx, core.NewOptions(core.Option{Key: "dir", Value: "/srv/app"}))`
func (s *CoreService) handleLaravelAppURL(_ core.Context, opts core.Options) core.Result {
	return core.Ok(GetLaravelAppURL(opts.String("dir")))
}
