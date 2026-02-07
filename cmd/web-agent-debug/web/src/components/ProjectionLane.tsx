import React from 'react';
import type { TimelineEntity } from '../types';

export interface ProjectionLaneProps {
  entities: TimelineEntity[];
  selectedEntityId?: string;
  onEntitySelect?: (entity: TimelineEntity) => void;
}

export function ProjectionLane({ entities, selectedEntityId, onEntitySelect }: ProjectionLaneProps) {
  if (entities.length === 0) {
    return (
      <div className="empty-lane">
        <span className="text-muted">No timeline entities</span>
      </div>
    );
  }

  return (
    <div className="projection-lane-content">
      {entities.map((entity) => (
        <EntityCard
          key={entity.id}
          entity={entity}
          selected={entity.id === selectedEntityId}
          onClick={() => onEntitySelect?.(entity)}
        />
      ))}

      <style>{`
        .projection-lane-content {
          display: flex;
          flex-direction: column;
          gap: 6px;
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

interface EntityCardProps {
  entity: TimelineEntity;
  selected?: boolean;
  onClick?: () => void;
}

function EntityCard({ entity, selected, onClick }: EntityCardProps) {
  const { id, kind, created_at, version, props } = entity;
  const time = new Date(created_at).toLocaleTimeString();
  const kindInfo = getKindInfo(kind);

  return (
    <div
      className={`entity-card ${selected ? 'selected' : ''}`}
      onClick={onClick}
      style={{ borderLeftColor: kindInfo.color }}
    >
      <div className="entity-card-header">
        <span className="entity-kind" style={{ color: kindInfo.color }}>
          {kindInfo.icon} {kind}
        </span>
        <div className="entity-header-right">
          {version && version > 1 && <span className="entity-version">v{version}</span>}
          <span className="entity-time">{time}</span>
        </div>
      </div>

      <div className="entity-card-id">{id.slice(0, 16)}</div>

      {/* Render summary based on kind */}
      {kind === 'message' && (
        <div className="entity-summary">
          <span className={`role-chip role-${String(props.role)}`}>{String(props.role)}</span>
          {Boolean(props.streaming) && <span className="streaming-badge">streaming</span>}
        </div>
      )}

      {kind === 'tool_call' && (
        <div className="entity-summary">
          <span className="tool-name">{String(props.name)}</span>
          {props.done ? (
            <span className="done-badge">✓</span>
          ) : (
            <span className="pending-badge">⏳</span>
          )}
        </div>
      )}

      {kind === 'tool_result' && (
        <div className="entity-summary">
          <span className="result-preview">
            {truncate(JSON.stringify(props.result), 30)}
          </span>
        </div>
      )}

      <style>{`
        .entity-card {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-left: 3px solid var(--accent-green);
          border-radius: 4px;
          padding: 8px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .entity-card:hover {
          background: var(--bg-hover);
          border-color: var(--border-light);
        }

        .entity-card.selected {
          border-color: var(--accent-green);
          background: rgba(16, 185, 129, 0.1);
        }

        .entity-card-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 4px;
        }

        .entity-kind {
          font-size: 11px;
          font-weight: 600;
        }

        .entity-header-right {
          display: flex;
          align-items: center;
          gap: 6px;
        }

        .entity-version {
          font-size: 9px;
          padding: 1px 4px;
          background: var(--bg-tertiary);
          border-radius: 3px;
          color: var(--text-secondary);
        }

        .entity-time {
          font-size: 10px;
          color: var(--text-muted);
        }

        .entity-card-id {
          font-size: 10px;
          font-family: monospace;
          color: var(--text-muted);
          margin-bottom: 6px;
        }

        .entity-summary {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 11px;
        }

        .role-chip {
          padding: 1px 6px;
          border-radius: 3px;
          font-weight: 500;
        }

        .role-user { background: rgba(59, 130, 246, 0.2); color: var(--accent-blue); }
        .role-assistant { background: rgba(16, 185, 129, 0.2); color: var(--accent-green); }
        .role-system { background: rgba(139, 92, 246, 0.2); color: var(--accent-purple); }

        .streaming-badge {
          font-size: 9px;
          padding: 1px 4px;
          background: rgba(245, 158, 11, 0.2);
          color: var(--accent-yellow);
          border-radius: 3px;
        }

        .tool-name {
          color: var(--accent-yellow);
          font-weight: 500;
        }

        .done-badge {
          color: var(--accent-green);
        }

        .pending-badge {
          color: var(--accent-yellow);
        }

        .result-preview {
          color: var(--text-secondary);
          font-family: monospace;
          font-size: 10px;
        }
      `}</style>
    </div>
  );
}

function getKindInfo(kind: string): { icon: string; color: string } {
  switch (kind) {
    case 'message':
      return { icon: '💬', color: 'var(--accent-blue)' };
    case 'tool_call':
      return { icon: '🔧', color: 'var(--accent-yellow)' };
    case 'tool_result':
      return { icon: '📤', color: 'var(--accent-cyan)' };
    case 'thinking_mode':
      return { icon: '💭', color: 'var(--accent-purple)' };
    case 'planning':
      return { icon: '📋', color: 'var(--accent-green)' };
    default:
      return { icon: '📦', color: 'var(--border-color)' };
  }
}

function truncate(str: string, maxLen: number): string {
  if (str.length <= maxLen) return str;
  return str.slice(0, maxLen) + '...';
}

export default ProjectionLane;
