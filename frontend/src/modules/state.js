import { reactive, computed } from 'vue';
import { pageSize } from './utils';

/*
workout events be like:
  {
    "id": "6ef6125f-1708-4977-8dad-598bfff2a3a4",
    "activityID": "49adfe46-854f-425e-9782-a1d9a5255e43",
    "date": 1719402960,
    "mood": 4,
    "motivation": 4,
    "energy": 2,
    "overall": 4,
    "notes": "",
    "exercises": {
      "0": {
        "typeID": "de2e77e0-c827-4bcb-a9fc-0255cbd0ec57",
        "index": 0,
        "parts": [
          {
            "intensity": 1,
            "volume": [
              [
                2.2
              ]
            ]
          }
        ]
      }
    }
  }
*/
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

  delete(e) {
    let index;
    for (let i = 0; i < this.events.length; i++) {
      if (this.events[i].id === e.id) {
        index = i;
        break;
      }
    }
    if (index) {
      this.events.splice(index, 1);
    }
  },

  clear() {
    this.events = [];
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
    if (programs) {
      for (const program of programs) {
        this.add(program);
      }
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
  // key is the activityID, value is an array of objects with fields programID, instanceID
  activeInstances: new Map(),
  // key is the activityID, value is an object with fields programID, instanceID
  // One current instance per activity is allowed
  currentInstances: new Map(),

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

  // Stores the programID and instanceID of the active instance for an activityID
  // An instance value of null indicates no active instance
  setCurrent(activityID, instance) {
    if (instance && Object.keys(instance).length > 0) {
      this.currentInstances.set(activityID, {
        programID: instance.programID,
        instanceID: instance.id,
      });

      //this.add(instance);
    } else {
      this.currentInstances.set(activityID, null);
    }
  },

  getCurrent(activityID) {
    if (activityID && this.currentInstances.has(activityID)) {
      const currentIDs = this.currentInstances.get(activityID);
      if (currentIDs) {
        const current = this.programInstances
          .get(currentIDs.programID)
          .get(currentIDs.instanceID);

        return current;
      }
    }
    return null;
  },

  // Called after the set of active instances are fetched from the server
  addAllActive(activityID, instances) {
    if (instances && instances.length > 0 && instances[0] != null) {
      let instancesToAdd = [];
      for (const inst of instances) {
        instancesToAdd.push({ programID: inst.programID, instanceID: inst.id });
      }
      this.activeInstances.set(activityID, instancesToAdd);
    } else {
      this.activeInstances.set(activityID, []);
    }
  },

  getActive(activityID) {
    const active = [];
    let programID;
    let instanceID;
    if (
      this.activeInstances.has(activityID) &&
      this.activeInstances.get(activityID)
    ) {
      const mapping = this.activeInstances.get(activityID);
      for (const inst of mapping) {
        // key is the activityID, value is an object with fields programID, instanceID
        programID = inst.programID;
        instanceID = inst.instanceID;
        active.push(this.get(instanceID, programID));
      }
      return active;
    }
  },

  removeActive(activityID) {
    this.activeInstances.set(activityID, null);
    if (this.currentInstances.get(activityID)) {
      this.currentInstances.set(activityID, null);
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

export const dailyStatsStore = reactive({
  dailyStats: new Array(),
  sleep: 0,
  bodyweight: 0,
  spirit: { mood: 0, stress: 0, energy: 0 },
  add(stat) {
    this.dailyStats.push(stat);
    if (stat.sleep) {
      this.sleep = stat;
    }
    if (stat.bodyweight) {
      this.bodyweight = stat;
    }
    if (stat.mood || stat.stress || stat.energy) {
      this.spirit = stat;
    }
  },
  update(stat) {
    for (let i; i < this.dailyStats.length; i++) {
      if (this.dailyStats[i].date == stat.date) {
        this.dailyStats.splice(i, 1);
        break;
      }
    }
    this.add(stat);
  },
  bulkAdd(stats) {
    for (const [key, value] of Object.entries(stats)) {
      this.add(value);
    }
  },
  getAll() {
    return this.dailyStats;
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
