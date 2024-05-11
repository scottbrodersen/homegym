<script setup>
  import ProgramDash from './ProgramDash.vue';
  import EventsGrid from './EventsGrid.vue';
  import { activityStore, programInstanceStore } from '../modules/state';
  import { provide, ref } from 'vue';
  import styles from './../style.module.css';

  const activityIDs = ref([]);
  const showAddWorkout = ref(true);
  const focusedEvent = ref('');

  const setFocusedEvent = (eventID) => {
    focusedEvent.value = eventID;
  };

  provide('focusedEvent', { focusedEvent, setFocusedEvent });

  activityStore.getAll().forEach((activity) => {
    activityIDs.value.push(activity.id);

    if (programInstanceStore.getActive(activity.id)) {
      showAddWorkout.value = false;
    }
  });
</script>

<template>
  <div>
    <div v-for="id in activityIDs" :key="id">
      <ProgramDash
        id="program-dash"
        :activityID="id"
        :class="[styles.pgmDash]"
      />
    </div>
    <div v-if="showAddWorkout" :class="[styles.blockPadSm]">
      <q-btn
        round
        size="0.65em"
        color="primary"
        icon="add"
        :to="{ name: 'event' }"
        id="addevent"
      />
    </div>
    <EventsGrid />
  </div>
</template>
