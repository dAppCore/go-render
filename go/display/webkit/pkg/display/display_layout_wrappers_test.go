package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/window"
)

type layoutResultWrapperCase struct {
	name      string
	action    string
	zero      any
	call      func(*Service) (any, resultFailure)
	setupGood func(*core.T, *core.Core)
	wantGood  func(*core.T, any)
}

func runLayoutResultWrapperCase(t *core.T, tc layoutResultWrapperCase) {
	t.Helper()

	t.Run("good", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		if tc.setupGood != nil {
			tc.setupGood(t, c)
		}

		got, err := tc.call(svc)
		core.RequireNoError(t, err)
		tc.wantGood(t, got)
	})

	t.Run("bad", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: core.AnError, OK: false}
		})

		got, err := tc.call(svc)
		core.AssertError(t, err)
		core.AssertEqual(t, tc.zero, got)
		core.AssertEqual(t, core.AnError, err)
	})

	t.Run("ugly-action", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: false}
		})

		got, err := tc.call(svc)
		core.AssertError(t, err)
		core.AssertEqual(t, tc.zero, got)
		core.AssertContains(t, err.Error(), tc.action)
	})

	t.Run("ugly-type", func(t *core.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: true}
		})

		got, err := tc.call(svc)
		core.AssertError(t, err)
		core.AssertEqual(t, tc.zero, got)
		core.AssertContains(t, err.Error(), "unexpected result type")
	})
}

func TestDisplay_LayoutDelegationWrappers(t *core.T) {
	errorCases := []errorOnlyWrapperCase{
		{
			name:   "DeleteLayout",
			action: "window.delete_layout",
			call: func(svc *Service) resultFailure {
				return svc.DeleteLayout("development")
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("window.delete_layout", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskDeleteLayout)
					core.AssertEqual(t, "development", task.Name)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "TileWindows",
			action: "window.tile_windows",
			call: func(svc *Service) resultFailure {
				return svc.TileWindows(window.TileModeGrid, []string{"editor", "terminal"})
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("window.tile_windows", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskTileWindows)
					core.AssertEqual(t, window.TileModeGrid.String(), task.Mode)
					core.AssertEqual(t, []string{"editor", "terminal"}, task.Windows)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SnapWindow",
			action: "window.snap_window",
			call: func(svc *Service) resultFailure {
				return svc.SnapWindow("preview", window.SnapCenter)
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("window.snap_window", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskSnapWindow)
					core.AssertEqual(t, "preview", task.Name)
					core.AssertEqual(t, window.SnapCenter.String(), task.Position)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "StackWindows",
			action: "window.stack_windows",
			call: func(svc *Service) resultFailure {
				return svc.StackWindows([]string{"editor", "preview"}, 24, 18)
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("window.stack_windows", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskStackWindows)
					core.AssertEqual(t, []string{"editor", "preview"}, task.Windows)
					core.AssertEqual(t, 24, task.OffsetX)
					core.AssertEqual(t, 18, task.OffsetY)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ApplyWorkflowLayout",
			action: "window.apply_workflow",
			call: func(svc *Service) resultFailure {
				return svc.ApplyWorkflowLayout(window.WorkflowCoding)
			},
			setupGood: func(t *core.T, c *core.Core) {
				t.Helper()
				c.Action("window.apply_workflow", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskApplyWorkflow)
					core.AssertEqual(t, window.WorkflowCoding.String(), task.Workflow)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *core.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "LayoutBesideEditor",
		action: "window.layout_beside_editor",
		zero:   window.LayoutBesideEditorResult{},
		call: func(svc *Service) (any, resultFailure) {
			return svc.LayoutBesideEditor("preview", "code", "right", 0.62)
		},
		setupGood: func(t *core.T, c *core.Core) {
			t.Helper()
			c.Action("window.layout_beside_editor", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskLayoutBesideEditor)
				core.AssertEqual(t, "preview", task.Name)
				core.AssertEqual(t, "code", task.Editor)
				core.AssertEqual(t, "right", task.Side)
				core.AssertInDelta(t, 0.62, task.Ratio, 0.0001)
				return core.Result{
					Value: window.LayoutBesideEditorResult{
						Editor: "code",
						EditorBounds: window.WindowBounds{
							X: 10, Y: 20, Width: 640, Height: 800,
						},
						WindowBounds: window.WindowBounds{
							X: 650, Y: 20, Width: 640, Height: 800,
						},
						Side:     "right",
						ScreenID: "screen-1",
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *core.T, got any) {
			t.Helper()
			result := got.(window.LayoutBesideEditorResult)
			core.AssertEqual(t, "code", result.Editor)
			core.AssertEqual(t, "right", result.Side)
			core.AssertEqual(t, "screen-1", result.ScreenID)
		},
	})

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "FindScreenSpace",
		action: "window.find_space",
		zero:   window.ScreenSpace{},
		call: func(svc *Service) (any, resultFailure) {
			return svc.FindScreenSpace("screen-1", 800, 600, 24)
		},
		setupGood: func(t *core.T, c *core.Core) {
			t.Helper()
			c.Action("window.find_space", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskScreenFindSpace)
				core.AssertEqual(t, "screen-1", task.ScreenID)
				core.AssertEqual(t, 800, task.Width)
				core.AssertEqual(t, 600, task.Height)
				core.AssertEqual(t, 24, task.Padding)
				return core.Result{
					Value: window.ScreenSpace{
						ScreenID: "screen-1",
						X:        100,
						Y:        120,
						Width:    800,
						Height:   600,
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *core.T, got any) {
			t.Helper()
			space := got.(window.ScreenSpace)
			core.AssertEqual(t, "screen-1", space.ScreenID)
			core.AssertEqual(t, 100, space.X)
			core.AssertEqual(t, 120, space.Y)
			core.AssertEqual(t, 800, space.Width)
			core.AssertEqual(t, 600, space.Height)
		},
	})

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "ArrangeWindowPair",
		action: "window.arrange_pair",
		zero:   window.PairArrangement{},
		call: func(svc *Service) (any, resultFailure) {
			return svc.ArrangeWindowPair("editor", "preview", "screen-1", 0.55)
		},
		setupGood: func(t *core.T, c *core.Core) {
			t.Helper()
			c.Action("window.arrange_pair", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskWindowArrangePair)
				core.AssertEqual(t, "editor", task.Primary)
				core.AssertEqual(t, "preview", task.Secondary)
				core.AssertEqual(t, "screen-1", task.ScreenID)
				core.AssertInDelta(t, 0.55, task.Ratio, 0.0001)
				return core.Result{
					Value: window.PairArrangement{
						Primary: window.WindowBounds{
							X: 0, Y: 0, Width: 800, Height: 600,
						},
						Secondary: window.WindowBounds{
							X: 800, Y: 0, Width: 800, Height: 600,
						},
						Orientation: "horizontal",
						ScreenID:    "screen-1",
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *core.T, got any) {
			t.Helper()
			arrangement := got.(window.PairArrangement)
			core.AssertEqual(t, "horizontal", arrangement.Orientation)
			core.AssertEqual(t, "screen-1", arrangement.ScreenID)
			core.AssertEqual(t, 800, arrangement.Primary.Width)
			core.AssertEqual(t, 800, arrangement.Secondary.Width)
		},
	})
}
