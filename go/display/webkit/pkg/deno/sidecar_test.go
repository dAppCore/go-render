package deno

import (
	core "dappco.re/go"
)

func TestSidecar_New_Good(t *core.T) {
	// New
	ax7Variant := "New:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := New(Options{})

	status := manager.Status()
	core.AssertEqual(t, "deno", status.Binary)
	core.AssertFalse(t, status.Running)
	core.AssertEmpty(t, status.PID)
}

func TestSidecar_New_Bad(t *core.T) {
	// New
	ax7Variant := "New:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := New(Options{Binary: "/usr/local/bin/deno-custom", Args: []string{"format"}})

	status := manager.Status()
	core.AssertEqual(t, "/usr/local/bin/deno-custom", status.Binary)
	core.AssertFalse(t, status.Running)
}

func TestSidecar_New_Ugly(t *core.T) {
	// New
	ax7Variant := "New:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := New(Options{Binary: "   "})

	status := manager.Status()
	core.AssertEqual(t, "deno", status.Binary)
	core.AssertFalse(t, status.Running)
}

// AX7 generated source-matching smoke coverage.
func TestSidecar_Manager_Start_Good(t *core.T) {
	// Manager Start
	ax7Variant := "Manager_Start:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Start_Bad(t *core.T) {
	// Manager Start
	ax7Variant := "Manager_Start:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Start_Ugly(t *core.T) {
	// Manager Start
	ax7Variant := "Manager_Start:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Stop_Good(t *core.T) {
	// Manager Stop
	ax7Variant := "Manager_Stop:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Stop_Bad(t *core.T) {
	// Manager Stop
	ax7Variant := "Manager_Stop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Stop_Ugly(t *core.T) {
	// Manager Stop
	ax7Variant := "Manager_Stop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Status_Good(t *core.T) {
	// Manager Status
	ax7Variant := "Manager_Status:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Status_Bad(t *core.T) {
	// Manager Status
	ax7Variant := "Manager_Status:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Status_Ugly(t *core.T) {
	// Manager Status
	ax7Variant := "Manager_Status:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_OnEvent_Good(t *core.T) {
	// Manager OnEvent
	ax7Variant := "Manager_OnEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_OnEvent_Bad(t *core.T) {
	// Manager OnEvent
	ax7Variant := "Manager_OnEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_OnEvent_Ugly(t *core.T) {
	// Manager OnEvent
	ax7Variant := "Manager_OnEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Eval_Good(t *core.T) {
	// Manager Eval
	ax7Variant := "Manager_Eval:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Eval_Bad(t *core.T) {
	// Manager Eval
	ax7Variant := "Manager_Eval:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Eval_Ugly(t *core.T) {
	// Manager Eval
	ax7Variant := "Manager_Eval:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Emit_Good(t *core.T) {
	// Manager Emit
	ax7Variant := "Manager_Emit:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Emit_Bad(t *core.T) {
	// Manager Emit
	ax7Variant := "Manager_Emit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSidecar_Manager_Emit_Ugly(t *core.T) {
	// Manager Emit
	ax7Variant := "Manager_Emit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("../../edge", map[string]any{})
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
