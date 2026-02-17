package html

import "testing"

func TestStripTags_Simple(t *testing.T) {
	got := StripTags(`<div>hello</div>`)
	want := "hello"
	if got != want {
		t.Errorf("StripTags(<div>hello</div>) = %q, want %q", got, want)
	}
}

func TestStripTags_Nested(t *testing.T) {
	got := StripTags(`<header role="banner"><h1>Title</h1></header>`)
	want := "Title"
	if got != want {
		t.Errorf("StripTags(nested) = %q, want %q", got, want)
	}
}

func TestStripTags_MultipleRegions(t *testing.T) {
	got := StripTags(`<header>Head</header><main>Body</main><footer>Foot</footer>`)
	want := "Head Body Foot"
	if got != want {
		t.Errorf("StripTags(multi) = %q, want %q", got, want)
	}
}

func TestStripTags_Empty(t *testing.T) {
	got := StripTags("")
	if got != "" {
		t.Errorf("StripTags(\"\") = %q, want empty", got)
	}
}

func TestStripTags_NoTags(t *testing.T) {
	got := StripTags("plain text")
	if got != "plain text" {
		t.Errorf("StripTags(plain) = %q, want %q", got, "plain text")
	}
}

func TestStripTags_Entities(t *testing.T) {
	got := StripTags(`&lt;script&gt;`)
	want := "&lt;script&gt;"
	if got != want {
		t.Errorf("StripTags should preserve entities, got %q, want %q", got, want)
	}
}
