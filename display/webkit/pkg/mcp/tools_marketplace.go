// pkg/mcp/tools_marketplace.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/marketplace"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MarketplaceListInput struct {
	URL string `json:"url,omitempty"`
}

type MarketplaceListOutput struct {
	RegistryURL string                 `json:"registry_url"`
	Manifests   []marketplace.Manifest `json:"manifests"`
}

func (s *Subsystem) marketplaceList(_ context.Context, _ *mcp.CallToolRequest, input MarketplaceListInput) (*mcp.CallToolResult, MarketplaceListOutput, resultFailure) {
	r := s.core.Action("display.marketplace.list").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, MarketplaceListOutput{}, e
		}
		return nil, MarketplaceListOutput{}, core.E("mcp.marketplaceList", "display.marketplace.list failed", nil)
	}
	payload, ok := r.Value.(map[string]any)
	if !ok {
		return nil, MarketplaceListOutput{}, core.E("mcp.marketplaceList", "unexpected result type", nil)
	}
	output := MarketplaceListOutput{RegistryURL: stringValue(payload, "registry_url")}
	if manifests, ok := payload["manifests"].([]marketplace.Manifest); ok {
		output.Manifests = manifests
	}
	return nil, output, nil
}

type MarketplaceFetchInput struct {
	URL string `json:"url"`
}

func (s *Subsystem) marketplaceFetch(_ context.Context, _ *mcp.CallToolRequest, input MarketplaceFetchInput) (*mcp.CallToolResult, marketplace.Manifest, resultFailure) {
	r := s.core.Action("display.marketplace.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, marketplace.Manifest{}, e
		}
		return nil, marketplace.Manifest{}, core.E("mcp.marketplaceFetch", "display.marketplace.fetch failed", nil)
	}
	manifest, ok := r.Value.(marketplace.Manifest)
	if !ok {
		return nil, marketplace.Manifest{}, core.E("mcp.marketplaceFetch", "unexpected result type", nil)
	}
	return nil, manifest, nil
}

type MarketplaceVerifyInput struct {
	URL string `json:"url"`
}

type MarketplaceVerifyOutput struct {
	Manifest marketplace.Manifest `json:"manifest"`
	Digest   string               `json:"digest"`
}

func (s *Subsystem) marketplaceVerify(_ context.Context, _ *mcp.CallToolRequest, input MarketplaceVerifyInput) (*mcp.CallToolResult, MarketplaceVerifyOutput, resultFailure) {
	r := s.core.Action("display.marketplace.verify").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, MarketplaceVerifyOutput{}, e
		}
		return nil, MarketplaceVerifyOutput{}, core.E("mcp.marketplaceVerify", "display.marketplace.verify failed", nil)
	}
	payload, ok := r.Value.(map[string]any)
	if !ok {
		return nil, MarketplaceVerifyOutput{}, core.E("mcp.marketplaceVerify", "unexpected result type", nil)
	}
	output := MarketplaceVerifyOutput{Digest: stringValue(payload, "digest")}
	if manifest, ok := payload["manifest"].(marketplace.Manifest); ok {
		output.Manifest = manifest
	}
	return nil, output, nil
}

type MarketplaceInstallInput struct {
	URL        string `json:"url"`
	InstallDir string `json:"install_dir,omitempty"`
	GitBinary  string `json:"git_binary,omitempty"`
}

type MarketplaceInstallOutput struct {
	Manifest   marketplace.Manifest `json:"manifest"`
	Digest     string               `json:"digest"`
	TargetDir  string               `json:"target_dir"`
	InstallDir string               `json:"install_dir"`
}

func (s *Subsystem) marketplaceInstall(_ context.Context, _ *mcp.CallToolRequest, input MarketplaceInstallInput) (*mcp.CallToolResult, MarketplaceInstallOutput, resultFailure) {
	r := s.core.Action("display.marketplace.install").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
		core.Option{Key: "install_dir", Value: input.InstallDir},
		core.Option{Key: "git_binary", Value: input.GitBinary},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, MarketplaceInstallOutput{}, e
		}
		return nil, MarketplaceInstallOutput{}, core.E("mcp.marketplaceInstall", "display.marketplace.install failed", nil)
	}
	payload, ok := r.Value.(map[string]any)
	if !ok {
		return nil, MarketplaceInstallOutput{}, core.E("mcp.marketplaceInstall", "unexpected result type", nil)
	}
	output := MarketplaceInstallOutput{
		Digest:     stringValue(payload, "digest"),
		TargetDir:  stringValue(payload, "target_dir"),
		InstallDir: stringValue(payload, "install_dir"),
	}
	if manifest, ok := payload["manifest"].(marketplace.Manifest); ok {
		output.Manifest = manifest
	}
	return nil, output, nil
}

func (s *Subsystem) registerMarketplaceTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "marketplace_list",
		Description: `List marketplace manifests from a registry. Example: {"url":"https://example.com/marketplace.yaml"}`,
	}, s.marketplaceList)
	addTool(s, server, &mcp.Tool{
		Name:        "marketplace_fetch",
		Description: `Fetch a marketplace manifest without installing it. Example: {"url":"https://example.com/core-ui.yaml"}`,
	}, s.marketplaceFetch)
	addTool(s, server, &mcp.Tool{
		Name:        "marketplace_verify",
		Description: `Fetch and verify a signed marketplace manifest. Example: {"url":"https://example.com/core-ui.yaml"}`,
	}, s.marketplaceVerify)
	addTool(s, server, &mcp.Tool{
		Name:        "marketplace_install",
		Description: `Fetch, verify, and install a marketplace manifest. Example: {"url":"https://example.com/core-ui.yaml","install_dir":"/Users/me/.core/apps"}`,
	}, s.marketplaceInstall)
}
