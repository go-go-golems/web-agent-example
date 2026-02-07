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
      <div className="empty-state">
        <h2>📊 Timeline View</h2>
        <p>Select a conversation to view its timeline.</p>
      </div>
    );
  }

  const isLoading = timelineLoading || eventsLoading || turnsLoading;

  if (isLoading) {
    return (
      <div className="loading-state">
        <p>Loading timeline...</p>
      </div>
    );
  }

  return (
    <div className="timeline-page">
      <div className="page-header">
        <h2>📊 Timeline</h2>
        <div className="header-meta">
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

      <style>{`
        .timeline-page {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .empty-state, .loading-state {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          height: 100%;
          text-align: center;
          color: var(--text-muted);
        }

        .page-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding-bottom: 16px;
          border-bottom: 1px solid var(--border-color);
        }

        .page-header h2 {
          margin: 0;
          font-size: 18px;
        }

        .header-meta {
          display: flex;
          gap: 16px;
          font-size: 13px;
          color: var(--text-secondary);
        }
      `}</style>
    </div>
  );
}

export default TimelinePage;
