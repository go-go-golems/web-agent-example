import type { Meta, StoryObj } from '@storybook/react';
import { AppShell } from './AppShell';
import type { Anomaly } from './AnomalyPanel';
import {
  createDefaultDebugHandlers,
  defaultHandlers,
} from '../mocks/msw/defaultHandlers';

const meta: Meta<typeof AppShell> = {
  title: 'Debug UI/AppShell',
  component: AppShell,
  parameters: {
    layout: 'fullscreen',
    msw: {
      handlers: defaultHandlers,
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
      handlers: createDefaultDebugHandlers({ conversations: [] }),
    },
  },
};

export const Loading: Story = {
  parameters: {
    msw: {
      handlers: createDefaultDebugHandlers({ conversations: [] }, { delayMs: { conversations: 10_000 } }),
    },
  },
};
