<script setup>
  import ProgramDash from './ProgramDash.vue';
  import EventsGrid from './EventsGrid.vue';
  import { activityStore } from '../modules/state';
  import { ref } from 'vue';
  const activityIDs = ref([]);
  activityStore.getAll().forEach((activity) => {
    activityIDs.value.push(activity.id);
  });
</script>

<template>
  <div>
    <div v-for="id in activityIDs" :key="id">
      <ProgramDash :activityID="id" />
    </div>
    <EventsGrid />
  </div>
</template>
