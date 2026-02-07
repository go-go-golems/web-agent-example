import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { SessionList } from './SessionList';
import { mockConversations } from '../mocks/data';

const meta: Meta<typeof SessionList> = {
  title: 'Debug UI/SessionList',
  component: SessionList,
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<typeof SessionList>;

export const WithData: Story = {
  args: {
    conversations: mockConversations,
    isLoading: false,
  },
};

export const Loading: Story = {
  args: {
    isLoading: true,
  },
};

export const Error: Story = {
  args: {
    error: 'Failed to connect to server',
    isLoading: false,
  },
};

export const Empty: Story = {
  args: {
    conversations: [],
    isLoading: false,
  },
};

export const SingleConversation: Story = {
  args: {
    conversations: [mockConversations[0]],
    isLoading: false,
  },
};

export const ManyConversations: Story = {
  args: {
    conversations: [
      ...mockConversations,
      { ...mockConversations[0], id: 'conv_extra1', session_id: 'sess_e1' },
      { ...mockConversations[1], id: 'conv_extra2', session_id: 'sess_e2' },
      { ...mockConversations[2], id: 'conv_extra3', session_id: 'sess_e3' },
      { ...mockConversations[0], id: 'conv_extra4', session_id: 'sess_e4' },
      { ...mockConversations[1], id: 'conv_extra5', session_id: 'sess_e5' },
    ],
    isLoading: false,
  },
};

// Story with MSW mocking
export const WithMSW: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get('/debug/conversations', () => {
          return HttpResponse.json({ conversations: mockConversations });
        }),
      ],
    },
  },
};
