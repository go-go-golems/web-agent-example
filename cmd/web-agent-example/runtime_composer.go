package main

import (
	"context"

	gepmiddleware "github.com/go-go-golems/geppetto/pkg/inference/middleware"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	infruntime "github.com/go-go-golems/pinocchio/pkg/inference/runtime"

	"github.com/go-go-golems/web-agent-example/pkg/discodialogue"
	"github.com/go-go-golems/web-agent-example/pkg/thinkingmode"
)

const (
	defaultRuntimeKey   = "default"
	defaultSystemPrompt = "You are a helpful assistant."
)

type webAgentRuntimeComposer struct {
	parsed *values.Values
}

func newWebAgentRuntimeComposer(parsed *values.Values) *webAgentRuntimeComposer {
	return &webAgentRuntimeComposer{parsed: parsed}
}

func (c *webAgentRuntimeComposer) Compose(
	ctx context.Context,
	req infruntime.ConversationRuntimeRequest,
) (infruntime.ComposedRuntime, error) {
	stepSettings, err := settings.NewStepSettingsFromParsedValues(c.parsed)
	if err != nil {
		return infruntime.ComposedRuntime{}, err
	}

	middlewares := []gepmiddleware.Middleware{
		thinkingmode.NewMiddleware(thinkingmode.DefaultConfig()),
		discodialogue.NewMiddleware(discodialogue.DefaultConfig()),
	}

	engine, err := infruntime.BuildEngineFromSettingsWithMiddlewares(
		ctx,
		stepSettings,
		defaultSystemPrompt,
		middlewares,
	)
	if err != nil {
		return infruntime.ComposedRuntime{}, err
	}

	runtimeKey := req.ProfileKey
	if runtimeKey == "" {
		runtimeKey = defaultRuntimeKey
	}

	return infruntime.ComposedRuntime{
		Engine:             engine,
		RuntimeKey:         runtimeKey,
		RuntimeFingerprint: runtimeKey,
		SeedSystemPrompt:   defaultSystemPrompt,
	}, nil
}
