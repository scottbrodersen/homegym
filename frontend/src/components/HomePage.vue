<script setup>
  /**
   * The home page of the app.
   * Shows the dashboard for the current program instance and a table of workout events that have been performed.
   * The event that is viewed in the dashboard is highlighted in the table.
   */
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
    const query = eventID ? '?event=' + eventID : '';
    const newURL = url.origin + url.pathname + query;
    history.replaceState(history.state, '', newURL);
  };

  // enable child components to get and set the focused event
  provide('focusedEvent', { focusedEvent, setFocusedEvent });

  // enable child components to get and set the selected event
  provide('selectedEvent', { selectedEvent, setSelectedEvent });

  activityStore.getAll().forEach((activity) => {
    activityIDs.value.push(activity.id);

    if (programInstanceStore.getCurrent(activity.id)) {
      showAddWorkout.value = false;
    }
  });
</script>

<template>
  <div :class="[styles.vert]">
    <div v-for="id in activityIDs" :key="id">
      <div v-if="id">
        <ProgramDash
          id="program-dash"
          :activityID="id"
          :class="[styles.pgmDash]"
        />
      </div>
    </div>

    <EventsGrid :eventID="selectedEvent" />
    <div :class="[styles.blockPadSm]">
      <q-btn
        size="0.65em"
        label="Add Event"
        :to="{ name: 'event' }"
        id="addevent"
        dark
        unelevated
      />
    </div>
  </div>
</template>
