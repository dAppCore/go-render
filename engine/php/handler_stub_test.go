package php

func TestHandlerStub_NewHandler_Good(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandlerStub_NewHandler_Bad(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandlerStub_NewHandler_Ugly(t *T) {
	subject := NewHandler
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandlerStub_Handler_LaravelRoot_Good(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandlerStub_Handler_LaravelRoot_Bad(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandlerStub_Handler_LaravelRoot_Ugly(t *T) {
	subject := (*Handler).LaravelRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandlerStub_Handler_DocRoot_Good(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandlerStub_Handler_DocRoot_Bad(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandlerStub_Handler_DocRoot_Ugly(t *T) {
	subject := (*Handler).DocRoot
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestHandlerStub_Handler_ServeHTTP_Good(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestHandlerStub_Handler_ServeHTTP_Bad(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestHandlerStub_Handler_ServeHTTP_Ugly(t *T) {
	subject := (*Handler).ServeHTTP
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
