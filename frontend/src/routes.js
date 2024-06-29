import * as VueRouter from 'vue-router';
import EventPage from './components/EventPage.vue';
import ActivityPage from './components/ActivityPage.vue';
import HomePage from './components/HomePage.vue';
import PageHeader from './components/PageHeader.vue';
import ExerciseTypesPage from './components/ExerciseTypesPage.vue';
import ProgramsPage from './components/ProgramsPage.vue';
import AnalyzePage from './components/AnalyzePage.vue';

const routes = [
  {
    path: '/homegym/home/:event?',
    components: { default: PageHeader, main: HomePage },
    name: 'home',
    props: {
      main: (route) => ({
        eventID: route.query.event,
      }),
    },
  },
  {
    path: '/homegym/activities/',
    components: { default: PageHeader, main: ActivityPage },
    name: 'activities',
  },
  {
    path: '/homegym/analyze/',
    components: { default: PageHeader, main: AnalyzePage },
    name: 'analyze',
  },
  {
    path: '/homegym/event/:eventId?',
    components: { default: PageHeader, main: EventPage },
    name: 'event',
    props: {
      main: (route) => ({
        eventId: route.params.eventId,
        programInstanceID: route.query.instance,
        dayIndex: route.query.day,
        blockIndex: route.query.block,
        microCycleIndex: route.query.cycle,
        workoutIndex: route.query.workout,
      }),
    },
  },
  {
    path: '/homegym/exercises/',
    components: { default: PageHeader, main: ExerciseTypesPage },
    name: 'exTypes',
  },
  {
    path: '/homegym/programs/',
    components: { default: PageHeader, main: ProgramsPage },
    name: 'programs',
    props: {
      main: (route) => ({
        activityID: route.query.activity,
        programID: route.query.program,
        instanceID: route.query.instance,
      }),
    },
  },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(),
  routes: routes,
});

export default router;
