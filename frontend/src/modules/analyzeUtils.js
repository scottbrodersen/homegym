import * as utils from '../modules/utils';

export const fetchMetrics = async (startDate, endDate, exercises) => {
  let url = '/homegym/api/events/metrics';
  const qParams = new Array();
  if (startDate) {
    qParams.push(`start=${startDate}`);
  }

  if (endDate) {
    qParams.push(`end=${endDate}`);
  }

  if (exercises.length > 0) {
    exercises.forEach((exID) => {
      qParams.push(`type=${exID}`);
    });
  }

  if (qParams.length > 0) {
    url = url + '?' + qParams.join('&');
  }

  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new ErrNotLoggedIn('unauthorized fetch of event metrics');
  }
  const metrics = await resp.json();

  return metrics;
};
