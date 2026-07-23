package php

func TestDetect_IsLaravelProject_Good(t *T) {
	subject := IsLaravelProject
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_IsLaravelProject_Bad(t *T) {
	subject := IsLaravelProject
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_IsLaravelProject_Ugly(t *T) {
	subject := IsLaravelProject
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_IsFrankenPHPProject_Good(t *T) {
	subject := IsFrankenPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_IsFrankenPHPProject_Bad(t *T) {
	subject := IsFrankenPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_IsFrankenPHPProject_Ugly(t *T) {
	subject := IsFrankenPHPProject
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_DetectServices_Good(t *T) {
	subject := DetectServices
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_DetectServices_Bad(t *T) {
	subject := DetectServices
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_DetectServices_Ugly(t *T) {
	subject := DetectServices
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_DetectPackageManager_Good(t *T) {
	subject := DetectPackageManager
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_DetectPackageManager_Bad(t *T) {
	subject := DetectPackageManager
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_DetectPackageManager_Ugly(t *T) {
	subject := DetectPackageManager
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_GetLaravelAppName_Good(t *T) {
	subject := GetLaravelAppName
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_GetLaravelAppName_Bad(t *T) {
	subject := GetLaravelAppName
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_GetLaravelAppName_Ugly(t *T) {
	subject := GetLaravelAppName
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_GetLaravelAppURL_Good(t *T) {
	subject := GetLaravelAppURL
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_GetLaravelAppURL_Bad(t *T) {
	subject := GetLaravelAppURL
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_GetLaravelAppURL_Ugly(t *T) {
	subject := GetLaravelAppURL
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDetect_ExtractDomainFromURL_Good(t *T) {
	subject := ExtractDomainFromURL
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDetect_ExtractDomainFromURL_Bad(t *T) {
	subject := ExtractDomainFromURL
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDetect_ExtractDomainFromURL_Ugly(t *T) {
	subject := ExtractDomainFromURL
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
