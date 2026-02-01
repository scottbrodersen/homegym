import { programInstanceStore } from '../modules/state';
import * as dateUtils from './../modules/dateUtils';
import * as utils from './utils';
import { Dialog } from 'quasar';
import WorkoutModal from '../components/WorkoutModal.vue';

// Valid values for a workout status
export const workoutStatuses = { FUTURE: 0, MISSED: 1, DONE: 2 };

/**
 * Associates the name of a colour style with each workout status value.
 * @param {String} status A workoutStatuses value.
 * @returns A String that stores the name of the colour style.
 */
export const getColorStyle = (status) => {
  if (status == workoutStatuses.FUTURE) {
    return 'futureWorkout';
  } else if (status == workoutStatuses.DONE) {
    return 'doneWorkout';
  }
  return 'missedWorkout';
};

/**
 * Determines the status of a workout event based on the properties of the event.
 * @param {Striing} eventID The ID of the event. No ID indicates the event has not occurred.
 * @param {Number} workoutIndex The index of the workout in the context of its program instance.
 * @param {Number} todayIndex The index of the current date in the context of the program instance.
 * @param {Boolean} isRestDay A value of true indicates that the event is a rest day, false otherwise.
 * @returns A workoutStatuses value.
 */
export const getWorkoutStatus = (
  eventID,
  workoutIndex,
  todayIndex,
  isRestDay,
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

/**
 * Determines the name and colour of an icon to associate with a workoutStatuses value.
 * @param {*} workoutStatus
 * @returns An object with properties name (the icon name) and colour (the icon colour).
 */
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

/**
 * Returns the planned workouts of the current program instance for an activity.
 * @param {String} activityID The ID of the activity.
 * @returns An array of workout objects in order of occurrence.
 */
export const getWorkouts = (activityID) => {
  const currentInstance = programInstanceStore.getCurrent(activityID);

  // walk the program to retrieve workouts
  let workouts = new Array();

  currentInstance.blocks.forEach((block) => {
    block.microCycles.forEach((cycle) => {
      workouts = workouts.concat(cycle.workouts);
    });
  });

  return workouts;
};

/**
 * Returns the dates of all planned non-rest day workouts for a program instance.
 * @param {Object} instance The program instance.
 * @returns An array of date values formatted for use with the Quasar QCalendar components.
 */
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

/**
 * Retrieves all of the workout dates for a program instance.
 * @param {Object} instance The program instance object.
 * @returns An array of dates in order of occurrence, formatted for use with the Quasar QCalender component.
 */
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

/**
 * Determines the date of a planned workout.
 * @param {Object} instance The program instance that the workout belongs to.
 * @param {[Number]} coords The coordinates of the workout, as [block_index, microcycle_index, workout_index].
 * @returns The workout date.
 */
export const getWorkoutDate = (instance, coords) => {
  const dayIndex = getDayIndex(instance, coords);
  const workoutDate = dateUtils.dateFromSeconds(instance.startDate);
  workoutDate.setDate(workoutDate.getDate() + dayIndex);
  return workoutDate;
};

/**
 * Determines the workout events that occurred on a date.
 * @param {String} activityID
 * @param {Number} date The workout date, in milliseconds since epoch.
 * @returns An array of event objects.
 */
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
  const eventPage = await utils.fetchEvents(
    null,
    Math.floor(nextDay.valueOf() / 1000),
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

/**
 * Determines the workout events that occurred on a day of a program instance.
 * @param {Object} instance The program instance.
 * @param {[Number]} coords The coordinates of the workout day, as [block_index, microcycle_index, workout_index].
 * @returns An array of event objects.
 */
export const getEventsOnWorkoutDay = async (instance, coords) => {
  const workoutDate = getWorkoutDate(instance, coords);
  return await getEventsForDate(instance.activityID, workoutDate);
};

/**
 * Determines the index of a day in a program or program instance.
 * @param {Object} program The program or program instance.
 * @param {[Number]} coords The coordinates of the day, as [block_index, microcycle_index, workout_index].
 * @returns The index (0-based) of the coordinates.
 */
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

/**
 * Determines the index of the day at the time of execution within a program instance.
 * @param {Object} instance The program instance.
 * @returns The index (0-based) of now.
 */
export const getTodayIndex = (instance) => {
  const now = new Date().valueOf();

  const startDate = dateUtils.dateFromSeconds(instance.startDate);

  // first day is 0th day
  return Math.floor((now - startDate.valueOf()) / 1000 / 60 / 60 / 24);
};

/**
 * Determines the coordinates for a day of a program or program instance
 * @param {Object} program The program or program instance.
 * @param {Number} programDayIndex The index of the day.
 * @returns The coordinates of the day, as [block_index, microcycle_index, workout_index].
 */
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

/**
 * Determines the number of days in a program.
 * @param {Object} program The program.
 * @returns The number of days in the program.
 */
export const getProgramLength = (program) => {
  let progLength = 0;
  program.blocks.forEach((block) => {
    block.microCycles.forEach((ms) => {
      progLength += ms.span;
    });
  });

  return progLength;
};

/**
 * Determines various status metrics for a program instance.
 * @param {String} instanceID The ID of the program instance.
 * @returns An object with properties percentComplete (the percentage of workout days that are in the past),
 *  adherence (the percentage of planned workouts that have been performed so far),
    coords (the coordinates of today's workout)
    and dayIndex (the index of today's workout).
 */
export const getProgramInstanceStatus = (instanceID) => {
  let percentComplete = 0;
  let adherence = 0;

  /*
  Coords is a 3x1 array that holds the coordinates of the workout for a date.
  E.g. [0,1,2] denotes the workout in the 3rd day of the 2nd microcycle in the 1st block.
  */
  let coords = [];

  const instance = programInstanceStore.get(instanceID);
  let dayIndex = getTodayIndex(instance);
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
      (numPerformed / (dayIndex + 1 - restDaysSoFar)) * 100,
    );
    coords = getWorkoutCoords(instance, dayIndex);
  }

  return {
    percentComplete: percentComplete,
    adherence: adherence,
    coords: coords,
    dayIndex: dayIndex,
  };
};

