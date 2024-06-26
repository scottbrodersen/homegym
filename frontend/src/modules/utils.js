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

const pageSize = () => {
  try {
    const availableHeight = window.innerHeight - 220 - 36;
    return Math.floor((availableHeight - 48 - 50 - 40 - 30) / 42);
  } catch (e) {
    return 8;
  }
};
const fetchPageSize = 20;

const intensityTypes = [
  'weight',
  'bodyweight',
  'rpe',
  'percentOfMax',
  'hrZone',
  'distance',
  'pace',
];

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

const volumeTypes = ['count', 'time', 'distance'];

const fetchEventPage = async (eventID = '', date = null) => {
  const startTime = !!date ? date : Math.floor(Date.now() / 1000);
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
  eventStore.addBulk(eventPage);
  return eventPage;
};

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
        lastProgram = programPage[programPage.length].id;
      }
    } else {
      done = true;
    }
  }
};

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
      lastInstance = instancePage[instancePage.length].id;
    }
  }
};

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

const fetchActiveProgramInstance = async (activityID) => {
  const url = `/homegym/api/activities/${activityID}/programs/instances/active/`;
  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of active program instance');
  } else if (resp.status == 404) {
    programInstanceStore.setActive(activityID, null);
    return;
  }

  const instance = await resp.json();

  programInstanceStore.setActive(activityID, instance);
};

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

const newActivityPrompt = () => {
  Dialog.create({
    component: NewActivityModal,
  })
    .onOk(() => {})
    .onCancel(() => {})
    .onDismiss(() => {});
};

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

const newProgramModal = (activityID, callback) => {
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

const openConfirmModal = async (message, callback) => {
  Dialog.create({
    component: ConfirmModal,
    componentProps: { message: message },
  })
    .onOk(async () => {
      await callback();
    })
    .onCancel(() => {});
};

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
  if (!isNewInstance) {
    programInstanceStore.setActive(rawInstance.activityID, rawInstance);
  } else {
    programInstanceStore.add(rawInstance);
  }
  return rawInstance.id;
};

const deactivateProgramInstance = async (activityID) => {
  const url = `/homegym/api/activities/${activityID}/programs/instances/active/`;

  const headers = new Headers();
  //headers.set('content-type', 'application/json');

  const options = {
    method: 'DELETE',
    //body: JSON.stringify(rawInstance),
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

class ErrNotUnique extends Error {
  constructor(message) {
    super(message);
    this.name = this.constructor.name;
  }
}

class ErrNotLoggedIn extends Error {
  constructor(message) {
    super(message);
    this.name = this.constructor.name;
  }
}

class OrderedList {
  list = [];
  static #actions = { add: 0, delete: 1, moveback: 2, moveahead: 3 };
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
   * Use when the action is issued from the context of the item component
   * and the component is unaware of its index.
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
    }
  }
  deleteItem(index) {
    this.list.splice(index, 1);
  }
  addItem(index) {
    if (!index) {
      this.list.splice(this.list.length, 0, {});
    } else {
      this.list.splice(index, 0, {});
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
}

// https://stackoverflow.com/a/21125098/2307622
const getCookieValue = (name) => {
  const regex = new RegExp(`(^| )${name}=([^;]+)`);
  const match = document.cookie.match(regex);
  if (match) {
    return match[2];
  }
};

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
  fetchActiveProgramInstance,
  fetchEventPage,
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
