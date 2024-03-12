<script setup>
  import { computed, inject, provide, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import {
    authPrompt,
    ErrNotLoggedIn,
    OrderedList,
    states,
    updateProgram,
  } from '../modules/utils';
  import { QBtn, QInput } from 'quasar';
  import styles from '../style.module.css';
  import { programsStore } from '../modules/state';

  const props = defineProps({ programID: String });
  const emit = defineEmits(['done']);

  const state = inject('state');
  const activity = inject('activity');
  const program = ref({});
  const changed = ref(false);
  const valid = ref(true);

  let blocks = new OrderedList(program.value.blocks);
  let baseline; // use to detect diff

  // clone the program in case of editing
  const cloneProgram = () => {
    baseline = JSON.stringify(
      programsStore.get(activity.value.id, props.programID)
    );
    program.value = JSON.parse(baseline);
    if (!program.value.blocks) {
      program.value.blocks = [{}];
    }
    blocks = new OrderedList(program.value.blocks);
  };

  watch(
    () => {
      return props.programID;
    },
    (newID, oldID) => {
      cloneProgram(newID);
    }
  );

  cloneProgram(props.programID);

  const saveProgram = async () => {
    try {
      const id = await updateProgram(program.value);
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(saveProgram);
      } else {
        throw e;
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  const cancel = () => {
    cloneProgram(program.value.id);
    emit('done', program.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  const programIsValid = () => {
    if (
      requiredField(program.value.title) !== true &&
      maxField(program.value.title) !== true
    ) {
      return false;
    }
    program.value.blocks.forEach((block) => {
      if (
        requiredField(block.title) !== true &&
        maxField(block.title) !== true
      ) {
        return false;
      }
      block.microCycles.forEach((cycle) => {
        if (
          requiredField(cycle.title) !== true &&
          maxField(cycle.title) !== true
        ) {
          return false;
        }
        cycle.workouts.forEach((workout) => {
          if (
            requiredField(workout.title) !== true &&
            maxField(workout.title) !== true
          ) {
            return false;
          }
          workout.segments.forEach((segment) => {
            if (requiredField(segment.exerciseTypeID) !== true) {
              return false;
            }
            if (
              requiredField(segment.prescription) !== true &&
              maxField(segment.prescription !== true)
            ) {
              return false;
            }
          });
        });
      });
    });
    return true;
  };

  // watch for changes and validate
  watch(
    () => {
      return program;
    },
    (newVal, oldVal) => {
      changed.value = baseline != JSON.stringify(newVal.value);
      valid.value = programIsValid();
    },
    { deep: true }
  );

  const updateButtonText = computed(() => {
    return !!program.value.id ? 'Update' : 'Add';
  });

  const doneButtonText = computed(() => {
    return changed.value ? 'Cancel' : 'Done';
  });

  const requiredField = (val) => {
    const result = (val && val.length > 0) || 'Required value.';
    return result;
  };

  const maxField = (val) => {
    const result = (val ? val.length < 256 : true) || 'Max 255 characters.';
    return result;
  };

  provide('requiredField', requiredField);
  provide('maxField', maxField);
</script>
<template>
  <div v-show="state != states.READ_ONLY">
    <q-input
      v-model="program.title"
      label="Program Name"
      stack-label
      dark
      :rules="[requiredField, maxField]"
    />
  </div>
  <ProgramBlock
    v-for="(block, index) of program.blocks"
    :key="index"
    :block="block"
    @update="(value) => updateBlocks(value, index)"
  />
  <div v-show="state != states.READ_ONLY" :class="[styles.buttonArray]">
    <q-btn
      :label="doneButtonText"
      color="accent"
      text-color="dark"
      @click="cancel"
    />
    <q-btn
      :label="updateButtonText"
      color="accent"
      text-color="dark"
      @click="saveProgram"
      :disable="!changed || !valid"
    />
  </div>
</template>
