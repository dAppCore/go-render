//go:build !js

// SPDX-License-Identifier: EUPL-1.2

package html

import (
	"testing"

	core "dappco.re/go"
)

// TestNewService_DefaultLocale — empty options falls back to "en".
func TestNewService_DefaultLocale(t *testing.T) {
	c := core.New(core.WithService(NewService(Options{})))
	r := c.Service("html")
	if !r.OK {
		t.Fatal("html service not registered")
	}
	svc := r.Value.(*Service)
	if svc.Context() == nil {
		t.Fatal("expected non-nil Context")
	}
	if svc.Context().Locale != "en" {
		t.Fatalf("expected default locale en, got %q", svc.Context().Locale)
	}
}

// TestNewService_LocaleOption — explicit Locale flows through to the Context.
func TestNewService_LocaleOption(t *testing.T) {
	c := core.New(core.WithService(NewService(Options{Locale: "de"})))
	svc := c.Service("html").Value.(*Service)
	if svc.Context().Locale != "de" {
		t.Fatalf("expected locale de, got %q", svc.Context().Locale)
	}
}

// TestRegister_Imperative — defaults shorthand.
func TestRegister_Imperative(t *testing.T) {
	c := core.New(core.WithService(Register))
	if !c.Service("html").OK {
		t.Fatal("html not registered via Register")
	}
}

// TestService_Context_NilReceiver — nil-receiver guard.
func TestService_Context_NilReceiver(t *testing.T) {
	var svc *Service
	if svc.Context() != nil {
		t.Fatal("expected nil Context on nil-receiver")
	}
}
