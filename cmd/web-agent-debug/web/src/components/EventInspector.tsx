import React, { useState } from 'react';
import type { SemEvent, ParsedBlock, TimelineEntity } from '../types';
import { CorrelationIdBar } from './CorrelationIdBar';

export type ViewMode = 'semantic' | 'sem' | 'raw';

export interface CorrelatedNodes {
  block?: ParsedBlock;
  prevEvent?: SemEvent;
  nextEvent?: SemEvent;
  entity?: TimelineEntity;
}

export interface TrustCheck {
  name: string;
  passed: boolean;
  message?: string;
}

export interface EventInspectorProps {
  event: SemEvent;
  correlatedNodes?: CorrelatedNodes;
  trustChecks?: TrustCheck[];
  onBlockClick?: (block: ParsedBlock) => void;
  onEventClick?: (event: SemEvent) => void;
  onEntityClick?: (entity: TimelineEntity) => void;
}

export function EventInspector({
  event,
  correlatedNodes,
  trustChecks,
  onBlockClick,
  onEventClick,
  onEntityClick,
}: EventInspectorProps) {
  const [viewMode, setViewMode] = useState<ViewMode>('semantic');

  // Extract correlation IDs from event data
  const sessionId = (event.data as Record<string, unknown>)?.session_id as string | undefined;
  const inferenceId = (event.data as Record<string, unknown>)?.inference_id as string | undefined;
  const turnId = (event.data as Record<string, unknown>)?.turn_id as string | undefined;

  return (
    <div className="event-inspector">
      {/* Correlation IDs */}
      <CorrelationIdBar
        sessionId={sessionId}
        inferenceId={inferenceId}
        turnId={turnId}
        seq={event.seq}
        streamId={event.stream_id}
      />

      {/* View mode tabs */}
      <ViewModeTabs activeMode={viewMode} onModeChange={setViewMode} />

      {/* Content based on view mode */}
      <div className="event-content">
        {viewMode === 'semantic' && <SemanticView event={event} />}
        {viewMode === 'sem' && <SemEnvelopeView event={event} />}
        {viewMode === 'raw' && <RawWireView event={event} />}
      </div>

      {/* Correlated nodes panel */}
      {correlatedNodes && (
        <CorrelatedNodesPanel
          nodes={correlatedNodes}
          onBlockClick={onBlockClick}
          onEventClick={onEventClick}
          onEntityClick={onEntityClick}
        />
      )}

      {/* Trust signals */}
      {trustChecks && trustChecks.length > 0 && (
        <TrustSignals checks={trustChecks} />
      )}

      <style>{`
        .event-inspector {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .event-content {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          padding: 16px;
        }
      `}</style>
    </div>
  );
}

interface ViewModeTabsProps {
  activeMode: ViewMode;
  onModeChange: (mode: ViewMode) => void;
}

function ViewModeTabs({ activeMode, onModeChange }: ViewModeTabsProps) {
  const modes: { mode: ViewMode; label: string; icon: string }[] = [
    { mode: 'semantic', label: 'Semantic', icon: '📖' },
    { mode: 'sem', label: 'SEM Frame', icon: '{ }' },
    { mode: 'raw', label: 'Raw Wire', icon: '⚡' },
  ];

  return (
    <div className="view-mode-tabs">
      {modes.map(({ mode, label, icon }) => (
        <button
          key={mode}
          className={`view-mode-tab ${activeMode === mode ? 'active' : ''}`}
          onClick={() => onModeChange(mode)}
        >
          <span className="tab-icon">{icon}</span>
          <span className="tab-label">{label}</span>
        </button>
      ))}

      <style>{`
        .view-mode-tabs {
          display: flex;
          gap: 4px;
          background: var(--bg-secondary);
          padding: 4px;
          border-radius: 6px;
        }

        .view-mode-tab {
          flex: 1;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 6px;
          padding: 8px 16px;
          border: none;
          background: transparent;
          color: var(--text-secondary);
          font-size: 13px;
          cursor: pointer;
          border-radius: 4px;
          transition: all 0.15s;
        }

        .view-mode-tab:hover {
          background: var(--bg-hover);
          color: var(--text-primary);
        }

        .view-mode-tab.active {
          background: var(--accent-blue);
          color: white;
        }

        .tab-icon {
          font-size: 14px;
        }
      `}</style>
    </div>
  );
}

interface SemanticViewProps {
  event: SemEvent;
}

