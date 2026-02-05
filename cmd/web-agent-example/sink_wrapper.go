package main

import (
	"strings"

	"github.com/go-go-golems/geppetto/pkg/events"
	"github.com/go-go-golems/geppetto/pkg/events/structuredsink"
	"github.com/go-go-golems/pinocchio/pkg/webchat"

	"github.com/go-go-golems/web-agent-example/pkg/discodialogue"
)

const discoMiddlewareName = "webagent-disco-dialogue"

func discoSinkWrapper() webchat.EventSinkWrapper {
	return func(convID string, cfg webchat.EngineConfig, sink events.EventSink) (events.EventSink, error) {
		if !hasMiddleware(cfg.Middlewares, discoMiddlewareName) {
			return sink, nil
		}
		extractors := []structuredsink.Extractor{
			discodialogue.NewDialogueLineExtractor(),
			discodialogue.NewDialogueCheckExtractor(),
			discodialogue.NewDialogueStateExtractor(),
		}
		wrapped := structuredsink.NewFilteringSink(
			sink,
			structuredsink.Options{Debug: false, Malformed: structuredsink.MalformedErrorEvents},
			extractors...,
		)
		return wrapped, nil
	}
}

func hasMiddleware(list []webchat.MiddlewareUse, name string) bool {
	needle := strings.TrimSpace(name)
	if needle == "" {
		return false
	}
	for _, mw := range list {
		if strings.TrimSpace(mw.Name) == needle {
			return true
		}
	}
	return false
}
