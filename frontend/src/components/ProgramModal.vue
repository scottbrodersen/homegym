<script setup>
  /**
   * A dialog for editing a program.
   *
   * Provides the following values to child components:
   *  state of EDIT (the edit/read-only state of a component)
   *  activity is the ID of the activity with which the program is associated
   *
   * The dialog output is the program object.
   */
  import { provide, ref } from 'vue';
  import * as styles from '../style.module.css';
  import {
    useDialogPluginComponent,
    QDialog,
    QCard,
    QCardActions,
    QBtn,
    QExpansionItem,
  } from 'quasar';
  import ProgramBlock from './ProgramBlock.vue';
  import ProgramMicrocycle from './ProgramMicrocycle.vue';
  import ProgramWorkout from './ProgramWorkout.vue';
  import * as utils from '../modules/utils';
  import ListActions from './ListActions.vue';

  defineEmits([...useDialogPluginComponent.emits]);
  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const onOKClick = () => {
    onDialogOK(program.value);
  };

  const props = defineProps({ program: Object });

  const program = ref(utils.deepToRaw(props.program));
  const blocks = new utils.OrderedList(program.value.blocks);
  const cycles = new Array();
  let workouts;

  program.value.blocks.forEach((block) => {
    cycles.push(new utils.OrderedList(block.microCycles));
  });
  const state = ref(utils.states.EDIT);

  provide('state', { state });
  provide('activity', props.program.activityID);

  const initWorkouts = (blockIndex, microCycleIndex) => {
    // called on expansion item show
    workouts = new utils.OrderedList(
      program.value.blocks[blockIndex].microCycles[microCycleIndex].workouts
    );
  };
</script>
<template>
  <q-dialog persistent full-width ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin">
      <div :class="[styles.blockPadSm]">
        <h1 :class="[styles.pgmTitle]">Editing {{ program.title }}</h1>
        <div v-for="(block, bix) of program.blocks" :key="bix">
          <div :class="[styles.horiz]">
            <ProgramBlock :block="block" />
            <div>
              <ListActions @update="(action) => blocks.update(action, bix)" />
            </div>
          </div>
          <div
            v-for="(cycle, cix) of block.microCycles"
            :key="cix"
            :class="[styles.horiz]"
          >
            <div :class="[styles.pgmMicrocycleEdit]">
              <ProgramMicrocycle :microcycle="cycle" />
              <q-expansion-item
                switch-toggle-side
                group="workouts"
                label="workouts"
                @before-show="(evt) => initWorkouts(bix, cix)"
              >
                <div
                  v-for="(workout, wix) of cycle.workouts"
                  :class="[styles.horiz]"
                >
                  <ProgramWorkout :workout="workout" />
                  <ListActions
                    @update="(action) => workouts.update(action, wix)"
                  />
                </div>
              </q-expansion-item>
            </div>
            <ListActions
              @update="(action) => cycles[bix].update(action, cix)"
            />
          </div>
        </div>
        <q-card-actions :class="[styles.pgmActions]" dark align="between">
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
      </div>
    </q-card>
  </q-dialog>
</template>
