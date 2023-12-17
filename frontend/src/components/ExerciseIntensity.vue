<script setup>
  import { unitsState } from '../modules/state.js';
  import styles from '../style.module.css';
  import { ref } from 'vue';
  import { intensityProps } from '../modules/utils';

  const props = defineProps({
    intensity: Number,
    type: String,
    writable: Boolean,
  });

  const formatProps = ref(intensityProps(props.type));

  const formatIntensity = (value) => {
    return value.toFixed(formatProps.value.decimals);
  };

  const emit = defineEmits(['update']);

  const intensity = ref(formatIntensity(props.intensity));

  // if this is a bodyweight execise, set intensity to 1
  if (props.type == 'bodyweight') {
    emit('update', 1);
  }
</script>
<template>
  <div v-if="props.writable && props.type != 'bodyweight'">
    <q-input
      :class="[styles.inputIntensity]"
      :label="props.type"
      stack-label
      v-model="intensity"
      :suffix="unitsState[props.type]"
      :mask="formatProps.mask"
      :rules="[formatProps.validate]"
      lazy-rules
      filled
      dense
      dark
      v-focus
      v-select
      no-error-icon
      hide-bottom-space
      @update:model-value="emit('update', intensity)"
    />
  </div>
  <div v-else-if="props.type != 'bodyweight'" :class="[styles.horiz]">
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ formatProps.prefix }}
      {{ props.intensity.toFixed(formatProps.decimals) }}
    </div>
    <div :class="[styles.sibSpxSmall]">{{ unitsState[props.type] }}</div>
  </div>
  <div v-else>
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ props.type }}
    </div>
  </div>
</template>
