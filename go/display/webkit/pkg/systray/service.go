package systray

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
	"dappco.re/go/render/display/webkit/pkg/notification"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	iconPath string
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	r := s.Core().QUERY(QueryConfig{})
	if r.OK {
		if trayConfig, ok := r.Value.(map[string]any); ok {
			s.applyConfig(trayConfig)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("systray.set_icon", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetTrayIconFromOptions(opts)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetIcon(t.Data))
	})
	s.Core().Action("systray.set_template_icon", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetTrayTemplateIconFromOptions(opts)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetTemplateIcon(t.Data))
	})
	s.Core().Action("systray.set_tooltip", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetTrayTooltipFromOptions(opts)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetTooltip(t.Tooltip))
	})
	s.Core().Action("systray.set_label", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetTrayLabelFromOptions(opts)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetLabel(t.Label))
	})
	s.Core().Action("systray.set_menu", func(_ context.Context, opts core.Options) core.Result {
		t := taskSetTrayMenuFromOptions(opts)
		return core.Result{Value: nil, OK: true}.New(s.taskSetTrayMenu(t))
	})
	s.Core().Action("systray.show_message", func(_ context.Context, opts core.Options) core.Result {
		t := taskShowMessageFromOptions(opts)
		if err := s.manager.ShowMessage(t.Title, t.Message); err == nil {
			return core.Result{OK: true}
		} else {
			fallback := s.Core().Action("notification.send").Run(context.Background(), core.NewOptions(
				core.Option{Key: "task", Value: notification.TaskSend{Options: notification.NotificationOptions{
					Title:   t.Title,
					Message: t.Message,
				}}},
			))
			if fallback.OK {
				return core.Result{OK: true}
			}
			if fallbackErr, ok := fallback.Value.(error); ok {
				return core.Result{Value: core.E("systray.show_message", "tray message failed and notification fallback failed", fallbackErr), OK: false}
			}
			return core.Result{Value: err, OK: false}
		}
	})
	s.Core().Action("gui.tray.showMessage", func(_ context.Context, opts core.Options) core.Result {
		t := taskShowMessageFromOptions(opts)
		if err := s.manager.ShowMessage(t.Title, t.Message); err == nil {
			return core.Result{OK: true}
		} else {
			fallback := s.Core().Action("notification.send").Run(context.Background(), core.NewOptions(
				core.Option{Key: "task", Value: notification.TaskSend{Options: notification.NotificationOptions{
					Title:   t.Title,
					Message: t.Message,
				}}},
			))
			if fallback.OK {
				return core.Result{OK: true}
			}
			if fallbackErr, ok := fallback.Value.(error); ok {
				return core.Result{Value: core.E("systray.show_message", "tray message failed and notification fallback failed", fallbackErr), OK: false}
			}
			return core.Result{Value: err, OK: false}
		}
	})
	s.Core().Action("systray.showPanel", func(_ context.Context, _ core.Options) core.Result {
		// Panel show — deferred (requires WindowHandle integration)
		return core.Result{OK: true}
	})
	s.Core().Action("systray.hidePanel", func(_ context.Context, _ core.Options) core.Result {
		// Panel hide — deferred (requires WindowHandle integration)
		return core.Result{OK: true}
	})
	s.Core().Action("systray.attach_window", func(_ context.Context, opts core.Options) core.Result {
		t := taskAttachWindowFromOptions(opts)
		if t.Name == "" {
			return core.Result{Value: core.E("systray.attach_window", "window name required", nil), OK: false}
		}
		if s.manager == nil {
			return core.Result{Value: core.E("systray.attach_window", "tray manager unavailable", nil), OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.manager.AttachWindow(namedWindowHandle{name: t.Name}, t.OffsetX, t.OffsetY))
	})
	return core.Result{OK: true}
}

// namedWindowHandle is the minimal WindowHandle implementation the
// platform uses to look up the underlying wails window by name.
type namedWindowHandle struct{ name string }

func (h namedWindowHandle) Name() string { return h.name }

func taskAttachWindowFromOptions(opts core.Options) TaskAttachWindow {
	if task := opts.Get("task"); task.OK {
		if t, ok := task.Value.(TaskAttachWindow); ok {
			return t
		}
	}
	return TaskAttachWindow{}
}

func (s *Service) applyConfig(configData map[string]any) {
	tooltip, _ := configData["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	if err := s.manager.Setup(tooltip, tooltip); err != nil {
		return
	}

	if iconPath, ok := configData["icon"].(string); ok && iconPath != "" {
		// Icon loading is deferred to when assets are available.
		// Store the path for later use.
		s.iconPath = iconPath
	}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryInfo:
		return core.Result{Value: s.manager.GetInfo(), OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) taskSetTrayMenu(t TaskSetTrayMenu) resultFailure {
	// Register IPC-emitting callbacks for each menu item
	for _, item := range t.Items {
		if item.ActionID != "" {
			actionID := item.ActionID
			s.manager.RegisterCallback(actionID, func() {
				coreutil.DispatchAction(s.Core(), "systray.taskSetTrayMenu", ActionTrayMenuItemClicked{ActionID: actionID})
			})
		}
	}
	return s.manager.SetMenu(t.Items)
}

func (s *Service) Manager() *Manager {
	return s.manager
}

func taskSetTrayIconFromOptions(opts core.Options) TaskSetTrayIcon {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetTrayIcon:
			return value
		case map[string]any:
			var decoded TaskSetTrayIcon
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskSetTrayIcon
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetTrayIcon{}
}

func taskSetTrayTemplateIconFromOptions(opts core.Options) TaskSetTrayTemplateIcon {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetTrayTemplateIcon:
			return value
		case map[string]any:
			var decoded TaskSetTrayTemplateIcon
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskSetTrayTemplateIcon
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetTrayTemplateIcon{}
}

func taskSetTrayTooltipFromOptions(opts core.Options) TaskSetTrayTooltip {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetTrayTooltip:
			return value
		case map[string]any:
			var decoded TaskSetTrayTooltip
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskSetTrayTooltip
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetTrayTooltip{}
}

func taskSetTrayLabelFromOptions(opts core.Options) TaskSetTrayLabel {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetTrayLabel:
			return value
		case map[string]any:
			var decoded TaskSetTrayLabel
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskSetTrayLabel
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetTrayLabel{}
}

func taskSetTrayMenuFromOptions(opts core.Options) TaskSetTrayMenu {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskSetTrayMenu:
			return value
		case map[string]any:
			var decoded TaskSetTrayMenu
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskSetTrayMenu
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskSetTrayMenu{}
}

func taskShowMessageFromOptions(opts core.Options) TaskShowMessage {
	if task := opts.Get("task"); task.OK {
		switch value := task.Value.(type) {
		case TaskShowMessage:
			return value
		case map[string]any:
			var decoded TaskShowMessage
			if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &decoded); result.OK {
				return decoded
			}
		}
	}
	var decoded TaskShowMessage
	if result := core.JSONUnmarshalString(core.JSONMarshalString(optsToMap(opts)), &decoded); result.OK {
		return decoded
	}
	return TaskShowMessage{}
}

func optsToMap(opts core.Options) map[string]any {
	items := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		items[item.Key] = item.Value
	}
	return items
}
