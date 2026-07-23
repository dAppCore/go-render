package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
)

func newClipboardToolsTestSubsystem(t *core.T, query func(core.Query) core.Result) *Subsystem {
	t.Helper()
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if query != nil {
			return query(q)
		}
		return core.Result{}
	})
	return New(c)
}

func TestToolsClipboard_clipboardRead_Good(t *core.T) {
	// clipboardRead
	ax7Variant := "clipboardRead:good"
	core.AssertContains(t, ax7Variant, "good")
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{
				Value: clipboard.ClipboardContent{
					Text:       "hello",
					HasContent: true,
				},
				OK: true,
			}
		}
		return core.Result{}
	})

	_, out, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	core.RequireNoError(t, err)
	core.AssertEqual(t, "hello", out.Content)
}

func TestToolsClipboard_clipboardRead_Bad(t *core.T) {
	// clipboardRead
	ax7Variant := "clipboardRead:bad"
	core.AssertContains(t, ax7Variant, "bad")
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{OK: false, Value: "clipboard backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "clipboard query failed")
}

func TestToolsClipboard_clipboardRead_Ugly(t *core.T) {
	// clipboardRead
	ax7Variant := "clipboardRead:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{OK: true, Value: core.NewError("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "unexpected result type")
}
