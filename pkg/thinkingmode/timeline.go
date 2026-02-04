package thinkingmode

import (
	"context"
	"strings"

	"github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/middleware"
	timelinepb "github.com/go-go-golems/pinocchio/pkg/sem/pb/proto/sem/timeline"
	"github.com/go-go-golems/pinocchio/pkg/webchat"
	"google.golang.org/protobuf/encoding/protojson"
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
		payload = &middleware.ThinkingModePayload{}
	}
	entity := &timelinepb.TimelineEntityV1{
		Id:   itemID,
		Kind: TimelineKind,
		Snapshot: &timelinepb.TimelineEntityV1_ThinkingMode{
			ThinkingMode: &timelinepb.ThinkingModeSnapshotV1{
				SchemaVersion: 1,
				Status:        status,
				Mode:          payload.Mode,
				Phase:         payload.Phase,
				Reasoning:     payload.Reasoning,
				Success:       success,
				Error:         errStr,
			},
		},
	}
	return p.Upsert(ctx, ev.Seq, entity)
}

func decodeThinkingEvent(ev webchat.TimelineSemEvent) (string, *middleware.ThinkingModePayload, string, bool, string) {
	switch ev.Type {
	case string(EventThinkingStarted):
		var pb middleware.ThinkingModeStarted
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "started", true, ""
	case string(EventThinkingUpdated):
		var pb middleware.ThinkingModeUpdate
		if err := protojson.Unmarshal(ev.Data, &pb); err != nil {
			return "", nil, "", false, ""
		}
		return pb.ItemId, pb.Data, "update", true, ""
	case string(EventThinkingCompleted):
		var pb middleware.ThinkingModeCompleted
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
