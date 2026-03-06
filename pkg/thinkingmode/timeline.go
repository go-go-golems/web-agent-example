package thinkingmode

import (
	"context"
	"strings"

	thinkingmodepb "github.com/go-go-golems/pinocchio/cmd/web-chat/thinkingmode/pb"
	timelinepb "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/timeline"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

const TimelineKind = "webagent_thinking_mode"

func init() {
	registerTimelineHandlers()
}

func registerTimelineHandlers() {
	webchat.RegisterTimelineHandler(string(EventThinkingStarted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertThinkingSnapshot(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventThinkingUpdated), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertThinkingSnapshot(ctx, p, ev)
	})
	webchat.RegisterTimelineHandler(string(EventThinkingCompleted), func(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent, _ int64) error {
		return upsertThinkingSnapshot(ctx, p, ev)
	})
}

func upsertThinkingSnapshot(ctx context.Context, p *webchat.TimelineProjector, ev webchat.TimelineSemEvent) error {
	itemID, payload, status, success, errStr := decodeThinkingEvent(ev)
	if strings.TrimSpace(itemID) == "" {
		itemID = ev.ID
	}
	if strings.TrimSpace(itemID) == "" {
		return nil
	}
	if payload == nil {
		payload = &thinkingmodepb.ThinkingModePayload{}
	}
	return p.Upsert(ctx, ev.Seq, timelineEntityFromMap(itemID, TimelineKind, map[string]any{
		"schemaVersion": 1,
		"status":        status,
		"mode":          payload.Mode,
		"phase":         payload.Phase,
		"reasoning":     payload.Reasoning,
		"success":       success,
		"error":         errStr,
	}))
}

func decodeThinkingEvent(ev webchat.TimelineSemEvent) (string, *thinkingmodepb.ThinkingModePayload, string, bool, string) {
	switch ev.Type {
	case string(EventThinkingStarted):
		var pb thinkingmodepb.ThinkingModeStarted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "started", true, ""
	case string(EventThinkingUpdated):
		var pb thinkingmodepb.ThinkingModeUpdate
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "update", true, ""
	case string(EventThinkingCompleted):
		var pb thinkingmodepb.ThinkingModeCompleted
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
