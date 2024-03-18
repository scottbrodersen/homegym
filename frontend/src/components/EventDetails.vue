<script setup>
  import EventExercises from './EventExercises.vue';
  import styles from '../style.module.css';
  import EventMeta from './EventMeta.vue';
  const props = defineProps({ rowProps: Object });
</script>
<template>
  <q-tr>
    <q-td :class="[styles.expanded, styles.blockPadSm]" colspan="5">
      <q-btn
        color="primary"
        round
        :to="{ name: 'event', params: { eventId: props.rowProps.key } }"
        icon="edit"
      />
      <EventMeta
        :mood="props.rowProps.row.mood"
        :energy="props.rowProps.row.energy"
        :motivation="props.rowProps.row.motivation"
        :readonly="true"
      />
      <Suspense timeout="0">
        <EventExercises :event-id="props.rowProps.row.id" />
        <template #fallback> Loading... </template>
      </Suspense>
      <EventMeta
        :overall="props.rowProps.row.overall"
        :notes="props.rowProps.row.notes"
        :readonly="true"
      />
    </q-td>
  </q-tr>
</template>
