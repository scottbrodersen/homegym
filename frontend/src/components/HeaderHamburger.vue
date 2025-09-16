<script setup>
  /**
   * A hamburger-style drop-down menu that enables users to enter and view daily statistics for today.
   */
  import * as dailyStatsUtils from '../modules/dailyStatsUtils';
  import * as styles from '../style.module.css';
  import * as state from '../modules/state';
  import {
    QBtn,
    QList,
    QItem,
    QItemLabel,
    QItemSection,
    QMenu,
    QSeparator,
    QToggle,
  } from 'quasar';
  const dailyStats = (statName) => {
    let existingStat;
    if (state.dailyStatsStore[statName]) {
      existingStat = state.dailyStatsStore[statName];
    }
    dailyStatsUtils
      .openDailyStatsModal(statName, existingStat)
      .then(async (stat) => {
        if (stat) {
          await dailyStatsUtils.saveDailyStat(stat);
          state.dailyStatsStore.update(stat);
        }
      });
  };

  const toggleMetric = (isMetric) => {
    if (isMetric) {
      state.metricState.setMetric();
    } else {
      state / metricState.setImperial();
    }
  };
</script>
<template>
  <q-btn icon="menu" :class="[styles.hgHamburger]">
    <q-menu>
      <q-list dense dark>
        <q-item-label header>Daily Stats</q-item-label>
        <q-item dark clickable v-close-popup @click="() => dailyStats('sleep')">
          <q-item-section>{{
            dailyStatsUtils.labels['sleep'][0]
          }}</q-item-section>
          <q-item-section side>{{
            state.dailyStatsStore.sleep.sleep
          }}</q-item-section>
        </q-item>
        <q-item
          dark
          clickable
          v-close-popup
          @click="() => dailyStats('spirit')"
        >
          <q-item-section>{{
            dailyStatsUtils.labels['mood'][0]
          }}</q-item-section>
          <q-item-section side>{{
            state.dailyStatsStore.spirit.mood
          }}</q-item-section>
        </q-item>
        <q-item
          dark
          clickable
          v-close-popup
          @click="() => dailyStats('spirit')"
        >
          <q-item-section>{{
            dailyStatsUtils.labels['stress'][0]
          }}</q-item-section>
          <q-item-section side>{{
            state.dailyStatsStore.spirit.stress
          }}</q-item-section></q-item
        >
        <q-item
          dark
          clickable
          v-close-popup
          @click="() => dailyStats('spirit')"
        >
          <q-item-section>{{
            dailyStatsUtils.labels['energy'][0]
          }}</q-item-section>
          <q-item-section side>{{
            state.dailyStatsStore.spirit.energy
          }}</q-item-section>
        </q-item>
        <q-item
          dark
          clickable
          v-close-popup
          @click="() => dailyStats('bodyweight')"
        >
          <q-item-section>{{
            dailyStatsUtils.labels['bodyweight'][0]
          }}</q-item-section>
          <q-item-section side>{{
            state.dailyStatsStore.bodyweight.bodyweight
          }}</q-item-section>
        </q-item>
        <q-separator dark spaced />

        <q-item dark clickable v-close-popup @click="() => dailyStats('bg')">
          <q-item-label>{{ dailyStatsUtils.labels['bg'][0] }}</q-item-label>
        </q-item>
        <q-item dark clickable v-close-popup @click="() => dailyStats('bp')">
          <q-item-label>{{ dailyStatsUtils.labels['bp'][0] }}</q-item-label>
        </q-item>
        <q-item dark clickable v-close-popup @click="() => dailyStats('food')">
          <q-item-label>{{ dailyStatsUtils.labels['food'][0] }}</q-item-label>
        </q-item>
        <q-item dark>
          <q-toggle
            label="Metric"
            label-left
            :model-value="state.metricState.metric"
            @update:model-value="toggleMetric"
          />
        </q-item>
      </q-list>
    </q-menu>
  </q-btn>
</template>
