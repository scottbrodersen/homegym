<script setup>
  import { computed, inject, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import {
    authPrompt,
    ErrNotLoggedIn,
    OrderedList,
    states,
    updateProgram,
  } from '../modules/utils.js';
  import { QBtn, QInput } from 'quasar';
  import styles from '../style.module.css';
  import { programsStore } from '../modules/state';

  const props = defineProps({ programID: String });
  const emit = defineEmits(['done']);

  const state = inject('state');
  const activity = inject('activity');
  const program = ref({});
  const changed = ref(false);

  let blocks = new OrderedList(program.value.blocks);
  let baseline;

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

  watch(
    () => {
      return program;
    },
    (newval, oldval) => {
      changed.value = baseline != JSON.stringify(newval.value);
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
  <div v-if="state != states.READ_ONLY">
    <q-input v-model="program.title" label="Program Name" stack-label dark />
  </div>
  <ProgramBlock
    v-for="(block, index) of program.blocks"
    :key="index"
    :block="block"
    @update="(value) => updateBlocks(value, index)"
  />
  <div :class="[styles.buttonArray]" v-show="state != states.READ_ONLY">
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
      :disable="!changed"
    />
  </div>
</template>
