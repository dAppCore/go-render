// pkg/systray/register.go
package systray

import core "dappco.re/go"

// Register(p) binds the systray service to a Core instance.
// core.WithService(systray.Register(wailsSystray))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, OK: true}
	}
}
