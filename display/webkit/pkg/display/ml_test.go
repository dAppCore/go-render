package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/chat"
)

func TestML_ModelState_GoodCase(t *core.T) {
	t.Setenv("CORE_ML_API_URL", "https://ml.example.test/api/")
	svc, c := newTestDisplayService(t)
	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{
			Value: []chat.ModelEntry{
				{Name: "alpha", Loaded: true, SizeBytes: 2048, Backend: "vulkan"},
				{Name: "beta", Loaded: false, SizeBytes: 4096},
			},
			OK: true,
		}
	})

	state := svc.modelState()
	core.AssertEqual(t, "https://ml.example.test/api", state.APIURL)
	core.AssertEqual(t, 1, len(state.Loaded))
	core.AssertEqual(t, int64(2048), state.VRAMBytes)
	core.AssertEqual(t, "vulkan", state.Backend)
	core.AssertEqual(t, "https://ml.example.test/api/v1/chat/completions", state.InferenceURL)
}

func TestML_ModelState_BadCase(t *core.T) {
	t.Setenv("CORE_ML_API_URL", "")
	svc, _ := newTestDisplayService(t)

	state := svc.modelState()
	core.AssertEqual(t, "http://localhost:8090", state.APIURL)
	core.AssertEmpty(t, state.Loaded)
	core.AssertEqual(t, int64(0), state.VRAMBytes)
	core.AssertEqual(t, "local", state.Backend)
}

func TestML_ModelState_UglyCase(t *core.T) {
	svc, c := newTestDisplayService(t)
	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{
			Value: []chat.ModelEntry{
				{Name: "unloaded", Loaded: false, SizeBytes: 8192},
				{Name: "blank-backend", Loaded: true, SizeBytes: 0},
			},
			OK: true,
		}
	})

	state := svc.modelState()
	core.AssertEqual(t, int64(0), state.VRAMBytes)
	core.AssertEqual(t, "local", state.Backend)
	core.AssertEqual(t, "\"line\\nquote\\\"slash\\\\\"", quoteJS("line\nquote\"slash\\"))
	core.AssertEqual(t, int64(0), estimateVRAM(state.Available))
}
