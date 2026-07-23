package dock

import core "dappco.re/go"

// Register(p) binds the dock service to a Core instance.
// core.WithService(dock.Register(wailsDock))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}
