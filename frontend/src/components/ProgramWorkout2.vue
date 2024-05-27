<script setup>
  import { inject, watch } from 'vue';
  import * as styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import ListActions from './ListActions.vue';
  import { QCheckbox, QInput } from 'quasar';
  import * as programUtils from '../modules/programUtils';

  const state = inject('state');
  const props = defineProps({ workout: Object });
  const emit = defineEmits(['update']);

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

  <div v-else>
    <ListActions @update="update" />
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
</template>
