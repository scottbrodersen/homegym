import DailyStatModal from '../components/DailyStatModal.vue';
import { Dialog } from 'quasar';
import * as utils from './utils';
import * as dateUtils from '../modules/dateUtils';
import Chart from 'chart.js/auto';
import 'date-fns';
import 'chartjs-adapter-date-fns';
import 'date-fns';
import { getRelativePosition } from 'chart.js/helpers';
import ChartDataLabels from 'chartjs-plugin-datalabels';

export const BLOODGLUCOSE = 'bg';
export const BLOODPRESSURE = 'bp';
export const BODYWEIGHT = 'bodyweight';
export const FOOD = 'food';
export const SLEEP = 'sleep';
export const MOOD = 'mood';
export const STRESS = 'stress';
export const ENERGY = 'energy';
export const SPIRIT = 'spirit';

export const statNames = [
  BLOODGLUCOSE,
  BLOODPRESSURE,
  BODYWEIGHT,
  FOOD,
  SLEEP,
  MOOD,
  STRESS,
  ENERGY,
];

export const spiritGroup = [MOOD, STRESS, ENERGY];

// Labels and units to use in the UI for daily stats
export const labels = {
  bg: ['Blood Glucose', 'mmol/L'],
  bp: ['Blood Pressure', 'Systolic', 'Diastolic', 'mm Hg'],
  bodyweight: ['Body Weight', 'lbs'],
  food: ['Food', 'Description', 'Protein', 'Carbs', 'Fat', 'Fiber', 'g'],
  sleep: ['Sleep', 'Hours'],
  mood: ['Mood', '/5'],
  stress: ['Stress', '/5'],
  energy: ['Energy', '/5'],
  spirit: ['Mood, Stress, Energy', '/5'],
};

/**
 * Generates an object of zero daily stats
 * @returns An object that has as properties the various daily stats we collect, all with zero values.
 */
export const emptyStats = () => {
  return {
    date: 0,
    bg: 0,
    bp: [0, 0],
    sleep: 0,
    food: {
      description: '',
      protein: 0,
      carbs: 0,
      fat: 0,
      fiber: 0,
    },
    bodyweight: 0,
    mood: 0,
    stress: 0,
    energy: 0,
  };
};

/**
 * Opens the DailyStatsModal dialog for editing a daily stat value
 * @param {string} statName The name of the daily stat.
 * @param {*} stats (Optional) The value of the stat.
 * @returns The stats from the dialog.
 */
export const openDailyStatsModal = (statName, stats) => {
  return new Promise((resolve, reject) => {
    Dialog.create({
      component: DailyStatModal,
      componentProps: { statName: statName, stats: stats },
    })
      .onOk((stats) => {
        resolve(stats);
      })
      .onCancel(() => {
        resolve();
      })
      .onDismiss(() => {});
  });
};

// Regular expressions for validating daily stats values
const bgRegex = new RegExp('^\\d{1,2}([.]\\d)?$');
const systolicRegex = new RegExp('^1\\d{2}$');
const diastolicRegex = new RegExp('^\\d{2,3}$');
const sleepRegex = new RegExp('^1?\\d([.]\\d{1,2})?$');
const bodyWeightRegex = new RegExp('^\\d{2,3}$');
const descriptionRegex = new RegExp('^[\\da-zA-Z -"\',@%$&]{5,256}$');
const nutrientRegex = new RegExp(`^\\d{0,3}$`);
const scaleOfFiveRegex = new RegExp(`^[1-5]$`);

/**
 * Generates a function for validating values against a regular expression.
 * @param {RegExp} regex The regular expression.
 * @returns A function that takes a value to validate. Returns true when valid, otherwise false.
 */
const fieldValidatorFactory = (regex) => {
  return (val) => {
    const result = regex.test(val) || 'Invalid value.';
    return result;
  };
};

// Create validator functions for each type of daily stat
export const bgValidator = fieldValidatorFactory(bgRegex);
export const systolicValidator = fieldValidatorFactory(systolicRegex);
export const diastolicValidator = fieldValidatorFactory(diastolicRegex);
export const sleepValidator = fieldValidatorFactory(sleepRegex);
export const bodyWeightValidator = fieldValidatorFactory(bodyWeightRegex);
export const foodDescriptionValidator = fieldValidatorFactory(descriptionRegex);
export const foodNutrientValidator = fieldValidatorFactory(nutrientRegex);
export const scaleOfFiveValidator = fieldValidatorFactory(scaleOfFiveRegex);

/**
 * Stores a daily stat value to the back end.
 * @param {Object} stat An object that contains the date of the stat and a value for the stat. See emptyStats() for the object structure.
 */
export const saveDailyStat = async (stat) => {
  const url = `/homegym/api/dailystats/?date=${stat.date}`;

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const options = {
    method: 'POST',
    body: JSON.stringify(stat),
    headers: headers,
  };

  const resp = await fetch(url, options);
  if (resp.status >= 200 && resp.status < 300) {
    utils.toast('Saved', 'positive');
  } else if (resp.status == 401) {
    console.log('unauthorized upsert of daily stat');
    await utils.authPromptAsync();
    saveDailyStat(stat);
  } else if (resp.status < 200 && resp.status >= 300) {
    const errBody = await resp.json();
    console.log(errBody.message);
    utils.toast('Error', 'negative');
  }
};

