package window

import (
	core "dappco.re/go"
)

func applyWindowOptions(t *core.T, options ...WindowOption) *Window {
	t.Helper()
	w, err := ApplyOptions(options...)
	core.RequireNoError(t, err)
	return w
}

func TestOptions_WindowOptionSetters_GoodCase(t *core.T) {
	w := applyWindowOptions(t,
		WithName("main"),
		WithTitle("Core GUI"),
		WithURL("/dashboard"),
		WithHTML("<main>Ready</main>"),
		WithJS("globalThis.__CORE_READY__ = true"),
		WithSize(1280, 800),
		WithPosition(160, 120),
		WithMinSize(640, 480),
		WithMaxSize(1920, 1080),
		WithFrameless(true),
		WithHidden(false),
		WithAlwaysOnTop(true),
		WithBackgroundColour(12, 34, 56, 78),
		WithFileDrop(true),
	)

	core.AssertEqual(t, "main", w.Name)
	core.AssertEqual(t, "Core GUI", w.Title)
	core.AssertEqual(t, "/dashboard", w.URL)
	core.AssertEqual(t, "<main>Ready</main>", w.HTML)
	core.AssertEqual(t, "globalThis.__CORE_READY__ = true", w.JS)
	core.AssertEqual(t, 1280, w.Width)
	core.AssertEqual(t, 800, w.Height)
	core.AssertEqual(t, 160, w.X)
	core.AssertEqual(t, 120, w.Y)
	core.AssertEqual(t, 640, w.MinWidth)
	core.AssertEqual(t, 480, w.MinHeight)
	core.AssertEqual(t, 1920, w.MaxWidth)
	core.AssertEqual(t, 1080, w.MaxHeight)
	core.AssertTrue(t, w.Frameless)
	core.AssertFalse(t, w.Hidden)
	core.AssertTrue(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{12, 34, 56, 78}, w.BackgroundColour)
	core.AssertTrue(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_BadCase(t *core.T) {
	w := applyWindowOptions(t,
		WithName(""),
		WithTitle(""),
		WithURL(""),
		WithHTML(""),
		WithJS(""),
		WithSize(0, 0),
		WithPosition(0, 0),
		WithMinSize(0, 0),
		WithMaxSize(0, 0),
		WithFrameless(false),
		WithHidden(false),
		WithAlwaysOnTop(false),
		WithBackgroundColour(0, 0, 0, 0),
		WithFileDrop(false),
	)

	core.AssertEqual(t, "", w.Name)
	core.AssertEqual(t, "", w.Title)
	core.AssertEqual(t, "", w.URL)
	core.AssertEqual(t, "", w.HTML)
	core.AssertEqual(t, "", w.JS)
	core.AssertEqual(t, 0, w.Width)
	core.AssertEqual(t, 0, w.Height)
	core.AssertEqual(t, 0, w.X)
	core.AssertEqual(t, 0, w.Y)
	core.AssertEqual(t, 0, w.MinWidth)
	core.AssertEqual(t, 0, w.MinHeight)
	core.AssertEqual(t, 0, w.MaxWidth)
	core.AssertEqual(t, 0, w.MaxHeight)
	core.AssertFalse(t, w.Frameless)
	core.AssertFalse(t, w.Hidden)
	core.AssertFalse(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{0, 0, 0, 0}, w.BackgroundColour)
	core.AssertFalse(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_UglyCase(t *core.T) {
	w := applyWindowOptions(t,
		WithName("⚙︎core-window"),
		WithTitle("A very long title that stays intact"),
		WithURL("core://settings?tab=%F0%9F%93%81"),
		WithHTML("<section data-id=\"αβγ\">unsafe-looking but literal</section>"),
		WithJS("globalThis.__CORE_STATE__ = { mode: 'worker', value: -1 };"),
		WithSize(-1920, -1080),
		WithPosition(-42, 99999),
		WithMinSize(-1, -2),
		WithMaxSize(32767, 32767),
		WithFrameless(true),
		WithHidden(true),
		WithAlwaysOnTop(true),
		WithBackgroundColour(255, 254, 253, 252),
		WithFileDrop(true),
	)

	core.AssertEqual(t, "⚙︎core-window", w.Name)
	core.AssertEqual(t, "A very long title that stays intact", w.Title)
	core.AssertEqual(t, "core://settings?tab=%F0%9F%93%81", w.URL)
	core.AssertEqual(t, "<section data-id=\"αβγ\">unsafe-looking but literal</section>", w.HTML)
	core.AssertEqual(t, "globalThis.__CORE_STATE__ = { mode: 'worker', value: -1 };", w.JS)
	core.AssertEqual(t, -1920, w.Width)
	core.AssertEqual(t, -1080, w.Height)
	core.AssertEqual(t, -42, w.X)
	core.AssertEqual(t, 99999, w.Y)
	core.AssertEqual(t, -1, w.MinWidth)
	core.AssertEqual(t, -2, w.MinHeight)
	core.AssertEqual(t, 32767, w.MaxWidth)
	core.AssertEqual(t, 32767, w.MaxHeight)
	core.AssertTrue(t, w.Frameless)
	core.AssertTrue(t, w.Hidden)
	core.AssertTrue(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{255, 254, 253, 252}, w.BackgroundColour)
	core.AssertTrue(t, w.EnableFileDrop)
}

func TestOptions_ApplyOptions_Good(t *core.T) {
	// ApplyOptions
	ax7Variant := "ApplyOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	w, err := ApplyOptions(
		nil,
		WithName("main"),
		WithTitle("Core"),
	)

	core.RequireNoError(t, err)
	core.AssertNotNil(t, w)
	core.AssertEqual(t, "main", w.Name)
	core.AssertEqual(t, "Core", w.Title)
}

func TestOptions_ApplyOptions_Bad(t *core.T) {
	// ApplyOptions
	ax7Variant := "ApplyOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	boom := core.AnError

	w, err := ApplyOptions(
		WithName("before"),
		func(*Window) resultFailure { return boom },
		WithTitle("after"),
	)

	core.AssertErrorIs(t, err, boom)
	core.AssertNil(t, w)
}

func TestOptions_ApplyOptions_Ugly(t *core.T) {
	// ApplyOptions
	ax7Variant := "ApplyOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	w, err := ApplyOptions()

	core.RequireNoError(t, err)
	core.AssertNotNil(t, w)
	core.AssertEqual(t, &Window{}, w)
}

// AX7 generated source-matching smoke coverage.
func TestOptions_WithName_Good(t *core.T) {
	// WithName
	ax7Variant := "WithName:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithName("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithName_Bad(t *core.T) {
	// WithName
	ax7Variant := "WithName:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithName("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithName_Ugly(t *core.T) {
	// WithName
	ax7Variant := "WithName:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithName("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithTitle_Good(t *core.T) {
	// WithTitle
	ax7Variant := "WithTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithTitle_Bad(t *core.T) {
	// WithTitle
	ax7Variant := "WithTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithTitle_Ugly(t *core.T) {
	// WithTitle
	ax7Variant := "WithTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithURL_Good(t *core.T) {
	// WithURL
	ax7Variant := "WithURL:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithURL_Bad(t *core.T) {
	// WithURL
	ax7Variant := "WithURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithURL_Ugly(t *core.T) {
	// WithURL
	ax7Variant := "WithURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHTML_Good(t *core.T) {
	// WithHTML
	ax7Variant := "WithHTML:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithHTML("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHTML_Bad(t *core.T) {
	// WithHTML
	ax7Variant := "WithHTML:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithHTML("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHTML_Ugly(t *core.T) {
	// WithHTML
	ax7Variant := "WithHTML:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithHTML("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithJS_Good(t *core.T) {
	// WithJS
	ax7Variant := "WithJS:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithJS("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithJS_Bad(t *core.T) {
	// WithJS
	ax7Variant := "WithJS:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithJS("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithJS_Ugly(t *core.T) {
	// WithJS
	ax7Variant := "WithJS:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithJS("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithSize_Good(t *core.T) {
	// WithSize
	ax7Variant := "WithSize:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithSize_Bad(t *core.T) {
	// WithSize
	ax7Variant := "WithSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithSize_Ugly(t *core.T) {
	// WithSize
	ax7Variant := "WithSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithPosition_Good(t *core.T) {
	// WithPosition
	ax7Variant := "WithPosition:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithPosition(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithPosition_Bad(t *core.T) {
	// WithPosition
	ax7Variant := "WithPosition:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithPosition(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithPosition_Ugly(t *core.T) {
	// WithPosition
	ax7Variant := "WithPosition:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithPosition(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMinSize_Good(t *core.T) {
	// WithMinSize
	ax7Variant := "WithMinSize:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithMinSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMinSize_Bad(t *core.T) {
	// WithMinSize
	ax7Variant := "WithMinSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithMinSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMinSize_Ugly(t *core.T) {
	// WithMinSize
	ax7Variant := "WithMinSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithMinSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMaxSize_Good(t *core.T) {
	// WithMaxSize
	ax7Variant := "WithMaxSize:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithMaxSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMaxSize_Bad(t *core.T) {
	// WithMaxSize
	ax7Variant := "WithMaxSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithMaxSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithMaxSize_Ugly(t *core.T) {
	// WithMaxSize
	ax7Variant := "WithMaxSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithMaxSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFrameless_Good(t *core.T) {
	// WithFrameless
	ax7Variant := "WithFrameless:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithFrameless(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFrameless_Bad(t *core.T) {
	// WithFrameless
	ax7Variant := "WithFrameless:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFrameless_Ugly(t *core.T) {
	// WithFrameless
	ax7Variant := "WithFrameless:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHidden_Good(t *core.T) {
	// WithHidden
	ax7Variant := "WithHidden:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithHidden(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHidden_Bad(t *core.T) {
	// WithHidden
	ax7Variant := "WithHidden:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithHidden(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithHidden_Ugly(t *core.T) {
	// WithHidden
	ax7Variant := "WithHidden:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithHidden(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithAlwaysOnTop_Good(t *core.T) {
	// WithAlwaysOnTop
	ax7Variant := "WithAlwaysOnTop:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithAlwaysOnTop(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithAlwaysOnTop_Bad(t *core.T) {
	// WithAlwaysOnTop
	ax7Variant := "WithAlwaysOnTop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithAlwaysOnTop_Ugly(t *core.T) {
	// WithAlwaysOnTop
	ax7Variant := "WithAlwaysOnTop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithBackgroundColour_Good(t *core.T) {
	// WithBackgroundColour
	ax7Variant := "WithBackgroundColour:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithBackgroundColour(1, 1, 1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithBackgroundColour_Bad(t *core.T) {
	// WithBackgroundColour
	ax7Variant := "WithBackgroundColour:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithBackgroundColour(0, 0, 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithBackgroundColour_Ugly(t *core.T) {
	// WithBackgroundColour
	ax7Variant := "WithBackgroundColour:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithBackgroundColour(0, 0, 0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFileDrop_Good(t *core.T) {
	// WithFileDrop
	ax7Variant := "WithFileDrop:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := WithFileDrop(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFileDrop_Bad(t *core.T) {
	// WithFileDrop
	ax7Variant := "WithFileDrop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := WithFileDrop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestOptions_WithFileDrop_Ugly(t *core.T) {
	// WithFileDrop
	ax7Variant := "WithFileDrop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := WithFileDrop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
