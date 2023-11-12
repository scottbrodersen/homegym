<script setup>
  import { unitsState } from '../modules/state.js';
  import styles from '../style.module.css';
  import { ref } from 'vue';
  const props = defineProps({
    intensity: Number,
    type: String,
    writable: Boolean,
  });
  const emit = defineEmits(['update']);
  const intensity = ref(props.intensity);
</script>
<template>
  <div v-if="props.writable">
    <q-input
      :class="[styles.inputIntensity]"
      :label="props.type"
      stack-label
      v-model.number="intensity"
      type="number"
      :suffix="unitsState[props.type]"
      filled
      dense
      dark
      v-focus
      v-select
      @update:model-value="emit('update', intensity)"
    />
  </div>
  <div v-else :class="[styles.horiz]">
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ props.intensity.toFixed(1) }}
    </div>
    <div :class="[styles.sibSpxSmall]">{{ unitsState[props.type] }}</div>
  </div>
</template>
