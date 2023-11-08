<script setup>
  import { reactive, ref, watch } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { QBtn, QBtnGroup, QToggle } from 'quasar';
  import { metricState } from '../modules/state.js';
  import styles from '../style.module.css';

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
    metric: 'Metric',
  });

  // field names match route names
  const active = reactive({
    home: false,
    event: false,
    activities: false,
    exTypes: false,
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
    <nav>
      <q-btn-group push square>
        <q-btn
          :class="{ active: active.home }"
          :unelevated="active.home"
          :glossy="!active.home"
          :label="labels.home"
          square
          :to="{ name: 'home' }"
        />
        <q-btn
          :class="{ active: active.activities }"
          :unelevated="active.activities"
          :glossy="!active.activities"
          :label="labels.activities"
          square
          :to="{ name: 'activities' }"
        />
        <q-btn
          :class="{ active: active.exTypes }"
          :unelevated="active.exTypes"
          :glossy="!active.exTypes"
          :label="labels.exTypes"
          square
          :to="{ name: 'exTypes' }"
        />

        <q-btn
          v-if="showTestTools"
          color="white"
          text-color="black"
          label="ZapToken"
          @click="removeToken"
        />
      </q-btn-group>
    </nav>
    <q-toggle
      :label="labels.metric"
      label-left
      :model-value="metricState.metric"
      @update:model-value="toggleMetric"
    />
  </header>
</template>