/**
 * Determines whether a value is not null.
 * @param {*} val The value.
 * @returns True if not null, otherwise false.
 */
export const requiredFieldValidator = (val) => {
  const result = (val && val.length > 0) || 'Required value.';
  return result;
};

/**
 * Determines whether the length of a value exceeds 255 characters.
 * @param {String} val The value.
 * @returns True when the value is less than 255 characters, otherwise false.
 */
export const maxFieldValidator = (val) => {
  const result = (val ? val.length < 256 : true) || 'Max 255 characters.';
  return result;
};

/**
 * Validates a program
 * @param {Object} program The program.
 * @returns True when the program is valid, otherwise false.
 */
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
    !(
      requiredFieldValidator(program.title) === true &&
      maxFieldValidator(program.title) === true
    )
  ) {
    return false;
  }
  for (let i = 0; i < program.blocks.length; i++) {
    if (
      !(
        requiredFieldValidator(program.blocks[i].title) === true &&
        maxFieldValidator(program.blocks[i].title) === true
      )
    ) {
      return false;
    }
    for (let j = 0; j < program.blocks[i].length; j++) {
      if (
        !(
          requiredFieldValidator(program.blocks[i].microCycles[j].title) ===
            true &&
          maxFieldValidator(program.blocks[i].microCycles[j].title) === true
        )
      ) {
        return false;
      }
      cycle.workouts.forEach((workout) => {
        if (!workoutValidator(workout)) {
          return false;
        }
      });
    }
  }
  return true;
};

/**
 * Validates a workout event.
 * @param {Object} workout The workout.
 * @returns True when the workout is valid, otherwise false.
 */
export const workoutValidator = (workout) => {
  if (
    !(
      requiredFieldValidator(workout.title) === true &&
      maxFieldValidator(workout.title) === true
    )
  ) {
    return false;
  }

  // If it's a rest day, we're done
  if (workout.restDay) {
    return true;
  }

  // No segments is ok but make it an empty array
  if (!workout.segments) {
    workout['segments'] = [];
  }

  for (let i = 0; i < workout.segments.length; i++) {
    if (requiredFieldValidator(workout.segments[i].exerciseTypeID) !== true) {
      return false;
    }
    if (
      !(
        requiredFieldValidator(workout.segments[i].prescription) === true &&
        maxFieldValidator(workout.segments[i].prescription === true)
      )
    ) {
      return false;
    }
  }

  return true;
};

/**
 * A dialog for entering the details of a workout event.
 * @param {Object} instance The program instance to which the workout belongs.
 * @param {[Number]} coords The coordinates of the workout that is being performed, as [block_index, microcycle_index, workout_index].
 * @returns The workout event object.
 */
export const newWorkoutModal = (instance, coords) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: WorkoutModal,
      componentProps: { instance: instance, coords: coords },
    })
      .onOk((event) => {
        resolve(event);
      })
      .onCancel(() => {
        resolve();
      })
      .onDismiss(() => {});
  });
};

/**
 * Finds the current active program instance based on date. The instance with the earliest start date is current.
 * @param {String} activityID The activity with which the instances are associated.
 * @returns The current program instance.
 */
export const selectCurrentProgramInstance = (activityID) => {
  let active = null;

  // Get the instances from the store
  let instances = [];
  let activeInstances = programInstanceStore.getActive(activityID);
  if (activeInstances) {
    for (let i = 0; i < activeInstances.length; i++) {
      instances.push(
        programInstanceStore.get(
          activeInstances[i].id,
          activeInstances[i].programID,
        ),
      );
    }
    let earliest = null;
    for (let i = 0; i < activeInstances.length; i++) {
      if (i == 0) {
        earliest = instances[i].startDate;
        active = instances[i];
      } else if (instances[i].startDate < earliest) {
        earliest = instances[i].startDate;
        active = instances[i];
      }
    }
  }
  // check if the start date is in the future
  if (active && dateUtils.nowInSeconds() < active.startDate) {
    return null;
  }
  return active;
};
