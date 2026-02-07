import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAppSelector } from '../store/hooks';
import { useGetTurnDetailQuery } from '../api/debugApi';
import { TurnInspector } from '../components/TurnInspector';

export function TurnDetailPage() {
  const { sessionId, turnId } = useParams();
  const navigate = useNavigate();
  const selectedConvId = useAppSelector((state) => state.ui.selectedConvId);

  const { data: turnDetail, isLoading, error } = useGetTurnDetailQuery(
    { 
      convId: selectedConvId ?? '', 
      sessionId: sessionId ?? '', 
      turnId: turnId ?? '' 
    },
    { skip: !selectedConvId || !sessionId || !turnId }
  );

  if (!selectedConvId || !sessionId || !turnId) {
    return (
      <div className="empty-state">
        <h2>Turn Not Found</h2>
        <p>Missing conversation, session, or turn ID.</p>
        <button className="btn" onClick={() => navigate('/')}>
          Go Back
        </button>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="loading-state">
        <p>Loading turn...</p>
      </div>
    );
  }

  if (error || !turnDetail) {
    return (
      <div className="empty-state">
        <h2>Failed to load turn</h2>
        <p>Could not load turn details.</p>
        <button className="btn" onClick={() => navigate(-1)}>
          Go Back
        </button>
      </div>
    );
  }

  return (
    <div className="turn-detail-page">
      <div className="page-header">
        <button className="btn btn-ghost" onClick={() => navigate(-1)}>
          ← Back
        </button>
        <h2>Turn: {turnId}</h2>
      </div>

      <TurnInspector turnDetail={turnDetail} />

      <style>{`
        .turn-detail-page {
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
          gap: 16px;
        }

        .page-header {
          display: flex;
          align-items: center;
          gap: 16px;
          padding-bottom: 16px;
          border-bottom: 1px solid var(--border-color);
        }

        .page-header h2 {
          margin: 0;
          font-size: 18px;
          font-family: monospace;
        }
      `}</style>
    </div>
  );
}

export default TurnDetailPage;
