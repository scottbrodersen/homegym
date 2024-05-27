<script setup>
  import { inject, watch } from 'vue';
  import * as styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import { QExpansionItem, QInput } from 'quasar';
  import ListActions from './ListActions.vue';
  import * as programUtils from '../modules/programUtils';

  const state = inject('state');
  const props = defineProps({ microcycle: Object });
  const emit = defineEmits(['update']);

  let workouts = new OrderedList(props.microcycle.workouts);

  if (!props.microcycle.workouts) {
    props.microcycle.workouts = [{}];
  }

  watch(
    () => {
      return props.microcycle.workouts;
    },
    () => {
      if (!props.microcycle.workouts) {
        props.microcycle.workouts = [{}];
      }
      workouts = new OrderedList(props.microcycle.workouts);
    }
  );

  // emits the action on the workouts ordered list
  const update = (action) => {
    emit('update', action);
  };

  // called when ProgramWorkout component emits an update
  const updateWorkouts = (action, index) => {
    workouts.update(action, index);
  };
</script>
<template>
  <div>
    <div v-if="state == states.READ_ONLY" :class="[styles.pgmMicrocycle]">
      <div :class="[styles.pgmMicrocycleTitle]">
        {{
          props.microcycle.title ? props.microcycle.title : '~~needs a title~~'
        }}
      </div>
      <div>
        {{ props.microcycle.description }}
      </div>
    </div>
    <div v-else>
      <ListActions @update="update" />
      <q-input
        v-model="props.microcycle.title"
        label="Microcycle Title"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
      <q-input
        v-model="props.microcycle.span"
        label="Days"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
      <q-input
        v-model="props.microcycle.description"
        label="Description"
        stack-label
        dark
        :rules="[programUtils.maxFieldValidator]"
      />
      <div :class="[styles.pgmChild]"></div>
    </div>
  </div>
</template>
