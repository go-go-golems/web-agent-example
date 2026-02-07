import type { TurnDetail, TurnSnapshot } from '../../types';
import { mockTurnDetail, mockTurns } from '../fixtures/turns';
import { pickByIndex } from './common';

export function makeTurnSnapshot(options: {
  index?: number;
  overrides?: Partial<TurnSnapshot>;
} = {}): TurnSnapshot {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockTurns, index),
    ...overrides,
  };
}

export function makeTurnSnapshots(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<TurnSnapshot>;
  } = {}
): TurnSnapshot[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeTurnSnapshot({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}

export function makeTurnDetail(options: {
  overrides?: Partial<TurnDetail>;
} = {}): TurnDetail {
  const { overrides = {} } = options;
  return {
    ...pickByIndex([mockTurnDetail], 0),
    ...overrides,
  };
}
