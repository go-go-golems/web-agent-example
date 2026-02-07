import type { TimelineEntity } from '../../types';
import { mockTimelineEntities } from '../fixtures/timeline';
import { pickByIndex } from './common';

export function makeTimelineEntity(options: {
  index?: number;
  overrides?: Partial<TimelineEntity>;
} = {}): TimelineEntity {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockTimelineEntities, index),
    ...overrides,
  };
}

export function makeTimelineEntities(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<TimelineEntity>;
  } = {}
): TimelineEntity[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeTimelineEntity({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}
