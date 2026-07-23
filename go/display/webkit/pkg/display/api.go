package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/systray"
)

const writeClipboardImageOp = "display.WriteClipboardImage"

// Screen is the public display-screen shape used by the display service API.
type Screen struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	ScaleFactor float64 `json:"scaleFactor"`
	IsPrimary   bool    `json:"isPrimary"`
}

// WorkArea is a usable screen rectangle.
type WorkArea = screen.Rect

// FileFilter describes a file type filter for the display service dialogs.
type FileFilter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

// OpenFileOptions configures an open-file dialog.
type OpenFileOptions struct {
	Title            string       `json:"title,omitempty"`
	DefaultDirectory string       `json:"defaultDirectory,omitempty"`
	DefaultFilename  string       `json:"defaultFilename,omitempty"`
	AllowMultiple    bool         `json:"allowMultiple,omitempty"`
	Filters          []FileFilter `json:"filters,omitempty"`
}

// SaveFileOptions configures a save-file dialog.
type SaveFileOptions struct {
	Title            string       `json:"title,omitempty"`
	DefaultDirectory string       `json:"defaultDirectory,omitempty"`
	DefaultFilename  string       `json:"defaultFilename,omitempty"`
	Filters          []FileFilter `json:"filters,omitempty"`
}

// OpenDirectoryOptions configures a folder picker dialog.
type OpenDirectoryOptions struct {
	Title            string `json:"title,omitempty"`
	DefaultDirectory string `json:"defaultDirectory,omitempty"`
}

// TrayMenuItem describes a tray menu entry.
type TrayMenuItem struct {
	Label       string         `json:"label,omitempty"`
	ActionID    string         `json:"actionId,omitempty"`
	IsSeparator bool           `json:"isSeparator,omitempty"`
	Children    []TrayMenuItem `json:"children,omitempty"`
}

// NotificationOptions configures a native notification.
type NotificationOptions struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

// Theme reports whether the active theme is dark.
type Theme struct {
	IsDark bool `json:"isDark"`
}

func unexpectedResultType(method string) resultFailure {
	return core.E(method, "unexpected result type", nil)
}

func failedQuery(method, query string) resultFailure {
	return core.E(method, query+" query failed", nil)
}

func (s *Service) GetScreens() []*Screen {
	r := s.Core().QUERY(screen.QueryAll{})
	if !r.OK {
		return []*Screen{}
	}
	screens, ok := r.Value.([]screen.Screen)
	if !ok {
		return []*Screen{}
	}
	result := make([]*Screen, 0, len(screens))
	for i := range screens {
		result = append(result, screenToDisplay(&screens[i]))
	}
	return result
}

func (s *Service) GetScreen(id string) (*Screen, resultFailure) {
	r := s.Core().QUERY(screen.QueryByID{ID: id})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, failedQuery("display.GetScreen", "screen.queryByID")
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, unexpectedResultType("display.GetScreen")
	}
	return screenToDisplay(scr), nil
}

func (s *Service) GetPrimaryScreen() (*Screen, resultFailure) {
	r := s.Core().QUERY(screen.QueryPrimary{})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, failedQuery("display.GetPrimaryScreen", "screen.queryPrimary")
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, unexpectedResultType("display.GetPrimaryScreen")
	}
	return screenToDisplay(scr), nil
}

func (s *Service) GetScreenAtPoint(x, y int) (*Screen, resultFailure) {
	r := s.Core().QUERY(screen.QueryAtPoint{X: x, Y: y})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, failedQuery("display.GetScreenAtPoint", "screen.queryAtPoint")
	}
	scr, ok := r.Value.(*screen.Screen)
	if !ok {
		return nil, unexpectedResultType("display.GetScreenAtPoint")
	}
	return screenToDisplay(scr), nil
}

func (s *Service) GetScreenForWindow(name string) (*Screen, resultFailure) {
	info, err := s.GetWindowInfo(name)
	if err != nil || info == nil {
		return nil, err
	}
	return s.GetScreenAtPoint(info.X+max(info.Width/2, 1), info.Y+max(info.Height/2, 1))
}

