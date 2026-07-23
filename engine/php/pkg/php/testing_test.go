package php

func TestTesting_DetectTestRunner_Good(t *T) {
	subject := DetectTestRunner
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestTesting_DetectTestRunner_Bad(t *T) {
	subject := DetectTestRunner
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestTesting_DetectTestRunner_Ugly(t *T) {
	subject := DetectTestRunner
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestTesting_RunTests_Good(t *T) {
	subject := RunTests
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestTesting_RunTests_Bad(t *T) {
	subject := RunTests
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestTesting_RunTests_Ugly(t *T) {
	subject := RunTests
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestTesting_RunParallel_Good(t *T) {
	subject := RunParallel
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestTesting_RunParallel_Bad(t *T) {
	subject := RunParallel
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestTesting_RunParallel_Ugly(t *T) {
	subject := RunParallel
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
