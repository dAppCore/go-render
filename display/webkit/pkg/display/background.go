package display

import (
	"context"
	"sync" // Note: AX-6 — sync.Mutex for registry guard, no core wrapper in pinned core module
	"time"

	core "dappco.re/go"
)

type BackgroundRegistry struct {
	mu                 sync.Mutex
	serviceWorkers     map[string]map[string]any
	backgroundFetches  map[string]map[string]any
	syncRegistrations  []map[string]any
	pushSubscriptions  []map[string]any
	paymentInstruments map[string]map[string]any
}

func NewBackgroundRegistry() *BackgroundRegistry {
	return &BackgroundRegistry{
		serviceWorkers:     make(map[string]map[string]any),
		backgroundFetches:  make(map[string]map[string]any),
		paymentInstruments: make(map[string]map[string]any),
	}
}

func (r *BackgroundRegistry) RegisterServiceWorker(scriptURL string, options map[string]any) map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := map[string]any{
		"script_url": scriptURL,
		"scope":      options["scope"],
		"updated_at": time.Now().UTC().Format(time.RFC3339),
		"active":     true,
	}
	r.serviceWorkers[scriptURL] = record
	return cloneMap(record)
}

func (r *BackgroundRegistry) AddFetch(id string, requests any, options map[string]any) map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := map[string]any{
		"id":         id,
		"requests":   requests,
		"options":    options,
		"state":      "registered",
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	}
	r.backgroundFetches[id] = record
	return cloneMap(record)
}

func (r *BackgroundRegistry) AddSync(record map[string]any) map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.syncRegistrations = append(r.syncRegistrations, cloneMap(record))
	return cloneMap(record)
}

func (r *BackgroundRegistry) AddPush(record map[string]any) map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pushSubscriptions = append(r.pushSubscriptions, cloneMap(record))
	return cloneMap(record)
}

// SyncRegistrationsCount returns the number of registered sync entries.
func (r *BackgroundRegistry) SyncRegistrationsCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.syncRegistrations)
}

// PushSubscriptionsCount returns the number of registered push subscriptions.
func (r *BackgroundRegistry) PushSubscriptionsCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.pushSubscriptions)
}

func (r *BackgroundRegistry) SetPaymentInstrument(key string, details map[string]any) map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()
	record := map[string]any{
		"key":        key,
		"details":    details,
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	}
	r.paymentInstruments[key] = record
	return cloneMap(record)
}

func cloneMap(values map[string]any) map[string]any {
	if values == nil {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = cloneValue(value)
	}
	return cloned
}

func cloneValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneMap(typed)
	case []any:
		cloned := make([]any, len(typed))
		for index, item := range typed {
			cloned[index] = cloneValue(item)
		}
		return cloned
	default:
		return typed
	}
}

func (s *Service) registerBackgroundActions() {
	s.Core().Action("core.background.service_worker.register", func(_ context.Context, opts core.Options) core.Result {
		scriptURL := core.Trim(opts.String("scriptURL"))
		record := s.background.RegisterServiceWorker(scriptURL, decodeMap(opts.Get("options").Value))
		return core.Result{Value: map[string]any{
			"scope":  record["scope"],
			"active": map[string]any{"scriptURL": scriptURL},
		}, OK: true}
	})
	s.Core().Action("core.background.fetch", func(_ context.Context, opts core.Options) core.Result {
		record := s.background.AddFetch(core.Trim(opts.String("id")), opts.Get("requests").Value, decodeMap(opts.Get("options").Value))
		return core.Result{Value: record, OK: true}
	})
	s.Core().Action("core.background.sync", func(_ context.Context, opts core.Options) core.Result {
		record := s.background.AddSync(map[string]any{
			"tag":        opts.String("tag"),
			"kind":       "sync",
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		})
		return core.Result{Value: record, OK: true}
	})
	s.Core().Action("core.background.periodicSync", func(_ context.Context, opts core.Options) core.Result {
		record := s.background.AddSync(map[string]any{
			"tag":        opts.String("tag"),
			"kind":       "periodic",
			"options":    decodeMap(opts.Get("options").Value),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		})
		return core.Result{Value: record, OK: true}
	})
	s.Core().Action("core.background.push.subscribe", func(_ context.Context, opts core.Options) core.Result {
		key := core.Trim(opts.String("applicationServerKey"))
		record := s.background.AddPush(map[string]any{
			"endpoint":             coreRouteURL("push", key),
			"applicationServerKey": key,
			"auth":                 "core-local",
			"updated_at":           time.Now().UTC().Format(time.RFC3339),
		})
		return core.Result{Value: record, OK: true}
	})
	s.Core().Action("core.payment.instrument.set", func(_ context.Context, opts core.Options) core.Result {
		record := s.background.SetPaymentInstrument(core.Trim(opts.String("key")), decodeMap(opts.Get("details").Value))
		return core.Result{Value: record, OK: true}
	})
}

func decodeMap(value any) map[string]any {
	decoded, ok := value.(map[string]any)
	if !ok || decoded == nil {
		return map[string]any{}
	}
	return cloneMap(decoded)
}
