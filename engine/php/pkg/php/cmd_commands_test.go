package php

func TestCmdCommands_AddCommands_Good(t *T) {
	subject := AddCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmdCommands_AddCommands_Bad(t *T) {
	subject := AddCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmdCommands_AddCommands_Ugly(t *T) {
	subject := AddCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
