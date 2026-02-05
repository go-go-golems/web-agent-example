package discodialogue

import (
	"context"

	gepevents "github.com/go-go-golems/geppetto/pkg/events"
	"github.com/go-go-golems/geppetto/pkg/events/structuredsink"
	"github.com/go-go-golems/geppetto/pkg/events/structuredsink/parsehelpers"
	"github.com/rs/zerolog/log"
)

const (
	tagPackage = "disco"
	verV1      = "v1"
)

// DialogueLineExtractor parses <disco:dialogue_line:v1> payloads.
type DialogueLineExtractor struct{ pkg, typ, ver string }

func NewDialogueLineExtractor() *DialogueLineExtractor {
	return &DialogueLineExtractor{pkg: tagPackage, typ: "dialogue_line", ver: verV1}
}

func (e *DialogueLineExtractor) TagPackage() string { return e.pkg }
func (e *DialogueLineExtractor) TagType() string    { return e.typ }
func (e *DialogueLineExtractor) TagVersion() string { return e.ver }

func (e *DialogueLineExtractor) NewSession(ctx context.Context, meta gepevents.EventMetadata, itemID string) structuredsink.ExtractorSession {
	return &dialogueLineSession{meta: meta, itemID: itemID, ctx: ctx}
}

// DialogueCheckExtractor parses <disco:dialogue_check:v1> payloads.
type DialogueCheckExtractor struct{ pkg, typ, ver string }

func NewDialogueCheckExtractor() *DialogueCheckExtractor {
	return &DialogueCheckExtractor{pkg: tagPackage, typ: "dialogue_check", ver: verV1}
}

func (e *DialogueCheckExtractor) TagPackage() string { return e.pkg }
func (e *DialogueCheckExtractor) TagType() string    { return e.typ }
func (e *DialogueCheckExtractor) TagVersion() string { return e.ver }

func (e *DialogueCheckExtractor) NewSession(ctx context.Context, meta gepevents.EventMetadata, itemID string) structuredsink.ExtractorSession {
	return &dialogueCheckSession{meta: meta, itemID: itemID, ctx: ctx}
}

// DialogueStateExtractor parses <disco:dialogue_state:v1> payloads.
type DialogueStateExtractor struct{ pkg, typ, ver string }

func NewDialogueStateExtractor() *DialogueStateExtractor {
	return &DialogueStateExtractor{pkg: tagPackage, typ: "dialogue_state", ver: verV1}
}

func (e *DialogueStateExtractor) TagPackage() string { return e.pkg }
func (e *DialogueStateExtractor) TagType() string    { return e.typ }
func (e *DialogueStateExtractor) TagVersion() string { return e.ver }

func (e *DialogueStateExtractor) NewSession(ctx context.Context, meta gepevents.EventMetadata, itemID string) structuredsink.ExtractorSession {
	return &dialogueStateSession{meta: meta, itemID: itemID, ctx: ctx}
}

type dialogueLineSession struct {
	ctx         context.Context
	itemID      string
	meta        gepevents.EventMetadata
	ctrl        *parsehelpers.YAMLController[DialogueLinePayload]
	lastValid   *DialogueLinePayload
	lastValidOK bool
}

func (s *dialogueLineSession) OnStart(ctx context.Context) []gepevents.Event {
	s.ctrl = nil
	s.lastValid = nil
	s.lastValidOK = false

	log.Info().Str("extractor", "disco.dialogue_line").Str("item_id", s.itemID).Msg("extractor: start")
	return []gepevents.Event{NewDialogueLineStarted(s.meta, s.itemID, nil)}
}

func (s *dialogueLineSession) OnRaw(ctx context.Context, chunk []byte) []gepevents.Event {
	if s.ctrl == nil {
		s.ctrl = parsehelpers.NewDebouncedYAML[DialogueLinePayload](parsehelpers.DebounceConfig{
			SnapshotEveryBytes: 512,
			SnapshotOnNewline:  true,
			MaxBytes:           64 << 10,
		})
	}

	snap, err := s.ctrl.FeedBytes(chunk)
	if err == nil && snap != nil {
		s.lastValid = snap
		s.lastValidOK = true
	}

	if s.lastValidOK {
		return []gepevents.Event{NewDialogueLineUpdate(s.meta, s.itemID, s.lastValid)}
	}
	return nil
}

func (s *dialogueLineSession) OnCompleted(ctx context.Context, raw []byte, success bool, err error) []gepevents.Event {
	payload := s.lastValid
	errStr := ""
	if err != nil {
		errStr = err.Error()
		success = false
	} else if raw != nil {
		if s.ctrl == nil {
			s.ctrl = parsehelpers.NewDebouncedYAML[DialogueLinePayload](parsehelpers.DebounceConfig{})
		}
		snap, perr := s.ctrl.FinalBytes(raw)
		if perr == nil {
			if snap != nil {
				payload = snap
				s.lastValid = payload
				s.lastValidOK = true
			}
		} else {
			errStr = perr.Error()
			success = false
		}
	}

	return []gepevents.Event{NewDialogueLineCompleted(s.meta, s.itemID, payload, success, errStr)}
}

