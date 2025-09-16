<script setup>
  /**
   * Provides a visual representation of the blocks, microcycles, and workouts of a program.
   * Establishes the coordinates (coords) of the workouts
   *
   * Coords is a 3x1 array that holds the coordinates of the workout for a date.
   * E.g. [0,1,2] denotes the workout in the 3rd day of the 2nd microcycle in the 1st block.
   *
   * Props:
   *  blocks is an array of program block objects.
   *
   * When a workout is clicked, the workout coordinates are emitted.
   */
  import * as styles from '../style.module.css';
  import { ref } from 'vue';
  import { QBtn } from 'quasar';

  const props = defineProps({ blocks: Array });

  const emit = defineEmits(['coords']);

  const selected = ref([0, 0, 0]);

  const setCoords = (coords) => {
    selected.value = coords;
    emit('coords', coords);
  };

  const isSelected = (coords) => {
    return (
      coords.length === selected.value.length &&
      coords.every((value, index) => value === selected.value[index])
    );
  };
  const keyNav = (event) => {
    switch (event.key) {
      case 'ArrowUp':
        if (selected.value[1] > 0) {
          selected.value[1]--;
        } else if (selected.value[0] > 0) {
          selected.value[0]--;
          selected.value[1] =
            props.blocks[selected.value[0]].microCycles.length - 1;
        }
        break;
      case 'ArrowDown':
        if (
          selected.value[1] <
          props.blocks[selected.value[0]].microCycles.length - 1
        ) {
          selected.value[1]++;
        } else if (selected.value[0] < props.blocks.length - 1) {
          selected.value[0]++;
          selected.value[1] = 0;
        }
        break;
      case 'ArrowLeft':
        if (selected.value[2] > 0) {
          selected.value[2]--;
        } else if (selected.value[1] > 0) {
          selected.value[1]--;
          selected.value[2] =
            props.blocks[selected.value[0]].microCycles[selected.value[1]]
              .workouts.length - 1;
        } else if (selected.value[0] > 0) {
          selected.value[0]--;
          selected.value[1] =
            props.blocks[selected.value[0]].microCycles.length - 1;
          selected.value[2] =
            props.blocks[selected.value[0]].microCycles[selected.value[1]]
              .workouts.length - 1;
        }
        break;
      case 'ArrowRight':
        if (
          selected.value[2] <
          props.blocks[selected.value[0]].microCycles[selected.value[1]]
            .workouts.length -
            1
        ) {
          selected.value[2]++;
        } else if (
          selected.value[1] <
          props.blocks[selected.value[0]].microCycles.length - 1
        ) {
          selected.value[1]++;
          selected.value[2] = 0;
        } else if (selected.value[0] < props.blocks.length - 1) {
          selected.value[0]++;
          selected.value[1] = 0;
          selected.value[2] = 0;
        }
        break;
    }
    emit('coords', selected.value);
  };
  document.addEventListener('keydown', (evt) => {
    keyNav(evt);
  });
  const blockStyle = (index) => {
    return index == selected.value[0]
      ? [styles.pgmMapBlock, styles.pgmSelected]
      : [styles.pgmMapBlock];
  };
</script>
<template>
  <div :class="[styles.pgmMap]" @keyDown="(evt) => keyNav(evt)">
    <div v-for="(block, bix) in props.blocks" :key="bix">
      <div :class="blockStyle(bix)">
        <div :class="styles.vert">
          <div
            v-for="(cycle, cix) in block.microCycles"
            :key="cix"
            :class="[styles.pgmMapCycle]"
          >
            <div
              v-for="(workout, wix) in cycle.workouts"
              :key="wix"
              :class="[styles.pgmMapWorkout]"
            >
              <q-btn
                icon="circle"
                :class="[
                  styles[isSelected([bix, cix, wix]) ? 'pgmMapSelected' : ''],
                ]"
                round
                @Click="setCoords([bix, cix, wix])"
                @Focus="setCoords([bix, cix, wix])"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
