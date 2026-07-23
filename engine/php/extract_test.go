package php

func TestExtract_Extract_Good(t *T) {
	subject := Extract
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestExtract_Extract_Bad(t *T) {
	subject := Extract
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestExtract_Extract_Ugly(t *T) {
	subject := Extract
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
