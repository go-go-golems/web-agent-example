import React, { useState } from 'react';
import { Outlet, NavLink, useParams } from 'react-router-dom';
import { SessionList } from './SessionList';
import { FilterBar, FilterState } from './FilterBar';
import { AnomalyPanel, Anomaly } from './AnomalyPanel';
import { useAppSelector } from '../store/hooks';

const defaultFilters: FilterState = {
  blockKinds: [],
  eventTypes: [],
  searchQuery: '',
  showEmpty: true,
};

export interface AppShellProps {
  /** Mock anomalies for Storybook */
  anomalies?: Anomaly[];
}

export function AppShell({ anomalies = [] }: AppShellProps) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [filterOpen, setFilterOpen] = useState(false);
  const [anomalyOpen, setAnomalyOpen] = useState(false);
  const [filters, setFilters] = useState<FilterState>(defaultFilters);

  const selectedConvId = useAppSelector((state) => state.ui.selectedConvId);
  const { sessionId } = useParams();

  const activeFilterCount = 
    filters.blockKinds.length + 
    filters.eventTypes.length + 
    (filters.searchQuery ? 1 : 0);

  return (
    <div className="app-shell">
      {/* Top nav bar */}
      <header className="app-header">
        <div className="header-left">
          <button 
            className="btn btn-icon"
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            title={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {sidebarCollapsed ? '☰' : '◀'}
          </button>
          <h1 className="app-title">🔍 Debug UI</h1>
        </div>

        <nav className="header-nav">
          <NavLink to="/" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`} end>
            Overview
          </NavLink>
          <NavLink to="/timeline" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            Timeline
          </NavLink>
          <NavLink to="/events" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            Events
          </NavLink>
        </nav>

        <div className="header-right">
          <button 
            className={`btn btn-icon ${activeFilterCount > 0 ? 'has-badge' : ''}`}
            onClick={() => setFilterOpen(!filterOpen)}
            title="Filters"
          >
            🔍
            {activeFilterCount > 0 && <span className="badge">{activeFilterCount}</span>}
          </button>
          <button 
            className={`btn btn-icon ${anomalies.length > 0 ? 'has-badge' : ''}`}
            onClick={() => setAnomalyOpen(!anomalyOpen)}
            title="Anomalies"
          >
            ⚠️
            {anomalies.length > 0 && <span className="badge">{anomalies.length}</span>}
          </button>
        </div>
      </header>

      <div className="app-body">
        {/* Sidebar */}
        <aside className={`app-sidebar ${sidebarCollapsed ? 'collapsed' : ''}`}>
          {!sidebarCollapsed && <SessionList />}
        </aside>

        {/* Main content */}
        <main className="app-main">
          {/* Breadcrumb */}
          <div className="breadcrumb">
            <span className="crumb">
              {selectedConvId ? `conv: ${selectedConvId.slice(0, 8)}...` : 'No conversation selected'}
            </span>
            {sessionId && (
              <>
                <span className="crumb-sep">/</span>
                <span className="crumb">session: {sessionId.slice(0, 8)}...</span>
              </>
            )}
          </div>

          {/* Router outlet */}
          <div className="main-content">
            <Outlet context={{ filters }} />
          </div>
        </main>

        {/* Filter sidebar (right) */}
        {filterOpen && (
          <aside className="filter-sidebar">
            <FilterBar 
              filters={filters} 
              onFiltersChange={setFilters}
              onClose={() => setFilterOpen(false)}
            />
          </aside>
        )}
      </div>

      {/* Anomaly panel overlay */}
      <AnomalyPanel
        anomalies={anomalies}
        isOpen={anomalyOpen}
        onClose={() => setAnomalyOpen(false)}
      />

      <style>{`
        .app-shell {
          display: flex;
          flex-direction: column;
          height: 100vh;
          background: var(--bg-primary);
          color: var(--text-primary);
        }

        .app-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 16px;
          height: 48px;
          background: var(--bg-secondary);
          border-bottom: 1px solid var(--border-color);
          flex-shrink: 0;
        }

        .header-left {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .app-title {
          font-size: 16px;
          font-weight: 600;
          margin: 0;
        }

        .header-nav {
          display: flex;
          gap: 4px;
        }

        .nav-link {
          padding: 8px 16px;
          color: var(--text-secondary);
          text-decoration: none;
          border-radius: 4px;
          font-size: 13px;
          transition: all 0.15s;
        }

        .nav-link:hover {
          background: var(--bg-hover);
          color: var(--text-primary);
        }

        .nav-link.active {
          background: var(--accent-blue);
          color: white;
        }

        .header-right {
          display: flex;
          gap: 8px;
        }

        .btn-icon {
          position: relative;
          width: 32px;
          height: 32px;
          padding: 0;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 16px;
        }

        .btn-icon .badge {
          position: absolute;
          top: -4px;
          right: -4px;
          min-width: 16px;
          height: 16px;
          padding: 0 4px;
          background: var(--accent-red);
          color: white;
          font-size: 10px;
          font-weight: 600;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .app-body {
          display: flex;
          flex: 1;
          overflow: hidden;
        }

        .app-sidebar {
          width: 280px;
          background: var(--bg-secondary);
          border-right: 1px solid var(--border-color);
          display: flex;
          flex-direction: column;
          transition: width 0.2s;
          overflow: hidden;
        }

        .app-sidebar.collapsed {
          width: 0;
          border-right: none;
        }

        .app-main {
          flex: 1;
          display: flex;
          flex-direction: column;
          overflow: hidden;
        }

        .breadcrumb {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 8px 16px;
          background: var(--bg-tertiary);
          border-bottom: 1px solid var(--border-color);
          font-size: 12px;
        }

        .crumb {
          color: var(--text-secondary);
          font-family: monospace;
        }

        .crumb-sep {
          color: var(--text-muted);
        }

        .main-content {
          flex: 1;
          overflow: auto;
          padding: 16px;
        }

        .filter-sidebar {
          width: 320px;
          background: var(--bg-secondary);
          border-left: 1px solid var(--border-color);
          overflow-y: auto;
        }
      `}</style>
    </div>
  );
}

export default AppShell;
