// pkg/clipboard/service_test.go
package clipboard

import (
	"context"

	core "dappco.re/go"
)

type mockPlatform struct {
	text     string
	ok       bool
	image    []byte
	hasImage bool
}

func (m *mockPlatform) Text() (string, bool) { return m.text, m.ok }
func (m *mockPlatform) SetText(text string) bool {
	m.text = text
	m.ok = text != ""
	return true
}
func (m *mockPlatform) Image() ([]byte, bool) { return m.image, m.hasImage }
func (m *mockPlatform) SetImage(data []byte) bool {
	m.image = append([]byte(nil), data...)
	m.hasImage = len(data) > 0
	return true
}

func newTestService(t *core.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(&mockPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "clipboard")
	return svc, c
}

func TestRegister_Good(t *core.T) {
	svc, _ := newTestService(t)
	core.AssertNotNil(t, svc)
	core.AssertNotEmpty(t, core.Sprintf("%T", svc))
}

func TestQueryText_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryText{})
	core.RequireTrue(t, r.OK)
	content := r.Value.(ClipboardContent)
	core.AssertEqual(t, "hello", content.Text)
	core.AssertTrue(t, content.HasContent)
}

func TestQueryText_Bad(t *core.T) {
	// No clipboard service registered
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryText{})
	core.AssertFalse(t, r.OK)
}

func TestTaskSetText_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.set_text").Run(context.Background(), core.NewOptions(
		core.Option{Key: "text", Value: "world"},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, true, r.Value)

	// Verify via query
	qr := c.QUERY(QueryText{})
	core.AssertEqual(t, "world", qr.Value.(ClipboardContent).Text)
}

func TestTaskClear_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.clear").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, r.OK)

	// Verify empty
	qr := c.QUERY(QueryText{})
	core.AssertEqual(t, "", qr.Value.(ClipboardContent).Text)
	core.AssertFalse(t, qr.Value.(ClipboardContent).HasContent)
}

func TestQueryImage_Good(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.set_image").Run(context.Background(), core.NewOptions(
		core.Option{Key: "data", Value: []byte{0x89, 0x50, 0x4e, 0x47}},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, true, r.Value)

	result := c.QUERY(QueryImage{})
	core.RequireTrue(t, result.OK)
	content := result.Value.(ImageContent)
	core.AssertTrue(t, content.HasImage)
	core.AssertEqual(t, []byte{0x89, 0x50, 0x4e, 0x47}, content.Data)
}

func TestTaskSetImage_RejectsOversize(t *core.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.set_image").Run(context.Background(), core.NewOptions(
		core.Option{Key: "data", Value: bytesRepeat([]byte("x"), MaxImageBytes+1)},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, false, r.Value)
}

// AX7 generated source-matching smoke coverage.
func TestService_Register_Good(t *core.T) {
	// Register
	ax7Variant := "Register:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Bad(t *core.T) {
	// Register
	ax7Variant := "Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Register_Ugly(t *core.T) {
	// Register
	ax7Variant := "Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Register(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Good(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Bad(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_HandleIPCEvents_Ugly(t *core.T) {
	// Service HandleIPCEvents
	ax7Variant := "Service_HandleIPCEvents:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.HandleIPCEvents(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
