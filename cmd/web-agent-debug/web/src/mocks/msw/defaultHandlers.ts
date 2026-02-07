import {
  mockConversationDetail,
  mockConversations,
  mockSessions,
} from '../fixtures/conversations';
import { mockEvents, mockMwTrace } from '../fixtures/events';
import { mockTimelineEntities } from '../fixtures/timeline';
import { mockTurnDetail, mockTurns } from '../fixtures/turns';
import { createDebugHandlers, type DebugHandlerData } from './createDebugHandlers';

export const defaultDebugHandlerData: DebugHandlerData = {
  conversations: mockConversations,
  conversationDetail: mockConversationDetail,
  sessions: mockSessions,
  turns: mockTurns,
  turnDetail: mockTurnDetail,
  events: mockEvents,
  timelineEntities: mockTimelineEntities,
  mwTrace: mockMwTrace,
};

export function createDefaultDebugHandlers(dataOverrides: Partial<DebugHandlerData> = {}) {
  return createDebugHandlers({
    data: {
      ...defaultDebugHandlerData,
      ...dataOverrides,
    },
  });
}

export const defaultHandlers = createDefaultDebugHandlers();
