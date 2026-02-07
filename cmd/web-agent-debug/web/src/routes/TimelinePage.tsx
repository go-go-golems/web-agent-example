import React from 'react';
import { useAppSelector } from '../store/hooks';
import { useGetTimelineQuery, useGetEventsQuery, useGetTurnsQuery } from '../api/debugApi';
import { TimelineLanes } from '../components/TimelineLanes';

export function TimelinePage() {
  const selectedConvId = useAppSelector((state) => state.ui.selectedConvId);

  const { data: timeline, isLoading: timelineLoading } = useGetTimelineQuery(
    { convId: selectedConvId ?? '' },
    { skip: !selectedConvId }
  );

  const { data: eventsData, isLoading: eventsLoading } = useGetEventsQuery(
    { convId: selectedConvId ?? '' },
    { skip: !selectedConvId }
  );

  const { data: turns, isLoading: turnsLoading } = useGetTurnsQuery(
    { convId: selectedConvId ?? '' },
    { skip: !selectedConvId }
  );

  if (!selectedConvId) {
    return (
      <div className="timeline-empty-state">
        <h2>📊 Timeline View</h2>
        <p>Select a conversation to view its timeline.</p>
      </div>
    );
  }

  const isLoading = timelineLoading || eventsLoading || turnsLoading;

  if (isLoading) {
    return (
      <div className="timeline-loading-state">
        <p>Loading timeline...</p>
      </div>
    );
  }

  return (
    <div className="timeline-page">
      <div className="timeline-page-header">
        <h2>📊 Timeline</h2>
        <div className="timeline-header-meta">
          <span>{turns?.length ?? 0} turns</span>
          <span>{eventsData?.events?.length ?? 0} events</span>
          <span>{timeline?.entities?.length ?? 0} entities</span>
        </div>
      </div>

      <TimelineLanes
        turns={turns ?? []}
        events={eventsData?.events ?? []}
        entities={timeline?.entities ?? []}
        isLive={false}
      />
    </div>
  );
}

export default TimelinePage;
