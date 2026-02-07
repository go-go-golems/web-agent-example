import type { Meta, StoryObj } from '@storybook/react';
import { TimelineLanes } from './TimelineLanes';
import { mockTurns, mockEvents, mockTimelineEntities } from '../mocks/data';

const meta: Meta<typeof TimelineLanes> = {
  title: 'Debug UI/TimelineLanes',
  component: TimelineLanes,
  parameters: {
    layout: 'fullscreen',
  },
  decorators: [
    (Story) => (
      <div style={{ height: '600px', background: 'var(--bg-primary)' }}>
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof TimelineLanes>;

export const Default: Story = {
  args: {
    turns: mockTurns,
    events: mockEvents,
    entities: mockTimelineEntities,
  },
};

export const WithSelection: Story = {
  args: {
    turns: mockTurns,
    events: mockEvents,
    entities: mockTimelineEntities,
    selectedTurnId: 'turn_01',
    selectedEventSeq: mockEvents[0].seq,
    selectedEntityId: mockTimelineEntities[0].id,
  },
};

export const Live: Story = {
  args: {
    turns: mockTurns,
    events: mockEvents,
    entities: mockTimelineEntities,
    isLive: true,
  },
};

export const Empty: Story = {
  args: {
    turns: [],
    events: [],
    entities: [],
  },
};

export const TurnsOnly: Story = {
  args: {
    turns: mockTurns,
    events: [],
    entities: [],
  },
};

export const EventsOnly: Story = {
  args: {
    turns: [],
    events: mockEvents,
    entities: [],
  },
};

export const ManyItems: Story = {
  args: {
    turns: [
      ...mockTurns,
      { ...mockTurns[0], turn_id: 'turn_03', phase: 'pre_inference' as const },
      { ...mockTurns[0], turn_id: 'turn_03', phase: 'post_inference' as const },
      { ...mockTurns[0], turn_id: 'turn_03', phase: 'final' as const },
      { ...mockTurns[1], turn_id: 'turn_04', phase: 'pre_inference' as const },
      { ...mockTurns[1], turn_id: 'turn_04', phase: 'final' as const },
    ],
    events: [
      ...mockEvents,
      ...mockEvents.map((e, i) => ({ ...e, seq: e.seq + 1000000 * (i + 1) })),
      ...mockEvents.map((e, i) => ({ ...e, seq: e.seq + 2000000 * (i + 1) })),
    ],
    entities: [
      ...mockTimelineEntities,
      { ...mockTimelineEntities[0], id: 'msg-extra1', created_at: Date.now() - 10000 },
      { ...mockTimelineEntities[1], id: 'tc-extra1', created_at: Date.now() - 8000 },
      { ...mockTimelineEntities[3], id: 'msg-extra2', created_at: Date.now() - 5000 },
    ],
  },
};
