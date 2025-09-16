<script setup>
  /**
   * A generic dialog for editing one or more values.
   * Props:
   *  values is an array of {label, value} objects, where label is the text to display and value is the value.
   *
   * Emits an array similar to that of the values props, but with the new values.
   */
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
