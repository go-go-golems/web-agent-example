import type { RenderEntity } from '@pwchat/webchat';

function badgeClass(status: string) {
  switch (status) {
    case 'completed':
      return 'bg-success';
    case 'error':
      return 'bg-danger';
    case 'update':
      return 'bg-info text-dark';
    default:
      return 'bg-secondary';
  }
}

export function WebAgentThinkingModeCard({ e }: { e: RenderEntity }) {
  const mode = String(e.props?.mode ?? '');
  const phase = String(e.props?.phase ?? '');
  const status = String(e.props?.status ?? '');
  const reasoning = e.props?.reasoning ? String(e.props.reasoning) : '';
  const error = e.props?.error ? String(e.props.error) : '';

  return (
    <div className="card webagent-thinking-card">
      <div className="card-header d-flex justify-content-between align-items-center">
        <div>
          <div className="fw-semibold">Thinking mode</div>
          <div className="text-muted small">{mode || 'unspecified'}</div>
        </div>
        <span className={`badge ${badgeClass(status)}`}>{status || 'pending'}</span>
      </div>
      <div className="card-body">
        <div className="d-flex gap-3 flex-wrap">
          <div>
            <div className="text-muted small">Phase</div>
            <div className="fw-semibold">{phase || 'n/a'}</div>
          </div>
          <div>
            <div className="text-muted small">Reasoning</div>
            <div className="webagent-thinking-reasoning">{reasoning || 'No reasoning provided.'}</div>
          </div>
        </div>
        {error ? <div className="alert alert-danger mt-3 mb-0">{error}</div> : null}
      </div>
    </div>
  );
}
