import type { ConversationDetail, ConversationSummary, SessionSummary } from '../../types';
import {
  mockConversationDetail,
  mockConversations,
  mockSessions,
} from '../fixtures/conversations';
import { pickByIndex } from './common';

export function makeConversation(options: {
  index?: number;
  overrides?: Partial<ConversationSummary>;
} = {}): ConversationSummary {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockConversations, index),
    ...overrides,
  };
}

export function makeConversations(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<ConversationSummary>;
  } = {}
): ConversationSummary[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeConversation({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}

export function makeConversationDetail(options: {
  index?: number;
  overrides?: Partial<ConversationDetail>;
} = {}): ConversationDetail {
  const { index = 0, overrides = {} } = options;
  const baseDetail = pickByIndex([mockConversationDetail], 0);
  const baseSummary = pickByIndex(mockConversations, index);

  return {
    ...baseDetail,
    ...baseSummary,
    ...overrides,
    engine_config: {
      ...baseDetail.engine_config,
      ...(overrides.engine_config ?? {}),
    },
  };
}

export function makeSession(options: {
  index?: number;
  overrides?: Partial<SessionSummary>;
} = {}): SessionSummary {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockSessions, index),
    ...overrides,
  };
}

export function makeSessions(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<SessionSummary>;
  } = {}
): SessionSummary[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeSession({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}
