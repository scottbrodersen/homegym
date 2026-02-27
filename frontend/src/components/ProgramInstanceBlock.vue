<script setup>
  /**
   * Displays the properties of a block of a program instance.
   * Enables editing of the properties.
   *
   * Props:
   *  block is the block object with properties title and description.
   *
   * Emits the edited block.
   */
  import { QBtn } from 'quasar';
  import * as styles from '../style.module.css';
  import { openEditValueModal } from '../modules/utils';

  const props = defineProps({ block: Object });
  const emit = defineEmits(['update']);
  const editBlock = () => {
    const values = [
      { label: 'Title', value: props.block.title },
      { label: 'Description', value: props.block.description },
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
    <div :class="[styles.pgmBlock, styles.horiz]">
      <div :class="[styles.pgmBlockTitle]">{{ props.block.title }}</div>
      <div>{{ props.block.description }}</div>
      <div>
        <q-btn icon="edit" color="primary" round dark @click="editBlock" />
      </div>
    </div>
  </div>
</template>
