package discodialogue

import (
	"context"
	"strings"

	semMw "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/middleware"
	timelinepb "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/timeline"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"google.golang.org/protobuf/encoding/protojson"
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
		payload = &semMw.DiscoDialogueLinePayload{}
	}

	entity := &timelinepb.TimelineEntityV1{
		Id:   itemID,
		Kind: TimelineKindLine,
		Snapshot: &timelinepb.TimelineEntityV1_DiscoDialogueLine{
			DiscoDialogueLine: &timelinepb.DiscoDialogueLineSnapshotV1{
				SchemaVersion: 1,
				Status:        status,
				DialogueId:    payload.DialogueId,
				LineId:        payload.LineId,
				Persona:       payload.Persona,
				Tone:          payload.Tone,
				Text:          payload.Text,
				Trigger:       payload.Trigger,
				Progress:      payload.Progress,
				Success:       success,
				Error:         errStr,
			},
		},
	}
	return p.Upsert(ctx, ev.Seq, entity)
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
		payload = &semMw.DiscoDialogueCheckPayload{}
	}

	entity := &timelinepb.TimelineEntityV1{
		Id:   itemID,
		Kind: TimelineKindCheck,
		Snapshot: &timelinepb.TimelineEntityV1_DiscoDialogueCheck{
			DiscoDialogueCheck: &timelinepb.DiscoDialogueCheckSnapshotV1{
				SchemaVersion: 1,
				Status:        status,
				DialogueId:    payload.DialogueId,
				LineId:        payload.LineId,
				CheckType:     payload.CheckType,
				Skill:         payload.Skill,
				Difficulty:    payload.Difficulty,
				Roll:          payload.Roll,
				Success:       success,
				Error:         errStr,
			},
		},
	}
	return p.Upsert(ctx, ev.Seq, entity)
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
		payload = &semMw.DiscoDialogueStatePayload{}
	}

	entity := &timelinepb.TimelineEntityV1{
		Id:   itemID,
		Kind: TimelineKindState,
		Snapshot: &timelinepb.TimelineEntityV1_DiscoDialogueState{
			DiscoDialogueState: &timelinepb.DiscoDialogueStateSnapshotV1{
				SchemaVersion: 1,
				Status:        status,
				DialogueId:    payload.DialogueId,
				Summary:       payload.Summary,
				Success:       success,
				Error:         errStr,
			},
		},
	}
	return p.Upsert(ctx, ev.Seq, entity)
}

func decodeLineEvent(ev webchat.TimelineSemEvent) (string, *semMw.DiscoDialogueLinePayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueLineStarted):
		var pb semMw.DiscoDialogueLineStarted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "started", true, ""
	case string(EventDialogueLineUpdate):
		var pb semMw.DiscoDialogueLineUpdate
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "update", true, ""
	case string(EventDialogueLineCompleted):
		var pb semMw.DiscoDialogueLineCompleted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemId, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}

func decodeCheckEvent(ev webchat.TimelineSemEvent) (string, *semMw.DiscoDialogueCheckPayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueCheckStarted):
		var pb semMw.DiscoDialogueCheckStarted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "started", true, ""
	case string(EventDialogueCheckUpdate):
		var pb semMw.DiscoDialogueCheckUpdate
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "update", true, ""
	case string(EventDialogueCheckCompleted):
		var pb semMw.DiscoDialogueCheckCompleted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemId, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}

func decodeStateEvent(ev webchat.TimelineSemEvent) (string, *semMw.DiscoDialogueStatePayload, string, bool, string) {
	switch ev.Type {
	case string(EventDialogueStateStarted):
		var pb semMw.DiscoDialogueStateStarted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "started", true, ""
	case string(EventDialogueStateUpdate):
		var pb semMw.DiscoDialogueStateUpdate
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "update", true, ""
	case string(EventDialogueStateCompleted):
		var pb semMw.DiscoDialogueStateCompleted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		status := "completed"
		if !pb.Success || strings.TrimSpace(pb.Error) != "" {
			status = "error"
		}
		return pb.ItemId, pb.Data, status, pb.Success, pb.Error
	}
	return "", nil, "", false, ""
}
