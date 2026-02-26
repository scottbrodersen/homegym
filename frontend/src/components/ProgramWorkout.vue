<script setup>
  /**
   * Displays the properties of a program workout in read-only or edit mode.
   *
   * Props:
   *  workout is the workout object to display.
   *
   * The injected state value indicates whether to display in read-only or edit mode.
   */
  import { inject, ref, watch } from 'vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';
  import * as styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import { QCheckbox, QInput } from 'quasar';
  import * as programUtils from '../modules/programUtils';
  import * as utils from '../modules/utils';

  const { state } = inject('state');
  const props = defineProps({ workout: Object });
  const rawWorkout = ref(utils.deepToRaw(props.workout));
  const emits = defineEmits('update');

  // todo: do not mutate props.workout

  if (!props.workout.segments) {
    rawWorkout.value.segments = [{}];
  }
  let segments = new OrderedList(rawWorkout.value.segments);

  const updateSegments = (action, index) => {
    segments.update(action, index);
    updateWorkout();
  };

  const updateSegmentExercise = (exerciseID, index) => {
    segments.list[index].exerciseTypeID = exerciseID;
    updateWorkout();
  };

  const updateWorkout = () => {
    emits('update', rawWorkout.value);
  };
</script>
<template>
  <div>
    <div v-if="state == states.READ_ONLY">
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
        <div v-show="!props.workout.restDay">
          <div v-for="(segment, ix) of segments.list" :key="ix">
            <ProgramWorkoutSegment :segment="segment" />
          </div>
        </div>
      </div>
    </div>
    <div v-if="state == states.EDIT">
      <q-input
        v-model="rawWorkout.title"
        label="Workout Title"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
        @update:model-value="() => updateWorkout()"
      />
      <q-checkbox
        v-model="rawWorkout.restDay"
        label="Rest Day"
        :toggle-indeterminate="false"
        indeterminate-value="never"
        dark
        @update:model-value="() => updateWorkout()"
      />
      <q-input
        v-model="rawWorkout.description"
        label="Description"
        stack-label
        dark
        :rules="[programUtils.maxFieldValidator]"
        @update:model-value="() => updateWorkout()"
      />
      <div v-show="!rawWorkout.restDay" :class="[styles.pgmChild]">
        <ProgramWorkoutSegment
          v-for="(segment, ix) of segments.list"
          :key="ix"
          :segment="segment"
          @update="
            (value) => {
              updateSegments(value, ix);
            }
          "
          @setExercise="
            (value) => {
              updateSegmentExercise(value, ix);
            }
          "
        />
      </div>
    </div>
  </div>
</template>
