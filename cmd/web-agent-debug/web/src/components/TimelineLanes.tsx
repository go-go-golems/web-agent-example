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
    </div>
  );
}

export default TimelineLanes;
