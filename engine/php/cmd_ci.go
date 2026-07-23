// cmd_ci.go implements the 'php ci' command for CI/CD pipeline integration.
//
// Usage:
//   core php ci                 # Run full CI pipeline
//   core php ci --json          # Output combined JSON report
//   core php ci --summary       # Output markdown summary
//   core php ci --sarif         # Generate SARIF files
//   core php ci --upload-sarif  # Upload SARIF to GitHub Security
//   core php ci --fail-on=high  # Only fail on high+ severity

package php

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"

	core "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// CI command flags
var (
	ciJSON        bool
	ciSummary     bool
	ciSARIF       bool
	ciUploadSARIF bool
	ciFailOn      string
)

// CIResult represents the overall CI pipeline result
type CIResult struct {
	Passed    bool            `json:"passed"`
	ExitCode  int             `json:"exit_code"`
	Duration  string          `json:"duration"`
	StartedAt time.Time       `json:"started_at"`
	Checks    []CICheckResult `json:"checks"`
	Summary   CISummary       `json:"summary"`
	Artifacts []string        `json:"artifacts,omitempty"`
}

// CICheckResult represents an individual check result
type CICheckResult struct {
	Name     string `json:"name"`
	Status   string `json:"status"` // passed, failed, warning, skipped
	Duration string `json:"duration"`
	Details  string `json:"details,omitempty"`
	Issues   int    `json:"issues,omitempty"`
	Errors   int    `json:"errors,omitempty"`
	Warnings int    `json:"warnings,omitempty"`
}

// CISummary contains aggregate statistics
type CISummary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Warnings int `json:"warnings"`
	Skipped  int `json:"skipped"`
}

type ciCheckDefinition struct {
	name  string
	run   func(context.Context, string) (CICheckResult, error)
	sarif bool
}

func addPHPCICommand(c *core.Core, prefix string) {
	path := phpCommandPath(prefix, "ci")
	phpFailureorCommand(c, path, phpT("cmd.php.ci.short"), func(opts core.Options) error {
		line := phpCommandLineFor(path, opts)
		ciJSON = line.Bool(`json`)
		ciSummary = line.Bool("summary")
		ciSARIF = line.Bool("sarif")
		ciUploadSARIF = line.Bool("upload-sarif")
		ciFailOn = line.String("fail-on", "error")
		return runPHPCI()
	})
}

func runPHPCI() error { // Result boundary
	cwdR := core.Getwd()
	if !cwdR.OK {
		return core.E("php", phpT(i18nFailGetKey, workingDirectorySubject), cwdR.Value.(error))
	}
	cwd := cwdR.Value.(string)

	if !IsPHPProject(cwd) {
		return core.E("php.runCI", phpT("cmd.php.error.not_php"), nil)
	}

	startTime := time.Now()
	ctx := context.Background()

	checks := defaultCIChecks()

	result := CIResult{
		StartedAt: startTime,
		Passed:    true,
		Checks:    make([]CICheckResult, 0, len(checks)),
	}

	printCIHeader()
	artifacts := runCIChecks(ctx, cwd, checks, &result)

	result.Duration = time.Since(startTime).Round(time.Millisecond).String()
	result.Artifacts = artifacts
	result.ExitCode = ciExitCode(result.Passed)

	return outputCIResult(ctx, cwd, result, artifacts)
}

func defaultCIChecks() []ciCheckDefinition {
	return []ciCheckDefinition{
		{"test", runCITest, false},
		{"stan", runCIStan, true},
		{"psalm", runCIPsalm, true},
		{`fmt`, runCIFmt, false},
		{"audit", runCIAudit, false},
		{"security", runCISecurity, false},
	}
}

func printCIHeader() {
	if ciJSON {
		return
	}

	cli.Print(cliSingleLineFormat, cli.BoldStyle.Render("core php ci - QA Pipeline"))
	cli.Print("%s\n\n", repeat("─", 40))
}

