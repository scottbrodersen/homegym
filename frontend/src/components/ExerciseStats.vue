<script setup>
  import { onBeforeMount, ref, watch } from 'vue';
  import {
    ErrNotLoggedIn,
    fetchExercise1RM,
    fetchExercisePR,
    openEditValueModal,
    updateExercisePR,
    updateExercise1RM,
    authPromptAsync,
  } from '../modules/utils.js';
  import * as styles from '../style.module.css';

  const props = defineProps({ exerciseID: String });

  const oneRM = ref('not set');
  const pr = ref('not set');

  const editPR = async () => {
    const value = [{ label: 'Enter your personal record', value: pr.value }];

    openEditValueModal(value).then(async (edited) => {
      if (edited) {
        try {
          await updateExercisePR(props.exerciseID, edited[0]);
          pr.value = edited[0];
        } catch (e) {
          if (e instanceof ErrNotLoggedIn) {
            console.log(e.message);
            await authPromptAsync();
            editPR();
          } else {
            console.log(e);
          }
        }
      }
    });
  };
  const edit1RM = async () => {
    const value = [{ label: 'Enter your 1 rep max', value: oneRM.value }];

    openEditValueModal(value).then(async (edited) => {
      if (edited) {
        try {
          await updateExercise1RM(props.exerciseID, edited[0]);
          oneRM.value = edited[0];
        } catch (e) {
          if (e instanceof ErrNotLoggedIn) {
            console.log(e.message);
            await authPromptAsync();
            edit1RM();
          } else {
            console.log(e);
          }
        }
      }
    });
  };
  const getStats = async (exerciseID) => {
    try {
      const raw1RM = await fetchExercise1RM(exerciseID);

      if (raw1RM > 0) {
        oneRM.value = raw1RM;
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
        getStats(exerciseID);
      } else {
        console.log(e);
      }
    }
    try {
      const rawPR = await fetchExercisePR(exerciseID);

      if (rawPR > 0) {
        pr.value = rawPR;
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
        getStats(exerciseID);
      } else {
        console.log(e);
      }
    }
  };

  watch(
    () => props.exerciseID,
    async () => {
      if (props.exerciseID) {
        pr.value = 'not set';
        oneRM.value = 'not set';
        await getStats(props.exerciseID);
      }
    },
  );

  onBeforeMount(async () => {
    if (props.exerciseID) {
      await getStats(props.exerciseID);
    }
  });
</script>
<template>
  <div>
    <div :class="[styles.exerciseStatsWrapper]">
      <div @click="editPR()">
        <span :class="[styles.hgBold]">PR:</span>
        <span :class="[styles.exerciseStatValue]">{{ pr }}</span>
      </div>
      <div @click="edit1RM()">
        <span :class="[styles.hgBold]">1RM:</span>
        <span :class="[styles.exerciseStatValue]">{{ oneRM }}</span>
      </div>
    </div>
  </div>
</template>
