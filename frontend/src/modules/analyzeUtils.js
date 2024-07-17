import * as utils from '../modules/utils';
import * as dateUtils from '../modules/dateUtils';
import Chart from 'chart.js/auto';
import 'date-fns';
import 'chartjs-adapter-date-fns';
import 'date-fns';

/**
 * Fetches exercise metrics from the server.
 *
 * @param {} startDate
 * @param {*} endDate
 * @param {*} exercises Array of exercise IDs. Empty value aborts the call.
 * @returns An object with fields date, load, and volume.
 */
export const fetchMetrics = async (startDate, endDate, exercises) => {
  if (!exercises || exercises.length == 0) {
    return;
  }

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
    throw new utils.ErrNotLoggedIn('unauthorized fetch of event metrics');
  }
  const metrics = await resp.json();

  return metrics;
};

/**
 * Amalgamates raw metrics into daily metrics.
 * @param {*} metrics
 * @returns
 */
export const getDailyTotals = (metrics) => {
  const dailyTotals = new Map();

  for (let i = 0; i < metrics.dates.length - 1; i++) {
    const midnight = dateUtils.setEpochToMidnight(metrics.dates[i]);
    if (dailyTotals.has(midnight)) {
      const volTot = dailyTotals.get(midnight)[0] + metrics.volume[i];
      const loadTot = dailyTotals.get(midnight)[1] + metrics.load[i];

      dailyTotals.set(midnight, [volTot != 0 ? loadTot / volTot : 0, loadTot]);
    } else {
      dailyTotals.set(midnight, [
        metrics.volume[i] != 0 ? metrics.load[i] / metrics.volume[i] : 0,
        metrics.load[i],
      ]);
    }
  }

  const dailyMetricStacks = {
    dates: [],
    lvRatio: [],
    load: [],
  };

  dailyTotals.forEach((value, key) => {
    dailyMetricStacks.dates.push(key);
    dailyMetricStacks.lvRatio.push(Math.round(value[0]));
    dailyMetricStacks.load.push(Math.round(value[1]));
  });

  return dailyMetricStacks;
};

export const getVolumeChart = (
  element,
  startDate,
  endDate,
  dates,
  lvRatioData,
  loadData
) => {
  const existing = Chart.getChart(element);
  if (existing) {
    existing.destroy();
  }

  return new Chart(element, {
    type: 'line',
    data: {
      labels: dates,
      datasets: [
        { data: lvRatioData, yAxisID: 'yVol', label: 'Load/Volume' },
        { data: loadData, yAxisID: 'yLoad', label: 'Load' },
      ],
    },
    options: {
      scales: {
        x: {
          type: 'time',
          max: startDate * 1000,
          min: endDate * 1000,
        },
        yLoad: {
          position: 'right',
          title: {
            text: 'Load',
          },
        },
        yVol: {
          title: {
            text: 'Load/Volume',
          },
          position: 'left',
        },
      },
    },
  });
};
