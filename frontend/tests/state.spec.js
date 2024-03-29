import { expect } from 'vitest';
import * as state from '../src/modules/state';
import * as testData from '../mocks/data';

describe('state', () => {
  it('sets active program instance', () => {
    const instance = testData.testProgramInstance(testData.testActivityID);

    state.programInstanceStore.setActive(instance.activityID, instance);

    const expectActive = {
      programID: instance.programID,
      instanceID: instance.id,
    };

    expect(state.programInstanceStore.activeInstances).toHaveLength(1);
    expect(
      state.programInstanceStore.activeInstances.has(instance.activityID)
    ).toEqual(true);
    expect(
      state.programInstanceStore.activeInstances.get(instance.activityID)
    ).toEqual(expectActive);
    console.log(state.programInstanceStore.programInstances);
    expect(
      state.programInstanceStore.getActive(testData.testActivityID)
    ).toEqual(instance);
    expect(
      state.programInstanceStore.get(instance.id, instance.programID)
    ).toEqual(instance);
    let inst = state.programInstanceStore.get(instance.id);
    expect(inst).toEqual(instance);
    inst = state.programInstanceStore.get('bad-id');
    expect(inst).toBeUndefined;
  });
  it('sets null active program instance', () => {
    const instance = null;

    state.programInstanceStore.setActive(testData.testActivityID, instance);

    const expectActive = null;

    expect(state.programInstanceStore.activeInstances).toHaveLength(1);
    expect(
      state.programInstanceStore.activeInstances.has(testData.testActivityID)
    ).toEqual(true);
    expect(
      state.programInstanceStore.activeInstances.get(testData.testActivityID)
    ).toEqual(expectActive);
    expect(
      state.programInstanceStore.getActive(testData.testActivityID)
    ).toEqual(instance);
  });
});
