<script async setup>
  import {
    authPrompt,
    fetchActivities,
    fetchExerciseTypes,
    ErrNotLoggedIn,
  } from '../modules/utils';
  import { exerciseTypeStore } from '../modules/state';
  import styles from '../style.module.css';

  try {
    if (exerciseTypeStore.exerciseTypes.size === 0) {
      await fetchExerciseTypes();
      await fetchActivities();
    }
  } catch (e) {
    if (e instanceof ErrNotLoggedIn) {
      console.log(e.message);
      authPrompt(setup);
    } else {
      console.log(e);
    }
  }
</script>
<template>
  <div :id="styles.wrapper">
    <router-view />
    <main :class="[styles.blockPadSm]">
      <router-view name="main" />
    </main>
  </div>
</template>
