// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"testing"

	i18n "dappco.re/go/core/i18n"
)

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
