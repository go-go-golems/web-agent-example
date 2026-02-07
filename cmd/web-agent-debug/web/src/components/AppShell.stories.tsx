import type { Meta, StoryObj } from '@storybook/react';
import { AppShell } from './AppShell';
import {
  createDefaultDebugHandlers,
  defaultHandlers,
} from '../mocks/msw/defaultHandlers';
import { makeAnomalies } from '../mocks/factories';

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

export const Default: Story = {
  args: {},
};

export const WithAnomalies: Story = {
  args: {
    anomalies: makeAnomalies(2),
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
