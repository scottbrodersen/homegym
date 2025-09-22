<script setup>
  /**
   * Displays the details of a workout event.
   * Enables editing the event.
   *
   * Props:
   *  eventId is the ID of the event to display.
   */
  import ExerciseInstance from './ExerciseInstance.vue';
  import * as styles from '../style.module.css';
  import EventMeta from './EventMeta.vue';
  import { eventStore } from '../modules/state';
  import { QBtn, QTr, QTd } from 'quasar';

  const props = defineProps({ eventId: String });
  const evt = eventStore.getByID(props.eventId);
</script>
<template>
  <q-tr>
    <q-td :class="[styles.expanded, styles.blockPadSm]" colspan="5">
      <q-btn
        color="primary"
        round
        :to="{ name: 'event', params: { eventId: props.eventId } }"
        icon="edit"
      />
      <div :class="[styles.blockPadXSm]">
        <ExerciseInstance
          v-for="(value, key) in evt.exercises"
          :exercise-instance="value"
          :activity-i-d="evt.activityID"
          :writable="false"
        />
      </div>
      <EventMeta :overall="evt.overall" :notes="evt.notes" :readonly="true" />
    </q-td>
  </q-tr>
</template>
