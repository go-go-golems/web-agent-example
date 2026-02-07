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
      <div className="overview-empty-state">
        <h2>👈 Select a conversation</h2>
        <p>Choose a conversation from the sidebar to view its details.</p>
      </div>
    );
  }

  if (convLoading || turnsLoading) {
    return (
      <div className="overview-loading-state">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="overview-page">
      {/* Conversation header */}
      <div className="overview-conv-header">
        <h2>Conversation {selectedConvId.slice(0, 8)}</h2>
        <div className="overview-conv-meta">
          <span>Session: {sessionId || conversation?.session_id || '—'}</span>
          <span>Turns: {turns?.length ?? 0}</span>
        </div>
      </div>

      {/* Timeline view */}
      <div className="overview-timeline-section">
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
        <div className="overview-turn-section">
          <h3>Turn Detail</h3>
          <TurnInspector turnDetail={turnDetail} />
        </div>
      )}
    </div>
  );
}

export default OverviewPage;
