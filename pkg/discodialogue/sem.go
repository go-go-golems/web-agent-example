package discodialogue

import (
	"encoding/json"

	semregistry "github.com/go-go-golems/pinocchio/pkg/sem/registry"
)

func init() {
	registerSemHandlers()
}

func registerSemHandlers() {
	semregistry.RegisterByType[*DialogueLineStartedEvent](func(ev *DialogueLineStartedEvent) ([][]byte, error) {
		m := &semDialogueLineStarted{ItemID: ev.ItemID, Data: linePayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueLineUpdateEvent](func(ev *DialogueLineUpdateEvent) ([][]byte, error) {
		m := &semDialogueLineUpdate{ItemID: ev.ItemID, Data: linePayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueLineCompletedEvent](func(ev *DialogueLineCompletedEvent) ([][]byte, error) {
		m := &semDialogueLineCompleted{ItemID: ev.ItemID, Data: linePayloadToSem(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineCompleted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckStartedEvent](func(ev *DialogueCheckStartedEvent) ([][]byte, error) {
		m := &semDialogueCheckStarted{ItemID: ev.ItemID, Data: checkPayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckUpdateEvent](func(ev *DialogueCheckUpdateEvent) ([][]byte, error) {
		m := &semDialogueCheckUpdate{ItemID: ev.ItemID, Data: checkPayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckCompletedEvent](func(ev *DialogueCheckCompletedEvent) ([][]byte, error) {
		m := &semDialogueCheckCompleted{ItemID: ev.ItemID, Data: checkPayloadToSem(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckCompleted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateStartedEvent](func(ev *DialogueStateStartedEvent) ([][]byte, error) {
		m := &semDialogueStateStarted{ItemID: ev.ItemID, Data: statePayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateUpdateEvent](func(ev *DialogueStateUpdateEvent) ([][]byte, error) {
		m := &semDialogueStateUpdate{ItemID: ev.ItemID, Data: statePayloadToSem(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateCompletedEvent](func(ev *DialogueStateCompletedEvent) ([][]byte, error) {
		m := &semDialogueStateCompleted{ItemID: ev.ItemID, Data: statePayloadToSem(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateCompleted), "id": ev.ItemID, "data": raw})}, nil
	})
}

func linePayloadToSem(p *DialogueLinePayload) *semDialogueLinePayload {
	if p == nil {
		return nil
	}
	return &semDialogueLinePayload{
		DialogueID: p.DialogueID,
		LineID:     p.LineID,
		Persona:    p.Persona,
		Tone:       p.Tone,
		Text:       p.Text,
		Trigger:    p.Trigger,
		Progress:   p.Progress,
		Status:     p.Status,
	}
}

func checkPayloadToSem(p *DialogueCheckPayload) *semDialogueCheckPayload {
	if p == nil {
		return nil
	}
	return &semDialogueCheckPayload{
		DialogueID: p.DialogueID,
		LineID:     p.LineID,
		CheckType:  p.CheckType,
		Skill:      p.Skill,
		Difficulty: p.Difficulty,
		Roll:       p.Roll,
		Success:    p.Success,
	}
}

func statePayloadToSem(p *DialogueStatePayload) *semDialogueStatePayload {
	if p == nil {
		return nil
	}
	return &semDialogueStatePayload{
		DialogueID: p.DialogueID,
		Status:     p.Status,
		Summary:    p.Summary,
	}
}

func wrapSem(ev map[string]any) []byte {
	b, _ := json.Marshal(map[string]any{"sem": true, "event": ev})
	return b
}

func marshalJSONRaw(v any) (json.RawMessage, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

type semDialogueLinePayload struct {
	DialogueID string  `json:"dialogueId,omitempty"`
	LineID     string  `json:"lineId,omitempty"`
	Persona    string  `json:"persona,omitempty"`
	Tone       string  `json:"tone,omitempty"`
	Text       string  `json:"text,omitempty"`
	Trigger    string  `json:"trigger,omitempty"`
	Progress   float64 `json:"progress,omitempty"`
	Status     string  `json:"status,omitempty"`
}

type semDialogueCheckPayload struct {
	DialogueID string `json:"dialogueId,omitempty"`
	LineID     string `json:"lineId,omitempty"`
	CheckType  string `json:"checkType,omitempty"`
	Skill      string `json:"skill,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
	Roll       int    `json:"roll,omitempty"`
	Success    bool   `json:"success"`
}

type semDialogueStatePayload struct {
	DialogueID string `json:"dialogueId,omitempty"`
	Status     string `json:"status,omitempty"`
	Summary    string `json:"summary,omitempty"`
}

type semDialogueLineStarted struct {
	ItemID string                  `json:"itemId,omitempty"`
	Data   *semDialogueLinePayload `json:"data,omitempty"`
}

type semDialogueLineUpdate struct {
	ItemID string                  `json:"itemId,omitempty"`
	Data   *semDialogueLinePayload `json:"data,omitempty"`
}

type semDialogueLineCompleted struct {
	ItemID  string                  `json:"itemId,omitempty"`
	Data    *semDialogueLinePayload `json:"data,omitempty"`
	Success bool                    `json:"success"`
	Error   string                  `json:"error,omitempty"`
}

type semDialogueCheckStarted struct {
	ItemID string                   `json:"itemId,omitempty"`
	Data   *semDialogueCheckPayload `json:"data,omitempty"`
}

type semDialogueCheckUpdate struct {
	ItemID string                   `json:"itemId,omitempty"`
	Data   *semDialogueCheckPayload `json:"data,omitempty"`
}

type semDialogueCheckCompleted struct {
	ItemID  string                   `json:"itemId,omitempty"`
	Data    *semDialogueCheckPayload `json:"data,omitempty"`
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
}

type semDialogueStateStarted struct {
	ItemID string                   `json:"itemId,omitempty"`
	Data   *semDialogueStatePayload `json:"data,omitempty"`
}

type semDialogueStateUpdate struct {
	ItemID string                   `json:"itemId,omitempty"`
	Data   *semDialogueStatePayload `json:"data,omitempty"`
}

type semDialogueStateCompleted struct {
	ItemID  string                   `json:"itemId,omitempty"`
	Data    *semDialogueStatePayload `json:"data,omitempty"`
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
}
