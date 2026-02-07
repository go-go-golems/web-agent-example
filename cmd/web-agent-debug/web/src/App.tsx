import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { SessionList } from './components/SessionList';

function App() {
  return (
    <div className="app-shell">
      <div className="sidebar">
        <div className="p-3" style={{ borderBottom: '1px solid var(--border-color)' }}>
          <h2>🔍 Debug UI</h2>
        </div>
        <SessionList />
      </div>
      <div className="main-content">
        <Routes>
          <Route path="/" element={<Welcome />} />
          <Route path="/turn/:convId/:sessionId/:turnId" element={<TurnRoute />} />
        </Routes>
      </div>
    </div>
  );
}

function Welcome() {
  return (
    <div className="text-center p-4 text-muted">
      <h2 className="mb-4">Welcome to Debug UI</h2>
      <p>Select a conversation from the sidebar to begin inspecting.</p>
    </div>
  );
}

function TurnRoute() {
  return (
    <div className="text-center p-4 text-muted">
      <p>Turn Inspector - Coming soon</p>
    </div>
  );
}

export default App;
