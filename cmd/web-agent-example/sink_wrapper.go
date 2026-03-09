package main

import (
	"strings"

	"github.com/go-go-golems/geppetto/pkg/events"
	"github.com/go-go-golems/geppetto/pkg/events/structuredsink"
	gepprofiles "github.com/go-go-golems/geppetto/pkg/profiles"
	infruntime "github.com/go-go-golems/pinocchio/pkg/inference/runtime"
	"github.com/go-go-golems/pinocchio/pkg/webchat"

	"github.com/go-go-golems/web-agent-example/pkg/discodialogue"
)

const discoMiddlewareName = "webagent-disco-dialogue"

func discoSinkWrapper() webchat.EventSinkWrapper {
	return func(convID string, req infruntime.ConversationRuntimeRequest, sink events.EventSink) (events.EventSink, error) {
		if req.ResolvedProfileRuntime == nil || !hasMiddleware(req.ResolvedProfileRuntime.Middlewares, discoMiddlewareName) {
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

func hasMiddleware(list []gepprofiles.MiddlewareUse, name string) bool {
	needle := strings.TrimSpace(name)
	if needle == "" {
		return false
	}
	for _, mw := range list {
		if mw.Enabled != nil && !*mw.Enabled {
			continue
		}
		if strings.TrimSpace(mw.Name) == needle {
			return true
		}
	}
	return false
}
