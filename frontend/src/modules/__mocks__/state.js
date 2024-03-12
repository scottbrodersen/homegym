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
  // getByActivity(activityID) {
  //   if (this.programs.get(activityID)) {
  //     const programs = [];
  //     const iter = this.programs.get(activityID).values();
  //     let p = iter.next();
  //     while (!p.done) {
  //       programs.push(p.value);
  //       p = iter.next();
  //     }
  //     return programs;
  //   }
  //   return undefined;
  // },
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
  // getAll() {
  //   const all = [];
  //   for (const activity of this.activities.values()) {
  //     all.push(activity);
  //   }
  //   return all;
  // },
};
