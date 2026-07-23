// pkg/dialog/service.go
package dialog

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/webview"
	"dappco.re/go/render/display/webkit/pkg/window"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the dialog service to a Core instance.
//
//	c.WithService(dialog.Register(wailsDialog))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	openFile := func(_ context.Context, opts core.Options) core.Result {
		openOpts, err := openFileOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		paths, err := s.platform.OpenFile(openOpts)
		return core.Result{}.New(paths, err)
	}
	saveFile := func(_ context.Context, opts core.Options) core.Result {
		saveOpts, err := saveFileOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		path, err := s.platform.SaveFile(saveOpts)
		return core.Result{}.New(path, err)
	}
	openDirectory := func(_ context.Context, opts core.Options) core.Result {
		openOpts, err := openDirectoryOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		path, err := s.platform.OpenDirectory(openOpts)
		return core.Result{}.New(path, err)
	}
	messageDialog := func(_ context.Context, opts core.Options) core.Result {
		messageOpts, err := messageDialogOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		button, err := s.platform.MessageDialog(messageOpts)
		return core.Result{}.New(button, err)
	}
	info := func(_ context.Context, opts core.Options) core.Result {
		infoOpts, err := infoDialogOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		button, err := s.platform.MessageDialog(infoOpts)
		return core.Result{}.New(button, err)
	}
	question := func(_ context.Context, opts core.Options) core.Result {
		questionOpts, err := questionDialogOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		button, err := s.platform.MessageDialog(questionOpts)
		return core.Result{}.New(button, err)
	}
	warning := func(_ context.Context, opts core.Options) core.Result {
		warningOpts, err := warningDialogOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		button, err := s.platform.MessageDialog(warningOpts)
		return core.Result{}.New(button, err)
	}
	errDialog := func(_ context.Context, opts core.Options) core.Result {
		errOpts, err := errorDialogOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		button, err := s.platform.MessageDialog(errOpts)
		return core.Result{}.New(button, err)
	}
	prompt := func(ctx context.Context, opts core.Options) core.Result {
		promptOpts, err := promptOptionsFrom(opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		windowName, err := s.promptWindowName()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		script := promptScript(promptOpts.Title, promptOpts.Message, promptOpts.DefaultValue)
		task := core.NewOptions(core.Option{Key: "task", Value: webview.TaskEvaluate{Window: windowName, Script: script}})
		result := s.Core().Action("webview.evaluate").Run(ctx, task)
		if !result.OK {
			// Keep the legacy GUI alias as a fallback for older startup wiring.
			result = s.Core().Action("gui.webview.eval").Run(ctx, task)
		}
		if !result.OK {
			if e, ok := result.Value.(error); ok {
				return core.Result{Value: e, OK: false}
			}
			return core.Result{OK: false}
		}
		switch value := result.Value.(type) {
		case nil:
			return core.Result{Value: PromptResult{Confirmed: false}, OK: true}
		case string:
			return core.Result{Value: PromptResult{Value: value, Confirmed: true}, OK: true}
		default:
			return core.Result{Value: PromptResult{Value: core.Sprint(value), Confirmed: true}, OK: true}
		}
	}
	s.Core().Action("dialog.open_file", openFile)
	s.Core().Action("gui.dialog.open", openFile)
	s.Core().Action("dialog.save_file", saveFile)
	s.Core().Action("gui.dialog.save", saveFile)
	s.Core().Action("dialog.open_directory", openDirectory)
	s.Core().Action("gui.dialog.openDirectory", openDirectory)
	s.Core().Action("dialog.message", messageDialog)
	s.Core().Action("gui.dialog.message", messageDialog)
	s.Core().Action("dialog.info", info)
	s.Core().Action("dialog.question", question)
	s.Core().Action("gui.dialog.confirm", question)
	s.Core().Action("dialog.warning", warning)
	s.Core().Action("dialog.error", errDialog)
	s.Core().Action("dialog.prompt", prompt)
	s.Core().Action("gui.dialog.prompt", prompt)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func openFileOptionsFrom(opts core.Options) (OpenFileOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskOpenFile:
			return v.Options, nil
		case TaskOpenFileWithOptions:
			if v.Options != nil {
				return *v.Options, nil
			}
		case OpenFileOptions:
			return v, nil
		}
	}
	return decodeOptions[OpenFileOptions](opts)
}

func saveFileOptionsFrom(opts core.Options) (SaveFileOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskSaveFile:
			return v.Options, nil
		case TaskSaveFileWithOptions:
			if v.Options != nil {
				return *v.Options, nil
			}
		case SaveFileOptions:
			return v, nil
		}
	}
	return decodeOptions[SaveFileOptions](opts)
}

func openDirectoryOptionsFrom(opts core.Options) (OpenDirectoryOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskOpenDirectory:
			return v.Options, nil
		case OpenDirectoryOptions:
			return v, nil
		}
	}
	return decodeOptions[OpenDirectoryOptions](opts)
}

