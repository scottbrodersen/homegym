<script setup>
  import { provide, ref } from 'vue';
  import { programInstanceStore } from '../modules/state';
  import * as styles from '../style.module.css';
  import WorkoutAgent from './WorkoutAgent.vue';
  import { getProgramInstanceStatus } from '../modules/programUtils';
  import { useRouter } from 'vue-router';
  import * as utils from '../modules/utils';
  const router = useRouter();

  const props = defineProps({ activityID: String });
  const activeInstance = props.activityID
    ? programInstanceStore.getActive(props.activityID)
    : null;
  const activity = ref(activeInstance ? props.activityID : null);
  provide('activity', activity);

  const state = ref(utils.states.READ_ONLY);
  provide('state', { state });

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
        v-if="percentComplete < 100"
        :activityID="props.activityID"
        :workoutCoords="workoutCoords"
        :dayIndex="dayIndex"
      />
      <div v-else>
        <div>Program is complete</div>
        <div>
          <q-btn
            label="Remove from dashboard"
            @click="
              () => {
                utils.deactivateProgramInstance(props.activityID);
              }
            "
            dark
          />
        </div>
      </div>
    </div>
  </div>
</template>
