import './patch-fetch.js';
import { expect, vi } from 'vitest';
import * as utils from '../src/modules/programUtils';
import * as data from '../mocks/data';

vi.mock('../src/modules/state');

const numOutstanding = 3;

const testProgram = data.testProgram();

describe('programUtils', () => {
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
    const expected = data.cycleSpan;
    const day = utils.getDayIndex(testProgram, [1, 0, 0]);
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

  it('gets the day index -- first day of first cycle of second block', () => {
    const expected = data.cycleSpan;
    const day = utils.getDayIndex(testProgram, [1, 1, 0]);
    expect(day).toEqual(expected);
  });

  it('returns the outstanding workouts', () => {
    const startDate =
      Date.now() - (data.numCompleted + numOutstanding) * 24 * 60 * 60 * 1000;
    const dayIndex = data.numCompleted + numOutstanding - 1; // 14 + 3 -1 = 16
    const testProgramInstance = data.testProgramInstance(startDate);

    const expectedOutstanding = [
      testProgramInstance.blocks[1].microCycles[0].workouts[0],
      testProgramInstance.blocks[1].microCycles[0].workouts[1],
      testProgramInstance.blocks[1].microCycles[0].workouts[2],
    ];

    const outstanding = utils.getOutstandingWorkouts(
      data.testActivityID,
      [1, 0, 2],
      16
    );

    expect(outstanding).toEqual(expectedOutstanding);
  });
});
