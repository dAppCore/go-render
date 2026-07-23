package php

import (
	"context"
	"io"
	"os"
	"os/exec"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// DockerBuildOptions configures Docker image building for PHP projects.
type DockerBuildOptions struct {
	// ProjectDir is the path to the PHP/Laravel project.
	ProjectDir string

	// ImageName is the name for the Docker image.
	ImageName string

	// Tag is the image tag (default: "latest").
	Tag string

	// Platform specifies the target platform (e.g., "linux/amd64", "linux/arm64").
	Platform string

	// Dockerfile is the path to a custom Dockerfile.
	// If empty, one will be auto-generated for FrankenPHP.
	Dockerfile string

	// NoBuildCache disables Docker build cache.
	NoBuildCache bool

	// BuildArgs are additional build arguments.
	BuildArgs map[string]string

	// Output is the writer for build output (default: os.Stdout).
	Output io.Writer
}

// LinuxKitBuildOptions configures LinuxKit image building for PHP projects.
type LinuxKitBuildOptions struct {
	// ProjectDir is the path to the PHP/Laravel project.
	ProjectDir string

	// OutputPath is the path for the output image.
	OutputPath string

	// Format is the output format: "iso", "qcow2", "raw", "vmdk".
	Format string

	// Template is the LinuxKit template name (default: server-php).
	Template string

	// Variables are template variables to apply.
	Variables map[string]string

	// Output is the writer for build output (default: os.Stdout).
	Output io.Writer
}

// ServeOptions configures running a production PHP container.
type ServeOptions struct {
	// ImageName is the Docker image to run.
	ImageName string

	// Tag is the image tag (default: "latest").
	Tag string

	// ContainerName is the name for the container.
	ContainerName string

	// Port is the host port to bind (default: 80).
	Port int

	// HTTPSPort is the host HTTPS port to bind (default: 443).
	HTTPSPort int

	// Detach runs the container in detached mode.
	Detach bool

	// EnvFile is the path to an environment file.
	EnvFile string

	// Volumes maps host paths to container paths.
	Volumes map[string]string

	// Output is the writer for output (default: os.Stdout).
	Output io.Writer
}

// BuildDocker builds a Docker image for the PHP project.
func BuildDocker(ctx context.Context, opts DockerBuildOptions) error { // Result boundary
	opts, err := normalizeDockerBuildOptions(opts)
	if err != nil {
		return err
	}

	dockerfilePath, cleanup, err := resolveDockerfilePath(opts)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	cmd := exec.CommandContext(ctx, "docker", dockerBuildArgs(opts, dockerfilePath)...)
	cmd.Dir = opts.ProjectDir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	if err := cmd.Run(); err != nil {
		return phpWrapMessage(err, "docker build failed")
	}

	return nil
}

func normalizeDockerBuildOptions(opts DockerBuildOptions) (DockerBuildOptions, error) { // Result boundary
	if opts.ProjectDir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return opts, phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.ProjectDir = cwdR.Value.(string)
	}

	// Validate project directory
	if !IsPHPProject(opts.ProjectDir) {
		return opts, phpFailure("not a PHP project: %s (missing composer.json)", opts.ProjectDir)
	}

	// Set defaults
	if opts.ImageName == "" {
		opts.ImageName = core.PathBase(opts.ProjectDir)
	}
	if opts.Tag == "" {
		opts.Tag = "latest"
	}
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	return opts, nil
}

func resolveDockerfilePath(opts DockerBuildOptions) (string, func(), error) { // Result boundary
	if opts.Dockerfile != "" {
		return opts.Dockerfile, nil, nil
	}

	content, err := GenerateDockerfile(opts.ProjectDir)
	if err != nil {
		return "", nil, phpWrapAction(err, "generate", "Dockerfile")
	}

	m := getMedium()
	tempDockerfile := core.PathJoin(opts.ProjectDir, "Dockerfile.core-generated")
	if err := m.Write(tempDockerfile, content); err != nil {
		return "", nil, phpWrapAction(err, "write", "Dockerfile")
	}

	return tempDockerfile, func() { _ = m.Delete(tempDockerfile) }, nil
}

