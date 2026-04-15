// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"reflect"
	"testing"

	i18n "dappco.re/go/core/i18n"
)

type recordingTranslator struct {
	key  string
	args []any
}

func (r *recordingTranslator) T(key string, args ...any) string {
	r.key = key
	r.args = append(r.args[:0], args...)
	return "translated"
}

func TestNewContext_OptionalLocale_Good(t *testing.T) {
	ctx := NewContext("en-GB")

	if ctx == nil {
		t.Fatal("NewContext returned nil")
	}
	if ctx.Locale != "en-GB" {
		t.Fatalf("NewContext locale = %q, want %q", ctx.Locale, "en-GB")
	}
	if ctx.Data == nil {
		t.Fatal("NewContext should initialise Data")
	}
}

func TestNewContextWithService_OptionalLocale_Good(t *testing.T) {
	svc, _ := i18n.New()
	ctx := NewContextWithService(svc, "fr-FR")

	if ctx == nil {
		t.Fatal("NewContextWithService returned nil")
	}
	if ctx.Locale != "fr-FR" {
		t.Fatalf("NewContextWithService locale = %q, want %q", ctx.Locale, "fr-FR")
	}
	if ctx.service == nil {
		t.Fatal("NewContextWithService should set translator service")
	}
}

func TestNewContextWithService_AppliesLocaleToService_Good(t *testing.T) {
	svc, _ := i18n.New()
	ctx := NewContextWithService(svc, "fr-FR")

	got := Text("prompt.yes").Render(ctx)
	if got != "o" {
		t.Fatalf("NewContextWithService locale translation = %q, want %q", got, "o")
	}
}

func TestTextNode_UsesMetadataAliasWhenDataNil_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)

	ctx := &Context{
		Metadata: map[string]any{"count": 1},
	}

	got := Text("i18n.count.file").Render(ctx)
	if got != "1 file" {
		t.Fatalf("Text with metadata-only count = %q, want %q", got, "1 file")
	}
}

func TestTextNode_CustomTranslatorReceivesCountArgs_Good(t *testing.T) {
	ctx := NewContextWithService(&recordingTranslator{})
	ctx.Metadata["count"] = 3

	got := Text("i18n.count.file", "ignored").Render(ctx)
	if got != "translated" {
		t.Fatalf("Text with custom translator = %q, want %q", got, "translated")
	}

	svc := ctx.service.(*recordingTranslator)
	if svc.key != "i18n.count.file" {
		t.Fatalf("custom translator key = %q, want %q", svc.key, "i18n.count.file")
	}

	wantArgs := []any{3, "ignored"}
	if !reflect.DeepEqual(svc.args, wantArgs) {
		t.Fatalf("custom translator args = %#v, want %#v", svc.args, wantArgs)
	}
}

func TestContext_SetService_AppliesLocale_Good(t *testing.T) {
	svc, _ := i18n.New()
	ctx := NewContext("fr-FR")

	if got := ctx.SetService(svc); got != ctx {
		t.Fatal("SetService should return the same context for chaining")
	}

	got := Text("prompt.yes").Render(ctx)
	if got != "o" {
		t.Fatalf("SetService locale translation = %q, want %q", got, "o")
	}
}

func TestContext_SetService_NilContext_Ugly(t *testing.T) {
	var ctx *Context
	if got := ctx.SetService(nil); got != nil {
		t.Fatal("SetService on nil context should return nil")
	}
}

func TestContext_SetLocale_AppliesLocale_Good(t *testing.T) {
	svc, _ := i18n.New()
	ctx := NewContextWithService(svc)

	if got := ctx.SetLocale("fr-FR"); got != ctx {
		t.Fatal("SetLocale should return the same context for chaining")
	}

	got := Text("prompt.yes").Render(ctx)
	if got != "o" {
		t.Fatalf("SetLocale translation = %q, want %q", got, "o")
	}
}

func TestContext_SetLocale_NilContext_Ugly(t *testing.T) {
	var ctx *Context
	if got := ctx.SetLocale("en-GB"); got != nil {
		t.Fatal("SetLocale on nil context should return nil")
	}
}

func TestCloneContext_PreservesMetadataAlias_Good(t *testing.T) {
	ctx := NewContext()
	ctx.Data["count"] = 3

	clone := cloneContext(ctx)
	if clone == nil {
		t.Fatal("cloneContext returned nil")
	}
	if clone.Data == nil || clone.Metadata == nil {
		t.Fatal("cloneContext should preserve non-nil metadata maps")
	}

	dataPtr := reflect.ValueOf(clone.Data).Pointer()
	metadataPtr := reflect.ValueOf(clone.Metadata).Pointer()
	if dataPtr != metadataPtr {
		t.Fatalf("cloneContext should keep Data and Metadata aliased, got %x and %x", dataPtr, metadataPtr)
	}
	if clone.Data["count"] != 3 || clone.Metadata["count"] != 3 {
		t.Fatalf("cloneContext should copy map contents, got Data=%v Metadata=%v", clone.Data, clone.Metadata)
	}
}
