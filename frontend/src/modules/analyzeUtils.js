/**
 * Functions for interacting with metrics.
 */
import * as utils from '../modules/utils';
import * as dateUtils from '../modules/dateUtils';
import * as state from '../modules/state';
import Chart from 'chart.js/auto';
import 'date-fns';
import 'chartjs-adapter-date-fns';
import 'date-fns';

/**
 * Fetches exercise metrics from the server.
 *
 * @param {integer} startDate The first day of the time period for which metrics are fetched.
 * @param {integer} endDate The final day of the time period.
 * @param {[string]} exercises Array of exercise IDs. Empty value aborts the call.
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
 * Fetches time series data from the server over a period of time. Time series events are those that can occur multiple times during a day.
 * Examples are blood pressure readings and food intake.
 * @param {number} startDate The start date of the period of time, in seconds since epoch.
 * @param {number} endDate The end date of the period of time, in seconds since epoch.
 * @returns The time series data formatted for graphing.
 */
export const getTimeSeriesData = async (startDate, endDate) => {
  const timeSeries = new Array();

  if (
    state.eventStore.events.length == 0 ||
    startDate > state.eventStore.events[0].date
  ) {
    // we need to get more recent events
    // easiest to just reset the cache
    state.eventStore.clear();
    await utils.fetchEvents('', startDate, endDate);
  } else {
    const lastCachedEvent = state.eventStore.getLast();

    if (endDate < lastCachedEvent.date) {
      // we need to get earlier events
      await utils.fetchEvents(
        lastCachedEvent.id,
        lastCachedEvent.date,
        endDate
      );
    }
  }
  const events = state.eventStore.getAll();

  events.forEach((evt) => {
    if (evt.date <= startDate && evt.date >= endDate) {
      const x = dateUtils.setEpochToMidnight(evt.date) * 1000;

      // stat times normalized to be within 24h of epoch
      const offsetMS = new Date().getTimezoneOffset() * 60 * 1000;
      const y = evt.date * 1000 - x + offsetMS;

      const statDate = new Date(y);
      const timeStr = dateUtils.formatTime(statDate);

      const description = new Array();
      description.push(state.activityStore.get(evt.activityID).name);

      for (const ex of Object.values(evt.exercises)) {
        description.push(state.exerciseTypeStore.get(ex.typeID).name);
      }

      timeSeries.push({
        x: x,
        y: y,
        description: description,
        timeStr: timeStr,
      });
    }
  });
  return timeSeries;
};

/**
 * Amalgamates raw metrics into daily totals.
 * @param {[Object]} metrics is an array of objects with properties of date, volume, and load that pertain to the performance of an exercise.
 * @returns An object with properties dates (array of dates), lvRatio (an array of load/volume ratios), and load (an array of load values).
 *  The arrays can be interpreted as stacks in that the related items in each array have the same index.
 */
export const getDailyTotals = (metrics) => {
  // key is date (at midnight), values an array of [total volume; total load/volume]
  const dailyTotals = new Map();

  // aggregate metrics for each day
  for (let i = 0; i < metrics.dates.length - 1; i++) {
    const midnight = dateUtils.setEpochToMidnight(metrics.dates[i]);
    if (dailyTotals.has(midnight)) {
      const volTot = dailyTotals.get(midnight)[0] + metrics.volume[i];
      const loadTot = dailyTotals.get(midnight)[1] + metrics.load[i];

      //      dailyTotals.set(midnight, [volTot != 0 ? loadTot / volTot : 0, loadTot]);
      dailyTotals.set(midnight, [volTot, loadTot]);
    } else {
      // dailyTotals.set(midnight, [
      //   metrics.volume[i] != 0 ? metrics.load[i] / metrics.volume[i] : 0,
      //   metrics.load[i],
      // ]);
      dailyTotals.set(midnight, [metrics.volume[i], metrics.load[i]]);
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

/**
 * Generates a chart of exercise metrics
 * @param {*} element The id of the HTML element in which to insert the chart.
 * @param {*} startDate The start date of the charted data.
 * @param {*} endDate The end date of the charted data.
 * @param {[string]} dates The labels for the x-axis.
 * @param {[number]} lvRatioData The load/volume ratios to chart.
 * @param {[number]} loadData The load values to chart
 * @returns
 */
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
        { data: lvRatioData, yAxisID: 'yVol', label: 'Volume' },
        { data: loadData, yAxisID: 'yLoad', label: 'Load' },
      ],
    },
    options: {
      maintainAspectRatio: false,
      aspectRatio: 1,
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
            text: 'Volume',
          },
          position: 'left',
        },
      },
    },
  });
};
