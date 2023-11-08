import { createApp } from "vue";
import { Quasar, Dialog } from "quasar";
import "@quasar/extras/material-icons/material-icons.css";
import iconSet from "quasar/icon-set/material-icons";
import "quasar/src/css/index.sass";
import "./style.module.css";
import App from "./App.vue";
import router from "./routes.js";

const app = createApp(App);

app.use(Quasar, {
  iconSet: iconSet,
  plugins: { Dialog }, // Quasar plugins
});
app.use(router);
app.mount("#app");
