import { useProgramInstanceStatus } from '../src/composables/programInstanceStatus';
import * as data from '../mocks/data.js';
import { vi, afterEach } from 'vitest';

vi.mock('../src/modules/state');
describe('programInstanceStatus composable', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns the expected status', () => {
    const now = new Date().valueOf();
    const started = 16; //days
    const startDate =
      new Date(now - started * 24 * 60 * 60 * 1000).valueOf() / 1000;
    const {
      percentComplete,
      adherence,
      blockIndex,
      microCycleIndex,
      workoutIndex,
    } = useProgramInstanceStatus(startDate);

    const programLength = data.numBlocks * data.numCycles * data.cycleSpan;
    expect(Math.floor(percentComplete.value)).toEqual(
      Math.floor((started / programLength) * 100)
    );
    expect(Math.floor(adherence.value)).toEqual(
      Math.floor((data.numCompleted / started) * 100)
    );
    const blockLength = data.cycleSpan * data.numCycles;
    expect(blockIndex.value).toEqual(Math.floor(started / blockLength));
    expect(microCycleIndex.value).toEqual(
      Math.floor((started - blockLength * blockIndex.value) / data.cycleSpan)
    );
    expect(workoutIndex.value).toEqual(
      started -
        blockLength * blockIndex.value -
        microCycleIndex.value * data.cycleSpan -
        1
    );
  });
});
