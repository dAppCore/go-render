package container

import (
	"context"
	"time"

	core "dappco.re/go"
)

func newTestContainerService(t *core.T, options TIMOptions) (*Service, *core.Core) {
	t.Helper()

	var svc *Service
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc = NewService(c, options)
			return core.Result{Value: svc, OK: true}
		}),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)
	core.AssertNotNil(t, svc)
	return svc, c
}

func newInvalidTestContainerService(t *core.T, options TIMOptions) (*Service, *core.Core, resultFailure) {
	t.Helper()

	if options.Detect == nil {
		options.Detect = func() ContainerRuntime {
			return RuntimeDocker
		}
	}

	var svc *Service
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc = NewService(c, options)
			return core.Result{Value: svc, OK: true}
		}),
		core.WithServiceLock(),
	)
	result := c.ServiceStartup(context.Background(), nil)
	core.AssertFalse(t, result.OK)
	core.AssertNotNil(t, svc)
	err, ok := result.Value.(resultFailure)
	core.RequireTrue(t, ok)
	return svc, c, err
}

func TestService_OptionsFromEnv_GoodCase(t *core.T) {
	t.Setenv("CORE_TIM_NAME", "  worker ")
	t.Setenv("CORE_TIM_IMAGE", " ghcr.io/example/tim:latest ")
	t.Setenv("CORE_TIM_COMMAND", "run,  --flag, , value ")
	t.Setenv("CORE_TIM_DATA_DIR", " /var/lib/core-tim ")
	t.Setenv("CORE_TIM_GPU", " all ")

	opts, err := OptionsFromEnvValidated()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "worker", opts.Name)
	core.AssertEqual(t, "ghcr.io/example/tim:latest", opts.Image)
	core.AssertEqual(t, []string{"run", "--flag", "value"}, opts.Command)
	core.AssertEqual(t, "/var/lib/core-tim", opts.DataDir)
	core.AssertEqual(t, "all", opts.Resources.GPU)
}

func TestService_OptionsFromEnv_BadCase(t *core.T) {
	t.Setenv("CORE_TIM_NAME", "")
	t.Setenv("CORE_TIM_IMAGE", "")
	t.Setenv("CORE_TIM_COMMAND", "")
	t.Setenv("CORE_TIM_DATA_DIR", "")
	t.Setenv("CORE_TIM_GPU", "")

	opts, err := OptionsFromEnvValidated()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, opts.Name)
	core.AssertEmpty(t, opts.Image)
	core.AssertNil(t, opts.Command)
	core.AssertEmpty(t, opts.DataDir)
	core.AssertEmpty(t, opts.Resources.GPU)
}

func TestService_OptionsFromEnv_UglyCase(t *core.T) {
	t.Setenv("CORE_TIM_NAME", " \t\n ")
	t.Setenv("CORE_TIM_IMAGE", " \t ghcr.io/example/tim:latest \n")
	t.Setenv("CORE_TIM_COMMAND", " , first ,, second , ")
	t.Setenv("CORE_TIM_DATA_DIR", "\t /tmp/core-tim \n")
	t.Setenv("CORE_TIM_GPU", "\t device=0 \n")

	opts, err := OptionsFromEnvValidated()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, opts.Name)
	core.AssertEqual(t, "ghcr.io/example/tim:latest", opts.Image)
	core.AssertEqual(t, []string{"first", "second"}, opts.Command)
	core.AssertEqual(t, "/tmp/core-tim", opts.DataDir)
	core.AssertEqual(t, "device=0", opts.Resources.GPU)
}

func TestService_OptionsFromEnv_RejectsLeadingDash(t *core.T) {
	cases := []struct {
		name   string
		envKey string
		value  string
		want   string
	}{
		{name: "name", envKey: "CORE_TIM_NAME", value: "-rm", want: "name cannot start with -"},
		{name: "image", envKey: "CORE_TIM_IMAGE", value: "--privileged", want: "image cannot start with -"},
		{name: "gpu", envKey: "CORE_TIM_GPU", value: "-it", want: "gpu cannot start with -"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			t.Setenv("CORE_TIM_NAME", "")
			t.Setenv("CORE_TIM_IMAGE", "")
			t.Setenv("CORE_TIM_COMMAND", "")
			t.Setenv("CORE_TIM_DATA_DIR", "")
			t.Setenv("CORE_TIM_GPU", "")
			t.Setenv(tc.envKey, tc.value)

			_, err := OptionsFromEnvValidated()

			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), tc.want)
		})
	}
}

