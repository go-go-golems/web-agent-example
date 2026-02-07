# Debug UI Storybook + Widget Authoring Playbook

## Goal

Provide a repeatable way to add or modify debug UI widgets and Storybook stories without reintroducing helper, style, or mock-data duplication.

This playbook is aligned with the PI-019 cleanup architecture in:

- `web-agent-example/cmd/web-agent-debug/web/src/ui`
- `web-agent-example/cmd/web-agent-debug/web/src/styles`
- `web-agent-example/cmd/web-agent-debug/web/src/mocks`

## Core rules (non-negotiable)

1. Do not add runtime `<style>{...}</style>` blocks in TSX runtime components.
2. Do not copy helper functions into components when an existing helper can be reused.
3. Do not build large story arrays inline in stories; use fixtures/factories/scenarios.
4. Prefer scenario-first story args (`make*Scenario`) over raw fixture wiring.
5. Use centralized MSW handlers (`defaultHandlers`, `createDefaultDebugHandlers`) instead of per-story endpoint rewrites.

## Directory map and responsibilities

### Runtime and presentation helpers

- `src/ui/format/*`: string/time/phase formatting helpers.
- `src/ui/presentation/*`: mapping helpers for event/block/timeline visual semantics.

Use these before creating any new ad hoc helper.

### Style layers

- `src/styles/tokens.css`: design tokens only.
- `src/styles/reset.css`: reset/base rules.
- `src/styles/primitives.css`: reusable primitives (cards/buttons/badges/inputs).
- `src/styles/layout.css`: app/page layout.
- `src/styles/components/*.css`: component-scoped rules.
- `src/index.css`: import orchestrator for all style layers.

### Mock architecture

- `src/mocks/fixtures/*`: canonical fixed datasets.
- `src/mocks/factories/*`: deterministic builders from fixtures.
- `src/mocks/scenarios/*`: reusable story contexts (args-level abstraction).
- `src/mocks/msw/createDebugHandlers.ts`: reusable debug endpoint handler builder.
- `src/mocks/msw/defaultHandlers.ts`: baseline handler bundle for most stories.

## Styling guidance for new widgets

## 1) Start from primitives/tokens

Before adding custom CSS:

- check if existing primitive classes cover 80% of needs;
- use tokens from `tokens.css` instead of hard-coded colors/spacings.

If you need a new repeated color/alpha/shadow value, add a token first.

## 2) Use stable style contracts

For component internals:

- prefer stable class names with a component prefix (`.widget-name__part`);
- use `data-part` for semantic subparts you expect to theme/override externally.

Recommended hybrid:

- class names for structure and local specificity;
- `data-part` markers for extension/theming hooks.

## 3) Keep runtime TSX style-free

Do this:

- TSX: structure + state logic + class/data-part wiring.
- CSS file: all style definitions.

Do not do this:

- TSX with inline style-tag blocks.

## 4) Import rules

- Add new component CSS in `src/styles/components/<Widget>.css`.
- Register it once via `src/index.css` import orchestration.

## Storybook story authoring guidance

## 1) Prefer scenario-first stories

Use scenario builders from `src/mocks/scenarios`:

```ts
import { makeTimelineScenario } from '../mocks/scenarios';

export const Default = {
  args: makeTimelineScenario('default').args,
};
```

If the scenario you need does not exist:

1. extend the scenario file,
2. reuse existing factories,
3. then consume from story.

## 2) When to use fixtures vs factories vs scenarios

- Fixtures: fixed canonical data samples.
- Factories: deterministic generation/variation (`makeEvents(n)`, overrides).
- Scenarios: story-ready contexts and named states.

Decision rule:

- If more than one story needs it: make a scenario.
- If story needs count/variant changes: use factory inside scenario.
- If story needs a single canonical sample only: fixture is acceptable.

## 3) Keep stories thin

Stories should mostly contain:

- metadata (`title`, `component`, `parameters`),
- small story exports selecting scenario keys,
- optional decorators for controlled layout wrappers.

Avoid embedding domain logic or payload construction in the story file.

## 4) Naming conventions

- Story titles: `Debug UI/<ComponentName>`.
- Scenario keys: behavior/state focused (`default`, `empty`, `withSelection`, `manyItems`, `errorsOnly`).
- Story exports: readable and minimal (`Default`, `Empty`, `ManyItems`).

## MSW guidance

## 1) Default path

For endpoint behavior equal to normal mock defaults:

