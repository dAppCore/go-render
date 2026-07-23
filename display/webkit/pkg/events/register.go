// pkg/events/register.go
package events

import core "dappco.re/go"

// Register binds the events service to a Core instance.
//
//	core.WithService(events.Register(wailsEventPlatform))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			listeners:      make(map[string][]func()),
			counts:         make(map[string]int),
		}, OK: true}
	}
}
