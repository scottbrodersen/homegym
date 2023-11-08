<script async setup>
  // Fetches event exerise instances and displays them.
  // To use, wrap in a Suspense component.
  import {
    fetchEventExercises,
    fetchActivityExercises,
  } from '../modules/utils';
  import { activityStore, eventStore } from '../modules/state';
  import ExerciseInstance from './ExerciseInstance.vue';
  import styles from '../style.module.css';
  import { onBeforeMount } from 'vue';

  const props = defineProps({ eventId: String });

  const getExercises = async (eventID) => {
    try {
      await fetchEventExercises(
        eventStore.getByID(eventID).date,
        eventStore.getByID(eventID).id
      );
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(getExercises, eventID);
      } else {
        throw e;
      }
    }
  };

  const getExerciseTypes = async (activityID) => {
    try {
      await fetchActivityExercises(activityID);
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(getExerciseTypes, activityID);
      } else {
        throw e;
      }
    }
  };

  // load state
  if (!eventStore.getByID(props.eventId).exInstances) {
    await getExercises(props.eventId);
  }
  const activityID = eventStore.getByID(props.eventId).activityID;
  if (activityStore.get(activityID).exercises == null) {
    await getExerciseTypes(activityID);
  }
</script>

<template>
  <div :class="[styles.blockPadXSm]">
    <ExerciseInstance
      v-for="(value, key) in eventStore.getByID(eventId).exInstances"
      :exercise-instance="value"
      :activity-id="eventStore.getByID(eventId).activityID"
      :writable="false"
    />
  </div>
</template>
