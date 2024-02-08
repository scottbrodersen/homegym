<script setup>
  import { exerciseTypeStore } from '../modules/state';
  import styles from '../style.module.css';
  import { inject } from 'vue';
  import { states } from '../modules/utils.js';
  import ExerciseSelect from './ExerciseSelect.vue';
  import { QInput } from 'quasar';
  import ListActions from './ListActions.vue';

  const state = inject('state');
  const activity = inject('activity');
  const props = defineProps({ segment: Object });
  const emit = defineEmits(['update']);

  const update = (action) => {
    emit('update', action);
  };
  const setExercise = (exerciseID) => {
    props.segment.exerciseTypeID = exerciseID;
  }
</script>
<template>
  <div v-if="state == states.READ_ONLY">
    <div :class="[styles.pgmSegment]">
      <span :class="[styles.exName]">
        {{
          props.segment.exerciseTypeID
            ? exerciseTypeStore.get(props.segment.exerciseTypeID).name
            : '<no exercise selected>'
        }}:
      </span>
      {{ props.segment.prescription }}</div>
  </div>
  <div v-else>
    <ListActions @update="update" />
    <Suspense>
      <ExerciseSelect
        :activityID="activity.id"
        :exerciseID="props.segment.exerciseTypeID"
        @selected-i-d="(id) => setExercise(id)"
      />
    </Suspense>

    <q-input
      v-model="props.segment.prescription"
      label="Prescription"
      stack-label
      dark
    />
  </div>
</template>
