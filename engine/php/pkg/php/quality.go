package php

import (
	"context"
	goio "io"
	"os"
	"os/exec"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// FormatOptions configures PHP code formatting.
type FormatOptions struct {
	// Dir is the project directory (defaults to current working directory).
	Dir string

	// Fix automatically fixes formatting issues.
	Fix bool

	// Diff shows a diff of changes instead of modifying files.
	Diff bool

	// JSON outputs results in JSON format.
	JSON bool

	// Paths limits formatting to specific paths.
	Paths []string

	// Output is the writer for output (defaults to os.Stdout).
	Output goio.Writer
}

// AnalyseOptions configures PHP static analysis.
type AnalyseOptions struct {
	// Dir is the project directory (defaults to current working directory).
	Dir string

	// Level is the PHPStan analysis level (0-9).
	Level int

	// Paths limits analysis to specific paths.
	Paths []string

	// Memory is the memory limit for analysis (e.g., "2G").
	Memory string

	// JSON outputs results in JSON format.
	JSON bool

	// SARIF outputs results in SARIF format for GitHub Security tab.
	SARIF bool

	// Output is the writer for output (defaults to os.Stdout).
	Output goio.Writer
}

// FormatterType represents the detected formatter.
type FormatterType string

// Formatter type constants.
const (
	// FormatterPint indicates Laravel Pint code formatter.
	FormatterPint FormatterType = "pint"
)

// AnalyserType represents the detected static analyser.
type AnalyserType string

// Static analyser type constants.
const (
	// AnalyserPHPStan indicates standard PHPStan analyser.
	AnalyserPHPStan AnalyserType = "phpstan"
	// AnalyserLarastan indicates Laravel-specific Larastan analyser.
	AnalyserLarastan AnalyserType = "larastan"
)

// DetectFormatter detects which formatter is available in the project.
func DetectFormatter(dir string) (FormatterType, bool) {
	m := getMedium()

	// Check for Pint config
	pintConfig := core.PathJoin(dir, "pint.json")
	if m.Exists(pintConfig) {
		return FormatterPint, true
	}

	// Check for vendor binary
	pintBin := core.PathJoin(dir, "vendor", "bin", "pint")
	if m.Exists(pintBin) {
		return FormatterPint, true
	}

	return "", false
}

// DetectAnalyser detects which static analyser is available in the project.
func DetectAnalyser(dir string) (AnalyserType, bool) {
	m := getMedium()

	// Check for PHPStan config
	phpstanConfig := core.PathJoin(dir, "phpstan.neon")
	phpstanDistConfig := core.PathJoin(dir, "phpstan.neon.dist")

	hasConfig := m.Exists(phpstanConfig) || m.Exists(phpstanDistConfig)

	// Check for vendor binary
	phpstanBin := core.PathJoin(dir, "vendor", "bin", "phpstan")
	hasBin := m.Exists(phpstanBin)

	if hasConfig || hasBin {
		// Check if it's Larastan (Laravel-specific PHPStan)
		larastanPath := core.PathJoin(dir, "vendor", "larastan", "larastan")
		if m.Exists(larastanPath) {
			return AnalyserLarastan, true
		}
		// Also check nunomaduro/larastan
		larastanPath2 := core.PathJoin(dir, "vendor", "nunomaduro", "larastan")
		if m.Exists(larastanPath2) {
			return AnalyserLarastan, true
		}
		return AnalyserPHPStan, true
	}

	return "", false
}

// Format runs Laravel Pint to format PHP code.
func Format(ctx context.Context, opts FormatOptions) error { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	// Check if formatter is available
	formatter, found := DetectFormatter(opts.Dir)
	if !found {
		return phpFailure("no formatter found (install Laravel Pint: composer require laravel/pint --dev)")
	}

	var cmdName string
	var args []string

	switch formatter {
	case FormatterPint:
		cmdName, args = buildPintCommand(opts)
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	return cmd.Run()
}

// Analyse runs PHPStan or Larastan for static analysis.
func Analyse(ctx context.Context, opts AnalyseOptions) error { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	// Check if analyser is available
	analyser, found := DetectAnalyser(opts.Dir)
	if !found {
		return phpFailure("no static analyser found (install PHPStan: composer require phpstan/phpstan --dev)")
	}

	var cmdName string
	var args []string

	switch analyser {
	case AnalyserPHPStan, AnalyserLarastan:
		cmdName, args = buildPHPStanCommand(opts)
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	return cmd.Run()
}

// buildPintCommand builds the command for running Laravel Pint.
func buildPintCommand(opts FormatOptions) (string, []string) {
	m := getMedium()

	// Check for vendor binary first
	vendorBin := core.PathJoin(opts.Dir, "vendor", "bin", "pint")
	cmdName := "pint"
	if m.Exists(vendorBin) {
		cmdName = vendorBin
	}

	var args []string

	if !opts.Fix {
		args = append(args, "--test")
	}

	if opts.Diff {
		args = append(args, "--diff")
	}

	if opts.JSON {
		args = append(args, "--format=json")
	}

	// Add specific paths if provided
	args = append(args, opts.Paths...)

	return cmdName, args
}

// buildPHPStanCommand builds the command for running PHPStan.
func buildPHPStanCommand(opts AnalyseOptions) (string, []string) {
	m := getMedium()

	// Check for vendor binary first
	vendorBin := core.PathJoin(opts.Dir, "vendor", "bin", "phpstan")
	cmdName := "phpstan"
	if m.Exists(vendorBin) {
		cmdName = vendorBin
	}

	args := []string{"analyse"}

	if opts.Level > 0 {
		args = append(args, "--level", cli.Sprintf("%d", opts.Level))
	}

	if opts.Memory != "" {
		args = append(args, "--memory-limit", opts.Memory)
	}

	// Output format - SARIF takes precedence over JSON
	if opts.SARIF {
		args = append(args, "--error-format=sarif")
	} else if opts.JSON {
		args = append(args, "--error-format=json")
	}

	// Add specific paths if provided
	args = append(args, opts.Paths...)

	return cmdName, args
}

// =============================================================================
// Psalm Static Analysis
// =============================================================================

// PsalmOptions configures Psalm static analysis.
type PsalmOptions struct {
	Dir      string
	Level    int  // Error level (1=strictest, 8=most lenient)
	Fix      bool // Auto-fix issues where possible
	Baseline bool // Generate/update baseline file
	ShowInfo bool // Show info-level issues
	JSON     bool // Output in JSON format
	SARIF    bool // Output in SARIF format for GitHub Security tab
	Output   goio.Writer
}

// PsalmType represents the detected Psalm configuration.
type PsalmType string

// Psalm configuration type constants.
const (
	// PsalmStandard indicates standard Psalm configuration.
	PsalmStandard PsalmType = "psalm"
)

// DetectPsalm checks if Psalm is available in the project.
func DetectPsalm(dir string) (PsalmType, bool) {
	m := getMedium()

	// Check for psalm.xml config
	psalmConfig := core.PathJoin(dir, "psalm.xml")
	psalmDistConfig := core.PathJoin(dir, "psalm.xml.dist")

	hasConfig := m.Exists(psalmConfig) || m.Exists(psalmDistConfig)

	// Check for vendor binary
	psalmBin := core.PathJoin(dir, "vendor", "bin", "psalm")
	if m.Exists(psalmBin) {
		return PsalmStandard, true
	}

	if hasConfig {
		return PsalmStandard, true
	}

	return "", false
}

// RunPsalm runs Psalm static analysis.
func RunPsalm(ctx context.Context, opts PsalmOptions) error { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	m := getMedium()

	// Build command
	vendorBin := core.PathJoin(opts.Dir, "vendor", "bin", "psalm")
	cmdName := "psalm"
	if m.Exists(vendorBin) {
		cmdName = vendorBin
	}

	args := []string{"--no-progress"}

	if opts.Level > 0 && opts.Level <= 8 {
		args = append(args, cli.Sprintf("--error-level=%d", opts.Level))
	}

	if opts.Fix {
		args = append(args, "--alter", "--issues=all")
	}

	if opts.Baseline {
		args = append(args, "--set-baseline=psalm-baseline.xml")
	}

	if opts.ShowInfo {
		args = append(args, "--show-info=true")
	}

	// Output format - SARIF takes precedence over JSON
	if opts.SARIF {
		args = append(args, "--output-format=sarif")
	} else if opts.JSON {
		args = append(args, "--output-format=json")
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	return cmd.Run()
}

// =============================================================================
// Security Audit
// =============================================================================

// AuditOptions configures dependency security auditing.
type AuditOptions struct {
	Dir    string
	JSON   bool // Output in JSON format
	Fix    bool // Auto-fix vulnerabilities (npm only)
	Output goio.Writer
}

// AuditResult holds the results of a security audit.
type AuditResult struct {
	Tool            string
	Vulnerabilities int
	Advisories      []AuditAdvisory
	Error           error
}

// AuditAdvisory represents a single security advisory.
type AuditAdvisory struct {
	Package     string
	Severity    string
	Title       string
	URL         string
	Identifiers []string
}

// RunAudit runs security audits on dependencies.
func RunAudit(ctx context.Context, opts AuditOptions) ([]AuditResult, error) { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return nil, phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	var results []AuditResult

	// Run composer audit
	composerResult := runComposerAudit(ctx, opts)
	results = append(results, composerResult)

	// Run npm audit if package.json exists
	if getMedium().Exists(core.PathJoin(opts.Dir, packageJSONFile)) {
		npmResult := runNpmAudit(ctx, opts)
		results = append(results, npmResult)
	}

	return results, nil
}

func runComposerAudit(ctx context.Context, opts AuditOptions) AuditResult {
	result := AuditResult{Tool: "composer"}

	args := []string{"audit", "--format=json"}

	cmd := exec.CommandContext(ctx, "composer", args...)
	cmd.Dir = opts.Dir

	output, err := cmd.Output()
	if err != nil {
		// composer audit returns non-zero if vulnerabilities found
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = append(output, exitErr.Stderr...)
		}
	}

	// Parse JSON output
	var auditData struct {
		Advisories map[string][]struct {
			Title          string `json:"title"`
			Link           string `json:"link"`
			CVE            string `json:"cve"`
			AffectedRanges string `json:"affectedVersions"`
		} `json:"advisories"`
	}

	if jsonR := core.JSONUnmarshal(output, &auditData); jsonR.OK {
		for pkg, advisories := range auditData.Advisories {
			for _, adv := range advisories {
				result.Advisories = append(result.Advisories, AuditAdvisory{
					Package:     pkg,
					Title:       adv.Title,
					URL:         adv.Link,
					Identifiers: []string{adv.CVE},
				})
			}
		}
		result.Vulnerabilities = len(result.Advisories)
	} else if err != nil {
		result.Error = err
	}

	return result
}

func runNpmAudit(ctx context.Context, opts AuditOptions) AuditResult {
	result := AuditResult{Tool: "npm"}

	args := []string{"audit", "--json"}
	if opts.Fix {
		args = []string{"audit", "fix"}
	}

	cmd := exec.CommandContext(ctx, "npm", args...)
	cmd.Dir = opts.Dir

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = append(output, exitErr.Stderr...)
		}
	}

	if !opts.Fix {
		// Parse JSON output
		var auditData struct {
			Metadata struct {
				Vulnerabilities struct {
					Total int `json:"total"`
				} `json:"vulnerabilities"`
			} `json:"metadata"`
			Vulnerabilities map[string]struct {
				Severity string `json:"severity"`
				Via      []any  `json:"via"`
			} `json:"vulnerabilities"`
		}

		if jsonR := core.JSONUnmarshal(output, &auditData); jsonR.OK {
			result.Vulnerabilities = auditData.Metadata.Vulnerabilities.Total
			for pkg, vuln := range auditData.Vulnerabilities {
				result.Advisories = append(result.Advisories, AuditAdvisory{
					Package:  pkg,
					Severity: vuln.Severity,
				})
			}
		} else if err != nil {
			result.Error = err
		}
	}

	return result
}

// =============================================================================
// Rector Automated Refactoring
// =============================================================================

// RectorOptions configures Rector code refactoring.
type RectorOptions struct {
	Dir        string
	Fix        bool // Apply changes (default is dry-run)
	Diff       bool // Show detailed diff
	ClearCache bool // Clear cache before running
	Output     goio.Writer
}

// DetectRector checks if Rector is available in the project.
func DetectRector(dir string) bool {
	m := getMedium()

	// Check for rector.php config
	rectorConfig := core.PathJoin(dir, "rector.php")
	if m.Exists(rectorConfig) {
		return true
	}

	// Check for vendor binary
	rectorBin := core.PathJoin(dir, "vendor", "bin", "rector")
	return m.Exists(rectorBin)
}

// RunRector runs Rector for automated code refactoring.
func RunRector(ctx context.Context, opts RectorOptions) error { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	m := getMedium()

	// Build command
	vendorBin := core.PathJoin(opts.Dir, "vendor", "bin", "rector")
	cmdName := "rector"
	if m.Exists(vendorBin) {
		cmdName = vendorBin
	}

	args := []string{"process"}

	if !opts.Fix {
		args = append(args, "--dry-run")
	}

	if opts.Diff {
		args = append(args, "--output-format", "diff")
	}

	if opts.ClearCache {
		args = append(args, "--clear-cache")
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	return cmd.Run()
}

// =============================================================================
// Infection Mutation Testing
// =============================================================================

// InfectionOptions configures Infection mutation testing.
type InfectionOptions struct {
	Dir           string
	MinMSI        int    // Minimum mutation score indicator (0-100)
	MinCoveredMSI int    // Minimum covered mutation score (0-100)
	Threads       int    // Number of parallel threads
	Filter        string // Filter files by pattern
	OnlyCovered   bool   // Only mutate covered code
	Output        goio.Writer
}

// DetectInfection checks if Infection is available in the project.
func DetectInfection(dir string) bool {
	m := getMedium()

	// Check for infection config files
	configs := []string{"infection.json", "infection.json5", "infection.json.dist"}
	for _, config := range configs {
		if m.Exists(core.PathJoin(dir, config)) {
			return true
		}
	}

	// Check for vendor binary
	infectionBin := core.PathJoin(dir, "vendor", "bin", "infection")
	return m.Exists(infectionBin)
}

// RunInfection runs Infection mutation testing.
func RunInfection(ctx context.Context, opts InfectionOptions) error { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	m := getMedium()

	// Build command
	vendorBin := core.PathJoin(opts.Dir, "vendor", "bin", "infection")
	cmdName := "infection"
	if m.Exists(vendorBin) {
		cmdName = vendorBin
	}

	var args []string

	// Set defaults
	minMSI := opts.MinMSI
	if minMSI == 0 {
		minMSI = 50
	}
	minCoveredMSI := opts.MinCoveredMSI
	if minCoveredMSI == 0 {
		minCoveredMSI = 70
	}
	threads := opts.Threads
	if threads == 0 {
		threads = 4
	}

	args = append(args, cli.Sprintf("--min-msi=%d", minMSI))
	args = append(args, cli.Sprintf("--min-covered-msi=%d", minCoveredMSI))
	args = append(args, cli.Sprintf("--threads=%d", threads))

	if opts.Filter != "" {
		args = append(args, "--filter="+opts.Filter)
	}

	if opts.OnlyCovered {
		args = append(args, "--only-covered")
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Output
	cmd.Stderr = opts.Output

	return cmd.Run()
}

// =============================================================================
// QA Pipeline
// =============================================================================

// QAOptions configures the full QA pipeline.
type QAOptions struct {
	Dir   string
	Quick bool // Only run quick checks
	Full  bool // Run all stages including slow checks
	Fix   bool // Auto-fix issues where possible
	JSON  bool // Output results as JSON
}

// QAStage represents a stage in the QA pipeline.
type QAStage string

// QA pipeline stage constants.
const (
	// QAStageQuick runs fast checks only (audit, fmt, stan).
	QAStageQuick QAStage = "quick"
	// QAStageStandard runs standard checks including tests.
	QAStageStandard QAStage = "standard"
	// QAStageFull runs all checks including slow security scans.
	QAStageFull QAStage = "full"
)

// QACheckResult holds the result of a single QA check.
type QACheckResult struct {
	Name     string
	Stage    QAStage
	Passed   bool
	Duration string
	Error    error
	Output   string
}

// QAResult holds the results of the full QA pipeline.
type QAResult struct {
	Stages  []QAStage
	Checks  []QACheckResult
	Passed  bool
	Summary string
}

// GetQAStages returns the stages to run based on options.
func GetQAStages(opts QAOptions) []QAStage {
	if opts.Quick {
		return []QAStage{QAStageQuick}
	}
	if opts.Full {
		return []QAStage{QAStageQuick, QAStageStandard, QAStageFull}
	}
	// Default: quick + standard
	return []QAStage{QAStageQuick, QAStageStandard}
}

// GetQAChecks returns the checks for a given stage.
func GetQAChecks(dir string, stage QAStage) []string {
	switch stage {
	case QAStageQuick:
		checks := []string{"audit", `fmt`, "stan"}
		return checks
	case QAStageStandard:
		checks := []string{}
		if _, found := DetectPsalm(dir); found {
			checks = append(checks, "psalm")
		}
		checks = append(checks, "test")
		return checks
	case QAStageFull:
		checks := []string{}
		if DetectRector(dir) {
			checks = append(checks, "rector")
		}
		if DetectInfection(dir) {
			checks = append(checks, "infection")
		}
		return checks
	}
	return nil
}

// =============================================================================
// Security Checks
// =============================================================================

// SecurityOptions configures security scanning.
type SecurityOptions struct {
	Dir      string
	Severity string // Minimum severity (critical, high, medium, low)
	JSON     bool   // Output in JSON format
	SARIF    bool   // Output in SARIF format
	URL      string // URL to check HTTP headers (optional)
	Output   goio.Writer
}

// SecurityResult holds the results of security scanning.
type SecurityResult struct {
	Checks  []SecurityCheck
	Summary SecuritySummary
}

// SecurityCheck represents a single security check result.
type SecurityCheck struct {
	ID          string
	Name        string
	Description string
	Severity    string
	Passed      bool
	Message     string
	Fix         string
	CWE         string
}

// SecuritySummary summarizes security check results.
type SecuritySummary struct {
	Total    int
	Passed   int
	Critical int
	High     int
	Medium   int
	Low      int
}

// RunSecurityChecks runs security checks on the project.
func RunSecurityChecks(ctx context.Context, opts SecurityOptions) (*SecurityResult, error) { // Result boundary
	if opts.Dir == "" {
		cwdR := core.Getwd()
		if !cwdR.OK {
			return nil, phpWrapAction(cwdR.Value.(error), "get", workingDirectorySubject)
		}
		opts.Dir = cwdR.Value.(string)
	}

	result := &SecurityResult{}

	// Run composer audit
	auditResults, _ := RunAudit(ctx, AuditOptions{Dir: opts.Dir})
	for _, audit := range auditResults {
		check := SecurityCheck{
			ID:          audit.Tool + "_audit",
			Name:        phpTitle(audit.Tool) + " Security Audit",
			Description: "Check " + audit.Tool + " dependencies for vulnerabilities",
			Severity:    "critical",
			Passed:      audit.Vulnerabilities == 0 && audit.Error == nil,
			CWE:         "CWE-1395",
		}
		if !check.Passed {
			check.Message = cli.Sprintf("Found %d vulnerabilities", audit.Vulnerabilities)
		}
		result.Checks = append(result.Checks, check)
	}

	// Check .env file for security issues
	envChecks := runEnvSecurityChecks(opts.Dir)
	result.Checks = append(result.Checks, envChecks...)

	// Check filesystem security
	fsChecks := runFilesystemSecurityChecks(opts.Dir)
	result.Checks = append(result.Checks, fsChecks...)

	// Calculate summary
	for _, check := range result.Checks {
		result.Summary.Total++
		if check.Passed {
			result.Summary.Passed++
		} else {
			switch check.Severity {
			case "critical":
				result.Summary.Critical++
			case "high":
				result.Summary.High++
			case "medium":
				result.Summary.Medium++
			case "low":
				result.Summary.Low++
			}
		}
	}

	return result, nil
}

func runEnvSecurityChecks(dir string) []SecurityCheck {
	envMap, ok := readEnvFileMap(dir)
	if !ok {
		return nil
	}

	var checks []SecurityCheck
	checks = appendEnvCheck(checks, envMap, "APP_DEBUG", debugModeCheck)
	checks = appendEnvCheck(checks, envMap, "APP_KEY", appKeyCheck)
	checks = appendEnvCheck(checks, envMap, "APP_URL", httpsEnforcedCheck)
	return checks
}

func readEnvFileMap(dir string) (map[string]string, bool) {
	m := getMedium()
	envPath := core.PathJoin(dir, ".env")
	envContent, err := m.Read(envPath)
	if err != nil {
		return nil, false
	}

	envLines := core.Split(envContent, "\n")
	envMap := make(map[string]string)
	for _, line := range envLines {
		line = core.Trim(line)
		if line == "" || core.HasPrefix(line, "#") {
			continue
		}
		parts := core.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	return envMap, true
}

func appendEnvCheck(checks []SecurityCheck, envMap map[string]string, key string, build func(string) SecurityCheck) []SecurityCheck {
	if value, ok := envMap[key]; ok {
		return append(checks, build(value))
	}
	return checks
}

func debugModeCheck(debug string) SecurityCheck {
	check := SecurityCheck{
		ID:          "debug_mode",
		Name:        "Debug Mode Disabled",
		Description: "APP_DEBUG should be false in production",
		Severity:    "critical",
		Passed:      core.Lower(debug) != "true",
		CWE:         "CWE-215",
	}
	if !check.Passed {
		check.Message = "Debug mode exposes sensitive information"
		check.Fix = "Set APP_DEBUG=false in .env"
	}
	return check
}

func appKeyCheck(key string) SecurityCheck {
	check := SecurityCheck{
		ID:          "app_key_set",
		Name:        "Application Key Set",
		Description: "APP_KEY must be set and valid",
		Severity:    "critical",
		Passed:      len(key) >= 32,
		CWE:         "CWE-321",
	}
	if !check.Passed {
		check.Message = "Missing or weak encryption key"
		check.Fix = "Run: php artisan key:generate"
	}
	return check
}

func httpsEnforcedCheck(url string) SecurityCheck {
	check := SecurityCheck{
		ID:          "https_enforced",
		Name:        "HTTPS Enforced",
		Description: "APP_URL should use HTTPS in production",
		Severity:    "high",
		Passed:      core.HasPrefix(url, "https://"),
		CWE:         "CWE-319",
	}
	if !check.Passed {
		check.Message = "Application not using HTTPS"
		check.Fix = "Update APP_URL to use https://"
	}
	return check
}

func runFilesystemSecurityChecks(dir string) []SecurityCheck {
	var checks []SecurityCheck
	m := getMedium()

	// Check .env not in public
	publicEnvPaths := []string{"public/.env", "public_html/.env"}
	for _, path := range publicEnvPaths {
		fullPath := core.PathJoin(dir, path)
		if m.Exists(fullPath) {
			checks = append(checks, SecurityCheck{
				ID:          "env_not_public",
				Name:        ".env Not Publicly Accessible",
				Description: ".env file should not be in public directory",
				Severity:    "critical",
				Passed:      false,
				Message:     "Environment file exposed to web at " + path,
				CWE:         "CWE-538",
			})
		}
	}

	// Check .git not in public
	publicGitPaths := []string{"public/.git", "public_html/.git"}
	for _, path := range publicGitPaths {
		fullPath := core.PathJoin(dir, path)
		if m.Exists(fullPath) {
			checks = append(checks, SecurityCheck{
				ID:          "git_not_public",
				Name:        ".git Not Publicly Accessible",
				Description: ".git directory should not be in public",
				Severity:    "critical",
				Passed:      false,
				Message:     "Git repository exposed to web (source code leak)",
				CWE:         "CWE-538",
			})
		}
	}

	return checks
}
