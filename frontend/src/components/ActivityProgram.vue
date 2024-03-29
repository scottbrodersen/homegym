<script setup>
  import { computed, inject, provide, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import {
    authPrompt,
    ErrNotLoggedIn,
    newProgramModal,
    OrderedList,
    states,
    toast,
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

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';

  let blocks; // = new OrderedList(program.value.blocks);
  let baseline = ''; // use to detect diff

  // Stores the stringified program as a baseline for detecting change
  // Clones the program so it can be edited without immediately changing the store
  const init = () => {
    if (!props.programID) {
      baseline = '';
      program.value = {};
    } else {
      baseline = JSON.stringify(
        programsStore.get(activity.value.id, props.programID)
      );
      program.value = JSON.parse(baseline);
      if (!program.value.blocks) {
        program.value.blocks = [{}];
      }
      blocks = new OrderedList(program.value.blocks);
    }
  };

  // Re-initialize when a different program is selected
  watch(
    () => {
      return props.programID;
    },
    (newID) => {
      init();
    }
  );

  // callback for new program modal
  const initProgram = (programProps) => {
    //activity.value = activityStore.get(programProps.activityID);
    if (!programProps) {
      state.value = states.READ_ONLY;
      return;
    }
    program.value = {
      id: null,
      title: programProps.title,
      activityID: activity.value.id,
      blocks: new Array(),
    };

    for (let i = 0; i < programProps.numBlocks; i++) {
      program.value.blocks.push({
        title: `${defaultBlockTitle} ${i + 1}`,
        intensity: null,
        microCycles: new Array(),
      });
      for (let j = 0; j < programProps.numCycles; j++) {
        program.value.blocks[i].microCycles.push({
          title: `${defaultMicroCycleTitle} ${j + 1}`,
          span: programProps.cycleSpan,
          intensity: null,
          workouts: new Array(),
        });
      }
    }
  };

  watch(
    () => {
      return state.value;
    },
    (newState, oldState) => {
      if (newState == states.NEW) {
        newProgramModal(activity.value.id, initProgram);
      }
    }
  );

  init();

  const saveProgram = async () => {
    try {
      const id = await updateProgram(program.value);
      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        authPrompt(saveProgram);
      } else {
        toast('Error', 'negative');
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  const cancel = () => {
    // init(program.value.id);
    emit('done', program.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  const programIsValid = () => {
    let noProps = true;
    for (const prop in program.value) {
      if (Object.hasOwn(program.value, prop)) {
        noProps = false;
        break;
      }
    }
    if (noProps) {
      return false;
    }

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
          if (!workout.segments) {
            workout['segments'] = [];
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
  <div>
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
    <div
      v-show="state != states.READ_ONLY && program.title"
      :class="[styles.buttonArray]"
    >
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
  </div>
</template>
