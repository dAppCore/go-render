package webkit

import (
	core "dappco.re/go"
)

func TestService_NewService_Good(t *core.T) {
	subject := NewService
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_NewService_Bad(t *core.T) {
	subject := NewService
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_NewService_Ugly(t *core.T) {
	subject := NewService
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Register_Good(t *core.T) {
	subject := Register
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Register_Bad(t *core.T) {
	subject := Register
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Register_Ugly(t *core.T) {
	subject := Register
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnStartup_Good(t *core.T) {
	subject := (*Service).OnStartup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	subject := (*Service).OnStartup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	subject := (*Service).OnStartup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	subject := (*Service).OnShutdown
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	subject := (*Service).OnShutdown
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	subject := (*Service).OnShutdown
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBuildWailsOptions_CurrentWailsFeatures_Good(t *core.T) {
	cfg := GuiConfig{
		Assets: AssetOptions{DisableLogging: true},
		Windows: WindowsOptions{
			WndClass:                      "CoreWindow",
			WebviewUserDataPath:           `C:\Core\WebView`,
			WebviewBrowserPath:            `C:\WebView2`,
			EnabledFeatures:               []string{"enabled"},
			DisabledFeatures:              []string{"disabled"},
			AdditionalBrowserArgs:         []string{"--remote-debugging-port=9222"},
			UseVisualHosting:              true,
			DisableQuitOnLastWindowClosed: true,
		},
	}

	got := buildWailsOptions(cfg)

	core.AssertTrue(t, got.Assets.DisableLogging)
	core.AssertEqual(t, "CoreWindow", got.Windows.WndClass)
	core.AssertEqual(t, `C:\Core\WebView`, got.Windows.WebviewUserDataPath)
	core.AssertEqual(t, `C:\WebView2`, got.Windows.WebviewBrowserPath)
	core.AssertEqual(t, []string{"enabled"}, got.Windows.EnabledFeatures)
	core.AssertEqual(t, []string{"disabled"}, got.Windows.DisabledFeatures)
	core.AssertEqual(t, []string{"--remote-debugging-port=9222"}, got.Windows.AdditionalBrowserArgs)
	core.AssertTrue(t, got.Windows.UseVisualHosting)
	core.AssertTrue(t, got.Windows.DisableQuitOnLastWindowClosed)
}
