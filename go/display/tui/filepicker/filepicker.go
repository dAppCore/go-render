// SPDX-Licence-Identifier: EUPL-1.2

// Package filepicker re-exports charmbracelet/bubbles/filepicker through
// go-html's tui seam — an interactive file/directory browser. Swap the
// import path (bubbles/filepicker → html/tui/filepicker) and keep every
// filepicker.Model / filepicker.New reference unchanged. Init, Update, View,
// SetHeight, Height, DidSelectFile, DidSelectDisabledFile and
// HighlightedPath are Model methods, not package functions, so they need no
// re-export here — they come along for free since Model is a genuine alias.
// There is no Option pattern: a Model is configured by setting its exported
// fields directly (CurrentDirectory, AllowedTypes, ShowHidden, DirAllowed,
// FileAllowed, AutoHeight, Cursor, KeyMap, Styles) before Init runs.
// IsHidden is the dotfile/hidden-attribute check the picker's own internal
// directory read applies whenever ShowHidden is false.
package filepicker

import "charm.land/bubbles/v2/filepicker"

// Model is the file picker itself. KeyMap and Styles are its two
// customisation seams, each with a matching Default* constructor below.
type (
	Model  = filepicker.Model
	KeyMap = filepicker.KeyMap
	Styles = filepicker.Styles
)

var (
	New           = filepicker.New
	DefaultKeyMap = filepicker.DefaultKeyMap
	DefaultStyles = filepicker.DefaultStyles
	IsHidden      = filepicker.IsHidden
)
