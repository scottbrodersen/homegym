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
  import { QItem, QSelect } from 'quasar';

  const props = defineProps({
    activityID: String,
    programID: String,
    hideCompleted: Boolean,
  });
  const emit = defineEmits(['selected']);

  const listItems = ref([]);
  const selectedObj = ref();

  // populates program list for selected activity
  const populatePrograms = (activityID) => {
    listItems.value = [];
    const programs = programsStore.getByActivity(activityID);
    if (programs) {
      programs.forEach((program) => {
        listItems.value.push(program);
        const instances = programInstanceStore.getByProgram(program.id);
        if (instances) {
          for (const instance of instances) {
            if (props.hideCompleted && isActive(activityID, instance.id)) {
              listItems.value.push(instance);
            } else if (!props.hideCompleted) {
              listItems.value.push(instance);
            }
          }
        }
      });
    }
  };

  const getPrograms = async (activityID) => {
    if (activityID && programsStore.getByActivity(activityID) === undefined) {
      try {
        await fetchPrograms(activityID);
        const programs = programsStore.getByActivity(activityID);
        if (programs) {
          for (const pgm of programs) {
            await fetchProgramInstances(pgm.id, activityID);
          }
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

  const isActive = (activityID, instanceID) => {
    if (!programInstanceStore.activeInstances.has(activityID)) {
      return false;
    }
    const activeInstances = programInstanceStore.getActive(activityID);
    for (const instance of activeInstances) {
      if (instance.id === instanceID) {
        return true;
      }
    }
    return false;
  };

  const isCurrent = (activityID, instanceID) => {
    const current = programInstanceStore.getCurrent(activityID);
    if (current != null && current.id == instanceID) {
      return true;
    }
    return false;
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
  watch(
    () => props.hideCompleted,
    () => {
      populatePrograms(props.activityID);
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
        <div
          v-if="
            scope.opt.programID && isCurrent(scope.opt.activityID, scope.opt.id)
          "
          :class="[styles.pgmInstItemCurrent]"
        >
          {{ scope.opt.title }}
        </div>
        <div
          v-else-if="
            scope.opt.programID && isActive(scope.opt.activityID, scope.opt.id)
          "
          :class="[styles.pgmInstItem]"
        >
          {{ scope.opt.title }}
        </div>
        <div
          v-else-if="scope.opt.programID && !props.hideCompleted"
          :class="[styles.pgmInstItemPast]"
        >
          {{ scope.opt.title }} (complete)
        </div>
        <div v-else-if="!scope.opt.programID">{{ scope.opt.title }}</div>
      </q-item>
    </template>
    <template v-slot:selected>
      <div v-if="selectedObj">
        {{ selectedObj.title }}
      </div>
    </template>
  </q-select>
</template>
