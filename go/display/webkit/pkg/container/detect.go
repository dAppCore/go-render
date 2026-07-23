package container

import (
	"runtime"
	"strconv"

	core "dappco.re/go"
)

// ContainerRuntime describes the preferred isolated workload runtime for CoreGUI.
//
//	runtime := container.Detect()
//	if runtime == container.RuntimeApple { /* prefer Apple Containers */ }
type ContainerRuntime string

const (
	RuntimeNone   ContainerRuntime = ""
	RuntimeApple  ContainerRuntime = "apple"
	RuntimeDocker ContainerRuntime = "docker"
	RuntimePodman ContainerRuntime = "podman"
)

// DetectEnvironment controls runtime detection for tests and other callers.
//
//	runtime := container.DetectWithEnvironment(container.DetectEnvironment{
//	    GOOS:           "darwin",
//	    ProductVersion: "26.0",
//	})
type DetectEnvironment struct {
	GOOS           string
	ProductVersion string
	LookPath       func(file string) (string, resultFailure)
}

// Detect prefers Apple Containers on macOS 26+, then Docker, then Podman.
//
//	runtime := container.Detect()
func Detect() ContainerRuntime {
	environment := DetectEnvironment{
		GOOS:           runtime.GOOS,
		ProductVersion: "",
		LookPath:       lookPath,
	}
	if runtime.GOOS == "darwin" {
		environment.ProductVersion = productVersion()
	}
	return DetectWithEnvironment(environment)
}

// DetectWithEnvironment applies the RFC runtime ordering using an explicit environment.
//
//	runtime := container.DetectWithEnvironment(env)
func DetectWithEnvironment(environment DetectEnvironment) ContainerRuntime {
	resolvePath := environment.LookPath
	if resolvePath == nil {
		resolvePath = lookPath
	}

	goos := core.Lower(core.Trim(environment.GOOS))
	if goos == "darwin" && majorVersion(environment.ProductVersion) >= 26 {
		if hasBinary(resolvePath, "container") || hasBinary(resolvePath, "apple-container") || hasBinary(resolvePath, "containerctl") {
			return RuntimeApple
		}
	}
	if hasBinary(resolvePath, "docker") {
		return RuntimeDocker
	}
	if hasBinary(resolvePath, "podman") {
		return RuntimePodman
	}
	return RuntimeNone
}

func hasBinary(lookPath func(string) (string, resultFailure), binary string) bool {
	if core.Trim(binary) == "" {
		return false
	}
	_, err := lookPath(binary)
	return err == nil
}

func majorVersion(productVersion string) int {
	productVersion = core.Trim(productVersion)
	if productVersion == "" {
		return 0
	}
	major, _, _ := cut(productVersion, ".")
	value, err := strconv.Atoi(major)
	if err != nil {
		return 0
	}
	return value
}

func productVersion() string {
	output, err := command("sw_vers", "-productVersion").Output()
	if err != nil {
		return ""
	}
	return core.Trim(string(output))
}
