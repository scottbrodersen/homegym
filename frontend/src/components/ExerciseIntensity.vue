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

  // if this is a bodyweight execise, set intensity to 1
  if (props.type == 'bodyweight') {
    emit('update', 1);
  }
</script>
<template>
  <div v-if="props.writable && props.type != 'bodyweight'">
    <q-input
      v-show="props.type != 'bodyweight'"
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
  <div v-else-if="props.type != 'bodyweight'" :class="[styles.horiz]">
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ props.intensity.toFixed(1) }}
    </div>
    <div :class="[styles.sibSpxSmall]">{{ unitsState[props.type] }}</div>
  </div>
  <div v-else>
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ props.type }}
    </div>
  </div>
</template>
