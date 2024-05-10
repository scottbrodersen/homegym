<script setup>
  import { programInstanceStore } from '../modules/state';
  import styles from '../style.module.css';
  import { provide, ref } from 'vue';
  import WorkoutAgent from './WorkoutAgent.vue';
  import { states } from '../modules/utils';
  import { getProgramInstanceStatus } from '../modules/programUtils';

  //provide('state', states.READ_ONLY);
  const props = defineProps({ activityID: String });
  const activeInstance = props.activityID
    ? programInstanceStore.getActive(props.activityID)
    : null;

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
