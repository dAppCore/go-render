//go:build !js

package main

import (
	"context"
	. "dappco.re/go"
	goio "io"
	"testing"
	"time"

	core "dappco.re/go"
	coreio "dappco.re/go/io"
)

func TestRun_WritesBundleGood(t *testing.T) {
	input := core.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	output := core.NewBuilder()

	if err := run(input, output, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
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

	err := run(input, output, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !core.Contains(err.Error(), "invalid JSON") {
		t.Fatalf("expected error to contain %q, got %v", "invalid JSON", err)
	}
}

func TestRun_InvalidTagBad(t *testing.T) {
	input := core.NewReader(`{"H":"notag"}`)
	output := core.NewBuilder()

	err := run(input, output, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !core.Contains(err.Error(), "hyphen") {
		t.Fatalf("expected error to contain %q, got %v", "hyphen", err)
	}
}

func TestRun_InvalidTagCharactersBad(t *testing.T) {
	input := core.NewReader(`{"H":"Nav-Bar","C":"nav bar"}`)
	output := core.NewBuilder()

	err := run(input, output, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !core.Contains(err.Error(), "lowercase hyphenated name") {
		t.Fatalf("expected error to contain %q, got %v", "lowercase hyphenated name", err)
	}
}

func TestRun_EmptySlotsGood(t *testing.T) {
	input := core.NewReader(`{}`)
	output := core.NewBuilder()

	if err := run(input, output, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := output.String(); got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestRun_WritesTypeScriptDefinitionsGood(t *testing.T) {
	input := core.NewReader(`{"H":"nav-bar","C":"main-content"}`)
	output := core.NewBuilder()

	if err := run(input, output, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
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

	done := make(chan error, 1)
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
	if err := <-done; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunDaemon_MissingPathsBad(t *testing.T) {
	err := runDaemon(context.Background(), "", "", false, time.Millisecond)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !core.Contains(err.Error(), "watch mode requires -input") {
		t.Fatalf("expected error to contain %q, got %v", "watch mode requires -input", err)
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

func TestMain_Writer_Write_Good(t *T) {
	writer := discardWriter{}
	n, err := writer.Write([]byte("agent"))
	AssertNoError(t, err)
	AssertEqual(t, 5, n)
}

func TestMain_Writer_Write_Bad(t *T) {
	writer := discardWriter{}
	n, err := writer.Write(nil)
	AssertNoError(t, err)
	AssertEqual(t, 0, n)
}

func TestMain_Writer_Write_Ugly(t *T) {
	writer := discardWriter{}
	n, err := writer.Write([]byte{0, 'x'})
	AssertNoError(t, err)
	AssertEqual(t, 2, n)
}