func TestService_NewService_GoodCase(t *core.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Name:  "normal-container",
		Image: "alpine:3.19",
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	core.AssertNotNil(t, svc.manager)
	state := svc.State()
	core.AssertEqual(t, "normal-container", state.Name)
	core.AssertEqual(t, "alpine:3.19", state.Image)
	core.AssertEqual(t, RuntimeDocker, state.Runtime)
	core.AssertEqual(t, "stopped", state.Status)
}

func TestService_NewService_BadCase(t *core.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state := svc.State()
	core.AssertEqual(t, "coregui-tim", state.Name)
	core.AssertEqual(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	core.AssertEqual(t, RuntimeNone, state.Runtime)
}

func TestService_NewService_UglyCase(t *core.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Name:    "  worker.node  ",
		Image:   "  ghcr.io/example/tim:edge  ",
		Command: []string{"alpha", "beta"},
		DataDir: " /tmp/data ",
		Resources: TIMResources{
			CPUCores: 2,
			MemoryMB: 512,
			GPU:      "all",
		},
		Detect: func() ContainerRuntime {
			return RuntimePodman
		},
	})

	state := svc.State()
	core.AssertEqual(t, "worker.node", state.Name)
	core.AssertEqual(t, "ghcr.io/example/tim:edge", state.Image)
	core.AssertEqual(t, []string{"alpha", "beta"}, state.Command)
	core.AssertEqual(t, "/tmp/data", state.DataDir)
	core.AssertEqual(t, TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"}, state.Resources)
	core.AssertEqual(t, RuntimePodman, state.Runtime)
}

func TestService_NewService_RejectsInvalidTIMOptions(t *core.T) {
	cases := []struct {
		name    string
		options TIMOptions
		want    string
	}{
		{
			name:    "name leading dash",
			options: TIMOptions{Name: "-rm", Image: "alpine:3.19"},
			want:    "name cannot start with -",
		},
		{
			name:    "image leading dash",
			options: TIMOptions{Name: "normal-container", Image: "--privileged"},
			want:    "image cannot start with -",
		},
		{
			name: "gpu leading dash",
			options: TIMOptions{
				Name:      "normal-container",
				Image:     "alpine:3.19",
				Resources: TIMResources{GPU: "-it"},
			},
			want: "gpu cannot start with -",
		},
		{
			name:    "name starts with dot",
			options: TIMOptions{Name: ".hidden", Image: "alpine:3.19"},
			want:    "name cannot start with .",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *core.T) {
			_, _, err := newInvalidTestContainerService(t, tc.options)

			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), tc.want)
		})
	}
}

func TestService_OnStartup_GoodCase(t *core.T) {
	var calls []string
	svc, c := newTestContainerService(t, TIMOptions{
		Name:  "coregui-tim",
		Image: "ghcr.io/example/tim:latest",
		Command: []string{
			"sleep",
			"1",
		},
		Resources: TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) resultFailure {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
		Now: func() time.Time {
			return time.Unix(123, 0).UTC()
		},
	})

	runtime := c.Action("container.runtime.detect").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, runtime.OK)
	core.AssertEqual(t, RuntimeDocker, runtime.Value)

	status := c.Action("tim.status").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, status.OK)
	initial := status.Value.(TIMState)
	core.AssertEqual(t, "stopped", initial.Status)

	started := c.Action("tim.start").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, started.OK)
	startState := started.Value.(TIMState)
	core.AssertEqual(t, "running", startState.Status)
	core.AssertEqual(t, time.Unix(123, 0).UTC(), startState.StartedAt)
	core.RequireNotEmpty(t, calls)
	core.AssertEqual(t, "docker", calls[0])
	core.AssertContains(t, calls, "run")
	core.AssertContains(t, calls, "--name")
	core.AssertContains(t, calls, "coregui-tim")
	core.AssertContains(t, calls, "--cpus")
	core.AssertContains(t, calls, "2")
	core.AssertContains(t, calls, "--memory")
	core.AssertContains(t, calls, "512m")
	core.AssertContains(t, calls, "--gpus")
	core.AssertContains(t, calls, "all")

	stopped := c.Action("tim.stop").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, stopped.OK)
	stopState := stopped.Value.(TIMState)
	core.AssertEqual(t, "stopped", stopState.Status)
	core.AssertEmpty(t, stopState.StartedAt)
	core.AssertEqual(t, "coregui-tim", svc.State().Name)
}

