import { createApp, ref } from 'vue';
import { Quasar, Dialog, Notify } from 'quasar';
import '@quasar/extras/material-icons/material-icons.css';
import iconSet from 'quasar/icon-set/material-icons';
import 'quasar/src/css/index.sass';
import './style.module.css';
import App from './App.vue';
import router from './routes.js';
import { select, focus } from './modules/directives.js';
import 'vite/modulepreload-polyfill';
import faviconUrl from './favicon.png';

// Dynamically create or update the favicon link
const link = document.createElement('link');
link.rel = 'icon';
link.type = 'image/png';
link.href = faviconUrl;
document.head.appendChild(link);

const app = createApp(App);

// context sensitive help
app.provide('docsRootURL', 'https://scottbrodersen.github.io/homegym/');
app.provide('docsContextQuery', 'context');
const docsContext = ref();
app.provide('docsContext', docsContext);

app.use(Quasar, {
  iconSet: iconSet,
  plugins: { Dialog, Notify }, // Quasar plugins
});

app.use(router);
app.directive('focus', focus);
app.directive('select', select);
app.mount('#app');
