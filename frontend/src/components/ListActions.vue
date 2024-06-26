<script setup>
  import { ref } from 'vue';
  import { QBtn } from 'quasar';
  import * as styles from '../style.module.css';
  import { OrderedList } from '../modules/utils';

  const emit = defineEmits(['update']);

  const update = (value) => {
    emit('update', value);
  };
</script>
<template>
  <q-btn icon="menu" :class="[styles.hgHamburger]">
    <q-menu>
      <q-list>
        <q-item clickable v-close-popup dark @click="update(OrderedList.ADD)">
          <q-icon name="add" size="1.65em" color="primary" />
        </q-item>
        <q-item
          clickable
          v-close-popup
          dark
          @click="update(OrderedList.MOVEBACK)"
        >
          <q-icon name="move_up" size="1.65em" color="primary" />
        </q-item>
        <q-item
          clickable
          v-close-popup
          dark
          @click="update(OrderedList.MOVEFWD)"
        >
          <q-icon name="move_down" size="1.65em" color="primary" />
        </q-item>
        <q-item
          clickable
          v-close-popup
          dark
          @click="update(OrderedList.DELETE)"
        >
          <q-icon name="delete" size="1.65em" color="primary" />
        </q-item>
      </q-list>
    </q-menu>
  </q-btn>
</template>
