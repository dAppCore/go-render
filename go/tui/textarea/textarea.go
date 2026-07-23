// Package textarea re-exports charmbracelet/bubbles/textarea through go-html's
// tui seam — a multi-line text editor. Swap the import path
// (bubbles/textarea → html/tui/textarea) and keep every textarea.Model /
// textarea.New reference unchanged.
package textarea

import "github.com/charmbracelet/bubbles/textarea"

type Model = textarea.Model

var New = textarea.New
