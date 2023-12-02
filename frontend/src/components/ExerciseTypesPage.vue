<script setup>
  import { exerciseTypeStore } from '../modules/state';
  import {
    fetchExerciseTypes,
    addExerciseType,
    updateExerciseType,
    intensityTypes,
    volumeTypes,
    openCompositionModal,
    openVariationModal,
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
      composition: {},
      basis: '',
    };
  };

  // local state for selected activity
  const currentExerciseType = ref(emptyType());
  const setComposition = (composition) => {
    currentExerciseType.value.composition = structuredClone(composition);
  };
  const setBasis = (id) => {
    currentExerciseType.value.basis = id;
    if (id == '') {
      isVariation.value = false;
    }
  };

  // model for composition checkbox
  const isComposite = ref(false);

  // model for variation checkbox
  const isVariation = ref(false);

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

    isComposite.value = !!currentExerciseType.value.composition ? true : false;
    isVariation.value = !!currentExerciseType.value.basis ? true : false;
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
    // if volume type is not count, make sure volume constraint is 0
    // and composite is empty
    if (currentExerciseType.value.volumeType !== 'count') {
      currentExerciseType.value.volumeConstraint = 0;
      currentExerciseType.value.composition = {};
    }

    // if not a composite or variation, make sure the associated references are empty
    if (!isComposite.value) {
      currentExerciseType.value.composition = {};
    }

    if (!isVariation.value) {
      currentExerciseType.value.basis = '';
    }

    // if id, update existing
    // no id, create new and handle the returned id
    try {
      if (!currentExerciseType.value.id) {
        currentExerciseType.value.id = await addExerciseType(
          currentExerciseType.value.name,
          currentExerciseType.value.intensityType,
          currentExerciseType.value.volumeType,
          currentExerciseType.value.volumeConstraint,
          currentExerciseType.value.basis
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
          dense
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
      <div v-else :class="[styles.vert]">
        <q-input
          v-model="currentExerciseType.name"
          filled
          type="text"
          label="Exercise Name"
          :readonly="state == states.READ_ONLY"
          dark
          v-focus
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
          :disable="state == states.READ_ONLY || isComposite"
          dark
        />

        <div :class="[styles.horiz]">
          <q-checkbox
            v-show="currentExerciseType.volumeType === 'count'"
            v-model="isVariation"
            label="Variation"
            :disable="state == states.READ_ONLY || isComposite"
            dark
          />
          <div
            v-show="isVariation && currentExerciseType.volumeType === 'count'"
            :class="[styles.maxRight, styles.blockPadSm]"
          >
            <q-btn
              round
              icon="arrow_right_alt"
              :disable="state == states.READ_ONLY || isComposite"
              color="primary"
              @click="
                openVariationModal(
                  currentExerciseType.id,
                  currentExerciseType.basis,
                  setBasis
                )
              "
              size="0.65em"
            />
          </div>
        </div>
        <div v-show="isVariation" :class="[styles.blockPadSm, styles.compShow]">
          <div v-if="!!currentExerciseType.basis">
            {{ exerciseTypeStore.get(currentExerciseType.basis).name }}
          </div>
        </div>

        <div :class="[styles.horiz]">
          <q-checkbox
            v-show="currentExerciseType.volumeType === 'count'"
            v-model="isComposite"
            label="Composite exercise"
            :disable="
              state == states.READ_ONLY ||
              currentExerciseType.volumeConstraint == 2 ||
              isVariation
            "
            dark
          />
          <div
            v-show="isComposite && currentExerciseType.volumeType === 'count'"
            :class="[styles.maxRight, styles.blockPadSm]"
          >
            <q-btn
              round
              icon="arrow_right_alt"
              :disable="
                state == states.READ_ONLY ||
                currentExerciseType.volumeConstraint == 2
              "
              color="primary"
              @click="
                openCompositionModal(
                  currentExerciseType.id,
                  currentExerciseType.composition,
                  setComposition
                )
              "
              size="0.65em"
            />
          </div>
        </div>
        <div
          v-show="isComposite && currentExerciseType.volumeType === 'count'"
          :class="[styles.blockPadSm, styles.compShow]"
        >
          <div v-for="(value, key) in currentExerciseType.composition">
            {{ exerciseTypeStore.get(key).name }} x
            {{ value }}
          </div>
        </div>
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
