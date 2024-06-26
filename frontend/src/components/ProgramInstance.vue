<script setup>
  import { inject, provide, ref, watch } from 'vue';
  import ProgramBlock2 from './ProgramBlock2.vue';
  import { programInstanceStore, programsStore } from './../modules/state';
  import { updateProgramInstance } from './../modules/utils';
  import {
    QBtn,
    QDialog,
    QIcon,
    QItem,
    QItemSection,
    QMenu,
    QList,
    QOptionGroup,
  } from 'quasar';
  import * as styles from '../style.module.css';
  import {
    authPromptAsync,
    deepToRaw,
    ErrNotLoggedIn,
    states,
    toast,
  } from '../modules/utils';
  import * as programUtils from '../modules/programUtils';
  import * as dateUtils from '../modules/dateUtils';
  import ProgramCalendar from './ProgramCalendar.vue';
  import ProgramMicrocycle2 from './ProgramMicrocycle2.vue';
  import ProgramWorkout2 from './ProgramWorkout2.vue';
  import ProgramWorkoutSegment from './ProgramWorkoutSegment.vue';
  import * as utils from '../modules/utils';

  const props = defineProps({ activityID: String, instanceID: String });
  const emit = defineEmits(['done']);

  const instance = ref();
  const coords = ref();
  const linkEventDialog = ref({ show: false, eventID: '', events: undefined });

  let today;

  let baseline = ''; // use to detect diff
  const changed = ref();
  const valid = ref(true);
  const programTitle = ref();

  const { editInstanceTitle, toggleInstanceTitle } =
    inject('editInstanceTitle');
  const { state, setState } = inject('state');

  provide('activity', props.activityID);

  const init = () => {
    if (!props.instanceID) {
      baseline = '';
      instance.value = {};
    } else {
      instance.value = deepToRaw(programInstanceStore.get(props.instanceID));
      baseline = JSON.stringify(instance.value);

      programTitle.value = instance.value.programID
        ? programsStore.get(props.activityID, instance.value.programID).title
        : '';

      today = programUtils.getTodayIndex(instance.value);
    }
  };

  // Re-initialize when a different instance is selected
  watch(
    () => props.instanceID,
    (newID) => {
      init();
    }
  );
  init();

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

  const saveInstance = async () => {
    try {
      const id = await updateProgramInstance(instance.value);
      toast('Saved', 'positive');
    } catch (e) {
      console.log(e.message);

      if (e instanceof ErrNotLoggedIn) {
        await authPromptAsync();
        saveInstance();
      } else {
        toast('Error', 'negative');
      }
    } finally {
      setState(states.READ_ONLY);
    }
  };

  const updateMicroCycle = async (
    microCycleProps,
    blockIndex,
    microCycleIndex
  ) => {
    instance.value.blocks[blockIndex].microCycles[microCycleIndex].title =
      microCycleProps.title;
    instance.value.blocks[blockIndex].microCycles[microCycleIndex].description =
      microCycleProps.description;
    await saveInstance();
  };

  const updateBlock = async (blockProps, index) => {
    instance.value.blocks[index].title = blockProps.title;
    instance.value.blocks[index].description = blockProps.description;
    await saveInstance();
  };

  const setCoords = (dayIndex) => {
    if (dayIndex != -1) {
      coords.value = programUtils.getWorkoutCoords(instance.value, dayIndex);
    }
  };

  watch(
    () => coords.value,
    (newCoords) => {
      const evtEl = document.getElementById(
        `workout${newCoords[0]}-${newCoords[1]}-${newCoords[2]}`
      );
      if (evtEl) {
        evtEl.scrollIntoView({
          behavior: 'smooth',
          block: 'start',
          inline: 'center',
        });
      }
    }
  );

  const isFuture = (coords) => {
    const day = programUtils.getDayIndex(instance.value, coords);
    return day > today ? true : false;
  };

  const hasEvent = (coords) => {
    const day = programUtils.getDayIndex(instance.value, coords);

    return instance.value.events[day] ? true : false;
  };

  const isRestDay = (coords) => {
    if (
      instance.value.blocks[coords[0]].microCycles[coords[1]].workouts[
        coords[2]
      ].restDay
    ) {
      return true;
    }
    return false;
  };

  const linkEvent = async (coords) => {
    if (coords === true) {
      instance.value.events[linkEventDialog.value.index] =
        linkEventDialog.value.eventID;
      await saveInstance();
    } else if (typeof coords === 'object') {
      const events = await programUtils.getEventsOnWorkoutDay(
        instance.value,
        coords
      );
      // set up dialog options
      if (events && events.length > 0) {
        linkEventDialog.value.events = [];
        events.forEach((evt) => {
          const eventTime = dateUtils.formatTime(
            dateUtils.dateFromSeconds(evt.date)
          );
          linkEventDialog.value.events.push({
            label: eventTime,
            value: evt.id,
          });
        });
      } else {
        linkEventDialog.value.events = undefined;
      }
      const workoutIndex = programUtils.getDayIndex(instance.value, coords);
      linkEventDialog.value.show = true;
      linkEventDialog.value.eventID =
        events.length > 0 ? events[0].id : undefined;
      linkEventDialog.value.index = workoutIndex;
    }
  };

  // Open modal to edit the instance title
  watch(
    () => editInstanceTitle.value,
    (newValue) => {
      if (newValue === true) {
        utils
          .openEditValueModal([
            {
              label: 'Instance Title',
              value: instance.value.title,
            },
          ])
          .then(async (newValue) => {
            if (newValue) {
              console.log(newValue);
              instance.value.title = newValue[0];
              await saveInstance();
            }
            toggleInstanceTitle();
          });
      }
    }
  );

  const editWorkout = async (coords) => {
    programUtils
      .newWorkoutModal(instance.value, coords)
      .then(async (workout) => {
        console.log('workout value: ' + workout);
        if (workout) {
          instance.value.blocks[coords[0]].microCycles[coords[1]].workouts[
            coords[2]
          ] = workout;
          try {
            await utils.updateProgramInstance(instance.value);
            toast('Saved', 'positive');
          } catch (e) {
            if (e instanceof ErrNotLoggedIn) {
              console.log(e.message);
              await authPromptAsync();
              editWorkout(coords);
            } else {
              console.log(e);
              toast('Error', 'negative');
            }
          }
        }
      });
  };
