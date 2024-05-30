<script setup>
  import { computed, inject, onBeforeMount, ref, watch } from 'vue';
  import ProgramBlock2 from './ProgramBlock2.vue';
  import { programInstanceStore, programsStore } from './../modules/state';
  import { updateProgramInstance } from './../modules/utils';
  import { QBtn, QIcon, QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import {
    authPrompt,
    deepToRaw,
    ErrNotLoggedIn,
    states,
    toast,
  } from '../modules/utils';
  import * as programUtils from '../modules/programUtils';
  import ProgramCalendar from './ProgramCalendar.vue';
  import ProgramMicrocycle2 from './ProgramMicrocycle2.vue';
  import ProgramWorkout2 from './ProgramWorkout2.vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';

  const props = defineProps({ instanceID: String });
  const emit = defineEmits(['done']);

  const instance = ref();
  const coords = ref();
  let today;

  let baseline = ''; // use to detect diff
  const changed = ref();
  const valid = ref(true);
  const programTitle = ref();

  const state = inject('state');
  const activityID = inject('activity').value;
  const init = (instanceID) => {
    instance.value = deepToRaw(programInstanceStore.get(instanceID));
    baseline = JSON.stringify(instance.value);

    programTitle.value = instance.value.programID
      ? programsStore.get(activityID, instance.value.programID).title
      : '';

    today = programUtils.getProgramInstanceStatus(instanceID)[3];
  };

  // Re-initialize when a different instance is selected
  watch(
    () => props.instanceID,
    (newID) => {
      init(newID);
    }
  );

  // watch for changes and validate
  watch(
    () => {
      return instance.value;
    },
    (newVal) => {
      if (state.value != states.READ_ONLY) {
        changed.value = baseline != JSON.stringify(newVal);
        valid.value = programUtils.programValidator(newVal);
      }
    },
    { deep: true }
  );

  onBeforeMount(() => {
    init(props.instanceID);
  });

  const saveInstance = async () => {
    try {
      const id = await updateProgramInstance(instance.value);
      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        authPrompt(saveInstance);
      } else {
        toast('Error', 'negative');
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  const cancel = () => {
    emit('done', instance.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  const doneButtonText = computed(() => {
    return changed.value ? 'Cancel' : 'Done';
  });

  const setCoords = (dayIndex) => {
    if (dayIndex != -1) {
      coords.value = programUtils.getWorkoutCoords(instance.value, dayIndex);
      scrollToWorkout(coords.value);
    }
  };

  const scrollToWorkout = (coords) => {
    const wo = document.getElementById(
      `workout${coords[0]}-${coords[1]}-${coords[2]}`
    );
    if (wo) {
      wo.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
        inline: 'center',
      });
    }
  };

  const statusIcon = (coords) => {
    const props = [{ name: '', colour: '' }];
    const day = programUtils.getDayIndex(instance.value, coords);
    if (day > today) {
      props[0].name = 'schedule';
      props[0].colour = 'blue';
    } else if (
      instance.value.blocks[coords[0]].microCycles[coords[1]].workouts[
        coords[2]
      ].restDay
    ) {
    } else if (instance.value.events[day]) {
      props[0].name = 'check_circle';
      props[0].colour = 'green';
      props[0].eventID = instance.value.events[day];
    } else {
      props[0].name = 'cancel';
      props[0].colour = 'red';
    }
    return props;
  };
</script>
<template>
  <div v-if="instance" :class="[styles.pgmInstance]">
    <ProgramCalendar
      :instance="instance"
      @dayIndex="setCoords"
      :class="[styles.centered]"
    />
    <div :class="[styles.instBase]">
      Base Program:
      {{ programTitle }}
    </div>
    <div v-show="state != states.READ_ONLY">
      <q-input
        v-model="instance.title"
        label="Name"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
    </div>
    <div :class="[styles.instInfo]">
      <ProgramBlock2
        v-if="coords"
        :block="instance.blocks[coords[0]]"
        @update="(value) => updateBlocks(value, bix)"
      />
      <ProgramMicrocycle2
        v-if="coords"
        :microcycle="instance.blocks[coords[0]].microCycles[coords[1]]"
      />
      <div v-else>The program was performed in the past. Select a date.</div>
    </div>

    <div id="inst-wrap" v-if="coords">
      <div
        v-for="(workout, wix) of instance.blocks[coords[0]].microCycles[
          coords[1]
        ].workouts"
        :key="wix"
        :class="[styles.instWorkout]"
      >
        <div :class="coords[2] == wix ? [styles.evtHighlight] : ''">
          <ProgramWorkout2
            :id="`workout${coords[0]}-${coords[1]}-${wix}`"
            :workout="workout"
          />
          <div v-show="!workout.restDay">
            <div v-for="(segment, six) of workout.segments" :key="six">
              <ProgramWorkoutSegment :segment="segment" />
            </div>
          </div>
        </div>
        <div>
          <div :class="[styles.instWorkoutStatus, styles.horiz]">
            <div
              v-for="(iconProps, ix) of statusIcon([coords[0], coords[1], wix])"
              :key="ix"
            >
              <q-icon
                v-if="iconProps.name != 'check_circle'"
                :name="iconProps.name"
                :color="iconProps.colour"
              />
              <q-btn
                v-if="iconProps.name == 'check_circle'"
                :icon="iconProps.name"
                :color="iconProps.colour"
                :to="{ name: 'home', query: { event: iconProps.eventID } }"
                round
                dense
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div
      v-show="state != states.READ_ONLY && instance.id"
      :class="[styles.buttonArray]"
    >
      <q-btn
        :label="doneButtonText"
        color="accent"
        text-color="dark"
        @click="cancel"
      />
      <q-btn
        label="Save"
        color="accent"
        text-color="dark"
        @click="saveInstance"
        :disable="!changed || !valid"
      />
    </div>
  </div>
</template>
