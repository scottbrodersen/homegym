<script setup>
  /**
   * Displays a program including the blocks, microcycles, and workouts.
   * Enables editing of the program.
   *
   * Props:
   *  activityID: The ID of an activity to preselect.
   *  programID: The ID of a program to preselect.
   */
  import { inject, onActivated, onMounted, provide, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import * as utils from '../modules/utils';
  import { QBtn } from 'quasar';
  import * as styles from '../style.module.css';
  import { programsStore } from '../modules/state';
  import * as programUtils from '../modules/programUtils';
  import ProgramMap from './ProgramMap.vue';
  import ProgramWorkout from './ProgramWorkout.vue';
  import ProgramMicrocycle from './ProgramMicrocycle.vue';

  const props = defineProps({ activityID: String, programID: String });
  const emit = defineEmits(['done']);

  const { editProgramTitle, toggleProgramTitle } = inject('editProgramTitle');
  const { newProgram, toggleNewProgram } = inject('newProgram');
  const { state } = inject('state');
  const { cloneProgram, toggleCloneProgram } = inject('cloneProgram');

  provide('activity', props.activityID);

  const program = ref({});
  const changed = ref(false);
  const valid = ref(true);

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';

  let baseline = ''; // use to detect diff
  /*
  Coords is a 3x1 array that holds the coordinates of the workout for a date.
  E.g. [0,1,2] denotes the workout in the 3rd day of the 2nd microcycle in the 1st block.
  */
  const coords = ref([0, 0, 0]);

  // Stores the stringified program as a baseline for detecting change
  // Clones the program so it can be edited without immediately changing the store
  const init = () => {
    if (!props.programID) {
      baseline = '';
      program.value = {};
    } else {
      program.value = utils.deepToRaw(
        programsStore.get(props.activityID, props.programID),
      );
      if (!program.value.blocks) {
        program.value.blocks = [{}];
      }
      baseline = JSON.stringify(program.value);
      changed.value = false;
    }
  };

  // Re-initialize when a different program is selected
  watch(
    () => {
      return props.programID;
    },
    (newID) => {
      init();
    },
  );

  // create new program
  const initProgram = async (programProps) => {
    if (!programProps) {
      return;
    }
    program.value = {
      id: null,
      title: programProps.title,
      activityID: props.activityID,
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
    await saveProgram(program.value);
  };

  // open the new program modal when newProgram is true
  watch(
    () => newProgram.value,
    (newProgramValue) => {
      if (newProgramValue) {
        utils.newProgramModal(props.activityID).then(async (programProps) => {
          if (programProps) {
            await initProgram(programProps);
          }
          toggleNewProgram();
        });
      }
    },
  );

  // Open modal to edit the instance title
  watch(
    () => editProgramTitle.value,
    (newValue) => {
      if (newValue === true) {
        utils
          .openEditValueModal([
            {
              label: 'Program Title',
              value: program.value.title,
            },
          ])
          .then(async (newValue) => {
            if (newValue) {
              program.value.title = newValue[0];
              await saveProgram(program.value);
            }
            toggleProgramTitle();
          });
      }
    },
  );

  // Open modal to clone the program
  watch(
    () => cloneProgram.value,
    (newValue) => {
      let clonedProgram = {};
      if (newValue === true) {
        // make a static copy of the program
        clonedProgram = utils.deepToRaw(
          programsStore.get(props.activityID, props.programID),
        );
        clonedProgram.id = '';
        utils
          .openEditValueModal([
            {
              label: 'Cloned Program Title',
              value: clonedProgram.title,
            },
          ])
          .then(async (newValue) => {
            if (newValue) {
              clonedProgram.title = newValue[0];
              await saveProgram(clonedProgram);
            }
            toggleCloneProgram();
          });
      }
    },
  );
  init();

  const saveProgram = async (p) => {
    try {
      const id = await utils.updateProgram(p);
      program.value = p;
      baseline = JSON.stringify(program.value);

      utils.toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof utils.ErrNotLoggedIn) {
        await utils.authPromptAsync();
        saveProgram(p);
      } else {
        utils.toast('Error', 'negative');
      }
    }
    if (state.value == utils.states.NEW) {
      state.value = utils.states.READ_ONLY;
    }
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
    { deep: true },
  );

  watch(
    () => coords.value,
    (newCoords) => {
      scrollToWorkout(newCoords);
    },
  );

  const scrollToWorkout = (coords) => {
    const wo = document.getElementById(
      `workout${coords[0]}-${coords[1]}-${coords[2]}`,
    );
    if (wo) {
      wo.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
        inline: 'center',
      });
    }
  };

  onActivated(() => {
    const docsContext = ref(inject('docsContext'));
    docsContext.value = 'programs';
  });

  onMounted(() => {
    // set height of workouts div
    const el = document.getElementById('wo-wrap');
    if (el) {
      document.getElementById('wo-wrap').style['max-height'] = `${
        window.innerHeight -
        document.getElementsByTagName('header')[0].offsetHeight -
        document.getElementById('program-map').offsetHeight -
        document.getElementById('pgm-context').offsetHeight -
        20
      }px`;
    }
  });

  const editBlocks = () => {
    utils.editProgramModal(program.value).then(async (updatedProgram) => {
      if (updatedProgram) {
        await saveProgram(updatedProgram);
      }
    });
  };
</script>
<template>
  <div>
    <div v-if="props.programID" :class="[styles.pgmWrap]">
      <div id="program-map" :class="[styles.programMap]">
        <ProgramMap
          :blocks="program.blocks"
          @coords="(value) => (coords = value)"
          :class="[styles.centered]"
        />
        <div>
          <q-btn
            icon="edit_note"
            round
            dark
            @click="editBlocks"
            color="primary"
          />
        </div>
      </div>
      <div :class="[styles.blockPadSm]">
        <div id="pgm-context">
          <div :class="[styles.pgmBlockTitle]">
            <div>Block:</div>
            <Transition>
              <ProgramBlock :block="program.blocks[coords[0]]" />
            </Transition>
          </div>
          <div :class="[styles.pgmCycleTitle]">
            <div>Cycle:</div>
            <Transition>
              <ProgramMicrocycle
                :microcycle="program.blocks[coords[0]].microCycles[coords[1]]"
                @update="(value) => updateCycles(value, ix)"
              />
            </Transition>
          </div>
        </div>
        <div id="wo-wrap">
          <div
            v-for="(workout, wix) of program.blocks[coords[0]].microCycles[
              coords[1]
            ].workouts"
            :key="wix"
            :class="
              wix == coords[2]
                ? [styles.horiz, styles.pgmSelected]
                : [styles.horiz]
            "
          >
            <ProgramWorkout
              :id="`workout${coords[0]}-${coords[1]}-${wix}`"
              :workout="workout"
              @click="
                () => {
                  coords[2] = wix;
                }
              "
            />
          </div>
        </div>
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
