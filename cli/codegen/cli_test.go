//go:build !js

package main

import (
	"testing"
	"time"

	core "dappco.re/go"
)

// TestInputPathFromOptions_Good — explicit -input wins over the stdin default.
func TestInputPathFromOptions_Good(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "input", Value: "/tmp/slots.json"})
	if got := inputPathFromOptions(opts); got != "/tmp/slots.json" {
		t.Fatalf("want /tmp/slots.json, got %q", got)
	}
}

// TestInputPathFromOptions_Ugly — absent -input falls back to the stdin path.
func TestInputPathFromOptions_Ugly(t *testing.T) {
	opts := core.NewOptions()
	if got := inputPathFromOptions(opts); got != defaultInputPath {
		t.Fatalf("want %q, got %q", defaultInputPath, got)
	}
}

// TestOutputPathFromOptions_Good — explicit -output wins over the stdout default.
func TestOutputPathFromOptions_Good(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "output", Value: "/tmp/out.js"})
	if got := outputPathFromOptions(opts); got != "/tmp/out.js" {
		t.Fatalf("want /tmp/out.js, got %q", got)
	}
}

// TestOutputPathFromOptions_Ugly — absent -output falls back to the stdout path.
func TestOutputPathFromOptions_Ugly(t *testing.T) {
	opts := core.NewOptions()
	if got := outputPathFromOptions(opts); got != defaultOutputPath {
		t.Fatalf("want %q, got %q", defaultOutputPath, got)
	}
}

// TestCodegenCommandFlags_Good — the shared flag set carries the documented defaults.
func TestCodegenCommandFlags_Good(t *testing.T) {
	flags := codegenCommandFlags()
	if flags.Bool("types") {
		t.Fatal("expected types default false")
	}
	if got := flags.String("input"); got != "" {
		t.Fatalf("expected empty input default, got %q", got)
	}
	if got := flags.String("output"); got != "" {
		t.Fatalf("expected empty output default, got %q", got)
	}
	if got := flags.String("poll"); got != defaultPollInterval.String() {
		t.Fatalf("expected poll default %q, got %q", defaultPollInterval.String(), got)
	}
}

// TestPollIntervalFromOptions_Good — a valid duration string parses through.
func TestPollIntervalFromOptions_Good(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "poll", Value: "500ms"})
	result := pollIntervalFromOptions(opts)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	got, _ := result.Value.(time.Duration)
	if got != 500*time.Millisecond {
		t.Fatalf("want 500ms, got %v", got)
	}
}

// TestPollIntervalFromOptions_Ugly — empty poll falls back to the default interval.
func TestPollIntervalFromOptions_Ugly(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "poll", Value: ""})
	result := pollIntervalFromOptions(opts)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	got, _ := result.Value.(time.Duration)
	if got != defaultPollInterval {
		t.Fatalf("want %v, got %v", defaultPollInterval, got)
	}
}

// TestPollIntervalFromOptions_Bad — an unparseable duration is a failed result.
func TestPollIntervalFromOptions_Bad(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "poll", Value: "not-a-duration"})
	result := pollIntervalFromOptions(opts)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "invalid poll interval") {
		t.Fatalf("expected error to contain %q, got %v", "invalid poll interval", result.Error())
	}
}

