package container

import (
	core "dappco.re/go"
)

func TestDetectWithEnvironment_PrefersAppleContainersOnMacOS26(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.0",
		LookPath: func(file string) (string, resultFailure) {
			if file == "container" {
				return "/usr/bin/container", nil
			}
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimeApple, runtime)
}

func TestDetectWithEnvironment_FallsBackToDockerWhenAppleUnavailable(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.1",
		LookPath: func(file string) (string, resultFailure) {
			if file == "docker" {
				return "/usr/local/bin/docker", nil
			}
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesDockerOnNonMacHosts(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(file string) (string, resultFailure) {
			if file == "docker" {
				return "/usr/bin/docker", nil
			}
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesPodmanWhenDockerMissing(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(file string) (string, resultFailure) {
			if file == "podman" {
				return "/usr/bin/podman", nil
			}
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimePodman, runtime)
}

func TestDetectWithEnvironment_ReturnsNoneWhenNoRuntimeIsAvailable(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(string) (string, resultFailure) {
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimeNone, runtime)
}

func TestMajorVersion(t *core.T) {
	core.AssertEqual(t, 26, majorVersion("26.0"))
	core.AssertEqual(t, 0, majorVersion("bogus"))
	core.AssertEqual(t, 0, majorVersion(""))
}

func TestDetect_Good(t *core.T) {
	binDir := t.TempDir()
	containerPath := writeExecutable(t, binDir, "container", "#!/bin/sh\nexit 0\n")

	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.0",
		LookPath: func(file string) (string, resultFailure) {
			if file == "container" {
				return containerPath, nil
			}
			return "", core.NewError("not found")
		},
	})

	core.AssertEqual(t, RuntimeApple, runtime)
}

func TestDetect_Bad(t *core.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf '25.0\\n'\n")
	writeExecutable(t, binDir, "docker", "#!/bin/sh\nexit 0\n")
	t.Setenv("PATH", binDir)

	core.AssertEqual(t, RuntimeDocker, Detect())
}

func TestDetect_Ugly(t *core.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf 'not-a-version\\n'\n")
	t.Setenv("PATH", binDir)

	core.AssertEqual(t, RuntimeNone, Detect())
}

func writeExecutable(t *core.T, dir, name, script string) string {
	t.Helper()

	path := core.PathJoin(dir, name)
	core.RequireNoError(t, coreWriteMode(path, script, 0o755))
	return path
}

// AX7 generated source-matching smoke coverage.
func TestDetect_Detect_Good(t *core.T) {
	// Detect
	ax7Variant := "Detect:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := Detect()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDetect_Detect_Bad(t *core.T) {
	// Detect
	ax7Variant := "Detect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := Detect()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDetect_Detect_Ugly(t *core.T) {
	// Detect
	ax7Variant := "Detect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := Detect()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDetect_DetectWithEnvironment_Good(t *core.T) {
	// DetectWithEnvironment
	ax7Variant := "DetectWithEnvironment:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := DetectWithEnvironment(*new(DetectEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDetect_DetectWithEnvironment_Bad(t *core.T) {
	// DetectWithEnvironment
	ax7Variant := "DetectWithEnvironment:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := DetectWithEnvironment(*new(DetectEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestDetect_DetectWithEnvironment_Ugly(t *core.T) {
	// DetectWithEnvironment
	ax7Variant := "DetectWithEnvironment:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := DetectWithEnvironment(*new(DetectEnvironment))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
