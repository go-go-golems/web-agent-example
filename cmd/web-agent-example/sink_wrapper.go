package main

import (
	"github.com/go-go-golems/geppetto/pkg/events"
	"github.com/go-go-golems/geppetto/pkg/events/structuredsink"
	infruntime "github.com/go-go-golems/pinocchio/pkg/inference/runtime"
	"github.com/go-go-golems/pinocchio/pkg/webchat"

	"github.com/go-go-golems/web-agent-example/pkg/discodialogue"
)

func discoSinkWrapper() webchat.EventSinkWrapper {
	return func(_ string, _ infruntime.ConversationRuntimeRequest, sink events.EventSink) (events.EventSink, error) {
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