func runCIChecks(ctx context.Context, cwd string, checks []ciCheckDefinition, result *CIResult) []string {
	var artifacts []string

	for _, check := range checks {
		printCICheckProgress(check.name)
		checkResult := executeCICheck(ctx, cwd, check)
		recordCICheckResult(result, checkResult)
		printCICheckResult(check.name, checkResult)
		artifacts = appendCIArtifact(ctx, cwd, check, artifacts)
	}

	return artifacts
}

func printCICheckProgress(name string) {
	if !ciJSON {
		cli.Print("  %s %s...", dimStyle.Render("→"), name)
	}
}

func executeCICheck(ctx context.Context, cwd string, check ciCheckDefinition) CICheckResult {
	checkResult, err := check.run(ctx, cwd)
	if err == nil {
		return checkResult
	}

	return CICheckResult{
		Name:    check.name,
		Status:  "failed",
		Details: err.Error(),
	}
}

func recordCICheckResult(result *CIResult, checkResult CICheckResult) {
	result.Checks = append(result.Checks, checkResult)
	result.Summary.Total++

	switch checkResult.Status {
	case "passed":
		result.Summary.Passed++
	case "failed":
		result.Summary.Failed++
		if shouldFailOn(checkResult, ciFailOn) {
			result.Passed = false
		}
	case "warning":
		result.Summary.Warnings++
	case "skipped":
		result.Summary.Skipped++
	}
}

func printCICheckResult(name string, checkResult CICheckResult) {
	if !ciJSON {
		cli.Print("\r  %s %s %s\n", getStatusIcon(checkResult.Status), name, dimStyle.Render(checkResult.Details))
	}
}

func appendCIArtifact(ctx context.Context, cwd string, check ciCheckDefinition, artifacts []string) []string {
	if !shouldGenerateSARIF(check) {
		return artifacts
	}

	sarifFile := core.PathJoin(cwd, check.name+".sarif")
	if generateSARIF(ctx, cwd, check.name, sarifFile) == nil {
		return append(artifacts, sarifFile)
	}

	return artifacts
}

func shouldGenerateSARIF(check ciCheckDefinition) bool {
	return (ciSARIF || ciUploadSARIF) && check.sarif
}

func ciExitCode(passed bool) int {
	if passed {
		return 0
	}
	return 1
}

func outputCIResult(ctx context.Context, cwd string, result CIResult, artifacts []string) error { // Result boundary
	if ciJSON {
		return outputCIJSONResult(result)
	}
	if ciSummary {
		return outputCISummaryResult(result)
	}
	return outputCIDefault(ctx, cwd, result, artifacts)
}

func outputCIJSONResult(result CIResult) error { // Result boundary
	if err := outputCIJSON(result); err != nil {
		return err
	}
	if !result.Passed {
		return phpExit(result.ExitCode, phpFailure(ciPipelineFailedMessage))
	}
	return nil
}

func outputCISummaryResult(result CIResult) error { // Result boundary
	if err := outputCISummary(result); err != nil {
		return err
	}
	if !result.Passed {
		return phpFailure(ciPipelineFailedMessage)
	}
	return nil
}

func outputCIDefault(ctx context.Context, cwd string, result CIResult, artifacts []string) error { // Result boundary
	cli.Print(cliSingleLineFormat, repeat("─", 40))

	if result.Passed {
		cli.Print(cliLabelValueFormat, successStyle.Render("✓ CI PASSED"), dimStyle.Render(result.Duration))
	} else {
		cli.Print(cliLabelValueFormat, errorStyle.Render("✗ CI FAILED"), dimStyle.Render(result.Duration))
	}

	if len(artifacts) > 0 {
		cli.Print(cliSingleLineFormat, dimStyle.Render("Artifacts:"))
		for _, a := range artifacts {
			cli.Print("  → %s\n", core.PathBase(a))
		}
	}

	// Upload SARIF if requested
	uploadCIArtifacts(ctx, artifacts)

	if !result.Passed {
		return phpFailure(ciPipelineFailedMessage)
	}
	return nil
}

