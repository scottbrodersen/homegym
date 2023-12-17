<script setup>
  import { reactive, ref, watch } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { QBtn, QBtnGroup, QToggle } from 'quasar';
  import { metricState } from '../modules/state.js';
  import styles from '../style.module.css';

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

  const router = useRouter();
  const route = useRoute();

  window.addEventListener('popstate', (event) => {
    router.replace({ path: event.state.current });
  });

  watch(
    () => route.name,
    (newname) => {
      setActive(newname);
    }
  );

  const removeToken = () => {
    document.cookie =
      'token=; Domain=localhost; Path=/homegym/; Max-Age=-99999999;';
  };

  const showTestTools = ref(process.env.TEST_TOOLS);
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
        :class="{ active: active.activities }"
        :unelevated="active.activities"
        :glossy="!active.activities"
        :label="labels.activities"
        square
        :padding="padding"
        :to="{ name: 'activities' }"
      />
      <q-btn
        :class="{ active: active.exTypes }"
        :unelevated="active.exTypes"
        :glossy="!active.exTypes"
        :label="labels.exTypes"
        square
        :padding="padding"
        :to="{ name: 'exTypes' }"
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
      <q-btn
        v-if="showTestTools"
        color="white"
        text-color="black"
        label="ZapToken"
        @click="removeToken"
      />
    </q-btn-group>
    <q-toggle
      :label="labels.metric"
      label-left
      :model-value="metricState.metric"
      @update:model-value="toggleMetric"
    />
  </header>
</template>
