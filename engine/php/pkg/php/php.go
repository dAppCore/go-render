package php

import (
	"context"
	"io"
	"sync"
	"time"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// Options configures the development server.
type Options struct {
	// Dir is the Laravel project directory.
	Dir string

	// Services specifies which services to start.
	// If empty, services are auto-detected.
	Services []DetectedService

	// NoVite disables the Vite dev server.
	NoVite bool

	// NoHorizon disables Laravel Horizon.
	NoHorizon bool

	// NoReverb disables Laravel Reverb.
	NoReverb bool

	// NoRedis disables the Redis server.
	NoRedis bool

	// HTTPS enables HTTPS with mkcert certificates.
	HTTPS bool

	// Domain is the domain for SSL certificates.
	// Defaults to APP_URL from .env or "localhost".
	Domain string

	// Ports for each service
	FrankenPHPPort int
	HTTPSPort      int
	VitePort       int
	ReverbPort     int
	RedisPort      int
}

// DevServer manages all development services.
type DevServer struct {
	opts     Options
	services []Service
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
	running  bool
}

// NewDevServer creates a new development server manager.
func NewDevServer(opts Options) *DevServer {
	return &DevServer{
		opts:     opts,
		services: make([]Service, 0),
	}
}

// Start starts all detected/configured services.
func (d *DevServer) Start(ctx context.Context, opts Options) error { // Result boundary
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.running {
		return phpFailure("dev server is already running")
	}

	if err := d.applyStartOptions(opts); err != nil {
		return err
	}

	// Verify this is a Laravel project
	if !IsLaravelProject(d.opts.Dir) {
		return phpFailure("not a Laravel project: %s", d.opts.Dir)
	}

	// Create cancellable context
	d.ctx, d.cancel = context.WithCancel(ctx)

	// Setup SSL if HTTPS is enabled
	certFile, keyFile, err := setupDevSSL(d.opts.Dir, opts)
	if err != nil {
		return err
	}

	// Create services
	d.services = d.createServices(opts, certFile, keyFile)

	// Start all services
	if err := d.startServices(); err != nil {
		return err
	}

	d.running = true
	return nil
}

func (d *DevServer) applyStartOptions(opts Options) error { // Result boundary
	if opts.Dir != "" {
		d.opts.Dir = opts.Dir
	}
	if d.opts.Dir != "" {
		return nil
	}

	cwdResult := core.Getwd()
	if !cwdResult.OK {
		err, _ := cwdResult.Value.(error)
		return phpWrapAction(err, "get", workingDirectorySubject)
	}
	d.opts.Dir, _ = cwdResult.Value.(string)
	return nil
}

func setupDevSSL(dir string, opts Options) (string, string, error) { // Result boundary
	if !opts.HTTPS {
		return "", "", nil
	}

	certFile, keyFile, err := SetupSSLIfNeeded(devSSLDomain(dir, opts.Domain), SSLOptions{})
	if err != nil {
		return "", "", phpWrapAction(err, "setup", "SSL")
	}

	return certFile, keyFile, nil
}

func devSSLDomain(dir, configuredDomain string) string {
	if configuredDomain != "" {
		return configuredDomain
	}
	if appURL := GetLaravelAppURL(dir); appURL != "" {
		return ExtractDomainFromURL(appURL)
	}
	return "localhost"
}

func (d *DevServer) createServices(opts Options, certFile, keyFile string) []Service {
	services := opts.Services
	if len(services) == 0 {
		services = DetectServices(d.opts.Dir)
	}
	services = d.filterServices(services, opts)

	result := make([]Service, 0, len(services))
	for _, svc := range services {
		if service := d.createService(svc, opts, certFile, keyFile); service != nil {
			result = append(result, service)
		}
	}

	return result
}

func (d *DevServer) createService(svc DetectedService, opts Options, certFile, keyFile string) Service {
	switch svc {
	case ServiceFrankenPHP:
		return NewFrankenPHPService(d.opts.Dir, FrankenPHPOptions{
			Port:      defaultPort(opts.FrankenPHPPort, 8000),
			HTTPSPort: defaultPort(opts.HTTPSPort, 443),
			HTTPS:     opts.HTTPS,
			CertFile:  certFile,
			KeyFile:   keyFile,
		})
	case ServiceVite:
		return NewViteService(d.opts.Dir, ViteOptions{Port: defaultPort(opts.VitePort, 5173)})
	case ServiceHorizon:
		return NewHorizonService(d.opts.Dir)
	case ServiceReverb:
		return NewReverbService(d.opts.Dir, ReverbOptions{Port: defaultPort(opts.ReverbPort, 8080)})
	case ServiceRedis:
		return NewRedisService(d.opts.Dir, RedisOptions{Port: defaultPort(opts.RedisPort, 6379)})
	default:
		return nil
	}
}

func defaultPort(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func (d *DevServer) startServices() error { // Result boundary
	var startErrors []error
	for _, svc := range d.services {
		if err := svc.Start(d.ctx); err != nil {
			startErrors = append(startErrors, phpFailure("%s: %v", svc.Name(), err))
		}
	}

	if len(startErrors) == 0 {
		return nil
	}

	for _, svc := range d.services {
		if err := svc.Stop(); err != nil {
			startErrors = append(startErrors, phpFailure("cleanup %s: %v", svc.Name(), err))
		}
	}
	return phpFailure("failed to start services: %v", startErrors)
}

// filterServices removes disabled services from the list.
func (d *DevServer) filterServices(services []DetectedService, opts Options) []DetectedService {
	filtered := make([]DetectedService, 0)

	for _, svc := range services {
		switch svc {
		case ServiceVite:
			if !opts.NoVite {
				filtered = append(filtered, svc)
			}
		case ServiceHorizon:
			if !opts.NoHorizon {
				filtered = append(filtered, svc)
			}
		case ServiceReverb:
			if !opts.NoReverb {
				filtered = append(filtered, svc)
			}
		case ServiceRedis:
			if !opts.NoRedis {
				filtered = append(filtered, svc)
			}
		default:
			filtered = append(filtered, svc)
		}
	}

	return filtered
}

// Stop stops all services gracefully.
func (d *DevServer) Stop() error { // Result boundary
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return nil
	}

	// Cancel context first
	if d.cancel != nil {
		d.cancel()
	}

	// Stop all services in reverse order
	var stopErrors []error
	for i := len(d.services) - 1; i >= 0; i-- {
		svc := d.services[i]
		if err := svc.Stop(); err != nil {
			stopErrors = append(stopErrors, phpFailure("%s: %v", svc.Name(), err))
		}
	}

	d.running = false

	if len(stopErrors) > 0 {
		return phpFailure("errors stopping services: %v", stopErrors)
	}

	return nil
}

// Logs returns a reader for the specified service's logs.
// If service is empty, returns unified logs from all services.
func (d *DevServer) Logs(service string, follow bool) (io.ReadCloser, error) { // Result boundary
	d.mu.RLock()
	defer d.mu.RUnlock()

	if service == "" {
		// Return unified logs
		return d.unifiedLogs(follow)
	}

	// Find specific service
	for _, svc := range d.services {
		if svc.Name() == service {
			return svc.Logs(follow)
		}
	}

	return nil, phpFailure("service not found: %s", service)
}

// unifiedLogs creates a reader that combines logs from all services.
func (d *DevServer) unifiedLogs(follow bool) (io.ReadCloser, error) { // Result boundary
	readers := make([]io.ReadCloser, 0)

	for _, svc := range d.services {
		reader, err := svc.Logs(follow)
		if err != nil {
			// Close any readers we already opened
			var closeErrors []error
			for _, r := range readers {
				if closeErr := r.Close(); closeErr != nil {
					closeErrors = append(closeErrors, closeErr)
				}
			}
			if len(closeErrors) > 0 {
				return nil, phpFailure("failed to get logs for %s: %v; failed to close readers: %v", svc.Name(), err, closeErrors)
			}
			return nil, phpFailure("failed to get logs for %s: %v", svc.Name(), err)
		}
		readers = append(readers, reader)
	}

	return newMultiServiceReader(d.services, readers, follow), nil
}

// Status returns the status of all services.
func (d *DevServer) Status() []ServiceStatus {
	d.mu.RLock()
	defer d.mu.RUnlock()

	statuses := make([]ServiceStatus, 0, len(d.services))
	for _, svc := range d.services {
		statuses = append(statuses, svc.Status())
	}

	return statuses
}

// IsRunning returns true if the dev server is running.
func (d *DevServer) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}

