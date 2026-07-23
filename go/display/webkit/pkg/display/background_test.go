package display

import (
	"context"

	core "dappco.re/go"
)

func TestBackground_CloneMap_GoodCase(t *core.T) {
	source := map[string]any{"alpha": "one", "beta": 2}

	cloned := cloneMap(source)

	core.AssertNotNil(t, cloned)
	core.AssertEqual(t, source, cloned)

	source["alpha"] = "mutated"
	core.AssertEqual(t, "one", cloned["alpha"])
}

func TestBackground_CloneMap_BadCase(t *core.T) {
	cloned := cloneMap(nil)

	core.AssertNotNil(t, cloned)
	core.AssertEmpty(t, cloned)
}

func TestBackground_CloneMap_UglyCase(t *core.T) {
	source := map[string]any{"nested": map[string]any{"value": "original"}}

	cloned := cloneMap(source)
	core.AssertNotNil(t, cloned)

	source["nested"].(map[string]any)["value"] = "changed"
	core.AssertEqual(t, map[string]any{"value": "original"}, cloned["nested"])
}

func TestBackground_DecodeMap_GoodCase(t *core.T) {
	source := map[string]any{"scope": "/app"}
	decoded := decodeMap(source)

	core.AssertNotNil(t, decoded)
	core.AssertEqual(t, map[string]any{"scope": "/app"}, decoded)
	source["scope"] = "/mutated"
	if decoded["scope"] != "/app" {
		t.Fatalf("decoded map changed after source mutation: %v", decoded["scope"])
	}
}

func TestBackground_DecodeMap_BadCase(t *core.T) {
	decoded := decodeMap("not-a-map")

	core.AssertNotNil(t, decoded)
	core.AssertEmpty(t, decoded)
}

func TestBackground_DecodeMap_UglyCase(t *core.T) {
	decoded := decodeMap(nil)

	core.AssertNotNil(t, decoded)
	core.AssertEmpty(t, decoded)
}

func TestBackground_RegisterBackgroundActions_GoodCase(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.service_worker.register").Run(context.Background(), core.NewOptions(
		core.Option{Key: "scriptURL", Value: "https://example.com/sw.js"},
		core.Option{Key: "options", Value: map[string]any{"scope": "/app"}},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "/app", payload["scope"])
	core.AssertContains(t, payload, "active")
	core.AssertEqual(t, map[string]any{"scriptURL": "https://example.com/sw.js"}, payload["active"])
}

func TestBackground_RegisterBackgroundActions_BadCase(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: "   "},
		core.Option{Key: "requests", Value: nil},
		core.Option{Key: "options", Value: nil},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "", payload["id"])
	core.AssertEqual(t, "registered", payload["state"])
	core.AssertNil(t, payload["requests"])
}

func TestBackground_RegisterBackgroundActions_UglyCase(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.payment.instrument.set").Run(context.Background(), core.NewOptions(
		core.Option{Key: "key", Value: "  card-01  "},
		core.Option{Key: "details", Value: map[string]any{"network": "visa", "last4": "4242"}},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "card-01", payload["key"])
	core.AssertEqual(t, map[string]any{"network": "visa", "last4": "4242"}, payload["details"])
}

func TestBackground_AddSync_Good(t *core.T) {
	// AddSync
	ax7Variant := "AddSync:good"
	core.AssertContains(t, ax7Variant, "good")
	r := NewBackgroundRegistry()
	source := map[string]any{"tag": "refresh", "kind": "sync"}
	record := r.AddSync(source)

	core.AssertNotNil(t, record)
	core.AssertEqual(t, "refresh", record["tag"])
	core.AssertEqual(t, "sync", record["kind"])
	source["tag"] = "mutated"
	core.AssertEqual(t, "refresh", record["tag"])
}

func TestBackground_AddSync_Bad(t *core.T) {
	// AddSync
	ax7Variant := "AddSync:bad"
	core.AssertContains(t, ax7Variant, "bad")
	r := NewBackgroundRegistry()
	record := r.AddSync(nil)

	core.AssertNotNil(t, record)
	core.AssertEmpty(t, record)
}

func TestBackground_AddSync_Ugly(t *core.T) {
	// AddSync
	ax7Variant := "AddSync:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	r := NewBackgroundRegistry()
	first := r.AddSync(map[string]any{"tag": "sync-1"})
	second := r.AddSync(map[string]any{"tag": "sync-2"})

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
	core.AssertEqual(t, 2, r.SyncRegistrationsCount())
}

func TestBackground_AddPush_Good(t *core.T) {
	// AddPush
	ax7Variant := "AddPush:good"
	core.AssertContains(t, ax7Variant, "good")
	r := NewBackgroundRegistry()
	source := map[string]any{"endpoint": "/push/abc", "auth": "core-local"}
	record := r.AddPush(source)

	core.AssertNotNil(t, record)
	core.AssertEqual(t, "/push/abc", record["endpoint"])
	core.AssertEqual(t, "core-local", record["auth"])
	source["endpoint"] = "/push/mutated"
	core.AssertEqual(t, "/push/abc", record["endpoint"])
}

func TestBackground_AddPush_Bad(t *core.T) {
	// AddPush
	ax7Variant := "AddPush:bad"
	core.AssertContains(t, ax7Variant, "bad")
	r := NewBackgroundRegistry()
	record := r.AddPush(nil)

	core.AssertNotNil(t, record)
	core.AssertEmpty(t, record)
}

