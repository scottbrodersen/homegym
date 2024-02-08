<script setup>
  import {
    useDialogPluginComponent,
    QCard,
    QDialog,
    QInput,
    QSelect,
  } from 'quasar';
  import { activityStore } from '../modules/state.js';
  import { computed, ref } from 'vue';
  import styles from '../style.module.css';

  const props = defineProps({ activityID: String });
  defineEmits([...useDialogPluginComponent.emits]);

  const activity = activityStore.get(props.activityID);
  const programTitle = ref('');
  const numBlocks = ref(0);
  const numCycles = ref(0);
  const cycleSpan = ref(7);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const nameIsInValid = computed(() => {
    if (programTitle.value.length < 3 && programTitle.length < 256) {
      return true;
    }
    return false;
  });

  const formIsValid = () => {
    if (!nameIsInValid && numBlocks > 0 && numCycles > 0 && cycleSpan > 0) {
      return true;
    }
    return false;
  };

  const onOKClick = () => {
    onDialogOK({
      title: programTitle.value,
      activityID: activity.id,
      numBlocks: numBlocks.value,
      numCycles: numCycles.value,
      cycleSpan: cycleSpan.value,
    });
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card
      class="q-dialog-plugin"
      dark
      :class="[styles.blockPadSm, styles.blockBorder]"
    >
      <div>
        <q-select
          label="Activity"
          stack-label
          v-model="activity"
          :options="activityStore.getAll()"
          option-label="name"
          option-value="id"
          dark
          :class="[styles.selActivity]"
        />
      </div>
      <div :v-show="!!activity">
        <q-input v-model="programTitle" label="Program Name" stack-label dark />
        <q-input
          v-model="numBlocks"
          mask="#"
          label="Number of blocks"
          stack-label
          dark
          :rules="[(val) => (val > 0 && val < 10) || 'Must be between 1 and 9']"
        />
        <q-input
          v-model="numCycles"
          mask="#"
          label="Number of microcycles"
          stack-label
          dark
          :rules="[(val) => (val > 0 && val < 10) || 'Must be between 1 and 9']"
        />
        <q-input
          v-model="cycleSpan"
          mask="#"
          label="Days in microcycles"
          stack-label
          dark
          :rules="[(val) => (val > 0 && val < 10) || 'Must be between 1 and 9']"
        />
      </div>
      <q-card-actions align="right">
        <q-btn
          color="accent"
          text-color="dark"
          label="Cancel"
          @click="onDialogCancel"
        />
        <q-btn
          color="accent"
          text-color="dark"
          label="Save"
          @click="onOKClick"
          :disabled="formIsValid()"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