function SemanticView({ event }: SemanticViewProps) {
  const { type, id, received_at, data } = event;
  const time = new Date(received_at).toLocaleString();
  const typeInfo = getEventTypeInfo(type);

  return (
    <div className="semantic-view">
      {/* Event header */}
      <div className="semantic-header">
        <span className="event-icon">{typeInfo.icon}</span>
        <span className="event-type" style={{ color: typeInfo.color }}>{type}</span>
      </div>

      {/* Event ID and time */}
      <div className="semantic-meta">
        <div className="meta-row">
          <span className="meta-label">ID:</span>
          <span className="meta-value mono">{id}</span>
        </div>
        <div className="meta-row">
          <span className="meta-label">Received:</span>
          <span className="meta-value">{time}</span>
        </div>
      </div>

      {/* Human-readable content */}
      <div className="semantic-content">
        {renderSemanticContent(type, data)}
      </div>

      <style>{`
        .semantic-view {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .semantic-header {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .event-icon {
          font-size: 28px;
        }

        .event-type {
          font-size: 18px;
          font-weight: 600;
        }

        .semantic-meta {
          display: flex;
          flex-direction: column;
          gap: 4px;
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 4px;
        }

        .meta-row {
          display: flex;
          gap: 8px;
          font-size: 13px;
        }

        .meta-label {
          color: var(--text-muted);
          min-width: 80px;
        }

        .meta-value {
          color: var(--text-primary);
        }

        .meta-value.mono {
          font-family: monospace;
        }

        .semantic-content {
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 4px;
        }
      `}</style>
    </div>
  );
}

function renderSemanticContent(type: string, data: unknown): React.ReactNode {
  const d = data as Record<string, unknown>;

  switch (type) {
    case 'llm.start':
      return (
        <div>
          <p>Started generating response as <strong>{String(d.role)}</strong></p>
        </div>
      );
    case 'llm.delta':
      return (
        <div>
          <p className="text-muted" style={{ marginBottom: '8px' }}>Streaming content:</p>
          <div style={{ whiteSpace: 'pre-wrap', fontSize: '14px' }}>
            {String(d.cumulative || '')}
          </div>
        </div>
      );
    case 'llm.final':
      return (
        <div>
          <p className="text-muted" style={{ marginBottom: '8px' }}>Final response:</p>
          <div style={{ whiteSpace: 'pre-wrap', fontSize: '14px' }}>
            {String(d.text || '')}
          </div>
        </div>
      );
    case 'tool.start':
      return (
        <div>
          <p>Calling tool: <strong style={{ color: 'var(--accent-yellow)' }}>{String(d.name)}</strong></p>
          {Boolean(d.input) && (
            <pre style={{ marginTop: '8px', fontSize: '12px' }}>
              {String(typeof d.input === 'string' ? d.input : JSON.stringify(d.input, null, 2))}
            </pre>
          )}
        </div>
      );
    case 'tool.result':
      return (
        <div>
          <p className="text-muted" style={{ marginBottom: '8px' }}>Tool result:</p>
          <pre style={{ fontSize: '12px' }}>
            {typeof d.result === 'string' ? d.result : JSON.stringify(d.result, null, 2)}
          </pre>
        </div>
      );
    default:
      return (
        <pre style={{ fontSize: '12px' }}>
          {JSON.stringify(d, null, 2)}
        </pre>
      );
  }
}

interface SemEnvelopeViewProps {
  event: SemEvent;
}

function SemEnvelopeView({ event }: SemEnvelopeViewProps) {
  return (
    <div className="sem-envelope-view">
      <h4>SEM Frame Envelope</h4>
      <pre className="json-view">
        {JSON.stringify(event, null, 2)}
      </pre>

      <style>{`
        .sem-envelope-view h4 {
          margin-bottom: 12px;
          font-size: 13px;
          color: var(--text-secondary);
        }

        .json-view {
          font-size: 12px;
          line-height: 1.5;
          overflow-x: auto;
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 4px;
        }
      `}</style>
    </div>
  );
}

interface RawWireViewProps {
  event: SemEvent;
}

function RawWireView({ event }: RawWireViewProps) {
  // In a real implementation, this would show the provider-native format
  const rawData = (event.data as Record<string, unknown>)?.raw_wire;

  return (
    <div className="raw-wire-view">
      <h4>Raw Wire Format (Provider-Native)</h4>
      <pre className="json-view">
        {rawData 
          ? JSON.stringify(rawData, null, 2)
          : '// Raw wire data not available for this event\n// This would show the original provider response format'}
      </pre>

      <style>{`
        .raw-wire-view h4 {
          margin-bottom: 12px;
          font-size: 13px;
          color: var(--text-secondary);
        }

        .json-view {
          font-size: 12px;
          line-height: 1.5;
          overflow-x: auto;
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 4px;
          color: var(--text-muted);
        }
      `}</style>
    </div>
  );
}

interface CorrelatedNodesPanelProps {
  nodes: CorrelatedNodes;
  onBlockClick?: (block: ParsedBlock) => void;
  onEventClick?: (event: SemEvent) => void;
  onEntityClick?: (entity: TimelineEntity) => void;
}

