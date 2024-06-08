<script setup>
  import { reactive, ref, watch } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import {
    QBtn,
    QBtnDropdown,
    QBtnGroup,
    QList,
    QItem,
    QItemLabel,
    QMenu,
    QToggle,
  } from 'quasar';
  import { metricState } from '../modules/state';
  import { getCookieValue } from '../modules/utils';
  import * as styles from '../style.module.css';

  const router = useRouter();
  const route = useRoute();

  const paths = {
    '/homegym/event/': 'event',
    '/homegym/activities/': 'activities',
    '/homegym/exercises/': 'exTypes',
    '/homegym/programs/': 'programs',
  };

  // field names match route names
  const active = reactive({
    home: true,
    event: false,
    activities: false,
    exTypes: false,
    programs: false,
  });

  const setActive = (activeId) => {
    for (const id in active) {
      if (id === activeId) {
        active[id] = true;
      } else {
        active[id] = false;
      }
    }
  };

  if (getCookieValue('followroute')) {
    const follow = getCookieValue('followroute');
    document.cookie =
      'followroute=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/homegym/';
    if (follow.startsWith('/homegym/event/')) {
      const parts = follow.split('/');
      const eventID = parts[parts.length - 1];
      await router.replace({ name: 'event', params: { eventId: eventID } });
    } else {
      setActive(paths[follow]);
      await router.replace(follow);
    }
  }

  const padding = ref('4px 10px');

  const toggleMetric = (isMetric) => {
    if (isMetric) {
      metricState.setMetric();
    } else {
      metricState.setImperial();
    }
  };

  const labels = reactive({
    home: 'Home',
    event: 'Event',
    activities: 'Activities',
    exTypes: 'Exercises',
    programs: 'Programs',
    metric: 'Metric',
  });

  window.addEventListener('popstate', (event) => {
    router.replace({ path: event.state.current });
  });

  watch(
    () => route.name,
    (newname) => {
      setActive(newname);
    }
  );
</script>

<template>
  <header>
    <q-btn-group push square stretch>
      <q-btn
        :class="{ active: active.home }"
        :unelevated="active.home"
        :glossy="!active.home"
        :label="labels.home"
        square
        :padding="padding"
        :to="{ name: 'home' }"
      />
      <q-btn
        :class="{ active: active.programs }"
        :unelevated="active.programs"
        :glossy="!active.programs"
        :label="labels.programs"
        square
        :padding="padding"
        :to="{ name: 'programs' }"
      />
      <q-btn-dropdown
        :class="{ active: active.activities }"
        :unelevated="active.activities"
        :glossy="!active.activities"
        :label="labels.activities"
        square
        :padding="padding"
      >
        <q-list>
          <q-item dark clickable :to="{ name: 'activities' }">
            <q-item-label dark>{{ labels.activities }}</q-item-label>
          </q-item>
          <q-item clickable :to="{ name: 'exTypes' }">
            <q-item-label>{{ labels.exTypes }}</q-item-label>
          </q-item>
        </q-list>
      </q-btn-dropdown>
      <q-btn icon="menu" :class="[styles.hgHamburger]">
        <q-menu>
          <q-list :class="[styles.hgMenu]">
            <q-item dark>
              <q-toggle
                :label="labels.metric"
                label-left
                :model-value="metricState.metric"
                @update:model-value="toggleMetric"
              />
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </q-btn-group>
  </header>
</template>
