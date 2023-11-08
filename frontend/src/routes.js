import * as VueRouter from 'vue-router';
import EventsGrid from './components/EventsGrid.vue';
import ActivityPage from './components/ActivityPage.vue';
import EventPage from './components/EventPage.vue';
import PageHeader from './components/PageHeader.vue';
import ExerciseTypesPage from './components/ExerciseTypesPage.vue';

const routes = [
  {
    path: '/homegym/home/',
    components: { default: PageHeader, main: EventsGrid },
    name: 'home',
  },
  {
    path: '/homegym/activities/',
    components: { default: PageHeader, main: ActivityPage },
    name: 'activities',
  },
  {
    path: '/homegym/event/:eventId?',
    components: { default: PageHeader, main: EventPage },
    name: 'event',
    props: { main: true },
  },
  {
    path: '/homegym/exercises/',
    components: { default: PageHeader, main: ExerciseTypesPage },
    name: 'exTypes',
  },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(),
  routes: routes,
});

export default router;
