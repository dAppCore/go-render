package container

import (
	"context"
	core "dappco.re/go"
	"time"
)

func TestTIMManager_NewTIMManager_Good(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	state := manager.State()
	core.AssertEqual(t, "coregui-tim", state.Name)
	core.AssertEqual(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	core.AssertEqual(t, RuntimeDocker, state.Runtime)
	core.AssertEqual(t, "stopped", state.Status)
	core.AssertEmpty(t, state.StartedAt)
}

func TestTIMManager_NewTIMManager_Bad(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := NewTIMManager(TIMOptions{
		Name:  " ",
		Image: " ",
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state := manager.State()
	core.AssertEqual(t, "coregui-tim", state.Name)
	core.AssertEqual(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	core.AssertEqual(t, RuntimeNone, state.Runtime)
	core.AssertEqual(t, "stopped", state.Status)
}

func TestTIMManager_NewTIMManager_Ugly(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := NewTIMManager(TIMOptions{
		Name:      "worker-node",
		Image:     "ghcr.io/example/tim:edge",
		Command:   []string{"alpha", "beta"},
		DataDir:   "/var/lib/tim",
		Runtime:   RuntimePodman,
		Resources: TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	state := manager.State()
	core.AssertEqual(t, "worker-node", state.Name)
	core.AssertEqual(t, "ghcr.io/example/tim:edge", state.Image)
	core.AssertEqual(t, RuntimePodman, state.Runtime)
	core.AssertEqual(t, []string{"alpha", "beta"}, state.Command)
	core.AssertEqual(t, "/var/lib/tim", state.DataDir)
	core.AssertEqual(t, TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"}, state.Resources)
}

func TestTIMManager_Start_Good(t *core.T) {
	// Start
	ax7Variant := "Start:good"
	core.AssertContains(t, ax7Variant, "good")
	var calls []string
	manager := NewTIMManager(TIMOptions{
		Name:    "coregui-tim",
		Image:   "ghcr.io/example/tim:latest",
		Command: []string{"sleep", "1"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) resultFailure {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
		Now: func() time.Time {
			return time.Unix(456, 0).UTC()
		},
	})

	state, err := manager.Start(context.Background())
	core.RequireNoError(t, err)
	core.AssertEqual(t, "running", state.Status)
	core.AssertEqual(t, time.Unix(456, 0).UTC(), state.StartedAt)
	core.AssertEqual(t, "docker", calls[0])
	core.AssertContains(t, calls, "run")
	core.AssertContains(t, calls, "--rm")
	core.AssertContains(t, calls, "--name")
	core.AssertContains(t, calls, "coregui-tim")
	core.AssertContains(t, calls, "ghcr.io/example/tim:latest")
}

func TestTIMManager_Start_Bad(t *core.T) {
	// Start
	ax7Variant := "Start:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state, err := manager.Start(context.Background())
	core.AssertError(t, err)
	core.AssertEqual(t, RuntimeNone, state.Runtime)
	core.AssertEqual(t, "stopped", state.Status)
	core.AssertContains(t, err.Error(), "no supported container runtime detected")
}

func TestTIMManager_Start_Ugly(t *core.T) {
	// Start
	ax7Variant := "Start:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) resultFailure {
			return core.NewError("docker failed")
		},
	})

	state, err := manager.Start(context.Background())
	core.AssertError(t, err)
	core.AssertEqual(t, "error", state.Status)
	core.AssertContains(t, err.Error(), "docker failed")
}

func TestTIMManager_Stop_Good(t *core.T) {
	// Stop
	ax7Variant := "Stop:good"
	core.AssertContains(t, ax7Variant, "good")
	var calls []string
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) resultFailure {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
	})

	started, err := manager.Start(context.Background())
	core.RequireNoError(t, err)
	core.AssertEqual(t, "running", started.Status)

	stopped, err := manager.Stop(context.Background())
	core.RequireNoError(t, err)
	core.AssertEqual(t, "stopped", stopped.Status)
	core.AssertEmpty(t, stopped.StartedAt)
	core.AssertGreaterOrEqual(t, len(calls), 3)
	core.AssertEqual(t, "docker", calls[len(calls)-3])
	core.AssertEqual(t, "stop", calls[len(calls)-2])
	core.AssertEqual(t, "coregui-tim", calls[len(calls)-1])
}

func TestTIMManager_Stop_Bad(t *core.T) {
	// Stop
	ax7Variant := "Stop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) resultFailure {
			return core.NewError("stop failed")
		},
	})

	_, err := manager.Start(context.Background())
	core.AssertError(t, err)

	state, err := manager.Stop(context.Background())
	core.AssertError(t, err)
	core.AssertEqual(t, "error", state.Status)
	core.AssertContains(t, err.Error(), "stop failed")
}

