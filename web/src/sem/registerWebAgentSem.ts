import { registerSem, type SemEvent } from '@pwchat/sem/registry';
import { timelineSlice, type TimelineEntity } from '@pwchat/store/timelineSlice';
import type { AppDispatch } from '@pwchat/store/store';

const EVENT_STARTED = 'webagent.thinking.started';
const EVENT_UPDATED = 'webagent.thinking.update';
const EVENT_COMPLETED = 'webagent.thinking.completed';

function buildEntity(ev: SemEvent, status: string): TimelineEntity | null {
  const data = (ev.data ?? {}) as any;
  const payload = data.data ?? {};
  const id = (data.itemId ?? ev.id ?? '').toString();
  if (!id) return null;

  const success = status === 'completed' ? Boolean(data.success) : undefined;
  const error = data.error ? String(data.error) : '';

  return {
    id,
    kind: 'webagent_thinking_mode',
    createdAt: Date.now(),
    updatedAt: Date.now(),
    props: {
      status,
      mode: payload.mode ?? '',
      phase: payload.phase ?? '',
      reasoning: payload.reasoning ?? '',
      success,
      error,
    },
  };
}

function register(evType: string, status: string) {
  registerSem(evType, (ev: SemEvent, dispatch: AppDispatch) => {
    const entity = buildEntity(ev, status);
    if (!entity) return;
    dispatch(timelineSlice.actions.upsertEntity(entity));
  });
}

export function registerWebAgentSem() {
  register(EVENT_STARTED, 'started');
  register(EVENT_UPDATED, 'update');
  register(EVENT_COMPLETED, 'completed');
}
