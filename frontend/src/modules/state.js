import { reactive, computed } from 'vue';
import { pageSize } from './utils';

export const eventStore = reactive({
  events: [],

  add(e) {
    this.events.push(e);
  },

  update(e) {
    for (let i = 0; i < this.events.length; i++) {
      if (this.events[i].id === e.id) {
        this.events[i] = e;
        return;
      }
    }
    throw Error('event not found');
  },

  addBulk(events) {
    this.events = this.events.concat(events);
  },

  getLast() {
    return this.events[this.events.length - 1];
  },

  getByID(eventID) {
    for (const e of this.events) {
      if (e.id === eventID) {
        return e;
      }
    }
    return null;
  },

  getPage(pageNumber = 0) {
    const start = pageNumber * pageSize();
    const end = start + pageSize();
    return this.events.slice(start, end);
  },
  getAll() {
    const all = [];
    for (const event of this.events.values()) {
      all.push(event);
    }
    return all;
  },
});

export const activityStore = reactive({
  activities: new Map(),
  add(activity) {
    this.activities.set(activity.id, activity);
  },
  get(activityID) {
    return this.activities.get(activityID);
  },
  getAll() {
    const all = [];
    for (const activity of this.activities.values()) {
      all.push(activity);
    }
    return all;
  },
});

export const exerciseTypeStore = reactive({
  exerciseTypes: new Map(),
  add(exerciseType) {
    this.exerciseTypes.set(exerciseType.id, exerciseType);
  },
  get(exerciseTypeID) {
    return this.exerciseTypes.get(exerciseTypeID);
  },
  getAll() {
    const all = [];
    for (const exerciseType of this.exerciseTypes.values()) {
      all.push(exerciseType);
    }
    return all;
  },
});

export const programsStore = reactive({
  // key is the activityID, value is a map of (programID, program)
  programs: new Map(),
  add(program) {
    if (this.programs.get(program.activityID)) {
      this.programs.get(program.activityID).set(program.id, program);
    } else {
      this.programs.set(program.activityID, new Map([[program.id, program]]));
    }
  },
  addBulk(programs) {
    for (const program of programs) {
      this.add(program);
    }
  },
  getByActivity(activityID) {
    if (this.programs.get(activityID)) {
      const programs = [];
      const iter = this.programs.get(activityID).values();
      let p = iter.next();
      while (!p.done) {
        programs.push(p.value);
        p = iter.next();
      }
      return programs;
    }
    return undefined;
  },
  get(activityID, programID) {
    if (this.programs.has(activityID)) {
      return this.programs.get(activityID).get(programID);
    }
    return undefined;
  },
});

export const programInstanceStore = reactive({
  // key is the programID, value is a map of (programInstanceID, programInstance)
  programInstances: new Map(),
  // key is the activityID, value is an object with fields programID, instanceID
  activeInstances: new Map(),
  add(instance) {
    if (this.programInstances.get(instance.programID)) {
      this.programInstances.get(instance.programID).set(instance.id, instance);
    } else {
      this.programInstances.set(
        instance.programID,
        new Map([[instance.id, instance]])
      );
    }
  },
  addBulk(instances) {
    for (const program of instances) {
      this.add(program);
    }
  },
  getByProgram(programID) {
    if (this.programInstances.has(programID)) {
      const instances = [];
      const iter = this.programInstances.get(programID).values();
      let p = iter.next();
      while (!p.done) {
        instances.push(p.value);
        p = iter.next();
      }
      return instances;
    }
    console.warn('Instances of program not yet added');
    return undefined;
  },
  get(instanceID, programID) {
    if (programID) {
      if (this.programInstances.has(programID)) {
        return this.programInstances.get(programID).get(instanceID);
      }
    } else {
      const instanceMaps = this.programInstances.values();
      let found = false;
      let inst;
      let imap = instanceMaps.next();
      while (!found && !imap.done) {
        inst = imap.value.get(instanceID);
        if (inst) found = true;
        imap = instanceMaps.next();
      }
      return inst;
    }

    return undefined;
  },
  // An instance value of null indicates no active instance
  setActive(activityID, instance) {
    if (instance) {
      this.activeInstances.set(activityID, {
        programID: instance.programID,
        instanceID: instance.id,
      });

      this.add(instance);
    } else {
      this.activeInstances.set(activityID, null);
    }
  },
  getActive(activityID) {
    let programID;
    let instanceID;
    if (this.activeInstances.has(activityID)) {
      const mapping = this.activeInstances.get(activityID);
      if (mapping) {
        programID = mapping.programID;
        instanceID = mapping.instanceID;
        return this.get(instanceID, programID);
      }
      return null;
    }
  },
});

export const eventMetricsStore = reactive({
  eventMetrics: new Map(),
  add(eventId, metrics) {
    this.eventMetrics.set(eventId, metrics);
  },
  get(eventId) {
    return !!this.eventMetrics.get(eventId)
      ? this.eventMetrics.get(eventId)
      : null;
  },
  setMetric(eventId, name, value) {
    let metric;
    if (this.eventMetrics.has(eventId)) {
      metric = this.eventMetrics.get(eventId);
    } else {
      metric = {};
    }
    metric[name] = value;
    this.add(eventId, metric);
  },
  getMetric(eventId, name) {
    return this.eventMetrics.has(eventId) &&
      !!this.eventMetrics.get(eventId)[name]
      ? this.eventMetrics.get(eventId)[name]
      : null;
  },
});

export const loginModalState = reactive({
  isOpen: false,
  opened() {
    this.isOpen = true;
  },
  closed() {
    this.isOpen = false;
  },
});

export const metricState = reactive({
  metric: true,
  setMetric() {
    this.metric = true;
  },
  setImperial() {
    this.metric = false;
  },
});

const weightUnit = computed(() => {
  return metricState.metric ? 'kg' : 'lbs';
});

export const unitsState = reactive({
  weight: weightUnit,
  longDistance: 'km',
  distance: 'm',
  time: 'min',
});
