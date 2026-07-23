// pkg/mcp/tools_dialog.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- dialog_open_file ---

type DialogOpenFileInput struct {
	Title                string              `json:"title,omitempty"`
	Directory            string              `json:"directory,omitempty"`
	Filters              []dialog.FileFilter `json:"filters,omitempty"`
	AllowMultiple        bool                `json:"allowMultiple,omitempty"`
	CanChooseDirectories bool                `json:"canChooseDirectories,omitempty"`
	CanChooseFiles       bool                `json:"canChooseFiles,omitempty"`
	ShowHiddenFiles      bool                `json:"showHiddenFiles,omitempty"`
}
type DialogOpenFileOutput struct {
	Paths []string `json:"paths"`
}

func (s *Subsystem) dialogOpenFile(_ context.Context, _ *mcp.CallToolRequest, input DialogOpenFileInput) (*mcp.CallToolResult, DialogOpenFileOutput, resultFailure) {
	r := s.core.Action("dialog.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskOpenFile{Options: dialog.OpenFileOptions{
			Title:                input.Title,
			Directory:            input.Directory,
			Filters:              input.Filters,
			AllowMultiple:        input.AllowMultiple,
			CanChooseDirectories: input.CanChooseDirectories,
			CanChooseFiles:       input.CanChooseFiles,
			ShowHiddenFiles:      input.ShowHiddenFiles,
		}}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogOpenFileOutput{}, e
		}
		return nil, DialogOpenFileOutput{}, nil
	}
	paths, ok := r.Value.([]string)
	if !ok {
		return nil, DialogOpenFileOutput{}, core.E("mcp.dialogOpenFile", "unexpected result type", nil)
	}
	return nil, DialogOpenFileOutput{Paths: paths}, nil
}

// --- dialog_save_file ---

type DialogSaveFileInput struct {
	Title           string              `json:"title,omitempty"`
	Directory       string              `json:"directory,omitempty"`
	Filename        string              `json:"filename,omitempty"`
	Filters         []dialog.FileFilter `json:"filters,omitempty"`
	ShowHiddenFiles bool                `json:"showHiddenFiles,omitempty"`
}
type DialogSaveFileOutput struct {
	Path string `json:"path,omitempty"`
}

func (s *Subsystem) dialogSaveFile(_ context.Context, _ *mcp.CallToolRequest, input DialogSaveFileInput) (*mcp.CallToolResult, DialogSaveFileOutput, resultFailure) {
	r := s.core.Action("dialog.save_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskSaveFile{Options: dialog.SaveFileOptions{
			Title:           input.Title,
			Directory:       input.Directory,
			Filename:        input.Filename,
			Filters:         input.Filters,
			ShowHiddenFiles: input.ShowHiddenFiles,
		}}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogSaveFileOutput{}, e
		}
		return nil, DialogSaveFileOutput{}, nil
	}
	path, ok := r.Value.(string)
	if !ok {
		return nil, DialogSaveFileOutput{}, core.E("mcp.dialogSaveFile", "unexpected result type", nil)
	}
	return nil, DialogSaveFileOutput{Path: path}, nil
}

// --- dialog_open_directory ---

type DialogOpenDirectoryInput struct {
	Title           string `json:"title,omitempty"`
	Directory       string `json:"directory,omitempty"`
	ShowHiddenFiles bool   `json:"showHiddenFiles,omitempty"`
}
type DialogOpenDirectoryOutput struct {
	Path string `json:"path,omitempty"`
}

func (s *Subsystem) dialogOpenDirectory(_ context.Context, _ *mcp.CallToolRequest, input DialogOpenDirectoryInput) (*mcp.CallToolResult, DialogOpenDirectoryOutput, resultFailure) {
	r := s.core.Action("dialog.open_directory").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskOpenDirectory{Options: dialog.OpenDirectoryOptions{
			Title:           input.Title,
			Directory:       input.Directory,
			ShowHiddenFiles: input.ShowHiddenFiles,
		}}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogOpenDirectoryOutput{}, e
		}
		return nil, DialogOpenDirectoryOutput{}, nil
	}
	path, ok := r.Value.(string)
	if !ok {
		return nil, DialogOpenDirectoryOutput{}, core.E("mcp.dialogOpenDirectory", "unexpected result type", nil)
	}
	return nil, DialogOpenDirectoryOutput{Path: path}, nil
}

// --- dialog_confirm ---

type DialogConfirmInput struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogConfirmOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogConfirm(_ context.Context, _ *mcp.CallToolRequest, input DialogConfirmInput) (*mcp.CallToolResult, DialogConfirmOutput, resultFailure) {
	r := s.core.Action("dialog.question").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskQuestion{
			Title:   input.Title,
			Message: input.Message,
			Buttons: input.Buttons,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogConfirmOutput{}, e
		}
		return nil, DialogConfirmOutput{}, nil
	}
	button, ok := r.Value.(string)
	if !ok {
		return nil, DialogConfirmOutput{}, core.E("mcp.dialogConfirm", "unexpected result type", nil)
	}
	return nil, DialogConfirmOutput{Button: button}, nil
}

// --- dialog_message ---

