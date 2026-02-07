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

## Helper Usage Rules

- Prefer shared helpers under `src/ui/format` and `src/ui/presentation` for formatting, icon/badge mapping, and timeline/event/block presentation logic.
- Do not duplicate formatter/mapping helpers inside components or routes when a shared helper already exists.
- Keep component-local helpers only for strictly local view concerns that are not reused elsewhere.
- When introducing new shared helper behavior, add or update unit tests in:
  - `src/ui/format/format.test.ts`
  - `src/ui/presentation/presentation.test.ts`
- If adding a new helper module, migrate call sites in the same change rather than leaving duplicate interim implementations.

## Style Contract

Styling follows a hybrid contract: namespaced classes as the primary API, with optional `data-part` hooks for targeted overrides in reusable shells.

- Root class:
  - Each component/page exposes a stable root class (for example `app-shell`, `timeline-lanes`, `event-inspector`, `overview-page`).
- Part naming:
  - Internal classes use component-prefixed names (for example `app-header-nav`, `timeline-lane-header`, `anomaly-detail-row`) to avoid global collisions.
- State modifiers:
  - State remains explicit via modifiers like `.active`, `.selected`, `.collapsed`, but only on namespaced elements.
- Token policy:
  - Repeated colors and alpha surfaces must come from `src/styles/tokens.css` variables.
  - Avoid hard-coded `rgba(...)`/hex values in component/primitives CSS unless introducing a new token.
- Override hook:
  - Use `data-part` only where external composition/theming needs a semantic hook beyond class selectors.

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
