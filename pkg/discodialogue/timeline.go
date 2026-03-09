package discodialogue

import (
	"context"
	"encoding/json"
	"strings"

	timelinepb "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/timeline"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	TimelineKindLine  = "disco_dialogue_line"
	TimelineKindCheck = "disco_dialogue_check"
	TimelineKindState = "disco_dialogue_state"
)

func init() {
	registerTimelineHandlers()
}

func registerTimelineHandlers() {
	webchat.RegisterTimelineHandler(string(EventDialogueLineStarted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueLine(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueLineUpdate), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueLine(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueLineCompleted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueLine(ctx, p, ev)
	})

	webchat.RegisterTimelineHandler(string(EventDialogueCheckStarted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueCheck(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueCheckUpdate), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueCheck(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueCheckCompleted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueCheck(ctx, p, ev)
	})

	webchat.RegisterTimelineHandler(string(EventDialogueStateStarted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueState(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueStateUpdate), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueState(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventDialogueStateCompleted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertDialogueState(ctx, p, ev)
	})
}

func upsertDialogueLine(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent) error {
	itemID, payload, status, success, errStr := decodeLineEvent(ev)
	if strings.TrimSpace(itemID) == "" {
		itemID = ev.ID
	}
	if strings.TrimSpace(itemID) == "" {
		return nil
	}
	if payload == nil {
		payload = &semDialogueLinePayload{}
	}

	return p.Upsert(ctx, ev.Seq, timelineEntityFromMap(itemID, TimelineKindLine, map[string]any{
		"schemaVersion": 1,
		"status":        status,
		"dialogueId":    payload.DialogueID,
		"lineId":        payload.LineID,
		"persona":       payload.Persona,
		"tone":          payload.Tone,
		"text":          payload.Text,
		"trigger":       payload.Trigger,
		"progress":      payload.Progress,
		"success":       success,
		"error":         errStr,
	}))
}

func upsertDialogueCheck(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent) error {
	itemID, payload, status, success, errStr := decodeCheckEvent(ev)
	if strings.TrimSpace(itemID) == "" {
		itemID = ev.ID
	}
	if strings.TrimSpace(itemID) == "" {
		return nil
	}
	if payload == nil {
		payload = &semDialogueCheckPayload{}
	}

	return p.Upsert(ctx, ev.Seq, timelineEntityFromMap(itemID, TimelineKindCheck, map[string]any{
		"schemaVersion": 1,
		"status":        status,
		"dialogueId":    payload.DialogueID,
		"lineId":        payload.LineID,
		"checkType":     payload.CheckType,
		"skill":         payload.Skill,
		"difficulty":    payload.Difficulty,
		"roll":          payload.Roll,
		"success":       success,
		"error":         errStr,
	}))
}

func upsertDialogueState(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent) error {
	itemID, payload, status, success, errStr := decodeStateEvent(ev)
	if strings.TrimSpace(itemID) == "" {
		itemID = ev.ID
	}
	if strings.TrimSpace(itemID) == "" {
		return nil
	}
	if payload == nil {
		payload = &semDialogueStatePayload{}
	}

	return p.Upsert(ctx, ev.Seq, timelineEntityFromMap(itemID, TimelineKindState, map[string]any{
		"schemaVersion": 1,
		"status":        status,
		"dialogueId":    payload.DialogueID,
		"summary":       payload.Summary,
		"success":       success,
		"error":         errStr,
	}))
}

func decodeLineEvent(ev webchat.TimelineSemEvent) (string, *semDialogueLinePayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueLineStarted):
		var pb semDialogueLineStarted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "started", true, ""
	case string(EventDialogueLineUpdate):
		var pb semDialogueLineUpdate
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "update", true, ""
	case string(EventDialogueLineCompleted):
		var pb semDialogueLineCompleted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemID, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}

func decodeCheckEvent(ev webchat.TimelineSemEvent) (string, *semDialogueCheckPayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueCheckStarted):
		var pb semDialogueCheckStarted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "started", true, ""
	case string(EventDialogueCheckUpdate):
		var pb semDialogueCheckUpdate
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "update", true, ""
	case string(EventDialogueCheckCompleted):
		var pb semDialogueCheckCompleted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemID, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}

func decodeStateEvent(ev webchat.TimelineSemEvent) (string, *semDialogueStatePayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueStateStarted):
		var pb semDialogueStateStarted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "started", true, ""
	case string(EventDialogueStateUpdate):
		var pb semDialogueStateUpdate
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemID, pb.Data, "update", true, ""
	case string(EventDialogueStateCompleted):
		var pb semDialogueStateCompleted
		if err := json.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemID, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}

func timelineEntityFromMap(id, kind string, props map[string]any) *timelinepb.TimelineEntityV2 {
	st, err := structpb.NewStruct(props)
	if err != nil {
		st = &structpb.Struct{Fields: map[string]*structpb.Value{}}
	}
	return &timelinepb.TimelineEntityV2{
		Id:    strings.TrimSpace(id),
		Kind:  strings.TrimSpace(kind),
		Props: st,
	}
}
