<script setup>
  import { useDialogPluginComponent, QForm, QInput } from 'quasar';
  import { ErrNotLoggedIn, ErrNotUnique } from '../modules/utils.js';
  import { ref, computed } from 'vue';
  import { activityStore } from '../modules/state.js';

  const name = ref('');

  const nameIsInValid = computed(() => {
    if (name.value.length < 3) {
      return true;
    }
    return false;
  });

  const addActivity = (newname) => {
    const url = `/homegym/api/activities/`;

    const activity = {
      name: newname,
    };

    const headers = new Headers();
    headers.set('content-type', 'application/json');

    const options = {
      method: 'POST',
      body: JSON.stringify(activity),
      headers: headers,
    };

    let respStatus;
    fetch(url, options)
      .then((resp) => {
        respStatus = resp.status;
        return resp.json();
      })
      .then((body) => {
        if (respStatus < 200 || respStatus > 299) {
          if (respStatus === 401) {
            throw new ErrNotLoggedIn('unauthorized request to upsert activity');
          } else if (
            respStatus === 400 &&
            body.message.includes('not unique')
          ) {
            throw new ErrNotUnique('activity name must be unique');
          } else {
            throw new Error();
          }
        }
        activity.id = body.id;
        // prevent unncessary fetch
        activity.exercises = [];
        activityStore.add(activity);
      })
      .catch((e) => {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          authPrompt(addActivity, [newname]);
        } else if (e instanceof ErrNotUnique) {
          console.log(e.message);
        }
        throw e;
      });
  };

  defineEmits([...useDialogPluginComponent.emits]);

  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  function onOKClick() {
    try {
      addActivity(name.value);
    } catch (e) {
      if (e instanceof ErrNotUnique) {
        alert('Activity name must be unique');
      } else {
        console.log(e);
      }
    }
    onDialogOK();
  }
</script>

<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.bgBlack]">
      <q-form>
        <q-input
          v-model="name"
          filled
          v-focus
          type="text"
          label="Activity Name"
          dark
        />
      </q-form>

      <q-card-actions align="right">
        <q-btn
          color="accent"
          text-color="dark"
          label="Save"
          @click="onOKClick"
          :disabled="nameIsInValid"
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
<style scoped></style>