func (s *Service) GetWorkAreas() []*WorkArea {
	r := s.Core().QUERY(screen.QueryWorkAreas{})
	if !r.OK {
		return []*WorkArea{}
	}
	areas, ok := r.Value.([]screen.Rect)
	if !ok {
		return []*WorkArea{}
	}
	result := make([]*WorkArea, 0, len(areas))
	for i := range areas {
		area := areas[i]
		result = append(result, &area)
	}
	return result
}

func (s *Service) OpenSingleFileDialog(opts OpenFileOptions) (string, resultFailure) {
	paths, err := s.OpenFileDialog(opts)
	if err != nil || len(paths) == 0 {
		return "", err
	}
	return paths[0], nil
}

func (s *Service) OpenFileDialog(opts OpenFileOptions) ([]string, resultFailure) {
	result := s.Core().Action("dialog.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskOpenFile{Options: toDialogOpenFileOptions(opts)}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, err
		}
		return nil, core.E("display.OpenFileDialog", "dialog.openFile action failed", nil)
	}
	paths, ok := result.Value.([]string)
	if !ok {
		return nil, unexpectedResultType("display.OpenFileDialog")
	}
	return paths, nil
}

func (s *Service) SaveFileDialog(opts SaveFileOptions) (string, resultFailure) {
	result := s.Core().Action("dialog.save_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskSaveFile{Options: toDialogSaveFileOptions(opts)}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return "", err
		}
		return "", core.E("display.SaveFileDialog", "dialog.saveFile action failed", nil)
	}
	path, ok := result.Value.(string)
	if !ok {
		return "", unexpectedResultType("display.SaveFileDialog")
	}
	return path, nil
}

func (s *Service) OpenDirectoryDialog(opts OpenDirectoryOptions) (string, resultFailure) {
	result := s.Core().Action("dialog.open_directory").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskOpenDirectory{Options: toDialogOpenDirectoryOptions(opts)}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return "", err
		}
		return "", core.E("display.OpenDirectoryDialog", "dialog.openDirectory action failed", nil)
	}
	path, ok := result.Value.(string)
	if !ok {
		return "", unexpectedResultType("display.OpenDirectoryDialog")
	}
	return path, nil
}

func (s *Service) ConfirmDialog(title, message string) (bool, resultFailure) {
	result := s.Core().Action("dialog.question").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskQuestion{
			Title:   title,
			Message: message,
			Buttons: []string{"Yes", "No"},
		}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return false, err
		}
		return false, core.E("display.ConfirmDialog", "dialog.question action failed", nil)
	}
	button, ok := result.Value.(string)
	if !ok {
		return false, unexpectedResultType("display.ConfirmDialog")
	}
	return button == "Yes", nil
}

func (s *Service) PromptDialog(title, message string) (string, bool, resultFailure) {
	result := s.Core().Action("dialog.prompt").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskPrompt{Title: title, Message: message}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return "", false, err
		}
		return "", false, core.E("display.PromptDialog", "dialog.prompt action failed", nil)
	}
	prompt, ok := result.Value.(dialog.PromptResult)
	if !ok {
		return "", false, core.E("display.PromptDialog", "unexpected result type", nil)
	}
	return prompt.Value, prompt.Confirmed, nil
}

func (s *Service) SetTrayIcon(icon []byte) resultFailure {
	result := s.Core().Action("systray.set_icon").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayIcon{Data: icon}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.SetTrayIcon", "systray.set_icon action failed", nil)
	}
	return nil
}

func (s *Service) SetTrayTooltip(tooltip string) resultFailure {
	result := s.Core().Action("systray.set_tooltip").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayTooltip{Tooltip: tooltip}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.SetTrayTooltip", "systray.set_tooltip action failed", nil)
	}
	return nil
}

func (s *Service) SetTrayLabel(label string) resultFailure {
	result := s.Core().Action("systray.set_label").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayLabel{Label: label}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.SetTrayLabel", "systray.set_label action failed", nil)
	}
	return nil
}

func (s *Service) SetTrayMenu(items []TrayMenuItem) resultFailure {
	result := s.Core().Action("systray.set_menu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayMenu{Items: trayMenuItemsToSystray(items)}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.SetTrayMenu", "systray.set_menu action failed", nil)
	}
	return nil
}

