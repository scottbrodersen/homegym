<script setup>
  /**
   * Enables users to select a category of exercises and then filters the list of exercises by the selection.
   * Emits an array of the IDs of filtered exercises.
   */

  import { onMounted, ref, watch } from 'vue';
  import { exerciseTypeStore } from '../modules/state.js';
  import { QOptionGroup, QSelect } from 'quasar';
  import * as styles from '../style.module.css';

  const props = defineProps({});
  const emit = defineEmits(['ids']);

  let exerciseOptions = exerciseTypeStore.getAll(); // reactive objects

  // exercise names that appear as tags
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

  const selectedFilter = ref('');

  // Emits the id's of a list of exercise names
  const setIDs = (exerciseNames) => {
    eTypeIDs = [];
    if (exerciseNames) {
      for (const name of exerciseNames) {
        eTypeIDs.push(getID(name));
      }
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

  watch(selectedFilter, (newValue) => {
    filterExercises(newValue);
  });

  onMounted(() => {
    selectedFilter.value = 'all';
  });

  // Filters the exercises that appear in the dropdown
  // Adds the IDs of the filtered exercises to the selected ID array
  // By default all filtered exercises are charted but none appear
  const filterExercises = (filter) => {
    selected.value = [];

    // IDs that are charted but do not appear as selected in the dropdown
    const ghostSelected = [];

    if (filter === 'all') {
      exerciseOptions = exerciseTypeStore.getAll();
    } else {
      exerciseOptions = [];
      exerciseTypeStore.getAll().forEach((type) => {
        if (filters[filter].includes(type.intensityType)) {
          exerciseOptions.push(type);
        }
      });
    }

    exerciseOptions.forEach((exerciseType) => {
      selected.value.push(exerciseType.name);
    });

    setIDs(selected.value);
  };
</script>
<template>
  <div :class="[styles.chartFilter]">
    <q-option-group
      v-model="selectedFilter"
      :options="filterOptions"
      dark
      inline
      dense
      :class="[styles.filterOptions]"
    />
    <q-select
      v-model="selected"
      multiple
      clearable
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
      @clear="emitIDs"
      hint="Select a category to narrow the options. Clear all selections to chart the all available options. "
    />
  </div>
</template>
