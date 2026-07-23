// pkg/menu/register.go
package menu

import core "dappco.re/go"

// Register(p) binds the menu service to a Core instance.
// core.WithService(menu.Register(wailsMenu))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, OK: true}
	}
}
