<script setup>
  import { onBeforeMount, ref, watch } from 'vue';
  import { QRating } from 'quasar';
  import styles from '../style.module.css';
  const props = defineProps([
    'mood',
    'energy',
    'motivation',
    'overall',
    'notes',
  ]);
  const emit = defineEmits([
    'mood',
    'energy',
    'motivation',
    'overall',
    'notes',
  ]);
  const mood = ref();
  const energy = ref();
  const motivation = ref();
  const overall = ref();
  const notes = ref();

  onBeforeMount(() => {
    mood.value = !!props.mood ? props.mood : 0;
    energy.value = !!props.energy ? props.energy : 0;
    motivation.value = !!props.motivation ? props.motivation : 0;
    overall.value = !!props.overall ? props.overall : 0;
    notes.value = !!props.notes ? props.notes : '';
  });

  watch(mood, (newMood, oldMood) => {
    if (newMood != oldMood) {
      emit('mood', newMood);
    }
  });
  watch(energy, (newEnergy, oldEnergy) => {
    if (newEnergy != oldEnergy) {
      emit('energy', newEnergy);
    }
  });
  watch(motivation, (newMotivation, oldMotivation) => {
    if (newMotivation != oldMotivation) {
      emit('motivation', newMotivation);
    }
  });
  watch(overall, (newOverall, oldOverall) => {
    if (newOverall != oldOverall) {
      emit('overall', newOverall);
    }
  });
  watch(notes, (newNotes, oldNotes) => {
    if (newNotes != oldNotes) {
      emit('notes', newNotes);
    }
  });
</script>
<template>
  <div :class="styles.metaDetails">
    <div :class="[styles.horiz, styles.maxSpacing]">
      <div>mood</div>
      <q-rating
        v-model="mood"
        size="1.5em"
        icon="sentiment_satisfied"
        color="primary"
      />
    </div>
    <div :class="[styles.horiz, styles.maxSpacing]">
      <div>motivation</div>
      <q-rating
        v-model="motivation"
        size="1.5em"
        icon="trending_up"
        color="primary"
      />
    </div>
    <div :class="[styles.horiz, styles.maxSpacing]">
      energy
      <q-rating v-model="energy" size="1.5em" icon="bolt" color="primary" />
    </div>
    <div :class="[styles.horiz, styles.maxSpacing]">
      <div>overall</div>
      <q-rating
        v-model="overall"
        size="1.5em"
        icon="thumb_up"
        color="primary"
      />
    </div>
    <div>
      <q-input v-model="notes" filled type="textarea" autogrow dark />
    </div>
  </div>
</template>