/**
 * Retrieves a page of daily stats from the backend.
 * @param {number} startDate (Optional) The start date of the page. No value starts the page at the first daily stat.
 * @param {number} endDate (Optional) The end date of the page. No value fills the page.
 * @param {number} pageSize (Optional) The number of items in the page.
 * @returns An array of dailyStats objects.
 */
const fetchDailyStatsPage = async (startDate, endDate, pageSize) => {
  let url = `/homegym/api/dailystats/`;
  const qParams = [];
  if (startDate) {
    qParams.push(`start=${startDate}`);
  }
  if (endDate) {
    qParams.push(`end=${endDate}`);
  }
  if (pageSize) {
    qParams.push(`pagesize=${pageSize}`);
  }

  if (qParams.length > 0) {
    url += `?${qParams.join('&')}`;
  }

  const resp = await fetch(url, {
    method: 'GET',
    mode: 'same-origin',
  });

  if (resp.status == 401) {
    throw new utils.ErrNotLoggedIn('unauthorized fetch of daily stats page');
  }
  const stats = await resp.json();
  return stats;
};

/**
 * Retrieves all daily stats from the back end over a period of time. In the background, stats are retrieved in pages.
 * @param {number} startDate (Optional) The start date of the page. No value starts the page at the first daily stat.
 * @param {number} endDate (Optional) The end date of the page. No value fills the page.
 * @param {number} pageSize (Optional) The number of items in the page.
 * @returns An array of stats retrieved for the defined period of time.
 */
export const fetchDailyStats = async (startDate, endDate, pageSize) => {
  let done = false;
  let stats = [];
  while (!done) {
    const statsPage = await fetchDailyStatsPage(startDate, endDate, pageSize);
    if (statsPage.length > 0) {
      stats = stats.concat(statsPage);
      if (statsPage.length < pageSize) {
        done = true;
      } else {
        startDate = stats[stats.length - 1].date;
      }
    }
    done = true;
  }
  return stats;
};

/**
 * Amalgamates a series of disparate daily stat objects into single objects, grouped date.
 * @param {[Object]} dailyStats An array of dailyStat objects that contain only one property value
 * @returns An array of objects, one for each day that a stat was collected, containing all stats for that day.
 */
export const toDayBuckets = (dailyStats) => {
  const dayBuckets = {};

  for (let i = 0; i < dailyStats.length; i++) {
    const midnight = dateUtils.setEpochToMidnight(dailyStats[i].date) * 1000;
    let dayBucket;
    if (dayBuckets[midnight]) {
      dayBucket = dayBuckets[midnight];
    } else {
      dayBucket = { date: midnight, sequential: [] };
    }
    if (
      dailyStats[i].sleep ||
      dailyStats[i].bodyweight ||
      dailyStats[i].mood ||
      dailyStats[i].energy ||
      dailyStats[i].stress
    ) {
      if (dailyStats[i].sleep) {
        dayBucket.sleep = dailyStats[i].sleep;
      }
      if (dailyStats[i].bodyweight) {
        dayBucket.bodyweight = dailyStats[i].bodyweight;
      }
      if (dailyStats[i].mood) {
        dayBucket.mood = dailyStats[i].mood;
      }
      if (dailyStats[i].energy) {
        dayBucket.energy = dailyStats[i].energy;
      }
      if (dailyStats[i].stress) {
        dayBucket.stress = dailyStats[i].stress;
      }
    } else {
      dayBucket.sequential = dayBucket.sequential.concat(dailyStats[i]);
    }

    dayBuckets[midnight] = dayBucket;
  }

  return dayBuckets;
};

/**
 * Generates a chart for daily stats.
 * @param {String} element The id of the HTML element in which the chart is inserted.
 * @param {Number} startDate The first date of the charted data, in seconds since epoch.
 * @param {Number} endDate The last date of the charted data, in seconds since epoch.
 * @param {[Object]} dataset An array of dailyStats objects, one object for each day that a stat was collected, each object containing all stats for that day.
 * @returns The Chart object.
 */
