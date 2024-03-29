<script setup>
  import { ref, computed, toRaw, watch } from 'vue';
  import styles from '../style.module.css';
  import DatePicker from './DatePicker.vue';
  import EventMeta from './EventMeta.vue';
  import ExerciseInstance from './ExerciseInstance.vue';
  import { eventStore, activityStore } from '../modules/state.js';
  import { QBtn } from 'quasar';
  import {
    authPrompt,
    storeEvent,
    storeEventExerciseInstances,
    ErrNotLoggedIn,
    toast,
  } from '../modules/utils.js';

  const props = defineProps({
    eventId: String,
  });

  const thisEvent = ref({});
  const thisEventActivityName = ref('');
  const activityNames = [];

  // populate state when opening an existing event
  if (!!props.eventId) {
    thisEvent.value = structuredClone(toRaw(eventStore.getByID(props.eventId)));
    if (!!thisEvent.value.activityID) {
      thisEventActivityName.value = activityStore.get(
        thisEvent.value.activityID
      ).name;
    }
  }

  for (const activity of activityStore.getAll()) {
    activityNames.push(activity.name);
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

  // Updates an event's exericse instance at a specific index.
  // An index of null adds a new instance
  // An empty updated instance removes it
  const setExerciseInstance = (index, updated) => {
    if (index == null) {
      if (!!!thisEvent.value.exInstances) {
        // initialize the exInstances object
        thisEvent.value.exInstances = {};
      }

      const newIndex = Object.keys(thisEvent.value.exInstances).length;
      const newInstance = {
        index: newIndex,
      };

      thisEvent.value.exInstances[newIndex] = newInstance;
    } else if (updated == {}) {
      delete thisEvent.value.exInstances[index];

      // normalize the indexes
      const normalized = {};
      Object.values(thisEvent.value.exInstances).forEach((exInst) => {
        normalized[exInst.index] = exInst;
      });

      thisEvent.value.exInstances = normalized;
    } else {
      thisEvent.value.exInstances[index] = updated;
    }
  };

  const updateButtonText = computed(() => {
    return !!thisEvent.value.id ? 'Update' : 'Add and continue';
  });

  const updateDateValue = (newDate) => {
    thisEvent.value.date = newDate;
  };

  const saveThisEvent = () => {
    // Use the stored date in the URL in case the date has been edited
    const url = !!thisEvent.value.id
      ? `/homegym/api/events/${eventStore.getByID(thisEvent.value.id).date}/${
          thisEvent.value.id
        }/`
      : '/homegym/api/events/';

    storeEvent(url, thisEvent.value)
      .then((responseEvent) => {
        if (!!thisEvent.value.id && thisEvent.value.id != responseEvent.id) {
          throw new Error('Event id mismatch');
        }
        thisEvent.value.id = responseEvent.id;
        if (!!eventStore.getByID(responseEvent.id)) {
          eventStore.update(thisEvent.value);
          saveExerciseInstances();
        } else {
          eventStore.add(thisEvent.value);
          // event was just created so no exercises to update
        }
        toast('Saved', 'positive');
      })
      .catch((e) => {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(saveThisEvent);
        } else {
          console.log(e);
        }
      });
  };

  const saveExerciseInstances = () => {
    storeEventExerciseInstances(
      thisEvent.value.id,
      thisEvent.value.date,
      thisEvent.value.exInstances
    )
      .then((responses) => {
        if (!!responses) {
          console.log(responses);
        }
      })
      .catch((e) => {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(saveExerciseInstances);
        } else {
          console.log(e);
        }
      });
  };
</script>

<template>
  <h1 :id="styles.event" :class="[styles.blockPadSm]">Edit Event</h1>
  <div :class="[styles.vert]">
    <div :class="[styles.eventTopRow]">
      <DatePicker
        :style="[styles.blockPadMed]"
        :date-value="thisEvent.date"
        @update="updateDateValue"
      />
      <q-select
        :class="[styles.selActivity]"
        :model-value="thisEventActivityName"
        @update:model-value="setActivity"
        :options="activityNames"
        label="Activity"
        dark
      />
    </div>
    <div :class="[styles.blockPadSm]">
      <EventMeta
        :mood="thisEvent.mood"
        :energy="thisEvent.energy"
        :motivation="thisEvent.motivation"
        :overall="thisEvent.overall"
        :notes="thisEvent.notes"
        v-show="thisEvent.activityID"
        @mood="
          (value) => {
            thisEvent.mood = value;
          }
        "
        @energy="
          (value) => {
            thisEvent.energy = value;
          }
        "
        @motivation="
          (value) => {
            thisEvent.motivation = value;
          }
        "
        @overall="
          (value) => {
            thisEvent.overall = value;
          }
        "
        @notes="
          (value) => {
            thisEvent.notes = value;
          }
        "
      />
    </div>

    <div
      :class="[styles.exInstContainer]"
      v-for="(value, index) in thisEvent.exInstances"
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

  <div :class="[styles.buttonArray]" v-show="!!thisEvent.activityID">
    <q-btn
      label="Add exercise"
      color="accent"
      text-color="dark"
      @click="setExerciseInstance(null, null)"
    />
    <q-btn
      :label="updateButtonText"
      color="accent"
      text-color="dark"
      @click="saveThisEvent"
    />
  </div>
</template>
