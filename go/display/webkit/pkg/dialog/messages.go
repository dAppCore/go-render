package dialog

// TaskOpenFile presents an open-file dialog with the given options.
//
//	result, _, err := c.PERFORM(dialog.TaskOpenFile{Options: dialog.OpenFileOptions{Title: "Pick file"}})
//	paths := result.([]string)
type TaskOpenFile struct{ Options OpenFileOptions }

// TaskOpenFileWithOptions presents an open-file dialog pre-configured from an options struct.
// Equivalent to TaskOpenFile but mirrors the stub DialogManager.OpenFileWithOptions API.
//
//	result, _, err := c.PERFORM(dialog.TaskOpenFileWithOptions{Options: &dialog.OpenFileOptions{Title: "Select log", AllowMultiple: true}})
type TaskOpenFileWithOptions struct{ Options *OpenFileOptions }

// TaskSaveFile presents a save-file dialog with the given options.
//
//	result, _, err := c.PERFORM(dialog.TaskSaveFile{Options: dialog.SaveFileOptions{Filename: "report.csv"}})
//	path := result.(string)
type TaskSaveFile struct{ Options SaveFileOptions }

// TaskSaveFileWithOptions presents a save-file dialog pre-configured from an options struct.
// Equivalent to TaskSaveFile but mirrors the stub DialogManager.SaveFileWithOptions API.
//
//	result, _, err := c.PERFORM(dialog.TaskSaveFileWithOptions{Options: &dialog.SaveFileOptions{Title: "Export data"}})
type TaskSaveFileWithOptions struct{ Options *SaveFileOptions }

// TaskOpenDirectory presents a directory picker dialog.
//
//	result, _, err := c.PERFORM(dialog.TaskOpenDirectory{Options: dialog.OpenDirectoryOptions{Title: "Choose folder"}})
//	path := result.(string)
type TaskOpenDirectory struct{ Options OpenDirectoryOptions }

// TaskMessageDialog presents a message dialog of the given type.
//
//	result, _, err := c.PERFORM(dialog.TaskMessageDialog{Options: dialog.MessageDialogOptions{Type: dialog.DialogQuestion, Title: "Confirm", Message: "Delete?", Buttons: []string{"Yes", "No"}}})
//	clicked := result.(string)
type TaskMessageDialog struct{ Options MessageDialogOptions }

// TaskInfo presents an information message dialog.
//
//	result, _, err := c.PERFORM(dialog.TaskInfo{Title: "Done", Message: "File saved successfully."})
//	clicked := result.(string)
type TaskInfo struct {
	Title   string
	Message string
	Buttons []string
}

// TaskQuestion presents a question message dialog.
//
//	result, _, err := c.PERFORM(dialog.TaskQuestion{Title: "Confirm", Message: "Delete file?", Buttons: []string{"Yes", "No"}})
//	if result.(string) == "Yes" { deleteFile() }
type TaskQuestion struct {
	Title   string
	Message string
	Buttons []string
}

// TaskWarning presents a warning message dialog.
//
//	result, _, err := c.PERFORM(dialog.TaskWarning{Title: "Low disk", Message: "Disk space is critically low."})
type TaskWarning struct {
	Title   string
	Message string
	Buttons []string
}

// TaskError presents an error message dialog.
//
//	result, _, err := c.PERFORM(dialog.TaskError{Title: "Operation failed", Message: err.Error()})
type TaskError struct {
	Title   string
	Message string
	Buttons []string
}

// TaskPrompt presents a text input prompt in the active application window.
//
//	result, _, err := c.PERFORM(dialog.TaskPrompt{Title: "Rename", Message: "Enter a new name", DefaultValue: "draft"})
//	value := result.(dialog.PromptResult)
type TaskPrompt struct {
	Title        string
	Message      string
	DefaultValue string
}

// PromptResult is the value returned by TaskPrompt.
//
//	result, _, err := c.PERFORM(dialog.TaskPrompt{Title: "Search", DefaultValue: "core"})
//	if result.(PromptResult).Confirmed { use(result.(PromptResult).Value) }
type PromptResult struct {
	Value     string
	Confirmed bool
}
