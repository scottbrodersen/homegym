<script setup>
  import { defineProps, ref } from 'vue';
  import * as styles from '../style.module.css';

  const props = defineProps({ writable: Boolean, distance: Number });
  const emit = defineEmits(['update']);

  const distance = ref(props.distance);

  const validate = (distValue) => {
    const regex = new RegExp('^[0-9](\.[0-9])?$');
    return regex.test(distValue);
  };

  const update = (newDistance) => {
    if (validate(newDistance)) {
      emit('update', parseFloat(newDistance));
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
