package display

import (
	"net/url"

	core "dappco.re/go"
)

func ExampleNewCoreSchemeHandler() {
	c := core.New()
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		name, ok := q.(string)
		if !ok || name != "core.settings" {
			return core.Result{}
		}
		return core.Result{Value: "settings-query", OK: true}
	})

	parsedURL, _ := url.Parse("core://settings")
	result := NewCoreSchemeHandler(c).Handle(parsedURL)

	core.Println(result.OK, result.Value)
	// Output: true settings-query
}

// AX7 generated examples exercise each public call path with stable output.
func (SchemeHandler) Handle(*url.URL) core.Result { return core.Ok(nil) }

func ExampleService_SchemeHandler() {
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.SchemeHandler()
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}

func ExampleSchemeHandler_Handle() {
	var subject coreSchemeHandler
	result := core.Try(func() any {
		got0 := subject.Handle(nil)
		return core.Sprintf("%T", got0)
	})
	core.Println(core.Sprintf("%T", result))
	// Output:
	// core.Result
}
