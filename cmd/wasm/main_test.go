//go:build js && wasm

package main

import (
	"testing"

	"syscall/js"
)

func TestRenderToString_Good(t *testing.T) {
	gotAny := renderToString(js.Value{}, []js.Value{
		js.ValueOf("C"),
		js.ValueOf("en-GB"),
		js.ValueOf(map[string]any{"C": "<strong>hello</strong>"}),
	})

	got, ok := gotAny.(string)
	if !ok {
		t.Fatalf("renderToString should return string, got %T", gotAny)
	}

	want := `<main role="main" data-block="C"><strong>hello</strong></main>`
	if got != want {
		t.Fatalf("renderToString(...) = %q, want %q", got, want)
	}
}

func TestRenderToString_EmptySlot_Good(t *testing.T) {
	gotAny := renderToString(js.Value{}, []js.Value{
		js.ValueOf("C"),
		js.ValueOf("en-GB"),
		js.ValueOf(map[string]any{"C": ""}),
	})

	got, ok := gotAny.(string)
	if !ok {
		t.Fatalf("renderToString should return string, got %T", gotAny)
	}

	want := `<main role="main" data-block="C"></main>`
	if got != want {
		t.Fatalf("renderToString empty slot = %q, want %q", got, want)
	}
}

func TestRenderToString_VariantTypeGuard(t *testing.T) {
	if got := renderToString(js.Value{}, []js.Value{js.ValueOf(123)}); got != "" {
		t.Fatalf("non-string variant should be empty, got %q", got)
	}

	if got := renderToString(js.Value{}, []js.Value{}); got != "" {
		t.Fatalf("missing variant should be empty, got %q", got)
	}
}

func TestRenderToString_LocaleTypeGuard(t *testing.T) {
	gotAny := renderToString(js.Value{}, []js.Value{
		js.ValueOf("C"),
		js.ValueOf(123),
		js.ValueOf(map[string]any{"C": "x"}),
	})

	got, ok := gotAny.(string)
	if !ok {
		t.Fatalf("renderToString should return string, got %T", gotAny)
	}

	want := `<main role="main" data-block="C">x</main>`
	if got != want {
		t.Fatalf("renderToString with non-string locale = %q, want %q", got, want)
	}
}
