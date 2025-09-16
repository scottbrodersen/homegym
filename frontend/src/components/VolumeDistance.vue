<script setup>
  /**
   * Displays a distance-based value for exercise volume in read-only or edit mode.
   *
   * Props:
   *  writable is True for edit mode, false for read-only
   *  distance (optional) is the value to display, in km
   *
   * Emits the distance in meters
   */
  import { defineProps, ref } from 'vue';
  import * as styles from '../style.module.css';

  const props = defineProps({ writable: Boolean, distance: Number });
  const emit = defineEmits(['update']);

  const distance = ref((props.distance / 1000).toFixed(2));

  // input is km
  const validate = (distValue) => {
    const regex = new RegExp('^[0-9](\.[0-9]{1,2})?$');
    return regex.test(distValue);
  };

  // emit meters
  const update = (newDistance) => {
    if (validate(newDistance)) {
      emit('update', parseFloat(newDistance * 1000));
    }
  };
</script>

<template>
  <div v-if="props.writable" :class="[styles.horiz]">
    <q-input
      v-model.number="distance"
      :class="[styles.timeInput]"
      focus
      select
      dark
      mask="#.#"
      label="Distance"
      stack-label
      :rules="[validate]"
      lazy-rules
      @update:model-value="update"
    />
  </div>
  <div v-else :class="[styles.distanceInput]">{{ distance }}</div>
</template>
