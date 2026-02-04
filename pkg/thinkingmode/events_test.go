package thinkingmode

import (
	"encoding/json"
	"testing"

	gepevents "github.com/go-go-golems/geppetto/pkg/events"
	"github.com/google/uuid"
)

func TestThinkingModeEventRoundTrip(t *testing.T) {
	meta := gepevents.EventMetadata{ID: uuid.New(), SessionID: "run-1", TurnID: "turn-1"}
	orig := NewThinkingModeStarted(meta, "item-1", &Payload{Mode: "fast", Phase: "start", Reasoning: "ok"})

	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	ev, err := gepevents.NewEventFromJson(b)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	typed, ok := ev.(*EventThinkingModeStarted)
	if !ok {
		t.Fatalf("expected EventThinkingModeStarted, got %T", ev)
	}
	if typed.Type() != EventThinkingStarted {
		t.Fatalf("unexpected type: %s", typed.Type())
	}
	if typed.ItemID != "item-1" {
		t.Fatalf("unexpected item id: %s", typed.ItemID)
	}
	if typed.Data == nil || typed.Data.Mode != "fast" {
		t.Fatalf("unexpected payload: %#v", typed.Data)
	}
}
