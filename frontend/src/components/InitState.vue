<script async setup>
  import { onBeforeMount } from 'vue';
  import {
    authPrompt,
    fetchActivities,
    fetchExerciseTypes,
    fetchActiveProgramInstance,
    ErrNotLoggedIn,
    fetchPrograms,
  } from '../modules/utils';
  import { activityStore, exerciseTypeStore } from '../modules/state';
  import styles from '../style.module.css';

  const init = async () => {
    try {
      if (exerciseTypeStore.exerciseTypes.size === 0) {
        await fetchExerciseTypes();
        await fetchActivities();

        const promises = [];
        for (const activity of activityStore.getAll()) {
          promises.push(fetchPrograms(activity.id));
          promises.push(fetchActiveProgramInstance(activity.id));
        }
        await Promise.all(promises);
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(init);
      } else {
        console.log(e);
      }
    }
  };
  await init();
</script>
<template>
  <div :id="styles.wrapper">
    <router-view />
    <main :class="[styles.blockPadSm]">
      <router-view name="main" />
    </main>
  </div>
</template>
