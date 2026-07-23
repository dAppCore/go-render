package php

import (
	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/io"
)

// DefaultMedium is the default filesystem medium used by the php package.
// It defaults to io.Local (unsandboxed filesystem access).
// Use SetMedium to change this for testing or sandboxed operation.
var DefaultMedium io.Medium = io.Local

// SetMedium sets the default medium for filesystem operations.
// This is primarily useful for testing with mock mediums.
func SetMedium(m io.Medium) {
	DefaultMedium = m
}

// getMedium returns the default medium for filesystem operations.
func getMedium() io.Medium {
	return DefaultMedium
}

// Style aliases from shared
var (
	successStyle = cli.SuccessStyle
	errorStyle   = cli.ErrorStyle
	dimStyle     = cli.DimStyle
	linkStyle    = cli.LinkStyle
)

// Service colors for log output (domain-specific, keep local)
var (
	phpFrankenPHPStyle = cli.NewStyle().Foreground(cli.ColourIndigo500)
	phpViteStyle       = cli.NewStyle().Foreground(cli.ColourYellow500)
	phpHorizonStyle    = cli.NewStyle().Foreground(cli.ColourOrange500)
	phpReverbStyle     = cli.NewStyle().Foreground(cli.ColourViolet500)
	phpRedisStyle      = cli.NewStyle().Foreground(cli.ColourRed500)
)

// Status styles (from shared)
var (
	phpStatusRunning = cli.SuccessStyle
	phpStatusStopped = cli.DimStyle
	phpStatusError   = cli.ErrorStyle
)

// QA command styles (from shared) — most moved to core/lint
var (
	phpQAWarningStyle = cli.WarningStyle
)

// AddPHPCommands adds PHP/Laravel development commands under the php namespace.
func AddPHPCommands(c *core.Core) {
	phpHelpCommand(c, "php", phpT("cmd.php.short"))
	addPHPCommandSet(c, "php")
}

// registerFrankenPHP is set by cmd_serve_frankenphp.go when CGO is enabled.
var registerFrankenPHP func(c *core.Core, prefix string)

// AddPHPRootCommands adds PHP commands directly to root (for standalone core-php binary).
func AddPHPRootCommands(c *core.Core) {
	addPHPCommandSet(c, "")
}

func addPHPCommandSet(c *core.Core, prefix string) {
	// Development
	addPHPDevCommand(c, prefix)
	addPHPLogsCommand(c, prefix)
	addPHPStopCommand(c, prefix)
	addPHPStatusCommand(c, prefix)
	addPHPSSLCommand(c, prefix)

	// Build & Deploy
	addPHPBuildCommand(c, prefix)
	addPHPServeCommand(c, prefix)
	addPHPShellCommand(c, prefix)

	// CI/CD Integration
	addPHPCICommand(c, prefix)

	// Package Management
	addPHPPackagesCommands(c, prefix)

	// Deployment
	addPHPDeployCommands(c, prefix)

	// FrankenPHP embedded commands (CGO only)
	if registerFrankenPHP != nil {
		registerFrankenPHP(c, prefix)
	}
}

func phpCommandPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
}

func activateWorkspacePackage() error { // Result boundary
	wsRoot, config, ok := loadActiveWorkspaceConfig()
	if !ok {
		return nil
	}

	targetDir := activeWorkspacePackageDir(wsRoot, config)
	if !getMedium().IsDir(targetDir) {
		cli.Warnf("Active package directory not found: %s", targetDir)
		return nil
	}

	if r := core.Chdir(targetDir); !r.OK {
		err, _ := r.Value.(error)
		return core.E("php", "failed to change directory to active package", err)
	}

	cli.Print(cliLabelValueFormat, dimStyle.Render("Workspace:"), config.Active)
	return nil
}

func loadActiveWorkspaceConfig() (string, *workspaceConfig, bool) {
	wsRoot, err := findWorkspaceRoot()
	if err != nil {
		return "", nil, false
	}

	config, err := loadWorkspaceConfig(wsRoot)
	if err != nil || config == nil || config.Active == "" {
		return "", nil, false
	}

	return wsRoot, config, true
}

func activeWorkspacePackageDir(wsRoot string, config *workspaceConfig) string {
	pkgDir := config.PackagesDir
	if pkgDir == "" {
		pkgDir = "./packages"
	}
	if !core.PathIsAbs(pkgDir) {
		pkgDir = core.PathJoin(wsRoot, pkgDir)
	}

	return core.PathJoin(pkgDir, config.Active)
}
