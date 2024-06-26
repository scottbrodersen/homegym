<script setup>
  import { activityStore } from '../modules/state';
  import {
    authPromptAsync,
    ErrNotLoggedIn,
    newProgramInstanceModal,
    states,
    toast,
    updateProgramInstance,
  } from '../modules/utils';
  import * as styles from '../style.module.css';
  import { onBeforeMount, provide, ref, Suspense } from 'vue';
  import { QSelect, QBtn } from 'quasar';
  import ActivityProgram from './ActivityProgram.vue';
  import ProgramSelect from './ProgramSelect.vue';
  import ProgramInstance from './ProgramInstance.vue';

  const props = defineProps({
    activityID: String,
    programID: String,
    instanceID: String,
  });

  const selectedProgram = ref(props.programID ? props.programID : '');
  const selectedProgramInstance = ref(props.instanceID ? props.instanceID : '');
  const activityID = ref(props.activityID);

  const state = ref(states.READ_ONLY);

  const setState = (value) => {
    state.value = value;
  };

  const editProgramTitle = ref(false);
  const editInstanceTitle = ref(false);
  const toggleProgramTitle = () => {
    editProgramTitle.value = !editProgramTitle.value;
  };
  const toggleInstanceTitle = () => {
    editInstanceTitle.value = !editInstanceTitle.value;
  };

  const newProgram = ref(false);
  const toggleNewProgram = () => {
    newProgram.value = !newProgram.value;
  };

  provide('editProgramTitle', { editProgramTitle, toggleProgramTitle });
  provide('editInstanceTitle', { editInstanceTitle, toggleInstanceTitle });
  provide('newProgram', { newProgram, toggleNewProgram });
  provide('state', { state, setState });

  const setActivitySelection = (id) => {
    activityID.value = id;
    const url = new URL(document.URL);
    const newURL = url.origin + url.pathname + '?activity=' + id;
    history.replaceState(history.state, '', newURL);
  };

  const setProgramSelection = (obj) => {
    if (obj) {
      const url = new URL(document.URL);
      let newURL = url.origin + url.pathname + '?activity=' + activityID.value;
      if (obj.programID) {
        selectedProgramInstance.value = '';
        selectedProgram.value = obj.programID;
        history.replaceState(
          history.state,
          '',
          newURL + '&program=' + obj.programID
        );
      } else {
        selectedProgram.value = '';
        selectedProgramInstance.value = obj.programInstanceID;
        history.replaceState(
          history.state,
          '',
          newURL + '&instance=' + obj.programInstanceID
        );
      }
    }
  };

  const startProgram = () => {
    newProgramInstanceModal(
      activityID.value,
      selectedProgram.value,
      saveProgramInstance
    );
  };

  const saveProgramInstance = async (programInstance) => {
    try {
      const id = await updateProgramInstance(programInstance);
      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        await authPromptAsync();
        saveProgramInstance();
      } else {
        toast('Error', 'negative');
      }
    } finally {
      setState(states.READ_ONLY);
    }
  };

  onBeforeMount(() => {
    setState(states.READ_ONLY);
  });
</script>
<template>
  <div>
    <div
      id="pgm-select-wrap"
      :class="[styles.hgCentered, styles.pgmSelectWrap]"
    >
      <div :class="[styles.pgmSelect]">
        <q-select
          id="activity"
          label="Activity"
          stack-label
          :model-value="activityID"
          :options="activityStore.getAll()"
          option-label="name"
          option-value="id"
          dark
          :class="[styles.selActivity]"
          emit-value
          map-options
          compact
          @Update:model-value="setActivitySelection"
        />
        <div>
          <q-btn
            id="new"
            @click="toggleNewProgram()"
            icon="add"
            round
            dark
            color="primary"
            :disable="newProgram || !activityID"
          />
        </div>
      </div>
      <div id="pgm-select" :class="[styles.pgmSelect]" v-show="activityID">
        <Suspense>
          <ProgramSelect
            :activityID="activityID ? activityID : ''"
            :programID="
              props.programID
                ? props.programID
                : props.instanceID
                ? props.instanceID
                : ''
            "
            @selected="setProgramSelection"
          />
        </Suspense>
        <div>
          <q-btn
            id="edit"
            @click="
              selectedProgram ? toggleProgramTitle() : toggleInstanceTitle()
            "
            icon="edit"
            round
            dark
            color="primary"
            :disabled="
              (!selectedProgram && !selectedProgramInstance) ||
              editInstanceTitle ||
              editProgramTitle
            "
          />
        </div>
      </div>
    </div>
    <div>
      <q-btn
        id="start"
        v-show="selectedProgram"
        @click="startProgram"
        label="Start"
        square
        flat
        dark
      />
    </div>
    <div>
      <ActivityProgram
        :programID="selectedProgram"
        :activityID="activityID"
        @done="setProgramSelection"
      />
      <ProgramInstance
        v-if="selectedProgramInstance && !newProgram"
        :instanceID="selectedProgramInstance"
        :activityID="activityID"
        @done="setProgramSelection"
      />
    </div>
  </div>
</template>
