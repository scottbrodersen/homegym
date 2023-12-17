<script setup>
  import { inject, watch } from 'vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';
  import styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import ListActions from './ListActions.vue';

  const state = inject('state');
  const props = defineProps({ workout: Object });
  const emit = defineEmits(['update']);

  let segments = new OrderedList(props.workout.segments);

  if (!props.workout.segments) {
    props.workout.segments = [{}];
  }

  watch(
    () => {
      return props.workout.segments;
    },
    () => {
      if (!props.workout.segments) {
        props.workout.segments = [{}];
      }
      segments = new OrderedList(props.workout.segments);
    }
  );

  const update = (action) => {
    emit('update', action);
  };

  const updateSegments = (action, index) => {
    segments.update(action, index);
  };
</script>
<template>
  <div v-if="state == states.READ_ONLY">
    <div :class="styles.pgmWorkout">
      <div>
        <span :class="[styles.hgBold]">{{ props.workout.title ? props.workout.title : '<needs a title>'}}:</span>
        {{ props.workout.intensity }}
      </div>
        <div v-for="(segment, ix) of segments.list" :key="ix">
          <ProgramWorkoutSegment :segment="segment" />
        </div>
    </div>
  </div>

  <div v-else>
    <ListActions @update="update" />
    <q-input
      v-model="props.workout.title"
      label="Workout Title"
      stack-label
      dark
    />
    <q-input
      v-model="props.workout.intensity"
      label="Intensity"
      stack-label
      dark
    />
    <div :class="[styles.pgmChild]">
      <ProgramWorkoutSegment v-for="(segment, ix) of props.workout.segments" :key="ix"
        :segment="segment"
        @update="
          (value) => {
            updateSegments(value, ix);
          }
        "
      />
      </div>
  </div>
</template>
