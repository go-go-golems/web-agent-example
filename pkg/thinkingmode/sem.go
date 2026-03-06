package thinkingmode

import (
	"encoding/json"

	thinkingmodepb "github.com/go-go-golems/pinocchio/cmd/web-chat/thinkingmode/pb"
	semregistry "github.com/go-go-golems/pinocchio/pkg/sem/registry"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	registerSemHandlers()
}

func registerSemHandlers() {
	semregistry.RegisterByType[*EventThinkingModeStarted](func(ev *EventThinkingModeStarted) ([][]byte, error) {
		data, err := payloadToProto(ev.Data)
		if err != nil {
			return nil, err
		}
		m := &thinkingmodepb.ThinkingModeStarted{ItemId: ev.ItemID, Data: data}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*EventThinkingModeUpdate](func(ev *EventThinkingModeUpdate) ([][]byte, error) {
		data, err := payloadToProto(ev.Data)
		if err != nil {
			return nil, err
		}
		m := &thinkingmodepb.ThinkingModeUpdate{ItemId: ev.ItemID, Data: data}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingUpdated), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*EventThinkingModeCompleted](func(ev *EventThinkingModeCompleted) ([][]byte, error) {
		data, err := payloadToProto(ev.Data)
		if err != nil {
			return nil, err
		}
		m := &thinkingmodepb.ThinkingModeCompleted{ItemId: ev.ItemID, Data: data, Success: ev.Success, Error: ev.Error}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventThinkingCompleted), "id": ev.ItemID, "data": raw})}, nil
	})
}

func payloadToProto(p *Payload) (*thinkingmodepb.ThinkingModePayload, error) {
	if p == nil {
		return nil, nil
	}
	var extra *structpb.Struct
	if len(p.ExtraData) > 0 {
		st, err := structpb.NewStruct(p.ExtraData)
		if err != nil {
			return nil, err
		}
		extra = st
	}
	return &thinkingmodepb.ThinkingModePayload{
		Mode:      p.Mode,
		Phase:     p.Phase,
		Reasoning: p.Reasoning,
		ExtraData: extra,
	}, nil
}

func wrapSem(ev map[string]any) []byte {
	b, _ := json.Marshal(map[string]any{"sem": true, "event": ev})
	return b
}

func protoToRaw(m proto.Message) (json.RawMessage, error) {
	if m == nil {
		return nil, nil
	}
	b, err := protojson.MarshalOptions{
		EmitUnpopulated: false,
		UseProtoNames:   false,
	}.Marshal(m)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