type dialogueCheckSession struct {
	ctx         context.Context
	itemID      string
	meta        gepevents.EventMetadata
	ctrl        *parsehelpers.YAMLController[DialogueCheckPayload]
	lastValid   *DialogueCheckPayload
	lastValidOK bool
}

func (s *dialogueCheckSession) OnStart(ctx context.Context) []gepevents.Event {
	s.ctrl = nil
	s.lastValid = nil
	s.lastValidOK = false

	log.Info().Str("extractor", "disco.dialogue_check").Str("item_id", s.itemID).Msg("extractor: start")
	return []gepevents.Event{NewDialogueCheckStarted(s.meta, s.itemID, nil)}
}

func (s *dialogueCheckSession) OnRaw(ctx context.Context, chunk []byte) []gepevents.Event {
	if s.ctrl == nil {
		s.ctrl = parsehelpers.NewDebouncedYAML[DialogueCheckPayload](parsehelpers.DebounceConfig{
			SnapshotEveryBytes: 512,
			SnapshotOnNewline:  true,
			MaxBytes:           32 << 10,
		})
	}

	snap, err := s.ctrl.FeedBytes(chunk)
	if err == nil && snap != nil {
		s.lastValid = snap
		s.lastValidOK = true
	}
	if s.lastValidOK {
		return []gepevents.Event{NewDialogueCheckUpdate(s.meta, s.itemID, s.lastValid)}
	}
	return nil
}

func (s *dialogueCheckSession) OnCompleted(ctx context.Context, raw []byte, success bool, err error) []gepevents.Event {
	payload := s.lastValid
	errStr := ""
	if err != nil {
		errStr = err.Error()
		success = false
	} else if raw != nil {
		if s.ctrl == nil {
			s.ctrl = parsehelpers.NewDebouncedYAML[DialogueCheckPayload](parsehelpers.DebounceConfig{})
		}
		snap, perr := s.ctrl.FinalBytes(raw)
		if perr == nil {
			if snap != nil {
				payload = snap
				s.lastValid = payload
				s.lastValidOK = true
			}
		} else {
			errStr = perr.Error()
			success = false
		}
	}

	return []gepevents.Event{NewDialogueCheckCompleted(s.meta, s.itemID, payload, success, errStr)}
}

type dialogueStateSession struct {
	ctx         context.Context
	itemID      string
	meta        gepevents.EventMetadata
	ctrl        *parsehelpers.YAMLController[DialogueStatePayload]
	lastValid   *DialogueStatePayload
	lastValidOK bool
}

func (s *dialogueStateSession) OnStart(ctx context.Context) []gepevents.Event {
	s.ctrl = nil
	s.lastValid = nil
	s.lastValidOK = false

	log.Info().Str("extractor", "disco.dialogue_state").Str("item_id", s.itemID).Msg("extractor: start")
	return []gepevents.Event{NewDialogueStateStarted(s.meta, s.itemID, nil)}
}

func (s *dialogueStateSession) OnRaw(ctx context.Context, chunk []byte) []gepevents.Event {
	if s.ctrl == nil {
		s.ctrl = parsehelpers.NewDebouncedYAML[DialogueStatePayload](parsehelpers.DebounceConfig{
			SnapshotEveryBytes: 512,
			SnapshotOnNewline:  true,
			MaxBytes:           32 << 10,
		})
	}

	snap, err := s.ctrl.FeedBytes(chunk)
	if err == nil && snap != nil {
		s.lastValid = snap
		s.lastValidOK = true
	}
	if s.lastValidOK {
		return []gepevents.Event{NewDialogueStateUpdate(s.meta, s.itemID, s.lastValid)}
	}
	return nil
}

func (s *dialogueStateSession) OnCompleted(ctx context.Context, raw []byte, success bool, err error) []gepevents.Event {
	payload := s.lastValid
	errStr := ""
	if err != nil {
		errStr = err.Error()
		success = false
	} else if raw != nil {
		if s.ctrl == nil {
			s.ctrl = parsehelpers.NewDebouncedYAML[DialogueStatePayload](parsehelpers.DebounceConfig{})
		}
		snap, perr := s.ctrl.FinalBytes(raw)
		if perr == nil {
			if snap != nil {
				payload = snap
				s.lastValid = payload
				s.lastValidOK = true
			}
		} else {
			errStr = perr.Error()
			success = false
		}
	}

	return []gepevents.Event{NewDialogueStateCompleted(s.meta, s.itemID, payload, success, errStr)}
}

var _ structuredsink.Extractor = (*DialogueLineExtractor)(nil)
var _ structuredsink.Extractor = (*DialogueCheckExtractor)(nil)
var _ structuredsink.Extractor = (*DialogueStateExtractor)(nil)
