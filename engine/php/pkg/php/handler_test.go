package php

func TestHandler_NewHandler_Good(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandler_NewHandler_Bad(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandler_NewHandler_Ugly(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandler_Handler_LaravelRoot_Good(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandler_Handler_LaravelRoot_Bad(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandler_Handler_LaravelRoot_Ugly(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandler_Handler_DocRoot_Good(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandler_Handler_DocRoot_Bad(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandler_Handler_DocRoot_Ugly(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandler_Handler_ServeHTTP_Good(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandler_Handler_ServeHTTP_Bad(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandler_Handler_ServeHTTP_Ugly(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
