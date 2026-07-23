package lifecycle

import core "dappco.re/go"

func TestMode_DetectMode_Good(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:good"
	core.AssertContains(t, ax7Variant, "good")
	t.Setenv(appModeEnv, "")
	t.Setenv(ciEnv, "")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}
}

func TestMode_DetectMode_Bad(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:bad"
	core.AssertContains(t, ax7Variant, "bad")
	t.Setenv(appModeEnv, "bogus")
	t.Setenv(ciEnv, "")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode after invalid env value, got %q", mode)
	}
}

func TestMode_DetectMode_Ugly(t *core.T) {
	// DetectMode
	ax7Variant := "DetectMode:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	t.Setenv(appModeEnv, "")
	t.Setenv(ciEnv, "true")

	mode := DetectMode()
	if mode != ModeWorker {
		t.Fatalf("expected worker mode in CI headless context, got %q", mode)
	}
}
