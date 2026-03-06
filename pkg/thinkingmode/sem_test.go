package thinkingmode

import (
	"encoding/json"
	"testing"

	thinkingmodepb "github.com/go-go-golems/pinocchio/cmd/web-chat/thinkingmode/pb"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestThinkingModeSemProtoRoundTrip(t *testing.T) {
	payload, err := payloadToProto(&Payload{
		Mode:      "deep",
		Phase:     "reasoning",
		Reasoning: "careful pass",
		ExtraData: map[string]any{"source": "unit-test"},
	})
	if err != nil {
		t.Fatalf("payloadToProto: %v", err)
	}
	if payload.Mode != "deep" || payload.Phase != "reasoning" || payload.Reasoning != "careful pass" {
		t.Fatalf("unexpected payload fields: %#v", payload)
	}
	if got := payload.ExtraData.GetFields()["source"].GetStringValue(); got != "unit-test" {
		t.Fatalf("unexpected extra_data.source: %q", got)
	}

	msg := &thinkingmodepb.ThinkingModeCompleted{
		ItemId:  "item-1",
		Data:    payload,
		Success: true,
	}
	raw, err := protoToRaw(msg)
	if err != nil {
		t.Fatalf("protoToRaw: %v", err)
	}

	var decoded thinkingmodepb.ThinkingModeCompleted
	if err := protojson.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("protojson.Unmarshal: %v", err)
	}
	if decoded.ItemId != "item-1" || decoded.Data.GetMode() != "deep" || !decoded.Success {
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
