//
/**
 * Creates a date object from a seconds since epoch value (i.e. UTC)
 * @param {Number} seconds The seconds since epoch.
 * @returns The date object in local time
 */
export const dateFromSeconds = (seconds) => {
  // javascript epoch is in milliseconds
  return seconds ? new Date(seconds * 1000) : new Date();
};

/**
 * Generates the current time in seconds since epoch.
 * @returns An integer that represents the seconds since epoch at the time of execution.
 */
export const nowInSeconds = () => {
  return Math.floor(Date.now() / 1000);
};

/**
 * Determines the time value for midnight of the day that a time value lands on
 * @param {number} seconds The time value for which midnight of that day is determined, in seconds since epoch.
 * @returns The time value of midnight, in milliseconds since midnight.
 */
export const setEpochToMidnight = (seconds) => {
  const date = dateFromSeconds(seconds);
  date.setHours(0, 0, 0);
  return Math.floor(date.valueOf() / 1000);
};

/**
 * Transforms a date string to a timestamp UTC in seconds.
 * @param {String} dateString A date value in the format YYYY-MM-DD HH-MM, eg. 2024-05-12 9:58
 * @returns An integer that represents the seconds since epoch.
 */
export const stringToEpoch = (dateString) => {
  const date = new Date(dateString); // local time zone
  const milliseconds = date.valueOf(); // UTC
  return Math.floor(milliseconds / 1000);
};

// Use to prefix single-digit date numbers
const prefixed = (dateNumber) => {
  const prefix = dateNumber < 10 ? '0' : '';
  return `${prefix}${dateNumber}`;
};

// Returns a formatted string as YYYY-MM-DD for a Date object
/**
 * Transforms a date object to a string in the format YYYY-MM-DD.
 * @param {Date} dateObj A Date object.
 * @returns The formatted string.
 */
export const formatDate = (dateObj) => {
  return `${dateObj.getFullYear()}-${prefixed(
    dateObj.getMonth() + 1
  )}-${prefixed(dateObj.getDate())}`;
};

/**
 * Transforms a date object to a string in the format HH:MM.
 * @param {*} dateObj A Date object.
 * @returns The formatted string.
 */
export const formatTime = (dateObj) => {
  const minutes =
    dateObj.getMinutes() < 10
      ? '0' + dateObj.getMinutes()
      : dateObj.getMinutes();
  return `${dateObj.getHours()}:${minutes}`;
};

//
/**
 * Transforms a date object to a formatted date and time string in the format YYYY-MM-DD HH:MM
 * @param {*} dateObj A date object.
 * @returns The formatted string.
 */
export const formatDateTime = (dateObj) => {
  return `${formatDate(dateObj)} ${formatTime(dateObj)}`;
};

export const dateMask = 'YYYY-MM-DD';
export const timeMask = 'HH:MM';
