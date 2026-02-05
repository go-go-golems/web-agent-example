package discodialogue

import (
	gepevents "github.com/go-go-golems/geppetto/pkg/events"
)

const (
	EventDialogueLineStarted   = gepevents.EventType("disco.dialogue.line.started")
	EventDialogueLineUpdate    = gepevents.EventType("disco.dialogue.line.update")
	EventDialogueLineCompleted = gepevents.EventType("disco.dialogue.line.completed")

	EventDialogueCheckStarted   = gepevents.EventType("disco.dialogue.check.started")
	EventDialogueCheckUpdate    = gepevents.EventType("disco.dialogue.check.update")
	EventDialogueCheckCompleted = gepevents.EventType("disco.dialogue.check.completed")

	EventDialogueStateStarted   = gepevents.EventType("disco.dialogue.state.started")
	EventDialogueStateUpdate    = gepevents.EventType("disco.dialogue.state.update")
	EventDialogueStateCompleted = gepevents.EventType("disco.dialogue.state.completed")
)

// DialogueLinePayload describes a single internal voice line.
type DialogueLinePayload struct {
	DialogueID string  `json:"dialogue_id" yaml:"dialogue_id"`
	LineID     string  `json:"line_id" yaml:"line_id"`
	Persona    string  `json:"persona" yaml:"persona"`
	Tone       string  `json:"tone" yaml:"tone"`
	Text       string  `json:"text" yaml:"text"`
	Trigger    string  `json:"trigger" yaml:"trigger"`
	Progress   float64 `json:"progress" yaml:"progress"`
	Status     string  `json:"status" yaml:"status"`
}

// DialogueCheckPayload describes a check/roll.
type DialogueCheckPayload struct {
	DialogueID string `json:"dialogue_id" yaml:"dialogue_id"`
	LineID     string `json:"line_id" yaml:"line_id"`
	CheckType  string `json:"check_type" yaml:"check_type"`
	Skill      string `json:"skill" yaml:"skill"`
	Difficulty int    `json:"difficulty" yaml:"difficulty"`
	Roll       int    `json:"roll" yaml:"roll"`
	Success    bool   `json:"success" yaml:"success"`
}

// DialogueStatePayload describes lifecycle status and summary.
type DialogueStatePayload struct {
	DialogueID string `json:"dialogue_id" yaml:"dialogue_id"`
	Status     string `json:"status" yaml:"status"`
	Summary    string `json:"summary" yaml:"summary"`
}

type DialogueLineStartedEvent struct {
	gepevents.EventImpl
	ItemID string               `json:"item_id"`
	Data   *DialogueLinePayload `json:"data,omitempty"`
}

type DialogueLineUpdateEvent struct {
	gepevents.EventImpl
	ItemID string               `json:"item_id"`
	Data   *DialogueLinePayload `json:"data,omitempty"`
}

type DialogueLineCompletedEvent struct {
	gepevents.EventImpl
	ItemID  string               `json:"item_id"`
	Data    *DialogueLinePayload `json:"data,omitempty"`
	Success bool                 `json:"success"`
	Error   string               `json:"error,omitempty"`
}

type DialogueCheckStartedEvent struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueCheckPayload `json:"data,omitempty"`
}

type DialogueCheckUpdateEvent struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueCheckPayload `json:"data,omitempty"`
}

type DialogueCheckCompletedEvent struct {
	gepevents.EventImpl
	ItemID  string                `json:"item_id"`
	Data    *DialogueCheckPayload `json:"data,omitempty"`
	Success bool                  `json:"success"`
	Error   string                `json:"error,omitempty"`
}

type DialogueStateStartedEvent struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueStatePayload `json:"data,omitempty"`
}

type DialogueStateUpdateEvent struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueStatePayload `json:"data,omitempty"`
}

type DialogueStateCompletedEvent struct {
	gepevents.EventImpl
	ItemID  string                `json:"item_id"`
	Data    *DialogueStatePayload `json:"data,omitempty"`
	Success bool                  `json:"success"`
	Error   string                `json:"error,omitempty"`
}

func NewDialogueLineStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload) *DialogueLineStartedEvent {
	return &DialogueLineStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueLineUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload) *DialogueLineUpdateEvent {
	return &DialogueLineUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueLineCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload, success bool, errStr string) *DialogueLineCompletedEvent {
	return &DialogueLineCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func NewDialogueCheckStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload) *DialogueCheckStartedEvent {
	return &DialogueCheckStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueCheckUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload) *DialogueCheckUpdateEvent {
	return &DialogueCheckUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueCheckCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload, success bool, errStr string) *DialogueCheckCompletedEvent {
	return &DialogueCheckCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func NewDialogueStateStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload) *DialogueStateStartedEvent {
	return &DialogueStateStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueStateUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload) *DialogueStateUpdateEvent {
	return &DialogueStateUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueStateCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload, success bool, errStr string) *DialogueStateCompletedEvent {
	return &DialogueStateCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func init() {
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineStarted), func() gepevents.Event {
		return &DialogueLineStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineUpdate), func() gepevents.Event {
		return &DialogueLineUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineCompleted), func() gepevents.Event {
		return &DialogueLineCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineCompleted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckStarted), func() gepevents.Event {
		return &DialogueCheckStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckUpdate), func() gepevents.Event {
		return &DialogueCheckUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckCompleted), func() gepevents.Event {
		return &DialogueCheckCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckCompleted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateStarted), func() gepevents.Event {
		return &DialogueStateStartedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateUpdate), func() gepevents.Event {
		return &DialogueStateUpdateEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateCompleted), func() gepevents.Event {
		return &DialogueStateCompletedEvent{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateCompleted}}
	})
}