</script>
<template>
  <div v-if="instance" :class="[styles.pgmInstance]">
    <ProgramCalendar
      :instance="instance"
      :coords="coords"
      @dayIndex="setCoords"
      :class="[styles.centered]"
    />
    <div :class="[styles.instBase]">
      Base Program:
      {{ programTitle }}
    </div>
    <div :class="[styles.instInfo]">
      <ProgramBlock2
        v-if="coords"
        :block="instance.blocks[coords[0]]"
        @update="(updated) => updateBlock(updated, coords[0])"
      />
      <ProgramMicrocycle2
        v-if="coords"
        :microcycle="instance.blocks[coords[0]].microCycles[coords[1]]"
        @update="(updated) => updateMicroCycle(updated, coords[0], coords[1])"
      />
      <div v-else>The program was performed in the past. Select a date.</div>
    </div>

    <div id="inst-wrap" v-if="coords">
      <div
        v-for="(workout, wix) of instance.blocks[coords[0]].microCycles[
          coords[1]
        ].workouts"
        :key="wix"
        :class="
          coords[2] == wix
            ? [styles.instWorkout, styles.evtHighlight]
            : [styles.instWorkout]
        "
        @click="
          () => {
            coords[2] = wix;
          }
        "
      >
        <div>
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
            <div>
              <q-icon
                v-if="isFuture([coords[0], coords[1], wix])"
                name="schedule"
                color="blue"
              />
              <q-btn
                v-else-if="hasEvent([coords[0], coords[1], wix])"
                icon="check_circle"
                color="green"
                :to="{
                  name: 'home',
                  query: {
                    event: programUtils.getDayIndex(instance, [
                      coords[0],
                      coords[1],
                      wix,
                    ]),
                  },
                }"
                round
                dense
              />
              <q-icon
                v-else-if="!isRestDay([coords[0], coords[1], wix])"
                name="cancel"
                color="red"
              />
            </div>
            <div :class="[styles.evtMenu]">
              <q-btn
                v-show="coords[2] == wix"
                icon="menu"
                :class="[styles.hgHamburger]"
              >
                <q-menu>
                  <q-list>
                    <q-item
                      clickable
                      v-close-popup
                      dark
                      @click="editWorkout([coords[0], coords[1], wix])"
                    >
                      Edit
                    </q-item>
                    <q-item
                      v-if="
                        !isRestDay([coords[0], coords[1], wix]) &&
                        hasEvent([coords[0], coords[1], wix])
                      "
                      clickable
                      v-close-popup
                      dark
                    >
                      <q-item-section>Unlink</q-item-section>
                    </q-item>
                    <q-item
                      v-if="
                        !isRestDay([coords[0], coords[1], wix]) &&
                        !isFuture([coords[0], coords[1], wix]) &&
                        !hasEvent([coords[0], coords[1], wix])
                      "
                      clickable
                      v-close-popup
                      dark
                      @click="linkEvent([coords[0], coords[1], wix])"
                    >
                      <q-item-section>Link</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-btn>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <q-dialog v-model="linkEventDialog.show" persistent>
    <q-card style="min-width: 350px" bordered dark>
      <q-card-section>Link an event with the planned workout:</q-card-section>
      <q-card-section>
        <div v-if="linkEventDialog.eventID">
          <q-option-group
            v-model="linkEventDialog.eventID"
            :options="linkEventDialog.events"
            dark
          />
        </div>
        <div v-else>There is no event on this day to link to.</div>
      </q-card-section>

      <q-card-actions
        v-if="linkEventDialog.eventID"
        align="right"
        class="text-primary"
      >
        <q-btn
          flat
          label="No"
          @Click="linkEventDialog = { show: false, eventID: '' }"
          v-close-popup
          dark
        />
        <q-btn flat label="OK" @Click="linkEvent(true)" v-close-popup dark />
      </q-card-actions>
      <q-card-actions v-else align="right" class="text-primary">
        <q-btn flat label="OK" v-close-popup dark />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
