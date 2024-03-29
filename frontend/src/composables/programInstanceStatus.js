import { ref } from 'vue';
import { programInstanceStore } from '../modules/state';

export const useProgramInstanceStatus = (instanceID) => {
  const percentComplete = ref(0);
  const adherence = ref(0);
  const blockIndex = ref(0);
  const microCycleIndex = ref(0);
  const workoutIndex = ref(0);

  const instance = programInstanceStore.get(instanceID);
  /* determine the current status:
    - next workout
    - % complete = (total number of workouts - future workouts)/(tot number)
    - % adherence = (tot previous planned - tot skipped) / (tot planned)

    use progInst.startDate to determine next workout
    - first workout is block[0].microcycle[0].workouts[0]
    - get day of first workout from progInst.startDate
    - calc number of days since startDate since today
    - walk the blocks to get next workout

    - workoutNumber = (blockIndex+1)*microCycleIndex*mcDays + workoutIndex
  */
  const now = new Date().valueOf();
  // javascript epoch is in milliseconds
  const epoch = instance.startDate * 1000;

  const daysIntoProgram = Math.floor((now - epoch) / 1000 / 60 / 60 / 24);

  let progLength = 0;
  instance.blocks.forEach((block) => {
    block.microCycles.forEach((ms) => {
      progLength += ms.span;
    });
  });

  let numPerformed = 0;
  if (instance.events) {
    // A workout is performed if the events map has a key for the program number && the value not undefined
    for (let i = 0; i < daysIntoProgram; i++) {
      if (instance.events.get(i)) {
        numPerformed++;
      }
    }
  }
  percentComplete.value = Math.floor((daysIntoProgram / progLength) * 100);
  adherence.value =
    daysIntoProgram > 0
      ? Math.floor((numPerformed / daysIntoProgram) * 100)
      : 0;

  let numWorkouts = 0;

  for (let i = 0; i < instance.blocks.length; i++) {
    for (let j = 0; j < instance.blocks[i].microCycles.length; j++) {
      const numWorkoutsInMicroCycle = instance.blocks[i].microCycles[j].span;
      const lastWorkoutInMicroCycle = numWorkouts + numWorkoutsInMicroCycle;
      if (lastWorkoutInMicroCycle < daysIntoProgram) {
        numWorkouts += numWorkoutsInMicroCycle;
        continue;
      }
      blockIndex.value = i - 1;
      microCycleIndex.value = j;
      workoutIndex.value = daysIntoProgram - numWorkouts - 1;
      break;
    }
  }
  return {
    percentComplete,
    adherence,
    blockIndex,
    microCycleIndex,
    workoutIndex,
  };
};
