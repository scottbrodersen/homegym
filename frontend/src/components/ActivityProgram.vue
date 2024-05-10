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
  const programIsValid = inject('programIsValid');
  const maxField = inject('maxField');
  const requiredField = inject('requiredField');

  const program = ref({});
  const changed = ref(false);
  const valid = ref(true);

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';

  let blocks; //
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
        for (let k = 0; k < programProps.cycleSpan; k++) {
          program.value.blocks[i].microCycles[j].workouts.push({
            title: `Day ${k + 1}`,
            segments: [{ exerciseTypeID: '', prescription: '' }],
          });
        }
      }
    }
  };

  watch(
    () => {
      return state.value;
    },
    (newState) => {
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
    emit('done', program.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  // watch for changes and validate
  watch(
    () => {
      return program.value;
    },
    (newVal) => {
      changed.value = baseline != JSON.stringify(newVal);
      valid.value = programIsValid(newVal);
    },
    { deep: true }
  );

  const updateButtonText = computed(() => {
    return !!program.value.id ? 'Update' : 'Add';
  });

  const doneButtonText = computed(() => {
    return changed.value ? 'Cancel' : 'Done';
  });
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
