import {
  eventStore,
  activityStore,
  exerciseTypeStore,
  loginModalState,
  programsStore,
  programInstanceStore,
} from './state';
import LoginModal from './../components/LoginModal.vue';
import { Dialog, Notify } from 'quasar';
import NewActivityModal from './../components/NewActivityModal.vue';
import VolumeModal from './../components/VolumeModal.vue';
import CompositionModal from './../components/CompositionModal.vue';
import VariationModal from './../components/VariationModal.vue';
import NewProgramModal from '../components/NewProgramModal.vue';
import ConfirmModal from './../components/ConfirmModal.vue';
import ProgramInstanceModal from './../components/ProgramInstanceModal.vue';
import { toRaw, isRef, isReactive, isProxy, unref } from 'vue';
import EditValueModal from '../components/EditValueModal.vue';
import ProgramModal from '../components/ProgramModal.vue';

/**
 * Determines the number of rows of events to include on the home page.
 * @returns
 */
const pageSize = () => {
  const defaultPage = 8;
  try {
    const availableHeight = window.innerHeight - 220 - 36;
    page = Math.floor((availableHeight - 48 - 50 - 40 - 30) / 42);
    if (page < 1) {
      return defaultPage;
    }
  } catch (e) {
    return defaultPage;
  }
};

// Page size when retrieving paged data from the back end.
const fetchPageSize = 100;

// Names of the types of intensities used to describe exercise performances.
const intensityTypes = [
  'weight',
  'bodyweight',
  'rpe',
  'percentOfMax',
  'hrZone',
  'distance',
  'pace',
];

/**
 * Returns an object of string and function properties for interpreting/handling values of intensities.
 * The intensity type determines how the object interprets/handles the values.
 * Note that intensity types are stored in the same format in the backend, regardless of type.
 * @param {*} intensityType The name of the intensity type.
 * @returns
 */
const intensityProps = (intensityType) => {
  if (intensityType == 'hrZone') {
    return {
      mask: '#',
      validate: (value) => {
        return /^[1-5]$/.test(value);
      },
      format: (value) => value.toFixed(0),
      prefix: 'HR Zone',
      value: (formatted) => Number(formatted).toFixed(0),
    };
  } else if (intensityType == 'rpe') {
    return {
      mask: '##',
      validate: (value) => {
        return /^[1-9]$/.test(value) || /^10$/.test(value);
      },
      format: (value) => value.toFixed(1),
      prefix: 'RPE (1-10)',
      value: (formatted) => Number(formatted).toFixed(0),
    };
  } else if (intensityType == 'pace') {
    return {
      mask: '##:##',
      format: (value) => {
        let minutes = `${(value / 60).toFixed(0)}`;
        let seconds = `${value % 60}`;
        while (minutes.length < 2) {
          minutes = '0' + minutes;
        }
        while (seconds.length < 2) {
          seconds = '0' + seconds;
        }
        return `${minutes}:${seconds}`;
      },
      value: (formatted) => {
        const parts = formatted.split(':');
        return Number(parts[0]) * 60 + Number(parts[1]);
      },
      prefix: '',
      validate: (value) => {
        return /^[0-9]{1,2}:[0-9]{2}$/.test(value);
      },
    };
  } else if (intensityType == 'bodyweight') {
    return {
      format: (value) => {
        return 'bodyweight';
      },
      value: (formatted) => {
        return '1';
      },
      validate: (value) => {
        return true;
      },
    };
  } else {
    return {
      mask: '',
      validate: (value) => {
        if (value.endsWith('.')) {
          return false;
        }
        return /^[0-9]+\.?[0-9]?$/.test(value);
      },
      format: (value) => value.toFixed(1),
      prefix: '',
      value: (formatted) => Number(formatted).toFixed(1),
    };
  }
};

// Names of the types of volume used to describe exercise performances.
const volumeTypes = ['count', 'time', 'distance'];

/**
 * Retrieves a page of workout events from the back end.
 * @param {String} eventID (Optional) The ID of the last event ID included in the previous page. No value for the first page.
 * @param {Number} date (Optional) The earliest date of the events to include in the page, in ms since epoch.
 * @returns An array of event objects.
 */
