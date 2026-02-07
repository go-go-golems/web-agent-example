import type { MwTrace, SemEvent } from '../../types';
import { mockEvents, mockMwTrace } from '../fixtures/events';
import { pickByIndex } from './common';

export function makeEvent(options: {
  index?: number;
  overrides?: Partial<SemEvent>;
} = {}): SemEvent {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockEvents, index),
    ...overrides,
  };
}

export function makeEvents(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<SemEvent>;
  } = {}
): SemEvent[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeEvent({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}

export function makeMwTrace(options: {
  overrides?: Partial<MwTrace>;
} = {}): MwTrace {
  const { overrides = {} } = options;
  return {
    ...pickByIndex([mockMwTrace], 0),
    ...overrides,
  };
}
