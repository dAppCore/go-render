package php

func TestQuality_DetectFormatter_Good(t *T) {
	subject := DetectFormatter
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_DetectFormatter_Bad(t *T) {
	subject := DetectFormatter
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_DetectFormatter_Ugly(t *T) {
	subject := DetectFormatter
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_DetectAnalyser_Good(t *T) {
	subject := DetectAnalyser
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_DetectAnalyser_Bad(t *T) {
	subject := DetectAnalyser
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_DetectAnalyser_Ugly(t *T) {
	subject := DetectAnalyser
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_Format_Good(t *T) {
	subject := Format
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_Format_Bad(t *T) {
	subject := Format
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_Format_Ugly(t *T) {
	subject := Format
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_Analyse_Good(t *T) {
	subject := Analyse
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_Analyse_Bad(t *T) {
	subject := Analyse
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_Analyse_Ugly(t *T) {
	subject := Analyse
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_DetectPsalm_Good(t *T) {
	subject := DetectPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_DetectPsalm_Bad(t *T) {
	subject := DetectPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_DetectPsalm_Ugly(t *T) {
	subject := DetectPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_RunPsalm_Good(t *T) {
	subject := RunPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_RunPsalm_Bad(t *T) {
	subject := RunPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_RunPsalm_Ugly(t *T) {
	subject := RunPsalm
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_RunAudit_Good(t *T) {
	subject := RunAudit
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_RunAudit_Bad(t *T) {
	subject := RunAudit
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_RunAudit_Ugly(t *T) {
	subject := RunAudit
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_DetectRector_Good(t *T) {
	subject := DetectRector
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_DetectRector_Bad(t *T) {
	subject := DetectRector
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_DetectRector_Ugly(t *T) {
	subject := DetectRector
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_RunRector_Good(t *T) {
	subject := RunRector
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_RunRector_Bad(t *T) {
	subject := RunRector
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_RunRector_Ugly(t *T) {
	subject := RunRector
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_DetectInfection_Good(t *T) {
	subject := DetectInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_DetectInfection_Bad(t *T) {
	subject := DetectInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_DetectInfection_Ugly(t *T) {
	subject := DetectInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_RunInfection_Good(t *T) {
	subject := RunInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_RunInfection_Bad(t *T) {
	subject := RunInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_RunInfection_Ugly(t *T) {
	subject := RunInfection
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_GetQAStages_Good(t *T) {
	subject := GetQAStages
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_GetQAStages_Bad(t *T) {
	subject := GetQAStages
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_GetQAStages_Ugly(t *T) {
	subject := GetQAStages
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_GetQAChecks_Good(t *T) {
	subject := GetQAChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_GetQAChecks_Bad(t *T) {
	subject := GetQAChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_GetQAChecks_Ugly(t *T) {
	subject := GetQAChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestQuality_RunSecurityChecks_Good(t *T) {
	subject := RunSecurityChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestQuality_RunSecurityChecks_Bad(t *T) {
	subject := RunSecurityChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestQuality_RunSecurityChecks_Ugly(t *T) {
	subject := RunSecurityChecks
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
