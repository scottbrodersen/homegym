<script setup>
  import { onMounted, ref, watch } from 'vue';
  import { QDate } from 'quasar';
  import * as dateUtils from '../modules/dateUtils';
  import * as programUtils from '../modules/programUtils';
  import ProgramMap from './ProgramMap.vue';

  const props = defineProps({
    coords: Object,
    instance: Object,
  });
  const emit = defineEmits(['dayIndex']);

  const nonRestDays = programUtils.getNonRestDates(props.instance);
  const formattedNonRestDays = formatForQuasar(nonRestDays);

  const workoutDates = programUtils.getInstanceWorkoutDates(props.instance);
  const formattedWorkoutDates = formatForQuasar(workoutDates);
  const date = ref(dateUtils.formatDate(dateUtils.dateFromSeconds()));
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
      // if (needSync()) {
      date.value = dateUtils.formatDate(
        programUtils.getWorkoutDate(props.instance, newCoords)
      );
      //    }
    },
    { deep: true }
  );
  onMounted(() => {
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
