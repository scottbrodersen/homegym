<script setup>
  import { provide, ref } from 'vue';
  import { programInstanceStore } from '../modules/state';
  import * as styles from '../style.module.css';
  import WorkoutAgent from './WorkoutAgent.vue';
  import { getProgramInstanceStatus } from '../modules/programUtils';
  import { useRouter } from 'vue-router';
  import * as utils from '../modules/utils';
  const router = useRouter();

  const percentComplete = ref(null);
  const adherence = ref(null);
  const coords = ref(null);
  const dayIndex = ref(null);

  const props = defineProps({ activityID: String });
  const currentInstance = ref(
    props.activityID ? programInstanceStore.getCurrent(props.activityID) : null
  );
  const activity = ref(currentInstance ? props.activityID : null);
  provide('activity', activity);

  const state = ref(utils.states.READ_ONLY);
  provide('state', { state });

  // get program stats
  const currentInstanceStatus = currentInstance.value
    ? getProgramInstanceStatus(currentInstance.value.id)
    : null;
  console.log('current instance status: ' + currentInstanceStatus);

  if (currentInstanceStatus) {
    percentComplete.value = currentInstanceStatus.percentComplete;
    adherence.value = currentInstanceStatus.adherence;
    coords.value = currentInstanceStatus.coords;
    dayIndex.value = currentInstanceStatus.dayIndex;
  }

  const goToProgram = () => {
    if (currentInstance.value) {
      router.push({
        name: 'programs',
        query: {
          activity: props.activityID,
          instance: currentInstance.value.id,
        },
      });
    }
  };
</script>
<template>
  <div v-if="currentInstance">
    <h1>
      Current Focus:<span @click="goToProgram">
        {{ currentInstance.title }}</span
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
        :workoutCoords="coords"
        :dayIndex="dayIndex"
      />
      <div v-else>
        <div>Program is complete</div>
        <div>
          <q-btn
            label="Remove from dashboard"
            @click="
              () => {
                utils.deactivateProgramInstance(
                  props.activityID,
                  currentInstance.id
                );
              }
            "
            dark
          />
        </div>
      </div>
    </div>
  </div>
</template>
