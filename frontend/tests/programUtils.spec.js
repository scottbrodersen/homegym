import './patch-fetch.js';
import { expect, beforeAll, afterAll, afterEach, vi } from 'vitest';
import * as programUtils from '../src/modules/programUtils';
import * as data from '../mocks/data';
import server from '../mocks/server.js';
import * as state from '../src/modules/state.js';

vi.mock('../src/modules/state');
vi.mock('../src/modules/utils');

const testProgram = data.testProgram();

describe('programUtils', () => {
  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' });
  });
  afterAll(() => server.close());
  afterEach(() => server.resetHandlers());

  it('gets the workout coordinates -- first day', () => {
    const expected = [0, 0, 0];
    const coords = programUtils.getWorkoutCoords(testProgram, 0);
    expect(coords).toEqual(expected);
  });
  it('gets the workout coordinates -- mid first cycle', () => {
    const expected = [0, 0, 2];
    const coords = programUtils.getWorkoutCoords(testProgram, 2);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- last day of first cycle', () => {
    const expected = [0, 0, 6];
    const coords = programUtils.getWorkoutCoords(testProgram, 6);
    expect(coords).toEqual(expected);
  });
  it('gets the workout coordinates -- first day of second cycle', () => {
    const expected = [0, 1, 0];
    const coords = programUtils.getWorkoutCoords(testProgram, 7);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- first day of second block', () => {
    const expected = [1, 0, 0];
    const coords = programUtils.getWorkoutCoords(testProgram, 14);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- middle of second cycle of second block', () => {
    const expected = [1, 1, 2];
    const coords = programUtils.getWorkoutCoords(testProgram, 23);
    expect(coords).toEqual(expected);
  });

  it('gets the day index -- first day', () => {
    const expected = 0;
    const day = programUtils.getDayIndex(testProgram, [0, 0, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first block', () => {
    const expected = 1;
    const day = programUtils.getDayIndex(testProgram, [0, 0, 1]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first day of second block', () => {
    const expected = data.cycleSpan * data.numCycles;
    const day = programUtils.getDayIndex(testProgram, [1, 0, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- third day of second block', () => {
    const expected = data.cycleSpan * data.numCycles + 2;
    const day = programUtils.getDayIndex(testProgram, [1, 0, 2]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- third day of second cycle of second block', () => {
    const expected = data.cycleSpan * (data.numCycles + 1) + 2;
    const day = programUtils.getDayIndex(testProgram, [1, 1, 2]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first day of second cycle of second block', () => {
    const expected = data.cycleSpan * (data.numCycles + 1);
    const day = programUtils.getDayIndex(testProgram, [1, 1, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the workouts for a date', async () => {
    const testProgramInstance = data.testProgramInstance(data.testDate);
    const events = await programUtils.getEventsOnWorkoutDay(
      testProgramInstance,
      [0, 0, 0]
    );
    expect(events.length).toEqual(1);
  });

  it('determines the current program instance', async () => {
    const testProgramInstance1 = data.testProgramInstance(data.testDate + 1);
    const testProgramInstance2 = data.testProgramInstance(data.testDate - 1);
    const testProgramInstance3 = data.testProgramInstance(data.testDate);

    testProgramInstance2.id = testProgramInstance2.id + '2';
    testProgramInstance3.id = testProgramInstance3.id + '3';

    const all = [
      testProgramInstance1,
      testProgramInstance2,
      testProgramInstance3,
    ];

    state.programInstanceStore.addAllActive(
      testProgramInstance1.activityID,
      all
    );

    expect(
      programUtils.selectCurrentProgramInstance(testProgramInstance1.activityID)
    ).toEqual(testProgramInstance2);
  });
});