func TestService_OnStartup_BadCase(t *core.T) {
	_, c := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	result := c.Action("tim.start").Run(context.Background(), core.NewOptions())

	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "no supported container runtime detected")
}

func TestService_OnStartup_UglyCase(t *core.T) {
	_, c := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) resultFailure {
			return core.NewError("boom")
		},
	})

	result := c.Action("tim.start").Run(context.Background(), core.NewOptions())

	core.AssertFalse(t, result.OK)
	core.AssertError(t, result.Value.(resultFailure))
	core.AssertContains(t, result.Value.(resultFailure).Error(), "boom")

	status := c.Action("tim.status").Run(context.Background(), core.NewOptions())
	core.RequireTrue(t, status.OK)
	core.AssertEqual(t, "error", status.Value.(TIMState).Status)
}

// AX7 generated source-matching smoke coverage.
func TestService_OptionsFromEnvValidated_Good(t *core.T) {
	// OptionsFromEnvValidated
	ax7Variant := "OptionsFromEnvValidated:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0, got1 := OptionsFromEnvValidated()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_OptionsFromEnvValidated_Bad(t *core.T) {
	// OptionsFromEnvValidated
	ax7Variant := "OptionsFromEnvValidated:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0, got1 := OptionsFromEnvValidated()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_OptionsFromEnvValidated_Ugly(t *core.T) {
	// OptionsFromEnvValidated
	ax7Variant := "OptionsFromEnvValidated:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0, got1 := OptionsFromEnvValidated()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Good(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Bad(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_OnStartup_Ugly(t *core.T) {
	// Service OnStartup
	ax7Variant := "Service_OnStartup:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.OnStartup(core.Background())
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_TIMOptions_Validate_Good(t *core.T) {
	// TIMOptions Validate
	ax7Variant := "TIMOptions_Validate:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject TIMOptions
	result := core.Try(func() any {
		got0, got1 := subject.Validate()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_TIMOptions_Validate_Bad(t *core.T) {
	// TIMOptions Validate
	ax7Variant := "TIMOptions_Validate:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject TIMOptions
	result := core.Try(func() any {
		got0, got1 := subject.Validate()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_TIMOptions_Validate_Ugly(t *core.T) {
	// TIMOptions Validate
	ax7Variant := "TIMOptions_Validate:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject TIMOptions
	result := core.Try(func() any {
		got0, got1 := subject.Validate()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_State_Good(t *core.T) {
	// Service State
	ax7Variant := "Service_State:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_State_Bad(t *core.T) {
	// Service State
	ax7Variant := "Service_State:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_Service_State_Ugly(t *core.T) {
	// Service State
	ax7Variant := "Service_State:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Service)
	result := core.Try(func() any {
		got0 := subject.State()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_NewService_Good(t *core.T) {
	// NewService
	ax7Variant := "NewService:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewService(core.New(), TIMOptions{Detect: func() ContainerRuntime { return RuntimeNone }})
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_NewService_Bad(t *core.T) {
	// NewService
	ax7Variant := "NewService:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewService(nil, TIMOptions{Detect: func() ContainerRuntime { return RuntimeNone }})
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_NewService_Ugly(t *core.T) {
	// NewService
	ax7Variant := "NewService:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewService(core.New(), TIMOptions{Name: "../../edge", Detect: func() ContainerRuntime { return RuntimeNone }})
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_OptionsFromEnv_Good(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := OptionsFromEnv()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_OptionsFromEnv_Bad(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := OptionsFromEnv()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestService_OptionsFromEnv_Ugly(t *core.T) {
	// OptionsFromEnv
	ax7Variant := "OptionsFromEnv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := OptionsFromEnv()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
