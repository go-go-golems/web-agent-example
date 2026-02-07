import React, { useState } from 'react';

export type AnomalyType = 'orphan_event' | 'missing_correlation' | 'timing_outlier' | 'sequence_gap' | 'schema_error';
export type AnomalySeverity = 'error' | 'warning' | 'info';

export interface Anomaly {
  id: string;
  type: AnomalyType;
  severity: AnomalySeverity;
  message: string;
  details?: string;
  timestamp: string;
  relatedIds?: {
    eventId?: string;
    turnId?: string;
    blockIndex?: number;
  };
}

export interface AnomalyPanelProps {
  anomalies: Anomaly[];
  isOpen: boolean;
  onClose: () => void;
  onAnomalyClick?: (anomaly: Anomaly) => void;
}

export function AnomalyPanel({ anomalies, isOpen, onClose, onAnomalyClick }: AnomalyPanelProps) {
  const [selectedAnomaly, setSelectedAnomaly] = useState<Anomaly | null>(null);
  const [filterSeverity, setFilterSeverity] = useState<AnomalySeverity | 'all'>('all');

  const filteredAnomalies = filterSeverity === 'all'
    ? anomalies
    : anomalies.filter(a => a.severity === filterSeverity);

  const counts = {
    error: anomalies.filter(a => a.severity === 'error').length,
    warning: anomalies.filter(a => a.severity === 'warning').length,
    info: anomalies.filter(a => a.severity === 'info').length,
  };

  const handleAnomalyClick = (anomaly: Anomaly) => {
    setSelectedAnomaly(anomaly);
    onAnomalyClick?.(anomaly);
  };

  if (!isOpen) return null;

  return (
    <div className="anomaly-panel-overlay">
      <div className="anomaly-panel">
        {/* Header */}
        <div className="anomaly-header">
          <div className="anomaly-title">
            <h3>⚠️ Anomalies</h3>
            <span className="anomaly-total">{anomalies.length} detected</span>
          </div>
          <button className="btn btn-ghost" onClick={onClose}>✕</button>
        </div>

        {/* Severity filter */}
        <div className="severity-filter">
          <button
            className={`severity-btn ${filterSeverity === 'all' ? 'active' : ''}`}
            onClick={() => setFilterSeverity('all')}
          >
            All ({anomalies.length})
          </button>
          <button
            className={`severity-btn severity-error ${filterSeverity === 'error' ? 'active' : ''}`}
            onClick={() => setFilterSeverity('error')}
          >
            🔴 Errors ({counts.error})
          </button>
          <button
            className={`severity-btn severity-warning ${filterSeverity === 'warning' ? 'active' : ''}`}
            onClick={() => setFilterSeverity('warning')}
          >
            🟡 Warnings ({counts.warning})
          </button>
          <button
            className={`severity-btn severity-info ${filterSeverity === 'info' ? 'active' : ''}`}
            onClick={() => setFilterSeverity('info')}
          >
            🔵 Info ({counts.info})
          </button>
        </div>

        {/* Anomaly list */}
        <div className="anomaly-list">
          {filteredAnomalies.length === 0 ? (
            <div className="empty-state">
              No anomalies in this category
            </div>
          ) : (
            filteredAnomalies.map(anomaly => (
              <AnomalyCard
                key={anomaly.id}
                anomaly={anomaly}
                selected={selectedAnomaly?.id === anomaly.id}
                onClick={() => handleAnomalyClick(anomaly)}
              />
            ))
          )}
        </div>

        {/* Detail view */}
        {selectedAnomaly && (
          <AnomalyDetail
            anomaly={selectedAnomaly}
            onClose={() => setSelectedAnomaly(null)}
          />
        )}

        <style>{`
          .anomaly-panel-overlay {
            position: fixed;
            top: 0;
            right: 0;
            bottom: 0;
            width: 400px;
            background: var(--bg-secondary);
            border-left: 1px solid var(--border-color);
            display: flex;
            flex-direction: column;
            z-index: 100;
            box-shadow: -4px 0 20px rgba(0, 0, 0, 0.3);
          }

          .anomaly-panel {
            display: flex;
            flex-direction: column;
            height: 100%;
          }

          .anomaly-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 16px;
            border-bottom: 1px solid var(--border-color);
          }

          .anomaly-title {
            display: flex;
            align-items: center;
            gap: 12px;
          }

          .anomaly-title h3 {
            margin: 0;
            font-size: 16px;
          }

          .anomaly-total {
            font-size: 12px;
            padding: 2px 8px;
            background: var(--accent-red);
            color: white;
            border-radius: 10px;
          }

          .severity-filter {
            display: flex;
            gap: 4px;
            padding: 12px 16px;
            border-bottom: 1px solid var(--border-color);
            overflow-x: auto;
          }

          .severity-btn {
            padding: 6px 12px;
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 6px;
            color: var(--text-secondary);
            font-size: 12px;
            cursor: pointer;
            white-space: nowrap;
            transition: all 0.15s;
          }

          .severity-btn:hover {
            background: var(--bg-hover);
          }

          .severity-btn.active {
            background: var(--bg-tertiary);
            border-color: var(--accent-blue);
            color: var(--text-primary);
          }

          .anomaly-list {
            flex: 1;
            overflow-y: auto;
            padding: 12px;
            display: flex;
            flex-direction: column;
            gap: 8px;
          }

          .empty-state {
            text-align: center;
            padding: 24px;
            color: var(--text-muted);
          }
        `}</style>
      </div>
    </div>
  );
}

