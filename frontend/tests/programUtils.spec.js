import './patch-fetch.js';
import { expect, beforeAll, afterAll, afterEach, vi } from 'vitest';
import * as utils from '../src/modules/programUtils';
import * as data from '../mocks/data';
import server from '../mocks/server.js';

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
    const coords = utils.getWorkoutCoords(testProgram, 0);
    expect(coords).toEqual(expected);
  });
  it('gets the workout coordinates -- mid first cycle', () => {
    const expected = [0, 0, 2];
    const coords = utils.getWorkoutCoords(testProgram, 2);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- last day of first cycle', () => {
    const expected = [0, 0, 6];
    const coords = utils.getWorkoutCoords(testProgram, 6);
    expect(coords).toEqual(expected);
  });
  it('gets the workout coordinates -- first day of second cycle', () => {
    const expected = [0, 1, 0];
    const coords = utils.getWorkoutCoords(testProgram, 7);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- first day of second block', () => {
    const expected = [1, 0, 0];
    const coords = utils.getWorkoutCoords(testProgram, 14);
    expect(coords).toEqual(expected);
  });

  it('gets the workout coordinates -- middle of second cycle of second block', () => {
    const expected = [1, 1, 2];
    const coords = utils.getWorkoutCoords(testProgram, 23);
    expect(coords).toEqual(expected);
  });

  it('gets the day index -- first day', () => {
    const expected = 0;
    const day = utils.getDayIndex(testProgram, [0, 0, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first block', () => {
    const expected = 1;
    const day = utils.getDayIndex(testProgram, [0, 0, 1]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first day of second block', () => {
    const expected = data.cycleSpan * data.numCycles;
    const day = utils.getDayIndex(testProgram, [1, 0, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- third day of second block', () => {
    const expected = data.cycleSpan * data.numCycles + 2;
    const day = utils.getDayIndex(testProgram, [1, 0, 2]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- third day of second cycle of second block', () => {
    const expected = data.cycleSpan * (data.numCycles + 1) + 2;
    const day = utils.getDayIndex(testProgram, [1, 1, 2]);
    expect(day).toEqual(expected);
  });

  it('gets the day index -- first day of second cycle of second block', () => {
    const expected = data.cycleSpan * (data.numCycles + 1);
    const day = utils.getDayIndex(testProgram, [1, 1, 0]);
    expect(day).toEqual(expected);
  });

  it('gets the workouts for a date', async () => {
    const testProgramInstance = data.testProgramInstance(data.testDate);
    const events = await utils.getEventsOnWorkoutDay(
      testProgramInstance,
      [0, 0, 0]
    );
    expect(events.length).toEqual(1);
  });
});
