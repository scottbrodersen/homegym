<script setup>
  import { unitsState } from '../modules/state.js';
  import * as styles from '../style.module.css';
  import { ref } from 'vue';
  import { intensityProps } from '../modules/utils';
  import { QInput } from 'quasar';

  const props = defineProps({
    intensity: Number,
    type: String,
    writable: Boolean,
  });

  const formatProps = intensityProps(props.type);

  const emit = defineEmits(['update']);

  const intensity = ref(formatProps.format(props.intensity));

  const update = (value) => {
    emit('update', formatProps.value(value));
  };
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
      :rules="[(value) => formatProps.validate(value) || 'fix it']"
      lazy-rules
      filled
      dense
      dark
      v-focus
      v-select
      no-error-icon
      hide-bottom-space
      @update:model-value="update"
    />
  </div>
  <div v-else :class="[styles.horiz]">
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ formatProps.prefix }}
      {{ formatProps.format(props.intensity) }}
    </div>
    <div v-if="props.type != 'bodyweight'" :class="[styles.sibSpxSmall]">
      {{ unitsState[props.type] }}
    </div>
  </div>
</template>
