import { createApp } from 'vue';
import { Quasar, Dialog, Notify } from 'quasar';
import '@quasar/extras/material-icons/material-icons.css';
import iconSet from 'quasar/icon-set/material-icons';
import 'quasar/src/css/index.sass';
import './style.module.css';
import App from './App.vue';
import router from './routes.js';

const app = createApp(App);

app.use(Quasar, {
  iconSet: iconSet,
  plugins: { Dialog, Notify }, // Quasar plugins
});

app.use(router);

app.directive('focus', {
  mounted: (el) => {
    const input =
      el.getElementsByTagName('input').length > 0
        ? el.getElementsByTagName('input')[0]
        : null;

    if (
      !!input &&
      input.hasAttribute('type') &&
      input.getAttribute('type') == 'number'
    ) {
      if (input.value == '0') {
        input.focus();
      }
    } else {
      el.focus();
    }
  },
});
app.directive('select', {
  mounted: (el) => {
    const input =
      el.getElementsByTagName('input').length > 0
        ? el.getElementsByTagName('input')[0]
        : null;
    if (!!input && input.value == '0') {
      input.select();
    }
  },
});
app.mount('#app');
