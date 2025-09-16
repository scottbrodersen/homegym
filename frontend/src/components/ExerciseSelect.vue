<script setup>
  /**
   * A drop down list of exercises for an activity.
   *
   * Props:
   *  activityID is the ID of the activity
   *  exerciseID is the ID of the pre-selected exercise. Provide no value if no exercise is pre-selected.
   *
   * Emits the ID of the selected exercise.
   */
  import { computed, onBeforeMount, ref, watch } from 'vue';
  import { activityStore, exerciseTypeStore } from '../modules/state.js';
  import { QSelect } from 'quasar';
  import * as styles from '../style.module.css';
  import {
    authPromptAsync,
    ErrNotLoggedIn,
    fetchActivityExercises,
  } from '../modules/utils';

  const props = defineProps({ activityID: String, exerciseID: String });
  const emits = defineEmits(['selectedID']);

  const exerciseName = props.exerciseID
    ? ref(exerciseTypeStore.get(props.exerciseID).name)
    : ref('');
  const eTypeIDs = [];

  const setExercise = (typeName) => {
    for (const id of eTypeIDs) {
      if (exerciseTypeStore.get(id).name == typeName) {
        emits('selectedID', id);
        exerciseName.value = typeName;
        break;
      }
    }
  };

  const exerciseNames = computed(() => {
    const names = [];
    if (props.activityID && activityStore.get(props.activityID).exercises) {
      activityStore.get(props.activityID).exercises.forEach((exerciseID) => {
        const eType = exerciseTypeStore.get(exerciseID);
        eTypeIDs.push(eType.id);
        names.push(eType.name);
      });
    }
    return names;
  });

  const getActivityExercises = async (activityID) => {
    // fetch activity exercises types if needed
    if (activityStore.get(activityID).exercises == null) {
      try {
        await fetchActivityExercises(activityID);
      } catch (e) {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          await authPromptAsync();
          getActivityExercises(activityID);
        } else {
          console.log(e);
        }
      }
    }
  };

  watch(
    () => {
      return props.activityID;
    },
    async (newID) => {
      await getActivityExercises(newID);
    }
  );
  onBeforeMount(async () => {
    await getActivityExercises(props.activityID);
  });
</script>
<template>
  <q-select
    :model-value="exerciseName"
    :options="exerciseNames"
    label="Exercise"
    stack-label
    :class="[styles.selExercise]"
    dark
    @update:model-value="setExercise"
  />
</template>
