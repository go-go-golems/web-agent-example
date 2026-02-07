import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type {
  ConversationSummary,
  ConversationDetail,
  SessionSummary,
  TurnSnapshot,
  TurnDetail,
  EventsResponse,
  MwTrace,
  TimelineSnapshot,
} from '../types';

export interface TurnQuery {
  convId: string;
  sessionId?: string;
  phase?: string;
  sinceMs?: number;
  limit?: number;
}

export interface TurnDetailQuery {
  convId: string;
  sessionId: string;
  turnId: string;
}

export interface EventQuery {
  convId: string;
  type?: string;
  sinceSeq?: number;
  limit?: number;
}

export interface MwTraceQuery {
  convId: string;
  inferenceId: string;
}

export interface TimelineQuery {
  convId: string;
  sinceVersion?: number;
  limit?: number;
}

export const debugApi = createApi({
  reducerPath: 'debugApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/debug/' }),
  tagTypes: ['Conversations', 'Turns', 'Events', 'Timeline', 'MwTrace'],
  endpoints: (builder) => ({
    // List all conversations
    getConversations: builder.query<ConversationSummary[], void>({
      query: () => 'conversations',
      transformResponse: (response: { conversations: ConversationSummary[] }) =>
        response.conversations,
      providesTags: ['Conversations'],
    }),

    // Get single conversation detail
    getConversation: builder.query<ConversationDetail, string>({
      query: (id) => `conversation/${id}`,
      providesTags: (_result, _error, id) => [{ type: 'Conversations', id }],
    }),

    // List sessions within a conversation
    getSessions: builder.query<SessionSummary[], string>({
      query: (convId) => `conversation/${convId}/sessions`,
      transformResponse: (response: { sessions: SessionSummary[] }) =>
        response.sessions,
    }),

    // Query turn snapshots
    getTurns: builder.query<TurnSnapshot[], TurnQuery>({
      query: ({ convId, sessionId, phase, sinceMs, limit }) => {
        const params = new URLSearchParams();
        if (convId) params.set('conv_id', convId);
        if (sessionId) params.set('session_id', sessionId);
        if (phase) params.set('phase', phase);
        if (sinceMs) params.set('since_ms', String(sinceMs));
        if (limit) params.set('limit', String(limit));
        return `turns?${params.toString()}`;
      },
      providesTags: ['Turns'],
    }),

    // Get turn detail with all phases
    getTurnDetail: builder.query<TurnDetail, TurnDetailQuery>({
      query: ({ convId, sessionId, turnId }) =>
        `turn/${convId}/${sessionId}/${turnId}`,
      providesTags: (_result, _error, { turnId }) => [{ type: 'Turns', id: turnId }],
    }),

    // Get events for a conversation
    getEvents: builder.query<EventsResponse, EventQuery>({
      query: ({ convId, type, sinceSeq, limit }) => {
        const params = new URLSearchParams();
        if (type) params.set('type', type);
        if (sinceSeq) params.set('since_seq', String(sinceSeq));
        if (limit) params.set('limit', String(limit));
        return `events/${convId}?${params.toString()}`;
      },
      providesTags: ['Events'],
    }),

    // Get timeline entities
    getTimeline: builder.query<TimelineSnapshot, TimelineQuery>({
      query: ({ convId, sinceVersion, limit }) => {
        const params = new URLSearchParams();
        params.set('conv_id', convId);
        if (sinceVersion) params.set('since_version', String(sinceVersion));
        if (limit) params.set('limit', String(limit));
        return `timeline?${params.toString()}`;
      },
      providesTags: ['Timeline'],
    }),

    // Get middleware trace
    getMwTrace: builder.query<MwTrace, MwTraceQuery>({
      query: ({ convId, inferenceId }) => `mw-trace/${convId}/${inferenceId}`,
      providesTags: ['MwTrace'],
    }),
  }),
});

export const {
  useGetConversationsQuery,
  useGetConversationQuery,
  useGetSessionsQuery,
  useGetTurnsQuery,
  useGetTurnDetailQuery,
  useGetEventsQuery,
  useGetTimelineQuery,
  useGetMwTraceQuery,
} = debugApi;
