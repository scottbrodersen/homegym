<script setup>
  import ProgramDash from './ProgramDash.vue';
  import EventsGrid from './EventsGrid.vue';
  import { activityStore, programInstanceStore } from '../modules/state';
  import { provide, ref } from 'vue';
  import * as styles from './../style.module.css';

  const props = defineProps({ eventID: String });
  const activityIDs = ref([]);
  const showAddWorkout = ref(true);
  const focusedEvent = ref(props.eventID ? props.eventID : '');
  const selectedEvent = ref(props.eventID ? props.eventID : '');

  const setFocusedEvent = (eventID) => {
    focusedEvent.value = eventID;
  };

  const setSelectedEvent = (eventID) => {
    selectedEvent.value = eventID;
    const url = new URL(document.URL);
    const newURL = url.origin + url.pathname + '?event=' + eventID;
    history.replaceState(history.state, '', newURL);
  };

  // enable child components to get and set the focused event
  provide('focusedEvent', { focusedEvent, setFocusedEvent });

  // enable child components to get and set the selected event
  provide('selectedEvent', { selectedEvent, setSelectedEvent });

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
    <EventsGrid :eventID="selectedEvent" />
  </div>
</template>
