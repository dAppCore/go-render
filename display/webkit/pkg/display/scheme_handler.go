package display

import (
	"context"
	"net/url"

	core "dappco.re/go"
)

type routeDispatchKind uint8

const (
	routeDispatchQuery routeDispatchKind = iota
	routeDispatchAction
)

var coreRouteDispatch = map[string]routeDispatchKind{
	"settings": routeDispatchQuery,
	"store":    routeDispatchQuery,
	"network":  routeDispatchQuery,
	"models":   routeDispatchQuery,
	"agent":    routeDispatchAction,
	"wallet":   routeDispatchAction,
	"identity": routeDispatchAction,
}

type coreSchemeHandler struct {
	core *core.Core
}

// CoreRouteQuery carries a core:// query route target and sanitized URL parameters.
type CoreRouteQuery struct {
	Target string
	Params url.Values
}

// NewCoreSchemeHandler returns the RFC route dispatcher for `core://` URLs.
//
//	handler := display.NewCoreSchemeHandler(c)
func NewCoreSchemeHandler(c *core.Core) RouteSchemeHandler {
	return coreSchemeHandler{core: c}
}

// SchemeHandler exposes the RFC route dispatcher for the active display service.
//
//	handler := svc.SchemeHandler()
func (s *Service) SchemeHandler() RouteSchemeHandler {
	if s == nil || s.ServiceRuntime == nil {
		return coreSchemeHandler{}
	}
	return NewCoreSchemeHandler(s.Core())
}

func (h coreSchemeHandler) Handle(rawURL *url.URL) core.Result {
	if h.core == nil {
		return core.Result{
			Value: core.E("display.coreSchemeHandler.Handle", "core runtime unavailable", nil),
			OK:    false,
		}
	}

	route, dispatch, result := resolveCoreSchemeRoute(rawURL)
	if !result.OK {
		return result
	}

	target := "core." + route
	params := cloneURLValues(rawURL.Query())
	switch dispatch {
	case routeDispatchAction:
		return h.core.Action(target).Run(context.Background(), optionsFromURLValues(params))
	case routeDispatchQuery:
		if len(params) > 0 {
			result = h.core.Query(CoreRouteQuery{Target: target, Params: params})
			if result.OK {
				return result
			}
		}
		result = h.core.Query(target)
		if result.OK {
			return result
		}
		return core.Result{
			Value: core.E("display.coreSchemeHandler.Handle", "query not handled: "+target, nil),
			OK:    false,
		}
	default:
		return core.Result{
			Value: core.E("display.coreSchemeHandler.Handle", "unsupported dispatch kind", nil),
			OK:    false,
		}
	}
}

func resolveCoreSchemeRoute(rawURL *url.URL) (string, routeDispatchKind, core.Result) {
	if rawURL == nil {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", "scheme URL is required", nil),
			OK:    false,
		}
	}
	if !equalFold(core.Trim(rawURL.Scheme), "core") {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", "unsupported scheme: "+rawURL.Scheme, nil),
			OK:    false,
		}
	}
	if core.Trim(rawURL.Opaque) != "" {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}
	if rawURL.User != nil || core.Trim(rawURL.Fragment) != "" || rawURL.Port() != "" {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}
	if path := core.Trim(rawURL.Path); path != "" && path != "/" {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}

	route := core.Lower(core.Trim(rawURL.Hostname()))
	if route == "" {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}

	dispatch, ok := coreRouteDispatch[route]
	if !ok {
		return "", routeDispatchQuery, core.Result{
			Value: core.E("display.resolveCoreSchemeRoute", "unknown core route: "+route, nil),
			OK:    false,
		}
	}

	return route, dispatch, core.Result{OK: true}
}

const malformedCoreURL = "malformed core URL"

func optionsFromURLValues(values url.Values) core.Options {
	opts := core.NewOptions()
	for key, items := range sanitizeCoreQuery(values) {
		if len(items) == 1 {
			opts.Set(key, items[0])
			continue
		}
		opts.Set(key, append([]string(nil), items...))
	}
	return opts
}