const fetchEventPage = async (eventID = '', date = null) => {
  const startTime = date ? date : Math.floor(Date.now() / 1000);
  const params = new URLSearchParams();

  params.append('count', fetchPageSize);

  if (eventID) params.append('previousID', eventID);
  if (startTime) params.append('date', startTime);

  const paramStr = params.toString();
  const url = `/homegym/api/events/?${paramStr}`;
  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of event page');
  }
  const eventPage = await resp.json();
  return eventPage;
};

/**
 * Retrieves workout events from the back end that occurred during a specific period of time.
 * Note that on the backend, events are stored in order of ID.
 * @param {String} eventID (Optional) If you are appending events to a list that you previously retrieved, the ID of the last-retrieved event.
 * @param {Number} startDate (Optional) The beginning of the period of time, in seconds since epoch. No value retrieves the earliest events.
 * @param {Number} endDate (Optional) The end of the period of time, in seconds since epoch. No value retrieves the latest events.
 * @returns
 */
const fetchEvents = async (eventID, startDate, endDate) => {
  let events = new Array();
  let done = false;
  let pageDate = startDate;
  let previousID = eventID;
  while (!done) {
    let page;
    try {
      page = await fetchEventPage(previousID, pageDate);
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
      } else {
        console.log(e);
        throw e;
      }
    }

    // if we had to log in after a 401, don't update the page boundaries
    if (page) {
      eventStore.addBulk(page);

      events = events.concat(page);

      const lastEvent = page[page.length - 1];
      if (
        page.length == 0 ||
        lastEvent.date < endDate ||
        page.length < fetchPageSize
      ) {
        done = true;
      } else {
        pageDate = lastEvent.date;
        previousID = lastEvent.id;
      }
    }
  }
  return events;
};

/**
 * Retrieves the programs from the back end for an activity.
 * Saves the results in the program store.
 * @param {String} activityID The ID of the activity.
 */
const fetchPrograms = async (activityID) => {
  let done = false;
  let lastProgram = '';
  while (!done) {
    const programPage = await fetchProgramPage(lastProgram, activityID);
    if (programPage) {
      programsStore.addBulk(programPage);

      if (programPage && programPage.length < pageSize()) {
        done = true;
      } else {
        lastProgram = programPage[programPage.length - 1].id;
      }
    } else {
      done = true;
    }
  }
};

/**
 * Retrieves a page of programs from the back end for an activity.
 * @param {String} programID The ID of the last program from the previous page. No value for the first page.
 * @param {String} activityID The ID of the activity.
 * @returns An array of program objects.
 */
const fetchProgramPage = async (programID, activityID) => {
  const params = new URLSearchParams();

  params.append('size', fetchPageSize);

  if (programID) params.append('previous', programID);

  const paramStr = params.toString();
  const url = `/homegym/api/activities/${activityID}/programs?${paramStr}`;
  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of program page');
  }
  const programPage = await resp.json();
  return programPage;
};

const fetchProgramInstances = async (programID, activityID) => {
  let done = false;
  let lastInstance = '';
  while (!done) {
    const instancePage = await fetchProgramInstancePage(
      lastInstance,
      programID,
      activityID
    );
    programInstanceStore.addBulk(instancePage);
    if (instancePage.length < pageSize()) {
      done = true;
    } else {
      lastInstance = instancePage[instancePage.length - 1].id;
    }
  }
};

/**
 * Retrieves a page of program instances from the back end for a program of an activity.
 * @param {String} programInstsanceID The ID of the last program instance from the previous page. No value for the first page.
 * @param {String} programID The ID of the program that the instance is based on.
 * @param {String} activityID The ID of the activity.
 * @returns An array of program instance objects.
 */
const fetchProgramInstancePage = async (
  programInstanceID,
  programID,
  activityID
) => {
  const params = new URLSearchParams();

  params.append('size', fetchPageSize);

  if (programInstanceID) params.append('previous', eventID);

  const paramStr = params.toString();
  const url = `/homegym/api/activities/${activityID}/programs/${programID}/instances?${paramStr}`;
  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of program instance page');
  }
  const instancePage = await resp.json();
  return instancePage;
};

