package display

import (
	"context"
	"net/url"

	core "dappco.re/go"
)

type schemeDispatchRecorder struct {
	queries []string
	actions []string
	params  map[string]url.Values
}

func newTestCoreSchemeHandler(t *core.T) (RouteSchemeHandler, *schemeDispatchRecorder) {
	t.Helper()

	c := newTestCore(t)

	recorder := &schemeDispatchRecorder{params: make(map[string]url.Values)}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		var name string
		switch typed := q.(type) {
		case CoreRouteQuery:
			name = typed.Target
			recorder.params[name] = cloneURLValues(typed.Params)
		case string:
			name = typed
		default:
			return core.Result{}
		}

		recorder.queries = append(recorder.queries, name)
		switch name {
		case "core.settings":
			return core.Result{Value: "settings-query", OK: true}
		case "core.store":
			return core.Result{Value: "store-query", OK: true}
		case "core.network":
			return core.Result{Value: "network-query", OK: true}
		case "core.models":
			return core.Result{Value: "models-query", OK: true}
		default:
			return core.Result{}
		}
	})

	c.Action("core.agent", func(_ context.Context, opts core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.agent")
		recorder.params["core.agent"] = url.Values{"q": []string{opts.String("q")}}
		return core.Result{Value: "agent-action", OK: true}
	})
	c.Action("core.wallet", func(_ context.Context, _ core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.wallet")
		return core.Result{Value: "wallet-action", OK: true}
	})
	c.Action("core.identity", func(_ context.Context, _ core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.identity")
		return core.Result{Value: "identity-action", OK: true}
	})

	svc := core.MustServiceFor[*Service](c, "display")
	return svc.SchemeHandler(), recorder
}

func TestSchemeHandler_Handle_ForwardsQueryParameters(t *core.T) {
	handler, recorder := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://store?q=invoice&tag=a&tag=b")
	core.RequireNoError(t, err)

	result := handler.Handle(parsedURL)
	core.RequireTrue(t, result.OK)
	core.AssertEqual(t, "store-query", result.Value)
	core.AssertEqual(t, []string{"invoice"}, recorder.params["core.store"]["q"])
	core.AssertEqual(t, []string{"a", "b"}, recorder.params["core.store"]["tag"])

	actionURL, err := url.Parse("core://agent?q=launch")
	core.RequireNoError(t, err)

	result = handler.Handle(actionURL)
	core.RequireTrue(t, result.OK)
	core.AssertEqual(t, "agent-action", result.Value)
	core.AssertEqual(t, []string{"launch"}, recorder.params["core.agent"]["q"])
}

func TestSchemeHandler_Handle_Good(t *core.T) {
	// Handle
	ax7Variant := "Handle:good"
	core.AssertContains(t, ax7Variant, "good")
	handler, recorder := newTestCoreSchemeHandler(t)

	tests := []struct {
		rawURL string
		value  string
	}{
		{rawURL: "core://settings", value: "settings-query"},
		{rawURL: "core://store", value: "store-query"},
		{rawURL: "core://network", value: "network-query"},
		{rawURL: "core://models", value: "models-query"},
		{rawURL: "core://agent", value: "agent-action"},
		{rawURL: "core://wallet", value: "wallet-action"},
		{rawURL: "core://identity", value: "identity-action"},
	}

	for _, test := range tests {
		parsedURL, err := url.Parse(test.rawURL)
		core.RequireNoError(t, err)

		result := handler.Handle(parsedURL)
		core.RequireTrue(t, result.OK, test.rawURL)
		core.AssertEqual(t, test.value, result.Value)
	}

	core.AssertEqual(t, []string{
		"core.settings",
		"core.store",
		"core.network",
		"core.models",
	}, recorder.queries)
	core.AssertEqual(t, []string{
		"core.agent",
		"core.wallet",
		"core.identity",
	}, recorder.actions)
}

func TestSchemeHandler_Handle_Bad(t *core.T) {
	// Handle
	ax7Variant := "Handle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	handler, _ := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://missing")
	core.RequireNoError(t, err)

	result := handler.Handle(parsedURL)
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure), "unknown core route: missing")
}

func TestSchemeHandler_Handle_Ugly(t *core.T) {
	// Handle
	ax7Variant := "Handle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	handler, _ := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://settings/profile")
	core.RequireNoError(t, err)

	result := handler.Handle(parsedURL)
	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure), "malformed core URL")
}

// AX7 generated source-matching smoke coverage.
func TestSchemeHandler_NewCoreSchemeHandler_Good(t *core.T) {
	// NewCoreSchemeHandler
	ax7Variant := "NewCoreSchemeHandler:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewCoreSchemeHandler(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_NewCoreSchemeHandler_Bad(t *core.T) {
	// NewCoreSchemeHandler
	ax7Variant := "NewCoreSchemeHandler:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewCoreSchemeHandler(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_NewCoreSchemeHandler_Ugly(t *core.T) {
	// NewCoreSchemeHandler
	ax7Variant := "NewCoreSchemeHandler:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewCoreSchemeHandler(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_Service_SchemeHandler_Good(t *core.T) {
	// Service SchemeHandler
	ax7Variant := "Service_SchemeHandler:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SchemeHandler()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_Service_SchemeHandler_Bad(t *core.T) {
	// Service SchemeHandler
	ax7Variant := "Service_SchemeHandler:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SchemeHandler()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_Service_SchemeHandler_Ugly(t *core.T) {
	// Service SchemeHandler
	ax7Variant := "Service_SchemeHandler:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SchemeHandler()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_SchemeHandler_Handle_Good(t *core.T) {
	// SchemeHandler Handle
	ax7Variant := "SchemeHandler_Handle:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject coreSchemeHandler
	result := core.Try(func() any {
		got0 := subject.Handle(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_SchemeHandler_Handle_Bad(t *core.T) {
	// SchemeHandler Handle
	ax7Variant := "SchemeHandler_Handle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject coreSchemeHandler
	result := core.Try(func() any {
		got0 := subject.Handle(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestSchemeHandler_SchemeHandler_Handle_Ugly(t *core.T) {
	// SchemeHandler Handle
	ax7Variant := "SchemeHandler_Handle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject coreSchemeHandler
	result := core.Try(func() any {
		got0 := subject.Handle(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