// Services returns the list of managed services.
func (d *DevServer) Services() []Service {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.services
}

// multiServiceReader combines multiple service log readers.
type multiServiceReader struct {
	services []Service
	readers  []io.ReadCloser
	follow   bool
	closed   bool
	mu       sync.RWMutex
}

func newMultiServiceReader(services []Service, readers []io.ReadCloser, follow bool) *multiServiceReader {
	return &multiServiceReader{
		services: services,
		readers:  readers,
		follow:   follow,
	}
}

func (m *multiServiceReader) Read(p []byte) (n int, err error) { // Result boundary
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return 0, io.EOF
	}
	m.mu.RUnlock()

	// Round-robin read from all readers
	for i, reader := range m.readers {
		buf := make([]byte, len(p))
		n, err := reader.Read(buf)
		if n > 0 {
			// Prefix with service name
			prefix := cli.Sprintf("[%s] ", m.services[i].Name())
			copy(p, prefix)
			copy(p[len(prefix):], buf[:n])
			return n + len(prefix), nil
		}
		if err != nil && err != io.EOF {
			return 0, err
		}
	}

	if m.follow {
		time.Sleep(100 * time.Millisecond)
		return 0, nil
	}

	return 0, io.EOF
}

func (m *multiServiceReader) Close() error { // Result boundary
	m.mu.Lock()
	m.closed = true
	m.mu.Unlock()

	var closeErr error
	for _, reader := range m.readers {
		if err := reader.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	return closeErr
}
