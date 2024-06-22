<script setup>
  import { inject, onMounted, ref, watch } from 'vue';
  import ProgramBlock from './ProgramBlock.vue';
  import * as utils from '../modules/utils';
  import { QBtn, QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import { programsStore } from '../modules/state';
  import * as programUtils from '../modules/programUtils';
  import ProgramMap from './ProgramMap.vue';
  import ProgramWorkout from './ProgramWorkout.vue';
  import ProgramMicrocycle from './ProgramMicrocycle.vue';
  import ListActions from './ListActions.vue';

  const props = defineProps({ programID: String });
  const emit = defineEmits(['done']);

  const { editTitle, toggleEditTitle } = inject('editTitle');
  const { state } = inject('state');
  const activityID = inject('activity');

  const program = ref({});
  const changed = ref(false);
  const valid = ref(true);

  const defaultBlockTitle = 'Block';
  const defaultMicroCycleTitle = 'MicroCycle';

  let workoutOrderedList;

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
      program.value = utils.deepToRaw(
        programsStore.get(activityID, props.programID)
      );
      if (!program.value.blocks) {
        program.value.blocks = [{}];
      }
      baseline = JSON.stringify(program.value);
      changed.value = false;
      //blocks = new utils.OrderedList(program.value.blocks);
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
      state.value = utils.states.READ_ONLY;
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

  // open the new program modal when state changes to NEW
  watch(
    () => {
      return state.value;
    },
    (newState) => {
      if (newState == utils.states.NEW) {
        utils.newProgramModal(activityID, initProgram);
      }
    }
  );

  // Open modal to edit the instance title
  watch(
    () => editTitle.value,
    (newValue) => {
      console.log('editTitle change: ' + newValue);
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
              console.log(newValue);
              program.value.title = newValue[0];
              await saveProgram(program.value);
            }
            toggleEditTitle();
          });
      }
    }
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
    { deep: true }
  );

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
    // set height of workouts div
    document.getElementById('wo-wrap').style['max-height'] = `${
      window.innerHeight -
      document.getElementsByTagName('header')[0].offsetHeight -
      document.getElementById('program-map').offsetHeight -
      document.getElementById('pgm-context').offsetHeight -
      20
    }px`;
  });

  const editBlocks = () => {
    utils.editProgramModal(program.value).then(async (updatedProgram) => {
      if (updatedProgram) {
        await saveProgram(updatedProgram);
      }
    });
  };
  const updateWorkouts = (action, workoutIndex) => {
    workoutOrderedList.update(action, workoutIndex);
  };
</script>
<template>
  <div :class="[styles.pgmWrap]">
    <div id="program-map" :class="[styles.programMap]">
      <ProgramMap
        :blocks="program.blocks"
        @coords="(value) => (coords = value)"
        :class="[styles.centered]"
      />
      <div :class="[styles.pgmEdit]">
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