function CorrelatedNodesPanel({ nodes, onBlockClick, onEventClick, onEntityClick }: CorrelatedNodesPanelProps) {
  const { block, prevEvent, nextEvent, entity } = nodes;

  return (
    <div className="correlated-nodes-panel">
      <h4>Correlated Nodes</h4>

      <div className="nodes-grid">
        {block && (
          <div className="node-link" onClick={() => onBlockClick?.(block)}>
            <span className="node-icon">📦</span>
            <span className="node-label">Linked Block</span>
            <span className="node-id">#{block.index} {block.kind}</span>
          </div>
        )}

        {prevEvent && (
          <div className="node-link" onClick={() => onEventClick?.(prevEvent)}>
            <span className="node-icon">⬅️</span>
            <span className="node-label">Previous Event</span>
            <span className="node-id">{prevEvent.type}</span>
          </div>
        )}

        {nextEvent && (
          <div className="node-link" onClick={() => onEventClick?.(nextEvent)}>
            <span className="node-icon">➡️</span>
            <span className="node-label">Next Event</span>
            <span className="node-id">{nextEvent.type}</span>
          </div>
        )}

        {entity && (
          <div className="node-link" onClick={() => onEntityClick?.(entity)}>
            <span className="node-icon">🎯</span>
            <span className="node-label">Timeline Entity</span>
            <span className="node-id">{entity.kind}</span>
          </div>
        )}
      </div>

      <style>{`
        .correlated-nodes-panel {
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 6px;
        }

        .correlated-nodes-panel h4 {
          margin-bottom: 12px;
          font-size: 13px;
          color: var(--text-secondary);
        }

        .nodes-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
          gap: 8px;
        }

        .node-link {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 4px;
          padding: 12px;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .node-link:hover {
          background: var(--bg-hover);
          border-color: var(--accent-blue);
        }

        .node-icon {
          font-size: 20px;
        }

        .node-label {
          font-size: 11px;
          color: var(--text-muted);
        }

        .node-id {
          font-size: 12px;
          font-weight: 500;
          color: var(--text-primary);
        }
      `}</style>
    </div>
  );
}

interface TrustSignalsProps {
  checks: TrustCheck[];
}

function TrustSignals({ checks }: TrustSignalsProps) {
  const passed = checks.filter(c => c.passed).length;
  const total = checks.length;

  return (
    <div className="trust-signals">
      <div className="trust-header">
        <h4>Trust Signals</h4>
        <span className={`trust-score ${passed === total ? 'all-pass' : 'has-fail'}`}>
          {passed}/{total} passed
        </span>
      </div>

      <div className="trust-checks">
        {checks.map((check, idx) => (
          <div key={idx} className={`trust-check ${check.passed ? 'passed' : 'failed'}`}>
            <span className="check-icon">{check.passed ? '✓' : '✗'}</span>
            <span className="check-name">{check.name}</span>
            {check.message && <span className="check-message">{check.message}</span>}
          </div>
        ))}
      </div>

      <style>{`
        .trust-signals {
          padding: 12px;
          background: var(--bg-secondary);
          border-radius: 6px;
        }

        .trust-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 12px;
        }

        .trust-header h4 {
          font-size: 13px;
          color: var(--text-secondary);
        }

        .trust-score {
          font-size: 12px;
          padding: 2px 8px;
          border-radius: 4px;
        }

        .trust-score.all-pass {
          background: rgba(16, 185, 129, 0.2);
          color: var(--accent-green);
        }

        .trust-score.has-fail {
          background: rgba(239, 68, 68, 0.2);
          color: var(--accent-red);
        }

        .trust-checks {
          display: flex;
          flex-direction: column;
          gap: 4px;
        }

        .trust-check {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 6px 8px;
          background: var(--bg-card);
          border-radius: 4px;
          font-size: 12px;
        }

        .trust-check.passed .check-icon {
          color: var(--accent-green);
        }

        .trust-check.failed .check-icon {
          color: var(--accent-red);
        }

        .check-name {
          color: var(--text-primary);
        }

        .check-message {
          color: var(--text-muted);
          font-size: 11px;
          margin-left: auto;
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
    return { icon: '🤖', color: 'var(--accent-blue)' };
  }
  if (type.startsWith('tool.')) {
    if (type === 'tool.start') return { icon: '🔧', color: 'var(--accent-yellow)' };
    if (type === 'tool.result') return { icon: '📤', color: 'var(--accent-cyan)' };
    return { icon: '🔧', color: 'var(--accent-yellow)' };
  }
  return { icon: '📦', color: 'var(--text-muted)' };
}

export default EventInspector;
