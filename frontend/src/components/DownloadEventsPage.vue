<script setup>
  import DatePicker from './DatePicker.vue';
  import { onBeforeMount, ref, watch } from 'vue';
  import * as dateUtils from '../modules/dateUtils';
  import * as utils from '../modules/utils';
  import * as state from '../modules/state';
  import { QInput, QBtn } from 'quasar';
  import * as styles from '../style.module.css';

  // default date range is 16 weeks since now
  const endDate = ref(dateUtils.nowInSeconds());
  const startDate = ref(endDate.value - 16 * 7 * 24 * 60 * 60);
  const content = ref('');

  const updateStartDate = (newStartDate) => {
    if (newStartDate > endDate.value) {
      utils.openConfirmModal('Start date must occur before end date');
    } else {
      startDate.value = newStartDate;
    }
  };

  /**
   * Sets the end date of the charted time period.
   * @param newEndDate The end date in seconds since epoch.
   */
  const updateEndDate = (newEndDate) => {
    if (newEndDate < startDate.value) {
      utils.openConfirmModal('End date must occur after start date');
    } else {
      endDate.value = newEndDate;
    }
  };

  const formatDate = (date) => {
    const dateObj = dateUtils.dateFromSeconds(date);
    return dateObj.toDateString();
  };

  const formatIntensity = (value, type) => {
    const intensityProps = utils.intensityProps(type);
    const prefix = intensityProps.prefix ? intensityProps.prefix + ': ' : '';
    return `${prefix}${intensityProps.format(value)}`;
  };

  const formatVolume = (value, type, volumeConstraint) => {
    let formatted = '';
    if (type === 'distance') {
      formatted = (value / 1000).toFixed(2);
    } else if (type === 'time') {
      const hours = Math.floor(value / 3600).toString();

      let minutes = Math.floor((value % 3600) / 60).toString();
      if (minutes.length == 1) {
        minutes = `0${minutes}`;
      }

      let seconds = Math.floor((value % 360) % 60).toString();
      if (seconds.length == 1) {
        seconds = `0${seconds}`;
      }
      formatted = `${hours}:${minutes}:${seconds}`;
    } else if (type === 'count') {
      for (const reps of value) {
        if (volumeConstraint == 0) {
          formatted += `${reps} / `;
        } else if (volumeConstraint == 1) {
          let count = 0;
          for (const rep of reps) {
            count = rep > 0 ? count + 1 : count;
          }
          formatted += `${count} / `;
        } else {
          for (const rep of reps) {
            formatted += `${rep}-`;
          }
          if (formatted.endsWith('-')) {
            formatted = formatted.slice(0, -1);
          }
          formatted += ' / ';
        }
      }
      if (formatted.endsWith('/ ')) {
        formatted = formatted.slice(0, -2);
      }
    }
    return formatted;
  };

  const getContent = () => {
    content.value = '';
    const events = state.eventStore.getAll();
    let line = '';
    for (const e of events) {
      if (e.date > endDate.value) {
        break;
      }
      if (e.date >= startDate.value) {
        line += `${formatDate(e.date)}`;
        line += ` (${state.activityStore.get(e.activityID).name})\n\n`;
        for (const exKey in e.exercises) {
          const exerciseType = state.exerciseTypeStore.get(
            e.exercises[exKey].typeID,
          );
          line += `\t${exerciseType.name}:\n`;

          for (const partKey in e.exercises[exKey].parts) {
            const intensityType = exerciseType.intensityType;
            const intensity = e.exercises[exKey].parts[partKey].intensity;
            line += `\t ${formatIntensity(intensity, intensityType)} x `;
            const volume = e.exercises[exKey].parts[partKey].volume;
            const constraint = exerciseType.volumeConstraint;
            const type = exerciseType.volumeType;
            line += `${formatVolume(volume, type, constraint)}\n`;
          }
        }
        line += `\n\tOverall: ${e.overall ? e.overall : 'N/A'}\n`;
        line += `\tNotes: ${e.notes}\n`;
        line += `------------------------------------\n`;
      }
    }
    content.value += `${line}`;
  };

  watch(
    () => {
      return startDate.value;
    },
    async () => {
      state.eventStore.clear();
      await utils.fetchEvents('', endDate.value, startDate.value);
      getContent();
    },
  );

  watch(
    () => {
      return endDate.value;
    },
    async () => {
      state.eventStore.clear();
      await utils.fetchEvents('', endDate.value, startDate.value);
      getContent();
    },
  );

  onBeforeMount(async () => {
    if (!state.eventStore.events.length) {
      await utils.fetchEvents('', endDate.value, startDate.value);
    }
    getContent();
  });

  function downloadFile() {
    // Create a Blob object from the content
    const blob = new Blob([content.value], { type: 'text/plain' });

    // Create a temporary URL for the blob
    const url = URL.createObjectURL(blob);

    // Use the anchor element method
    const link = document.createElement('a');
    link.href = url;
    link.download = 'homegymLog.txt';
    link.style.display = 'none';

    document.body.appendChild(link);
    link.click();

    // Clean up by revoking the object URL and removing the link
    setTimeout(() => {
      URL.revokeObjectURL(url);
      document.body.removeChild(link);
    }, 0); // Timeout needed for Firefox compatibility
  }
</script>

<template>
  <div :class="[styles.downloadPage]">
    <div>
      <div :class="[styles.downloadDates]">
        <div>
          <div>Start Date:</div>
          <DatePicker
            :dateValue="startDate"
            :hideTime="true"
            @update="updateStartDate"
          />
        </div>
        <div>
          <div>End Date:</div>
          <DatePicker
            :dateValue="endDate"
            :hideTime="true"
            @update="updateEndDate"
          />
        </div>
      </div>
      <div><q-btn dark label="Download" @click="downloadFile" /></div>
    </div>
    <div :class="[styles.downloadContent]">
      <q-input v-model="content" dark autogrow disabled />
    </div>
  </div>
</template>
