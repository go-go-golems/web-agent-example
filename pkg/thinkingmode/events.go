package thinkingmode

import (
	gepevents "github.com/go-go-golems/geppetto/pkg/events"
)

const (
	EventThinkingStarted   = gepevents.EventType("webagent.thinking.started")
	EventThinkingUpdated   = gepevents.EventType("webagent.thinking.update")
	EventThinkingCompleted = gepevents.EventType("webagent.thinking.completed")
)

// Payload carries custom thinking-mode data for the web-agent-example pipeline.
type Payload struct {
	Mode      string         `json:"mode" yaml:"mode"`
	Phase     string         `json:"phase" yaml:"phase"`
	Reasoning string         `json:"reasoning" yaml:"reasoning"`
	ExtraData map[string]any `json:"extra_data,omitempty" yaml:"extra_data,omitempty"`
}

type EventThinkingModeStarted struct {
	gepevents.EventImpl
	ItemID string   `json:"item_id"`
	Data   *Payload `json:"data,omitempty"`
}

func NewThinkingModeStarted(metadata gepevents.EventMetadata, itemID string, data *Payload) *EventThinkingModeStarted {
	return &EventThinkingModeStarted{
		EventImpl: gepevents.EventImpl{Type_: EventThinkingStarted, Metadata_: metadata},
		ItemID:    itemID,
		Data:      data,
	}
}

type EventThinkingModeUpdate struct {
	gepevents.EventImpl
	ItemID string   `json:"item_id"`
	Data   *Payload `json:"data,omitempty"`
}

func NewThinkingModeUpdate(metadata gepevents.EventMetadata, itemID string, data *Payload) *EventThinkingModeUpdate {
	return &EventThinkingModeUpdate{
		EventImpl: gepevents.EventImpl{Type_: EventThinkingUpdated, Metadata_: metadata},
		ItemID:    itemID,
		Data:      data,
	}
}

type EventThinkingModeCompleted struct {
	gepevents.EventImpl
	ItemID  string   `json:"item_id"`
	Data    *Payload `json:"data,omitempty"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

func NewThinkingModeCompleted(metadata gepevents.EventMetadata, itemID string, data *Payload, success bool, errStr string) *EventThinkingModeCompleted {
	return &EventThinkingModeCompleted{
		EventImpl: gepevents.EventImpl{Type_: EventThinkingCompleted, Metadata_: metadata},
		ItemID:    itemID,
		Data:      data,
		Success:   success,
		Error:     errStr,
	}
}

func init() {
	_ = gepevents.RegisterEventFactory(string(EventThinkingStarted), func() gepevents.Event {
		return &EventThinkingModeStarted{EventImpl: gepevents.EventImpl{Type_: EventThinkingStarted}}
	})
	_ = gepevents.RegisterEventFactory(string(EventThinkingUpdated), func() gepevents.Event {
		return &EventThinkingModeUpdate{EventImpl: gepevents.EventImpl{Type_: EventThinkingUpdated}}
	})
	_ = gepevents.RegisterEventFactory(string(EventThinkingCompleted), func() gepevents.Event {
		return &EventThinkingModeCompleted{EventImpl: gepevents.EventImpl{Type_: EventThinkingCompleted}}
	})
}
