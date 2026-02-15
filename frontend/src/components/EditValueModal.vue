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
  import { computed, ref } from 'vue';
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

  // returns true if none of the values has not changed
  const disableSave = computed(() => {
    let disable = true;
    for (let i = 0; i < props.values.length; i++) {
      if (
        newValues.value[i] != props.values[i].value &&
        newValues.value[i].length > 3
      ) {
        disable = false;
        break;
      }
      if (newValues.value[i].length < 3) {
        disable = true;
        break;
      }
    }
    return disable;
  });
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.bgBlack]">
      <div v-for="(valueObj, ix) in props.values">
        <q-input
          v-model="newValues[ix]"
          :label="props.values[ix].label"
          dark
          :rules="[
            (val) =>
              (val != props.values[ix].value && val.length >= 3) ||
              'Value must be different. Min length 3',
          ]"
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
          label="Done"
          @click="onOKClick"
          :class="[styles.maxRight]"
          :disable="disableSave"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
