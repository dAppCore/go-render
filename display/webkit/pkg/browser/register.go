package browser

import core "dappco.re/go"

// Register(p) binds the browser service to a Core instance.
// core.WithService(browser.Register(wailsBrowser))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}
