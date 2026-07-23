package php

func TestDockerfile_GenerateDockerfile_Good(t *T) {
	subject := GenerateDockerfile
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDockerfile_GenerateDockerfile_Bad(t *T) {
	subject := GenerateDockerfile
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDockerfile_GenerateDockerfile_Ugly(t *T) {
	subject := GenerateDockerfile
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDockerfile_DetectDockerfileConfig_Good(t *T) {
	subject := DetectDockerfileConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDockerfile_DetectDockerfileConfig_Bad(t *T) {
	subject := DetectDockerfileConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDockerfile_DetectDockerfileConfig_Ugly(t *T) {
	subject := DetectDockerfileConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDockerfile_GenerateDockerfileFromConfig_Good(t *T) {
	subject := GenerateDockerfileFromConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDockerfile_GenerateDockerfileFromConfig_Bad(t *T) {
	subject := GenerateDockerfileFromConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDockerfile_GenerateDockerfileFromConfig_Ugly(t *T) {
	subject := GenerateDockerfileFromConfig
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestDockerfile_GenerateDockerignore_Good(t *T) {
	subject := GenerateDockerignore
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestDockerfile_GenerateDockerignore_Bad(t *T) {
	subject := GenerateDockerignore
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestDockerfile_GenerateDockerignore_Ugly(t *T) {
	subject := GenerateDockerignore
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
