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
	"time"

	core "dappco.re/go"
	"dappco.re/go/html/codegen"
	coreio "dappco.re/go/io"
	log "dappco.re/go/log"
)

const defaultPollInterval = 250 * time.Millisecond
const defaultInputPath = "/dev/stdin"
const defaultOutputPath = "/dev/stdout"

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

func run(input, output any, emitTypes bool) error {
	data, err := readInput(input)
	if err != nil {
		return log.E("codegen", "reading input", err)
	}

	out, err := generate(data, emitTypes)
	if err != nil {
		return err
	}

	if err := writeOutput(output, out); err != nil {
		return log.E("codegen", "writing output", err)
	}
	return nil
}

func readInput(input any) ([]byte, error) {
	if path, ok := input.(string); ok {
		return readLocalFile(path)
	}

	result := core.ReadAll(input)
	if !result.OK {
		return nil, resultError("codegen", "reading input stream", result)
	}
	content, _ := result.Value.(string)
	return []byte(content), nil
}

func writeOutput(output any, content string) error {
	if path, ok := output.(string); ok {
		return writeLocalFile(path, content)
	}

	result := core.WriteAll(output, content)
	if !result.OK {
		return resultError("codegen", "writing output stream", result)
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
	if path == "" {
		path = defaultInputPath
	}

	if path == defaultInputPath {
		f, err := coreio.Local.Open(path)
		if err != nil {
			return nil, err
		}

		result := core.ReadAll(f)
		if !result.OK {
			return nil, resultError("codegen", "reading stdin", result)
		}
		content, _ := result.Value.(string)
		return []byte(content), nil
	}

	content, err := coreio.Local.Read(path)
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func writeLocalFile(path, content string) error {
	if path == "" {
		path = defaultOutputPath
	}

	if path == defaultOutputPath {
		f, err := coreio.Local.Create(path)
		if err != nil {
			f, err = coreio.Local.Append(path)
			if err != nil {
				core.Print(nil, "%s", content)
				return nil
			}
		}

		result := core.WriteAll(f, content)
		if !result.OK {
			return resultError("codegen", "writing stdout", result)
		}
		return nil
	}

	return coreio.Local.Write(path, content)
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
	return run(defaultInputPath, defaultOutputPath, emitTypes)
}

func newCodegenApp() *core.Core {
	c := core.New(core.WithOption("name", "codegen"))
	if cli := c.Cli(); cli != nil {
		cli.SetOutput(discardWriter{})
	}

	registerCodegenCommands(c)
	return c
}

func registerCodegenCommands(c *core.Core) {
	c.Command("generate", core.Command{
		Description: "Generate JavaScript or TypeScript from a JSON slot map on stdin",
		Flags:       codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runGenerateCommand(opts, opts.Bool("types")))
		},
	})
	c.Command("types", core.Command{
		Description: "Generate TypeScript declarations from a JSON slot map on stdin",
		Flags:       codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runGenerateCommand(opts, true))
		},
	})
	c.Command("-types", core.Command{
		Hidden: true,
		Flags:  codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runGenerateCommand(opts, true))
		},
	})
	c.Command("--types", core.Command{
		Hidden: true,
		Flags:  codegenCommandFlags(),
		Action: func(opts core.Options) core.Result {
			return resultFromError(runGenerateCommand(opts, true))
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

func runGenerateCommand(opts core.Options, emitTypes bool) error {
	return run(inputPathFromOptions(opts), outputPathFromOptions(opts), emitTypes)
}

func inputPathFromOptions(opts core.Options) string {
	if input := opts.String("input"); input != "" {
		return input
	}
	return defaultInputPath
}

func outputPathFromOptions(opts core.Options) string {
	if output := opts.String("output"); output != "" {
		return output
	}
	return defaultOutputPath
}

func runWatchCommand(c *core.Core, opts core.Options, emitTypes bool) error {
	pollInterval, err := pollIntervalFromOptions(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if c != nil {
		ctx = c.Context()
	}
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

func runCodegenApp(c *core.Core) error {
	if c == nil {
		return log.E("codegen.main", "core app is required", nil)
	}

	defer func() {
		if result := c.ServiceShutdown(context.Background()); !result.OK {
			log.Warn("codegen shutdown failed", "scope", "codegen.main", "err", result.Error())
		}
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
		return core.Fail(err)
	}
	return core.Ok(nil)
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

type discardWriter struct{}

func (discardWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func main() {
	c := newCodegenApp()
	if err := runCodegenApp(c); err != nil {
		log.Error("codegen failed", "scope", "codegen.main", "err", err)
	}
}