func uploadCIArtifacts(ctx context.Context, artifacts []string) {
	if !ciUploadSARIF || len(artifacts) == 0 {
		return
	}

	cli.Blank()
	for _, sarifFile := range artifacts {
		if err := uploadSARIFToGitHub(ctx, sarifFile); err != nil {
			cli.Print("  %s %s: %s\n", errorStyle.Render("✗"), core.PathBase(sarifFile), err)
		} else {
			cli.Print("  %s %s uploaded\n", successStyle.Render("✓"), core.PathBase(sarifFile))
		}
	}
}

// runCITest runs Pest/PHPUnit tests
func runCITest(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: "test", Status: "passed"}

	opts := TestOptions{
		Dir:    dir,
		Output: nil, // Suppress output
	}

	if err := RunTests(ctx, opts); err != nil {
		result.Status = "failed"
		result.Details = err.Error()
	} else {
		result.Details = "all tests passed"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// runCIStan runs PHPStan
func runCIStan(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: "stan", Status: "passed"}

	_, found := DetectAnalyser(dir)
	if !found {
		result.Status = "skipped"
		result.Details = "PHPStan not configured"
		return result, nil
	}

	opts := AnalyseOptions{
		Dir:    dir,
		Output: nil,
	}

	if err := Analyse(ctx, opts); err != nil {
		result.Status = "failed"
		result.Details = "errors found"
	} else {
		result.Details = "0 errors"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// runCIPsalm runs Psalm
func runCIPsalm(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: "psalm", Status: "passed"}

	_, found := DetectPsalm(dir)
	if !found {
		result.Status = "skipped"
		result.Details = "Psalm not configured"
		return result, nil
	}

	opts := PsalmOptions{
		Dir:    dir,
		Output: nil,
	}

	if err := RunPsalm(ctx, opts); err != nil {
		result.Status = "failed"
		result.Details = "errors found"
	} else {
		result.Details = "0 errors"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// runCIFmt checks code formatting
func runCIFmt(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: `fmt`, Status: "passed"}

	_, found := DetectFormatter(dir)
	if !found {
		result.Status = "skipped"
		result.Details = "no formatter configured"
		return result, nil
	}

	opts := FormatOptions{
		Dir:    dir,
		Fix:    false, // Check only
		Output: nil,
	}

	if err := Format(ctx, opts); err != nil {
		result.Status = "warning"
		result.Details = "formatting issues"
	} else {
		result.Details = "code style OK"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// runCIAudit runs composer audit
func runCIAudit(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: "audit", Status: "passed"}

	results, err := RunAudit(ctx, AuditOptions{
		Dir:    dir,
		Output: nil,
	})
	if err != nil {
		result.Status = "failed"
		result.Details = err.Error()
		result.Duration = time.Since(start).Round(time.Millisecond).String()
		return result, nil
	}

	totalVulns := 0
	for _, r := range results {
		totalVulns += r.Vulnerabilities
	}

	if totalVulns > 0 {
		result.Status = "failed"
		result.Details = core.Sprintf("%d vulnerabilities", totalVulns)
		result.Issues = totalVulns
	} else {
		result.Details = "no vulnerabilities"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// runCISecurity runs security checks
func runCISecurity(ctx context.Context, dir string) (CICheckResult, error) { // Result boundary
	start := time.Now()
	result := CICheckResult{Name: "security", Status: "passed"}

	secResult, err := RunSecurityChecks(ctx, SecurityOptions{
		Dir:    dir,
		Output: nil,
	})
	if err != nil {
		result.Status = "failed"
		result.Details = err.Error()
		result.Duration = time.Since(start).Round(time.Millisecond).String()
		return result, nil
	}

	if secResult.Summary.Critical > 0 || secResult.Summary.High > 0 {
		result.Status = "failed"
		result.Details = core.Sprintf("%d critical, %d high", secResult.Summary.Critical, secResult.Summary.High)
		result.Issues = secResult.Summary.Critical + secResult.Summary.High
	} else if secResult.Summary.Medium > 0 {
		result.Status = "warning"
		result.Details = core.Sprintf("%d medium issues", secResult.Summary.Medium)
		result.Warnings = secResult.Summary.Medium
	} else {
		result.Details = "no issues"
	}

	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

// shouldFailOn determines if a check should cause CI failure based on --fail-on
func shouldFailOn(check CICheckResult, level string) bool {
	switch level {
	case "critical":
		return check.Status == "failed" && check.Issues > 0
	case "high", "error":
		return check.Status == "failed"
	case "warning":
		return check.Status == "failed" || check.Status == "warning"
	default:
		return check.Status == "failed"
	}
}

// getStatusIcon returns the icon for a check status
func getStatusIcon(status string) string {
	switch status {
	case "passed":
		return successStyle.Render("✓")
	case "failed":
		return errorStyle.Render("✗")
	case "warning":
		return phpQAWarningStyle.Render("⚠")
	case "skipped":
		return dimStyle.Render("-")
	default:
		return dimStyle.Render("?")
	}
}

// outputCIJSON outputs the result as JSON
func outputCIJSON(result CIResult) error { // Result boundary
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	core.Println(string(data))
	return nil
}

// outputCISummary outputs a markdown summary
func outputCISummary(result CIResult) error { // Result boundary
	var sb stringBuilder

	sb.WriteString("## CI Pipeline Results\n\n")

	if result.Passed {
		sb.WriteString("**Status:** ✅ Passed\n\n")
	} else {
		sb.WriteString("**Status:** ❌ Failed\n\n")
	}

	sb.WriteString("| Check | Status | Details |\n")
	sb.WriteString("|-------|--------|----------|\n")

	for _, check := range result.Checks {
		icon := "✅"
		switch check.Status {
		case "failed":
			icon = "❌"
		case "warning":
			icon = "⚠️"
		case "skipped":
			icon = "⏭️"
		}
		sb.WriteString(core.Sprintf("| %s | %s | %s |\n", check.Name, icon, check.Details))
	}

	sb.WriteString(core.Sprintf("\n**Duration:** %s\n", result.Duration))

	core.Println(sb.String())
	return nil
}

// generateSARIF generates a SARIF file for a specific check
func generateSARIF(ctx context.Context, dir, checkName, outputFile string) error { // Result boundary
	var args []string

	switch checkName {
	case "stan":
		args = []string{"vendor/bin/phpstan", "analyse", "--error-format=sarif", "--no-progress"}
	case "psalm":
		args = []string{"vendor/bin/psalm", "--output-format=sarif"}
	default:
		return core.E("php.SARIF", core.Sprintf("SARIF not supported for %s", checkName), nil)
	}

	cmd := exec.CommandContext(ctx, "php", args...)
	cmd.Dir = dir

	// Capture output - command may exit non-zero when issues are found
	// but still produce valid SARIF output
	output, err := cmd.CombinedOutput()
	if len(output) == 0 {
		if err != nil {
			return core.E("php.SARIF", "failed to generate SARIF", err)
		}
		return core.E("php.SARIF", "no SARIF output generated", nil)
	}

	// Validate output is valid JSON
	var js json.RawMessage
	if err := json.Unmarshal(output, &js); err != nil {
		return core.E("php.SARIF", "invalid SARIF output", err)
	}

	return getMedium().Write(outputFile, string(output))
}

// uploadSARIFToGitHub uploads a SARIF file to GitHub Security tab
func uploadSARIFToGitHub(ctx context.Context, sarifFile string) error { // Result boundary
	// Validate commit SHA before calling API
	sha := getGitSHA()
	if sha == "" {
		return core.E("php.SARIF", "cannot upload SARIF: git commit SHA not available (ensure you're in a git repository)", nil)
	}

	// Use gh CLI to upload
	cmd := exec.CommandContext(ctx, "gh", "api",
		"repos/{owner}/{repo}/code-scanning/sarifs",
		"-X", "POST",
		"-F", "sarif=@"+sarifFile,
		"-F", "ref="+getGitRef(),
		"-F", "commit_sha="+sha,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return core.E("php.SARIF", core.Sprintf("%s: %s", err, string(output)), nil)
	}
	return nil
}

// getGitRef returns the current git ref
func getGitRef() string {
	cmd := exec.Command("git", "symbolic-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "refs/heads/main"
	}
	return core.Trim(string(output))
}

// getGitSHA returns the current git commit SHA
func getGitSHA() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return core.Trim(string(output))
}
