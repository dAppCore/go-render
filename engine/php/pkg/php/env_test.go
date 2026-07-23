package php

func TestEnv_PrepareRuntimeEnvironment_Good(t *T) {
	subject := PrepareRuntimeEnvironment
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestEnv_PrepareRuntimeEnvironment_Bad(t *T) {
	subject := PrepareRuntimeEnvironment
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestEnv_PrepareRuntimeEnvironment_Ugly(t *T) {
	subject := PrepareRuntimeEnvironment
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestEnv_AppendEnv_Good(t *T) {
	subject := AppendEnv
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestEnv_AppendEnv_Bad(t *T) {
	subject := AppendEnv
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestEnv_AppendEnv_Ugly(t *T) {
	subject := AppendEnv
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