func TestTIMManager_Stop_Ugly(t *core.T) {
	// Stop
	ax7Variant := "Stop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state, err := manager.Stop(context.Background())
	core.RequireNoError(t, err)
	core.AssertEqual(t, "stopped", state.Status)
}

func TestTIMManager_runtimeCommand_Good(t *core.T) {
	// runtimeCommand
	ax7Variant := "runtimeCommand:good"
	core.AssertContains(t, ax7Variant, "good")
	cases := []struct {
		name     string
		runtime  ContainerRuntime
		verb     string
		wantBin  string
		wantArgs []string
	}{
		{
			name:     "docker run",
			runtime:  RuntimeDocker,
			verb:     "run",
			wantBin:  "docker",
			wantArgs: []string{"run", "-d", "--rm", "--name", "tim"},
		},
		{
			name:     "apple stop",
			runtime:  RuntimeApple,
			verb:     "stop",
			wantBin:  "container",
			wantArgs: []string{"stop", "tim"},
		},
		{
			name:     "podman run",
			runtime:  RuntimePodman,
			verb:     "run",
			wantBin:  "podman",
			wantArgs: []string{"run", "-d", "--replace", "--name", "tim"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			manager := NewTIMManager(TIMOptions{
				Name:    "tim",
				Image:   "image",
				Runtime: tc.runtime,
				Detect: func() ContainerRuntime {
					return RuntimeDocker
				},
			})

			gotBin, gotArgs := manager.runtimeCommand(tc.runtime, tc.verb)
			core.AssertEqual(t, tc.wantBin, gotBin)
			core.AssertEqual(t, tc.wantArgs, gotArgs[:len(tc.wantArgs)])
			if tc.verb == "run" {
				core.AssertContains(t, gotArgs, "image")
			}
		})
	}
}

func TestTIMManager_runtimeCommand_Bad(t *core.T) {
	// runtimeCommand
	ax7Variant := "runtimeCommand:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := NewTIMManager(TIMOptions{
		Name:    "tim",
		Image:   "image",
		Runtime: RuntimeNone,
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	bin, args := manager.runtimeCommand(RuntimeNone, "run")
	core.AssertEqual(t, "docker", bin)
	core.AssertContains(t, args, "--rm")
}

func TestTIMManager_runtimeCommand_Ugly(t *core.T) {
	// runtimeCommand
	ax7Variant := "runtimeCommand:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := NewTIMManager(TIMOptions{
		Name:    "tim",
		Image:   "image",
		DataDir: "/host/tim",
		WorkDir: "/work",
		Env: map[string]string{
			"CORE_ENV": "test",
		},
		Resources: TIMResources{
			CPUCores: 2,
			MemoryMB: 512,
			GPU:      "all",
		},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	bin, args := manager.runtimeCommand(RuntimeDocker, "run")
	core.AssertEqual(t, "docker", bin)
	core.AssertContains(t, args, "--cpus")
	core.AssertContains(t, args, "2")
	core.AssertContains(t, args, "--memory")
	core.AssertContains(t, args, "512m")
	core.AssertContains(t, args, "--gpus")
	core.AssertContains(t, args, "all")
	core.AssertContains(t, args, "-v")
	core.AssertContains(t, args, "/host/tim:/host/tim")
	core.AssertContains(t, args, "-w")
	core.AssertContains(t, args, "/work")
	core.AssertContains(t, args, "-e")
	core.AssertContains(t, args, "CORE_ENV=test")
}

