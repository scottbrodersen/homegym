<script setup>
  import { useDialogPluginComponent, QCheckbox, QList } from 'quasar';
  import { ref } from 'vue';
  import { exerciseTypeStore } from '../modules/state';
  import styles from '../style.module.css';

  const props = defineProps({ exerciseID: String, composition: Object });
  const composition = ref({});
  const exerciseList = ref([]);
  const selectedIDs = ref([]);

  for (const id in props.composition) {
    selectedIDs.value.push(id);
    composition.value[id] = props.composition[id];
  }

  exerciseTypeStore.exerciseTypes.forEach((value, key) => {
    if (key != props.exerciseID && !!!value.composition) {
      exerciseList.value.push(value);
    }
  });

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const onOKClick = () => {
    // composition can contain values for exercises that were selected then deselected
    // return only the values for selected exercises
    const comps = {};
    selectedIDs.value.forEach((exerciseID) => {
      if (!!composition.value[exerciseID]) {
        comps[exerciseID] = composition.value[exerciseID];
      }
    });
    onDialogOK(comps);
  };
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card
      dark
      class="q-dialog-plugin"
      :class="[styles.bgBlack, styles.compositeDialog]"
    >
      <q-list :class="[styles.listStd]" dense>
        <q-item v-for="eType in exerciseList" :key="eType.id">
          <div :class="[styles.horiz, styles.listStd]">
            <q-checkbox
              v-model="selectedIDs"
              :val="eType.id"
              :label="eType.name"
              dark
            />
            <q-input
              v-if="selectedIDs.includes(eType.id)"
              mask="#"
              :model-value="composition[eType.id]"
              @update:model-value="
                (value) => (composition[eType.id] = Number(value))
              "
              label="Reps"
              stack-label
              dark
              :class="[styles.compRep, styles.maxRight]"
            />
          </div>
        </q-item>
      </q-list>
      <q-card-actions align="right">
        <q-btn
          color="accent"
          label="Cancel"
          text-color="dark"
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