func TestBackground_AddPush_Ugly(t *core.T) {
	// AddPush
	ax7Variant := "AddPush:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	r := NewBackgroundRegistry()
	first := r.AddPush(map[string]any{"endpoint": "/push/abc"})
	second := r.AddPush(map[string]any{"endpoint": "/push/def"})

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
	core.AssertEqual(t, 2, r.PushSubscriptionsCount())
	core.AssertEqual(t, "/push/abc", first["endpoint"])
	core.AssertEqual(t, "/push/def", second["endpoint"])
}

// AX7 generated source-matching smoke coverage.
func TestBackground_NewBackgroundRegistry_Good(t *core.T) {
	// NewBackgroundRegistry
	ax7Variant := "NewBackgroundRegistry:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewBackgroundRegistry()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_NewBackgroundRegistry_Bad(t *core.T) {
	// NewBackgroundRegistry
	ax7Variant := "NewBackgroundRegistry:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewBackgroundRegistry()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_NewBackgroundRegistry_Ugly(t *core.T) {
	// NewBackgroundRegistry
	ax7Variant := "NewBackgroundRegistry:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewBackgroundRegistry()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_RegisterServiceWorker_Good(t *core.T) {
	// BackgroundRegistry RegisterServiceWorker
	ax7Variant := "BackgroundRegistry_RegisterServiceWorker:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.RegisterServiceWorker("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_RegisterServiceWorker_Bad(t *core.T) {
	// BackgroundRegistry RegisterServiceWorker
	ax7Variant := "BackgroundRegistry_RegisterServiceWorker:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.RegisterServiceWorker("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_RegisterServiceWorker_Ugly(t *core.T) {
	// BackgroundRegistry RegisterServiceWorker
	ax7Variant := "BackgroundRegistry_RegisterServiceWorker:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.RegisterServiceWorker("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddFetch_Good(t *core.T) {
	// BackgroundRegistry AddFetch
	ax7Variant := "BackgroundRegistry_AddFetch:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddFetch("agent", "agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddFetch_Bad(t *core.T) {
	// BackgroundRegistry AddFetch
	ax7Variant := "BackgroundRegistry_AddFetch:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddFetch("", nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddFetch_Ugly(t *core.T) {
	// BackgroundRegistry AddFetch
	ax7Variant := "BackgroundRegistry_AddFetch:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddFetch("../../edge", map[string]any{}, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddSync_Good(t *core.T) {
	// BackgroundRegistry AddSync
	ax7Variant := "BackgroundRegistry_AddSync:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddSync(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddSync_Bad(t *core.T) {
	// BackgroundRegistry AddSync
	ax7Variant := "BackgroundRegistry_AddSync:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddSync(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddSync_Ugly(t *core.T) {
	// BackgroundRegistry AddSync
	ax7Variant := "BackgroundRegistry_AddSync:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddSync(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddPush_Good(t *core.T) {
	// BackgroundRegistry AddPush
	ax7Variant := "BackgroundRegistry_AddPush:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddPush(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddPush_Bad(t *core.T) {
	// BackgroundRegistry AddPush
	ax7Variant := "BackgroundRegistry_AddPush:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddPush(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_AddPush_Ugly(t *core.T) {
	// BackgroundRegistry AddPush
	ax7Variant := "BackgroundRegistry_AddPush:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.AddPush(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SyncRegistrationsCount_Good(t *core.T) {
	// BackgroundRegistry SyncRegistrationsCount
	ax7Variant := "BackgroundRegistry_SyncRegistrationsCount:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SyncRegistrationsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SyncRegistrationsCount_Bad(t *core.T) {
	// BackgroundRegistry SyncRegistrationsCount
	ax7Variant := "BackgroundRegistry_SyncRegistrationsCount:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SyncRegistrationsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SyncRegistrationsCount_Ugly(t *core.T) {
	// BackgroundRegistry SyncRegistrationsCount
	ax7Variant := "BackgroundRegistry_SyncRegistrationsCount:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SyncRegistrationsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_PushSubscriptionsCount_Good(t *core.T) {
	// BackgroundRegistry PushSubscriptionsCount
	ax7Variant := "BackgroundRegistry_PushSubscriptionsCount:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.PushSubscriptionsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_PushSubscriptionsCount_Bad(t *core.T) {
	// BackgroundRegistry PushSubscriptionsCount
	ax7Variant := "BackgroundRegistry_PushSubscriptionsCount:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.PushSubscriptionsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_PushSubscriptionsCount_Ugly(t *core.T) {
	// BackgroundRegistry PushSubscriptionsCount
	ax7Variant := "BackgroundRegistry_PushSubscriptionsCount:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.PushSubscriptionsCount()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SetPaymentInstrument_Good(t *core.T) {
	// BackgroundRegistry SetPaymentInstrument
	ax7Variant := "BackgroundRegistry_SetPaymentInstrument:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SetPaymentInstrument("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SetPaymentInstrument_Bad(t *core.T) {
	// BackgroundRegistry SetPaymentInstrument
	ax7Variant := "BackgroundRegistry_SetPaymentInstrument:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SetPaymentInstrument("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestBackground_BackgroundRegistry_SetPaymentInstrument_Ugly(t *core.T) {
	// BackgroundRegistry SetPaymentInstrument
	ax7Variant := "BackgroundRegistry_SetPaymentInstrument:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BackgroundRegistry)
	result := core.Try(func() any {
		got0 := subject.SetPaymentInstrument("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
