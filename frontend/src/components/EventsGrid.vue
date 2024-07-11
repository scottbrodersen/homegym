<script async setup>
  import {
    eventStore,
    activityStore,
    eventMetricsStore,
  } from '../modules/state.js';
  import { computed, inject, onMounted, ref, watch } from 'vue';
  import { QTable, QTr, QTd, QBtn } from 'quasar';
  import {
    authPromptAsync,
    fetchEventPage,
    pageSize,
    ErrNotLoggedIn,
  } from '../modules/utils.js';
  import * as styles from '../style.module.css';
  import EventDetails from './EventDetails.vue';

  const props = defineProps({ eventID: String });
  const table = ref();
  const expanded = ref([]);

  // array of eventstore events
  let rows = ref([]);

  // pagination state
  const pagination = ref({
    sortBy: 'date',
    page: 0,
    rowsPerPage: pageSize(),
    rowsNumber: 0,
  });

  const pagesNumber = computed(() => {
    return eventStore.events.length / pageSize();
  });

  const { focusedEvent, setFocusedEvent } = inject('focusedEvent');
  const { selectedEvent, setSelectedEvent } = inject('selectedEvent');

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
    {
      name: 'volume',
      required: true,
      label: 'Volume',
      align: 'left',
      field: 'id',
      format: (id) => {
        return eventMetricsStore.getMetric(id, 'volume');
      },
    },
    {
      name: 'load',
      required: true,
      label: 'Load',
      align: 'left',
      field: 'id',
      format: (id) => {
        return eventMetricsStore.getMetric(id, 'load');
      },
    },
    {
      name: 'meta',
      required: true,
      field: 'id',
      format: () => {
        return null;
      },
      style: (row) => background(row.id),
    },
  ];

  // handles table pagination
  // tops up the eventStore when the last page of the store is requested
  const setPage = async (props) => {
    try {
      // fetch a page if we are showing the last page
      if (
        eventStore.events.length === 0 ||
        props.pagination.page >= eventStore.events.length / pageSize()
      ) {
        const lastEvent = eventStore.getLast();
        const lastEventID = lastEvent ? lastEvent.id : 0;
        const lastEventDate = lastEvent ? lastEvent.date : 0;

        const events = await fetchEventPage(lastEventID, lastEventDate);

        // initialize objects in metrics store
        events.forEach((event) => {
          setMetrics(event);
        });
      }
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
        setPage(props);
      } else {
        console.log(e);
      }
    }

    rows.value = eventStore.getPage(props.pagination.page - 1);
    pagination.value.rowsNumber = eventStore.events.length;
    pagination.value.page = props.pagination.page;
  };

  const setMetrics = (event) => {
    // volume is the total reps done in an event
    let volume = 0;
    // load is total work done in an event
    let load = 0;

    for (const index of Object.keys(event.exercises)) {
      event.exercises[index].parts.forEach((part) => {
        part.volume.forEach((set) => {
          set.forEach((rep) => {
            if (rep != 0) {
              volume++;
              load += Math.floor(part.intensity * rep);
            }
          });
        });
      });
    }

    eventMetricsStore.setMetric(event.id, 'volume', volume);
    eventMetricsStore.setMetric(event.id, 'load', load);
  };

  await setPage({ pagination: { page: 1 } });

  // custom expand row function to allow only one row to be expanded at a time
  const expandRow = async (rowID) => {
    // close the currently-expanded row or expand the row to expand
    const expandedRowID = expanded.value.pop();
    expanded.value = expandedRowID === rowID ? [] : [rowID];

    await toRowPage(rowID);

    setSelectedEvent(expanded.value[0]);
  };

  onMounted(() => {
    if (props.eventID) {
      expandRow(props.eventID);
    }
  });

  watch(
    () => {
      return props.eventID;
    },
    async (newID) => {
      if (expanded.value.length > 0 && expanded.value[0] != newID) {
        expandRow(newID);
      }
    }
  );

  watch(
    () => {
      return focusedEvent.value;
    },
    async (newID) => {
      await toRowPage(newID);
    }
  );

  // turns to the page that contains the event
  const toRowPage = async (eventID) => {
    let pageNumber;
    let events = eventStore.getAll();
    for (let i = 0; i < events.length; i++) {
      if (events[i].id == eventID) {
        pageNumber = Math.floor(i / pagination.value.rowsPerPage) + 1;
        break;
      }
    }
    if (pageNumber && pageNumber != pagination.value.page) {
      await setPage({ pagination: { page: pageNumber } });
    }
  };
</script>

<template>
  <div>
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
        <template v-slot:body="props">
          <q-tr
            :props="props"
            :id="props.key"
            :class="props.key == focusedEvent ? styles.evtHighlight : ''"
          >
            <q-td
              v-for="col in props.cols"
              :key="col.name"
              :props="props"
              @click="expandRow(props.row.id)"
            >
              {{ col.value }}
            </q-td>
          </q-tr>
          <Transition name="scale">
            <EventDetails
              v-show="props.expand"
              :event-id="props.key"
              :key="props.key"
              class="slider"
            />
          </Transition>
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
  </div>
</template>
<style>
  .scale-enter-active {
    transition: 0.25s;
  }
  .scale-enter-from {
    transform: scaleY(1.05);
  }
</style>