func (s *Service) GetTrayInfo() map[string]any {
	r := s.Core().QUERY(systray.QueryInfo{})
	if !r.OK {
		return nil
	}
	info, _ := r.Value.(map[string]any)
	return info
}

func (s *Service) ShowTrayMessage(title, message string) resultFailure {
	result := s.Core().Action("systray.show_message").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskShowMessage{Title: title, Message: message}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.ShowTrayMessage", "systray.show_message action failed", nil)
	}
	return nil
}

func (s *Service) ReadClipboard() (string, resultFailure) {
	r := s.Core().QUERY(clipboard.QueryText{})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return "", err
		}
		return "", nil
	}
	content, ok := r.Value.(clipboard.ClipboardContent)
	if !ok {
		return "", core.E("display.ReadClipboard", "unexpected result type", nil)
	}
	return content.Text, nil
}

func (s *Service) WriteClipboard(text string) resultFailure {
	result := s.Core().Action("clipboard.set_text").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: clipboard.TaskSetText{Text: text}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.WriteClipboard", "clipboard.set_text action failed", nil)
	}
	return nil
}

func (s *Service) HasClipboard() bool {
	text, err := s.ReadClipboard()
	return err == nil && text != ""
}

func (s *Service) ClearClipboard() resultFailure {
	result := s.Core().Action("clipboard.clear").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.ClearClipboard", "clipboard.clear action failed", nil)
	}
	return nil
}

func (s *Service) ReadClipboardImage() ([]byte, resultFailure) {
	r := s.Core().QUERY(clipboard.QueryImage{})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return nil, err
		}
		return nil, failedQuery("display.ReadClipboardImage", "clipboard.queryImage")
	}
	content, ok := r.Value.(clipboard.ImageContent)
	if !ok {
		return nil, unexpectedResultType("display.ReadClipboardImage")
	}
	if !content.HasImage {
		return nil, nil
	}
	return append([]byte(nil), content.Data...), nil
}

func (s *Service) WriteClipboardImage(data []byte) resultFailure {
	if len(data) == 0 {
		return core.E(writeClipboardImageOp, "clipboard image data is required", nil)
	}
	if len(data) > clipboard.MaxImageBytes {
		return core.E(writeClipboardImageOp, "clipboard image exceeds maximum size", nil)
	}
	result := s.Core().Action("clipboard.set_image").Run(context.Background(), core.NewOptions(
		core.Option{Key: "data", Value: append([]byte(nil), data...)},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E(writeClipboardImageOp, "clipboard.setImage action failed", nil)
	}
	return nil
}

func (s *Service) ShowNotification(opts NotificationOptions) resultFailure {
	return s.sendNotification(notification.NotificationOptions{
		ID:       opts.ID,
		Title:    opts.Title,
		Message:  opts.Message,
		Subtitle: opts.Subtitle,
	})
}

func (s *Service) ShowInfoNotification(title, message string) resultFailure {
	return s.sendNotification(notification.NotificationOptions{Title: title, Message: message})
}

func (s *Service) ShowWarningNotification(title, message string) resultFailure {
	return s.sendNotification(notification.NotificationOptions{
		Title:    title,
		Message:  message,
		Severity: notification.SeverityWarning,
	})
}

func (s *Service) ShowErrorNotification(title, message string) resultFailure {
	return s.sendNotification(notification.NotificationOptions{
		Title:    title,
		Message:  message,
		Severity: notification.SeverityError,
	})
}

func (s *Service) RequestNotificationPermission() (bool, resultFailure) {
	r := s.Core().Action("notification.request_permission").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return false, err
		}
		return false, core.E("display.RequestNotificationPermission", "notification.requestPermission action failed", nil)
	}
	granted, ok := r.Value.(bool)
	if !ok {
		return false, unexpectedResultType("display.RequestNotificationPermission")
	}
	return granted, nil
}

func (s *Service) CheckNotificationPermission() (bool, resultFailure) {
	r := s.Core().QUERY(notification.QueryPermission{})
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return false, err
		}
		return false, core.E("display.CheckNotificationPermission", "notification query failed", nil)
	}
	status, ok := r.Value.(notification.PermissionStatus)
	if !ok {
		return false, unexpectedResultType("display.CheckNotificationPermission")
	}
	return status.Granted, nil
}

