import { http, HttpResponse } from 'msw';
import {
  mockConversations,
  mockConversationDetail,
  mockSessions,
  mockTurns,
  mockTurnDetail,
  mockEvents,
  mockMwTrace,
  mockTimelineEntities,
} from './data';

export const handlers = [
  // GET /debug/conversations
  http.get('/debug/conversations', () => {
    return HttpResponse.json({ conversations: mockConversations });
  }),

  // GET /debug/conversation/:id
  http.get('/debug/conversation/:id', ({ params }) => {
    const { id } = params;
    if (id === 'conv_8a3f') {
      return HttpResponse.json(mockConversationDetail);
    }
    return HttpResponse.json(
      { ...mockConversations.find((c) => c.id === id), engine_config: mockConversationDetail.engine_config },
      { status: 200 }
    );
  }),

  // GET /debug/conversation/:id/sessions
  http.get('/debug/conversation/:id/sessions', () => {
    return HttpResponse.json({ sessions: mockSessions });
  }),

  // GET /debug/turns
  http.get('/debug/turns', ({ request }) => {
    const url = new URL(request.url);
    const convId = url.searchParams.get('conv_id');
    const sessionId = url.searchParams.get('session_id');
    
    let filtered = mockTurns;
    if (convId) {
      filtered = filtered.filter((t) => t.conv_id === convId);
    }
    if (sessionId) {
      filtered = filtered.filter((t) => t.session_id === sessionId);
    }
    return HttpResponse.json(filtered);
  }),

  // GET /debug/turn/:convId/:sessionId/:turnId
  http.get('/debug/turn/:convId/:sessionId/:turnId', ({ params }) => {
    const { turnId } = params;
    if (turnId === 'turn_01') {
      return HttpResponse.json(mockTurnDetail);
    }
    // Generate a simple detail for other turns
    const turn = mockTurns.find((t) => t.turn_id === turnId);
    if (turn) {
      return HttpResponse.json({
        conv_id: turn.conv_id,
        session_id: turn.session_id,
        turn_id: turn.turn_id,
        phases: {
          final: { captured_at: new Date().toISOString(), turn: turn.turn },
        },
      });
    }
    return HttpResponse.json({ error: 'Turn not found' }, { status: 404 });
  }),

  // GET /debug/events/:convId
  http.get('/debug/events/:convId', () => {
    return HttpResponse.json({
      events: mockEvents,
      total: mockEvents.length,
      buffer_capacity: 1000,
    });
  }),

  // GET /debug/timeline
  http.get('/debug/timeline', () => {
    return HttpResponse.json({
      entities: mockTimelineEntities,
      version: Date.now(),
    });
  }),

  // GET /debug/mw-trace/:convId/:inferenceId
  http.get('/debug/mw-trace/:convId/:inferenceId', () => {
    return HttpResponse.json(mockMwTrace);
  }),
];
