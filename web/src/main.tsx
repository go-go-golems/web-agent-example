import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { ErrorBoundary } from '@pwchat/components/ErrorBoundary';
import { store } from '@pwchat/store/store';
import { registerWebAgentSem } from './sem/registerWebAgentSem';
import { App } from './App';
import 'bootstrap/dist/css/bootstrap.min.css';
import './styles.css';

const root = document.getElementById('root');
if (!root) {
  throw new Error('Missing #root element');
}

registerWebAgentSem();

ReactDOM.createRoot(root).render(
  <React.StrictMode>
    <Provider store={store}>
      <ErrorBoundary>
        <App />
      </ErrorBoundary>
    </Provider>
  </React.StrictMode>,
);
