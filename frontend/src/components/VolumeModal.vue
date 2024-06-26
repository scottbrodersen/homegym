<script setup>
  import { ref, toRaw } from 'vue';
  import {
    useDialogPluginComponent,
    QDialog,
    QCard,
    QCardActions,
    QBtn,
    QInput,
  } from 'quasar';
  import * as styles from '../style.module.css';
  import VolumeReps from './VolumeReps.vue';
  import { exerciseTypeStore, unitsState } from '../modules/state';

  const emit = defineEmits([...useDialogPluginComponent.emits, 'update']);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const props = defineProps({
    exerciseTypeID: String,
    intensity: Number,
    segmentIndex: Number,
    volume: Array,
  });

  // model for the edit exercise dialog
  const volume = ref([]);

  if (!!props.volume) {
    volume.value = JSON.parse(JSON.stringify(props.volume));
  }

  // adds a set to the volume array
  const incrementVolume = () => {
    volume.value.push([]);
  };

  if (volume.value.length == 0) {
    incrementVolume();
  }

  // reactive
  const exerciseType = exerciseTypeStore.get(props.exerciseTypeID);

  function onOKClick() {
    const perfObj = {
      segmentIndex: props.segmentIndex,
    };

    perfObj.volume = volume.value;

    onDialogOK(perfObj);
  }

  const countReps = (volArray) => {
    let count = volArray.length;
    volArray.forEach((rep) => {
      if (rep === 0) {
        count--;
      }
    });
    return count;
  };

  // updates the volume array for control events
  const syncVolume = (index, value) => {
    if (
      exerciseType.volumeConstraint === 2 ||
      exerciseType.volumeConstraint === 0
    ) {
      if (value === -1) {
        volume.value[index].pop();
      } else {
        volume.value[index].push(value);
      }
    } else if (exerciseType.volumeConstraint === 1) {
      const vol = [];
      for (let i = 0; i < value; i++) {
        vol.push(1);
      }
      volume.value[index] = vol;
    }
  };
</script>

<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin">
      <div :class="[styles.horiz, styles.blockPadSm, styles.dark]">
        <div :class="[styles.exName]">
          {{ exerciseType.name }}
        </div>
        <div
          v-if="exerciseType.intensityType != 'bodyweight'"
          :class="[styles.exName]"
        >
          {{ exerciseType.intensityType }}:
          {{ props.intensity }}
          {{ unitsState[exerciseType.intensityType] }}
        </div>
        <div :class="[styles.maxRight]">
          <q-btn round color="primary" icon="add" @Click="incrementVolume" />
        </div>
      </div>
      <div :class="[styles.blockPadSm, styles.dark]">
        <div :class="styles.blockPadSm" v-for="(v, i) in volume" :key="i">
          <q-input
            v-if="exerciseType.volumeConstraint === 1"
            :model-value="countReps(volume[i])"
            type="number"
            filled
            dark
            v-focus
            v-select
            style="width: 100px"
            @update:model-value="
              (value) => {
                syncVolume(i, value);
              }
            "
          />
          <div
            :class="[styles.horiz]"
            v-if="exerciseType.volumeConstraint === 2"
          >
            <div :class="[styles.horiz]">
              <div :class="styles.sibSpSmall">
                <q-btn
                  round
                  color="primary"
                  icon="thumb_up"
                  @Click="syncVolume(i, 1)"
                />
                <q-btn
                  round
                  color="primary"
                  icon="thumb_down"
                  @Click="syncVolume(i, 0)"
                />
              </div>
              <VolumeReps
                :reps="v"
                :volume-constraint="2"
                :class="[styles.sibSpSmall]"
              />
            </div>
            <div :class="[styles.maxRight]">
              <q-btn
                v-show="v.length > 0"
                round
                color="primary"
                icon="backspace"
                @Click="syncVolume(i, -1)"
              />
            </div>
          </div>
        </div>
      </div>
      <q-card-actions dark align="between" :class="[styles.blockBorder]">
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
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
