import { programInstanceStore } from '../modules/state';

export const workoutStatuses = { FUTURE: 0, MISSED: 1, DONE: 2 };
export const getColorStyle = (status) => {
  if (status == workoutStatuses.FUTURE) {
    return 'futureWorkout';
  } else if (status == workoutStatuses.DONE) {
    return 'doneWorkout';
  }
  return 'missedWorkout';
};

export const getWorkoutStatus = (
  eventID,
  workoutIndex,
  todayIndex,
  isRestDay
) => {
  if (eventID && eventID != '') {
    return workoutStatuses.DONE;
  } else if (workoutIndex >= todayIndex) {
    return workoutStatuses.FUTURE;
  } else if (isRestDay) {
    return workoutStatuses.DONE;
  } else if (eventID == '') {
    return workoutStatuses.MISSED;
  } else if (!eventID) {
    return workoutStatuses.MISSED;
  }
};

export const getStatusIcons = (workoutStatus) => {
  if (workoutStatus == workoutStatuses.DONE) {
    return { name: 'check_circle', colour: 'green' };
  } else if (workoutStatus == workoutStatuses.MISSED) {
    return { name: 'cancel', colour: 'red' };
  } else if (workoutStatus == workoutStatuses.FUTURE) {
    return { name: '', colour: 'white' };
  }
  return { name: 'yard', colour: 'yellow' };
};

export const getWorkouts = (activityID) => {
  const activeInstance = programInstanceStore.getActive(activityID);

  // walk the program to retrieve workouts
  let workouts = new Array();

  activeInstance.blocks.forEach((block) => {
    block.microCycles.forEach((cycle) => {
      workouts = workouts.concat(cycle.workouts);
    });
  });

  return workouts;
};

export const getWorkoutCoords = (program, programDayIndex) => {
  let blockIndex;
  let microCycleIndex;
  let workoutIndex;

  let numDays = 0;

  // walk the program
  let done = false;
  for (let i = 0; i < program.blocks.length; i++) {
    for (let j = 0; j < program.blocks[i].microCycles.length; j++) {
      const numDaysInMicroCycle = program.blocks[i].microCycles[j].span;
      const lastDayIndexOfMicroCycle = numDays + numDaysInMicroCycle - 1;
      if (lastDayIndexOfMicroCycle < programDayIndex) {
        numDays += numDaysInMicroCycle;
        continue;
      }
      blockIndex = i;
      microCycleIndex = j;
      workoutIndex = programDayIndex - numDays;
      done = true;
      break;
    }
    if (done) {
      break;
    }
  }

  return [blockIndex, microCycleIndex, workoutIndex];
};

export const getProgramInstanceStatus = (instanceID) => {
  let percentComplete = 0;
  let adherence = 0;
  let coords = [];
  let dayIndex = 0;

  const instance = programInstanceStore.get(instanceID);

  const nowDate = new Date();

  // set to midnight
  nowDate.setUTCHours(0);
  nowDate.setUTCMinutes(0);
  nowDate.setUTCSeconds(0);
  nowDate.setUTCMilliseconds(0);

  const now = nowDate.valueOf();

  // javascript timestamp is in msec
  const startDate = new Date(instance.startDate * 1000);

  // make sure start date is at midnight
  startDate.setUTCHours(0);
  startDate.setUTCMinutes(0);
  startDate.setUTCSeconds(0);
  startDate.setUTCMilliseconds(0);

  // first day is 0th day
  dayIndex = Math.floor((now - startDate.valueOf()) / 1000 / 60 / 60 / 24);

  let progLength = 0;
  instance.blocks.forEach((block) => {
    block.microCycles.forEach((ms) => {
      progLength += ms.span;
    });
  });

  let restDaysSoFar = 0;
  let dayCount = 0;
  for (let i = 0; i < instance.blocks.length; i++) {
    for (let j = 0; j < instance.blocks[i].microCycles.length; j++) {
      for (
        let k = 0;
        k < instance.blocks[i].microCycles[j].workouts.length;
        k++
      ) {
        dayCount++;
        if (dayCount > dayIndex) {
          break;
        }
        if (instance.blocks[i].microCycles[j].workouts[k].restDay) {
          restDaysSoFar++;
        }
      }
    }
  }

  let numPerformed = 0;
  if (instance.events) {
    // A workout is performed if the events object has a key for the program day && the value is an event ID
    for (let i = 0; i < dayIndex + 1; i++) {
      if (instance.events[i] != undefined && instance.events[i] != '') {
        numPerformed++;
      }
    }
  }

  percentComplete = Math.floor((dayIndex / progLength) * 100);
  adherence = Math.floor((numPerformed / (dayIndex + 1 - restDaysSoFar)) * 100);

  coords = getWorkoutCoords(instance, dayIndex);

  return [percentComplete, adherence, coords, dayIndex];
};
