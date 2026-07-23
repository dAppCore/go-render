//go:build !js

package main

import (
	"context"
	goio "io"
	"testing"
	"time"

	core "dappco.re/go"
	coreio "dappco.re/go/io"
)

func TestRun_WritesBundleGood(t *testing.T) {
	input := core.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	output := core.NewBuilder()

	if result := run(input, output, false); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	js := output.String()
	if !core.Contains(js, "NavBar") {
		t.Fatal("expected js to contain NavBar")
	}
	if !core.Contains(js, "MainContent") {
		t.Fatal("expected js to contain MainContent")
	}
	if !core.Contains(js, "customElements.define") {
		t.Fatal("expected js to contain customElements.define")
	}
	if got := countSubstr(js, "extends HTMLElement"); got != 2 {
		t.Fatalf("want 2 extends HTMLElement, got %d", got)
	}
}

func TestRun_InvalidJSONBad(t *testing.T) {
	input := core.NewReader(`not json`)
	output := core.NewBuilder()

	result := run(input, output, false)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "invalid JSON") {
		t.Fatalf("expected error to contain %q, got %v", "invalid JSON", result.Error())
	}
}

func TestRun_InvalidTagBad(t *testing.T) {
	input := core.NewReader(`{"H":"notag"}`)
	output := core.NewBuilder()

	result := run(input, output, false)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "hyphen") {
		t.Fatalf("expected error to contain %q, got %v", "hyphen", result.Error())
	}
}

func TestRun_InvalidTagCharactersBad(t *testing.T) {
	input := core.NewReader(`{"H":"Nav-Bar","C":"nav bar"}`)
	output := core.NewBuilder()

	result := run(input, output, false)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "lowercase hyphenated name") {
		t.Fatalf("expected error to contain %q, got %v", "lowercase hyphenated name", result.Error())
	}
}

func TestRun_EmptySlotsGood(t *testing.T) {
	input := core.NewReader(`{}`)
	output := core.NewBuilder()

	if result := run(input, output, false); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	if got := output.String(); got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestRun_WritesTypeScriptDefinitionsGood(t *testing.T) {
	input := core.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	output := core.NewBuilder()

	if result := run(input, output, true); !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}

	dts := output.String()
	for _, want := range []string{
		"declare global",
		`"nav-bar": NavBar;`,
		`"main-content": MainContent;`,
		"export declare class NavBar extends HTMLElement",
		"export declare class MainContent extends HTMLElement",
	} {
		if !core.Contains(dts, want) {
			t.Fatalf("expected dts to contain %q", want)
		}
	}
}

func TestRunDaemon_WritesUpdatedBundleGood(t *testing.T) {
	dir := t.TempDir()
	inputPath := core.Path(dir, "slots.json")
	outputPath := core.Path(dir, "bundle.js")

	if err := writeTextFile(inputPath, `{"H":"nav-bar","C":"main-content"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan core.Result, 1)
	go func() {
		done <- runDaemon(ctx, inputPath, outputPath, false, 5*time.Millisecond)
	}()

	deadline := time.Now().Add(time.Second)
	ok := false
	for time.Now().Before(deadline) {
		got, err := readTextFile(outputPath)
		if err == nil && core.Contains(got, "NavBar") && core.Contains(got, "MainContent") {
			ok = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !ok {
		t.Fatal("expected bundle file to contain NavBar and MainContent within 1s")
	}

	cancel()
	if result := <-done; !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
}

func TestRunDaemon_MissingPathsBad(t *testing.T) {
	result := runDaemon(context.Background(), "", "", false, time.Millisecond)
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
	if !core.Contains(result.Error(), "watch mode requires -input") {
		t.Fatalf("expected error to contain %q, got %v", "watch mode requires -input", result.Error())
	}
}

func countSubstr(s, substr string) int {
	if substr == "" {
		return len(s) + 1
	}

	count := 0
	for i := 0; i <= len(s)-len(substr); {
		j := indexSubstr(s[i:], substr)
		if j < 0 {
			return count
		}
		count++
		i += j + len(substr)
	}

	return count
}

func indexSubstr(s, substr string) int {
	if substr == "" {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

func writeTextFile(path, content string) error {
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

func readTextFile(path string) (string, error) {
	f, err := coreio.Local.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err := goio.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
