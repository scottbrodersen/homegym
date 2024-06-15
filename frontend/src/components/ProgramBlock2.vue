<script setup>
  import { inject, watch } from 'vue';
  import { QExpansionItem, QInput } from 'quasar';
  import * as styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import ListActions from './ListActions.vue';
  import * as programUtils from '../modules/programUtils';

  const { state } = inject('state');
  const props = defineProps({ block: Object });
  const emit = defineEmits(['update']);

  if (!props.block.microCycles) {
    props.block.microCycles = [{}];
  }

  let cycles = new OrderedList(props.block.microCycles);

  watch(
    () => {
      return props.block.microCycles;
    },
    () => {
      cycles = new OrderedList(props.block.microCycles);

      if (!props.block.microCycles) {
        props.block.microCycles = [{}];
      }
    }
  );

  // emit the action from the ListActions buttons
  const update = (action) => {
    emit('update', action);
  };

  // perform the action emitted from a microcycle component
  const updateCycles = (action, index) => {
    cycles.update(action, index);
  };
</script>
<template>
  <div>
    <div v-show="state == states.READ_ONLY" :class="[styles.pgmBlock]">
      <div :class="[styles.pgmBlockTitle]">{{ props.block.title }}</div>
      <div>{{ props.block.description }}</div>
    </div>
    <div v-show="state == states.EDIT">
      <div :class="[styles.horiz]">
        <div :class="[styles.pgmEditbles]">
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
        <ListActions @update="update" />
      </div>
    </div>
  </div>
</template>
