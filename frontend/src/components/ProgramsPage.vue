<script setup>
  /**
   * Enables a user to select a program or program instance for an activity.
   * Displays the selected item.
   *
   * Props:
   *  activityID is the activity with which the program is associated
   *  programID (optional) is the ID of the program that is pre-selected
   *  instanceID (optional) is the ID of the program instance that is pre-selected
   *
   * Either programID or instanceID can be provided, but not both.
   *
   * Provides the following values to child components:
   *  editProgramTitle is true when the edit button is clicked for a program title.
   *  editInstanceTitle is true when the edit button is clicked for a program instance title.
   *  newProgram is true when a new program is being created
   *  state indicates whether to display the program or instance in read-only or edit mode.
   */
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
  import { inject, onBeforeMount, provide, ref, Suspense } from 'vue';
  import { QSelect, QBtn, QCheckbox } from 'quasar';
  import ActivityProgram from './ActivityProgram.vue';
  import ProgramSelect from './ProgramSelect.vue';
  import ProgramInstance from './ProgramInstance.vue';

  const docsContext = ref(inject('docsContext'));
  docsContext.value = 'programs';

  const props = defineProps({
    activityID: String,
    programID: String,
    instanceID: String,
  });

  const activities = ref([]);
  const selectedProgram = ref(props.programID ? props.programID : '');
  const selectedProgramInstance = ref(props.instanceID ? props.instanceID : '');
  const activityID = ref(props.activityID);

  const hideCompleted = ref(true);

  const state = ref(states.READ_ONLY);

  const setState = (value) => {
    state.value = value;
  };

  const editProgramTitle = ref(false);
  const editInstanceTitle = ref(false);
  const cloneProgram = ref(false);

  // listener for the edit program title button
  const toggleProgramTitle = () => {
    editProgramTitle.value = !editProgramTitle.value;
  };

  // listener for the edit instance title button
  const toggleInstanceTitle = () => {
    editInstanceTitle.value = !editInstanceTitle.value;
  };

  // listener for the new program button
  const newProgram = ref(false);
  const toggleNewProgram = () => {
    newProgram.value = !newProgram.value;
  };

  //listener for the clone program button
  const toggleCloneProgram = () => {
    cloneProgram.value = !cloneProgram.value;
  };

  // Mutates the program title
  provide('editProgramTitle', { editProgramTitle, toggleProgramTitle });
  // Mutates the program instance
  provide('editInstanceTitle', { editInstanceTitle, toggleInstanceTitle });
  // Mutates the new program state
  provide('newProgram', { newProgram, toggleNewProgram });
  // mutates the edit/read only state
  provide('state', { state, setState });
  // mutates the clone program state
  provide('cloneProgram', { cloneProgram, toggleCloneProgram });

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
          newURL + '&program=' + obj.programID,
        );
      } else {
        selectedProgram.value = '';
        selectedProgramInstance.value = obj.programInstanceID;
        history.replaceState(
          history.state,
          '',
          newURL + '&instance=' + obj.programInstanceID,
        );
      }
    }
  };

  const startProgram = () => {
    newProgramInstanceModal(
      activityID.value,
      selectedProgram.value,
      saveProgramInstance,
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

  const getActivities = () => {
    activities.value = activityStore.getAll();
    activities.value.sort((a, b) => a.name.localeCompare(b.name));
  };

  onBeforeMount(() => {
    setState(states.READ_ONLY);
    getActivities();
  });
</script>
<template>
  <div>
    <div
      id="pgm-select-wrap"
      :class="[styles.pgmCentered, styles.pgmSelectWrap]"
    >
      <div :class="[styles.selActivityWrap]">
        <q-select
          id="activity"
          label="Activity"
          stack-label
          :model-value="activityID"
          :options="activities"
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
            :hideCompleted="hideCompleted"
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
        <div>
          <q-btn
            id="clone"
            @click="toggleCloneProgram()"
            icon="content_copy"
            round
            dark
            color="primary"
            :disabled="
              !selectedProgram ||
              selectedProgramInstance ||
              editInstanceTitle ||
              editProgramTitle
            "
          />
        </div>
      </div>
      <div v-show="activityID">
        <q-checkbox v-model="hideCompleted" label="Hide completed" dark />
      </div>
    </div>
    <div :class="[styles.pgmStartWrap]">
      <q-btn
        id="start"
        v-show="selectedProgram"
        @click="startProgram"
        label="Start"
        square
        flat
        dark
        outline
      />
    </div>
    <div>
      <ActivityProgram
        v-show="selectedProgram || newProgram"
        :programID="selectedProgram"
        :activityID="activityID"
        @done="setProgramSelection"
      />
      <ProgramInstance
        v-if="selectedProgramInstance && !newProgram"
        :instanceID="selectedProgramInstance"
        :activityID="activityID"
      />
    </div>
  </div>
</template>
