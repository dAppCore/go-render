//go:build !js

// Package main provides a build-time CLI for generating Web Component bundles.
// Reads a JSON slot map from stdin, writes the generated JS or TypeScript to stdout.
//
// Usage:
//
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ > components.js
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ -types > components.d.ts
//	echo '{"H":"nav-bar","C":"main-content"}' | go run ./cmd/codegen/ -watch -input slots.json -output components.js
package main

import (
	"context"
	"errors"
	"flag"
	goio "io"
	"os"
	"os/signal"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/core/html/codegen"
	coreio "dappco.re/go/core/io"
	log "dappco.re/go/core/log"
)

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
		pollInterval = 250 * time.Millisecond
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
			if errors.Is(ctx.Err(), context.Canceled) {
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

func main() {
	emitWatch := flag.Bool("watch", false, "poll input and rewrite output when the JSON changes")
	inputPath := flag.String("input", "", "path to the JSON slot map used by -watch")
	outputPath := flag.String("output", "", "path to the generated bundle written by -watch")
	emitTypes := flag.Bool("types", false, "emit TypeScript declarations instead of JavaScript")
	pollInterval := flag.Duration("poll", 250*time.Millisecond, "poll interval used by -watch")
	flag.Parse()

	if *emitWatch {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if err := runDaemon(ctx, *inputPath, *outputPath, *emitTypes, *pollInterval); err != nil {
			log.Error("codegen failed", "scope", "codegen.main", "err", err)
			os.Exit(1)
		}
		return
	}

	stdin, err := coreio.Local.Open("/dev/stdin")
	if err != nil {
		log.Error("failed to open stdin", "scope", "codegen.main", "err", log.E("codegen.main", "open stdin", err))
		os.Exit(1)
	}

	stdout, err := coreio.Local.Create("/dev/stdout")
	if err != nil {
		_ = stdin.Close()
		log.Error("failed to open stdout", "scope", "codegen.main", "err", log.E("codegen.main", "open stdout", err))
		os.Exit(1)
	}
	defer func() {
		_ = stdin.Close()
		_ = stdout.Close()
	}()

	if err := run(stdin, stdout, *emitTypes); err != nil {
		log.Error("codegen failed", "scope", "codegen.main", "err", err)
		os.Exit(1)
	}
}
