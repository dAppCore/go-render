package php

import (
	"os/exec"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

const (
	// DefaultSSLDir is the default directory for SSL certificates.
	DefaultSSLDir = ".core/ssl"
)

// SSLOptions configures SSL certificate generation.
type SSLOptions struct {
	// Dir is the directory to store certificates.
	// Defaults to ~/.core/ssl/
	Dir string
}

// GetSSLDir returns the SSL directory, creating it if necessary.
func GetSSLDir(opts SSLOptions) (string, error) { // Result boundary
	m := getMedium()
	dir := opts.Dir
	if dir == "" {
		homeR := core.UserHomeDir()
		if !homeR.OK {
			return "", phpWrapAction(homeR.Value.(error), "get", "home directory")
		}
		dir = core.PathJoin(homeR.Value.(string), DefaultSSLDir)
	}

	if err := m.EnsureDir(dir); err != nil {
		return "", phpWrapAction(err, "create", "SSL directory")
	}

	return dir, nil
}

// CertPaths returns the paths to the certificate and key files for a domain.
func CertPaths(domain string, opts SSLOptions) (certFile, keyFile string, err error) { // Result boundary
	dir, err := GetSSLDir(opts)
	if err != nil {
		return "", "", err
	}

	certFile = core.PathJoin(dir, cli.Sprintf("%s.pem", domain))
	keyFile = core.PathJoin(dir, cli.Sprintf("%s-key.pem", domain))

	return certFile, keyFile, nil
}

// CertsExist checks if SSL certificates exist for the given domain.
func CertsExist(domain string, opts SSLOptions) bool {
	m := getMedium()
	certFile, keyFile, err := CertPaths(domain, opts)
	if err != nil {
		return false
	}

	if !m.IsFile(certFile) {
		return false
	}

	if !m.IsFile(keyFile) {
		return false
	}

	return true
}

// SetupSSL creates local SSL certificates using mkcert.
// It installs the local CA if not already installed and generates
// certificates for the given domain.
func SetupSSL(domain string, opts SSLOptions) error { // Result boundary
	// Check if mkcert is installed
	if _, err := exec.LookPath("mkcert"); err != nil {
		return phpFailure("mkcert is not installed. Install it with: brew install mkcert (macOS) or see https://github.com/FiloSottile/mkcert")
	}

	dir, err := GetSSLDir(opts)
	if err != nil {
		return err
	}

	// Install local CA (idempotent operation)
	installCmd := exec.Command("mkcert", "-install")
	if output, err := installCmd.CombinedOutput(); err != nil {
		return phpFailure("failed to install mkcert CA: %v\n%s", err, output)
	}

	// Generate certificates
	certFile := core.PathJoin(dir, cli.Sprintf("%s.pem", domain))
	keyFile := core.PathJoin(dir, cli.Sprintf("%s-key.pem", domain))

	// mkcert generates cert and key with specific naming
	genCmd := exec.Command("mkcert",
		"-cert-file", certFile,
		"-key-file", keyFile,
		domain,
		"localhost",
		"127.0.0.1",
		"::1",
	)

	if output, err := genCmd.CombinedOutput(); err != nil {
		return phpFailure("failed to generate certificates: %v\n%s", err, output)
	}

	return nil
}

// SetupSSLIfNeeded checks if certificates exist and creates them if not.
func SetupSSLIfNeeded(domain string, opts SSLOptions) (certFile, keyFile string, err error) { // Result boundary
	certFile, keyFile, err = CertPaths(domain, opts)
	if err != nil {
		return "", "", err
	}

	if !CertsExist(domain, opts) {
		if err := SetupSSL(domain, opts); err != nil {
			return "", "", err
		}
	}

	return certFile, keyFile, nil
}

// IsMkcertInstalled checks if mkcert is available in PATH.
func IsMkcertInstalled() bool {
	_, err := exec.LookPath("mkcert")
	return err == nil
}

// InstallMkcertCA installs the local CA for mkcert.
func InstallMkcertCA() error { // Result boundary
	if !IsMkcertInstalled() {
		return phpFailure("mkcert is not installed")
	}

	cmd := exec.Command("mkcert", "-install")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return phpFailure("failed to install mkcert CA: %v\n%s", err, output)
	}

	return nil
}

// GetMkcertCARoot returns the path to the mkcert CA root directory.
func GetMkcertCARoot() (string, error) { // Result boundary
	if !IsMkcertInstalled() {
		return "", phpFailure("mkcert is not installed")
	}

	cmd := exec.Command("mkcert", "-CAROOT")
	output, err := cmd.Output()
	if err != nil {
		return "", phpWrapAction(err, "get", "mkcert CA root")
	}

	return core.Trim(string(output)), nil
}
