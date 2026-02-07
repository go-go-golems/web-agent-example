import React from 'react';
import { useParams } from 'react-router-dom';
import { useAppSelector } from '../store/hooks';
import { useGetConversationQuery, useGetTurnsQuery, useGetTurnDetailQuery } from '../api/debugApi';
import { TimelineLanes } from '../components/TimelineLanes';
import { TurnInspector } from '../components/TurnInspector';

export function OverviewPage() {
  const { sessionId } = useParams();
  const selectedConvId = useAppSelector((state) => state.ui.selectedConvId);
  const selectedTurnId = useAppSelector((state) => state.ui.selectedTurnId);

  const { data: conversation, isLoading: convLoading } = useGetConversationQuery(
    selectedConvId ?? '',
    { skip: !selectedConvId }
  );

  const { data: turns, isLoading: turnsLoading } = useGetTurnsQuery(
    { convId: selectedConvId ?? '', sessionId },
    { skip: !selectedConvId }
  );

  const { data: turnDetail, isLoading: turnDetailLoading } = useGetTurnDetailQuery(
    { 
      convId: selectedConvId ?? '', 
      sessionId: sessionId ?? conversation?.session_id ?? '', 
      turnId: selectedTurnId ?? '' 
    },
    { skip: !selectedConvId || !selectedTurnId }
  );

  if (!selectedConvId) {
    return (
      <div className="empty-state">
        <h2>👈 Select a conversation</h2>
        <p>Choose a conversation from the sidebar to view its details.</p>
      </div>
    );
  }

  if (convLoading || turnsLoading) {
    return (
      <div className="loading-state">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="overview-page">
      {/* Conversation header */}
      <div className="conv-header">
        <h2>Conversation {selectedConvId.slice(0, 8)}</h2>
        <div className="conv-meta">
          <span>Session: {sessionId || conversation?.session_id || '—'}</span>
          <span>Turns: {turns?.length ?? 0}</span>
        </div>
      </div>

      {/* Timeline view */}
      <div className="timeline-section">
        <h3>Timeline</h3>
        <TimelineLanes
          turns={turns ?? []}
          events={[]}
          entities={[]}
          isLive={false}
        />
      </div>

      {/* Turn inspector (if turn selected) */}
      {selectedTurnId && turnDetail && !turnDetailLoading && (
        <div className="turn-section">
          <h3>Turn Detail</h3>
          <TurnInspector turnDetail={turnDetail} />
        </div>
      )}

      <style>{`
        .overview-page {
          display: flex;
          flex-direction: column;
          gap: 24px;
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

        .empty-state h2 {
          font-size: 24px;
          margin-bottom: 8px;
        }

        .conv-header {
          padding-bottom: 16px;
          border-bottom: 1px solid var(--border-color);
        }

        .conv-header h2 {
          margin: 0 0 8px 0;
          font-size: 20px;
        }

        .conv-meta {
          display: flex;
          gap: 16px;
          font-size: 13px;
          color: var(--text-secondary);
        }

        .timeline-section, .turn-section {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 8px;
          padding: 16px;
        }

        .timeline-section h3, .turn-section h3 {
          margin: 0 0 16px 0;
          font-size: 14px;
          color: var(--text-secondary);
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }
      `}</style>
    </div>
  );
}

export default OverviewPage;
