<script async setup>
  import {
    authPromptAsync,
    fetchActivities,
    fetchExerciseTypes,
    fetchActiveProgramInstances,
    fetchProgramInstances,
    ErrNotLoggedIn,
    fetchPrograms,
  } from '../modules/utils';
  import { selectCurrentProgramInstance } from '../modules/programUtils';
  import * as dailyStatsUtils from '../modules/dailyStatsUtils';
  import * as dateUtils from '../modules/dateUtils';
  import {
    activityStore,
    exerciseTypeStore,
    dailyStatsStore,
    programInstanceStore,
    programsStore,
  } from '../modules/state';
  import * as styles from '../style.module.css';
  import { getCurrentInstance } from 'vue';

  const init = async () => {
    const endDate = dateUtils.setEpochToMidnight(dateUtils.nowInSeconds());
    const startDate = endDate + 24 * 60 * 60;

    try {
      // Get daily stats from today only
      const dailyStats = await dailyStatsUtils.fetchDailyStats(
        startDate,
        endDate
      );

      dailyStatsStore.bulkAdd(dailyStats);

      if (exerciseTypeStore.exerciseTypes.size === 0) {
        await fetchExerciseTypes();
        await fetchActivities();

        const promises = [];
        for (const activity of activityStore.getAll()) {
          promises.push(fetchPrograms(activity.id));
          promises.push(fetchActiveProgramInstances(activity.id));
        }
        await Promise.all(promises);

        for (const activity of activityStore.getAll()) {
          const promises2 = [];
          for (const program of programsStore.getByActivity(activity.id)) {
            promises2.push(fetchProgramInstances(program.id, activity.id));
          }
          await Promise.all(promises2);
        }

        for (const activity of activityStore.getAll()) {
          programInstanceStore.setCurrent(
            activity.id,
            selectCurrentProgramInstance(activity.id)
          );
        }
        console.log('done init state');
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();

        await init();
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
