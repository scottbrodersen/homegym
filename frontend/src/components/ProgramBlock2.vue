<script setup>
  import { inject } from 'vue';
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
    <div :class="[styles.pgmBlock]">
      <div :class="[styles.pgmBlockTitle]">{{ props.block.title }}</div>
      <div>{{ props.block.description }}</div>
      <div>
        <q-btn icon="edit" color="primary" round dark @click="editBlock" />
      </div>
    </div>
  </div>
</template>
