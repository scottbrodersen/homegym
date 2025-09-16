<script setup>
  /**
   * Displays the properties of a program instance workout in read-only or edit mode.
   *
   * Props:
   *  workout is the workout object to display.
   *
   * The injected state value indicates whether to display in read-only or edit mode.
   */
  import { inject, watch } from 'vue';
  import * as styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import { QCheckbox, QInput } from 'quasar';
  import * as programUtils from '../modules/programUtils';

  const { state } = inject('state');
  const props = defineProps({ workout: Object });

  if (!props.workout.segments) {
    props.workout.segments = [{}];
  }
  let segments = new OrderedList(props.workout.segments);

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
</script>
<template>
  <div>
    <div v-show="state == states.READ_ONLY">
      <div :class="styles.pgmWorkout">
        <div>
          <span :class="[styles.hgBold]"
            >{{
              props.workout.title ? props.workout.title : '~~ needs a title ~~'
            }}:
          </span>
          {{ props.workout.description }}
        </div>
        <div v-show="props.workout.restDay">REST DAY</div>
      </div>
    </div>
    <div v-show="state == states.EDIT">
      <q-input
        v-model="props.workout.title"
        label="Workout Title"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
      <q-checkbox
        v-model="props.workout.restDay"
        label="Rest Day"
        :toggle-indeterminate="false"
        indeterminate-value="never"
        dark
      />
      <q-input
        v-model="props.workout.description"
        label="Description"
        stack-label
        dark
        :rules="[programUtils.maxFieldValidator]"
      />
    </div>
  </div>
</template>
