<script setup>
  import { ref, toRaw } from 'vue';
  import { useDialogPluginComponent } from 'quasar';
  import styles from '../style.module.css';
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

  // reactive
  const exerciseType = exerciseTypeStore.get(props.exerciseTypeID);

  function onOKClick() {
    const perfObj = {
      segmentIndex: props.segmentIndex,
    };

    // transform volume values as needed before emitting
    // volume.value.forEach((vol, i) => {
    //   if (exerciseType.volumeConstraint === 1) {
    //     volume.value[i] = repCountToBinaryReps(vol);
    //   }
    // });

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

  // adds a set to the volume array
  const incrementVolume = () => {
    volume.value.push([]);
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

  // Transforms a rep count to the format used for storing
  const repCountToBinaryReps = (count) => {
    const arr = new Array();
    for (let i = 0; i < count; i++) {
      arr.push(1);
    }

    return arr;
  };
</script>

<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin">
      <div
        :class="[
          styles.horiz,
          styles.maxSpacing,
          styles.alignCenter,
          styles.blockPadSm,
          styles.dark,
        ]"
      >
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
        <div>
          <q-btn
            round
            color="primary"
            icon="add_circle"
            @Click="incrementVolume"
          />
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
            :class="[styles.horiz, styles.maxSpacing]"
            v-if="exerciseType.volumeConstraint === 2"
          >
            <div :class="[styles.horiz, styles.maxSpacing]">
              <div :class="styles.sibSpSmall">
                <q-btn
                  :class="styles.sibSpxSmall"
                  round
                  color="primary"
                  icon="thumb_up"
                  @Click="syncVolume(i, 1)"
                />
                <q-btn
                  :class="styles.sibSpxSmall"
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
            <div>
              <q-btn
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
        <q-btn color="primary" label="Done" @click="onOKClick" />
        <q-btn color="primary" label="Cancel" @click="onDialogCancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
