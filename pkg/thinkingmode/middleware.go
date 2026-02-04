package thinkingmode

import (
	"context"
	"strings"

	gepevents "github.com/go-go-golems/geppetto/pkg/events"
	rootmw "github.com/go-go-golems/geppetto/pkg/inference/middleware"
	"github.com/go-go-golems/geppetto/pkg/turns"
	"github.com/google/uuid"
)

type Config struct {
	Mode      string
	Phase     string
	Reasoning string
}

func DefaultConfig() Config {
	return Config{
		Mode:      "fast",
		Phase:     "select",
		Reasoning: "thinking mode selected",
	}
}

func ConfigFromAny(cfg any) Config {
	out := DefaultConfig()
	switch v := cfg.(type) {
	case Config:
		return v
	case *Config:
		if v != nil {
			return *v
		}
	case map[string]any:
		if mode, ok := v["mode"].(string); ok && strings.TrimSpace(mode) != "" {
			out.Mode = strings.TrimSpace(mode)
		}
		if phase, ok := v["phase"].(string); ok && strings.TrimSpace(phase) != "" {
			out.Phase = strings.TrimSpace(phase)
		}
		if reason, ok := v["reasoning"].(string); ok && strings.TrimSpace(reason) != "" {
			out.Reasoning = strings.TrimSpace(reason)
		}
	}
	return out
}

func metadataFromTurn(t *turns.Turn) gepevents.EventMetadata {
	md := gepevents.EventMetadata{ID: uuid.New()}
	if t == nil {
		return md
	}
	md.TurnID = t.ID
	if sid, ok, err := turns.KeyTurnMetaSessionID.Get(t.Metadata); err == nil && ok {
		md.SessionID = sid
	}
	if inf, ok, err := turns.KeyTurnMetaInferenceID.Get(t.Metadata); err == nil && ok {
		md.InferenceID = inf
	}
	return md
}

// NewMiddleware emits custom thinking-mode events around inference.
func NewMiddleware(cfg Config) rootmw.Middleware {
	return func(next rootmw.HandlerFunc) rootmw.HandlerFunc {
		return func(ctx context.Context, t *turns.Turn) (*turns.Turn, error) {
			if t == nil {
				return next(ctx, t)
			}

			meta := metadataFromTurn(t)
			payload := &Payload{
				Mode:      cfg.Mode,
				Phase:     cfg.Phase,
				Reasoning: cfg.Reasoning,
			}

			gepevents.PublishEventToContext(ctx, NewThinkingModeStarted(meta, t.ID, payload))
			out, err := next(ctx, t)
			if err != nil {
				gepevents.PublishEventToContext(ctx, NewThinkingModeCompleted(meta, t.ID, payload, false, err.Error()))
				return out, err
			}
			gepevents.PublishEventToContext(ctx, NewThinkingModeCompleted(meta, t.ID, payload, true, ""))
			return out, nil
		}
	}
}
