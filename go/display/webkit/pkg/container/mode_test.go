package container

import core "dappco.re/go"

func TestDetectModeWithEnvironment(t *core.T) {
	mode := DetectModeWithEnvironment(ModeEnvironment{
		Args: []string{"--mode=worker"},
	})
	if mode != ModeWorker {
		t.Fatalf("expected worker mode, got %q", mode)
	}

	mode = DetectModeWithEnvironment(ModeEnvironment{
		LookupEnv: func(key string) (string, bool) {
			if key == "CORE_GUI_MODE" {
				return "manager", true
			}
			return "", false
		},
	})
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}

	mode = DetectModeWithEnvironment(ModeEnvironment{
		LookupEnv: func(key string) (string, bool) {
			if key == "CORE_GUI_MODE" {
				return "worker", true
			}
			return "", false
		},
	})
	if mode != ModeWorker {
		t.Fatalf("expected worker mode from environment, got %q", mode)
	}
}

func TestMode_DetectMode_Good(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:good"
	core.AssertContains(t, ax7Variant, "good")
	oldArgsFunc := argsFunc
	argsFunc = func() []string { return []string{"core-gui", "--mode=worker"} }
	t.Cleanup(func() { argsFunc = oldArgsFunc })
	t.Setenv("CORE_GUI_MODE", "manager")

	mode := DetectMode()
	if mode != ModeWorker {
		t.Fatalf("expected worker mode, got %q", mode)
	}
}

func TestMode_DetectMode_Bad(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:bad"
	core.AssertContains(t, ax7Variant, "bad")
	oldArgsFunc := argsFunc
	argsFunc = func() []string { return []string{"core-gui"} }
	t.Cleanup(func() { argsFunc = oldArgsFunc })
	t.Setenv("CORE_GUI_MODE", "bogus")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}
}

func TestMode_DetectMode_Ugly(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	oldArgsFunc := argsFunc
	argsFunc = func() []string { return []string{"core-gui", "--unexpected=flag"} }
	t.Cleanup(func() { argsFunc = oldArgsFunc })
	t.Setenv("CORE_GUI_MODE", " \tbogus\n")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode after malformed input, got %q", mode)
	}
}

// AX7 generated source-matching smoke coverage.
func TestMode_DetectModeWithEnvironment_Good(t *core.T) {
	// DetectModeWithEnvironment
	ax7Variant := "DetectModeWithEnvironment:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := DetectModeWithEnvironment(*new(ModeEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMode_DetectModeWithEnvironment_Bad(t *core.T) {
	// DetectModeWithEnvironment
	ax7Variant := "DetectModeWithEnvironment:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := DetectModeWithEnvironment(*new(ModeEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMode_DetectModeWithEnvironment_Ugly(t *core.T) {
	// DetectModeWithEnvironment
	ax7Variant := "DetectModeWithEnvironment:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := DetectModeWithEnvironment(*new(ModeEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