```ts
import { defaultHandlers } from '../mocks/msw/defaultHandlers';

parameters: {
  msw: { handlers: defaultHandlers },
}
```

## 2) Override path

For per-story behavior adjustments (empty/loading/custom payload):

```ts
import { createDefaultDebugHandlers } from '../mocks/msw/defaultHandlers';

parameters: {
  msw: {
    handlers: createDefaultDebugHandlers(
      { conversations: [] },
      { delayMs: { conversations: 10_000 } },
    ),
  },
}
```

Rules:

- only override what the story needs;
- keep overrides declarative and local to the story variant;
- do not rewrite entire route trees in component stories unless unavoidable.

## 3) Adding a new endpoint

When a widget needs new API behavior:

1. extend `createDebugHandlers` with typed inputs;
2. thread defaults through `defaultHandlers`;
3. expose override knobs via `createDefaultDebugHandlers` options;
4. update scenarios/stories to use override options, not custom ad hoc handlers.

## Factory and scenario reuse patterns

## 1) Add factory first when repetition appears

If you copy-paste data composition twice, stop and create/extend a factory helper.

Good pattern:

- factory builds deterministic base data;
- scenario composes factory outputs into story args.

## 2) Keep deterministic behavior

Use existing deterministic helpers from `src/mocks/factories/deterministic.ts` for synthetic ids/timestamps/seq values.

Do not use random values in stories. Determinism keeps snapshots and visual comparisons stable.

## 3) Scenarios should encode intent

A scenario should describe UI state, not implementation details.

Examples:

- `withFailedChecks`
- `turnsOnly`
- `manyAnomalies`

Not good:

- `case3`
- `variantB`

## Adding a new widget: step-by-step checklist

1. Create widget component file under `src/components/<Widget>.tsx`.
2. Reuse existing `src/ui/format` / `src/ui/presentation` helpers; add new helper only if needed in multiple components.
3. Create `src/styles/components/<Widget>.css` and wire class/data-part hooks.
4. Import CSS through `src/index.css` orchestration.
5. Add/extend fixture + factory + scenario layers as needed.
6. Create `src/components/<Widget>.stories.tsx`:
   - scenario-first args,
   - default MSW handlers,
   - minimal story-local logic.
7. Run validation:
   - `npm run check:helpers:dedupe`
   - `npm run check:styles:inline-runtime`
   - `npm run build`
   - `npm run build-storybook`
8. If story adds endpoint behavior, ensure it is implemented through central handler builders.
9. Update docs (`src/mocks/README.md` or frontend `README.md`) when contracts change.

## PR review checklist for widget/story changes

Use this checklist during review:

- No runtime inline `<style>` blocks added.
- No duplicate helper signatures introduced in components/routes.
- Story uses scenarios/factories/fixtures correctly (not local large arrays).
- MSW handlers are centralized and override only deltas.
- Styling uses tokens/primitives/layers, with minimal hard-coded values.
- New story states are deterministic and reproducible.
- Build and Storybook build pass.

## Common anti-patterns to reject

- Local helper copy in one component because "it is small".
- Story with a long inline mock payload block.
- Per-story handcrafted endpoint handlers for common routes.
- Direct hard-coded color literals repeated across component CSS.
- New CSS injected in TSX render tree.

## Quick templates

## New scenario template

```ts
// src/mocks/scenarios/myWidgetScenarios.ts
import { makeEvents } from '../factories';

export const myWidgetScenarios = {
  default: {
    args: {
      items: makeEvents(4),
    },
  },
  empty: {
    args: {
      items: [],
    },
  },
} as const;

export type MyWidgetScenarioKey = keyof typeof myWidgetScenarios;

export function makeMyWidgetScenario(key: MyWidgetScenarioKey) {
  return myWidgetScenarios[key];
}
```

## Story usage template

```ts
import type { Meta, StoryObj } from '@storybook/react';
import { MyWidget } from './MyWidget';
import { makeMyWidgetScenario } from '../mocks/scenarios/myWidgetScenarios';

const meta: Meta<typeof MyWidget> = {
  title: 'Debug UI/MyWidget',
  component: MyWidget,
};

export default meta;

type Story = StoryObj<typeof MyWidget>;

export const Default: Story = {
  args: makeMyWidgetScenario('default').args,
};

export const Empty: Story = {
  args: makeMyWidgetScenario('empty').args,
};
```

## Final note

When in doubt, optimize for shared architecture over local convenience. Small local shortcuts are exactly how helper/style/mock duplication returns.
