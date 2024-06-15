<script async setup>
  import { programsStore, programInstanceStore } from '../modules/state';
  import { onBeforeMount, ref, watch } from 'vue';
  import * as styles from '../style.module.css';
  import {
    authPromptAsync,
    fetchPrograms,
    fetchProgramInstances,
    ErrNotLoggedIn,
  } from '../modules/utils';
  const props = defineProps({ activityID: String, programID: String });
  const emit = defineEmits(['selected']);

  const listItems = ref([]);
  const selectedObj = ref();

  // populates program list for selected activity
  const populatePrograms = (activityID) => {
    programsStore.getByActivity(activityID).forEach((program) => {
      listItems.value.push(program);
      const instances = programInstanceStore.getByProgram(program.id);
      if (instances) {
        listItems.value.push(...instances);
      }
    });
  };

  const getPrograms = async (activityID) => {
    if (activityID && programsStore.getByActivity(activityID) === undefined) {
      try {
        await fetchPrograms(activityID);
        for (const pgm of programsStore.getByActivity(activityID)) {
          await fetchProgramInstances(pgm.id, activityID);
        }
      } catch (e) {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          await authPromptAsync();
          getPrograms(activityID);
        } else {
          console.log(e.message);
        }
      }
    }
    populatePrograms(activityID);
  };

  watch(
    () => props.activityID,
    async (newID) => {
      if (newID) {
        await getPrograms(newID);
      }
    }
  );

  const addNew = (programID, done) => {
    // update selector
    done(programsStore.get(activity.value.id, programID), 'add-unique');
  };

  watch(
    () => selectedObj.value,
    (newSelected) => {
      const idObj = newSelected.programID
        ? { programInstanceID: selectedObj.value.id }
        : { programID: selectedObj.value.id };
      emit('selected', idObj);
    }
  );

  onBeforeMount(async () => {
    if (props.activityID) {
      await getPrograms(props.activityID);
    }

    if (props.programID) {
      for (let i = 0; i < listItems.value.length; i++) {
        if (props.programID == listItems.value[i].id) {
          selectedObj.value = listItems.value[i];
          break;
        }
      }
    }
  });
</script>
<template>
  <q-select
    id="program"
    label="Program"
    stack-label
    v-model="selectedObj"
    :options="listItems"
    option-label="title"
    option-value="id"
    @new-value="addNew"
    dark
    :class="[styles.selProgram]"
  >
    <template v-slot:option="scope">
      <q-item v-bind="scope.itemProps">
        <div v-if="scope.opt.programID" :class="[styles.pgmInstItem]">
          {{ scope.opt.title }}
        </div>
        <div v-else :class="[styles.pgmItem]">{{ scope.opt.title }}</div>
      </q-item>
    </template>
    <template v-slot:selected>
      <div v-if="selectedObj">
        {{ selectedObj.title }}
      </div>
    </template>
  </q-select>
</template>
