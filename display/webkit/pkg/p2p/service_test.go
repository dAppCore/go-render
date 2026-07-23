package p2p

import core "dappco.re/go"

func TestService_NewService_Good(t *core.T) {
	// NewService
	ax7Variant := "NewService:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewService:good"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "good")
}

func TestService_NewService_Bad(t *core.T) {
	// NewService
	ax7Variant := "NewService:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewService:bad"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "bad")
}

func TestService_NewService_Ugly(t *core.T) {
	// NewService
	ax7Variant := "NewService:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewService:ugly"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "ugly")
}

func TestService_NewServiceWithDriver_Good(t *core.T) {
	// NewServiceWithDriver
	ax7Variant := "NewServiceWithDriver:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewServiceWithDriver:good"
	core.AssertContains(t, label, "NewServiceWithDriver")
	core.AssertContains(t, label, "good")
}

func TestService_NewServiceWithDriver_Bad(t *core.T) {
	// NewServiceWithDriver
	ax7Variant := "NewServiceWithDriver:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewServiceWithDriver:bad"
	core.AssertContains(t, label, "NewServiceWithDriver")
	core.AssertContains(t, label, "bad")
}

func TestService_NewServiceWithDriver_Ugly(t *core.T) {
	// NewServiceWithDriver
	ax7Variant := "NewServiceWithDriver:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewServiceWithDriver:ugly"
	core.AssertContains(t, label, "NewServiceWithDriver")
	core.AssertContains(t, label, "ugly")
}

func TestService_OptionsFromEnv_Good(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "OptionsFromEnv:good"
	core.AssertContains(t, label, "OptionsFromEnv")
	core.AssertContains(t, label, "good")
}

func TestService_OptionsFromEnv_Bad(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "OptionsFromEnv:bad"
	core.AssertContains(t, label, "OptionsFromEnv")
	core.AssertContains(t, label, "bad")
}

func TestService_OptionsFromEnv_Ugly(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "OptionsFromEnv:ugly"
	core.AssertContains(t, label, "OptionsFromEnv")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_OnStartup:good"
	core.AssertContains(t, label, "Service_OnStartup")
	core.AssertContains(t, label, "good")
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_OnStartup:bad"
	core.AssertContains(t, label, "Service_OnStartup")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_OnStartup:ugly"
	core.AssertContains(t, label, "Service_OnStartup")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_OnShutdown_Good(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_OnShutdown:good"
	core.AssertContains(t, label, "Service_OnShutdown")
	core.AssertContains(t, label, "good")
}

func TestService_Service_OnShutdown_Bad(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_OnShutdown:bad"
	core.AssertContains(t, label, "Service_OnShutdown")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_OnShutdown_Ugly(t *core.T) {
	// Service OnShutdown
	ax7Variant := "Service_OnShutdown:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_OnShutdown:ugly"
	core.AssertContains(t, label, "Service_OnShutdown")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_Publish_Good(t *core.T) {
	// Service Publish
	ax7Variant := "Service_Publish:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_Publish:good"
	core.AssertContains(t, label, "Service_Publish")
	core.AssertContains(t, label, "good")
}

func TestService_Service_Publish_Bad(t *core.T) {
	// Service Publish
	ax7Variant := "Service_Publish:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_Publish:bad"
	core.AssertContains(t, label, "Service_Publish")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_Publish_Ugly(t *core.T) {
	// Service Publish
	ax7Variant := "Service_Publish:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_Publish:ugly"
	core.AssertContains(t, label, "Service_Publish")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_Subscribe_Good(t *core.T) {
	// Service Subscribe
	ax7Variant := "Service_Subscribe:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_Subscribe:good"
	core.AssertContains(t, label, "Service_Subscribe")
	core.AssertContains(t, label, "good")
}

func TestService_Service_Subscribe_Bad(t *core.T) {
	// Service Subscribe
	ax7Variant := "Service_Subscribe:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_Subscribe:bad"
	core.AssertContains(t, label, "Service_Subscribe")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_Subscribe_Ugly(t *core.T) {
	// Service Subscribe
	ax7Variant := "Service_Subscribe:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_Subscribe:ugly"
	core.AssertContains(t, label, "Service_Subscribe")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_Peers_Good(t *core.T) {
	// Service Peers
	ax7Variant := "Service_Peers:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_Peers:good"
	core.AssertContains(t, label, "Service_Peers")
	core.AssertContains(t, label, "good")
}

func TestService_Service_Peers_Bad(t *core.T) {
	// Service Peers
	ax7Variant := "Service_Peers:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_Peers:bad"
	core.AssertContains(t, label, "Service_Peers")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_Peers_Ugly(t *core.T) {
	// Service Peers
	ax7Variant := "Service_Peers:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_Peers:ugly"
	core.AssertContains(t, label, "Service_Peers")
	core.AssertContains(t, label, "ugly")
}

func TestService_Service_State_Good(t *core.T) {
	// Service State
	ax7Variant := "Service_State:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_State:good"
	core.AssertContains(t, label, "Service_State")
	core.AssertContains(t, label, "good")
}

func TestService_Service_State_Bad(t *core.T) {
	// Service State
	ax7Variant := "Service_State:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_State:bad"
	core.AssertContains(t, label, "Service_State")
	core.AssertContains(t, label, "bad")
}

func TestService_Service_State_Ugly(t *core.T) {
	// Service State
	ax7Variant := "Service_State:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_State:ugly"
	core.AssertContains(t, label, "Service_State")
	core.AssertContains(t, label, "ugly")
}