func dockerBuildArgs(opts DockerBuildOptions, dockerfilePath string) []string {
	imageRef := cli.Sprintf("%s:%s", opts.ImageName, opts.Tag)

	args := []string{"build", "-t", imageRef, "-f", dockerfilePath}

	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	if opts.NoBuildCache {
		args = append(args, "--no-cache")
	}

	for key, value := range opts.BuildArgs {
		args = append(args, "--build-arg", cli.Sprintf("%s=%s", key, value))
	}

	args = append(args, opts.ProjectDir)
	return args
}

// BuildLinuxKit builds a LinuxKit image for the PHP project.
func BuildLinuxKit(ctx context.Context, opts LinuxKitBuildOptions) error { // Result boundary
	if opts.ProjectDir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.ProjectDir = cwdR.Value.(string)
	}

	// Validate project directory
	if !IsPHPProject(opts.ProjectDir) {
		return phpFailure("not a PHP project: %s (missing composer.json)", opts.ProjectDir)
	}

	// Set defaults
	if opts.Template == "" {
		opts.Template = defaultLinuxKitTemplateName
	}
	if opts.Format == "" {
		opts.Format = "qcow2"
	}
	if opts.OutputPath == "" {
		opts.OutputPath = core.PathJoin(opts.ProjectDir, "dist", core.PathBase(opts.ProjectDir))
	}
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	// Ensure output directory exists
	m := getMedium()
	outputDir := core.PathDir(opts.OutputPath)
	if err := m.EnsureDir(outputDir); err != nil {
		return phpWrapAction(err, "create", "output directory")
	}

	// Find linuxkit binary
	linuxkitPath, err := lookupLinuxKit()
	if err != nil {
		return err
	}

	// Get template content
	templateContent, err := getLinuxKitTemplate(opts.Template)
	if err != nil {
		return phpWrapAction(err, "get", "template")
	}

	// Apply variables
	if opts.Variables == nil {
		opts.Variables = make(map[string]string)
	}
	// Add project-specific variables
	opts.Variables["PROJECT_DIR"] = opts.ProjectDir
	opts.Variables["PROJECT_NAME"] = core.PathBase(opts.ProjectDir)

	content, err := applyTemplateVariables(templateContent, opts.Variables)
	if err != nil {
		return phpWrapAction(err, "apply", "template variables")
	}

	// Write template to temp file
	tempYAML := core.PathJoin(opts.ProjectDir, ".core-linuxkit.yml")
	if err := m.Write(tempYAML, content); err != nil {
		return phpWrapAction(err, "write", "template")
	}
	defer func() { _ = m.Delete(tempYAML) }()

	// Build LinuxKit image
	args := []string{
		"build",
		"--format", opts.Format,
		"--name", opts.OutputPath,
		tempYAML,
	}

	cmd := exec.CommandContext(ctx, linuxkitPath, args...)
	cmd.Dir = opts.ProjectDir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	if err := cmd.Run(); err != nil {
		return phpWrapMessage(err, "linuxkit build failed")
	}

	return nil
}

// ServeProduction runs a production PHP container.
func ServeProduction(ctx context.Context, opts ServeOptions) error { // Result boundary
	if opts.ImageName == "" {
		return phpFailure("image name is required")
	}

	// Set defaults
	if opts.Tag == "" {
		opts.Tag = "latest"
	}
	if opts.Port == 0 {
		opts.Port = 80
	}
	if opts.HTTPSPort == 0 {
		opts.HTTPSPort = 443
	}
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	imageRef := cli.Sprintf("%s:%s", opts.ImageName, opts.Tag)

	args := []string{"run"}

	if opts.Detach {
		args = append(args, "-d")
	} else {
		args = append(args, "--rm")
	}

	if opts.ContainerName != "" {
		args = append(args, "--name", opts.ContainerName)
	}

	// Port mappings
	args = append(args, "-p", cli.Sprintf("%d:80", opts.Port))
	args = append(args, "-p", cli.Sprintf("%d:443", opts.HTTPSPort))

	// Environment file
	if opts.EnvFile != "" {
		args = append(args, "--env-file", opts.EnvFile)
	}

	// Volume mounts
	for hostPath, containerPath := range opts.Volumes {
		args = append(args, "-v", cli.Sprintf("%s:%s", hostPath, containerPath))
	}

	args = append(args, imageRef)

	cmd := exec.CommandContext(ctx, "docker", args...)

	if opts.Detach {
		cmd.Stderr = opts.Output
		output, err := cmd.Output()
		if err != nil {
			return phpWrapAction(err, "start", "container")
		}
		containerID := core.Trim(string(output))
		cli.Print("Container started: %s\n", containerID[:12])
		return nil
	}

	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output
	return cmd.Run()
}

