<script setup>
  import { computed, ref, Suspense } from 'vue';
  import { exerciseTypeStore } from '../modules/state';
  import ExerciseIntensity from './ExerciseIntensity.vue';
  import ExerciseSelect from './ExerciseSelect.vue';
  import VolumeReps from './VolumeReps.vue';
  import VolumeTime from './VolumeTime.vue';
  import VolumeDistance from './VolumeDistance.vue';
  import * as styles from '../style.module.css';
  import { openVolumeModal } from '../modules/utils';
  import { QDialog, QBtn, QCard, QCardActions, QCardSection } from 'quasar';
  /*
  interface exerciseInstance = {
    index: Number
    typeID: String
    parts: Array[
      {
        intensity: Number,
        // volume is an array of sets, and a set is an array of reps
        volume: Array[Array[Number]]
      }]
  }
  */

  const props = defineProps({
    exerciseInstance: Object,
    activityID: String,
    writable: Boolean,
  });

  const emit = defineEmits(['update']);

  if (!props.activityID && !props.exerciseInstance.index) {
    throw Error(
      'ExerciseInstance requires an activity id and an indexed exercise instance'
    );
  }

  // model
  const instance = ref(JSON.parse(JSON.stringify(props.exerciseInstance)));

  // initialize an empty instance
  if (!instance.value.typeID) {
    instance.value.typeID = '';
    instance.value.parts = [];
  }

  const exerciseName = ref(
    props.exerciseInstance.typeID
      ? exerciseTypeStore.get(props.exerciseInstance.typeID).name
      : ''
  );

  const isCountReps = computed(() => {
    const volumeConstraint = exerciseTypeStore.get(
      instance.value.typeID
    ).volumeConstraint;
    return volumeConstraint === 1;
  });

  // index of instance segment to delete
  const toDelete = ref(null);

  // model for delete confirmation dialog
  const confirmDelete = computed(() => {
    if (toDelete.value != null) {
      return true;
    }
    return false;
  });

  const setExerciseType = (typeID) => {
    instance.value.typeID = typeID;
    if (instance.value.parts.length == 0) {
      addSegments();
    }
    exerciseName.value = exerciseTypeStore.get(typeID).name;
  };

  const addSegments = () => {
    if (
      exerciseTypeStore.get(instance.value.typeID).intensityType == 'hrZone'
    ) {
      for (let i = 1; i < 6; i++) {
        instance.value.parts.push({
          intensity: i,
          volume: [0],
        });
      }
    } else {
      const segment = { intensity: 0, volume: [] };
      const volumeType = exerciseTypeStore.get(
        instance.value.typeID
      ).volumeType;
      if (volumeType == 'time' || volumeType == 'distance') {
        segment.volume.push(0);
      }
      instance.value.parts.push(segment);
    }
  };

  const deleteSegment = () => {
    instance.value.parts.splice(toDelete.value, 1);
    toDelete.value = null;
    emit('update', instance.value);
  };

  const updateVolume = (volumeObj) => {
    instance.value.parts[volumeObj.segmentIndex].volume = volumeObj.volume;
    emit('update', instance.value);
  };

  const updateIntensity = (value, segmentIndex) => {
    instance.value.parts[segmentIndex].intensity = Number(value);
    emit('update', instance.value);
  };
</script>

<template>
  <div v-if="props.writable" :class="[styles.horiz]">
    <div>
      <Suspense>
        <ExerciseSelect
          :activityID="props.activityID"
          :exerciseID="instance.typeID"
          @selectedID="
            (value) => {
              setExerciseType(value);
            }
          "
      /></Suspense>
    </div>
    <div :class="[styles.maxRight]">
      <q-btn
        v-show="!!exerciseName"
        round
        color="primary"
        icon="add"
        @click="addSegments"
      />
    </div>
  </div>
  <div v-else :class="[styles.exName]">{{ exerciseName }}</div>
  <div>
    <div
      :class="[styles.exInstRow]"
      v-for="(part, partIndex) in instance.parts"
      :key="partIndex"
    >
      <div :class="[styles.sibSpMed]" v-if="props.writable">
        <q-btn
          round
          icon="delete"
          color="primary"
          @click="toDelete = partIndex"
        />
      </div>
      <div :class="[styles.sibSpMed]">
        <ExerciseIntensity
          :intensity="Number(part.intensity)"
          :type="exerciseTypeStore.get(instance.typeID).intensityType"
          :writable="props.writable"
          @update="(value) => updateIntensity(value, partIndex)"
        />
      </div>
      <div
        v-if="exerciseTypeStore.get(instance.typeID).volumeType == 'count'"
        :class="[styles.sibSpMed, styles.horiz, styles.hg100wide]"
      >
        <div :class="[styles.volume, isCountReps ? styles.repCountSet : '']">
          <VolumeReps
            v-for="(set, index) in part.volume"
            :key="index"
            :reps="set"
            :volume-constraint="
              exerciseTypeStore.get(instance.typeID).volumeConstraint
            "
          />
        </div>
        <div v-if="props.writable" :class="[styles.maxRight]">
          <q-btn
            round
            icon="arrow_right_alt"
            color="primary"
            @click="
              openVolumeModal(
                instance.typeID,
                part.intensity,
                partIndex,
                part.volume,
                updateVolume
              )
            "
          />
        </div>
      </div>
      <div
        v-else-if="exerciseTypeStore.get(instance.typeID).volumeType == 'time'"
        :class="[styles.sibSpMed]"
      >
        <VolumeTime
          :time="part.volume.length > 0 ? part.volume[0][0] : ''"
          :writable="props.writable"
          @update="
            (value) => {
              updateVolume({ volume: [[value]], segmentIndex: partIndex });
            }
          "
        />
      </div>
      <div v-else :class="[styles.sibSpMed]">
        <VolumeDistance
          :distance="part.volume.length > 0 ? part.volume[0][0] : Number(0.0)"
          :writable="props.writable"
          @update="
            (value) => {
              updateVolume({ volume: [[value]], segmentIndex: partIndex });
            }
          "
        />
      </div>
    </div>
  </div>
  <q-dialog v-model="confirmDelete">
    <q-card dark>
      <q-card-section class="q-pt-none">
        Are you sure you want to delete the exercise?
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          flat
          label="No"
          color="accent"
          @click="toDelete = null"
          v-close-popup
          dark
        />
        <q-btn
          flat
          label="Yes"
          color="accent"
          dark
          @click="deleteSegment()"
          v-close-popup
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
