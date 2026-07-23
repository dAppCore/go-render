package php

func TestCmd_SetMedium_Good(t *T) {
	subject := SetMedium
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmd_SetMedium_Bad(t *T) {
	subject := SetMedium
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmd_SetMedium_Ugly(t *T) {
	subject := SetMedium
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestCmd_AddPHPCommands_Good(t *T) {
	subject := AddPHPCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmd_AddPHPCommands_Bad(t *T) {
	subject := AddPHPCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmd_AddPHPCommands_Ugly(t *T) {
	subject := AddPHPCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestCmd_AddPHPRootCommands_Good(t *T) {
	subject := AddPHPRootCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestCmd_AddPHPRootCommands_Bad(t *T) {
	subject := AddPHPRootCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestCmd_AddPHPRootCommands_Ugly(t *T) {
	subject := AddPHPRootCommands
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
