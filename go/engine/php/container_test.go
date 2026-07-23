package php

func TestContainer_BuildDocker_Good(t *T) {
	subject := BuildDocker
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestContainer_BuildDocker_Bad(t *T) {
	subject := BuildDocker
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestContainer_BuildDocker_Ugly(t *T) {
	subject := BuildDocker
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestContainer_BuildLinuxKit_Good(t *T) {
	subject := BuildLinuxKit
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestContainer_BuildLinuxKit_Bad(t *T) {
	subject := BuildLinuxKit
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestContainer_BuildLinuxKit_Ugly(t *T) {
	subject := BuildLinuxKit
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestContainer_ServeProduction_Good(t *T) {
	subject := ServeProduction
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestContainer_ServeProduction_Bad(t *T) {
	subject := ServeProduction
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestContainer_ServeProduction_Ugly(t *T) {
	subject := ServeProduction
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestContainer_Shell_Good(t *T) {
	subject := Shell
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestContainer_Shell_Bad(t *T) {
	subject := Shell
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestContainer_Shell_Ugly(t *T) {
	subject := Shell
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestContainer_IsPHPProject_Good(t *T) {
	subject := IsPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestContainer_IsPHPProject_Bad(t *T) {
	subject := IsPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestContainer_IsPHPProject_Ugly(t *T) {
	subject := IsPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