func messageDialogOptionsFrom(opts core.Options) (MessageDialogOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskMessageDialog:
			return v.Options, nil
		case MessageDialogOptions:
			return v, nil
		}
	}
	return decodeOptions[MessageDialogOptions](opts)
}

func infoDialogOptionsFrom(opts core.Options) (MessageDialogOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskInfo:
			return MessageDialogOptions{
				Type:    DialogInfo,
				Title:   v.Title,
				Message: v.Message,
				Buttons: v.Buttons,
			}, nil
		case MessageDialogOptions:
			v.Type = DialogInfo
			return v, nil
		default:
			return MessageDialogOptions{}, core.E("dialog.infoDialogOptionsFrom", "failed to decode info dialog options", nil)
		}
	}
	return typedMessageDialogOptionsFrom(opts, DialogInfo, "dialog.infoDialogOptionsFrom")
}

func warningDialogOptionsFrom(opts core.Options) (MessageDialogOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskWarning:
			return MessageDialogOptions{
				Type:    DialogWarning,
				Title:   v.Title,
				Message: v.Message,
				Buttons: v.Buttons,
			}, nil
		case MessageDialogOptions:
			v.Type = DialogWarning
			return v, nil
		default:
			return MessageDialogOptions{}, core.E("dialog.warningDialogOptionsFrom", "failed to decode warning dialog options", nil)
		}
	}
	return typedMessageDialogOptionsFrom(opts, DialogWarning, "dialog.warningDialogOptionsFrom")
}

func errorDialogOptionsFrom(opts core.Options) (MessageDialogOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskError:
			return MessageDialogOptions{
				Type:    DialogError,
				Title:   v.Title,
				Message: v.Message,
				Buttons: v.Buttons,
			}, nil
		case MessageDialogOptions:
			v.Type = DialogError
			return v, nil
		default:
			return MessageDialogOptions{}, core.E("dialog.errorDialogOptionsFrom", "failed to decode error dialog options", nil)
		}
	}
	return typedMessageDialogOptionsFrom(opts, DialogError, "dialog.errorDialogOptionsFrom")
}

func typedMessageDialogOptionsFrom(opts core.Options, dialogType DialogType, op string) (MessageDialogOptions, resultFailure) {
	if !hasDirectDialogOptions(opts) {
		return MessageDialogOptions{}, core.E(op, "failed to decode dialog options", nil)
	}
	decoded, err := decodeOptions[MessageDialogOptions](opts)
	if err != nil {
		return MessageDialogOptions{}, err
	}
	decoded.Type = dialogType
	return decoded, nil
}

func hasDirectDialogOptions(opts core.Options) bool {
	for _, item := range opts.Items() {
		if item.Key != "task" {
			return true
		}
	}
	return false
}

func questionDialogOptionsFrom(opts core.Options) (MessageDialogOptions, resultFailure) {
	if task := opts.Get("task"); task.OK {
		switch v := task.Value.(type) {
		case TaskQuestion:
			return MessageDialogOptions{
				Type:    DialogQuestion,
				Title:   v.Title,
				Message: v.Message,
				Buttons: v.Buttons,
			}, nil
		case MessageDialogOptions:
			v.Type = DialogQuestion
			return v, nil
		}
	}
	if direct, err := decodeOptions[TaskQuestion](opts); err == nil {
		return MessageDialogOptions{
			Type:    DialogQuestion,
			Title:   direct.Title,
			Message: direct.Message,
			Buttons: direct.Buttons,
		}, nil
	}
	decoded, err := decodeOptions[MessageDialogOptions](opts)
	if err != nil {
		return MessageDialogOptions{}, err
	}
	decoded.Type = DialogQuestion
	return decoded, nil
}

func promptOptionsFrom(opts core.Options) (TaskPrompt, resultFailure) {
	if task := opts.Get("task"); task.OK {
		if v, ok := task.Value.(TaskPrompt); ok {
			return v, nil
		}
	}
	return decodeOptions[TaskPrompt](opts)
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
		return input, core.E("dialog.decodeOptions", "failed to decode dialog options", nil)
	}
	return input, nil
}

func (s *Service) promptWindowName() (string, resultFailure) {
	r := s.Core().QUERY(window.QueryWindowList{})
	if !r.OK {
		return "", core.E("dialog.promptWindowName", "window service unavailable", nil)
	}
	windows, ok := r.Value.([]window.WindowInfo)
	if !ok {
		return "", core.E("dialog.promptWindowName", "unexpected window list result type", nil)
	}
	for _, info := range windows {
		if info.Focused {
			return info.Name, nil
		}
	}
	if len(windows) > 0 {
		return windows[0].Name, nil
	}
	return "", core.E("dialog.promptWindowName", "no application window available for prompt", nil)
}

func promptScript(title, message, defaultValue string) string {
	promptText := title
	if message != "" {
		if promptText != "" {
			promptText += "\n\n"
		}
		promptText += message
	}
	return core.Sprintf(`(() => {
  const value = window.prompt(%s, %s);
  return value === null ? null : value;
})()`, core.JSONMarshalString(promptText), core.JSONMarshalString(defaultValue))
}
