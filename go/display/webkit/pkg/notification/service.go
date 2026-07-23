// pkg/notification/service.go
package notification

import (
	"context"
	"strconv"
	"sync"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/internal/coreutil"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform   Platform
	categories map[string]NotificationCategory
	mu         sync.Mutex
	active     map[string]NotificationOptions
}

func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			categories:     make(map[string]NotificationCategory),
			active:         make(map[string]NotificationOptions),
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	if platform, ok := s.platform.(responsePlatform); ok {
		platform.OnResponse(s.handleResponse)
	}
	send := func(_ context.Context, opts core.Options) core.Result {
		options, err := notificationOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(s.send(options))
	}
	s.Core().Action("notification.request_permission", func(_ context.Context, _ core.Options) core.Result {
		granted, err := s.platform.RequestPermission()
		return core.Result{}.New(granted, err)
	})
	s.Core().Action("gui.notification.requestPermission", func(_ context.Context, _ core.Options) core.Result {
		granted, err := s.platform.RequestPermission()
		return core.Result{}.New(granted, err)
	})
	s.Core().Action("notification.revoke_permission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: nil, OK: true}.New(s.platform.RevokePermission())
	})
	s.Core().Action("gui.notification.revokePermission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: nil, OK: true}.New(s.platform.RevokePermission())
	})
	registerCategory := func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskRegisterCategory)
		return core.Result{Value: nil, OK: true}.New(s.registerCategory(t.Category))
	}
	s.Core().Action("notification.registerCategory", registerCategory)
	s.Core().Action("gui.notification.registerCategory", registerCategory)
	s.Core().Action("notification.clear", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClear)
		return core.Result{Value: nil, OK: true}.New(s.clear(t.ID))
	})
	s.Core().Action("gui.notification.clear", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClear)
		return core.Result{Value: nil, OK: true}.New(s.clear(t.ID))
	})
	s.Core().Action("notification.send", send)
	s.Core().Action("gui.notification.send", send)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryPermission:
		granted, err := s.platform.CheckPermission()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: PermissionStatus{Granted: granted}, OK: true}
	default:
		return core.Result{}
	}
}

// send attempts native notification, falls back to dialog via IPC.
func (s *Service) send(options NotificationOptions) resultFailure {
	// Generate ID if not provided
	if options.ID == "" {
		options.ID = "core-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	options = s.applyCategoryActions(options)

	if options.Update {
		platform, ok := s.platform.(updatePlatform)
		if !ok {
			return core.E("notification.send", "platform does not support notification updates", nil)
		}
		if err := platform.Update(options); err != nil {
			return core.E("notification.send", "native notification update failed", err)
		}
	} else if err := s.platform.Send(options); err != nil {
		// Fallback: show as dialog via IPC
		if err := s.fallbackDialog(options); err != nil {
			return err
		}
	}
	s.mu.Lock()
	s.active[options.ID] = options
	s.mu.Unlock()
	return nil
}

func (s *Service) registerCategory(category NotificationCategory) resultFailure {
	if category.ID == "" {
		return core.E("notification.registerCategory", "category ID is required", nil)
	}
	if platform, ok := s.platform.(categoryPlatform); ok {
		if err := platform.RegisterCategory(category); err != nil {
			return core.E("notification.registerCategory", "native category registration failed", err)
		}
	}
	s.mu.Lock()
	s.categories[category.ID] = category
	s.mu.Unlock()
	return nil
}

func (s *Service) handleResponse(notificationID, actionID, userText string) {
	switch actionID {
	case "", "DEFAULT_ACTION":
		coreutil.DispatchAction(s.Core(), "notification.click", ActionNotificationClicked{ID: notificationID})
	default:
		coreutil.DispatchAction(s.Core(), "notification.action", ActionNotificationActionTriggered{
			NotificationID: notificationID,
			ActionID:       actionID,
			UserText:       userText,
		})
	}
}

func (s *Service) applyCategoryActions(options NotificationOptions) NotificationOptions {
	if options.CategoryID == "" || len(options.Actions) > 0 {
		return options
	}

	s.mu.Lock()
	category, ok := s.categories[options.CategoryID]
	s.mu.Unlock()
	if !ok || len(category.Actions) == 0 {
		return options
	}

	options.Actions = append([]NotificationAction(nil), category.Actions...)
	return options
}

// fallbackDialog shows a dialog via IPC when native notifications fail.
func (s *Service) fallbackDialog(options NotificationOptions) resultFailure {
	// Map severity to dialog type
	var dt dialog.DialogType
	switch options.Severity {
	case SeverityWarning:
		dt = dialog.DialogWarning
	case SeverityError:
		dt = dialog.DialogError
	default:
		dt = dialog.DialogInfo
	}

	message := options.Message
	if options.Subtitle != "" {
		message = options.Subtitle + "\n\n" + message
	}

	r := s.Core().Action("dialog.message").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskMessageDialog{
			Options: dialog.MessageDialogOptions{
				Type:    dt,
				Title:   options.Title,
				Message: message,
				Buttons: []string{"OK"},
			},
		}},
	))
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return err
		}
	}
	return nil
}

func (s *Service) clear(id string) resultFailure {
	if clearer, ok := s.platform.(ClearPlatform); ok {
		if err := clearer.Clear(id); err != nil {
			return err
		}
	}

	ids := s.removeActive(id)
	for _, notificationID := range ids {
		coreutil.DispatchAction(s.Core(), "notification.dismiss", ActionNotificationDismissed{ID: notificationID})
	}
	return nil
}

func (s *Service) removeActive(id string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id != "" {
		if _, ok := s.active[id]; !ok {
			return nil
		}
		delete(s.active, id)
		return []string{id}
	}

	ids := make([]string, 0, len(s.active))
	for notificationID := range s.active {
		ids = append(ids, notificationID)
	}
	clear(s.active)
	return ids
}

func notificationOptionsFrom(opts core.Options) (NotificationOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskSend:
			return v.Options, nil
		case NotificationOptions:
			return v, nil
		}
	}
	return decodeOptions[NotificationOptions](opts)
}

func decodeOptions[T any](opts core.Options) (T, resultFailure) {
	var input T
	items := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		items[item.Key] = item.Value
	}
	if len(items) == 0 {
		return input, nil
	}
	result := core.JSONUnmarshalString(core.JSONMarshalString(items), &input)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return input, err
		}
		return input, core.E("notification.decodeOptions", "failed to decode notification options", nil)
	}
	return input, nil
}
