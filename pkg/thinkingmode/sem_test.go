package thinkingmode

import (
	"encoding/json"
	"testing"
)

func TestThinkingModeSemProtoRoundTrip(t *testing.T) {
	payload := semPayloadFromEventData(&Payload{
		Mode:      "deep",
		Phase:     "reasoning",
		Reasoning: "careful pass",
		ExtraData: map[string]any{"source": "unit-test"},
	})
	if payload.Mode != "deep" || payload.Phase != "reasoning" || payload.Reasoning != "careful pass" {
		t.Fatalf("unexpected payload fields: %#v", payload)
	}
	if got := payload.ExtraData["source"]; got != "unit-test" {
		t.Fatalf("unexpected extra_data.source: %q", got)
	}

	msg := &semThinkingModeCompleted{
		ItemID:  "item-1",
		Data:    payload,
		Success: true,
	}
	raw, err := marshalJSONRaw(msg)
	if err != nil {
		t.Fatalf("marshalJSONRaw: %v", err)
	}

	var decoded semThinkingModeCompleted
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if decoded.ItemID != "item-1" || decoded.Data.Mode != "deep" || !decoded.Success {
		t.Fatalf("unexpected decoded payload: %#v", decoded)
	}

	frame := wrapSem(map[string]any{
		"type": string(EventThinkingCompleted),
		"id":   "item-1",
		"data": raw,
	})
	var envelope map[string]any
	if err := json.Unmarshal(frame, &envelope); err != nil {
		t.Fatalf("frame json unmarshal: %v", err)
	}
	if sem, _ := envelope["sem"].(bool); !sem {
		t.Fatalf("expected sem envelope, got: %#v", envelope["sem"])
	}
}
