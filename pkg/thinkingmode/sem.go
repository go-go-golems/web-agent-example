package thinkingmode

import (
	"encoding/json"

	semregistry "github.com/go-go-golems/pinocchio/pkg/sem/registry"
)

func init() {
	registerSemHandlers()
}

func registerSemHandlers() {
	semregistry.RegisterByType[*EventThinkingModeStarted](func(ev *EventThinkingModeStarted) ([][]byte, error) {
		m := &semThinkingModeStarted{ItemID: ev.ItemID, Data: semPayloadFromEventData(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*EventThinkingModeUpdate](func(ev *EventThinkingModeUpdate) ([][]byte, error) {
		m := &semThinkingModeUpdate{ItemID: ev.ItemID, Data: semPayloadFromEventData(ev.Data)}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingUpdated), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*EventThinkingModeCompleted](func(ev *EventThinkingModeCompleted) ([][]byte, error) {
		m := &semThinkingModeCompleted{
			ItemID:  ev.ItemID,
			Data:    semPayloadFromEventData(ev.Data),
			Success: ev.Success,
			Error:   ev.Error,
		}
		raw, err := marshalJSONRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingCompleted), "id": ev.ItemID, "data": raw})}, nil
	})
}

func semPayloadFromEventData(p *Payload) *semThinkingModePayload {
	if p == nil {
		return nil
	}
	return &semThinkingModePayload{
		Mode:      p.Mode,
		Phase:     p.Phase,
		Reasoning: p.Reasoning,
		ExtraData: cloneExtraData(p.ExtraData),
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

func cloneExtraData(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

type semThinkingModePayload struct {
	Mode      string         `json:"mode,omitempty"`
	Phase     string         `json:"phase,omitempty"`
	Reasoning string         `json:"reasoning,omitempty"`
	ExtraData map[string]any `json:"extraData,omitempty"`
}

type semThinkingModeStarted struct {
	ItemID string                  `json:"itemId,omitempty"`
	Data   *semThinkingModePayload `json:"data,omitempty"`
}

type semThinkingModeUpdate struct {
	ItemID string                  `json:"itemId,omitempty"`
	Data   *semThinkingModePayload `json:"data,omitempty"`
}

type semThinkingModeCompleted struct {
	ItemID  string                  `json:"itemId,omitempty"`
	Data    *semThinkingModePayload `json:"data,omitempty"`
	Success bool                    `json:"success"`
	Error   string                  `json:"error,omitempty"`
}
