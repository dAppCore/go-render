package contextmenu

import core "dappco.re/go"

// Register(p) binds the context menu service to a Core instance.
// core.WithService(contextmenu.Register(wailsContextMenu))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime:  core.NewServiceRuntime[Options](c, Options{}),
			platform:        p,
			registeredMenus: make(map[string]ContextMenuDef),
		}, OK: true}
	}
}
