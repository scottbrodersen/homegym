<script>
  /**
   * Displays the metadata of a workout event.
   *
   * Props:
   * An object that contains the values of each metadata item.
   *
   * Emits updated values.
   */
  export const labels = {
    overall: 'Overall Performance',
    notes: 'Notes',
  };
</script>
<script setup>
  import { onBeforeMount, ref, watch } from 'vue';
  import { QRating } from 'quasar';
  import * as styles from '../style.module.css';
  const props = defineProps({
    overall: Number,
    notes: String,
    readonly: Boolean,
  });
  const emit = defineEmits(['update']);
  const overall = ref();
  const notes = ref();

  onBeforeMount(() => {
    overall.value = props.overall ? props.overall : 0;
    notes.value = props.notes ? props.notes : '';
  });

  watch(overall, (newOverall, oldOverall) => {
    if (newOverall != oldOverall) {
      emit('update', 'overall', newOverall);
    }
  });
  watch(notes, (newNotes, oldNotes) => {
    if (newNotes != oldNotes) {
      emit('update', 'notes', newNotes);
    }
  });
</script>
<template>
  <div v-if="props.readonly" :class="[styles.horiz]">
    <div :class="styles.sibSpSmall" v-for="(value, key) in labels">
      <div v-if="props[key]" :class="styles.eventMetaItem">
        <div :class="[styles.hgBold]">{{ value }}: &nbsp;</div>
        <div>{{ props[key] }}</div>
      </div>
    </div>
  </div>
  <div v-else>
    <div :class="[styles.metaDetails]">
      <div :class="[styles.metaRating]">
        <div>Overall Performance</div>
        <q-rating
          v-model="overall"
          size="1.5em"
          icon="thumb_up"
          color="secondary"
        />
      </div>
    </div>
    <div :class="[styles.hg100wide]">
      <q-input
        v-model="notes"
        filled
        type="textarea"
        autogrow
        dark
        label="Notes"
        stack-label
      />
    </div>
  </div>
</template>
