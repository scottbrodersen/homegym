<script setup>
  /**
   * Displays or edits the properties of a program block.
   *
   * Props:
   *  block is the block object with properties title and description.
   *
   * Injected:
   *  state indicates whether to present the block in read-only or edit mode.
   */
  import { inject } from 'vue';
  import { QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import { states } from '../modules/utils.js';
  import * as programUtils from '../modules/programUtils';

  const { state } = inject('state');
  const props = defineProps({ block: Object });
</script>
<template>
  <div :class="[styles.pgmBlock]">
    <div v-if="state == states.READ_ONLY">
      <div>{{ props.block.title }}</div>
      <div>{{ props.block.description }}</div>
    </div>
    <div v-else>
      <div :class="[styles.horiz]">
        <div>
          <q-input
            v-model="props.block.title"
            label-slot
            stack-label
            dark
            :rules="[
              programUtils.requiredFieldValidator,
              programUtils.maxFieldValidator,
            ]"
          >
            <template v-slot:label>
              <div :class="[styles.pgmBlockLabel]">Block Title</div>
            </template></q-input
          >
          <q-input
            v-model="props.block.description"
            label="Description"
            stack-label
            dark
            @focus="(event) => console.log(event)"
            :rules="[programUtils.maxFieldValidator]"
          />
        </div>
      </div>
    </div>
  </div>
</template>
