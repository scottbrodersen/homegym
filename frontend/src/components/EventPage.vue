<script async setup>
  import { ref, computed } from 'vue';
  import styles from '../style.module.css';
  import DatePicker from './DatePicker.vue';
  import EventMeta from './EventMeta.vue';
  import ExerciseInstance from './ExerciseInstance.vue';
  import { eventStore, activityStore } from '../modules/state';
  import { fetchEventPage, openConfirmModal } from '../modules/utils';
  import { QBtn, QSelect } from 'quasar';
  import {
    authPrompt,
    storeEvent,
    storeEventExerciseInstances,
    ErrNotLoggedIn,
    toast,
  } from '../modules/utils.js';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  const props = defineProps({
    eventId: String,
  });

  const thisEvent = ref({});
  const thisEventActivityName = ref('');
  const activityNames = [];
  let baseline = ref();

  const setBaseline = (eventID) => {
    baseline.value = JSON.stringify(eventStore.getByID(props.eventId));
  };

  // populate state when opening an existing event
  if (props.eventId) {
    if (!eventStore.getByID(props.eventId)) {
      await fetchEventPage(props.eventId);
    }
    setBaseline();
    thisEvent.value = JSON.parse(baseline.value);

    if (thisEvent.value.activityID) {
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

  // Updates an event's exercise instance at a specific index.
  // An index of null adds a new instance
  // An empty updated instance removes it
  const setExerciseInstance = (index, updated) => {
    if (index == null) {
      if (!!!thisEvent.value.exercises) {
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

  const updateButtonText = computed(() => {
    return !!thisEvent.value.id ? 'Update' : 'Add and continue';
  });

  const updateDateValue = (newDate) => {
    thisEvent.value.date = newDate;
  };

  const saveThisEvent = () => {
    // Use the stored date for the URL path in case the date has been edited
    const url = thisEvent.value.id
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
        setBaseline();
        //baseline = JSON.stringify(thisEvent.value);
        //changed.value = false;
        toast('Saved', 'positive');
      })
      .catch((e) => {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(saveThisEvent);
        } else {
          console.log(e);
          toast('Error', 'negative');
        }
      });
  };

  const saveExerciseInstances = () => {
    storeEventExerciseInstances(
      thisEvent.value.id,
      thisEvent.value.date,
      thisEvent.value.exercises
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
  const changed = computed(() => {
    return baseline.value != JSON.stringify(thisEvent.value);
  });

  const cancel = async () => {
    const route = { name: 'home' };
    if (changed.value) {
      openConfirmModal('Lose unsaved changes?', route, router);
    } else {
      await router.replace(route);
    }
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
        @update="meta, (value) => (thisEvent[meta] = value)"
      />
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
  <div>
    <q-btn
      label="Add exercise"
      color="primary"
      @click="setExerciseInstance(null, null)"
    />
  </div>
  <div :class="[styles.buttonArray]" v-show="!!thisEvent.activityID">
    <q-btn label="Cancel" color="accent" text-color="dark" @click="cancel" />
    <q-btn
      :label="updateButtonText"
      color="accent"
      text-color="dark"
      :disabled="!changed"
      @click="saveThisEvent"
    />
  </div>
</template>
