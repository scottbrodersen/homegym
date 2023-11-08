<script setup>
  import { exerciseTypeStore } from '../modules/state';
  import {
    fetchExerciseTypes,
    addExerciseType,
    updateExerciseType,
    intensityTypes,
    volumeTypes,
  } from '../modules/utils';
  import styles from '../style.module.css';
  import { computed, ref, onBeforeMount } from 'vue';

  const emptyType = () => {
    return {
      id: '',
      name: '',
      intensityType: '',
      volumeType: '',
      volumeConstraint: 1,
    };
  };

  // local state for selected activity
  const currentExerciseType = ref(emptyType());

  const states = {
    READ_ONLY: 0,
    EDIT: 1,
    NEW: 2,
  };

  const state = ref(states.READ_ONLY);

  const isTypeValid = computed(() => {
    if (state.value == states.READ_ONLY) {
      return true;
    }

    if (currentExerciseType.value.name.length < 3) {
      return false;
    }
    if (
      !currentExerciseType.value.intensityType ||
      !currentExerciseType.value.volumeType
    ) {
      return false;
    }

    if (
      currentExerciseType.value.volumeType == 'count' &&
      (currentExerciseType.value.volumeConstraint !== 1 ||
        currentExerciseType.value.volumeConstraint !== 2)
    ) {
      return false;
    }

    return true;
  });

  const isChanged = computed(() => {
    if (state.value == states.READ_ONLY) {
      return false;
    }
    // true for new type
    if (state.value == states.NEW) {
      return true;
    }

    // true when any value has changed from stored
    const storedEntry = exerciseTypeStore.get(currentExerciseType.value.id);
    for (const [key, value] of Object.entries(storedEntry)) {
      if (value != currentExerciseType.value[key]) {
        return true;
      }
    }
    return false;
  });

  const setCurrentExercise = (exerciseType) => {
    currentExerciseType.value = Object.assign(
      currentExerciseType.value,
      exerciseType
    );
  };

  const setNewState = () => {
    state.value = states.NEW;
    currentExerciseType.value = emptyType();
  };

  const saveType = async () => {
    // volumeConstraint value is fuzzy when not 2
    if (currentExerciseType.value.volumeConstraint != 2) {
      currentExerciseType.value.volumeConstraint =
        currentExerciseType.value.volumeType === 'count' ? 1 : 0;
    }
    // if id, update existing
    // no id, create new and handle the returned id
    try {
      if (!currentExerciseType.value.id) {
        currentExerciseType.value.id = await addExerciseType(
          currentExerciseType.value.name,
          currentExerciseType.value.intensityType,
          currentExerciseType.value.volumeType,
          currentExerciseType.value.volumeConstraint
        );
      } else {
        await updateExerciseType(currentExerciseType.value);
      }
    } catch (error) {
      if (error instanceof ErrNotLoggedIn) {
        console.log(error.message);
        authPrompt(saveType);
      } else {
        throw error;
      }
    }
  };

  const resetValues = () => {
    // if id, keep selected and reset from store
    // if no id, clear values
    if (!!currentExerciseType.value.id) {
      setCurrentExercise(exerciseTypeStore.get(currentExerciseType.value.id));
    } else {
      setCurrentExercise(emptyType());
    }
  };

  const resetAndCancel = () => {
    resetValues();
    state.value = states.READ_ONLY;
  };

  onBeforeMount(async () => {
    if (exerciseTypeStore.exerciseTypes.length === 0) {
      try {
        await fetchExerciseTypes();
      } catch (e) {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(seetup);
        } else {
          console.log(e.message);
        }
      }
    }
  });
</script>
<template>
  <div :class="[styles.grid2Col]">
    <div :class="[styles.colTitleWrapper, styles.leftColumn]">
      <div :class="[styles.listTitle, styles.sibSpSmall]">Exercises</div>
      <div :class="[styles.sibSpSmall]">
        <q-btn
          size="0.65em"
          round
          icon="add"
          color="primary"
          :disable="state != states.READ_ONLY"
          @Click="setNewState"
        />
      </div>
    </div>
    <div :class="[styles.leftColumn, styles.horiz, styles.blockBorder]">
      <q-list :class="[styles.listStd]" bordered separator>
        <q-item
          clickable
          :disable="state != states.READ_ONLY"
          v-for="[id, exType] in exerciseTypeStore.exerciseTypes"
          :key="id"
          @Click.stop="setCurrentExercise(exType)"
        >
          <q-item-section>
            <q-item-label>{{ exType.name }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-icon
              name="start"
              color="green"
              size="sm"
              v-show="currentExerciseType.id == id"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </div>
    <div :class="[styles.colTitleWrapper, styles.rightColumn]">
      <div :class="[styles.listTitle, styles.sibSpSmall]">Properties</div>
      <div :class="[styles.sibSpSmall]">
        <q-btn
          size="0.65em"
          round
          icon="edit"
          color="primary"
          :disable="state != states.READ_ONLY || !currentExerciseType.id"
          @Click="state = states.EDIT"
        />
      </div>
    </div>
    <div :class="[styles.rightColumn, styles.blockBorder]">
      <div
        v-if="currentExerciseType.id == '' && state == states.READ_ONLY"
        :class="[styles.blockPadSm]"
      >
        Select an exercise
      </div>
      <div v-else>
        <q-input
          v-model="currentExerciseType.name"
          filled
          type="text"
          label="Exercise Name"
          :readonly="state == states.READ_ONLY"
          dark
        />
        <q-select
          v-model="currentExerciseType.intensityType"
          :options="intensityTypes"
          filled
          label="Intensity"
          emit-value
          :readonly="state == states.READ_ONLY"
          dark
        />
        <q-select
          v-model="currentExerciseType.volumeType"
          :options="volumeTypes"
          filled
          label="Volume"
          :readonly="state == states.READ_ONLY"
          emit-value
          dark
        />
        <q-checkbox
          v-show="currentExerciseType.volumeType === 'count'"
          v-model.number="currentExerciseType.volumeConstraint"
          :true-value="Number('2')"
          :false-value="Number('1')"
          :toggle-indeterminate="false"
          label="Track Failed reps"
          :disable="state == states.READ_ONLY"
          dark
        />
        <div v-show="state != states.READ_ONLY" :class="[styles.horiz]">
          <q-btn
            color="primary"
            icon="save"
            @click="saveType"
            :disable="!isChanged && isTypeValid"
          />
          <q-btn
            color="primary"
            icon="restart_alt"
            @click="resetValues"
            :disable="
              (!isChanged && state == states.READ_ONLY) || state != states.NEW
            "
          />
          <q-btn color="primary" label="done" @click="resetAndCancel" />
        </div>
      </div>
    </div>
  </div>
</template>
