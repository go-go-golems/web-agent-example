import { fromJson, type Message } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import {
  type ThinkingModeCompleted,
  ThinkingModeCompletedSchema,
  type ThinkingModeStarted,
  ThinkingModeStartedSchema,
  type ThinkingModeUpdate,
  ThinkingModeUpdateSchema,
} from '@pwchat/sem/pb/proto/sem/middleware/thinking_mode_pb';
import { registerSem, type SemEvent } from '@pwchat/sem/registry';
import { timelineSlice, type TimelineEntity } from '@pwchat/store/timelineSlice';
import type { AppDispatch } from '@pwchat/store/store';

const EVENT_STARTED = 'webagent.thinking.started';
const EVENT_UPDATED = 'webagent.thinking.update';
const EVENT_COMPLETED = 'webagent.thinking.completed';

function decodeProto<T extends Message>(schema: GenMessage<T>, raw: unknown): T | null {
  if (!raw || typeof raw !== 'object') return null;
  try {
    return fromJson(schema as any, raw as any, { ignoreUnknownFields: true }) as T;
  } catch {
    return null;
  }
}

function buildEntity(
  ev: SemEvent,
  status: string,
  decoded: ThinkingModeStarted | ThinkingModeUpdate | ThinkingModeCompleted | null,
): TimelineEntity | null {
  const payload = decoded?.data;
  const id = (decoded?.itemId ?? ev.id ?? '').toString();
  if (!id) return null;

  const completed = decoded as ThinkingModeCompleted | null;
  const success = status === 'completed' ? Boolean(completed?.success) : undefined;
  const error = completed?.error ? String(completed.error) : '';

  return {
    id,
    kind: 'webagent_thinking_mode',
    createdAt: Date.now(),
    updatedAt: Date.now(),
    props: {
      status,
      mode: payload?.mode ?? '',
      phase: payload?.phase ?? '',
      reasoning: payload?.reasoning ?? '',
      success,
      error,
    },
  };
}

function register(evType: string, status: string) {
  registerSem(evType, (ev: SemEvent, dispatch: AppDispatch) => {
    let decoded: ThinkingModeStarted | ThinkingModeUpdate | ThinkingModeCompleted | null = null;
    if (evType === EVENT_STARTED) decoded = decodeProto<ThinkingModeStarted>(ThinkingModeStartedSchema, ev.data);
    if (evType === EVENT_UPDATED) decoded = decodeProto<ThinkingModeUpdate>(ThinkingModeUpdateSchema, ev.data);
    if (evType === EVENT_COMPLETED) decoded = decodeProto<ThinkingModeCompleted>(ThinkingModeCompletedSchema, ev.data);
    const entity = buildEntity(ev, status, decoded);
    if (!entity) return;
    dispatch(timelineSlice.actions.upsertEntity(entity));
  });
}

export function registerWebAgentSem() {
  register(EVENT_STARTED, 'started');
  register(EVENT_UPDATED, 'update');
  register(EVENT_COMPLETED, 'completed');
}