/**
 * Retrieves the program instances for all programs of an activity.
 * Stores the results in the program instance store.
 * todo: 404 should indicate the activityID was not found instead of no instances found for the activity
 * @param {*} activityID The ID of the activity
 */
const fetchActiveProgramInstances = async (activityID) => {
  const url = `/homegym/api/activities/${activityID}/programs/instances/active/`;
  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of active program instance');
  } else if (resp.status == 404) {
    // No active instances so set current to null
    //programInstanceStore.setCurrent(activityID, null);
    return;
  }

  const instances = await resp.json();
  // bug: result is array of string json instead of array of objects
  if (instances) {
    programInstanceStore.addAllActive(activityID, instances);
  }
};

/**
 * Retrieves all activities from the back end.
 * Stores the results in the activity store.
 */
const fetchActivities = async () => {
  const resp = await fetch('/homegym/api/activities/', {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of activities');
  }

  const activities = await resp.json();
  activities.forEach(async (a) => {
    activityStore.add(a);
    await fetchActivityExercises(a.id);
  });
};

/**
 * Retrieves all exercise types from the back end.
 * Stores results in the exercise types store.
 */
const fetchExerciseTypes = async () => {
  const resp = await fetch('/homegym/api/exercises/', {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of exercise types');
  }

  const exerciseTypes = await resp.json();
  exerciseTypes.forEach((a) => {
    exerciseTypeStore.add(a);
  });
};

/**
 * Retrieves the exercise types that are associated with an activity.
 * Stores the results in the activity store.
 * @param {String} activityID The ID of the activity.
 */
const fetchActivityExercises = async (activityID) => {
  const resp = await fetch(`/homegym/api/activities/${activityID}/exercises/`, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of activity exercises');
  }
  const exercises = await resp.json();

  const activity = activityStore.get(activityID);
  activity.exercises = exercises;
};

/**
 * Opens a dialog for logging in.
 * To be used when a fetch is sent to the back end and the authentication token is expired.
 * @returns
 */
const authPromptAsync = () => {
  return new Promise((resolve, reject) => {
    if (!loginModalState.isOpen) {
      loginModalState.opened();
      Dialog.create({
        component: LoginModal,
      })
        .onOk(() => {})
        .onCancel(() => {})
        .onDismiss(() => {
          loginModalState.closed();
          resolve();
        });
    } else {
      reject('login prompt already open');
    }
  });
};

/**
 * Opens a dialog that displays a message and enables the user to edit a value.
 * Both the message and value to edit are provided by the caller.
 * @param {Object} valueObjects An object with properties label and value, eg.
 *            {
 *             label: 'Program Title',
 *             value: program.value.title,
 *             }
 * @returns The changed value.
 */
const openEditValueModal = (valueObjects) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: EditValueModal,
      componentProps: {
        values: valueObjects,
      },
    })
      .onOk((newValue) => {
        resolve(newValue);
      })
      .onCancel(() => {})
      .onDismiss(() => {
        resolve();
      });
  });
};

/**
 * A dialog for creating an activity.
 */
const newActivityPrompt = () => {
  Dialog.create({
    component: NewActivityModal,
  })
    .onOk(() => {})
    .onCancel(() => {})
    .onDismiss(() => {});
};

/**
 * A dialog for expressing the volume of a performed exercise.
 * Optionally executes a function that uses the resulting exercise instance as a parameter.
 * @param {*} exerciseTypeID The ID of the exercise
 * @param {*} intensity  (Optional) The intensity to display (e.g. weight, perceived exertion, etc)
 * @param {*} segmentIndex A 0-based index that denotes the place in a list of performed exercises.
 * @param {*} volume The volume performed. Each item in the array represents a set.
 * @param {*} callback (Optional) The function to execute.
 */
