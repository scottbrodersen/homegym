<script setup>
  /**
   * Dialog for editing the workout of a program instance.
   *
   * Props:
   *  instance is the program instance to which the workout belongs.
   *  coords is the 3x1 array of coordinates that indicates where the workout occurs within the structure of the instance.
   *    E.g. [0,1,2] denotes the workout in the 3rd day of the 2nd microcycle in the 1st block.
   *
   *  Returns the workout object.
   */
  import { provide, ref, watch } from 'vue';
  import { useDialogPluginComponent, QCard, QDialog } from 'quasar';
  import * as utils from '../modules/utils';
  import * as programUtils from '../modules/programUtils';
  import ProgramWorkout2 from './ProgramWorkout2.vue';
  import * as styles from '../style.module.css';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();
  const emit = defineEmits([...useDialogPluginComponent.emits]);

  const props = defineProps({ instance: Object, coords: Array });
  const workout = ref(
    props.instance.blocks[props.coords[0]].microCycles[props.coords[1]]
      .workouts[props.coords[2]]
  );

  const workoutIsValid = ref(true);

  const segments = new utils.OrderedList(workout.value.segments);

  const state = ref(utils.states.EDIT);

  provide('state', { state });
  provide('activity', props.instance.activityID);

  const onOKClick = () => {
    onDialogOK(workout.value);
  };

  watch(
    () => workout.value,
    (newValue) => {
      workoutIsValid.value = programUtils.workoutValidator(newValue);
    },
    { deep: true }
  );

  const updateSegments = (action, index) => {
    segments.update(action, index);
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card
      class="q-dialog-plugin"
      dark
      :class="[styles.blockPadSm, styles.blockBorder]"
    >
      <ProgramWorkout2 :workout="workout" />
      <ProgramWorkoutSegment
        v-for="(segment, ix) in workout.segments"
        :key="ix"
        :segment="segment"
        @update="(action) => updateSegments(action, ix)"
      />
      <q-card-actions align="right">
        <q-btn color="primary" icon="close" round @click="onDialogCancel" />
        <q-btn
          color="primary"
          icon="done"
          round
          @click="onOKClick"
          :disabled="!workoutIsValid"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
