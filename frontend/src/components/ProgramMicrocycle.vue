<script setup>
/**
 * Displays a program microcycle in read-only or edit mode.
 *
 * Props:
 *  microcycle: A program microcycle object.
 *
 * Injected:
 *  state indicates whether to present the block in read-only or edit mode.
 */
  import { inject, ref } from 'vue';
  import * as styles from '../style.module.css';
  import {  states } from '../modules/utils.js';
  import {  QInput } from 'quasar';
  import * as programUtils from '../modules/programUtils';
  import * as utils from '../modules/utils';

  const {state} = inject('state');
  const props = defineProps({ microcycle: Object });
  const emit = defineEmits(['update']);

  const rawCycle = ref(utils.deepToRaw(props.microcycle));

  const updateCycle = () => {
    emit('update', rawCycle.value);
  }

</script>
<template>
  <div :class="[styles.pgmMicrocycle]">
    <div v-if="state == states.READ_ONLY">
      <div :class="[styles.horiz]">
        <div :class="[styles.hgBold, styles.sibSpxSmall]">
          {{ props.microcycle.title ? props.microcycle.title : '<no microcycle title>' }}
        </div>
      </div>
        <div :class="[styles.sibSpxSmall]">
          {{ props.microcycle.description }}
        </div>
    </div>
    <div v-else>
      <div :class="[styles.horiz]">
        <div>
          <q-input
            v-model="rawCycle.title"
            label="Microcycle Title"
            stack-label
            dark
            :rules="[programUtils.requiredFieldValidator, programUtils.maxFieldValidator,]"
            @update:model-value="() => updateCycle()"/>
          <q-input v-model.number="rawCycle.span" type="number" label="Days" stack-label dense dark :rules="[
              programUtils.requiredFieldValidator,
              programUtils.maxFieldValidator,
            ]"
            @update:model-value="() => updateCycle()"/>
          <q-input
            v-model="rawCycle.description"
            label="Description"
            stack-label
            dark
            :rules="[programUtils.maxFieldValidator]"
            @update:model-value="() => updateCycle()"/>
        </div>
      </div>
    </div>
  </div>
</template>
