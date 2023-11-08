<script setup>
  import { defineProps, computed } from 'vue';
  import styles from '../style.module.css';
  import BinaryRep from './BinaryRep.vue';
  const props = defineProps({ reps: Array, volumeConstraint: Number });

  const repCount = computed(() => {
    if (props.volumeConstraint === 1) {
      let count = 0;
      for (const rep of props.reps) {
        count = rep > 0 ? count + 1 : count;
      }

      return count;
    }
  });
</script>

<template>
  <div
    v-if="props.volumeConstraint === 2"
    :class="[styles.centered, styles.binarySet, styles.horiz]"
  >
    <BinaryRep
      :class="styles.sibSpSmall"
      v-for="rep in props.reps"
      :success="!!rep"
    />
  </div>

  <div v-else-if="props.volumeConstraint === 0">
    <span v-for="(rep, index) in props.reps">{{ rep }}</span>
  </div>

  <div v-else :class="[styles.repCount]">{{ repCount }}</div>
</template>
