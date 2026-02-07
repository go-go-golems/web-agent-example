import React, { useState } from 'react';
import { useAppSelector } from '../store/hooks';
import { useGetEventsQuery } from '../api/debugApi';
import { EventCard } from '../components/EventCard';
import { EventInspector } from '../components/EventInspector';
import type { SemEvent } from '../types';

export function EventsPage() {
  const selectedConvId = useAppSelector((state) => state.ui.selectedConvId);
  const [selectedEvent, setSelectedEvent] = useState<SemEvent | null>(null);

  const { data: eventsData, isLoading } = useGetEventsQuery(
    { convId: selectedConvId ?? '' },
    { skip: !selectedConvId }
  );

  if (!selectedConvId) {
    return (
      <div className="empty-state">
        <h2>⚡ Events</h2>
        <p>Select a conversation to view its events.</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="loading-state">
        <p>Loading events...</p>
      </div>
    );
  }

  const events = eventsData?.events ?? [];

  return (
    <div className="events-page">
      <div className="page-header">
        <h2>⚡ Events</h2>
        <div className="header-meta">
          <span>{events.length} events</span>
          <span>Buffer: {eventsData?.buffer_capacity ?? 0}</span>
        </div>
      </div>

      <div className="events-layout">
        {/* Event list */}
        <div className="events-list">
          {events.length === 0 ? (
            <div className="empty-list">No events recorded</div>
          ) : (
            events.map((event) => (
              <EventCard
                key={event.id}
                event={event}
                onClick={() => setSelectedEvent(event)}
                selected={selectedEvent?.id === event.id}
              />
            ))
          )}
        </div>

        {/* Event detail */}
        {selectedEvent && (
          <div className="event-detail">
            <EventInspector event={selectedEvent} />
          </div>
        )}
      </div>

      <style>{`
        .events-page {
          display: flex;
          flex-direction: column;
          gap: 16px;
          height: 100%;
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
          flex-shrink: 0;
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

        .events-layout {
          display: flex;
          gap: 16px;
          flex: 1;
          overflow: hidden;
        }

        .events-list {
          width: 300px;
          overflow-y: auto;
          display: flex;
          flex-direction: column;
          gap: 8px;
          flex-shrink: 0;
        }

        .empty-list {
          text-align: center;
          padding: 24px;
          color: var(--text-muted);
        }

        .event-detail {
          flex: 1;
          overflow-y: auto;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 8px;
          padding: 16px;
        }
      `}</style>
    </div>
  );
}

export default EventsPage;
