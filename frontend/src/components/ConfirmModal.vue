<script setup>
  import {
    useDialogPluginComponent,
    QBtn,
    QCard,
    QCardActions,
    QDialog,
  } from 'quasar';
  import * as styles from '../style.module.css';

  const props = defineProps({ message: String });

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const onOKClick = () => {
    onDialogOK();
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.bgBlack]">
      <div>{{ props.message }}</div>
      <q-card-actions align="right">
        <q-btn
          color="accent"
          label="No"
          text-color="dark"
          @click="onDialogCancel"
        />
        <q-btn
          color="accent"
          text-color="dark"
          label="OK"
          @click="onOKClick"
          :class="[styles.maxRight]"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
