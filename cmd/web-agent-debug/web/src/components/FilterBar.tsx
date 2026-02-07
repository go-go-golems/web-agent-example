import React from 'react';
import type { BlockKind } from '../types';

export interface FilterState {
  blockKinds: BlockKind[];
  eventTypes: string[];
  timeRange?: { start: number; end: number };
  searchQuery: string;
  showEmpty: boolean;
}

export interface FilterBarProps {
  filters: FilterState;
  onFiltersChange: (filters: FilterState) => void;
  onClose?: () => void;
}

const ALL_BLOCK_KINDS: BlockKind[] = ['system', 'user', 'llm_text', 'tool_call', 'tool_use', 'reasoning'];
const ALL_EVENT_TYPES = ['llm.start', 'llm.delta', 'llm.final', 'tool.start', 'tool.result', 'tool.done', 'log'];

export function FilterBar({ filters, onFiltersChange, onClose }: FilterBarProps) {
  const toggleBlockKind = (kind: BlockKind) => {
    const current = filters.blockKinds;
    const updated = current.includes(kind)
      ? current.filter(k => k !== kind)
      : [...current, kind];
    onFiltersChange({ ...filters, blockKinds: updated });
  };

  const toggleEventType = (type: string) => {
    const current = filters.eventTypes;
    const updated = current.includes(type)
      ? current.filter(t => t !== type)
      : [...current, type];
    onFiltersChange({ ...filters, eventTypes: updated });
  };

  const setSearchQuery = (query: string) => {
    onFiltersChange({ ...filters, searchQuery: query });
  };

  const toggleShowEmpty = () => {
    onFiltersChange({ ...filters, showEmpty: !filters.showEmpty });
  };

  const clearAll = () => {
    onFiltersChange({
      blockKinds: [],
      eventTypes: [],
      timeRange: undefined,
      searchQuery: '',
      showEmpty: true,
    });
  };

  const selectAllBlocks = () => {
    onFiltersChange({ ...filters, blockKinds: [...ALL_BLOCK_KINDS] });
  };

  const selectAllEvents = () => {
    onFiltersChange({ ...filters, eventTypes: [...ALL_EVENT_TYPES] });
  };

  const activeFilterCount = 
    filters.blockKinds.length + 
    filters.eventTypes.length + 
    (filters.searchQuery ? 1 : 0) +
    (filters.timeRange ? 1 : 0);

  return (
    <div className="filter-bar">
      {/* Header */}
      <div className="filter-header">
        <div className="filter-title">
          <h3>Filters</h3>
          {activeFilterCount > 0 && (
            <span className="filter-count">{activeFilterCount} active</span>
          )}
        </div>
        <div className="filter-actions">
          <button className="btn btn-ghost" onClick={clearAll}>Clear All</button>
          {onClose && (
            <button className="btn btn-ghost" onClick={onClose}>✕</button>
          )}
        </div>
      </div>

      {/* Search */}
      <div className="filter-section">
        <label className="filter-label">Search</label>
        <input
          type="text"
          className="filter-input"
          placeholder="Search blocks, events, metadata..."
          value={filters.searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {/* Block kinds */}
      <div className="filter-section">
        <div className="filter-section-header">
          <label className="filter-label">Block Kinds</label>
          <button className="btn btn-ghost text-xs" onClick={selectAllBlocks}>Select All</button>
        </div>
        <div className="filter-chips">
          {ALL_BLOCK_KINDS.map(kind => (
            <FilterChip
              key={kind}
              label={kind}
              icon={getKindIcon(kind)}
              active={filters.blockKinds.includes(kind)}
              onClick={() => toggleBlockKind(kind)}
              colorClass={`kind-${kind}`}
            />
          ))}
        </div>
      </div>

      {/* Event types */}
      <div className="filter-section">
        <div className="filter-section-header">
          <label className="filter-label">Event Types</label>
          <button className="btn btn-ghost text-xs" onClick={selectAllEvents}>Select All</button>
        </div>
        <div className="filter-chips">
          {ALL_EVENT_TYPES.map(type => (
            <FilterChip
              key={type}
              label={type}
              active={filters.eventTypes.includes(type)}
              onClick={() => toggleEventType(type)}
              colorClass={`event-${type.replace('.', '-')}`}
            />
          ))}
        </div>
      </div>

      {/* Options */}
      <div className="filter-section">
        <label className="filter-label">Options</label>
        <label className="filter-checkbox">
          <input
            type="checkbox"
            checked={filters.showEmpty}
            onChange={toggleShowEmpty}
          />
          <span>Show empty phases/lanes</span>
        </label>
      </div>

      <style>{`
        .filter-bar {
          display: flex;
          flex-direction: column;
          gap: 16px;
          padding: 16px;
          background: var(--bg-secondary);
          border-radius: 8px;
          min-width: 280px;
        }

        .filter-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .filter-title {
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .filter-title h3 {
          font-size: 16px;
          margin: 0;
        }

        .filter-count {
          font-size: 11px;
          padding: 2px 8px;
          background: var(--accent-blue);
          color: white;
          border-radius: 10px;
        }

        .filter-actions {
          display: flex;
          gap: 4px;
        }

        .filter-section {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .filter-section-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .filter-label {
          font-size: 12px;
          font-weight: 600;
          color: var(--text-secondary);
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }

        .filter-input {
          padding: 8px 12px;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          color: var(--text-primary);
          font-size: 13px;
        }

        .filter-input:focus {
          outline: none;
          border-color: var(--accent-blue);
        }

        .filter-input::placeholder {
          color: var(--text-muted);
        }

        .filter-chips {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }

        .filter-checkbox {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 13px;
          color: var(--text-secondary);
          cursor: pointer;
        }

        .filter-checkbox input {
          width: 16px;
          height: 16px;
          cursor: pointer;
        }
      `}</style>
    </div>
  );
}

interface FilterChipProps {
  label: string;
  icon?: string;
  active: boolean;
  onClick: () => void;
  colorClass?: string;
}

function FilterChip({ label, icon, active, onClick, colorClass }: FilterChipProps) {
  return (
    <button
      className={`filter-chip ${active ? 'active' : ''} ${colorClass || ''}`}
      onClick={onClick}
    >
      {icon && <span className="chip-icon">{icon}</span>}
      <span className="chip-label">{label}</span>

      <style>{`
        .filter-chip {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          padding: 4px 10px;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: 14px;
          color: var(--text-secondary);
          font-size: 12px;
          cursor: pointer;
          transition: all 0.15s;
        }

        .filter-chip:hover {
          border-color: var(--border-light);
          color: var(--text-primary);
        }

        .filter-chip.active {
          background: var(--accent-blue);
          border-color: var(--accent-blue);
          color: white;
        }

        .filter-chip.active.kind-system { background: var(--block-system); border-color: var(--block-system); }
        .filter-chip.active.kind-user { background: var(--block-user); border-color: var(--block-user); }
        .filter-chip.active.kind-llm_text { background: var(--block-assistant); border-color: var(--block-assistant); }
        .filter-chip.active.kind-tool_call { background: var(--block-tool-call); border-color: var(--block-tool-call); }
        .filter-chip.active.kind-tool_use { background: var(--block-tool-use); border-color: var(--block-tool-use); }
        .filter-chip.active.kind-reasoning { background: var(--block-reasoning); border-color: var(--block-reasoning); }

        .chip-icon {
          font-size: 12px;
        }
      `}</style>
    </button>
  );
}

function getKindIcon(kind: BlockKind): string {
  switch (kind) {
    case 'system': return '⚙️';
    case 'user': return '👤';
    case 'llm_text': return '🤖';
    case 'tool_call': return '🔧';
    case 'tool_use': return '📤';
    case 'reasoning': return '💭';
    default: return '📦';
  }
}

export default FilterBar;
