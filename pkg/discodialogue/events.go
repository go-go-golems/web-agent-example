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

type EventDialogueLineStarted struct {
	gepevents.EventImpl
	ItemID string               `json:"item_id"`
	Data   *DialogueLinePayload `json:"data,omitempty"`
}

type EventDialogueLineUpdate struct {
	gepevents.EventImpl
	ItemID string               `json:"item_id"`
	Data   *DialogueLinePayload `json:"data,omitempty"`
}

type EventDialogueLineCompleted struct {
	gepevents.EventImpl
	ItemID  string               `json:"item_id"`
	Data    *DialogueLinePayload `json:"data,omitempty"`
	Success bool                 `json:"success"`
	Error   string               `json:"error,omitempty"`
}

type EventDialogueCheckStarted struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueCheckPayload `json:"data,omitempty"`
}

type EventDialogueCheckUpdate struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueCheckPayload `json:"data,omitempty"`
}

type EventDialogueCheckCompleted struct {
	gepevents.EventImpl
	ItemID  string                `json:"item_id"`
	Data    *DialogueCheckPayload `json:"data,omitempty"`
	Success bool                  `json:"success"`
	Error   string                `json:"error,omitempty"`
}

type EventDialogueStateStarted struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueStatePayload `json:"data,omitempty"`
}

type EventDialogueStateUpdate struct {
	gepevents.EventImpl
	ItemID string                `json:"item_id"`
	Data   *DialogueStatePayload `json:"data,omitempty"`
}

type EventDialogueStateCompleted struct {
	gepevents.EventImpl
	ItemID  string                `json:"item_id"`
	Data    *DialogueStatePayload `json:"data,omitempty"`
	Success bool                  `json:"success"`
	Error   string                `json:"error,omitempty"`
}

func NewDialogueLineStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload) *EventDialogueLineStarted {
	return &EventDialogueLineStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueLineUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload) *EventDialogueLineUpdate {
	return &EventDialogueLineUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueLineCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueLinePayload, success bool, errStr string) *EventDialogueLineCompleted {
	return &EventDialogueLineCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func NewDialogueCheckStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload) *EventDialogueCheckStarted {
	return &EventDialogueCheckStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueCheckUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload) *EventDialogueCheckUpdate {
	return &EventDialogueCheckUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueCheckCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueCheckPayload, success bool, errStr string) *EventDialogueCheckCompleted {
	return &EventDialogueCheckCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func NewDialogueStateStarted(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload) *EventDialogueStateStarted {
	return &EventDialogueStateStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateStarted, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueStateUpdate(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload) *EventDialogueStateUpdate {
	return &EventDialogueStateUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateUpdate, Metadata_: metadata}, ItemID: itemID, Data: data}
}

func NewDialogueStateCompleted(metadata gepevents.EventMetadata, itemID string, data *DialogueStatePayload, success bool, errStr string) *EventDialogueStateCompleted {
	return &EventDialogueStateCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateCompleted, Metadata_: metadata}, ItemID: itemID, Data: data, Success: success, Error: errStr}
}

func init() {
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineStarted), func() gepevents.Event {
		return &EventDialogueLineStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineUpdate), func() gepevents.Event {
		return &EventDialogueLineUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueLineCompleted), func() gepevents.Event {
		return &EventDialogueLineCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueLineCompleted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckStarted), func() gepevents.Event {
		return &EventDialogueCheckStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckUpdate), func() gepevents.Event {
		return &EventDialogueCheckUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueCheckCompleted), func() gepevents.Event {
		return &EventDialogueCheckCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueCheckCompleted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateStarted), func() gepevents.Event {
		return &EventDialogueStateStarted{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateUpdate), func() gepevents.Event {
		return &EventDialogueStateUpdate{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateUpdate}}
	})
	_ = gepevents.RegisterEventFactory(string(EventDialogueStateCompleted), func() gepevents.Event {
		return &EventDialogueStateCompleted{EventImpl: gepevents.EventImpl{Type_: EventDialogueStateCompleted}}
	})
}