const openVolumeModal = (
  exerciseTypeID,
  intensity,
  segmentIndex,
  volume,
  callback
) => {
  Dialog.create({
    component: VolumeModal,
    componentProps: {
      exerciseTypeID: exerciseTypeID,
      intensity: intensity,
      segmentIndex: segmentIndex,
      volume: volume,
    },
  })
    .onOk((perfObj) => {
      callback(perfObj);
    })
    .onCancel(() => {})
    .onDismiss(() => {});
};

/**
 * A dialog that creates a program.
 * @param {String} activityID The ID of the activity with which the program is associated.
 * @returns The program object.
 */
const newProgramModal = (activityID) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: NewProgramModal,
      componentProps: { activityID: activityID },
    })
      .onOk((programProps) => {
        resolve(programProps);
      })
      .onCancel(() => {})
      .onDismiss(() => {
        resolve();
      });
  });
};

/**
 * A dialog that creates a program instance.
 * @param {String} activityID The ID of the activity with which the program instance is associated.
 * @param {String} programID The ID of the program that is being instantiated.
 * @param {Function} callback A function to call after the instance is created. Takes an instance object as an argument.
 */
const newProgramInstanceModal = (activityID, programID, callback) => {
  Dialog.create({
    component: ProgramInstanceModal,
    componentProps: { activityID: activityID, programID: programID },
  })
    .onOk((instance) => {
      callback(instance);
    })
    .onCancel(() => {})
    .onDismiss(() => {});
};

/**
 * A dialog for editing the properties of a program.
 * @param {Object} program The program object
 * @returns The editing program object.
 */
const editProgramModal = (program) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: ProgramModal,
      componentProps: {
        program: program,
      },
    })
      .onOk((newValue) => {
        resolve(newValue);
      })
      .onCancel(() => {})
      .onDismiss(() => {
        resolve();
      });
  });
};

/**
 * A dialog for defining or editing the composition of an exercise. A composite exercise is derived from 2 or more other exercises, such as a clean and jerk.
 * @param {String} exerciseTypeID The ID of the composite exercise.
 * @param {Object} composition (Optional) A map where the keys are the exercise ID's and the values are the number of reps that are performed of that exercise. Provide no value when first defining the composition.
 * @param {Function} callback A function that takes the edited composition as an argument.
 */
const openCompositionModal = (exerciseTypeID, composition, callback) => {
  Dialog.create({
    component: CompositionModal,
    componentProps: { exerciseID: exerciseTypeID, composition: composition },
  })
    .onOk((composition) => {
      callback(composition);
    })
    .onCancel(() => {});
};

/**
 * A dialog that sets or edits, for an exercise, the ID of an exercise upon which it is based. For example, a clean deadlift is a variation of a deadlift.
 * @param {String} exerciseTypeID The ID of the exercise that is a variation of another exercise.
 * @param {String} basisID (Optional) The ID of the exercise that is the basis. Provide a value when editing an existing variation.
 * @param {Function} callback A function to execute after the variation is defined. The function takes the variation ID as an argument.
 */
const openVariationModal = (exerciseTypeID, basisID, callback) => {
  Dialog.create({
    component: VariationModal,
    componentProps: { exerciseID: exerciseTypeID, basisID: basisID },
  })
    .onOk((id) => {
      callback(id);
    })
    .onCancel(() => {});
};

/**
 * A generic dialog for confirming a change.
 * @param {String} message A message to display on the dialog that indicates the change that is being confirmed.
 * @returns A boolean that is true when the change is confirmed and false when the change is cancelled.
 */
const openConfirmModal = (message) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: ConfirmModal,
      componentProps: { message: message },
    })
      .onOk(() => {
        resolve(true);
      })
      .onCancel(() => {
        resolve(false);
      });
  });
};

/**
 * Authenticates with the back end using a user name and password.
 * @param {String} id The user name.
 * @param {String} pwd The password.
 */
