<script setup>
  /**
   * A calendar that flags the days on which a workout occurs for a program instance.
   *
   * Props:
   *  Coords is a 3x1 array that holds the coordinates of the workout for a date.
   *      E.g. [0,1,2] denotes the workout in the 3rd day of the 2nd microcycle in the 1st block.
   *  instance is a program instance object.
   *
   * Emits the index of the workout that occurs on the selected date.
   */
  import { onMounted, ref, watch } from 'vue';
  import { QDate } from 'quasar';
  import * as dateUtils from '../modules/dateUtils';
  import * as programUtils from '../modules/programUtils';

  const props = defineProps({
    coords: Object,
    instance: Object,
  });
  const emit = defineEmits(['dayIndex']);

  let nonRestDays = [];
  let formattedNonRestDays = [];

  let workoutDates = [];
  let formattedWorkoutDates = [];

  const date = ref();

  const init = () => {
    nonRestDays = programUtils.getNonRestDates(props.instance);
    formattedNonRestDays = formatForQuasar(nonRestDays);

    workoutDates = programUtils.getInstanceWorkoutDates(props.instance);
    formattedWorkoutDates = formatForQuasar(workoutDates);

    date.value = dateUtils.formatDate(
      dateUtils.dateFromSeconds(props.instance.startDate)
    );
  };

  let dayIndex;

  const emitDayIndex = (date) => {
    dayIndex = workoutDates.findIndex((workoutDate) => workoutDate == date);
    if (needSync()) {
      emit('dayIndex', dayIndex);
    }
  };

  function formatForQuasar(dateArr) {
    const formatted = new Array();
    dateArr.forEach((dateStr) => {
      formatted.push(dateStr.replaceAll('-', '/'));
    });
    return formatted;
  }

  watch(
    () => props.coords,
    (newCoords) => {
      date.value = dateUtils.formatDate(
        programUtils.getWorkoutDate(props.instance, newCoords)
      );
    },
    { deep: true }
  );

  watch(
    () => props.instance,
    (newInstance) => {
      init();
    }
  );

  onMounted(() => {
    init();
    emitDayIndex(date.value);
  });

  const needSync = () => {
    const indexFromCoords = props.coords
      ? programUtils.getDayIndex(props.instance, props.coords)
      : undefined;
    return indexFromCoords != dayIndex;
  };
</script>
<template>
  <div>
    <q-date
      v-model="date"
      :events="formattedNonRestDays"
      :options="formattedWorkoutDates"
      :mask="dateUtils.dateMask"
      event-color="primary"
      dark
      flat
      minimal
      @update:model-value="(value) => emitDayIndex(value)"
    />
  </div>
</template>
