<script setup>
  /**
   * Dialog for logging in.
   */
  import {
    useDialogPluginComponent,
    QDialog,
    QCard,
    QInput,
    QIcon,
    QCardActions,
    QBtn,
  } from 'quasar';
  import { login } from './../modules/utils.js';
  import { ref } from 'vue';
  let password = ref('');
  let isPwd = ref(true);
  let username = ref('');

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

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
        <q-btn
          color="accent"
          text-color="dark"
          label="Log in"
          @click="onOKClick"
        />
        <q-btn
          color="accent"
          text-color="dark"
          label="Cancel"
          @click="onDialogCancel"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
