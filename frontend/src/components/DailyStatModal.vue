<script setup>
  import {
    useDialogPluginComponent,
    QBtn,
    QCard,
    QCardActions,
    QDialog,
    QInput,
    QSlider,
  } from 'quasar';
  import * as styles from '../style.module.css';
  import { ref, watch } from 'vue';
  import DatePicker from './DatePicker.vue';
  import * as utils from '../modules/utils';
  import * as dailyStatsUtils from '../modules/dailyStatsUtils';
  import * as dateUtils from '../modules/dateUtils';

  defineEmits([...useDialogPluginComponent.emits]);
  const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
    useDialogPluginComponent();

  const props = defineProps({ statName: String, stats: Object });

  if (!dailyStatsUtils.labels[props.statName]) {
    throw new Error('Unknown stat');
  }

  const title = ref(dailyStatsUtils.labels[props.statName][0]);

  const stats = props.stats
    ? ref(utils.deepToRaw(props.stats))
    : ref(dailyStatsUtils.emptyStats());

  const onOKClick = () => {
    onDialogOK(stats.value);
  };

  const disableSave = ref(true);

  watch(
    () => {
      return stats.value;
    },
    (newStat) => {
      if (props.statName == 'bg') {
        disableSave.value =
          typeof dailyStatsUtils.bgValidator(newStat.bg) == 'string';
      } else if (props.statName == 'bp') {
        disableSave.value =
          typeof dailyStatsUtils.systolicValidator(newStat.bp[0]) == 'string' ||
          typeof dailyStatsUtils.diastolicValidator(newStat.bp[1]) == 'string';
      } else if (props.statName == 'bodyweight') {
        disableSave.value =
          typeof dailyStatsUtils.bodyWeightValidator(newStat.bodyweight) ==
          'string';
      } else if (props.statName == 'sleep') {
        disableSave.value =
          typeof dailyStatsUtils.sleepValidator(newStat.sleep) == 'string';
      } else if (props.statName == 'food') {
        disableSave.value =
          typeof dailyStatsUtils.foodDescriptionValidator(
            newStat.food.description
          ) == 'string';
      } else if (props.statName == 'spirit') {
        disableSave.value = false;
      }
    },
    { deep: true }
  );
</script>
<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card dark class="q-dialog-plugin" :class="[styles.dailyStatsWrap]">
      <h1>{{ title }}</h1>
      <DatePicker
        v-if="!props.stats || !props.stats.date"
        @update="(value) => (stats.date = value)"
      />
      <div v-else>
        {{ dateUtils.formatDate(dateUtils.dateFromSeconds(props.stats.date)) }}
      </div>

      <div v-if="props.statName == 'bg'">
        <q-input
          v-model.number="stats.bg"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][0]"
          :suffix="dailyStatsUtils.labels[props.statName][1]"
          dark
          :rules="[dailyStatsUtils.bgValidator]"
        />
      </div>

      <div v-if="props.statName == 'bp'">
        <q-input
          v-model.number="stats.bp[0]"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][1]"
          :suffix="dailyStatsUtils.labels[props.statName][3]"
          dark
          :rules="[dailyStatsUtils.systolicValidator]"
          v-focus
        />
        <q-input
          v-model.number="stats.bp[1]"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][2]"
          :suffix="dailyStatsUtils.labels[props.statName][3]"
          dark
          :rules="[dailyStatsUtils.diastolicValidator]"
        />
      </div>

      <div v-if="props.statName == 'sleep'">
        <q-input
          v-model.number="stats.sleep"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][0]"
          :suffix="dailyStatsUtils.labels[props.statName][1]"
          dark
          :rules="[dailyStatsUtils.sleepValidator]"
          v-focus
        />
      </div>

      <div v-if="props.statName == 'spirit'">
        <div :class="[styles.dailySliders]">
          <div :class="[styles.dailySlider]">
            <q-slider
              v-model="stats.mood"
              vertical
              :min-="0"
              :max="5"
              :step="1"
              snap
              reverse
              dark
              label
              label-always
            />
            <div>Mood</div>
          </div>
          <div :class="[styles.dailySlider]">
            <q-slider
              v-model="stats.stress"
              vertical
              :min-="0"
              :max="5"
              :step="1"
              snap
              reverse
              dark
              label
              label-always
            />
            <div>Stress</div>
          </div>
          <div :class="[styles.dailySlider]">
            <q-slider
              v-model="stats.energy"
              vertical
              :min-="0"
              :max="5"
              :step="1"
              snap
              reverse
              dark
              label
              label-always
            />
            <div>Energy</div>
          </div>
        </div>
      </div>

      <div v-if="props.statName == 'bodyweight'">
        <q-input
          v-model.number="stats.bodyweight"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][0]"
          :suffix="dailyStatsUtils.labels[props.statName][1]"
          dark
          :rules="[dailyStatsUtils.bodyWeightValidator]"
          v-focus
        />
      </div>

      <div v-if="props.statName == 'food'">
        <q-input
          v-model="stats.food.description"
          :label="dailyStatsUtils.labels[props.statName][1]"
          dark
          :rules="[dailyStatsUtils.foodDescriptionValidator]"
        />
        <q-input
          v-model.number="stats.food.protein"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][2]"
          :suffix="dailyStatsUtils.labels[props.statName][6]"
          dark
          :rules="[dailyStatsUtils.foodNutrientValidator]"
        />
        <q-input
          v-model.number="stats.food.carbs"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][3]"
          :suffix="dailyStatsUtils.labels[props.statName][6]"
          dark
          :rules="[dailyStatsUtils.foodNutrientValidator]"
        />
        <q-input
          v-model.number="stats.food.fat"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][4]"
          :suffix="dailyStatsUtils.labels[props.statName][6]"
          dark
          :rules="[dailyStatsUtils.foodNutrientValidator]"
        />
        <q-input
          v-model.number="stats.food.fiber"
          type="number"
          :label="dailyStatsUtils.labels[props.statName][5]"
          :suffix="dailyStatsUtils.labels[props.statName][6]"
          dark
          :rules="[dailyStatsUtils.foodNutrientValidator]"
        />
      </div>

      <q-card-actions align="right">
        <q-btn
          color="accent"
          text-color="dark"
          label="Cancel"
          @click="onDialogCancel"
        />
        <q-btn
          color="accent"
          text-color="dark"
          label="Save"
          :disable="disableSave"
          @click="onOKClick"
          :class="[styles.maxRight]"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
