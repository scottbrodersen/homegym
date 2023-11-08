<script setup>
  import styles from '../style.module.css';

  const props = defineProps({
    meta: Object,
  });

  const labels = {
    mood: 'Mood',
    energy: 'Energy',
    motivation: 'Motivation',
    overall: 'Overall Performance',
    notes: 'Notes',
  };
</script>

<template>
  <div :class="[styles.horiz]">
    <div :class="styles.sibSpSmall" v-for="(m, i) in meta">
      <span :class="[styles.hgBold]">{{ labels[i] }}</span
      >: <span>{{ m }}</span>
    </div>
  </div>
</template>
