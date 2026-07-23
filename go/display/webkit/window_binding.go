// SPDX-License-Identifier: EUPL-1.2

package webkit

import (
	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/window"
)

// WindowBindingService is the canonical Wails IPC binding for
// frontend-driven window verbs (Open / Hide / List / SetSize). Add it
// to GuiConfig.Bindings via webkit.Bind(webkit.NewWindowBindingService(c))
// and frontend code calls:
//
//	import * as WindowService from "@gui/windowbindingservice";
//	WindowService.Open("chat");
//	WindowService.Hide("chat");
//	WindowService.List();
//	WindowService.SetSize("chat", 900, 700);
//
// Internally the verbs route to webkit.OpenWindow / webkit.HideWindow /
// the GuiConfig.WindowRegistry / window.set_size action — every consumer
// gets the same dispatch path. Replaces the per-app WindowService that
// each consumer would otherwise hand-roll.
type WindowBindingService struct {
	core *core.Core
}

// NewWindowBindingService builds a WindowBindingService bound to the
// supplied Core. Pass the result via webkit.Bind into GuiConfig.Bindings.
func NewWindowBindingService(c *core.Core) *WindowBindingService {
	return &WindowBindingService{core: c}
}

// ServiceName lets Wails name the binding "Window" regardless of the
// Go struct name. Frontend bindings land at @gui/window.
func (s *WindowBindingService) ServiceName() string { return "Window" }

// ServiceStartup is the Wails lifecycle no-op. The Core is captured at
// construction; nothing else to wire here.
func (s *WindowBindingService) ServiceStartup(context, _ any) core.Result {
	return core.Ok(nil)
}

// ServiceShutdown is the Wails lifecycle no-op.
func (s *WindowBindingService) ServiceShutdown() core.Result { return core.Ok(nil) }

// Open shows + focuses the named window. Wraps webkit.OpenWindow. Returns
// a Fail Result when the window is not in the registry so the
// frontend can surface a clear error.
func (s *WindowBindingService) Open(name string) core.Result {
	if s == nil || s.core == nil {
		return core.Fail(core.NewError("gui: window binding not bound to a Core"))
	}
	if !OpenWindow(s.core, name) {
		return core.Fail(core.NewError("gui: no registered window named: " + name))
	}
	return core.Ok(nil)
}

// Hide hides the named window. Wraps webkit.HideWindow. Returns a Fail
// Result on unregistered name.
func (s *WindowBindingService) Hide(name string) core.Result {
	if s == nil || s.core == nil {
		return core.Fail(core.NewError("gui: window binding not bound to a Core"))
	}
	if !HideWindow(s.core, name) {
		return core.Fail(core.NewError("gui: no registered window named: " + name))
	}
	return core.Ok(nil)
}

// List returns the names of every registered window. Frontend can
// render a switcher / jump-list from this. Walks GuiConfig.WindowRegistry
// directly via the webkit.Service lookup.
func (s *WindowBindingService) List() core.Result {
	if s == nil || s.core == nil {
		return core.Fail(core.NewError("gui: window binding not bound to a Core"))
	}
	svc, ok := core.ServiceFor[*Service](s.core, "gui")
	if !ok || svc == nil {
		return core.Ok([]string{})
	}
	registry := svc.Options().WindowRegistry
	names := make([]string, 0, len(registry))
	for _, w := range registry {
		if w == nil || w.Name == "" {
			continue
		}
		names = append(names, w.Name)
	}
	return core.Ok(names)
}

// SetSize resizes the named window. Element-driven sizing — Lit
// components call this from firstUpdated() with their declared w/h so
// the WebView's content owns the dimensions rather than the Go boot-
// time spec being the only source of truth.
func (s *WindowBindingService) SetSize(name string, width, height int) core.Result {
	if s == nil || s.core == nil {
		return core.Fail(core.NewError("gui: window binding not bound to a Core"))
	}
	if name == "" {
		return core.Fail(core.NewError("gui: SetSize requires a window name"))
	}
	if width <= 0 || height <= 0 {
		return core.Fail(core.NewError("gui: SetSize requires positive dimensions"))
	}
	return s.core.Action("window.set_size").Run(core.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetSize{Name: name, Width: width, Height: height}},
	))
}