const login = async (id, pwd) => {
  const url = '/homegym/login';
  const body = `{"username": "${id}", "password": "${pwd}"}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: body,
    headers: headers,
  };

  const resp = await fetch(url, options);
  if (resp.status == 401) {
    throw new ErrNotLoggedIn('logged in using invalid credentials');
  } else if (resp.status < 200 || resp.status >= 300) {
    throw new Error('failed to login');
  }
};

/**
 * Stores a new exercise type on the back end. Stores the exercise type in the exercise type store.
 * @param {String} name The name of the exercise type.
 * @param {String} intensityType The name of the intensity that the exercise uses. Valid values are defined in the intensityTypes constant.
 * @param {String} VolumeType The name of the type of volume that the exercise uses. Valid values are defined in the volumeTypes constant.
 * @param {Number} volumeConstraint A number that indicates how to interpret volume values. See the VolumeReps component for more information.
 * @param {String} basisID The ID of the exercise that forms the basis of this exercise type.
 * @returns The ID of the new exercise type.
 */
const addExerciseType = async (
  name,
  intensityType,
  VolumeType,
  volumeConstraint,
  basisID
) => {
  const url = '/homegym/api/exercises/';

  const exerciseType = {
    name: name,
    intensityType: intensityType,
    volumeType: VolumeType,
    volumeConstraint: volumeConstraint,
    basis: basisID,
  };

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(exerciseType),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized upsert of exercise type');
  } else if (resp.status == 400) {
    //TODO: check error  message before claiming ErrNotUnique
    throw new Error(ErrNotUnique);
  } else if (resp.status < 200 && resp.status >= 300) {
    throw new Error();
  }

  const respBody = await resp.json();

  exerciseType.id = respBody.id;
  exerciseTypeStore.add(exerciseType);

  return respBody.id;
};

/**
 * Updates an exercise type on the back end. Stores the exercise type in the exercise type store.
 * @param {Object} exerciseType The exerciseType object.
 * @returns
 */
const updateExerciseType = async (exerciseType) => {
  const url = `/homegym/api/exercises/${exerciseType.id}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(exerciseType),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized upsert of exercise type');
  } else if (resp.status == 400) {
    throw new Error('bad request');
  } else if (resp.status < 200 && resp.status >= 300) {
    throw new Error();
  }

  exerciseTypeStore.add(exerciseType);

  return;
};

/**
 * Updates the list of exercises that are associated with an activity. Stores the updated activity in the activity store.
 * @param {Object} activity The activity object that contains the changed list of exercises.
 */
const updateActivityExercises = async (activity) => {
  const url = `/homegym/api/activities/${activity.id}/exercises/`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(activity),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of activity exercises');
  } else if (resp.status < 200 || resp.status >= 300) {
    console.log('failed to update activity exercises');
    throw new Error();
  }

  activityStore.add(activity);
};

/**
 * Adds or updates a program on the back end. Stores the program in the program store.
 * When the program object has an ID, the existing program is updated. No ID stores a new program.
 * @param {Object} program The program object.
 * @returns The program ID.
 */
const updateProgram = async (program) => {
  if (!program.activityID) {
    throw new Error('missing activity ID');
  }
  const id = program.id ? program.id : '';
  const url = `/homegym/api/activities/${program.activityID}/programs/${id}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(program),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of program');
  } else if (resp.status < 200 || resp.status >= 300) {
    const errBody = await resp.json();
    throw new Error(errBody.message);
  }

  if (!!!program.id) {
    const respBody = await resp.json();

    program.id = respBody.id;
  }
  programsStore.add(program);
  return program.id;
};

/**
 * Adds or updates a program instance on the back end. Stores the instance in the program instance store.
 * When the instance object has an ID, the existing program instance is updated. No ID stores a new program instance.
 * @param {Object} instance The program instance object.
 * @returns The program instance ID.
 */
const updateProgramInstance = async (instance) => {
  const rawInstance = deepToRaw(instance);
  if (!rawInstance.activityID) {
    throw new Error('missing activity ID');
  }
  if (!rawInstance.programID) {
    throw new Error('missing program ID');
  }
  const isNewInstance = rawInstance.id ? true : false;
  const url = `/homegym/api/activities/${rawInstance.activityID}/programs/${
    rawInstance.programID
  }/instances/${rawInstance.id ? rawInstance.id : ''}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(rawInstance),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of program');
  } else if (resp.status < 200 || resp.status >= 300) {
    const errBody = await resp.json();
    throw new Error(errBody.message);
  }

  if (!rawInstance.id) {
    const respBody = await resp.json();

    rawInstance.id = respBody.id;
  }

  // New instances are set as active
  if (isNewInstance) {
    programInstanceStore.add(rawInstance);
  }
  return rawInstance.id;
};

