import type { RenderEntity } from '@pwchat/webchat';

function badgeClass(status: string) {
  switch (status) {
    case 'completed':
      return 'bg-success';
    case 'error':
      return 'bg-danger';
    case 'update':
      return 'bg-info text-dark';
    case 'started':
      return 'bg-primary';
    default:
      return 'bg-secondary';
  }
}

function titleForKind(kind: string) {
  switch (kind) {
    case 'disco_dialogue_line':
      return 'Internal Dialogue';
    case 'disco_dialogue_check':
      return 'Disco Check';
    case 'disco_dialogue_state':
      return 'Dialogue State';
    default:
      return 'Disco Dialogue';
  }
}

export function DiscoDialogueCard({ e }: { e: RenderEntity }) {
  const kind = String(e.kind ?? '');
  const props = (e.props ?? {}) as any;
  const status = String(props.status ?? '');
  const dialogueId = String(props.dialogueId ?? '');
  const lineId = String(props.lineId ?? '');
  const error = props.error ? String(props.error) : '';
  const success = typeof props.success === 'boolean' ? props.success : undefined;
  const progress = typeof props.progress === 'number' ? props.progress : Number(props.progress ?? NaN);
  const progressPct = Number.isFinite(progress) ? Math.max(0, Math.min(100, Math.round(progress * 100))) : null;

  const headerMeta = [dialogueId ? `dialogue:${dialogueId}` : '', lineId ? `line:${lineId}` : '']
    .filter(Boolean)
    .join(' • ');

  return (
    <div className={`card disco-dialogue-card disco-dialogue-${kind}`}>
      <div className="card-header d-flex justify-content-between align-items-center">
        <div>
          <div className="fw-semibold">{titleForKind(kind)}</div>
          {headerMeta ? <div className="text-muted small disco-dialogue-meta">{headerMeta}</div> : null}
        </div>
        <span className={`badge ${badgeClass(status)}`}>{status || 'pending'}</span>
      </div>
      <div className="card-body">
        {kind === 'disco_dialogue_line' ? (
          <>
            <div className="d-flex flex-wrap gap-3 align-items-start">
              <div>
                <div className="text-muted small">Persona</div>
                <div className="fw-semibold disco-dialogue-persona">{props.persona || 'Unknown'}</div>
              </div>
              <div>
                <div className="text-muted small">Tone</div>
                <div className="fw-semibold">{props.tone || 'n/a'}</div>
              </div>
              <div>
                <div className="text-muted small">Trigger</div>
                <div className="fw-semibold">{props.trigger || 'n/a'}</div>
              </div>
            </div>
            <div className="disco-dialogue-line-text mt-3">{props.text || 'No dialogue text yet.'}</div>
            {progressPct !== null ? (
              <div className="mt-3">
                <div className="text-muted small">Progress</div>
                <div className="progress disco-dialogue-progress">
                  <div className="progress-bar" role="progressbar" style={{ width: `${progressPct}%` }} />
                </div>
              </div>
            ) : null}
          </>
        ) : null}

        {kind === 'disco_dialogue_check' ? (
          <>
            <div className="d-flex flex-wrap gap-3">
              <div>
                <div className="text-muted small">Check</div>
                <div className="fw-semibold">{props.checkType || 'n/a'}</div>
              </div>
              <div>
                <div className="text-muted small">Skill</div>
                <div className="fw-semibold">{props.skill || 'n/a'}</div>
              </div>
              <div>
                <div className="text-muted small">Difficulty</div>
                <div className="fw-semibold">{Number.isFinite(props.difficulty) ? props.difficulty : props.difficulty ?? 'n/a'}</div>
              </div>
              <div>
                <div className="text-muted small">Roll</div>
                <div className="fw-semibold">{Number.isFinite(props.roll) ? props.roll : props.roll ?? 'n/a'}</div>
              </div>
            </div>
            {typeof success === 'boolean' ? (
              <div className={`mt-3 alert ${success ? 'alert-success' : 'alert-danger'} mb-0`}>
                {success ? 'Check succeeded' : 'Check failed'}
              </div>
            ) : null}
          </>
        ) : null}

        {kind === 'disco_dialogue_state' ? (
          <>
            <div className="text-muted small">Summary</div>
            <div className="disco-dialogue-summary">{props.summary || 'No summary yet.'}</div>
          </>
        ) : null}

        {error ? <div className="alert alert-danger mt-3 mb-0">{error}</div> : null}
      </div>
    </div>
  );
}
