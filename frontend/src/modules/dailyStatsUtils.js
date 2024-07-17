import DailyStatModal from '../components/DailyStatModal.vue';
import { Dialog } from 'quasar';
import * as utils from './utils';
import * as dateUtils from '../modules/dateUtils';
import Chart from 'chart.js/auto';
import 'date-fns';
import 'chartjs-adapter-date-fns';
import 'date-fns';

export const statNames = [
  'bg',
  'bp',
  'bodyweight',
  'food',
  'sleep',
  'mood',
  'stress',
  'energy',
];

export const spiritGroup = ['mood', 'stress', 'energy'];

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

const bgRegex = new RegExp('^[0-9]{1,2}[.][0-9]$');
const systolicRegex = new RegExp('^1[0-9]{2}$');
const diastolicRegex = new RegExp('^[0-9]{2,3}$');
const sleepRegex = new RegExp('^1?[0-9][.][0-9]$');
const bodyWeightRegex = new RegExp('^[0-9]{3}');
const descriptionRegex = new RegExp('^[0-9a-zA-Z -"\',@%$&]{5,256}$');
const nutrientRegex = new RegExp(`^[0-9]{0,3}$`);
const scaleOfFiveRegex = new RegExp(`^[1-5]$`);

const fieldValidatorFactory = (regex) => {
  return (val) => {
    const result = regex.test(val) || 'Invalid value.';
    return result;
  };
};

export const bgValidator = fieldValidatorFactory(bgRegex);
export const systolicValidator = fieldValidatorFactory(systolicRegex);
export const diastolicValidator = fieldValidatorFactory(diastolicRegex);
export const sleepValidator = fieldValidatorFactory(sleepRegex);
export const bodyWeightValidator = fieldValidatorFactory(bodyWeightRegex);
export const foodDescriptionValidator = fieldValidatorFactory(descriptionRegex);
export const foodNutrientValidator = fieldValidatorFactory(nutrientRegex);
export const scaleOfFiveValidator = fieldValidatorFactory(scaleOfFiveRegex);

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

      if (stat.bp && stat.bp[0] && stat.bp[1]) {
        bpData.push({
          x: x,
          y: y,
          description: [`${stat.bp[0]}/${stat.bp[1]}`],
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
            stat.description,
            `protein: ${stat.protein}`,
            `fat: ${stat.fat}`,
            `carbs: ${stat.carbs}`,
          ],
        });
      }
    }
  });
  return { bg: bgData, bp: bpData, food: foodData };
};

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
    datasets.push({
      label: key,
      data: value,
    });
  }

  return new Chart(element, {
    type: 'scatter',
    data: {
      datasets: datasets,
    },
    options: {
      elements: {
        point: {
          radius: 3,
        },
      },
      scales: {
        x: {
          type: 'time',
        },
        y: {
          type: 'time',
          time: { unit: 'hour' },
          min: '00:00:00',
          max: '23:59:59',
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
      },
    },
  });
};