interface AnomalyCardProps {
  anomaly: Anomaly;
  selected: boolean;
  onClick: () => void;
}

function AnomalyCard({ anomaly, selected, onClick }: AnomalyCardProps) {
  const { type, severity, message, timestamp } = anomaly;
  const time = new Date(timestamp).toLocaleTimeString();

  return (
    <div
      className={`anomaly-card severity-${severity} ${selected ? 'selected' : ''}`}
      onClick={onClick}
    >
      <div className="anomaly-card-header">
        <span className="anomaly-type">{getTypeLabel(type)}</span>
        <span className="anomaly-time">{time}</span>
      </div>
      <div className="anomaly-message">{message}</div>

      <style>{`
        .anomaly-card {
          padding: 12px;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .anomaly-card:hover {
          background: var(--bg-hover);
        }

        .anomaly-card.selected {
          border-color: var(--accent-blue);
          background: rgba(59, 130, 246, 0.1);
        }

        .anomaly-card.severity-error {
          border-left: 3px solid var(--accent-red);
        }

        .anomaly-card.severity-warning {
          border-left: 3px solid var(--accent-yellow);
        }

        .anomaly-card.severity-info {
          border-left: 3px solid var(--accent-blue);
        }

        .anomaly-card-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 6px;
        }

        .anomaly-type {
          font-size: 11px;
          font-weight: 600;
          color: var(--text-secondary);
          text-transform: uppercase;
        }

        .anomaly-time {
          font-size: 10px;
          color: var(--text-muted);
        }

        .anomaly-message {
          font-size: 13px;
          color: var(--text-primary);
        }
      `}</style>
    </div>
  );
}

interface AnomalyDetailProps {
  anomaly: Anomaly;
  onClose: () => void;
}

function AnomalyDetail({ anomaly, onClose }: AnomalyDetailProps) {
  const { type, severity, message, details, timestamp, relatedIds } = anomaly;

  return (
    <div className="anomaly-detail">
      <div className="detail-header">
        <h4>Anomaly Details</h4>
        <button className="btn btn-ghost" onClick={onClose}>✕</button>
      </div>

      <div className="detail-content">
        <div className="detail-row">
          <span className="detail-label">Type:</span>
          <span className="detail-value">{getTypeLabel(type)}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Severity:</span>
          <span className={`severity-badge severity-${severity}`}>{severity}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Time:</span>
          <span className="detail-value">{new Date(timestamp).toLocaleString()}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Message:</span>
          <span className="detail-value">{message}</span>
        </div>
        {details && (
          <div className="detail-row">
            <span className="detail-label">Details:</span>
            <pre className="detail-pre">{details}</pre>
          </div>
        )}
        {relatedIds && (
          <div className="detail-row">
            <span className="detail-label">Related:</span>
            <div className="related-ids">
              {relatedIds.eventId && <span>Event: {relatedIds.eventId}</span>}
              {relatedIds.turnId && <span>Turn: {relatedIds.turnId}</span>}
              {relatedIds.blockIndex !== undefined && <span>Block: #{relatedIds.blockIndex}</span>}
            </div>
          </div>
        )}
      </div>

      <style>{`
        .anomaly-detail {
          border-top: 1px solid var(--border-color);
          background: var(--bg-tertiary);
        }

        .detail-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 12px 16px;
          border-bottom: 1px solid var(--border-color);
        }

        .detail-header h4 {
          margin: 0;
          font-size: 14px;
        }

        .detail-content {
          padding: 16px;
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        .detail-row {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }

        .detail-label {
          font-size: 11px;
          font-weight: 600;
          color: var(--text-muted);
          text-transform: uppercase;
        }

        .detail-value {
          font-size: 13px;
          color: var(--text-primary);
        }

        .severity-badge {
          display: inline-block;
          padding: 2px 8px;
          border-radius: 4px;
          font-size: 11px;
          font-weight: 600;
          text-transform: uppercase;
        }

        .severity-badge.severity-error {
          background: rgba(239, 68, 68, 0.2);
          color: var(--accent-red);
        }

        .severity-badge.severity-warning {
          background: rgba(245, 158, 11, 0.2);
          color: var(--accent-yellow);
        }

        .severity-badge.severity-info {
          background: rgba(59, 130, 246, 0.2);
          color: var(--accent-blue);
        }

        .detail-pre {
          font-size: 11px;
          padding: 8px;
          background: var(--bg-secondary);
          border-radius: 4px;
          overflow-x: auto;
        }

        .related-ids {
          display: flex;
          flex-direction: column;
          gap: 4px;
          font-size: 12px;
          font-family: monospace;
          color: var(--accent-cyan);
        }
      `}</style>
    </div>
  );
}

function getTypeLabel(type: AnomalyType): string {
  switch (type) {
    case 'orphan_event': return 'Orphan Event';
    case 'missing_correlation': return 'Missing Correlation';
    case 'timing_outlier': return 'Timing Outlier';
    case 'sequence_gap': return 'Sequence Gap';
    case 'schema_error': return 'Schema Error';
    default: return type;
  }
}

export default AnomalyPanel;