// Shell opens a shell in a running container.
func Shell(ctx context.Context, containerID string) error { // Result boundary
	if containerID == "" {
		return phpFailure("container ID is required")
	}

	// Resolve partial container ID
	fullID, err := resolveDockerContainerID(ctx, containerID)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "docker", "exec", "-it", fullID, "/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsPHPProject checks if the given directory is a PHP project.
func IsPHPProject(dir string) bool {
	composerPath := core.PathJoin(dir, composerJSONFile)
	return getMedium().IsFile(composerPath)
}

// commonLinuxKitPaths defines default search locations for linuxkit.
var commonLinuxKitPaths = []string{
	"/usr/local/bin/linuxkit",
	"/opt/homebrew/bin/linuxkit",
}

// lookupLinuxKit finds the linuxkit binary.
func lookupLinuxKit() (string, error) { // Result boundary
	// Check PATH first
	if path, err := exec.LookPath("linuxkit"); err == nil {
		return path, nil
	}

	m := getMedium()
	for _, p := range commonLinuxKitPaths {
		if m.IsFile(p) {
			return p, nil
		}
	}

	return "", phpFailure("linuxkit not found. Install with: brew install linuxkit (macOS) or see https://github.com/linuxkit/linuxkit")
}

// getLinuxKitTemplate retrieves a LinuxKit template by name.
func getLinuxKitTemplate(name string) (string, error) { // Result boundary
	// Default server-php template for PHP projects
	if name == defaultLinuxKitTemplateName {
		return defaultServerPHPTemplate, nil
	}

	// Try to load from container package templates
	// This would integrate with forge.lthn.ai/core/go/pkg/container
	return "", phpFailure("template not found: %s", name)
}

// applyTemplateVariables applies variable substitution to template content.
func applyTemplateVariables(content string, vars map[string]string) (string, error) { // Result boundary
	result := content
	for key, value := range vars {
		placeholder := "${" + key + "}"
		result = core.Replace(result, placeholder, value)
	}
	return result, nil
}

// resolveDockerContainerID resolves a partial container ID to a full ID.
func resolveDockerContainerID(ctx context.Context, partialID string) (string, error) { // Result boundary
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--no-trunc", "--format", "{{.ID}}")
	output, err := cmd.Output()
	if err != nil {
		return "", phpWrapAction(err, "list", "containers")
	}

	lines := core.Split(core.Trim(string(output)), "\n")
	var matches []string

	for _, line := range lines {
		if core.HasPrefix(line, partialID) {
			matches = append(matches, line)
		}
	}

	switch len(matches) {
	case 0:
		return "", phpFailure("no container found matching: %s", partialID)
	case 1:
		return matches[0], nil
	default:
		return "", phpFailure("multiple containers match '%s', be more specific", partialID)
	}
}

// defaultServerPHPTemplate is the default LinuxKit template for PHP servers.
const defaultServerPHPTemplate = `# LinuxKit configuration for PHP/FrankenPHP server
kernel:
  image: linuxkit/kernel:6.6.13
  cmdline: "console=tty0 console=ttyS0"
init:
  - linuxkit/init:v1.0.1
  - linuxkit/runc:v1.0.1
  - linuxkit/containerd:v1.0.1
onboot:
  - name: sysctl
    image: linuxkit/sysctl:v1.0.1
  - name: dhcpcd
    image: linuxkit/dhcpcd:v1.0.1
    command: ["/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf"]
services:
  - name: getty
    image: linuxkit/getty:v1.0.1
    env:
      - INSECURE=true
  - name: sshd
    image: linuxkit/sshd:v1.0.1
files:
  - path: etc/ssh/authorized_keys
    contents: |
      ${SSH_KEY:-}
`
