// Creates a date object from a seconds since epoch value (i.e. UTC)
// Note that the date object is in local time
export const dateFromSeconds = (seconds) => {
  // javascript epoch is in milliseconds
  return seconds ? new Date(seconds * 1000) : new Date();
};

export const nowInSeconds = () => {
  return Math.floor(Date.now() / 1000);
};

export const setEpochToMidnight = (seconds) => {
  const date = dateFromSeconds(seconds);
  date.setHours(0, 0, 0);
  return Math.floor(date.valueOf() / 1000);
};

// transforms the date string to timestamp UTC in seconds
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
export const formatDate = (dateObj) => {
  return `${dateObj.getFullYear()}-${prefixed(
    dateObj.getMonth() + 1
  )}-${prefixed(dateObj.getDate())}`;
};

// Returns a formatted string as HH:mm
export const formatTime = (dateObj) => {
  const minutes =
    dateObj.getMinutes() < 10
      ? '0' + dateObj.getMinutes()
      : dateObj.getMinutes();
  return `${dateObj.getHours()}:${minutes}`;
};

export const formatDateTime = (dateObj) => {
  return `${formatDate(dateObj)} ${formatTime(dateObj)}`;
};

export const dateMask = 'YYYY-MM-DD';
export const timeMask = 'HH:MM';