func (s *Service) ClearNotifications() resultFailure {
	result := s.Core().Action("notification.clear").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskClear{}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.ClearNotifications", "notification.clear action failed", nil)
	}
	return nil
}

func (s *Service) SetTheme(theme string) resultFailure {
	result := s.Core().Action("environment.set_theme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: environment.TaskSetTheme{Theme: theme}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.SetTheme", "environment.set_theme action failed", nil)
	}
	return nil
}

func (s *Service) GetTheme() *Theme {
	r := s.Core().QUERY(environment.QueryTheme{})
	info, ok := themeInfoFromQueryResult(s, "display.GetTheme", r)
	if !ok {
		return nil
	}
	return &Theme{IsDark: info.IsDark}
}

func (s *Service) GetSystemTheme() string {
	r := s.Core().QUERY(environment.QueryTheme{})
	info, ok := themeInfoFromQueryResult(s, "display.GetSystemTheme", r)
	if !ok {
		return ""
	}
	return info.Theme
}

func themeInfoFromQueryResult(s *Service, method string, r core.Result) (environment.ThemeInfo, bool) {
	if !r.OK {
		return environment.ThemeInfo{}, false
	}
	info, ok := r.Value.(environment.ThemeInfo)
	if ok {
		return info, true
	}
	if s != nil && s.Core() != nil {
		s.Core().LogWarn(core.Errorf("query=environment.QueryTheme value_type=%T", r.Value), method, "malformed theme query payload")
	}
	return environment.ThemeInfo{}, false
}

func (s *Service) sendNotification(opts notification.NotificationOptions) resultFailure {
	result := s.Core().Action("notification.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskSend{Options: opts}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return core.E("display.sendNotification", "notification.send action failed", nil)
	}
	return nil
}

func screenToDisplay(scr *screen.Screen) *Screen {
	if scr == nil {
		return nil
	}
	return &Screen{
		ID:          scr.ID,
		Name:        scr.Name,
		X:           scr.Bounds.X,
		Y:           scr.Bounds.Y,
		Width:       scr.Bounds.Width,
		Height:      scr.Bounds.Height,
		ScaleFactor: scr.ScaleFactor,
		IsPrimary:   scr.IsPrimary,
	}
}

func toDialogOpenFileOptions(opts OpenFileOptions) dialog.OpenFileOptions {
	return dialog.OpenFileOptions{
		Title:         opts.Title,
		Directory:     opts.DefaultDirectory,
		Filename:      opts.DefaultFilename,
		AllowMultiple: opts.AllowMultiple,
		Filters:       toDialogFileFilters(opts.Filters),
	}
}

func toDialogSaveFileOptions(opts SaveFileOptions) dialog.SaveFileOptions {
	return dialog.SaveFileOptions{
		Title:     opts.Title,
		Directory: opts.DefaultDirectory,
		Filename:  opts.DefaultFilename,
		Filters:   toDialogFileFilters(opts.Filters),
	}
}

func toDialogOpenDirectoryOptions(opts OpenDirectoryOptions) dialog.OpenDirectoryOptions {
	return dialog.OpenDirectoryOptions{
		Title:     opts.Title,
		Directory: opts.DefaultDirectory,
	}
}

func toDialogFileFilters(filters []FileFilter) []dialog.FileFilter {
	if len(filters) == 0 {
		return nil
	}
	result := make([]dialog.FileFilter, 0, len(filters))
	for _, filter := range filters {
		result = append(result, dialog.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		})
	}
	return result
}

func trayMenuItemsToSystray(items []TrayMenuItem) []systray.TrayMenuItem {
	if len(items) == 0 {
		return nil
	}
	result := make([]systray.TrayMenuItem, 0, len(items))
	for _, item := range items {
		converted := systray.TrayMenuItem{
			Label:    item.Label,
			ActionID: item.ActionID,
		}
		if item.IsSeparator {
			converted.Type = "separator"
		}
		if len(item.Children) > 0 {
			converted.Submenu = trayMenuItemsToSystray(item.Children)
		}
		result = append(result, converted)
	}
	return result
}
