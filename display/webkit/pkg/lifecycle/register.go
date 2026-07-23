package lifecycle

import core "dappco.re/go"

// Register(p) binds the lifecycle service to a Core instance.
// core.WithService(lifecycle.Register(wailsLifecycle))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}
