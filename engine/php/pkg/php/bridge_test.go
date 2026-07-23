package php

func TestBridge_NewBridge_Good(t *T) {
	subject := NewBridge
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestBridge_NewBridge_Bad(t *T) {
	subject := NewBridge
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestBridge_NewBridge_Ugly(t *T) {
	subject := NewBridge
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestBridge_Bridge_Port_Good(t *T) {
	subject := (*Bridge).Port
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestBridge_Bridge_Port_Bad(t *T) {
	subject := (*Bridge).Port
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestBridge_Bridge_Port_Ugly(t *T) {
	subject := (*Bridge).Port
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestBridge_Bridge_URL_Good(t *T) {
	subject := (*Bridge).URL
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestBridge_Bridge_URL_Bad(t *T) {
	subject := (*Bridge).URL
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestBridge_Bridge_URL_Ugly(t *T) {
	subject := (*Bridge).URL
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestBridge_Bridge_Shutdown_Good(t *T) {
	subject := (*Bridge).Shutdown
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestBridge_Bridge_Shutdown_Bad(t *T) {
	subject := (*Bridge).Shutdown
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestBridge_Bridge_Shutdown_Ugly(t *T) {
	subject := (*Bridge).Shutdown
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