/**
 * Deactivates a specific program instance by ID. Called when the instance is completed.
 * @param {String} activityID The ID of the activity with which the instance is associated.
 * @param {String} instanceID The ID of the instance.
 */
const deactivateProgramInstance = async (activityID, instanceID) => {
  const url = `/homegym/api/activities/${activityID}/programs/instances/active?instanceid=${instanceID}`;

  const headers = new Headers();

  const options = {
    method: 'DELETE',
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of program');
  } else if (resp.status < 200 || resp.status >= 300) {
    throw new Error('error deactivating program instance');
  }

  programInstanceStore.removeActive(activityID);
};

// event param has no id if it is new
/**
 * Adds or updates a workout event on the backend. When the event object contains an ID, the existing event is updated. No ID creates a new event.
 * @param {String} url The URL of the REST operation.
 * @param {Object} event The event object.
 * @returns The event object.
 */
const storeEvent = async (url, event) => {
  const headers = new Headers();

  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(event),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized upsert of event');
  } else if (resp.status < 200 || resp.status >= 300) {
    console.log('failed to upsert event');
    throw new Error();
  }

  if (resp.status != 204) {
    const body = await resp.json();
    event.id = body.id;
  }

  return event;
};

/**
 * Deletes an event from the back end.
 * @param {Object} event The event object.
 * @returns
 */
const deleteEvent = async (event) => {
  const url = `/homegym/api/events/${event.date}/${event.id}/`;

  const headers = new Headers();

  headers.set('content-type', 'application/json');

  const options = {
    method: 'DELETE',
    body: JSON.stringify(event),
    headers: headers,
  };

  const resp = await fetch(url, options);

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized delete of event');
  } else if (resp.status < 200 || resp.status >= 300) {
    console.log('failed to delete event');
    throw new Error();
  }

  if (resp.status != 204) {
    const body = await resp.json();
    console.log(body.message);
  }

  return;
};

/**
 * Displays a toast notification in the UI to indicate a successful or failed operation.
 * @param {String} message The message to display on the toast
 * @param {String} type A value of 'positive' for a success, and any other value for a failure.
 */
const toast = (message, type) => {
  const color = type == 'positive' ? 'green' : 'red';
  const icon = type == 'positive' ? 'checkmark' : 'error';
  Notify.create({
    type: type,
    color: color,
    message: message,
    icon: icon,
    position: 'top-right',
    timeout: 2000,
  });
};

const states = {
  READ_ONLY: 0,
  EDIT: 1,
  NEW: 2,
};

/**
 * An error to use when a certain value is not unique.
 */
class ErrNotUnique extends Error {
  constructor(message) {
    super(message);
    this.name = this.constructor.name;
  }
}

/**
 * The error to throw when a back end operation is called but authentication fails.
 */
class ErrNotLoggedIn extends Error {
  constructor(message) {
    super(message);
    this.name = this.constructor.name;
  }
}

/**
 * An ordered list supported by a number of actions that manipulate the list.
 * The ListActions component provides a UI for the list actions.
 */
