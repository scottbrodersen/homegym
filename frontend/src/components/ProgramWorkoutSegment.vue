<script async setup>
  import { exerciseTypeStore } from '../modules/state';
  import * as styles from '../style.module.css';
  import { inject, Suspense } from 'vue';
  import { states } from '../modules/utils.js';
  import ExerciseSelect from './ExerciseSelect.vue';
  import { QInput } from 'quasar';
  import ListActions from './ListActions.vue';
  import * as programUtils from '../modules/programUtils';

  const props = defineProps({ segment: Object });
  const emit = defineEmits(['update']);

  const { state } = inject('state');
  const activityID = inject('activity');

  const update = (action) => {
    emit('update', action);
  };
  const setExercise = (exerciseID) => {
    props.segment.exerciseTypeID = exerciseID;
  };
</script>
<template>
  <div>
    <div v-if="state == states.READ_ONLY">
      <div :class="[styles.pgmSegment]">
        <span :class="[styles.exName]">
          {{
            props.segment.exerciseTypeID
              ? exerciseTypeStore.get(props.segment.exerciseTypeID).name
              : '~~no exercise selected~~'
          }}:
        </span>
        {{ props.segment.prescription }}
      </div>
    </div>
    <div v-if="state == states.EDIT" :class="[styles.horiz]">
      <div>
        <Suspense>
          <ExerciseSelect
            :activityID="activityID"
            :exerciseID="props.segment.exerciseTypeID"
            @selected-i-d="(id) => setExercise(id)"
          />
        </Suspense>
        <q-input
          v-model="props.segment.prescription"
          :class="[styles.pgmSegmentPn]"
          label="Prescription"
          stack-label
          dark
          :rules="[
            programUtils.requiredFieldValidator,
            programUtils.maxFieldValidator,
          ]"
          type="textarea"
        />
      </div>
      <ListActions @update="update" />
    </div>
  </div>
</template>
