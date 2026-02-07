import type { Meta, StoryObj } from '@storybook/react';
import { EventTrackLane } from './EventTrackLane';
import { mockEvents } from '../mocks/data';

const meta: Meta<typeof EventTrackLane> = {
  title: 'Debug UI/EventTrackLane',
  component: EventTrackLane,
  parameters: {
    layout: 'padded',
  },
  decorators: [
    (Story) => (
      <div style={{ maxWidth: '300px', background: 'var(--bg-primary)', padding: '8px' }}>
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof EventTrackLane>;

export const Default: Story = {
  args: {
    events: mockEvents,
  },
};

export const WithSelection: Story = {
  args: {
    events: mockEvents,
    selectedSeq: mockEvents[0].seq,
  },
};

export const Empty: Story = {
  args: {
    events: [],
  },
};

export const SingleEvent: Story = {
  args: {
    events: [mockEvents[0]],
  },
};

export const ManyEvents: Story = {
  args: {
    events: [
      ...mockEvents,
      ...mockEvents.map((e, i) => ({ ...e, seq: e.seq + 1000000 * (i + 1), id: `${e.id}-dup${i}` })),
    ],
  },
};
