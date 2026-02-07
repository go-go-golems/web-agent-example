import type { Meta, StoryObj } from '@storybook/react';
import { EventInspector } from './EventInspector';
import { mockEvents, mockTimelineEntities } from '../mocks/data';
import type { ParsedBlock } from '../types';

const meta: Meta<typeof EventInspector> = {
  title: 'Debug UI/EventInspector',
  component: EventInspector,
  parameters: {
    layout: 'padded',
  },
};

export default meta;
type Story = StoryObj<typeof EventInspector>;

const mockBlock: ParsedBlock = {
  index: 2,
  id: 'tc_001',
  kind: 'tool_call',
  payload: { id: 'tc_001', name: 'get_weather', args: { location: 'Paris' } },
  metadata: {},
};

export const LLMStart: Story = {
  args: {
    event: mockEvents[0],
  },
};

export const LLMDelta: Story = {
  args: {
    event: mockEvents[3],
  },
};

export const LLMFinal: Story = {
  args: {
    event: mockEvents[5],
  },
};

export const ToolStart: Story = {
  args: {
    event: mockEvents[1],
  },
};

export const ToolResult: Story = {
  args: {
    event: mockEvents[2],
  },
};

export const WithCorrelatedNodes: Story = {
  args: {
    event: mockEvents[1],
    correlatedNodes: {
      block: mockBlock,
      prevEvent: mockEvents[0],
      nextEvent: mockEvents[2],
      entity: mockTimelineEntities[1],
    },
  },
};

export const WithTrustChecks: Story = {
  args: {
    event: mockEvents[0],
    trustChecks: [
      { name: 'Correlation ID present', passed: true },
      { name: 'Sequence monotonic', passed: true },
      { name: 'Timestamp valid', passed: true },
      { name: 'Schema valid', passed: true },
    ],
  },
};

export const WithFailedChecks: Story = {
  args: {
    event: mockEvents[0],
    trustChecks: [
      { name: 'Correlation ID present', passed: true },
      { name: 'Sequence monotonic', passed: false, message: 'Gap detected' },
      { name: 'Timestamp valid', passed: true },
      { name: 'Schema valid', passed: false, message: 'Missing required field' },
    ],
  },
};

export const FullExample: Story = {
  args: {
    event: {
      ...mockEvents[1],
      data: {
        ...(mockEvents[1].data as Record<string, unknown>),
        session_id: 'sess_01234567',
        inference_id: 'inf_abcdef12',
        turn_id: 'turn_001',
      },
    },
    correlatedNodes: {
      block: mockBlock,
      prevEvent: mockEvents[0],
      nextEvent: mockEvents[2],
      entity: mockTimelineEntities[1],
    },
    trustChecks: [
      { name: 'Correlation ID present', passed: true },
      { name: 'Sequence monotonic', passed: true },
      { name: 'Timestamp valid', passed: true },
    ],
  },
};
