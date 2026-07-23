package php

func TestGoCliHelpers_CommandLine_Bool_Good(t *T) {
	subject := phpCommandLine.Bool
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestGoCliHelpers_CommandLine_Bool_Bad(t *T) {
	subject := phpCommandLine.Bool
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestGoCliHelpers_CommandLine_Bool_Ugly(t *T) {
	subject := phpCommandLine.Bool
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestGoCliHelpers_CommandLine_String_Good(t *T) {
	subject := phpCommandLine.String
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestGoCliHelpers_CommandLine_String_Bad(t *T) {
	subject := phpCommandLine.String
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestGoCliHelpers_CommandLine_String_Ugly(t *T) {
	subject := phpCommandLine.String
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestGoCliHelpers_CommandLine_Int_Good(t *T) {
	subject := phpCommandLine.Int
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestGoCliHelpers_CommandLine_Int_Bad(t *T) {
	subject := phpCommandLine.Int
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestGoCliHelpers_CommandLine_Int_Ugly(t *T) {
	subject := phpCommandLine.Int
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestGoCliHelpers_CommandLine_Args_Good(t *T) {
	subject := phpCommandLine.Args
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestGoCliHelpers_CommandLine_Args_Bad(t *T) {
	subject := phpCommandLine.Args
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestGoCliHelpers_CommandLine_Args_Ugly(t *T) {
	subject := phpCommandLine.Args
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
