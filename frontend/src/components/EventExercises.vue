<script async setup>
  // Fetches event exerise instances and displays them.
  // To use, wrap in a Suspense component.
  import {
    authPrompt,
    fetchEventExercises,
    fetchActivityExercises,
    ErrNotLoggedIn,
  } from '../modules/utils';
  import {
    activityStore,
    eventStore,
    eventMetricsStore,
  } from '../modules/state';
  import ExerciseInstance from './ExerciseInstance.vue';
  import styles from '../style.module.css';

  const props = defineProps({ eventId: String });

  const getExercises = async (eventID) => {
    try {
      await fetchEventExercises(
        eventStore.getByID(eventID).date,
        eventStore.getByID(eventID).id
      );

      // calculate metrics
      calculateVolume(eventStore.getByID(eventID).exInstances);
      calculateLoad(eventStore.getByID(eventID).exInstances);
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

  // volume is the total number of reps
  const calculateVolume = (exerciseInstances) => {
    let volume = 0;

    for (const index of Object.keys(exerciseInstances)) {
      let instVolume = 0;
      exerciseInstances[index].parts.forEach((part) => {
        part.volume.forEach((set) => {
          set.forEach((rep) => {
            if (rep != 0) {
              instVolume++;
            }
          });
        });
      });
      volume += instVolume;
    }

    eventMetricsStore.setMetric(props.eventId, 'volume', volume);
  };

  const calculateLoad = (exerciseInstances) => {
    let load = 0;
    for (const index of Object.keys(exerciseInstances)) {
      let instLoad = 0;
      exerciseInstances[index].parts.forEach((part) => {
        let partVolume = 0;
        part.volume.forEach((set) => {
          set.forEach((rep) => {
            if (rep != 0) {
              partVolume++;
            }
          });
        });
        instLoad += Math.floor(part.intensity * partVolume);
      });

      load += instLoad;
    }
    eventMetricsStore.setMetric(props.eventId, 'load', load);
  };

  const instanceLoad = (instance) => {
    let instLoad = 0;
    instance.parts.forEach((part) => {
      let partVolume = 0;
      part.volume.forEach((set) => {
        set.forEach((rep) => {
          if (rep != 0) {
            partVolume++;
          }
        });
      });
      instLoad += Math.floor(part.intensity * partVolume);
    });
    return instLoad;
  };

  // load state
  if (!!!eventStore.getByID(props.eventId).exInstances) {
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
