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
        <div className="app-header-left">
          <button 
            className="btn app-btn-icon"
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            title={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {sidebarCollapsed ? '☰' : '◀'}
          </button>
          <h1 className="app-title">🔍 Debug UI</h1>
        </div>

        <nav className="app-header-nav">
          <NavLink to="/" className={({ isActive }) => `app-nav-link ${isActive ? 'active' : ''}`} end>
            Overview
          </NavLink>
          <NavLink to="/timeline" className={({ isActive }) => `app-nav-link ${isActive ? 'active' : ''}`}>
            Timeline
          </NavLink>
          <NavLink to="/events" className={({ isActive }) => `app-nav-link ${isActive ? 'active' : ''}`}>
            Events
          </NavLink>
        </nav>

        <div className="app-header-right">
          <button 
            className={`btn app-btn-icon ${activeFilterCount > 0 ? 'has-badge' : ''}`}
            onClick={() => setFilterOpen(!filterOpen)}
            title="Filters"
          >
            🔍
            {activeFilterCount > 0 && <span className="app-btn-badge">{activeFilterCount}</span>}
          </button>
          <button 
            className={`btn app-btn-icon ${anomalies.length > 0 ? 'has-badge' : ''}`}
            onClick={() => setAnomalyOpen(!anomalyOpen)}
            title="Anomalies"
          >
            ⚠️
            {anomalies.length > 0 && <span className="app-btn-badge">{anomalies.length}</span>}
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
          <div className="app-breadcrumb">
            <span className="app-breadcrumb-crumb">
              {selectedConvId ? `conv: ${selectedConvId.slice(0, 8)}...` : 'No conversation selected'}
            </span>
            {sessionId && (
              <>
                <span className="app-breadcrumb-sep">/</span>
                <span className="app-breadcrumb-crumb">session: {sessionId.slice(0, 8)}...</span>
              </>
            )}
          </div>

          {/* Router outlet */}
          <div className="app-main-content">
            <Outlet context={{ filters }} />
          </div>
        </main>

        {/* Filter sidebar (right) */}
        {filterOpen && (
          <aside className="app-filter-sidebar">
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
    </div>
  );
}

export default AppShell;
