package php

func TestCmdServeFrankenphp_ResponseWriter_Header_Good(t *T) {
	subject := (*execResponseWriter).Header
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmdServeFrankenphp_ResponseWriter_Header_Bad(t *T) {
	subject := (*execResponseWriter).Header
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmdServeFrankenphp_ResponseWriter_Header_Ugly(t *T) {
	subject := (*execResponseWriter).Header
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestCmdServeFrankenphp_ResponseWriter_Write_Good(t *T) {
	subject := (*execResponseWriter).Write
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmdServeFrankenphp_ResponseWriter_Write_Bad(t *T) {
	subject := (*execResponseWriter).Write
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmdServeFrankenphp_ResponseWriter_Write_Ugly(t *T) {
	subject := (*execResponseWriter).Write
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestCmdServeFrankenphp_ResponseWriter_WriteHeader_Good(t *T) {
	subject := (*execResponseWriter).WriteHeader
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmdServeFrankenphp_ResponseWriter_WriteHeader_Bad(t *T) {
	subject := (*execResponseWriter).WriteHeader
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmdServeFrankenphp_ResponseWriter_WriteHeader_Ugly(t *T) {
	subject := (*execResponseWriter).WriteHeader
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
