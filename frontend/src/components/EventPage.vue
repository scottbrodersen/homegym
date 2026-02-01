<script async setup>
  /**
   * A page for displaying, editing, and deleting an existing workout event or creating a new workout event.
   *
   * Props:
   *  eventID is the ID of the existing event to display. Do not provide a value for new events.
   *  programInstanceID is the ID of the workout in a program instance that this event fulfills
   *  blockIndex is the index of the block of the workout. 0-based.
   *  microCycleIndex is the index of the microcycle in the block. 0-based.
   *  workoutIndex is the index of the workout in the microcycle. 0-based.
   *  dayIndex is the index of the program day that the workout falls on. 0-based.
   */
  import { inject, ref, computed } from 'vue';
  import * as styles from '../style.module.css';
  import DatePicker from './DatePicker.vue';
  import EventMeta from './EventMeta.vue';
  import ExerciseInstance from './ExerciseInstance.vue';
  import {
    eventStore,
    activityStore,
    programInstanceStore,
  } from '../modules/state';
  import { QBtn, QSelect, QSpinner } from 'quasar';
  import {
    authPromptAsync,
    storeEvent,
    deleteEvent,
    ErrNotLoggedIn,
    toast,
    fetchEvents,
    openConfirmModal,
    updateProgramInstance,
  } from '../modules/utils.js';
  import { useRouter } from 'vue-router';

  const docsContext = ref(inject('docsContext'));
  docsContext.value = 'event';

  const router = useRouter();

  const props = defineProps({
    eventId: String,
    programInstanceID: String,
    blockIndex: String,
    microCycleIndex: String,
    workoutIndex: String,
    dayIndex: String,
  });

  const thisEvent = ref({});
  const thisEventActivityName = ref('');
  const activityNames = [];
  let programInstance;
  let updateProgram = false;
  let baseline = ref();

  const setBaseline = (eventID) => {
    baseline.value = JSON.stringify(eventStore.getByID(props.eventId));
  };

  // populate state when opening an existing event
  if (props.eventId) {
    if (!eventStore.getByID(props.eventId)) {
      await fetchEvents(props.eventId);
    }
    setBaseline();
    thisEvent.value = JSON.parse(baseline.value);

    if (thisEvent.value.activityID) {
      thisEventActivityName.value = activityStore.get(
        thisEvent.value.activityID,
      ).name;
    }
  } else if (props.programInstanceID && props.dayIndex) {
    // prefill the form when creating an event from a program workout

    updateProgram = true;
    baseline.value = '';
    programInstance = programInstanceStore.get(props.programInstanceID);

    thisEvent.value = { activityID: programInstance.activityID };

    thisEventActivityName.value = activityStore.get(
      thisEvent.value.activityID,
    ).name;

    thisEvent.value.exercises = {};
    const workout =
      programInstance.blocks[props.blockIndex].microCycles[
        props.microCycleIndex
      ].workouts[props.workoutIndex];

    for (let i = 0; i < workout.segments.length; i++) {
      thisEvent.value.exercises[i] = {
        index: i,
        typeID: workout.segments[i].exerciseTypeID,
        parts: [{ intensity: 0, volume: [] }],
      };
    }
  }

  for (const activity of activityStore.getAll()) {
    activityNames.push(activity.name);
    activityNames.sort();
  }

  const setActivity = (activityName) => {
    for (const activity of activityStore.getAll()) {
      if (activity.name == activityName) {
        thisEvent.value.activityID = activity.id;
        thisEventActivityName.value = activityName;
        break;
      }
    }
  };

  // Updates an event's exercise instance at a specific index.
  // An index of null adds a new instance
  // An empty updated instance removes it
  const setExerciseInstance = (index, updated) => {
    if (index == null) {
      if (!thisEvent.value.exercises) {
        // initialize the exercises object
        thisEvent.value.exercises = {};
      }

      const newIndex = Object.keys(thisEvent.value.exercises).length;
      const newInstance = {
        index: newIndex,
      };

      thisEvent.value.exercises[newIndex] = newInstance;
    } else if (updated == {}) {
      delete thisEvent.value.exercises[index];

      // normalize the indexes
      const normalized = {};
      Object.values(thisEvent.value.exercises).forEach((exInst) => {
        normalized[exInst.index] = exInst;
      });

      thisEvent.value.exercises = normalized;
    } else {
      thisEvent.value.exercises[index] = updated;
    }
  };

  const updateDateValue = (newDate) => {
    thisEvent.value.date = newDate;
  };

  const saveThisEvent = async () => {
    disableSave.value = true;
    showSpinner.value = true;

    // Use the stored date for the URL path in case the date has been edited
    // todo: move the URL construction to the storeEvent function
    const url = thisEvent.value.id
      ? `/homegym/api/events/${eventStore.getByID(thisEvent.value.id).date}/${
          thisEvent.value.id
        }/`
      : '/homegym/api/events/';
    try {
      const responseEvent = await storeEvent(url, thisEvent.value);

      thisEvent.value.id = responseEvent.id;

      if (eventStore.getByID(responseEvent.id)) {
        eventStore.update(thisEvent.value);
      } else {
        eventStore.add(thisEvent.value);
      }

      setBaseline();
      showSpinner.value = false;
      toast('Saved', 'positive');
      disableSave.value = false;

      // update the program instance if props.instanceID
      if (updateProgram) {
        programInstance.events[props.dayIndex] = thisEvent.value.id;
        await updateProgramInstance(programInstance);

        updateProgram = false;
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
        saveThisEvent();
      } else {
        console.log(e);
        toast('Error', 'negative');
      }
    }
  };

  const deleteThisEvent = () => {
    openConfirmModal('Delete event forever?').then(async (confirmed) => {
      if (confirmed === true) {
        try {
          await deleteEvent(thisEvent.value);
          toast('Saved', 'positive');
          eventStore.delete(thisEvent.value);
          await router.replace({ name: 'home' });
        } catch (e) {
          if (e instanceof ErrNotLoggedIn) {
            console.log(e.message);
            await authPromptAsync();
            deleteThisEvent(thisEvent.value);
          } else {
            console.log(e);
            toast('Error', 'negative');
          }
        }
      }
    });
  };

  // controls whether the Save button is enabled
  const changed = computed(() => {
    return baseline.value != JSON.stringify(thisEvent.value);
  });

  // use to override changed
  const disableSave = ref(false);

  const showSpinner = ref(false);

  const cancel = async () => {
    const route = { name: 'home' };
    if (changed.value) {
      openConfirmModal('Lose unsaved changes?').then(async (confirmed) => {
        if (confirmed === true) {
          await router.replace(route);
        }
      });
    } else {
      await router.replace(route);
    }
  };
