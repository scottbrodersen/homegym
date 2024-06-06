<script setup>
  import { programInstanceStore } from '../modules/state';
  import * as styles from '../style.module.css';
  import WorkoutAgent from './WorkoutAgent.vue';
  import { getProgramInstanceStatus } from '../modules/programUtils';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  const props = defineProps({ activityID: String });
  const activeInstance = props.activityID
    ? programInstanceStore.getActive(props.activityID)
    : null;

  // get program stats
  const [percentComplete, adherence, workoutCoords, dayIndex] = activeInstance
    ? getProgramInstanceStatus(activeInstance.id)
    : [null, null, null, null];

  const goToProgram = () => {
    if (activeInstance) {
      router.push({
        name: 'programs',
        query: {
          activity: props.activityID,
          instance: activeInstance.id,
        },
      });
    }
  };
</script>
<template>
  <div v-if="activeInstance">
    <h1>
      Current Focus:<span @click="goToProgram">
        {{ activeInstance.title }}</span
      >
    </h1>
    <div :class="[styles.pgmStatus]">
      <div>Progress: {{ percentComplete }}%</div>
      <div>Adherence: {{ adherence }}%</div>
    </div>
    <div>
      <WorkoutAgent
        v-if="dayIndex"
        :activityID="props.activityID"
        :workoutCoords="workoutCoords"
        :dayIndex="dayIndex"
      />
      <div v-else>
        Program is complete
        <q-btn
          round
          color="primary"
          icon="visibility"
          :to="{
            name: 'programs',
            query: {
              activity: props.activityID,
              instance: activeInstance.id,
            },
          }"
        />
      </div>
    </div>
  </div>
</template>
