import React from 'react';
import type { SemEvent } from '../types';

export interface EventTrackLaneProps {
  events: SemEvent[];
  selectedSeq?: number;
  onEventSelect?: (event: SemEvent) => void;
}

export function EventTrackLane({ events, selectedSeq, onEventSelect }: EventTrackLaneProps) {
  if (events.length === 0) {
    return (
      <div className="empty-lane">
        <span className="text-muted">No events</span>
      </div>
    );
  }

  return (
    <div className="event-track-lane">
      {events.map((event) => (
        <EventDot
          key={event.seq}
          event={event}
          selected={event.seq === selectedSeq}
          onClick={() => onEventSelect?.(event)}
        />
      ))}

      <style>{`
        .event-track-lane {
          display: flex;
          flex-direction: column;
          gap: 4px;
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

interface EventDotProps {
  event: SemEvent;
  selected?: boolean;
  onClick?: () => void;
}

function EventDot({ event, selected, onClick }: EventDotProps) {
  const { type, id, stream_id, received_at } = event;
  const time = new Date(received_at).toLocaleTimeString();
  const typeInfo = getEventTypeInfo(type);

  return (
    <div
      className={`event-dot ${selected ? 'selected' : ''}`}
      onClick={onClick}
      style={{ borderLeftColor: typeInfo.color }}
    >
      <div className="event-dot-header">
        <span className="event-type" style={{ color: typeInfo.color }}>
          {typeInfo.icon} {type}
        </span>
        <span className="event-time">{time}</span>
      </div>

      <div className="event-dot-meta">
        <span className="event-id" title={id}>{id.slice(0, 12)}</span>
        {stream_id && <span className="event-stream">#{stream_id}</span>}
      </div>

      <style>{`
        .event-dot {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-left: 3px solid var(--accent-blue);
          border-radius: 4px;
          padding: 8px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .event-dot:hover {
          background: var(--bg-hover);
          border-color: var(--border-light);
        }

        .event-dot.selected {
          border-color: var(--accent-blue);
          background: rgba(59, 130, 246, 0.1);
        }

        .event-dot-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 4px;
        }

        .event-type {
          font-size: 11px;
          font-weight: 600;
        }

        .event-time {
          font-size: 10px;
          color: var(--text-muted);
        }

        .event-dot-meta {
          display: flex;
          gap: 8px;
          font-size: 10px;
          font-family: monospace;
          color: var(--text-muted);
        }

        .event-id {
          color: var(--text-secondary);
        }

        .event-stream {
          color: var(--accent-cyan);
        }
      `}</style>
    </div>
  );
}

function getEventTypeInfo(type: string): { icon: string; color: string } {
  if (type.startsWith('llm.')) {
    if (type === 'llm.start') return { icon: '▶️', color: 'var(--accent-green)' };
    if (type === 'llm.delta') return { icon: '📝', color: 'var(--accent-blue)' };
    if (type === 'llm.final') return { icon: '✅', color: 'var(--accent-green)' };
    if (type.includes('thinking')) return { icon: '💭', color: 'var(--accent-purple)' };
    return { icon: '🤖', color: 'var(--accent-blue)' };
  }
  if (type.startsWith('tool.')) {
    if (type === 'tool.start') return { icon: '🔧', color: 'var(--accent-yellow)' };
    if (type === 'tool.result') return { icon: '📤', color: 'var(--accent-cyan)' };
    if (type === 'tool.done') return { icon: '✓', color: 'var(--accent-green)' };
    return { icon: '🔧', color: 'var(--accent-yellow)' };
  }
  if (type === 'log') return { icon: '📋', color: 'var(--text-muted)' };
  return { icon: '📦', color: 'var(--border-color)' };
}

export default EventTrackLane;
