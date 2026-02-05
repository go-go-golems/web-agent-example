package discodialogue

import (
	"encoding/json"

	semMw "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/middleware"
	semregistry "github.com/go-go-golems/pinocchio/pkg/sem/registry"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func init() {
	registerSemHandlers()
}

func registerSemHandlers() {
	semregistry.RegisterByType[*DialogueLineStartedEvent](func(ev *DialogueLineStartedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueLineStarted{ItemId: ev.ItemID, Data: linePayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueLineUpdateEvent](func(ev *DialogueLineUpdateEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueLineUpdate{ItemId: ev.ItemID, Data: linePayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueLineCompletedEvent](func(ev *DialogueLineCompletedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueLineCompleted{ItemId: ev.ItemID, Data: linePayloadToProto(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueLineCompleted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckStartedEvent](func(ev *DialogueCheckStartedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueCheckStarted{ItemId: ev.ItemID, Data: checkPayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckUpdateEvent](func(ev *DialogueCheckUpdateEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueCheckUpdate{ItemId: ev.ItemID, Data: checkPayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueCheckCompletedEvent](func(ev *DialogueCheckCompletedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueCheckCompleted{ItemId: ev.ItemID, Data: checkPayloadToProto(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueCheckCompleted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateStartedEvent](func(ev *DialogueStateStartedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueStateStarted{ItemId: ev.ItemID, Data: statePayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateStarted), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateUpdateEvent](func(ev *DialogueStateUpdateEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueStateUpdate{ItemId: ev.ItemID, Data: statePayloadToProto(ev.Data)}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateUpdate), "id": ev.ItemID, "data": raw})}, nil
	})

	semregistry.RegisterByType[*DialogueStateCompletedEvent](func(ev *DialogueStateCompletedEvent) ([][]byte, error) {
		m := &semMw.DiscoDialogueStateCompleted{ItemId: ev.ItemID, Data: statePayloadToProto(ev.Data), Success: ev.Success, Error: ev.Error}
		raw, err := protoToRaw(m)
		if err != nil {
			return nil, err
		}
		return [][]byte{wrapSem(map[string]any{"type": string(EventDialogueStateCompleted), "id": ev.ItemID, "data": raw})}, nil
	})
}

func linePayloadToProto(p *DialogueLinePayload) *semMw.DiscoDialogueLinePayload {
	if p == nil {
		return nil
	}
	return &semMw.DiscoDialogueLinePayload{
		DialogueId: p.DialogueID,
		LineId:     p.LineID,
		Persona:    p.Persona,
		Tone:       p.Tone,
		Text:       p.Text,
		Trigger:    p.Trigger,
		Progress:   p.Progress,
		Status:     p.Status,
	}
}

func checkPayloadToProto(p *DialogueCheckPayload) *semMw.DiscoDialogueCheckPayload {
	if p == nil {
		return nil
	}
	return &semMw.DiscoDialogueCheckPayload{
		DialogueId: p.DialogueID,
		LineId:     p.LineID,
		CheckType:  p.CheckType,
		Skill:      p.Skill,
		Difficulty: int32(p.Difficulty),
		Roll:       int32(p.Roll),
		Success:    p.Success,
	}
}

func statePayloadToProto(p *DialogueStatePayload) *semMw.DiscoDialogueStatePayload {
	if p == nil {
		return nil
	}
	return &semMw.DiscoDialogueStatePayload{
		DialogueId: p.DialogueID,
		Status:     p.Status,
		Summary:    p.Summary,
	}
}

func wrapSem(ev map[string]any) []byte {
	b, _ := json.Marshal(map[string]any{"sem": true, "event": ev})
	return b
}

func protoToRaw(m proto.Message) (json.RawMessage, error) {
	if m == nil {
		return nil, nil
	}
	b, err := protojson.MarshalOptions{EmitUnpopulated: false, UseProtoNames: false}.Marshal(m)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
