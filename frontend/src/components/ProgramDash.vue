<script setup>
  import { programInstanceStore } from '../modules/state';
  import * as styles from '../style.module.css';
  import WorkoutAgent from './WorkoutAgent.vue';
  import { getProgramInstanceStatus } from '../modules/programUtils';

  const props = defineProps({ activityID: String });
  const activeInstance = props.activityID
    ? programInstanceStore.getActive(props.activityID)
    : null;

  // get program stats
  const [percentComplete, adherence, workoutCoords, dayIndex] = activeInstance
    ? getProgramInstanceStatus(activeInstance.id)
    : [null, null, null, null];
</script>
<template>
  <div v-if="activeInstance">
    <h1>Current Focus: {{ activeInstance.title }}</h1>
    <div :class="[styles.pgmStatus]">
      <div>Progress: {{ percentComplete }}%</div>
      <div>Adherence: {{ adherence }}%</div>
    </div>
    <div>
      <WorkoutAgent
        :activityID="props.activityID"
        :workoutCoords="workoutCoords"
        :dayIndex="dayIndex"
      />
    </div>
  </div>
</template>
