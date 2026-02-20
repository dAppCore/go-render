package codegen

import "testing"

func BenchmarkGenerateClass(b *testing.B) {
	for b.Loop() {
		GenerateClass("photo-grid", "C")
	}
}

func BenchmarkTagToClassName(b *testing.B) {
	for b.Loop() {
		TagToClassName("my-super-widget-component")
	}
}

func BenchmarkGenerateBundle_Small(b *testing.B) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "main-content",
	}
	b.ResetTimer()
	for b.Loop() {
		GenerateBundle(slots)
	}
}

func BenchmarkGenerateBundle_Full(b *testing.B) {
	slots := map[string]string{
		"H": "nav-bar",
		"L": "side-panel",
		"C": "main-content",
		"R": "aside-widget",
		"F": "page-footer",
	}
	b.ResetTimer()
	for b.Loop() {
		GenerateBundle(slots)
	}
}

func BenchmarkGenerateRegistration(b *testing.B) {
	for b.Loop() {
		GenerateRegistration("photo-grid", "PhotoGrid")
	}
}