// TestRunGenerateCommand_Good — generate writes a bundle to an explicit output file.
func TestRunGenerateCommand_Good(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.js")
	if err := writeTextFile(inputPath, `{"H":"nav-bar"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := core.NewOptions(
		core.Option{Key: "input", Value: inputPath},
		core.Option{Key: "output", Value: outputPath},
	)
	if result := runGenerateCommand(opts, false); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	got, err := readTextFile(outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !core.Contains(got, "NavBar") {
		t.Fatalf("expected bundle to contain NavBar, got %q", got)
	}
}

// TestRunGenerateCommand_Bad — invalid JSON in the input file surfaces an error.
func TestRunGenerateCommand_Bad(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.js")
	if err := writeTextFile(inputPath, `not json`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := core.NewOptions(
		core.Option{Key: "input", Value: inputPath},
		core.Option{Key: "output", Value: outputPath},
	)
	if result := runGenerateCommand(opts, false); result.OK {
		t.Fatal("expected error result, got OK")
	}
}

// TestRunGenerateCommand_TypesGood — types mode emits TypeScript declarations.
func TestRunGenerateCommand_TypesGood(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.d.ts")
	if err := writeTextFile(inputPath, `{"H":"nav-bar"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts := core.NewOptions(
		core.Option{Key: "input", Value: inputPath},
		core.Option{Key: "output", Value: outputPath},
	)
	if result := runGenerateCommand(opts, true); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	got, err := readTextFile(outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !core.Contains(got, "declare global") {
		t.Fatalf("expected declarations, got %q", got)
	}
}

// TestRunWatchCommand_Bad — watch without -input is rejected before polling.
func TestRunWatchCommand_Bad(t *testing.T) {
	opts := core.NewOptions(
		core.Option{Key: "input", Value: ""},
		core.Option{Key: "output", Value: ""},
		core.Option{Key: "poll", Value: "5ms"},
	)
	result := runWatchCommand(nil, opts, false)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "watch mode requires -input") {
		t.Fatalf("expected error to contain %q, got %v", "watch mode requires -input", result.Error())
	}
}

// TestRunWatchCommand_BadPoll — an unparseable poll interval short-circuits.
func TestRunWatchCommand_BadPoll(t *testing.T) {
	opts := core.NewOptions(core.Option{Key: "poll", Value: "nope"})
	result := runWatchCommand(nil, opts, false)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "invalid poll interval") {
		t.Fatalf("expected error to contain %q, got %v", "invalid poll interval", result.Error())
	}
}

// TestNewCodegenApp_Good — the app registers every documented command path.
func TestNewCodegenApp_Good(t *testing.T) {
	c := newCodegenApp()
	if c == nil {
		t.Fatal("expected non-nil core app")
	}
	if c.Cli() == nil {
		t.Fatal("expected a CLI to be present")
	}
}

// TestRegisterCodegenCommands_Good — registration is idempotent on a bare core.
func TestRegisterCodegenCommands_Good(t *testing.T) {
	// core.WithCli() opts into the CLI service (default since dappco.re/go
	// v0.12.0 removed the previously-unconditional CLI registration).
	c := core.New(core.WithOption("name", "codegen"), core.WithCli())
	registerCodegenCommands(c)
	if c.Cli() == nil {
		t.Fatal("expected a CLI after registration")
	}
}

// TestGenerateCommandAction_Good — the registered "generate" action runs end-to-end.
func TestGenerateCommandAction_Good(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.js")
	if err := writeTextFile(inputPath, `{"H":"nav-bar"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := newCodegenApp()
	cmd, ok := c.Command("generate").Value.(*core.Command)
	if !ok || cmd.Action == nil {
		t.Fatal("expected a generate command with an action")
	}

	opts := core.NewOptions(
		core.Option{Key: "input", Value: inputPath},
		core.Option{Key: "output", Value: outputPath},
	)
	if result := cmd.Action(opts); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	got, err := readTextFile(outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !core.Contains(got, "NavBar") {
		t.Fatalf("expected bundle to contain NavBar, got %q", got)
	}
}

// TestTypesCommandAction_Good — the registered "types" action forces declarations.
func TestTypesCommandAction_Good(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.d.ts")
	if err := writeTextFile(inputPath, `{"H":"nav-bar"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := newCodegenApp()
	cmd, ok := c.Command("types").Value.(*core.Command)
	if !ok || cmd.Action == nil {
		t.Fatal("expected a types command with an action")
	}

	opts := core.NewOptions(
		core.Option{Key: "input", Value: inputPath},
		core.Option{Key: "output", Value: outputPath},
	)
	if result := cmd.Action(opts); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	got, err := readTextFile(outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !core.Contains(got, "declare global") {
		t.Fatalf("expected declarations, got %q", got)
	}
}

// TestWatchCommandAction_Bad — the registered "watch" action rejects missing -input.
func TestWatchCommandAction_Bad(t *testing.T) {
	c := newCodegenApp()
	cmd, ok := c.Command("watch").Value.(*core.Command)
	if !ok || cmd.Action == nil {
		t.Fatal("expected a watch command with an action")
	}
	opts := core.NewOptions(
		core.Option{Key: "input", Value: ""},
		core.Option{Key: "output", Value: ""},
	)
	if result := cmd.Action(opts); result.OK {
		t.Fatal("expected error result, got OK")
	}
}

// TestReadLocalFile_Bad — a missing named file is a failed result.
func TestReadLocalFile_Bad(t *testing.T) {
	dir := t.TempDir()
	missing := core.Path(dir, "does-not-exist.json")
	result := readLocalFile(missing)
	if result.OK {
		t.Fatal("expected error result for missing file, got OK")
	}
}

// TestReadLocalFile_Good — a named file reads back its bytes.
func TestReadLocalFile_Good(t *testing.T) {
	dir := t.TempDir()
	path := core.Path(dir, "slots.json")
	if err := writeTextFile(path, `{"H":"nav-bar"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := readLocalFile(path)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	data, _ := result.Value.([]byte)
	if !core.Contains(string(data), "nav-bar") {
		t.Fatalf("expected file contents, got %q", string(data))
	}
}

// TestWriteLocalFile_Good — a named file is written with the content.
func TestWriteLocalFile_Good(t *testing.T) {
	dir := t.TempDir()
	path := core.Path(dir, "out.js")
	if result := writeLocalFile(path, "hello"); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	got, err := readTextFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("want hello, got %q", got)
	}
}

// TestWriteLocalFile_Bad — writing under a path whose parent is a file fails.
func TestWriteLocalFile_Bad(t *testing.T) {
	dir := t.TempDir()
	parent := core.Path(dir, "afile")
	if err := writeTextFile(parent, "x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// parent is a regular file, so treating it as a directory must fail.
	nested := core.Path(parent, "child.js")
	if result := writeLocalFile(nested, "data"); result.OK {
		t.Fatal("expected error result writing under a file, got OK")
	}
}

// TestResultFromError_Good — a nil error becomes an OK result.
func TestResultFromError_Good(t *testing.T) {
	if result := resultFromError(nil); !result.OK {
		t.Fatalf("expected OK, got %v", result.Error())
	}
}

// TestResultFromError_Bad — a non-nil error becomes a failed result.
func TestResultFromError_Bad(t *testing.T) {
	if result := resultFromError(core.E("codegen.test", "boom", nil)); result.OK {
		t.Fatal("expected failed result, got OK")
	}
}

// TestResultError_Good — an OK input result stays OK regardless of op/msg.
func TestResultError_Good(t *testing.T) {
	if result := resultError("codegen.test", "wrap", core.Ok(nil)); !result.OK {
		t.Fatalf("expected OK passthrough, got %v", result.Error())
	}
}

// TestResultError_Bad — a failed input result is wrapped with op + msg.
func TestResultError_Bad(t *testing.T) {
	in := core.Fail(core.E("inner", "cause", nil))
	result := resultError("codegen.test", "wrap", in)
	if result.OK {
		t.Fatal("expected failed result, got OK")
	}
	if !core.Contains(result.Error(), "wrap") {
		t.Fatalf("expected wrapped message, got %v", result.Error())
	}
}

// TestResultError_Ugly — a failed result with a non-error Value still wraps cleanly.
func TestResultError_Ugly(t *testing.T) {
	in := core.Result{Value: "not-an-error", OK: false}
	result := resultError("codegen.test", "wrap", in)
	if result.OK {
		t.Fatal("expected failed result, got OK")
	}
	if !core.Contains(result.Error(), "wrap") {
		t.Fatalf("expected wrapped message, got %v", result.Error())
	}
}
