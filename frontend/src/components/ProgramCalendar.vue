<script setup>
  import { onMounted, ref } from 'vue';
  import { QDate } from 'quasar';
  import * as dateUtils from '../modules/dateUtils';
  import * as programUtils from '../modules/programUtils';

  const props = defineProps({ dateValue: Number, instance: Object });
  const emit = defineEmits(['dayIndex']);

  const nonRestDays = programUtils.getNonRestDates(props.instance);
  const formattedNonRestDays = formatForQuasar(nonRestDays);

  const workoutDates = programUtils.getInstanceWorkoutDates(props.instance);
  const formattedWorkoutDates = formatForQuasar(workoutDates);
  const date = ref(
    dateUtils.formatDate(dateUtils.dateFromSeconds(props.dateValue))
  );

  const setDate = (date) => {
    emit(
      'dayIndex',
      workoutDates.findIndex((workoutDate) => workoutDate == date)
    );
  };

  function formatForQuasar(dateArr) {
    const formatted = new Array();
    dateArr.forEach((dateStr) => {
      formatted.push(dateStr.replaceAll('-', '/'));
    });
    return formatted;
  }

  onMounted(() => {
    setDate(date.value);
  });
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
      minimal
      bordered
      @update:model-value="(value) => setDate(value)"
    />
  </div>
</template>
