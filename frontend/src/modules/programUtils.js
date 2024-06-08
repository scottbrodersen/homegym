import { programInstanceStore } from '../modules/state';
import * as dateUtils from './../modules/dateUtils';
import * as utils from './utils';

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

// Returns the planned workouts of the active instance
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

// Returns the dates of all planned non-rest day workouts
export const getNonRestDates = (instance) => {
  const dates = new Array();

  let date = dateUtils.dateFromSeconds(instance.startDate);

  // set to day before the start date
  date.setDate(date.getDate() - 1);

  // walk the program
  instance.blocks.forEach((block) => {
    block.microCycles.forEach((cycle) => {
      cycle.workouts.forEach((workout) => {
        date.setDate(date.getDate() + 1);
        if (!workout.restDay) {
          dates.push(dateUtils.formatDate(date));
        }
      });
    });
  });

  return dates;
};

export const getInstanceWorkoutDates = (instance) => {
  let startDate = dateUtils.dateFromSeconds(instance.startDate);
  const numDays = getProgramLength(instance);
  const dates = new Array();

  for (let i = 0; i < numDays; i++) {
    let date = new Date(startDate.getTime());
    date.setDate(startDate.getDate() + i);
    dates.push(dateUtils.formatDate(date));
  }

  return dates;
};

export const getWorkoutDate = (instance, coords) => {
  const dayIndex = getDayIndex(instance, coords);
  const workoutDate = dateUtils.dateFromSeconds(instance.startDate);
  workoutDate.setDate(workoutDate.getDate() + dayIndex);
  return workoutDate;
};

const getEventsForDate = async (activityID, date) => {
  const eventsForDate = new Array();

  // Find the limits of the date window for the workouts we want to get
  const midnight = new Date(date);
  midnight.setHours(0);
  midnight.setMinutes(0);
  midnight.setSeconds(0);

  // the end of the day on the event date is midnight of the next day
  const nextDay = new Date(midnight);

  nextDay.setDate(nextDay.getDate() + 1);

  // Get page of events that occur earlier than nextDay
  const eventPage = await utils.fetchEventPage(
    null,
    Math.floor(nextDay.valueOf() / 1000)
  );

  for (let i = 0; i < eventPage.length; i++) {
    if ((eventPage[i].activityID = activityID)) {
      if (eventPage[i].date > midnight.valueOf() / 1000) {
        eventsForDate.push(eventPage[i]);
      } else {
        // remaining event dates are earlier than midnight
        break;
      }
    }
  }
  return eventsForDate;
};

export const getEventsOnWorkoutDay = async (instance, coords) => {
  const workoutDate = getWorkoutDate(instance, coords);
  return await getEventsForDate(instance.activityID, workoutDate);
};

export const getDayIndex = (program, coords) => {
  let day = 0;
  // add days for previous  blocks
  for (let i = 0; i < coords[0]; i++) {
    const block = program.blocks[i];
    for (let j = 0; j < block.microCycles.length; j++) {
      day += block.microCycles[j].span;
    }
  }

  // get days for previous cycles in the current block
  for (let i = 0; i < coords[1]; i++) {
    day += program.blocks[coords[0]].microCycles[i].span;
  }

  // add previous days in current microcycle
  day += coords[2];

  return day;
};

export const getTodayIndex = (instance) => {
  const now = new Date().valueOf();

  const startDate = dateUtils.dateFromSeconds(instance.startDate);

  // first day is 0th day
  return Math.floor((now - startDate.valueOf()) / 1000 / 60 / 60 / 24);
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

export const getProgramLength = (program) => {
  let progLength = 0;
  program.blocks.forEach((block) => {
    block.microCycles.forEach((ms) => {
      progLength += ms.span;
    });
  });

  return progLength;
};

export const getProgramInstanceStatus = (instanceID) => {
  let percentComplete = 0;
  let adherence = 0;
  let coords = [];

  const instance = programInstanceStore.get(instanceID);
  const dayIndex = getTodayIndex(instance);
  const progLength = getProgramLength(instance);

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

  if (dayIndex > progLength) {
    percentComplete = 100;
    adherence = Math.floor((numPerformed / (progLength - restDaysSoFar)) * 100);
    dayIndex = undefined;
  } else {
    percentComplete = Math.floor((dayIndex / progLength) * 100);
    adherence = Math.floor(
      (numPerformed / (dayIndex + 1 - restDaysSoFar)) * 100
    );
    coords = getWorkoutCoords(instance, dayIndex);
  }

  return [percentComplete, adherence, coords, dayIndex];
};

export const requiredFieldValidator = (val) => {
  const result = (val && val.length > 0) || 'Required value.';
  return result;
};

export const maxFieldValidator = (val) => {
  const result = (val ? val.length < 256 : true) || 'Max 255 characters.';
  return result;
};

export const programValidator = (program) => {
  let noProps = true;
  for (const prop in program) {
    if (Object.hasOwn(program, prop)) {
      noProps = false;
      break;
    }
  }
  if (noProps) {
    return false;
  }

  if (
    requiredFieldValidator(program.title) !== true &&
    maxFieldValidator(program.title) !== true
  ) {
    return false;
  }
  program.blocks.forEach((block) => {
    if (
      requiredFieldValidator(block.title) !== true &&
      maxFieldValidator(block.title) !== true
    ) {
      return false;
    }
    block.microCycles.forEach((cycle) => {
      if (
        requiredFieldValidator(cycle.title) !== true &&
        maxFieldValidator(cycle.title) !== true
      ) {
        return false;
      }
      cycle.workouts.forEach((workout) => {
        if (
          requiredFieldValidator(workout.title) !== true &&
          maxFieldValidator(workout.title) !== true
        ) {
          return false;
        }
        if (!workout.segments) {
          workout['segments'] = [];
        }
        workout.segments.forEach((segment) => {
          if (requiredFieldValidator(segment.exerciseTypeID) !== true) {
            return false;
          }
          if (
            requiredFieldValidator(segment.prescription) !== true &&
            maxFieldValidator(segment.prescription !== true)
          ) {
            return false;
          }
        });
      });
    });
  });
  return true;
};
