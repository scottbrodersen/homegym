<script setup>
  import { onMounted, ref } from 'vue';
  import * as utils from '../modules/utils';
  import * as dateUtils from '../modules/dateUtils';
  import * as state from '../modules/state';
  import DatePicker from './DatePicker.vue';
  import * as analyzeUtils from '../modules/analyzeUtils';
  import { QBtn, QSelect } from 'quasar';
  import Chart from 'chart.js/auto';
  import 'date-fns';
  import 'chartjs-adapter-date-fns';
  import 'date-fns';
  import ExerciseFilter from './ExerciseFilter.vue';

  // default date range is 4 weeks since now
  const startDate = ref(dateUtils.nowInSecondsUTC());
  const endDate = ref(startDate.value - 28 * 24 * 60 * 60);
  let exerciseTypes = [];

  const metrics = ref([]);

  let volChart;

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

  const chart = async () => {
    if (typeof volChart == 'object' && volChart.hasOwnProperty('id')) {
      volChart.destroy();
    }
    const volData = Array();
    const loadData = Array();

    for (let i = 0; i < metrics.value.dates.length; i++) {
      volData.push({
        x: metrics.value.dates[i] * 1000,
        y: metrics.value.volume[i],
      });
      loadData.push({
        x: metrics.value.dates[i] * 1000,
        y: metrics.value.load[i],
      });
    }

    volChart = new Chart(document.getElementById('volume'), {
      type: 'line',
      data: {
        labels: metrics.value.dates,
        datasets: [
          { data: volData, yAxisID: 'yVol', label: 'Volume' },
          { data: loadData, yAxisID: 'yLoad', label: 'Load' },
        ],
      },
      options: {
        scales: {
          x: {
            type: 'time',
            time: {
              unit: 'day',
            },
            min: volData[volData.length - 1].x,
            max: volData[0].x,
          },
          yLoad: {
            position: 'right',
          },
          yVol: {
            position: 'left',
          },
        },
      },
    });
  };

  const setExerciseTypes = (types) => {
    exerciseTypes = types;
    getMetrics();
  };
  const getMetrics = async () => {
    try {
      if (startDate.value && endDate.value) {
        metrics.value = await analyzeUtils.fetchMetrics(
          startDate.value,
          endDate.value,
          exerciseTypes
        );
      }
    } catch (e) {
      if (e instanceof utils.ErrNotLoggedIn) {
        console.log(e.message);
        await utils.authPromptAsync();
        getMetrics();
      } else {
        console.log(e);
      }
    }
    chart();
  };
  onMounted(async () => {
    await getMetrics();
  });
</script>
<template>
  <div>
    <DatePicker
      :dateValue="startDate"
      :hideTime="true"
      @update="updateStartDate"
    />
    <DatePicker :dateValue="endDate" :hideTime="true" @update="updateEndDate" />
    <ExerciseFilter @ids="(val) => setExerciseTypes(val)" />
    <canvas id="volume"></canvas>
  </div>
</template>
