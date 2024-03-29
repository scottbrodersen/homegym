import * as VueRouter from 'vue-router';
import EventPage from './components/EventPage.vue';
import ActivityPage from './components/ActivityPage.vue';
import HomePage from './components/HomePage.vue';
import PageHeader from './components/PageHeader.vue';
import ExerciseTypesPage from './components/ExerciseTypesPage.vue';
import ProgramsPage from './components/ProgramsPage.vue';

const routes = [
  {
    path: '/homegym/home/',
    components: { default: PageHeader, main: HomePage },
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
  {
    path: '/homegym/programs/',
    components: { default: PageHeader, main: ProgramsPage },
    name: 'programs',
  },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(),
  routes: routes,
});

export default router;
