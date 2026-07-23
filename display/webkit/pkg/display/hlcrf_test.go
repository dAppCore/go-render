package display

import (
	core "dappco.re/go"
)

func TestHLCRF_DefaultHLCRFTag_GoodCase(t *core.T) {
	core.AssertEqual(t, "core-widget", defaultHLCRFTag("Widget.ts"))
	observedType := core.Sprintf("%T", defaultHLCRFTag("Widget.ts"))
	core.AssertNotEmpty(t, observedType)
}

func TestHLCRF_DefaultHLCRFTag_BadCase(t *core.T) {
	core.AssertEqual(t, "feature-card", defaultHLCRFTag("feature_card.html"))
	observedType := core.Sprintf("%T", defaultHLCRFTag("feature_card.html"))
	core.AssertNotEmpty(t, observedType)
}

func TestHLCRF_DefaultHLCRFTag_UglyCase(t *core.T) {
	core.AssertEqual(t, "core-", defaultHLCRFTag(""))
	observedType := core.Sprintf("%T", defaultHLCRFTag(""))
	core.AssertNotEmpty(t, observedType)
}

func TestHLCRF_BuildHLCRFComponents_GoodCase(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(root, ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "card.html"), "<article>Card</article>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, ".core", "view.yaml"), core.Join("\n", []string{
		"hlcrf:",
		"  - name: card.html",
		"  - tag: core-inline",
		"    template: <section>Inline</section>",
	}...), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(core.PathJoin(root, "index.html"))

	core.RequireNoError(t, err)
	core.RequireNotEmpty(t, script)
	core.AssertContains(t, script, "customElements.define")
	core.AssertContains(t, script, "article>Card</article>")
	core.AssertContains(t, script, "<section>Inline</section>")
	core.AssertContains(t, script, "core-card")
	core.AssertContains(t, script, "core-inline")
}

func TestHLCRF_CompileHLCRFTemplate_GoodCase(t *core.T) {
	compiled := compileHLCRFTemplate(`<section data-slot="H">{{slot "H"}}</section><main>{{ slot "L-C" }}</main><footer>{{ slot "" }}{{ slot "default" }}</footer>`)

	core.AssertContains(t, compiled, `<slot name="H"></slot>`)
	core.AssertContains(t, compiled, `<slot name="L-C"></slot>`)
	core.AssertContains(t, compiled, `<footer><slot></slot><slot></slot></footer>`)
}

func TestHLCRF_BuildHLCRFComponents_BadCase(t *core.T) {
	svc := &Service{}

	script, err := svc.buildHLCRFComponents(core.PathJoin(t.TempDir(), "missing.html"))

	core.RequireNoError(t, err)
	core.AssertEmpty(t, script)
}

func TestHLCRF_BuildHLCRFComponents_UglyCase(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(root, ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, ".core", "view.yaml"), core.Join("\n", []string{
		"hlcrf:",
		"  - name: missing.html",
		"  - template: <span>Fallback</span>",
	}...), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(core.PathJoin(root, "index.html"))

	core.RequireNoError(t, err)
	core.RequireNotEmpty(t, script)
	core.AssertContains(t, script, "<span>Fallback</span>")
	core.AssertNotContains(t, script, "missing.html")
}

func TestHLCRF_BuildHLCRFComponents_RejectsTraversal(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreEnsureDir(core.PathJoin(root, ".core")))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, "outside.html"), "<span>Outside</span>", 0o644))
	core.RequireNoError(t, coreWriteMode(core.PathJoin(root, ".core", "view.yaml"), core.Join("\n", []string{
		"hlcrf:",
		"  - name: ../outside.html",
	}...), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(core.PathJoin(root, "index.html"))

	core.AssertError(t, err)
	core.AssertEmpty(t, script)
}
