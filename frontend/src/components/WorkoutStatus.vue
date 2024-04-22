<script setup>
  import { computed, ref } from 'vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';
  import styles from '../style.module.css';
  import * as utils from '../modules/programUtils';
  import { QIcon } from 'quasar';

  const props = defineProps({
    eventID: String,
    todayIndex: Number,
    workoutIndex: Number,
    workout: Object,
  });
  const fromNow = 0 - (props.todayIndex - props.workoutIndex);
  const when = ref('');
  const colourStyle = ref('');

  if (fromNow == 0) {
    when.value = 'Today';
  } else if (fromNow == -1) {
    when.value = 'Yesterday';
  } else if (fromNow < -1) {
    when.value = `${-1 * fromNow} days ago`;
  } else if (fromNow == 1) {
    when.value = 'Tomorrow';
  } else if (fromNow > 1) {
    when.value = `${fromNow} days from now`;
  }

  const status = ref(
    utils.getWorkoutStatus(
      props.eventID,
      props.workoutIndex,
      props.todayIndex,
      props.workout.restDay
    )
  );

  const icon = ref(utils.getStatusIcons(status.value));

  colourStyle.value = utils.getColorStyle(status);
</script>
<template>
  <div>
    <div :class="[styles.horiz]">
      <div>{{ when }}</div>
      <div>
        <q-icon :name="icon.name" :color="icon.colour" right size="25px" />
      </div>
    </div>
    <div>
      <span :class="[styles.hgBold]"
        >{{
          props.workout.title ? props.workout.title : '~~ needs a title ~~'
        }}:</span
      >
      {{ props.workout.description }}
      <div v-if="props.workout.restDay">REST DAY</div>
    </div>
    <div v-if="props.workout.segments">
      <ProgramWorkoutSegment
        v-for="(segment, ix) of props.workout.segments"
        :key="ix"
        :segment="segment"
      />
    </div>
  </div>
</template>
