<script setup>
  import {
    useDialogPluginComponent,
    QInput,
    QDialog,
    QCard,
    QCardActions,
    QBtn,
  } from 'quasar';
  import { ref } from 'vue';
  import * as styles from '../style.module.css';

  const props = defineProps({ fieldValue: String, fieldLabel: String });

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const newValue = ref(props.fieldValue);

  const onOKClick = () => {
    onDialogOK(newValue.value);
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.bgBlack]">
      <q-input v-model="newValue" label="props.fieldLabel" dark />
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
          label="Done"
          @click="onOKClick"
          :class="[styles.maxRight]"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
