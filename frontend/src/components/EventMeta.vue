<script>
  /**
   * Displays the metadata of a workout event.
   *
   * Props:
   * An object that contains the values of each metadata item.
   *
   * Emits updated values of the values.
   */
  export const labels = {
    mood: 'Mood',
    energy: 'Energy',
    motivation: 'Motivation',
    overall: 'Overall Performance',
    notes: 'Notes',
  };
</script>
<script setup>
  import { onBeforeMount, ref, watch } from 'vue';
  import { QRating } from 'quasar';
  import * as styles from '../style.module.css';
  const props = defineProps({
    mood: Number,
    energy: Number,
    motivation: Number,
    overall: Number,
    notes: String,
    readonly: Boolean,
  });
  const emit = defineEmits(['update']);
  const mood = ref();
  const energy = ref();
  const motivation = ref();
  const overall = ref();
  const notes = ref();

  onBeforeMount(() => {
    mood.value = props.mood ? props.mood : 0;
    energy.value = props.energy ? props.energy : 0;
    motivation.value = props.motivation ? props.motivation : 0;
    overall.value = props.overall ? props.overall : 0;
    notes.value = props.notes ? props.notes : '';
  });

  watch(mood, (newMood, oldMood) => {
    if (newMood != oldMood) {
      emit('update', 'mood', newMood);
    }
  });
  watch(energy, (newEnergy, oldEnergy) => {
    if (newEnergy != oldEnergy) {
      emit('update', 'energy', newEnergy);
    }
  });
  watch(motivation, (newMotivation, oldMotivation) => {
    if (newMotivation != oldMotivation) {
      emit('update', 'motivation', newMotivation);
    }
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
      <div v-if="props[key]">
        <span :class="[styles.hgBold]">{{ value }}</span
        >: <span>{{ props[key] }}</span>
      </div>
    </div>
  </div>
  <div v-else>
    <div :class="[styles.metaDetails]">
      <div :class="[styles.metaRating]">
        <div>mood</div>
        <q-rating
          v-model="mood"
          size="1.5em"
          icon="sentiment_satisfied"
          color="secondary"
        />
      </div>
      <div :class="[styles.metaRating]">
        <div>motivation</div>
        <q-rating
          v-model="motivation"
          size="1.5em"
          icon="trending_up"
          color="secondary"
        />
      </div>
      <div :class="[styles.metaRating]">
        <div>energy</div>
        <q-rating v-model="energy" size="1.5em" icon="bolt" color="secondary" />
      </div>
      <div :class="[styles.metaRating]">
        <div>overall</div>
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
