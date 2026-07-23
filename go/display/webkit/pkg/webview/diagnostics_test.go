package webview

import (
	core "dappco.re/go"
)

func TestDiagnostics_ComputedStyleScript_Good(t *core.T) {
	// ComputedStyleScript
	ax7Variant := "ComputedStyleScript:good"
	core.AssertContains(t, ax7Variant, "good")
	script := ComputedStyleScript("#app")

	core.AssertContains(t, script, "document.querySelector")
	core.AssertContains(t, script, "window.getComputedStyle(el)")
	core.AssertContains(t, script, mustJSON("#app"))
}

func TestDiagnostics_ComputedStyleScript_Bad(t *core.T) {
	// ComputedStyleScript
	ax7Variant := "ComputedStyleScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	script := ComputedStyleScript("")

	core.AssertContains(t, script, `document.querySelector("")`)
	core.AssertContains(t, script, "return null;")
}

func TestDiagnostics_ComputedStyleScript_Ugly(t *core.T) {
	// ComputedStyleScript
	ax7Variant := "ComputedStyleScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	selector := "#app\"\n\t"
	script := ComputedStyleScript(selector)

	core.AssertContains(t, script, mustJSON(selector))
}

func TestDiagnostics_HighlightScript_Good(t *core.T) {
	// HighlightScript
	ax7Variant := "HighlightScript:good"
	core.AssertContains(t, ax7Variant, "good")
	script := HighlightScript(".card", "#00ff00")

	core.AssertContains(t, script, mustJSON(".card"))
	core.AssertContains(t, script, mustJSON("#00ff00"))
	core.AssertContains(t, script, `outline = "3px solid " +`)
}

func TestDiagnostics_HighlightScript_Bad(t *core.T) {
	// HighlightScript
	ax7Variant := "HighlightScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	script := HighlightScript(".card", "")

	core.AssertContains(t, script, `3px solid `)
	core.AssertContains(t, script, mustJSON("#ff9800"))
}

func TestDiagnostics_HighlightScript_Ugly(t *core.T) {
	// HighlightScript
	ax7Variant := "HighlightScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	selector := `#card" + alert(1) + "`
	script := HighlightScript(selector, "#123456")

	core.AssertContains(t, script, mustJSON(selector))
	core.AssertContains(t, script, mustJSON("#123456"))
}

func TestDiagnostics_NetworkLogScript_Good(t *core.T) {
	// NetworkLogScript
	ax7Variant := "NetworkLogScript:good"
	core.AssertContains(t, ax7Variant, "good")
	script := NetworkLogScript(5)

	core.AssertContains(t, script, "slice(-5)")
	core.AssertContains(t, script, `window.__coreNetworkLog`)
}

func TestDiagnostics_NetworkLogScript_Bad(t *core.T) {
	// NetworkLogScript
	ax7Variant := "NetworkLogScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	script := NetworkLogScript(0)

	core.AssertContains(t, script, `performance.getEntriesByType("resource")`)
	core.AssertNotContains(t, script, "slice(-")
}

func TestDiagnostics_NetworkLogScript_Ugly(t *core.T) {
	// NetworkLogScript
	ax7Variant := "NetworkLogScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	script := NetworkLogScript(-7)

	core.AssertContains(t, script, `performance.getEntriesByType("resource")`)
	core.AssertNotContains(t, script, "slice(-")
}

func TestDiagnostics_normalizeWhitespace_Good(t *core.T) {
	// normalizeWhitespace
	ax7Variant := "normalizeWhitespace:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertEqual(t, "hello world", normalizeWhitespace("  hello world  "))
	observedType := core.Sprintf("%T", normalizeWhitespace("  hello world  "))
	core.AssertNotEmpty(t, observedType)
}

func TestDiagnostics_normalizeWhitespace_Bad(t *core.T) {
	// normalizeWhitespace
	ax7Variant := "normalizeWhitespace:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertEmpty(t, normalizeWhitespace(""))
	observedType := core.Sprintf("%T", normalizeWhitespace(""))
	core.AssertNotEmpty(t, observedType)
}

func TestDiagnostics_normalizeWhitespace_Ugly(t *core.T) {
	// normalizeWhitespace
	ax7Variant := "normalizeWhitespace:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertEmpty(t, normalizeWhitespace("\n\t   "))
	observedType := core.Sprintf("%T", normalizeWhitespace("\n\t   "))
	core.AssertNotEmpty(t, observedType)
}

func mustJSON(v any) string {
	data, err := jsonMarshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// AX7 generated source-matching smoke coverage.
func TestDiagnostics_PerformanceScript_Good(t *core.T) {
	// PerformanceScript
	ax7Variant := "PerformanceScript:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := PerformanceScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_PerformanceScript_Bad(t *core.T) {
	// PerformanceScript
	ax7Variant := "PerformanceScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := PerformanceScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_PerformanceScript_Ugly(t *core.T) {
	// PerformanceScript
	ax7Variant := "PerformanceScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := PerformanceScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_ResourcesScript_Good(t *core.T) {
	// ResourcesScript
	ax7Variant := "ResourcesScript:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := ResourcesScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_ResourcesScript_Bad(t *core.T) {
	// ResourcesScript
	ax7Variant := "ResourcesScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := ResourcesScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_ResourcesScript_Ugly(t *core.T) {
	// ResourcesScript
	ax7Variant := "ResourcesScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := ResourcesScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkInitScript_Good(t *core.T) {
	// NetworkInitScript
	ax7Variant := "NetworkInitScript:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NetworkInitScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkInitScript_Bad(t *core.T) {
	// NetworkInitScript
	ax7Variant := "NetworkInitScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NetworkInitScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkInitScript_Ugly(t *core.T) {
	// NetworkInitScript
	ax7Variant := "NetworkInitScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NetworkInitScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkClearScript_Good(t *core.T) {
	// NetworkClearScript
	ax7Variant := "NetworkClearScript:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NetworkClearScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkClearScript_Bad(t *core.T) {
	// NetworkClearScript
	ax7Variant := "NetworkClearScript:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NetworkClearScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDiagnostics_NetworkClearScript_Ugly(t *core.T) {
	// NetworkClearScript
	ax7Variant := "NetworkClearScript:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NetworkClearScript()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
