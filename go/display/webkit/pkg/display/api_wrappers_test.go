package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"dappco.re/go/render/display/webkit/pkg/systray"
)

type errorOnlyWrapperCase struct {
	name      string
	action    string
	call      func(*Service) resultFailure
	setupGood func(*core.T, *core.Core)
}

func runErrorOnlyWrapperCase(t *core.T, tc errorOnlyWrapperCase) {
	t.Helper()

	t.Run("good", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		if tc.setupGood != nil {
			tc.setupGood(t, c)
		}
		core.RequireNoError(t, tc.call(svc))
	})

	t.Run("bad", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: core.AnError, OK: false}
		})

		err := tc.call(svc)
		core.AssertError(t, err)
		core.AssertEqual(t, core.AnError, err)
	})

	t.Run("ugly", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: false}
		})

		err := tc.call(svc)
		core.AssertError(t, err)
		core.AssertContains(t, err.Error(), tc.action)
	})
}

func TestDisplayAPI_TrayWrappers(t *core.T) {
	cases := []errorOnlyWrapperCase{
		{
			name:   "SetTrayTooltip",
			action: "systray.set_tooltip",
			call: func(svc *Service) resultFailure {
				return svc.SetTrayTooltip("Helper tooltip")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("systray.set_tooltip", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayTooltip)
					core.AssertEqual(t, "Helper tooltip", task.Tooltip)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SetTrayLabel",
			action: "systray.set_label",
			call: func(svc *Service) resultFailure {
				return svc.SetTrayLabel("Launcher")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("systray.set_label", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayLabel)
					core.AssertEqual(t, "Launcher", task.Label)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SetTrayMenu",
			action: "systray.set_menu",
			call: func(svc *Service) resultFailure {
				return svc.SetTrayMenu([]TrayMenuItem{
					{Label: "Open", ActionID: "open"},
					{
						Label:    "More",
						ActionID: "more",
						Children: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}},
					},
				})
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("systray.set_menu", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayMenu)
					core.AssertLen(t, task.Items, 2)
					core.AssertEqual(t, "Open", task.Items[0].Label)
					core.AssertLen(t, task.Items[1].Submenu, 1)
					core.AssertEqual(t, "nested", task.Items[1].Submenu[0].ActionID)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowTrayMessage",
			action: "systray.show_message",
			call: func(svc *Service) resultFailure {
				return svc.ShowTrayMessage("Status", "Task complete")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("systray.show_message", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskShowMessage)
					core.AssertEqual(t, "Status", task.Title)
					core.AssertEqual(t, "Task complete", task.Message)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}
}

func TestDisplayAPI_ClipboardWrappers(t *core.T) {
	t.Run("ClearClipboard", func(t *core.T) {
		runErrorOnlyWrapperCase(t, errorOnlyWrapperCase{
			name:   "ClearClipboard",
			action: "clipboard.clear",
			call: func(svc *Service) resultFailure {
				return svc.ClearClipboard()
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("clipboard.clear", func(_ context.Context, opts core.Options) core.Result {
					core.AssertEqual(t, 0, opts.Len())
					return core.Result{OK: true}
				})
			},
		})
	})

	t.Run("HasClipboard", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{
					Value: clipboard.ClipboardContent{
						Text:       "present",
						HasContent: true,
					},
					OK: true,
				}
			default:
				return core.Result{}
			}
		})
		core.AssertTrue(t, svc.HasClipboard())

		svc, c = newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{
					Value: clipboard.ClipboardContent{
						Text:       "",
						HasContent: false,
					},
					OK: true,
				}
			default:
				return core.Result{}
			}
		})
		core.AssertFalse(t, svc.HasClipboard())

		svc, c = newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{OK: false}
			default:
				return core.Result{}
			}
		})
		core.AssertFalse(t, svc.HasClipboard())
	})
}

func TestDisplayAPI_NotificationWrappers(t *core.T) {
	cases := []errorOnlyWrapperCase{
		{
			name:   "ShowNotification",
			action: "notification.send",
			call: func(svc *Service) resultFailure {
				return svc.ShowNotification(NotificationOptions{
					ID:       "alert-1",
					Title:    "Deploy",
					Message:  "Done",
					Subtitle: "CI",
				})
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					core.AssertEqual(t, notification.NotificationOptions{
						ID:       "alert-1",
						Title:    "Deploy",
						Message:  "Done",
						Subtitle: "CI",
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowInfoNotification",
			action: "notification.send",
			call: func(svc *Service) resultFailure {
				return svc.ShowInfoNotification("Info", "Ready")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					core.AssertEqual(t, notification.NotificationOptions{
						Title:   "Info",
						Message: "Ready",
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowWarningNotification",
			action: "notification.send",
			call: func(svc *Service) resultFailure {
				return svc.ShowWarningNotification("Warn", "Careful")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					core.AssertEqual(t, notification.NotificationOptions{
						Title:    "Warn",
						Message:  "Careful",
						Severity: notification.SeverityWarning,
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowErrorNotification",
			action: "notification.send",
			call: func(svc *Service) resultFailure {
				return svc.ShowErrorNotification("Error", "Failed")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					core.AssertEqual(t, notification.NotificationOptions{
						Title:    "Error",
						Message:  "Failed",
						Severity: notification.SeverityError,
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ClearNotifications",
			action: "notification.clear",
			call: func(svc *Service) resultFailure {
				return svc.ClearNotifications()
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("notification.clear", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskClear)
					core.AssertEmpty(t, task.ID)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}
}

func TestDisplayAPI_ThemeWrapper(t *core.T) {
	runErrorOnlyWrapperCase(t, errorOnlyWrapperCase{
		name:   "SetTheme",
		action: "environment.set_theme",
		call: func(svc *Service) resultFailure {
			return svc.SetTheme("system")
		},
		setupGood: func(t *core.T, c *core.Core) {
			t.Helper()
			c.Action("environment.set_theme", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(environment.TaskSetTheme)
				core.AssertEqual(t, "system", task.Theme)
				return core.Result{OK: true}
			})
		},
	})
}

func TestDisplayAPI_GetTrayInfo(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{
				Value: map[string]any{
					"tooltip": "Ready",
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	info := svc.GetTrayInfo()
	core.AssertNotNil(t, info)
	core.AssertEqual(t, "Ready", info["tooltip"])

	svc, c = newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})
	core.AssertNil(t, svc.GetTrayInfo())

	svc, c = newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})
	core.AssertNil(t, svc.GetTrayInfo())
}
