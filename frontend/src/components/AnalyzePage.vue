<script setup>
  /**
   * Displays data for the metrics that we track:
   *  - exercise volume
   *  - daily health markers
   *  - food intake and blood glucose
   *
   * Exercise volume can be filtered by exercise type and date range.
   */
  import { onMounted, ref } from 'vue';
  import * as utils from '../modules/utils';
  import * as dateUtils from '../modules/dateUtils';
  import * as styles from '../style.module.css';
  import * as dailyStatsUtils from '../modules/dailyStatsUtils';
  import DatePicker from './DatePicker.vue';
  import * as analyzeUtils from '../modules/analyzeUtils';
  import Chart from 'chart.js/auto';
  import ExerciseFilter from './ExerciseFilter.vue';
  import { QExpansionItem } from 'quasar';

  // default date range is 16 weeks since now
  const endDate = ref(dateUtils.nowInSeconds());
  const startDate = ref(endDate.value - 16 * 7 * 24 * 60 * 60);
  let exerciseTypes = [];

  let rawMetrics = { dates: [], load: [], volume: [] };
  // chart data for exercise metrics
  const metrics = ref({});
  // chart data for daily stats
  const dailyStats = ref([]);

  /**
   * Sets the start date of the charted time period.
   * @param newStartDate The start date in seconds since epoch.
   */
  const updateStartDate = (newStartDate) => {
    if (newStartDate > endDate.value) {
      utils.openConfirmModal('Start date must occur before end date');
    } else {
      startDate.value = newStartDate;
    }
    getMetrics();
  };

  /**
   * Sets the end date of the charted time period.
   * @param newEndDate The end date in seconds since epoch.
   */
  const updateEndDate = (newEndDate) => {
    if (newEndDate < startDate.value) {
      utils.openConfirmModal('End date must occur after start date');
    } else {
      endDate.value = newEndDate;
    }
    getMetrics();
  };

  Chart.defaults.color = 'rgb(252,252,252)';
  Chart.defaults.elements.line.borderWidth = 1;
  Chart.defaults.elements.point.radius = 1;
  Chart.defaults.plugins.legend.labels.boxHeight = 1;
  Chart.defaults.plugins.legend.labels.color = 'rgb(252,252,252)';
  Chart.defaults.scale.title.display = true;
  Chart.defaults.scale.border.color = 'rgb(72,72,72)';
  Chart.defaults.scales.time.time.unit = 'day';
  Chart.defaults.scales.time.time.tooltipFormat = 'MMM d, yyyy';

  /**
   * Creates a graph of total load/volume and total load over time for a set of performed exercises
   */
  const exerciseChart = () => {
    const lvRatioData = Array();
    const loadData = Array();

    for (let i = 0; i < metrics.value.dates.length; i++) {
      lvRatioData.push({
        x: metrics.value.dates[i] * 1000,
        y: metrics.value.lvRatio[i],
      });
      loadData.push({
        x: metrics.value.dates[i] * 1000,
        y: metrics.value.load[i],
      });
    }

    analyzeUtils.getVolumeChart(
      document.getElementById('chartvolume'),
      endDate.value,
      startDate.value,
      metrics.value.dates,
      lvRatioData,
      loadData
    );
  };

  /**
   * Creates a graph of daily stats (e.g. sleep, mood, etc) over time
   */
  const dailyStatsChart = () => {
    const dayBuckets = dailyStatsUtils.toDayBuckets(dailyStats.value);
    dailyStatsUtils.getDailyChart(
      document.getElementById('chartdaily'),
      endDate.value,
      startDate.value,
      dayBuckets
    );
  };

  /**
   * Creates a graph of time series data (i.e. measurements that can occur multiple times in a day such as blood pressure, food intake, etc)
   */
  const timeSeriesChart = async () => {
    const datasets = dailyStatsUtils.getTimeSeriesDataSets(dailyStats.value);
    const eventDataset = await analyzeUtils.getTimeSeriesData(
      endDate.value,
      startDate.value
    );

    datasets.events = eventDataset;

    dailyStatsUtils.getTimeSeriesChart(
      document.getElementById('chartglucose'),
      endDate.value,
      startDate.value,
      datasets
    );
  };

  /**
   * Sets the exercise types for which metrics are graphed.
   * @param types An array of exercise type IDs.
   */
  const setExerciseTypes = (types) => {
    exerciseTypes = types;
    getMetrics(true);
  };

  /**
   * Retrieves exercise metrics from the server and creates the exercise metrics, daily stats, and time series charts.
   * @param {boolean} updatedTypes True if the chart is updated due to a change in the exercise filter.
   */
  const getMetrics = async (updatedTypes = false) => {
    resetMetrics();
    try {
      if (endDate.value && startDate.value) {
        if (exerciseTypes.length > 0) {
          rawMetrics = await analyzeUtils.fetchMetrics(
            endDate.value,
            startDate.value,
            exerciseTypes
          );

          metrics.value = analyzeUtils.getDailyTotals(rawMetrics);
          exerciseChart();
        }

        // Get daily stats only on date changes
        if (!updatedTypes) {
          dailyStats.value = await dailyStatsUtils.fetchDailyStats(
            endDate.value,
            startDate.value
          );
          dailyStatsChart();
          timeSeriesChart();
        }
      }
    } catch (e) {
      if (e instanceof utils.ErrNotLoggedIn) {
        console.log(e.message);
        await utils.authPromptAsync();
        await getMetrics();
      } else {
        console.log(e);
      }
    }
  };

  const resetMetrics = () => {
    metrics.value = { dates: [], lvRatios: [], load: [] };
  };

  onMounted(async () => {
    await getMetrics();
  });
</script>
<template>
  <div :class="[styles.analyzePage]">
    <div :class="[styles.analyzeDates]">
      <div>
        <div>Start Date:</div>
        <DatePicker
          :dateValue="startDate"
          :hideTime="true"
          @update="updateStartDate"
        />
      </div>
      <div>
        <div>End Date:</div>
        <DatePicker
          :dateValue="endDate"
          :hideTime="true"
          @update="updateEndDate"
        />
      </div>
    </div>
    <h3>Exercise Volume</h3>
    <q-expansion-item
      :class="[styles.analyzeFilter]"
      label="Exercise Filters"
      dark
      dense
    >
      <ExerciseFilter @ids="(val) => setExerciseTypes(val)" />
    </q-expansion-item>
    <div :class="[styles.analyzeChart]">
      <canvas id="chartvolume"></canvas>
    </div>
    <h3>Health Markers</h3>

    <div :class="[styles.analyzeChart]">
      <canvas id="chartdaily"></canvas>
    </div>

    <h3>Blood Glucose Markers</h3>

    <div :class="[styles.analyzeChart]">
      <canvas id="chartglucose"></canvas>
    </div>
  </div>
</template>
