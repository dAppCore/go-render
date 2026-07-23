// pkg/dialog/platform.go
package dialog

// Platform abstracts the native dialog backend.
type Platform interface {
	OpenFile(options OpenFileOptions) ([]string, resultFailure)
	SaveFile(options SaveFileOptions) (string, resultFailure)
	OpenDirectory(options OpenDirectoryOptions) (string, resultFailure)
	MessageDialog(options MessageDialogOptions) (string, resultFailure)
}

// DialogType represents the type of message dialog.
type DialogType int

const (
	DialogInfo DialogType = iota
	DialogWarning
	DialogError
	DialogQuestion
)

// OpenFileOptions contains options for the open file dialog.
//
//	opts := OpenFileOptions{Title: "Select image", Filters: []FileFilter{{DisplayName: "Images", Pattern: "*.png;*.jpg"}}}
type OpenFileOptions struct {
	Title                string       `json:"title,omitempty"`
	Directory            string       `json:"directory,omitempty"`
	Filename             string       `json:"filename,omitempty"`
	Filters              []FileFilter `json:"filters,omitempty"`
	AllowMultiple        bool         `json:"allowMultiple,omitempty"`
	CanChooseDirectories bool         `json:"canChooseDirectories,omitempty"`
	CanChooseFiles       bool         `json:"canChooseFiles,omitempty"`
	ShowHiddenFiles      bool         `json:"showHiddenFiles,omitempty"`
}

// SaveFileOptions contains options for the save file dialog.
//
//	opts := SaveFileOptions{Title: "Export", Filename: "report.pdf", ShowHiddenFiles: false}
type SaveFileOptions struct {
	Title           string       `json:"title,omitempty"`
	Directory       string       `json:"directory,omitempty"`
	Filename        string       `json:"filename,omitempty"`
	Filters         []FileFilter `json:"filters,omitempty"`
	ShowHiddenFiles bool         `json:"showHiddenFiles,omitempty"`
}

// OpenDirectoryOptions contains options for the directory picker.
//
//	opts := OpenDirectoryOptions{Title: "Choose folder", ShowHiddenFiles: true}
type OpenDirectoryOptions struct {
	Title           string `json:"title,omitempty"`
	Directory       string `json:"directory,omitempty"`
	AllowMultiple   bool   `json:"allowMultiple,omitempty"`
	ShowHiddenFiles bool   `json:"showHiddenFiles,omitempty"`
}

// MessageDialogOptions contains options for a message dialog.
type MessageDialogOptions struct {
	Type    DialogType `json:"type"`
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Buttons []string   `json:"buttons,omitempty"`
}

// FileFilter represents a file type filter for dialogs.
type FileFilter struct {
	DisplayName string   `json:"displayName"`
	Pattern     string   `json:"pattern"`
	Extensions  []string `json:"extensions,omitempty"`
}
