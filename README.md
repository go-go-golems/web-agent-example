# web-agent-example

A reference webchat implementation that showcases custom web agent middleware and UI widgets on top of the Pinocchio webchat stack.

This example currently ships two custom features:
- `webagent-thinking-mode`: injects a system prompt that controls the assistant's reasoning mode.
- `webagent-disco-dialogue`: injects a Disco Elysium-style internal dialogue prompt and streams structured dialogue/check/state events into a custom widget.

## Quickstart

1. Start the backend server:

```bash
go run ./cmd/web-agent-example serve --addr :8080 --log-level debug
```

2. Start the frontend dev server:

```bash
cd web
npm install
npm run dev
```

3. Open the UI at `http://localhost:5174`.

## How It Works

- The backend registers middleware hooks and exposes `POST /chat`, `GET /hydrate`, `GET /timeline`, and `GET /ws` for streaming.
- The frontend uses `@pwchat` components (aliased to `pinocchio/cmd/web-chat/web/src`) and registers custom renderers for disco dialogue timeline entities.
- Structured sink extractors parse tagged YAML blocks emitted by the model and emit typed events. The timeline projector maps those events into snapshot entities that are streamed to the UI.

## Disco Dialogue UI

When the Disco dialogue middleware is enabled, you should see the following entity types appear in the timeline:
- `disco_dialogue_line`
- `disco_dialogue_check`
- `disco_dialogue_state`

Each entity is rendered by `DiscoDialogueCard` in the frontend and styled as a system entry.

## Middleware Overrides

The UI sends overrides as part of the `/chat` request. Example payload:

```json
{
  "conv_id": "...",
  "prompt": "hello",
  "overrides": {
    "middlewares": [
      { "name": "webagent-thinking-mode", "config": { "mode": "fast" } },
      { "name": "webagent-disco-dialogue", "config": { "personas": ["Logic", "Empathy"], "tone": "noir" } }
    ]
  }
}
```

The default profile also enables both middlewares, so the UI toggle only needs to disable the disco middleware when desired.

## Development Notes

- The frontend build outputs to `cmd/web-agent-example/static/dist`.
- Timeline snapshots are stored in memory by default. Use `--timeline-dsn` or `--timeline-db` to persist timeline data.
- Use `--log-level debug` when diagnosing event flow and websocket traffic.
