package container

import (
	"context"
	"regexp"
	"unicode"

	core "dappco.re/go"
)

type Service struct {
	*core.ServiceRuntime[TIMOptions]
	manager    *TIMManager
	startupErr error
}

var timContainerNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

func NewService(c *core.Core, options TIMOptions) *Service {
	options, err := normalizeTIMOptions(options)
	if options.Exec == nil && c != nil {
		options.Exec = func(ctx context.Context, name string, args ...string) resultFailure {
			result := c.Process().Run(ctx, name, args...)
			if result.OK {
				return nil
			}
			if err, ok := result.Value.(error); ok {
				return err
			}
			return core.E("container.TIMManager.Exec", result.Error(), nil)
		}
	}
	return &Service{
		ServiceRuntime: core.NewServiceRuntime(c, options),
		manager:        NewTIMManager(options),
		startupErr:     err,
	}
}

func OptionsFromEnv() TIMOptions {
	options, _ := OptionsFromEnvValidated()
	return options
}

func OptionsFromEnvValidated() (TIMOptions, resultFailure) {
	return TIMOptions{
		Name:    core.Trim(core.Env("CORE_TIM_NAME")),
		Image:   core.Trim(core.Env("CORE_TIM_IMAGE")),
		Command: splitCSV(core.Trim(core.Env("CORE_TIM_COMMAND"))),
		DataDir: core.Trim(core.Env("CORE_TIM_DATA_DIR")),
		Resources: TIMResources{
			GPU: core.Trim(core.Env("CORE_TIM_GPU")),
		},
	}.Validate()
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	if s.startupErr != nil {
		return core.Result{Value: s.startupErr, OK: false}
	}
	s.Core().Action("container.runtime.detect", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: coalesceRuntime(s.manager.options.Runtime, s.manager.options.Detect()), OK: true}
	})
	s.Core().Action("tim.start", func(ctx context.Context, _ core.Options) core.Result {
		state, err := s.manager.Start(ctx)
		return core.Result{}.New(state, err)
	})
	s.Core().Action("tim.stop", func(ctx context.Context, _ core.Options) core.Result {
		state, err := s.manager.Stop(ctx)
		return core.Result{}.New(state, err)
	})
	s.Core().Action("tim.status", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.manager.State(), OK: true}
	})
	return core.Result{OK: true}
}

func (options TIMOptions) Validate() (TIMOptions, resultFailure) {
	return normalizeTIMOptions(options)
}

func normalizeTIMOptions(options TIMOptions) (TIMOptions, resultFailure) {
	options.Name = core.Trim(options.Name)
	options.Image = core.Trim(options.Image)
	options.DataDir = core.Trim(options.DataDir)
	options.Resources.GPU = core.Trim(options.Resources.GPU)

	if err := validateTIMContainerName(options.Name); err != nil {
		return options, err
	}
	if err := validateTIMArgValue("image", options.Image); err != nil {
		return options, err
	}
	if err := validateTIMArgValue("gpu", options.Resources.GPU); err != nil {
		return options, err
	}
	return options, nil
}

func validateTIMContainerName(value string) resultFailure {
	if value == "" {
		return nil
	}
	if core.HasPrefix(value, "-") {
		return core.E("container.validateTIMOptions", "name cannot start with -", nil)
	}
	if core.HasPrefix(value, ".") {
		return core.E("container.validateTIMOptions", "name cannot start with .", nil)
	}
	if !timContainerNamePattern.MatchString(value) {
		return core.E("container.validateTIMOptions", "name must contain only letters, digits, underscores, dots, and hyphens", nil)
	}
	return nil
}

func validateTIMArgValue(label, value string) resultFailure {
	if value == "" {
		return nil
	}
	if core.HasPrefix(value, "-") {
		return core.E("container.validateTIMOptions", label+" cannot start with -", nil)
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return core.E("container.validateTIMOptions", label+" contains invalid control characters", nil)
		}
		if unicode.IsSpace(r) {
			return core.E("container.validateTIMOptions", label+" cannot contain whitespace", nil)
		}
	}
	return nil
}

func (s *Service) State() TIMState {
	return s.manager.State()
}

func splitCSV(value string) []string {
	if core.Trim(value) == "" {
		return nil
	}
	parts := core.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = core.Trim(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
