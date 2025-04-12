import { expect, test } from 'vitest';
import * as dsu from '../src/modules/dailyStatsUtils';

describe('dailyStatsUtils', () => {
  test('correctly tests blood glucose values', () => {
    expect(dsu.bgValidator(5.0)).toBeTruthy;
    expect(dsu.bgValidator(5)).toBeTruthy;
    expect(dsu.bgValidator(5.12)).toBe('Invalid value.');
    expect(dsu.bgValidator('text')).toBe('Invalid value.');
  });

  test('correctly tests systolic pressure values', () => {
    expect(dsu.systolicValidator(123)).toBe(true);
    expect(dsu.systolicValidator(1230)).toBe('Invalid value.');
    expect(dsu.systolicValidator(5.0)).toBe('Invalid value.');
    expect(dsu.systolicValidator(5)).toBe('Invalid value.');
    expect(dsu.systolicValidator(250)).toBe('Invalid value.');
    expect(dsu.systolicValidator('text')).toBe('Invalid value.');
  });

  test('correctly tests diastolic pressure values', () => {
    expect(dsu.diastolicValidator(12)).toBe(true);
    expect(dsu.diastolicValidator(135)).toBe(true);
    expect(dsu.diastolicValidator(1230)).toBe('Invalid value.');
    expect(dsu.diastolicValidator(5.0)).toBe('Invalid value.');
    expect(dsu.diastolicValidator(5)).toBe('Invalid value.');
    expect(dsu.diastolicValidator('text')).toBe('Invalid value.');
  });

  test('correctly tests sleep values', () => {
    expect(dsu.sleepValidator(5.0)).toBe(true);
    expect(dsu.sleepValidator(15.0)).toBe(true);
    expect(dsu.sleepValidator(5)).toBe(true);
    expect(dsu.sleepValidator(5.12)).toBe(true);
    expect(dsu.sleepValidator(25.0)).toBe('Invalid value.');
    expect(dsu.sleepValidator(150)).toBe('Invalid value.');
    expect(dsu.sleepValidator('text')).toBe('Invalid value.');
  });

  test('correctly tests body weight values', () => {
    expect(dsu.bodyWeightValidator(12)).toBe(true);
    expect(dsu.bodyWeightValidator(135)).toBe(true);
    expect(dsu.bodyWeightValidator(1230)).toBe('Invalid value.');
    expect(dsu.bodyWeightValidator(123.4)).toBe('Invalid value.');
    expect(dsu.bodyWeightValidator(5)).toBe('Invalid value.');
    expect(dsu.bodyWeightValidator('text')).toBe('Invalid value.');
  });

  test('correctly tests description values', () => {
    expect(dsu.foodDescriptionValidator('bLahblah 0,@%$&')).toBe(true);
    expect(dsu.foodDescriptionValidator('sho')).toBe('Invalid value.');
    expect(
      dsu.foodDescriptionValidator(
        '256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256charac'
      )
    ).toBe(true);
    expect(
      dsu.foodDescriptionValidator(
        'morethan256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters256characters25'
      )
    ).toBe('Invalid value.');
  });

  test('correctly tests nutrient values', () => {
    expect(dsu.foodNutrientValidator('')).toBe(true);
    expect(dsu.foodNutrientValidator('1')).toBe(true);
    expect(dsu.foodNutrientValidator('12')).toBe(true);
    expect(dsu.foodNutrientValidator('123')).toBe(true);
    expect(dsu.foodNutrientValidator('1234')).toBe('Invalid value.');
    expect(dsu.foodNutrientValidator('text')).toBe('Invalid value.');
    expect(dsu.foodNutrientValidator('1.0')).toBe('Invalid value.');
  });
  test('correctly tests scale of five values', () => {
    expect(dsu.scaleOfFiveValidator('1')).toBe(true);
    expect(dsu.scaleOfFiveValidator('2')).toBe(true);
    expect(dsu.scaleOfFiveValidator('3')).toBe(true);
    expect(dsu.scaleOfFiveValidator('4')).toBe(true);
    expect(dsu.scaleOfFiveValidator('5')).toBe(true);
    expect(dsu.scaleOfFiveValidator('0')).toBe('Invalid value.');
    expect(dsu.scaleOfFiveValidator('6')).toBe('Invalid value.');
    expect(dsu.scaleOfFiveValidator('text')).toBe('Invalid value.');
  });

  const testMood = 1;
  const testStress = 4;
  const testEnergy = 3;
  const testFat = 3;
  const testFibre = 4;
  const testCarbs = 5;
  const testProtein = 6;
  const testFoodDescription = 'test meal';
  const testBP = [135, 80];
  const testBodyWeight = 130;
  const testSleep = 7.0;
  const testBloodGlucose = 6.1;
  const testFood = {
    protein: testProtein,
    carbs: testCarbs,
    fat: testFat,
    fiber: testFibre,
    description: testFoodDescription,
  };

  const testDates = [1741773971, 1754452371];

  const testDailyStat = (date, numMeals) => {
    const stat = [
      {
        date: date,
        bp: [0, 0],
        food: testFood,
      },
      {
        date: date + 60,
        food: {
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
        },
        mood: testMood,
        stress: testStress,
        energy: testEnergy,
      },
      {
        date: date + 120,
        bp: testBP,
        food: {
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
        },
      },
      {
        date: date + 180,
        bg: testBloodGlucose,
        bp: [0, 0],
        food: {
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
        },
      },
      {
        date: date + 240,
        bp: [0, 0],
        food: {
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
        },
        bodyweight: testBodyWeight,
      },
      {
        date: date + 320,
        bp: [0, 0],
        sleep: testSleep,
        food: {
          protein: 0,
          carbs: 0,
          fat: 0,
          fiber: 0,
        },
      },
    ];

    // add meals as needed
    for (let i = 1; i < numMeals; i++) {
      stat.push({
        date: date + i * 360,
        bp: [0, 0],
        food: testFood,
      });
    }
    return stat;
  };

  test('Correctly split one day of daily stats with one meal into day buckets', () => {
    const oneDayStats = testDailyStat(testDates[0], 1);
    const oneDayBuckets = dsu.toDayBuckets(oneDayStats);
    const midnights = Object.keys(oneDayBuckets);

    expect(midnights.length).toBe(1);

    const bucket = oneDayBuckets[midnights[0]];

    expect(bucket.sleep).toBe(testSleep);
    expect(bucket.bodyweight).toBe(testBodyWeight);
    expect(bucket.mood).toBe(testMood);
    expect(bucket.energy).toBe(testEnergy);
    expect(bucket.stress).toBe(testStress);
    expect(bucket.sequential.length).toBe(3);
  });

  test('Correctly split one day of daily stats with 2 meals into day buckets', () => {
    const oneDayStats = testDailyStat(testDates[0], 2);
    const oneDayBuckets = dsu.toDayBuckets(oneDayStats);
    const midnights = Object.keys(oneDayBuckets);

    expect(midnights.length).toBe(1);

    const bucket = oneDayBuckets[midnights[0]];

    expect(bucket.sleep).toBe(testSleep);
    expect(bucket.bodyweight).toBe(testBodyWeight);
    expect(bucket.mood).toBe(testMood);
    expect(bucket.energy).toBe(testEnergy);
    expect(bucket.stress).toBe(testStress);
    expect(bucket.sequential.length).toBe(4);
  });

  test('Correctly split 2 days of daily stats with 1 meal into day buckets', () => {
    const twoDayStats = testDailyStat(testDates[0], 1);
    const dayTwoStats = testDailyStat(testDates[1], 1);
    for (const stat of dayTwoStats) {
      twoDayStats.push(stat);
    }

    const buckets = dsu.toDayBuckets(twoDayStats);

    const midnights = Object.keys(buckets);

    expect(midnights.length).toBe(2);

    for (const midnight of midnights) {
      const bucket = buckets[midnight];

      //console.log(bucket);
      expect(bucket.sleep).toBe(testSleep);
      expect(bucket.bodyweight).toBe(testBodyWeight);
      expect(bucket.mood).toBe(testMood);
      expect(bucket.energy).toBe(testEnergy);
      expect(bucket.stress).toBe(testStress);
      expect(bucket.sequential.length).toBe(3);
    }
  });
});
