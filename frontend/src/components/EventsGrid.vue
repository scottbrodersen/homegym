<script setup>
  import {
    eventStore,
    activityStore,
    exerciseTypeStore,
    eventMetricsStore,
  } from '../modules/state.js';
  import { computed, ref, onBeforeMount, Suspense } from 'vue';
  import { QTable, QTr, QTh, QTd, QBtn } from 'quasar';
  import {
    authPrompt,
    fetchEventPage,
    pageSize,
    ErrNotLoggedIn,
  } from '../modules/utils.js';
  import EventExercises from './EventExercises.vue';
  import styles from '../style.module.css';
  import EventMeta from './EventMeta.vue';

  const table = ref();
  const expanded = ref([]);

  // array of eventstore events
  let rows = ref([]);

  // pagination state
  const pagination = ref({
    sortBy: 'date',
    page: 0,
    rowsPerPage: pageSize,
    rowsNumber: 0,
  });

  const pagesNumber = computed(() => {
    return eventStore.events.length / pageSize;
  });

  const columns = [
    {
      name: 'date',
      required: true,
      label: 'Date',
      align: 'left',
      field: 'date',
      format: (val) => {
        // stored time is in seconds
        const date = new Date(val * 1000);
        return `${date.toLocaleDateString()}`;
      },
      sortable: false,
    },
    {
      name: 'activity',
      required: true,
      label: 'Activity',
      align: 'left',
      field: 'activityID',
      format: (val) => {
        return activityStore.get(val).name;
      },
      sortable: false,
    },
  ];

  // handles table pagination
  // tops up the eventStore when the last page of the store is requested
  const setPage = async (props) => {
    try {
      // fetch a page if we are showing the last page
      if (
        eventStore.events.length === 0 ||
        props.pagination.page >= eventStore.events.length / pageSize
      ) {
        const lastEvent = eventStore.getLast();
        const lastEventID = lastEvent ? lastEvent.id : 0;
        const lastEventDate = lastEvent ? lastEvent.date : 0;

        const events = await fetchEventPage(lastEventID, lastEventDate);

        // initialize objects in metrics store
        events.forEach((event) => {
          eventMetricsStore.add(event.id, {});
        });
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(setPage, [props]);
      } else {
        console.log(e);
      }
    }

    rows.value = eventStore.getPage(props.pagination.page - 1);
    pagination.value.rowsNumber = eventStore.events.length;
    pagination.value.page = props.pagination.page;
  };

  const initState = async () => {
    try {
      await setPage({ pagination: { page: 1 } });
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        authPrompt(initState);
      } else {
        console.log(e);
      }
    }
  };

  onBeforeMount(initState);

  // custom expand row function to allow only one row to be expanded at a time
  const expandRow = (props) => {
    // close the currently-expanded row or expand the row to expand
    const expandedRowID = expanded.value.pop();
    expanded.value = expandedRowID === props.row.id ? [] : [props.row.id];
  };

  // row gradient
  const background = (eventID) => {
    let count = 0;
    let total = 0;
    const e = eventStore.getByID(eventID);

    if (e.mood) {
      count++;
      total += e.mood;
    }
    if (e.energy) {
      count++;
      total += e.energy;
    }
    if (e.motivation) {
      count++;
      total += e.motivation;
    }

    const start = count > 0 ? Math.round(total / count) : 0;
    const end = e.overall > 0 ? e.overall : 0;

    const colours = {
      0: '--mood0',
      1: '--mood1',
      2: '--mood2',
      3: '--mood3',
      4: '--mood4',
      5: '--mood5',
    };
    const layer1 = '90';
    const layer2 = '50';

    const mask =
      ', linear-gradient(to right, #c0c0c050, #ffffff50, #c0c0c050), radial-gradient(#ffffff, #c0c0c0)';
    return `background:linear-gradient(var(${colours[end]}) 15%, var(${colours[start]})) 85% ${mask};`;
  };
</script>

<template>
  <div :class="[styles.horiz, styles.alignCenter, styles.blockPadSm]">
    <div :class="[styles.sibSpMed]">
      <h1>Event Log</h1>
    </div>
    <div :class="[styles.sibSpMed]">
      <q-btn
        round
        size="0.65em"
        color="primary"
        icon="add"
        :to="{ name: 'event' }"
        id="addevent"
      />
    </div>
  </div>
  <div :class="styles.eventsTable">
    <q-table
      ref="table"
      :rows="rows"
      :columns="columns"
      row-key="id"
      v-model:pagination="pagination"
      v-model:expanded="expanded"
      :rowsPerPageOptions="[]"
      @request="setPage"
      dark
    >
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th></q-th>
          <q-th :props="props" v-for="col in props.cols" :key="col.name">
            {{ col.label }}
          </q-th>
          <q-th>Volume</q-th>
          <q-th>Load</q-th>
          <q-th></q-th>
        </q-tr>
      </template>
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td>
            <q-btn
              no-hover
              flat
              dense
              round
              :ripple="false"
              @click="expandRow(props)"
              size="1.5em"
              :icon="props.expand ? 'arrow_drop_down' : 'arrow_right'"
            />
          </q-td>
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.value }}
          </q-td>
          <q-td>{{ eventMetricsStore.getMetric(props.row.id, 'volume') }}</q-td>
          <q-td>{{ eventMetricsStore.getMetric(props.row.id, 'load') }}</q-td>
          <q-td :style="background(props.row.id)" />
        </q-tr>
        <q-tr v-show="props.expand" :props="props">
          <q-td :class="[styles.expanded, styles.blockPadSm]" colspan="5">
            <q-btn
              color="primary"
              round
              :to="{ name: 'event', params: { eventId: props.key } }"
              icon="edit"
            />
            <EventMeta
              :mood="props.row.mood"
              :energy="props.row.energy"
              :motivation="props.row.motivation"
              :readonly="true"
            />
            <Suspense :key="pagination.page" timeout="0">
              <EventExercises :event-id="props.row.id" />
              <template #fallback> Loading... </template>
            </Suspense>
            <EventMeta
              :overall="props.row.overall"
              :notes="props.row.notes"
              :readonly="true"
            />
          </q-td>
        </q-tr>
      </template>
      <template v-slot:bottom="scope">
        <q-btn
          :class="[styles.maxRight]"
          v-if="scope.pagesNumber > 2"
          icon="first_page"
          color="grey-8"
          round
          dense
          flat
          :disable="scope.isFirstPage"
          @click="scope.firstPage"
        />
        <q-btn
          :class="scope.pagesNumber > 2 ? '' : [styles.maxRight]"
          icon="chevron_left"
          color="grey-8"
          round
          dense
          flat
          :disable="scope.isFirstPage"
          @click="scope.prevPage"
        />
        <div>Page {{ pagination.page }} of {{ Math.ceil(pagesNumber) }}</div>
        <q-btn
          :class="scope.pagesNumber > 2 ? '' : [styles.maxLeft]"
          icon="chevron_right"
          color="grey-8"
          round
          dense
          flat
          :disable="scope.isLastPage"
          @click="scope.nextPage"
        />
        <q-btn
          :class="[styles.maxLeft]"
          v-if="scope.pagesNumber > 2"
          icon="last_page"
          color="grey-8"
          round
          dense
          flat
          :disable="scope.isLastPage"
          @click="scope.lastPage"
        />
      </template>
    </q-table>
  </div>
</template>
