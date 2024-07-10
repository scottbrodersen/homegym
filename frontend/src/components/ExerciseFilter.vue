<script setup>
  import { ref, watch } from 'vue';
  import { exerciseTypeStore } from '../modules/state.js';
  import { QOptionGroup, QSelect } from 'quasar';
  import * as styles from '../style.module.css';

  const props = defineProps({});
  const emit = defineEmits(['ids']);

  let exerciseOptions = exerciseTypeStore.getAll(); // reactive

  // names
  const selected = ref([]); // model

  let eTypeIDs = [];

  const filters = {
    weight: ['weight', 'percentOfMax'],
    cardio: ['hrZone', 'distance', 'pace'],
    rpe: ['rpe'],
    bodyweight: ['bodyweight'],
  };

  const filterOptions = [
    {
      label: 'All',
      value: 'all',
    },
    { label: 'Weight-based', value: 'weight' },
    { label: 'Cardio', value: 'cardio' },
    { label: 'Bodyweight', value: 'bodyweight' },
    { label: 'RPE', value: 'rpe' },
  ];

  const selectedOption = ref('all');

  const setIDs = (modelValue) => {
    eTypeIDs = [];
    for (const name of modelValue) {
      eTypeIDs.push(getID(name));
    }
    emitIDs();
  };

  const getID = (typeName) => {
    for (const exercise of exerciseOptions) {
      if (exercise.name == typeName) {
        return exercise.id;
      }
    }
    return null;
  };

  const emitIDs = () => {
    emit('ids', eTypeIDs);
  };

  watch(selectedOption, (newValue) => {
    filterExercises(newValue);
    emitIDs();
  });

  const filterExercises = (filter) => {
    selected.value = [];
    if (filter === 'all') {
      exerciseOptions = exerciseTypeStore.getAll();
      setIDs(selected.value);
    } else {
      const ghostSelected = [];
      exerciseOptions = [];
      exerciseTypeStore.getAll().forEach((type) => {
        if (filters[filter].includes(type.intensityType)) {
          exerciseOptions.push(type);
          ghostSelected.push(type.name);
        }
      });
      setIDs(ghostSelected);
    }
  };
</script>
<template>
  <div :class="[styles.chartFilter]">
    <q-option-group
      v-model="selectedOption"
      :options="filterOptions"
      dark
      inline
      dense
      :class="[styles.filterOptions]"
    />
    <q-select
      v-model="selected"
      multiple
      :options="exerciseOptions"
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
      hint="Select a category to narrow the options. Clear all selections to chart the all available options. "
    />
  </div>
</template>
