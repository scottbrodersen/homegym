<script setup>
  import { onMounted, ref } from 'vue';
  import * as utils from '../modules/utils';
  import * as dateUtils from '../modules/dateUtils';
  import * as styles from '../style.module.css';
  import * as dailyStatsUtils from '../modules/dailyStatsUtils';
  import DatePicker from './DatePicker.vue';
  import * as analyzeUtils from '../modules/analyzeUtils';
  import Chart from 'chart.js/auto';
  import ExerciseFilter from './ExerciseFilter.vue';

  // default date range is 16 weeks since now
  const startDate = ref(dateUtils.nowInSeconds());
  const endDate = ref(startDate.value - 16 * 7 * 24 * 60 * 60);
  let exerciseTypes = [];

  let rawMetrics = { dates: [], load: [], volume: [] };
  const metrics = ref({});
  const dailyStats = ref([]);

  const updateStartDate = (newDate) => {
    if (newDate <= endDate.value) {
      utils.openConfirmModal('Start date must occur after end date');
    } else {
      startDate.value = newDate;
    }
    getMetrics();
  };

  const updateEndDate = (newDate) => {
    if (newDate >= startDate.value) {
      utils.openConfirmModal('End date must occur before start date');
    } else {
      endDate.value = newDate;
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
      startDate.value,
      endDate.value,
      metrics.value.dates,
      lvRatioData,
      loadData
    );
  };

  const dailyStatsChart = () => {
    const dayBuckets = dailyStatsUtils.toDayBuckets(dailyStats.value);
    dailyStatsUtils.getDailyChart(
      document.getElementById('chartdaily'),
      startDate.value,
      endDate.value,
      dayBuckets
    );
  };

  const timeSeriesChart = () => {
    const datasets = dailyStatsUtils.getTimeSeriesDataSets(dailyStats.value);
    dailyStatsUtils.getTimeSeriesChart(
      document.getElementById('charttimeseries'),
      startDate.value,
      endDate.value,
      datasets
    );
  };

  const setExerciseTypes = (types) => {
    exerciseTypes = types;
    getMetrics(true);
  };

  const getMetrics = async (updatedTypes = false) => {
    resetMetrics();
    try {
      if (startDate.value && endDate.value) {
        if (exerciseTypes.length > 0) {
          rawMetrics = await analyzeUtils.fetchMetrics(
            startDate.value,
            endDate.value,
            exerciseTypes
          );

          metrics.value = analyzeUtils.getDailyTotals(rawMetrics);
          exerciseChart();
        }

        // Get daily stats only on date changes
        if (!updatedTypes) {
          dailyStats.value = await dailyStatsUtils.fetchDailyStats(
            startDate.value,
            endDate.value
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
    <ExerciseFilter @ids="(val) => setExerciseTypes(val)" />
    <canvas id="chartvolume"></canvas>
    <canvas id="chartdaily"></canvas>
    <canvas id="charttimeseries"></canvas>
  </div>
</template>
