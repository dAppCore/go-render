// pkg/systray/messages.go
package systray

type QueryConfig struct{}

type QueryInfo struct{}

type TaskSetTrayIcon struct{ Data []byte }

// TaskSetTrayTemplateIcon updates the tray icon as a macOS template
// image (auto-inverted by the OS to match light/dark menu bar). On
// non-darwin backends it falls back to a regular icon set.
//
//	c.Action("systray.set_template_icon").Run(ctx, NewOptions(
//	    Option{Key: "task", Value: TaskSetTrayTemplateIcon{Data: iconBytes}},
//	))
type TaskSetTrayTemplateIcon struct{ Data []byte }

type TaskSetTrayTooltip struct{ Tooltip string }

type TaskSetTrayLabel struct{ Label string }

type TaskSetTrayMenu struct{ Items []TrayMenuItem }

type TaskShowMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type TaskShowPanel struct{}

type TaskHidePanel struct{}

type TaskSaveConfig struct{ Config map[string]any }

// TaskAttachWindow anchors a previously-created window (by Name) to the
// tray. OffsetY is in platform pixels along the tray's natural axis;
// OffsetX is reserved for backends that expose two-axis offsets
// (Wails 3 alpha.91 uses a single offset only).
//
//	c.Action("systray.attach_window").Run(ctx, NewOptions(
//	    Option{Key: "task", Value: TaskAttachWindow{Name: "tray", OffsetY: 5}},
//	))
type TaskAttachWindow struct {
	Name    string
	OffsetX int
	OffsetY int
}

type ActionTrayClicked struct{}

type ActionTrayMenuItemClicked struct{ ActionID string }
