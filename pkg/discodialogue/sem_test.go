package discodialogue

import (
	"testing"

	thinkingmodepb "github.com/go-go-golems/pinocchio/cmd/web-chat/thinkingmode/pb"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestDiscoDialogueLineSemProtoRoundTrip(t *testing.T) {
	payload := linePayloadToProto(&DialogueLinePayload{
		DialogueID: "dlg-1",
		LineID:     "line-1",
		Persona:    "Logic",
		Tone:       "noir",
		Text:       "Check first principles.",
		Trigger:    "passive",
		Progress:   0.5,
		Status:     "update",
	})
	if payload.DialogueId != "dlg-1" || payload.LineId != "line-1" || payload.Persona != "Logic" {
		t.Fatalf("unexpected line payload: %#v", payload)
	}

	msg := &thinkingmodepb.DiscoDialogueLineCompleted{
		ItemId:  "line-item-1",
		Data:    payload,
		Success: true,
	}
	raw, err := protoToRaw(msg)
	if err != nil {
		t.Fatalf("protoToRaw: %v", err)
	}

	var decoded thinkingmodepb.DiscoDialogueLineCompleted
	if err := protojson.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("protojson.Unmarshal: %v", err)
	}
	if decoded.ItemId != "line-item-1" || decoded.Data.GetDialogueId() != "dlg-1" || decoded.Data.GetPersona() != "Logic" {
		t.Fatalf("unexpected decoded line payload: %#v", decoded)
	}
}

func TestDiscoDialogueCheckSemProtoRoundTrip(t *testing.T) {
	payload := checkPayloadToProto(&DialogueCheckPayload{
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

	msg := &thinkingmodepb.DiscoDialogueCheckCompleted{
		ItemId:  "check-item-1",
		Data:    payload,
		Success: false,
		Error:   "missed threshold",
	}
	raw, err := protoToRaw(msg)
	if err != nil {
		t.Fatalf("protoToRaw: %v", err)
	}

	var decoded thinkingmodepb.DiscoDialogueCheckCompleted
	if err := protojson.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("protojson.Unmarshal: %v", err)
	}
	if decoded.ItemId != "check-item-1" || decoded.Data.GetSkill() != "volition" || decoded.Success {
		t.Fatalf("unexpected decoded check payload: %#v", decoded)
	}
}