export const getDailyChart = (element, startDate, endDate, dataset) => {
  const existing = Chart.getChart(element);
  if (existing) {
    existing.destroy();
  }

  return new Chart(element, {
    type: 'line',
    data: {
      datasets: [
        {
          label: 'Sleep',
          data: Object.values(dataset),
          parsing: {
            xAxisKey: 'date',
            yAxisKey: 'sleep',
          },
          yAxisID: 'yBio',
        },
        {
          label: 'Mood',
          data: Object.values(dataset),
          parsing: { xAxisKey: 'date', yAxisKey: 'mood' },
          yAxisID: 'yBio',
        },
        {
          label: 'Energy',
          data: Object.values(dataset),
          parsing: { xAxisKey: 'date', yAxisKey: 'energy' },
          yAxisID: 'yBio',
        },
        {
          label: 'Stress',
          data: Object.values(dataset),
          yAxisID: 'yBio',
          parsing: { xAxisKey: 'date', yAxisKey: 'stress' },
        },
        {
          label: 'BodyWeight',
          data: Object.values(dataset),
          parsing: {
            xAxisKey: 'date',
            yAxisKey: 'bodyweight',
          },
          yAxisID: 'yBW',

          spanGaps: true,
        },
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
        yBW: {
          position: 'right',
          min: 0,
          title: {
            text: 'Body Weight',
          },
        },
        yBio: {
          position: 'left',
          min: 0,
          max: 12,
          title: {
            text: 'Sleep, Mood, Stress, Energy',
          },
        },
      },
    },
  });
};

/**
 * Amalgamates a series of disparate time series stat objects into single objects, grouped date.
 * Time series stats are those that can occur multiple times a day.
 * @param {[Object]} dailyStats An array of dailyStats objects.
 * @returns An object that contains a property for each time series stat. Each property is an array of data formatted for graphing.
 */
export const getTimeSeriesDataSets = (dailyStats) => {
  const bgData = new Array();
  const foodData = new Array();
  const bpData = new Array();

  dailyStats.forEach((stat) => {
    const x = dateUtils.setEpochToMidnight(stat.date) * 1000;

    // stat times normalized to be within 24h of epoch
    const offsetMS = new Date().getTimezoneOffset() * 60 * 1000;
    const y = stat.date * 1000 - x + offsetMS;

    const statDate = new Date(y);
    const timeStr = dateUtils.formatTime(statDate);

    if (stat.bg) {
      bgData.push({
        x: x,
        y: y,
        description: [stat.bg],
        timeStr: timeStr,
      });
    }
    if (stat.bp && stat.bp[0] && stat.bp[1]) {
      bpData.push({
        x: x,
        y: y,
        description: [`${stat.bp[0]}/${stat.bp[1]}`],
        timeStr: timeStr,
      });
    }

    if (
      stat.food.description ||
      stat.food.protein ||
      stat.food.fat ||
      stat.food.carbs
    ) {
      foodData.push({
        x: x,
        y: y,
        description: [
          stat.food.description,
          `protein: ${stat.food.protein}`,
          `fat: ${stat.food.fat}`,
          `carbs: ${stat.food.carbs}`,
        ],
        timeStr: timeStr,
      });
    }
  });
  return { bg: bgData, bp: bpData, food: foodData };
};

/**
 * Generates a chart for time series stats.
 * @param {String} element The id of the HTML element in which the chart is inserted.
 * @param {Number} startDate The first date of the charted data, in seconds since epoch.
 * @param {Number} endDate The last date of the charted data, in seconds since epoch.
 * @param {[Object]} datasetsObj An array of Objects that contains the data to chart.
 * @returns The new Chart object.
 */
export const getTimeSeriesChart = (
  element,
  startDate,
  endDate,
  datasetsObj
) => {
  const existing = Chart.getChart(element);
  if (existing) {
    existing.destroy();
  }
  const datasets = new Array();

  for (const [key, value] of Object.entries(datasetsObj)) {
    const ds = {
      label: key,
      data: value,
    };
    if (key == 'bg') {
      ds.datalabels = {
        display: true,
        formatter: function (value, context) {
          return Number(value.description[0]) > 6 ? value.description : null;
        },
        color: '#FF0000',
        align: 'right',
        offset: 2,
      };
    }
    datasets.push(ds);
  }

  const chart = new Chart(element, {
    plugins: [ChartDataLabels],
    type: 'scatter',
    data: {
      datasets: datasets,
    },
    options: {
      maintainAspectRatio: false,
      aspectRatio: 1,
      elements: {
        point: {
          radius: 3,
        },
      },
      scales: {
        x: {
          type: 'time',
          max: startDate * 1000,
          min: endDate * 1000,
        },
        y: {
          type: 'time',
          time: { unit: 'hour' },
          min: '1970-01-01T01',
          max: '1970-01-01T23',
        },
      },
      plugins: {
        tooltip: {
          callbacks: {
            title: (context) => {
              if (context[0].raw.timeStr) {
                return context[0].raw.timeStr;
              }
            },
            label: (context) => {
              let label = context.dataset.label || '';

              if (label) {
                label += ': ';
              }
              if (context.raw.description) {
                label += context.raw.description;
              }
              return label;
            },
          },
        },
        datalabels: {
          display: false,
        },
      },
      onClick: (e) => {
        const canvasPosition = getRelativePosition(e, chart);

        // Substitute the appropriate scale IDs
        const dataX = chart.scales.x.getValueForPixel(canvasPosition.x);
        const dataY = chart.scales.y.getValueForPixel(canvasPosition.y);
        console.log(`clicked: (${dataX}, ${dataY})`);
        console.log(chart.tooltip.dataPoints);
        console.log(e);
      },
    },
  });

  return chart;
};
