<script setup>
  import { ref } from 'vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';
  import * as styles from '../style.module.css';
  import * as programUtils from '../modules/programUtils';
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
    programUtils.getWorkoutStatus(
      props.eventID,
      props.workoutIndex,
      props.todayIndex,
      props.workout.restDay
    )
  );

  const icon = ref(programUtils.getStatusIcons(status.value));

  colourStyle.value = programUtils.getColorStyle(status);
</script>
<template>
  <div>
    <div :class="[styles.workoutStatusWrap]">
      <div>
        <div>
          <div :class="[styles.hgBold]">{{ when }}</div>
          <div>
            <span :class="[styles.hgBold]"
              >{{
                props.workout.title
                  ? props.workout.title
                  : '~~ needs a title ~~'
              }}:</span
            >
            {{ props.workout.description }}
          </div>
          <div v-if="props.workout.restDay">REST DAY</div>
        </div>
        <div v-if="!props.workout.restDay && props.workout.segments">
          <ProgramWorkoutSegment
            v-for="(segment, ix) of props.workout.segments"
            :key="ix"
            :segment="segment"
          />
        </div>
      </div>
      <div :class="[styles.workoutStatusIcon]">
        <q-icon :name="icon.name" :color="icon.colour" right />
      </div>
    </div>
  </div>
</template>
