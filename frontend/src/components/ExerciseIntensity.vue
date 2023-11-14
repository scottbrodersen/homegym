<script setup>
  import { unitsState } from '../modules/state.js';
  import styles from '../style.module.css';
  import { ref } from 'vue';
  import { intensityTypeProps } from '../modules/utils';

  const props = defineProps({
    intensity: Number,
    type: String,
    writable: Boolean,
  });

  const emit = defineEmits(['update']);

  const intensity = ref(props.intensity);

  const mask = ref();
  const validate = ref(() => {
    return null;
  });
  const decimals = ref(1);
  const prefx = ref('');

  if (
    !!intensityTypeProps[props.type] &&
    !!intensityTypeProps[props.type].mask
  ) {
    mask.value = intensityTypeProps[props.type].mask;
  } else if (!!intensityTypeProps.default.mask) {
    mask.value = intensityTypeProps.default.mask;
  }

  if (
    !!intensityTypeProps[props.type] &&
    !!intensityTypeProps[props.type].validate
  ) {
    validate.value = intensityTypeProps[props.type].validate;
  } else if (!!intensityTypeProps.default.validate) {
    validate.value = intensityTypeProps.default.validate;
  }

  if (
    !!intensityTypeProps[props.type] &&
    intensityTypeProps[props.type].decimals != undefined
  ) {
    decimals.value = intensityTypeProps[props.type].decimals;
  } else if (!!intensityTypeProps.default.decimals) {
    decimals.value = intensityTypeProps.default.decimals;
  }

  if (
    !!intensityTypeProps[props.type] &&
    !!intensityTypeProps[props.type].prefx
  ) {
    prefx.value = intensityTypeProps[props.type].prefx;
  }

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
      :mask="mask"
      :rules="[validate]"
      lazy-rules
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
      {{ prefx }} {{ props.intensity.toFixed(decimals) }}
    </div>
    <div :class="[styles.sibSpxSmall]">{{ unitsState[props.type] }}</div>
  </div>
  <div v-else>
    <div :class="[styles.intensity, styles.sibSpxSmall]">
      {{ props.type }}
    </div>
  </div>
</template>
