# Debug UI Frontend

Debug visualization UI for the geppetto/pinocchio webchat system.

## Setup

```bash
npm install
npx msw init public --save
npm run storybook
```

## Architecture

- **Storybook 8** for component development
- **MSW 2** for API mocking
- **RTK Query** for data fetching
- **Redux Toolkit** for state management
- **React Router v6** for navigation

## Component Roadmap

See PI-013 tasks.md for full breakdown:

### Done (need recreation after git mishap)
- ConversationCard
- BlockCard (with expandable metadata)
- CorrelationIdBar
- SessionList
- TurnInspector (phase tabs, turn metadata)
- EventCard
- MiddlewareChainView
- TimelineEntityCard

### TODO
- Screen 1: Session Overview (three-lane timeline)
- Screen 3: Snapshot Diff
- Screen 5: Event Inspector
- Screen 6: Structured Sink View
- Screen 7: FilterBar
- Screen 8: AnomalyPanel
- AppShell & Routing
- Live WebSocket features

## API Types

All types use `session_id` terminology (not `run_id`) per PI-017.