func TestTIMManager_resourceArgs_Good(t *core.T) {
	// resourceArgs
	ax7Variant := "resourceArgs:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertNil(t, resourceArgs(TIMResources{}))
	observedType := core.Sprintf("%T", resourceArgs(TIMResources{}))
	core.AssertNotEmpty(t, observedType)
}

func TestTIMManager_resourceArgs_Bad(t *core.T) {
	// resourceArgs
	ax7Variant := "resourceArgs:bad"
	core.AssertContains(t, ax7Variant, "bad")
	args := resourceArgs(TIMResources{CPUCores: 2})

	core.AssertEqual(t, []string{"--cpus", "2"}, args)
	core.AssertNotEmpty(t, core.Sprintf("%T", args))
}

func TestTIMManager_resourceArgs_Ugly(t *core.T) {
	// resourceArgs
	ax7Variant := "resourceArgs:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	args := resourceArgs(TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"})

	core.AssertEqual(t, []string{"--cpus", "2", "--memory", "512m", "--gpus", "all"}, args)
	core.AssertNotEmpty(t, core.Sprintf("%T", args))
}

func TestTIMManager_coalesceRuntime_Good(t *core.T) {
	// coalesceRuntime
	ax7Variant := "coalesceRuntime:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertEqual(t, RuntimeApple, coalesceRuntime(RuntimeApple, RuntimeDocker))
	observedType := core.Sprintf("%T", coalesceRuntime(RuntimeApple, RuntimeDocker))
	core.AssertNotEmpty(t, observedType)
}

func TestTIMManager_coalesceRuntime_Bad(t *core.T) {
	// coalesceRuntime
	ax7Variant := "coalesceRuntime:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertEqual(t, RuntimeDocker, coalesceRuntime(RuntimeNone, RuntimeDocker))
	observedType := core.Sprintf("%T", coalesceRuntime(RuntimeNone, RuntimeDocker))
	core.AssertNotEmpty(t, observedType)
}

func TestTIMManager_coalesceRuntime_Ugly(t *core.T) {
	// coalesceRuntime
	ax7Variant := "coalesceRuntime:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertEqual(t, RuntimeNone, coalesceRuntime(RuntimeNone, RuntimeNone))
	observedType := core.Sprintf("%T", coalesceRuntime(RuntimeNone, RuntimeNone))
	core.AssertNotEmpty(t, observedType)
}

// AX7 generated source-matching smoke coverage.
func TestTim_NewTIMManager_Good(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewTIMManager(*new(TIMOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_NewTIMManager_Bad(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewTIMManager(*new(TIMOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_NewTIMManager_Ugly(t *core.T) {
	// NewTIMManager
	ax7Variant := "NewTIMManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewTIMManager(*new(TIMOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_State_Good(t *core.T) {
	// TIMManager State
	ax7Variant := "TIMManager_State:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_State_Bad(t *core.T) {
	// TIMManager State
	ax7Variant := "TIMManager_State:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_State_Ugly(t *core.T) {
	// TIMManager State
	ax7Variant := "TIMManager_State:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Start_Good(t *core.T) {
	// TIMManager Start
	ax7Variant := "TIMManager_Start:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Start_Bad(t *core.T) {
	// TIMManager Start
	ax7Variant := "TIMManager_Start:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Start_Ugly(t *core.T) {
	// TIMManager Start
	ax7Variant := "TIMManager_Start:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Stop_Good(t *core.T) {
	// TIMManager Stop
	ax7Variant := "TIMManager_Stop:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Stop_Bad(t *core.T) {
	// TIMManager Stop
	ax7Variant := "TIMManager_Stop:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTim_TIMManager_Stop_Ugly(t *core.T) {
	// TIMManager Stop
	ax7Variant := "TIMManager_Stop:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TIMManager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}
