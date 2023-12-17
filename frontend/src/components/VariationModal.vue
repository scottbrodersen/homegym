<script setup>
  import { useDialogPluginComponent, QOptionGroup } from 'quasar';
  import { ref, toValue } from 'vue';
  import { exerciseTypeStore } from '../modules/state';
  import styles from '../style.module.css';

  const props = defineProps({ exerciseID: String, basisID: String });
  const selectedID = ref(props.basisID);
  const options = ref([]);

  exerciseTypeStore.exerciseTypes.forEach((value, key) => {
    if (key != props.exerciseID) {
      options.value.push({ label: value.name, value: key });
    }
  });

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const onOKClick = () => {
    onDialogOK(toValue(selectedID));
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card
      dark
      class="q-dialog-plugin"
      :class="[styles.bgBlack, styles.compositeDialog]"
    >
      <q-option-group
        :class="[styles.listStd]"
        v-model="selectedID"
        :options="options"
        dark
      />
      <q-card-actions align="right">
        <q-btn
          color="accent"
          text-color="dark"
          label="Clear"
          @click="selectedID = ''"
          :class="[styles.maxLeft]"
        />
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
