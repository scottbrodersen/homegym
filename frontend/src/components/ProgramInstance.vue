<script setup>
  import { onBeforeMount, ref, watch } from 'vue';
  import ActivityProgram from './ActivityProgram.vue';
  import { programInstanceStore, programsStore } from './../modules/state';

  const props = defineProps({ instanceID: String, programID: String });

  const getInstance = (instanceID) => {
    return instanceID ? programInstanceStore.get(props.instanceID) : null;
  };

  const getProgramTitle = (activityID, programID) => {
    return programsStore.get(activityID, programID).title;
  };

  let instance;
  const programTitle = ref();

  watch(
    () => props.instanceID,
    (newID) => {
      instance = getInstance(newID);
      programTitle.value = getProgramTitle(
        instance.activityID,
        instance.programID
      );
    }
  );
  onBeforeMount(() => {
    instance = getInstance(props.instanceID);
    programTitle.value = getProgramTitle(
      instance.activityID,
      instance.programID
    );
  });
</script>
<template>
  <div v-if="instance">
    <div>Start Date: {{ instance.startDate }}</div>
    <div>
      Base Program:
      {{ programTitle }}
    </div>
    <div>Events:</div>
    <div v-for="(eventID, dayIndex) of instance.events" :key="dayIndex">
      {{ dayIndex }}: {{ eventID }}
    </div>
    <ActivityProgram :programID="instance.id" />
  </div>
</template>
