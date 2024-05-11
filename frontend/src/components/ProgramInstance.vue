<script setup>
  import { computed, inject, onBeforeMount, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import { programInstanceStore, programsStore } from './../modules/state';
  import { updateProgramInstance } from './../modules/utils';
  import { QBtn, QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import {
    authPrompt,
    deepToRaw,
    ErrNotLoggedIn,
    states,
    toast,
  } from '../modules/utils';
  import * as programUtils from '../modules/programUtils';

  const props = defineProps({ instanceID: String });
  const emit = defineEmits(['done']);

  const instance = ref();
  let baseline = ''; // use to detect diff
  const changed = ref();
  const valid = ref(true);
  const programTitle = ref();

  const state = inject('state');
  const activityID = inject('activity').value;

  const init = (instanceID) => {
    instance.value = deepToRaw(programInstanceStore.get(instanceID));
    baseline = JSON.stringify(instance.value);

    programTitle.value = instance.value.programID
      ? programsStore.get(activityID, instance.value.programID).title
      : '';
  };

  // Re-initialize when a different instance is selected
  watch(
    () => props.instanceID,
    (newID) => {
      init(newID);
    }
  );

  // watch for changes and validate
  watch(
    () => {
      return instance.value;
    },
    (newVal) => {
      if (state.value != states.READ_ONLY) {
        changed.value = baseline != JSON.stringify(newVal);
        valid.value = programUtils.programValidator(newVal);
      }
    },
    { deep: true }
  );

  onBeforeMount(() => {
    init(props.instanceID);
  });

  const saveInstance = async () => {
    try {
      const id = await updateProgramInstance(instance.value);
      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        authPrompt(saveInstance);
      } else {
        toast('Error', 'negative');
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  const cancel = () => {
    emit('done', instance.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  const doneButtonText = computed(() => {
    return changed.value ? 'Cancel' : 'Done';
  });
</script>
<template>
  <div v-if="instance">
    <div>Start Date: {{ instance.startDate }}</div>
    <div>
      Base Program:
      {{ programTitle }}
    </div>
    <div>Events:</div>
    <div v-for="(eventID, dayIndex) of instance.events" :key="dayIndex">
      {{ dayIndex }}: {{ eventID }}
    </div>
    <div v-show="state != states.READ_ONLY">
      <q-input
        v-model="instance.title"
        label="Name"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
    </div>
    <ProgramBlock
      v-for="(block, index) of instance.blocks"
      :key="index"
      :block="block"
      @update="(value) => updateBlocks(value, index)"
    />
    <div
      v-show="state != states.READ_ONLY && instance.id"
      :class="[styles.buttonArray]"
    >
      <q-btn
        :label="doneButtonText"
        color="accent"
        text-color="dark"
        @click="cancel"
      />
      <q-btn
        label="Save"
        color="accent"
        text-color="dark"
        @click="saveInstance"
        :disable="!changed || !valid"
      />
    </div>
  </div>
</template>