type DialogMessageInput struct {
	Type    string   `json:"type,omitempty"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogMessageOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogMessage(_ context.Context, _ *mcp.CallToolRequest, input DialogMessageInput) (*mcp.CallToolResult, DialogMessageOutput, resultFailure) {
	dialogType := dialog.DialogInfo
	switch input.Type {
	case "", "info":
		dialogType = dialog.DialogInfo
	case "question":
		dialogType = dialog.DialogQuestion
	case "warning":
		dialogType = dialog.DialogWarning
	case "error":
		dialogType = dialog.DialogError
	default:
		return nil, DialogMessageOutput{}, core.E("mcp.dialogMessage", "invalid dialog type: "+input.Type, nil)
	}

	r := s.core.Action("dialog.message").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskMessageDialog{
			Options: dialog.MessageDialogOptions{
				Type:    dialogType,
				Title:   input.Title,
				Message: input.Message,
				Buttons: input.Buttons,
			},
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogMessageOutput{}, e
		}
		return nil, DialogMessageOutput{}, nil
	}
	button, ok := r.Value.(string)
	if !ok {
		return nil, DialogMessageOutput{}, core.E("mcp.dialogMessage", "unexpected result type", nil)
	}
	return nil, DialogMessageOutput{Button: button}, nil
}

// --- dialog_prompt ---

type DialogPromptInput struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	DefaultValue string `json:"defaultValue,omitempty"`
}
type DialogPromptOutput struct {
	Value     string `json:"value"`
	Confirmed bool   `json:"confirmed"`
}

func (s *Subsystem) dialogPrompt(_ context.Context, _ *mcp.CallToolRequest, input DialogPromptInput) (*mcp.CallToolResult, DialogPromptOutput, resultFailure) {
	r := s.core.Action("dialog.prompt").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskPrompt{
			Title:        input.Title,
			Message:      input.Message,
			DefaultValue: input.DefaultValue,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogPromptOutput{}, e
		}
		return nil, DialogPromptOutput{}, nil
	}
	result, ok := r.Value.(dialog.PromptResult)
	if !ok {
		return nil, DialogPromptOutput{}, core.E("mcp.dialogPrompt", "unexpected result type", nil)
	}
	return nil, DialogPromptOutput{Value: result.Value, Confirmed: result.Confirmed}, nil
}

// --- dialog_info ---

type DialogInfoInput struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogInfoOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogInfo(_ context.Context, _ *mcp.CallToolRequest, input DialogInfoInput) (*mcp.CallToolResult, DialogInfoOutput, resultFailure) {
	r := s.core.Action("dialog.info").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskInfo{
			Title:   input.Title,
			Message: input.Message,
			Buttons: input.Buttons,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogInfoOutput{}, e
		}
		return nil, DialogInfoOutput{}, nil
	}
	button, ok := r.Value.(string)
	if !ok {
		return nil, DialogInfoOutput{}, core.E("mcp.dialogInfo", "unexpected result type", nil)
	}
	return nil, DialogInfoOutput{Button: button}, nil
}

// --- dialog_warning ---

type DialogWarningInput struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogWarningOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogWarning(_ context.Context, _ *mcp.CallToolRequest, input DialogWarningInput) (*mcp.CallToolResult, DialogWarningOutput, resultFailure) {
	r := s.core.Action("dialog.warning").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskWarning{
			Title:   input.Title,
			Message: input.Message,
			Buttons: input.Buttons,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogWarningOutput{}, e
		}
		return nil, DialogWarningOutput{}, nil
	}
	button, ok := r.Value.(string)
	if !ok {
		return nil, DialogWarningOutput{}, core.E("mcp.dialogWarning", "unexpected result type", nil)
	}
	return nil, DialogWarningOutput{Button: button}, nil
}

// --- dialog_error ---

type DialogErrorInput struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogErrorOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogError(_ context.Context, _ *mcp.CallToolRequest, input DialogErrorInput) (*mcp.CallToolResult, DialogErrorOutput, resultFailure) {
	r := s.core.Action("dialog.error").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskError{
			Title:   input.Title,
			Message: input.Message,
			Buttons: input.Buttons,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DialogErrorOutput{}, e
		}
		return nil, DialogErrorOutput{}, nil
	}
	button, ok := r.Value.(string)
	if !ok {
		return nil, DialogErrorOutput{}, core.E("mcp.dialogError", "unexpected result type", nil)
	}
	return nil, DialogErrorOutput{Button: button}, nil
}

// --- Registration ---

func (s *Subsystem) registerDialogTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "dialog_open_file", Description: "Show an open file dialog"}, s.dialogOpenFile)
	addTool(s, server, &mcp.Tool{Name: "dialog_save_file", Description: "Show a save file dialog"}, s.dialogSaveFile)
	addTool(s, server, &mcp.Tool{Name: "dialog_open_directory", Description: "Show a directory picker dialog"}, s.dialogOpenDirectory)
	addTool(s, server, &mcp.Tool{Name: "dialog_confirm", Description: "Show a question/confirmation dialog"}, s.dialogConfirm)
	addTool(s, server, &mcp.Tool{Name: "dialog_message", Description: "Show a message dialog"}, s.dialogMessage)
	addTool(s, server, &mcp.Tool{
		Name:        "dialog_prompt",
		Description: `Show an input prompt dialog in the active window. Example: {"title":"Rename","message":"New name","defaultValue":"project-a"}`,
	}, s.dialogPrompt)
	addTool(s, server, &mcp.Tool{Name: "dialog_info", Description: "Show an information message dialog"}, s.dialogInfo)
	addTool(s, server, &mcp.Tool{Name: "dialog_warning", Description: "Show a warning message dialog"}, s.dialogWarning)
	addTool(s, server, &mcp.Tool{Name: "dialog_error", Description: "Show an error message dialog"}, s.dialogError)
}
