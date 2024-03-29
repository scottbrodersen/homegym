<script setup>
  import { inject, watch } from 'vue';
  import ProgramWorkout from './ProgramWorkout.vue';
  import styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import { QExpansionItem, QInput } from 'quasar';
  import ListActions from './ListActions.vue';

  const state = inject('state');
  const requiredField = inject('requiredField');
  const maxField = inject('maxField');
  const props = defineProps({ microcycle: Object });
  const emit = defineEmits(['update']);

  let workouts = new OrderedList(props.microcycle.workouts);

  if (!props.microcycle.workouts) {
    props.microcycle.workouts = [{}]
  }

  watch(
    () => {
      return props.microcycle.workouts;
    },
    () => {
      if (!props.microcycle.workouts) {
        props.microcycle.workouts = [{}]
      }
      workouts = new OrderedList(props.microcycle.workouts);
    }
  );
  const update = (action) => {
    emit('update', action);
  };

  const updateWorkouts = (action, index) => {
    workouts.update(action, index);
  };
</script>
<template>
  <div :class="[styles.pgmMicrocycle]">
    <div v-if="state == states.READ_ONLY">
      <div :class="[styles.horiz]">
        <div :class="[styles.hgBold, styles.sibSpxSmall]">
          {{ props.microcycle.title ? props.microcycle.title : '<no microcycle title>' }}
        </div>
        <div :class="[, styles.sibSpxSmall]">
          ({{ props.microcycle.span }} days):
        </div>

      </div>
        <div :class="[styles.sibSpxSmall]">
          {{ props.microcycle.description }}
        </div>
              <div :class="[styles.pgmMcWorkouts]">
        <ProgramWorkout
          v-for="(workout, ix) of props.microcycle.workouts"
          :key="ix"
          :workout="workout"
        />
      </div>
    </div>
    <div v-else>
      <ListActions @update="update" />
      <q-input
        v-model="props.microcycle.title"
        label="Microcycle Title"
        stack-label
        dark
        :rules="[requiredField, maxField]"
     />
      <q-input v-model="props.microcycle.span" label="Days" stack-label dark :rules="[requiredField, maxField]"
/>
      <q-input
        v-model="props.microcycle.description"
        label="Description"
        stack-label
        dark
        :rules="[maxField]"
      />
      <div :class="[styles.pgmChild]">
        <q-expansion-item
          label="Workouts"
          default-opened
          expand-separator
          expand-icon="arrow_right"
          expanded-icon="arrow_drop_down"
          :expand-icon-class="styles.pgmChildExpander"
          switch-toggle-side
          dense
        >
          <div :class="[styles.pgmMcWorkoutsEdit]">
            <ProgramWorkout
              v-for="(workout, ix) of props.microcycle.workouts"
              :key="ix"
              :workout="workout"
              @update="
                (value) => {
                  updateWorkouts(value, ix);
                }
              "
              :class="[styles.pgmChild]"
            />
          </div>
        </q-expansion-item>
      </div>
    </div>
  </div>
</template>
