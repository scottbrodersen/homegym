import { expect } from 'vitest';
import * as state from '../src/modules/state';
import * as testData from '../mocks/data';

describe('state', () => {
  it('sets active program instance', () => {
    const instance = testData.testProgramInstance(testData.testActivityID);
    state.programInstanceStore.add(instance);
    state.programInstanceStore.setCurrent(instance.activityID, instance);
    state.programInstanceStore.addAllActive(testData.testActivityID, [
      instance,
    ]);

    const expectCurrent = {
      programID: instance.programID,
      instanceID: instance.id,
    };

    expect(state.programInstanceStore.activeInstances).toHaveLength(1);
    expect(
      state.programInstanceStore.activeInstances.has(testData.testActivityID)
    ).toEqual(true);
    expect(
      state.programInstanceStore.activeInstances.get(instance.activityID)
    ).toEqual([expectCurrent]);

    expect(
      state.programInstanceStore.getActive(testData.testActivityID)
    ).toEqual([instance]);

    expect(
      state.programInstanceStore.getCurrent(testData.testActivityID)
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
    state.programInstanceStore.addAllActive(testData.testActivityID, [
      instance,
    ]);
    state.programInstanceStore.setCurrent(testData.testActivityID, instance);

    expect(state.programInstanceStore.activeInstances).toHaveLength(1);
    expect(
      state.programInstanceStore.activeInstances.has(testData.testActivityID)
    ).toEqual(true);
    expect(
      state.programInstanceStore.activeInstances.get(testData.testActivityID)
    ).toEqual([]);
    expect(
      state.programInstanceStore.getCurrent(testData.testActivityID)
    ).toEqual(instance);
  });
});
