import React from 'react';
import type { TurnSnapshot } from '../types';
import { formatPhaseShort } from '../ui/format/phase';
import { formatTimeShort } from '../ui/format/time';
import { getBlockPresentation } from '../ui/presentation/blocks';

export interface StateTrackLaneProps {
  turns: TurnSnapshot[];
  selectedTurnId?: string;
  onTurnSelect?: (turn: TurnSnapshot) => void;
}

export function StateTrackLane({ turns, selectedTurnId, onTurnSelect }: StateTrackLaneProps) {
  if (turns.length === 0) {
    return (
      <div className="empty-lane">
        <span className="text-muted">No turn snapshots</span>
      </div>
    );
  }

  return (
    <div className="state-track-lane">
      {turns.map((turn) => (
        <TurnCard
          key={`${turn.turn_id}-${turn.phase}`}
          turn={turn}
          selected={turn.turn_id === selectedTurnId}
          onClick={() => onTurnSelect?.(turn)}
        />
      ))}

      <style>{`
        .state-track-lane {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .empty-lane {
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 24px;
          color: var(--text-muted);
        }
      `}</style>
    </div>
  );
}

interface TurnCardProps {
  turn: TurnSnapshot;
  selected?: boolean;
  onClick?: () => void;
}

function TurnCard({ turn, selected, onClick }: TurnCardProps) {
  const { turn_id, session_id, phase, turn: turnData, created_at_ms } = turn;
  const blockCount = turnData.blocks.length;
  const time = formatTimeShort(created_at_ms);

  // Get block kind summary
  const kindCounts = turnData.blocks.reduce((acc, block) => {
    acc[block.kind] = (acc[block.kind] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  return (
    <div
      className={`turn-card ${selected ? 'selected' : ''}`}
      onClick={onClick}
    >
      <div className="turn-card-header">
        <span className={`phase-badge phase-${phase}`}>{formatPhaseShort(phase)}</span>
        <span className="turn-time">{time}</span>
      </div>

      <div className="turn-card-id">
        <span className="turn-id" title={turn_id}>
          {turn_id.slice(0, 12)}
        </span>
        <span className="session-id" title={session_id}>
          sess:{session_id.slice(0, 8)}
        </span>
      </div>

      <div className="turn-card-blocks">
        {Object.entries(kindCounts).map(([kind, count]) => (
          <span key={kind} className={`block-chip block-kind-${kind}`}>
            {getBlockPresentation(kind).icon} {count}
          </span>
        ))}
      </div>

      <div className="turn-card-footer">
        <span className="block-count">{blockCount} blocks</span>
      </div>

      <style>{`
        .turn-card {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          padding: 10px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .turn-card:hover {
          background: var(--bg-hover);
          border-color: var(--border-light);
        }

        .turn-card.selected {
          border-color: var(--accent-purple);
          background: rgba(139, 92, 246, 0.1);
        }

        .turn-card-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 6px;
        }

        .phase-badge {
          font-size: 10px;
          font-weight: 600;
          text-transform: uppercase;
          padding: 2px 6px;
          border-radius: 3px;
        }

        .phase-pre_inference { background: rgba(245, 158, 11, 0.2); color: var(--accent-yellow); }
        .phase-post_inference { background: rgba(59, 130, 246, 0.2); color: var(--accent-blue); }
        .phase-post_tools { background: rgba(6, 182, 212, 0.2); color: var(--accent-cyan); }
        .phase-final { background: rgba(16, 185, 129, 0.2); color: var(--accent-green); }

        .turn-time {
          font-size: 10px;
          color: var(--text-muted);
        }

        .turn-card-id {
          display: flex;
          gap: 8px;
          font-size: 11px;
          font-family: monospace;
          color: var(--text-secondary);
          margin-bottom: 8px;
        }

        .turn-id {
          color: var(--accent-purple);
        }

        .session-id {
          color: var(--text-muted);
        }

        .turn-card-blocks {
          display: flex;
          flex-wrap: wrap;
          gap: 4px;
          margin-bottom: 6px;
        }

        .block-chip {
          font-size: 10px;
          padding: 2px 6px;
          border-radius: 3px;
          background: var(--bg-secondary);
        }

        .block-kind-system { border-left: 2px solid var(--block-system); }
        .block-kind-user { border-left: 2px solid var(--block-user); }
        .block-kind-llm_text { border-left: 2px solid var(--block-assistant); }
        .block-kind-tool_call { border-left: 2px solid var(--block-tool-call); }
        .block-kind-tool_use { border-left: 2px solid var(--block-tool-use); }
        .block-kind-reasoning { border-left: 2px solid var(--block-reasoning); }

        .turn-card-footer {
          font-size: 10px;
          color: var(--text-muted);
        }
      `}</style>
    </div>
  );
}

export default StateTrackLane;
