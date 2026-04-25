//go:build !js

// Package main provides a build-time CLI for generating Web Component bundles.
// Reads a JSON slot map from stdin, writes the generated JS or TypeScript to stdout.
//
// Usage:
//
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ -types > components.d.ts
//	go run ./cmd/codegen/ -watch --input=slots.json --output=components.js
package main

import (
	"context"
	goio "io"
	"syscall"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/html/codegen"
	coreio "dappco.re/go/io"
	log "dappco.re/go/log"
)

const defaultPollInterval = 250 * time.Millisecond

func generate(data []byte, emitTypes bool) (string, error) {
	var slots map[string]string
	if result := core.JSONUnmarshal(data, &slots); !result.OK {
		err, _ := result.Value.(error)
		return "", log.E("codegen", "invalid JSON", err)
	}

	if emitTypes {
		return codegen.GenerateTypeScriptDefinitions(slots), nil
	}

	out, err := codegen.GenerateBundle(slots)
	if err != nil {
		return "", log.E("codegen", "generate bundle", err)
	}
	return out, nil
}

func run(r goio.Reader, w goio.Writer, emitTypes bool) error {
	data, err := goio.ReadAll(r)
	if err != nil {
		return log.E("codegen", "reading stdin", err)
	}

	out, err := generate(data, emitTypes)
	if err != nil {
		return err
	}

	_, err = goio.WriteString(w, out)
	if err != nil {
		return log.E("codegen", "writing output", err)
	}
	return nil
}

func runDaemon(ctx context.Context, inputPath, outputPath string, emitTypes bool, pollInterval time.Duration) error {
	if inputPath == "" {
		return log.E("codegen", "watch mode requires -input", nil)
	}
	if outputPath == "" {
		return log.E("codegen", "watch mode requires -output", nil)
	}
	if pollInterval <= 0 {
		pollInterval = defaultPollInterval
	}

	var lastInput []byte
	for {
		input, err := readLocalFile(inputPath)
		if err != nil {
			return log.E("codegen", "reading input file", err)
		}

		if !sameBytes(input, lastInput) {
			out, err := generate(input, emitTypes)
			if err != nil {
				return err
			}
			if err := writeLocalFile(outputPath, out); err != nil {
				return log.E("codegen", "writing output file", err)
			}
			lastInput = append(lastInput[:0], input...)
		}

		select {
		case <-ctx.Done():
			if core.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

func readLocalFile(path string) ([]byte, error) {
	f, err := coreio.Local.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	return goio.ReadAll(f)
}

func writeLocalFile(path, content string) error {
	f, err := coreio.Local.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = goio.WriteString(f, content)
	return err
}

func sameBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range len(a) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func runStdio(emitTypes bool) error {
	stdin, err := coreio.Local.Open("/dev/stdin")
	if err != nil {
		return log.E("codegen.main", "open stdin", err)
	}
	defer func() {
		_ = stdin.Close()
	}()

	return run(stdin, stdoutWriter{}, emitTypes)
}

type stdoutWriter struct{}

func (stdoutWriter) Write(data []byte) (int, error) {
	total := 0
	for len(data) > 0 {
		n, err := syscall.Write(1, data)
		total += n
		data = data[n:]
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, goio.ErrShortWrite
		}
	}
	return total, nil
}

func newCodegenApp() *core.Core {
	c := core.New(core.WithOption("name", "codegen"))
	if cli := c.Cli(); cli != nil {
		cli.SetOutput(goio.Discard)
	}

	registerCodegenCommands(c)
	return c
}

func registerCodegenCommands(c *core.Core) {
	c.Command("generate", core.Command{
		Description: "Generate JavaScript or TypeScript from a JSON slot map on stdin",
		Flags:       codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runStdio(opts.Bool("types")))
		},
	})
	c.Command("types", core.Command{
		Description: "Generate TypeScript declarations from a JSON slot map on stdin",
		Action: func(core.Options) core.Result {
			return resultFromError(runStdio(true))
		},
	})
	c.Command("-types", core.Command{
		Hidden: true,
		Action: func(core.Options) core.Result {
			return resultFromError(runStdio(true))
		},
	})
	c.Command("--types", core.Command{
		Hidden: true,
		Action: func(core.Options) core.Result {
			return resultFromError(runStdio(true))
		},
	})
	c.Command("watch", core.Command{
		Description: "Poll an input JSON file and rewrite the generated output",
		Flags:       codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runWatchCommand(c, opts, opts.Bool("types")))
		},
	})
	c.Command("-watch", core.Command{
		Hidden: true,
		Flags:  codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runWatchCommand(c, opts, opts.Bool("types")))
		},
	})
	c.Command("--watch", core.Command{
		Hidden: true,
		Flags:  codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runWatchCommand(c, opts, opts.Bool("types")))
		},
	})
}

func codegenCommandFlags() core.Options {
	return core.NewOptions(
		core.Option{Key: "types", Value: false},
		core.Option{Key: "input", Value: ""},
		core.Option{Key: "output", Value: ""},
		core.Option{Key: "poll", Value: defaultPollInterval.String()},
	)
}

func runWatchCommand(c *core.Core, opts core.Options, emitTypes bool) error {
	pollInterval, err := pollIntervalFromOptions(opts)
	if err != nil {
		return err
	}

	ctx, cancel := codegenContext(c)
	defer cancel()

	return runDaemon(ctx, opts.String("input"), opts.String("output"), emitTypes, pollInterval)
}

func pollIntervalFromOptions(opts core.Options) (time.Duration, error) {
	raw := opts.String("poll")
	if raw == "" {
		return defaultPollInterval, nil
	}

	pollInterval, err := time.ParseDuration(raw)
	if err != nil {
		return 0, log.E("codegen", "invalid poll interval", err)
	}
	return pollInterval, nil
}

func codegenContext(c *core.Core) (context.Context, context.CancelFunc) {
	if c == nil {
		return context.WithCancel(context.Background())
	}

	ctx, cancel := context.WithCancel(c.Context())
	c.Action("signal.received", func(_ context.Context, opts core.Options) core.Result {
		switch opts.String("name") {
		case "SIGINT", "SIGTERM", "interrupt":
			cancel()
		}
		return core.Result{OK: true}
	})
	_ = c.Action("signal.start").Run(c.Context(), core.NewOptions(
		core.Option{Key: "signals", Value: []string{"SIGINT", "SIGTERM"}},
	))

	return ctx, cancel
}

func runCodegenApp(c *core.Core) error {
	if c == nil {
		return log.E("codegen.main", "core app is required", nil)
	}

	defer func() {
		_ = c.ServiceShutdown(context.Background())
	}()

	if result := c.ServiceStartup(c.Context(), nil); !result.OK {
		return resultError("codegen.main", "startup failed", result)
	}

	cli := c.Cli()
	if cli == nil {
		return runStdio(false)
	}

	result := cli.Run()
	if result.OK {
		return nil
	}
	if err, ok := result.Value.(error); ok && err != nil {
		return err
	}

	return runStdio(false)
}

func resultFromError(err error) core.Result {
	if err != nil {
		return core.Result{Value: err, OK: false}
	}
	return core.Result{OK: true}
}

func resultError(op, msg string, result core.Result) error {
	if result.OK {
		return nil
	}
	if err, ok := result.Value.(error); ok && err != nil {
		return err
	}
	return log.E(op, msg, nil)
}

func main() {
	c := newCodegenApp()
	if err := runCodegenApp(c); err != nil {
		log.Error("codegen failed", "scope", "codegen.main", "err", err)
		syscall.Exit(1)
	}
}
