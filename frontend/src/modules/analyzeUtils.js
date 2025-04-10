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
 * Amalgamates raw metrics into daily metrics.
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
            text: 'Load/Volume',
          },
          position: 'left',
        },
      },
    },
  });
};
