<script setup>
  import { computed, inject, onMounted, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import {
    authPromptAsync,
    ErrNotLoggedIn,
    newProgramModal,
    OrderedList,
    states,
    toast,
    updateProgram,
    deepToRaw,
  } from '../modules/utils';
  import { QBtn, QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import { programsStore } from '../modules/state';
  import * as programUtils from '../modules/programUtils';
  import ProgramMap from './ProgramMap.vue';
  import ProgramWorkout from './ProgramWorkout.vue';

  const props = defineProps({ programID: String });
  const emit = defineEmits(['done']);

  const { state } = inject('state');
  const activityID = inject('activity').value;

  const program = ref({});
  const changed = ref(false);
  const valid = ref(true);

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';

  let blocks; //
  let baseline = ''; // use to detect diff
  const coords = ref([0, 0, 0]);

  // Stores the stringified program as a baseline for detecting change
  // Clones the program so it can be edited without immediately changing the store
  const init = () => {
    if (!props.programID) {
      baseline = '';
      program.value = {};
    } else {
      program.value = deepToRaw(programsStore.get(activityID, props.programID));
      if (!program.value.blocks) {
        program.value.blocks = [{}];
      }
      baseline = JSON.stringify(program.value);
      changed.value = false;
      blocks = new OrderedList(program.value.blocks);
    }
  };

  // Re-initialize when a different program is selected
  watch(
    () => {
      return props.programID;
    },
    (newID) => {
      init();
    }
  );

  // callback for new program modal
  const initProgram = (programProps) => {
    if (!programProps) {
      state.value = states.READ_ONLY;
      return;
    }
    program.value = {
      id: null,
      title: programProps.title,
      activityID: activityID,
      blocks: new Array(),
    };

    for (let i = 0; i < programProps.numBlocks; i++) {
      program.value.blocks.push({
        title: `${defaultBlockTitle} ${i + 1}`,
        intensity: null,
        microCycles: new Array(),
      });
      for (let j = 0; j < programProps.numCycles; j++) {
        program.value.blocks[i].microCycles.push({
          title: `${defaultMicroCycleTitle} ${j + 1}`,
          span: programProps.cycleSpan,
          intensity: null,
          workouts: new Array(),
        });
        for (let k = 0; k < programProps.cycleSpan; k++) {
          program.value.blocks[i].microCycles[j].workouts.push({
            title: `Day ${k + 1}`,
            segments: [{ exerciseTypeID: '', prescription: '' }],
          });
        }
      }
    }
  };

  watch(
    () => {
      return state.value;
    },
    (newState) => {
      if (newState == states.NEW) {
        newProgramModal(activityID, initProgram);
      }
    }
  );

  init();

  const saveProgram = async () => {
    try {
      const id = await updateProgram(program.value);

      baseline = JSON.stringify(program.value);

      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        await authPromptAsync();
        saveProgram();
      } else {
        toast('Error', 'negative');
      }
    }
    if (state.value == states.NEW) {
      state.value = states.EDIT;
    }
  };

  const cancel = () => {
    emit('done', program.value.id);
    changed.value = false;
  };

  const updateBlocks = (action, index) => {
    blocks.update(action, index);
  };

  // watch for changes and validate
  watch(
    () => {
      return program.value;
    },
    (newVal) => {
      changed.value = baseline != JSON.stringify(newVal);
      valid.value = programUtils.programValidator(newVal);
    },
    { deep: true }
  );

  const updateButtonText = computed(() => {
    return !!program.value.id ? 'Update' : 'Add';
  });

  const doneButtonText = computed(() => {
    return changed.value ? 'Cancel' : 'Done';
  });

  watch(
    () => coords.value,
    (newCoords) => {
      scrollToWorkout(newCoords);
    }
  );

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

  onMounted(() => {
    document.getElementById('wo-wrap').style['max-height'] = `${
      window.innerHeight -
      document.getElementsByTagName('header')[0].offsetHeight -
      document.getElementById('program-map').offsetHeight -
      document.getElementById('pgm-context').offsetHeight -
      20
    }px`;
  });
</script>
<template>
  <div :class="[styles.pgmWrap]">
    <div v-show="state != states.READ_ONLY">
      <q-input
        v-model="program.title"
        label="Program Name"
        stack-label
        dark
        :rules="[
          programUtils.requiredFieldValidator,
          programUtils.maxFieldValidator,
        ]"
      />
    </div>
    <div id="program-map" :class="[styles.hgCentered]">
      <ProgramMap
        :blocks="program.blocks"
        @coords="(value) => (coords = value)"
      />
    </div>
    <div v-if="state == states.READ_ONLY" :class="[styles.blockPadSm]">
      <div id="pgm-context">
        <div :class="[styles.pgmBlockTitle]">
          <div>Block:</div>
          <Transition>
            <div :class="[styles.vert]" :key="coords[0]">
              <div>{{ program.blocks[coords[0]].title }}</div>
              <div>{{ program.blocks[coords[0]].description }}</div>
            </div>
          </Transition>
        </div>
        <div :class="[styles.pgmCycleTitle]">
          <div>Cycle:</div>
          <Transition>
            <div :class="[styles.vert]" :key="coords[1]">
              <div>
                {{ program.blocks[coords[0]].microCycles[coords[1]].title }}
              </div>
              <div>
                {{
                  program.blocks[coords[0]].microCycles[coords[1]].description
                }}
              </div>
            </div>
          </Transition>
        </div>
      </div>
      <div id="wo-wrap">
        <ProgramWorkout
          v-for="(workout, wix) of program.blocks[coords[0]].microCycles[
            coords[1]
          ].workouts"
          :key="wix"
          :id="`workout${coords[0]}-${coords[1]}-${wix}`"
          :workout="workout"
          :class="wix == coords[2] ? [styles.pgmSelected] : ''"
        />
      </div>
    </div>
    <div v-show="state != states.READ_ONLY && program.title">
      <ProgramBlock
        v-for="(block, index) of program.blocks"
        :key="index"
        :block="block"
        @update="(value) => updateBlocks(value, index)"
      />
      <div :class="[styles.buttonArray]">
        <q-btn
          :label="doneButtonText"
          color="accent"
          text-color="dark"
          @click="cancel"
        />
        <q-btn
          :label="updateButtonText"
          color="accent"
          text-color="dark"
          @click="saveProgram"
          :disable="!changed || !valid"
        />
      </div>
    </div>
  </div>
</template>
<style>
  .v-enter-active {
    transition: opacity 0.5s ease;
  }

  .v-enter-from,
  .v-leave-to {
    opacity: 0;
  }
</style>
