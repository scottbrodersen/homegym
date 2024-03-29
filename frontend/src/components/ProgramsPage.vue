<script setup>
  import { activityStore } from '../modules/state';
  import {
    authPrompt,
    ErrNotLoggedIn,
    newProgramInstanceModal,
    states,
    toast,
    updateProgramInstance,
  } from '../modules/utils';
  import styles from '../style.module.css';
  import { computed, onBeforeMount, provide, ref, Suspense } from 'vue';
  import { QSelect, QBtn } from 'quasar';
  import ActivityProgram from './ActivityProgram.vue';
  import ProgramSelect from './ProgramSelect.vue';
  import ProgramInstance from './ProgramInstance.vue';

  // props used only in unit tests
  const props = defineProps({ activity: Object, programID: String });

  const selectedProgram = ref(props.programID ? props.programID : '');
  const selectedProgramInstance = ref('');
  const activity = ref(props.activity);

  const state = ref(states.READ_ONLY);

  provide('activity', activity);
  provide('state', state);

  const setState = (value) => {
    state.value = value;
  };

  const setSelection = (obj) => {
    if (obj.programID) {
      selectedProgramInstance.value = '';
      selectedProgram.value = obj.programID;
    } else {
      selectedProgram.value = '';
      selectedProgramInstance.value = obj.programInstanceID;
    }
    setState(states.READ_ONLY);
  };

  const disableEdit = computed(() => {
    return state != states.READ_ONLY && !selectedProgram.value;
  });

  const startProgram = () => {
    newProgramInstanceModal(
      activity.value.id,
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
        authPrompt(saveProgramInstance);
      } else {
        toast('Error', 'negative');
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  onBeforeMount(() => {
    setState(states.READ_ONLY);
  });
</script>
<template>
  <div>
    <div :class="[styles.pgmSelect]">
      <q-select
        id="activity"
        label="Activity"
        stack-label
        v-model="activity"
        :options="activityStore.getAll()"
        option-label="name"
        dark
        :class="[styles.selActivity]"
      />
      <div>
        <q-btn
          id="new"
          @click="setState(states.NEW)"
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
      <Suspense>
        <ProgramSelect
          :activityID="activity ? activity.id : ''"
          @selected="setSelection"
        />
      </Suspense>
      <div>
        <q-btn
          id="edit"
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
      <q-btn
        id="start"
        v-show="selectedProgram"
        @click="startProgram"
        icon="play_arrow"
        round
        dark
        color="primary"
      />
    </div>
    <div>
      <ActivityProgram
        v-show="selectedProgram"
        :programID="selectedProgram"
        @done="setSelection"
      />
      <ProgramInstance
        v-show="selectedProgramInstance"
        :instanceID="selectedProgramInstance"
      />
    </div>
  </div>
</template>
