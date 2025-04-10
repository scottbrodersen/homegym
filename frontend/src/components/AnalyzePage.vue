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
  import { QExpansionItem } from 'quasar';

  // default date range is 16 weeks since now
  const endDate = ref(dateUtils.nowInSeconds());
  const startDate = ref(endDate.value - 16 * 7 * 24 * 60 * 60);
  let exerciseTypes = [];

  let rawMetrics = { dates: [], load: [], volume: [] };
  const metrics = ref({});
  const dailyStats = ref([]);

  const updateStartDate = (newStartDate) => {
    if (newStartDate > endDate.value) {
      utils.openConfirmModal('Start date must occur before end date');
    } else {
      startDate.value = newStartDate;
    }
    getMetrics();
  };

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

  const dailyStatsChart = () => {
    const dayBuckets = dailyStatsUtils.toDayBuckets(dailyStats.value);
    dailyStatsUtils.getDailyChart(
      document.getElementById('chartdaily'),
      endDate.value,
      startDate.value,
      dayBuckets
    );
  };

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

  const setExerciseTypes = (types) => {
    exerciseTypes = types;
    getMetrics(true);
  };

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
