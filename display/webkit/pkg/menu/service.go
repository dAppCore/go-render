// pkg/menu/service.go
package menu

import (
	"context"

	core "dappco.re/go"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager      *Manager
	platform     Platform
	menuItems    []MenuItem
	showDevTools bool
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	r := s.Core().QUERY(QueryConfig{})
	if r.OK {
		if menuConfig, ok := r.Value.(map[string]any); ok {
			s.applyConfig(menuConfig)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("menu.set_app_menu", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetAppMenuFromOptions(opts)
		if s.manager == nil || s.manager.Platform() == nil {
			return core.Result{Value: core.E("menu.set_app_menu", "menu manager unavailable", nil), OK: false}
		}
		s.menuItems = t.Items
		s.manager.SetApplicationMenu(t.Items)
		return core.Result{OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) applyConfig(configData map[string]any) {
	if v, ok := configData["show_dev_tools"]; ok {
		if show, ok := v.(bool); ok {
			s.showDevTools = show
		}
	}
}

func (s *Service) ShowDevTools() bool {
	return s.showDevTools
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryGetAppMenu:
		return core.Result{Value: s.menuItems, OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) Manager() *Manager {
	return s.manager
}

func taskSetAppMenuFromOptions(opts core.Options) TaskSetAppMenu {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetAppMenu:
			return value
		case map[string]any:
			var decoded TaskSetAppMenu
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}

	var decoded TaskSetAppMenu
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetAppMenu{}
}

func optsToMap(opts core.Options) map[string]any {
	items := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		items[item.Key] = item.Value
	}
	return items
}
