<script setup>
  import { watch } from 'vue';
  import ActivityProgram from './ActivityProgram.vue';
  import { programInstanceStore, programsStore } from './../modules/state';

  const props = defineProps({ instanceID: String, programID: String });
  let instance = props.instanceID ? programInstanceStore.get(newID) : null;
  watch(
    () => props.instanceID,
    (newID) => {
      instance = newID ? programInstanceStore.get(newID) : null;
    }
  );
</script>
<template>
  <div v-if="instance">
    Start Date: {{ instance.startDate }} Base Program:
    {{ programsStore.get(instance.activityID, instance.programID).title }}
    <ActivityProgram :programID="instance.programID" />
  </div>
</template>
