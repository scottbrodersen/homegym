<script setup>
  /**
   * The top navigation bar.
   */
  import { reactive, ref, watch, inject } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import {
    QBtn,
    QBtnDropdown,
    QBtnGroup,
    QList,
    QItem,
    QItemLabel,
  } from 'quasar';
  import { getCookieValue } from '../modules/utils';
  import * as styles from '../style.module.css';
  import HeaderHamburger from './HeaderHamburger.vue';

  const router = useRouter();
  const route = useRoute();

  const paths = {
    '/homegym/event/': 'event',
    '/homegym/activities/': 'activities',
    '/homegym/analyze/': 'analyze',
    '/homegym/exercises/': 'exTypes',
    '/homegym/programs/': 'programs',
  };

  // docs context
  const rootDocsURL = inject('docsRootURL');
  const docsContextQuery = inject('docsContextQuery');
  const docsContext = inject('docsContext');
  const docsURL = ref(
    rootDocsURL + '?' + docsContextQuery + '=' + docsContext.value
  );
  const hgdocs = ref('hgdocs');

  watch(
    () => docsContext.value,
    (context) => {
      docsURL.value =
        rootDocsURL + '?' + docsContextQuery + '=' + docsContext.value;
    }
  );

  // field names match route names
  const active = reactive({
    home: true,
    event: false,
    activities: false,
    analyze: false,
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

  const padding = ref('4px 4px');

  const labels = reactive({
    home: 'Home',
    event: 'Event',
    activities: 'Activities',
    analyze: 'Analyze',
    exTypes: 'Exercises',
    programs: 'Program',
    activity: 'Activity',
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
        :to="{ name: 'home' }"
      />
      <q-btn
        :class="{ active: active.programs }"
        :unelevated="active.programs"
        :glossy="!active.programs"
        :label="labels.programs"
        square
        :to="{ name: 'programs' }"
      />
      <q-btn-dropdown
        :class="{ active: active.activities }"
        :unelevated="active.activities"
        :glossy="!active.activities"
        :label="labels.activity"
        square
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
      <q-btn
        :class="{ active: active.analyze }"
        :unelevated="active.analyze"
        :glossy="!active.analyze"
        :label="labels.analyze"
        square
        :to="{ name: 'analyze' }"
      />
      <HeaderHamburger />
      <q-btn icon="help" :href="docsURL" target="hgdocs" />
    </q-btn-group>
  </header>
</template>
