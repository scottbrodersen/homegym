<script setup>
  import { programsStore, activityStore } from '../modules/state';
  import {
    fetchPrograms,
    ErrNotLoggedIn,
    newProgramModal,
    states,
  } from '../modules/utils';
  import styles from '../style.module.css';
  import { computed, onBeforeMount, provide, ref } from 'vue';
  import { QSelect, QBtn } from 'quasar';
  import Program from './Program.vue';

  const selectedProgram = ref('');
  const activity = ref();
  const programs = ref([]);

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';
  const state = ref(states.READ_ONLY);

  provide('activity', activity);
  provide('state', state);

  const setState = (value) => {
    state.value = value;
  };

  // populates program list for selected activity
  const populatePrograms = (activityID) => {
    programs.value = programsStore.getByActivity(activityID);
  };

  const getPrograms = async (activityID) => {
    if (activityID && programsStore.getByActivity(activityID) === undefined) {
      try {
        await fetchPrograms(activityID);
      } catch (e) {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(getPrograms, activityID);
        } else {
          console.log(e.message);
        }
      }
    }
    populatePrograms(activityID);
  };

  // callback for new program modal
  const initProgram = (programProps) => {
    activity.value = activityStore.get(programProps.activityID);

    program.value = {
      id: null,
      title: programProps.title,
      activityID: programProps.activityID,
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
    setState(states.NEW);
    selectedProgram.value = '';
  };

  const setSelection = (programID) => {
    selectedProgram.value = programID;
    setState(states.READ_ONLY);
  };

  const addNew = (programID, done) => {
    // update selector
    done(programsStore.get(activity.value.id, programID), 'add-unique');
  };

  const disableEdit = computed(() => {
    return state != states.READ_ONLY && !selectedProgram.value;
  });

  onBeforeMount(() => {
    setState(states.READ_ONLY);
  });
</script>
<template>
  <div :class="[styles.pgmSelect]">
    <q-select
      label="Activity"
      stack-label
      v-model="activity"
      :options="activityStore.getAll()"
      option-label="name"
      dark
      :class="[styles.selActivity]"
      @update:model-value="(value) => getPrograms(value.id)"
    />
    <div>
      <q-btn
        @click="newProgramModal(activity.id, initProgram)"
        icon="add"
        round
        dark
        color="primary"
        :disable="state != states.READ_ONLY || !activity"
      />
    </div>
  </div>
  <div
    :class="[styles.pgmSelect]"
    v-show="!!activity"
    v-if="state == states.READ_ONLY"
  >
    <q-select
      label="Program"
      stack-label
      v-model="selectedProgram"
      :options="programs"
      option-label="title"
      option-value="id"
      emit-value
      @new-value="addNew"
      dark
      :class="[styles.selProgram]"
    >
      <template v-slot:selected>
        <div v-if="selectedProgram">
          {{ programsStore.get(activity.id, selectedProgram).title }}
        </div>
      </template>
    </q-select>
    <div>
      <q-btn
        @click="setState(states.EDIT)"
        icon="edit"
        round
        dark
        color="primary"
        :disable="disableEdit"
      />
    </div>
  </div>

  <div>
    <Program
      v-if="selectedProgram"
      :programID="selectedProgram"
      @done="setSelection"
    />
  </div>
</template>
