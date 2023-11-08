<script setup>
import { useDialogPluginComponent } from "quasar";
import { login } from "./../modules/utils.js";
import { ref } from "vue";
import { loginModalState } from "../modules/state";
let password = ref("");
let isPwd = ref(true);
let username = ref("");

defineEmits([
  // REQUIRED; need to specify some events that your
  // component will emit through useDialogPluginComponent()
  ...useDialogPluginComponent.emits,
]);

const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
  useDialogPluginComponent();
// dialogRef      - Vue ref to be applied to QDialog
// onDialogHide   - Function to be used as handler for @hide on QDialog
// onDialogOK     - Function to call to settle dialog with "ok" outcome
//                    example: onDialogOK() - no payload
//                    example: onDialogOK({ /*...*/ }) - with payload
// onDialogCancel - Function to call to settle dialog with "cancel" outcome

function onOKClick() {
  login(username.value, password.value)
    .then(() => {
      onDialogOK();
    })
    .catch((e) => {
      console.log(e);
    });
}
</script>
<template>
  <q-dialog persistent ref="dialogRef" @hide="onDialogHide">
    <q-card class="q-dialog-plugin">
      <div>Please sign in</div>
      <form>
        <q-input
          v-model="username"
          filled
          type="text"
          label="User Name"
          autocomplete="username"
        />

        <q-input
          v-model="password"
          filled
          :type="isPwd ? 'password' : 'text'"
          label="Password"
          autocomplete="current-password"
        >
          <template v-slot:append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>
      </form>

      <q-card-actions align="right">
        <q-btn color="primary" label="Log in" @click="onOKClick" />
        <q-btn color="primary" label="Cancel" @click="onDialogCancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
<style scoped></style>
