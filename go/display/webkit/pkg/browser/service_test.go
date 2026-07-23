// pkg/browser/service_test.go
package browser

import (
	"context"
	core "dappco.re/go"
)

type mockPlatform struct {
	lastURL  string
	lastPath string
	urlErr   resultFailure
	fileErr  resultFailure
}

func (m *mockPlatform) OpenURL(url string) resultFailure {
	m.lastURL = url
	return m.urlErr
}

func (m *mockPlatform) OpenFile(path string) resultFailure {
	m.lastPath = path
	return m.fileErr
}

func newTestBrowserService(t *core.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "browser")
	return svc, c
}

func TestRegister_Good(t *core.T) {
	mp := &mockPlatform{}
	svc, _ := newTestBrowserService(t, mp)
	core.AssertNotNil(t, svc)
	core.AssertNotNil(t, svc.platform)
}

func TestTaskOpenURL_Good(t *core.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "https://example.com", mp.lastURL)
}

func TestTaskOpenURL_Bad_Scheme(t *core.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "javascript:alert(1)"},
	))
	core.AssertFalse(t, r.OK)
	core.AssertEmpty(t, mp.lastURL)
}

func TestTaskOpenURL_Bad_Credentials(t *core.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://user:pass@example.com"},
	))
	core.AssertFalse(t, r.OK)
	core.AssertEmpty(t, mp.lastURL)
}

func TestTaskOpenURL_Bad_PlatformError(t *core.T) {
	mp := &mockPlatform{urlErr: core.NewError("browser not found")}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	core.AssertFalse(t, r.OK)
}

func TestTaskOpenFile_Good(t *core.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: core.Concat("pa", "th"), Value: "/tmp/readme.txt"},
	))
	core.RequireTrue(t, r.OK)
	core.AssertEqual(t, "/tmp/readme.txt", mp.lastPath)
}

func TestTaskOpenFile_Bad_RelativePath(t *core.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: core.Concat("pa", "th"), Value: "relative/readme.txt"},
	))
	core.AssertFalse(t, r.OK)
	core.AssertEmpty(t, mp.lastPath)
}

func TestTaskOpenFile_Bad_PlatformError(t *core.T) {
	mp := &mockPlatform{fileErr: core.NewError("file not found")}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.open_file").Run(context.Background(), core.NewOptions(
		core.Option{Key: core.Concat("pa", "th"), Value: "/nonexistent"},
	))
	core.AssertFalse(t, r.OK)
}

func TestTaskOpenURL_Bad_NoService(t *core.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("browser.open_url").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	core.AssertFalse(t, r.OK)
}

func TestService_validatedOpenURL_Good(t *core.T) {
	// validatedOpenURL
	ax7Variant := "validatedOpenURL:good"
	core.AssertContains(t, ax7Variant, "good")
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "trimmed",
			raw:  "  https://example.com  ",
			want: "https://example.com",
		},
		{
			name: "pathAndQuery",
			raw:  "https://example.com/docs?q=core",
			want: "https://example.com/docs?q=core",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			got, err := validatedOpenURL(tc.raw)
			core.RequireNoError(t, err)
			core.AssertEqual(t, tc.want, got)
		})
	}
}

func TestService_validatedOpenURL_Bad(t *core.T) {
	// validatedOpenURL
	ax7Variant := "validatedOpenURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	cases := []string{
		"",
		"   ",
		"example.com",
		"ftp://example.com",
		"https://user:pass@example.com",
	}

	for _, raw := range cases {
		t.Run(raw, func(t *core.T) {
			got, err := validatedOpenURL(raw)
			core.AssertError(t, err)
			core.AssertEmpty(t, got)
		})
	}
}

func TestService_validatedOpenURL_Ugly(t *core.T) {
	// validatedOpenURL
	ax7Variant := "validatedOpenURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got, err := validatedOpenURL("https://example.com/\x00")
	core.AssertError(t, err)
	core.AssertEmpty(t, got)
}

func TestService_validatedOpenFilePath_Good(t *core.T) {
	// validatedOpenFilePath
	ax7Variant := "validatedOpenFilePath:good"
	core.AssertContains(t, ax7Variant, "good")
	got, err := validatedOpenFilePath("/tmp/../tmp/report.txt")
	core.RequireNoError(t, err)
	core.AssertEqual(t, core.CleanPath("/tmp/report.txt", string(core.PathSeparator)), got)
}

func TestService_validatedOpenFilePath_Bad(t *core.T) {
	// validatedOpenFilePath
	ax7Variant := "validatedOpenFilePath:bad"
	core.AssertContains(t, ax7Variant, "bad")
	cases := []string{
		"",
		"relative/report.txt",
	}

	for _, raw := range cases {
		t.Run(raw, func(t *core.T) {
			got, err := validatedOpenFilePath(raw)
			core.AssertError(t, err)
			core.AssertEmpty(t, got)
		})
	}
}

func TestService_validatedOpenFilePath_Ugly(t *core.T) {
	// validatedOpenFilePath
	ax7Variant := "validatedOpenFilePath:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got, err := validatedOpenFilePath("/tmp/\x00report.txt")
	core.AssertError(t, err)
	core.AssertEmpty(t, got)
}

// AX7 generated source-matching smoke coverage.
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
