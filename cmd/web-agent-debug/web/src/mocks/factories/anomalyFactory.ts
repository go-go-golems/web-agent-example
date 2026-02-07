import type { Anomaly } from '../../components/AnomalyPanel';
import {
  mockAnomalies,
  mockAppShellAnomalies,
} from '../fixtures/anomalies';
import { pickByIndex } from './common';

export function makeAnomaly(options: {
  index?: number;
  overrides?: Partial<Anomaly>;
} = {}): Anomaly {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockAnomalies, index),
    ...overrides,
  };
}

export function makeAnomalies(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<Anomaly>;
  } = {}
): Anomaly[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeAnomaly({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}

export function makeAppShellAnomaly(options: {
  index?: number;
  overrides?: Partial<Anomaly>;
} = {}): Anomaly {
  const { index = 0, overrides = {} } = options;
  return {
    ...pickByIndex(mockAppShellAnomalies, index),
    ...overrides,
  };
}

export function makeAppShellAnomalies(
  count: number,
  options: {
    startIndex?: number;
    mapOverrides?: (listIndex: number) => Partial<Anomaly>;
  } = {}
): Anomaly[] {
  const { startIndex = 0, mapOverrides } = options;
  return Array.from({ length: count }, (_, listIndex) =>
    makeAppShellAnomaly({
      index: startIndex + listIndex,
      overrides: mapOverrides?.(listIndex),
    })
  );
}
