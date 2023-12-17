<script setup>
  import { computed, ref } from 'vue';
  import { activityStore, exerciseTypeStore } from '../modules/state.js';
  import { QSelect } from 'quasar';
  import styles from '../style.module.css';

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
    if (!!props.activityID) {
      activityStore.get(props.activityID).exercises.forEach((exerciseID) => {
        const eType = exerciseTypeStore.get(exerciseID);
        eTypeIDs.push(eType.id);
        names.push(eType.name);
      });
    }
    return names;
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
