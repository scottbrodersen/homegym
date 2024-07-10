import DailyStatModal from '../components/DailyStatModal.vue';
import { Dialog } from 'quasar';
import * as utils from './utils';

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
    await authPromptAsync();
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
