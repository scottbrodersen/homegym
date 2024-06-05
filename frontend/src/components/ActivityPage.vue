<script setup>
  import { activityStore, exerciseTypeStore } from '../modules/state';
  import { computed, reactive, ref, toRaw } from 'vue';
  import {
    authPromptAsync,
    updateActivityExercises,
    fetchActivityExercises,
    ErrNotLoggedIn,
    newActivityPrompt,
    states,
    toast,
  } from '../modules/utils';
  import {
    QBtn,
    QCheckbox,
    QPopupEdit,
    QIcon,
    QInput,
    QList,
    QItem,
    QItemSection,
    QItemLabel,
  } from 'quasar';
  import * as styles from '../style.module.css';

  // selected activity
  const currentActivity = reactive({
    id: '',
    name: '',
  });

  const selectedTypeIDs = ref([]);

  const state = ref(states.READ_ONLY);

  const setCurrentActivity = async (activity) => {
    if (activity.exercises == null) {
      try {
        await fetchActivityExercises(activity.id);
      } catch (e) {
        if (e instanceof ErrNotLoggedIn) {
          console.log(e.message);
          await authPromptAsync();
          setCurrentActivity(activity);
        } else {
          console.log(e);
        }
      }
    }
    currentActivity.id = activity.id;
    currentActivity.name = activity.name;
    selectedTypeIDs.value = toRaw(activity.exercises);
  };

  const isChanged = computed(() => {
    const storedTypeIDs = activityStore.get(currentActivity.id).exercises;

    if (storedTypeIDs == null) {
      return false;
    }

    if (selectedTypeIDs.value.length != storedTypeIDs.length) {
      return true;
    }

    storedTypeIDs.forEach((id) => {
      if (!selectedTypeIDs.value.includes(id)) {
        return true;
      }
    });

    return false;
  });

  const resetValues = () => {
    if (!!currentActivity.id) {
      selectedTypeIDs.value = toRaw(
        activityStore.get(currentActivity.id).exercises
      );
    }
  };

  // called when activity name is changed
  const updateActivityName = async (newName) => {
    // make sure we have a change
    const storedActivity = activityStore.get(currentActivity.id);

    if (storedActivity.name != newName) {
      const url = `/homegym/api/activities/${currentActivity.id}/`;

      const update = {
        name: newName,
        id: currentActivity.id,
      };

      const headers = new Headers();
      headers.set('content-type', 'application/json');

      const options = {
        method: 'POST',
        body: JSON.stringify(update),
        headers: headers,
      };

      const resp = await fetch(url, options);

      if (resp.status === 401) {
        console.log('unauthorized request to upsert activity name');
        await authPromptAsync();
        updateActivityName(newName);
      } else if (resp.status === 400) {
        const body400 = await resp.json();
        if (body400.message.includes('not unique')) {
          // TODO: use a dialog instead of alert
          alert(body400.message);
        }
        return;
      } else if (resp.status < 200 || resp.status > 299) {
        throw new Error();
      }

      const activity = toRaw(activityStore.get(currentActivity.id));
      activity.name = newName;
      activityStore.add(activity);

      currentActivity.name = newName;
    }
  };

  const saveExerciseIDs = async () => {
    const updatedActivity = {
      id: currentActivity.id,
      name: currentActivity.name,
      exercises: [...selectedTypeIDs.value],
    };

    try {
      await updateActivityExercises(updatedActivity);
      toast('Saved', 'positive');
    } catch (e) {
      if (e instanceof ErrNotLoggedIn) {
        console.log(e.message);
        await authPromptAsync();
        saveExerciseIDs();
      } else {
        console.log(e);
        toast('Error', 'negative');
      }
    }
  };
</script>

<template>
  <div :class="[styles.grid2Col]">
    <div :class="[styles.colTitleWrapper, styles.leftColumn]">
      <div :class="[styles.listTitle, styles.sibSpSmall]">Activities</div>
      <div :class="[styles.sibSpSmall]">
        <q-btn
          size="0.65em"
          round
          :disable="state == states.EDIT"
          icon="add"
          name="add-activity"
          color="primary"
          @Click="newActivityPrompt"
        />
      </div>
      <div :class="[styles.sibSpSmall]">
        <q-btn
          name="edit-activity"
          size="0.65em"
          round
          :disable="state == states.EDIT || currentActivity.id == ''"
          icon="edit"
          color="primary"
        />
        <q-popup-edit
          :model-value="currentActivity.name"
          persistent
          buttons
          label-set="Save"
          @update:model-value="updateActivityName"
          v-slot="scope"
        >
          <q-input
            v-model="scope.value"
            dense
            autofocus
            counter
            v-focus
            @keyup.enter="scope.set"
          />
        </q-popup-edit>
      </div>
    </div>

    <div :class="[styles.leftColumn]">
      <q-list
        :class="[styles.listStd, styles.blockBorder]"
        bordered
        separator
        dense
      >
        <q-item
          clickable
          v-for="[id, activity] in activityStore.activities"
          :key="id"
          :disable="state == states.EDIT && id != currentActivity.id"
          @Click.stop="setCurrentActivity(activity)"
        >
          <q-item-section>
            <q-item-label>{{ activity.name }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-icon
              name="start"
              color="green"
              size="sm"
              v-show="currentActivity.id == id"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <div :class="[styles.colTitleWrapper, styles.rightColumn]">
      <div :class="[styles.listTitle, styles.sibSpSmall]">Exercises</div>
      <div :class="[styles.sibSpSmall]">
        <q-btn
          size="0.65em"
          round
          :disable="state == states.EDIT || currentActivity.id == ''"
          icon="edit"
          name="edit-exercises"
          color="primary"
          @Click="state = states.EDIT"
        />
      </div>
    </div>
    <div :class="[styles.rightColumn]">
      <div v-if="state == states.EDIT">
        <q-list :class="[styles.listStd]" compact>
          <q-item
            v-for="[id, eType] in exerciseTypeStore.exerciseTypes"
            :key="id"
          >
            <q-item-section>
              <q-checkbox
                v-model="selectedTypeIDs"
                :val="eType.id"
                :label="eType.name"
                dark
              />
            </q-item-section>
          </q-item>
        </q-list>
        <div :class="[styles.buttonArray]">
          <q-btn
            color="accent"
            text-color="dark"
            icon="save"
            @click="saveExerciseIDs"
            :disabled="!isChanged"
          />
          <q-btn
            color="accent"
            text-color="dark"
            icon="restart_alt"
            @click="resetValues"
            :disabled="!isChanged"
          />
          <q-btn
            color="accent"
            text-color="dark"
            icon="done"
            :disabled="isChanged"
            @click="state = states.READ_ONLY"
            :class="[styles.maxRight]"
          />
        </div>
      </div>
      <div v-if="state == states.READ_ONLY">
        <q-list
          :class="[styles.listStd, styles.blockBorder]"
          dense
          bordered
          separator
          ><q-item v-if="currentActivity.id == ''"
            ><q-item-section
              ><q-item-label>Select an activity</q-item-label></q-item-section
            ></q-item
          >
          <q-item v-for="e in selectedTypeIDs" :key="e.id">
            <q-item-section>
              <q-item-label>{{ exerciseTypeStore.get(e).name }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </div>
    </div>
  </div>
</template>
