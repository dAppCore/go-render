package marketplace

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	core "dappco.re/go"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name        string            `yaml:"name" json:"name"`
	Version     string            `yaml:"version" json:"version"`
	Repository  string            `yaml:"repository" json:"repository"`
	Ref         string            `yaml:"ref" json:"ref"`
	Description string            `yaml:"description" json:"description"`
	Files       map[string]string `yaml:"files" json:"files"`
	Signature   Signature         `yaml:"signature" json:"signature"`
}

type Signature struct {
	Algorithm string `yaml:"algorithm" json:"algorithm"`
	PublicKey string `yaml:"public_key" json:"public_key"`
	Value     string `yaml:"value" json:"value"`
}

type Installer struct {
	HTTPClient *http.Client
	GitBinary  string
	InstallDir string
	GitRunner  func(context.Context, string, ...string) ([]byte, error)
}

const maxManifestBytes = 1 << 20

var credentialRedactionPattern = regexp.MustCompile(`(?i)\b([a-z][a-z0-9+.-]*://)([^@\s/]+)@`)

func (i Installer) FetchManifest(ctx context.Context, manifestURL string) (Manifest, resultFailure) {
	client := i.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL, nil)
	if err != nil {
		return Manifest{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Manifest{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return Manifest{}, core.Errorf("manifest fetch failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManifestBytes+1))
	if err != nil {
		return Manifest{}, err
	}
	if len(body) > maxManifestBytes {
		return Manifest{}, core.Errorf("manifest fetch failed: manifest exceeds %d bytes", maxManifestBytes)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return Manifest{}, err
	}
	if manifest.Repository == "" {
		manifest.Repository = manifestURL
	}
	return manifest, nil
}

func VerifyManifest(manifest Manifest) resultFailure {
	if core.Lower(core.Trim(manifest.Signature.Algorithm)) != "ed25519" {
		return core.NewError("manifest signature algorithm must be ed25519")
	}
	if manifest.Signature.Value == "" || manifest.Signature.PublicKey == "" {
		return core.NewError("manifest signature is required")
	}
	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature, err := base64.StdEncoding.DecodeString(manifest.Signature.Value)
	if err != nil {
		return err
	}
	publicKey, err := base64.StdEncoding.DecodeString(manifest.Signature.PublicKey)
	if err != nil {
		return err
	}
	if len(signature) != ed25519.SignatureSize {
		return core.NewError("manifest signature has invalid size")
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return core.NewError("manifest public key has invalid size")
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), []byte(payload), signature) {
		return core.NewError("manifest signature verification failed")
	}
	return nil
}

func (i Installer) Verify(ctx context.Context, manifestURL string) (Manifest, resultFailure) {
	manifest, err := i.FetchManifest(ctx, manifestURL)
	if err != nil {
		return Manifest{}, err
	}
	if err := VerifyManifest(manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (i Installer) Install(ctx context.Context, manifest Manifest) (string, resultFailure) {
	if core.Trim(i.InstallDir) == "" {
		return "", core.NewError("install dir is required")
	}
	if err := VerifyManifest(manifest); err != nil {
		return "", err
	}
	if err := validateManifestName(manifest.Name); err != nil {
		return "", err
	}
	if err := validateRepositorySource(manifest.Repository); err != nil {
		return "", err
	}
	if err := validateCloneArgOptional("ref", manifest.Ref); err != nil {
		return "", err
	}
	if err := coreMkdirAll(i.InstallDir, 0o755); err != nil {
		return "", err
	}
	rootAbs, err := pathAbs(i.InstallDir)
	if err != nil {
		return "", err
	}
	rootResolved, err := pathEvalSymlinks(rootAbs)
	if err != nil {
		return "", err
	}
	targetDir := core.PathJoin(rootResolved, safeName(manifest.Name))
	targetAbs, err := pathAbs(targetDir)
	if err != nil {
		return "", err
	}
	cleanupTarget := true
	defer func() {
		if cleanupTarget {
			if err := coreRemoveAll(targetDir); err != nil {
				return
			}
		}
	}()

	rel, err := pathRel(rootResolved, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || core.HasPrefix(rel, ".."+string(core.PathSeparator)) {
		return "", core.NewError("install path escapes install dir")
	}
	if err := coreRemoveAll(targetDir); err != nil {
		return "", err
	}
	args := []string{"clone", "--depth", "1"}
	if manifest.Ref != "" {
		args = append(args, "--branch", manifest.Ref)
	}
	args = append(args, "--", manifest.Repository, targetDir)
	binary := i.GitBinary
	if core.Trim(binary) == "" {
		binary = "git"
	}
	runGit := i.GitRunner
	if runGit == nil {
		runGit = func(ctx context.Context, binary string, args ...string) ([]byte, error) {
			return runGitCommand(ctx, binary, args...)
		}
	}
	if output, err := runGit(ctx, binary, args...); err != nil {
		return "", core.Errorf("git clone failed: %w: %s", err, sanitizeCommandOutput(output))
	}
	if err := writeInstalledManifest(targetDir, manifest); err != nil {
		return "", err
	}
	cleanupTarget = false
	return targetDir, nil
}

func runGitCommand(ctx context.Context, binary string, args ...string) ([]byte, resultFailure) {
	cmd := commandContext(ctx, binary, args...)
	return cmd.CombinedOutput()
}

func (i Installer) List(ctx context.Context, registryURL string) ([]Manifest, resultFailure) {
	client := i.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, registryURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, core.Errorf("marketplace list failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManifestBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxManifestBytes {
		return nil, core.Errorf("marketplace list failed: payload exceeds %d bytes", maxManifestBytes)
	}
	return decodeManifestList(body)
}

func validateManifestName(value string) resultFailure {
	trimmed := core.Trim(value)
	if trimmed == "" {
		return core.NewError("manifest name is required")
	}
	if containsAny(trimmed, `/\`) {
		return core.NewError("manifest name must not contain path separators")
	}
	if core.Contains(trimmed, "..") {
		return core.NewError("manifest name must not contain path traversal segments")
	}
	return nil
}

func validateCloneArg(label, value string) resultFailure {
	trimmed := core.Trim(value)
	if trimmed == "" {
		return core.Errorf("%s is required", label)
	}
	if core.HasPrefix(trimmed, "-") {
		return core.Errorf("%s must not begin with a dash", label)
	}
	if containsAny(trimmed, "\x00\r\n") {
		return core.Errorf("%s contains invalid control characters", label)
	}
	return nil
}

func validateCloneArgOptional(label, value string) resultFailure {
	trimmed := core.Trim(value)
	if trimmed == "" {
		return nil
	}
	return validateCloneArg(label, trimmed)
}

func validateRepositorySource(value string) resultFailure {
	trimmed := core.Trim(value)
	if trimmed == "" {
		return core.NewError("repository is required")
	}
	if containsAny(trimmed, "\x00\r\n") {
		return core.NewError("repository contains invalid control characters")
	}
	if core.HasPrefix(core.Lower(trimmed), "ext::") {
		return core.NewError("repository must not use git remote helper protocols")
	}
	if core.HasPrefix(trimmed, "-") {
		return core.NewError("repository must not begin with a dash")
	}
	if core.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return err
		}
		switch core.Lower(parsed.Scheme) {
		case "http", "https", "ssh", "git":
		default:
			return core.Errorf("repository scheme %q is not allowed", parsed.Scheme)
		}
		return nil
	}
	if core.Contains(trimmed, string(core.PathSeparator)) || core.PathIsAbs(trimmed) {
		return core.NewError("repository path clones are not allowed")
	}
	if !core.Contains(trimmed, ":") {
		return core.NewError("repository must be a URL or scp-style remote")
	}
	return nil
}

func DigestManifest(manifest Manifest) string {
	hash := sha256.Sum256([]byte(manifest.Name + ":" + manifest.Version + ":" + manifest.Repository + ":" + manifest.Ref))
	return hex.EncodeToString(hash[:])
}

func safeName(value string) string {
	original := value
	value = core.Trim(core.Lower(value))
	if value == "" {
		return fallbackSafeName(original)
	}
	builder := core.NewBuilder()
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == '.':
			builder.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				builder.WriteRune('-')
				lastDash = true
			}
		}
	}
	cleaned := trimRunes(builder.String(), "-._")
	if cleaned == "" {
		return fallbackSafeName(original)
	}
	return cleaned
}

func fallbackSafeName(value string) string {
	hash := sha256.Sum256([]byte(value))
	return "module-" + hex.EncodeToString(hash[:])[:8]
}

func decodeManifestList(body []byte) ([]Manifest, resultFailure) {
	trimmed := core.Trim(string(body))
	if trimmed == "" {
		return nil, nil
	}
	var manifests []Manifest
	if core.HasPrefix(trimmed, "[") {
		if err := jsonUnmarshal(body, &manifests); err != nil {
			return nil, err
		}
		return manifests, nil
	}
	var wrapped struct {
		Manifests []Manifest `json:"manifests" yaml:"manifests"`
	}
	if err := yaml.Unmarshal(body, &wrapped); err == nil && wrapped.Manifests != nil {
		return wrapped.Manifests, nil
	}
	if err := yaml.Unmarshal(body, &manifests); err != nil {
		return nil, err
	}
	return manifests, nil
}

func writeInstalledManifest(targetDir string, manifest Manifest) resultFailure {
	manifestDir := core.PathJoin(targetDir, ".core")
	if err := coreMkdirAll(manifestDir, 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	return coreWriteFile(core.PathJoin(manifestDir, "marketplace.yaml"), data, 0o644)
}

func sanitizeCommandOutput(output []byte) string {
	trimmed := core.Trim(string(output))
	if trimmed == "" {
		return "command produced no output"
	}
	sanitized := credentialRedactionPattern.ReplaceAllString(trimmed, "$1[redacted]@")
	const maxOutputChars = 512
	if len(sanitized) > maxOutputChars {
		sanitized = sanitized[:maxOutputChars] + "..."
	}
	return sanitized
}
