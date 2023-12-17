<script setup>
  import { ref } from 'vue';
  import { QBtn, QExpansionItem } from 'quasar';
  import styles from '../style.module.css';
  import { OrderedList } from '../modules/utils';

  const emit = defineEmits(['update']);

  const expanded = ref(false);

  const update = (value) => {
    emit('update', value);
    expanded.value = false;
  };
</script>
<template>
  <div>
    <q-expansion-item
      v-model="expanded"
      :class="[styles.listActions]"
      hide-expand-icon
    >
      <template v-slot:header>
        <div :class="[styles.listActionsHeader]">
          <q-icon name="menu" size="1.65em" />
        </div>
      </template>
      <div :class="styles.listActionsList">
        <div>
          <q-btn
            round
            icon="add"
            color="primary"
            size="0.65em"
            @Click="update(OrderedList.ADD)"
          />
        </div>
        <div>
          <q-btn
            round
            icon="move_up"
            color="primary"
            size="0.65em"
            @Click="update(OrderedList.MOVEBACK)"
          />
        </div>
        <div>
          <q-btn
            round
            icon="move_down"
            color="primary"
            size="0.65em"
            @Click="update(OrderedList.MOVEFWD)"
          />
        </div>
        <div>
          <q-btn
            round
            icon="delete"
            color="primary"
            size="0.65em"
            @Click="update(OrderedList.DELETE)"
          />
        </div>
      </div>
    </q-expansion-item>
  </div>
</template>
