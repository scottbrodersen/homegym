<script setup>
  import { inject, watch } from 'vue';
  import { QExpansionItem, QInput } from 'quasar';
  import ProgramMicrocycle from './ProgramMicrocycle.vue';
  import styles from '../style.module.css';
  import { OrderedList, states } from '../modules/utils.js';
  import ListActions from './ListActions.vue';

  const state = inject('state');
  const requiredField = inject('requiredField');
  const maxField = inject('maxField');
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
  <div :class="[styles.pgmBlock]">
    <div v-if="state == states.READ_ONLY">
      <span :class="[styles.pgmBlockTitle]">{{ props.block.title }}:</span>
      {{ props.block.description }}
      <div
        :class="[styles.pgmCycles]"
        v-for="(cycle, ix) of props.block.microCycles"
        :key="ix"
      >
        <ProgramMicrocycle
          :microcycle="cycle"
          @update="(value) => updateCycles(value, ix)"
        />
      </div>
    </div>
    <div v-else>
      <div :class="[styles.horiz]">
        <div :class="[styles.pgmEditbles]">
          <q-input
            v-model="props.block.title"
            label-slot
            stack-label
            dark
            :rules="[requiredField, maxField]"
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
            :rules="[maxField]"
          />
        </div>
        <ListActions @update="update" />
      </div>
      <div :class="[styles.pgmChild]">
        <q-expansion-item
          v-if="state != states.READ_ONLY"
          default-opened
          expand-separator
          label="MicroCycles"
          :expand-icon-class="styles.pgmChildExpander"
          expand-icon="arrow_right"
          expanded-icon="arrow_drop_down"
          switch-toggle-side
          dense
        >
          <ProgramMicrocycle
            v-for="(cycle, ix) of props.block.microCycles"
            :key="ix"
            :microcycle="cycle"
            @update="(value) => updateCycles(value, ix)"
          />
        </q-expansion-item>
      </div>
    </div>
  </div>
</template>
