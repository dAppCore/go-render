// pkg/window/register.go
package window

import core "dappco.re/go"

// Register(p) binds the window service to a Core instance.
// core.WithService(window.Register(window.NewWailsPlatform(app)))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
			specs:          make(map[string]registeredSpec),
		}, OK: true}
	}
}
