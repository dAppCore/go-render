package keybinding

import core "dappco.re/go"

// Register(p) binds the keybinding service to a Core instance.
// core.WithService(keybinding.Register(wailsKeybinding))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime:     core.NewServiceRuntime[Options](c, Options{}),
			platform:           p,
			registeredBindings: make(map[string]BindingInfo),
		}, OK: true}
	}
}
