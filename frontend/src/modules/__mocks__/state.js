import * as data from '../../../mocks/data';

export const programsStore = {
  // key is the activityID, value is a map of (activityID, program)
  // programs: new Map(),
  // add(program) {
  //   if (this.programs.get(program.activityID)) {
  //     this.programs.get(program.activityID).set(program.id, program);
  //   } else {
  //     this.programs.set(program.activityID, new Map([[program.id, program]]));
  //   }
  // },
  // addBulk(programs) {
  //   for (const program of programs) {
  //     this.add(program);
  //   }
  // },
  getByActivity(activityID) {
    return data.testProgram();
  },
  get(activityID, programID) {
    return data.testProgram();
  },
};
export const exerciseTypeStore = {
  // exerciseTypes: new Map(),
  // add(exerciseType) {},
  get(exerciseTypeID) {
    return data.fetchedEventExercises[0];
  },
  getAll() {
    return data.fetchedEventExercises;
  },
};

export const activityStore = {
  // activities: new Map(),
  // add(activity) {
  //   this.activities.set(activity.id, activity);
  // },
  get(activityID) {
    return data.fetchedTestActivities[0];
  },
  getAll() {
    return data.fetchedTestActivities;
  },
};

export const programInstanceStore = {
  // key is the programID, value is a map of (programInstanceID, programInstance)
  programInstances: new Map(),
  // key is the activityID, value is an object with fields programID, instanceID
  activeInstances: new Map(),
  add(instance) {},
  addBulk(instances) {},
  getByProgram(programID) {
    if (programID == data.testProgramID) {
      return [data.testProgramInstance];
    }
    return undefined;
  },
  get(programID, instanceID) {
    if (programID) {
      return data.testProgramInstance(programID);
    }
    return undefined;
  },
  setActive(activityID, instance) {
    return;
  },
  getActive(activityID) {
    if (activityID == data.testActivityID) {
      return data.testProgramInstance();
    } else {
      // assume the ID is the instance start date
      return data.testProgramInstance(activityID);
    }
  },
};
