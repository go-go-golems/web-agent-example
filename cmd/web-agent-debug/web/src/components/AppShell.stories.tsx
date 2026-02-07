import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { AppShell } from './AppShell';
import { mockConversations, mockTurns, mockEvents, mockTimelineEntities } from '../mocks/data';
import type { Anomaly } from './AnomalyPanel';

const meta: Meta<typeof AppShell> = {
  title: 'Debug UI/AppShell',
  component: AppShell,
  parameters: {
    layout: 'fullscreen',
    msw: {
      handlers: [
        http.get('/debug/conversations', () => {
          return HttpResponse.json({ conversations: mockConversations });
        }),
        http.get('/debug/conversation/:id', () => {
          return HttpResponse.json(mockConversations[0]);
        }),
        http.get('/debug/turns', () => {
          return HttpResponse.json(mockTurns);
        }),
        http.get('/debug/events/:id', () => {
          return HttpResponse.json({ events: mockEvents, total: mockEvents.length, buffer_capacity: 1000 });
        }),
        http.get('/debug/timeline', () => {
          return HttpResponse.json({ entities: mockTimelineEntities, version: Date.now() });
        }),
      ],
    },
  },
};

export default meta;
type Story = StoryObj<typeof AppShell>;

const mockAnomalies: Anomaly[] = [
  {
    id: 'anom_001',
    type: 'orphan_event',
    severity: 'error',
    message: 'Event has no matching turn',
    timestamp: new Date().toISOString(),
  },
  {
    id: 'anom_002',
    type: 'timing_outlier',
    severity: 'warning',
    message: 'Slow inference detected',
    timestamp: new Date().toISOString(),
  },
];

export const Default: Story = {
  args: {},
};

export const WithAnomalies: Story = {
  args: {
    anomalies: mockAnomalies,
  },
};

export const EmptyState: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get('/debug/conversations', () => {
          return HttpResponse.json({ conversations: [] });
        }),
      ],
    },
  },
};

export const Loading: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get('/debug/conversations', async () => {
          await new Promise(resolve => setTimeout(resolve, 10000));
          return HttpResponse.json({ conversations: [] });
        }),
      ],
    },
  },
};