</script>

<template>
  <div>
    <h1 :id="styles.event" :class="[styles.blockPadSm]">Edit Event</h1>
    <div :class="[styles.vert]">
      <div :class="[styles.eventTopRow]">
        <div>
          <DatePicker
            :style="[styles.blockPadMed]"
            :date-value="thisEvent.date"
            @update="updateDateValue"
          />
        </div>
        <div :class="[styles.activitySelect]">
          <q-select
            :class="[styles.selActivity]"
            :model-value="thisEventActivityName"
            @update:model-value="setActivity"
            :options="activityNames"
            label="Activity"
            dark
          />
        </div>
      </div>
      <div
        :class="[styles.exInstContainer]"
        v-for="(value, index) in thisEvent.exercises"
        :key="index"
      >
        <ExerciseInstance
          :exercise-instance="value"
          :activity-i-d="thisEvent.activityID"
          :writable="true"
          @update="(updated) => setExerciseInstance(index, updated)"
        />
      </div>
    </div>
    <div :class="[styles.blockPadSm]">
      <q-btn
        v-if="thisEventActivityName"
        label="Add exercise"
        color="primary"
        @click="setExerciseInstance(null, null)"
      />
    </div>
    <div :class="[styles.blockPadSm]">
      <EventMeta
        :overall="thisEvent.overall"
        :notes="thisEvent.notes"
        v-show="thisEvent.activityID"
        @update="(meta, value) => (thisEvent[meta] = value)"
      />
    </div>
    <div
      :class="[styles.buttonArray, styles.stickyBottom]"
      v-show="thisEvent.activityID"
    >
      <q-btn label="Cancel" color="accent" text-color="dark" @click="cancel" />
      <q-btn label="Delete" color="negative" dark @click="deleteThisEvent" />
      <q-btn
        label="Save"
        color="accent"
        text-color="dark"
        :disabled="!changed || disableSave"
        @click="saveThisEvent"
      />
    </div>
    <div v-show="showSpinner" :class="[styles.spinner, styles.horiz]">
      <q-spinner size="3em" />
    </div>
  </div>
</template>
