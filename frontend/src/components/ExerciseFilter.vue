<script setup>
  import { ref } from 'vue';
  import { exerciseTypeStore } from '../modules/state.js';
  import { QSelect } from 'quasar';
  import * as styles from '../style.module.css';

  const props = defineProps({});
  const emit = defineEmits(['ids']);

  const exercises = exerciseTypeStore.getAll(); // reactive

  // names
  const selected = ref([]); // model

  let eTypeIDs = [];

  const setIDs = (modelValue) => {
    eTypeIDs = [];
    for (const name of modelValue) {
      eTypeIDs.push(getID(name));
    }
    emitIDs();
  };

  const getID = (typeName) => {
    for (const exercise of exercises) {
      if (exercise.name == typeName) {
        return exercise.id;
      }
    }
    return null;
  };

  const emitIDs = () => {
    emit('ids', eTypeIDs);
  };
</script>
<template>
  <q-select
    v-model="selected"
    multiple
    :options="exercises"
    option-label="name"
    option-value="name"
    emit-value
    use-chips
    label="Exercises"
    stack-label
    :class="[styles.selExercise]"
    dark
    @update:model-value="setIDs"
    @popup-hide="emitIDs"
  />
</template>
