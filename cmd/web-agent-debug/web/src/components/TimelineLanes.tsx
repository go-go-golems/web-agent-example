import React, { useRef, useEffect } from 'react';
import type { TurnSnapshot, SemEvent, TimelineEntity } from '../types';
import { StateTrackLane } from './StateTrackLane';
import { EventTrackLane } from './EventTrackLane';
import { ProjectionLane } from './ProjectionLane';
import { NowMarker } from './NowMarker';

export interface TimelineLanesProps {
  turns: TurnSnapshot[];
  events: SemEvent[];
  entities: TimelineEntity[];
  isLive?: boolean;
  onTurnSelect?: (turn: TurnSnapshot) => void;
  onEventSelect?: (event: SemEvent) => void;
  onEntitySelect?: (entity: TimelineEntity) => void;
  selectedTurnId?: string;
  selectedEventSeq?: number;
  selectedEntityId?: string;
}

export function TimelineLanes({
  turns,
  events,
  entities,
  isLive = false,
  onTurnSelect,
  onEventSelect,
  onEntitySelect,
  selectedTurnId,
  selectedEventSeq,
  selectedEntityId,
}: TimelineLanesProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const stateRef = useRef<HTMLDivElement>(null);
  const eventRef = useRef<HTMLDivElement>(null);
  const projRef = useRef<HTMLDivElement>(null);

  // Sync scroll across lanes
  useEffect(() => {
    const handleScroll = (source: HTMLDivElement | null) => {
      if (!source) return;
      const scrollTop = source.scrollTop;
      
      [stateRef, eventRef, projRef].forEach(ref => {
        if (ref.current && ref.current !== source) {
          ref.current.scrollTop = scrollTop;
        }
      });
    };

    const refs = [stateRef, eventRef, projRef];
    const handlers = refs.map(ref => {
      const handler = () => handleScroll(ref.current);
      ref.current?.addEventListener('scroll', handler);
      return { ref, handler };
    });

    return () => {
      handlers.forEach(({ ref, handler }) => {
        ref.current?.removeEventListener('scroll', handler);
      });
    };
  }, []);

  return (
    <div className="timeline-lanes" ref={containerRef}>
      {/* Header row */}
      <div className="timeline-header">
        <div className="lane-header">
          <span className="lane-title">📋 State Track</span>
          <span className="lane-count">{turns.length} turns</span>
        </div>
        <div className="lane-header">
          <span className="lane-title">⚡ Events</span>
          <span className="lane-count">{events.length} events</span>
        </div>
        <div className="lane-header">
          <span className="lane-title">🎯 Projection</span>
          <span className="lane-count">{entities.length} entities</span>
        </div>
      </div>

      {/* Lane columns */}
      <div className="timeline-body">
        <div className="lane state-lane" ref={stateRef}>
          <StateTrackLane
            turns={turns}
            selectedTurnId={selectedTurnId}
            onTurnSelect={onTurnSelect}
          />
          {isLive && <NowMarker />}
        </div>

        <div className="lane event-lane" ref={eventRef}>
          <EventTrackLane
            events={events}
            selectedSeq={selectedEventSeq}
            onEventSelect={onEventSelect}
          />
          {isLive && <NowMarker />}
        </div>

        <div className="lane projection-lane" ref={projRef}>
          <ProjectionLane
            entities={entities}
            selectedEntityId={selectedEntityId}
            onEntitySelect={onEntitySelect}
          />
          {isLive && <NowMarker />}
        </div>
      </div>

      <style>{`
        .timeline-lanes {
          display: flex;
          flex-direction: column;
          height: 100%;
          background: var(--bg-primary);
        }

        .timeline-header {
          display: flex;
          border-bottom: 1px solid var(--border-color);
          background: var(--bg-secondary);
        }

        .lane-header {
          flex: 1;
          padding: 8px 12px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          border-right: 1px solid var(--border-color);
        }

        .lane-header:last-child {
          border-right: none;
        }

        .lane-title {
          font-weight: 600;
          font-size: 13px;
        }

        .lane-count {
          font-size: 11px;
          color: var(--text-muted);
        }

        .timeline-body {
          display: flex;
          flex: 1;
          overflow: hidden;
        }

        .lane {
          flex: 1;
          overflow-y: auto;
          overflow-x: hidden;
          padding: 8px;
          border-right: 1px solid var(--border-color);
          position: relative;
        }

        .lane:last-child {
          border-right: none;
        }

        .state-lane {
          background: rgba(139, 92, 246, 0.03);
        }

        .event-lane {
          background: rgba(59, 130, 246, 0.03);
        }

        .projection-lane {
          background: rgba(16, 185, 129, 0.03);
        }
      `}</style>
    </div>
  );
}

export default TimelineLanes;
