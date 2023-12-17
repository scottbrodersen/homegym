import {
  eventStore,
  activityStore,
  exerciseTypeStore,
  loginModalState,
  programsStore,
} from './state';
import LoginModal from './../components/LoginModal.vue';
import { Dialog, Notify } from 'quasar';
import NewActivityModal from './../components/NewActivityModal.vue';
import VolumeModal from './../components/VolumeModal.vue';
import CompositionModal from './../components/CompositionModal.vue';
import VariationModal from './../components/VariationModal.vue';
import ProgramModal from './../components/ProgramModal.vue';
const pageSize = 10;
const fetchPageSize = 20;

const intensityTypes = [
  'weight',
  'bodyweight',
  'RPE',
  'percentOfMax',
  'hrZone',
  'distance',
];

const intensityProps = (intensityType) => {
  if (intensityType == 'hrZone') {
    return {
      mask: '#',
      validate: (value) => {
        return /^[1-5]$/.test(value);
      },
      decimals: 0,
      prefx: 'HR Zone',
    };
  } else {
    return {
      mask: '',
      validate: (value) => {
        return /^[0-9]+\.?[0-9]?$/.test(value);
      },
      decimals: 1,
      prefix: '',
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
    programsStore.addBulk(programPage);
    if (programPage.length < pageSize) {
      done = true;
    } else {
      lastProgram = programPage[programPage.length].id;
    }
  }
};

const fetchProgramPage = async (programID = '', activityID) => {
  const params = new URLSearchParams();

  params.append('size', fetchPageSize);

  if (programID) params.append('previous', eventID);

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

const fetchActivities = async () => {
  const resp = await fetch('/homegym/api/activities/', {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of activities');
  }

  const activities = await resp.json();
  activities.forEach((a) => {
    activityStore.add(a);
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

const fetchEventExercises = async (eventDate, eventID) => {
  const resp = await fetch(
    `/homegym/api/events/${eventDate}/${eventID}/exercises/`,
    {
      method: 'GET',
      mode: 'same-origin',
    }
  );
  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of event exercises');
  }
  const exercises = await resp.json();

  // store in an object to use numbers as keys
  const sorted = {};

  if (!!exercises) {
    const items = Object.values(exercises);

    items.sort((a, b) => {
      return a.index - b.index;
    });

    items.forEach((exInst) => {
      sorted[exInst.index] = exInst;
    });
  }
  try {
    eventStore.setEventExercises(eventID, sorted);
  } catch (e) {
    console.log(e);
  }
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
  activityStore.add(activity);
  return exercises;
};

const authPrompt = (after, args) => {
  if (!loginModalState.isOpen) {
    loginModalState.opened();
    Dialog.create({
      component: LoginModal,
    })
      .onOk(() => {})
      .onCancel(() => {})
      .onDismiss(() => {
        loginModalState.closed();
        if (after) {
          if (Array.isArray(args)) {
            after(...args);
          } else {
            after(args);
          }
        }
      });
  }
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

const newProgramModal = (callback) => {
  Dialog.create({
    component: ProgramModal,
    componentProps: {},
  })
    .onOk((programProps) => {
      callback(programProps);
    })
    .onCancel(() => {})
    .onDismiss(() => {});
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
    //TODO: check error  message before claiming errnotunique
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
    console.log('failed to update program');
    throw new Error();
  }

  if (!!!program.id) {
    const respBody = await resp.json();

    program.id = respBody.id;
  }
  programsStore.add(program);
  return program.id;
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

const storeEventExerciseInstances = async (
  eventID,
  eventDate,
  exerciseInstances
) => {
  const headers = new Headers();
  const options = {
    method: 'POST',
    body: JSON.stringify(exerciseInstances),
    headers: headers,
  };

  const url = `/homegym/api/events/${eventDate}/${eventID}/exercises/`;

  try {
    const resp = await fetch(url, options);

    if (resp.status == 401) {
      throw new ErrNotLoggedIn('unauthorized upsert of event');
    } else if (resp.status < 200 || resp.status >= 300) {
      console.log('failed to store exercise instance');
      throw new Error();
    }
  } catch (e) {
    throw e;
  }
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

export {
  authPrompt,
  fetchActivities,
  fetchPrograms,
  fetchEventPage,
  fetchExerciseTypes,
  fetchEventExercises,
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
  newActivityPrompt,
  newProgramModal,
  storeEvent,
  storeEventExerciseInstances,
  openVolumeModal,
  toast,
  openCompositionModal,
  openVariationModal,
  states,
  OrderedList,
};
