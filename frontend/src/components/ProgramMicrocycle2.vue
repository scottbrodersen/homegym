<script setup>
  /**
   * Displays a program instance microcycle.
   * Enables users to edit the microcycle.
   *
   * Props:
   *  microcycle: A program instance microcycle object.
   */
  import * as styles from '../style.module.css';
  import { openEditValueModal } from '../modules/utils.js';

  const props = defineProps({ microcycle: Object });
  const emit = defineEmits(['update']);

  const editMicroCycle = () => {
    const values = [
      { label: 'Title', value: props.microcycle.title },
      { label: 'Description', value: props.microcycle.description },
    ];
    openEditValueModal(values).then((edited) => {
      if (edited) {
        const editedBlock = { title: edited[0], description: edited[1] };
        emit('update', editedBlock);
      }
    });
  };
</script>
<template>
  <div>
    <div :class="[styles.pgmMicrocycle, styles.horiz]">
      <div :class="[styles.pgmMicrocycleTitle]">
        {{
          props.microcycle.title ? props.microcycle.title : '~~needs a title~~'
        }}
      </div>
      <div>
        {{ props.microcycle.description }}
      </div>
      <div>
        <q-btn icon="edit" color="primary" round dark @click="editMicroCycle" />
      </div>
    </div>
  </div>
</template>