class OrderedList {
  // model
  list = [];
  copiedIndex = -1;
  // define the verbs for items in the list
  static #actions = {
    add: 0,
    delete: 1,
    moveback: 2,
    moveahead: 3,
    copy: 4,
    paste: 5,
  };
  static get ADD() {
    return OrderedList.#actions.add;
  }
  static get DELETE() {
    return OrderedList.#actions.delete;
  }
  static get MOVEBACK() {
    return OrderedList.#actions.moveback;
  }
  static get MOVEFWD() {
    return OrderedList.#actions.moveahead;
  }
  static get COPY() {
    return OrderedList.#actions.copy;
  }
  static get PASTE() {
    return OrderedList.#actions.paste;
  }
  get list() {
    return this.list;
  }

  constructor(arr) {
    if (arr) {
      this.list = arr;
    } else {
      this.list = [];
    }
  }

  /**
   * Use update when the action is issued from the context of a component
   * such that the component is unaware of the list item index.
   * The component emits the action value.
   */
  update(action, index) {
    switch (action) {
      case OrderedList.#actions.add:
        this.addItem(index);
        break;
      case OrderedList.#actions.delete:
        this.deleteItem(index);
        break;
      case OrderedList.#actions.moveback:
        this.shiftItemBack(index);
        break;
      case OrderedList.#actions.moveahead:
        this.shiftItemForward(index);
        break;
      case OrderedList.#actions.copy:
        this.copyItem(index);
        break;
      case OrderedList.#actions.paste:
        this.pasteItem(index);
        break;
    }
  }
  deleteItem(index) {
    this.list.splice(index, 1);
  }
  addItem(index) {
    if (!index) {
      this.list.splice(this.list.length, 0, {});
    } else {
      // add the item after the one on which add was invoked
      this.list.splice(index + 1, 0, {});
    }
  }
  shiftItemBack(index) {
    if (index > 0 && index < this.list.length) {
      const moved = this.list.splice(index, 1);
      this.list.splice(index - 1, 0, moved[0]);
    }
  }
  shiftItemForward(index) {
    if (index >= 0 && index < this.list.length - 1) {
      const moved = this.list.splice(index, 1);
      this.list.splice(index + 1, 0, moved[0]);
    }
  }
  copyItem(index) {
    this.copiedIndex = index;
    console.log(this.list[index]);
  }
  pasteItem(index) {
    if (this.copiedIndex > -1) {
      const copy = Object.assign({}, this.list[this.copiedIndex]);
      console.log(copy);
      this.list[index] = copy;
    }
  }
}

// https://stackoverflow.com/a/21125098/2307622
const getCookieValue = (name) => {
  const regex = new RegExp(`(^| )${name}=([^;]+)`);
  const match = document.cookie.match(regex);
  if (match) {
    return match[2];
  }
};

//todo: use when logout is implemented
const deleteCookie = (name) => {
  const regex = new RegExp(`(^| )${name}=([^;]+)`);
  const match = document.cookie.match(regex);
  if (match) {
    return match[2];
  }
};

//https://github.com/vuejs/core/issues/5303#issuecomment-1543596383
const deepToRaw = (sourceObj) => {
  const objectIterator = (input) => {
    if (Array.isArray(input)) {
      return input.map((item) => objectIterator(item));
    }
    if (isRef(input) || isReactive(input) || isProxy(input)) {
      return objectIterator(toRaw(input));
    }
    if (input && typeof input === 'object') {
      return Object.keys(input).reduce((acc, key) => {
        acc[key] = objectIterator(input[key]);
        return acc;
      }, {});
    }
    return input;
  };

  return objectIterator(unref(sourceObj));
};

export {
  authPromptAsync,
  fetchActivities,
  fetchPrograms,
  fetchProgramInstances,
  fetchActiveProgramInstances,
  fetchEventPage,
  fetchEvents,
  fetchExerciseTypes,
  pageSize,
  login,
  fetchActivityExercises,
  intensityTypes,
  intensityProps,
  volumeTypes,
  addExerciseType,
  updateExerciseType,
  updateProgram,
  ErrNotUnique,
  updateActivityExercises,
  ErrNotLoggedIn,
  openEditValueModal,
  newActivityPrompt,
  newProgramModal,
  newProgramInstanceModal,
  editProgramModal,
  updateProgramInstance,
  deactivateProgramInstance,
  storeEvent,
  deleteEvent,
  openVolumeModal,
  toast,
  openCompositionModal,
  openVariationModal,
  openConfirmModal,
  states,
  OrderedList,
  getCookieValue,
  deepToRaw,
};
