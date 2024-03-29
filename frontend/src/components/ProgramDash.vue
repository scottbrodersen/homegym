<script setup>
  import { programInstanceStore } from '../modules/state';
  import { useProgramInstanceStatus } from '../composables/programInstanceStatus';
  import styles from '../style.module.css';
  import { ref } from 'vue';

  const props = defineProps({ activityID: String });
  const activeInstance = ref(
    props.activityID ? programInstanceStore.getActive(props.activityID) : null
  );

  const {
    percentComplete,
    adherence,
    blockIndex,
    microCycleIndex,
    workoutIndex,
  } = activeInstance.value
    ? useProgramInstanceStatus(activeInstance.value.id)
    : {
        percentComplete: null,
        adherence: null,
        blockIndex: null,
        microCycleIndex: null,
        workoutIndex: null,
      };
</script>
<template>
  <div v-if="activeInstance">
    <div>{{ activeInstance.title }}</div>
    <div :class="[styles.horiz]">
      <div>Progress: {{ percentComplete }}%</div>
      <div>Adherence: {{ adherence }}%</div>
    </div>
  </div>
</template>
