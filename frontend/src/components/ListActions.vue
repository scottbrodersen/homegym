<script setup>
  /**
   * A hamburger-style drop down menu of the operations that the OrderedList class defines.
   *
   * Emits the selected operation.
   */
  import { QBtn, QMenu, QList, QItem, QIcon } from 'quasar';
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
        <q-item clickable v-close-popup dark @click="update(OrderedList.COPY)">
          <q-icon name="content_copy" size="1.65em" color="primary" />
        </q-item>
        <q-item clickable v-close-popup dark @click="update(OrderedList.PASTE)">
          <q-icon name="content_paste" size="1.65em" color="primary" />
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
