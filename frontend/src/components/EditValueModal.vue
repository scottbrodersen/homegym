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

  defineEmits([...useDialogPluginComponent.emits]);
  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  // an array of {label, value} objects
  const props = defineProps({ values: Array });

  const newValues = ref(new Array());
  props.values.forEach((valueObj) => {
    newValues.value.push(valueObj.value);
  });

  const onOKClick = () => {
    onDialogOK(newValues.value);
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.bgBlack]">
      <div v-for="(valueObj, ix) in props.values">
        <q-input v-model="newValues[ix]" :label="props.values[ix].label" dark />
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
          label="Done"
          @click="onOKClick"
          :class="[styles.maxRight]"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
