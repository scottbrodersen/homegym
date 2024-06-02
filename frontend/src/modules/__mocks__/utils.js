import * as data from '../../../mocks/data';

export const fetchActivityExercises = async (activityID) => {
  return data.fetchedExercises;
};

export const fetchEventPage = async (eventID, date) => {
  return data.fetchedEvents(10);
};
