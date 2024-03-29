<script setup>
  import {
    useDialogPluginComponent,
    QCard,
    QCardActions,
    QDialog,
    QInput,
  } from 'quasar';
  import { programsStore } from '../modules/state.js';
  import { computed, ref } from 'vue';
  import styles from '../style.module.css';
  import DatePicker from './DatePicker.vue';

  const props = defineProps({ activityID: String, programID: String });
  defineEmits([...useDialogPluginComponent.emits]);

  const program = programsStore.get(props.activityID, props.programID);
  const instanceTitle = ref('');
  const startDate = ref();

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const updateDateValue = (newDate) => {
    startDate.value = newDate;
  };

  const titleIsValid = (title) => {
    if (title.length > 2 && title.length < 257) {
      return false;
    }
    return true;
  };

  const formIsValid = computed(() => {
    return titleIsValid(instanceTitle.value) && startDate.value;
  });

  const onOKClick = () => {
    const instance = JSON.parse(JSON.stringify(program));

    instance.title = instanceTitle.value;
    instance.startDate = startDate.value;
    instance.programID = props.programID;
    delete instance.id;

    onDialogOK(instance);
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card
      class="q-dialog-plugin"
      dark
      :class="[styles.blockPadSm, styles.blockBorder]"
    >
      <div>Start the {{ program.title }} program</div>
      <div>
        <q-input
          v-model="instanceTitle"
          label="Name"
          stack-label
          dark
          :rules="[
            (val) => !titleIsValid(val) || 'Length must be between 3 and 256',
          ]"
          lazy-rules
        />
        <DatePicker
          :style="[styles.blockPadMed]"
          :hideTime="true"
          @update="updateDateValue"
        />
      </div>
      <q-card-actions align="right">
        <q-btn color="primary" icon="close" round @click="onDialogCancel" />
        <q-btn
          color="primary"
          icon="done"
          round
          @click="onOKClick"
          :disabled="formIsValid"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
