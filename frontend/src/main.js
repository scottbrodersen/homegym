import { createApp } from 'vue';
import { Quasar, Dialog, Notify } from 'quasar';
import '@quasar/extras/material-icons/material-icons.css';
import iconSet from 'quasar/icon-set/material-icons';
import 'quasar/src/css/index.sass';
import './style.module.css';
import App from './App.vue';
import router from './routes.js';
import { select, focus } from './modules/directives.js';
import 'vite/modulepreload-polyfill';

const app = createApp(App);

app.use(Quasar, {
  iconSet: iconSet,
  plugins: { Dialog, Notify }, // Quasar plugins
});

app.use(router);
app.directive('focus', focus);
app.directive('select', select);
app.mount('#app');
