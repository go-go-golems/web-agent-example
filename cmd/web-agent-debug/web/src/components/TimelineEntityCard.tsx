import React from 'react';
import type { TimelineEntity } from '../types';

export interface TimelineEntityCardProps {
  entity: TimelineEntity;
  selected?: boolean;
  onClick?: () => void;
  compact?: boolean;
}

export function TimelineEntityCard({ entity, selected = false, onClick, compact = false }: TimelineEntityCardProps) {
  const { id, kind, created_at, version, props } = entity;
  const time = new Date(created_at).toLocaleTimeString();
  const kindInfo = getKindInfo(kind);

  return (
    <div
      className={`card ${selected ? 'selected' : ''}`}
      onClick={onClick}
      style={{
        cursor: onClick ? 'pointer' : 'default',
        padding: compact ? '8px' : '12px',
        borderLeft: `3px solid ${kindInfo.color}`,
      }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-1">
        <div className="flex items-center gap-2">
          <span className={`badge ${kindInfo.badgeClass}`}>
            {kindInfo.icon} {kind}
          </span>
          {version && (
            <span className="text-xs text-muted">v{version}</span>
          )}
        </div>
        <span className="text-xs text-muted">{time}</span>
      </div>

      {/* ID */}
      <div className="text-xs text-secondary mb-1" style={{ fontFamily: 'monospace' }}>
        {truncateId(id)}
      </div>

      {/* Content based on kind */}
      {!compact && renderEntityContent(kind, props)}
    </div>
  );
}

function getKindInfo(kind: string): { icon: string; color: string; badgeClass: string } {
  switch (kind) {
    case 'message':
      return { icon: '💬', color: 'var(--accent-blue)', badgeClass: 'badge-blue' };
    case 'tool_call':
      return { icon: '🔧', color: 'var(--accent-yellow)', badgeClass: 'badge-yellow' };
    case 'tool_result':
      return { icon: '📤', color: 'var(--accent-cyan)', badgeClass: 'badge-cyan' };
    case 'thinking_mode':
      return { icon: '💭', color: 'var(--accent-purple)', badgeClass: 'badge-purple' };
    case 'planning':
      return { icon: '📋', color: 'var(--accent-green)', badgeClass: 'badge-green' };
    case 'log':
      return { icon: '📝', color: 'var(--text-muted)', badgeClass: 'badge-blue' };
    default:
      return { icon: '📦', color: 'var(--border-color)', badgeClass: 'badge-blue' };
  }
}

function truncateId(id: string): string {
  if (id.length <= 20) return id;
  return id.slice(0, 8) + '...' + id.slice(-8);
}

function renderEntityContent(kind: string, props: Record<string, unknown>): React.ReactNode {
  switch (kind) {
    case 'message': {
      const role = props.role as string;
      const content = props.content as string;
      const streaming = props.streaming as boolean;
      return (
        <div className="mt-2">
          <div className="flex items-center gap-2 mb-1">
            <span className={`badge ${role === 'user' ? 'badge-blue' : 'badge-green'}`}>
              {role}
            </span>
            {streaming && <span className="badge badge-yellow">streaming</span>}
          </div>
          <div className="text-sm" style={{ maxHeight: '60px', overflow: 'hidden' }}>
            {truncateText(content, 100)}
          </div>
        </div>
      );
    }
    case 'tool_call': {
      const name = props.name as string;
      const done = props.done as boolean;
      return (
        <div className="mt-2 flex items-center gap-2">
          <span className="badge badge-yellow">{name}</span>
          {done ? (
            <span className="badge badge-green">done</span>
          ) : (
            <span className="badge badge-blue">pending</span>
          )}
        </div>
      );
    }
    case 'tool_result': {
      const result = props.result;
      return (
        <pre className="mt-2 text-xs" style={{ maxHeight: '60px', overflow: 'hidden' }}>
          {truncateText(JSON.stringify(result, null, 2), 80)}
        </pre>
      );
    }
    default:
      if (Object.keys(props).length > 0) {
        return (
          <pre className="mt-2 text-xs" style={{ maxHeight: '40px', overflow: 'hidden' }}>
            {truncateText(JSON.stringify(props), 60)}
          </pre>
        );
      }
      return null;
  }
}

function truncateText(text: string | undefined, maxLength: number): string {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + '...';
}

export default TimelineEntityCard;
