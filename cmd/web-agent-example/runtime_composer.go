package main

import (
	"context"
	"strings"

	gepmiddleware "github.com/go-go-golems/geppetto/pkg/inference/middleware"
	gepprofiles "github.com/go-go-golems/geppetto/pkg/profiles"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"

	infruntime "github.com/go-go-golems/pinocchio/pkg/inference/runtime"

	"github.com/go-go-golems/web-agent-example/pkg/discodialogue"
	"github.com/go-go-golems/web-agent-example/pkg/thinkingmode"
)

type defaultRuntimeComposer struct {
	base *settings.StepSettings
}

func newDefaultRuntimeComposer(parsed *values.Values) (infruntime.RuntimeBuilder, error) {
	base, err := settings.NewStepSettingsFromParsedValues(parsed)
	if err != nil {
		return nil, errors.Wrap(err, "build step settings from parsed values")
	}
	return &defaultRuntimeComposer{base: base}, nil
}

func (c *defaultRuntimeComposer) Compose(ctx context.Context, req infruntime.ConversationRuntimeRequest) (infruntime.ComposedRuntime, error) {
	if ctx == nil {
		return infruntime.ComposedRuntime{}, errors.New("compose context is nil")
	}

	stepSettings := c.base
	if stepSettings == nil {
		var err error
		stepSettings, err = settings.NewStepSettings()
		if err != nil {
			return infruntime.ComposedRuntime{}, err
		}
	} else {
		stepSettings = stepSettings.Clone()
	}

	spec := cloneRuntimeSpec(req.ResolvedProfileRuntime)
	if spec == nil {
		spec = defaultRuntimeSpec()
	}

	systemPrompt := strings.TrimSpace(spec.SystemPrompt)
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant."
	}

	resolvedMiddlewares, err := resolveRuntimeMiddlewares(spec.Middlewares)
	if err != nil {
		return infruntime.ComposedRuntime{}, err
	}

	eng, err := infruntime.BuildEngineFromSettingsWithMiddlewares(ctx, stepSettings, systemPrompt, resolvedMiddlewares)
	if err != nil {
		return infruntime.ComposedRuntime{}, err
	}

	runtimeKey := strings.TrimSpace(req.ProfileKey)
	if runtimeKey == "" {
		runtimeKey = defaultRuntimeKey
	}
	runtimeFingerprint := strings.TrimSpace(req.ResolvedProfileFingerprint)
	if runtimeFingerprint == "" {
		runtimeFingerprint = defaultRuntimeFingerprint
	}

	return infruntime.ComposedRuntime{
		Engine:             eng,
		RuntimeKey:         runtimeKey,
		RuntimeFingerprint: runtimeFingerprint,
		SeedSystemPrompt:   systemPrompt,
	}, nil
}

func resolveRuntimeMiddlewares(uses []gepprofiles.MiddlewareUse) ([]gepmiddleware.Middleware, error) {
	middlewares := make([]gepmiddleware.Middleware, 0, len(uses))
	for _, use := range uses {
		if use.Enabled != nil && !*use.Enabled {
			continue
		}

		switch strings.TrimSpace(use.Name) {
		case "webagent-thinking-mode":
			middlewares = append(middlewares, thinkingmode.NewMiddleware(thinkingmode.ConfigFromAny(use.Config)))
		case "webagent-disco-dialogue":
			middlewares = append(middlewares, discodialogue.NewMiddleware(discodialogue.ConfigFromAny(use.Config)))
		case "":
			continue
		default:
			return nil, errors.Errorf("unknown middleware %q", use.Name)
		}
	}
	return middlewares, nil
}

func defaultRuntimeSpec() *gepprofiles.RuntimeSpec {
	return &gepprofiles.RuntimeSpec{
		SystemPrompt: "You are a helpful assistant.",
		Middlewares: []gepprofiles.MiddlewareUse{
			{
				Name:   "webagent-thinking-mode",
				Config: thinkingmode.DefaultConfig(),
			},
			{
				Name:   "webagent-disco-dialogue",
				Config: discodialogue.DefaultConfig(),
			},
		},
	}
}

func cloneRuntimeSpec(spec *gepprofiles.RuntimeSpec) *gepprofiles.RuntimeSpec {
	if spec == nil {
		return nil
	}

	cloned := &gepprofiles.RuntimeSpec{
		SystemPrompt: spec.SystemPrompt,
		Tools:        append([]string(nil), spec.Tools...),
	}
	if len(spec.StepSettingsPatch) > 0 {
		cloned.StepSettingsPatch = cloneMap(spec.StepSettingsPatch)
	}
	if len(spec.Middlewares) > 0 {
		cloned.Middlewares = make([]gepprofiles.MiddlewareUse, 0, len(spec.Middlewares))
		for _, use := range spec.Middlewares {
			copied := gepprofiles.MiddlewareUse{
				Name:   use.Name,
				ID:     use.ID,
				Config: cloneAny(use.Config),
			}
			if use.Enabled != nil {
				enabled := *use.Enabled
				copied.Enabled = &enabled
			}
			cloned.Middlewares = append(cloned.Middlewares, copied)
		}
	}
	return cloned
}

func cloneMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = cloneAny(v)
	}
	return out
}

func cloneAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return cloneMap(t)
	case []any:
		out := make([]any, 0, len(t))
		for _, item := range t {
			out = append(out, cloneAny(item))
		}
		return out
	default:
		return v
	}
}
