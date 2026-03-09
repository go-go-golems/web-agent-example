package discodialogue

import (
	"testing"

	"encoding/json"
)

func TestDiscoDialogueLineSemProtoRoundTrip(t *testing.T) {
	payload := linePayloadToSem(&DialogueLinePayload{
		DialogueID: "dlg-1",
		LineID:     "line-1",
		Persona:    "Logic",
		Tone:       "noir",
		Text:       "Check first principles.",
		Trigger:    "passive",
		Progress:   0.5,
		Status:     "update",
	})
	if payload.DialogueID != "dlg-1" || payload.LineID != "line-1" || payload.Persona != "Logic" {
		t.Fatalf("unexpected line payload: %#v", payload)
	}

	msg := &semDialogueLineCompleted{
		ItemID:  "line-item-1",
		Data:    payload,
		Success: true,
	}
	raw, err := marshalJSONRaw(msg)
	if err != nil {
		t.Fatalf("marshalJSONRaw: %v", err)
	}

	var decoded semDialogueLineCompleted
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if decoded.ItemID != "line-item-1" || decoded.Data.DialogueID != "dlg-1" || decoded.Data.Persona != "Logic" {
		t.Fatalf("unexpected decoded line payload: %#v", decoded)
	}
}

func TestDiscoDialogueCheckSemProtoRoundTrip(t *testing.T) {
	payload := checkPayloadToSem(&DialogueCheckPayload{
		DialogueID: "dlg-1",
		LineID:     "line-2",
		CheckType:  "active",
		Skill:      "volition",
		Difficulty: 12,
		Roll:       10,
		Success:    false,
	})
	if payload.Skill != "volition" || payload.Difficulty != 12 || payload.Roll != 10 {
		t.Fatalf("unexpected check payload: %#v", payload)
	}

	msg := &semDialogueCheckCompleted{
		ItemID:  "check-item-1",
		Data:    payload,
		Success: false,
		Error:   "missed threshold",
	}
	raw, err := marshalJSONRaw(msg)
	if err != nil {
		t.Fatalf("marshalJSONRaw: %v", err)
	}

	var decoded semDialogueCheckCompleted
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if decoded.ItemID != "check-item-1" || decoded.Data.Skill != "volition" || decoded.Success {
		t.Fatalf("unexpected decoded check payload: %#v", decoded)
	}
}
